package api_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func TestSysRekey_Verification(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client
	client.SetMaxRetries(0)

	// This first block verifies that if we are using recovery keys to force a
	// rekey of a stored-shares barrier that verification is not allowed since
	// the keys aren't returned
	vault.DefaultSealPretendsToAllowRecoveryKeys = true
	vault.DefaultSealPretendsToAllowStoredShares = true
	status, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
		StoredShares:        1,
		RequireVerification: true,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "requiring verification not supported") {
		t.Fatalf("unexpected error: %v", err)
	}

	// Now we set things back and start a normal rekey with the verification process
	vault.DefaultSealPretendsToAllowRecoveryKeys = false
	vault.DefaultSealPretendsToAllowStoredShares = false

	var verificationNonce string
	var newKeys []string
	doRekeyInitialSteps := func() {
		status, err = client.Sys().RekeyInit(&api.RekeyInitRequest{
			SecretShares:        5,
			SecretThreshold:     3,
			RequireVerification: true,
		})
		if err != nil {
			t.Fatal(err)
		}
		if status == nil {
			t.Fatal("nil status")
		}
		if !status.VerificationRequired {
			t.Fatal("expected verification required")
		}

		var resp *api.RekeyUpdateResponse
		for i := 0; i < 3; i++ {
			resp, err = client.Sys().RekeyUpdate(base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i]), status.Nonce)
			if err != nil {
				t.Fatal(err)
			}
		}
		switch {
		case !resp.Complete:
			t.Fatal("expected completion")
		case !resp.VerificationRequired:
			t.Fatal("expected verification required")
		case resp.VerificationNonce == "":
			t.Fatal("verification nonce expected")
		}
		verificationNonce = resp.VerificationNonce
		newKeys = resp.KeysB64
		t.Logf("verification nonce: %q", verificationNonce)
	}

	doRekeyInitialSteps()

	// We are still going, so should not be able to init again
	_, err = client.Sys().RekeyInit(&api.RekeyInitRequest{
		SecretShares:        5,
		SecretThreshold:     3,
		RequireVerification: true,
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Sealing should clear state, so after this we should be able to perform
	// the above again
	cluster.EnsureCoresSealed(t)
	cluster.UnsealCores(t)
	doRekeyInitialSteps()

	doStartVerify := func() {
		// Start the process
		for i := 0; i < 2; i++ {
			status, err := client.Sys().RekeyVerificationUpdate(newKeys[i], verificationNonce)
			if err != nil {
				t.Fatal(err)
			}
			switch {
			case status.Nonce != verificationNonce:
				t.Fatalf("unexpected nonce, expected %q, got %q", verificationNonce, status.Nonce)
			case status.Complete:
				t.Fatal("unexpected completion")
			}
		}

		// Check status
		vStatus, err := client.Sys().RekeyVerificationStatus()
		if err != nil {
			t.Fatal(err)
		}
		switch {
		case vStatus.Nonce != verificationNonce:
			t.Fatalf("unexpected nonce, expected %q, got %q", verificationNonce, vStatus.Nonce)
		case vStatus.T != 3:
			t.Fatal("unexpected threshold")
		case vStatus.N != 5:
			t.Fatal("unexpected number of new keys")
		case vStatus.Progress != 2:
			t.Fatal("unexpected progress")
		}
	}

	doStartVerify()

	// Cancel; this should still keep the rekey process going but just cancel
	// the verification operation
	err = client.Sys().RekeyVerificationCancel()
	if err != nil {
		t.Fatal(err)
	}
	// Verify cannot init again
	_, err = client.Sys().RekeyInit(&api.RekeyInitRequest{
		SecretShares:        5,
		SecretThreshold:     3,
		RequireVerification: true,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	vStatus, err := client.Sys().RekeyVerificationStatus()
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case vStatus.Nonce == verificationNonce:
		t.Fatalf("unexpected nonce, expected not-%q but got it", verificationNonce)
	case vStatus.T != 3:
		t.Fatal("unexpected threshold")
	case vStatus.N != 5:
		t.Fatal("unexpected number of new keys")
	case vStatus.Progress != 0:
		t.Fatal("unexpected progress")
	}

	verificationNonce = vStatus.Nonce
	doStartVerify()

	// Sealing should clear state, but we never actually finished, so it should
	// still be the old keys (which are still currently set)
	cluster.EnsureCoresSealed(t)
	cluster.UnsealCores(t)

	// Should be able to init again and get back to where we were
	doRekeyInitialSteps()
	doStartVerify()

	// Provide the final new key
	vuStatus, err := client.Sys().RekeyVerificationUpdate(newKeys[2], verificationNonce)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case vuStatus.Nonce != verificationNonce:
		t.Fatalf("unexpected nonce, expected %q, got %q", verificationNonce, vuStatus.Nonce)
	case !vuStatus.Complete:
		t.Fatal("expected completion")
	}

	// Seal and unseal -- it should fail to unseal because the key has now been
	// rotated
	cluster.EnsureCoresSealed(t)
	if err := cluster.UnsealCoresWithError(); err == nil {
		t.Fatal("expected error")
	}

	// Swap out the keys with our new ones and try again
	var newKeyBytes [][]byte
	for _, key := range newKeys {
		val, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Fatal(err)
		}
		newKeyBytes = append(newKeyBytes, val)
	}
	cluster.BarrierKeys = newKeyBytes
	if err := cluster.UnsealCoresWithError(); err != nil {
		t.Fatal("expected error")
	}
}
