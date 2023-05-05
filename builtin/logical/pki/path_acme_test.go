// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"

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
		{"issuer", "/issuer/int-ca"},
		{"issuer_role", "/issuer/int-ca/roles/test-role"},
	}
	testCtx := context.Background()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseAcmeURL := "/v1/pki" + tc.prefixUrl + "/acme/"
			accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err, "failed creating rsa key")

			acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

			t.Logf("Testing discover on %s", baseAcmeURL)
			discovery, err := acmeClient.Discover(testCtx)
			require.NoError(t, err, "failed acme discovery call")

			discoveryBaseUrl := client.Address() + baseAcmeURL
			require.Equal(t, discoveryBaseUrl+"new-nonce", discovery.NonceURL)
			require.Equal(t, discoveryBaseUrl+"new-account", discovery.RegURL)
			require.Equal(t, discoveryBaseUrl+"new-order", discovery.OrderURL)
			require.Equal(t, discoveryBaseUrl+"revoke-cert", discovery.RevokeURL)
			require.Equal(t, discoveryBaseUrl+"key-change", discovery.KeyChangeURL)
			require.False(t, discovery.ExternalAccountRequired, "bad value for external account required in directory")

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

			// Make sure order's do not accept dates
			_, err = acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{{Type: "dns", Value: "localhost"}},
				acme.WithOrderNotBefore(time.Now().Add(10*time.Minute)))
			require.Error(t, err, "should have rejected a new order with NotBefore set")

			_, err = acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{{Type: "dns", Value: "localhost"}},
				acme.WithOrderNotAfter(time.Now().Add(10*time.Minute)))
			require.Error(t, err, "should have rejected a new order with NotAfter set")

			// Make sure DNS identifiers cannot include IP addresses
			_, err = acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{{Type: "dns", Value: "127.0.0.1"}},
				acme.WithOrderNotAfter(time.Now().Add(10*time.Minute)))
			require.Error(t, err, "should have rejected a new order with IP-like DNS-type identifier")
			_, err = acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{{Type: "dns", Value: "*.127.0.0.1"}},
				acme.WithOrderNotAfter(time.Now().Add(10*time.Minute)))
			require.Error(t, err, "should have rejected a new order with IP-like DNS-type identifier")

			// Create an order
			t.Logf("Testing Authorize Order on %s", baseAcmeURL)
			identifiers := []string{"localhost.localdomain", "*.localdomain"}
			createOrder, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
				{Type: "dns", Value: identifiers[0]},
				{Type: "dns", Value: identifiers[1]},
			})
			require.NoError(t, err, "failed creating order")
			require.Equal(t, acme.StatusPending, createOrder.Status)
			require.Empty(t, createOrder.CertURL)
			require.Equal(t, createOrder.URI+"/finalize", createOrder.FinalizeURL)
			require.Len(t, createOrder.AuthzURLs, 2, "expected two authzurls")

			// Get order
			t.Logf("Testing GetOrder on %s", baseAcmeURL)
			getOrder, err := acmeClient.GetOrder(testCtx, createOrder.URI)
			require.NoError(t, err, "failed fetching order")
			require.Equal(t, acme.StatusPending, createOrder.Status)
			if diffs := deep.Equal(createOrder, getOrder); diffs != nil {
				t.Fatalf("Differences exist between create and get order: \n%v", strings.Join(diffs, "\n"))
			}

			// Make sure the identifiers returned in the order contain the original values
			var ids []string
			for _, id := range getOrder.Identifiers {
				require.Equal(t, "dns", id.Type)
				ids = append(ids, id.Value)
			}
			require.ElementsMatch(t, identifiers, ids, "order responses should have all original identifiers")

			// Load authorizations
			var authorizations []*acme.Authorization
			for _, authUrl := range getOrder.AuthzURLs {
				auth, err := acmeClient.GetAuthorization(testCtx, authUrl)
				require.NoError(t, err, "failed fetching authorization: %s", authUrl)

				authorizations = append(authorizations, auth)
			}

			// We should have 2 separate auth challenges as we have two separate identifier
			require.Len(t, authorizations, 2, "expected 2 authorizations in order")

			var wildcardAuth *acme.Authorization
			var domainAuth *acme.Authorization
			for _, auth := range authorizations {
				if auth.Wildcard {
					wildcardAuth = auth
				} else {
					domainAuth = auth
				}
			}

			// Test the values for the domain authentication
			require.Equal(t, acme.StatusPending, domainAuth.Status)
			require.Equal(t, "dns", domainAuth.Identifier.Type)
			require.Equal(t, "localhost.localdomain", domainAuth.Identifier.Value)
			require.False(t, domainAuth.Wildcard, "should not be a wildcard")
			require.True(t, domainAuth.Expires.IsZero(), "authorization should only have expiry set on valid status")

			require.Len(t, domainAuth.Challenges, 2, "expected two challenges")
			require.Equal(t, acme.StatusPending, domainAuth.Challenges[0].Status)
			require.True(t, domainAuth.Challenges[0].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "http-01", domainAuth.Challenges[0].Type)
			require.NotEmpty(t, domainAuth.Challenges[0].Token, "missing challenge token")
			require.Equal(t, acme.StatusPending, domainAuth.Challenges[1].Status)
			require.True(t, domainAuth.Challenges[1].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "dns-01", domainAuth.Challenges[1].Type)
			require.NotEmpty(t, domainAuth.Challenges[1].Token, "missing challenge token")

			// Test the values for the wilcard authentication
			require.Equal(t, acme.StatusPending, wildcardAuth.Status)
			require.Equal(t, "dns", wildcardAuth.Identifier.Type)
			require.Equal(t, "localdomain", wildcardAuth.Identifier.Value) // Make sure we strip the *. in auth responses
			require.True(t, wildcardAuth.Wildcard, "should be a wildcard")
			require.True(t, wildcardAuth.Expires.IsZero(), "authorization should only have expiry set on valid status")

			require.Len(t, wildcardAuth.Challenges, 1, "expected two challenges")
			require.Equal(t, acme.StatusPending, domainAuth.Challenges[0].Status)
			require.True(t, wildcardAuth.Challenges[0].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "dns-01", wildcardAuth.Challenges[0].Type)
			require.NotEmpty(t, domainAuth.Challenges[0].Token, "missing challenge token")

			// Load a challenge directly; this triggers validation to start.
			challenge, err := acmeClient.GetChallenge(testCtx, domainAuth.Challenges[0].URI)
			require.NoError(t, err, "failed to load challenge")
			require.Equal(t, acme.StatusProcessing, challenge.Status)
			require.True(t, challenge.Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "http-01", challenge.Type)

			require.NotEmpty(t, challenge.Token, "missing challenge token")

			// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow
			//       test.
			pkiMount := findStorageMountUuid(t, client, "pki")
			accountId := acct.URI[strings.LastIndex(acct.URI, "/"):]
			for _, authURI := range getOrder.AuthzURLs {
				authId := authURI[strings.LastIndex(authURI, "/"):]

				rawPath := path.Join("/sys/raw/logical/", pkiMount, getAuthorizationPath(accountId, authId))
				resp, err := client.Logical().ReadWithContext(testCtx, rawPath)
				require.NoError(t, err, "failed looking up authorization storage")
				require.NotNil(t, resp, "sys raw response was nil")
				require.NotEmpty(t, resp.Data["value"], "no value field in sys raw response")

				var authz ACMEAuthorization
				err = jsonutil.DecodeJSON([]byte(resp.Data["value"].(string)), &authz)
				require.NoError(t, err, "error decoding authorization: %w", err)
				authz.Status = ACMEAuthorizationValid
				for _, challenge := range authz.Challenges {
					challenge.Status = ACMEChallengeValid
				}

				encodeJSON, err := jsonutil.EncodeJSON(authz)
				require.NoError(t, err, "failed encoding authz json")
				_, err = client.Logical().WriteWithContext(testCtx, rawPath, map[string]interface{}{
					"value":    base64.StdEncoding.EncodeToString(encodeJSON),
					"encoding": "base64",
				})
				require.NoError(t, err, "failed writing authorization storage")
			}

			// Make sure sending a CSR with the account key gets rejected.
			goodCr := &x509.CertificateRequest{
				Subject:  pkix.Name{CommonName: identifiers[1]},
				DNSNames: []string{identifiers[0], identifiers[1]},
			}
			t.Logf("csr: %v", goodCr)

			// We want to make sure people are not using the same keys for CSR/Certs and their ACME account.
			csrSignedWithAccountKey, err := x509.CreateCertificateRequest(rand.Reader, goodCr, accountKey)
			require.NoError(t, err, "failed generating csr")
			_, _, err = acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csrSignedWithAccountKey, true)
			require.Error(t, err, "should not be allowed to use the account key for a CSR")

			csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			require.NoError(t, err, "failed generated key for CSR")

			// Validate we reject CSRs that contain names that aren't in the original order
			badCr := &x509.CertificateRequest{
				Subject:  pkix.Name{CommonName: createOrder.Identifiers[0].Value},
				DNSNames: []string{"www.notinorder.com"},
			}

			csrWithBadName, err := x509.CreateCertificateRequest(rand.Reader, badCr, csrKey)
			require.NoError(t, err, "failed generating csr with bad name")

			_, _, err = acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csrWithBadName, true)
			require.Error(t, err, "should not be allowed to csr with different names than order")

			// Validate we reject CSRs that contains fewer names than in the original order.
			badCr = &x509.CertificateRequest{
				Subject: pkix.Name{CommonName: identifiers[0]},
			}

			csrWithBadName, err = x509.CreateCertificateRequest(rand.Reader, badCr, csrKey)
			require.NoError(t, err, "failed generating csr with bad name")

			_, _, err = acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csrWithBadName, true)
			require.Error(t, err, "should not be allowed to csr with different names than order")

			// Finally test a proper CSR, with the correct name and signed with a different key works.
			csr, err := x509.CreateCertificateRequest(rand.Reader, goodCr, csrKey)
			require.NoError(t, err, "failed generating csr")

			certs, _, err := acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csr, true)
			require.NoError(t, err, "failed finalizing order")
			require.Len(t, certs, 3, "expected three items within the returned certs")

			testAcmeCertSignedByCa(t, client, certs, "int-ca")

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

