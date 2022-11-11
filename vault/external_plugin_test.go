package vault

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin"
	"github.com/hashicorp/vault/sdk/plugin/mock"
)

const vaultTestingMockPluginEnv = "VAULT_TESTING_MOCK_PLUGIN"

var (
	pluginCacheLock sync.Mutex
	pluginCache     = map[string][]byte{}
)

type testPlugin struct {
	name     string
	typ      consts.PluginType
	version  string
	fileName string
	sha256   string
}

// version is used to override the plugin's self-reported version
func testCoreWithPlugins(t *testing.T, typ consts.PluginType, versions ...string) (*Core, []testPlugin) {
	t.Helper()
	pluginDir, cleanup := MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })

	var plugins []testPlugin
	for _, version := range versions {
		plugins = append(plugins, compilePlugin(t, typ, version, pluginDir))
	}
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	core := TestCoreWithSealAndUI(t, conf)
	core, _, _ = testCoreUnsealed(t, core)
	return core, plugins
}

func getPlugin(t *testing.T, typ consts.PluginType) (string, string, string, string) {
	t.Helper()
	var pluginName string
	var pluginType string
	var pluginMain string
	var pluginVersionLocation string

	switch typ {
	case consts.PluginTypeCredential:
		pluginType = "approle"
		pluginName = "vault-plugin-auth-" + pluginType
		pluginMain = filepath.Join("builtin", "credential", pluginType, "cmd", pluginType, "main.go")
		pluginVersionLocation = fmt.Sprintf("github.com/hashicorp/vault/builtin/credential/%s.ReportedVersion", pluginType)
	case consts.PluginTypeSecrets:
		pluginType = "consul"
		pluginName = "vault-plugin-secrets-" + pluginType
		pluginMain = filepath.Join("builtin", "logical", pluginType, "cmd", pluginType, "main.go")
		pluginVersionLocation = fmt.Sprintf("github.com/hashicorp/vault/builtin/logical/%s.ReportedVersion", pluginType)
	case consts.PluginTypeDatabase:
		pluginType = "postgresql"
		pluginName = "vault-plugin-database-" + pluginType
		pluginMain = filepath.Join("plugins", "database", pluginType, fmt.Sprintf("%s-database-plugin", pluginType), "main.go")
		pluginVersionLocation = fmt.Sprintf("github.com/hashicorp/vault/plugins/database/%s.ReportedVersion", pluginType)
	default:
		t.Fatal(typ.String())
	}
	return pluginName, pluginType, pluginMain, pluginVersionLocation
}

// to mount a plugin, we need a working binary plugin, so we compile one here.
// pluginVersion is used to override the plugin's self-reported version
func compilePlugin(t *testing.T, typ consts.PluginType, pluginVersion string, pluginDir string) testPlugin {
	t.Helper()

	pluginName, pluginType, pluginMain, pluginVersionLocation := getPlugin(t, typ)

	pluginCacheLock.Lock()
	defer pluginCacheLock.Unlock()

	var pluginBytes []byte

	dir := ""
	// detect if we are in the "vault/" or the root directory and compensate
	if _, err := os.Stat("builtin"); os.IsNotExist(err) {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		dir = filepath.Dir(wd)
	}

	pluginPath := path.Join(pluginDir, pluginName)
	if pluginVersion != "" {
		pluginPath += "-" + pluginVersion
	}

	key := fmt.Sprintf("%s %s %s", pluginName, pluginType, pluginVersion)
	// cache the compilation to only run once
	var ok bool
	pluginBytes, ok = pluginCache[key]
	if !ok {
		// we need to compile
		line := []string{"build"}
		if pluginVersion != "" {
			line = append(line, "-ldflags", fmt.Sprintf("-X %s=%s", pluginVersionLocation, pluginVersion))
		}
		line = append(line, "-o", pluginPath, pluginMain)
		cmd := exec.Command("go", line...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("error running go build %v output: %s", err, output))
		}
		pluginCache[key], err = os.ReadFile(pluginPath)
		if err != nil {
			t.Fatal(err)
		}
		pluginBytes = pluginCache[key]
	}

	// write the cached plugin if necessary
	var err error
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		err = os.WriteFile(pluginPath, pluginBytes, 0o755)
	}
	if err != nil {
		t.Fatal(err)
	}

	sha := sha256.New()
	_, err = sha.Write(pluginBytes)
	if err != nil {
		t.Fatal(err)
	}
	return testPlugin{
		name:     pluginName,
		typ:      typ,
		version:  pluginVersion,
		fileName: path.Base(pluginPath),
		sha256:   fmt.Sprintf("%x", sha.Sum(nil)),
	}
}

