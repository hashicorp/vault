package approle

import (
	"context"
	"testing"
	"time"

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
	resp, err = b.HandleRequest(context.Background(), roleRoleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	roleID := resp.Data["role_id"]

	roleSecretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
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
	loginResp, err := b.HandleRequest(context.Background(), loginReq)
	if err != nil || (loginResp != nil && loginResp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, loginResp)
	}

	if loginResp.Auth == nil {
		t.Fatalf("expected a non-nil auth object in the response")
	}

	// Test renewal
	renewReq := generateRenewRequest(storage, loginResp.Auth)

	resp, err = b.HandleRequest(context.Background(), renewReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Auth.TTL != 400*time.Second {
		t.Fatalf("expected period value from response to be 400s, got: %s", resp.Auth.TTL)
	}

	// Test renewal with period
	loginResp.Auth.Period = 600 * time.Second
	renewReq = generateRenewRequest(storage, loginResp.Auth)

	resp, err = b.HandleRequest(context.Background(), renewReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Auth.Period != 600*time.Second {
		t.Fatalf("expected period value from response to be 600s, got: %s", resp.Auth.Period)
	}
}

func generateRenewRequest(s logical.Storage, auth *logical.Auth) *logical.Request {
	renewReq := &logical.Request{
		Operation: logical.RenewOperation,
		Storage:   s,
		Auth:      &logical.Auth{},
	}
	renewReq.Auth.InternalData = auth.InternalData
	renewReq.Auth.Metadata = auth.Metadata
	renewReq.Auth.LeaseOptions = auth.LeaseOptions
	renewReq.Auth.Policies = auth.Policies
	renewReq.Auth.IssueTime = time.Now()
	renewReq.Auth.Period = auth.Period

	return renewReq
}
