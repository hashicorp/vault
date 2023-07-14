// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"golang.org/x/crypto/acme"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"

	"github.com/armon/go-metrics"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"

	"github.com/stretchr/testify/require"
)

func TestTidyConfigs(t *testing.T) {
	t.Parallel()

	var cfg tidyConfig
	operations := strings.Split(cfg.AnyTidyConfig(), " / ")
	t.Logf("Got tidy operations: %v", operations)

	for _, operation := range operations {
		b, s := CreateBackendWithStorage(t)

		resp, err := CBWrite(b, s, "config/auto-tidy", map[string]interface{}{
			"enabled": true,
			operation: true,
		})
		requireSuccessNonNilResponse(t, resp, err, "expected to be able to enable auto-tidy operation "+operation)

		resp, err = CBRead(b, s, "config/auto-tidy")
		requireSuccessNonNilResponse(t, resp, err, "expected to be able to read auto-tidy operation for operation "+operation)
		require.True(t, resp.Data[operation].(bool), "expected operation to be enabled after reading auto-tidy config "+operation)

		resp, err = CBWrite(b, s, "tidy", map[string]interface{}{
			operation: true,
		})
		requireSuccessNonNilResponse(t, resp, err, "expected to be able to start tidy operation with "+operation)
		if len(resp.Warnings) > 0 {
			t.Logf("got warnings while starting manual tidy: %v", resp.Warnings)
			for _, warning := range resp.Warnings {
				if strings.Contains(warning, "Manual tidy requested but no tidy operations were set.") {
					t.Fatalf("expected to be able to enable tidy operation with just %v but got warning: %v / (resp=%v)", operation, warning, resp)
				}
			}
		}
	}
}

func TestAutoTidy(t *testing.T) {
	t.Parallel()

	// While we'd like to reduce this duration, we need to wait until
	// the rollback manager timer ticks. With the new helper, we can
	// modify the rollback manager timer period directly, allowing us
	// to shorten the total test time significantly.
	//
	// We set the delta CRL time to ensure it executes prior to the
	// main CRL rebuild, and the new CRL doesn't rebuild until after
	// we're done.
	newPeriod := 1 * time.Second

	// This test requires the periodicFunc to trigger, which requires we stand
	// up a full test cluster.
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
		// See notes below about usage of /sys/raw for reading cluster
		// storage without barrier encryption.
		EnableRaw:      true,
		RollbackPeriod: newPeriod,
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// Mount PKI
	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "10m",
			MaxLeaseTTL:     "60m",
		},
	})
	require.NoError(t, err)

	// Generate root.
	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root X1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data["issuer_id"])
	issuerId := resp.Data["issuer_id"]

	// Run tidy so status is not empty when we run it later...
	_, err = client.Logical().Write("pki/tidy", map[string]interface{}{
		"tidy_revoked_certs": true,
	})
	require.NoError(t, err)

	// Setup a testing role.
	_, err = client.Logical().Write("pki/roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
	})
	require.NoError(t, err)

	// Write the auto-tidy config.
	_, err = client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"enabled":            true,
		"interval_duration":  "1s",
		"tidy_cert_store":    true,
		"tidy_revoked_certs": true,
		"safety_buffer":      "1s",
	})
	require.NoError(t, err)

	// Issue a cert and revoke it.
	resp, err = client.Logical().Write("pki/issue/local-testing", map[string]interface{}{
		"common_name": "example.com",
		"ttl":         "10s",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["serial_number"])
	require.NotEmpty(t, resp.Data["certificate"])
	leafSerial := resp.Data["serial_number"].(string)
	leafCert := parseCert(t, resp.Data["certificate"].(string))

	// Read cert before revoking
	resp, err = client.Logical().Read("pki/cert/" + leafSerial)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])
	revocationTime, err := (resp.Data["revocation_time"].(json.Number)).Int64()
	require.Equal(t, int64(0), revocationTime, "revocation time was not zero")
	require.Empty(t, resp.Data["revocation_time_rfc3339"], "revocation_time_rfc3339 was not empty")
	require.Empty(t, resp.Data["issuer_id"], "issuer_id was not empty")

	_, err = client.Logical().Write("pki/revoke", map[string]interface{}{
		"serial_number": leafSerial,
	})
	require.NoError(t, err)

	// Cert should still exist.
	resp, err = client.Logical().Read("pki/cert/" + leafSerial)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])
	revocationTime, err = (resp.Data["revocation_time"].(json.Number)).Int64()
	require.NoError(t, err, "failed converting %s to int", resp.Data["revocation_time"])
	revTime := time.Unix(revocationTime, 0)
	now := time.Now()
	if !(now.After(revTime) && now.Add(-10*time.Minute).Before(revTime)) {
		t.Fatalf("parsed revocation time not within the last 10 minutes current time: %s, revocation time: %s", now, revTime)
	}
	utcLoc, err := time.LoadLocation("UTC")
	require.NoError(t, err, "failed to parse UTC location?")

	rfc3339RevocationTime, err := time.Parse(time.RFC3339Nano, resp.Data["revocation_time_rfc3339"].(string))
	require.NoError(t, err, "failed parsing revocation_time_rfc3339 field: %s", resp.Data["revocation_time_rfc3339"])

	require.Equal(t, revTime.In(utcLoc), rfc3339RevocationTime.Truncate(time.Second),
		"revocation times did not match revocation_time: %s, "+"rfc3339 time: %s", revTime, rfc3339RevocationTime)
	require.Equal(t, issuerId, resp.Data["issuer_id"], "issuer_id on leaf cert did not match")

	// Wait for cert to expire and the safety buffer to elapse.
	time.Sleep(time.Until(leafCert.NotAfter) + 3*time.Second)

	// Wait for auto-tidy to run afterwards.
	var foundTidyRunning string
	var foundTidyFinished bool
	timeoutChan := time.After(120 * time.Second)
	for {
		if foundTidyRunning != "" && foundTidyFinished {
			break
		}

		select {
		case <-timeoutChan:
			t.Fatalf("expected auto-tidy to run (%v) and finish (%v) before 120 seconds elapsed", foundTidyRunning, foundTidyFinished)
		default:
			time.Sleep(250 * time.Millisecond)

			resp, err = client.Logical().Read("pki/tidy-status")
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Data)
			require.NotEmpty(t, resp.Data["state"])
			require.NotEmpty(t, resp.Data["time_started"])
			state := resp.Data["state"].(string)
			started := resp.Data["time_started"].(string)
			t.Logf("Resp: %v", resp.Data)

			// We want the _next_ tidy run after the cert expires. This
			// means if we're currently finished when we hit this the
			// first time, we want to wait for the next run.
			if foundTidyRunning == "" {
				foundTidyRunning = started
			} else if foundTidyRunning != started && !foundTidyFinished && state == "Finished" {
				foundTidyFinished = true
			}
		}
	}

	// Cert should no longer exist.
	resp, err = client.Logical().Read("pki/cert/" + leafSerial)
	require.Nil(t, err)
	require.Nil(t, resp)
}

