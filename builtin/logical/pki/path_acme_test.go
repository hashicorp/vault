// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki/dnstest"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/acme"
	"golang.org/x/net/http2"
)

// TestAcmeBasicWorkflow a test that will validate a basic ACME workflow using the Golang ACME client.
func TestAcmeBasicWorkflow(t *testing.T) {
	t.Parallel()
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()
	cases := []struct {
		name      string
		prefixUrl string
	}{
		{"root", "acme/"},
		{"role", "roles/test-role/acme/"},
		{"issuer", "issuer/int-ca/acme/"},
		{"issuer_role", "issuer/int-ca/roles/test-role/acme/"},
	}
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseAcmeURL := "/v1/pki/" + tc.prefixUrl
			accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
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

			require.Len(t, domainAuth.Challenges, 3, "expected three challenges")
			require.Equal(t, acme.StatusPending, domainAuth.Challenges[0].Status)
			require.True(t, domainAuth.Challenges[0].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "http-01", domainAuth.Challenges[0].Type)
			require.NotEmpty(t, domainAuth.Challenges[0].Token, "missing challenge token")
			require.Equal(t, acme.StatusPending, domainAuth.Challenges[1].Status)
			require.True(t, domainAuth.Challenges[1].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "dns-01", domainAuth.Challenges[1].Type)
			require.NotEmpty(t, domainAuth.Challenges[1].Token, "missing challenge token")
			require.Equal(t, acme.StatusPending, domainAuth.Challenges[2].Status)
			require.True(t, domainAuth.Challenges[2].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "tls-alpn-01", domainAuth.Challenges[2].Type)
			require.NotEmpty(t, domainAuth.Challenges[2].Token, "missing challenge token")

			// Test the values for the wildcard authentication
			require.Equal(t, acme.StatusPending, wildcardAuth.Status)
			require.Equal(t, "dns", wildcardAuth.Identifier.Type)
			require.Equal(t, "localdomain", wildcardAuth.Identifier.Value) // Make sure we strip the *. in auth responses
			require.True(t, wildcardAuth.Wildcard, "should be a wildcard")
			require.True(t, wildcardAuth.Expires.IsZero(), "authorization should only have expiry set on valid status")

			require.Len(t, wildcardAuth.Challenges, 1, "expected one challenge")
			require.Equal(t, acme.StatusPending, domainAuth.Challenges[0].Status)
			require.True(t, wildcardAuth.Challenges[0].Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "dns-01", wildcardAuth.Challenges[0].Type)
			require.NotEmpty(t, domainAuth.Challenges[0].Token, "missing challenge token")

			// Make sure that getting a challenge does not start it.
			challenge, err := acmeClient.GetChallenge(testCtx, domainAuth.Challenges[0].URI)
			require.NoError(t, err, "failed to load challenge")
			require.Equal(t, acme.StatusPending, challenge.Status)
			require.True(t, challenge.Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "http-01", challenge.Type)

			// Accept a challenge; this triggers validation to start.
			challenge, err = acmeClient.Accept(testCtx, domainAuth.Challenges[0])
			require.NoError(t, err, "failed to load challenge")
			require.Equal(t, acme.StatusProcessing, challenge.Status)
			require.True(t, challenge.Validated.IsZero(), "validated time should be 0 on challenge")
			require.Equal(t, "http-01", challenge.Type)

			require.NotEmpty(t, challenge.Token, "missing challenge token")

			// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow
			//       test.
			markAuthorizationSuccess(t, client, acmeClient, acct, getOrder)

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

			// Validate we reject CSRs that contain CN that aren't in the original order
			badCr := &x509.CertificateRequest{
				Subject:  pkix.Name{CommonName: "not-in-original-order.com"},
				DNSNames: []string{identifiers[0], identifiers[1]},
			}
			t.Logf("csr: %v", badCr)

			csrWithBadCName, err := x509.CreateCertificateRequest(rand.Reader, badCr, csrKey)
			require.NoError(t, err, "failed generating csr with bad common name")

			_, _, err = acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csrWithBadCName, true)
			require.Error(t, err, "should not be allowed to csr with different common names than order")

			// Validate we reject CSRs that contain DNS names that aren't in the original order
			badCr = &x509.CertificateRequest{
				Subject:  pkix.Name{CommonName: createOrder.Identifiers[0].Value},
				DNSNames: []string{"www.notinorder.com"},
			}

			csrWithBadName, err := x509.CreateCertificateRequest(rand.Reader, badCr, csrKey)
			require.NoError(t, err, "failed generating csr with bad name")

			_, _, err = acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csrWithBadName, true)
			require.Error(t, err, "should not be allowed to csr with different names than order")

			// Validate we reject CSRs that contain IP addresses that weren't in the original order
			badCr = &x509.CertificateRequest{
				Subject:     pkix.Name{CommonName: createOrder.Identifiers[0].Value},
				IPAddresses: []net.IP{{127, 0, 0, 1}},
			}

			csrWithBadIP, err := x509.CreateCertificateRequest(rand.Reader, badCr, csrKey)
			require.NoError(t, err, "failed generating csr with bad name")

			_, _, err = acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csrWithBadIP, true)
			require.Error(t, err, "should not be allowed to csr with different ip address than order")

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

			// Make sure the certificate has a NotAfter date of a maximum of 90 days
			acmeCert, err := x509.ParseCertificate(certs[0])
			require.NoError(t, err, "failed parsing acme cert bytes")
			maxAcmeNotAfter := time.Now().Add(defaultAcmeMaxTTL)
			if maxAcmeNotAfter.Before(acmeCert.NotAfter) {
				require.Fail(t, fmt.Sprintf("certificate has a NotAfter value %v greater than ACME max ttl %v", acmeCert.NotAfter, maxAcmeNotAfter))
			}

			// Can we revoke it using the account key revocation
			err = acmeClient.RevokeCert(ctx, nil, certs[0], acme.CRLReasonUnspecified)
			require.NoError(t, err, "failed to revoke certificate through account key")

			// Make sure it was actually revoked
			certResp, err := client.Logical().ReadWithContext(ctx, "pki/cert/"+serialFromCert(acmeCert))
			require.NoError(t, err, "failed to read certificate status")
			require.NotNil(t, certResp, "certificate status response was nil")
			revocationTime := certResp.Data["revocation_time"].(json.Number)
			revocationTimeInt, err := revocationTime.Int64()
			require.NoError(t, err, "failed converting revocation_time value: %v", revocationTime)
			require.Greater(t, revocationTimeInt, int64(0),
				"revocation time was not greater than 0, revocation did not work value was: %v", revocationTimeInt)

			// Make sure we can revoke an authorization as a client
			err = acmeClient.RevokeAuthorization(ctx, authorizations[0].URI)
			require.NoError(t, err, "failed revoking authorization status")

			revokedAuth, err := acmeClient.GetAuthorization(ctx, authorizations[0].URI)
			require.NoError(t, err, "failed fetching authorization")
			require.Equal(t, acme.StatusDeactivated, revokedAuth.Status)

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
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Enable EAB
	_, err := client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled":    true,
		"eab_policy": "always-required",
	})
	require.NoError(t, err)

	cases := []struct {
		name      string
		prefixUrl string
	}{
		{"root", "acme/"},
		{"role", "roles/test-role/acme/"},
		{"issuer", "issuer/int-ca/acme/"},
		{"issuer_role", "issuer/int-ca/roles/test-role/acme/"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseAcmeURL := "/v1/pki/" + tc.prefixUrl
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

			// Test fetch, list, delete workflow
			kid, _ := getEABKey(t, client, tc.prefixUrl)
			resp, err := client.Logical().ListWithContext(testCtx, "pki/eab")
			require.NoError(t, err, "failed to list eab tokens")
			require.NotNil(t, resp, "list response for eab tokens should not be nil")
			require.Contains(t, resp.Data, "keys")
			require.Contains(t, resp.Data, "key_info")
			require.Len(t, resp.Data["keys"], 1)
			require.Contains(t, resp.Data["keys"], kid)

			_, err = client.Logical().DeleteWithContext(testCtx, "pki/eab/"+kid)
			require.NoError(t, err, "failed to delete eab")

			// List eabs should return zero results
			resp, err = client.Logical().ListWithContext(testCtx, "pki/eab")
			require.NoError(t, err, "failed to list eab tokens")
			require.Nil(t, resp, "list response for eab tokens should have been nil")

			// fetch a new EAB
			kid, eabKeyBytes := getEABKey(t, client, tc.prefixUrl)
			acct := &acme.Account{
				ExternalAccountBinding: &acme.ExternalAccountBinding{
					KID: kid,
					Key: eabKeyBytes,
				},
			}

			// Make sure we can list our key
			resp, err = client.Logical().ListWithContext(testCtx, "pki/eab")
			require.NoError(t, err, "failed to list eab tokens")
			require.NotNil(t, resp, "list response for eab tokens should not be nil")
			require.Contains(t, resp.Data, "keys")
			require.Contains(t, resp.Data, "key_info")
			require.Len(t, resp.Data["keys"], 1)
			require.Contains(t, resp.Data["keys"], kid)

			keyInfo := resp.Data["key_info"].(map[string]interface{})
			require.Contains(t, keyInfo, kid)

			infoForKid := keyInfo[kid].(map[string]interface{})
			require.Equal(t, "hs", infoForKid["key_type"])
			require.Equal(t, tc.prefixUrl+"directory", infoForKid["acme_directory"])

			// Create new account with EAB
			t.Logf("Testing register on %s", baseAcmeURL)
			_, err = acmeClient.Register(testCtx, acct, func(tosURL string) bool { return true })
			require.NoError(t, err, "failed registering new account with eab")

			// Make sure our EAB is no longer available
			resp, err = client.Logical().ListWithContext(context.Background(), "pki/eab")
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

	// Go sneaky, sneaky and update the acme configuration through sys/raw to bypass config/cluster path checks
	pkiMount := findStorageMountUuid(t, client, "pki")
	rawPath := path.Join("/sys/raw/logical/", pkiMount, storageAcmeConfig)
	_, err := client.Logical().WriteWithContext(context.Background(), rawPath, map[string]interface{}{
		"value": "{\"enabled\": true, \"eab_policy_name\": \"not-required\"}",
	})
	require.NoError(t, err, "failed updating acme config through sys/raw")

	// Force reload the plugin so we read the new config we slipped in.
	_, err = client.Sys().ReloadPluginWithContext(context.Background(), &api.ReloadPluginInput{Mounts: []string{"pki"}})
	require.NoError(t, err, "failed reloading plugin")

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
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

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
	accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
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

// TestAcmeEabCrossingDirectoryPath make sure that if an account attempts to use a different ACME
// directory path that an EAB was created within we get an error.
func TestAcmeEabCrossingDirectoryPath(t *testing.T) {
	t.Parallel()
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	// Enable EAB
	_, err := client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled":    true,
		"eab_policy": "always-required",
	})
	require.NoError(t, err)

	baseAcmeURL := "/v1/pki/acme/"
	accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// fetch a new EAB
	kid, eabKeyBytes := getEABKey(t, client, "roles/test-role/acme/")
	acct := &acme.Account{
		ExternalAccountBinding: &acme.ExternalAccountBinding{
			KID: kid,
			Key: eabKeyBytes,
		},
	}

	// Create new account
	_, err = acmeClient.Register(testCtx, acct, func(tosURL string) bool { return true })
	require.ErrorContains(t, err, "failed to verify eab", "should have failed as EAB is for a different directory")
}

