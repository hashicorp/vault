package approle

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

func (b *backend) requestNoErr(t *testing.T, req *logical.Request) *logical.Response {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route(req.Path), req.Operation), resp, true)
	return resp
}

func TestAppRole_LocalSecretIDsRead(t *testing.T) {
	b, storage := createBackendWithStorage(t)

	roleData := map[string]interface{}{
		"local_secret_ids": true,
		"bind_secret_id":   true,
	}

	b.requestNoErr(t, &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/testrole",
		Storage:   storage,
		Data:      roleData,
	})

	resp := b.requestNoErr(t, &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   storage,
		Path:      "role/testrole/local-secret-ids",
	})

	if !resp.Data["local_secret_ids"].(bool) {
		t.Fatalf("expected local_secret_ids to be returned")
	}
}

func TestAppRole_LocalNonLocalSecretIDs(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	// Create a role with local_secret_ids set
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole1",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"policies":         []string{"default", "role1policy"},
			"bind_secret_id":   true,
			"local_secret_ids": true,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\n resp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	// Create another role without setting local_secret_ids
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole2",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"policies":       []string{"default", "role1policy"},
			"bind_secret_id": true,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\n resp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	count := 10
	// Create secret IDs on testrole1
	for i := 0; i < count; i++ {
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/testrole1/secret-id",
			Operation: logical.UpdateOperation,
			Storage:   storage,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)
	}

	// Check the number of secret IDs generated
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole1/secret-id",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if len(resp.Data["keys"].([]string)) != count {
		t.Fatalf("failed to list secret IDs")
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.ListOperation), resp, true)

	// Create secret IDs on testrole1
	for i := 0; i < count; i++ {
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Path:      "role/testrole2/secret-id",
			Operation: logical.UpdateOperation,
			Storage:   storage,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole2/secret-id",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if len(resp.Data["keys"].([]string)) != count {
		t.Fatalf("failed to list secret IDs")
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.ListOperation), resp, true)
}

func TestAppRole_UpgradeSecretIDPrefix(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	// Create a role entry directly in storage without SecretIDPrefix
	err = b.setRoleEntry(context.Background(), storage, "testrole", &roleStorageEntry{
		RoleID:           "testroleid",
		HMACKey:          "testhmackey",
		Policies:         []string{"default"},
		BindSecretID:     true,
		BoundCIDRListOld: "127.0.0.1/18,192.178.1.2/24",
	}, "")
	if err != nil {
		t.Fatal(err)
	}

	// Reading the role entry should upgrade it to contain SecretIDPrefix
	role, err := b.roleEntry(context.Background(), storage, "testrole")
	if err != nil {
		t.Fatal(err)
	}
	if role.SecretIDPrefix == "" {
		t.Fatalf("expected SecretIDPrefix to be set")
	}

	// Ensure that the API response contains local_secret_ids
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\n resp: %#v", err, resp)
	}
	_, ok := resp.Data["local_secret_ids"]
	if !ok {
		t.Fatalf("expected local_secret_ids to be present in the response")
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)
}

func TestAppRole_LocalSecretIDImmutability(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"policies":         []string{"default"},
		"bind_secret_id":   true,
		"bound_cidr_list":  []string{"127.0.0.1/18", "192.178.1.2/24"},
		"local_secret_ids": true,
	}

	// Create a role with local_secret_ids set
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data:      roleData,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	// Attempt to modify local_secret_ids should fail
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data:      roleData,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error since local_secret_ids can't be overwritten")
	}
}

func TestAppRole_UpgradeBoundCIDRList(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"policies":        []string{"default"},
		"bind_secret_id":  true,
		"bound_cidr_list": []string{"127.0.0.1/18", "192.178.1.2/24"},
	}

	// Create a role with bound_cidr_list set
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data:      roleData,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	// Read the role and check that the bound_cidr_list is set properly
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	expected := []string{"127.0.0.1/18", "192.178.1.2/24"}
	actual := resp.Data["secret_id_bound_cidrs"].([]string)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: secret_id_bound_cidrs; expected: %#v\nactual: %#v\n", expected, actual)
	}

	// Modify the storage entry of the role to hold the old style string typed bound_cidr_list
	role := &roleStorageEntry{
		RoleID:           "testroleid",
		HMACKey:          "testhmackey",
		Policies:         []string{"default"},
		BindSecretID:     true,
		BoundCIDRListOld: "127.0.0.1/18,192.178.1.2/24",
		SecretIDPrefix:   secretIDPrefix,
	}
	err = b.setRoleEntry(context.Background(), storage, "testrole", role, "")
	if err != nil {
		t.Fatal(err)
	}

	// Read the role. The upgrade code should have migrated the old type to the new type
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: bound_cidr_list; expected: %#v\nactual: %#v\n", expected, actual)
	}

	// Create a secret-id by supplying a subset of the role's CIDR blocks with the new type
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole/secret-id",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"cidr_list": []string{"127.0.0.1/24"},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	if resp.Data["secret_id"].(string) == "" {
		t.Fatalf("failed to generate secret-id")
	}

	// Check that the backwards compatibility for the string type is not broken
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole/secret-id",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"cidr_list": "127.0.0.1/24",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	if resp.Data["secret_id"].(string) == "" {
		t.Fatalf("failed to generate secret-id")
	}
}

