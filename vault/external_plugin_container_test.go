// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

func testClusterWithContainerPlugins(t *testing.T, types []consts.PluginType) (*TestClusterCore, []pluginhelpers.TestPlugin) {
	var plugins []*TestPluginConfig
	for _, typ := range types {
		plugins = append(plugins, &TestPluginConfig{
			Typ:       typ,
			Versions:  []string{"v1.0.0"},
			Container: true,
		})
	}
	// Use os.MkdirTemp because t.TempDir() exceeds the Unix socket length limit.
	// See https://www.man7.org/linux/man-pages/man7/unix.7.html for details.
	tmpdir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Fatal(err)
		}
	})
	cluster := NewTestCluster(t, &CoreConfig{
		PluginTmpdir: tmpdir,
	}, &TestClusterOptions{
		Plugins: plugins,
	})

	cluster.Start()
	t.Cleanup(cluster.Cleanup)

	core := cluster.Cores[0]
	TestWaitActive(t, core.Core)

	return core, cluster.Plugins
}

func TestExternalPluginInContainer_MountAndUnmount(t *testing.T) {
	t.Run("rootful docker runtimes", func(t *testing.T) {
		t.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
		c, plugins := testClusterWithContainerPlugins(t, []consts.PluginType{
			consts.PluginTypeCredential,
			consts.PluginTypeSecrets,
		})

		for _, plugin := range plugins {
			t.Run(plugin.Typ.String(), func(t *testing.T) {
				t.Run("default runtime", func(t *testing.T) {
					if _, err := exec.LookPath("runsc"); err != nil {
						t.Skip("Skipping test as runsc not found on path")
					}
					mountAndUnmountContainerPlugin_WithRuntime(t, c, plugin, "", false)
				})

				t.Run("runc", func(t *testing.T) {
					mountAndUnmountContainerPlugin_WithRuntime(t, c, plugin, "runc", false)
				})

				t.Run("runsc", func(t *testing.T) {
					if _, err := exec.LookPath("runsc"); err != nil {
						t.Skip("Skipping test as runsc not found on path")
					}
					mountAndUnmountContainerPlugin_WithRuntime(t, c, plugin, "runsc", false)
				})
			})
		}
	})

	t.Run("rootless runsc", func(t *testing.T) {
		if _, err := exec.LookPath("runsc"); err != nil {
			t.Skip("Skipping test as runsc not found on path")
		}

		t.Setenv("DOCKER_HOST", fmt.Sprintf("unix:///run/user/%d/docker.sock", os.Getuid()))
		c, plugins := testClusterWithContainerPlugins(t, []consts.PluginType{consts.PluginTypeCredential})
		mountAndUnmountContainerPlugin_WithRuntime(t, c, plugins[0], "runsc", true)

	})
}

func mountAndUnmountContainerPlugin_WithRuntime(t *testing.T, c *TestClusterCore, plugin pluginhelpers.TestPlugin, ociRuntime string, rootless bool) {
	if ociRuntime != "" {
		registerPluginRuntime(t, c.systemBackend, ociRuntime, rootless)
	}
	registerContainerPlugin(t, c.systemBackend, plugin.Name, plugin.Typ.String(), "1.0.0", plugin.ImageSha256, plugin.Image, ociRuntime)

	mountPlugin(t, c.systemBackend, plugin.Name, plugin.Typ, "v1.0.0", "")

	expectTmpdirEntries := func(expected int) {
		t.Helper()
		entries, err := os.ReadDir(c.CoreConfig.PluginTmpdir)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != expected {
			t.Fatalf("expected %d in tmpdir, got %v", expected, entries)
		}
	}
	expectTmpdirEntries(1)

	routeRequest := func(expectMatch bool) {
		pluginPath := "foo/bar"
		if plugin.Typ == consts.PluginTypeCredential {
			pluginPath = "auth/foo/bar"
		}
		match := c.router.MatchingMount(namespace.RootContext(context.Background()), pluginPath)
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
	expectTmpdirEntries(0)
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
	resp, err := sys.HandleRequest(namespace.RootContext(context.Background()), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}

func registerPluginRuntime(t *testing.T, sys *SystemBackend, ociRuntime string, rootless bool) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("plugins/runtimes/catalog/%s/%s", consts.PluginRuntimeTypeContainer, ociRuntime))
	req.Data = map[string]interface{}{
		"oci_runtime": ociRuntime,
		"rootless":    rootless,
	}
	resp, err := sys.HandleRequest(namespace.RootContext(context.Background()), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}
