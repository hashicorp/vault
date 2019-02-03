package storagepacker

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/shamir"
	"github.com/hashicorp/vault/vault"
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

	// Step 1: write something into Identity so that we create storage paths
	// and know where to put things
	_, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "foobar",
	})

	// Step 2: seal, so we can modify data without Vault
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	// Step 3: Unseal the barrier so we can write legit stuff into the data
	// store
	barrierKey, err := shamir.Combine(cluster.BarrierKeys[0:3])
	if err != nil {
		t.Fatal(err)
	}
	if barrierKey == nil {
		t.Fatal("nil barrier key")
	}

	if core.UnderlyingStorage == nil {
		t.Fatal("underlying storage is nil")
	}

	barrier, err := vault.NewAESGCMBarrier(core.UnderlyingStorage)
	if err != nil {
		t.Fatal(err)
	}
	if barrier == nil {
		t.Fatal("nil barrier")
	}

	if err := barrier.Unseal(ctx, barrierKey); err != nil {
		t.Fatal(err)
	}

	// Step 4: Remove exisitng packer data, create a legacy packer, write
	// stuff, ensure that all buckets are created.
	bes, err := barrier.List(ctx, "logical/")
	if err != nil {
		t.Fatal(err)
	}

	if len(bes) > 1 {
		t.Fatalf("expected only identity logical area, got %v", bes)
	}

	entityPackerLogger := logger.Named("storagepacker").Named("entities")
	groupPackerLogger := logger.Named("storagepacker").Named("groups")
	storage := logical.NewStorageView(barrier, "logical/"+bes[0])

	numEntries := 10000

	if err := logical.ClearView(ctx, storage); err != nil {
		t.Fatal(err)
	}

	entityPacker, err := NewLegacyStoragePacker(ctx, &storagepacker.Config{
		BucketStorageView: storage.SubView("packer/buckets/"),
		Logger:            entityPackerLogger,
	})
	if err != nil {
		t.Fatal(err)
	}

	groupPacker, err := NewLegacyStoragePacker(ctx, &storagepacker.Config{
		BucketStorageView: storage.SubView("packer/group/buckets/"),
		Logger:            groupPackerLogger,
	})
	if err != nil {
		t.Fatal(err)
	}

	var entity identity.Entity
	var group identity.Group
	var item storagepacker.Item
	for i := 0; i < numEntries; i++ {
		entity.ID, _ = uuid.GenerateUUID()
		entity.Name = fmt.Sprintf("%d", i)
		entityAsAny, err := ptypes.MarshalAny(&entity)
		if err != nil {
			t.Fatal(err)
		}
		item.ID = entity.ID
		item.Message = entityAsAny
		if err := entityPacker.PutItem(ctx, &item); err != nil {
			t.Fatal(err)
		}

		group.ID, _ = uuid.GenerateUUID()
		group.Name = fmt.Sprintf("%d", i)
		groupAsAny, err := ptypes.MarshalAny(&group)
		if err != nil {
			t.Fatal(err)
		}
		item.ID = group.ID
		item.Message = groupAsAny
		if err := groupPacker.PutItem(ctx, &item); err != nil {
			t.Fatal(err)
		}
	}

	buckets, err := barrier.List(ctx, "logical/"+bes[0]+"packer/buckets/")
	if err != nil {
		t.Fatal(err)
	}
	if len(buckets) != 256 {
		t.Fatalf("%d", len(buckets))
	}

	buckets, err = barrier.List(ctx, "logical/"+bes[0]+"packer/group/buckets/")
	if err != nil {
		t.Fatal(err)
	}
	if len(buckets) != 256 {
		t.Fatalf("%d", len(buckets))
	}

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

		buckets, err = barrier.List(ctx, "logical/"+bes[0]+"packer/buckets/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 1 {
			t.Fatalf("%d", len(buckets))
		}

		buckets, err = barrier.List(ctx, "logical/"+bes[0]+"packer/buckets/v2/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 256 {
			t.Fatalf("%d", len(buckets))
		}

		buckets, err = barrier.List(ctx, "logical/"+bes[0]+"packer/group/buckets/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 1 {
			t.Fatalf("%d", len(buckets))
		}

		buckets, err = barrier.List(ctx, "logical/"+bes[0]+"packer/group/buckets/v2/")
		if err != nil {
			t.Fatal(err)
		}
		if len(buckets) != 256 {
			t.Fatalf("%d", len(buckets))
		}
	}
	step5()

	// Step 6: seal and unseal to make sure we're not just reading cache; IOW repeat step 5
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	step5()

}
