package stepwise

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/hashicorp/vault/api"
)

// TestHelper is a package global that plugins will use to extract Vault
// and Docker Clients after setup
var TestHelper *Helper

// Helper is intended as a per-package singleton created in TestMain which
// other tests in a package can use to create Terraform execution contexts
type Helper struct {
	// api client for use
	Client *api.Client
	// Cluster *dockerDriver.DockerCluster
	// name for plugin in test
	Name string
	// sourceDir is the dir containing the plugin test binary
	SourceDir string
	// temp dir where plugin is compiled
	BuildDir string
}

// Cleanup cals the Cluster Cleanup method, if Cluster is not nil
func (h *Helper) Cleanup() {
	// if h.Cluster != nil {
	// 	h.Cluster.Cleanup()
	// }
	if h.BuildDir != "" {
		_ = os.RemoveAll(h.BuildDir) // clean up
	}
}

// UseDocker setups docker, copying the plugin test binary
func UseDocker(name, src string) *Helper {
	return &Helper{
		Name:      name,
		SourceDir: src,
	}
}

// TODO remove this Run function
// Run runs tests after setup is complete. If the package test helper is not
// nil, Run will call the cleanup after tests complete.
// func Run(m *testing.M) {
// 	stat := m.Run()
// 	if TestHelper != nil {
// 		TestHelper.Cleanup()
// 	}
// 	os.Exit(stat)
// }

// CompilePlugin is a helper method to compile a sourc plugin
func CompilePlugin(name, srcDir, tmpDir string) (string, string, error) {
	binPath := path.Join(tmpDir, name)

	cmd := exec.Command("go", "build", "-o", path.Join(tmpDir, name), path.Join(srcDir, fmt.Sprintf("cmd/%s/main.go", name)))
	var out bytes.Buffer
	cmd.Stdout = &out

	// match the target architecture of the docker container
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}

	// calculate sha256
	f, err := os.Open(binPath)
	if err != nil {
		return "", "", err
	}

	h := sha256.New()
	if _, ioErr := io.Copy(h, f); ioErr != nil {
		panic(ioErr)
	}

	_ = f.Close()

	sha256value := fmt.Sprintf("%x", h.Sum(nil))
	// q.Q("=> compiled, sha:", name, sha256value)

	return binPath, sha256value, err
}
