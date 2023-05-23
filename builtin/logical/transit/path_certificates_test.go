// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"testing"
)

func TestTransit_SignEmptyCSR(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	// Create the policy
	policyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/existing_key",
		Storage:   s,
	}
	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	csrSignReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/existing_key/csr",
		Storage:   s,
		Data:      map[string]interface{}{},
	}
	resp, err = b.HandleRequest(context.Background(), csrSignReq)
	if resp == nil || !resp.IsError() {
		// FIXME: Set an error message
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
}
