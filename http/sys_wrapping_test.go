// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/vault"
)

// Test wrapping functionality
func TestHTTP_Wrapping(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	client := cores[0].Client
	client.SetToken(cluster.RootToken)

	// Write a value that we will use with wrapping for lookup
	_, err := client.Logical().Write("secret/foo", map[string]interface{}{
		"zip": "zap",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set a wrapping lookup function for reads on that path
	client.SetWrappingLookupFunc(func(operation, path string) string {
		if operation == "GET" && path == "secret/foo" {
			return "5m"
		}

		return api.DefaultWrappingLookupFunc(operation, path)
	})

	// First test: basic things that should fail, lookup edition
	// Root token isn't a wrapping token
	_, err = client.Logical().Write("sys/wrapping/lookup", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	// Not supplied
	_, err = client.Logical().Write("sys/wrapping/lookup", map[string]interface{}{
		"foo": "bar",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	// Nonexistent token isn't a wrapping token
	_, err = client.Logical().Write("sys/wrapping/lookup", map[string]interface{}{
		"token": "bar",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Second: basic things that should fail, unwrap edition
	// Root token isn't a wrapping token
	_, err = client.Logical().Unwrap(cluster.RootToken)
	if err == nil {
		t.Fatal("expected error")
	}
	// Root token isn't a wrapping token
	_, err = client.Logical().Write("sys/wrapping/unwrap", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	// Not supplied
	_, err = client.Logical().Write("sys/wrapping/unwrap", map[string]interface{}{
		"foo": "bar",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	// Nonexistent token isn't a wrapping token
	_, err = client.Logical().Write("sys/wrapping/unwrap", map[string]interface{}{
		"token": "bar",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	//
	// Test lookup
	//

	// Create a wrapping token
	secret, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.WrapInfo == nil {
		t.Fatal("secret or wrap info is nil")
	}
	wrapInfo := secret.WrapInfo

	// Test this twice to ensure no ill effect to the wrapping token as a result of the lookup
	for i := 0; i < 2; i++ {
		secret, err = client.Logical().Write("sys/wrapping/lookup", map[string]interface{}{
			"token": wrapInfo.Token,
		})
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil || secret.Data == nil {
			t.Fatal("secret or secret data is nil")
		}
		creationTTL, _ := secret.Data["creation_ttl"].(json.Number).Int64()
		if int(creationTTL) != wrapInfo.TTL {
			t.Fatalf("mismatched ttls: %d vs %d", creationTTL, wrapInfo.TTL)
		}
		if secret.Data["creation_time"].(string) != wrapInfo.CreationTime.Format(time.RFC3339Nano) {
			t.Fatalf("mismatched creation times: %q vs %q", secret.Data["creation_time"].(string), wrapInfo.CreationTime.Format(time.RFC3339Nano))
		}
	}

	//
	// Test unwrap
	//

	// Create a wrapping token
	secret, err = client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.WrapInfo == nil {
		t.Fatal("secret or wrap info is nil")
	}
	wrapInfo = secret.WrapInfo

	// Test unwrap via the client token
	client.SetToken(wrapInfo.Token)
	secret, err = client.Logical().Write("sys/wrapping/unwrap", nil)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Warnings != nil {
		t.Fatalf("Warnings found: %v", secret.Warnings)
	}
	if secret == nil || secret.Data == nil {
		t.Fatal("secret or secret data is nil")
	}
	ret1 := secret
	// Should be expired and fail
	_, err = client.Logical().Write("sys/wrapping/unwrap", nil)
	if err == nil {
		t.Fatal("expected err")
	}

	// Create a wrapping token
	client.SetToken(cluster.RootToken)
	secret, err = client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.WrapInfo == nil {
		t.Fatal("secret or wrap info is nil")
	}
	wrapInfo = secret.WrapInfo

	// Test as a separate token
	secret, err = client.Logical().Write("sys/wrapping/unwrap", map[string]interface{}{
		"token": wrapInfo.Token,
	})
	if err != nil {
		t.Fatal(err)
	}
	ret2 := secret
	// Should be expired and fail
	_, err = client.Logical().Write("sys/wrapping/unwrap", map[string]interface{}{
		"token": wrapInfo.Token,
	})
	if err == nil {
		t.Fatal("expected err")
	}

	// Create a wrapping token
	secret, err = client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.WrapInfo == nil {
		t.Fatal("secret or wrap info is nil")
	}
	wrapInfo = secret.WrapInfo

	// Read response directly
	client.SetToken(wrapInfo.Token)
	secret, err = client.Logical().Read("cubbyhole/response")
	if err != nil {
		t.Fatal(err)
	}
	ret3 := secret
	// Should be expired and fail
	_, err = client.Logical().Write("cubbyhole/response", nil)
	if err == nil {
		t.Fatal("expected err")
	}

	// Create a wrapping token
	client.SetToken(cluster.RootToken)
	secret, err = client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.WrapInfo == nil {
		t.Fatal("secret or wrap info is nil")
	}
	wrapInfo = secret.WrapInfo

	// Read via Unwrap method
	secret, err = client.Logical().Unwrap(wrapInfo.Token)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Warnings != nil {
		t.Fatalf("Warnings found: %v", secret.Warnings)
	}
	ret4 := secret
	// Should be expired and fail
	_, err = client.Logical().Unwrap(wrapInfo.Token)
	if err == nil {
		t.Fatal("expected err")
	}

	if !reflect.DeepEqual(ret1.Data, map[string]interface{}{
		"zip": "zap",
	}) {
		t.Fatalf("ret1 data did not match expected: %#v", ret1.Data)
	}
	if !reflect.DeepEqual(ret2.Data, map[string]interface{}{
		"zip": "zap",
	}) {
		t.Fatalf("ret2 data did not match expected: %#v", ret2.Data)
	}
	var ret3Secret api.Secret
	err = jsonutil.DecodeJSON([]byte(ret3.Data["response"].(string)), &ret3Secret)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ret3Secret.Data, map[string]interface{}{
		"zip": "zap",
	}) {
		t.Fatalf("ret3 data did not match expected: %#v", ret3Secret.Data)
	}
	if !reflect.DeepEqual(ret4.Data, map[string]interface{}{
		"zip": "zap",
	}) {
		t.Fatalf("ret4 data did not match expected: %#v", ret4.Data)
	}

	//
	// Custom wrapping
	//

	client.SetToken(cluster.RootToken)
	data := map[string]interface{}{
		"zip":   "zap",
		"three": json.Number("2"),
	}

	// Don't set a request TTL on that path, should fail
	client.SetWrappingLookupFunc(func(operation, path string) string {
		return ""
	})
	secret, err = client.Logical().Write("sys/wrapping/wrap", data)
	if err == nil {
		t.Fatal("expected error")
	}

	// Re-set the lookup function
	client.SetWrappingLookupFunc(func(operation, path string) string {
		if operation == "GET" && path == "secret/foo" {
			return "5m"
		}

		return api.DefaultWrappingLookupFunc(operation, path)
	})
	secret, err = client.Logical().Write("sys/wrapping/wrap", data)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Warnings != nil {
		t.Fatalf("Warnings found: %v", secret.Warnings)
	}
	secret, err = client.Logical().Unwrap(secret.WrapInfo.Token)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Warnings != nil {
		t.Fatalf("Warnings found: %v", secret.Warnings)
	}
	if !reflect.DeepEqual(data, secret.Data) {
		t.Fatalf("custom wrap did not match expected: %#v", secret.Data)
	}

	//
	// Test rewrap
	//

	// Create a wrapping token
	secret, err = client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.WrapInfo == nil {
		t.Fatal("secret or wrap info is nil")
	}
	wrapInfo = secret.WrapInfo

	// Check for correct CreationPath before rewrap
	if wrapInfo.CreationPath != "secret/foo" {
		t.Fatalf("error on wrapInfo.CreationPath: expected: secret/foo, got: %s", wrapInfo.CreationPath)
	}

	// Test rewrapping
	secret, err = client.Logical().Write("sys/wrapping/rewrap", map[string]interface{}{
		"token": wrapInfo.Token,
	})
	if err != nil {
		t.Fatal(err)
	}
	if secret.Warnings != nil {
		t.Fatalf("Warnings found: %v", secret.Warnings)
	}

	// Check for correct Creation path after rewrap
	if wrapInfo.CreationPath != "secret/foo" {
		t.Fatalf("error on wrapInfo.CreationPath: expected: secret/foo, got: %s", wrapInfo.CreationPath)
	}

	// Should be expired and fail
	_, err = client.Logical().Write("sys/wrapping/unwrap", map[string]interface{}{
		"token": wrapInfo.Token,
	})
	if err == nil {
		t.Fatal("expected err")
	}

	// Attempt unwrapping the rewrapped token
	wrapToken := secret.WrapInfo.Token
	secret, err = client.Logical().Unwrap(wrapToken)
	if err != nil {
		t.Fatal(err)
	}
	// Should be expired and fail
	_, err = client.Logical().Unwrap(wrapToken)
	if err == nil {
		t.Fatal("expected err")
	}

	if !reflect.DeepEqual(secret.Data, map[string]interface{}{
		"zip": "zap",
	}) {
		t.Fatalf("secret data did not match expected: %#v", secret.Data)
	}

	// Ensure that wrapping lookup without a client token responds correctly
	client.ClearToken()
	secret, err = client.Logical().Read("sys/wrapping/lookup")
	if secret != nil {
		t.Fatalf("expected no response: %#v", secret)
	}

	if err == nil {
		t.Fatal("expected error")
	}

	var respError *api.ResponseError
	if errors.As(err, &respError); respError.StatusCode != 403 {
		t.Fatalf("expected 403 response, actual: %d", respError.StatusCode)
	}
}
