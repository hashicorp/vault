// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"

	"github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/plugins/database/postgresql"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginruntimeutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	backendplugin "github.com/hashicorp/vault/sdk/plugin"
)

func testPluginCatalog(t *testing.T) *PluginCatalog {
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Trace,
	})
	storage, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	testDir, err := filepath.EvalSymlinks(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Fatal(err)
	}
	pluginRuntimeCatalog := testPluginRuntimeCatalog(t)
	pluginCatalog, err := SetupPluginCatalog(
		context.Background(),
		&PluginCatalogInput{
			Logger:               logger,
			BuiltinRegistry:      corehelpers.NewMockBuiltinRegistry(),
			CatalogView:          logical.NewLogicalStorage(storage),
			PluginDirectory:      testDir,
			EnableMlock:          false,
			PluginRuntimeCatalog: pluginRuntimeCatalog,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	return pluginCatalog
}

type warningCountLogger struct {
	log.Logger
	warnings int
}

func (l *warningCountLogger) Warn(msg string, args ...interface{}) {
	l.warnings++
	l.Logger.Warn(msg, args...)
}

func (l *warningCountLogger) reset() {
	l.warnings = 0
}

// TestPluginCatalog_SetupPluginCatalog_WarningsWithLegacyEnvSetting ensures we
// log the correct number of warnings during plugin catalog setup (which is run
// during unseal) if users have set the flag to keep old behavior. This is to
// help users migrate safely to the new default behavior.
func TestPluginCatalog_SetupPluginCatalog_WarningsWithLegacyEnvSetting(t *testing.T) {
	logger := &warningCountLogger{
		Logger: log.New(&hclog.LoggerOptions{
			Level: hclog.Trace,
		}),
	}
	storage, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	logicalStorage := logical.NewLogicalStorage(storage)

	// prefix to avoid collisions with other tests.
	const prefix = "TEST_PLUGIN_CATALOG_ENV_"
	plugin := &pluginutil.PluginRunner{
		Name:    "mysql-database-plugin",
		Type:    consts.PluginTypeDatabase,
		Version: "1.0.0",
		Env: []string{
			prefix + "A=1",
			prefix + "VALUE_WITH_EQUALS=1=2",
			prefix + "EMPTY_VALUE=",
		},
	}

	// Insert a plugin into storage before the catalog is setup.
	buf, err := json.Marshal(plugin)
	if err != nil {
		t.Fatal(err)
	}
	logicalEntry := logical.StorageEntry{
		Key:   path.Join(plugin.Type.String(), plugin.Name, plugin.Version),
		Value: buf,
	}
	if err := logicalStorage.Put(context.Background(), &logicalEntry); err != nil {
		t.Fatal(err)
	}

	for name, tc := range map[string]struct {
		sysEnv           map[string]string
		expectedWarnings int
	}{
		"no env": {},
		"colliding env, no flag": {
			sysEnv: map[string]string{
				prefix + "A": "10",
			},
			expectedWarnings: 0,
		},
		"colliding env, with flag": {
			sysEnv: map[string]string{
				pluginutil.PluginUseLegacyEnvLayering: "true",
				prefix + "A":                          "10",
			},
			expectedWarnings: 2,
		},
		"all colliding env, with flag": {
			sysEnv: map[string]string{
				pluginutil.PluginUseLegacyEnvLayering: "true",
				prefix + "A":                          "10",
				prefix + "VALUE_WITH_EQUALS":          "1=2",
				prefix + "EMPTY_VALUE":                "",
			},
			expectedWarnings: 4,
		},
	} {
		t.Run(name, func(t *testing.T) {
			logger.reset()
			for k, v := range tc.sysEnv {
				t.Setenv(k, v)
			}

			_, err := SetupPluginCatalog(
				context.Background(),
				&PluginCatalogInput{
					Logger:               logger,
					BuiltinRegistry:      corehelpers.NewMockBuiltinRegistry(),
					CatalogView:          logicalStorage,
					PluginDirectory:      "",
					Tmpdir:               "",
					EnableMlock:          false,
					PluginRuntimeCatalog: nil,
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expectedWarnings != logger.warnings {
				t.Fatalf("expected %d warnings, got %d", tc.expectedWarnings, logger.warnings)
			}
		})
	}
}

func TestPluginCatalog_CRUD(t *testing.T) {
	const pluginName = "mysql-database-plugin"

	pluginCatalog := testPluginCatalog(t)

	// Register a fake plugin in the catalog.
	file, err := os.CreateTemp(pluginCatalog.directory, "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    pluginName,
		Type:    consts.PluginTypeDatabase,
		Version: "1.0.0",
		Command: filepath.Base(file.Name()),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Register a pinned version, should not affect anything below.
	err = pluginCatalog.SetPinnedVersion(context.Background(), &pluginutil.PinnedVersion{
		Name:    pluginName,
		Type:    consts.PluginTypeDatabase,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get builtin plugin
	p, err := pluginCatalog.Get(context.Background(), pluginName, consts.PluginTypeDatabase, "")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// Get it again, explicitly specifying builtin version
	builtinVersion := versions.GetBuiltinVersion(consts.PluginTypeDatabase, pluginName)
	p2, err := pluginCatalog.Get(context.Background(), pluginName, consts.PluginTypeDatabase, builtinVersion)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedBuiltin := &pluginutil.PluginRunner{
		Name:    pluginName,
		Type:    consts.PluginTypeDatabase,
		Builtin: true,
		Version: builtinVersion,
	}
	expectedBuiltin.BuiltinFactory, _ = builtinplugins.Registry.Get(pluginName, consts.PluginTypeDatabase)

	if &(p.BuiltinFactory) == &(expectedBuiltin.BuiltinFactory) {
		t.Fatal("expected BuiltinFactory did not match actual")
	}
	expectedBuiltin.BuiltinFactory = nil
	p.BuiltinFactory = nil
	p2.BuiltinFactory = nil
	if !reflect.DeepEqual(p, expectedBuiltin) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", p, expectedBuiltin)
	}
	if !reflect.DeepEqual(p2, expectedBuiltin) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", p2, expectedBuiltin)
	}

	// Set a plugin, test overwriting a builtin plugin
	command := filepath.Base(file.Name())
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    pluginName,
		Type:    consts.PluginTypeDatabase,
		Version: "",
		Command: command,
		Args:    []string{"--test"},
		Env:     []string{"FOO=BAR"},
		Sha256:  []byte{'1'},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get the plugin
	p, err = pluginCatalog.Get(context.Background(), pluginName, consts.PluginTypeDatabase, "")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// Get it again, explicitly specifying builtin version.
	// This time it should fail because it was overwritten.
	p2, err = pluginCatalog.Get(context.Background(), pluginName, consts.PluginTypeDatabase, builtinVersion)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if p2 != nil {
		t.Fatalf("expected no result, got: %#v", p2)
	}

	expected := &pluginutil.PluginRunner{
		Name:    pluginName,
		Type:    consts.PluginTypeDatabase,
		Command: filepath.Join(pluginCatalog.directory, filepath.Base(file.Name())),
		Args:    []string{"--test"},
		Env:     []string{"FOO=BAR"},
		Sha256:  []byte{'1'},
		Builtin: false,
		Version: "",
	}

	if !reflect.DeepEqual(p, expected) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", p, expected)
	}

	// Delete the plugin
	err = pluginCatalog.Delete(context.Background(), pluginName, consts.PluginTypeDatabase, "")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	// Get builtin plugin
	p, err = pluginCatalog.Get(context.Background(), pluginName, consts.PluginTypeDatabase, "")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedBuiltin = &pluginutil.PluginRunner{
		Name:    pluginName,
		Type:    consts.PluginTypeDatabase,
		Builtin: true,
		Version: versions.GetBuiltinVersion(consts.PluginTypeDatabase, pluginName),
	}
	expectedBuiltin.BuiltinFactory, _ = builtinplugins.Registry.Get(pluginName, consts.PluginTypeDatabase)

	if &(p.BuiltinFactory) == &(expectedBuiltin.BuiltinFactory) {
		t.Fatal("expected BuiltinFactory did not match actual")
	}
	expectedBuiltin.BuiltinFactory = nil
	p.BuiltinFactory = nil
	if !reflect.DeepEqual(p, expectedBuiltin) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", p, expectedBuiltin)
	}
}

func TestPluginCatalog_VersionedCRUD(t *testing.T) {
	pluginCatalog := testPluginCatalog(t)

	// Set a versioned plugin.
	file, err := os.CreateTemp(pluginCatalog.directory, "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	const name = "mysql-database-plugin"
	const version = "1.0.0"
	command := fmt.Sprintf("%s", filepath.Base(file.Name()))
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    name,
		Type:    consts.PluginTypeDatabase,
		Version: version,
		Command: command,
		Args:    []string{"--test"},
		Env:     []string{"FOO=BAR"},
		Sha256:  []byte{'1'},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get the plugin
	plugin, err := pluginCatalog.Get(context.Background(), name, consts.PluginTypeDatabase, version)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expected := &pluginutil.PluginRunner{
		Name:    name,
		Type:    consts.PluginTypeDatabase,
		Version: version,
		Command: filepath.Join(pluginCatalog.directory, filepath.Base(file.Name())),
		Args:    []string{"--test"},
		Env:     []string{"FOO=BAR"},
		Sha256:  []byte{'1'},
		Builtin: false,
	}

	if !reflect.DeepEqual(plugin, expected) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugin, expected)
	}

	// Also get the builtin version to check we can still access that.
	builtinVersion := versions.GetBuiltinVersion(consts.PluginTypeDatabase, name)
	plugin, err = pluginCatalog.Get(context.Background(), name, consts.PluginTypeDatabase, builtinVersion)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expected = &pluginutil.PluginRunner{
		Name:    name,
		Type:    consts.PluginTypeDatabase,
		Version: builtinVersion,
		Builtin: true,
	}

	// Check by marshalling to JSON to avoid messing with BuiltinFactory function field.
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	actualBytes, err := json.Marshal(plugin)
	if err != nil {
		t.Fatal(err)
	}
	if string(expectedBytes) != string(actualBytes) {
		t.Fatalf("expected %s, got %s", string(expectedBytes), string(actualBytes))
	}
	if !plugin.Builtin {
		t.Fatal("expected builtin true but got false")
	}

	// Delete the plugin
	err = pluginCatalog.Delete(context.Background(), name, consts.PluginTypeDatabase, version)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	// Get plugin - should fail
	plugin, err = pluginCatalog.Get(context.Background(), name, consts.PluginTypeDatabase, version)
	if err != nil {
		t.Fatal(err)
	}
	if plugin != nil {
		t.Fatalf("expected no plugin with this version to be in the catalog, but found %+v", plugin)
	}
}

func TestPluginCatalog_List(t *testing.T) {
	pluginCatalog := testPluginCatalog(t)

	// Get builtin plugins and sort them
	builtinKeys := builtinplugins.Registry.Keys(consts.PluginTypeDatabase)
	sort.Strings(builtinKeys)

	// List only builtin plugins
	plugins, err := pluginCatalog.List(context.Background(), consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sort.Strings(plugins)

	if len(plugins) != len(builtinKeys) {
		t.Fatalf("unexpected length of plugin list, expected %d, got %d", len(builtinKeys), len(plugins))
	}

	if !reflect.DeepEqual(plugins, builtinKeys) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugins, builtinKeys)
	}

	// Set a plugin, test overwriting a builtin plugin
	file, err := os.CreateTemp(pluginCatalog.directory, "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	command := filepath.Base(file.Name())
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    "mysql-database-plugin",
		Type:    consts.PluginTypeDatabase,
		Version: "",
		Command: command,
		Args:    []string{"--test"},
		Env:     []string{},
		Sha256:  []byte{'1'},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set another plugin
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    "aaaaaaa",
		Type:    consts.PluginTypeDatabase,
		Version: "",
		Command: command,
		Args:    []string{"--test"},
		Env:     []string{},
		Sha256:  []byte{'1'},
	})
	if err != nil {
		t.Fatal(err)
	}

	// List the plugins
	plugins, err = pluginCatalog.List(context.Background(), consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sort.Strings(plugins)

	// plugins has a test-added plugin called "aaaaaaa" that is not built in
	if len(plugins) != len(builtinKeys)+1 {
		t.Fatalf("unexpected length of plugin list, expected %d, got %d", len(builtinKeys)+1, len(plugins))
	}

	// verify the first plugin is the one we just created.
	if !reflect.DeepEqual(plugins[0], "aaaaaaa") {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugins[0], "aaaaaaa")
	}

	// verify the builtin plugins are correct
	if !reflect.DeepEqual(plugins[1:], builtinKeys) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugins[1:], builtinKeys)
	}
}

