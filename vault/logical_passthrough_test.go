package vault

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestPassthroughBackend_RootPaths(t *testing.T) {
	b := testPassthroughBackend()
	root := b.SpecialPaths()
	if root != nil {
		t.Fatalf("unexpected: %v", root)
	}
}

func TestPassthroughBackend_Write(t *testing.T) {
	b := testPassthroughBackend()
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	out, err := req.Storage.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("failed to write to view")
	}
}

func TestPassthroughBackend_Read_Lease(t *testing.T) {
	b := testPassthroughBackend()
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"
	req.Data["lease"] = "1h"
	storage := req.Storage

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       time.Hour,
			},
		},
		Data: map[string]interface{}{
			"raw":   "test",
			"lease": "1h",
		},
	}

	resp.Secret.InternalData = nil
	resp.Secret.LeaseID = ""
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}
}

func TestPassthroughBackend_Read_TTL(t *testing.T) {
	b := testPassthroughBackend()
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"
	req.Data["ttl"] = "1h"
	storage := req.Storage

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       time.Hour,
			},
		},
		Data: map[string]interface{}{
			"raw": "test",
			"ttl": "1h",
		},
	}

	resp.Secret.InternalData = nil
	resp.Secret.LeaseID = ""
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}
}

func TestPassthroughBackend_Delete(t *testing.T) {
	b := testPassthroughBackend()
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"
	storage := req.Storage

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.DeleteOperation, "foo")
	req.Storage = storage
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestPassthroughBackend_List(t *testing.T) {
	b := testPassthroughBackend()
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"
	storage := req.Storage

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ListOperation, "")
	req.Storage = storage
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := &logical.Response{
		Data: map[string]interface{}{
			"keys": []string{"foo"},
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}
}

func testPassthroughBackend() logical.Backend {
	b, _ := PassthroughBackendFactory(&logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 30,
		},
	})
	return b
}
