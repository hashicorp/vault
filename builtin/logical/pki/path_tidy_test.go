package pki

import (
	"encoding/json"
	"errors"
	"fmt"
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
	_, err = CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_cert_store": true,
		"safety_buffer":   "1s",
		"pause_duration":  "1s",
	})

	// If we wait six seconds, the operation should still be running. That's
	// how we check that pause_duration works.
	time.Sleep(3 * time.Second)

	resp, err := CBRead(b, s, "tidy-status")

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
		2*newPeriod,  // A short time period is ideal here to test metrics are emitted every periodic func
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
	mostRecentInterval := inmemSink.Data()[len(inmemSink.Data())-1]
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
	_, err = client.Logical().Write("/sys/plugins/reload/backend", map[string]interface{}{
		"plugin": "pki",
	})

	// By reading the auto-tidy endpoint, we ensure that initialize has completed (which has a write lock on auto-tidy)
	_, err = client.Logical().Read("/pki/config/auto-tidy")
	if err != nil {
		t.Fatal(err)
	}

	// Since publish_stored_certificate_count_metrics is still false, these metrics should still not exist yet
	mostRecentInterval = inmemSink.Data()[len(inmemSink.Data())-1]
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
	if err != nil {
		t.Fatal(err)
	}
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

	// We read the auto-tidy endpoint again, to ensure any metrics logic has completed (lock on config)
	_, err = client.Logical().Read("/pki/config/auto-tidy")
	if err != nil {
		t.Fatal(err)
	}

	// Check Metrics After Cert Has Be Created and Revoked
	tidyStatus, err = client.Logical().Read("pki/tidy-status")
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatalf("Revoked one certificate, but got a revoked cert store count of %v", revokedCertCount)
	}
	// This should now be initialized
	certCountError, ok := tidyStatus.Data["certificate_counting_error"]
	if ok && certCountError.(string) != "" {
		t.Fatalf("Expected certificate count error to disappear after initialization, but got error %v", certCountError)
	}

	testhelpers.RetryUntil(t, newPeriod*5, func() error {
		mostRecentInterval = inmemSink.Data()[len(inmemSink.Data())-1]
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
		mostRecentInterval = inmemSink.Data()[len(inmemSink.Data())-1]
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