// TestAcmeDisabledWithEnvVar verifies if VAULT_DISABLE_PUBLIC_ACME is set that we completely
// disable the ACME service
func TestAcmeDisabledWithEnvVar(t *testing.T) {
	// Setup a cluster with the configuration set to not-required, initially as the
	// configuration will validate if the environment var is set
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	// Seal setup the environment variable, and unseal which now means we have a cluster
	// with ACME configuration saying it is enabled with a bad EAB policy.
	cluster.EnsureCoresSealed(t)
	t.Setenv("VAULT_DISABLE_PUBLIC_ACME", "true")
	cluster.UnsealCores(t)

	// Make sure that ACME is disabled now.
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

// TestAcmeConfigChecksPublicAcmeEnv verifies certain EAB policy values can not be set if ENV var is enabled
func TestAcmeConfigChecksPublicAcmeEnv(t *testing.T) {
	t.Setenv("VAULT_DISABLE_PUBLIC_ACME", "true")
	cluster, client := setupTestPkiCluster(t)
	defer cluster.Cleanup()

	_, err := client.Logical().WriteWithContext(context.Background(), "pki/config/cluster", map[string]interface{}{
		"path": "https://dadgarcorp.com/v1/pki",
	})
	require.NoError(t, err)

	_, err = client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled":    true,
		"eab_policy": string(eabPolicyAlwaysRequired),
	})
	require.NoError(t, err)

	for _, policyName := range []EabPolicyName{eabPolicyNewAccountRequired, eabPolicyNotRequired} {
		_, err = client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
			"enabled":    true,
			"eab_policy": string(policyName),
		})
		require.Error(t, err, "eab policy %s should have not been allowed to be set")
	}

	// Make sure we can disable ACME and the eab policy is not checked
	_, err = client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled":    false,
		"eab_policy": string(eabPolicyNotRequired),
	})
	require.NoError(t, err)
}

// TestAcmeHonorsAlwaysEnforceErr verifies that we get an error and not truncated if the issuer's
// leaf_not_after_behavior is set to always_enforce_err
func TestAcmeHonorsAlwaysEnforceErr(t *testing.T) {
	t.Parallel()

	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	mount := "pki"
	resp, err := client.Logical().WriteWithContext(context.Background(), mount+"/issuers/generate/intermediate/internal",
		map[string]interface{}{
			"key_name":    "short-key",
			"key_type":    "ec",
			"common_name": "test.com",
		})
	require.NoError(t, err, "failed creating intermediary CSR")
	intermediateCSR := resp.Data["csr"].(string)

	// Sign the intermediate CSR using /pki
	resp, err = client.Logical().Write(mount+"/issuer/root-ca/sign-intermediate", map[string]interface{}{
		"csr":     intermediateCSR,
		"ttl":     "10m",
		"max_ttl": "1h",
	})
	require.NoError(t, err, "failed signing intermediary CSR")
	intermediateCertPEM := resp.Data["certificate"].(string)

	// Configure the intermediate cert as the CA in /pki2
	resp, err = client.Logical().Write(mount+"/issuers/import/cert", map[string]interface{}{
		"pem_bundle": intermediateCertPEM,
	})
	require.NoError(t, err, "failed importing intermediary cert")
	importedIssuersRaw := resp.Data["imported_issuers"].([]interface{})
	require.Len(t, importedIssuersRaw, 1)
	shortCaUuid := importedIssuersRaw[0].(string)

	_, err = client.Logical().Write(mount+"/issuer/"+shortCaUuid, map[string]interface{}{
		"leaf_not_after_behavior": "always_enforce_err",
		"issuer_name":             "short-ca",
	})
	require.NoError(t, err, "failed updating issuer name")

	baseAcmeURL := "/v1/pki/issuer/short-ca/acme/"
	accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// Create new account
	t.Logf("Testing register on %s", baseAcmeURL)
	acct, err := acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Create an order
	t.Logf("Testing Authorize Order on %s", baseAcmeURL)
	identifiers := []string{"*.localdomain"}
	order, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.NoError(t, err, "failed creating order")

	// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow
	//       test.
	markAuthorizationSuccess(t, client, acmeClient, acct, order)

	// Build a proper CSR, with the correct name and signed with a different key works.
	goodCr := &x509.CertificateRequest{DNSNames: []string{identifiers[0]}}
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed generated key for CSR")
	csr, err := x509.CreateCertificateRequest(rand.Reader, goodCr, csrKey)
	require.NoError(t, err, "failed generating csr")

	_, _, err = acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.ErrorContains(t, err, "cannot satisfy request, as TTL would result in notAfter", "failed finalizing order")
}