func TestCore_EnableExternalPlugin(t *testing.T) {
	for name, tc := range map[string]struct {
		pluginType    consts.PluginType
		routerPath    string
		expectedMatch string
	}{
		"enable external credential plugin": {
			pluginType:    consts.PluginTypeCredential,
			routerPath:    "auth/foo/bar",
			expectedMatch: "auth/foo/",
		},
		"enable external secrets plugin": {
			pluginType:    consts.PluginTypeSecrets,
			routerPath:    "foo/bar",
			expectedMatch: "foo/",
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, plugins := testCoreWithPlugins(t, tc.pluginType, "")
			registerPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType.String(), "1.0.0", plugins[0].sha256, plugins[0].fileName)

			mountPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType, "v1.0.0", "")

			match := c.router.MatchingMount(namespace.RootContext(nil), tc.routerPath)
			if match != tc.expectedMatch {
				t.Fatalf("missing mount, match: %q", match)
			}
		})
	}
}

func TestCore_EnableExternalPlugin_MultipleVersions(t *testing.T) {
	for name, tc := range map[string]struct {
		pluginType       consts.PluginType
		registerVersions []string
		mountVersion     string
		expectedVersion  string
		routerPath       string
		expectedMatch    string
	}{
		"enable external credential plugin, multiple versions available": {
			pluginType:       consts.PluginTypeCredential,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.0",
			expectedVersion:  "v1.0.0",
			routerPath:       "auth/foo/bar",
			expectedMatch:    "auth/foo/",
		},
		"enable external secrets plugin, multiple versions available": {
			pluginType:       consts.PluginTypeSecrets,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.0",
			expectedVersion:  "v1.0.0",
			routerPath:       "foo/bar",
			expectedMatch:    "foo/",
		},
		"enable external credential plugin, multiple versions available, select other version": {
			pluginType:       consts.PluginTypeCredential,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.1",
			expectedVersion:  "v1.0.1",
			routerPath:       "auth/foo/bar",
			expectedMatch:    "auth/foo/",
		},
		"enable external secrets plugin, multiple versions available, select other version": {
			pluginType:       consts.PluginTypeSecrets,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.1",
			expectedVersion:  "v1.0.1",
			routerPath:       "foo/bar",
			expectedMatch:    "foo/",
		},
		"enable external credential plugin, selects latest when version not specified": {
			pluginType:       consts.PluginTypeCredential,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "",
			expectedVersion:  "v1.0.1",
			routerPath:       "auth/foo/bar",
			expectedMatch:    "auth/foo/",
		},
		"enable external secrets plugin, selects latest when version not specified": {
			pluginType:       consts.PluginTypeSecrets,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "",
			expectedVersion:  "v1.0.1",
			routerPath:       "foo/bar",
			expectedMatch:    "foo/",
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, plugins := testCoreWithPlugins(t, tc.pluginType, "")
			for _, version := range tc.registerVersions {
				registerPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType.String(), version, plugins[0].sha256, plugins[0].fileName)
			}

			mountPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType, tc.mountVersion, "")

			match := c.router.MatchingMount(namespace.RootContext(nil), tc.routerPath)
			if match != tc.expectedMatch {
				t.Fatalf("missing mount, match: %q", match)
			}

			raw, _ := c.router.root.Get(match)
			if raw.(*routeEntry).mountEntry.Version != tc.expectedVersion {
				t.Errorf("Expected mount to be version %s but got %s", tc.expectedVersion, raw.(*routeEntry).mountEntry.Version)
			}

			if raw.(*routeEntry).mountEntry.RunningVersion != tc.expectedVersion {
				t.Errorf("Expected mount running version to be %s but got %s", tc.expectedVersion, raw.(*routeEntry).mountEntry.RunningVersion)
			}

			if raw.(*routeEntry).mountEntry.RunningSha256 == "" {
				t.Errorf("Expected RunningSha256 to be present: %+v", raw.(*routeEntry).mountEntry.RunningSha256)
			}
		})
	}
}

