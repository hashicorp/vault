package token

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// Tests the regression in
// https://github.com/hashicorp/vault/pull/6920
func TestRecoverFromPanic(t *testing.T) {
	logger := hclog.New(nil)

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"noop": vault.NoopBackendFactory,
		},
		EnableRaw: true,
		Logger:    logger,
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	err := client.Sys().Mount("noop", &api.MountInput{
		Type: "noop",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("noop/panic")
	if err == nil {
		t.Fatal("expected error")
	}

	// This will deadlock the test if we hit the condition
	cluster.EnsureCoresSealed(t)
}