// TestAcmeBasicWorkflowWithEab verify that new accounts require EAB's if enforced by configuration.
func TestAcmeBasicWorkflowWithEab(t *testing.T) {
	t.Parallel()
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()
	testCtx := context.Background()

	// Enable EAB
	_, err := client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled":    true,
		"eab_policy": "always-required",
	})
	require.NoError(t, err)

	baseAcmeURL := "/v1/pki/acme/"
	accountKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed creating ec key")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	t.Logf("Testing discover on %s", baseAcmeURL)
	discovery, err := acmeClient.Discover(testCtx)
	require.NoError(t, err, "failed acme discovery call")
	require.True(t, discovery.ExternalAccountRequired, "bad value for external account required in directory")

	// Create new account without EAB, should fail
	t.Logf("Testing register on %s", baseAcmeURL)
	_, err = acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.ErrorContains(t, err, "urn:ietf:params:acme:error:externalAccountRequired",
		"expected failure creating an account without eab")

	kid, eabKeyBytes := getEABKey(t, client)
	acct := &acme.Account{
		ExternalAccountBinding: &acme.ExternalAccountBinding{
			KID: kid,
			Key: eabKeyBytes,
		},
	}

	// Create new account with EAB
	t.Logf("Testing register on %s", baseAcmeURL)
	_, err = acmeClient.Register(testCtx, acct, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering new account with eab")

	// Make sure our EAB is no longer available
	resp, err := client.Logical().ListWithContext(context.Background(), "pki/acme/eab")
	require.NoError(t, err, "failed to list eab tokens")
	require.Nil(t, resp, "list response for eab tokens should have been nil due to empty list")

	// Attempt to create another account with the same EAB as before -- should fail
	accountKey2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed creating ec key")

	acmeClient2 := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey2)
	acct2 := &acme.Account{
		ExternalAccountBinding: &acme.ExternalAccountBinding{
			KID: kid,
			Key: eabKeyBytes,
		},
	}

	_, err = acmeClient2.Register(testCtx, acct2, func(tosURL string) bool { return true })
	require.ErrorContains(t, err, "urn:ietf:params:acme:error:unauthorized", "should fail due to EAB re-use")

	// We can lookup/find an existing account without EAB if we have the account key
	_, err = acmeClient.GetReg(testCtx /* unused url */, "")
	require.NoError(t, err, "expected to lookup existing account without eab")
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

	// Enable ACME but don't set a path.
	_, err := client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled": true,
	})
	require.NoError(t, err)

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