func TestTidyCancellation(t *testing.T) {
	t.Parallel()

	numLeaves := 100

	b, s := CreateBackendWithStorage(t)

	// Create a root, a role, and a bunch of leaves.
	_, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"ttl":         "20m",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	_, err = CBWrite(b, s, "roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
	})
	require.NoError(t, err)
	for i := 0; i < numLeaves; i++ {
		_, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
			"common_name": "testing",
			"ttl":         "1s",
		})
		require.NoError(t, err)
	}

	// Kick off a tidy operation (which runs in the background), but with
	// a slow-ish pause between certificates.
	resp, err := CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_cert_store": true,
		"safety_buffer":   "1s",
		"pause_duration":  "1s",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("tidy"), logical.UpdateOperation), resp, true)

	// If we wait six seconds, the operation should still be running. That's
	// how we check that pause_duration works.
	time.Sleep(3 * time.Second)

	resp, err = CBRead(b, s, "tidy-status")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Equal(t, resp.Data["state"], "Running")

	// If we now cancel the operation, the response should say Cancelling.
	cancelResp, err := CBWrite(b, s, "tidy-cancel", map[string]interface{}{})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("tidy-cancel"), logical.UpdateOperation), resp, true)
	require.NoError(t, err)
	require.NotNil(t, cancelResp)
	require.NotNil(t, cancelResp.Data)
	state := cancelResp.Data["state"].(string)
	howMany := cancelResp.Data["cert_store_deleted_count"].(uint)

	if state == "Cancelled" {
		// Rest of the test can't run; log and exit.
		t.Log("Went to cancel the operation but response was already cancelled")
		return
	}

	require.Equal(t, state, "Cancelling")

	// Wait a little longer, and ensure we only processed at most 2 more certs
	// after the cancellation respon.
	time.Sleep(3 * time.Second)

	statusResp, err := CBRead(b, s, "tidy-status")
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("tidy-status"), logical.ReadOperation), resp, true)
	require.NoError(t, err)
	require.NotNil(t, statusResp)
	require.NotNil(t, statusResp.Data)
	require.Equal(t, statusResp.Data["state"], "Cancelled")
	nowMany := statusResp.Data["cert_store_deleted_count"].(uint)
	if howMany+3 <= nowMany {
		t.Fatalf("expected to only process at most 3 more certificates, but processed (%v >>> %v) certs", nowMany, howMany)
	}
}

