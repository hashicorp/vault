// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"os"
)

func TestBackend_PathRoles(t *testing.T) {

	// Skip tests if we are not running acceptance tests
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}

	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend()
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"role":            "devrole",
		"description":     "My dev role",
		"add_ocid_list":   "ocid1,ocid2",
		"add_policy_list": "policy1,policy2",
		"ttl":             1500,
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   config.StorageView,
		Data:      roleData,
	}

	numRoles := 10
	baseRolePath := "role/devrole"

	//first create the roles
	for i := 1; i <= numRoles; i++ {
		roleReq.Path = baseRolePath + strconv.Itoa(i)
		resp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("Role creation failed. resp:%#v\n err:%v", resp, err)
		}
	}

	//now read the roles
	for i := 1; i <= numRoles; i++ {
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.ReadOperation,
			Path:      baseRolePath + strconv.Itoa(i),
			Storage:   config.StorageView,
		})

		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("Read roles failed. resp:%#v\n err:%v", resp, err)
		}
	}

	//now update the roles
	roleDataUpdate := map[string]interface{}{
		"role":               "devrole",
		"description":        "My developer role",
		"add_ocid_list":      "ocid3",
		"remove_ocid_list":   "ocid1",
		"add_policy_list":    "policy1,policy3",
		"remove_policy_list": "policy4,policy3",
		"ttl":                1000,
	}

	roleReqUpdate := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Data:      roleDataUpdate,
	}
	for i := 1; i <= numRoles; i++ {
		roleReqUpdate.Path = baseRolePath + strconv.Itoa(i)
		resp, err = b.HandleRequest(context.Background(), roleReqUpdate)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("Role update failed. resp:%#v\n err:%v", resp, err)
		}
	}

	//now read the roles again
	for i := 1; i <= numRoles; i++ {
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.ReadOperation,
			Path:      baseRolePath + strconv.Itoa(i),
			Storage:   config.StorageView,
		})

		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("Read roles failed. resp:%#v\n err:%v", resp, err)
		}
	}

	//now list the roles
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "role/",
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("Listing roles failed. resp:%#v\n err:%v", resp, err)
	}

	if len(resp.Data["keys"].([]string)) != numRoles {
		t.Fatalf("Failed to list all the roles")
	}

	//now delete half the roles
	for i := 1; i <= 5; i++ {
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.DeleteOperation,
			Path:      baseRolePath + strconv.Itoa(i),
			Storage:   config.StorageView,
		})

		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("Read roles failed. resp:%#v\n err:%v", resp, err)
		}
	}

	//now list the roles again
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "role/",
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("Listing roles failed. resp:%#v\n err:%v", resp, err)
	}

	roleCount := len(resp.Data["keys"].([]string))
	if roleCount != 5 {
		t.Fatalf("Failed to list the expected number of roles")
	}
}
