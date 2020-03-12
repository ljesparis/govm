package cmd

import (
	"github.com/ljesparis/govm/pkg/utils"
	"github.com/spf13/cobra"
)

type contextKey uint8

const (
	ctxKey contextKey = iota
)

func getCorrectPackage(version string, cmd *cobra.Command) (p string, err error) {
	os, _ := cmd.Flags().GetString("os")
	arch, _ := cmd.Flags().GetString("arch")
	pt, _ := cmd.Flags().GetString("package")

	p, err = utils.GetPackageFilename(version, os, arch, pt)
	return
}