func TestTidyIssuers(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	// Create a root that expires quickly and one valid for longer.
	_, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root1 example.com",
		"issuer_name": "root-expired",
		"ttl":         "1s",
		"key_type":    "ec",
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root2 example.com",
		"issuer_name": "root-valid",
		"ttl":         "60m",
		"key_type":    "rsa",
	})
	require.NoError(t, err)

	// Sleep long enough to expire the root.
	time.Sleep(2 * time.Second)

	// First tidy run shouldn't remove anything; too long of safety buffer.
	_, err = CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_expired_issuers": true,
		"issuer_safety_buffer": "60m",
	})
	require.NoError(t, err)

	// Wait for tidy to finish.
	time.Sleep(2 * time.Second)

	// Expired issuer should exist.
	resp, err := CBRead(b, s, "issuer/root-expired")
	requireSuccessNonNilResponse(t, resp, err, "expired should still be present")
	resp, err = CBRead(b, s, "issuer/root-valid")
	requireSuccessNonNilResponse(t, resp, err, "valid should still be present")

	// Second tidy run with shorter safety buffer shouldn't remove the
	// expired one, as it should be the default issuer.
	_, err = CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_expired_issuers": true,
		"issuer_safety_buffer": "1s",
	})
	require.NoError(t, err)

	// Wait for tidy to finish.
	time.Sleep(2 * time.Second)

	// Expired issuer should still exist.
	resp, err = CBRead(b, s, "issuer/root-expired")
	requireSuccessNonNilResponse(t, resp, err, "expired should still be present")
	resp, err = CBRead(b, s, "issuer/root-valid")
	requireSuccessNonNilResponse(t, resp, err, "valid should still be present")

	// Update the default issuer.
	_, err = CBWrite(b, s, "config/issuers", map[string]interface{}{
		"default": "root-valid",
	})
	require.NoError(t, err)

	// Third tidy run should remove the expired one.
	_, err = CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_expired_issuers": true,
		"issuer_safety_buffer": "1s",
	})
	require.NoError(t, err)

	// Wait for tidy to finish.
	time.Sleep(2 * time.Second)

	// Valid issuer should exist still; other should be removed.
	resp, err = CBRead(b, s, "issuer/root-expired")
	require.Error(t, err)
	require.Nil(t, resp)
	resp, err = CBRead(b, s, "issuer/root-valid")
	requireSuccessNonNilResponse(t, resp, err, "valid should still be present")

	// Finally, one more tidy should cause no changes.
	_, err = CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_expired_issuers": true,
		"issuer_safety_buffer": "1s",
	})
	require.NoError(t, err)

	// Wait for tidy to finish.
	time.Sleep(2 * time.Second)

	// Valid issuer should exist still; other should be removed.
	resp, err = CBRead(b, s, "issuer/root-expired")
	require.Error(t, err)
	require.Nil(t, resp)
	resp, err = CBRead(b, s, "issuer/root-valid")
	requireSuccessNonNilResponse(t, resp, err, "valid should still be present")

	// Ensure we have safety buffer and expired issuers set correctly.
	statusResp, err := CBRead(b, s, "tidy-status")
	require.NoError(t, err)
	require.NotNil(t, statusResp)
	require.NotNil(t, statusResp.Data)
	require.Equal(t, statusResp.Data["issuer_safety_buffer"], 1)
	require.Equal(t, statusResp.Data["tidy_expired_issuers"], true)
}