func TestAppRole_RoleNameLowerCasing(t *testing.T) {
	var resp *logical.Response
	var err error
	var roleID, secretID string

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	// Save a role with out LowerCaseRoleName set
	role := &roleStorageEntry{
		RoleID:         "testroleid",
		HMACKey:        "testhmackey",
		Policies:       []string{"default"},
		BindSecretID:   true,
		SecretIDPrefix: secretIDPrefix,
	}
	err = b.setRoleEntry(context.Background(), storage, "testRoleName", role, "")
	if err != nil {
		t.Fatal(err)
	}

	secretIDReq := &logical.Request{
		Path:      "role/testRoleName/secret-id",
		Operation: logical.UpdateOperation,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	secretID = resp.Data["secret_id"].(string)
	roleID = "testroleid"

	// Regular login flow. This should succeed.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Lower case the role name when generating the secret id
	secretIDReq.Path = "role/testrolename/secret-id"
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	secretID = resp.Data["secret_id"].(string)

	// Login should fail
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
	})
	if err != nil && err != logical.ErrInvalidCredentials {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error")
	}

	// Delete the role and create it again. This time don't directly persist
	// it, but route the request to the creation handler so that it sets the
	// LowerCaseRoleName to true.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testRoleName",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.DeleteOperation), resp, true)

	roleReq := &logical.Request{
		Path:      "role/testRoleName",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"bind_secret_id": true,
		},
	}
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	// Create secret id with lower cased role name
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrolename/secret-id",
		Operation: logical.UpdateOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	secretID = resp.Data["secret_id"].(string)

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrolename/role-id",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.ReadOperation), resp, true)

	roleID = resp.Data["role_id"].(string)

	// Login should pass
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}

	// Lookup of secret ID should work in case-insensitive manner
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrolename/secret-id/lookup",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"secret_id": secretID,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil {
		t.Fatalf("failed to lookup secret IDs")
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 16, logical.UpdateOperation), resp, true)

	// Listing of secret IDs should work in case-insensitive manner
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrolename/secret-id",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.ListOperation), resp, true)

	if len(resp.Data["keys"].([]string)) != 1 {
		t.Fatalf("failed to list secret IDs")
	}
}

func TestAppRole_RoleReadSetIndex(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleReq := &logical.Request{
		Path:      "role/testrole",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"bind_secret_id": true,
		},
	}

	// Create a role
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	roleIDReq := &logical.Request{
		Path:      "role/testrole/role-id",
		Operation: logical.ReadOperation,
		Storage:   storage,
	}

	// Get the role ID
	resp, err = b.HandleRequest(context.Background(), roleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.ReadOperation), resp, true)

	roleID := resp.Data["role_id"].(string)

	// Delete the role ID index
	err = b.roleIDEntryDelete(context.Background(), storage, roleID)
	if err != nil {
		t.Fatal(err)
	}

	// Read the role again. This should add the index and return a warning
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}

	// Check if the warning is being returned
	if !strings.Contains(resp.Warnings[0], "Role identifier was missing an index back to role name.") {
		t.Fatalf("bad: expected a warning in the response")
	}

	roleIDIndex, err := b.roleIDEntry(context.Background(), storage, roleID)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the index has been successfully created
	if roleIDIndex == nil || roleIDIndex.Name != "testrole" {
		t.Fatalf("bad: expected role to have an index")
	}

	roleReq.Operation = logical.UpdateOperation
	roleReq.Data = map[string]interface{}{
		"bind_secret_id": true,
		"policies":       "default",
	}

	// Check if updating and reading of roles work and that there are no lock
	// contentions dangling due to previous operation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)
}

func TestAppRole_CIDRSubset(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"role_id":         "role-id-123",
		"policies":        "a,b",
		"bound_cidr_list": "127.0.0.1/24",
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/testrole1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	secretIDData := map[string]interface{}{
		"cidr_list": "127.0.0.1/16",
	}
	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/testrole1/secret-id",
		Data:      secretIDData,
	}

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if resp != nil {
		t.Fatalf("resp:%#v", resp)
	}
	if err == nil {
		t.Fatal("expected an error")
	}

	roleData["bound_cidr_list"] = "192.168.27.29/16,172.245.30.40/24,10.20.30.40/30"
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)

	secretIDData["cidr_list"] = "192.168.27.29/20,172.245.30.40/25,10.20.30.40/32"
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("resp: %#v", resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)
}

func TestAppRole_TokenBoundCIDRSubset32Mask(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"role_id":           "role-id-123",
		"policies":          "a,b",
		"token_bound_cidrs": "127.0.0.1/32",
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/testrole1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	secretIDData := map[string]interface{}{
		"token_bound_cidrs": "127.0.0.1/32",
	}
	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/testrole1/secret-id",
		Data:      secretIDData,
	}

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	secretIDData = map[string]interface{}{
		"token_bound_cidrs": "127.0.0.1/24",
	}
	secretIDReq = &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/testrole1/secret-id",
		Data:      secretIDData,
	}

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if resp != nil {
		t.Fatalf("resp:%#v", resp)
	}

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestAppRole_RoleConstraints(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"role_id":  "role-id-123",
		"policies": "a,b",
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/testrole1",
		Storage:   storage,
		Data:      roleData,
	}

	// Set bind_secret_id, which is enabled by default
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	// Set bound_cidr_list alone by explicitly disabling bind_secret_id
	roleReq.Operation = logical.UpdateOperation
	roleData["bind_secret_id"] = false
	roleData["bound_cidr_list"] = "0.0.0.0/0"
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)

	// Remove both constraints
	roleReq.Operation = logical.UpdateOperation
	roleData["bound_cidr_list"] = ""
	roleData["bind_secret_id"] = false
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if resp != nil && resp.IsError() {
		t.Fatalf("err:%v, resp:%#v", err, resp)
	}
	if err == nil {
		t.Fatalf("expected an error")
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)
}

func TestAppRole_RoleIDUpdate(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"role_id":            "role-id-123",
		"policies":           "a,b",
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/testrole1",
		Storage:   storage,
		Data:      roleData,
	}
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	roleIDUpdateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/testrole1/role-id",
		Storage:   storage,
		Data: map[string]interface{}{
			"role_id": "customroleid",
		},
	}
	resp, err = b.HandleRequest(context.Background(), roleIDUpdateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.UpdateOperation), resp, true)

	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/testrole1/secret-id",
	}
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	secretID := resp.Data["secret_id"].(string)

	loginData := map[string]interface{}{
		"role_id":   "customroleid",
		"secret_id": secretID,
	}
	loginReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
		Connection: &logical.Connection{
			RemoteAddr: "127.0.0.1",
		},
	}
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Auth == nil {
		t.Fatalf("expected a non-nil auth object in the response")
	}
}

