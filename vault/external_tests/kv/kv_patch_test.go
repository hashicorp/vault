package kv

import (
	"context"
	"testing"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// Verifies that patching works by default with the root token
func TestKV_Patch_RootToken(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	// make sure this client is using the root token
	client.SetToken(cluster.RootToken)

	// Enable KVv2
	err := client.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write a kv value and patch it
	_, err = client.Logical().Write("kv/data/foo", map[string]interface{}{"data": map[string]interface{}{"bar": "baz"}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().JSONMergePatch(context.Background(), "kv/data/foo", map[string]interface{}{"data": map[string]interface{}{"bar": "quux"}})
	if err != nil {
		t.Fatal(err)
	}

	secret, err := client.Logical().Read("kv/data/foo")
	bar := secret.Data["data"].(map[string]interface{})["bar"]
	if bar != "quux" {
		t.Fatalf("expected bar to be quux but it was %q", bar)
	}
}
