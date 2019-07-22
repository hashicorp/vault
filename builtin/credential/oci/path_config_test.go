// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"testing"

	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"os"
)

func TestBackend_PathConfig(t *testing.T) {

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
	configPath := BASE_CONFIG_PATH + HOME_TENANCY_ID_CONFIG_NAME

	//First create the config
	configData := map[string]interface{}{
		"configName":  HOME_TENANCY_ID_CONFIG_NAME,
		"configValue": "ocid1.tenancy.oc1..dummy",
	}

	configReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   config.StorageView,
		Data:      configData,
	}

	configReq.Path = configPath
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("Config creation failed. resp:%#v\n err:%v", resp, err)
	}

	//Now try to create a different config (should fail)
	invalidConfigData := map[string]interface{}{
		"configName":  "mydummyconfig",
		"configValue": "ocid1.tenancy.oc1..dummy",
	}

	invalidConfigReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   config.StorageView,
		Data:      invalidConfigData,
	}

	invalidConfigReq.Path = "config/mydummyconfig"
	resp, err = b.HandleRequest(context.Background(), invalidConfigReq)
	if err != nil {
		t.Fatalf("Config creation failed. resp:%#v\n err:%v", resp, err)
	}
	if resp == nil || resp.IsError() == false {
		t.Fatalf("Config creation succeded while it should'nt have.")
	}

	//now read the config
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      configReq.Path,
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("Read config failed. resp:%#v\n err:%v", resp, err)
	}

	//now try to update the config (should pass)
	configUpdate := map[string]interface{}{
		"configName":  HOME_TENANCY_ID_CONFIG_NAME,
		"configValue": "ocid1.tenancy.oc1..dummy",
	}

	configReqUpdate := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Data:      configUpdate,
	}

	configReqUpdate.Path = configPath
	resp, err = b.HandleRequest(context.Background(), configReqUpdate)
	if err != nil {
		t.Fatalf("bad: config update failed. resp:%#v\n err:%v", resp, err)
	}

	if resp != nil && resp.IsError() == true {
		t.Fatalf("Config update succeded while it should'nt have.")
	}

	//now list the configs
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "config/",
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("listing configs failed. resp:%#v\n err:%v", resp, err)
	}
	if len(resp.Data["keys"].([]string)) != 1 {
		t.Fatalf("failed to list all configs")
	}

	//now try to delete the config (should fail)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      configPath,
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("Config delete failed. resp:%#v\n err:%v", resp, err)
	}
	fmt.Println("All tests completed successfully")
}
