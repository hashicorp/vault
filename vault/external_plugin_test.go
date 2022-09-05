package vault

import (
	"context"
	"crypto/sha256"
	"errors"
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
)

var (
	compileAuthOnce   sync.Once
	compileSecretOnce sync.Once
	pluginBytes       []byte
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
	var once sync.Once
	switch typ {
	case consts.PluginTypeCredential:
		pluginType = "approle"
		pluginName = "vault-plugin-auth-" + pluginType
		builtinDirectory = "credential"
		once = compileAuthOnce
	case consts.PluginTypeSecrets:
		pluginType = "consul"
		pluginName = "vault-plugin-secrets-" + pluginType
		builtinDirectory = "logical"
		once = compileSecretOnce
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
		pluginBytes, err = os.ReadFile(pluginPath)
		if err != nil {
			t.Fatal(err)
		}
	})

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
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType)
			d := &framework.FieldData{
				Raw: map[string]interface{}{
					"name":    pluginName,
					"sha256":  pluginSHA256,
					"version": "v1.0.0",
					"command": pluginName,
				},
				Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
			}
			resp, err := c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Error() != nil {
				t.Fatalf("%#v", resp)
			}

			me := &MountEntry{
				Table:   mountTable(tc.pluginType),
				Path:    "foo",
				Type:    pluginName,
				Version: "v1.0.0",
			}
			enable := enableFunc(c, tc.pluginType)
			err = enable(namespace.RootContext(nil), me)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

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
		routerPath       string
		expectedMatch    string
	}{
		"enable external credential plugin, multiple versions available": {
			pluginType:       consts.PluginTypeCredential,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.0",
			routerPath:       "auth/foo/bar",
			expectedMatch:    "auth/foo/",
		},
		"enable external secrets plugin, multiple versions available": {
			pluginType:       consts.PluginTypeSecrets,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.0",
			routerPath:       "foo/bar",
			expectedMatch:    "foo/",
		},
		"enable external credential plugin, multiple versions available, select other version": {
			pluginType:       consts.PluginTypeCredential,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.1",
			routerPath:       "auth/foo/bar",
			expectedMatch:    "auth/foo/",
		},
		"enable external secrets plugin, multiple versions available, select other version": {
			pluginType:       consts.PluginTypeSecrets,
			registerVersions: []string{"v1.0.0", "v1.0.1"},
			mountVersion:     "v1.0.1",
			routerPath:       "foo/bar",
			expectedMatch:    "foo/",
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, pluginName, pluginSHA256 := testCoreWithPlugin(t, tc.pluginType)
			for _, version := range tc.registerVersions {
				d := &framework.FieldData{
					Raw: map[string]interface{}{
						"name":    pluginName,
						"sha256":  pluginSHA256,
						"version": version,
						"command": pluginName,
					},
					Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
				}
				resp, err := c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
				if err != nil {
					t.Fatal(err)
				}
				if resp.Error() != nil {
					t.Fatalf("%#v", resp)
				}
			}

			me := &MountEntry{
				Table:   mountTable(tc.pluginType),
				Path:    "foo",
				Type:    pluginName,
				Version: tc.mountVersion,
			}
			enable := enableFunc(c, tc.pluginType)
			err := enable(namespace.RootContext(nil), me)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			match := c.router.MatchingMount(namespace.RootContext(nil), tc.routerPath)
			if match != tc.expectedMatch {
				t.Fatalf("missing mount, match: %q", match)
			}

			raw, _ := c.router.root.Get(match)
			if raw.(*routeEntry).mountEntry.Version != tc.mountVersion {
				t.Errorf("Expected mount to be version %s but got %s", tc.mountVersion, raw.(*routeEntry).mountEntry.Version)
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
			d := &framework.FieldData{
				Raw: map[string]interface{}{
					"name":    pluginName,
					"sha256":  pluginSHA256,
					"command": pluginName,
				},
				Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
			}
			resp, err := c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Error() != nil {
				t.Fatalf("%#v", resp)
			}

			me := &MountEntry{
				Table: mountTable(tc.pluginType),
				Path:  "foo",
				Type:  pluginName,
			}
			enable := enableFunc(c, tc.pluginType)
			err = enable(namespace.RootContext(nil), me)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

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
			d := &framework.FieldData{
				Raw: map[string]interface{}{
					"name":    pluginName,
					"sha256":  pluginSHA256,
					"command": pluginName,
				},
				Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
			}
			resp, err := c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Error() != nil {
				t.Fatalf("%#v", resp)
			}

			me := &MountEntry{
				Table:   mountTable(tc.pluginType),
				Path:    "foo",
				Type:    pluginName,
				Version: "v1.0.0",
			}
			enable := enableFunc(c, tc.pluginType)
			err = enable(namespace.RootContext(nil), me)
			if err == nil || !errors.Is(err, ErrPluginNotFound) {
				t.Fatalf("Expected to get plugin not found but got: %v", err)
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

func mountTable(pluginType consts.PluginType) string {
	switch pluginType {
	case consts.PluginTypeCredential:
		return credentialTableType
	case consts.PluginTypeSecrets:
		return mountTableType
	default:
		panic("test does not support plugin type yet")
	}
}

func enableFunc(c *Core, pluginType consts.PluginType) func(context.Context, *MountEntry) error {
	switch pluginType {
	case consts.PluginTypeCredential:
		return c.enableCredential
	case consts.PluginTypeSecrets:
		return c.mount
	default:
		panic(pluginType.String())
	}
}
