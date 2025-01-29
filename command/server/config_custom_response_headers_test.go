// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"testing"

	"github.com/go-test/deep"
)

var defaultCustomHeaders = map[string]string{
	"Strict-Transport-Security": "max-age=1; domains",
	"Content-Security-Policy":   "default-src 'others'",
	"X-Vault-Ignored":           "ignored",
	"X-Custom-Header":           "Custom header value default",
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

var defaultCustomHeadersMultiListener = map[string]string{
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	"Content-Security-Policy":   "default-src 'others'",
	"X-Vault-Ignored":           "ignored",
	"X-Custom-Header":           "Custom header value default",
}

var defaultSTS = map[string]string{
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
}

func TestCustomResponseHeadersConfigs(t *testing.T) {
	expectedCustomResponseHeader := map[string]map[string]string{
		"default": defaultCustomHeaders,
		"307":     customHeaders307,
		"3xx":     customHeader3xx,
		"200":     customHeaders200,
		"2xx":     customHeader2xx,
		"400":     customHeader400,
	}

	config, err := LoadConfigFile("./test-fixtures/config_custom_response_headers_1.hcl")
	if err != nil {
		t.Fatalf("Error encountered when loading config %+v", err)
	}
	if diff := deep.Equal(expectedCustomResponseHeader, config.Listeners[0].CustomResponseHeaders); diff != nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}
}

func TestCustomResponseHeadersConfigsMultipleListeners(t *testing.T) {
	expectedCustomResponseHeader := map[string]map[string]string{
		"default": defaultCustomHeadersMultiListener,
		"307":     customHeaders307,
		"3xx":     customHeader3xx,
		"200":     customHeaders200,
		"2xx":     customHeader2xx,
		"400":     customHeader400,
	}

	config, err := LoadConfigFile("./test-fixtures/config_custom_response_headers_multiple_listeners.hcl")
	if err != nil {
		t.Fatalf("Error encountered when loading config %+v", err)
	}
	if diff := deep.Equal(expectedCustomResponseHeader, config.Listeners[0].CustomResponseHeaders); diff != nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}

	if diff := deep.Equal(expectedCustomResponseHeader, config.Listeners[1].CustomResponseHeaders); diff == nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}
	if diff := deep.Equal(expectedCustomResponseHeader["default"], config.Listeners[1].CustomResponseHeaders["default"]); diff != nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}

	if diff := deep.Equal(expectedCustomResponseHeader, config.Listeners[2].CustomResponseHeaders); diff == nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}

	if diff := deep.Equal(defaultSTS, config.Listeners[2].CustomResponseHeaders["default"]); diff != nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}

	if diff := deep.Equal(expectedCustomResponseHeader, config.Listeners[3].CustomResponseHeaders); diff == nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}

	if diff := deep.Equal(defaultSTS, config.Listeners[3].CustomResponseHeaders["default"]); diff != nil {
		t.Fatalf("parsed custom headers do not match the expected ones, difference: %v", diff)
	}
}
