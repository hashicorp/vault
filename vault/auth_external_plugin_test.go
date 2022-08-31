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
)

var (
	compileOnce sync.Once
	pluginBytes []byte
)

func testCoreWithPlugin(t *testing.T) (*Core, string, string) {
	t.Helper()
	pluginName, pluginSha256, pluginDir := compilePlugin(t)
	conf := &CoreConfig{
		BuiltinRegistry: NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}
	core := TestCoreWithSealAndUI(t, conf)
	core, _, _ = testCoreUnsealed(t, core)
	return core, pluginName, pluginSha256
}

// to mount a plugin, we need a working binary plugin, so we compile one here.
func compilePlugin(t *testing.T) (string, string, string) {
	pluginType := "approle"
	pluginName := "vault-plugin-auth-" + pluginType

	dir := ""
	// detect if we are in the "vault/" or the root directory and compensate
	if _, err := os.Stat("builtin"); os.IsNotExist(err) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dir = filepath.Dir(wd)
	}

	pluginDir, cleanup := MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })

	pluginPath := path.Join(pluginDir, pluginName)

	// cache the compilation to only run once
	compileOnce.Do(func() {
		cmd := exec.Command("go", "build", "-o", pluginPath, fmt.Sprintf("builtin/credential/%s/cmd/%s/main.go", pluginType, pluginType))
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("error running go build %v output: %s", err, output))
		}
		pluginBytes, err = os.ReadFile(pluginPath)
		if err != nil {
			panic(err)
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

func TestCore_EnableExternalCredentialPlugin(t *testing.T) {
	c, pluginName, pluginSha256 := testCoreWithPlugin(t)
	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
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
		Table:   credentialTableType,
		Path:    "foo",
		Type:    pluginName,
		Version: "v1.0.0",
	}
	err = c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}
}

func TestCore_EnableExternalCredentialPlugin_MultipleVersions(t *testing.T) {
	c, pluginName, pluginSha256 := testCoreWithPlugin(t)
	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
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

	d = &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
			"version": "v1.0.1",
			"command": pluginName,
		},
		Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
	}
	resp, err = c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() != nil {
		t.Fatalf("%#v", resp)
	}

	me := &MountEntry{
		Table:   credentialTableType,
		Path:    "foo",
		Type:    pluginName,
		Version: "v1.0.0",
	}
	err = c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}

	raw, _ := c.router.root.Get(match)
	if raw.(*routeEntry).mountEntry.Version != "v1.0.0" {
		t.Errorf("Expected mount to be version v1.0.0 but got %s", raw.(*routeEntry).mountEntry.Version)
	}
}

func TestCore_EnableExternalCredentialPlugin_MultipleVersions_MountSecond(t *testing.T) {
	c, pluginName, pluginSha256 := testCoreWithPlugin(t)
	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
			"command": pluginName,
			"version": "v1.0.0",
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

	d = &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
			"version": "v1.0.1",
			"command": pluginName,
		},
		Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
	}
	resp, err = c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error() != nil {
		t.Fatalf("%#v", resp)
	}

	me := &MountEntry{
		Table:   credentialTableType,
		Path:    "foo",
		Type:    pluginName,
		Version: "v1.0.1",
	}
	err = c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}

	raw, _ := c.router.root.Get(match)
	if raw.(*routeEntry).mountEntry.Version != "v1.0.1" {
		t.Errorf("Expected mount to be version v1.0.1 but got %s", raw.(*routeEntry).mountEntry.Version)
	}
}

func TestCore_EnableExternalCredentialPlugin_NoVersionsOkay(t *testing.T) {
	c, pluginName, pluginSha256 := testCoreWithPlugin(t)
	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
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
		Table: credentialTableType,
		Path:  "foo",
		Type:  pluginName,
	}
	err = c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}
}

func TestCore_EnableExternalCredentialPlugin_NoVersionOnRegister(t *testing.T) {
	c, pluginName, pluginSha256 := testCoreWithPlugin(t)
	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
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
		Table:   credentialTableType,
		Path:    "foo",
		Type:    pluginName,
		Version: "v1.0.0",
	}
	err = c.enableCredential(namespace.RootContext(nil), me)
	if err == nil || !errors.Is(err, ErrPluginNotFound) {
		t.Fatalf("Expected to get plugin not found but got: %v", err)
	}
}

func TestCore_EnableExternalCredentialPlugin_InvalidName(t *testing.T) {
	c, pluginName, pluginSha256 := testCoreWithPlugin(t)
	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":    pluginName,
			"sha256":  pluginSha256,
			"version": "v1.0.0",
			"command": pluginName + "xyz",
		},
		Schema: c.systemBackend.pluginsCatalogCRUDPath().Fields,
	}
	_, err := c.systemBackend.handlePluginCatalogUpdate(context.Background(), nil, d)
	if err == nil || !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("should have gotten a no such file or directory error inserting the plugin: %v", err)
	}
}
