// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSysRekey_Verification(t *testing.T) {
	testcases := []struct {
		recovery     bool
		legacyShamir bool
	}{
		{recovery: true, legacyShamir: false},
		{recovery: false, legacyShamir: false},
		{recovery: false, legacyShamir: true},
	}

	for _, tc := range testcases {
		recovery, legacy := tc.recovery, tc.legacyShamir
		t.Run(fmt.Sprintf("recovery=%v,legacyShamir=%v", recovery, legacy), func(t *testing.T) {
			t.Parallel()
			testSysRekey_Verification(t, recovery, legacy)
		})
	}
}

func testSysRekey_Verification(t *testing.T, recovery bool, legacyShamir bool) {
	opts := &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	}
	switch {
	case recovery:
		if legacyShamir {
			panic("invalid case")
		}
		opts.SealFunc = func() vault.Seal {
			return vault.NewTestSeal(t, &seal.TestSealOpts{
				StoredKeys: seal.StoredKeysSupportedGeneric,
			})
		}
	case legacyShamir:
		opts.SealFunc = func() vault.Seal {
			return vault.NewTestSeal(t, &seal.TestSealOpts{
				StoredKeys: seal.StoredKeysNotSupported,
			})
		}
	}
	inm, err := inmem.NewInmemHA(nil, logging.NewVaultLogger(hclog.Debug))
	if err != nil {
		t.Fatal(err)
	}
	conf := vault.CoreConfig{
		Physical: inm,
	}
	cluster := vault.NewTestCluster(t, &conf, opts)
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

		keys := cluster.BarrierKeys
		if recovery {
			keys = cluster.RecoveryKeys
		}
		var resp *api.RekeyUpdateResponse
		for i := 0; i < 3; i++ {
			resp, err = updateFunc(base64.StdEncoding.EncodeToString(keys[i]), status.Nonce)
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
	_, err = initFunc(&api.RekeyInitRequest{
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
	if err := cluster.UnsealCoresWithError(t, recovery); err != nil {
		t.Fatal(err)
	}
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
		vault.TestWaitActive(t, cluster.Cores[0].Core)

		// Should be able to init again and get back to where we were
		doRekeyInitialSteps()
		doStartVerify()
	} else {
		// We haven't finished, so generating a root token should still be the
		// old keys (which are still currently set)
		testhelpers.GenerateRoot(t, cluster, testhelpers.GenerateRootRegular)
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

		// Simulate restarting Vault rather than just a seal/unseal, because
		// the standbys may not have had time to learn about the new key before
		// we sealed them.  We could sleep, but that's unreliable.
		oldKeys := cluster.BarrierKeys
		opts.SkipInit = true
		opts.SealFunc = nil // post rekey we should use the barrier config on disk
		cluster = vault.NewTestCluster(t, &conf, opts)
		cluster.BarrierKeys = oldKeys
		cluster.Start()
		defer cluster.Cleanup()

		if err := cluster.UnsealCoresWithError(t, false); err == nil {
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
		if err := cluster.UnsealCoresWithError(t, false); err != nil {
			t.Fatal(err)
		}
	} else {
		// The old keys should no longer work
		_, err := testhelpers.GenerateRootWithError(t, cluster, testhelpers.GenerateRootRegular)
		if err == nil {
			t.Fatal("expected error")
		}

		// Put the new keys in place and run again
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
		testhelpers.GenerateRoot(t, cluster, testhelpers.GenerateRootRegular)
	}
}
