// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
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