func TestTidyIssuerConfig(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	// Ensure the default auto-tidy config matches expectations
	resp, err := CBRead(b, s, "config/auto-tidy")
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/auto-tidy"), logical.ReadOperation), resp, true)
	requireSuccessNonNilResponse(t, resp, err)

	jsonBlob, err := json.Marshal(&defaultTidyConfig)
	require.NoError(t, err)
	var defaultConfigMap map[string]interface{}
	err = json.Unmarshal(jsonBlob, &defaultConfigMap)
	require.NoError(t, err)

	// Coerce defaults to API response types.
	defaultConfigMap["interval_duration"] = int(time.Duration(defaultConfigMap["interval_duration"].(float64)) / time.Second)
	defaultConfigMap["issuer_safety_buffer"] = int(time.Duration(defaultConfigMap["issuer_safety_buffer"].(float64)) / time.Second)
	defaultConfigMap["safety_buffer"] = int(time.Duration(defaultConfigMap["safety_buffer"].(float64)) / time.Second)
	defaultConfigMap["pause_duration"] = time.Duration(defaultConfigMap["pause_duration"].(float64)).String()
	defaultConfigMap["revocation_queue_safety_buffer"] = int(time.Duration(defaultConfigMap["revocation_queue_safety_buffer"].(float64)) / time.Second)
	defaultConfigMap["acme_account_safety_buffer"] = int(time.Duration(defaultConfigMap["acme_account_safety_buffer"].(float64)) / time.Second)

	require.Equal(t, defaultConfigMap, resp.Data)

	// Ensure setting issuer-tidy related fields stick.
	resp, err = CBWrite(b, s, "config/auto-tidy", map[string]interface{}{
		"tidy_expired_issuers": true,
		"issuer_safety_buffer": "5s",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/auto-tidy"), logical.UpdateOperation), resp, true)

	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, true, resp.Data["tidy_expired_issuers"])
	require.Equal(t, 5, resp.Data["issuer_safety_buffer"])
}

// TestCertStorageMetrics ensures that when enabled, metrics are able to count the number of certificates in storage and
// number of revoked certificates in storage.  Moreover, this test ensures that the gauge is emitted periodically, so
// that the metric does not disappear or go stale.
func TestCertStorageMetrics(t *testing.T) {
	// This tests uses the same setup as TestAutoTidy
	newPeriod := 1 * time.Second

	// We set up a metrics accumulator
	inmemSink := metrics.NewInmemSink(
		2*newPeriod, // A short time period is ideal here to test metrics are emitted every periodic func
		10*newPeriod) // Do not keep a huge amount of metrics in the sink forever, clear them out to save memory usage.

	metricsConf := metrics.DefaultConfig("")
	metricsConf.EnableHostname = false
	metricsConf.EnableHostnameLabel = false
	metricsConf.EnableServiceLabel = false
	metricsConf.EnableTypePrefix = false

	_, err := metrics.NewGlobal(metricsConf, inmemSink)
	if err != nil {
		t.Fatal(err)
	}

	// This test requires the periodicFunc to trigger, which requires we stand
	// up a full test cluster.
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
		// See notes below about usage of /sys/raw for reading cluster
		// storage without barrier encryption.
		EnableRaw:      true,
		RollbackPeriod: newPeriod,
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// Mount PKI
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "10m",
			MaxLeaseTTL:     "60m",
		},
	})
	require.NoError(t, err)

	// Generate root.
	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root X1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data["issuer_id"])

	// Set up a testing role.
	_, err = client.Logical().Write("pki/roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
	})
	require.NoError(t, err)

	// Run tidy so that tidy-status is not empty
	_, err = client.Logical().Write("pki/tidy", map[string]interface{}{
		"tidy_revoked_certs": true,
	})
	require.NoError(t, err)

	// Since certificate counts are off by default, we shouldn't see counts in the tidy status
	tidyStatus, err := client.Logical().Read("pki/tidy-status")
	if err != nil {
		t.Fatal(err)
	}
	// backendUUID should exist, we need this for metrics
	backendUUID := tidyStatus.Data["internal_backend_uuid"].(string)
	// "current_cert_store_count", "current_revoked_cert_count"
	countData, ok := tidyStatus.Data["current_cert_store_count"]
	if ok && countData != nil {
		t.Fatalf("Certificate counting should be off by default, but current cert store count %v appeared in tidy status in unconfigured mount", countData)
	}
	revokedCountData, ok := tidyStatus.Data["current_revoked_cert_count"]
	if ok && revokedCountData != nil {
		t.Fatalf("Certificate counting should be off by default, but revoked cert count %v appeared in tidy status in unconfigured mount", revokedCountData)
	}

	// Since certificate counts are off by default, those metrics should not exist yet
	stableMetric := inmemSink.Data()
	mostRecentInterval := stableMetric[len(stableMetric)-1]
	_, ok = mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_revoked_certificates_stored"]
	if ok {
		t.Fatalf("Certificate counting should be off by default, but revoked cert count was emitted as a metric in an unconfigured mount")
	}
	_, ok = mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_certificates_stored"]
	if ok {
		t.Fatalf("Certificate counting should be off by default, but total certificate count was emitted as a metric in an unconfigured mount")
	}

	// Write the auto-tidy config.
	_, err = client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"enabled":                                  true,
		"interval_duration":                        "1s",
		"tidy_cert_store":                          true,
		"tidy_revoked_certs":                       true,
		"safety_buffer":                            "1s",
		"maintain_stored_certificate_counts":       true,
		"publish_stored_certificate_count_metrics": false,
	})
	require.NoError(t, err)

	// Reload the Mount - Otherwise Stored Certificate Counts Will Not Be Populated
	// Sealing cores as plugin reload triggers the race detector - VAULT-13635
	testhelpers.EnsureCoresSealed(t, cluster)
	testhelpers.EnsureCoresUnsealed(t, cluster)

	// Wait until a tidy run has completed.
	testhelpers.RetryUntil(t, 5*time.Second, func() error {
		resp, err = client.Logical().Read("pki/tidy-status")
		if err != nil {
			return fmt.Errorf("error reading tidy status: %w", err)
		}
		if finished, ok := resp.Data["time_finished"]; !ok || finished == "" || finished == nil {
			return fmt.Errorf("tidy time_finished not run yet: %v", finished)
		}
		return nil
	})

	// Since publish_stored_certificate_count_metrics is still false, these metrics should still not exist yet
	stableMetric = inmemSink.Data()
	mostRecentInterval = stableMetric[len(stableMetric)-1]
	_, ok = mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_revoked_certificates_stored"]
	if ok {
		t.Fatalf("Certificate counting should be off by default, but revoked cert count was emitted as a metric in an unconfigured mount")
	}
	_, ok = mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_certificates_stored"]
	if ok {
		t.Fatalf("Certificate counting should be off by default, but total certificate count was emitted as a metric in an unconfigured mount")
	}

	// But since certificate counting is on, the metrics should exist on tidyStatus endpoint:
	tidyStatus, err = client.Logical().Read("pki/tidy-status")
	require.NoError(t, err, "failed reading tidy-status endpoint")

	// backendUUID should exist, we need this for metrics
	backendUUID = tidyStatus.Data["internal_backend_uuid"].(string)
	// "current_cert_store_count", "current_revoked_cert_count"
	certStoreCount, ok := tidyStatus.Data["current_cert_store_count"]
	if !ok {
		t.Fatalf("Certificate counting has been turned on, but current cert store count does not appear in tidy status")
	}
	if certStoreCount != json.Number("1") {
		t.Fatalf("Only created one certificate, but a got a certificate count of %v", certStoreCount)
	}
	revokedCertCount, ok := tidyStatus.Data["current_revoked_cert_count"]
	if !ok {
		t.Fatalf("Certificate counting has been turned on, but revoked cert store count does not appear in tidy status")
	}
	if revokedCertCount != json.Number("0") {
		t.Fatalf("Have not yet revoked a certificate, but got a revoked cert store count of %v", revokedCertCount)
	}

	// Write the auto-tidy config, again, this time turning on metrics
	_, err = client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"enabled":                                  true,
		"interval_duration":                        "1s",
		"tidy_cert_store":                          true,
		"tidy_revoked_certs":                       true,
		"safety_buffer":                            "1s",
		"maintain_stored_certificate_counts":       true,
		"publish_stored_certificate_count_metrics": true,
	})
	require.NoError(t, err, "failed updating auto-tidy configuration")

	// Issue a cert and revoke it.
	resp, err = client.Logical().Write("pki/issue/local-testing", map[string]interface{}{
		"common_name": "example.com",
		"ttl":         "10s",
	})
	require.NoError(t, err, "failed to issue leaf certificate")
	require.NotNil(t, resp, "nil response without error on issuing leaf certificate")
	require.NotNil(t, resp.Data, "empty Data without error on issuing leaf certificate")
	require.NotEmpty(t, resp.Data["serial_number"])
	require.NotEmpty(t, resp.Data["certificate"])
	leafSerial := resp.Data["serial_number"].(string)
	leafCert := parseCert(t, resp.Data["certificate"].(string))

	// Read cert before revoking
	resp, err = client.Logical().Read("pki/cert/" + leafSerial)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])
	revocationTime, err := (resp.Data["revocation_time"].(json.Number)).Int64()
	require.Equal(t, int64(0), revocationTime, "revocation time was not zero")
	require.Empty(t, resp.Data["revocation_time_rfc3339"], "revocation_time_rfc3339 was not empty")
	require.Empty(t, resp.Data["issuer_id"], "issuer_id was not empty")

	revokeResp, err := client.Logical().Write("pki/revoke", map[string]interface{}{
		"serial_number": leafSerial,
	})
	require.NoError(t, err, "failed revoking serial number: %s", leafSerial)

	for _, warning := range revokeResp.Warnings {
		if strings.Contains(warning, "already expired; refusing to add to CRL") {
			t.Skipf("Skipping test as we missed the revocation window of our leaf cert")
		}
	}

	// We read the auto-tidy endpoint again, to ensure any metrics logic has completed (lock on config)
	_, err = client.Logical().Read("/pki/config/auto-tidy")
	require.NoError(t, err, "failed to read auto-tidy configuration")

	// Check Metrics After Cert Has Be Created and Revoked
	tidyStatus, err = client.Logical().Read("pki/tidy-status")
	require.NoError(t, err, "failed to read tidy-status")

	backendUUID = tidyStatus.Data["internal_backend_uuid"].(string)
	certStoreCount, ok = tidyStatus.Data["current_cert_store_count"]
	if !ok {
		t.Fatalf("Certificate counting has been turned on, but current cert store count does not appear in tidy status")
	}
	if certStoreCount != json.Number("2") {
		t.Fatalf("Created root and leaf certificate, but a got a certificate count of %v", certStoreCount)
	}
	revokedCertCount, ok = tidyStatus.Data["current_revoked_cert_count"]
	if !ok {
		t.Fatalf("Certificate counting has been turned on, but revoked cert store count does not appear in tidy status")
	}
	if revokedCertCount != json.Number("1") {
		t.Fatalf("Revoked one certificate, but got a revoked cert store count of %v\n:%v", revokedCertCount, tidyStatus)
	}
	// This should now be initialized
	certCountError, ok := tidyStatus.Data["certificate_counting_error"]
	if ok && certCountError.(string) != "" {
		t.Fatalf("Expected certificate count error to disappear after initialization, but got error %v", certCountError)
	}

	testhelpers.RetryUntil(t, newPeriod*5, func() error {
		stableMetric = inmemSink.Data()
		mostRecentInterval = stableMetric[len(stableMetric)-1]
		revokedCertCountGaugeValue, ok := mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_revoked_certificates_stored"]
		if !ok {
			return errors.New("turned on metrics, but revoked cert count was not emitted")
		}
		if revokedCertCountGaugeValue.Value != 1 {
			return fmt.Errorf("revoked one certificate, but metrics emitted a revoked cert store count of %v", revokedCertCountGaugeValue)
		}
		certStoreCountGaugeValue, ok := mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_certificates_stored"]
		if !ok {
			return errors.New("turned on metrics, but total certificate count was not emitted")
		}
		if certStoreCountGaugeValue.Value != 2 {
			return fmt.Errorf("stored two certificiates, but total certificate count emitted was %v", certStoreCountGaugeValue.Value)
		}
		return nil
	})

	// Wait for cert to expire and the safety buffer to elapse.
	sleepFor := time.Until(leafCert.NotAfter) + 3*time.Second
	t.Logf("%v: Sleeping for %v, leaf certificate expires: %v", time.Now().Format(time.RFC3339), sleepFor, leafCert.NotAfter)
	time.Sleep(sleepFor)

	// Wait for auto-tidy to run afterwards.
	var foundTidyRunning string
	var foundTidyFinished bool
	timeoutChan := time.After(120 * time.Second)
	for {
		if foundTidyRunning != "" && foundTidyFinished {
			break
		}

		select {
		case <-timeoutChan:
			t.Fatalf("expected auto-tidy to run (%v) and finish (%v) before 120 seconds elapsed", foundTidyRunning, foundTidyFinished)
		default:
			time.Sleep(250 * time.Millisecond)

			resp, err = client.Logical().Read("pki/tidy-status")
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Data)
			require.NotEmpty(t, resp.Data["state"])
			require.NotEmpty(t, resp.Data["time_started"])
			state := resp.Data["state"].(string)
			started := resp.Data["time_started"].(string)

			t.Logf("%v: Resp: %v", time.Now().Format(time.RFC3339), resp.Data)

			// We want the _next_ tidy run after the cert expires. This
			// means if we're currently finished when we hit this the
			// first time, we want to wait for the next run.
			if foundTidyRunning == "" {
				foundTidyRunning = started
			} else if foundTidyRunning != started && !foundTidyFinished && state == "Finished" {
				foundTidyFinished = true
			}
		}
	}

	// After Tidy, Cert Store Count Should Still Be Available, and Be Updated:
	// Check Metrics After Cert Has Be Created and Revoked
	tidyStatus, err = client.Logical().Read("pki/tidy-status")
	if err != nil {
		t.Fatal(err)
	}
	backendUUID = tidyStatus.Data["internal_backend_uuid"].(string)
	// "current_cert_store_count", "current_revoked_cert_count"
	certStoreCount, ok = tidyStatus.Data["current_cert_store_count"]
	if !ok {
		t.Fatalf("Certificate counting has been turned on, but current cert store count does not appear in tidy status")
	}
	if certStoreCount != json.Number("1") {
		t.Fatalf("Created root and leaf certificate, deleted leaf, but a got a certificate count of %v", certStoreCount)
	}
	revokedCertCount, ok = tidyStatus.Data["current_revoked_cert_count"]
	if !ok {
		t.Fatalf("Certificate counting has been turned on, but revoked cert store count does not appear in tidy status")
	}
	if revokedCertCount != json.Number("0") {
		t.Fatalf("Revoked certificate has been tidied, but got a revoked cert store count of %v", revokedCertCount)
	}

	testhelpers.RetryUntil(t, newPeriod*5, func() error {
		stableMetric = inmemSink.Data()
		mostRecentInterval = stableMetric[len(stableMetric)-1]
		revokedCertCountGaugeValue, ok := mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_revoked_certificates_stored"]
		if !ok {
			return errors.New("turned on metrics, but revoked cert count was not emitted")
		}
		if revokedCertCountGaugeValue.Value != 0 {
			return fmt.Errorf("revoked certificate has been tidied, but metrics emitted a revoked cert store count of %v", revokedCertCountGaugeValue)
		}
		certStoreCountGaugeValue, ok := mostRecentInterval.Gauges["secrets.pki."+backendUUID+".total_certificates_stored"]
		if !ok {
			return errors.New("turned on metrics, but total certificate count was not emitted")
		}
		if certStoreCountGaugeValue.Value != 1 {
			return fmt.Errorf("only one of two certificates left after tidy, but total certificate count emitted was %v", certStoreCountGaugeValue.Value)
		}
		return nil
	})
}

