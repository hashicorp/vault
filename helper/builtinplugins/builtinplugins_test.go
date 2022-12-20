package builtinplugins

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
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
			Logger:                      logging.NewVaultLogger(hclog.Trace),
			BuiltinRegistry:             Registry,
			PendingRemovalMountsAllowed: true,
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

	for _, authType := range Registry.Keys(consts.PluginTypeCredential) {
		deprecationStatus, ok := Registry.DeprecationStatus(authType, consts.PluginTypeCredential)
		if !ok || deprecationStatus == consts.Removed {
			continue
		}

		if err := client.Sys().EnableAuthWithOptions(authType, &api.EnableAuthOptions{
			Type: authType,
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := client.Logical().ReadWithData(
			"auth/"+authType,
			map[string][]string{"help": {"1"}},
		); err != nil {
			t.Fatal(err)
		}
	}

	for _, secretsType := range Registry.Keys(consts.PluginTypeSecrets) {
		deprecationStatus, ok := Registry.DeprecationStatus(secretsType, consts.PluginTypeSecrets)
		if !ok || deprecationStatus == consts.Removed {
			continue
		}

		if err := client.Sys().Mount(secretsType, &api.MountInput{
			Type: secretsType,
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := client.Logical().ReadWithData(
			secretsType,
			map[string][]string{"help": {"1"}},
		); err != nil {
			t.Fatal(err)
		}
	}
}
