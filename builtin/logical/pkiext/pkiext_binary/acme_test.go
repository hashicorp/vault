// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pkiext_binary

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"net"
	"net/http"
	"path"
	"testing"
	"time"

	"golang.org/x/crypto/acme"

	"github.com/hashicorp/vault/builtin/logical/pkiext"
	hDocker "github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/stretchr/testify/require"
)

// Test_ACME will start a Vault cluster using the docker based binary, and execute
// a bunch of sub-tests against that cluster. It is up to each sub-test to run/configure
// a new pki mount within the cluster to not interfere with each other.
func Test_ACME(t *testing.T) {
	cluster := NewVaultPkiClusterWithDNS(t)
	defer cluster.Cleanup()

	tc := map[string]func(t *testing.T, cluster *VaultPkiCluster){
		"certbot":           SubtestACMECertbot,
		"acme ip sans":      SubtestACMEIPAndDNS,
		"acme wildcard":     SubtestACMEWildcardDNS,
		"acme prevents ica": SubtestACMEPreventsICADNS,
	}

	// Wrap the tests within an outer group, so that we run all tests
	// in parallel, but still wait for all tests to finish before completing
	// and running the cleanup of the Vault cluster.
	t.Run("group", func(gt *testing.T) {
		for testName := range tc {
			// Trap the function to be embedded later in the run so it
			// doesn't get clobbered on the next for iteration
			testFunc := tc[testName]

			gt.Run(testName, func(st *testing.T) {
				st.Parallel()
				testFunc(st, cluster)
			})
		}
	})
}

func SubtestACMECertbot(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("pki")
	require.NoError(t, err, "failed setting up acme mount")

	directory := "https://" + pki.GetActiveContainerIP() + ":8200/v1/pki/acme/directory"
	vaultNetwork := pki.GetContainerNetworkName()

	logConsumer, logStdout, logStderr := getDockerLog(t)

	t.Logf("creating on network: %v", vaultNetwork)
	runner, err := hDocker.NewServiceRunner(hDocker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/certbot/certbot",
		ImageTag:      "latest",
		ContainerName: "vault_pki_certbot_test",
		NetworkName:   vaultNetwork,
		Entrypoint:    []string{"sleep", "45"},
		LogConsumer:   logConsumer,
		LogStdout:     logStdout,
		LogStderr:     logStderr,
	})
	require.NoError(t, err, "failed creating service runner")

	ctx := context.Background()
	result, err := runner.Start(ctx, true, false)
	require.NoError(t, err, "could not start container")
	require.NotNil(t, result, "could not start container")

	defer runner.Stop(context.Background(), result.Container.ID)

	networks, err := runner.GetNetworkAndAddresses(result.Container.ID)
	require.NoError(t, err, "could not read container's IP address")
	require.Contains(t, networks, vaultNetwork, "expected to contain vault network")

	ipAddr := networks[vaultNetwork]
	hostname := "certbot-acme-client.dadgarcorp.com"

	err = pki.AddHostname(hostname, ipAddr)
	require.NoError(t, err, "failed to update vault host files")

	certbotCmd := []string{
		"certbot",
		"certonly",
		"--no-eff-email",
		"--email", "certbot.client@dadgarcorp.com",
		"--agree-tos",
		"--no-verify-ssl",
		"--standalone",
		"--non-interactive",
		"--server", directory,
		"-d", hostname,
	}
	logCatCmd := []string{"cat", "/var/log/letsencrypt/letsencrypt.log"}

	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, result.Container.ID, certbotCmd)
	t.Logf("Certbot Issue Command: %v\nstdout: %v\nstderr: %v\n", certbotCmd, string(stdout), string(stderr))
	if err != nil || retcode != 0 {
		logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
		t.Logf("Certbot logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
	}
	require.NoError(t, err, "got error running issue command")
	require.Equal(t, 0, retcode, "expected zero retcode issue command result")

	certbotRevokeCmd := []string{
		"certbot",
		"revoke",
		"--no-eff-email",
		"--email", "certbot.client@dadgarcorp.com",
		"--agree-tos",
		"--no-verify-ssl",
		"--non-interactive",
		"--no-delete-after-revoke",
		"--cert-name", hostname,
	}

	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotRevokeCmd)
	t.Logf("Certbot Revoke Command: %v\nstdout: %v\nstderr: %v\n", certbotRevokeCmd, string(stdout), string(stderr))
	if err != nil || retcode != 0 {
		logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
		t.Logf("Certbot logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
	}
	require.NoError(t, err, "got error running revoke command")
	require.Equal(t, 0, retcode, "expected zero retcode revoke command result")

	// Revoking twice should fail.
	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotRevokeCmd)
	t.Logf("Certbot Double Revoke Command: %v\nstdout: %v\nstderr: %v\n", certbotRevokeCmd, string(stdout), string(stderr))
	if err != nil || retcode == 0 {
		logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
		t.Logf("Certbot logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
	}

	require.NoError(t, err, "got error running double revoke command")
	require.NotEqual(t, 0, retcode, "expected non-zero retcode double revoke command result")
}

