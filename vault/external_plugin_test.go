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
)

var (
	pluginCacheLock sync.Mutex
	pluginCache     = map[string][]byte{}
)

// version is used to override the plugin's self-reported version
func testCoreWithPlugin(t *testing.T, typ consts.PluginType, version string) (*Core, string, string) {
	t.Helper()
	pluginName, pluginSHA256, pluginDir := compilePlugin(t, typ, version)
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	core := TestCoreWithSealAndUI(t, conf)
	core, _, _ = testCoreUnsealed(t, core)
	return core, pluginName, pluginSHA256
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
func compilePlugin(t *testing.T, typ consts.PluginType, pluginVersion string) (pluginName string, shasum string, pluginDir string) {
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

	pluginDir, cleanup := MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })

	pluginPath := path.Join(pluginDir, pluginName)

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
		err = os.WriteFile(pluginPath, pluginBytes, 0o777)
	}
	if err != nil {
		t.Fatal(err)
	}

	sha := sha256.New()
	_, err = sha.Write(pluginBytes)
	if err != nil {
		t.Fatal(err)
	}
	return pluginName, fmt.Sprintf("%x", sha.Sum(nil)), pluginDir
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType, "")
			registerPlugin(t, c.systemBackend, pluginName, tc.pluginType.String(), "1.0.0", pluginSHA256)

			mountPlugin(t, c.systemBackend, pluginName, tc.pluginType, "v1.0.0")

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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType, "")
			for _, version := range tc.registerVersions {
				registerPlugin(t, c.systemBackend, pluginName, tc.pluginType.String(), version, pluginSHA256)
			}

			mountPlugin(t, c.systemBackend, pluginName, tc.pluginType, tc.mountVersion)

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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType, "")
			// When an unversioned plugin is registered, mounting a plugin with no
			// version specified should mount the unversioned plugin even if there
			// are versioned plugins available.
			for _, version := range []string{"", "v1.0.0"} {
				registerPlugin(t, c.systemBackend, pluginName, tc.pluginType.String(), version, pluginSHA256)
			}

			mountPlugin(t, c.systemBackend, pluginName, tc.pluginType, "")

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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType, "")
			registerPlugin(t, c.systemBackend, pluginName, tc.pluginType.String(), "", pluginSHA256)

			req := logical.TestRequest(t, logical.UpdateOperation, mountTable(tc.pluginType))
			req.Data = map[string]interface{}{
				"type": pluginName,
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType, "")
			d := &framework.FieldData{
				Raw: map[string]interface{}{
					"name":    pluginName,
					"sha256":  pluginSHA256,
					"version": "v1.0.0",
					"command": pluginName + "xyz",
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType, tc.setRunningVersion)
			registerPlugin(t, c.systemBackend, pluginName, tc.pluginType.String(), tc.setRunningVersion, pluginSHA256)

			shaBytes, _ := hex.DecodeString(pluginSHA256)
			commandFull := filepath.Join(c.pluginCatalog.directory, pluginName)
			entry := &pluginutil.PluginRunner{
				Name:    pluginName,
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

func registerPlugin(t *testing.T, sys *SystemBackend, pluginName, pluginType, version, sha string) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("plugins/catalog/%s/%s", pluginType, pluginName))
	req.Data = map[string]interface{}{
		"name":    pluginName,
		"command": pluginName,
		"sha256":  sha,
		"version": version,
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() != nil {
		t.Fatalf("%#v", resp)
	}
}

func mountPlugin(t *testing.T, sys *SystemBackend, pluginName string, pluginType consts.PluginType, version string) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, mountTable(pluginType))
	req.Data = map[string]interface{}{
		"type": pluginName,
	}
	if version != "" {
		req.Data["config"] = map[string]interface{}{
			"plugin_version": version,
		}
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() != nil {
		t.Fatalf("%#v", resp)
	}
}

func mountTable(pluginType consts.PluginType) string {
	switch pluginType {
	case consts.PluginTypeCredential:
		return "auth/foo"
	case consts.PluginTypeSecrets:
		return "mounts/foo"
	default:
		panic("test does not support plugin type yet")
	}
}
