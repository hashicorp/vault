package vault

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/plugins/database/postgresql"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	backendplugin "github.com/hashicorp/vault/sdk/plugin"

	"github.com/hashicorp/vault/helper/builtinplugins"
)

func TestPluginCatalog_CRUD(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	sym, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	core.pluginCatalog.directory = sym

	// Get builtin plugin
	p, err := core.pluginCatalog.Get(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedBuiltin := &pluginutil.PluginRunner{
		Name:    "mysql-database-plugin",
		Type:    consts.PluginTypeDatabase,
		Builtin: true,
	}
	expectedBuiltin.BuiltinFactory, _ = builtinplugins.Registry.Get("mysql-database-plugin", consts.PluginTypeDatabase)

	if &(p.BuiltinFactory) == &(expectedBuiltin.BuiltinFactory) {
		t.Fatal("expected BuiltinFactory did not match actual")
	}
	expectedBuiltin.BuiltinFactory = nil
	p.BuiltinFactory = nil
	if !reflect.DeepEqual(p, expectedBuiltin) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", p, expectedBuiltin)
	}

	// Set a plugin, test overwriting a builtin plugin
	file, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	command := fmt.Sprintf("%s", filepath.Base(file.Name()))
	err = core.pluginCatalog.Set(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase, command, []string{"--test"}, []string{"FOO=BAR"}, []byte{'1'})
	if err != nil {
		t.Fatal(err)
	}

	// Get the plugin
	p, err = core.pluginCatalog.Get(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expected := &pluginutil.PluginRunner{
		Name:    "mysql-database-plugin",
		Type:    consts.PluginTypeDatabase,
		Command: filepath.Join(sym, filepath.Base(file.Name())),
		Args:    []string{"--test"},
		Env:     []string{"FOO=BAR"},
		Sha256:  []byte{'1'},
		Builtin: false,
	}

	if !reflect.DeepEqual(p, expected) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", p, expected)
	}

	// Delete the plugin
	err = core.pluginCatalog.Delete(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	// Get builtin plugin
	p, err = core.pluginCatalog.Get(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedBuiltin = &pluginutil.PluginRunner{
		Name:    "mysql-database-plugin",
		Type:    consts.PluginTypeDatabase,
		Builtin: true,
	}
	expectedBuiltin.BuiltinFactory, _ = builtinplugins.Registry.Get("mysql-database-plugin", consts.PluginTypeDatabase)

	if &(p.BuiltinFactory) == &(expectedBuiltin.BuiltinFactory) {
		t.Fatal("expected BuiltinFactory did not match actual")
	}
	expectedBuiltin.BuiltinFactory = nil
	p.BuiltinFactory = nil
	if !reflect.DeepEqual(p, expectedBuiltin) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", p, expectedBuiltin)
	}
}

func TestPluginCatalog_List(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	sym, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	core.pluginCatalog.directory = sym

	// Get builtin plugins and sort them
	builtinKeys := builtinplugins.Registry.Keys(consts.PluginTypeDatabase)
	sort.Strings(builtinKeys)

	// List only builtin plugins
	plugins, err := core.pluginCatalog.List(context.Background(), consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if len(plugins) != len(builtinKeys) {
		t.Fatalf("unexpected length of plugin list, expected %d, got %d", len(builtinKeys), len(plugins))
	}

	for i, p := range builtinKeys {
		if !reflect.DeepEqual(plugins[i], p) {
			t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugins[i], p)
		}
	}

	// Set a plugin, test overwriting a builtin plugin
	file, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	command := filepath.Base(file.Name())
	err = core.pluginCatalog.Set(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase, command, []string{"--test"}, []string{}, []byte{'1'})
	if err != nil {
		t.Fatal(err)
	}

	// Set another plugin
	err = core.pluginCatalog.Set(context.Background(), "aaaaaaa", consts.PluginTypeDatabase, command, []string{"--test"}, []string{}, []byte{'1'})
	if err != nil {
		t.Fatal(err)
	}

	// List the plugins
	plugins, err = core.pluginCatalog.List(context.Background(), consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// plugins has a test-added plugin called "aaaaaaa" that is not built in
	if len(plugins) != len(builtinKeys)+1 {
		t.Fatalf("unexpected length of plugin list, expected %d, got %d", len(builtinKeys)+1, len(plugins))
	}

	// verify the first plugin is the one we just created.
	if !reflect.DeepEqual(plugins[0], "aaaaaaa") {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugins[0], "aaaaaaa")
	}

	// verify the builtin plugins are correct
	for i, p := range builtinKeys {
		if !reflect.DeepEqual(plugins[i+1], p) {
			t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugins[i+1], p)
		}
	}
}