func SubtestACMEIPAndDNS(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("pki-ip-dns-sans")
	require.NoError(t, err, "failed setting up acme mount")

	// Since we interact with ACME from outside the container network the ACME
	// configuration needs to be updated to use the host port and not the internal
	// docker ip.
	basePath, err := pki.UpdateClusterConfigLocalAddr()
	require.NoError(t, err, "failed updating cluster config")

	logConsumer, logStdout, logStderr := getDockerLog(t)

	// Setup an nginx container that we can have respond the queries for ips
	runner, err := hDocker.NewServiceRunner(hDocker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/nginx",
		ImageTag:      "latest",
		ContainerName: "vault_pki_ipsans_test",
		NetworkName:   pki.GetContainerNetworkName(),
		LogConsumer:   logConsumer,
		LogStdout:     logStdout,
		LogStderr:     logStderr,
	})
	require.NoError(t, err, "failed creating service runner")

	ctx := context.Background()
	result, err := runner.Start(ctx, true, false)
	require.NoError(t, err, "could not start container")
	require.NotNil(t, result, "could not start container")

	nginxContainerId := result.Container.ID
	defer runner.Stop(context.Background(), nginxContainerId)
	networks, err := runner.GetNetworkAndAddresses(nginxContainerId)

	challengeFolder := "/usr/share/nginx/html/.well-known/acme-challenge/"
	createChallengeFolderCmd := []string{
		"sh", "-c",
		"mkdir -p '" + challengeFolder + "'",
	}
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, nginxContainerId, createChallengeFolderCmd)
	require.NoError(t, err, "failed to create folder in nginx container")
	t.Logf("Update host file command: %v\nstdout: %v\nstderr: %v", createChallengeFolderCmd, string(stdout), string(stderr))
	require.Equal(t, 0, retcode, "expected zero retcode from mkdir in nginx container")

	ipAddr := networks[pki.GetContainerNetworkName()]
	hostname := "go-lang-acme-client.dadgarcorp.com"

	err = pki.AddHostname(hostname, ipAddr)
	require.NoError(t, err, "failed to update vault host files")

	// Perform an ACME lifecycle with an order that contains both an IP and a DNS name identifier
	err = pki.UpdateRole("ip-dns-sans", map[string]interface{}{
		"key_type":                    "any",
		"allowed_domains":             "dadgarcorp.com",
		"allow_subdomains":            "true",
		"allow_wildcard_certificates": "false",
	})
	require.NoError(t, err, "failed creating role ip-dns-sans")

	directoryUrl := basePath + "/roles/ip-dns-sans/acme/directory"
	acmeOrderIdentifiers := []acme.AuthzID{
		{Type: "ip", Value: ipAddr},
		{Type: "dns", Value: hostname},
	}
	cr := &x509.CertificateRequest{
		Subject:     pkix.Name{CommonName: hostname},
		DNSNames:    []string{hostname},
		IPAddresses: []net.IP{net.ParseIP(ipAddr)},
	}

	provisioningFunc := func(acmeClient *acme.Client, auths []*acme.Authorization) []*acme.Challenge {
		// For each http-01 challenge, generate the file to place underneath the nginx challenge folder
		acmeCtx := hDocker.NewBuildContext()
		var challengesToAccept []*acme.Challenge
		for _, auth := range auths {
			for _, challenge := range auth.Challenges {
				if challenge.Status != acme.StatusPending {
					t.Logf("ignoring challenge not in status pending: %v", challenge)
					continue
				}

				if challenge.Type == "http-01" {
					challengeBody, err := acmeClient.HTTP01ChallengeResponse(challenge.Token)
					require.NoError(t, err, "failed generating challenge response")

					challengePath := acmeClient.HTTP01ChallengePath(challenge.Token)
					require.NoError(t, err, "failed generating challenge path")

					challengeFile := path.Base(challengePath)

					acmeCtx[challengeFile] = hDocker.PathContentsFromString(challengeBody)

					challengesToAccept = append(challengesToAccept, challenge)
				}
			}
		}

		require.GreaterOrEqual(t, len(challengesToAccept), 1, "Need at least one challenge, got none")

		// Copy all challenges within the nginx container
		err = runner.CopyTo(nginxContainerId, challengeFolder, acmeCtx)
		require.NoError(t, err, "failed copying challenges to container")

		return challengesToAccept
	}

	acmeCert := doAcmeValidationWithGoLibrary(t, directoryUrl, acmeOrderIdentifiers, cr, provisioningFunc, "")

	require.Len(t, acmeCert.IPAddresses, 1, "expected only a single ip address in cert")
	require.Equal(t, ipAddr, acmeCert.IPAddresses[0].String())
	require.Equal(t, []string{hostname}, acmeCert.DNSNames)
	require.Equal(t, hostname, acmeCert.Subject.CommonName)

	// Perform an ACME lifecycle with an order that contains just an IP identifier
	err = pki.UpdateRole("ip-sans", map[string]interface{}{
		"key_type":            "any",
		"use_csr_common_name": "false",
		"require_cn":          "false",
	})
	require.NoError(t, err, "failed creating role ip-sans")

	directoryUrl = basePath + "/roles/ip-sans/acme/directory"
	acmeOrderIdentifiers = []acme.AuthzID{
		{Type: "ip", Value: ipAddr},
	}
	cr = &x509.CertificateRequest{
		IPAddresses: []net.IP{net.ParseIP(ipAddr)},
	}

	acmeCert = doAcmeValidationWithGoLibrary(t, directoryUrl, acmeOrderIdentifiers, cr, provisioningFunc, "")

	require.Len(t, acmeCert.IPAddresses, 1, "expected only a single ip address in cert")
	require.Equal(t, ipAddr, acmeCert.IPAddresses[0].String())
	require.Empty(t, acmeCert.DNSNames, "acme cert dns name field should have been empty")
	require.Equal(t, "", acmeCert.Subject.CommonName)
}

