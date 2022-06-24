package cert

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestCert_RoleResolve(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	roleName := "roleName"

	loginData := map[string]interface{}{
		"name": roleName,
	}
	loginReq := &logical.Request{
		Operation: logical.ResolveRoleOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
		Connection: &logical.Connection{
			RemoteAddr: "127.0.0.1",
		},
	}

	resp, err := b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["role"] != roleName {
		t.Fatalf("Role was not as expected. Expected %s, received %s", roleName, resp.Data["role"])
	}
}
