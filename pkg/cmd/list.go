package cmd

import (
	"io/ioutil"
	"os"
	"strings"

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
				oldName := el.Name()
				newName := strings.ReplaceAll(oldName, "go", "")
				if oldName == cv {
					cmd.Printf("=> %s\n", newName)
				} else {
					cmd.Println("  ", newName)
				}
			}
		}
	} else {
		cmd.Println("[*] there's not downloaded go sources")
	}

	cmd.Println()
}
