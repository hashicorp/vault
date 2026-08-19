// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build race

package vault

import (
	"context"
	"sync"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
)

// TestIdentityStoreInvalidate_LocalAliasesDataRace is a regression test for a
// data race (VAULT-45771) in invalidateLocalAliasesBucket. The invalidation
// used to read (via reflect.DeepEqual) and mutate (via UpsertAlias) the shared
// *identity.Entity and *identity.Alias objects stored in MemDB directly. Other
// goroutines - for example login or token validation on a performance standby -
// read or clone those same objects from MemDB using lock-free read
// transactions without holding i.lock, which races with the invalidation.
//
// This test is meaningful only under the race detector, so it lives in a file
// gated behind the "race" build tag (which the Go toolchain sets automatically
// when -race is passed) and is excluded from ordinary builds.
//
// This is deliberately an internal (white-box) test: it calls Invalidate
// directly and reads and clones the same MemDB entity and alias via
// MemDBEntityByID and MemDBAliasByID. That guarantees the racy condition on
// every run - two goroutines touching the same shared, unsynchronized objects.
// An external test cannot: driven through the public API, invalidations arrive
// asynchronously over the performance-replication WAL, so their timing cannot
// be lined up with the concurrent reads.
func TestIdentityStoreInvalidate_LocalAliasesDataRace(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	is := c.identityStore
	ctx := context.Background()

	entityID, err := uuid.GenerateUUID()
	require.NoError(t, err)

	aliasID, err := uuid.GenerateUUID()
	require.NoError(t, err)

	// Persist an entity in storage and invalidate so that MemDB picks it up.
	entity := &identity.Entity{
		Name:        "test",
		NamespaceID: namespace.RootNamespaceID,
		ID:          entityID,
		Aliases:     []*identity.Alias{},
		BucketKey:   is.entityPacker.BucketKey(entityID),
	}

	entityAsAny, err := anypb.New(entity)
	require.NoError(t, err)

	require.NoError(t, is.entityPacker.PutItem(ctx, &storagepacker.Item{
		ID:      entityID,
		Message: entityAsAny,
	}))
	is.Invalidate(ctx, is.entityPacker.BucketKey(entityID))

	// buildLocalAliasItem builds a storage item for the entity's single local
	// alias, varying only the alias Name. LocalBucketKey is always set to the
	// bucket key we invalidate below so the alias stays discoverable by
	// MemDBLocalAliasesByBucketKeyInTxn; that lookup is what makes
	// invalidateLocalAliasesBucket run its proto.Equal comparison against the
	// shared MemDB alias. (If LocalBucketKey were left unset the alias would
	// never be found, the comparison would be short-circuited, and the alias
	// would take the update path unconditionally.)
	localAliasBucketKey := is.localAliasPacker.BucketKey(entityID)

	buildLocalAliasItem := func(name string) *storagepacker.Item {
		localAliases := &identity.LocalAliases{
			Aliases: []*identity.Alias{
				{
					ID:             aliasID,
					Name:           name,
					NamespaceID:    namespace.RootNamespaceID,
					CanonicalID:    entityID,
					MountAccessor:  "userpass-000000",
					Local:          true,
					LocalBucketKey: localAliasBucketKey,
				},
			},
		}

		message, err := anypb.New(localAliases)
		require.NoError(t, err)

		return &storagepacker.Item{ID: entityID, Message: message}
	}

	// Prime MemDB with an initial version of the local alias and its association
	// to the entity.
	require.NoError(t, is.localAliasPacker.PutItem(ctx, buildLocalAliasItem("test-initial")))
	is.Invalidate(ctx, localAliasBucketKey)

	// Two variants of the alias that differ only by Name.
	aliasItems := []*storagepacker.Item{
		buildLocalAliasItem("test-a"),
		buildLocalAliasItem("test-b"),
	}

	// Arranging the racy condition does not by itself guarantee the race
	// detector reports it on a single run: whether the conflicting accesses
	// interleave (and whether the runtime happens to impose an ordering between
	// them) is up to the goroutine scheduler. Repeating the work many times
	// makes the detector reliably observe an unsynchronized access.
	const iterations = 200

	var wg sync.WaitGroup
	wg.Add(2)

	// Repeatedly swap the stored alias, then invalidate the bucket. Varying the
	// alias every iteration keeps invalidateLocalAliasesBucket on its "changed"
	// branch, so each pass both compares (proto.Equal) the shared MemDB alias
	// and updates its entity - the accesses that, before the fix, were an
	// unsafe reflect.DeepEqual read and an in-place UpsertAlias on the shared
	// objects the reader goroutine below is concurrently cloning.
	go func() {
		defer wg.Done()
		for n := 0; n < iterations; n++ {
			if !assert.NoError(t, is.localAliasPacker.PutItem(ctx, aliasItems[n%2])) {
				return
			}
			is.Invalidate(ctx, localAliasBucketKey)
		}
	}()

	// Concurrently read and clone the same entity and alias from MemDB without
	// holding i.lock, mimicking a login or token-validation path on a
	// performance standby. Cloning marshals the underlying protobuf messages,
	// which is what races with the invalidation's reads and mutations of the
	// shared objects.
	go func() {
		defer wg.Done()
		for n := 0; n < iterations; n++ {
			if e, err := is.MemDBEntityByID(entityID, true); err == nil && e != nil {
				_, _ = e.Clone()
			}
			if a, err := is.MemDBAliasByID(aliasID, true, false); err == nil && a != nil {
				_, _ = a.Clone()
			}
		}
	}()

	wg.Wait()

	// Sanity check that the entity and its local alias are still present and
	// correctly associated after all the concurrent invalidations.
	memDBEntity, err := is.MemDBEntityByID(entityID, false)
	require.NoError(t, err)
	require.NotNil(t, memDBEntity)
	require.Len(t, memDBEntity.Aliases, 1)
	require.Equal(t, aliasID, memDBEntity.Aliases[0].ID)

	memDBAlias, err := is.MemDBAliasByID(aliasID, false, false)
	require.NoError(t, err)
	require.NotNil(t, memDBAlias)
	require.Equal(t, entityID, memDBAlias.CanonicalID)
}
