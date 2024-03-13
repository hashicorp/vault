// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package builtinplugins

import (
	"testing"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	logicalDb "github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// TestBuiltinPluginsWork exists to confirm that all the credential and secrets plugins in Registry can successfully be
// initialized. Database plugins are excluded as there is no general way to initialize them - they require
// plugin-specific configuration at the time of initialization.
//
// This detects coding errors which would cause the plugins to panic on initialization - various aspects of the
// configuration of a framework.Backend are checked during Backend.init(), which runs as a sync.Once function triggered
// upon first request.
//
// In this test, a help request is used to trigger that initialization, since it is valid for all plugins.
func TestBuiltinPluginsWork(t *testing.T) {
	cluster := vault.NewTestCluster(
		t,
		&vault.CoreConfig{
			BuiltinRegistry: Registry,
			LogicalBackends: map[string]logical.Factory{
				// This needs to be here for madly overcomplicated reasons, otherwise we end up mounting a KV v1 even
				// when we try to explicitly mount a KV v2...
				//
				// vault.NewCore hardcodes "kv" to vault.PassthroughBackendFactory if no explicit entry is configured,
				// and this hardcoding is re-overridden in command.logicalBackends to point back to the real KV plugin.
				// As far as I can tell, nothing at all relies upon the definition of "kv" in builtinplugins.Registry,
				// as it always gets resolved via the logicalBackends map and the pluginCatalog is never queried.
				"kv": logicalKv.Factory,
				// Semi-similarly, "database" is added in command.logicalBackends and not at all in
				// builtinplugins.Registry, so we need to add it here to be able to test it!
				"database": logicalDb.Factory,
			},
			PendingRemovalMountsAllowed: true,
			// Specifying at least one audit backend factory will prevent NewTestCluster
			// from attempting to enable a noop audit, and audit isn't required for this test.
			AuditBackends: map[string]audit.Factory{
				"noop": corehelpers.NoopAuditFactory(nil),
			},
		},
		&vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
			NumCores:    1,
		},
	)

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	for _, authType := range append(
		Registry.Keys(consts.PluginTypeCredential),
		"token",
	) {
		deprecationStatus, _ := Registry.DeprecationStatus(authType, consts.PluginTypeCredential)
		if deprecationStatus == consts.Removed {
			continue
		}

		t.Run("Auth Method "+authType, func(t *testing.T) {
			// This builtin backend is automatically mounted and should not be mounted again
			if authType != "token" {
				if err := client.Sys().EnableAuthWithOptions(authType, &api.EnableAuthOptions{
					Type: authType,
				}); err != nil {
					t.Fatal(err)
				}
			}

			if _, err := client.Logical().ReadWithData(
				"auth/"+authType,
				map[string][]string{"help": {"1"}},
			); err != nil {
				t.Fatal(err)
			}
		})
	}

	for _, secretsType := range append(
		Registry.Keys(consts.PluginTypeSecrets),
		"database",
		"cubbyhole",
		"identity",
		"sys",
	) {
		deprecationStatus, _ := Registry.DeprecationStatus(secretsType, consts.PluginTypeSecrets)
		if deprecationStatus == consts.Removed {
			continue
		}

		t.Run("Secrets Engine "+secretsType, func(t *testing.T) {
			switch secretsType {
			// These three builtin backends are automatically mounted and should not be mounted again
			case "cubbyhole":
			case "identity":
			case "sys":

			default:
				if err := client.Sys().Mount(secretsType, &api.MountInput{
					Type: secretsType,
				}); err != nil {
					t.Fatal(err)
				}
			}

			if _, err := client.Logical().ReadWithData(
				secretsType,
				map[string][]string{"help": {"1"}},
			); err != nil {
				t.Fatal(err)
			}
		})
	}

	t.Run("Secrets Engine kv v2", func(t *testing.T) {
		if err := client.Sys().Mount("kv-v2", &api.MountInput{
			Type: "kv",
			Options: map[string]string{
				"version": "2",
			},
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := client.Logical().ReadWithData(
			"kv-v2",
			map[string][]string{"help": {"1"}},
		); err != nil {
			t.Fatal(err)
		}
	})

	// This last part is not strictly necessary for original purpose of this test (checking the plugins initialize
	// without errors), but whilst we have a test Vault with one of everything mounted, let's also test that the full
	// OpenAPI document generation succeeds too.
	t.Run("Whole OpenAPI document", func(t *testing.T) {
		if _, err := client.Logical().Read("sys/internal/specs/openapi"); err != nil {
			t.Fatal(err)
		}
	})
}
