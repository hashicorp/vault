// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package system

import (
	"testing"

	"github.com/hashicorp/vault/command/server"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestCustomResponseHeadersConfigInteractUiConfig verifies that headers already
// present in the listener custom response headers configuration cannot be set
// via the UI headers API, while headers not present in the listener config can.
func TestCustomResponseHeadersConfigInteractUiConfig(t *testing.T) {
	t.Parallel()

	customResponseHeaders := map[string]map[string]string{
		"default": {
			"Strict-Transport-Security": "max-age=1; domains",
			"Content-Security-Policy":   "default-src 'others'",
			"X-Vault-Ignored":           "ignored",
			"X-Custom-Header":           "Custom header value default",
			"X-Frame-Options":           "Deny",
			"X-Content-Type-Options":    "nosniff",
			"Content-Type":              "text/plain; charset=utf-8",
			"X-XSS-Protection":          "1; mode=block",
		},
		"307": {"X-Custom-Header": "Custom header value 307"},
		"3xx": {
			"X-Vault-Ignored-3xx": "Ignored 3xx",
			"X-Custom-Header":     "Custom header value 3xx",
		},
		"200": {
			"Someheader-200":  "200",
			"X-Custom-Header": "Custom header value 200",
		},
		"2xx": {"X-Custom-Header": "Custom header value 2xx"},
		"400": {"Someheader-400": "400"},
	}

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		RawConfig: &server.Config{
			SharedConfig: &configutil.SharedConfig{
				Listeners: []*configutil.Listener{
					{
						Type:                  "tcp",
						Address:               "127.0.0.1:443",
						CustomResponseHeaders: customResponseHeaders,
					},
				},
			},
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	client := cluster.Cores[0].Client

	t.Run("cannot set header that exists in listener custom response headers", func(t *testing.T) {
		_, err := client.Logical().Write("sys/config/ui/headers/X-Custom-Header", map[string]interface{}{
			"values": []string{"UI Custom Header"},
		})
		require.ErrorContains(t, err, "This header already exists in the server configuration and cannot be set in the UI.")
	})

	t.Run("cannot set status-code-specific header that exists in listener config", func(t *testing.T) {
		_, err := client.Logical().Write("sys/config/ui/headers/Someheader-400", map[string]interface{}{
			"values": []string{"400"},
		})
		require.ErrorContains(t, err, "This header already exists in the server configuration and cannot be set in the UI.")

		// Confirm the header was not persisted.
		resp, err := client.Logical().Read("sys/config/ui/headers/Someheader-400")
		require.NoError(t, err)
		require.Nil(t, resp)
	})

	t.Run("can set header that does not exist in listener custom response headers", func(t *testing.T) {
		_, err := client.Logical().Write("sys/config/ui/headers/X-CustomUiHeader", map[string]interface{}{
			"values": []string{"Ui header value"},
		})
		require.NoError(t, err)

		resp, err := client.Logical().ReadWithData("sys/config/ui/headers/X-CustomUiHeader", map[string][]string{
			"multivalue": {"true"},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, []interface{}{"Ui header value"}, resp.Data["values"])
	})
}
