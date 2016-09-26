package approle

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestAppRole_RoleLogin(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	createRole(t, b, storage, "role1", "a,b,c")
	roleRoleIDReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/role1/role-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(roleRoleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	roleID := resp.Data["role_id"]

	roleSecretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	secretID := resp.Data["secret_id"]

	loginData := map[string]interface{}{
		"role_id":   roleID,
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