// TestAcmeTruncatesToIssuerExpiry make sure that if the selected issuer's expiry is shorter than the
// CSR's selected TTL value in ACME and the issuer's leaf_not_after_behavior setting is set to Err,
// we will override the configured behavior and truncate to the issuer's NotAfter
func TestAcmeTruncatesToIssuerExpiry(t *testing.T) {
	t.Parallel()

	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	mount := "pki"
	resp, err := client.Logical().WriteWithContext(context.Background(), mount+"/issuers/generate/intermediate/internal",
		map[string]interface{}{
			"key_name":    "short-key",
			"key_type":    "ec",
			"common_name": "test.com",
		})
	require.NoError(t, err, "failed creating intermediary CSR")
	intermediateCSR := resp.Data["csr"].(string)

	// Sign the intermediate CSR using /pki
	resp, err = client.Logical().Write(mount+"/issuer/root-ca/sign-intermediate", map[string]interface{}{
		"csr":     intermediateCSR,
		"ttl":     "10m",
		"max_ttl": "1h",
	})
	require.NoError(t, err, "failed signing intermediary CSR")
	intermediateCertPEM := resp.Data["certificate"].(string)

	shortCa := parseCert(t, intermediateCertPEM)

	// Configure the intermediate cert as the CA in /pki2
	resp, err = client.Logical().Write(mount+"/issuers/import/cert", map[string]interface{}{
		"pem_bundle": intermediateCertPEM,
	})
	require.NoError(t, err, "failed importing intermediary cert")
	importedIssuersRaw := resp.Data["imported_issuers"].([]interface{})
	require.Len(t, importedIssuersRaw, 1)
	shortCaUuid := importedIssuersRaw[0].(string)

	_, err = client.Logical().Write(mount+"/issuer/"+shortCaUuid, map[string]interface{}{
		"leaf_not_after_behavior": "err",
		"issuer_name":             "short-ca",
	})
	require.NoError(t, err, "failed updating issuer name")

	baseAcmeURL := "/v1/pki/issuer/short-ca/acme/"
	accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// Create new account
	t.Logf("Testing register on %s", baseAcmeURL)
	acct, err := acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Create an order
	t.Logf("Testing Authorize Order on %s", baseAcmeURL)
	identifiers := []string{"*.localdomain"}
	order, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.NoError(t, err, "failed creating order")

	// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow
	//       test.
	markAuthorizationSuccess(t, client, acmeClient, acct, order)

	// Build a proper CSR, with the correct name and signed with a different key works.
	goodCr := &x509.CertificateRequest{DNSNames: []string{identifiers[0]}}
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed generated key for CSR")
	csr, err := x509.CreateCertificateRequest(rand.Reader, goodCr, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.NoError(t, err, "failed finalizing order")
	require.Len(t, certs, 3, "expected full acme chain")

	testAcmeCertSignedByCa(t, client, certs, "short-ca")

	acmeCert, err := x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert")

	require.Equal(t, shortCa.NotAfter, acmeCert.NotAfter, "certificate times aren't the same")
}

// TestAcmeRoleExtKeyUsage verify that ACME by default ignores the role's various ExtKeyUsage flags,
// but if the ACME configuration override of allow_role_ext_key_usage is set that we then honor
// the role's flag.
func TestAcmeRoleExtKeyUsage(t *testing.T) {
	t.Parallel()

	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	roleName := "test-role"

	roleOpt := map[string]interface{}{
		"ttl":                         "365h",
		"max_ttl":                     "720h",
		"key_type":                    "any",
		"allowed_domains":             "localdomain",
		"allow_subdomains":            "true",
		"allow_wildcard_certificates": "true",
		"require_cn":                  "true", /* explicit default */
		"server_flag":                 "true",
		"client_flag":                 "true",
		"code_signing_flag":           "true",
		"email_protection_flag":       "true",
	}

	_, err := client.Logical().Write("pki/roles/"+roleName, roleOpt)

	baseAcmeURL := "/v1/pki/roles/" + roleName + "/acme/"
	accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	require.NoError(t, err, "failed creating role test-role")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// Create new account
	t.Logf("Testing register on %s", baseAcmeURL)
	acct, err := acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Create an order
	t.Logf("Testing Authorize Order on %s", baseAcmeURL)
	identifiers := []string{"*.localdomain"}
	order, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.NoError(t, err, "failed creating order")

	// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow test.
	markAuthorizationSuccess(t, client, acmeClient, acct, order)

	// Build a proper CSR, with the correct name and signed with a different key works.
	goodCr := &x509.CertificateRequest{DNSNames: []string{identifiers[0]}}
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed generated key for CSR")
	csr, err := x509.CreateCertificateRequest(rand.Reader, goodCr, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.NoError(t, err, "order finalization failed")
	require.GreaterOrEqual(t, len(certs), 1, "expected at least one cert in bundle")
	acmeCert, err := x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert")

	require.Equal(t, 1, len(acmeCert.ExtKeyUsage), "mis-match on expected ExtKeyUsages")
	require.ElementsMatch(t, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, acmeCert.ExtKeyUsage,
		"mismatch of ExtKeyUsage flags")

	// Now turn the ACME configuration allow_role_ext_key_usage and retest to make sure we get a certificate
	// with them all
	_, err = client.Logical().WriteWithContext(context.Background(), "pki/config/acme", map[string]interface{}{
		"enabled":                  true,
		"eab_policy":               "not-required",
		"allow_role_ext_key_usage": true,
	})
	require.NoError(t, err, "failed updating ACME configuration")

	t.Logf("Testing Authorize Order on %s", baseAcmeURL)
	order, err = acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.NoError(t, err, "failed creating order")

	// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow test.
	markAuthorizationSuccess(t, client, acmeClient, acct, order)

	certs, _, err = acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.NoError(t, err, "order finalization failed")
	require.GreaterOrEqual(t, len(certs), 1, "expected at least one cert in bundle")
	acmeCert, err = x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert")

	require.Equal(t, 4, len(acmeCert.ExtKeyUsage), "mis-match on expected ExtKeyUsages")
	require.ElementsMatch(t, []x509.ExtKeyUsage{
		x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageCodeSigning, x509.ExtKeyUsageEmailProtection,
	},
		acmeCert.ExtKeyUsage, "mismatch of ExtKeyUsage flags")
}

func TestIssuerRoleDirectoryAssociations(t *testing.T) {
	t.Parallel()

	// This creates two issuers for us (root-ca, int-ca) and two
	// roles (test-role, acme) that we can use with various directory
	// configurations.
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	// Setup DNS for validations.
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dns := dnstest.SetupResolver(t, "dadgarcorp.com")
	defer dns.Cleanup()
	_, err := client.Logical().WriteWithContext(testCtx, "pki/config/acme", map[string]interface{}{
		"dns_resolver": dns.GetLocalAddr(),
	})
	require.NoError(t, err, "failed to specify dns resolver")

	// 1. Use a forbidden role should fail.
	resp, err := client.Logical().WriteWithContext(testCtx, "pki/config/acme", map[string]interface{}{
		"enabled":       true,
		"allowed_roles": []string{"acme"},
	})
	require.NoError(t, err, "failed to write config")
	require.NotNil(t, resp)

	_, err = client.Logical().ReadWithContext(testCtx, "pki/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/default/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role under default issuer")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/int-ca/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role under int-ca issuer")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/root-ca/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role under root-ca issuer")

	_, err = client.Logical().ReadWithContext(testCtx, "pki/roles/acme/acme/directory")
	require.NoError(t, err, "failed to allow usage of acme")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/default/roles/acme/acme/directory")
	require.NoError(t, err, "failed to allow usage of acme under default issuer")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/int-ca/roles/acme/acme/directory")
	require.NoError(t, err, "failed to allow usage of acme under int-ca issuer")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/root-ca/roles/acme/acme/directory")
	require.NoError(t, err, "failed to allow usage of acme under root-ca issuer")

	// 2. Use a forbidden issuer should fail.
	resp, err = client.Logical().WriteWithContext(testCtx, "pki/config/acme", map[string]interface{}{
		"allowed_roles":   []string{"acme"},
		"allowed_issuers": []string{"int-ca"},
	})
	require.NoError(t, err, "failed to write config")
	require.NotNil(t, resp)

	_, err = client.Logical().ReadWithContext(testCtx, "pki/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/default/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role under default issuer")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/int-ca/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role under int-ca issuer")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/root-ca/roles/test-role/acme/directory")
	require.Error(t, err, "failed to forbid usage of test-role under root-ca issuer")

	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/root-ca/roles/acme/acme/directory")
	require.Error(t, err, "failed to forbid usage of acme under root-ca issuer")

	_, err = client.Logical().ReadWithContext(testCtx, "pki/roles/acme/acme/directory")
	require.NoError(t, err, "failed to allow usage of acme")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/default/roles/acme/acme/directory")
	require.NoError(t, err, "failed to allow usage of acme under default issuer")
	_, err = client.Logical().ReadWithContext(testCtx, "pki/issuer/int-ca/roles/acme/acme/directory")
	require.NoError(t, err, "failed to allow usage of acme under int-ca issuer")

	// 3. Setting the default directory to be a sign-verbatim policy and
	// using two different CAs should result in certs signed by each CA.
	resp, err = client.Logical().WriteWithContext(testCtx, "pki/config/acme", map[string]interface{}{
		"allowed_roles":            []string{"*"},
		"allowed_issuers":          []string{"*"},
		"default_directory_policy": "sign-verbatim",
	})
	require.NoError(t, err, "failed to write config")
	require.NotNil(t, resp)

	// default == int-ca
	acmeClientDefault := getAcmeClientForCluster(t, cluster, "/v1/pki/issuer/default/acme/", nil)
	defaultLeafCert := doACMEForDomainWithDNS(t, dns, acmeClientDefault, []string{"default-ca.dadgarcorp.com"})
	requireSignedByAtPath(t, client, defaultLeafCert, "pki/issuer/int-ca")

	acmeClientIntCA := getAcmeClientForCluster(t, cluster, "/v1/pki/issuer/int-ca/acme/", nil)
	intCALeafCert := doACMEForDomainWithDNS(t, dns, acmeClientIntCA, []string{"int-ca.dadgarcorp.com"})
	requireSignedByAtPath(t, client, intCALeafCert, "pki/issuer/int-ca")

	acmeClientRootCA := getAcmeClientForCluster(t, cluster, "/v1/pki/issuer/root-ca/acme/", nil)
	rootCALeafCert := doACMEForDomainWithDNS(t, dns, acmeClientRootCA, []string{"root-ca.dadgarcorp.com"})
	requireSignedByAtPath(t, client, rootCALeafCert, "pki/issuer/root-ca")

	// 4. Using a role-based default directory should allow us to control leaf
	// issuance on the base and issuer-specific directories.
	resp, err = client.Logical().WriteWithContext(testCtx, "pki/config/acme", map[string]interface{}{
		"allowed_roles":            []string{"*"},
		"allowed_issuers":          []string{"*"},
		"default_directory_policy": "role:acme",
	})
	require.NoError(t, err, "failed to write config")
	require.NotNil(t, resp)

	resp, err = client.Logical().JSONMergePatch(testCtx, "pki/roles/acme", map[string]interface{}{
		"ou":             "IT Security",
		"organization":   []string{"Dadgar Corporation, Limited"},
		"allow_any_name": true,
	})
	require.NoError(t, err, "failed to write role differentiator")
	require.NotNil(t, resp)

	for _, issuer := range []string{"", "default", "int-ca", "root-ca"} {
		// Path should override role.
		directory := "/v1/pki/issuer/" + issuer + "/acme/"
		issuerPath := "/pki/issuer/" + issuer
		if issuer == "" {
			directory = "/v1/pki/acme/"
			issuerPath = "/pki/issuer/int-ca"
		} else if issuer == "default" {
			issuerPath = "/pki/issuer/int-ca"
		}

		t.Logf("using directory: %v / issuer: %v", directory, issuerPath)

		acmeClient := getAcmeClientForCluster(t, cluster, directory, nil)
		leafCert := doACMEForDomainWithDNS(t, dns, acmeClient, []string{"role-restricted.dadgarcorp.com"})
		require.Contains(t, leafCert.Subject.Organization, "Dadgar Corporation, Limited", "on directory: %v", directory)
		require.Contains(t, leafCert.Subject.OrganizationalUnit, "IT Security", "on directory: %v", directory)
		requireSignedByAtPath(t, client, leafCert, issuerPath)
	}
}

