package utils

import (
	"bytes"
	"errors"
	"runtime"
	"text/template"
)

const (
	TarBinary uint64 = iota + 1
	MsiBinary
	PkgBinary
	ZipBinary
)

var (
	UnknowPackage = errors.New("unknow package")
)

func GetPackageType(os string) (pts string, pti uint64) {
	if os == "linux" || os == "darwing" || os == "freebsd" {
		pts = "tar.gz"
		pti = TarBinary
	} else {
		pts = "zip"
		pti = PkgBinary
	}

	return
}

func DefaultSystemPackageType() (pts string, pti uint64) {
	pts, pti = GetPackageType(runtime.GOOS)
	return
}

func String2Int(pt string) (pti uint64, err error) {

	if pt == "tar.gz" {
		pti = TarBinary
	} else if pt == "msi" {
		pti = MsiBinary
	} else if pt == "tar.gz" {
		pti = MsiBinary
	} else if pt == "pkg" {
		pti = PkgBinary
	} else if pt == "zip" {
		pti = ZipBinary
	} else {
		pti = 0
		err = UnknowPackage
	}

	return
}

func Int2String(pType uint64) (pt string, err error) {

	switch pType {
	case TarBinary:
		pt = "tar.gz"
	case MsiBinary:
		pt = "msi"
	case PkgBinary:
		pt = "pkg"
	case ZipBinary:
		pt = "zip"
	default:
		err = UnknowPackage
	}

	return
}

func PackageFilename1(goversion, os, arch string, pType uint64) (string, error) {

	pt, err := Int2String(pType)
	if err != nil {
		return "", err
	}

	tmpl := template.New("compiledSource")
	tmpl, err = tmpl.Parse("go{{.GO_VERSION}}.{{.OS}}-{{.ARCH}}.{{.PACKAGE}}")
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, map[string]string{
		"GO_VERSION": goversion,
		"OS":         os,
		"ARCH":       arch,
		"PACKAGE":    pt,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func PackageFilename2(goversion, os, arch string, pType string) (string, error) {

	pt, err := String2Int(pType)
	if err != nil {
		return "", err
	}

	return PackageFilename1(goversion, os, arch, pt)
}
