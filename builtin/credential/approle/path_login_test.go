package approle

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestAppRole_BoundCIDRLogin(t *testing.T) {
	var resp *logical.Response
	var err error
	b, s := createBackendWithStorage(t)

	// Create a role with secret ID binding disabled and only bound cidr list
	// enabled
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"bind_secret_id":    false,
			"bound_cidr_list":   []string{"127.0.0.1/8"},
			"token_bound_cidrs": []string{"10.0.0.0/8"},
		},
		Storage: s,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Read the role ID
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole/role-id",
		Operation: logical.ReadOperation,
		Storage:   s,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleID := resp.Data["role_id"]

	// Fill in the connection information and login with just the role ID
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"role_id": roleID,
		},
		Storage:    s,
		Connection: &logical.Connection{RemoteAddr: "127.0.0.1"},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp.Auth == nil {
		t.Fatal("expected login to succeed")
	}
	if len(resp.Auth.BoundCIDRs) != 1 {
		t.Fatal("bad token bound cidrs")
	}
	if resp.Auth.BoundCIDRs[0].String() != "10.0.0.0/8" {
		t.Fatalf("bad: %s", resp.Auth.BoundCIDRs[0].String())
	}

	// Override with a secret-id value, verify it doesn't pass
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "role/testrole",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"bind_secret_id": true,
		},
		Storage: s,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleSecretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/testrole/secret-id",
		Storage:   s,
		Data: map[string]interface{}{
			"token_bound_cidrs": []string{"11.0.0.0/24"},
		},
	}
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
	if err == nil {
		t.Fatal("expected error due to mismatching subnet relationship")
	}
	roleSecretIDReq.Data["token_bound_cidrs"] = "10.0.0.0/24"
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	secretID := resp.Data["secret_id"]

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
		Storage:    s,
		Connection: &logical.Connection{RemoteAddr: "127.0.0.1"},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp.Auth == nil {
		t.Fatal("expected login to succeed")
	}
	if len(resp.Auth.BoundCIDRs) != 1 {
		t.Fatal("bad token bound cidrs")
	}
	if resp.Auth.BoundCIDRs[0].String() != "10.0.0.0/24" {
		t.Fatalf("bad: %s", resp.Auth.BoundCIDRs[0].String())
	}
}

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

	///
	// Test renewal with period
	///

	// Create role
	period := 600 * time.Second
	roleData := map[string]interface{}{
		"policies": "a,b,c",
		"period":   period.String(),
	}
	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/" + "role-period",
		Storage:   storage,
		Data:      roleData,
	}
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	roleRoleIDReq = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/role-period/role-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), roleRoleIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	roleID = resp.Data["role_id"]

	roleSecretIDReq = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role-period/secret-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	secretID = resp.Data["secret_id"]

	loginData["role_id"] = roleID
	loginData["secret_id"] = secretID

	loginResp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (loginResp != nil && loginResp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, loginResp)
	}

	if loginResp.Auth == nil {
		t.Fatalf("expected a non-nil auth object in the response")
	}

	renewReq = generateRenewRequest(storage, loginResp.Auth)

	resp, err = b.HandleRequest(context.Background(), renewReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Auth.Period != period {
		t.Fatalf("expected period value of %d in the response, got: %s", period, resp.Auth.Period)
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
	renewReq.Auth.Period = auth.Period

	return renewReq
}
