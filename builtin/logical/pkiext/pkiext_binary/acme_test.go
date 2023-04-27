// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pkiext_binary

import (
	"context"
	"fmt"
	"os"
	"testing"

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

	runner.Stop(ctx, result.Container.ID)
}

func RunACMERootTest(t *testing.T, caKeyType string, caKeyBits int, caUsePSS bool, roleKeyType string, roleKeyBits int, roleUsePSS bool) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}

	testSuffix := fmt.Sprintf(" - %v %v %v - %v %v %v", caKeyType, caKeyType, caUsePSS, roleKeyType, roleKeyBits, roleUsePSS)

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
		},
	}

	cluster := tcDocker.NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()

	var vaultNetwork string
	var vaultAddr string
	var vaultNodeID string
	for index, rawNode := range cluster.Nodes() {
		node, ok := rawNode.(*tcDocker.DockerClusterNode)
		require.True(t, ok, "failed to cast NewTestDockerCluster's Node to DockerClusterNode")
		t.Logf("[%d] Cluster Node %v - %v / %v", index, node.Name(), node.ContainerNetworkName, node.ContainerIPAddress)
		if index == 0 {
			vaultNodeID = node.Container.ID
			vaultNetwork = node.ContainerNetworkName
			vaultAddr = node.ContainerIPAddress
		}
	}

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
		"path":     "https://" + vaultAddr + ":8200/v1/pki",
		"aia_path": "http://" + vaultAddr + ":8200/v1/pki",
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

	directory := "https://" + vaultAddr + ":8200/v1/pki/acme/directory"
	CheckCertBot(t, vaultNetwork, vaultNodeID, directory)
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
