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
	"os"
	"path"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"golang.org/x/crypto/acme"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pkiext"
	hDocker "github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	tcDocker "github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	"github.com/stretchr/testify/require"
)

func CheckCertBot(t *testing.T, vaultNetwork string, vaultNodeID string, directory string) {
	logConsumer := func(s string) {
		t.Logf(s)
	}

	logStdout := &pkiext.LogConsumerWriter{Consumer: logConsumer}
	logStderr := &pkiext.LogConsumerWriter{Consumer: logConsumer}

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
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, vaultNodeID, updateHostsCmd)
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

func RunACMERootTest(t *testing.T, caKeyType string, caKeyBits int, caUsePSS bool, roleKeyType string, roleKeyBits int, roleUsePSS bool) {
	cluster, vaultNetwork, vaultAddr, vaultNodeID, _ := setupVaultDocker(t)
	defer cluster.Cleanup()

	setupAcme(t, cluster, vaultAddr, "8200", caKeyType, caKeyBits, caUsePSS, roleKeyType, roleKeyBits, roleUsePSS)

	directory := "https://" + vaultAddr + ":8200/v1/pki/acme/directory"
	CheckCertBot(t, vaultNetwork, vaultNodeID, directory)
}

func setupAcme(t *testing.T, cluster *tcDocker.DockerCluster, vaultAddr string, vaultPort string, caKeyType string,
	caKeyBits int, caUsePSS bool, roleKeyType string, roleKeyBits int, roleUsePSS bool,
) {
	testSuffix := fmt.Sprintf(" - %v %v %v - %v %v %v", caKeyType, caKeyType, caUsePSS, roleKeyType, roleKeyBits, roleUsePSS)

	client := cluster.Nodes()[0].APIClient()
	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
			AllowedResponseHeaders: []string{
				"Last-Modified", "Replay-Nonce",
				"Link", "Location",
			},
		},
	})
	require.NoError(t, err, "failed mounting pki endpoint")

	// Set URLs pointing to the issuer.
	_, err = client.Logical().Write("pki/config/cluster", map[string]interface{}{
		"path":     "https://" + vaultAddr + ":" + vaultPort + "/v1/pki",
		"aia_path": "http://" + vaultAddr + ":" + vaultPort + "/v1/pki",
	})
	require.NoError(t, err)

	// Setup root+intermediate CA hierarchy within this mount.
	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name":  "Root X1" + testSuffix,
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"key_bits":     caKeyBits,
		"use_pss":      caUsePSS,
		"issuer_name":  "root",
	})
	require.NoError(t, err, "failed to create root cert")
	require.NotNil(t, resp, "failed to create root cert")
	require.NotEmpty(t, resp.Data, "failed to create root cert")
	// rootCert := resp.Data["certificate"].(string)
	resp, err = client.Logical().Write("pki/intermediate/generate/internal", map[string]interface{}{
		"common_name":  "Intermediate I1" + testSuffix,
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"key_bits":     caKeyBits,
		"use_pss":      caUsePSS,
	})
	require.NoError(t, err, "failed to create int csr")
	require.NotNil(t, resp, "failed to create int csr")
	require.NotEmpty(t, resp.Data, "failed to create int csr")
	resp, err = client.Logical().Write("pki/issuer/default/sign-intermediate", map[string]interface{}{
		"common_name":  "Intermediate I1",
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"csr":          resp.Data["csr"],
	})
	require.NoError(t, err, "failed to create sign int")
	require.NotNil(t, resp, "failed to create sign int")
	require.NotEmpty(t, resp.Data, "failed to create sign int")
	intCert := resp.Data["certificate"].(string)
	resp, err = client.Logical().Write("pki/issuers/import/bundle", map[string]interface{}{
		"pem_bundle": intCert,
	})
	require.NoError(t, err, "failed to create import int")
	require.NotNil(t, resp, "failed to create import int")
	require.NotEmpty(t, resp.Data, "failed to create import int")
	_, err = client.Logical().Write("pki/config/issuers", map[string]interface{}{
		"default": resp.Data["imported_issuers"].([]interface{})[0],
	})
	require.NoError(t, err, "failed to set intermediate as default")
	resp, err = client.Logical().JSONMergePatch(context.Background(), "pki/issuer/default", map[string]interface{}{
		"leaf_not_after_behavior": "truncate",
	})
	require.NoError(t, err, "failed to update intermediate ttl behavior")
	t.Logf("got response from updating int issuer: %v", resp)
	_, err = client.Logical().JSONMergePatch(context.Background(), "pki/issuer/root", map[string]interface{}{
		"leaf_not_after_behavior": "truncate",
	})
	require.NoError(t, err, "failed to update root ttl behavior")
	t.Logf("got response from updating root issuer: %v", resp)
}

