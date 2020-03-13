package cmd

import (
	"io/ioutil"
	"os"

	"github.com/ljesparis/govm/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	listSources = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List golang sources",
		Args:    cobra.ExactValidArgs(0),
		Run:     listSourcesCmd,
	}
)

func listSourcesCmd(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context().Value(ctxKey).(map[string]string)
	cv, _ := utils.GetCurrentGoVersion()

	dirs, err := ioutil.ReadDir(ctx["sources"])
	if err != nil {
		cmd.Println("error listing sources content")
		os.Exit(1)
	}

	cmd.Println("available sources:")
	cmd.Println()

	if len(dirs) > 0 {
		for _, el := range dirs {
			if el.IsDir() {
				if el.Name() == cv {
					cmd.Printf("=> %s\n", el.Name())
				} else {
					cmd.Println("  ", el.Name())
				}
			}
		}
	} else {
		cmd.Println("There's not downloaded go sources")
	}

	cmd.Println()
}
