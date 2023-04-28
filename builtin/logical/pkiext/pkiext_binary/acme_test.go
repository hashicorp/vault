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
	"fmt"
	"net"
	"net/http"
	"path"
	"testing"
	"time"

	"golang.org/x/crypto/acme"

	"github.com/hashicorp/vault/builtin/logical/pkiext"
	hDocker "github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

func Test_ACME(t *testing.T) {
	t.Parallel()

	cluster := NewVaultPkiCluster(t)
	defer cluster.Cleanup()

	t.Run("certbot", func(st *testing.T) { SubtestACMECertbot(st, cluster) })
	t.Run("acme ip sans", func(st *testing.T) { SubTestACMEIPSans(st, cluster) })
}

func SubtestACMECertbot(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("pki")
	require.NoError(t, err, "failed setting up acme mount")

	directory := "https://" + pki.GetActiveContainerIP() + ":8200/v1/pki/acme/directory"
	vaultNetwork := pki.GetContainerNetworkName()

	logConsumer := func(s string) {
		t.Logf(s)
	}

	logStdout := &pkiext.LogConsumerWriter{logConsumer}
	logStderr := &pkiext.LogConsumerWriter{logConsumer}

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

	ctx := context.Background()
	result, err := runner.Start(ctx, true, false)
	require.NoError(t, err, "could not start container")
	require.NotNil(t, result, "could not start container")

	defer runner.Stop(context.Background(), result.Container.ID)

	networks, err := runner.GetNetworkAndAddresses(result.Container.ID)
	require.NoError(t, err, "could not read container's IP address")
	require.Contains(t, networks, vaultNetwork, "expected to contain vault network")

	ipAddr := networks[vaultNetwork]
	hostname := "acme-client.dadgarcorp.com"

	updateHostsCmd := []string{
		"sh", "-c",
		"echo '" + ipAddr + " " + hostname + "' >> /etc/hosts",
	}
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, pki.GetActiveContainerID(), updateHostsCmd)
	require.NoError(t, err, "failed to update vault host file")
	t.Logf("Update host file command: %v\nstdout: %v\nstderr: %v", updateHostsCmd, string(stdout), string(stderr))
	require.Equal(t, 0, retcode, "expected zero retcode from updating vault host file")

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

	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotCmd)
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

