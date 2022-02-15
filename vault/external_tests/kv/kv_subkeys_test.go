package kv

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestKV_Subkeys_NotFound(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	// Mount a KVv2 backend
	err := c.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := c.NewRequest("GET", "/v1/kv/subkeys/foo")
	apiResp, err := c.RawRequestWithContext(context.Background(), req)

	if err == nil || apiResp == nil {
		t.Fatalf("expected subkeys request to fail, err :%v, resp: %#v", err, apiResp)
	}

	if apiResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected subkeys request to fail with %d status code, resp: %#v", http.StatusNotFound, apiResp)
	}
}

func TestKV_Subkeys_Deleted(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	// Mount a KVv2 backend
	err := c.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	kvData := map[string]interface{}{
		"data": map[string]interface{}{
			"bar": "a",
		},
	}

	resp, err := c.Logical().Write("kv/data/foo", kvData)
	if err != nil {
		t.Fatalf("write failed, err :%v, resp: %#v\n", err, resp)
	}

	resp, err = c.Logical().Delete("kv/data/foo")
	if err != nil {
		t.Fatalf("delete failed, err :%v, resp: %#v\n", err, resp)
	}

	req := c.NewRequest("GET", "/v1/kv/subkeys/foo")
	apiResp, err := c.RawRequestWithContext(context.Background(), req)
	if resp != nil {
		defer apiResp.Body.Close()
	}

	if err == nil || apiResp == nil {
		t.Fatalf("expected subkeys request to fail, err :%v, resp: %#v", err, apiResp)
	}

	if apiResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected subkeys request to fail with %d status code, resp: %#v", http.StatusNotFound, apiResp)
	}

	secret, err := api.ParseSecret(apiResp.Body)
	if err != nil {
		t.Fatalf("failed to parse resp body, err: %v", err)
	}

	if subkeys, ok := secret.Data["subkeys"]; !ok || subkeys != nil {
		t.Fatalf("expected nil subkeys, actual: %#v", subkeys)
	}

	metadata, ok := secret.Data["metadata"].(map[string]interface{})

	if !ok {
		t.Fatalf("metadata not present in response or invalid, metadata: %#v", secret.Data["metadata"])
	}

	if deletionTime, ok := metadata["deletion_time"].(string); !ok || deletionTime == "" {
		t.Fatalf("metadata does not contain deletion time, metadata: %#v", metadata)
	}
}

func TestKV_Subkeys_Destroyed(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	// Mount a KVv2 backend
	err := c.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	kvData := map[string]interface{}{
		"data": map[string]interface{}{
			"bar": "a",
		},
	}

	resp, err := c.Logical().Write("kv/data/foo", kvData)
	if err != nil {
		t.Fatalf("write failed, err :%v, resp: %#v\n", err, resp)
	}

	destroyVersions := map[string]interface{}{
		"versions": []int{1},
	}

	resp, err = c.Logical().Write("kv/destroy/foo", destroyVersions)
	if err != nil {
		t.Fatalf("destroy failed, err :%v, resp: %#v\n", err, resp)
	}

	req := c.NewRequest("GET", "/v1/kv/subkeys/foo")
	apiResp, err := c.RawRequestWithContext(context.Background(), req)
	if resp != nil {
		defer apiResp.Body.Close()
	}

	if err == nil || apiResp == nil {
		t.Fatalf("expected subkeys request to fail, err :%v, resp: %#v", err, apiResp)
	}

	if apiResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected subkeys request to fail with %d status code, resp: %#v", http.StatusNotFound, apiResp)
	}

	secret, err := api.ParseSecret(apiResp.Body)
	if err != nil {
		t.Fatalf("failed to parse resp body, err: %v", err)
	}

	if subkeys, ok := secret.Data["subkeys"]; !ok || subkeys != nil {
		t.Fatalf("expected nil subkeys, actual: %#v", subkeys)
	}

	metadata, ok := secret.Data["metadata"].(map[string]interface{})

	if !ok {
		t.Fatalf("metadata not present in response or invalid, metadata: %#v", secret.Data["metadata"])
	}

	if destroyed, ok := metadata["destroyed"].(bool); !ok || !destroyed {
		t.Fatalf("expected destroyed to be true, metadata: %#v", metadata)
	}
}

func TestKV_Subkeys_CurrentVersion(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	// Mount a KVv2 backend
	err := c.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	kvData := map[string]interface{}{
		"data": map[string]interface{}{
			"foo": "does-not-matter",
			"bar": map[string]interface{}{
				"a": map[string]interface{}{
					"c": "does-not-matter",
				},
				"b": map[string]interface{}{},
			},
		},
	}

	resp, err := c.Logical().Write("kv/data/foo", kvData)
	if err != nil {
		t.Fatalf("write failed - err :%v, resp: %#v\n", err, resp)
	}

	resp, err = c.Logical().Read("kv/subkeys/foo")

	if err != nil {
		t.Fatalf("read failed - err :%v", err)
	}

	expectedSubkeys := map[string]interface{}{
		"foo": nil,
		"bar": map[string]interface{}{
			"a": map[string]interface{}{
				"c": nil,
			},
			"b": nil,
		},
	}

	if diff := deep.Equal(resp.Data["subkeys"], expectedSubkeys); len(diff) > 0 {
		t.Fatalf("resp and expected data mismatch, diff: %#v", diff)
	}
}
