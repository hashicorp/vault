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

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/plugins/database/postgresql"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"

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
	testPlugin := "mux-postgresql-database-plugin"

	sym, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	core.pluginCatalog.directory = sym

	TestAddTestPlugin(t, core, testPlugin, consts.PluginTypeDatabase, "TestPluginCatalog_PluginMain_PostgresMultiplexed", []string{}, "")

	config := pluginutil.PluginClientConfig{
		Name:            testPlugin,
		PluginType:      consts.PluginTypeDatabase,
		PluginSets:      v5.PluginSets,
		HandshakeConfig: v5.HandshakeConfig,
		Logger:          log.NewNullLogger(),
		IsMetadataMode:  false,
		AutoMTLS:        true,
	}

	resp, err := core.pluginCatalog.NewPluginClient(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\n\nresp: %#+v\n\n", resp)

	externalPlugins := core.pluginCatalog.externalPlugins
	extPluginLen := len(externalPlugins)
	if extPluginLen != 1 {
		t.Fatalf("expected 1 external plugin but got %d", extPluginLen)
	}
	if !externalPlugins[testPlugin].multiplexingSupport {
		t.Fatalf("expected external plugin to be multiplexed")
	}
}

// func TestPluginCatalog_PluginMain_Postgres(t *testing.T) {
// 	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
// 		return
// 	}

// 	v5.Serve(postgresql.New)
// }

func TestPluginCatalog_PluginMain_PostgresMultiplexed(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	v5.ServeMultiplex(postgresql.New)
}