func TestAppRole_RoleIDUniqueness(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"role_id":            "role-id-123",
		"policies":           "a,b",
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/testrole1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	roleReq.Path = "role/testrole2"
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleData["role_id"] = "role-id-456"
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	roleReq.Operation = logical.UpdateOperation
	roleData["role_id"] = "role-id-123"
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)

	roleReq.Path = "role/testrole1"
	roleData["role_id"] = "role-id-456"
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleIDData := map[string]interface{}{
		"role_id": "role-id-456",
	}
	roleIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/testrole1/role-id",
		Storage:   storage,
		Data:      roleIDData,
	}
	resp, err = b.HandleRequest(context.Background(), roleIDReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleIDData["role_id"] = "role-id-123"
	roleIDReq.Path = "role/testrole2/role-id"
	resp, err = b.HandleRequest(context.Background(), roleIDReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleIDData["role_id"] = "role-id-2000"
	resp, err = b.HandleRequest(context.Background(), roleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.UpdateOperation), resp, true)

	roleIDData["role_id"] = "role-id-1000"
	roleIDReq.Path = "role/testrole1/role-id"
	resp, err = b.HandleRequest(context.Background(), roleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.UpdateOperation), resp, true)
}

func TestAppRole_RoleDeleteSecretID(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	createRole(t, b, storage, "role1", "a,b")
	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	// Create 3 secrets on the role
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(context.Background(), listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.ListOperation), resp, true)

	secretIDAccessors := resp.Data["keys"].([]string)
	if len(secretIDAccessors) != 3 {
		t.Fatalf("bad: len of secretIDAccessors: expected:3 actual:%d", len(secretIDAccessors))
	}

	roleReq := &logical.Request{
		Operation: logical.DeleteOperation,
		Storage:   storage,
		Path:      "role/role1",
	}
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.DeleteOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), listReq)
	if err != nil || resp == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("expected an error. err:%v resp:%#v", err, resp)
	}
}

func TestAppRole_RoleSecretIDReadDelete(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	createRole(t, b, storage, "role1", "a,b")
	secretIDCreateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(context.Background(), secretIDCreateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	secretID := resp.Data["secret_id"].(string)
	if secretID == "" {
		t.Fatal("expected non empty secret ID")
	}

	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id/lookup",
		Data: map[string]interface{}{
			"secret_id": secretID,
		},
	}
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 16, logical.UpdateOperation), resp, true)

	if resp.Data == nil {
		t.Fatal(err)
	}

	deleteSecretIDReq := &logical.Request{
		Operation: logical.DeleteOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id/destroy",
		Data: map[string]interface{}{
			"secret_id": secretID,
		},
	}
	resp, err = b.HandleRequest(context.Background(), deleteSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 17, logical.DeleteOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if resp != nil && resp.IsError() {
		t.Fatalf("error response:%#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestAppRole_RoleSecretIDAccessorReadDelete(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	createRole(t, b, storage, "role1", "a,b")
	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(context.Background(), listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.ListOperation), resp, true)

	hmacSecretID := resp.Data["keys"].([]string)[0]

	hmacReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id-accessor/lookup",
		Data: map[string]interface{}{
			"secret_id_accessor": hmacSecretID,
		},
	}
	resp, err = b.HandleRequest(context.Background(), hmacReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp.Data == nil {
		t.Fatal(err)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 18, logical.UpdateOperation), resp, true)

	hmacReq.Path = "role/role1/secret-id-accessor/destroy"
	resp, err = b.HandleRequest(context.Background(), hmacReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 19, logical.UpdateOperation), resp, true)

	hmacReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), hmacReq)
	if resp != nil && resp.IsError() {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestAppRoleSecretIDLookup(t *testing.T) {
	b, storage := createBackendWithStorage(t)
	createRole(t, b, storage, "role1", "a,b")

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id-accessor/lookup",
		Data: map[string]interface{}{
			"secret_id_accessor": "invalid",
		},
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := &logical.Response{
		Data: map[string]interface{}{
			"http_content_type": "application/json",
			"http_raw_body":     `{"request_id":"","lease_id":"","renewable":false,"lease_duration":0,"data":{"error":"failed to find accessor entry for secret_id_accessor: \"invalid\""},"wrap_info":null,"warnings":null,"auth":null}`,
			"http_status_code":  404,
		},
	}
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("resp:%#v expected:%#v", resp, expected)
	}
}

func TestAppRoleRoleListSecretID(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	createRole(t, b, storage, "role1", "a,b")

	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	// Create 5 'secret_id's
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id/",
	}
	resp, err = b.HandleRequest(context.Background(), listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.ListOperation), resp, true)

	secrets := resp.Data["keys"].([]string)
	if len(secrets) != 5 {
		t.Fatalf("bad: len of secrets: expected:5 actual:%d", len(secrets))
	}
}

func TestAppRole_RoleList(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	createRole(t, b, storage, "role1", "a,b")
	createRole(t, b, storage, "role2", "c,d")
	createRole(t, b, storage, "role3", "e,f")
	createRole(t, b, storage, "role4", "g,h")
	createRole(t, b, storage, "role5", "i,j")

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "role",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 1, logical.ListOperation), resp, true)

	actual := resp.Data["keys"].([]string)
	expected := []string{"role1", "role2", "role3", "role4", "role5"}
	if !policyutil.EquivalentPolicies(actual, expected) {
		t.Fatalf("bad: listed roles: expected:%s\nactual:%s", expected, actual)
	}
}

