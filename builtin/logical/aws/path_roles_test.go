package aws

import (
	"strconv"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_PathListRoles(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend()
	if err := b.Setup(config); err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"arn": "testarn",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Data:      roleData,
	}

	for i := 1; i <= 10; i++ {
		roleReq.Path = "roles/testrole" + strconv.Itoa(i)
		resp, err = b.HandleRequest(roleReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: role creation failed. resp:%#v\n err:%v", resp, err)
		}
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles",
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: listing roles failed. resp:%#v\n err:%v", resp, err)
	}

	if len(resp.Data["keys"].([]string)) != 10 {
		t.Fatalf("failed to list all 10 roles")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles/",
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: listing roles failed. resp:%#v\n err:%v", resp, err)
	}

	if len(resp.Data["keys"].([]string)) != 10 {
		t.Fatalf("failed to list all 10 roles")
	}
}
