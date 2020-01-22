// +build !enterprise

package command

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSealMigrationAutoToShamir(t *testing.T) {
	tcluster := newTransitSealServer(t)
	defer tcluster.Cleanup()

	logger := logging.NewVaultLogger(hclog.Trace).Named(t.Name())
	phys, err := physInmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	haPhys, err := physInmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	autoSeal := tcluster.makeKeyAndSeal(t, "key1")
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		Seal:            autoSeal,
		Physical:        phys,
		HAPhysical:      haPhys.(physical.HABackend),
		DisableSealWrap: true,
	}, &vault.TestClusterOptions{
		Logger:      logger,
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    1,
	})
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

	testhelpers.WaitForActiveNode(t, cluster)

	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	shamirSeal := vault.NewDefaultSeal(&seal.Access{
		Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
			Logger: logger.Named("shamir"),
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