func TestAppRole_RoleSecretIDWithoutFields(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"policies":           "p,q,r,s",
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/role1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	roleSecretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

	if resp.Data["secret_id"].(string) == "" {
		t.Fatalf("failed to generate secret_id")
	}
	if resp.Data["secret_id_ttl"].(int64) != int64(roleData["secret_id_ttl"].(int)) {
		t.Fatalf("secret_id_ttl has not defaulted to the role's secret id ttl")
	}
	if resp.Data["secret_id_num_uses"].(int) != roleData["secret_id_num_uses"].(int) {
		t.Fatalf("secret_id_num_uses has not defaulted to the role's secret id num_uses")
	}

	roleSecretIDReq.Path = "role/role1/custom-secret-id"
	roleCustomSecretIDData := map[string]interface{}{
		"secret_id": "abcd123",
	}
	roleSecretIDReq.Data = roleCustomSecretIDData
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 20, logical.UpdateOperation), resp, true)

	if resp.Data["secret_id"] != "abcd123" {
		t.Fatalf("failed to set specific secret_id to role")
	}
	if resp.Data["secret_id_ttl"].(int64) != int64(roleData["secret_id_ttl"].(int)) {
		t.Fatalf("secret_id_ttl has not defaulted to the role's secret id ttl")
	}
	if resp.Data["secret_id_num_uses"].(int) != roleData["secret_id_num_uses"].(int) {
		t.Fatalf("secret_id_num_uses has not defaulted to the role's secret id num_uses")
	}
}

func TestAppRole_RoleSecretIDWithValidFields(t *testing.T) {
	type testCase struct {
		name    string
		payload map[string]interface{}
	}

	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"policies":           "p,q,r,s",
		"secret_id_num_uses": 0,
		"secret_id_ttl":      0,
		"token_ttl":          400,
		"token_max_ttl":      500,
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/role1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	testCases := []testCase{
		{
			name:    "finite num_uses ttl",
			payload: map[string]interface{}{"secret_id": "finite", "ttl": 5, "num_uses": 5},
		},
		{
			name:    "infinite num_uses and ttl",
			payload: map[string]interface{}{"secret_id": "infinite", "ttl": 0, "num_uses": 0},
		},
		{
			name:    "finite num_uses and infinite ttl",
			payload: map[string]interface{}{"secret_id": "mixed1", "ttl": 0, "num_uses": 5},
		},
		{
			name:    "infinite num_uses and finite ttl",
			payload: map[string]interface{}{"secret_id": "mixed2", "ttl": 5, "num_uses": 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			roleSecretIDReq := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "role/role1/secret-id",
				Storage:   storage,
			}
			roleCustomSecretIDData := tc.payload
			roleSecretIDReq.Data = roleCustomSecretIDData

			resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%v resp:%#v", err, resp)
			}
			schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

			if resp.Data["secret_id"].(string) == "" {
				t.Fatalf("failed to generate secret_id")
			}
			if resp.Data["secret_id_ttl"].(int64) != int64(tc.payload["ttl"].(int)) {
				t.Fatalf("secret_id_ttl has not been set by the 'ttl' field")
			}
			if resp.Data["secret_id_num_uses"].(int) != tc.payload["num_uses"].(int) {
				t.Fatalf("secret_id_num_uses has not been set by the 'num_uses' field")
			}

			roleSecretIDReq.Path = "role/role1/custom-secret-id"
			roleSecretIDReq.Data = roleCustomSecretIDData
			resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%v resp:%#v", err, resp)
			}
			schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 20, logical.UpdateOperation), resp, true)

			if resp.Data["secret_id"] != tc.payload["secret_id"] {
				t.Fatalf("failed to set specific secret_id to role")
			}
			if resp.Data["secret_id_ttl"].(int64) != int64(tc.payload["ttl"].(int)) {
				t.Fatalf("secret_id_ttl has not been set by the 'ttl' field")
			}
			if resp.Data["secret_id_num_uses"].(int) != tc.payload["num_uses"].(int) {
				t.Fatalf("secret_id_num_uses has not been set by the 'num_uses' field")
			}
		})
	}
}

