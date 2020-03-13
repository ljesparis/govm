package cmd

import (
	"context"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var (
	govmContext = context.TODO()

	//Govm represents the cli command to execute
	Govm = &cobra.Command{
		Use:                   "govm [flags] [command]",
		Version:               "v1.0.0-alpha.0",
		Short:                 "govm is a go version manager",
		DisableFlagsInUseLine: true,
		Run:                   govCmd,
	}
)

func init() {
	home, _ := os.UserHomeDir()
	govmHomeDir := path.Join(home, ".govm")
	// govmBinDir := path.Join(govHomeDir, "bin")
	govmSourcesDir := path.Join(govmHomeDir, "sources")
	govmCacheDir := path.Join(govmHomeDir, "cache")

	govmContext = context.WithValue(govmContext, ctxKey, map[string]string{
		"home":    govmHomeDir,
		"sources": govmSourcesDir,
		"cache":   govmCacheDir,
	})

	Govm.AddCommand(listSources)
	Govm.AddCommand(selectSource)

	cobra.OnInitialize(func() {
		if _, err := os.Stat(govmHomeDir); os.IsNotExist(err) {
			if err = os.MkdirAll(govmHomeDir, 0777); err != nil {
				Govm.Println("cannot create main gov directories like: ", govmHomeDir, govmSourcesDir)
				os.Exit(1)
			}
		} else if os.IsPermission(err) {
			Govm.Println("cannot access gov home, please check directory permissions")
			os.Exit(1)
		}

		if err := os.MkdirAll(govmSourcesDir, 0777); err != nil {
			Govm.Println("cannot create main gov directories like: ", govmHomeDir, govmSourcesDir)
			os.Exit(1)
		}
		if err := os.MkdirAll(govmCacheDir, 0777); err != nil {
			Govm.Println("cannot create main gov directories like: ", govmHomeDir, govmCacheDir)
			os.Exit(1)
		}
	})
}

func govCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		if err := cmd.Help(); err != nil {
			cmd.Println(err)
		}
	}
}

// Run will with govm app with context
func Run() {
	if err := Govm.ExecuteContext(govmContext); err != nil {
		Govm.Println(err)
	}
}