// TestAcmeAccountsCrossingDirectoryPath make sure that if an account attempts to use a different ACME
// directory path that we get an error.
func TestAcmeAccountsCrossingDirectoryPath(t *testing.T) {
	t.Parallel()
	cluster, _, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	baseAcmeURL := "/v1/pki/acme/"
	accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	testCtx := context.Background()
	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// Create new account
	acct, err := acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Try to update the account under another ACME directory
	baseAcmeURL2 := "/v1/pki/roles/test-role/acme/"
	acmeClient2 := getAcmeClientForCluster(t, cluster, baseAcmeURL2, accountKey)
	acct.Contact = []string{"mailto:test3@example.com"}
	_, err = acmeClient2.UpdateReg(testCtx, acct)
	require.Error(t, err, "successfully updated account when we should have failed due to different directory")
	// We don't test for the specific error about using the wrong directory, as the golang library
	// swallows the error we are sending back to a no account error
}

// TestAcmeDisabledWithEnvVar verifies if VAULT_DISABLE_PUBLIC_ACME is set that we completely
// disable the ACME service
func TestAcmeDisabledWithEnvVar(t *testing.T) {
	t.Setenv("VAULT_DISABLE_PUBLIC_ACME", "true")

	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	for _, method := range []string{http.MethodHead, http.MethodGet} {
		t.Run(fmt.Sprintf("%s", method), func(t *testing.T) {
			req := client.NewRequest(method, "/v1/pki/acme/new-nonce")
			_, err := client.RawRequestWithContext(ctx, req)
			require.Error(t, err, "should have received an error as ACME should have been disabled")

			if apiError, ok := err.(*api.ResponseError); ok {
				require.Equal(t, 404, apiError.StatusCode)
			}
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

	_, err = client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled":    true,
		"eab_policy": "not-required",
	})
	require.NoError(t, err)

	// Allow certain headers to pass through for ACME support
	_, err = client.Logical().WriteWithContext(context.Background(), "sys/mounts/pki/tune", map[string]interface{}{
		"allowed_response_headers": []string{"Last-Modified", "Replay-Nonce", "Link", "Location"},
		"max_lease_ttl":            "920000h",
	})
	require.NoError(t, err, "failed tuning mount response headers")

	resp, err := client.Logical().WriteWithContext(context.Background(), "/pki/issuers/generate/root/internal",
		map[string]interface{}{
			"issuer_name": "root-ca",
			"key_name":    "root-key",
			"key_type":    "ec",
			"common_name": "root.com",
			"ttl":         "7200h",
			"max_ttl":     "920000h",
		})
	require.NoError(t, err, "failed creating root CA")

	resp, err = client.Logical().WriteWithContext(context.Background(), "/pki/issuers/generate/intermediate/internal",
		map[string]interface{}{
			"key_name":    "int-key",
			"key_type":    "ec",
			"common_name": "test.com",
		})
	require.NoError(t, err, "failed creating intermediary CSR")
	intermediateCSR := resp.Data["csr"].(string)

	// Sign the intermediate CSR using /pki
	resp, err = client.Logical().Write("pki/issuer/root-ca/sign-intermediate", map[string]interface{}{
		"csr":     intermediateCSR,
		"ttl":     "720h",
		"max_ttl": "7200h",
	})
	require.NoError(t, err, "failed signing intermediary CSR")
	intermediateCertPEM := resp.Data["certificate"].(string)

	// Configure the intermediate cert as the CA in /pki2
	resp, err = client.Logical().Write("/pki/issuers/import/cert", map[string]interface{}{
		"pem_bundle": intermediateCertPEM,
	})
	require.NoError(t, err, "failed importing intermediary cert")
	importedIssuersRaw := resp.Data["imported_issuers"].([]interface{})
	require.Len(t, importedIssuersRaw, 1)
	intCaUuid := importedIssuersRaw[0].(string)

	_, err = client.Logical().Write("/pki/issuer/"+intCaUuid, map[string]interface{}{
		"issuer_name": "int-ca",
	})
	require.NoError(t, err, "failed updating issuer name")

	_, err = client.Logical().Write("/pki/config/issuers", map[string]interface{}{
		"default": "int-ca",
	})
	require.NoError(t, err, "failed updating default issuer")

	_, err = client.Logical().Write("/pki/roles/test-role", map[string]interface{}{
		"ttl_duration":                "365h",
		"max_ttl_duration":            "720h",
		"key_type":                    "any",
		"allowed_domains":             "localdomain",
		"allow_subdomains":            "true",
		"allow_wildcard_certificates": "true",
	})
	require.NoError(t, err, "failed creating role test-role")

	_, err = client.Logical().Write("/pki/roles/acme", map[string]interface{}{
		"ttl_duration":     "365h",
		"max_ttl_duration": "720h",
		"key_type":         "any",
	})
	require.NoError(t, err, "failed creating role acme")

	return cluster, client, pathConfig
}

func testAcmeCertSignedByCa(t *testing.T, client *api.Client, derCerts [][]byte, issuerRef string) {
	t.Helper()
	require.NotEmpty(t, derCerts)
	acmeCert, err := x509.ParseCertificate(derCerts[0])
	require.NoError(t, err, "failed parsing acme cert bytes")

	resp, err := client.Logical().ReadWithContext(context.Background(), "pki/issuer/"+issuerRef)
	require.NoError(t, err, "failed reading issuer with name %s", issuerRef)
	issuerCert := parseCert(t, resp.Data["certificate"].(string))
	issuerChainRaw := resp.Data["ca_chain"].([]interface{})

	err = acmeCert.CheckSignatureFrom(issuerCert)
	require.NoError(t, err, "issuer %s did not sign provided cert", issuerRef)

	expectedCerts := [][]byte{derCerts[0]}

	for _, entry := range issuerChainRaw {
		chainCert := parseCert(t, entry.(string))
		expectedCerts = append(expectedCerts, chainCert.Raw)
	}

	if diffs := deep.Equal(expectedCerts, derCerts); diffs != nil {
		t.Fatalf("diffs were found between the acme chain returned and the expected value: \n%v", diffs)
	}
}

func setupTestPkiCluster(t *testing.T) (*vault.TestCluster, *api.Client) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
		EnableRaw: true,
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

func getEABKey(t *testing.T, client *api.Client) (string, []byte) {
	resp, err := client.Logical().WriteWithContext(ctx, "pki/acme/eab", map[string]interface{}{})
	require.NoError(t, err, "failed getting eab key")
	require.NotNil(t, resp, "eab key returned nil response")
	require.NotEmpty(t, resp.Data["id"], "eab key response missing id field")
	kid := resp.Data["id"].(string)

	require.NotEmpty(t, resp.Data["private_key"], "eab key response missing private_key field")
	base64Key := resp.Data["private_key"].(string)
	privateKeyBytes, err := base64.RawURLEncoding.DecodeString(base64Key)
	require.NoError(t, err, "failed base 64 decoding eab key response")

	return kid, privateKeyBytes
}
