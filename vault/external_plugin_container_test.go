// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
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
		pluginType    consts.PluginType
		routerPath    string
		expectedMatch string
		listRolesPath string
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
			c, plugin := testClusterWithContainerPlugin(t, tc.pluginType, "v1.0.0")

			registerContainerPlugin(t, c.systemBackend, plugin.Name, tc.pluginType.String(), "1.0.0", plugin.ImageSha256, plugin.Image)

			mountPlugin(t, c.systemBackend, plugin.Name, tc.pluginType, "v1.0.0", "")

			match := c.router.MatchingMount(namespace.RootContext(nil), tc.routerPath)
			if match != tc.expectedMatch {
				t.Fatalf("missing mount, match: %q", match)
			}

			unmountPlugin(t, c.systemBackend, plugin.Name, tc.pluginType, "v1.0.0", "foo")
		})
	}
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
			registerContainerPlugin(t, c.systemBackend, plugin.Name, tc.pluginType.String(), tc.setRunningVersion, plugin.ImageSha256, plugin.Image)

			shaBytes, _ := hex.DecodeString(plugin.ImageSha256)
			entry := &pluginutil.PluginRunner{
				Name:     plugin.Name,
				OCIImage: plugin.Image,
				Args:     nil,
				Sha256:   shaBytes,
				Builtin:  false,
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

func registerContainerPlugin(t *testing.T, sys *SystemBackend, pluginName, pluginType, version, sha, image string) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("plugins/catalog/%s/%s", pluginType, pluginName))
	req.Data = map[string]interface{}{
		"oci_image": image,
		"sha256":    sha,
		"version":   version,
	}
	resp, err := sys.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}