func TestTidyAcme(t *testing.T) {

	// This would still be way easier if I could do both sides
	// b, s := CreateBackendWithStorage(t)
	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()
	testCtx := context.Background()

	// Configure Tidy to Tidy Acme (by Writing the auto-tidy config.)
	_, err := client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"enabled":           true,
		"interval_duration": "1s",
		"tidy_acme":         true,
		"safety_buffer":     "1s",
	})
	require.NoError(t, err)

	// Register an Account, do nothing with it
	baseAcmeURL := "/v1/pki/acme/"
	accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed creating rsa key")

	acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

	// Create new account
	t.Logf("Testing register on %s", baseAcmeURL)
	acct, err := acmeClient.Register(testCtx, &acme.Account{}, func(tosURL string) bool { return true })
	require.NoError(t, err, "failed registering account")

	// Run Tidy
	_, err = client.Logical().Write("pki/tidy", map[string]interface{}{
		"tidy_acme": true,
	})
	require.NoError(t, err)

	// Check that the Account is Still There, Still Valid
	listResp, err := client.Logical().ListWithContext(testCtx, "pki"+acmeThumbprintPrefix)
	require.NotNil(t, listResp)
	thumbprintEntries := listResp.Data["keys"].([]interface{})
	require.Equal(t, len(thumbprintEntries), 1)
	thumbprint := thumbprintEntries[0].(string)

	// Let "Time Pass"; this is a HACK, this function sys-writes to overwrite the date on objects in storage
	backDateAcmeAccountSys(t, testCtx, client, thumbprint, 30*24*time.Hour, "pki")

	// Run Tidy
	_, err = client.Logical().Write("pki/tidy", map[string]interface{}{
		"tidy_acme": true,
	})
	require.NoError(t, err)

	// Let "Time Pass"; this is a HACK, this function sys-writes to overwrite the date on objects in storage
	backDateAcmeAccountSys(t, testCtx, client, thumbprint, 30*24*time.Hour, "pki")

	// Run Tidy
	_, err = client.Logical().Write("pki/tidy", map[string]interface{}{
		"tidy_acme": true,
	})
	require.NoError(t, err)

	// Check Account No Longer Appears
	listResp, err = client.Logical().ListWithContext(testCtx, "pki"+acmeThumbprintPrefix)
	require.NotNil(t, listResp)
	thumbprintEntries = listResp.Data["keys"].([]interface{})
	require.Equal(t, len(thumbprintEntries), 0)

	// Nor Under Account
	getResp, err := client.Logical().ReadWithContext(testCtx, "pki"+acmeAccountPrefix+acct.URI)
	require.Nil(t, getResp)

}

