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
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSealMigrationAutoToShamir(t *testing.T) {
	t.Parallel()
	t.Run("inmem", func(t *testing.T) {
		t.Parallel()
		testSealMigrationAutoToShamir(t, teststorage.InmemBackendSetup)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()
		testSealMigrationAutoToShamir(t, teststorage.FileBackendSetup)
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()
		testSealMigrationAutoToShamir(t, teststorage.ConsulBackendSetup)
	})

	t.Run("raft", func(t *testing.T) {
		t.Parallel()
		testSealMigrationAutoToShamir(t, teststorage.RaftBackendSetup)
	})
}

func testSealMigrationAutoToShamir(t *testing.T, setup teststorage.ClusterSetupMutator) {
	tcluster := newTransitSealServer(t)
	defer tcluster.Cleanup()

	autoSeal := tcluster.makeKeyAndSeal(t, "key1")
	conf, opts := teststorage.ClusterSetup(&vault.CoreConfig{
		Seal:            autoSeal,
		DisableSealWrap: true,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    1,
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

	logger := cluster.Logger.Named("shamir")
	shamirSeal := vault.NewDefaultSeal(&seal.Access{
		Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
			Logger: logger,
		}),
	})

	if err := adjustCoreForSealMigration(logger, cluster.Cores[0].Core, shamirSeal, autoSeal); err != nil {
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

	tcluster.Cleanup()
	// Assign nil to Cores so the deferred Cleanup doesn't break.
	tcluster.Cores = nil

	// Now the recovery keys are actually the barrier unseal keys.
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

	b, r, err := cluster.Cores[0].Core.PhysicalSealConfigs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	verifyBarrierConfig(t, b, wrapping.Shamir, 5, 3, 1)
	if r != nil {
		t.Fatalf("expected nil recovery config, got: %#v", r)
	}
}
