package api

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/shamir"
	"github.com/hashicorp/vault/vault"
)

func TestSysRekey_Verification(t *testing.T) {
	testSysRekey_Verification(t, false)
	testSysRekey_Verification(t, true)
}

func testSysRekey_Verification(t *testing.T, recovery bool) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client
	client.SetMaxRetries(0)

	initFunc := client.Sys().RekeyInit
	updateFunc := client.Sys().RekeyUpdate
	verificationUpdateFunc := client.Sys().RekeyVerificationUpdate
	verificationStatusFunc := client.Sys().RekeyVerificationStatus
	verificationCancelFunc := client.Sys().RekeyVerificationCancel
	if recovery {
		initFunc = client.Sys().RekeyRecoveryKeyInit
		updateFunc = client.Sys().RekeyRecoveryKeyUpdate
		verificationUpdateFunc = client.Sys().RekeyRecoveryKeyVerificationUpdate
		verificationStatusFunc = client.Sys().RekeyRecoveryKeyVerificationStatus
		verificationCancelFunc = client.Sys().RekeyRecoveryKeyVerificationCancel
	}

	sealAccess := cluster.Cores[0].Core.SealAccess()
	sealTestingParams := &vault.SealAccessTestingParams{}

	// This first block verifies that if we are using recovery keys to force a
	// rekey of a stored-shares barrier that verification is not allowed since
	// the keys aren't returned
	if !recovery {
		sealTestingParams.PretendToAllowRecoveryKeys = true
		sealTestingParams.PretendToAllowStoredShares = true
		if err := sealAccess.SetTestingParams(sealTestingParams); err != nil {
			t.Fatal(err)
		}

		_, err := initFunc(&api.RekeyInitRequest{
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
		sealTestingParams.PretendToAllowRecoveryKeys = false
		sealTestingParams.PretendToAllowStoredShares = false
		if err := sealAccess.SetTestingParams(sealTestingParams); err != nil {
			t.Fatal(err)
		}
	} else {
		cluster.RecoveryKeys = cluster.BarrierKeys
		sealTestingParams.PretendToAllowRecoveryKeys = true
		recoveryKey, err := shamir.Combine(cluster.BarrierKeys)
		if err != nil {
			t.Fatal(err)
		}
		sealTestingParams.PretendRecoveryKey = recoveryKey
		if err := sealAccess.SetTestingParams(sealTestingParams); err != nil {
			t.Fatal(err)
		}
	}

	var verificationNonce string
	var newKeys []string
	doRekeyInitialSteps := func() {
		status, err := initFunc(&api.RekeyInitRequest{
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
			resp, err = updateFunc(base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i]), status.Nonce)
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
	_, err := initFunc(&api.RekeyInitRequest{
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
			status, err := verificationUpdateFunc(newKeys[i], verificationNonce)
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
		vStatus, err := verificationStatusFunc()
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
	err = verificationCancelFunc()
	if err != nil {
		t.Fatal(err)
	}
	// Verify cannot init again
	_, err = initFunc(&api.RekeyInitRequest{
		SecretShares:        5,
		SecretThreshold:     3,
		RequireVerification: true,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	vStatus, err := verificationStatusFunc()
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

	if !recovery {
		// Sealing should clear state, but we never actually finished, so it should
		// still be the old keys (which are still currently set)
		cluster.EnsureCoresSealed(t)
		cluster.UnsealCores(t)

		// Should be able to init again and get back to where we were
		doRekeyInitialSteps()
		doStartVerify()
	} else {
		// We haven't finished, so generating a root token should still be the
		// old keys (which are still currently set)
		testhelpers.GenerateRoot(t, cluster, false)
	}

	// Provide the final new key
	vuStatus, err := verificationUpdateFunc(newKeys[2], verificationNonce)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case vuStatus.Nonce != verificationNonce:
		t.Fatalf("unexpected nonce, expected %q, got %q", verificationNonce, vuStatus.Nonce)
	case !vuStatus.Complete:
		t.Fatal("expected completion")
	}

	if !recovery {
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
	} else {
		// The old keys should no longer work
		_, err := testhelpers.GenerateRootWithError(t, cluster, false)
		if err == nil {
			t.Fatal("expected error")
		}

		// Put tne new keys in place and run again
		cluster.RecoveryKeys = nil
		for _, key := range newKeys {
			dec, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				t.Fatal(err)
			}
			cluster.RecoveryKeys = append(cluster.RecoveryKeys, dec)
		}
		if err := client.Sys().GenerateRootCancel(); err != nil {
			t.Fatal(err)
		}
		testhelpers.GenerateRoot(t, cluster, false)
	}
}
