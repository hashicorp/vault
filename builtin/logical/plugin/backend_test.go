package plugin

import (
	"io/ioutil"
	stdhttp "net/http"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/plugins/backend/mock"
	"github.com/hashicorp/vault/vault"
)

func TestBackend(t *testing.T) {
	config := testConfig(t)

	_, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_Factory(t *testing.T) {
	config := testConfig(t)

	_, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMain(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	content := []byte(vault.TestClusterCACert)
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
	plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
}

func testConfig(t *testing.T) *logical.BackendConfig {
	coreConfig := &vault.CoreConfig{}

	handler1 := stdhttp.NewServeMux()
	handler2 := stdhttp.NewServeMux()
	handler3 := stdhttp.NewServeMux()

	// Chicken-and-egg: Handler needs a core. So we create handlers first, then
	// add routes chained to a Handler-created handler.
	cores := vault.TestCluster(t, []stdhttp.Handler{handler1, handler2, handler3}, coreConfig, false)
	handler1.Handle("/", http.Handler(cores[0].Core))
	handler2.Handle("/", http.Handler(cores[1].Core))
	handler3.Handle("/", http.Handler(cores[2].Core))

	core := cores[0]

	sys := vault.TestDynamicSystemView(core.Core)

	config := &logical.BackendConfig{
		Logger: nil,
		System: sys,
		Config: map[string]string{
			"plugin_name": "mock-plugin",
		},
	}

	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestBackend_PluginMain")

	return config
}
