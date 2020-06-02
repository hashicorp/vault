package stepwise

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

const pluginPrefix = "vault-plugin-"

// CompilePlugin is a helper method to compile a sourc plugin
// TODO refactor compile plugin input and output to be types
func CompilePlugin(name, pluginName, srcDir, tmpDir string) (string, string, string, error) {
	binName := name
	if !strings.HasPrefix(binName, pluginPrefix) {
		binName = fmt.Sprintf("%s%s", pluginPrefix, binName)
	}
	binPath := path.Join(tmpDir, binName)

	cmd := exec.Command("go", "build", "-o", binPath, path.Join(srcDir, fmt.Sprintf("cmd/%s/main.go", pluginName)))
	var out bytes.Buffer
	cmd.Stdout = &out

	// match the target architecture of the docker container
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	err := cmd.Run()
	if err != nil {
		return "", "", "", err
	}

	// calculate sha256
	f, err := os.Open(binPath)
	if err != nil {
		return "", "", "", err
	}

	h := sha256.New()
	if _, ioErr := io.Copy(h, f); ioErr != nil {
		panic(ioErr)
	}

	_ = f.Close()

	sha256value := fmt.Sprintf("%x", h.Sum(nil))

	return binName, binPath, sha256value, err
}
