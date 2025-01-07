// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki/dnstest"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/acme"
)

// TestACMEBilling is a basic test that will validate client counts created via ACME workflows.
func TestACMEBilling(t *testing.T) {
	t.Parallel()
	timeutil.SkipAtEndOfMonth(t)

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
	expectedCount := validateClientCount(t, client, "", -1, "initial fetch")

	// Unique identifier: should increase by one.
	doACMEForDomainWithDNS(t, dns, acmeClientPKI, []string{"dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "pki", expectedCount+1, "new certificate")

	// Different identifier; should increase by one.
	doACMEForDomainWithDNS(t, dns, acmeClientPKI, []string{"example.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "pki", expectedCount+1, "new certificate")

	// While same identifiers, used together and so thus are unique; increase by one.
	doACMEForDomainWithDNS(t, dns, acmeClientPKI, []string{"example.dadgarcorp.com", "dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "pki", expectedCount+1, "new certificate")

	// Same identifiers in different order are not unique; keep the same.
	doACMEForDomainWithDNS(t, dns, acmeClientPKI, []string{"dadgarcorp.com", "example.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "pki", expectedCount, "different order; same identifiers")

	// Using a different mount shouldn't affect counts.
	doACMEForDomainWithDNS(t, dns, acmeClientPKI2, []string{"dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "", expectedCount, "different mount; same identifiers")

	// But using a different identifier should.
	doACMEForDomainWithDNS(t, dns, acmeClientPKI2, []string{"pki2.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "pki2", expectedCount+1, "different mount with different identifiers")

	// A new identifier in a unique namespace will affect results.
	doACMEForDomainWithDNS(t, dns, acmeClientPKINS1, []string{"unique.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "ns1/pki", expectedCount+1, "unique identifier in a namespace")

	// But in a different namespace with the existing identifier will not.
	doACMEForDomainWithDNS(t, dns, acmeClientPKINS2, []string{"unique.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "", expectedCount, "existing identifier in a namespace")
	doACMEForDomainWithDNS(t, dns, acmeClientPKI2, []string{"unique.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "", expectedCount, "existing identifier outside of a namespace")

	// Creating a unique identifier in a namespace with a mount with the
	// same name as another namespace should increase counts as well.
	doACMEForDomainWithDNS(t, dns, acmeClientPKINS2, []string{"very-unique.dadgarcorp.com"})
	expectedCount = validateClientCount(t, client, "ns2/pki", expectedCount+1, "unique identifier in a different namespace")

	// Check the current fragment
	fragment := cluster.Cores[0].Core.ResetActivityLog()[0]
	if fragment == nil {
		t.Fatal("no fragment created")
	}
	validateAcmeClientTypes(t, fragment, expectedCount)
}

func validateAcmeClientTypes(t *testing.T, fragment *activity.LogFragment, expectedCount int64) {
	t.Helper()
	if int64(len(fragment.Clients)) != expectedCount {
		t.Fatalf("bad number of entities, expected %v: got %v, entities are: %v", expectedCount, len(fragment.Clients), fragment.Clients)
	}

	for _, ac := range fragment.Clients {
		if ac.ClientType != vault.ACMEActivityType {
			t.Fatalf("Couldn't find expected '%v' client_type in %v", vault.ACMEActivityType, fragment.Clients)
		}
	}
}

func validateClientCount(t *testing.T, client *api.Client, mount string, expected int64, message string) int64 {
	resp, err := client.Logical().Read("/sys/internal/counters/activity/monthly")
	require.NoError(t, err, "failed to fetch client count values")
	t.Logf("got client count numbers: %v", resp)

	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Contains(t, resp.Data, "acme_clients")
	require.Contains(t, resp.Data, "months")

	rawCount := resp.Data["acme_clients"].(json.Number)
	count, err := rawCount.Int64()
	require.NoError(t, err, "failed to parse number as int64: "+rawCount.String())

	if expected != -1 {
		require.Equal(t, expected, count, "value of client counts did not match expectations: "+message)
	}

	if mount == "" {
		return count
	}

	months := resp.Data["months"].([]interface{})
	if len(months) > 1 {
		t.Fatalf("running across a month boundary despite using SkipAtEndOfMonth(...); rerun test from start fully in the next month instead")
	}

	require.Equal(t, 1, len(months), "expected only a single month when running this test")

	monthlyInfo := months[0].(map[string]interface{})

	// Validate this month's aggregate counts match the overall value.
	require.Contains(t, monthlyInfo, "counts", "expected monthly info to contain a count key")
	monthlyCounts := monthlyInfo["counts"].(map[string]interface{})
	require.Contains(t, monthlyCounts, "acme_clients", "expected month[0].counts to contain a non_entity_clients key")
	monthlyCountNonEntityRaw := monthlyCounts["acme_clients"].(json.Number)
	monthlyCountNonEntity, err := monthlyCountNonEntityRaw.Int64()
	require.NoError(t, err, "failed to parse number as int64: "+monthlyCountNonEntityRaw.String())
	require.Equal(t, count, monthlyCountNonEntity, "expected equal values for non entity client counts")

	// Validate this mount's namespace is included in the namespaces list,
	// if this is enterprise. Otherwise, if its OSS or we don't have a
	// namespace, we default to the value root.
	mountNamespace := ""
	mountPath := mount + "/"
	if constants.IsEnterprise && strings.Contains(mount, "/") {
		pieces := strings.Split(mount, "/")
		require.Equal(t, 2, len(pieces), "we do not support nested namespaces in this test")
		mountNamespace = pieces[0] + "/"
		mountPath = pieces[1] + "/"
	}

	require.Contains(t, monthlyInfo, "namespaces", "expected monthly info to contain a namespaces key")
	monthlyNamespaces := monthlyInfo["namespaces"].([]interface{})
	foundNamespace := false
	for index, namespaceRaw := range monthlyNamespaces {
		namespace := namespaceRaw.(map[string]interface{})
		require.Contains(t, namespace, "namespace_path", "expected monthly.namespaces[%v] to contain a namespace_path key", index)
		namespacePath := namespace["namespace_path"].(string)

		if namespacePath != mountNamespace {
			t.Logf("skipping non-matching namespace %v: %v != %v / %v", index, namespacePath, mountNamespace, namespace)
			continue
		}

		foundNamespace = true

		// This namespace must have a non-empty aggregate non-entity count.
		require.Contains(t, namespace, "counts", "expected monthly.namespaces[%v] to contain a counts key", index)
		namespaceCounts := namespace["counts"].(map[string]interface{})
		require.Contains(t, namespaceCounts, "acme_clients", "expected namespace counts to contain a non_entity_clients key")
		namespaceCountNonEntityRaw := namespaceCounts["acme_clients"].(json.Number)
		namespaceCountNonEntity, err := namespaceCountNonEntityRaw.Int64()
		require.NoError(t, err, "failed to parse number as int64: "+namespaceCountNonEntityRaw.String())
		require.Greater(t, namespaceCountNonEntity, int64(0), "expected at least one non-entity client count value in the namespace")

		require.Contains(t, namespace, "mounts", "expected monthly.namespaces[%v] to contain a mounts key", index)
		namespaceMounts := namespace["mounts"].([]interface{})
		foundMount := false
		for mountIndex, mountRaw := range namespaceMounts {
			mountInfo := mountRaw.(map[string]interface{})
			require.Contains(t, mountInfo, "mount_path", "expected monthly.namespaces[%v].mounts[%v] to contain a mount_path key", index, mountIndex)
			mountInfoPath := mountInfo["mount_path"].(string)
			if mountPath != mountInfoPath {
				t.Logf("skipping non-matching mount path %v in namespace %v: %v != %v / %v of %v", mountIndex, index, mountPath, mountInfoPath, mountInfo, namespace)
				continue
			}

			foundMount = true

			// This mount must also have a non-empty non-entity client count.
			require.Contains(t, mountInfo, "counts", "expected monthly.namespaces[%v].mounts[%v] to contain a counts key", index, mountIndex)
			mountCounts := mountInfo["counts"].(map[string]interface{})
			require.Contains(t, mountCounts, "acme_clients", "expected mount counts to contain a non_entity_clients key")
			mountCountNonEntityRaw := mountCounts["acme_clients"].(json.Number)
			mountCountNonEntity, err := mountCountNonEntityRaw.Int64()
			require.NoError(t, err, "failed to parse number as int64: "+mountCountNonEntityRaw.String())
			require.Greater(t, mountCountNonEntity, int64(0), "expected at least one non-entity client count value in the mount")
		}

		require.True(t, foundMount, "expected to find the mount "+mountPath+" in the list of mounts for namespace, but did not")
	}

	require.True(t, foundNamespace, "expected to find the namespace "+mountNamespace+" in the list of namespaces, but did not")

	return count
}

func doACMEForDomainWithDNS(t *testing.T, dns *dnstest.TestServer, acmeClient *acme.Client, domains []string) *x509.Certificate {
	cr := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: domains[0]},
		DNSNames: domains,
	}

	return doACMEForCSRWithDNS(t, dns, acmeClient, domains, cr)
}

func doACMEForCSRWithDNS(t *testing.T, dns *dnstest.TestServer, acmeClient *acme.Client, domains []string, cr *x509.CertificateRequest) *x509.Certificate {
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
