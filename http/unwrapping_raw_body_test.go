package http

import (
	"testing"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestUnwrapping_Raw_Body(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// Mount a k/v backend, version 2
	err := client.Sys().Mount("kv", &api.MountInput{
		Type:    "kv",
		Options: map[string]string{"version": "2"},
	})
	if err != nil {
		t.Fatal(err)
	}

	client.SetWrappingLookupFunc(func(operation, path string) string {
		return "5m"
	})
	secret, err := client.Logical().Write("kv/foo/bar", map[string]interface{}{
		"a": "b",
	})
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil {
		t.Fatal("nil secret")
	}
	if secret.WrapInfo == nil {
		t.Fatal("nil wrap info")
	}
	wrapToken := secret.WrapInfo.Token

	client.SetWrappingLookupFunc(nil)
	secret, err = client.Logical().Unwrap(wrapToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(secret.Warnings) != 1 {
		t.Fatal("expected 1 warning")
	}
}