func setupVaultDocker(t *testing.T) (*tcDocker.DockerCluster, string, string, string, nat.PortBinding) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}
	opts := &tcDocker.DockerClusterOptions{
		ImageRepo: "docker.mirror.hashicorp.services/hashicorp/vault",
		// We're replacing the binary anyway, so we're not too particular about
		// the docker image version tag.
		ImageTag:    "latest",
		VaultBinary: binary,
		ClusterOptions: testcluster.ClusterOptions{
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
			},
			NumCores: 1,
		},
	}

	cluster := tcDocker.NewTestDockerCluster(t, opts)

	var vaultNetwork string
	var vaultAddr string
	var vaultNodeID string
	var vaultBindedPort nat.PortBinding
	for index, rawNode := range cluster.Nodes() {
		node, ok := rawNode.(*tcDocker.DockerClusterNode)
		require.True(t, ok, "failed to cast NewTestDockerCluster's Node to DockerClusterNode")
		t.Logf("[%d] Cluster Node %v - %v / %v", index, node.Name(), node.ContainerNetworkName, node.ContainerIPAddress)
		if index == 0 {
			vaultNodeID = node.Container.ID
			vaultNetwork = node.ContainerNetworkName
			vaultAddr = node.ContainerIPAddress
			vaultBindedPort = node.Container.NetworkSettings.Ports["8200/tcp"][0]
		}
	}
	return cluster, vaultNetwork, vaultAddr, vaultNodeID, vaultBindedPort
}

func Test_ACMERSAPure(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "rsa", 2048, false, "rsa", 2048, false)
}

func Test_ACMERSAPurePSS(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "rsa", 2048, false, "rsa", 2048, true)
}

func Test_ACMERSAPSSPure(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "rsa", 2048, true, "rsa", 2048, false)
}

func Test_ACMERSAPSSPurePSS(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "rsa", 2048, true, "rsa", 2048, true)
}

func Test_ACMEECDSA256Pure(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "ec", 256, false, "ec", 256, false)
}

func Test_ACMEECDSAHybrid(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "ec", 256, false, "rsa", 2048, false)
}

func Test_ACMEECDSAHybridPSS(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "ec", 256, false, "rsa", 2048, true)
}

func Test_ACMERSAHybrid(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "rsa", 2048, false, "ec", 256, false)
}

func Test_ACMERSAPSSHybrid(t *testing.T) {
	t.Parallel()

	RunACMERootTest(t, "rsa", 2048, true, "ec", 256, false)
}

// Test_ACMEIPSans verify that we can perform ACME validations on IP and DNS identifiers using the Golang ACME library.
func Test_ACMEIPSans(t *testing.T) {
	cluster, vaultNetwork, _, vaultNodeID, vaultBindedPort := setupVaultDocker(t)
	defer cluster.Cleanup()

	setupAcme(t, cluster, "127.0.0.1", vaultBindedPort.HostPort, "rsa", 2048, true, "ec", 256, false)

	logConsumer := func(s string) {
		t.Logf(s)
	}

	logStdout := &pkiext.LogConsumerWriter{Consumer: logConsumer}
	logStderr := &pkiext.LogConsumerWriter{Consumer: logConsumer}

	// Setup an nginx container that we can have respond the queries for ips
	runner, err := hDocker.NewServiceRunner(hDocker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/nginx",
		ImageTag:      "latest",
		ContainerName: "vault_pki_ipsans_test",
		NetworkName:   vaultNetwork,
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
	require.NoError(t, err, "could not read container's IP address")
	require.Contains(t, networks, vaultNetwork, "expected to contain vault network")

	createChallengeFolderCmd := []string{
		"sh", "-c",
		"mkdir -p '/usr/share/nginx/html/.well-known/acme-challenge/'",
	}
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, nginxContainerId, createChallengeFolderCmd)
	require.NoError(t, err, "failed to create folder in nginx container")
	t.Logf("Update host file command: %v\nstdout: %v\nstderr: %v", createChallengeFolderCmd, string(stdout), string(stderr))
	require.Equal(t, 0, retcode, "expected zero retcode from mkdir in nginx container")

	ipAddr := networks[vaultNetwork]
	hostname := "go-lang-acme-client.dadgarcorp.com"

	updateHostsCmd := []string{
		"sh", "-c",
		"echo '" + ipAddr + " " + hostname + "' >> /etc/hosts",
	}
	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, vaultNodeID, updateHostsCmd)
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

	directoryUrl := fmt.Sprintf("https://%s:%s/%s", vaultBindedPort.HostIP, vaultBindedPort.HostPort, "v1/pki/acme/directory")
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
