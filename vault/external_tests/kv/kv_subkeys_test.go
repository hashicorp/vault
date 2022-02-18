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

// TestKV_Subkeys_NotFound issues a read to the subkeys endpoint for a path
// that does not exist. A 400 status should be returned.
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

	apiRespRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		req := c.NewRequest("GET", "/v1/kv/subkeys/foo")
		return c.RawRequestWithContext(context.Background(), req)
	})

	apiResp, ok := apiRespRaw.(*api.Response)
	if !ok {
		t.Fatalf("response not an api.Response, actual: %#v", apiRespRaw)
	}

	if err == nil || apiResp == nil {
		t.Fatalf("expected subkeys request to fail, err :%v, resp: %#v", err, apiResp)
	}

	if apiResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected subkeys request to fail with %d status code, resp: %#v", http.StatusNotFound, apiResp)
	}
}

// TestKV_Subkeys_Deleted writes a single version of a secret to the KVv2
// secret engine. The secret is subsequently deleted. A read to the subkeys
// endpoint should return a 400 status with a nil "subkeys" value and the
// "deletion_time" key in the "metadata" key should be not be empty.
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

	resp, err := kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Write("kv/data/foo", kvData)
	})
	if err != nil {
		t.Fatalf("write failed, err :%v, resp: %#v", err, resp)
	}

	secretRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Delete("kv/data/foo")
	})
	if err != nil {
		t.Fatalf("delete failed, err :%v, resp: %#v", err, secretRaw)
	}

	apiRespRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		req := c.NewRequest("GET", "/v1/kv/subkeys/foo")
		return c.RawRequestWithContext(context.Background(), req)
	})

	apiResp, ok := apiRespRaw.(*api.Response)
	if !ok {
		t.Fatalf("response not a api.Response, actual: %#v", apiRespRaw)
	}

	if apiResp != nil {
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

	subkeys, ok := secret.Data["subkeys"]
	if !ok {
		t.Fatalf("key \"subkeys\" not found in response")
	}

	if subkeys != nil {
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

// TestKV_Subkeys_Destroyed writes a single version of a secret to the KVv2
// secret engine. The secret is subsequently destroyed. A read to the subkeys
// endpoint should return a 400 status with a nil "subkeys" value and the
// "destroyed" key in the "metadata" key should be set to true.
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

	secretRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Write("kv/data/foo", kvData)
	})
	if err != nil {
		t.Fatalf("write failed, err :%v, resp: %#v", err, secretRaw)
	}

	destroyVersions := map[string]interface{}{
		"versions": []int{1},
	}

	secretRaw, err = kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Write("kv/destroy/foo", destroyVersions)
	})
	if err != nil {
		t.Fatalf("destroy failed, err :%v, resp: %#v", err, secretRaw)
	}

	secret, ok := secretRaw.(*api.Secret)
	if !ok {
		t.Fatalf("response not an api.Secret, actual: %#v", secretRaw)
	}

	apiRespRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		req := c.NewRequest("GET", "/v1/kv/subkeys/foo")
		return c.RawRequestWithContext(context.Background(), req)
	})

	apiResp, ok := apiRespRaw.(*api.Response)
	if !ok {
		t.Fatalf("response not a api.Response, actual: %#v", apiRespRaw)
	}

	if apiResp != nil {
		defer apiResp.Body.Close()
	}

	if err == nil || apiResp == nil {
		t.Fatalf("expected subkeys request to fail, err :%v, resp: %#v", err, apiResp)
	}

	if apiResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected subkeys request to fail with %d status code, resp: %#v", http.StatusNotFound, apiResp)
	}

	secret, err = api.ParseSecret(apiResp.Body)
	if err != nil {
		t.Fatalf("failed to parse resp body, err: %v", err)
	}

	subkeys, ok := secret.Data["subkeys"]
	if !ok {
		t.Fatalf("key \"subkeys\" not found in response")
	}

	if subkeys != nil {
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

// TestKV_Subkeys_CurrentVersion writes multiples versions of a secret to the
// KVv2 secret engine. It ensures that the subkeys endpoint returns a 200 status
// and current version of the secret.
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

	secretRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Write("kv/data/foo", kvData)
	})
	if err != nil {
		t.Fatalf("write failed, err :%v, resp: %#v", err, secretRaw)
	}

	kvData = map[string]interface{}{
		"data": map[string]interface{}{
			"baz": "does-not-matter",
		},
	}

	secretRaw, err = kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().JSONMergePatch(context.Background(), "kv/data/foo", kvData)
	})
	if err != nil {
		t.Fatalf("patch failed, err :%v, resp: %#v", err, secretRaw)
	}

	apiRespRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		req := c.NewRequest("GET", "/v1/kv/subkeys/foo")
		return c.RawRequestWithContext(context.Background(), req)
	})

	apiResp, ok := apiRespRaw.(*api.Response)
	if !ok {
		t.Fatalf("response not a api.Response, actual: %#v", apiRespRaw)
	}

	if apiResp != nil {
		defer apiResp.Body.Close()
	}

	if err != nil || apiResp == nil {
		t.Fatalf("subkeys request failed, err :%v, resp: %#v", err, apiResp)
	}

	if apiResp.StatusCode != http.StatusOK {
		t.Fatalf("expected subkeys request to succeed with %d status code, resp: %#v", http.StatusOK, apiResp)
	}

	secret, err := api.ParseSecret(apiResp.Body)
	if err != nil {
		t.Fatalf("failed to parse resp body, err: %v", err)
	}

	expectedSubkeys := map[string]interface{}{
		"foo": nil,
		"bar": map[string]interface{}{
			"a": map[string]interface{}{
				"c": nil,
			},
			"b": nil,
		},
		"baz": nil,
	}

	if diff := deep.Equal(secret.Data["subkeys"], expectedSubkeys); len(diff) > 0 {
		t.Fatalf("resp and expected data mismatch, diff: %#v", diff)
	}
}
