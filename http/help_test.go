// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"
	"testing"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestHelp(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// request without /v1/ prefix
	resp := testHttpGet(t, token, addr+"/?help=1")
	testResponseStatus(t, resp, 404)

	resp = testHttpGet(t, "", addr+"/v1/sys/mounts?help=1")
	if resp.StatusCode != http.StatusForbidden {
		t.Fatal("expected permission denied with no token")
	}

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts?help=1")

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if _, ok := actual["help"]; !ok {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestHelp_AuditRequestID(t *testing.T) {
	noop := audit.TestNoopAudit(t, "noop/", nil)
	c, _, root := vault.TestCoreUnsealedWithConfig(t, &vault.CoreConfig{
		AuditBackends: map[string]audit.Factory{
			"noop": func(config *audit.BackendConfig, _ audit.HeaderFormatter) (audit.Backend, error) {
				return noop, nil
			},
		},
	})
	ln, addr := TestServer(t, c)
	defer ln.Close()
	TestServerAuth(t, addr, root)

	// Enable the audit backend
	resp := testHttpPost(t, root, addr+"/v1/sys/audit/noop", map[string]interface{}{
		"type": "noop",
	})
	testResponseStatus(t, resp, 204)

	// Make a help request
	resp = testHttpGet(t, root, addr+"/v1/sys/mounts?help=1")
	testResponseStatus(t, resp, 200)

	// Find the help request in the audit trail
	var helpReq *logical.Request
	for _, r := range noop.Req {
		if r.Operation == logical.HelpOperation {
			helpReq = r
			break
		}
	}
	if helpReq == nil {
		t.Fatalf("no help request found in audit trail; got %d requests", len(noop.Req))
	}
	if helpReq.ID == "" {
		t.Fatal("help request in audit trail has empty request ID")
	}

	// Find the matching response entry and verify it carries the same request ID
	var helpRespReq *logical.Request
	for _, r := range noop.RespReq {
		if r.Operation == logical.HelpOperation {
			helpRespReq = r
			break
		}
	}
	if helpRespReq == nil {
		t.Fatalf("no help response found in audit trail; got %d responses", len(noop.RespReq))
	}
	if helpRespReq.ID == "" {
		t.Fatal("help response in audit trail has empty request ID")
	}
	if helpRespReq.ID != helpReq.ID {
		t.Fatalf("request ID mismatch: request has %q, response has %q", helpReq.ID, helpRespReq.ID)
	}
}
