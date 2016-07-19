package http

import (
	"testing"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/vault"
)

func TestSysRenew(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// write secret
	resp := testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data":  "bar",
		"lease": "1h",
	})
	testResponseStatus(t, resp, 204)

	// read secret
	resp = testHttpGet(t, token, addr+"/v1/secret/foo")
	var result struct {
		LeaseId string `json:"lease_id"`
	}
	if err := jsonutil.DecodeJSONFromReader(resp.Body, &result); err != nil {
		t.Fatalf("bad: %s", err)
	}

	resp = testHttpPut(t, token, addr+"/v1/sys/renew/"+result.LeaseId, nil)
	testResponseStatus(t, resp, 200)
}

func TestSysRevoke(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/revoke/secret/foo/1234", nil)
	testResponseStatus(t, resp, 204)
}

func TestSysRevokePrefix(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/revoke-prefix/secret/foo/1234", nil)
	testResponseStatus(t, resp, 204)
}
