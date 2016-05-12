package mfa

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMFAMethods_CRUD_TOTP(t *testing.T) {
	b, err := MFABackendFactory(logical.TestBackendConfig())
	if err != nil {
		t.Fatal(err)
	}

	storage := &logical.InmemStorage{}
	b.(*MFABackend).SetStorage(storage)

	req := logical.TestRequest(t, logical.ReadOperation, "methods/test")

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("should not see a method")
	}

	req.Operation = logical.UpdateOperation
	checkExists, found, err := b.HandleExistenceCheck(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if !checkExists {
		t.Fatal("no existence check found")
	}
	if found {
		t.Fatal("entry should not have been found")
	}

	//
	// Creation
	//
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"type": "somefin",
	}

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatal("got a nil response, expected an error response")
	}
	if !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", *resp)
	}

	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"type": "totp",
	}

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response, got %#v", *resp)
	}

	checkExists, found, err = b.HandleExistenceCheck(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if !checkExists {
		t.Fatal("no existence check found")
	}
	if !found {
		t.Fatal("entry should have been found")
	}

	//
	// Updating
	//
	req.Operation = logical.UpdateOperation
	req.Data = map[string]interface{}{
		"type":                "totp",
		"totp_hash_algorithm": "sha28182",
	}

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatal("got a nil response, expected an error response")
	}
	if !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", *resp)
	}

	req.Operation = logical.UpdateOperation
	req.Data = map[string]interface{}{
		"type":                "totp",
		"totp_hash_algorithm": "sha256",
	}

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response, got %#v", *resp)
	}

	req.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected := map[string]interface{}{
		"name":                "test",
		"type":                "totp",
		"totp_hash_algorithm": "sha256",
	}

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("expected:\n%v\nactual:\n%v\n", expected, resp.Data)
	}

	//
	// Listing
	//
	req.Path = "methods"
	req.Operation = logical.ListOperation
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expectedList := []string{
		"test",
	}

	if !reflect.DeepEqual(expectedList, resp.Data["keys"].([]string)) {
		t.Fatalf("expected:\n%v\nactual:\n%v\n", expectedList, resp.Data["keys"].([]string))
	}

	//
	// Deletion
	//
	req = logical.TestRequest(t, logical.DeleteOperation, "methods/test")

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("got an error response deleting the test method")
	}

	req.Operation = logical.UpdateOperation
	checkExists, found, err = b.HandleExistenceCheck(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if !checkExists {
		t.Fatal("no existence check found")
	}
	if found {
		t.Fatal("entry should not have been found")
	}
}
