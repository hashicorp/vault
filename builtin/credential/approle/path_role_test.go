package approle

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
)

func TestAppRole_RoleReadSetIndex(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	roleReq := &logical.Request{
		Path:      "role/testrole",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"bind_secret_id": true,
		},
	}

	// Create a role
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}

	roleIDReq := &logical.Request{
		Path:      "role/testrole/role-id",
		Operation: logical.ReadOperation,
		Storage:   storage,
	}

	// Get the role ID
	resp, err = b.HandleRequest(roleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}
	roleID := resp.Data["role_id"].(string)

	// Delete the role ID index
	err = b.roleIDEntryDelete(storage, roleID)
	if err != nil {
		t.Fatal(err)
	}

	// Read the role again. This should add the index and return a warning
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}

	// Check if the warning is being returned
	if !strings.Contains(resp.Warnings[0], "Role identifier was missing an index back to role name.") {
		t.Fatalf("bad: expected a warning in the response")
	}

	roleIDIndex, err := b.roleIDEntry(storage, roleID)
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
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v\n", resp, err)
	}
}

func TestAppRole_CIDRSubset(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

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

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}

	secretIDData := map[string]interface{}{
		"cidr_list": "127.0.0.1/16",
	}
	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/testrole1/secret-id",
		Data:      secretIDData,
	}

	resp, err = b.HandleRequest(secretIDReq)
	if resp != nil || resp.IsError() {
		t.Fatalf("resp:%#v", resp)
	}
	if err == nil {
		t.Fatal("expected an error")
	}

	roleData["bound_cidr_list"] = "192.168.27.29/16,172.245.30.40/24,10.20.30.40/30"
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}

	secretIDData["cidr_list"] = "192.168.27.29/20,172.245.30.40/25,10.20.30.40/32"
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("resp: %#v", resp)
	}
}

