package cmd

import (
	"context"
	"os"
	"path"

	"github.com/spf13/cobra"
)

const (
	defaultContextKey = "gov_key"
)

var (
	GovmContext = context.TODO()

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
	govHomeDir := path.Join(home, ".govm")
	// govConfig := path.Join(govHomeDir, ".config")
	govSourcesDir := path.Join(govHomeDir, "sources")
	govCacheDir := path.Join(govHomeDir, "cache")

	GovmContext = context.WithValue(GovmContext, defaultContextKey, map[string]string{
		"home":    govHomeDir,
		"sources": govSourcesDir,
		"cache":   govCacheDir,
	})

	// Gov.Flags().Int("verbose", 1, "verbosity level {1,2, 3}. Default: 1")
	// Gov.Flags().String("cache", defaultCacheDir, "default cache directory")
	// Gov.Flags().String("sources", defaultSourcesDir, "default sources directory")

	Govm.AddCommand(currentGo)
	Govm.AddCommand(listSources)
	Govm.AddCommand(selectSource)

	cobra.OnInitialize(func() {
		if _, err := os.Stat(govHomeDir); os.IsNotExist(err) {
			if err = os.MkdirAll(govHomeDir, 0777); err != nil {
				Govm.Println("cannot create main gov directories like: ", govHomeDir, govSourcesDir)
				os.Exit(1)
			}
		} else if os.IsPermission(err) {
			Govm.Println("cannot access gov home, please check directory permissions")
			os.Exit(1)
		}

		if err := os.MkdirAll(govSourcesDir, 0777); err != nil {
			Govm.Println("cannot create main gov directories like: ", govHomeDir, govSourcesDir)
			os.Exit(1)
		}
		if err := os.MkdirAll(govCacheDir, 0777); err != nil {
			Govm.Println("cannot create main gov directories like: ", govHomeDir, govCacheDir)
			os.Exit(1)
		}

		/*if _, err := os.Stat(govConfig); os.IsNotExist(err) {
			f, err := os.Create(govConfig)
			if err != nil {
				Gov.Println(err)
			}
			f.Close()
		} else if os.IsPermission(err) {
			Gov.Println("cannot access gov config file, please check file permissions")
			os.Exit(1)
		}*/
	})
}

func govCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}
}

func ExecuteWithContext() {
	Govm.ExecuteContext(GovmContext)
}