type acmeGoValidatorProvisionerFunc func(acmeClient *acme.Client, auths []*acme.Authorization) []*acme.Challenge

func doAcmeValidationWithGoLibrary(t *testing.T, directoryUrl string, acmeOrderIdentifiers []acme.AuthzID, cr *x509.CertificateRequest, provisioningFunc acmeGoValidatorProvisionerFunc, expectedFailure string) *x509.Certificate {
	// Since we are contacting Vault through the host ip/port, the certificate will not validate properly
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa account key")

	t.Logf("Using the following url for the ACME directory: %s", directoryUrl)
	acmeClient := &acme.Client{
		Key:          accountKey,
		HTTPClient:   httpClient,
		DirectoryURL: directoryUrl,
	}

	testCtx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelFunc()

	// Create new account
	_, err = acmeClient.Register(testCtx, &acme.Account{Contact: []string{"mailto:ipsans@dadgarcorp.com"}},
		func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Create an ACME order
	order, err := acmeClient.AuthorizeOrder(testCtx, acmeOrderIdentifiers)
	require.NoError(t, err, "failed creating ACME order")

	var auths []*acme.Authorization
	for _, authUrl := range order.AuthzURLs {
		authorization, err := acmeClient.GetAuthorization(testCtx, authUrl)
		require.NoError(t, err, "failed to lookup authorization at url: %s", authUrl)
		auths = append(auths, authorization)
	}

	// Handle the validation using the external validation mechanism.
	challengesToAccept := provisioningFunc(acmeClient, auths)
	require.NotEmpty(t, challengesToAccept, "provisioning function failed to return any challenges to accept")

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

	t.Logf("[TEST-LOG] Created CSR: %v", hex.EncodeToString(csr))

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, false)
	if err != nil {
		if expectedFailure != "" {
			require.Contains(t, err.Error(), expectedFailure, "got a unexpected failure not matching expected value")
			return nil
		}

		require.NoError(t, err, "failed to get a certificate back from ACME")
	} else if expectedFailure != "" {
		t.Fatalf("expected failure containing: %s got none", expectedFailure)
	}

	acmeCert, err := x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert bytes")

	return acmeCert
}

