package cmd

import (
	"context"
	"os"
	"path"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	// common permission
	perm os.FileMode = 0770
)

var (
	govmContext = context.TODO()

	//Govm represents the cli command to execute
	Govm = &cobra.Command{
		Use:                   "govm [flags] [command]",
		Version:               "v1.0.0-alpha.1",
		Short:                 "govm is a go version manager",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if err := cmd.Help(); err != nil {
					cmd.Println(err)
				}
			}
		},
	}
)

func init() {
	Govm.AddCommand(listSources)
	Govm.AddCommand(selectSource)
	Govm.AddCommand(deleteSource)

	home, _ := os.UserHomeDir()
	govmHomeDir := path.Join(home, ".govm")
	govmSourcesDir := path.Join(govmHomeDir, "sources")
	govmCacheDir := path.Join(govmHomeDir, "cache")
	
	var govmBinDir string
	if runtime.GOOS == "windows" {
		govmBinDir = path.Join(govmHomeDir, "bin")
	} else {
		govmBinDir = path.Join(home, defaultGoBinDir)
	}

	govmContext = context.WithValue(govmContext, ctxKey, map[string]string{
		"home":    govmHomeDir,
		"sources": govmSourcesDir,
		"cache":   govmCacheDir,
		"bin":     govmBinDir,
	})

	cobra.OnInitialize(func() {
		if err := os.MkdirAll(govmHomeDir, perm); os.IsPermission(err) {
			Govm.Println("cannot create govm home directory, please check permissions.")
			os.Exit(1)
		}

		if err := os.MkdirAll(govmSourcesDir, perm); os.IsPermission(err) {
			Govm.Println("cannot create sources directory, please check home folder permissions")
			os.Exit(1)
		}
		if err := os.MkdirAll(govmCacheDir, perm); os.IsPermission(err) {
			Govm.Println("cannot create cache directory, please check home folder permissions")
			os.Exit(1)
		}
		if runtime.GOOS == "windows" {
			if err := os.MkdirAll(govmBinDir, perm); os.IsPermission(err) {
				Govm.Println("cannot create binary directory, please check home folder permissions")
				os.Exit(1)
			}
		}
	})
}

// Run will with govm app with context
func Run() {
	if err := Govm.ExecuteContext(govmContext); err != nil {
		Govm.Println(err)
	}
}
