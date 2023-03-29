package pki

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2/json"
)

// TestAcmeDirectory a basic test that will validate the various directory APIs
// are available and produce the correct responses.
func TestAcmeDirectory(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Setting templated AIAs should succeed.
	pathConfig := "https://localhost:8200/v1/pki"

	_, err := CBWrite(b, s, "config/cluster", map[string]interface{}{
		"path":     pathConfig,
		"aia_path": "http://localhost:8200/cdn/pki",
	})
	require.NoError(t, err)

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
