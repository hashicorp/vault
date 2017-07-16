package vault_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	lplugin "github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
	"github.com/hashicorp/vault/vault"
	log "github.com/mgutz/logxi/v1"
)

func TestSystemBackend_enableAuth_plugin(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, true)
	cluster.StartListeners()
	defer cluster.CloseListeners()
	cores := cluster.Cores

	cores[0].Handler.Handle("/", http.Handler(cores[0].Core))
	cores[1].Handler.Handle("/", http.Handler(cores[1].Core))
	cores[2].Handler.Handle("/", http.Handler(cores[2].Core))

	core := cores[0]

	b := vault.NewSystemBackend(core.Core)
	logger := logformat.NewVaultLogger(log.LevelTrace)
	bc := &logical.BackendConfig{
		Logger: logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}

	err := b.Backend.Setup(bc)
	if err != nil {
		t.Fatal(err)
	}

	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestBackend_PluginMain")

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/mock-plugin")
	req.Data["type"] = "plugin"
	req.Data["plugin_name"] = "mock-plugin"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
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

	factoryFunc := mock.FactoryType(logical.TypeCredential)

	args := []string{"--ca-cert=" + tmpfile.Name()}

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)
	err = lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}
