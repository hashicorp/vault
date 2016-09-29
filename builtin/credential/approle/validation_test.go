package approle

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestAppRole_SecretIDNumUsesUpgrade(t *testing.T) {
	var resp *logical.Response
	var err error

	b, storage := createBackendWithStorage(t)

	roleData := map[string]interface{}{
		"secret_id_num_uses": 10,
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

	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}

	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	secretIDReq.Operation = logical.UpdateOperation
	secretIDReq.Path = "role/role1/secret-id/lookup"
	secretIDReq.Data = map[string]interface{}{
		"secret_id": resp.Data["secret_id"].(string),
	}
	resp, err = b.HandleRequest(secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Check if the response contains the value set for secret_id_num_uses
	// and not SecretIDNumUses
	if resp.Data["secret_id_num_uses"] != 10 ||
		resp.Data["SecretIDNumUses"] != 0 {
		t.Fatal("invalid secret_id_num_uses")
	}
}
