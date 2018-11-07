package http

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	bplugin "github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
	"github.com/hashicorp/vault/vault"
)

func getPluginClusterAndCore(t testing.TB, logger log.Logger) (*vault.TestCluster, *vault.TestClusterCore) {
	inm, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	coreConfig := &vault.CoreConfig{
		Physical:   inm,
		HAPhysical: inmha.(physical.HABackend),
		LogicalBackends: map[string]logical.Factory{
			"plugin": bplugin.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
		Logger:      logger.Named("testclusteroptions"),
	})
	cluster.Start()

	cores := cluster.Cores
	core := cores[0]

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	vault.TestWaitActive(t, core.Core)
	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "TestPlugin_PluginMain", []string{}, "")

	// Mount the mock plugin
	err = core.Client.Sys().Mount("mock", &api.MountInput{
		Type: "mock-plugin",
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
	logger := log.New(&log.LoggerOptions{
		Mutex: &sync.Mutex{},
	})
	cluster, core := getPluginClusterAndCore(t, logger)
	defer cluster.Cleanup()

	_, err := core.Client.Logical().Write("mock/kv/foo", map[string]interface{}{
		"value": "baz",
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
		"value": "baz",
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
	logger := log.New(&log.LoggerOptions{
		Mutex: &sync.Mutex{},
	})
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

func TestPlugin_GetParams(t *testing.T) {
	logger := log.New(&log.LoggerOptions{
		Mutex: &sync.Mutex{},
	})
	cluster, core := getPluginClusterAndCore(t, logger)
	defer cluster.Cleanup()

	_, err := core.Client.Logical().Write("mock/kv/foo", map[string]interface{}{
		"value": "baz",
	})
	if err != nil {
		t.Fatal(err)
	}

	r := core.Client.NewRequest("GET", "/v1/mock/kv/foo")
	r.Params.Add("version", "12")
	resp, err := core.Client.RawRequest(r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]interface{}{
		"value":   "baz",
		"version": json.Number("12"),
	}

	if !reflect.DeepEqual(secret.Data, expected) {
		t.Fatal(secret.Data)
	}
}
