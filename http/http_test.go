package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/jsonutil"
)

func testHttpGet(t *testing.T, token string, addr string) *http.Response {
	t.Logf("Token is %s", token)
	return testHttpData(t, "GET", token, addr, nil)
}

func testHttpDelete(t *testing.T, token string, addr string) *http.Response {
	return testHttpData(t, "DELETE", token, addr, nil)
}

func testHttpPost(t *testing.T, token string, addr string, body interface{}) *http.Response {
	return testHttpData(t, "POST", token, addr, body)
}

func testHttpPut(t *testing.T, token string, addr string, body interface{}) *http.Response {
	return testHttpData(t, "PUT", token, addr, body)
}

func testHttpData(t *testing.T, method string, token string, addr string, body interface{}) *http.Response {
	bodyReader := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(bodyReader)
		if err := enc.Encode(body); err != nil {
			t.Fatalf("err:%s", err)
		}
	}

	req, err := http.NewRequest(method, addr, bodyReader)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if len(token) != 0 {
		req.Header.Set("X-Vault-Token", token)
	}

	client := cleanhttp.DefaultClient()
	client.Timeout = 60 * time.Second

	// From https://github.com/michiwend/gomusicbrainz/pull/4/files
	defaultRedirectLimit := 30

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > defaultRedirectLimit {
			return fmt.Errorf("%d consecutive requests(redirects)", len(via))
		}
		if len(via) == 0 {
			// No redirects
			return nil
		}
		// mutate the subsequent redirect requests with the first Header
		if token := via[0].Header.Get("X-Vault-Token"); len(token) != 0 {
			req.Header.Set("X-Vault-Token", token)
		}
		return nil
	}

	resp, err := client.Do(req)
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

	if err := jsonutil.DecodeJSONFromReader(resp.Body, out); err != nil {
		t.Fatalf("err: %s", err)
	}
}
