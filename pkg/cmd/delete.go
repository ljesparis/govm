package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

var (
	deleteSource = &cobra.Command{
		Use:     "delete [version]",
		Aliases: []string{"d", "dd"},
		Short:   "Delete golang source",
		Args:    cobra.ExactValidArgs(1),
		Run:     deleteSourceCmd,
	}
)

func deleteSourceCmd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context().Value(ctxKey).(map[string]string)
	sourcesDir := ctx["sources"]
	goVersion := args[0]

	if isVersionInUse(goVersion) {
		cmd.Println("Delete go version in use is not allowed")
		os.Exit(1)
	}

	rSource := path.Join(sourcesDir, goVersion)
	if _, err := os.Stat(rSource); os.IsNotExist(err) {
		cmd.Println("Go source does not exists")
		os.Exit(1)
	} else if os.IsPermission(err) {
		cmd.Println("Permission denied while trying to delete go source")
		os.Exit(1)
	}

	_ = os.RemoveAll(rSource)
}
