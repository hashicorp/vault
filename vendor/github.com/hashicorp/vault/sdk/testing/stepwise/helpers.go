package stepwise

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

const pluginPrefix = "vault-plugin-"

// CompilePlugin is a helper method to compile a source plugin
// TODO refactor compile plugin input and output to be types
func CompilePlugin(name, pluginName, srcDir, tmpDir string) (string, string, string, error) {
	binName := name
	if !strings.HasPrefix(binName, pluginPrefix) {
		binName = fmt.Sprintf("%s%s", pluginPrefix, binName)
	}
	binPath := path.Join(tmpDir, binName)

	cmd := exec.Command("go", "build", "-o", binPath, path.Join(srcDir, fmt.Sprintf("cmd/%s/main.go", pluginName)))
	cmd.Stdout = &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	cmd.Stderr = errOut

	// match the target architecture of the docker container
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	if err := cmd.Run(); err != nil {
		// if err here is not nil, it's typically a generic "exit status 1" error
		// message. Return the stderr instead
		return "", "", "", errors.New(errOut.String())
	}

	// calculate sha256
	f, err := os.Open(binPath)
	if err != nil {
		return "", "", "", err
	}

	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", "", "", err
	}

	sha256value := fmt.Sprintf("%x", h.Sum(nil))

	return binName, binPath, sha256value, nil
}