func TestAppRole_RoleConstraints(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Set bound_cidr_list alone by explicitly disabling bind_secret_id
	roleReq.Operation = logical.UpdateOperation
	roleData["bind_secret_id"] = false
	roleData["bound_cidr_list"] = "0.0.0.0/0"
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Remove both constraints
	roleReq.Operation = logical.UpdateOperation
	roleData["bound_cidr_list"] = ""
	roleData["bind_secret_id"] = false
	resp, err = b.HandleRequest(roleReq)
	if resp != nil && resp.IsError() {
		t.Fatalf("err:%v, resp:%#v", err, resp)
	}
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestAppRole_RoleIDUpdate(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleIDUpdateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/testrole1/role-id",
		Storage:   storage,
		Data: map[string]interface{}{
			"role_id": "customroleid",
		},
	}
	resp, err = b.HandleRequest(roleIDUpdateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/testrole1/secret-id",
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
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
	resp, err = b.HandleRequest(loginReq)
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

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Path = "role/testrole2"
	resp, err = b.HandleRequest(roleReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleData["role_id"] = "role-id-456"
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.UpdateOperation
	roleData["role_id"] = "role-id-123"
	resp, err = b.HandleRequest(roleReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleReq.Path = "role/testrole1"
	roleData["role_id"] = "role-id-456"
	resp, err = b.HandleRequest(roleReq)
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
	resp, err = b.HandleRequest(roleIDReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleIDData["role_id"] = "role-id-123"
	roleIDReq.Path = "role/testrole2/role-id"
	resp, err = b.HandleRequest(roleIDReq)
	if err == nil && !(resp != nil && resp.IsError()) {
		t.Fatalf("expected an error: got resp:%#v", resp)
	}

	roleIDData["role_id"] = "role-id-2000"
	resp, err = b.HandleRequest(roleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleIDData["role_id"] = "role-id-1000"
	roleIDReq.Path = "role/testrole1/role-id"
	resp, err = b.HandleRequest(roleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}

func TestAppRole_RoleDeleteSecretID(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	createRole(t, b, storage, "role1", "a,b")
	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	// Create 3 secrets on the role
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	secretIDAccessors := resp.Data["keys"].([]string)
	if len(secretIDAccessors) != 3 {
		t.Fatalf("bad: len of secretIDAccessors: expected:3 actual:%d", len(secretIDAccessors))
	}

	roleReq := &logical.Request{
		Operation: logical.DeleteOperation,
		Storage:   storage,
		Path:      "role/role1",
	}
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	resp, err = b.HandleRequest(listReq)
	if err != nil || resp == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("expected an error. err:%v resp:%#v", err, resp)
	}
}

func TestAppRole_RoleSecretIDReadDelete(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	createRole(t, b, storage, "role1", "a,b")
	secretIDCreateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(secretIDCreateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

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
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
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
	resp, err = b.HandleRequest(deleteSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	resp, err = b.HandleRequest(secretIDReq)
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

	createRole(t, b, storage, "role1", "a,b")
	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	resp, err = b.HandleRequest(listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	hmacSecretID := resp.Data["keys"].([]string)[0]

	hmacReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id-accessor/lookup",
		Data: map[string]interface{}{
			"secret_id_accessor": hmacSecretID,
		},
	}
	resp, err = b.HandleRequest(hmacReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp.Data == nil {
		t.Fatal(err)
	}

	hmacReq.Path = "role/role1/secret-id-accessor/destroy"
	resp, err = b.HandleRequest(hmacReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	hmacReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(hmacReq)
	if resp != nil && resp.IsError() {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestAppRoleRoleListSecretID(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	createRole(t, b, storage, "role1", "a,b")

	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id",
	}
	// Create 5 'secret_id's
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Storage:   storage,
		Path:      "role/role1/secret-id/",
	}
	resp, err = b.HandleRequest(listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	secrets := resp.Data["keys"].([]string)
	if len(secrets) != 5 {
		t.Fatalf("bad: len of secrets: expected:5 actual:%d", len(secrets))
	}
}

func TestAppRole_RoleList(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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
	resp, err = b.HandleRequest(listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	actual := resp.Data["keys"].([]string)
	expected := []string{"role1", "role2", "role3", "role4", "role5"}
	if !policyutil.EquivalentPolicies(actual, expected) {
		t.Fatalf("bad: listed roles: expected:%s\nactual:%s", expected, actual)
	}
}

func TestAppRole_RoleSecretID(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleSecretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["secret_id"].(string) == "" {
		t.Fatalf("failed to generate secret_id")
	}

	roleSecretIDReq.Path = "role/role1/custom-secret-id"
	roleCustomSecretIDData := map[string]interface{}{
		"secret_id": "abcd123",
	}
	roleSecretIDReq.Data = roleCustomSecretIDData
	roleSecretIDReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["secret_id"] != "abcd123" {
		t.Fatalf("failed to set specific secret_id to role")
	}
}

func TestAppRole_RoleCRUD(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	roleData := map[string]interface{}{
		"policies":           "p,q,r,s",
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
		"token_num_uses":     600,
		"bound_cidr_list":    "127.0.0.1/32,127.0.0.1/16",
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/role1",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	expected := map[string]interface{}{
		"bind_secret_id":     true,
		"policies":           []string{"p", "q", "r", "s"},
		"secret_id_num_uses": 10,
		"secret_id_ttl":      300,
		"token_ttl":          400,
		"token_max_ttl":      500,
		"token_num_uses":     600,
		"bound_cidr_list":    "127.0.0.1/32,127.0.0.1/16",
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

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

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
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp.Data["role_id"].(string) != "test_role_id" {
		t.Fatalf("bad: role_id: expected:test_role_id actual:%s\n", resp.Data["role_id"].(string))
	}

	roleReq.Data = map[string]interface{}{"role_id": "custom_role_id"}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp.Data["role_id"].(string) != "custom_role_id" {
		t.Fatalf("bad: role_id: expected:custom_role_id actual:%s\n", resp.Data["role_id"].(string))
	}

	// RUD for bind_secret_id field
	roleReq.Path = "role/role1/bind-secret-id"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Data = map[string]interface{}{"bind_secret_id": false}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["bind_secret_id"].(bool) {
		t.Fatalf("bad: bind_secret_id: expected:false actual:%t\n", resp.Data["bind_secret_id"].(bool))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if !resp.Data["bind_secret_id"].(bool) {
		t.Fatalf("expected the default value of 'true' to be set")
	}

	// RUD for policies field
	roleReq.Path = "role/role1/policies"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Data = map[string]interface{}{"policies": "a1,b1,c1,d1"}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if !reflect.DeepEqual(resp.Data["policies"].([]string), []string{"a1", "b1", "c1", "d1"}) {
		t.Fatalf("bad: policies: actual:%s\n", resp.Data["policies"].([]string))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	expectedPolicies := []string{"default"}
	actualPolicies := resp.Data["policies"].([]string)
	if !policyutil.EquivalentPolicies(expectedPolicies, actualPolicies) {
		t.Fatalf("bad: policies: expected:%s actual:%s", expectedPolicies, actualPolicies)
	}

	// RUD for secret-id-num-uses field
	roleReq.Path = "role/role1/secret-id-num-uses"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Data = map[string]interface{}{"secret_id_num_uses": 200}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["secret_id_num_uses"].(int) != 200 {
		t.Fatalf("bad: secret_id_num_uses: expected:200 actual:%d\n", resp.Data["secret_id_num_uses"].(int))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["secret_id_num_uses"].(int) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for secret_id_ttl field
	roleReq.Path = "role/role1/secret-id-ttl"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Data = map[string]interface{}{"secret_id_ttl": 3001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["secret_id_ttl"].(time.Duration) != 3001 {
		t.Fatalf("bad: secret_id_ttl: expected:3001 actual:%d\n", resp.Data["secret_id_ttl"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["secret_id_ttl"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for secret-id-num-uses field
	roleReq.Path = "role/role1/token-num-uses"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp.Data["token_num_uses"].(int) != 600 {
		t.Fatalf("bad: token_num_uses: expected:600 actual:%d\n", resp.Data["token_num_uses"].(int))
	}

	roleReq.Data = map[string]interface{}{"token_num_uses": 60}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["token_num_uses"].(int) != 60 {
		t.Fatalf("bad: token_num_uses: expected:60 actual:%d\n", resp.Data["token_num_uses"].(int))
	}

	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["token_num_uses"].(int) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for 'period' field
	roleReq.Path = "role/role1/period"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Data = map[string]interface{}{"period": 9001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["period"].(time.Duration) != 9001 {
		t.Fatalf("bad: period: expected:9001 actual:%d\n", resp.Data["9001"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["period"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for token_ttl field
	roleReq.Path = "role/role1/token-ttl"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Data = map[string]interface{}{"token_ttl": 4001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["token_ttl"].(time.Duration) != 4001 {
		t.Fatalf("bad: token_ttl: expected:4001 actual:%d\n", resp.Data["token_ttl"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["token_ttl"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// RUD for token_max_ttl field
	roleReq.Path = "role/role1/token-max-ttl"
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Data = map[string]interface{}{"token_max_ttl": 5001}
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["token_max_ttl"].(time.Duration) != 5001 {
		t.Fatalf("bad: token_max_ttl: expected:5001 actual:%d\n", resp.Data["token_max_ttl"].(time.Duration))
	}
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["token_max_ttl"].(time.Duration) != 0 {
		t.Fatalf("expected value to be reset")
	}

	// Delete test for role
	roleReq.Path = "role/role1"
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
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

	resp, err := b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}
