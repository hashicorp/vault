package database

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_CA(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}
	defer b.Cleanup(context.Background())

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "ca",
		Storage:   config.StorageView,
	}
	resp, err := b.HandleRequest(namespace.RootContext(context.TODO()), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
	serial := resp.Data["serial_number"].(string)

	resp, err = b.HandleRequest(namespace.RootContext(context.TODO()), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	if serial != resp.Data["serial_number"].(string) {
		t.Fatalf("ca has not been cached: %q, %q", serial, resp.Data["serial_number"].(string))
	}

	req.Path = "ca/rotate"
	req.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
	if serial == resp.Data["serial_number"].(string) {
		t.Fatalf("ca has not been changed: %q, %q", serial, resp.Data["serial_number"].(string))
	}
}