func TestACMESubjectFieldsAndExtensionsIgnored(t *testing.T) {
	t.Parallel()

	// This creates two issuers for us (root-ca, int-ca) and two
	// roles (test-role, acme) that we can use with various directory
	// configurations.
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	// Setup DNS for validations.
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dns := dnstest.SetupResolver(t, "dadgarcorp.com")
	defer dns.Cleanup()
	_, err := client.Logical().WriteWithContext(testCtx, "pki/config/acme", map[string]interface{}{
		"dns_resolver": dns.GetLocalAddr(),
	})
	require.NoError(t, err, "failed to specify dns resolver")

	// Use the default sign-verbatim policy and ensure OU does not get set.
	directory := "/v1/pki/acme/"
	domains := []string{"no-ou.dadgarcorp.com"}
	acmeClient := getAcmeClientForCluster(t, cluster, directory, nil)
	cr := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: domains[0], OrganizationalUnit: []string{"DadgarCorp IT"}},
		DNSNames: domains,
	}
	cert := doACMEForCSRWithDNS(t, dns, acmeClient, domains, cr)
	t.Logf("Got certificate: %v", cert)
	require.Empty(t, cert.Subject.OrganizationalUnit)

	// Use the default sign-verbatim policy and ensure extension does not get set.
	domains = []string{"no-ext.dadgarcorp.com"}
	extension, err := certutil.CreateDeltaCRLIndicatorExt(12345)
	require.NoError(t, err)
	cr = &x509.CertificateRequest{
		Subject:         pkix.Name{CommonName: domains[0]},
		DNSNames:        domains,
		ExtraExtensions: []pkix.Extension{extension},
	}
	cert = doACMEForCSRWithDNS(t, dns, acmeClient, domains, cr)
	t.Logf("Got certificate: %v", cert)
	for _, ext := range cert.Extensions {
		require.False(t, ext.Id.Equal(certutil.DeltaCRLIndicatorOID))
	}
	require.NotEmpty(t, cert.Extensions)
}

