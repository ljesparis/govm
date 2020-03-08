package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/ljesparis/govm/pkg/utils"
	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
)

const (
	// Url where golang can be downloaded
	goSourceUrl string = "https://dl.google.com/go"
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
	pt, _ := utils.DefaultSystemPackageType()
	selectSource.Example = "  gov select 1.14\n  gov s 1.14 --arch 386\n  gov sl 1.14 --os windows"

	selectSource.Flags().BoolP("cache", "c", false, "cache downloaded sources")
	selectSource.Flags().String("os", runtime.GOOS, "os compatible sources")
	selectSource.Flags().String("arch", runtime.GOARCH, "os compatible architecture")
	selectSource.Flags().String("package", pt, "default package type")
}

func selectGoVersionCmd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context().Value(defaultContextKey).(map[string]string)
	cacheDir := ctx["cache"]
	sourcesDir := ctx["sources"]

	userSelectedGoVersion := args[0]
	sourceName := fmt.Sprintf("go%s", userSelectedGoVersion)
	sourceCompletePath := path.Join(sourcesDir, sourceName)
	compressedSourceName, err := getCorrectPackage(userSelectedGoVersion, cmd)
	if err != nil && err == utils.UnknowPackage {
		cmd.Println("Selected package is not supported")
		return
	}
	sourcesUrl := fmt.Sprintf("%s/%s", goSourceUrl, compressedSourceName)
	compressedCompletePath := path.Join(cacheDir, compressedSourceName)

	useCache, _ := cmd.Flags().GetBool("cache")

	gv, _ := utils.GetCurrentGoVersion()
	if gv == sourceName {
		cmd.Println("Go version already selected")
		return
	}

GoSources:
	// Section to download and decompress go source.

	// Check if go source path exists
	if _, err := os.Stat(sourceCompletePath); os.IsNotExist(err) || useCache {
		cmd.Printf("Go with version %s does not exists\n", userSelectedGoVersion)

		if _, err := os.Stat(compressedCompletePath); os.IsNotExist(err) {

			err = utils.DefaultClient(sourcesUrl, compressedCompletePath).Download(func(s *utils.Stats) {
				cmd.Printf("\rdownloading %s/100%s", s.DownloadedPercentage, "%")
			})

			cmd.Printf("\n")

			if err != nil {
				// check if http error is 404 or anything else
				if err == utils.SourceDoesNotExists {
					cmd.Println("Error downloading the file, go version used does not exists")
				} else {
					cmd.Println("Error downloading the file, please check internet connection")
				}
				return
			}

		} else if os.IsPermission(err) {
			cmd.Printf("Does not have permission to read the following package: %s\n", compressedSourceName)
			return
		} else {
			if useCache {
				cmd.Println("Using compressed source from cache")
			}
		}

		if err := archiver.Unarchive(compressedCompletePath, sourceCompletePath); err != nil && os.IsNotExist(err) {
			cmd.PrintErrln("Error decompressing sources")
			return
		}

		// Checking if cache flag was set to save compressed source
		if !useCache {
			if err := os.Remove(compressedCompletePath); err != nil {
				cmd.Printf("Error deleting %s source", compressedCompletePath)
			}
		}

	} else if os.IsPermission(err) {
		cmd.Printf("Does not have permission to access the following path: %s\n", sourceCompletePath)
		return
	} else {
		// Check if go source path does contain go binary.
		if _, err := utils.GetGoversion(path.Join(sourceCompletePath, goSourceBin)); err != nil {

			// If go binary does not exists, source path will be deleted
			if err = os.RemoveAll(sourceCompletePath); err != nil && os.IsPermission(err) {
				cmd.Printf("Does not have permission to delete %s directory, please check", sourceCompletePath)
				return
			} else {
				goto GoSources
			}
		} else {
			goto GoSymLinks
		}
	}

GoSymLinks:
	// Section to create symbolic links.
	// If symbolic link already exists, will be deleted and then created.
	// Otherwise, will be created

	goDst := path.Join(defaultGoBinDir, "go")
	gofmtDst := path.Join(defaultGoBinDir, "gofmt")
	goSrc := path.Join(sourceCompletePath, goSourceBin)
	gofmtSrc := path.Join(sourceCompletePath, gofmtSourceBin)

	os.Remove(goDst) // delete go symbolic link if exists
	os.Remove(gofmtDst) // delete gofmt symbolic link if exists

	if err := os.Symlink(goSrc, goDst); err != nil && os.IsPermission(err) {
		cmd.PrintErrln("Cannot create go symlink, please check permissions")
		return
	}

	if err := os.Symlink(gofmtSrc, gofmtDst); err != nil && os.IsPermission(err) {
		cmd.PrintErrln("Cannot create gofmt symlink, please check permissions")
	}

	cmd.Println("Done!")
	cmd.Println("Go and try 'go version' command")
}