func TestCore_EnableExternalPlugin_PendingRemoval(t *testing.T) {
	pluginDir, cleanup := MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })

	// create an external plugin to shadow the builtin "app-id"
	plugin := compilePlugin(t, consts.PluginTypeCredential, "v1.2.3", pluginDir)
	err := os.Link(path.Join(pluginDir, plugin.fileName), path.Join(pluginDir, "app-id"))
	if err != nil {
		t.Fatal(err)
	}
	pluginName := "app-id"
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	c := TestCoreWithSealAndUI(t, conf)
	c, _, _ = testCoreUnsealed(t, c)

	pendingRemovalString := "pending removal"

	// Create a new auth method with builtin app-id
	resp, err := mountPluginWithResponse(t, c.systemBackend, pluginName, consts.PluginTypeCredential, "", "")
	if err == nil {
		t.Fatalf("expected error when mounting deprecated backend")
	}
	if resp == nil || resp.Data == nil || !strings.Contains(resp.Data["error"].(string), pendingRemovalString) {
		t.Fatalf("expected error response to contain %q but got %+v", pendingRemovalString, resp)
	}

	// Register a shadow plugin
	registerPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential.String(), "v1.2.3", plugin.sha256, plugin.fileName)
	mountPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential, "", "")
}

func TestCore_EnableExternalPlugin_ShadowBuiltin(t *testing.T) {
	pluginDir, cleanup := MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })

	// create an external plugin to shadow the builtin "approle"
	plugin := compilePlugin(t, consts.PluginTypeCredential, "v1.2.3", pluginDir)
	err := os.Link(path.Join(pluginDir, plugin.fileName), path.Join(pluginDir, "approle"))
	if err != nil {
		t.Fatal(err)
	}
	pluginName := "approle"
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	c := TestCoreWithSealAndUI(t, conf)
	c, _, _ = testCoreUnsealed(t, c)

	verifyAuthListDeprecationStatus := func(authName string, checkExists bool) error {
		req := logical.TestRequest(t, logical.ReadOperation, mountTable(consts.PluginTypeCredential))
		req.Data = map[string]interface{}{
			"type": authName,
		}
		resp, err := c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			return err
		}
		status := resp.Data["deprecation_status"]
		if checkExists && status == nil {
			return fmt.Errorf("expected deprecation status but found none")
		} else if !checkExists && status != nil {
			return fmt.Errorf("expected nil deprecation status but found %q", status)
		}
		return nil
	}

	// Create a new auth method with builtin approle
	mountPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential, "", "")

	// Read the auth table to verify deprecation status
	if err := verifyAuthListDeprecationStatus(pluginName, true); err != nil {
		t.Fatal(err)
	}

	// Register a shadow plugin
	registerPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential.String(), "v1.2.3", plugin.sha256, plugin.fileName)

	// Verify auth table hasn't changed
	if err := verifyAuthListDeprecationStatus(pluginName, true); err != nil {
		t.Fatal(err)
	}

	// Remount auth method using registered shadow plugin
	unmountPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential, "", "")
	mountPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential, "", "")

	// Verify auth table has changed
	if err := verifyAuthListDeprecationStatus(pluginName, false); err != nil {
		t.Fatal(err)
	}

	// Deregister shadow plugin
	deregisterPlugin(t, c.systemBackend, pluginName, consts.PluginTypeSecrets.String(), "v1.2.3", plugin.sha256, plugin.fileName)

	// Verify auth table hasn't changed
	if err := verifyAuthListDeprecationStatus(pluginName, false); err != nil {
		t.Fatal(err)
	}

	// Remount auth method
	unmountPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential, "", "")
	mountPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential, "", "")

	// Verify auth table has changed
	if err := verifyAuthListDeprecationStatus(pluginName, false); err != nil {
		t.Fatal(err)
	}
}

