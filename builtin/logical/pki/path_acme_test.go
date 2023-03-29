package pki

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2/json"
)

// TestAcmeDirectory a basic test that will validate the various directory APIs
// are available and produce the correct responses.
func TestAcmeDirectory(t *testing.T) {
	t.Parallel()
	b, s, pathConfig := setupAcmeBackend(t)

	cases := []struct {
		name         string
		prefixUrl    string
		directoryUrl string
	}{
		{"root", "", "acme/directory"},
		{"role", "/roles/test-role", "roles/test-role/acme/directory"},
		{"issuer", "/issuer/default", "issuer/default/acme/directory"},
		{"issuer_role", "/issuer/default/roles/test-role", "issuer/default/roles/test-role/acme/directory"},
		{"issuer_role_acme", "/issuer/acme/roles/acme", "issuer/acme/roles/acme/acme/directory"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dirResp, err := CBRead(b, s, tc.directoryUrl)
			require.NoError(t, err, "failed reading ACME directory configuration")

			require.Contains(t, dirResp.Data, "http_content_type", "missing Content-Type header")
			require.Contains(t, dirResp.Data["http_content_type"], "application/json",
				"missing appropriate content type in header")

			requiredUrls := map[string]string{
				"newNonce":   pathConfig + tc.prefixUrl + "/acme/new-nonce",
				"newAccount": pathConfig + tc.prefixUrl + "/acme/new-account",
				"newOrder":   pathConfig + tc.prefixUrl + "/acme/new-order",
				"revokeCert": pathConfig + tc.prefixUrl + "/acme/revoke-cert",
				"keyChange":  pathConfig + tc.prefixUrl + "/acme/key-change",
			}

			rawBodyBytes := dirResp.Data["http_raw_body"].([]byte)
			respType := map[string]interface{}{}
			err = json.Unmarshal(rawBodyBytes, &respType)
			require.NoError(t, err, "failed unmarshalling ACME directory response body")

			for key, expectedUrl := range requiredUrls {
				require.Contains(t, respType, key, "missing required value %s from data", key)
				require.Equal(t, expectedUrl, respType[key], "different URL returned for %s", key)
			}
		})
	}
}

func TestAcmeNonce(t *testing.T) {
	t.Parallel()
	b, s, pathConfig := setupAcmeBackend(t)

	cases := []struct {
		name         string
		prefixUrl    string
		directoryUrl string
	}{
		{"root", "", "acme/new-nonce"},
		{"role", "/roles/test-role", "roles/test-role/acme/new-nonce"},
		{"issuer", "/issuer/default", "issuer/default/acme/new-nonce"},
		{"issuer_role", "/issuer/default/roles/test-role", "issuer/default/roles/test-role/acme/new-nonce"},
	}
	for _, tc := range cases {
		for _, httpOp := range []string{"get", "header"} {
			t.Run(fmt.Sprintf("%s-%s", tc.name, httpOp), func(t *testing.T) {
				var resp *logical.Response
				var err error
				switch httpOp {
				case "get":
					resp, err = CBRead(b, s, tc.directoryUrl)
				case "header":
					resp, err = CBHeader(b, s, tc.directoryUrl)
				}
				require.NoError(t, err, "failed %s op for new-nouce", httpOp)

				// Proper Status Code
				require.Equal(t, http.StatusOK, resp.Data["http_status_code"])

				// Make sure we return the Cache-Control header
				require.Contains(t, resp.Headers, "Cache-Control", "missing Cache-Control header")
				require.Contains(t, resp.Headers["Cache-Control"], "no-store",
					"missing Cache-Control header with no-store header value")
				require.Len(t, resp.Headers["Cache-Control"], 1,
					"Cache-Control header should have only a single header")

				// Test for our nonce header value
				require.Contains(t, resp.Headers, "Replay-Nonce", "missing Replay-Nonce header")
				require.NotEmpty(t, resp.Headers["Replay-Nonce"], "missing Replay-Nonce header with an actual value")
				require.Len(t, resp.Headers["Replay-Nonce"], 1,
					"Replay-Nonce header should have only a single header")

				// Test Link header value
				require.Contains(t, resp.Headers, "Link", "missing Link header")
				expectedLinkHeader := fmt.Sprintf("<%s>;rel=\"index\"", pathConfig+tc.prefixUrl+"/acme/directory")
				require.Contains(t, resp.Headers["Link"], expectedLinkHeader,
					"different value for link header than expected")
				require.Len(t, resp.Headers["Link"], 1, "Link header should have only a single header")
			})
		}
	}
}

// TestAcmeClusterPathNotConfigured basic testing of the ACME error handler.
func TestAcmeClusterPathNotConfigured(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Do not fill in the path option within the local cluster configuration
	cases := []struct {
		name         string
		directoryUrl string
	}{
		{"root", "acme/directory"},
		{"role", "roles/test-role/acme/directory"},
		{"issuer", "issuer/default/acme/directory"},
		{"issuer_role", "issuer/default/roles/test-role/acme/directory"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dirResp, err := CBRead(b, s, tc.directoryUrl)
			require.NoError(t, err, "failed reading ACME directory configuration")

			require.Contains(t, dirResp.Data, "http_content_type", "missing Content-Type header")
			require.Contains(t, dirResp.Data["http_content_type"], "application/problem+json",
				"missing appropriate content type in header")

			require.Equal(t, http.StatusInternalServerError, dirResp.Data["http_status_code"])

			require.Contains(t, dirResp.Data, "http_raw_body", "missing http_raw_body from data")
			rawBodyBytes := dirResp.Data["http_raw_body"].([]byte)
			respType := map[string]interface{}{}
			err = json.Unmarshal(rawBodyBytes, &respType)
			require.NoError(t, err, "failed unmarshalling ACME directory response body")

			require.Equal(t, "urn:ietf:params:acme:error:serverInternal", respType["type"])
			require.NotEmpty(t, respType["detail"])
		})
	}
}

func setupAcmeBackend(t *testing.T) (*backend, logical.Storage, string) {
	b, s := CreateBackendWithStorage(t)

	// Setting templated AIAs should succeed.
	pathConfig := "https://localhost:8200/v1/pki"

	_, err := CBWrite(b, s, "config/cluster", map[string]interface{}{
		"path":     pathConfig,
		"aia_path": "http://localhost:8200/cdn/pki",
	})
	require.NoError(t, err)
	return b, s, pathConfig
}
