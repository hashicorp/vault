// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"log"
	"testing"
)

var ecdsa_csr = `
-----BEGIN CERTIFICATE REQUEST-----
MIIBcTCB9wIBADB4MQswCQYDVQQGEwJBVTEOMAwGA1UECAwFU3RhdGUxDjAMBgNV
BAcMBUxvY2FsMQwwCgYDVQQKDANPcmcxEDAOBgNVBAsMB1NlY3Rpb24xDTALBgNV
BAMMBE5hbWUxGjAYBgkqhkiG9w0BCQEWC21lQG1haWwuY29tMHYwEAYHKoZIzj0C
AQYFK4EEACIDYgAEp9ZME8XMVDsJ/dxJpgY40HwgCX2gOmZ/vMbl3NwwKgvrbhIx
nfFIlmB+iOWHCG0r5r9Skjg+WJsvX2xUyte/ojj79Vu76GOjfarlVempYIqIQEDt
Ivf5JzoJdObyiMSNoAAwCgYIKoZIzj0EAwMDaQAwZgIxAMtO9tP8KDMgbfnIkQ+v
uMd36nUzk16Eteo6x+8ZOGCtpNXGRCaJzVTHtOhwtde6/QIxALHawm9r2oaeq1Fe
lWH+0Qu3QpFhIZkJUdW8jL7cIRlKzJGkqKx5P9A6IMddYfxbCA==
-----END CERTIFICATE REQUEST-----`

func TestTransit_SignCSR(t *testing.T) {
	var resp *logical.Response
	var err error
	b, s := createBackendWithStorage(t)

	// Create the policy
	policyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/rsa-key",
		Storage:   s,
		Data: map[string]interface{}{
			"type": "rsa-2048",
		},
	}
	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	csrSignReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/rsa-key/csr",
		Storage:   s,
		Data: map[string]interface{}{
			"csr": ecdsa_csr,
		},
	}
	resp, err = b.HandleRequest(context.Background(), csrSignReq)
	// FIXME: Also this check?
	if err != nil || (resp != nil && resp.IsError()) {
		// FIXME: Set an error message
		t.Fatalf("Failed to sign CSR, err:%v resp:%#v", err, resp)
	}

	log.Printf("CSR: %s", resp.Data["csr"])
	t.Fatal("Fail")
}