func TestAppRole_ErrorsRoleSecretIDWithInvalidFields(t *testing.T) {
	type testCase struct {
		name     string
		payload  map[string]interface{}
		expected string
	}

	type roleTestCase struct {
		name    string
		options map[string]interface{}
		cases   []testCase
	}

	infiniteTestCases := []testCase{
		{
			name:     "infinite ttl",
			payload:  map[string]interface{}{"secret_id": "abcd123", "num_uses": 1, "ttl": 0},
			expected: "ttl cannot be longer than the role's secret_id_ttl",
		},
		{
			name:     "infinite num_uses",
			payload:  map[string]interface{}{"secret_id": "abcd123", "num_uses": 0, "ttl": 1},
			expected: "num_uses cannot be higher than the role's secret_id_num_uses",
		},
	}

	negativeTestCases := []testCase{
		{
			name:     "negative num_uses",
			payload:  map[string]interface{}{"secret_id": "abcd123", "num_uses": -1, "ttl": 0},
			expected: "num_uses cannot be negative",
		},
	}

	roleTestCases := []roleTestCase{
		{
			name: "infinite role secret id ttl",
			options: map[string]interface{}{
				"secret_id_num_uses": 1,
				"secret_id_ttl":      0,
			},
			cases: []testCase{
				{
					name:     "higher num_uses",
					payload:  map[string]interface{}{"secret_id": "abcd123", "ttl": 0, "num_uses": 2},
					expected: "num_uses cannot be higher than the role's secret_id_num_uses",
				},
			},
		},
		{
			name: "infinite role num_uses",
			options: map[string]interface{}{
				"secret_id_num_uses": 0,
				"secret_id_ttl":      1,
			},
			cases: []testCase{
				{
					name:     "longer ttl",
					payload:  map[string]interface{}{"secret_id": "abcd123", "ttl": 2, "num_uses": 0},
					expected: "ttl cannot be longer than the role's secret_id_ttl",
				},
			},
		},
		{
			name: "finite role ttl and num_uses",
			options: map[string]interface{}{
				"secret_id_num_uses": 2,
				"secret_id_ttl":      2,
			},
			cases: infiniteTestCases,
		},
		{
			name: "mixed role ttl and num_uses",
			options: map[string]interface{}{
				"secret_id_num_uses": 400,
				"secret_id_ttl":      500,
			},
			cases: negativeTestCases,
		},
	}

	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	for i, rc := range roleTestCases {
		roleData := map[string]interface{}{
			"policies":      "p,q,r,s",
			"token_ttl":     400,
			"token_max_ttl": 500,
		}
		roleData["secret_id_num_uses"] = rc.options["secret_id_num_uses"]
		roleData["secret_id_ttl"] = rc.options["secret_id_ttl"]

		roleReq := &logical.Request{
			Operation: logical.CreateOperation,
			Path:      fmt.Sprintf("role/role%d", i),
			Storage:   storage,
			Data:      roleData,
		}

		resp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
		schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

		for _, tc := range rc.cases {
			t.Run(fmt.Sprintf("%s/%s", rc.name, tc.name), func(t *testing.T) {
				roleSecretIDReq := &logical.Request{
					Operation: logical.UpdateOperation,
					Path:      fmt.Sprintf("role/role%d/secret-id", i),
					Storage:   storage,
				}
				roleSecretIDReq.Data = tc.payload
				resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
				if err != nil || (resp != nil && !resp.IsError()) {
					t.Fatalf("err:%v resp:%#v", err, resp)
				}
				if resp.Data["error"].(string) != tc.expected {
					t.Fatalf("expected: %q, got: %q", tc.expected, resp.Data["error"].(string))
				}

				roleSecretIDReq.Path = fmt.Sprintf("role/role%d/custom-secret-id", i)
				resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
				if err != nil || (resp != nil && !resp.IsError()) {
					t.Fatalf("err:%v resp:%#v", err, resp)
				}
				if resp.Data["error"].(string) != tc.expected {
					t.Fatalf("expected: %q, got: %q", tc.expected, resp.Data["error"].(string))
				}
			})
		}
	}
}