func TestPluginCatalog_NewPluginClient(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	sym, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	core.pluginCatalog.directory = sym

	if extPlugins := len(core.pluginCatalog.externalPlugins); extPlugins != 0 {
		t.Fatalf("expected externalPlugins map to be of len 0 but got %d", extPlugins)
	}

	// register plugins
	TestAddTestPlugin(t, core, "mux-postgres", consts.PluginTypeUnknown, "TestPluginCatalog_PluginMain_PostgresMultiplexed", []string{}, "")
	TestAddTestPlugin(t, core, "single-postgres-1", consts.PluginTypeUnknown, "TestPluginCatalog_PluginMain_Postgres", []string{}, "")
	TestAddTestPlugin(t, core, "single-postgres-2", consts.PluginTypeUnknown, "TestPluginCatalog_PluginMain_Postgres", []string{}, "")

	TestAddTestPlugin(t, core, "mux-userpass", consts.PluginTypeUnknown, "TestPluginCatalog_PluginMain_UserpassMultiplexed", []string{}, "")
	TestAddTestPlugin(t, core, "single-userpass-1", consts.PluginTypeUnknown, "TestPluginCatalog_PluginMain_Userpass", []string{}, "")
	TestAddTestPlugin(t, core, "single-userpass-2", consts.PluginTypeUnknown, "TestPluginCatalog_PluginMain_Userpass", []string{}, "")

	var cleanupFuncs []func()
	// run plugins
	// run "mux-postgres" twice which will start a single plugin for 2
	// distinct connections
	c := TestRunTestPlugin(t, core, consts.PluginTypeDatabase, "mux-postgres")
	cleanupFuncs = append(cleanupFuncs, c)
	c = TestRunTestPlugin(t, core, consts.PluginTypeDatabase, "mux-postgres")
	cleanupFuncs = append(cleanupFuncs, c)
	c = TestRunTestPlugin(t, core, consts.PluginTypeDatabase, "single-postgres-1")
	cleanupFuncs = append(cleanupFuncs, c)
	c = TestRunTestPlugin(t, core, consts.PluginTypeDatabase, "single-postgres-2")
	cleanupFuncs = append(cleanupFuncs, c)

	// run "mux-userpass" twice which will start a single plugin for 2
	// distinct connections
	c = TestRunTestPlugin(t, core, consts.PluginTypeCredential, "mux-userpass")
	cleanupFuncs = append(cleanupFuncs, c)
	c = TestRunTestPlugin(t, core, consts.PluginTypeCredential, "mux-userpass")
	cleanupFuncs = append(cleanupFuncs, c)
	c = TestRunTestPlugin(t, core, consts.PluginTypeCredential, "single-userpass-1")
	cleanupFuncs = append(cleanupFuncs, c)
	c = TestRunTestPlugin(t, core, consts.PluginTypeCredential, "single-userpass-2")
	cleanupFuncs = append(cleanupFuncs, c)

	externalPlugins := core.pluginCatalog.externalPlugins
	if len(externalPlugins) != 6 {
		t.Fatalf("expected externalPlugins map to be of len 6 but got %d", len(externalPlugins))
	}

	// check connections map
	expectConnectionLen(t, 2, externalPlugins["mux-postgres"].connections)
	expectConnectionLen(t, 1, externalPlugins["single-postgres-1"].connections)
	expectConnectionLen(t, 1, externalPlugins["single-postgres-2"].connections)
	expectConnectionLen(t, 2, externalPlugins["mux-userpass"].connections)
	expectConnectionLen(t, 1, externalPlugins["single-userpass-1"].connections)
	expectConnectionLen(t, 1, externalPlugins["single-userpass-2"].connections)

	// check multiplexing support
	expectMultiplexingSupport(t, true, externalPlugins["mux-postgres"].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins["single-postgres-1"].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins["single-postgres-2"].multiplexingSupport)
	expectMultiplexingSupport(t, true, externalPlugins["mux-userpass"].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins["single-userpass-1"].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins["single-userpass-2"].multiplexingSupport)

	// cleanup all of the external plugin processes
	for _, cleanupFunc := range cleanupFuncs {
		cleanupFunc()
	}

	// check connections map
	if len(externalPlugins) != 0 {
		t.Fatalf("expected external plugin map to be of len 0 but got %d", len(externalPlugins))
	}
}

func TestPluginCatalog_PluginMain_Userpass(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := backendplugin.Serve(
		&backendplugin.ServeOpts{
			BackendFactoryFunc: userpass.Factory,
			TLSProviderFunc:    tlsProviderFunc,
		},
	)
	if err != nil {
		t.Fatalf("Failed to initialize userpass: %s", err)
	}
}

func TestPluginCatalog_PluginMain_UserpassMultiplexed(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := backendplugin.ServeMultiplex(
		&backendplugin.ServeOpts{
			BackendFactoryFunc: userpass.Factory,
			TLSProviderFunc:    tlsProviderFunc,
		},
	)
	if err != nil {
		t.Fatalf("Failed to initialize userpass: %s", err)
	}
}

func TestPluginCatalog_PluginMain_Postgres(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	dbType, err := postgresql.New()
	if err != nil {
		t.Fatalf("Failed to initialize postgres: %s", err)
	}

	v5.Serve(dbType.(v5.Database))
}

func TestPluginCatalog_PluginMain_PostgresMultiplexed(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	v5.ServeMultiplex(postgresql.New)
}

// expectConnectionLen asserts that the PluginCatalog's externalPlugin
// connections map has a length of expectedLen
func expectConnectionLen(t *testing.T, expectedLen int, connections map[string]*pluginClient) {
	if len(connections) != expectedLen {
		t.Fatalf("expected external plugin's connections map to be of len %d but got %d", expectedLen, len(connections))
	}
}

func expectMultiplexingSupport(t *testing.T, expected, actual bool) {
	if expected != actual {
		t.Fatalf("expected external plugin multiplexing support to be %t", expected)
	}
}