func SubTestACMEIPSans(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("pki-ip-sans")
	require.NoError(t, err, "failed setting up acme mount")

	basePath := fmt.Sprintf("https://127.0.0.1:%s/v1/%s", pki.GetActiveContainerExposedPort(), pki.mount)
	err = pki.UpdateClusterConfig(map[string]interface{}{"path": basePath})
	require.NoError(t, err, "failed updating cluster config")

	logConsumer := func(s string) {
		t.Logf(s)
	}

	logStdout := &pkiext.LogConsumerWriter{logConsumer}
	logStderr := &pkiext.LogConsumerWriter{logConsumer}

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

	ctx := context.Background()
	result, err := runner.Start(ctx, true, false)
	require.NoError(t, err, "could not start container")
	require.NotNil(t, result, "could not start container")

	nginxContainerId := result.Container.ID
	defer runner.Stop(context.Background(), nginxContainerId)
	networks, err := runner.GetNetworkAndAddresses(nginxContainerId)

	createChallengeFolderCmd := []string{
		"sh", "-c",
		"mkdir -p '/usr/share/nginx/html/.well-known/acme-challenge/'",
	}
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, nginxContainerId, createChallengeFolderCmd)
	require.NoError(t, err, "failed to create folder in nginx container")
	t.Logf("Update host file command: %v\nstdout: %v\nstderr: %v", createChallengeFolderCmd, string(stdout), string(stderr))
	require.Equal(t, 0, retcode, "expected zero retcode from mkdir in nginx container")

	ipAddr := networks[pki.GetContainerNetworkName()]
	hostname := "go-lang-acme-client.dadgarcorp.com"

	updateHostsCmd := []string{
		"sh", "-c",
		"echo '" + ipAddr + " " + hostname + "' >> /etc/hosts",
	}
	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, pki.GetActiveContainerID(), updateHostsCmd)
	require.NoError(t, err, "failed to update vault host file")
	t.Logf("Update host file command: %v\nstdout: %v\nstderr: %v", updateHostsCmd, string(stdout), string(stderr))
	require.Equal(t, 0, retcode, "expected zero retcode from updating vault host file")

	accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	// Since we are contacting Vault through the host ip/port, the certificate will not validate properly
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	directoryUrl := basePath + "/acme/directory"
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

	order, err := acmeClient.AuthorizeOrder(testCtx, []acme.AuthzID{
		{Type: "ip", Value: ipAddr},
		{Type: "dns", Value: hostname},
	})
	require.NoError(t, err, "failed creating ACME order")
	require.Len(t, order.AuthzURLs, 2)

	var auths []*acme.Authorization
	for _, authUrl := range order.AuthzURLs {
		authorization, err := acmeClient.GetAuthorization(testCtx, authUrl)
		require.NoError(t, err, "failed to lookup authorization at url: %s", authUrl)
		auths = append(auths, authorization)
	}

	acmeCtx := hDocker.NewBuildContext()
	containerPathForChallenges := ""
	var challengesToAccept []*acme.Challenge
	for _, auth := range auths {
		var types []string
		for _, challenge := range auth.Challenges {
			types = append(types, challenge.Type)
		}
		require.Contains(t, types, "http-01")

		if auth.Identifier.Type == "ip" {
			require.Len(t, auth.Challenges, 1, "expected only a single challenge type for ip identifier: %v", auth.Challenges)
		} else {
			require.Len(t, auth.Challenges, 2, "expected multiple challenges for dns name: %v", auth.Challenges)
			require.Contains(t, types, "dns-01")
		}

		for _, challenge := range auth.Challenges {
			if challenge.Status != acme.StatusPending {
				t.Logf("ignoring challenge not in status pending: %v", challenge)
				continue
			}
			if challenge.Type == "http-01" {
				t.Logf("Performing challenge in nginx: %v", challenge)

				challengeBody, err := acmeClient.HTTP01ChallengeResponse(challenge.Token)
				require.NoError(t, err, "failed generating challenge response")

				challengePath := acmeClient.HTTP01ChallengePath(challenge.Token)
				require.NoError(t, err, "failed generating challenge path")

				containerPath := path.Join("/usr/share/nginx/html/", challengePath)
				challengeFile := path.Base(containerPath)
				containerPathForChallenges = path.Dir(containerPath)

				acmeCtx[challengeFile] = hDocker.PathContentsFromString(challengeBody)

				challengesToAccept = append(challengesToAccept, challenge)
			}
		}
	}
	err = runner.CopyTo(nginxContainerId, containerPathForChallenges, acmeCtx)
	require.NoError(t, err, "failed copying challenges to container")

	for _, challenge := range challengesToAccept {
		_, err = acmeClient.Accept(testCtx, challenge)
		require.NoError(t, err, "failed to accept challenge: %v", challenge)
	}

	cr := &x509.CertificateRequest{
		DNSNames:    []string{hostname},
		IPAddresses: []net.IP{net.ParseIP(ipAddr)},
	}
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	csr, err := x509.CreateCertificateRequest(rand.Reader, cr, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, true)
	require.NoError(t, err, "failed to get a certificate back from ACME")

	acmeCert, err := x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert bytes")

	require.Len(t, acmeCert.IPAddresses, 1, "expected only a single ip address in cert")
	require.Equal(t, ipAddr, acmeCert.IPAddresses[0].String())
	require.Equal(t, []string{hostname}, acmeCert.DNSNames)
	require.Equal(t, "", acmeCert.Subject.CommonName)
}
