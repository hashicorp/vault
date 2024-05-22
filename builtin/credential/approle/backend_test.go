// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package approle

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func createBackendWithStorage(t *testing.T) (*backend, logical.Storage) {
	t.Helper()
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil {
		t.Fatalf("failed to create backend")
	}
	err = b.Backend.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	return b, config.StorageView
}

func TestAppRole_RoleServiceToBatchNumUses(t *testing.T) {
	b, s := createBackendWithStorage(t)

	requestFunc := func(operation logical.Operation, data map[string]interface{}) {
		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/testrole",
			Operation: operation,
			Storage:   s,
			Data:      data,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %#v\nresp: %#v", err, resp)
		}
	}

	data := map[string]interface{}{
		"bind_secret_id":     true,
		"secret_id_num_uses": 0,
		"secret_id_ttl":      "10m",
		"token_policies":     "policy",
		"token_ttl":          "5m",
		"token_max_ttl":      "10m",
		"token_num_uses":     2,
		"token_type":         "default",
	}
	requestFunc(logical.CreateOperation, data)

	data["token_num_uses"] = 0
	data["token_type"] = "batch"
	requestFunc(logical.UpdateOperation, data)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole/role-id",
		Operation: logical.ReadOperation,
		Storage:   s,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	roleID := resp.Data["role_id"]

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole/secret-id",
		Operation: logical.UpdateOperation,
		Storage:   s,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	secretID := resp.Data["secret_id"]

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
		Storage: s,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	require.NotNil(t, resp.Auth)
}

func TestAppRole_RoleNameCaseSensitivity(t *testing.T) {
	testFunc := func(t *testing.T, roleName string) {
		var resp *logical.Response
		var err error
		b, s := createBackendWithStorage(t)

		// Create the role
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName,
			Operation: logical.CreateOperation,
			Storage:   s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
		}

		// Get the role-id
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/role-id",
			Operation: logical.ReadOperation,
			Storage:   s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		roleID := resp.Data["role_id"]

		// Create a secret-id
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/secret-id",
			Operation: logical.UpdateOperation,
			Storage:   s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		secretID := resp.Data["secret_id"]
		secretIDAccessor := resp.Data["secret_id_accessor"]

		// Ensure login works
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		if resp.Auth == nil {
			t.Fatalf("failed to perform login")
		}

		// Destroy secret ID accessor
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/secret-id-accessor/destroy",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"secret_id_accessor": secretIDAccessor,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}

		// Login again using the accessor's corresponding secret ID should fail
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil && err != logical.ErrInvalidCredentials {
			t.Fatal(err)
		}
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error due to invalid secret ID")
		}

		// Generate another secret ID
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/secret-id",
			Operation: logical.UpdateOperation,
			Storage:   s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		secretID = resp.Data["secret_id"]
		secretIDAccessor = resp.Data["secret_id_accessor"]

		// Ensure login works
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		if resp.Auth == nil {
			t.Fatalf("failed to perform login")
		}

		// Destroy the secret ID
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/secret-id/destroy",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"secret_id": secretID,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}

		// Login again using the same secret ID should fail
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil && err != logical.ErrInvalidCredentials {
			t.Fatal(err)
		}
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error due to invalid secret ID")
		}

		// Generate another secret ID
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/secret-id",
			Operation: logical.UpdateOperation,
			Storage:   s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		secretID = resp.Data["secret_id"]
		secretIDAccessor = resp.Data["secret_id_accessor"]

		// Ensure login works
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		if resp.Auth == nil {
			t.Fatalf("failed to perform login")
		}

		// Destroy the secret ID using lower cased role name
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + strings.ToLower(roleName) + "/secret-id/destroy",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"secret_id": secretID,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}

		// Login again using the same secret ID should fail
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil && err != logical.ErrInvalidCredentials {
			t.Fatal(err)
		}
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error due to invalid secret ID")
		}

		// Generate another secret ID
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/secret-id",
			Operation: logical.UpdateOperation,
			Storage:   s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		secretID = resp.Data["secret_id"]
		secretIDAccessor = resp.Data["secret_id_accessor"]

		// Ensure login works
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		if resp.Auth == nil {
			t.Fatalf("failed to perform login")
		}

		// Destroy the secret ID using upper cased role name
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + strings.ToUpper(roleName) + "/secret-id/destroy",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"secret_id": secretID,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}

		// Login again using the same secret ID should fail
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil && err != logical.ErrInvalidCredentials {
			t.Fatal(err)
		}
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error due to invalid secret ID")
		}

		// Generate another secret ID
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/" + roleName + "/secret-id",
			Operation: logical.UpdateOperation,
			Storage:   s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		secretID = resp.Data["secret_id"]
		secretIDAccessor = resp.Data["secret_id_accessor"]

		// Ensure login works
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		if resp.Auth == nil {
			t.Fatalf("failed to perform login")
		}

		// Destroy the secret ID using mixed case name
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/saMpleRolEnaMe/secret-id/destroy",
			Operation: logical.UpdateOperation,
			Storage:   s,
			Data: map[string]interface{}{
				"secret_id": secretID,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}

		// Login again using the same secret ID should fail
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "login",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			},
			Storage: s,
		})
		if err != nil && err != logical.ErrInvalidCredentials {
			t.Fatal(err)
		}
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error due to invalid secret ID")
		}
	}

	// Lower case role name
	testFunc(t, "samplerolename")
	// Upper case role name
	testFunc(t, "SAMPLEROLENAME")
	// Mixed case role name
	testFunc(t, "SampleRoleName")
}
