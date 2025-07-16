// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
)

var defaultCustomHeaders = map[string]string{
	"Strict-Transport-Security": "max-age=1; domains",
	"Content-Security-Policy":   "default-src 'others'",
	"X-Vault-Ignored":           "ignored",
	"X-Custom-Header":           "Custom header value default",
	"X-Frame-Options":           "Deny",
	"X-Content-Type-Options":    "nosniff",
	"Content-Type":              "text/plain; charset=utf-8",
	"X-XSS-Protection":          "1; mode=block",
}

var customHeaders307 = map[string]string{
	"X-Custom-Header": "Custom header value 307",
}

var customHeader3xx = map[string]string{
	"X-Vault-Ignored-3xx": "Ignored 3xx",
	"X-Custom-Header":     "Custom header value 3xx",
}

var customHeaders200 = map[string]string{
	"Someheader-200":  "200",
	"X-Custom-Header": "Custom header value 200",
}

var customHeader2xx = map[string]string{
	"X-Custom-Header": "Custom header value 2xx",
}

var customHeader400 = map[string]string{
	"Someheader-400": "400",
}

func TestConfigCustomHeaders(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
	phys, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	logl := &logical.InmemStorage{}
	uiConfig := NewUIConfig(true, phys, logl)

	rawListenerConfig := []*configutil.Listener{
		{
			Type:    "tcp",
			Address: "127.0.0.1:443",
			CustomResponseHeaders: map[string]map[string]string{
				"default": defaultCustomHeaders,
				"307":     customHeaders307,
				"3xx":     customHeader3xx,
				"200":     customHeaders200,
				"2xx":     customHeader2xx,
				"400":     customHeader400,
			},
		},
	}

	uiHeaders, err := uiConfig.Headers(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	listenerCustomHeaders := NewListenerCustomHeader(rawListenerConfig, logger, uiHeaders)
	if listenerCustomHeaders == nil || len(listenerCustomHeaders) != 1 {
		t.Fatalf("failed to get custom header configuration")
	}

	lch := listenerCustomHeaders[0]

	if lch.ExistCustomResponseHeader("X-Vault-Ignored-307") {
		t.Fatalf("header name with X-Vault prefix is not valid")
	}
	if lch.ExistCustomResponseHeader("X-Vault-Ignored-3xx") {
		t.Fatalf("header name with X-Vault prefix is not valid")
	}

	if !lch.ExistCustomResponseHeader("X-Custom-Header") {
		t.Fatalf("header name with X-Vault prefix is not valid")
	}
}

func TestCustomResponseHeadersConfigInteractUiConfig(t *testing.T) {
	b := testSystemBackend(t)
	paths := b.(*SystemBackend).configPaths()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "")
	b.(*SystemBackend).Core.systemBarrierView = view

	logger := logging.NewVaultLogger(log.Trace)
	rawListenerConfig := []*configutil.Listener{
		{
			Type:    "tcp",
			Address: "127.0.0.1:443",
			CustomResponseHeaders: map[string]map[string]string{
				"default": defaultCustomHeaders,
				"307":     customHeaders307,
				"3xx":     customHeader3xx,
				"200":     customHeaders200,
				"2xx":     customHeader2xx,
				"400":     customHeader400,
			},
		},
	}
	uiHeaders, err := b.(*SystemBackend).Core.uiConfig.Headers(context.Background())
	if err != nil {
		t.Fatalf("failed to get headers from ui config")
	}
	customListenerHeader := NewListenerCustomHeader(rawListenerConfig, logger, uiHeaders)
	if customListenerHeader == nil {
		t.Fatalf("custom header config should be configured")
	}
	b.(*SystemBackend).Core.customListenerHeader.Store(customListenerHeader)
	clh := b.(*SystemBackend).Core.customListenerHeader
	if clh == nil {
		t.Fatalf("custom header config should be configured in core")
	}

	w := httptest.NewRecorder()
	hw := logical.NewHTTPResponseWriter(w)

	// setting a header that already exist in custom headers
	req := logical.TestRequest(t, logical.UpdateOperation, "config/ui/headers/X-Custom-Header")
	req.Data["values"] = []string{"UI Custom Header"}
	req.ResponseWriter = hw

	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatal("request did not fail on setting a header that is present in custom response headers")
	}
	schema.ValidateResponse(
		t,
		schema.FindResponseSchema(t, paths, 3, req.Operation),
		resp,
		true,
	)

	if !strings.Contains(resp.Data["error"].(string), fmt.Sprintf("This header already exists in the server configuration and cannot be set in the UI.")) {
		t.Fatalf("failed to get the expected error")
	}

	// setting a header that already exist in custom headers
	req = logical.TestRequest(t, logical.UpdateOperation, "config/ui/headers/Someheader-400")
	req.Data["values"] = []string{"400"}
	req.ResponseWriter = hw

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatal("request did not fail on setting a header that is present in custom response headers")
	}
	schema.ValidateResponse(
		t,
		schema.FindResponseSchema(t, paths, 3, req.Operation),
		resp,
		true,
	)

	h, err := b.(*SystemBackend).Core.uiConfig.Headers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if h.Get("Someheader-400") == "400" {
		t.Fatalf("should not be able to set a header that is in custom response headers")
	}

	// setting an ui specific header
	req = logical.TestRequest(t, logical.UpdateOperation, "config/ui/headers/X-CustomUiHeader")
	req.Data["values"] = []string{"Ui header value"}
	req.ResponseWriter = hw

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal("request failed on setting a header that is not present in custom response headers.", "error:", err)
	}
	schema.ValidateResponse(
		t,
		schema.FindResponseSchema(t, paths, 3, req.Operation),
		resp,
		true,
	)

	h, err = b.(*SystemBackend).Core.uiConfig.Headers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if h.Get("X-CustomUiHeader") != "Ui header value" {
		t.Fatalf("failed to set a header that is not in custom response headers")
	}
}
