package mfa

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMFAIdentifiers_CRD_TOTP(t *testing.T) {
	b, err := MFABackendFactory(logical.TestBackendConfig())
	if err != nil {
		t.Fatal(err)
	}

	storage := &logical.InmemStorage{}
	b.(*MFABackend).SetStorage(storage)

	// Create the TOTP role
	req := logical.TestRequest(t, logical.CreateOperation, "methods/test")
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"type":                "totp",
		"totp_hash_algorithm": "sha256",
	}

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatal("got response, expected a nil response")
	}

	req = logical.TestRequest(t, logical.ReadOperation, "methods/test/jeff@hashicorp.com")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("should not see an identifier")
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
	// Note that this function just does CRUD tests, the actual TOTP validation
	// is in other tests
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{}

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

	req.Data = map[string]interface{}{
		"totp_account_name": "jeff@hashicorp.com",
	}

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
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
	// Reading
	//
	req.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected := map[string]interface{}{
		"totp_account_name": "jeff@hashicorp.com",
		"identifier":        "jeff@hashicorp.com",
	}
	delete(resp.Data, "creation_time_utc")
	delete(resp.Data, "creation_time_string")

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("expected:\n%v\nactual:\n%v\n", expected, resp.Data)
	}

	//
	// Listing
	//
	req.Path = "methods/test/"
	req.Operation = logical.ListOperation
	req.Data = map[string]interface{}{}
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expectedList := []string{
		"jeff@hashicorp.com",
	}

	if !reflect.DeepEqual(expectedList, resp.Data["keys"].([]string)) {
		t.Fatalf("expected:\n%v\nactual:\n%v\n", expectedList, resp.Data["keys"].([]string))
	}

	//
	// Deletion
	//
	req = logical.TestRequest(t, logical.DeleteOperation, "methods/test/jeff@hashicorp.com")

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