// TestAcmeWithCsrIncludingBasicConstraintExtension verify that we error out for a CSR that is requesting a
// certificate with the IsCA set to true, false is okay, within the basic constraints extension and that no matter what
// the extension is not present on the returned certificate.
func TestAcmeWithCsrIncludingBasicConstraintExtension(t *testing.T) {
	t.Parallel()

	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	baseAcmeURL := "/v1/pki/acme/"
	accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// Create new account
	t.Logf("Testing register on %s", baseAcmeURL)
	acct, err := acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Create an order
	t.Logf("Testing Authorize Order on %s", baseAcmeURL)
	identifiers := []string{"*.localdomain"}
	order, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.NoError(t, err, "failed creating order")

	// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow test.
	markAuthorizationSuccess(t, client, acmeClient, acct, order)

	// Build a CSR with IsCA set to true, making sure we reject it
	extension, err := certutil.CreateBasicConstraintExtension(true, -1)
	require.NoError(t, err, "failed generating basic constraint extension")

	isCATrueCSR := &x509.CertificateRequest{
		DNSNames:        []string{identifiers[0]},
		ExtraExtensions: []pkix.Extension{extension},
	}
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed generated key for CSR")
	csr, err := x509.CreateCertificateRequest(rand.Reader, isCATrueCSR, csrKey)
	require.NoError(t, err, "failed generating csr")

	_, _, err = acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.Error(t, err, "order finalization should have failed with IsCA set to true")

	extension, err = certutil.CreateBasicConstraintExtension(false, -1)
	require.NoError(t, err, "failed generating basic constraint extension")
	isCAFalseCSR := &x509.CertificateRequest{
		DNSNames:   []string{identifiers[0]},
		Extensions: []pkix.Extension{extension},
	}

	csr, err = x509.CreateCertificateRequest(rand.Reader, isCAFalseCSR, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.NoError(t, err, "order finalization should have failed with IsCA set to false")

	require.GreaterOrEqual(t, len(certs), 1, "expected at least one cert in bundle")
	acmeCert, err := x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert")

	// Make sure we don't have any basic constraint extension within the returned cert
	for _, ext := range acmeCert.Extensions {
		if ext.Id.Equal(certutil.ExtensionBasicConstraintsOID) {
			// We shouldn't have this extension in our cert
			t.Fatalf("acme csr contained a basic constraints extension")
		}
	}
}

func markAuthorizationSuccess(t *testing.T, client *api.Client, acmeClient *acme.Client, acct *acme.Account, order *acme.Order) {
	testCtx := context.Background()

	pkiMount := findStorageMountUuid(t, client, "pki")

	// Delete any and all challenge validation entries to stop the engine from overwriting our hack here
	i := 0
	for {
		deleteCvEntries(t, client, pkiMount)

		accountId := acct.URI[strings.LastIndex(acct.URI, "/"):]
		for _, authURI := range order.AuthzURLs {
			authId := authURI[strings.LastIndex(authURI, "/"):]

			// sys/raw does not work with namespaces
			baseClient := client.WithNamespace("")

			values, err := baseClient.Logical().ListWithContext(testCtx, "sys/raw/logical/")
			require.NoError(t, err)
			require.True(t, true, "values: %v", values)

			rawPath := path.Join("sys/raw/logical/", pkiMount, getAuthorizationPath(accountId, authId))
			resp, err := baseClient.Logical().ReadWithContext(testCtx, rawPath)
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
			_, err = baseClient.Logical().WriteWithContext(testCtx, rawPath, map[string]interface{}{
				"value":    base64.StdEncoding.EncodeToString(encodeJSON),
				"encoding": "base64",
			})
			require.NoError(t, err, "failed writing authorization storage")
		}

		// Give some time
		time.Sleep(200 * time.Millisecond)

		// Check to see if we have fixed up the status and no new entries have appeared.
		if !deleteCvEntries(t, client, pkiMount) {
			// No entries found
			// Look to see if we raced against the engine
			orderLookup, err := acmeClient.GetOrder(testCtx, order.URI)
			require.NoError(t, err, "failed loading order status after manually ")

			if orderLookup.Status == string(ACMEOrderReady) {
				// Our order seems to be in the proper status, should be safe-ish to go ahead now
				break
			} else {
				t.Logf("order status was not ready, retrying")
			}
		} else {
			t.Logf("new challenge entries appeared after deletion, retrying")
		}

		if i > 5 {
			t.Fatalf("We are constantly deleting cv entries or order status is not changing, something is wrong")
		}

		i++
	}
}

func deleteCvEntries(t *testing.T, client *api.Client, pkiMount string) bool {
	testCtx := context.Background()

	baseClient := client.WithNamespace("")

	cvPath := path.Join("sys/raw/logical/", pkiMount, acmeValidationPrefix)
	resp, err := baseClient.Logical().ListWithContext(testCtx, cvPath)
	require.NoError(t, err, "failed listing cv path items")

	deletedEntries := false
	if resp != nil {
		cvEntries := resp.Data["keys"].([]interface{})
		for _, cvEntry := range cvEntries {
			cvEntryPath := path.Join(cvPath, cvEntry.(string))
			_, err = baseClient.Logical().DeleteWithContext(testCtx, cvEntryPath)
			require.NoError(t, err, "failed to delete cv entry")
			deletedEntries = true
		}
	}

	return deletedEntries
}

func setupAcmeBackend(t *testing.T) (*vault.TestCluster, *api.Client, string) {
	cluster, client := setupTestPkiCluster(t)

	return setupAcmeBackendOnClusterAtPath(t, cluster, client, "pki")
}

func setupAcmeBackendOnClusterAtPath(t *testing.T, cluster *vault.TestCluster, client *api.Client, mount string) (*vault.TestCluster, *api.Client, string) {
	mount = strings.Trim(mount, "/")

	// Setting templated AIAs should succeed.
	pathConfig := client.Address() + "/v1/" + mount

	namespace := ""
	mountName := mount
	if mount != "pki" {
		if strings.Contains(mount, "/") && constants.IsEnterprise {
			ns_pieces := strings.Split(mount, "/")
			c := len(ns_pieces)
			// mount is c-1
			ns_name := ns_pieces[c-2]
			if len(ns_pieces) > 2 {
				// Parent's namespaces
				parent := strings.Join(ns_pieces[0:c-2], "/")
				_, err := client.WithNamespace(parent).Logical().Write("/sys/namespaces/"+ns_name, nil)
				require.NoError(t, err, "failed to create nested namespaces "+parent+" -> "+ns_name)
			} else {
				_, err := client.Logical().Write("/sys/namespaces/"+ns_name, nil)
				require.NoError(t, err, "failed to create nested namespace "+ns_name)
			}
			namespace = strings.Join(ns_pieces[0:c-1], "/")
			mountName = ns_pieces[c-1]
		}

		err := client.WithNamespace(namespace).Sys().Mount(mountName, &api.MountInput{
			Type: "pki",
			Config: api.MountConfigInput{
				DefaultLeaseTTL: "3000h",
				MaxLeaseTTL:     "600000h",
			},
		})
		require.NoError(t, err, "failed to mount new PKI instance at "+mount)
	}

	err := client.Sys().TuneMountWithContext(ctx, mount, api.MountConfigInput{
		DefaultLeaseTTL: "3000h",
		MaxLeaseTTL:     "600000h",
	})
	require.NoError(t, err, "failed updating mount lease times "+mount)

	_, err = client.Logical().WriteWithContext(context.Background(), mount+"/config/cluster", map[string]interface{}{
		"path":     pathConfig,
		"aia_path": "http://localhost:8200/cdn/" + mount,
	})
	require.NoError(t, err)

	_, err = client.Logical().WriteWithContext(context.Background(), mount+"/config/acme", map[string]interface{}{
		"enabled":    true,
		"eab_policy": "not-required",
	})
	require.NoError(t, err)

	// Allow certain headers to pass through for ACME support
	_, err = client.WithNamespace(namespace).Logical().WriteWithContext(context.Background(), "sys/mounts/"+mountName+"/tune", map[string]interface{}{
		"allowed_response_headers": []string{"Last-Modified", "Replay-Nonce", "Link", "Location"},
		"max_lease_ttl":            "920000h",
	})
	require.NoError(t, err, "failed tuning mount response headers")

	resp, err := client.Logical().WriteWithContext(context.Background(), mount+"/issuers/generate/root/internal",
		map[string]interface{}{
			"issuer_name": "root-ca",
			"key_name":    "root-key",
			"key_type":    "ec",
			"common_name": "Test Root R1 " + mount,
			"ttl":         "7200h",
			"max_ttl":     "920000h",
		})
	require.NoError(t, err, "failed creating root CA")

	resp, err = client.Logical().WriteWithContext(context.Background(), mount+"/issuers/generate/intermediate/internal",
		map[string]interface{}{
			"key_name":    "int-key",
			"key_type":    "ec",
			"common_name": "Test Int X1 " + mount,
		})
	require.NoError(t, err, "failed creating intermediary CSR")
	intermediateCSR := resp.Data["csr"].(string)

	// Sign the intermediate CSR using /pki
	resp, err = client.Logical().Write(mount+"/issuer/root-ca/sign-intermediate", map[string]interface{}{
		"csr":     intermediateCSR,
		"ttl":     "7100h",
		"max_ttl": "910000h",
	})
	require.NoError(t, err, "failed signing intermediary CSR")
	intermediateCertPEM := resp.Data["certificate"].(string)

	// Configure the intermediate cert as the CA in /pki2
	resp, err = client.Logical().Write(mount+"/issuers/import/cert", map[string]interface{}{
		"pem_bundle": intermediateCertPEM,
	})
	require.NoError(t, err, "failed importing intermediary cert")
	importedIssuersRaw := resp.Data["imported_issuers"].([]interface{})
	require.Len(t, importedIssuersRaw, 1)
	intCaUuid := importedIssuersRaw[0].(string)

	_, err = client.Logical().Write(mount+"/issuer/"+intCaUuid, map[string]interface{}{
		"issuer_name": "int-ca",
	})
	require.NoError(t, err, "failed updating issuer name")

	_, err = client.Logical().Write(mount+"/config/issuers", map[string]interface{}{
		"default": "int-ca",
	})
	require.NoError(t, err, "failed updating default issuer")

	_, err = client.Logical().Write(mount+"/roles/test-role", map[string]interface{}{
		"ttl":                         "168h",
		"max_ttl":                     "168h",
		"key_type":                    "any",
		"allowed_domains":             "localdomain",
		"allow_subdomains":            "true",
		"allow_wildcard_certificates": "true",
	})
	require.NoError(t, err, "failed creating role test-role")

	_, err = client.Logical().Write(mount+"/roles/acme", map[string]interface{}{
		"ttl":      "3650h",
		"max_ttl":  "7200h",
		"key_type": "any",
	})
	require.NoError(t, err, "failed creating role acme")

	return cluster, client, pathConfig
}

func testAcmeCertSignedByCa(t *testing.T, client *api.Client, derCerts [][]byte, issuerRef string) *x509.Certificate {
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

	return acmeCert
}

// TestAcmeValidationError make sure that we properly return errors on validation errors.
func TestAcmeValidationError(t *testing.T) {
	t.Parallel()
	cluster, _, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	baseAcmeURL := "/v1/pki/acme/"
	accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// Create new account
	t.Logf("Testing register on %s", baseAcmeURL)
	_, err = acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Create an order
	t.Logf("Testing Authorize Order on %s", baseAcmeURL)
	identifiers := []string{"www.dadgarcorp.com"}
	order, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.NoError(t, err, "failed creating order")

	// Load authorizations
	var authorizations []*acme.Authorization
	for _, authUrl := range order.AuthzURLs {
		auth, err := acmeClient.GetAuthorization(testCtx, authUrl)
		require.NoError(t, err, "failed fetching authorization: %s", authUrl)

		authorizations = append(authorizations, auth)
	}
	require.Len(t, authorizations, 1, "expected a certain number of authorizations")
	require.Len(t, authorizations[0].Challenges, 3, "expected a certain number of challenges associated with authorization")

	acceptedAuth, err := acmeClient.Accept(testCtx, authorizations[0].Challenges[0])
	require.NoError(t, err, "Should have been allowed to accept challenge 1")
	require.Equal(t, string(ACMEChallengeProcessing), acceptedAuth.Status)

	_, err = acmeClient.Accept(testCtx, authorizations[0].Challenges[1])
	require.Error(t, err, "Should have been prevented to accept challenge 2")

	// Make sure our challenge returns errors
	testhelpers.RetryUntil(t, 30*time.Second, func() error {
		challenge, err := acmeClient.GetChallenge(testCtx, authorizations[0].Challenges[0].URI)
		if err != nil {
			return err
		}

		if challenge.Error == nil {
			return fmt.Errorf("no error set in challenge yet")
		}

		acmeError, ok := challenge.Error.(*acme.Error)
		if !ok {
			return fmt.Errorf("unexpected error back: %v", err)
		}

		if acmeError.ProblemType != "urn:ietf:params:acme:error:incorrectResponse" {
			return fmt.Errorf("unexpected ACME error back: %v", acmeError)
		}

		return nil
	})

	// Make sure our challenge,auth and order status change.
	// This takes a little too long to run in CI properly, we need the ability to influence
	// how long the validations take before CI can go wild on this.
	if os.Getenv("CI") == "" {
		testhelpers.RetryUntil(t, 10*time.Minute, func() error {
			challenge, err := acmeClient.GetChallenge(testCtx, authorizations[0].Challenges[0].URI)
			if err != nil {
				return fmt.Errorf("failed to load challenge: %w", err)
			}

			if challenge.Status != string(ACMEChallengeInvalid) {
				return fmt.Errorf("challenge state was not changed to invalid: %v", challenge)
			}

			authz, err := acmeClient.GetAuthorization(testCtx, authorizations[0].URI)
			if err != nil {
				return fmt.Errorf("failed to load authorization: %w", err)
			}

			if authz.Status != string(ACMEAuthorizationInvalid) {
				return fmt.Errorf("authz state was not changed to invalid: %v", authz)
			}

			myOrder, err := acmeClient.GetOrder(testCtx, order.URI)
			if err != nil {
				return fmt.Errorf("failed to load order: %w", err)
			}

			if myOrder.Status != string(ACMEOrderInvalid) {
				return fmt.Errorf("order state was not changed to invalid: %v", order)
			}

			return nil
		})
	}
}

// TestAcmeRevocationAcrossAccounts makes sure that we can revoke certificates using different accounts if
// we have another ACME account or not but access to the certificate key. Also verifies we can't revoke
// certificates across account keys.
func TestAcmeRevocationAcrossAccounts(t *testing.T) {
	t.Parallel()

	cluster, vaultClient, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	baseAcmeURL := "/v1/pki/acme/"
	accountKey1, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient1 := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey1)

	_, leafKey, certs := doACMEWorkflow(t, vaultClient, acmeClient1)
	acmeCert, err := x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert bytes")

	// Make sure our cert is not revoked
	certResp, err := vaultClient.Logical().ReadWithContext(ctx, "pki/cert/"+serialFromCert(acmeCert))
	require.NoError(t, err, "failed to read certificate status")
	require.NotNil(t, certResp, "certificate status response was nil")
	revocationTime := certResp.Data["revocation_time"].(json.Number)
	revocationTimeInt, err := revocationTime.Int64()
	require.NoError(t, err, "failed converting revocation_time value: %v", revocationTime)
	require.Equal(t, revocationTimeInt, int64(0),
		"revocation time was not 0, cert was already revoked: %v", revocationTimeInt)

	// Test that we can't revoke the certificate with another account's key
	accountKey2, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient2 := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey2)
	_, err = acmeClient2.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering second account")

	err = acmeClient2.RevokeCert(ctx, nil, certs[0], acme.CRLReasonUnspecified)
	require.Error(t, err, "should have failed revoking the certificate with a different account")

	// Make sure our cert is not revoked
	certResp, err = vaultClient.Logical().ReadWithContext(ctx, "pki/cert/"+serialFromCert(acmeCert))
	require.NoError(t, err, "failed to read certificate status")
	require.NotNil(t, certResp, "certificate status response was nil")
	revocationTime = certResp.Data["revocation_time"].(json.Number)
	revocationTimeInt, err = revocationTime.Int64()
	require.NoError(t, err, "failed converting revocation_time value: %v", revocationTime)
	require.Equal(t, revocationTimeInt, int64(0),
		"revocation time was not 0, cert was already revoked: %v", revocationTimeInt)

	// But we can revoke if we sign the request with the certificate's key and a different account
	err = acmeClient2.RevokeCert(ctx, leafKey, certs[0], acme.CRLReasonUnspecified)
	require.NoError(t, err, "should have been allowed to revoke certificate with csr key across accounts")

	// Make sure our cert is now revoked
	certResp, err = vaultClient.Logical().ReadWithContext(ctx, "pki/cert/"+serialFromCert(acmeCert))
	require.NoError(t, err, "failed to read certificate status")
	require.NotNil(t, certResp, "certificate status response was nil")
	revocationTime = certResp.Data["revocation_time"].(json.Number)
	revocationTimeInt, err = revocationTime.Int64()
	require.NoError(t, err, "failed converting revocation_time value: %v", revocationTime)
	require.Greater(t, revocationTimeInt, int64(0),
		"revocation time was not greater than 0, cert was not revoked: %v", revocationTimeInt)

	// Make sure we can revoke a certificate without a registered ACME account
	_, leafKey2, certs2 := doACMEWorkflow(t, vaultClient, acmeClient1)

	acmeClient3 := getAcmeClientForCluster(t, cluster, baseAcmeURL, nil)
	err = acmeClient3.RevokeCert(ctx, leafKey2, certs2[0], acme.CRLReasonUnspecified)
	require.NoError(t, err, "should be allowed to revoke a cert with no ACME account but with cert key")

	// Make sure our cert is now revoked
	acmeCert2, err := x509.ParseCertificate(certs2[0])
	require.NoError(t, err, "failed parsing acme cert 2 bytes")

	certResp, err = vaultClient.Logical().ReadWithContext(ctx, "pki/cert/"+serialFromCert(acmeCert2))
	require.NoError(t, err, "failed to read certificate status")
	require.NotNil(t, certResp, "certificate status response was nil")
	revocationTime = certResp.Data["revocation_time"].(json.Number)
	revocationTimeInt, err = revocationTime.Int64()
	require.NoError(t, err, "failed converting revocation_time value: %v", revocationTime)
	require.Greater(t, revocationTimeInt, int64(0),
		"revocation time was not greater than 0, cert was not revoked: %v", revocationTimeInt)
}