// These calls are for tests using the "storage" and "backend" direct setup
func backDateAcmeAccount(t *testing.T, testContext context.Context, b *backend, s logical.Storage, thumbprintString string, backdateAmount time.Duration) {
	thumbprintPath := acmeThumbprintPrefix + thumbprintString
	thumbprintEntry, err := s.Get(testContext, thumbprintPath)
	if err != nil {
		t.Fatalf("unable to fetch thumbprint entry %v", err)
	}

	var thumbprint acmeThumbprint
	err = thumbprintEntry.DecodeJSON(&thumbprint)
	if err != nil {
		t.Fatalf("unable to decode thumbprint entry %v to find account entry: %v", thumbprintEntry, err)
	}

	accountPath := acmeAccountPrefix + thumbprint.Kid
	accountEntry, err := s.Get(testContext, accountPath)
	if err != nil {
		t.Fatalf("unable to fetch account entry %v: %v", thumbprint.Kid, err)
	}

	var account acmeAccount
	err = accountEntry.DecodeJSON(&account)
	if err != nil {
		t.Fatalf("unable to decode acme account %v: %v", accountEntry, err)
	}

	account.AccountCreatedDate = backDate(account.AccountCreatedDate, backdateAmount)
	account.MaxCertExpiry = backDate(account.MaxCertExpiry, backdateAmount)
	account.AccountRevokedDate = backDate(account.AccountRevokedDate, backdateAmount)

	json, err := logical.StorageEntryJSON(accountPath, &account)
	err = s.Put(testContext, json)
	if err != nil {
		t.Fatalf("error saving backdated account entry at %v: %v", accountPath, err)
	}

	orders, err := s.List(testContext, acmeAccountPrefix+thumbprint.Kid+"/orders/")
	if err != nil {
		t.Fatalf("unable to list orders on account %v", thumbprint.Kid)
	}

	for _, order := range orders {
		backDateAcmeOrder(t, testContext, b, s, thumbprint.Kid, order, backdateAmount)
	}

	// No need to change certificates entries here - no time is stored on AcmeCertEntry

}