func TestPluginCatalog_ListVersionedPlugins(t *testing.T) {
	pluginCatalog := testPluginCatalog(t)

	// Get builtin plugins and sort them
	builtinKeys := builtinplugins.Registry.Keys(consts.PluginTypeDatabase)
	sort.Strings(builtinKeys)

	// List only builtin plugins
	plugins, err := pluginCatalog.ListVersionedPlugins(context.Background(), consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sortVersionedPlugins(plugins)

	if len(plugins) != len(builtinKeys) {
		t.Fatalf("unexpected length of plugin list, expected %d, got %d", len(builtinKeys), len(plugins))
	}

	for i, plugin := range plugins {
		if plugin.Name != builtinKeys[i] {
			t.Fatalf("expected plugin list with names %v but got %+v", builtinKeys, plugins)
		}
	}

	// Set a plugin, test overwriting a builtin plugin
	file, err := ioutil.TempFile(pluginCatalog.directory, "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	command := filepath.Base(file.Name())
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    "mysql-database-plugin",
		Type:    consts.PluginTypeDatabase,
		Version: "",
		Command: command,
		Args:    []string{"--test"},
		Env:     []string{},
		Sha256:  []byte{'1'},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set another plugin, with version information
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    "aaaaaaa",
		Type:    consts.PluginTypeDatabase,
		Version: "1.1.0",
		Command: command,
		Args:    []string{"--test"},
		Env:     []string{},
		Sha256:  []byte{'1'},
	})
	if err != nil {
		t.Fatal(err)
	}

	// List the plugins
	plugins, err = pluginCatalog.ListVersionedPlugins(context.Background(), consts.PluginTypeDatabase)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sortVersionedPlugins(plugins)

	// plugins has a test-added plugin called "aaaaaaa" that is not built in
	if len(plugins) != len(builtinKeys)+1 {
		t.Fatalf("unexpected length of plugin list, expected %d, got %d", len(builtinKeys)+1, len(plugins))
	}

	// verify the first plugin is the one we just created.
	if !reflect.DeepEqual(plugins[0].Name, "aaaaaaa") {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", plugins[0], "aaaaaaa")
	}
	if plugins[0].SemanticVersion == nil {
		t.Fatalf("expected non-nil semantic version for %v", plugins[0].Name)
	}

	// verify the builtin plugins are correct
	for i, plugin := range plugins[1:] {
		if plugin.Name != builtinKeys[i] {
			t.Fatalf("expected plugin list with names %v but got %+v", builtinKeys, plugins)
		}
		switch plugin.Name {
		case "mysql-database-plugin":
			if plugin.Builtin {
				t.Fatalf("expected %v plugin to be an unversioned external plugin", plugin)
			}
			if plugin.Version != "" {
				t.Fatalf("expected no version information for %v but got %s", plugin, plugin.Version)
			}
		default:
			if !plugin.Builtin {
				t.Fatalf("expected %v plugin to be builtin", plugin)
			}
			if !versions.IsBuiltinVersion(plugin.Version) {
				t.Fatalf("expected +builtin metadata but got %s", plugin.Version)
			}
		}

		if plugin.SemanticVersion == nil {
			t.Fatalf("expected non-nil semantic version for %v", plugin)
		}
	}
}

