package utils

import (
	"bytes"
	"errors"
	"runtime"
	"strings"
	"text/template"
)

// Represent the package extension of the compressed go source code
// that are available at https://golang.org/dl/
const (
	TarBinary uint64 = iota + 1
	MsiBinary
	PkgBinary
	ZipBinary
)

var (
	// ErrUnknowPackage error used when package is not Tar,Msi,pkg or zip
	ErrUnknowPackage = errors.New("unknow package")
)

// IsValidOS check if os passed as parameter is supported
func IsOSSupported(os string) bool {
	oss := []string{"linux", "darwin", "freebsd", "windows"}
	for _, tmp := range oss {
		if strings.Compare(tmp, os) == 0 {
			return true
		}
	}
	return false
}

// DefaultPackageType return the default package type
// for every supported system.
func DefaultPackageType() (pts string, pti uint64) {
	os := runtime.GOOS
	if os == "linux" || os == "darwin" || os == "freebsd" {
		pts = "tar.gz"
		pti = TarBinary
	} else {
		pts = "zip"
		pti = PkgBinary
	}
	return
}

// IsPackageTypeSupported check if package type is supported
func IsPackageTypeSupported(pt string) (pti uint64, err error) {

	if pt == "tar.gz" {
		pti = TarBinary
	} else if pt == "msi" {
		pti = MsiBinary
	} else if pt == "pkg" {
		pti = PkgBinary
	} else if pt == "zip" {
		pti = ZipBinary
	} else {
		pti = 0
		err = ErrUnknowPackage
	}

	return
}

// GetPackageFilename return the compressed go source filename with
// the go version, operating system, architecture and package type.
func GetPackageFilename(goversion, os, arch string, pType string) (string, error) {

	_, err := IsPackageTypeSupported(pType)
	if err != nil {
		return "", err
	}

	tmpl := template.New("compiledSourceTmpl")
	tmpl, err = tmpl.Parse("go{{.GO_VERSION}}.{{.OS}}-{{.ARCH}}.{{.PACKAGE}}")
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, map[string]string{
		"GO_VERSION": goversion,
		"OS":         os,
		"ARCH":       arch,
		"PACKAGE":    pType,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
