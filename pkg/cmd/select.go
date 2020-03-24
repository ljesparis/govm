package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ljesparis/govm/pkg/utils"
	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
)

const (
	// Url where golang can be downloaded
	goSourceURL string = "https://dl.google.com/go"
)

var (
	selectSource = &cobra.Command{
		Use:                   "select [flags] [go version]",
		Aliases:               []string{"s", "sl"},
		Short:                 "Select golang source",
		Long:                  "Select golang source and download it if does not exists in local system",
		Args:                  cobra.ExactValidArgs(1),
		DisableFlagsInUseLine: true,
		Run:                   selectGoVersionCmd,
	}
)

func init() {
	pt, _ := utils.DefaultPackageType()
	selectSource.Example = "  gov select 1.14\n  gov s 1.14 --arch 386\n  gov sl 1.14 --os windows"

	selectSource.Flags().BoolP("cache", "c", false, "cache downloaded sources")
	selectSource.Flags().String("os", runtime.GOOS, "os compatible sources")
	selectSource.Flags().String("arch", runtime.GOARCH, "os compatible architecture")
	selectSource.Flags().String("package", pt, "default package type")
}

func selectGoVersionCmd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context().Value(ctxKey).(map[string]string)
	cacheDir := ctx["cache"]
	sourcesDir := ctx["sources"]

	// check if selected source is compatible with current
	// operating system
	tmp, _ := cmd.Flags().GetString("os")
	if utils.IsOSSupported(tmp) && strings.Compare(tmp, runtime.GOOS) != 0 {
		cmd.Println("current operating system does not support '" + tmp + "' binaries.")
		os.Exit(1)
	} else if !utils.IsOSSupported(tmp) && strings.Compare(tmp, runtime.GOOS) != 0 {
		cmd.Println("unknown operating system.")
		os.Exit(1)
	}

	corrupted := false
	userSelectedGoVersion := args[0]
	sourceName := fmt.Sprintf("go%s", userSelectedGoVersion)
	sourceCompletePath := path.Join(sourcesDir, sourceName)
	compressedSourceName, err := getCorrectPackage(userSelectedGoVersion, cmd)
	if err != nil && err == utils.ErrUnknowPackage {
		cmd.Println("selected package is not supported")
		os.Exit(1)
	}
	sourcesURL := fmt.Sprintf("%s/%s", goSourceURL, compressedSourceName)
	compressedCompletePath := path.Join(cacheDir, compressedSourceName)

	if isVersionInUse(sourceName) {
		cmd.Println("go version already selected")
		os.Exit(1)
	}

GoSources:
	// Check if go source path exists
	if _, err := os.Stat(sourceCompletePath); os.IsNotExist(err) {

		if !corrupted {
			cmd.Printf("go with version %s does not exists\n", userSelectedGoVersion)
		}

		// download compressed source and exit if an err is returned
		if err = downloadCompressedSource(sourcesURL, compressedCompletePath, cmd); err != nil {
			os.Exit(1)
		}

		// trying to decompress source file
		if err := deCompressFile(compressedCompletePath, sourceCompletePath, cmd); err != nil {
			os.Exit(1)
		}

	} else if os.IsPermission(err) {
		cmd.Printf("does not have permission to access the following path: %s\n", sourceCompletePath)
		os.Exit(1)
	} else {
		// Check if go source path does contain go binary.
		if _, err := utils.GetGoVersion(path.Join(sourceCompletePath, goSourceBin)); err != nil {

			// If go binary does not exists, source path will be deleted
			if err = os.RemoveAll(sourceCompletePath); err != nil && os.IsPermission(err) {
				cmd.Printf("does not have permission to delete %s directory, please check", sourceCompletePath)
				os.Exit(1)
			} else {
				corrupted = true
				cmd.Println("go source is corrupted")
				goto GoSources
			}
		}
	}

	// If symbolic link already exists, will be deleted and then created.
	// Otherwise, will be created
	if err = createSymbolicLink(sourceCompletePath, cmd); err != nil {
		os.Exit(1)
	}
}

func downloadCompressedSource(url, compressedDstPath string, cmd *cobra.Command) error {

	// Check if compressed go source file does not exists
	if _, err := os.Stat(compressedDstPath); os.IsNotExist(err) {
		err := utils.DefaultClient(url, compressedDstPath, perm).Download(func(s *utils.Stats) {
			cmd.Printf("\rdownloading: %s/100%s", s.DownloadedPercentage, "%")
		})

		cmd.Printf("\n")

		if err != nil {

			// check if http error is 404
			if err == utils.ErrSourceDoesNotExists {
				cmd.Println("error downloading the file, go version does not exists")
			} else if os.IsPermission(err) {
				cmd.Println("error downloading the file, please check govm home permissions.")
			} else {
				cmd.Println("error downloading the file, please check internet connection")
			}

			return err
		}
	} else if os.IsPermission(err) {
		cmd.Println("does not have permission to read the package ", compressedDstPath)
		return err
	} else {
		cmd.Println("using compressed source from cache")
	}

	return nil
}

func deCompressFile(compressedSource, sourceDst string, cmd *cobra.Command) error {
	useCache, _ := cmd.Flags().GetBool("cache")

	cmd.Println("decompressing source")

	// decompress file
	if err := archiver.Unarchive(compressedSource, sourceDst); err != nil && os.IsNotExist(err) {
		cmd.PrintErrln("compressed source does not exists")
		return err
	} else if os.IsPermission(err) {
		cmd.PrintErrln("does not have permission to read file ", compressedSource)
		return err
	}

	// Checking if cache flag was set to save compressed source
	if !useCache {
		if err := os.Remove(compressedSource); err != nil {
			cmd.Printf("error deleting %s source\n", compressedSource)
		}
	}

	return nil
}

func createSymbolicLink(sourcePath string, cmd *cobra.Command) error {
	goDst := path.Join(defaultGoBinDir, "go")
	gofmtDst := path.Join(defaultGoBinDir, "gofmt")
	goSrc := path.Join(sourcePath, goSourceBin)
	gofmtSrc := path.Join(sourcePath, gofmtSourceBin)

	os.Remove(goDst)    // delete go symbolic link if exists
	os.Remove(gofmtDst) // delete gofmt symbolic link if exists

	// check if a symbolic link can be created
	if err := os.Symlink(goSrc, goDst); err != nil && os.IsPermission(err) {
		cmd.PrintErrln("Cannot create go symlink, please check permissions")
		return err
	}

	if err := os.Symlink(gofmtSrc, gofmtDst); err != nil && os.IsPermission(err) {
		cmd.PrintErrln("Cannot create gofmt symlink, please check permissions")
	}

	st, _ := utils.GetCurrentGoVersionAll()
	cmd.Printf("version selected ==> %s\n", st)

	return nil
}
