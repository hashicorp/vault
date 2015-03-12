package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
)

func testCore(t *testing.T) *vault.Core {
	physicalBackend := physical.NewInmem()
	c, err := vault.NewCore(&vault.CoreConfig{
		Physical: physicalBackend,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return c
}

func testCoreInit(t *testing.T, core *vault.Core) [][]byte {
	result, err := core.Initialize(&vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return result.SecretShares
}

func testServer(t *testing.T, core *vault.Core) (net.Listener, string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	addr := "http://" + ln.Addr().String()

	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: Handler(core),
	}
	go server.Serve(ln)

	return ln, addr
}

func testHttpPut(t *testing.T, addr string, body interface{}) *http.Response {
	bodyReader := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(bodyReader)
		if err := enc.Encode(body); err != nil {
			t.Fatalf("err:%s", err)
		}
	}

	req, err := http.NewRequest("PUT", addr, bodyReader)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return resp
}

func testResponseStatus(t *testing.T, resp *http.Response, code int) {
	if resp.StatusCode != code {
		body := new(bytes.Buffer)
		io.Copy(body, resp.Body)
		resp.Body.Close()

		t.Fatalf(
			"Expected status %d, got %d. Body:\n\n%s",
			code, resp.StatusCode, body.String())
	}
}

func testResponseBody(t *testing.T, resp *http.Response, out interface{}) {
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		t.Fatalf("err: %s", err)
	}
}