func TestPluginCatalog_ListHandlesPluginNamesWithSlashes(t *testing.T) {
	pluginCatalog := testPluginCatalog(t)

	file, err := os.CreateTemp(pluginCatalog.directory, "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	command := filepath.Base(file.Name())
	ctx := context.Background()

	pluginsToRegister := []pluginutil.PluginRunner{
		{
			Name: "unversioned-plugin",
		},
		{
			Name: "unversioned-plugin/with-slash",
		},
		{
			Name: "unversioned-plugin/with-two/slashes",
		},
		{
			Name:    "versioned-plugin",
			Version: "v1.0.0",
		},
		{
			Name:    "versioned-plugin/with-slash",
			Version: "v1.0.0",
		},
		{
			Name:    "versioned-plugin/with-two/slashes",
			Version: "v1.0.0",
		},
	}
	for _, entry := range pluginsToRegister {
		err = pluginCatalog.Set(ctx, pluginutil.SetPluginInput{
			Name:    entry.Name,
			Type:    consts.PluginTypeCredential,
			Version: entry.Version,
			Command: command,
			Args:    nil,
			Env:     nil,
			Sha256:  nil,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	plugins, err := pluginCatalog.ListVersionedPlugins(ctx, consts.PluginTypeCredential)
	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range pluginsToRegister {
		found := false
		for _, plugin := range plugins {
			if expected.Name == plugin.Name && expected.Version == plugin.Version {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Did not find %#v in %#v", expected, plugins)
		}
	}
}

func TestPluginCatalog_NewPluginClient(t *testing.T) {
	pluginCatalog := testPluginCatalog(t)

	if extPlugins := len(pluginCatalog.externalPlugins); extPlugins != 0 {
		t.Fatalf("expected externalPlugins map to be of len 0 but got %d", extPlugins)
	}

	// register plugins
	TestAddTestPlugin(t, pluginCatalog, "mux-postgres", consts.PluginTypeUnknown, "", "TestPluginCatalog_PluginMain_PostgresMultiplexed", []string{})
	TestAddTestPlugin(t, pluginCatalog, "single-postgres-1", consts.PluginTypeUnknown, "", "TestPluginCatalog_PluginMain_Postgres", []string{})
	TestAddTestPlugin(t, pluginCatalog, "single-postgres-2", consts.PluginTypeUnknown, "", "TestPluginCatalog_PluginMain_Postgres", []string{})

	TestAddTestPlugin(t, pluginCatalog, "mux-userpass", consts.PluginTypeUnknown, "", "TestPluginCatalog_PluginMain_UserpassMultiplexed", []string{})
	TestAddTestPlugin(t, pluginCatalog, "single-userpass-1", consts.PluginTypeUnknown, "", "TestPluginCatalog_PluginMain_Userpass", []string{})
	TestAddTestPlugin(t, pluginCatalog, "single-userpass-2", consts.PluginTypeUnknown, "", "TestPluginCatalog_PluginMain_Userpass", []string{})

	getKey := func(pluginName string, pluginType consts.PluginType) externalPluginsKey {
		t.Helper()
		ctx := context.Background()
		plugin, err := pluginCatalog.Get(ctx, pluginName, pluginType, "")
		if err != nil {
			t.Fatal(err)
		}
		if plugin == nil {
			t.Fatal("did not find " + pluginName)
		}
		key, err := makeExternalPluginsKey(plugin)
		if err != nil {
			t.Fatal(err)
		}
		return key
	}

	var pluginClients []*pluginClient
	// run plugins
	// run "mux-postgres" twice which will start a single plugin for 2
	// distinct connections
	c := testRunTestPlugin(t, pluginCatalog, consts.PluginTypeDatabase, "mux-postgres")
	pluginClients = append(pluginClients, c)
	c = testRunTestPlugin(t, pluginCatalog, consts.PluginTypeDatabase, "mux-postgres")
	pluginClients = append(pluginClients, c)
	c = testRunTestPlugin(t, pluginCatalog, consts.PluginTypeDatabase, "single-postgres-1")
	pluginClients = append(pluginClients, c)
	c = testRunTestPlugin(t, pluginCatalog, consts.PluginTypeDatabase, "single-postgres-2")
	pluginClients = append(pluginClients, c)

	// run "mux-userpass" twice which will start a single plugin for 2
	// distinct connections
	c = testRunTestPlugin(t, pluginCatalog, consts.PluginTypeCredential, "mux-userpass")
	pluginClients = append(pluginClients, c)
	c = testRunTestPlugin(t, pluginCatalog, consts.PluginTypeCredential, "mux-userpass")
	pluginClients = append(pluginClients, c)
	c = testRunTestPlugin(t, pluginCatalog, consts.PluginTypeCredential, "single-userpass-1")
	pluginClients = append(pluginClients, c)
	c = testRunTestPlugin(t, pluginCatalog, consts.PluginTypeCredential, "single-userpass-2")
	pluginClients = append(pluginClients, c)

	externalPlugins := pluginCatalog.externalPlugins
	if len(externalPlugins) != 6 {
		t.Fatalf("expected externalPlugins map to be of len 6 but got %d", len(externalPlugins))
	}

	// check connections map
	expectConnectionLen(t, 2, externalPlugins[getKey("mux-postgres", consts.PluginTypeDatabase)].connections)
	expectConnectionLen(t, 1, externalPlugins[getKey("single-postgres-1", consts.PluginTypeDatabase)].connections)
	expectConnectionLen(t, 1, externalPlugins[getKey("single-postgres-2", consts.PluginTypeDatabase)].connections)
	expectConnectionLen(t, 2, externalPlugins[getKey("mux-userpass", consts.PluginTypeCredential)].connections)
	expectConnectionLen(t, 1, externalPlugins[getKey("single-userpass-1", consts.PluginTypeCredential)].connections)
	expectConnectionLen(t, 1, externalPlugins[getKey("single-userpass-2", consts.PluginTypeCredential)].connections)

	// check multiplexing support
	expectMultiplexingSupport(t, true, externalPlugins[getKey("mux-postgres", consts.PluginTypeDatabase)].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins[getKey("single-postgres-1", consts.PluginTypeDatabase)].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins[getKey("single-postgres-2", consts.PluginTypeDatabase)].multiplexingSupport)
	expectMultiplexingSupport(t, true, externalPlugins[getKey("mux-userpass", consts.PluginTypeCredential)].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins[getKey("single-userpass-1", consts.PluginTypeCredential)].multiplexingSupport)
	expectMultiplexingSupport(t, false, externalPlugins[getKey("single-userpass-2", consts.PluginTypeCredential)].multiplexingSupport)

	// cleanup all of the external plugin processes
	for _, client := range pluginClients {
		client.Close()
	}

	// check that externalPlugins map is cleaned up
	if len(externalPlugins) != 0 {
		t.Fatalf("expected external plugin map to be of len 0 but got %d", len(externalPlugins))
	}
}

func TestPluginCatalog_MakeExternalPluginsKey_Comparable(t *testing.T) {
	var plugins []pluginutil.PluginRunner
	hasher := sha256.New()
	hasher.Write([]byte("Some random input"))

	for i := 0; i < 2; i++ {
		plugins = append(plugins, pluginutil.PluginRunner{
			Name:    "Name",
			Type:    consts.PluginTypeDatabase,
			Version: "Version",
			Command: "Command",
			Args:    []string{"Some", "Args"},
			Env:     []string{"Env=foo", "bar=", "baz=foo"},
			Sha256:  hasher.Sum(nil),
			Builtin: true,
		})
	}

	var keys []externalPluginsKey
	for _, plugin := range plugins {
		key, err := makeExternalPluginsKey(&plugin)
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key)
	}

	if keys[0] != keys[1] {
		t.Fatal("expected equality")
	}
}

// TestPluginCatalog_ErrDirectoryNotConfigured ensures we correctly report an
// error when registering a binary plugin without a directory configured, and
// always allow registration of container plugins (rejecting on non-Linux happens
// in the logical system API handler).
func TestPluginCatalog_ErrDirectoryNotConfigured(t *testing.T) {
	catalog := testPluginCatalog(t)
	tempDir := catalog.directory
	catalog.directory = ""

	const pluginRuntime = "custom-runtime"
	const ociRuntime = "runc"
	err := catalog.runtimeCatalog.Set(context.Background(), &pluginruntimeutil.PluginRuntimeConfig{
		Name:       pluginRuntime,
		Type:       consts.PluginRuntimeTypeContainer,
		OCIRuntime: ociRuntime,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]func(t *testing.T){
		"set binary plugin": func(t *testing.T) {
			file, err := os.CreateTemp(tempDir, "temp")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			command := filepath.Base(file.Name())
			// Should error if directory not set.
			err = catalog.Set(context.Background(), pluginutil.SetPluginInput{
				Name:    "binary",
				Type:    consts.PluginTypeDatabase,
				Command: command,
			})
			dirSet := catalog.directory != ""
			if dirSet {
				if err != nil {
					t.Fatal(err)
				}
				p, err := catalog.Get(context.Background(), "binary", consts.PluginTypeDatabase, "")
				if err != nil {
					t.Fatal(err)
				}
				expectedCommand := filepath.Join(tempDir, command)
				if p.Command != expectedCommand {
					t.Fatalf("Expected %s, got %s", expectedCommand, p.Command)
				}
			}
			if !dirSet && err == nil {
				t.Fatal("expected error without directory set")
			}
			// Make sure we can still get builtins too
			_, err = catalog.Get(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase, "")
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
		},
		"set container plugin": func(t *testing.T) {
			if runtime.GOOS != "linux" {
				t.Skip("Containerized plugins only supported on Linux")
			}

			// Should never error.
			plugin := pluginhelpers.CompilePlugin(t, consts.PluginTypeDatabase, "", tempDir)
			plugin.Image, plugin.ImageSha256 = pluginhelpers.BuildPluginContainerImage(t, plugin, tempDir)

			err := catalog.Set(context.Background(), pluginutil.SetPluginInput{
				Name:     "container",
				Type:     consts.PluginTypeDatabase,
				OCIImage: plugin.Image,
				Runtime:  pluginRuntime,
			})
			if err != nil {
				t.Fatal(err)
			}
			// Check we can get it back ok.
			p, err := catalog.Get(context.Background(), "container", consts.PluginTypeDatabase, "")
			if err != nil {
				t.Fatal(err)
			}
			if p.OCIImage != plugin.Image {
				t.Fatalf("Expected %s, got %s", plugin.Image, p.OCIImage)
			}
			// Make sure we can still get builtins too
			_, err = catalog.Get(context.Background(), "mysql-database-plugin", consts.PluginTypeDatabase, "")
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
		},
	}

	t.Run("directory not set", func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test)
		}
	})

	catalog.directory = tempDir

	t.Run("directory set", func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test)
		}
	})
}

