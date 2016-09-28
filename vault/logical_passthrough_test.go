package vault

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestPassthroughBackend_RootPaths(t *testing.T) {
	b := testPassthroughBackend()
	test := func(b logical.Backend) {
		root := b.SpecialPaths()
		if root != nil {
			t.Fatalf("unexpected: %v", root)
		}
	}
	test(b)
	b = testPassthroughLeasedBackend()
	test(b)
}

func TestPassthroughBackend_Write(t *testing.T) {
	test := func(b logical.Backend) {
		req := logical.TestRequest(t, logical.UpdateOperation, "foo")
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
	b := testPassthroughBackend()
	test(b)
	b = testPassthroughLeasedBackend()
	test(b)
}

func TestPassthroughBackend_Read(t *testing.T) {
	test := func(b logical.Backend, ttlType string, leased bool) {
		req := logical.TestRequest(t, logical.UpdateOperation, "foo")
		req.Data["raw"] = "test"
		req.Data[ttlType] = "1h"
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
				ttlType: "1h",
			},
		}

		if !leased {
			expected.Secret.Renewable = false
		}
		resp.Secret.InternalData = nil
		resp.Secret.LeaseID = ""
		if !reflect.DeepEqual(resp, expected) {
			t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
		}
	}
	b := testPassthroughLeasedBackend()
	test(b, "lease", true)
	test(b, "ttl", true)
	b = testPassthroughBackend()
	test(b, "lease", false)
	test(b, "ttl", false)
}

func TestPassthroughBackend_Delete(t *testing.T) {
	test := func(b logical.Backend) {
		req := logical.TestRequest(t, logical.UpdateOperation, "foo")
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
	b := testPassthroughBackend()
	test(b)
	b = testPassthroughLeasedBackend()
	test(b)
}

func TestPassthroughBackend_List(t *testing.T) {
	test := func(b logical.Backend) {
		req := logical.TestRequest(t, logical.UpdateOperation, "foo")
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
	b := testPassthroughBackend()
	test(b)
	b = testPassthroughLeasedBackend()
	test(b)
}

func TestPassthroughBackend_Revoke(t *testing.T) {
	test := func(b logical.Backend) {
		req := logical.TestRequest(t, logical.RevokeOperation, "generic")
		req.Secret = &logical.Secret{
			InternalData: map[string]interface{}{
				"secret_type": "generic",
			},
		}

		if _, err := b.HandleRequest(req); err != nil {
			t.Fatalf("err: %v", err)
		}
	}
	b := testPassthroughBackend()
	test(b)
	b = testPassthroughLeasedBackend()
	test(b)
}

func testPassthroughBackend() logical.Backend {
	b, _ := PassthroughBackendFactory(&logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	})
	return b
}

func testPassthroughLeasedBackend() logical.Backend {
	b, _ := LeasedPassthroughBackendFactory(&logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	})
	return b
}