func TestCore_EnableExternalKv_MultipleVersions(t *testing.T) {
	pluginDir, cleanup := MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })

	// new kv plugin can be registered but not mounted
	plugin := compilePlugin(t, consts.PluginTypeSecrets, "v1.2.3", pluginDir)
	err := os.Link(path.Join(pluginDir, plugin.fileName), path.Join(pluginDir, "kv"))
	if err != nil {
		t.Fatal(err)
	}
	pluginName := "kv"
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	c := TestCoreWithSealAndUI(t, conf)
	c, _, _ = testCoreUnsealed(t, c)

	registerPlugin(t, c.systemBackend, pluginName, consts.PluginTypeSecrets.String(), "v1.2.3", plugin.sha256, plugin.fileName)
	req := logical.TestRequest(t, logical.ReadOperation, "plugins/catalog")
	resp, err := c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() != nil {
		t.Fatalf("%#v", resp)
	}
	found := false
	for _, plugin := range resp.Data["detailed"].([]pluginutil.VersionedPlugin) {
		if plugin.Name == pluginName && plugin.Version == "v1.2.3" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Expected to find v1.2.3 kv plugin but did not")
	}
	req = logical.TestRequest(t, logical.UpdateOperation, mountTable(consts.PluginTypeSecrets))
	req.Data = map[string]interface{}{
		"type": pluginName,
	}
	req.Data["config"] = map[string]interface{}{
		"plugin_version": "v1.2.3",
	}
	resp, err = c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() == nil {
		t.Fatal("Expected resp error but got successful response")
	}
}

func TestCore_EnableExternalNoop_MultipleVersions(t *testing.T) {
	pluginDir, cleanup := MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })

	// new noop plugin can be registered but not mounted
	plugin := compilePlugin(t, consts.PluginTypeCredential, "v1.2.3", pluginDir)
	err := os.Link(path.Join(pluginDir, plugin.fileName), path.Join(pluginDir, "noop"))
	if err != nil {
		t.Fatal(err)
	}
	pluginName := "noop"
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	c := TestCoreWithSealAndUI(t, conf)
	c, _, _ = testCoreUnsealed(t, c)

	registerPlugin(t, c.systemBackend, pluginName, consts.PluginTypeCredential.String(), "v1.2.3", plugin.sha256, plugin.fileName)
	req := logical.TestRequest(t, logical.ReadOperation, "plugins/catalog")
	resp, err := c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() != nil {
		t.Fatalf("%#v", resp)
	}
	found := false
	for _, plugin := range resp.Data["detailed"].([]pluginutil.VersionedPlugin) {
		if plugin.Name == "noop" && plugin.Version == "v1.2.3" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Expected to find v1.2.3 noop plugin but did not")
	}
	req = logical.TestRequest(t, logical.UpdateOperation, mountTable(consts.PluginTypeCredential))
	req.Data = map[string]interface{}{
		"type": pluginName,
	}
	req.Data["config"] = map[string]interface{}{
		"plugin_version": "v1.2.3",
	}
	resp, err = c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() == nil {
		t.Fatal("Expected resp error but got successful response")
	}
}

func TestCore_EnableExternalPlugin_NoVersionsOkay(t *testing.T) {
	for name, tc := range map[string]struct {
		pluginType    consts.PluginType
		routerPath    string
		expectedMatch string
	}{
		"enable external credential plugin with no version": {
			pluginType:    consts.PluginTypeCredential,
			routerPath:    "auth/foo/bar",
			expectedMatch: "auth/foo/",
		},
		"enable external secrets plugin with no version": {
			pluginType:    consts.PluginTypeSecrets,
			routerPath:    "foo/bar",
			expectedMatch: "foo/",
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, plugins := testCoreWithPlugins(t, tc.pluginType, "")
			// When an unversioned plugin is registered, mounting a plugin with no
			// version specified should mount the unversioned plugin even if there
			// are versioned plugins available.
			for _, version := range []string{"", "v1.0.0"} {
				registerPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType.String(), version, plugins[0].sha256, plugins[0].fileName)
			}

			mountPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType, "", "")

			match := c.router.MatchingMount(namespace.RootContext(nil), tc.routerPath)
			if match != tc.expectedMatch {
				t.Fatalf("missing mount, match: %q", match)
			}

			raw, _ := c.router.root.Get(match)
			if raw.(*routeEntry).mountEntry.Version != "" {
				t.Errorf("Expected mount to be empty version but got %s", raw.(*routeEntry).mountEntry.Version)
			}
		})
	}
}

