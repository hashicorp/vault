package vault

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/logical"
)

func TestPassthroughBackend_RootPaths(t *testing.T) {
	b := testPassthroughBackend()
	test := func(b logical.Backend) {
		root := b.SpecialPaths()
		if len(root.Root) != 0 {
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

		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %v", resp)
		}

		out, err := req.Storage.Get(context.Background(), "foo")
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
	test := func(b logical.Backend, ttlType string, ttl interface{}, leased bool) {
		req := logical.TestRequest(t, logical.UpdateOperation, "foo")
		req.Data["raw"] = "test"
		var reqTTL interface{}
		switch ttl.(type) {
		case int64:
			reqTTL = ttl.(int64)
		case string:
			reqTTL = ttl.(string)
		default:
			t.Fatal("unknown ttl type")
		}
		req.Data[ttlType] = reqTTL
		storage := req.Storage

		if _, err := b.HandleRequest(context.Background(), req); err != nil {
			t.Fatalf("err: %v", err)
		}

		req = logical.TestRequest(t, logical.ReadOperation, "foo")
		req.Storage = storage

		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		expectedTTL, err := parseutil.ParseDurationSecond(ttl)
		if err != nil {
			t.Fatal(err)
		}

		// What comes back if an int is passed in is a json.Number which is
		// actually aliased as a string so to make the deep equal happy if it's
		// actually a number we set it to an int64
		var respTTL interface{} = resp.Data[ttlType]
		_, ok := respTTL.(json.Number)
		if ok {
			respTTL, err = respTTL.(json.Number).Int64()
			if err != nil {
				t.Fatal(err)
			}
			resp.Data[ttlType] = respTTL
		}

		expected := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					Renewable: true,
					TTL:       expectedTTL,
				},
			},
			Data: map[string]interface{}{
				"raw":   "test",
				ttlType: reqTTL,
			},
		}

		if !leased {
			expected.Secret.Renewable = false
		}
		resp.Secret.InternalData = nil
		resp.Secret.LeaseID = ""
		if !reflect.DeepEqual(resp, expected) {
			t.Fatalf("bad response.\n\nexpected:\n%#v\n\nGot:\n%#v", expected, resp)
		}
	}
	b := testPassthroughLeasedBackend()
	test(b, "lease", "1h", true)
	test(b, "ttl", "5", true)
	b = testPassthroughBackend()
	test(b, "lease", int64(10), false)
	test(b, "ttl", "40s", false)
}

func TestPassthroughBackend_Delete(t *testing.T) {
	test := func(b logical.Backend) {
		req := logical.TestRequest(t, logical.UpdateOperation, "foo")
		req.Data["raw"] = "test"
		storage := req.Storage

		if _, err := b.HandleRequest(context.Background(), req); err != nil {
			t.Fatalf("err: %v", err)
		}

		req = logical.TestRequest(t, logical.DeleteOperation, "foo")
		req.Storage = storage
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %v", resp)
		}

		req = logical.TestRequest(t, logical.ReadOperation, "foo")
		req.Storage = storage
		resp, err = b.HandleRequest(context.Background(), req)
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

		if _, err := b.HandleRequest(context.Background(), req); err != nil {
			t.Fatalf("err: %v", err)
		}

		req = logical.TestRequest(t, logical.ListOperation, "")
		req.Storage = storage
		resp, err := b.HandleRequest(context.Background(), req)
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
		req := logical.TestRequest(t, logical.RevokeOperation, "kv")
		req.Secret = &logical.Secret{
			InternalData: map[string]interface{}{
				"secret_type": "kv",
			},
		}

		if _, err := b.HandleRequest(context.Background(), req); err != nil {
			t.Fatalf("err: %v", err)
		}
	}
	b := testPassthroughBackend()
	test(b)
	b = testPassthroughLeasedBackend()
	test(b)
}

func testPassthroughBackend() logical.Backend {
	b, _ := PassthroughBackendFactory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	})
	return b
}

func testPassthroughLeasedBackend() logical.Backend {
	b, _ := LeasedPassthroughBackendFactory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	})
	return b
}
