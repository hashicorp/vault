package ssh

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestSSH_ConfigCAUpdateDelete(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	caReq := &logical.Request{
		Path:      "config/ca",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
	}

	// Auto-generate the keys
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(caReq)
	if err == nil {
		t.Fatalf("expected an error")
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = map[string]interface{}{
		"public_key":  publicKey,
		"private_key": privateKey,
	}

	// Successfully create a new one
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(caReq)
	if err == nil {
		t.Fatalf("expected an error")
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = nil

	// Successfully create a new one
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}
}
