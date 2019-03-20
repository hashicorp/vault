package plugin_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/pluginutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	logicalPlugin "github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
	"github.com/hashicorp/vault/vault"
)

func TestBackend_impl(t *testing.T) {
	var _ logical.Backend = &plugin.PluginBackend{}
}

func TestBackend(t *testing.T) {
	config, cleanup := testConfig(t)
	defer cleanup()

	_, err := plugin.Backend(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_Factory(t *testing.T) {
	config, cleanup := testConfig(t)
	defer cleanup()

	_, err := plugin.Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMain(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	err := logicalPlugin.Serve(&logicalPlugin.ServeOpts{
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
		Logger: logging.NewVaultLogger(log.Debug),
		System: sys,
		Config: map[string]string{
			"plugin_name": "mock-plugin",
			"plugin_type": "database",
		},
	}

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeDatabase, "TestBackend_PluginMain", []string{}, "")

	return config, func() {
		cluster.Cleanup()
	}
}