// TestAcmeMaxTTL verify that we can update the ACME configuration's max_ttl value and
// get a certificate that has a higher notAfter beyond the 90 day original limit
func TestAcmeMaxTTL(t *testing.T) {
	t.Parallel()
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	numHours := 140 * 24 // The ACME role has a TTL of 152 days
	acmeConfig := map[string]interface{}{
		"enabled":                  true,
		"allowed_issuers":          "*",
		"allowed_roles":            "*",
		"default_directory_policy": "role:acme",
		"dns_resolver":             "",
		"eab_policy_name":          "",
		"max_ttl":                  fmt.Sprintf("%dh", numHours),
	}
	_, err := client.Logical().WriteWithContext(testCtx, "pki/config/acme", acmeConfig)
	require.NoError(t, err, "error configuring acme")

	// First Create Our Client
	accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")
	acmeClient := getAcmeClientForCluster(t, cluster, "/v1/pki/acme/", accountKey)

	discovery, err := acmeClient.Discover(testCtx)
	require.NoError(t, err, "failed acme discovery call")
	t.Logf("%v", discovery)

	acct, err := acmeClient.Register(testCtx, &acme.Account{
		Contact: []string{"mailto:test@example.com"},
	}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")
	require.Equal(t, acme.StatusValid, acct.Status)
	require.Contains(t, acct.Contact, "mailto:test@example.com")
	require.Len(t, acct.Contact, 1)

	authorizations := []acme.AuthzID{
		{"dns", "localhost"},
	}
	// Create an order
	identifiers := make([]string, len(authorizations))
	for index, auth := range authorizations {
		identifiers[index] = auth.Value
	}

	createOrder, err := acmeClient.AuthorizeOrder(testCtx, authorizations)
	require.NoError(t, err, "failed creating order")
	require.Equal(t, acme.StatusPending, createOrder.Status)
	require.Empty(t, createOrder.CertURL)
	require.Equal(t, createOrder.URI+"/finalize", createOrder.FinalizeURL)
	require.Len(t, createOrder.AuthzURLs, len(authorizations), "expected same number of authzurls as identifiers")

	// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow
	//       test.
	markAuthorizationSuccess(t, client, acmeClient, acct, createOrder)

	// Submit the CSR
	requestCSR := x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "localhost"},
	}
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed generated key for CSR")
	csr, err := x509.CreateCertificateRequest(rand.Reader, &requestCSR, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csr, true)
	require.NoError(t, err, "failed finalizing order")

	// Validate we get a signed cert back
	acmeCert := testAcmeCertSignedByCa(t, client, certs, "int-ca")
	duration := time.Duration(numHours) * time.Hour
	maxTTL := time.Now().Add(duration)
	buffer := time.Duration(24) * time.Hour
	dayTruncate := time.Duration(24) * time.Hour

	acmeCertNotAfter := acmeCert.NotAfter.Truncate(dayTruncate)

	// Make sure we are in the ballpark of our max_ttl value.
	require.Greaterf(t, acmeCertNotAfter, maxTTL.Add(-1*buffer), "ACME cert: %v should have been greater than max TTL was %v", acmeCert.NotAfter, maxTTL)
	require.Less(t, acmeCertNotAfter, maxTTL.Add(buffer), "ACME cert: %v should have been less than max TTL was %v", acmeCert.NotAfter, maxTTL)
}

