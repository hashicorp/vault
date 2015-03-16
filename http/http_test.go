package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func testHttpDelete(t *testing.T, addr string) *http.Response {
	return testHttpData(t, "DELETE", addr, nil)
}

func testHttpPost(t *testing.T, addr string, body interface{}) *http.Response {
	return testHttpData(t, "POST", addr, body)
}

func testHttpPut(t *testing.T, addr string, body interface{}) *http.Response {
	return testHttpData(t, "PUT", addr, body)
}

func testHttpData(t *testing.T, method string, addr string, body interface{}) *http.Response {
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