func TestAppRole_RoleCRUD(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"policies":              "p,q,r,s",
		"secret_id_num_uses":    10,
		"secret_id_ttl":         300,
		"token_ttl":             400,
		"token_max_ttl":         500,
		"token_num_uses":        600,
		"secret_id_bound_cidrs": "127.0.0.1/32,127.0.0.1/16",
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/role1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	expected := map[string]interface{}{
		"bind_secret_id":        true,
		"policies":              []string{"p", "q", "r", "s"},
		"secret_id_num_uses":    10,
		"secret_id_ttl":         300,
		"token_ttl":             400,
		"token_max_ttl":         500,
		"token_num_uses":        600,
		"secret_id_bound_cidrs": []string{"127.0.0.1/32", "127.0.0.1/16"},
		"token_bound_cidrs":     []string{},
		"token_type":            "default",
	}

	var expectedStruct roleStorageEntry
	err = mapstructure.Decode(expected, &expectedStruct)
	if err != nil {
		t.Fatal(err)
	}

	var actualStruct roleStorageEntry
	err = mapstructure.Decode(resp.Data, &actualStruct)
	if err != nil {
		t.Fatal(err)
	}

	expectedStruct.RoleID = actualStruct.RoleID
	if diff := deep.Equal(expectedStruct, actualStruct); diff != nil {
		t.Fatal(diff)
	}

	roleData = map[string]interface{}{
		"role_id":            "test_role_id",
		"policies":           "a,b,c,d",
		"secret_id_num_uses": 100,
		"secret_id_ttl":      3000,
		"token_ttl":          4000,
		"token_max_ttl":      5000,
	}
	roleReq.Data = roleData
	roleReq.Operation = logical.UpdateOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	expected = map[string]interface{}{
		"policies":           []string{"a", "b", "c", "d"},
		"secret_id_num_uses": 100,
		"secret_id_ttl":      3000,
		"token_ttl":          4000,
		"token_max_ttl":      5000,
	}
	err = mapstructure.Decode(expected, &expectedStruct)
	if err != nil {
		t.Fatal(err)
	}

	err = mapstructure.Decode(resp.Data, &actualStruct)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedStruct, actualStruct) {
		t.Fatalf("bad:\nexpected:%#v\nactual:%#v\n", expectedStruct, actualStruct)
	}

	// RU for role_id field
	roleReq.Path = "role/role1/role-id"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.ReadOperation), resp, true)

	if resp.Data["role_id"].(string) != "test_role_id" {
		t.Fatalf("bad: role_id: expected:test_role_id actual:%s\n", resp.Data["role_id"].(string))
	}

	roleReq.Data = map[string]interface{}{"role_id": "custom_role_id"}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 14, logical.ReadOperation), resp, true)

	if resp.Data["role_id"].(string) != "custom_role_id" {
		t.Fatalf("bad: role_id: expected:custom_role_id actual:%s\n", resp.Data["role_id"].(string))
	}

	// RUD for bind_secret_id field
	roleReq.Path = "role/role1/bind-secret-id"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 7, logical.ReadOperation), resp, true)

	roleReq.Data = map[string]interface{}{"bind_secret_id": false}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 7, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 7, logical.ReadOperation), resp, true)

	if resp.Data["bind_secret_id"].(bool) {
		t.Fatalf("bad: bind_secret_id: expected:false actual:%t\n", resp.Data["bind_secret_id"].(bool))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 7, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 7, logical.ReadOperation), resp, true)

	if !resp.Data["bind_secret_id"].(bool) {
		t.Fatalf("expected the default value of 'true' to be set")
	}

	// RUD for policies field
	roleReq.Path = "role/role1/policies"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 3, logical.ReadOperation), resp, true)

	roleReq.Data = map[string]interface{}{"policies": "a1,b1,c1,d1"}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 3, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 3, logical.ReadOperation), resp, true)

	if !reflect.DeepEqual(resp.Data["policies"].([]string), []string{"a1", "b1", "c1", "d1"}) {
		t.Fatalf("bad: policies: actual:%s\n", resp.Data["policies"].([]string))
	}
	if !reflect.DeepEqual(resp.Data["token_policies"].([]string), []string{"a1", "b1", "c1", "d1"}) {
		t.Fatalf("bad: policies: actual:%s\n", resp.Data["policies"].([]string))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 3, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 3, logical.ReadOperation), resp, true)

	expectedPolicies := []string{}
	actualPolicies := resp.Data["token_policies"].([]string)
	if !policyutil.EquivalentPolicies(expectedPolicies, actualPolicies) {
		t.Fatalf("bad: token_policies: expected:%s actual:%s", expectedPolicies, actualPolicies)
	}

	// RUD for secret-id-num-uses field
	roleReq.Path = "role/role1/secret-id-num-uses"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 8, logical.ReadOperation), resp, true)

	roleReq.Data = map[string]interface{}{"secret_id_num_uses": 200}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 8, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 8, logical.ReadOperation), resp, true)

	if resp.Data["secret_id_num_uses"].(int) != 200 {
		t.Fatalf("bad: secret_id_num_uses: expected:200 actual:%d\n", resp.Data["secret_id_num_uses"].(int))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 8, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 8, logical.ReadOperation), resp, true)

	if resp.Data["secret_id_num_uses"].(int) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for secret_id_ttl field
	roleReq.Path = "role/role1/secret-id-ttl"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 9, logical.ReadOperation), resp, true)

	roleReq.Data = map[string]interface{}{"secret_id_ttl": 3001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 9, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 9, logical.ReadOperation), resp, true)

	if resp.Data["secret_id_ttl"].(time.Duration) != 3001 {
		t.Fatalf("bad: secret_id_ttl: expected:3001 actual:%d\n", resp.Data["secret_id_ttl"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 9, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 9, logical.ReadOperation), resp, true)

	if resp.Data["secret_id_ttl"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for secret-id-num-uses field
	roleReq.Path = "role/role1/token-num-uses"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 11, logical.ReadOperation), resp, true)

	if resp.Data["token_num_uses"].(int) != 600 {
		t.Fatalf("bad: token_num_uses: expected:600 actual:%d\n", resp.Data["token_num_uses"].(int))
	}

	roleReq.Data = map[string]interface{}{"token_num_uses": 60}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 11, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 11, logical.ReadOperation), resp, true)

	if resp.Data["token_num_uses"].(int) != 60 {
		t.Fatalf("bad: token_num_uses: expected:60 actual:%d\n", resp.Data["token_num_uses"].(int))
	}

	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 11, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 11, logical.ReadOperation), resp, true)

	if resp.Data["token_num_uses"].(int) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for 'period' field
	roleReq.Path = "role/role1/period"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 10, logical.ReadOperation), resp, true)

	roleReq.Data = map[string]interface{}{"period": 9001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 10, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 10, logical.ReadOperation), resp, true)

	if resp.Data["period"].(time.Duration) != 9001 {
		t.Fatalf("bad: period: expected:9001 actual:%d\n", resp.Data["9001"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 10, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 10, logical.ReadOperation), resp, true)

	if resp.Data["token_period"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for token_ttl field
	roleReq.Path = "role/role1/token-ttl"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 12, logical.ReadOperation), resp, true)

	roleReq.Data = map[string]interface{}{"token_ttl": 4001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 12, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 12, logical.ReadOperation), resp, true)

	if resp.Data["token_ttl"].(time.Duration) != 4001 {
		t.Fatalf("bad: token_ttl: expected:4001 actual:%d\n", resp.Data["token_ttl"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 12, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 12, logical.ReadOperation), resp, true)

	if resp.Data["token_ttl"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for token_max_ttl field
	roleReq.Path = "role/role1/token-max-ttl"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 13, logical.ReadOperation), resp, true)

	roleReq.Data = map[string]interface{}{"token_max_ttl": 5001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 13, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 13, logical.ReadOperation), resp, true)

	if resp.Data["token_max_ttl"].(time.Duration) != 5001 {
		t.Fatalf("bad: token_max_ttl: expected:5001 actual:%d\n", resp.Data["token_max_ttl"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 13, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 13, logical.ReadOperation), resp, true)

	if resp.Data["token_max_ttl"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// Delete test for role
	roleReq.Path = "role/role1"
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func TestAppRole_RoleWithTokenBoundCIDRsCRUD(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"policies":              "p,q,r,s",
		"secret_id_num_uses":    10,
		"secret_id_ttl":         300,
		"token_ttl":             400,
		"token_max_ttl":         500,
		"token_num_uses":        600,
		"secret_id_bound_cidrs": "127.0.0.1/32,127.0.0.1/16",
		"token_bound_cidrs":     "127.0.0.1/32,127.0.0.1/16",
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/role1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	expected := map[string]interface{}{
		"bind_secret_id":        true,
		"policies":              []string{"p", "q", "r", "s"},
		"secret_id_num_uses":    10,
		"secret_id_ttl":         300,
		"token_ttl":             400,
		"token_max_ttl":         500,
		"token_num_uses":        600,
		"token_bound_cidrs":     []string{"127.0.0.1/32", "127.0.0.1/16"},
		"secret_id_bound_cidrs": []string{"127.0.0.1/32", "127.0.0.1/16"},
		"token_type":            "default",
	}

	var expectedStruct roleStorageEntry
	err = mapstructure.Decode(expected, &expectedStruct)
	if err != nil {
		t.Fatal(err)
	}

	var actualStruct roleStorageEntry
	err = mapstructure.Decode(resp.Data, &actualStruct)
	if err != nil {
		t.Fatal(err)
	}

	expectedStruct.RoleID = actualStruct.RoleID
	if !reflect.DeepEqual(expectedStruct, actualStruct) {
		t.Fatalf("bad:\nexpected:%#v\nactual:%#v\n", expectedStruct, actualStruct)
	}

	roleData = map[string]interface{}{
		"role_id":            "test_role_id",
		"policies":           "a,b,c,d",
		"secret_id_num_uses": 100,
		"secret_id_ttl":      3000,
		"token_ttl":          4000,
		"token_max_ttl":      5000,
	}
	roleReq.Data = roleData
	roleReq.Operation = logical.UpdateOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	expected = map[string]interface{}{
		"policies":           []string{"a", "b", "c", "d"},
		"secret_id_num_uses": 100,
		"secret_id_ttl":      3000,
		"token_ttl":          4000,
		"token_max_ttl":      5000,
	}
	err = mapstructure.Decode(expected, &expectedStruct)
	if err != nil {
		t.Fatal(err)
	}

	err = mapstructure.Decode(resp.Data, &actualStruct)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedStruct, actualStruct) {
		t.Fatalf("bad:\nexpected:%#v\nactual:%#v\n", expectedStruct, actualStruct)
	}

	// RUD for secret-id-bound-cidrs field
	roleReq.Path = "role/role1/secret-id-bound-cidrs"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 5, logical.ReadOperation), resp, true)

	if resp.Data["secret_id_bound_cidrs"].([]string)[0] != "127.0.0.1/32" ||
		resp.Data["secret_id_bound_cidrs"].([]string)[1] != "127.0.0.1/16" {
		t.Fatalf("bad: secret_id_bound_cidrs: expected:127.0.0.1/32,127.0.0.1/16 actual:%d\n", resp.Data["secret_id_bound_cidrs"].(int))
	}

	roleReq.Data = map[string]interface{}{"secret_id_bound_cidrs": []string{"127.0.0.1/20"}}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 5, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 5, logical.ReadOperation), resp, true)

	if resp.Data["secret_id_bound_cidrs"].([]string)[0] != "127.0.0.1/20" {
		t.Fatalf("bad: secret_id_bound_cidrs: expected:127.0.0.1/20 actual:%s\n", resp.Data["secret_id_bound_cidrs"].([]string)[0])
	}

	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 5, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 5, logical.ReadOperation), resp, true)

	if len(resp.Data["secret_id_bound_cidrs"].([]string)) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for token-bound-cidrs field
	roleReq.Path = "role/role1/token-bound-cidrs"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 6, logical.ReadOperation), resp, true)

	if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1" ||
		resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[1].String() != "127.0.0.1/16" {
		m, err := json.Marshal(resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler))
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("bad: token_bound_cidrs: expected:127.0.0.1/32,127.0.0.1/16 actual:%s\n", string(m))
	}

	roleReq.Data = map[string]interface{}{"token_bound_cidrs": []string{"127.0.0.1/20"}}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 6, logical.UpdateOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 6, logical.ReadOperation), resp, true)

	if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1/20" {
		t.Fatalf("bad: token_bound_cidrs: expected:127.0.0.1/20 actual:%s\n", resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0])
	}

	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 6, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 6, logical.ReadOperation), resp, true)

	if len(resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// Delete test for role
	roleReq.Path = "role/role1"
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func TestAppRole_RoleWithTokenTypeCRUD(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	roleData := map[string]interface{}{
		"policies":           "p,q,r,s",
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
		"token_num_uses":     600,
		"token_type":         "default-service",
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/role1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

	if 0 == len(resp.Warnings) {
		t.Fatalf("bad:\nexpected warning in resp:%#v\n", resp.Warnings)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	expected := map[string]interface{}{
		"bind_secret_id":     true,
		"policies":           []string{"p", "q", "r", "s"},
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
		"token_num_uses":     600,
		"token_type":         "service",
	}

	var expectedStruct roleStorageEntry
	err = mapstructure.Decode(expected, &expectedStruct)
	if err != nil {
		t.Fatal(err)
	}

	var actualStruct roleStorageEntry
	err = mapstructure.Decode(resp.Data, &actualStruct)
	if err != nil {
		t.Fatal(err)
	}

	expectedStruct.RoleID = actualStruct.RoleID
	if !reflect.DeepEqual(expectedStruct, actualStruct) {
		t.Fatalf("bad:\nexpected:%#v\nactual:%#v\n", expectedStruct, actualStruct)
	}

	roleData = map[string]interface{}{
		"role_id":            "test_role_id",
		"policies":           "a,b,c,d",
		"secret_id_num_uses": 100,
		"secret_id_ttl":      3000,
		"token_ttl":          4000,
		"token_max_ttl":      5000,
		"token_type":         "default-service",
	}
	roleReq.Data = roleData
	roleReq.Operation = logical.UpdateOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation), resp, true)

	if 0 == len(resp.Warnings) {
		t.Fatalf("bad:\nexpected a warning in resp:%#v\n", resp.Warnings)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.ReadOperation), resp, true)

	expected = map[string]interface{}{
		"policies":           []string{"a", "b", "c", "d"},
		"secret_id_num_uses": 100,
		"secret_id_ttl":      3000,
		"token_ttl":          4000,
		"token_max_ttl":      5000,
		"token_type":         "service",
	}
	err = mapstructure.Decode(expected, &expectedStruct)
	if err != nil {
		t.Fatal(err)
	}

	err = mapstructure.Decode(resp.Data, &actualStruct)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedStruct, actualStruct) {
		t.Fatalf("bad:\nexpected:%#v\nactual:%#v\n", expectedStruct, actualStruct)
	}

	// Delete test for role
	roleReq.Path = "role/role1"
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.DeleteOperation), resp, true)

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func createRole(t *testing.T, b *backend, s logical.Storage, roleName, policies string) {
	roleData := map[string]interface{}{
		"policies":           policies,
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/" + roleName,
		Storage:   s,
		Data:      roleData,
	}

	paths := rolePaths(b)

	resp, err := b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)
}

