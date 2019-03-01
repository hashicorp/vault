// +build !enterprise

package command

import (
	"context"
	"encoding/base64"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical"
	physInmem "github.com/hashicorp/vault/physical/inmem"
	"github.com/hashicorp/vault/shamir"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSealMigration(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	phys, err := physInmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	haPhys, err := physInmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	shamirSeal := vault.NewDefaultSeal()
	coreConfig := &vault.CoreConfig{
		Seal:            shamirSeal,
		Physical:        phys,
		HAPhysical:      haPhys.(physical.HABackend),
		DisableSealWrap: true,
	}
	clusterConfig := &vault.TestClusterOptions{
		Logger:      logger,
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    1,
	}

	ctx := context.Background()
	var keys []string
	var rootToken string

	// First: start up as normal with shamir seal, init it
	{
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		client := cluster.Cores[0].Client
		coreConfig = cluster.Cores[0].CoreConfig

		// Init
		resp, err := client.Sys().Init(&api.InitRequest{
			SecretShares:    2,
			SecretThreshold: 2,
		})
		if err != nil {
			t.Fatal(err)
		}
		keys = resp.KeysB64
		rootToken = resp.RootToken

		// Now seal
		cluster.Cleanup()
		// This will prevent cleanup from running again on the defer
		cluster.Cores = nil
	}

	// Second: start up as normal with shamir seal and unseal, make sure
	// everything is normal
	{
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		for _, key := range keys {
			resp, err = client.Sys().Unseal(key)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatal("expected unsealed state")
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	var autoSeal vault.Seal

	// Third: create an autoseal and activate migration
	{
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core

		newSeal := vault.NewAutoSeal(seal.NewTestSeal(nil))
		newSeal.SetCore(core)
		autoSeal = newSeal
		if err := adjustCoreForSealMigration(core, newSeal, nil); err != nil {
			t.Fatal(err)
		}

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		unsealOpts := &api.UnsealOpts{}
		for _, key := range keys {
			unsealOpts.Key = key
			unsealOpts.Migrate = false
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err == nil {
				t.Fatal("expected error due to lack of migrate parameter")
			}
			unsealOpts.Migrate = true
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	// Fourth: verify autoseal and recovery key usage
	{
		coreConfig.Seal = autoSeal
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core
		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		if err := core.UnsealWithStoredKeys(ctx); err != nil {
			t.Fatal(err)
		}
		resp, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		keyParts := [][]byte{}
		for _, key := range keys {
			raw, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				t.Fatal(err)
			}
			keyParts = append(keyParts, raw)
		}
		recoveredKey, err := shamir.Combine(keyParts)
		if err != nil {
			t.Fatal(err)
		}
		sealAccess := core.SealAccess()
		if err := sealAccess.VerifyRecoveryKey(ctx, recoveredKey); err != nil {
			t.Fatal(err)
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	// We should see stored barrier keys; after the sixth test, we shouldn't
	if entry, err := phys.Get(ctx, vault.StoredBarrierKeysPath); err != nil || entry == nil {
		t.Fatalf("expected nil error and non-nil entry, got error %#v and entry %#v", err, entry)
	}

	altTestSeal := seal.NewTestSeal(nil)
	altTestSeal.Type = "test-alternate"
	altSeal := vault.NewAutoSeal(altTestSeal)

	// Fifth: migrate from auto-seal to auto-seal
	{
		coreConfig.Seal = autoSeal
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core

		if err := adjustCoreForSealMigration(core, altSeal, autoSeal); err != nil {
			t.Fatal(err)
		}

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		unsealOpts := &api.UnsealOpts{}
		for _, key := range keys {
			unsealOpts.Key = key
			unsealOpts.Migrate = true
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	// Sixth: create an Shamir seal and activate migration. Verify it doesn't work
	// if disabled isn't set.
	{
		coreConfig.Seal = altSeal
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core

		if err := adjustCoreForSealMigration(core, shamirSeal, altSeal); err != nil {
			t.Fatal(err)
		}

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		unsealOpts := &api.UnsealOpts{}
		for _, key := range keys {
			unsealOpts.Key = key
			unsealOpts.Migrate = false
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err == nil {
				t.Fatal("expected error due to lack of migrate parameter")
			}
			unsealOpts.Migrate = true
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	if entry, err := phys.Get(ctx, vault.StoredBarrierKeysPath); err != nil || entry != nil {
		t.Fatalf("expected nil error and nil entry, got error %#v and entry %#v", err, entry)
	}

	// Seventh: verify autoseal is off and the expected key shares work
	{
		coreConfig.Seal = shamirSeal
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core
		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		if err := core.UnsealWithStoredKeys(ctx); err != nil {
			t.Fatal(err)
		}
		resp, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if !resp.Sealed {
			t.Fatalf("expected sealed state; got %#v", *resp)
		}

		for _, key := range keys {
			resp, err = client.Sys().Unseal(key)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatal("expected unsealed state")
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}
}