func TestCore_EnableExternalCredentialPlugin_NoVersionOnRegister(t *testing.T) {
	for name, tc := range map[string]struct {
		pluginType    consts.PluginType
		routerPath    string
		expectedMatch string
	}{
		"enable external credential plugin with version, but no version was provided on registration": {
			pluginType:    consts.PluginTypeCredential,
			routerPath:    "auth/foo/bar",
			expectedMatch: "auth/foo/",
		},
		"enable external secrets plugin with version, but no version was provided on registration": {
			pluginType:    consts.PluginTypeSecrets,
			routerPath:    "foo/bar",
			expectedMatch: "foo/",
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, plugins := testCoreWithPlugins(t, tc.pluginType, "")
			registerPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType.String(), "", plugins[0].sha256, plugins[0].fileName)

			req := logical.TestRequest(t, logical.UpdateOperation, mountTable(tc.pluginType))
			req.Data = map[string]interface{}{
				"type": plugins[0].name,
				"config": map[string]interface{}{
					"plugin_version": "v1.0.0",
				},
			}
			resp, _ := c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
			if resp == nil || !resp.IsError() || !strings.Contains(resp.Error().Error(), ErrPluginNotFound.Error()) {
				t.Fatalf("Expected to get plugin not found but got: %v", resp.Error())
			}
		})
	}
}