// TestAppRole_TokenutilUpgrade ensures that when we read values out that are
// values with upgrade logic we see the correct struct entries populated
func TestAppRole_TokenutilUpgrade(t *testing.T) {
	tests := []struct {
		name              string
		storageValMissing bool
		storageVal        string
		expectedTokenType logical.TokenType
	}{
		{
			"token_type_missing",
			true,
			"",
			logical.TokenTypeDefault,
		},
		{
			"token_type_empty",
			false,
			"",
			logical.TokenTypeDefault,
		},
		{
			"token_type_service",
			false,
			"service",
			logical.TokenTypeService,
		},
	}

	s := &logical.InmemStorage{}

	config := logical.TestBackendConfig()
	config.StorageView = s

	ctx := context.Background()

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil {
		t.Fatalf("failed to create backend")
	}
	if err := b.Setup(ctx, config); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Construct the storage entry object based on our test case.
			tokenTypeKV := ""
			if !tt.storageValMissing {
				tokenTypeKV = fmt.Sprintf(`, "token_type": "%s"`, tt.storageVal)
			}
			entryVal := fmt.Sprintf(`{"policies": ["foo"], "period": 300000000000, "token_bound_cidrs": ["127.0.0.1", "10.10.10.10/24"]%s}`, tokenTypeKV)

			// Hand craft JSON because there is overlap between fields
			if err := s.Put(ctx, &logical.StorageEntry{
				Key:   "role/" + tt.name,
				Value: []byte(entryVal),
			}); err != nil {
				t.Fatal(err)
			}

			resEntry, err := b.roleEntry(ctx, s, tt.name)
			if err != nil {
				t.Fatal(err)
			}

			exp := &roleStorageEntry{
				SecretIDPrefix: "secret_id/",
				Policies:       []string{"foo"},
				Period:         300 * time.Second,
				TokenParams: tokenutil.TokenParams{
					TokenPolicies: []string{"foo"},
					TokenPeriod:   300 * time.Second,
					TokenBoundCIDRs: []*sockaddr.SockAddrMarshaler{
						{SockAddr: sockaddr.MustIPAddr("127.0.0.1")},
						{SockAddr: sockaddr.MustIPAddr("10.10.10.10/24")},
					},
					TokenType: tt.expectedTokenType,
				},
			}
			if diff := deep.Equal(resEntry, exp); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}

