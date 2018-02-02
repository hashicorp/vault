package http

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	bplugin "github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/logbridge"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
	"github.com/hashicorp/vault/physical/inmem"
	"github.com/hashicorp/vault/vault"
)

func getPluginClusterAndCore(t testing.TB, logger *logbridge.Logger) (*vault.TestCluster, *vault.TestClusterCore) {
	inmha, err := inmem.NewInmemHA(nil, logger.LogxiLogger())
	if err != nil {
		t.Fatal(err)
	}

	coreConfig := &vault.CoreConfig{
		Physical: inmha,
		LogicalBackends: map[string]logical.Factory{
			"plugin": bplugin.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
		RawLogger:   logger,
	})
	cluster.Start()

	cores := cluster.Cores
	core := cores[0]

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	vault.TestWaitActive(t, core.Core)
	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestPlugin_PluginMain")

	// Mount the mock plugin
	err = core.Client.Sys().Mount("mock", &api.MountInput{
		Type:       "plugin",
		PluginName: "mock-plugin",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cluster, core
}

func TestPlugin_PluginMain(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args := []string{"--ca-cert=" + caPEM}

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	factoryFunc := mock.FactoryType(logical.TypeLogical)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal("Why are we here")
}

func TestPlugin_MockList(t *testing.T) {
	logger := logbridge.NewLogger(hclog.New(&hclog.LoggerOptions{
		Mutex: &sync.Mutex{},
	}))
	cluster, core := getPluginClusterAndCore(t, logger)
	defer cluster.Cleanup()

	_, err := core.Client.Logical().Write("mock/kv/foo", map[string]interface{}{
		"bar": "baz",
	})
	if err != nil {
		t.Fatal(err)
	}

	keys, err := core.Client.Logical().List("mock/kv/")
	if err != nil {
		t.Fatal(err)
	}
	if keys.Data["keys"].([]interface{})[0].(string) != "foo" {
		t.Fatal(keys)
	}

	_, err = core.Client.Logical().Write("mock/kv/zoo", map[string]interface{}{
		"bar": "baz",
	})
	if err != nil {
		t.Fatal(err)
	}

	keys, err = core.Client.Logical().List("mock/kv/")
	if err != nil {
		t.Fatal(err)
	}
	if keys.Data["keys"].([]interface{})[0].(string) != "foo" || keys.Data["keys"].([]interface{})[1].(string) != "zoo" {
		t.Fatal(keys)
	}
}

func TestPlugin_MockRawResponse(t *testing.T) {
	logger := logbridge.NewLogger(hclog.New(&hclog.LoggerOptions{
		Mutex: &sync.Mutex{},
	}))
	cluster, core := getPluginClusterAndCore(t, logger)
	defer cluster.Cleanup()

	resp, err := core.Client.RawRequest(core.Client.NewRequest("GET", "/v1/mock/raw"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body[:]) != "Response" {
		t.Fatal("bad body")
	}

	if resp.StatusCode != 200 {
		t.Fatal("bad status")
	}

}
