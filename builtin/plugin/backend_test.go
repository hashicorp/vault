// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	logicalPlugin "github.com/hashicorp/vault/sdk/plugin"
	"github.com/hashicorp/vault/sdk/plugin/mock"
	"github.com/hashicorp/vault/vault"
)

func TestBackend_impl(t *testing.T) {
	var _ logical.Backend = &plugin.PluginBackend{}
}

func TestBackend(t *testing.T) {
	pluginCmds := []string{"TestBackend_PluginMain", "TestBackend_PluginMain_Multiplexed"}

	for _, pluginCmd := range pluginCmds {
		t.Run(pluginCmd, func(t *testing.T) {
			config, cleanup := testConfig(t, pluginCmd)
			defer cleanup()

			_, err := plugin.Backend(context.Background(), config)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBackend_Factory(t *testing.T) {
	pluginCmds := []string{"TestBackend_PluginMain", "TestBackend_PluginMain_Multiplexed"}

	for _, pluginCmd := range pluginCmds {
		t.Run(pluginCmd, func(t *testing.T) {
			config, cleanup := testConfig(t, pluginCmd)
			defer cleanup()

			_, err := plugin.Factory(context.Background(), config)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBackend_PluginMain(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := logicalPlugin.Serve(&logicalPlugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMain_Multiplexed(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := logicalPlugin.ServeMultiplex(&logicalPlugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func testConfig(t *testing.T, pluginCmd string) (*logical.BackendConfig, func()) {
	t.Helper()
	pluginDir := corehelpers.MakeTestPluginDir(t)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		PluginDirectory: pluginDir,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	cores := cluster.Cores

	core := cores[0]

	sys := vault.TestDynamicSystemView(core.Core, nil)

	config := &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Debug),
		System: sys,
		Config: map[string]string{
			"plugin_name":    "mock-plugin",
			"plugin_type":    "secret",
			"plugin_version": "v0.0.0+mock",
		},
	}

	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "", pluginCmd,
		[]string{fmt.Sprintf("%s=%s", pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)})

	return config, func() {
		cluster.Cleanup()
	}
}
