package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysRenew(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// write secret
	resp := testHttpPut(t, addr+"/v1/secret/foo", map[string]interface{}{
		"data":  "bar",
		"lease": "1h",
	})
	testResponseStatus(t, resp, 204)

	// read secret
	resp, err := http.Get(addr + "/v1/secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var result struct {
		LeaseId string `json:"lease_id"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("bad: %s", err)
	}

	resp = testHttpPut(t, addr+"/v1/sys/renew/"+result.LeaseId, nil)
	testResponseStatus(t, resp, 200)
}

func TestSysRevoke(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, addr+"/v1/sys/revoke/secret/foo/1234", nil)
	testResponseStatus(t, resp, 204)
}

func TestSysRevokePrefix(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, addr+"/v1/sys/revoke-prefix/secret/foo/1234", nil)
	testResponseStatus(t, resp, 204)
}