func TestCore_EnableExternalCredentialPlugin_InvalidName(t *testing.T) {
	for name, tc := range map[string]struct {
		pluginType consts.PluginType
	}{
		"enable external credential plugin with the wrong name": {
			pluginType: consts.PluginTypeCredential,
		},
		"enable external secrets plugin with the wrong name": {
			pluginType: consts.PluginTypeSecrets,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, plugins := testCoreWithPlugins(t, tc.pluginType, "")
			d := &framework.FieldData{
				Raw: map[string]interface{}{
					"name":    plugins[0].name,
					"sha256":  plugins[0].sha256,
					"version": "v1.0.0",
					"command": plugins[0].name + "xyz",
				},
				Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
			}
			_, err := c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
			if err == nil || !strings.Contains(err.Error(), "no such file or directory") {
				t.Fatalf("should have gotten a no such file or directory error inserting the plugin: %v", err)
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
			c, plugins := testCoreWithPlugins(t, tc.pluginType, tc.setRunningVersion)
			registerPlugin(t, c.systemBackend, plugins[0].name, tc.pluginType.String(), tc.setRunningVersion, plugins[0].sha256, plugins[0].fileName)

			shaBytes, _ := hex.DecodeString(plugins[0].sha256)
			commandFull := filepath.Join(c.pluginCatalog.directory, plugins[0].fileName)
			entry := &pluginutil.PluginRunner{
				Name:    plugins[0].name,
				Command: commandFull,
				Args:    nil,
				Sha256:  shaBytes,
				Builtin: false,
			}

			var version logical.PluginVersion
			var err error
			if tc.pluginType == consts.PluginTypeDatabase {
				version, err = c.pluginCatalog.getDatabaseRunningVersion(context.Background(), entry)
			} else {
				version, err = c.pluginCatalog.getBackendRunningVersion(context.Background(), entry)
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

func TestExternalPlugin_CheckFilePermissions(t *testing.T) {
	// Turn on the check.
	if err := os.Setenv(consts.VaultEnableFilePermissionsCheckEnv, "true"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv(consts.VaultEnableFilePermissionsCheckEnv); err != nil {
			t.Fatal(err)
		}
	}()

	for name, tc := range map[string]struct {
		pluginNameFmt string
		pluginType    consts.PluginType
		pluginVersion string
	}{
		"plugin name and file name match": {
			pluginNameFmt: "%s",
			pluginType:    consts.PluginTypeCredential,
		},
		"plugin name and file name mismatch": {
			pluginNameFmt: "%s-foo",
			pluginType:    consts.PluginTypeSecrets,
		},
		"plugin name has slash": {
			pluginNameFmt: "%s/foo",
			pluginType:    consts.PluginTypeCredential,
		},
		"plugin with version": {
			pluginNameFmt: "%s/foo",
			pluginType:    consts.PluginTypeCredential,
			pluginVersion: "v1.2.3",
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, plugins := testCoreWithPlugins(t, tc.pluginType, tc.pluginVersion)
			registeredPluginName := fmt.Sprintf(tc.pluginNameFmt, plugins[0].name)

			// Permissions will be checked once during registration.
			req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("plugins/catalog/%s/%s", tc.pluginType.String(), registeredPluginName))
			req.Data = map[string]interface{}{
				"command": plugins[0].fileName,
				"sha256":  plugins[0].sha256,
				"version": tc.pluginVersion,
			}
			resp, err := c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Error() != nil {
				t.Fatal(resp.Error())
			}

			// Now attempt to mount the plugin, which should trigger checking the permissions again.
			req = logical.TestRequest(t, logical.UpdateOperation, mountTable(tc.pluginType))
			req.Data = map[string]interface{}{
				"type": registeredPluginName,
			}
			if tc.pluginVersion != "" {
				req.Data["config"] = map[string]interface{}{
					"plugin_version": tc.pluginVersion,
				}
			}
			resp, err = c.systemBackend.HandleRequest(namespace.RootContext(nil), req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Error() != nil {
				t.Fatal(resp.Error())
			}
		})
	}
}

func TestExternalPlugin_DifferentVersionsAndArgs_AreNotMultiplexed(t *testing.T) {
	env := []string{fmt.Sprintf("%s=yes", vaultTestingMockPluginEnv)}
	core, _, _ := TestCoreUnsealed(t)

	for i, tc := range []struct {
		version  string
		testName string
	}{
		{"v1.2.3", "TestBackend_PluginMain_Multiplexed_Logical_v123"},
		{"v1.2.4", "TestBackend_PluginMain_Multiplexed_Logical_v124"},
	} {
		// Register and mount plugins.
		TestAddTestPlugin(t, core, "mux-secret", consts.PluginTypeSecrets, tc.version, tc.testName, env, "")
		mountPlugin(t, core.systemBackend, "mux-secret", consts.PluginTypeSecrets, tc.version, fmt.Sprintf("foo%d", i))
	}

	if len(core.pluginCatalog.externalPlugins) != 2 {
		t.Fatalf("expected 2 external plugins, but got %d", len(core.pluginCatalog.externalPlugins))
	}
}

func TestExternalPlugin_DifferentTypes_AreNotMultiplexed(t *testing.T) {
	const version = "v1.2.3"
	env := []string{fmt.Sprintf("%s=yes", vaultTestingMockPluginEnv)}
	core, _, _ := TestCoreUnsealed(t)

	// Register and mount plugins.
	TestAddTestPlugin(t, core, "mux-aws", consts.PluginTypeSecrets, version, "TestBackend_PluginMain_Multiplexed_Logical_v123", env, "")
	TestAddTestPlugin(t, core, "mux-aws", consts.PluginTypeCredential, version, "TestBackend_PluginMain_Multiplexed_Credential_v123", env, "")

	mountPlugin(t, core.systemBackend, "mux-aws", consts.PluginTypeSecrets, version, "")
	mountPlugin(t, core.systemBackend, "mux-aws", consts.PluginTypeCredential, version, "")

	if len(core.pluginCatalog.externalPlugins) != 2 {
		t.Fatalf("expected 2 external plugins, but got %d", len(core.pluginCatalog.externalPlugins))
	}
}

func TestExternalPlugin_DifferentEnv_AreNotMultiplexed(t *testing.T) {
	const version = "v1.2.3"
	baseEnv := []string{
		fmt.Sprintf("%s=yes", vaultTestingMockPluginEnv),
	}
	alteredEnv := []string{
		fmt.Sprintf("%s=yes", vaultTestingMockPluginEnv),
		"FOO=BAR",
	}

	core, _, _ := TestCoreUnsealed(t)

	// Register and mount plugins.
	for i, env := range [][]string{baseEnv, alteredEnv} {
		TestAddTestPlugin(t, core, "mux-secret", consts.PluginTypeSecrets, version, "TestBackend_PluginMain_Multiplexed_Logical_v123", env, "")
		mountPlugin(t, core.systemBackend, "mux-secret", consts.PluginTypeSecrets, version, fmt.Sprintf("foo%d", i))
	}

	if len(core.pluginCatalog.externalPlugins) != 2 {
		t.Fatalf("expected 2 external plugins, but got %d", len(core.pluginCatalog.externalPlugins))
	}
}

// Used to run a mock multiplexed secrets plugin
func TestBackend_PluginMain_Multiplexed_Logical_v123(t *testing.T) {
	if os.Getenv(vaultTestingMockPluginEnv) == "" {
		return
	}

	os.Setenv(mock.MockPluginVersionEnv, "v1.2.3")

	err := plugin.ServeMultiplex(&plugin.ServeOpts{
		BackendFactoryFunc: mock.FactoryType(logical.TypeLogical),
	})
	if err != nil {
		t.Fatal(err)
	}
}

// Used to run a mock multiplexed secrets plugin
func TestBackend_PluginMain_Multiplexed_Logical_v124(t *testing.T) {
	if os.Getenv(vaultTestingMockPluginEnv) == "" {
		return
	}

	os.Setenv(mock.MockPluginVersionEnv, "v1.2.4")

	err := plugin.ServeMultiplex(&plugin.ServeOpts{
		BackendFactoryFunc: mock.FactoryType(logical.TypeLogical),
	})
	if err != nil {
		t.Fatal(err)
	}
}

// Used to run a mock multiplexed auth plugin
func TestBackend_PluginMain_Multiplexed_Credential_v123(t *testing.T) {
	if os.Getenv(vaultTestingMockPluginEnv) == "" {
		return
	}

	os.Setenv(mock.MockPluginVersionEnv, "v1.2.3")

	err := plugin.ServeMultiplex(&plugin.ServeOpts{
		BackendFactoryFunc: mock.FactoryType(logical.TypeCredential),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func registerPlugin(t *testing.T, sys *SystemBackend, pluginName, pluginType, version, sha, command string) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("plugins/catalog/%s/%s", pluginType, pluginName))
	req.Data = map[string]interface{}{
		"command": command,
		"sha256":  sha,
		"version": version,
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}

func mountPluginWithResponse(t *testing.T, sys *SystemBackend, pluginName string, pluginType consts.PluginType, version, path string) (*logical.Response, error) {
	t.Helper()
	var mountPath string
	if path == "" {
		mountPath = mountTable(pluginType)
	} else {
		mountPath = mountTableWithPath(consts.PluginTypeSecrets, path)
	}
	req := logical.TestRequest(t, logical.UpdateOperation, mountPath)
	req.Data = map[string]interface{}{
		"type": pluginName,
	}
	if version != "" {
		req.Data["config"] = map[string]interface{}{
			"plugin_version": version,
		}
	}
	return sys.HandleRequest(namespace.RootContext(nil), req)
}

func mountPlugin(t *testing.T, sys *SystemBackend, pluginName string, pluginType consts.PluginType, version, path string) {
	t.Helper()
	resp, err := mountPluginWithResponse(t, sys, pluginName, pluginType, version, path)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}

func unmountPlugin(t *testing.T, sys *SystemBackend, pluginName string, pluginType consts.PluginType, version, path string) {
	t.Helper()
	var mountPath string
	if path == "" {
		mountPath = mountTable(pluginType)
	} else {
		mountPath = mountTableWithPath(consts.PluginTypeSecrets, path)
	}
	req := logical.TestRequest(t, logical.DeleteOperation, mountPath)
	req.Data = map[string]interface{}{
		"type": pluginName,
	}
	if version != "" {
		req.Data["config"] = map[string]interface{}{
			"plugin_version": version,
		}
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}

func deregisterPlugin(t *testing.T, sys *SystemBackend, pluginName, pluginType, version, sha, command string) {
	t.Helper()
	req := logical.TestRequest(t, logical.DeleteOperation, fmt.Sprintf("plugins/catalog/%s/%s", pluginType, pluginName))
	req.Data = map[string]interface{}{
		"command": command,
		"sha256":  sha,
		"version": version,
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}

func mountTable(pluginType consts.PluginType) string {
	return mountTableWithPath(pluginType, "foo")
}

func mountTableWithPath(pluginType consts.PluginType, path string) string {
	switch pluginType {
	case consts.PluginTypeCredential:
		return "auth/" + path
	case consts.PluginTypeSecrets:
		return "mounts/" + path
	default:
		panic("test does not support mounting plugin type yet: " + pluginType.String())
	}
}
