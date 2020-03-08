package cmd

import (
	"github.com/ljesparis/govm/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	currentGo = &cobra.Command{
		Use:     "current",
		Aliases: []string{"c", "cr"},
		Short:   "Show current go version",
		Long:    "Show current go version that is used by the system",
		Args:    cobra.ExactValidArgs(0),
		Run:     currentGoFunc,
	}
)

func currentGoFunc(cmd *cobra.Command, _ []string) {
	cv, err := utils.GetCurrentGoVersion()
	if err != nil || len(cv) == 0 {
		cmd.Println("there's not active go version yet")
		return
	}

	cmd.Println(cv)
}
