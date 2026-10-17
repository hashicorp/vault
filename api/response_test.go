// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func testErrorResponse(target string, requestAddress string) *Response {
	req, _ := http.NewRequest("GET", target, nil)
	return &Response{
		Response: &http.Response{
			StatusCode: 403,
			Body:       io.NopCloser(strings.NewReader(`{"errors":["permission denied"]}`)),
			Request:    req,
		},
		requestAddress: requestAddress,
	}
}

func TestResponseErrorReportsRequestURL(t *testing.T) {
	resp := testErrorResponse("http://127.0.0.1:8200/v1/sys/key-status", "")
	err := resp.Error()
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "URL: GET http://127.0.0.1:8200/v1/sys/key-status") {
		t.Fatalf("expected request URL in error, got:\n%s", err.Error())
	}
}

func TestResponseErrorReportsConfiguredUnixAddress(t *testing.T) {
	resp := testErrorResponse(
		"http://localhost/v1/sys/key-status",
		"unix:///run/user/137/ca-vault.sock",
	)
	err := resp.Error()
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "URL: GET unix:///run/user/137/ca-vault.sock/v1/sys/key-status") {
		t.Fatalf("expected configured unix address in error, got:\n%s", err.Error())
	}
	if strings.Contains(err.Error(), "localhost") {
		t.Fatalf("did not expect synthetic localhost URL in error, got:\n%s", err.Error())
	}
}
