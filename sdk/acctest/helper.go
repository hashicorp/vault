package acctest

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
	"github.com/y0ssar1an/q"
)

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
	q.Q("==> acctest.Run start")
	stat := m.Run()
	if TestHelper != nil {
		q.Q("==> ==> acctest.Run Cleanup")
		if TestHelper.Cluster != nil {
			q.Q("==> ==> acctest.Run Cleanup")
			TestHelper.Cluster.Cleanup()
		} else {
			q.Q("==> ==> acctest.Run Cluster was nil")
		}
	} else {
		q.Q("==> ==> acctest Helper nil")
	}
	q.Q("==> acctest.Run finish")
	os.Exit(stat)
}

// Setup creates any temp dir and compiles the binary for copying to Docker
func Setup(name string) error {
	if os.Getenv("VAULT_ACC") == "1" {
		absPluginExecPath, _ := filepath.Abs(os.Args[0])
		pluginName := path.Base(absPluginExecPath)
		os.Link(absPluginExecPath, path.Join("/Users/clint/Desktop/plugins", pluginName))
		// setup docker, send src and name
		// run tests
		coreConfig := &vault.CoreConfig{
			DisableMlock: true,
		}
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		wd = path.Join(wd, "vault/plugins/uuid")
		cmd := exec.Command("go", "build", "-o", "./vault/plugins/uuid", "/Users/clint/go-src/github.com/catsby/vault-plugin-secrets-uuid/cmd/uuid/main.go")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		// cluster, err := acctest.NewDockerCluster(t.Name(), coreConfig, nil)
		// dOpts := &acctest.DockerClusterOptions{PluginTestBin: absPluginExecPath}
		//TODO: cleanup working dir
		dOpts := &DockerClusterOptions{PluginTestBin: wd}
		cluster, err := NewDockerCluster("test-uuid", coreConfig, dOpts)
		if err != nil {
			panic(err)
		}

		cores := cluster.ClusterNodes
		client := cores[0].Client
		// calculate sha256 of binary/vault/plugins/uuid
		pPath := "/Users/clint/go-src/github.com/catsby/vault-plugin-secrets-uuid/vault/plugins/uuid"

		f, err := os.Open(pPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			panic(err)
		}

		sha256value := fmt.Sprintf("%x", h.Sum(nil))

		TestHelper = &Helper{
			Client:  client,
			Cluster: cluster,
		}
		// use client to mount plugin

		err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
			Name:    "uuid",
			Type:    consts.PluginTypeSecrets,
			Command: "uuid",
			SHA256:  sha256value,
		})
		if err != nil {
			panic(err)
		}

		// run tests
		// stat := m.Run()
		// cluster.Cleanup()
		// os.Exit(stat)
	}
	return nil
}