// TestVaultOperatorACMEDisableWorkflow validates that the Vault management API for ACME accounts works as expected.
func TestVaultOperatorACMEDisableWorkflow(t *testing.T) {
	t.Parallel()
	cluster, vaultClient, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Make sure we can call list on an empty ACME mount
	resp, err := vaultClient.Logical().ListWithContext(testCtx, "pki/acme/mgmt/account/keyid")
	require.NoError(t, err, "failed listing acme accounts")
	require.Nil(t, resp, "expected nil, nil response on list of an empty mount")

	// Make sure we get nil, nil response when trying to read a non-existent key (the API returns this for a 404)
	resp, err = vaultClient.Logical().ReadWithContext(testCtx, "pki/acme/mgmt/account/keyid/doesnotexist")
	require.NoError(t, err, "failed reading non-existent ACME key")
	require.Nil(t, resp, "expected nil, nil response on a non-existent ACME key")

	// Make sure we get an error response when trying to write to a non-existent key, unlike read the write call returns an error.
	_, err = vaultClient.Logical().WriteWithContext(testCtx, "pki/acme/mgmt/account/keyid/doesnotexist", map[string]interface{}{"status": "valid"})
	require.ErrorContains(t, err, "did not exist", "failed writing non-existent ACME key")

	baseAcmeURL := "/v1/pki/acme/"
	accountKey1, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey1)
	acct, _, _ := doACMEWorkflow(t, vaultClient, acmeClient)

	// ACME client KID is formatted as https://127.0.0.1:60777/v1/pki/acme/account/45f52b66-a3ed-4080-dec7-cdcff1ef189f
	acmeClientKid := acmeClient.KID

	// Make sure we can call list on an empty ACME mount
	resp, err = vaultClient.Logical().ListWithContext(testCtx, "pki/acme/mgmt/account/keyid")
	require.NoError(t, err, "failed listing acme accounts")
	require.NotNil(t, resp, "expected non-nil response on list on ACME keyid")
	keysFromList := resp.Data["keys"].([]interface{})
	require.Len(t, keysFromList, 1, "expected one key in the list")
	kid := keysFromList[0].(string)
	require.True(t, strings.HasSuffix(string(acmeClientKid), kid), "expected key to match the one the ACME client has")

	// Make sure we can read the key we just listed
	resp, err = vaultClient.Logical().ReadWithContext(testCtx, "pki/acme/mgmt/account/keyid/"+kid)
	require.NoError(t, err, "failed reading ACME with account key")
	require.NotNil(t, resp, "read response was nil on ACME keyid")
	require.Equal(t, "acme/", resp.Data["directory"], "expected directory field in response")
	require.Equal(t, kid, resp.Data["key_id"], "expected key_id field in response")
	require.Equal(t, "valid", resp.Data["status"], "expected status field in response")
	require.NotEmpty(t, resp.Data["created_time"], "expected created_time field in response")
	require.NotEmpty(t, resp.Data["orders"], "expected orders field in response")
	require.Empty(t, resp.Data["revoked_time"], "expected revoked_time field in response")
	require.Empty(t, resp.Data["eab"], "expected eab field in response to be empty")
	orders := resp.Data["orders"].([]interface{})
	require.Len(t, orders, 1, "expected one order in the list")
	order := orders[0].(map[string]interface{})
	require.NotEmpty(t, order["order_id"], "expected order_id field in response")
	require.NotEmpty(t, order["cert_expiry"], "expected cert_expiry field in response")
	require.NotEmpty(t, order["cert_serial_number"], "expected cert_serial_number field in response")

	// Make sure we can update the status of the account to revoked and the revoked_time field is set
	resp, err = vaultClient.Logical().WriteWithContext(testCtx, "pki/acme/mgmt/account/keyid/"+kid, map[string]interface{}{"status": "revoked"})
	require.NoError(t, err, "failed updating writing ACME with account key")
	require.NotNil(t, resp, "expected non-nil response on write to ACME keyid")
	require.Equal(t, "acme/", resp.Data["directory"], "expected directory field in response")
	require.Equal(t, kid, resp.Data["key_id"], "expected key_id field in response")
	require.Equal(t, "revoked", resp.Data["status"], "expected status field in response")
	require.NotEmpty(t, resp.Data["created_time"], "expected created_time field in response")
	require.NotEmpty(t, resp.Data["revoked_time"], "expected revoked_time field in response")
	require.Empty(t, resp.Data["eab"], "expected eab field in response to be empty")
	require.Empty(t, resp.Data["orders"], "write response should not contain orders")

	// Now make sure that we can't use the ACME account anymore
	identifiers := []string{"*.localdomain"}
	_, err = acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.ErrorContains(t, err, "account in status: revoked", "Requesting an order with a revoked account should have failed")

	// Switch the account back to valid and make sure we can use it again
	resp, err = vaultClient.Logical().WriteWithContext(testCtx, "pki/acme/mgmt/account/keyid/"+kid, map[string]interface{}{"status": "valid"})
	require.NoError(t, err, "failed updating writing ACME with account key")
	require.Empty(t, resp.Data["revoked_time"], "revoked_time should have been reset")
	require.Equal(t, "valid", resp.Data["status"], "status should have been reset to valid")

	doACMEOrderWorkflow(t, vaultClient, acmeClient, acct)
}

