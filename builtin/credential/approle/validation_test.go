package approle

import (
	"context"
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

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	secretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}

	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	secretIDReq.Operation = logical.UpdateOperation
	secretIDReq.Path = "role/role1/secret-id/lookup"
	secretIDReq.Data = map[string]interface{}{
		"secret_id": resp.Data["secret_id"].(string),
	}
	resp, err = b.HandleRequest(context.Background(), secretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Check if the response contains the value set for secret_id_num_uses
	if resp.Data["secret_id_num_uses"] != 10 {
		t.Fatal("invalid secret_id_num_uses")
	}
}
