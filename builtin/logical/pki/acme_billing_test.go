// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"testing"
	"time"

	"golang.org/x/crypto/acme"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki/dnstest"

	"github.com/stretchr/testify/require"
)

// TestACMEBilling is a basic test that will validate client counts created via ACME workflows.
func TestACMEBilling(t *testing.T) {
	t.Parallel()
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	dns := dnstest.SetupResolver(t, "dadgarcorp.com")
	defer dns.Cleanup()

	// Enable additional mounts.
	setupAcmeBackendOnClusterAtPath(t, cluster, client, "pki2")
	setupAcmeBackendOnClusterAtPath(t, cluster, client, "ns1/pki")
	setupAcmeBackendOnClusterAtPath(t, cluster, client, "ns2/pki")

	// Enable custom DNS resolver for testing.
	for _, mount := range []string{"pki", "pki2", "ns1/pki", "ns2/pki"} {
		_, err := client.Logical().Write(mount+"/config/acme", map[string]interface{}{
			"dns_resolver": dns.GetLocalAddr(),
		})
		require.NoError(t, err, "failed to set local dns resolver address for testing on mount: "+mount)
	}

	// Enable client counting.
	_, err := client.Logical().Write("/sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err, "failed to enable client counting")

	// Setup ACME clients. We refresh account keys each time for consistency.
	acmeClientPKI := getAcmeClientForCluster(t, cluster, "/v1/pki/acme/", nil)
	acmeClientPKI2 := getAcmeClientForCluster(t, cluster, "/v1/pki2/acme/", nil)
	acmeClientPKINS1 := getAcmeClientForCluster(t, cluster, "/v1/ns1/pki/acme/", nil)
	acmeClientPKINS2 := getAcmeClientForCluster(t, cluster, "/v1/ns2/pki/acme/", nil)

	// Get our initial count.
	expectedCount := validateClientCount(t, client, -1, "initial fetch")

	// Unique identifier: should increase by one.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKI, []string{"dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount+1, "new certificate")

	// Different identifier; should increase by one.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKI, []string{"example.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount+1, "new certificate")

	// While same identifiers, used together and so thus are unique; increase by one.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKI, []string{"example.dadgarcorp.com", "dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount+1, "new certificate")

	// Same identifiers in different order are not unique; keep the same.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKI, []string{"dadgarcorp.com", "example.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount, "new certificate")

	// Using a different mount shouldn't affect counts.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKI2, []string{"dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount, "different mount; same identifiers")

	// But using a different identifier should.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKI2, []string{"pki2.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount+1, "different mount with different identifiers")

	// A new identifier in a unique namespace will affect results.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKINS1, []string{"unique.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount+1, "unique identifier in a namespace")

	// But in a different namespace with the existing identifier will not.
	doACMEForDomainWithDNS(t, dns, &acmeClientPKINS2, []string{"unique.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount, "existing identifier in a namespace")
	doACMEForDomainWithDNS(t, dns, &acmeClientPKI2, []string{"unique.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, expectedCount, "existing identifier outside of a namespace")
}

func validateClientCount(t *testing.T, client *api.Client, expected int64, message string) int64 {
	resp, err := client.Logical().Read("/sys/internal/counters/activity/monthly")
	require.NoError(t, err, "failed to fetch client count values")
	t.Logf("got client count numbers: %v", resp)

	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Contains(t, resp.Data, "non_entity_clients")

	rawCount := resp.Data["non_entity_clients"].(json.Number)
	count, err := rawCount.Int64()
	require.NoError(t, err, "failed to parse number as int64: "+rawCount.String())

	if expected != -1 {
		require.Equal(t, expected, count, "value of client counts did not match expectations: "+message)
	}

	return count
}

func doACMEForDomainWithDNS(t *testing.T, dns *dnstest.TestServer, acmeClient *acme.Client, domains []string) *x509.Certificate {
	cr := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: domains[0]},
		DNSNames: domains,
	}

	accountKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed to generate account key")
	acmeClient.Key = accountKey

	testCtx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelFunc()

	// Register the client.
	_, err = acmeClient.Register(testCtx, &acme.Account{Contact: []string{"mailto:ipsans@dadgarcorp.com"}}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Create the Order
	var orderIdentifiers []acme.AuthzID
	for _, domain := range domains {
		orderIdentifiers = append(orderIdentifiers, acme.AuthzID{Type: "dns", Value: domain})
	}
	order, err := acmeClient.AuthorizeOrder(testCtx, orderIdentifiers)
	require.NoError(t, err, "failed creating ACME order")

	// Fetch its authorizations.
	var auths []*acme.Authorization
	for _, authUrl := range order.AuthzURLs {
		authorization, err := acmeClient.GetAuthorization(testCtx, authUrl)
		require.NoError(t, err, "failed to lookup authorization at url: %s", authUrl)
		auths = append(auths, authorization)
	}

	// For each dns-01 challenge, place the record in the associated DNS resolver.
	var challengesToAccept []*acme.Challenge
	for _, auth := range auths {
		for _, challenge := range auth.Challenges {
			if challenge.Status != acme.StatusPending {
				t.Logf("ignoring challenge not in status pending: %v", challenge)
				continue
			}

			if challenge.Type == "dns-01" {
				challengeBody, err := acmeClient.DNS01ChallengeRecord(challenge.Token)
				require.NoError(t, err, "failed generating challenge response")

				dns.AddRecord("_acme-challenge."+auth.Identifier.Value, "TXT", challengeBody)
				defer dns.RemoveRecord("_acme-challenge."+auth.Identifier.Value, "TXT", challengeBody)

				require.NoError(t, err, "failed setting DNS record")

				challengesToAccept = append(challengesToAccept, challenge)
			}
		}
	}

	dns.PushConfig()
	require.GreaterOrEqual(t, len(challengesToAccept), 1, "Need at least one challenge, got none")

	// Tell the ACME server, that they can now validate those challenges.
	for _, challenge := range challengesToAccept {
		_, err = acmeClient.Accept(testCtx, challenge)
		require.NoError(t, err, "failed to accept challenge: %v", challenge)
	}

	// Wait for the order/challenges to be validated.
	_, err = acmeClient.WaitOrder(testCtx, order.URI)
	require.NoError(t, err, "failed waiting for order to be ready")

	// Create/sign the CSR and ask ACME server to sign it returning us the final certificate
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	csr, err := x509.CreateCertificateRequest(rand.Reader, cr, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, false)
	require.NoError(t, err, "failed to get a certificate back from ACME")

	acmeCert, err := x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert bytes")

	return acmeCert
}
