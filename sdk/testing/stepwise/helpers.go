package stepwise

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

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	dockerDriver "github.com/hashicorp/vault/sdk/testing/stepwise/drivers/docker"
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
	Cluster *dockerDriver.DockerCluster
	// name for plugin in test
	Name string
	// sourceDir is the dir containing the plugin test binary
	SourceDir string
	// temp dir where plugin is compiled
	buildDir string
}

// Cleanup cals the Cluster Cleanup method, if Cluster is not nil
func (h *Helper) Cleanup() {
	if h.Cluster != nil {
		h.Cluster.Cleanup()
	}
	if h.buildDir != "" {
		_ = os.RemoveAll(h.buildDir) // clean up
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

func compilePlugin(name, srcDir, tmpDir string) (string, string, error) {
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

// Setup creates any temp dir and compiles the binary for copying to Docker
func Setup(name string) error {
	if os.Getenv("VAULT_ACC") == "1" {
		// get the working directory of the plugin being tested.
		srcDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// tmpDir gets cleaned up when the cluster is cleaned up
		tmpDir, err := ioutil.TempDir("", "bin")
		if err != nil {
			log.Fatal(err)
		}

		binPath, sha256value, err := compilePlugin(name, srcDir, tmpDir)
		if err != nil {
			panic(err)
		}

		coreConfig := &vault.CoreConfig{
			DisableMlock: true,
		}

		dOpts := &dockerDriver.DockerClusterOptions{PluginTestBin: binPath}
		cluster, err := dockerDriver.NewDockerCluster(fmt.Sprintf("test-%s", name), coreConfig, dOpts)
		if err != nil {
			panic(err)
		}

		cores := cluster.ClusterNodes
		client := cores[0].Client

		TestHelper = &Helper{
			Client:   client,
			Cluster:  cluster,
			buildDir: tmpDir,
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
