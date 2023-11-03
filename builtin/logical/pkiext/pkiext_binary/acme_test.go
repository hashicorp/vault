// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"path"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/acme"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/logical/pkiext"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	hDocker "github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/caddy_http.json
var caddyConfigTemplateHTTP string

//go:embed testdata/caddy_http_eab.json
var caddyConfigTemplateHTTPEAB string

//go:embed testdata/caddy_tls_alpn.json
var caddyConfigTemplateTLSALPN string

// Test_ACME will start a Vault cluster using the docker based binary, and execute
// a bunch of sub-tests against that cluster. It is up to each sub-test to run/configure
// a new pki mount within the cluster to not interfere with each other.
func Test_ACME(t *testing.T) {
	cluster := NewVaultPkiClusterWithDNS(t)
	defer cluster.Cleanup()

	tc := map[string]func(t *testing.T, cluster *VaultPkiCluster){
		"caddy http":        SubtestACMECaddy(caddyConfigTemplateHTTP, false),
		"caddy http eab":    SubtestACMECaddy(caddyConfigTemplateHTTPEAB, true),
		"caddy tls-alpn":    SubtestACMECaddy(caddyConfigTemplateTLSALPN, false),
		"certbot":           SubtestACMECertbot,
		"certbot eab":       SubtestACMECertbotEab,
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

	// Do not run these tests in parallel.
	t.Run("step down", func(gt *testing.T) { SubtestACMEStepDownNode(gt, cluster) })
}

// caddyConfig contains information used to render a Caddy configuration file from a template.
type caddyConfig struct {
	Hostname  string
	Directory string
	CACert    string
	EABID     string
	EABKey    string
}

// SubtestACMECaddy returns an ACME test for Caddy using the provided template.
func SubtestACMECaddy(configTemplate string, enableEAB bool) func(*testing.T, *VaultPkiCluster) {
	return func(t *testing.T, cluster *VaultPkiCluster) {
		ctx := context.Background()

		// Roll a random run ID for mount and hostname uniqueness.
		runID, err := uuid.GenerateUUID()
		require.NoError(t, err, "failed to generate a unique ID for test run")
		runID = strings.Split(runID, "-")[0]

		// Create the PKI mount with ACME enabled
		pki, err := cluster.CreateAcmeMount(runID)
		require.NoError(t, err, "failed to set up ACME mount")

		// Conditionally enable EAB and retrieve the key.
		var eabID, eabKey string
		if enableEAB {
			err = pki.UpdateAcmeConfig(true, map[string]interface{}{
				"eab_policy": "new-account-required",
			})
			require.NoError(t, err, "failed to configure EAB policy in PKI mount")

			eabID, eabKey, err = pki.GetEabKey("acme/")
			require.NoError(t, err, "failed to retrieve EAB key from PKI mount")
		}

		directory := fmt.Sprintf("https://%s:8200/v1/%s/acme/directory", pki.GetActiveContainerIP(), runID)
		vaultNetwork := pki.GetContainerNetworkName()
		t.Logf("dir: %s", directory)

		logConsumer, logStdout, logStderr := getDockerLog(t)

		sleepTimer := "45"

		// Kick off Caddy container.
		t.Logf("creating on network: %v", vaultNetwork)
		caddyRunner, err := hDocker.NewServiceRunner(hDocker.RunOptions{
			ImageRepo:     "docker.mirror.hashicorp.services/library/caddy",
			ImageTag:      "2.6.4",
			ContainerName: fmt.Sprintf("caddy_test_%s", runID),
			NetworkName:   vaultNetwork,
			Ports:         []string{"80/tcp", "443/tcp", "443/udp"},
			Entrypoint:    []string{"sleep", sleepTimer},
			LogConsumer:   logConsumer,
			LogStdout:     logStdout,
			LogStderr:     logStderr,
		})
		require.NoError(t, err, "failed creating caddy service runner")

		caddyResult, err := caddyRunner.Start(ctx, true, false)
		require.NoError(t, err, "could not start Caddy container")
		require.NotNil(t, caddyResult, "could not start Caddy container")

		defer caddyRunner.Stop(ctx, caddyResult.Container.ID)

		networks, err := caddyRunner.GetNetworkAndAddresses(caddyResult.Container.ID)
		require.NoError(t, err, "could not read caddy container's IP address")
		require.Contains(t, networks, vaultNetwork, "expected to contain vault network")

		ipAddr := networks[vaultNetwork]
		hostname := fmt.Sprintf("%s.dadgarcorp.com", runID)

		err = pki.AddHostname(hostname, ipAddr)
		require.NoError(t, err, "failed to update vault host files")

		// Render the Caddy configuration from the specified template.
		tmpl, err := template.New("config").Parse(configTemplate)
		require.NoError(t, err, "failed to parse Caddy config template")
		var b strings.Builder
		err = tmpl.Execute(
			&b,
			caddyConfig{
				Hostname:  hostname,
				Directory: directory,
				CACert:    "/tmp/vault_ca_cert.crt",
				EABID:     eabID,
				EABKey:    eabKey,
			},
		)
		require.NoError(t, err, "failed to render Caddy config template")

		// Push the Caddy config and the cluster listener's CA certificate over to the docker container.
		cpCtx := hDocker.NewBuildContext()
		cpCtx["caddy_config.json"] = hDocker.PathContentsFromString(b.String())
		cpCtx["vault_ca_cert.crt"] = hDocker.PathContentsFromString(string(cluster.GetListenerCACertPEM()))
		err = caddyRunner.CopyTo(caddyResult.Container.ID, "/tmp/", cpCtx)
		require.NoError(t, err, "failed to copy Caddy config and Vault listener CA certificate to container")

		// Start the Caddy server.
		caddyCmd := []string{
			"caddy",
			"start",
			"--config", "/tmp/caddy_config.json",
		}
		stdout, stderr, retcode, err := caddyRunner.RunCmdWithOutput(ctx, caddyResult.Container.ID, caddyCmd)
		t.Logf("Caddy Start Command: %v\nstdout: %v\nstderr: %v\n", caddyCmd, string(stdout), string(stderr))
		require.NoError(t, err, "got error running Caddy start command")
		require.Equal(t, 0, retcode, "expected zero retcode Caddy start command result")

		// Start a cURL container.
		curlRunner, err := hDocker.NewServiceRunner(hDocker.RunOptions{
			ImageRepo:     "docker.mirror.hashicorp.services/curlimages/curl",
			ImageTag:      "8.4.0",
			ContainerName: fmt.Sprintf("curl_test_%s", runID),
			NetworkName:   vaultNetwork,
			Entrypoint:    []string{"sleep", sleepTimer},
			LogConsumer:   logConsumer,
			LogStdout:     logStdout,
			LogStderr:     logStderr,
		})
		require.NoError(t, err, "failed creating cURL service runner")

		curlResult, err := curlRunner.Start(ctx, true, false)
		require.NoError(t, err, "could not start cURL container")
		require.NotNil(t, curlResult, "could not start cURL container")

		// Retrieve the PKI mount CA cert and copy it over to the cURL container.
		mountCACert, err := pki.GetCACertPEM()
		require.NoError(t, err, "failed to retrieve PKI mount CA certificate")

		mountCACertCtx := hDocker.NewBuildContext()
		mountCACertCtx["ca_cert.crt"] = hDocker.PathContentsFromString(mountCACert)
		err = curlRunner.CopyTo(curlResult.Container.ID, "/tmp/", mountCACertCtx)
		require.NoError(t, err, "failed to copy PKI mount CA certificate to cURL container")

		// Use cURL to hit the Caddy server and validate that a certificate was retrieved successfully.
		curlCmd := []string{
			"curl",
			"-L",
			"--cacert", "/tmp/ca_cert.crt",
			"--resolve", hostname + ":443:" + ipAddr,
			"https://" + hostname + "/",
		}
		stdout, stderr, retcode, err = curlRunner.RunCmdWithOutput(ctx, curlResult.Container.ID, curlCmd)
		t.Logf("cURL Command: %v\nstdout: %v\nstderr: %v\n", curlCmd, string(stdout), string(stderr))
		require.NoError(t, err, "got error running cURL command")
		require.Equal(t, 0, retcode, "expected zero retcode cURL command result")
	}
}

func SubtestACMECertbot(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("pki")
	require.NoError(t, err, "failed setting up acme mount")

	directory := "https://" + pki.GetActiveContainerIP() + ":8200/v1/pki/acme/directory"
	vaultNetwork := pki.GetContainerNetworkName()

	logConsumer, logStdout, logStderr := getDockerLog(t)

	// Default to 45 second timeout, but bump to 120 when running locally or if nightly regression
	// flag is provided.
	sleepTimer := "45"
	if testhelpers.IsLocalOrRegressionTests() {
		sleepTimer = "120"
	}

	t.Logf("creating on network: %v", vaultNetwork)
	runner, err := hDocker.NewServiceRunner(hDocker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/certbot/certbot",
		ImageTag:      "latest",
		ContainerName: "vault_pki_certbot_test",
		NetworkName:   vaultNetwork,
		Entrypoint:    []string{"sleep", sleepTimer},
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

	// Sinkhole a domain that's invalid just in case it's registered in the future.
	cluster.Dns.AddDomain("armoncorp.com")
	cluster.Dns.AddRecord("armoncorp.com", "A", "127.0.0.1")

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

	// N.B. We're using the `certonly` subcommand here because it seems as though the `renew` command
	// attempts to install the cert for you. This ends up hanging and getting killed by docker, but is
	// also not desired behavior. The certbot docs suggest using `certonly` to renew as seen here:
	// https://eff-certbot.readthedocs.io/en/stable/using.html#renewing-certificates
	certbotRenewCmd := []string{
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
		"--cert-name", hostname,
		"--force-renewal",
	}

	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotRenewCmd)
	t.Logf("Certbot Renew Command: %v\nstdout: %v\nstderr: %v\n", certbotRenewCmd, string(stdout), string(stderr))
	if err != nil || retcode != 0 {
		logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
		t.Logf("Certbot logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
	}
	require.NoError(t, err, "got error running renew command")
	require.Equal(t, 0, retcode, "expected zero retcode renew command result")

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

	// Attempt to issue against a domain that doesn't match the challenge.
	// N.B. This test only runs locally or when the nightly regression env var is provided to CI.
	if testhelpers.IsLocalOrRegressionTests() {
		certbotInvalidIssueCmd := []string{
			"certbot",
			"certonly",
			"--no-eff-email",
			"--email", "certbot.client@dadgarcorp.com",
			"--agree-tos",
			"--no-verify-ssl",
			"--standalone",
			"--non-interactive",
			"--server", directory,
			"-d", "armoncorp.com",
			"--issuance-timeout", "10",
		}

		stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotInvalidIssueCmd)
		t.Logf("Certbot Invalid Issue Command: %v\nstdout: %v\nstderr: %v\n", certbotInvalidIssueCmd, string(stdout), string(stderr))
		if err != nil || retcode != 0 {
			logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
			t.Logf("Certbot logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
		}
		require.NoError(t, err, "got error running issue command")
		require.NotEqual(t, 0, retcode, "expected non-zero retcode issue command result")
	}

	// Attempt to close out our ACME account
	certbotUnregisterCmd := []string{
		"certbot",
		"unregister",
		"--no-verify-ssl",
		"--non-interactive",
		"--server", directory,
	}

	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotUnregisterCmd)
	t.Logf("Certbot Unregister Command: %v\nstdout: %v\nstderr: %v\n", certbotUnregisterCmd, string(stdout), string(stderr))
	if err != nil || retcode != 0 {
		logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
		t.Logf("Certbot logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
	}
	require.NoError(t, err, "got error running unregister command")
	require.Equal(t, 0, retcode, "expected zero retcode unregister command result")

	// Attempting to close out our ACME account twice should fail
	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotUnregisterCmd)
	t.Logf("Certbot double Unregister Command: %v\nstdout: %v\nstderr: %v\n", certbotUnregisterCmd, string(stdout), string(stderr))
	if err != nil || retcode != 0 {
		logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
		t.Logf("Certbot double logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
	}
	require.NoError(t, err, "got error running double unregister command")
	require.Equal(t, 1, retcode, "expected non-zero retcode double unregister command result")
}

func SubtestACMECertbotEab(t *testing.T, cluster *VaultPkiCluster) {
	mountName := "pki-certbot-eab"
	pki, err := cluster.CreateAcmeMount(mountName)
	require.NoError(t, err, "failed setting up acme mount")

	err = pki.UpdateAcmeConfig(true, map[string]interface{}{
		"eab_policy": "new-account-required",
	})
	require.NoError(t, err)

	eabId, base64EabKey, err := pki.GetEabKey("acme/")

	directory := "https://" + pki.GetActiveContainerIP() + ":8200/v1/" + mountName + "/acme/directory"
	vaultNetwork := pki.GetContainerNetworkName()

	logConsumer, logStdout, logStderr := getDockerLog(t)

	t.Logf("creating on network: %v", vaultNetwork)
	runner, err := hDocker.NewServiceRunner(hDocker.RunOptions{
		ImageRepo:     "docker.mirror.hashicorp.services/certbot/certbot",
		ImageTag:      "latest",
		ContainerName: "vault_pki_certbot_eab_test",
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
	hostname := "certbot-eab-acme-client.dadgarcorp.com"

	err = pki.AddHostname(hostname, ipAddr)
	require.NoError(t, err, "failed to update vault host files")

	certbotCmd := []string{
		"certbot",
		"certonly",
		"--no-eff-email",
		"--email", "certbot.client@dadgarcorp.com",
		"--eab-kid", eabId,
		"--eab-hmac-key='" + base64EabKey + "'",
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

	certbotRenewCmd := []string{
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
		"--cert-name", hostname,
		"--force-renewal",
	}

	stdout, stderr, retcode, err = runner.RunCmdWithOutput(ctx, result.Container.ID, certbotRenewCmd)
	t.Logf("Certbot Renew Command: %v\nstdout: %v\nstderr: %v\n", certbotRenewCmd, string(stdout), string(stderr))
	if err != nil || retcode != 0 {
		logsStdout, logsStderr, _, _ := runner.RunCmdWithOutput(ctx, result.Container.ID, logCatCmd)
		t.Logf("Certbot logs\nstdout: %v\nstderr: %v\n", string(logsStdout), string(logsStderr))
	}
	require.NoError(t, err, "got error running renew command")
	require.Equal(t, 0, retcode, "expected zero retcode renew command result")

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
		"allow_subdomains":            true,
		"allow_wildcard_certificates": false,
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
		"use_csr_common_name": false,
		"require_cn":          false,
		"client_flag":         false,
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
		"client_flag":                 false,
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
		"client_flag":                 false,
	})
	require.NoError(t, err, "failed creating role wildcard")
	directoryUrl = basePath + "/roles/ica/acme/directory"

	doAcmeValidationWithGoLibrary(t, directoryUrl, acmeOrderIdentifiers, cr, provisioningFunc, "refusing to accept CSR with Basic Constraints extension")
	pki.RemoveDNSRecordsForDomain(hostname)
}

// SubtestACMEStepDownNode Verify that we can properly run an ACME session through a
// secondary node, and midway through the challenge verification process, seal the
// active node and make sure we can complete the ACME session on the new active node.
func SubtestACMEStepDownNode(t *testing.T, cluster *VaultPkiCluster) {
	pki, err := cluster.CreateAcmeMount("stepdown-test")
	require.NoError(t, err)

	// Since we interact with ACME from outside the container network the ACME
	// configuration needs to be updated to use the host port and not the internal
	// docker ip. We also grab the non-active node here on purpose to verify
	// ACME related APIs are properly forwarded across standby hosts.
	nonActiveNodes := pki.GetNonActiveNodes()
	require.GreaterOrEqual(t, len(nonActiveNodes), 1, "Need at least one non-active node")

	nonActiveNode := nonActiveNodes[0]

	basePath := fmt.Sprintf("https://%s/v1/%s", nonActiveNode.HostPort, pki.mount)
	err = pki.UpdateClusterConfig(map[string]interface{}{
		"path": basePath,
	})

	hostname := "go-lang-stepdown-client.dadgarcorp.com"

	acmeOrderIdentifiers := []acme.AuthzID{
		{Type: "dns", Value: hostname},
	}
	cr := &x509.CertificateRequest{
		DNSNames: []string{hostname, hostname},
	}

	accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa account key")

	acmeClient := &acme.Client{
		Key: accountKey,
		HTTPClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
		DirectoryURL: basePath + "/acme/directory",
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

	require.Len(t, order.AuthzURLs, 1, "expected a single authz url")
	authUrl := order.AuthzURLs[0]

	authorization, err := acmeClient.GetAuthorization(testCtx, authUrl)
	require.NoError(t, err, "failed to lookup authorization at url: %s", authUrl)

	dnsTxtRecordsToAdd := map[string]string{}

	var challengesToAccept []*acme.Challenge
	for _, challenge := range authorization.Challenges {
		if challenge.Status != acme.StatusPending {
			t.Logf("ignoring challenge not in status pending: %v", challenge)
			continue
		}

		if challenge.Type == "dns-01" {
			challengeBody, err := acmeClient.DNS01ChallengeRecord(challenge.Token)
			require.NoError(t, err, "failed generating challenge response")

			// Collect the challenges for us to add the DNS records after step-down
			dnsTxtRecordsToAdd["_acme-challenge."+authorization.Identifier.Value] = challengeBody
			challengesToAccept = append(challengesToAccept, challenge)
		}
	}

	// Tell the ACME server, that they can now validate those challenges, this will cause challenge
	// verification failures on the main node as the DNS records do not exist.
	for _, challenge := range challengesToAccept {
		_, err = acmeClient.Accept(testCtx, challenge)
		require.NoError(t, err, "failed to accept challenge: %v", challenge)
	}

	// Now wait till we start seeing the challenge engine start failing the lookups.
	testhelpers.RetryUntil(t, 10*time.Second, func() error {
		myAuth, err := acmeClient.GetAuthorization(testCtx, authUrl)
		require.NoError(t, err, "failed to lookup authorization at url: %s", authUrl)

		for _, challenge := range myAuth.Challenges {
			if challenge.Error != nil {
				// The engine failed on one of the challenges, we are done waiting
				return nil
			}
		}

		return fmt.Errorf("no challenges for auth %v contained any errors", myAuth.Identifier)
	})

	// Seal the active node now and wait for the next node to appear
	previousActiveNode := pki.GetActiveClusterNode()
	t.Logf("Stepping down node id: %s", previousActiveNode.NodeID)

	haStatus, _ := previousActiveNode.APIClient().Sys().HAStatus()
	t.Logf("Node: %v HaStatus: %v\n", previousActiveNode.NodeID, haStatus)

	testhelpers.RetryUntil(t, 2*time.Minute, func() error {
		state, err := previousActiveNode.APIClient().Sys().RaftAutopilotState()
		if err != nil {
			return err
		}

		t.Logf("Node: %v Raft AutoPilotState: %v\n", previousActiveNode.NodeID, state)

		if !state.Healthy {
			return fmt.Errorf("raft auto pilot state is not healthy")
		}

		// Make sure that we have at least one node that can take over prior to sealing the current active node.
		if state.FailureTolerance < 1 {
			msg := fmt.Sprintf("there is no fault tolerance within raft state yet: %d", state.FailureTolerance)
			t.Log(msg)
			return errors.New(msg)
		}

		return nil
	})

	t.Logf("Sealing active node")
	err = previousActiveNode.APIClient().Sys().Seal()
	require.NoError(t, err, "failed stepping down node")

	// Add our DNS records now
	t.Logf("Adding DNS records")
	for dnsHost, dnsValue := range dnsTxtRecordsToAdd {
		err = pki.AddDNSRecord(dnsHost, "TXT", dnsValue)
		require.NoError(t, err, "failed adding DNS record: %s:%s", dnsHost, dnsValue)
	}

	// Wait for our new active node to come up
	testhelpers.RetryUntil(t, 2*time.Minute, func() error {
		newNode := pki.GetActiveClusterNode()
		if newNode.NodeID == previousActiveNode.NodeID {
			return fmt.Errorf("existing node is still the leader after stepdown: %s", newNode.NodeID)
		}

		t.Logf("New active node has node id: %v", newNode.NodeID)
		return nil
	})

	// Wait for the order/challenges to be validated.
	_, err = acmeClient.WaitOrder(testCtx, order.URI)
	if err != nil {
		// We failed waiting for the order to become ready, lets print out current challenge statuses to help debugging
		myAuth, authErr := acmeClient.GetAuthorization(testCtx, authUrl)
		require.NoError(t, authErr, "failed to lookup authorization at url: %s and wait order failed with: %v", authUrl, err)

		t.Logf("Authorization Status: %s", myAuth.Status)
		for _, challenge := range myAuth.Challenges {
			// The engine failed on one of the challenges, we are done waiting
			t.Logf("challenge: %v state: %v Error: %v", challenge.Type, challenge.Status, challenge.Error)
		}

		require.NoError(t, err, "failed waiting for order to be ready")
	}

	// Create/sign the CSR and ask ACME server to sign it returning us the final certificate
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	csr, err := x509.CreateCertificateRequest(rand.Reader, cr, csrKey)
	require.NoError(t, err, "failed generating csr")

	certs, _, err := acmeClient.CreateOrderCert(testCtx, order.FinalizeURL, csr, false)
	require.NoError(t, err, "failed to get a certificate back from ACME")

	_, err = x509.ParseCertificate(certs[0])
	require.NoError(t, err, "failed parsing acme cert bytes")
}

func getDockerLog(t *testing.T) (func(s string), *pkiext.LogConsumerWriter, *pkiext.LogConsumerWriter) {
	logConsumer := func(s string) {
		t.Logf(s)
	}

	logStdout := &pkiext.LogConsumerWriter{logConsumer}
	logStderr := &pkiext.LogConsumerWriter{logConsumer}
	return logConsumer, logStdout, logStderr
}