func doACMEWorkflow(t *testing.T, vaultClient *api.Client, acmeClient *acme.Client) (*acme.Account, *ecdsa.PrivateKey, [][]byte) {
	testCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create new account
	acct, err := acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	if err != nil {
		if strings.Contains(err.Error(), "acme: account already exists") {
			acct, err = acmeClient.GetReg(testCtx, "")
			require.NoError(t, err, "failed looking up account after account exists error?")
		} else {
			require.NoError(t, err, "failed registering account")
		}
	}

	csrKey, certs := doACMEOrderWorkflow(t, vaultClient, acmeClient, acct)
	return acct, csrKey, certs
}

func doACMEOrderWorkflow(t *testing.T, vaultClient *api.Client, acmeClient *acme.Client, acct *acme.Account) (*ecdsa.PrivateKey, [][]byte) {
	testCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an order
	identifiers := []string{"*.localdomain"}
	order, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "dns", Value: identifiers[0]},
	})
	require.NoError(t, err, "failed creating order")

	// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow
	//       test.
	markAuthorizationSuccess(t, vaultClient, acmeClient, acct, order)

	// Build a proper CSR, with the correct name and signed with a different key works.
	goodCr := &x509.CertificateRequest{DNSNames: []string{identifiers[0]}}
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed generated key for CSR")
	csr, err := x509.CreateCertificateRequest(rand.Reader, goodCr, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.NoError(t, err, "failed finalizing order")
	require.Len(t, certs, 3, "expected full acme chain")

	return csrKey, certs
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

func getAcmeClientForCluster(t *testing.T, cluster *vault.TestCluster, baseUrl string, key crypto.Signer) *acme.Client {
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
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl + "/"
	}
	baseAcmeURL := fmt.Sprintf("https://%s/%s", coreAddr.String(), baseUrl)
	return &acme.Client{
		Key:          key,
		HTTPClient:   httpClient,
		DirectoryURL: baseAcmeURL + "directory",
	}
}

func getEABKey(t *testing.T, client *api.Client, baseUrl string) (string, []byte) {
	t.Helper()

	resp, err := client.Logical().WriteWithContext(ctx, path.Join("pki/", baseUrl, "/new-eab"), map[string]interface{}{})
	require.NoError(t, err, "failed getting eab key")
	require.NotNil(t, resp, "eab key returned nil response")
	require.NotEmpty(t, resp.Data["id"], "eab key response missing id field")
	kid := resp.Data["id"].(string)

	require.NotEmpty(t, resp.Data["key"], "eab key response missing private_key field")
	base64Key := resp.Data["key"].(string)
	require.True(t, strings.HasPrefix(base64Key, "vault-eab-0-"), "%s should have had a prefix of vault-eab-0-", base64Key)
	privateKeyBytes, err := base64.RawURLEncoding.DecodeString(base64Key)
	require.NoError(t, err, "failed base 64 decoding eab key response")

	require.Equal(t, "hs", resp.Data["key_type"], "eab key_type field mis-match")
	require.Equal(t, baseUrl+"directory", resp.Data["acme_directory"], "eab acme_directory field mis-match")
	require.NotEmpty(t, resp.Data["created_on"], "empty created_on field")
	_, err = time.Parse(time.RFC3339, resp.Data["created_on"].(string))
	require.NoError(t, err, "failed parsing eab created_on field")

	return kid, privateKeyBytes
}

func TestACMEClientRequestLimits(t *testing.T) {
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	cases := []struct {
		name           string
		authorizations []acme.AuthzID
		requestCSR     x509.CertificateRequest
		valid          bool
	}{
		{
			"validate-only-cn",
			[]acme.AuthzID{
				{"dns", "localhost"},
			},
			x509.CertificateRequest{
				Subject: pkix.Name{CommonName: "localhost"},
			},
			true,
		},
		{
			"validate-only-san",
			[]acme.AuthzID{
				{"dns", "localhost"},
			},
			x509.CertificateRequest{
				DNSNames: []string{"localhost"},
			},
			true,
		},
		{
			"validate-only-ip-address",
			[]acme.AuthzID{
				{"ip", "127.0.0.1"},
			},
			x509.CertificateRequest{
				IPAddresses: []net.IP{{127, 0, 0, 1}},
			},
			true,
		},
	}

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	acmeConfig := map[string]interface{}{
		"enabled":                  true,
		"allowed_issuers":          "*",
		"allowed_roles":            "*",
		"default_directory_policy": "sign-verbatim",
		"dns_resolver":             "",
		"eab_policy_name":          "",
	}
	_, err := client.Logical().WriteWithContext(testCtx, "pki/config/acme", acmeConfig)
	require.NoError(t, err, "error configuring acme")

	for _, tc := range cases {

		// First Create Our Client
		accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
		require.NoError(t, err, "failed creating rsa key")
		acmeClient := getAcmeClientForCluster(t, cluster, "/v1/pki/acme/", accountKey)

		discovery, err := acmeClient.Discover(testCtx)
		require.NoError(t, err, "failed acme discovery call")
		t.Logf("%v", discovery)

		acct, err := acmeClient.Register(testCtx, &acme.Account{
			Contact: []string{"mailto:test@example.com"},
		}, func(tosURL string) bool { return true })
		require.NoError(t, err, "failed registering account")
		require.Equal(t, acme.StatusValid, acct.Status)
		require.Contains(t, acct.Contact, "mailto:test@example.com")
		require.Len(t, acct.Contact, 1)

		// Create an order
		t.Logf("Testing Authorize Order on %s", "pki/acme")
		identifiers := make([]string, len(tc.authorizations))
		for index, auth := range tc.authorizations {
			identifiers[index] = auth.Value
		}

		createOrder, err := acmeClient.AuthorizeOrder(testCtx, tc.authorizations)
		require.NoError(t, err, "failed creating order")
		require.Equal(t, acme.StatusPending, createOrder.Status)
		require.Empty(t, createOrder.CertURL)
		require.Equal(t, createOrder.URI+"/finalize", createOrder.FinalizeURL)
		require.Len(t, createOrder.AuthzURLs, len(tc.authorizations), "expected same number of authzurls as identifiers")

		// HACK: Update authorization/challenge to completed as we can't really do it properly in this workflow
		//       test.
		markAuthorizationSuccess(t, client, acmeClient, acct, createOrder)

		// Submit the CSR
		csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err, "failed generated key for CSR")
		csr, err := x509.CreateCertificateRequest(rand.Reader, &tc.requestCSR, csrKey)
		require.NoError(t, err, "failed generating csr")

		certs, _, err := acmeClient.CreateOrderCert(testCtx, createOrder.FinalizeURL, csr, true)

		if tc.valid {
			require.NoError(t, err, "failed finalizing order")

			// Validate we get a signed cert back
			testAcmeCertSignedByCa(t, client, certs, "int-ca")
		} else {
			require.Error(t, err, "Not a valid CSR, should err")
		}
	}
}
