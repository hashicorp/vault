package pki

import (
	"testing"
	"time"

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
		EnableRaw: true,
	}
	cluster := vault.CreateTestClusterWithRollbackPeriod(t, newPeriod, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
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

	b, s := createBackendWithStorage(t)

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
	require.NoError(t, err)
	require.NotNil(t, statusResp)
	require.NotNil(t, statusResp.Data)
	require.Equal(t, statusResp.Data["state"], "Cancelled")
	nowMany := statusResp.Data["cert_store_deleted_count"].(uint)
	if howMany+3 <= nowMany {
		t.Fatalf("expected to only process at most 3 more certificates, but processed (%v >>> %v) certs", nowMany, howMany)
	}
}
