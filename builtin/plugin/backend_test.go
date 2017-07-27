package plugin

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/pluginutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
	"github.com/hashicorp/vault/vault"
)

func TestBackend(t *testing.T) {
	config, cleanup := testConfig(t)
	defer cleanup()

	_, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_Factory(t *testing.T) {
	config, cleanup := testConfig(t)
	defer cleanup()

	_, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMain(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	caPem := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPem == "" {
		t.Fatal("CA cert not passed in")
	}

	content := []byte(caPem)
	tmpfile, err := ioutil.TempFile("", "test-cacert")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	args := []string{"--ca-cert=" + tmpfile.Name()}

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	err = plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func testConfig(t *testing.T) (*logical.BackendConfig, func()) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	cores := cluster.Cores

	core := cores[0]

	sys := vault.TestDynamicSystemView(core.Core)

	config := &logical.BackendConfig{
		Logger: nil,
		System: sys,
		Config: map[string]string{
			"plugin_name": "mock-plugin",
		},
	}

	os.Setenv(pluginutil.PluginCACertPEMEnv, string(cluster.CACertPEM))

	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestBackend_PluginMain")

	return config, func() {
		cluster.Cleanup()
	}
}
