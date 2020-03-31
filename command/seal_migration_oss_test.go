// +build !enterprise

package command

import (
	"context"
	"encoding/base64"
	"testing"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSealMigration_TransitToShamir(t *testing.T) {
	t.Parallel()
	t.Run("inmem", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToShamir(t, teststorage.InmemBackendSetup)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToShamir(t, teststorage.FileBackendSetup)
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToShamir(t, teststorage.ConsulBackendSetup)
	})

	t.Run("raft", func(t *testing.T) {
		t.Parallel()
		testSealMigrationTransitToShamir(t, teststorage.RaftBackendSetup)
	})
}

func testSealMigrationTransitToShamir(t *testing.T, setup teststorage.ClusterSetupMutator) {

	// Create the transit server.
	tcluster := sealhelper.NewTransitSealServer(t)
	defer tcluster.Cleanup()
	tcluster.MakeKey(t, "key1")
	var transitSeal vault.Seal

	// Create a cluster that uses transit.
	conf, opts := teststorage.ClusterSetup(&vault.CoreConfig{
		DisableSealWrap: true,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    3,
		SealFunc: func() vault.Seal {
			transitSeal = tcluster.MakeSeal(t, "key1")
			return transitSeal
		},
	},
		setup,
	)
	opts.SetupFunc = nil
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	// Initialize the cluster, and fetch the recovery keys.
	client := cluster.Cores[0].Client
	initResp, err := client.Sys().Init(&api.InitRequest{
		RecoveryShares:    5,
		RecoveryThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range initResp.RecoveryKeysB64 {
		b, _ := base64.RawStdEncoding.DecodeString(k)
		cluster.RecoveryKeys = append(cluster.RecoveryKeys, b)
	}
	testhelpers.WaitForActiveNode(t, cluster)

	rootToken := initResp.RootToken
	client.SetToken(rootToken)
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	// Create a Shamir seal.
	logger := cluster.Logger.Named("shamir")
	shamirSeal := vault.NewDefaultSeal(&seal.Access{
		Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
			Logger: logger,
		}),
	})

	// Transition to Shamir seal.
	if err := adjustCoreForSealMigration(logger, cluster.Cores[0].Core, shamirSeal, transitSeal); err != nil {
		t.Fatal(err)
	}

	// Unseal Shamir.
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

	// Seal the cluster.
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	// Nuke the transit server; assign nil to Cores so the deferred Cleanup
	// doesn't break.
	tcluster.Cleanup()
	tcluster.Cores = nil

	// Unseal the cluster. Now the recovery keys are actually the barrier
	// unseal keys.
	for _, key := range initResp.RecoveryKeysB64 {
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
	verifyBarrierConfig(t, b, wrapping.Shamir, 5, 3, 1)
	if r != nil {
		t.Fatalf("expected nil recovery config, got: %#v", r)
	}
}
