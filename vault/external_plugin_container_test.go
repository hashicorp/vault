// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginruntimeutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func testClusterWithContainerPlugin(t *testing.T, pluginType consts.PluginType, version string) (*Core, pluginhelpers.TestPlugin) {
	coreConfig := &CoreConfig{
		CredentialBackends: map[string]logical.Factory{},
	}

	cluster := NewTestCluster(t, coreConfig, &TestClusterOptions{
		Plugins: &TestPluginConfig{
			Typ:       pluginType,
			Versions:  []string{version},
			Container: true,
		},
	})

	cluster.Start()
	t.Cleanup(cluster.Cleanup)

	c := cluster.Cores[0].Core
	TestWaitActive(t, c)
	plugins := cluster.Plugins

	return c, plugins[0]
}

func TestExternalPluginInContainer_MountAndUnmount(t *testing.T) {
	for name, tc := range map[string]struct {
		pluginType consts.PluginType
	}{
		"auth": {
			pluginType: consts.PluginTypeCredential,
		},
		"secrets": {
			pluginType: consts.PluginTypeSecrets,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, plugin := testClusterWithContainerPlugin(t, tc.pluginType, "v1.0.0")

			t.Run("default", func(t *testing.T) {
				if _, err := exec.LookPath("runsc"); err != nil {
					t.Skip("Skipping test as runsc not found on path")
				}
				mountAndUnmountContainerPlugin_WithRuntime(t, c, plugin, "")
			})

			t.Run("runc", func(t *testing.T) {
				mountAndUnmountContainerPlugin_WithRuntime(t, c, plugin, "runc")
			})

			t.Run("runsc", func(t *testing.T) {
				if _, err := exec.LookPath("runsc"); err != nil {
					t.Skip("Skipping test as runsc not found on path")
				}
				mountAndUnmountContainerPlugin_WithRuntime(t, c, plugin, "runsc")
			})
		})
	}
}

func mountAndUnmountContainerPlugin_WithRuntime(t *testing.T, c *Core, plugin pluginhelpers.TestPlugin, ociRuntime string) {
	if ociRuntime != "" {
		registerPluginRuntime(t, c.systemBackend, ociRuntime, ociRuntime)
	}
	registerContainerPlugin(t, c.systemBackend, plugin.Name, plugin.Typ.String(), "1.0.0", plugin.ImageSha256, plugin.Image, ociRuntime)

	mountPlugin(t, c.systemBackend, plugin.Name, plugin.Typ, "v1.0.0", "")

	routeRequest := func(expectMatch bool) {
		pluginPath := "foo/bar"
		if plugin.Typ == consts.PluginTypeCredential {
			pluginPath = "auth/foo/bar"
		}
		match := c.router.MatchingMount(namespace.RootContext(nil), pluginPath)
		if expectMatch && match != strings.TrimSuffix(pluginPath, "bar") {
			t.Fatalf("missing mount, match: %q", match)
		}
		if !expectMatch && match != "" {
			t.Fatalf("expected no match for path, but got %q", match)
		}
	}

	routeRequest(true)
	unmountPlugin(t, c.systemBackend, plugin.Typ, "foo")
	routeRequest(false)
}

func TestExternalPluginInContainer_GetBackendTypeVersion(t *testing.T) {
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
			c, plugin := testClusterWithContainerPlugin(t, tc.pluginType, tc.setRunningVersion)
			for _, ociRuntime := range []string{"runc", "runsc"} {
				t.Run(ociRuntime, func(t *testing.T) {
					if _, err := exec.LookPath(ociRuntime); err != nil {
						t.Skipf("Skipping test as %s not found on path", ociRuntime)
					}
					shaBytes, _ := hex.DecodeString(plugin.ImageSha256)
					entry := &pluginutil.PluginRunner{
						Name:     plugin.Name,
						OCIImage: plugin.Image,
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
		})
	}
}

func registerContainerPlugin(t *testing.T, sys *SystemBackend, pluginName, pluginType, version, sha, image, runtime string) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("plugins/catalog/%s/%s", pluginType, pluginName))
	req.Data = map[string]interface{}{
		"oci_image": image,
		"sha256":    sha,
		"version":   version,
		"runtime":   runtime,
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}

func registerPluginRuntime(t *testing.T, sys *SystemBackend, name, ociRuntime string) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("plugins/runtimes/catalog/%s/%s", consts.PluginRuntimeTypeContainer, name))
	req.Data = map[string]interface{}{
		"oci_runtime": ociRuntime,
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}
