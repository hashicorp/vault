package vault

import (
	"context"
	"crypto/sha256"
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
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	compileAuthOnce   sync.Once
	compileSecretOnce sync.Once
	authPluginBytes   []byte
	secretPluginBytes []byte
)

func testCoreWithPlugin(t *testing.T, typ consts.PluginType) (*Core, string, string) {
	t.Helper()
	pluginName, pluginSHA256, pluginDir := compilePlugin(t, typ)
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	core := TestCoreWithSealAndUI(t, conf)
	core, _, _ = testCoreUnsealed(t, core)
	return core, pluginName, pluginSHA256
}

// to mount a plugin, we need a working binary plugin, so we compile one here.
func compilePlugin(t *testing.T, typ consts.PluginType) (name string, shasum string, pluginDir string) {
	t.Helper()

	var pluginType, pluginName, builtinDirectory string
	var once *sync.Once
	var pluginBytes *[]byte
	switch typ {
	case consts.PluginTypeCredential:
		pluginType = "approle"
		pluginName = "vault-plugin-auth-" + pluginType
		builtinDirectory = "credential"
		once = &compileAuthOnce
		pluginBytes = &authPluginBytes
	case consts.PluginTypeSecrets:
		pluginType = "consul"
		pluginName = "vault-plugin-secrets-" + pluginType
		builtinDirectory = "logical"
		once = &compileSecretOnce
		pluginBytes = &secretPluginBytes
	default:
		t.Fatal(typ.String())
	}

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

	// cache the compilation to only run once
	once.Do(func() {
		cmd := exec.Command("go", "build", "-o", pluginPath, fmt.Sprintf("builtin/%s/%s/cmd/%s/main.go", builtinDirectory, pluginType, pluginType))
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("error running go build %v output: %s", err, output))
		}
		*pluginBytes, err = os.ReadFile(pluginPath)
		if err != nil {
			t.Fatal(err)
		}
	})

	// write the cached plugin if necessary
	var err error
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		err = os.WriteFile(pluginPath, *pluginBytes, 0o777)
	}
	if err != nil {
		t.Fatal(err)
	}

	sha := sha256.New()
	_, err = sha.Write(*pluginBytes)
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType)
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType)
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

			// we don't override the running version of non-builtins, and they don't have the version set explicitly (yet)
			if raw.(*routeEntry).mountEntry.RunningVersion != "" {
				t.Errorf("Expected mount to have no running version but got %s", raw.(*routeEntry).mountEntry.RunningVersion)
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType)
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType)
			registerPlugin(t, c.systemBackend, pluginName, tc.pluginType.String(), "", pluginSHA256)

			req := logical.TestRequest(t, logical.UpdateOperation, mountPath(tc.pluginType))
			req.Data = map[string]interface{}{
				"type":    pluginName,
				"version": "v1.0.0",
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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType)
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
	req := logical.TestRequest(t, logical.UpdateOperation, mountPath(pluginType))
	req.Data = map[string]interface{}{
		"type": pluginName,
	}
	if version != "" {
		req.Data["version"] = version
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() != nil {
		t.Fatalf("%#v", resp)
	}
}

func mountPath(pluginType consts.PluginType) string {
	switch pluginType {
	case consts.PluginTypeCredential:
		return "auth/foo"
	case consts.PluginTypeSecrets:
		return "mounts/foo"
	default:
		panic("test does not support plugin type yet")
	}
}
