package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

const targetCacheSize = 12345

func TestTransit_CacheConfig(t *testing.T) {
	b1, storage := createBackendWithSysView(t)

	doReq := func(b *backend, req *logical.Request) *logical.Response {
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("got err:\n%#v\nreq:\n%#v\n", err, *req)
		}
		return resp
	}

	doErrReq := func(b *backend, req *logical.Request) {
		resp, err := b.HandleRequest(context.Background(), req)
		if err == nil {
			if resp == nil || !resp.IsError() {
				t.Fatalf("expected error; req:\n%#v\n", *req)
			}
		}
	}

	validateResponse := func(resp *logical.Response, expectedCacheSize int, expectedWarning bool) {
		actualCacheSize, ok := resp.Data["size"].(int)
		if !ok {
			t.Fatalf("No size returned")
		}
		if expectedCacheSize != actualCacheSize {
			t.Fatalf("testAccReadCacheConfig expected: %d got: %d", expectedCacheSize, actualCacheSize)
		}
		// check for the presence/absence of warnings - warnings are expected if a cache size has been
		// configured but not yet applied by reloading the plugin
		warningCheckPass := expectedWarning == (len(resp.Warnings) > 0)
		if !warningCheckPass {
			t.Fatalf(
				"testAccSteporeadCacheConfig warnings error.\n"+
					"expect warnings: %t but number of warnings was: %d",
				expectedWarning, len(resp.Warnings),
			)
		}
	}

	writeReq := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "cache-config",
		Data: map[string]interface{}{
			"size": targetCacheSize,
		},
	}

	readReq := &logical.Request{
		Storage:   storage,
		Operation: logical.ReadOperation,
		Path:      "cache-config",
	}

	// test steps
	// b1 should spin up with an unlimited cache
	validateResponse(doReq(b1, readReq), 0, false)
	doReq(b1, writeReq)
	validateResponse(doReq(b1, readReq), targetCacheSize, true)

	// b2 should spin up with a configured cache
	b2 := createBackendWithSysViewWithStorage(t, storage)
	validateResponse(doReq(b2, readReq), targetCacheSize, false)

	// b3 enables transit without a cache, trying to read it should error
	b3 := createBackendWithForceNoCacheWithSysViewWithStorage(t, storage)
	doErrReq(b3, readReq)
}
