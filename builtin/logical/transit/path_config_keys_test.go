// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestTransit_ConfigKeys(t *testing.T) {
	b, s := createBackendWithSysView(t)

	doReq := func(req *logical.Request) *logical.Response {
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("got err:\n%#v\nreq:\n%#v\n", err, *req)
		}
		return resp
	}
	doErrReq := func(req *logical.Request) {
		resp, err := b.HandleRequest(context.Background(), req)
		if err == nil {
			if resp == nil || !resp.IsError() {
				t.Fatalf("expected error; req:\n%#v\n", *req)
			}
		}
	}

	// First read the global config
	req := &logical.Request{
		Storage:   s,
		Operation: logical.ReadOperation,
		Path:      "config/keys",
	}
	resp := doReq(req)
	if resp.Data["disable_upsert"].(bool) != false {
		t.Fatalf("expected disable_upsert to be false; got: %v", resp)
	}

	// Ensure we can upsert.
	req.Operation = logical.CreateOperation
	req.Path = "encrypt/upsert-1"
	req.Data = map[string]interface{}{
		"plaintext": "aGVsbG8K",
	}
	doReq(req)

	// Disable upserting.
	req.Operation = logical.UpdateOperation
	req.Path = "config/keys"
	req.Data = map[string]interface{}{
		"disable_upsert": true,
	}
	doReq(req)

	// Attempt upserting again, it should fail.
	req.Operation = logical.CreateOperation
	req.Path = "encrypt/upsert-2"
	req.Data = map[string]interface{}{
		"plaintext": "aGVsbG8K",
	}
	doErrReq(req)

	// Redoing this with the first key should succeed.
	req.Path = "encrypt/upsert-1"
	doReq(req)
}
