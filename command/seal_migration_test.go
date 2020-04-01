package command

import (
	"context"
	"testing"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	vaultseal "github.com/hashicorp/vault/vault/seal"
)

func verifyBarrierConfig(t *testing.T, cfg *vault.SealConfig, sealType string, shares, threshold, stored int) {
	t.Helper()
	if cfg.Type != sealType {
		t.Fatalf("bad seal config: %#v, expected type=%q", cfg, sealType)
	}
	if cfg.SecretShares != shares {
		t.Fatalf("bad seal config: %#v, expected SecretShares=%d", cfg, shares)
	}
	if cfg.SecretThreshold != threshold {
		t.Fatalf("bad seal config: %#v, expected SecretThreshold=%d", cfg, threshold)
	}
	if cfg.StoredShares != stored {
		t.Fatalf("bad seal config: %#v, expected StoredShares=%d", cfg, stored)
	}
}

func TestSealMigration_ShamirToTransit(t *testing.T) {
	t.Parallel()
	t.Run("inmem", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTransit(t, teststorage.InmemBackendSetup)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTransit(t, teststorage.FileBackendSetup)
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTransit(t, teststorage.ConsulBackendSetup)
	})

	t.Run("raft", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTransit(t, teststorage.RaftBackendSetup)
	})
}

func testSealMigrationShamirToTransit(t *testing.T, setup teststorage.ClusterSetupMutator) {

	// Create a cluster that uses shamir.
	conf, opts := teststorage.ClusterSetup(&vault.CoreConfig{
		DisableSealWrap: true,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    3,
	},
		setup,
	)
	opts.SetupFunc = nil
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	// Initialize the cluster, and unseal it using the shamir keys.
	client := cluster.Cores[0].Client
	initResp, err := client.Sys().Init(&api.InitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	var resp *api.SealStatusResponse
	for _, key := range initResp.KeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if resp == nil || resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", resp)
	}

	testhelpers.WaitForActiveNode(t, cluster)

	rootToken := initResp.RootToken
	client.SetToken(rootToken)
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	// Create the transit server.
	tcluster := sealhelper.NewTransitSealServer(t)
	defer tcluster.Cleanup()
	tcluster.MakeKey(t, "key1")
	transitSeal := tcluster.MakeSeal(t, "key1")

	// Transition to transit seal.
	if err := adjustCoreForSealMigration(cluster.Logger, cluster.Cores[0].Core, transitSeal, nil); err != nil {
		t.Fatal(err)
	}

	// Unseal and migrate to transit.
	for _, key := range initResp.KeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err == nil {
			t.Fatal("expected error due to lack of migrate parameter")
		}
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key, Migrate: true})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if resp == nil || resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", resp)
	}

	testhelpers.WaitForActiveNode(t, cluster)

	// Seal the cluster.
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	// Seal the transit cluster; we expect the unseal of our main cluster
	// to fail as a result.
	tcluster.EnsureCoresSealed(t)

	// Verify that we cannot unseal.  Now the barrier unseal keys are actually
	// the recovery keys.
	for _, key := range initResp.KeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err != nil {
			break
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if err == nil || resp != nil {
		t.Fatalf("expected sealed state; got %#v", resp)
	}

	// Unseal the transit server; we expect the unseal to work now on our main
	// cluster.
	tcluster.UnsealCores(t)

	// Verify that we can unseal.
	for _, key := range initResp.KeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if resp == nil || resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", resp)
	}

	// Make sure the seal configs were updated correctly.
	b, r, err := cluster.Cores[0].Core.PhysicalSealConfigs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	verifyBarrierConfig(t, b, wrapping.Transit, 1, 1, 1)
	verifyBarrierConfig(t, r, wrapping.Shamir, 5, 3, 0)
}

func TestSealMigration_ShamirToTestSeal(t *testing.T) {
	t.Parallel()
	t.Run("inmem", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTestSeal(t, teststorage.InmemBackendSetup)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTestSeal(t, teststorage.FileBackendSetup)
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTestSeal(t, teststorage.ConsulBackendSetup)
	})

	t.Run("raft", func(t *testing.T) {
		t.Parallel()
		testSealMigrationShamirToTestSeal(t, teststorage.RaftBackendSetup)
	})
}

