package openpgp

import (
	"context"
	"github.com/hashicorp/vault/logical"
	"testing"
)

func TestPGP_ExportNotExistingKeyReturnsNotFound(t *testing.T) {
	storage := &logical.InmemStorage{}

	b := Backend()

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.ReadOperation,
		Path:      "export/test",
	}
	rsp, err := b.HandleRequest(context.Background(), req)

	if !(rsp == nil && err == nil) {
		t.Fatal("Key does not exist but does not return not found")
	}
}

func TestPGP_ExportNotExportableKeyReturnsNotFound(t *testing.T) {
	storage := &logical.InmemStorage{}

	b := Backend()

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/test",
		Data: map[string]interface{}{
			"real_name": "Vault PGP test",
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	reqExp := &logical.Request{
		Storage:   storage,
		Operation: logical.ReadOperation,
		Path:      "export/test",
	}
	resp, err := b.HandleRequest(context.Background(), reqExp)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatal("Key not exportable but was exported")
	}
}

func TestPGP_ExportKey(t *testing.T) {
	storage := &logical.InmemStorage{}

	b := Backend()

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/test",
		Data: map[string]interface{}{
			"real_name":  "Vault PGP test",
			"exportable": true,
		},
	}
	_, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	reqExp := &logical.Request{
		Storage:   storage,
		Operation: logical.ReadOperation,
		Path:      "export/test",
	}
	resp, err := b.HandleRequest(context.Background(), reqExp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("not expected error response: %#v", *resp)
	}
	name, ok := resp.Data["name"]
	if !ok {
		t.Fatalf("no name key found in response data %#v", resp.Data)
	}
	if name != "test" {
		t.Fatalf("not expected name, expected test got: %s", name)
	}
}
