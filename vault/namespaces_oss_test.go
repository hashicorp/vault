//go:build !enterprise

package vault

import (
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestSystemBackend_Namespaces_OSS_CRD(t *testing.T) {
	var resp *logical.Response
	var err error
	// Create the backend
	b := testSystemBackend(t)

	// Create a namespace
	resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "namespaces/ns1",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
	}
}
