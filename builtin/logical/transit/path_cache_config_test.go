// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	targetCacheSize = 12345
	smallCacheSize  = 3
)

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

	writeSmallCacheSizeReq := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "cache-config",
		Data: map[string]interface{}{
			"size": smallCacheSize,
		},
	}

	readReq := &logical.Request{
		Storage:   storage,
		Operation: logical.ReadOperation,
		Path:      "cache-config",
	}

	polReq := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/aes256",
		Data: map[string]interface{}{
			"derived": true,
		},
	}

	// test steps
	// b1 should spin up with an unlimited cache
	validateResponse(doReq(b1, readReq), 0, false)

	// Change cache size to targetCacheSize 12345 and validate that cache size is updated
	doReq(b1, writeReq)
	validateResponse(doReq(b1, readReq), targetCacheSize, false)
	b1.invalidate(context.Background(), "cache-config/")

	// Change the cache size to 1000 to mock the scenario where
	// current cache size and stored cache size are different and
	// a cache update is needed
	b1.lm.InitCache(1000)

	// Write a new policy which in its code path detects that cache size has changed
	// and refreshes the cache to 12345
	doReq(b1, polReq)

	// Validate that cache size is updated to 12345
	validateResponse(doReq(b1, readReq), targetCacheSize, false)

	// b2 should spin up with a configured cache
	b2 := createBackendWithSysViewWithStorage(t, storage)
	validateResponse(doReq(b2, readReq), targetCacheSize, false)

	// b3 enables transit without a cache, trying to read it should error
	b3 := createBackendWithForceNoCacheWithSysViewWithStorage(t, storage)
	doErrReq(b3, readReq)

	// b4 should spin up with a size less than minimum cache size (10)
	b4, storage := createBackendWithSysView(t)
	doErrReq(b4, writeSmallCacheSizeReq)
}
