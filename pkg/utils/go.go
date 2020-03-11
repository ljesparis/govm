package utils

import (
	"bytes"
	"os/exec"
	"regexp"
)

// GetGoVersion returns the current version of go which is
// in the path passed as parameter.
func GetGoVersion(goBinPath string) (string, error) {
	buffer := bytes.NewBuffer([]byte{})
	cmd := exec.Command(goBinPath, "version")
	cmd.Stdout = buffer
	cmd.Stderr = buffer
	if err := cmd.Run(); err != nil {
		return "", err
	}

	re, err := regexp.Compile(`go[0-9]+(.[0-9]+)*`)
	if err != nil {
		return "", err
	}

	return string(re.Find(buffer.Bytes())), nil
}

// GetCurrentGoVersion returns the current version of go which is
// in the global path:
//
// unix -> $PATH
// windows ->
func GetCurrentGoVersion() (string, error) {
	return GetGoVersion("go")
}