func testSealMigrationShamirToTestSeal(t *testing.T, setup teststorage.ClusterSetupMutator) {

	// Create a cluster that uses shamir.
	conf, opts := teststorage.ClusterSetup(&vault.CoreConfig{
		DisableSealWrap: true,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    3,
	},
		setup,
	)
	opts.SetupFunc = nil
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	// Initialize the cluster, and unseal it using the shamir keys.
	client := cluster.Cores[0].Client
	initResp, err := client.Sys().Init(&api.InitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	var resp *api.SealStatusResponse
	for _, key := range initResp.KeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if resp == nil || resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", resp)
	}

	testhelpers.WaitForActiveNode(t, cluster)

	rootToken := initResp.RootToken
	client.SetToken(rootToken)
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	// Create the test seal
	testSeal := vault.NewAutoSeal(vaultseal.NewTestSeal(&vaultseal.TestSealOpts{}))

	// Transition to test seal.
	if err := adjustCoreForSealMigration(cluster.Logger, cluster.Cores[0].Core, testSeal, nil); err != nil {
		t.Fatal(err)
	}

	// Unseal and migrate to test seal.
	for _, key := range initResp.KeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err == nil {
			t.Fatal("expected error due to lack of migrate parameter")
		}
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key, Migrate: true})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if resp == nil || resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", resp)
	}
	testhelpers.WaitForActiveNode(t, cluster)

	// Seal the cluster.
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	////// Seal the transit cluster; we expect the unseal of our main cluster
	////// to fail as a result.
	////tcluster.EnsureCoresSealed(t)

	////// Verify that we cannot unseal.  Now the barrier unseal keys are actually
	////// the recovery keys.
	////for _, key := range initResp.KeysB64 {
	////	resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
	////	if err != nil {
	////		break
	////	}
	////	if resp == nil || !resp.Sealed {
	////		break
	////	}
	////}
	////if err == nil || resp != nil {
	////	t.Fatalf("expected sealed state; got %#v", resp)
	////}

	////// Unseal the transit server; we expect the unseal to work now on our main
	////// cluster.
	////tcluster.UnsealCores(t)

	// Verify that we can unseal.
	for _, key := range initResp.KeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if resp == nil || resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", resp)
	}

	// Make sure the seal configs were updated correctly.
	b, r, err := cluster.Cores[0].Core.PhysicalSealConfigs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	verifyBarrierConfig(t, b, wrapping.Test, 1, 1, 1)
	verifyBarrierConfig(t, r, wrapping.Shamir, 5, 3, 0)
}

/*
func TestSealMigration_TransitToTransit(t *testing.T) {
	t.Parallel()
	t.Run("inmem", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToTransit(t, teststorage.InmemBackendSetup)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToTransit(t, teststorage.FileBackendSetup)
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToTransit(t, teststorage.ConsulBackendSetup)
	})

	t.Run("raft", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToTransit(t, teststorage.RaftBackendSetup)
	})
}

func testSealMigrationTransitToTransit(t *testing.T, setup teststorage.ClusterSetupMutator) {
	tcluster := sealhelper.NewTransitSealServer(t)
	defer tcluster.Cleanup()
	tcluster.MakeKey(t, "key1")
	tcluster.MakeKey(t, "key2")
	var seals []vault.Seal

	conf, opts := teststorage.ClusterSetup(&vault.CoreConfig{
		DisableSealWrap: true,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    3,
		SealFunc: func() vault.Seal {
			tseal := tcluster.MakeSeal(t, "key1")
			seals = append(seals, tseal)
			return tseal
		},
	},
		setup,
	)
	opts.SetupFunc = nil
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	initResp, err := client.Sys().Init(&api.InitRequest{
		RecoveryShares:    5,
		RecoveryThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	rootToken := initResp.RootToken
	client.SetToken(rootToken)
	for _, k := range initResp.RecoveryKeysB64 {
		b, _ := base64.RawStdEncoding.DecodeString(k)
		cluster.RecoveryKeys = append(cluster.RecoveryKeys, b)
	}

	testhelpers.WaitForActiveNode(t, cluster)

	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	logger := cluster.Logger.Named("shamir")
	autoSeal2 := tcluster.MakeSeal(t, "key2")
	if err := adjustCoreForSealMigration(logger, cluster.Cores[0].Core, autoSeal2, seals[0]); err != nil {
		t.Fatal(err)
	}

	// Although we're unsealing using the recovery keys, this is still an
	// autounseal; if we stopped the transit cluster this would fail.
	var resp *api.SealStatusResponse
	for _, key := range initResp.RecoveryKeysB64 {
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key})
		if err == nil {
			t.Fatal("expected error due to lack of migrate parameter")
		}
		resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: key, Migrate: true})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.Sealed {
			break
		}
	}
	if resp == nil || resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", resp)
	}

	testhelpers.WaitForActiveNode(t, cluster)

	// Seal and unseal again to verify that things are working fine
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	// Delete the original seal's transit key.
	_, err = tcluster.Cores[0].Client.Logical().Delete(path.Join("transit", "keys", "key1"))
	if err != nil {
		t.Fatal(err)
	}

	err = cluster.Cores[0].Core.UnsealWithStoredKeys(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
*/