func backDateAcmeOrder(t *testing.T, testContext context.Context, b *backend, s logical.Storage, accountKid string, orderId string, backdateAmount time.Duration) {
	orderPath := acmeAccountPrefix + accountKid + "orders" + orderId

	orderEntry, err := s.Get(testContext, orderPath)
	if err != nil {
		t.Fatalf("unable to fetch order entry %v on account %v", orderId, accountKid)
	}

	var order *acmeOrder
	err = orderEntry.DecodeJSON(order)
	if err != nil {
		t.Fatalf("error decoding order entry %v on account %v, %v produced: %v", orderId, accountKid, orderEntry, err)
	}

	order.Expires = backDate(order.Expires, backdateAmount)
	order.CertificateExpiry = backDate(order.CertificateExpiry, backdateAmount)

	json, err := logical.StorageEntryJSON(orderPath, &order)
	err = s.Put(testContext, json)
	if err != nil {
		t.Fatalf("error saving backdated order entry %v on account %v : %v", orderId, accountKid, err)
	}

	for _, authId := range order.AuthorizationIds {
		backDateAcmeAuthorization(t, testContext, s, accountKid, authId, backdateAmount)
	}

}

func backDateAcmeAuthorization(t *testing.T, testContext context.Context, s logical.Storage, accountKid string, authId string, backdateAmount time.Duration) {
	authPath := acmeAccountPrefix + accountKid + "/authorizations/" + authId

	authEntry, err := s.Get(testContext, authPath)
	if err != nil {
		t.Fatalf("unable to fetch authorization %v : %v", authPath, err)
	}

	var auth *ACMEAuthorization
	err = authEntry.DecodeJSON(auth)
	if err != nil {
		t.Fatalf("error decoding auth %v, auth entry %v produced %v", authPath, authEntry, err)
	}

	expiry, err := auth.GetExpires()
	if err != nil {
		t.Fatalf("could not get expiry on %v: %v", authPath, err)
	}
	newExpiry := backDate(expiry, backdateAmount)
	auth.Expires = time.Time.Format(newExpiry, time.RFC3339)

	json, err := logical.StorageEntryJSON(authPath, auth)
	err = s.Put(testContext, json)
	if err != nil {
		t.Fatalf("error updating authorization date on %v: %v", authPath, err)
	}

}