func TestAppRole_SecretID_WithTTL(t *testing.T) {
	tests := []struct {
		name      string
		roleName  string
		ttl       int64
		sysTTLCap bool
	}{
		{
			"zero ttl",
			"role-zero-ttl",
			0,
			false,
		},
		{
			"custom ttl",
			"role-custom-ttl",
			60,
			false,
		},
		{
			"system ttl capped",
			"role-sys-ttl-cap",
			700000000,
			true,
		},
	}

	b, storage := createBackendWithStorage(t)

	paths := rolePaths(b)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create role
			roleData := map[string]interface{}{
				"policies":      "default",
				"secret_id_ttl": tt.ttl,
			}

			roleReq := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "role/" + tt.roleName,
				Storage:   storage,
				Data:      roleData,
			}
			resp, err := b.HandleRequest(context.Background(), roleReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%v resp:%#v", err, resp)
			}
			schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 0, logical.CreateOperation), resp, true)

			// Generate secret ID
			secretIDReq := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "role/" + tt.roleName + "/secret-id",
				Storage:   storage,
			}
			resp, err = b.HandleRequest(context.Background(), secretIDReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%v resp:%#v", err, resp)
			}
			schema.ValidateResponse(t, schema.FindResponseSchema(t, paths, 15, logical.UpdateOperation), resp, true)

			// Extract the "ttl" value from the response data if it exists
			ttlRaw, okTTL := resp.Data["secret_id_ttl"]
			if !okTTL {
				t.Fatalf("expected TTL value in response")
			}

			var (
				respTTL int64
				ok      bool
			)
			respTTL, ok = ttlRaw.(int64)
			if !ok {
				t.Fatalf("expected ttl to be an integer, got: %s", err)
			}

			// Verify secret ID response for different cases
			switch {
			case tt.sysTTLCap:
				if respTTL != int64(b.System().MaxLeaseTTL().Seconds()) {
					t.Fatalf("expected TTL value to be system's max lease TTL, got: %d", respTTL)
				}
			default:
				if respTTL != tt.ttl {
					t.Fatalf("expected TTL value to be %d, got: %d", tt.ttl, respTTL)
				}
			}
		})
	}
}

func TestAppRole_RoleSecretIDAccessorCrossDelete(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	// Create First Role
	createRole(t, b, storage, "role1", "a,b")
	_ = b.requestNoErr(t, &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	})

	// Create Second Role
	createRole(t, b, storage, "role2", "a,b")
	_ = b.requestNoErr(t, &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role2/secret-id",
	})

	// Get role2 secretID Accessor
	resp = b.requestNoErr(t, &logical.Request{
		Operation: logical.ListOperation,
		Storage:   storage,
		Path:      "role/role2/secret-id",
	})

	// Read back role2 secretID Accessor information
	hmacSecretID := resp.Data["keys"].([]string)[0]
	_ = b.requestNoErr(t, &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role2/secret-id-accessor/lookup",
		Data: map[string]interface{}{
			"secret_id_accessor": hmacSecretID,
		},
	})

	// Attempt to destroy role2 secretID accessor using role1 path
	_ = b.requestNoErr(t, &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id-accessor/destroy",
		Data: map[string]interface{}{
			"secret_id_accessor": hmacSecretID,
		},
	})

	// Check to see if role2 secretID accessor was deleted
	accessorHashes, err := storage.List(context.Background(), "accessor/")
	if err != nil {
		t.Fatal(err)
	}
	if len(accessorHashes) != 2 {
		t.Fatalf("unexpected accessor deletion; expected 2 accessors, got %d", len(accessorHashes))
	}
}
