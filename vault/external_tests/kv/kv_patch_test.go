// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kv

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestKV_Patch_BadContentTypeHeader(t *testing.T) {
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
		t.Fatalf("write failed - err :%#v, resp: %#v\n", err, secretRaw)
	}

	secretRaw, err = kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Read("kv/data/foo")
	})
	if err != nil {
		t.Fatalf("read failed - err :%#v, resp: %#v\n", err, secretRaw)
	}

	apiRespRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		req := c.NewRequest("PATCH", "/v1/kv/data/foo")
		req.Headers = http.Header{
			"Content-Type": []string{"application/json"},
		}

		if err := req.SetJSONBody(kvData); err != nil {
			t.Fatal(err)
		}

		return c.RawRequestWithContext(context.Background(), req)
	})

	apiResp, ok := apiRespRaw.(*api.Response)
	if !ok {
		t.Fatalf("response not an api.Response, actual: %#v", apiRespRaw)
	}

	if err == nil || apiResp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected PATCH request to fail with %d status code - err :%#v, resp: %#v\n", http.StatusUnsupportedMediaType, err, apiResp)
	}
}

func kvRequestWithRetry(t *testing.T, req func() (interface{}, error)) (interface{}, error) {
	t.Helper()

	var err error
	var resp interface{}

	// Loop until return message does not indicate upgrade, or timeout.
	timeout := time.After(20 * time.Second)
	ticker := time.Tick(time.Second)

	for {
		select {
		case <-timeout:
			t.Error("timeout expired waiting for upgrade")
		case <-ticker:
			resp, err = req()

			if err == nil {
				return resp, nil
			}

			responseError := err.(*api.ResponseError)
			if !strings.Contains(responseError.Error(), "Upgrading from non-versioned to versioned data") {
				return resp, err
			}
		}
	}
}

func TestKV_Patch_Audit(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
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

	if err := c.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	auditLogFile, err := ioutil.TempFile("", "httppatch")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": auditLogFile.Name(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	writeData := map[string]interface{}{
		"data": map[string]interface{}{
			"bar": "a",
		},
	}

	resp, err := kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Write("kv/data/foo", writeData)
	})
	if err != nil {
		t.Fatalf("write request failed, err: %#v, resp: %#v\n", err, resp)
	}

	patchData := map[string]interface{}{
		"data": map[string]interface{}{
			"baz": "b",
		},
	}

	resp, err = kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().JSONMergePatch(context.Background(), "kv/data/foo", patchData)
	})

	if err != nil {
		t.Fatalf("patch request failed, err: %#v, resp: %#v\n", err, resp)
	}

	patchRequestLogCount := 0
	patchResponseLogCount := 0
	decoder := json.NewDecoder(auditLogFile)

	var auditRecord map[string]interface{}
	for decoder.Decode(&auditRecord) == nil {
		auditRequest := map[string]interface{}{}

		if req, ok := auditRecord["request"]; ok {
			auditRequest = req.(map[string]interface{})
		}

		if auditRequest["operation"] == "patch" && auditRecord["type"] == "request" {
			patchRequestLogCount += 1
		} else if auditRequest["operation"] == "patch" && auditRecord["type"] == "response" {
			patchResponseLogCount += 1
		}
	}

	if patchRequestLogCount != 1 {
		t.Fatalf("expected 1 patch request audit log record, saw %d\n", patchRequestLogCount)
	}

	if patchResponseLogCount != 1 {
		t.Fatalf("expected 1 patch response audit log record, saw %d\n", patchResponseLogCount)
	}
}

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
	_, err = kvRequestWithRetry(t, func() (interface{}, error) {
		data := map[string]interface{}{
			"data": map[string]interface{}{
				"bar": "baz",
				"foo": "qux",
			},
		}

		return client.Logical().Write("kv/data/foo", data)
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = kvRequestWithRetry(t, func() (interface{}, error) {
		data := map[string]interface{}{
			"data": map[string]interface{}{
				"bar": "quux",
				"foo": nil,
			},
		}
		return client.Logical().JSONMergePatch(context.Background(), "kv/data/foo", data)
	})

	if err != nil {
		t.Fatal(err)
	}

	secretRaw, err := kvRequestWithRetry(t, func() (interface{}, error) {
		return client.Logical().Read("kv/data/foo")
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, ok := secretRaw.(*api.Secret)
	if !ok {
		t.Fatalf("response not an api.Secret, actual: %#v", secretRaw)
	}

	bar := secret.Data["data"].(map[string]interface{})["bar"]
	if bar != "quux" {
		t.Fatalf("expected bar to be quux but it was %q", bar)
	}

	if _, ok := secret.Data["data"].(map[string]interface{})["foo"]; ok {
		t.Fatalf("expected data not to include foo")
	}
}
