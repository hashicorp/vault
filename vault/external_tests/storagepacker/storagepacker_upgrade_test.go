package storagepacker

import (
	"context"
	"fmt"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/kr/pretty"
)

func TestIdentityStore_StoragePacker_UpgradeFromLegacy(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
		Logger:      logger,
	})

	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client
	ctx := context.Background()
	numEntries := 10000

	storage := logical.NewLogicalStorage(core.UnderlyingStorage)

	// Step 1: Seal, so we can swap out the packer creation func
	cluster.EnsureCoresSealed(t)

	// Step 2: Start with a legacy packer
	vault.StoragePackerCreationFunc.Store(storagepacker.StoragePackerFactory(NewLegacyStoragePacker))

	// Step 3: Unseal with legacy, write stuff
	testhelpers.EnsureCoresUnsealed(t, cluster)
	vault.TestWaitActive(t, core.Core)

	for i := 0; i < numEntries; i++ {
		secret, err := client.Logical().Write("identity/entity", map[string]interface{}{
			"name": fmt.Sprintf("%d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			t.Fatal("nil secret")
		}
		if secret.Data["name"] != fmt.Sprintf("%d", i) {
			t.Fatalf("bad name, secret is %s", pretty.Sprint(secret))
		}

		secret, err = client.Logical().Write("identity/group", map[string]interface{}{
			"name": fmt.Sprintf("%d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			t.Fatal("nil secret")
		}
		if secret.Data["name"] != fmt.Sprintf("%d", i) {
			t.Fatal("bad name")
		}
	}

	// Step 4: Seal Vault again, check that the values we expect exist, swap to new storage packer
	cluster.EnsureCoresSealed(t)

	bes, err := storage.List(ctx, "logical/")
	if err != nil {
		t.Fatal(err)
	}

	if len(bes) > 1 {
		t.Fatalf("expected only identity logical area, got %v", bes)
	}

	buckets, err := storage.List(ctx, "logical/"+bes[0]+"packer/buckets/")
	if err != nil {
		t.Fatal(err)
	}
	if len(buckets) != 256 {
		t.Fatalf("%d", len(buckets))
	}
	t.Log(buckets)

	buckets, err = storage.List(ctx, "logical/"+bes[0]+"packer/group/buckets/")
	if err != nil {
		t.Fatal(err)
	}
	if len(buckets) != 256 {
		t.Fatalf("%d", len(buckets))
	}
	t.Log(buckets)

	vault.StoragePackerCreationFunc.Store(storagepacker.StoragePackerFactory(storagepacker.NewStoragePackerV2))

	// Step 5: Unseal Vault, make sure we can fetch every one of the created
	// identities, and that storage looks as we expect
	step5 := func() {
		testhelpers.EnsureCoresUnsealed(t, cluster)

		for i := 0; i < numEntries; i++ {
			secret, err := client.Logical().Read(fmt.Sprintf("identity/entity/name/%d", i))
			if err != nil {
				t.Fatal(err)
			}
			if secret == nil {
				t.Fatal("nil secret")
			}
			if secret.Data["name"] != fmt.Sprintf("%d", i) {
				t.Fatal("bad name")
			}

			secret, err = client.Logical().Read(fmt.Sprintf("identity/group/name/%d", i))
			if err != nil {
				t.Fatal(err)
			}
			if secret == nil {
				t.Fatal("nil secret")
			}
			if secret.Data["name"] != fmt.Sprintf("%d", i) {
				t.Fatal("bad name")
			}
		}

		buckets, err = storage.List(ctx, "logical/"+bes[0]+"packer/buckets/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 1 {
			t.Fatalf("%d", len(buckets))
		}

		t.Log(buckets)

		buckets, err = storage.List(ctx, "logical/"+bes[0]+"packer/buckets/v2/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 256 {
			t.Fatalf("%d", len(buckets))
		}

		t.Log(buckets)

		buckets, err = storage.List(ctx, "logical/"+bes[0]+"packer/group/buckets/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 1 {
			t.Fatalf("%d", len(buckets))
		}

		t.Log(buckets)

		buckets, err = storage.List(ctx, "logical/"+bes[0]+"packer/group/buckets/v2/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 256 {
			t.Fatalf("%d", len(buckets))
		}

		t.Log(buckets)
	}
	step5()

	// Step 6: seal and unseal to make sure we're not just reading cache; IOW repeat step 5
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	step5()

}
