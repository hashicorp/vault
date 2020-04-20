package acctest

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
)

// TestHelper is a package global that plugins will use to extract Vault
// and Docker Clients after setup
var TestHelper *Helper

// Helper is intended as a per-package singleton created in TestMain which
// other tests in a package can use to create Terraform execution contexts
type Helper struct {
	// api client for use
	Client  *api.Client
	Cluster *DockerCluster
	// name for plugin in test
	Name string
	// sourceDir is the dir containing the plugin test binary
	SourceDir string
}

// Cleanup cals the Cluster Cleanup method, if Cluster is not nil
func (h *Helper) Cleanup() {
	if h.Cluster != nil {
		h.Cluster.Cleanup()
	}
}

// UseDocker setups docker, copying the plugin test binary
func UseDocker(name, src string) *Helper {
	return &Helper{
		Name:      name,
		SourceDir: src,
	}
}

// Run runs tests after setup is complete. If the package test helper is not
// nil, Run will call the cleanup after tests complete.
func Run(m *testing.M) {
	stat := m.Run()
	if TestHelper != nil {
		TestHelper.Cleanup()
	}
	os.Exit(stat)
}

func compilePlugin(name, srcDir, tmpDir string) (string, string, error) {
	binPath := path.Join(tmpDir, name)

	cmd := exec.Command("go", "build", "-o", path.Join(tmpDir, name), path.Join(srcDir, fmt.Sprintf("cmd/%s/main.go", name)))
	var out bytes.Buffer
	cmd.Stdout = &out
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

	defer func() {
		_ = f.Close()
	}()

	h := sha256.New()
	if _, ioErr := io.Copy(h, f); ioErr != nil {
		panic(ioErr)
	}

	sha256value := fmt.Sprintf("%x", h.Sum(nil))

	return binPath, sha256value, err
}

// Setup creates any temp dir and compiles the binary for copying to Docker
func Setup(name string) error {
	if os.Getenv("VAULT_ACC") == "1" {
		// get the working directory of the plugin being tested.
		srcDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		tmpDir, err := ioutil.TempDir("", "bin")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(tmpDir) // clean up

		binPath, sha256value, err := compilePlugin(name, srcDir, tmpDir)
		if err != nil {
			panic(err)
		}

		coreConfig := &vault.CoreConfig{
			DisableMlock: true,
		}

		dOpts := &DockerClusterOptions{PluginTestBin: binPath}
		cluster, err := NewDockerCluster(fmt.Sprintf("test-%s", name), coreConfig, dOpts)
		if err != nil {
			panic(err)
		}

		cores := cluster.ClusterNodes
		client := cores[0].Client

		TestHelper = &Helper{
			Client:  client,
			Cluster: cluster,
		}

		// use client to mount plugin
		err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
			Name:    name,
			Type:    consts.PluginTypeSecrets,
			Command: name,
			SHA256:  sha256value,
		})
		if err != nil {
			panic(err)
		}
	}
	return nil
}