func SubtestACMEWildcardDNS(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("pki-dns-wildcards")
	require.NoError(t, err, "failed setting up acme mount")

	// Since we interact with ACME from outside the container network the ACME
	// configuration needs to be updated to use the host port and not the internal
	// docker ip.
	basePath, err := pki.UpdateClusterConfigLocalAddr()
	require.NoError(t, err, "failed updating cluster config")

	hostname := "go-lang-wildcard-client.dadgarcorp.com"
	wildcard := "*." + hostname

	// Do validation without a role first.
	directoryUrl := basePath + "/acme/directory"
	acmeOrderIdentifiers := []acme.AuthzID{
		{Type: "dns", Value: hostname},
		{Type: "dns", Value: wildcard},
	}
	cr := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: wildcard},
		DNSNames: []string{hostname, wildcard},
	}

	provisioningFunc := func(acmeClient *acme.Client, auths []*acme.Authorization) []*acme.Challenge {
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

					err = pki.AddDNSRecord("_acme-challenge."+auth.Identifier.Value, "TXT", challengeBody)
					require.NoError(t, err, "failed setting DNS record")

					challengesToAccept = append(challengesToAccept, challenge)
				}
			}
		}

		require.GreaterOrEqual(t, len(challengesToAccept), 1, "Need at least one challenge, got none")
		return challengesToAccept
	}

	acmeCert := doAcmeValidationWithGoLibrary(t, directoryUrl, acmeOrderIdentifiers, cr, provisioningFunc, "")
	require.Contains(t, acmeCert.DNSNames, hostname)
	require.Contains(t, acmeCert.DNSNames, wildcard)
	require.Equal(t, wildcard, acmeCert.Subject.CommonName)
	pki.RemoveDNSRecordsForDomain(hostname)

	// Redo validation with a role this time.
	err = pki.UpdateRole("wildcard", map[string]interface{}{
		"key_type":                    "any",
		"allowed_domains":             "go-lang-wildcard-client.dadgarcorp.com",
		"allow_subdomains":            true,
		"allow_bare_domains":          true,
		"allow_wildcard_certificates": true,
	})
	require.NoError(t, err, "failed creating role wildcard")
	directoryUrl = basePath + "/roles/wildcard/acme/directory"

	acmeCert = doAcmeValidationWithGoLibrary(t, directoryUrl, acmeOrderIdentifiers, cr, provisioningFunc, "")
	require.Contains(t, acmeCert.DNSNames, hostname)
	require.Contains(t, acmeCert.DNSNames, wildcard)
	require.Equal(t, wildcard, acmeCert.Subject.CommonName)
	pki.RemoveDNSRecordsForDomain(hostname)
}

func SubtestACMEPreventsICADNS(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("pki-dns-ica")
	require.NoError(t, err, "failed setting up acme mount")

	// Since we interact with ACME from outside the container network the ACME
	// configuration needs to be updated to use the host port and not the internal
	// docker ip.
	basePath, err := pki.UpdateClusterConfigLocalAddr()
	require.NoError(t, err, "failed updating cluster config")

	hostname := "go-lang-intermediate-ca-cert.dadgarcorp.com"

	// Do validation without a role first.
	directoryUrl := basePath + "/acme/directory"
	acmeOrderIdentifiers := []acme.AuthzID{
		{Type: "dns", Value: hostname},
	}
	cr := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: hostname},
		DNSNames: []string{hostname},
		ExtraExtensions: []pkix.Extension{
			// Basic Constraint with IsCA asserted to true.
			{
				Id:       certutil.ExtensionBasicConstraintsOID,
				Critical: true,
				Value:    []byte{0x30, 0x03, 0x01, 0x01, 0xFF},
			},
		},
	}

	provisioningFunc := func(acmeClient *acme.Client, auths []*acme.Authorization) []*acme.Challenge {
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

					err = pki.AddDNSRecord("_acme-challenge."+auth.Identifier.Value, "TXT", challengeBody)
					require.NoError(t, err, "failed setting DNS record")

					challengesToAccept = append(challengesToAccept, challenge)
				}
			}
		}

		require.GreaterOrEqual(t, len(challengesToAccept), 1, "Need at least one challenge, got none")
		return challengesToAccept
	}

	doAcmeValidationWithGoLibrary(t, directoryUrl, acmeOrderIdentifiers, cr, provisioningFunc, "refusing to accept CSR with Basic Constraints extension")
	pki.RemoveDNSRecordsForDomain(hostname)

	// Redo validation with a role this time.
	err = pki.UpdateRole("ica", map[string]interface{}{
		"key_type":                    "any",
		"allowed_domains":             "go-lang-intermediate-ca-cert.dadgarcorp.com",
		"allow_subdomains":            true,
		"allow_bare_domains":          true,
		"allow_wildcard_certificates": true,
	})
	require.NoError(t, err, "failed creating role wildcard")
	directoryUrl = basePath + "/roles/ica/acme/directory"

	doAcmeValidationWithGoLibrary(t, directoryUrl, acmeOrderIdentifiers, cr, provisioningFunc, "refusing to accept CSR with Basic Constraints extension")
	pki.RemoveDNSRecordsForDomain(hostname)
}

func getDockerLog(t *testing.T) (func(s string), *pkiext.LogConsumerWriter, *pkiext.LogConsumerWriter) {
	logConsumer := func(s string) {
		t.Logf(s)
	}

	logStdout := &pkiext.LogConsumerWriter{logConsumer}
	logStderr := &pkiext.LogConsumerWriter{logConsumer}
	return logConsumer, logStdout, logStderr
}
