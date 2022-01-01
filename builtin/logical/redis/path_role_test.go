package redis

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_Role(t *testing.T) {
	b, ctx, s, _, stop := getBackendAndSetConfig(t)
	defer stop()

	data := map[string]interface{}{
		"rules": []string{"on", "allkeys", "+set"},
	}

	req := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/test",
		Storage:   s,
		Data:      data,
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data, resp.Data) {
		t.Fatalf("Expected: %#v\nActual: %#v", data, resp.Data)
	}

	req.Operation = logical.ReadOperation
	req.Data = nil
	resp, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data, resp.Data) {
		t.Fatalf("Expected: %#v\nActual: %#v", data, resp.Data)
	}

	req.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal(resp)
	}

	req.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("Expected error, got: %#v", resp)
	}
	if resp.Error().Error() != "No role found" {
		t.Fatalf("Wrong error: %s", resp.Error())
	}
}
