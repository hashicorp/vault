package pki

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2/json"
)

// TestAcmeDirectory a basic test that will validate the various directory APIs
// are available and produce the correct responses.
func TestAcmeDirectory(t *testing.T) {
	t.Parallel()
	cluster, client, pathConfig := setupAcmeBackend(t)
	defer cluster.Cleanup()

	cases := []struct {
		name         string
		prefixUrl    string
		directoryUrl string
	}{
		{"root", "", "pki/acme/directory"},
		{"role", "/roles/test-role", "pki/roles/test-role/acme/directory"},
		{"issuer", "/issuer/default", "pki/issuer/default/acme/directory"},
		{"issuer_role", "/issuer/default/roles/test-role", "pki/issuer/default/roles/test-role/acme/directory"},
		{"issuer_role_acme", "/issuer/acme/roles/acme", "pki/issuer/acme/roles/acme/acme/directory"},
	}
	testCtx := context.Background()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dirResp, err := client.Logical().ReadRawWithContext(testCtx, tc.directoryUrl)
			require.NoError(t, err, "failed reading ACME directory configuration")

			require.Equal(t, 200, dirResp.StatusCode)
			require.Equal(t, "application/json", dirResp.Header.Get("Content-Type"))

			requiredUrls := map[string]string{
				"newNonce":   pathConfig + tc.prefixUrl + "/acme/new-nonce",
				"newAccount": pathConfig + tc.prefixUrl + "/acme/new-account",
				"newOrder":   pathConfig + tc.prefixUrl + "/acme/new-order",
				"revokeCert": pathConfig + tc.prefixUrl + "/acme/revoke-cert",
				"keyChange":  pathConfig + tc.prefixUrl + "/acme/key-change",
			}

			rawBodyBytes, err := io.ReadAll(dirResp.Body)
			require.NoError(t, err, "failed reading from directory response body")
			_ = dirResp.Body.Close()

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

// TestAcmeNonce a basic test that will validate we get back a nonce with the proper status codes
// based on the
func TestAcmeNonce(t *testing.T) {
	t.Parallel()
	cluster, client, pathConfig := setupAcmeBackend(t)
	defer cluster.Cleanup()

	cases := []struct {
		name         string
		prefixUrl    string
		directoryUrl string
	}{
		{"root", "", "pki/acme/new-nonce"},
		{"role", "/roles/test-role", "pki/roles/test-role/acme/new-nonce"},
		{"issuer", "/issuer/default", "pki/issuer/default/acme/new-nonce"},
		{"issuer_role", "/issuer/default/roles/test-role", "pki/issuer/default/roles/test-role/acme/new-nonce"},
	}

	for _, tc := range cases {
		for _, httpOp := range []string{"get", "header"} {
			t.Run(fmt.Sprintf("%s-%s", tc.name, httpOp), func(t *testing.T) {
				var req *api.Request
				switch httpOp {
				case "get":
					req = client.NewRequest(http.MethodGet, "/v1/"+tc.directoryUrl)
				case "header":
					req = client.NewRequest(http.MethodHead, "/v1/"+tc.directoryUrl)
				}
				res, err := client.RawRequestWithContext(ctx, req)
				require.NoError(t, err, "failed sending raw request")
				_ = res.Body.Close()

				// Proper Status Code
				switch httpOp {
				case "get":
					require.Equal(t, http.StatusNoContent, res.StatusCode)
				case "header":
					require.Equal(t, http.StatusOK, res.StatusCode)
				}

				// Make sure we don't have a Content-Type header.
				require.Equal(t, "", res.Header.Get("Content-Type"))

				// Make sure we return the Cache-Control header
				require.Contains(t, res.Header.Get("Cache-Control"), "no-store",
					"missing Cache-Control header with no-store header value")

				// Test for our nonce header value
				require.NotEmpty(t, res.Header.Get("Replay-Nonce"), "missing Replay-Nonce header with an actual value")

				// Test Link header value
				expectedLinkHeader := fmt.Sprintf("<%s>;rel=\"index\"", pathConfig+tc.prefixUrl+"/acme/directory")
				require.Contains(t, res.Header.Get("Link"), expectedLinkHeader,
					"different value for link header than expected")
			})
		}
	}
}

// TestAcmeClusterPathNotConfigured basic testing of the ACME error handler.
func TestAcmeClusterPathNotConfigured(t *testing.T) {
	t.Parallel()
	cluster, client := setupTestPkiCluster(t)
	defer cluster.Cleanup()

	// Do not fill in the path option within the local cluster configuration
	cases := []struct {
		name         string
		directoryUrl string
	}{
		{"root", "pki/acme/directory"},
		{"role", "pki/roles/test-role/acme/directory"},
		{"issuer", "pki/issuer/default/acme/directory"},
		{"issuer_role", "pki/issuer/default/roles/test-role/acme/directory"},
	}
	testCtx := context.Background()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dirResp, err := client.Logical().ReadRawWithContext(testCtx, tc.directoryUrl)
			require.Error(t, err, "expected failure reading ACME directory configuration got none")

			require.Equal(t, "application/problem+json", dirResp.Header.Get("Content-Type"))
			require.Equal(t, http.StatusInternalServerError, dirResp.StatusCode)

			rawBodyBytes, err := io.ReadAll(dirResp.Body)
			require.NoError(t, err, "failed reading from directory response body")
			_ = dirResp.Body.Close()

			respType := map[string]interface{}{}
			err = json.Unmarshal(rawBodyBytes, &respType)
			require.NoError(t, err, "failed unmarshalling ACME directory response body")

			require.Equal(t, "urn:ietf:params:acme:error:serverInternal", respType["type"])
			require.NotEmpty(t, respType["detail"])
		})
	}
}

func setupAcmeBackend(t *testing.T) (*vault.TestCluster, *api.Client, string) {
	cluster, client := setupTestPkiCluster(t)

	// Setting templated AIAs should succeed.
	pathConfig := "https://localhost:8200/v1/pki"

	_, err := client.Logical().WriteWithContext(context.Background(), "pki/config/cluster", map[string]interface{}{
		"path":     pathConfig,
		"aia_path": "http://localhost:8200/cdn/pki",
	})
	require.NoError(t, err)

	// Allow certain headers to pass through for ACME support
	_, err = client.Logical().WriteWithContext(context.Background(), "sys/mounts/pki/tune", map[string]interface{}{
		"allowed_response_headers": []string{"Last-Modified", "Replay-Nonce", "Link"},
	})
	require.NoError(t, err, "failed tuning mount response headers")

	return cluster, client, pathConfig
}

func setupTestPkiCluster(t *testing.T) (*vault.TestCluster, *api.Client) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	client := cluster.Cores[0].Client
	mountPKIEndpoint(t, client, "pki")
	return cluster, client
}