// TestRuntimeConfigPopulatedIfSpecified ensures plugins read from the catalog
// are returned with their container runtime config populated if it was
// specified.
func TestRuntimeConfigPopulatedIfSpecified(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Containerized plugins only supported on Linux")
	}

	pluginCatalog := testPluginCatalog(t)

	plugin := pluginhelpers.CompilePlugin(t, consts.PluginTypeDatabase, "", pluginCatalog.directory)
	plugin.Image, plugin.ImageSha256 = pluginhelpers.BuildPluginContainerImage(t, plugin, pluginCatalog.directory)

	const runtime = "custom-runtime"
	err := pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:     "container",
		Type:     consts.PluginTypeDatabase,
		OCIImage: plugin.Image,
		Runtime:  runtime,
	})
	if err == nil {
		t.Fatal("specified runtime doesn't exist yet, should have failed")
	}

	const ociRuntime = "runc"
	err = pluginCatalog.runtimeCatalog.Set(context.Background(), &pluginruntimeutil.PluginRuntimeConfig{
		Name:       runtime,
		Type:       consts.PluginRuntimeTypeContainer,
		OCIRuntime: ociRuntime,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Now setting the plugin with a runtime should succeed.
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:     "container",
		Type:     consts.PluginTypeDatabase,
		OCIImage: plugin.Image,
		Runtime:  runtime,
	})
	if err != nil {
		t.Fatal(err)
	}

	p, err := pluginCatalog.Get(context.Background(), "container", consts.PluginTypeDatabase, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Runtime != runtime {
		t.Errorf("expected %s, got %s", runtime, p.Runtime)
	}
	if p.RuntimeConfig == nil {
		t.Fatal()
	}
	if p.RuntimeConfig.OCIRuntime != ociRuntime {
		t.Errorf("expected %s, got %s", ociRuntime, p.RuntimeConfig.OCIRuntime)
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

func TestPluginCatalog_PluginMain_PostgresMultiplexed(_ *testing.T) {
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

func TestSortVersionedPlugins(t *testing.T) {
	versionedPlugin := func(typ consts.PluginType, name string, pluginVersion string, builtin bool) pluginutil.VersionedPlugin {
		return pluginutil.VersionedPlugin{
			Type:    typ.String(),
			Name:    name,
			Version: pluginVersion,
			SHA256:  "",
			Builtin: builtin,
			SemanticVersion: func() *version.Version {
				if pluginVersion != "" {
					return version.Must(version.NewVersion(pluginVersion))
				}

				return version.Must(version.NewVersion("0.0.0"))
			}(),
		}
	}

	differingTypes := []pluginutil.VersionedPlugin{
		versionedPlugin(consts.PluginTypeSecrets, "c", "1.0.0", false),
		versionedPlugin(consts.PluginTypeDatabase, "c", "1.0.0", false),
		versionedPlugin(consts.PluginTypeCredential, "c", "1.0.0", false),
	}
	differingNames := []pluginutil.VersionedPlugin{
		versionedPlugin(consts.PluginTypeCredential, "c", "1.0.0", false),
		versionedPlugin(consts.PluginTypeCredential, "b", "1.0.0", false),
		versionedPlugin(consts.PluginTypeCredential, "a", "1.0.0", false),
	}
	differingVersions := []pluginutil.VersionedPlugin{
		versionedPlugin(consts.PluginTypeCredential, "c", "10.0.0", false),
		versionedPlugin(consts.PluginTypeCredential, "c", "2.0.1", false),
		versionedPlugin(consts.PluginTypeCredential, "c", "2.1.0", false),
	}
	versionedUnversionedAndBuiltin := []pluginutil.VersionedPlugin{
		versionedPlugin(consts.PluginTypeCredential, "c", "1.0.0", false),
		versionedPlugin(consts.PluginTypeCredential, "c", "", false),
		versionedPlugin(consts.PluginTypeCredential, "c", "1.0.0", true),
	}

	for name, tc := range map[string][]pluginutil.VersionedPlugin{
		"ascending types":    differingTypes,
		"ascending names":    differingNames,
		"ascending versions": differingVersions,
		// Include differing versions twice so we can test out equality too.
		"differing types, names and versions": append(differingTypes,
			append(differingNames,
				append(differingVersions, differingVersions...)...)...),
		"mix of unversioned, versioned, and builtin": versionedUnversionedAndBuiltin,
	} {
		t.Run(name, func(t *testing.T) {
			sortVersionedPlugins(tc)
			for i := 1; i < len(tc); i++ {
				previous := tc[i-1]
				current := tc[i]
				if current.Type > previous.Type {
					continue
				}
				if current.Name > previous.Name {
					continue
				}
				if current.SemanticVersion.GreaterThan(previous.SemanticVersion) {
					continue
				}
				if current.Type == previous.Type && current.Name == previous.Name && current.SemanticVersion.Equal(previous.SemanticVersion) {
					continue
				}

				t.Fatalf("versioned plugins at index %d and %d were not properly sorted: %+v, %+v", i-1, i, previous, current)
			}
		})
	}
}

func TestExternalPlugin_getBackendTypeVersion(t *testing.T) {
	for name, tc := range map[string]struct {
		pluginType        consts.PluginType
		setRunningVersion string
	}{
		"external credential plugin": {
			pluginType:        consts.PluginTypeCredential,
			setRunningVersion: "v1.2.3",
		},
		"external secrets plugin": {
			pluginType:        consts.PluginTypeSecrets,
			setRunningVersion: "v1.2.3",
		},
		"external database plugin": {
			pluginType:        consts.PluginTypeDatabase,
			setRunningVersion: "v1.2.3",
		},
	} {
		t.Run(name, func(t *testing.T) {
			pluginCatalog := testPluginCatalog(t)
			plugin := pluginhelpers.CompilePlugin(t, tc.pluginType, tc.setRunningVersion, pluginCatalog.directory)

			shaBytes, _ := hex.DecodeString(plugin.Sha256)
			commandFull := filepath.Join(pluginCatalog.directory, plugin.FileName)
			entry := &pluginutil.PluginRunner{
				Name:    plugin.Name,
				Command: commandFull,
				Args:    nil,
				Sha256:  shaBytes,
				Builtin: false,
			}

			var version logical.PluginVersion
			var err error
			if tc.pluginType == consts.PluginTypeDatabase {
				version, err = pluginCatalog.getDatabaseRunningVersion(context.Background(), entry)
			} else {
				version, err = pluginCatalog.getBackendRunningVersion(context.Background(), entry)
			}
			if err != nil {
				t.Fatal(err)
			}
			if version.Version != tc.setRunningVersion {
				t.Errorf("Expected to get version %v but got %v", tc.setRunningVersion, version.Version)
			}
		})
	}
}

func TestExternalPluginInContainer_GetBackendTypeVersion(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Containerized plugins only supported on Linux")
	}

	pluginCatalog := testPluginCatalog(t)

	type testCase struct {
		plugin      pluginhelpers.TestPlugin
		expectedErr error
	}
	var testCases []testCase

	for _, pluginType := range []consts.PluginType{
		consts.PluginTypeCredential,
		consts.PluginTypeSecrets,
		consts.PluginTypeDatabase,
	} {
		plugin := pluginhelpers.CompilePlugin(t, pluginType, "v1.2.3", pluginCatalog.directory)
		plugin.Image, plugin.ImageSha256 = pluginhelpers.BuildPluginContainerImage(t, plugin, pluginCatalog.directory)

		testCases = append(testCases, testCase{
			plugin:      plugin,
			expectedErr: nil,
		})

		plugin.Image += "-will-not-be-found"
		testCases = append(testCases, testCase{
			plugin:      plugin,
			expectedErr: ErrPluginUnableToRun,
		})
	}

	for _, tc := range testCases {
		t.Run(tc.plugin.Typ.String(), func(t *testing.T) {
			expectedErrTestName := "nil err"
			if tc.expectedErr != nil {
				expectedErrTestName = tc.expectedErr.Error()
			}

			t.Run(expectedErrTestName, func(t *testing.T) {
				for _, ociRuntime := range []string{"runc", "runsc"} {
					t.Run(ociRuntime, func(t *testing.T) {
						if _, err := exec.LookPath(ociRuntime); err != nil {
							t.Skipf("Skipping test as %s not found on path", ociRuntime)
						}

						shaBytes, _ := hex.DecodeString(tc.plugin.ImageSha256)
						entry := &pluginutil.PluginRunner{
							Name:     tc.plugin.Name,
							OCIImage: tc.plugin.Image,
							Args:     nil,
							Sha256:   shaBytes,
							Builtin:  false,
							Runtime:  ociRuntime,
							RuntimeConfig: &pluginruntimeutil.PluginRuntimeConfig{
								OCIRuntime: ociRuntime,
							},
						}

						var version logical.PluginVersion
						var err error
						if tc.plugin.Typ == consts.PluginTypeDatabase {
							version, err = pluginCatalog.getDatabaseRunningVersion(context.Background(), entry)
						} else {
							version, err = pluginCatalog.getBackendRunningVersion(context.Background(), entry)
						}

						if tc.expectedErr == nil {
							if err != nil {
								t.Fatalf("Expected successful get backend type version but got: %v", err)
							}
							if version.Version != tc.plugin.Version {
								t.Errorf("Expected to get version %v but got %v", tc.plugin.Version, version.Version)
							}

						} else if !errors.Is(err, tc.expectedErr) {
							t.Errorf("Expected to get err %s but got %s", tc.expectedErr, err)
						}
					})
				}
			})
		})
	}
}

// TestPluginCatalog_CannotDeletePinnedVersion ensures we cannot delete a
// plugin which is referred to in an active pinned version.
func TestPluginCatalog_CannotDeletePinnedVersion(t *testing.T) {
	pluginCatalog := testPluginCatalog(t)

	// Register a fake plugin in the catalog.
	file, err := os.CreateTemp(pluginCatalog.directory, "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    "my-plugin",
		Type:    consts.PluginTypeSecrets,
		Version: "1.0.0",
		Command: filepath.Base(file.Name()),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Pin a version and check we can't delete it.
	err = pluginCatalog.SetPinnedVersion(context.Background(), &pluginutil.PinnedVersion{
		Name:    "my-plugin",
		Type:    consts.PluginTypeSecrets,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = pluginCatalog.Delete(context.Background(), "my-plugin", consts.PluginTypeSecrets, "1.0.0")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrPinnedVersion) {
		t.Fatal(err)
	}

	// Now delete the pinned version and we should be able to delete the plugin.
	err = pluginCatalog.DeletePinnedVersion(context.Background(), consts.PluginTypeSecrets, "my-plugin")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	err = pluginCatalog.Delete(context.Background(), "my-plugin", consts.PluginTypeSecrets, "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
}

// testRunTestPlugin runs the testFunc which has already been registered to the
// plugin catalog and returns a pluginClient. This can be called after calling
// TestAddTestPlugin.
func testRunTestPlugin(t *testing.T, pluginCatalog *PluginCatalog, pluginType consts.PluginType, pluginName string) *pluginClient {
	t.Helper()
	config := testPluginClientConfig(pluginCatalog, pluginType, pluginName)
	client, err := pluginCatalog.NewPluginClient(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func testPluginClientConfig(pluginCatalog *PluginCatalog, pluginType consts.PluginType, pluginName string) pluginutil.PluginClientConfig {
	config := pluginutil.PluginClientConfig{
		Name:           pluginName,
		PluginType:     pluginType,
		Logger:         log.NewNullLogger(),
		AutoMTLS:       true,
		IsMetadataMode: false,
		Wrapper: pluginCatalogStaticSystemView{
			pluginCatalog: pluginCatalog,
			StaticSystemView: logical.StaticSystemView{
				VersionString: "testVersion",
			},
		},
	}

	switch pluginType {
	case consts.PluginTypeCredential, consts.PluginTypeSecrets:
		config.PluginSets = backendplugin.PluginSet
		config.HandshakeConfig = backendplugin.HandshakeConfig
	case consts.PluginTypeDatabase:
		config.PluginSets = v5.PluginSets
		config.HandshakeConfig = v5.HandshakeConfig
	}

	return config
}

type pluginCatalogStaticSystemView struct {
	logical.StaticSystemView
	pluginCatalog *PluginCatalog
}

func (p pluginCatalogStaticSystemView) NewPluginClient(ctx context.Context, config pluginutil.PluginClientConfig) (pluginutil.PluginClient, error) {
	return p.pluginCatalog.NewPluginClient(ctx, config)
}
