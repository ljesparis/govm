// +build aix darwin dragonfly freebsd linux  netbsd openbsd solaris

package cmd

import (
	"path"
	"os"
)

const (
	goBin = "go"
	gofmtBin = "gofmt"
	goSourceBin    = "go/bin/go"
	gofmtSourceBin = "go/bin/gofmt"
)

func binDir(_ string) string {
	home, _ := os.UserHomeDir()
	return path.Join(home, "/.local/bin")
}