// The sys tests refer to all of the tests using sys/raw/logical which work off of a client
func backDateAcmeAccountSys(t *testing.T, testContext context.Context, client *api.Client, thumbprintString string, backdateAmount time.Duration, mount string) {
	rawThumbprintPath := path.Join("sys/raw/logical/", mount, acmeThumbprintPrefix+thumbprintString)
	thumbprintResp, err := client.Logical().ReadWithContext(testContext, rawThumbprintPath)
	if err != nil {
		t.Fatalf("unable to fetch thumbprint response at %v: %v", rawThumbprintPath, err)
	}

	var thumbprint acmeThumbprint
	err = jsonutil.DecodeJSON([]byte(thumbprintResp.Data["value"].(string)), &thumbprint)
	if err != nil {
		t.Fatalf("unable to decode thumbprint response %v to find account entry: %v", thumbprintResp.Data, err)
	}

	accountPath := path.Join("sys/raw/logical", mount, acmeAccountPrefix+thumbprint.Kid)
	accountResp, err := client.Logical().ReadWithContext(testContext, accountPath)
	if err != nil {
		t.Fatalf("unable to fetch account entry %v: %v", thumbprint.Kid, err)
	}

	var account acmeAccount
	err = jsonutil.DecodeJSON([]byte(accountResp.Data["value"].(string)), &account)
	if err != nil {
		t.Fatalf("unable to decode acme account %v: %v", accountResp, err)
	}

	account.AccountCreatedDate = backDate(account.AccountCreatedDate, backdateAmount)
	account.MaxCertExpiry = backDate(account.MaxCertExpiry, backdateAmount)
	account.AccountRevokedDate = backDate(account.AccountRevokedDate, backdateAmount)

	encodeJSON, err := jsonutil.EncodeJSON(account)
	_, err = client.Logical().WriteWithContext(context.Background(), accountPath, map[string]interface{}{
		"value":    base64.StdEncoding.EncodeToString(encodeJSON),
		"encoding": "base64",
	})
	if err != nil {
		t.Fatalf("error saving backdated account entry at %v: %v", accountPath, err)
	}

	ordersPath := path.Join("sys/raw/logical", mount, acmeAccountPrefix, thumbprint.Kid, "/orders/")
	ordersRaw, err := client.Logical().ListWithContext(context.Background(), ordersPath)
	if err != nil {
		t.Fatalf("unable to list orders on account %v", thumbprint.Kid)
	}

	orders := ordersRaw.Data

	for _, orderInterface := range orders {
		orderId := orderInterface.(string)
		backDateAcmeOrderSys(t, testContext, client, thumbprint.Kid, orderId, backdateAmount, mount)
	}

	// No need to change certificates entries here - no time is stored on AcmeCertEntry

}

func backDateAcmeOrderSys(t *testing.T, testContext context.Context, client *api.Client, accountKid string, orderId string, backdateAmount time.Duration, mount string) {
	rawOrderPath := path.Join("sys/raw/logical/", mount, acmeAccountPrefix, accountKid, "orders", orderId)
	orderResp, err := client.Logical().ReadWithContext(testContext, rawOrderPath)
	if err != nil {
		t.Fatalf("unable to fetch order entry %v on account %v at %v", orderId, accountKid, rawOrderPath)
	}

	var order *acmeOrder
	err = jsonutil.DecodeJSON([]byte(orderResp.Data["value"].(string)), &order)
	if err != nil {
		t.Fatalf("error decoding order entry %v on account %v, %v produced: %v", orderId, accountKid, orderResp, err)
	}

	order.Expires = backDate(order.Expires, backdateAmount)
	order.CertificateExpiry = backDate(order.CertificateExpiry, backdateAmount)

	encodeJSON, err := jsonutil.EncodeJSON(order)
	_, err = client.Logical().WriteWithContext(context.Background(), rawOrderPath, map[string]interface{}{
		"value":    base64.StdEncoding.EncodeToString(encodeJSON),
		"encoding": "base64",
	})
	if err != nil {
		t.Fatalf("error saving backdated order entry %v on account %v : %v", orderId, accountKid, err)
	}

	for _, authId := range order.AuthorizationIds {
		backDateAcmeAuthorizationSys(t, testContext, client, accountKid, authId, backdateAmount, mount)
	}

}

func backDateAcmeAuthorizationSys(t *testing.T, testContext context.Context, client *api.Client, accountKid string, authId string, backdateAmount time.Duration, mount string) {
	rawAuthPath := path.Join("sys/raw/logical/", mount, acmeAccountPrefix, accountKid, "/authorizations/", authId)

	authResp, err := client.Logical().ReadWithContext(testContext, rawAuthPath)
	if err != nil {
		t.Fatalf("unable to fetch authorization %v : %v", rawAuthPath, err)
	}

	var auth *ACMEAuthorization
	err = jsonutil.DecodeJSON([]byte(authResp.Data["value"].(string)), &auth)
	if err != nil {
		t.Fatalf("error decoding auth %v, auth entry %v produced %v", rawAuthPath, authResp, err)
	}

	expiry, err := auth.GetExpires()
	if err != nil {
		t.Fatalf("could not get expiry on %v: %v", rawAuthPath, err)
	}
	newExpiry := backDate(expiry, backdateAmount)
	auth.Expires = time.Time.Format(newExpiry, time.RFC3339)

	encodeJSON, err := jsonutil.EncodeJSON(auth)
	_, err = client.Logical().WriteWithContext(context.Background(), rawAuthPath, map[string]interface{}{
		"value":    base64.StdEncoding.EncodeToString(encodeJSON),
		"encoding": "base64",
	})
	if err != nil {
		t.Fatalf("error updating authorization date on %v: %v", rawAuthPath, err)
	}
}

func backDate(original time.Time, change time.Duration) time.Time {
	if original.IsZero() {
		return original
	}

	zeroTime := time.Time{}

	if original.Before(zeroTime.Add(change)) {
		return zeroTime
	}

	return original.Add(-change)

}
