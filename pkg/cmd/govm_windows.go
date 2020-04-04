package cmd

import (
	"path"
)

const (
	goBin = "go.exe"
	gofmtBin = "gofmt.exe"
	goSourceBin    = "go\\bin\\go.exe"
	gofmtSourceBin = "go\\bin\\gofmt.exe"
)

func binDir(home string) string {
	return path.Join(home, "bin")
}
