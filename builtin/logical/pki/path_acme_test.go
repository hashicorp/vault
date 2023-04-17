// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/crypto/acme"
	"golang.org/x/net/http2"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2/json"
)

// TestAcmeBasicWorkflow a basic test that will validate a basic ACME workflow using the Golang ACME client.
func TestAcmeBasicWorkflow(t *testing.T) {
	t.Parallel()
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()
	cases := []struct {
		name      string
		prefixUrl string
	}{
		{"root", ""},
		{"role", "/roles/test-role"},
		{"issuer", "/issuer/default"},
		{"issuer_role", "/issuer/default/roles/test-role"},
		{"issuer_role_acme", "/issuer/acme/roles/acme"},
	}
	testCtx := context.Background()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseAcmeURL := "/v1/pki" + tc.prefixUrl + "/acme/"
			key, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err, "failed creating rsa key")

			acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, key)

			t.Logf("Testing discover on %s", baseAcmeURL)
			discovery, err := acmeClient.Discover(testCtx)
			require.NoError(t, err, "failed acme discovery call")

			discoveryBaseUrl := client.Address() + baseAcmeURL
			require.Equal(t, discoveryBaseUrl+"new-nonce", discovery.NonceURL)
			require.Equal(t, discoveryBaseUrl+"new-account", discovery.RegURL)
			require.Equal(t, discoveryBaseUrl+"new-order", discovery.OrderURL)
			require.Equal(t, discoveryBaseUrl+"revoke-cert", discovery.RevokeURL)
			require.Equal(t, discoveryBaseUrl+"key-change", discovery.KeyChangeURL)

			// Attempt to update prior to creating an account
			t.Logf("Testing updates with no proper account fail on %s", baseAcmeURL)
			_, err = acmeClient.UpdateReg(testCtx, &acme.Account{Contact: []string{"mailto:shouldfail@example.com"}})
			require.ErrorIs(t, err, acme.ErrNoAccount, "expected failure attempting to update prior to account registration")

			// Create new account
			t.Logf("Testing register on %s", baseAcmeURL)
			acct, err := acmeClient.Register(testCtx, &acme.Account{
				Contact: []string{"mailto:test@example.com", "mailto:test2@test.com"},
			}, func(tosURL string) bool { return true })
			require.NoError(t, err, "failed registering account")
			require.Equal(t, acme.StatusValid, acct.Status)
			require.Contains(t, acct.Contact, "mailto:test@example.com")
			require.Contains(t, acct.Contact, "mailto:test2@test.com")
			require.Len(t, acct.Contact, 2)

			// Call register again we should get existing account
			t.Logf("Testing duplicate register returns existing account on %s", baseAcmeURL)
			_, err = acmeClient.Register(testCtx, acct, func(tosURL string) bool { return true })
			require.ErrorIs(t, err, acme.ErrAccountAlreadyExists,
				"We should have returned a 200 status code which would have triggered an error in the golang acme"+
					" library")

			// Update contact
			t.Logf("Testing Update account contacts on %s", baseAcmeURL)
			acct.Contact = []string{"mailto:test3@example.com"}
			acct2, err := acmeClient.UpdateReg(testCtx, acct)
			require.NoError(t, err, "failed updating account")
			require.Equal(t, acme.StatusValid, acct2.Status)
			// We should get this back, not the original values.
			require.Contains(t, acct2.Contact, "mailto:test3@example.com")
			require.Len(t, acct2.Contact, 1)

			// Create an order
			t.Logf("Testing Authorize Order on %s", baseAcmeURL)
			createOrder, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{{Type: "dns", Value: "www.test.com"}},
				acme.WithOrderNotBefore(time.Now().Add(10*time.Minute)),
				acme.WithOrderNotAfter(time.Now().Add(7*24*time.Hour)))
			require.NoError(t, err, "failed creating order")
			require.Equal(t, acme.StatusPending, createOrder.Status)
			require.Empty(t, createOrder.CertURL)
			require.Equal(t, createOrder.URI+"/finalize", createOrder.FinalizeURL)
			require.Len(t, createOrder.AuthzURLs, 1, "expected one authzurls")

			// Get order
			t.Logf("Testing GetOrder on %s", baseAcmeURL)
			getOrder, err := acmeClient.GetOrder(testCtx, createOrder.URI)
			require.NoError(t, err, "failed fetching order")
			require.Equal(t, acme.StatusPending, createOrder.Status)
			if diffs := deep.Equal(createOrder, getOrder); diffs != nil {
				t.Fatalf("Differences exist between create and get order: \n%v", strings.Join(diffs, "\n"))
			}

			// Load authorization
			auth, err := acmeClient.GetAuthorization(testCtx, getOrder.AuthzURLs[0])
			require.NoError(t, err, "failed fetching authorization")
			require.Equal(t, acme.StatusPending, auth.Status)
			require.Equal(t, "dns", auth.Identifier.Type)
			require.Equal(t, "www.test.com", auth.Identifier.Value)
			require.False(t, auth.Wildcard, "should not be a wildcard")
			require.True(t, auth.Expires.IsZero(), "authorization should only have expiry set on valid status")

			require.Len(t, auth.Challenges, 1, "expected one challenge")
			require.Equal(t, acme.StatusPending, auth.Challenges[0].Status)
			require.True(t, auth.Challenges[0].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "http-01", auth.Challenges[0].Type)

			// TODO: This currently does fail
			// require.NotEmpty(t, auth.Challenges[0].Token, "missing challenge token")

			// Load a challenge directly
			challenge, err := acmeClient.GetChallenge(testCtx, auth.Challenges[0].URI)
			require.NoError(t, err, "failed to load challenge")
			require.Equal(t, acme.StatusPending, challenge.Status)
			require.True(t, challenge.Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "http-01", challenge.Type)

			// TODO: This currently does fail
			// require.NotEmpty(t, challenge.Token, "missing challenge token")

			// Deactivate account
			t.Logf("Testing deactivate account on %s", baseAcmeURL)
			err = acmeClient.DeactivateReg(testCtx)
			require.NoError(t, err, "failed deactivating account")

			// Make sure we get an unauthorized error trying to update the account again.
			t.Logf("Testing update on deactivated account fails on %s", baseAcmeURL)
			_, err = acmeClient.UpdateReg(testCtx, acct)
			require.Error(t, err, "expected account to be deactivated")
			require.IsType(t, &acme.Error{}, err, "expected acme error type")
			acmeErr := err.(*acme.Error)
			require.Equal(t, "urn:ietf:params:acme:error:unauthorized", acmeErr.ProblemType)
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
	pathConfig := client.Address() + "/v1/pki"

	_, err := client.Logical().WriteWithContext(context.Background(), "pki/config/cluster", map[string]interface{}{
		"path":     pathConfig,
		"aia_path": "http://localhost:8200/cdn/pki",
	})
	require.NoError(t, err)

	// Allow certain headers to pass through for ACME support
	_, err = client.Logical().WriteWithContext(context.Background(), "sys/mounts/pki/tune", map[string]interface{}{
		"allowed_response_headers": []string{"Last-Modified", "Replay-Nonce", "Link", "Location"},
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

func getAcmeClientForCluster(t *testing.T, cluster *vault.TestCluster, baseUrl string, key crypto.Signer) acme.Client {
	coreAddr := cluster.Cores[0].Listeners[0].Address
	tlsConfig := cluster.Cores[0].TLSConfig()

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = tlsConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	httpClient := &http.Client{Transport: transport}
	if baseUrl[0] == '/' {
		baseUrl = baseUrl[1:]
	}
	if !strings.HasPrefix(baseUrl, "v1/") {
		baseUrl = "v1/" + baseUrl
	}
	baseAcmeURL := fmt.Sprintf("https://%s/%s", coreAddr.String(), baseUrl)
	return acme.Client{
		Key:          key,
		HTTPClient:   httpClient,
		DirectoryURL: baseAcmeURL + "directory",
	}
}
