// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"fmt"
	"math/rand"
	"testing"

	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

// TODO:
// Property 2:
// After merging, each alias should appear exactly once in toEntity and MemDB.
// There should be no duplicate aliases in either, and alias ordering must remain consistent with pre-merge constraints.
func TestUnit_EntityMerge(t *testing.T) {
	seedval := rand.Int63()
	seed := rand.New(rand.NewSource(seedval)) // Seed for deterministic test
	defer t.Logf("Test generated with seed: %d", seedval)

	logger := corehelpers.NewTestLogger(t)
	ims, err := inmem.NewTransactionalInmemHA(nil, logger)
	require.NoError(t, err)

	cfg := &CoreConfig{
		Physical:        ims,
		HAPhysical:      ims.(physical.HABackend),
		Logger:          logger,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}

	c, _, _ := TestCoreUnsealedWithConfig(t, cfg)

	upme, err := TestUserpassMount(c, false)
	require.NoError(t, err)

	ctx := context.Background()

	identityStore := c.identityStore

	name := fmt.Sprintf("entity-%d", 0)
	alias := fmt.Sprintf("alias-%d", 0)
	e1 := makeEntityForPacker(t, name, identityStore.entityPacker, seed)
	attachAlias(t, e1, alias, upme, seed)
	err = identityStore.persistEntity(ctx, e1)
	require.NoError(t, err)

	name2 := fmt.Sprintf("entity-%d", 1)
	alias2 := fmt.Sprintf("alias-%d", 1)
	e2 := makeEntityForPacker(t, name2, identityStore.entityPacker, seed)
	attachAlias(t, e2, alias, upme, seed)
	attachAlias(t, e2, alias2, upme, seed)
	err = identityStore.persistEntity(ctx, e2)
	require.NoError(t, err)

	txn := identityStore.db.Txn(true)
	defer txn.Abort()
	err = identityStore.MemDBUpsertEntityInTxn(txn, e1)
	require.NoError(t, err)
	err = identityStore.MemDBUpsertEntityInTxn(txn, e2)
	require.NoError(t, err)
	err = identityStore.MemDBUpsertAliasInTxn(txn, e1.Aliases[0], false)
	require.NoError(t, err)
	err = identityStore.MemDBUpsertAliasInTxn(txn, e2.Aliases[1], false)
	require.NoError(t, err)
	txn.Commit()

	ctx = namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	persist := true

	// Property 1:
	// After merging, ToEntity should retain all its original aliases
	// and all aliases from FromEntity, ensuring no alias is lost. Additionally, FromEntity should be deleted.
	// Merging toEntity and fromEntity should be idempotent

	txn = identityStore.db.Txn(true)
	defer txn.Abort()
	err1, err2 := identityStore.mergeEntityAsPartOfUpsert(ctx, txn, e1, e2.ID, persist)
	require.NoError(t, err1)
	require.NoError(t, err2)
	txn.Commit()

	entity1, err := identityStore.MemDBEntityByID(e1.ID, true)
	require.NoError(t, err)
	require.NotNil(t, entity1)
	require.Contains(t, entity1.MergedEntityIDs, e2.ID)

	entity2, err := identityStore.MemDBEntityByID(e2.ID, true)
	require.NoError(t, err)
	require.Nil(t, entity2)

	require.Len(t, entity1.Aliases, 2)

	for _, a := range entity1.Aliases {
		alias, err := identityStore.MemDBAliasByID(a.ID, true, false)
		require.NoError(t, err)
		require.NotNil(t, alias)
		require.Equal(t, entity1.ID, a.CanonicalID)
	}

	err = identityStore.resetDB()
	require.NoError(t, err)

	err = identityStore.loadEntities(ctx, true)
	require.NoError(t, err)

	entity1, err = identityStore.MemDBEntityByID(e1.ID, true)
	require.NoError(t, err)
	require.NotNil(t, entity1)
	require.Contains(t, entity1.MergedEntityIDs, e2.ID)

	entity2, err = identityStore.MemDBEntityByID(e2.ID, true)
	require.NoError(t, err)
	require.Nil(t, entity2)

	require.Len(t, entity1.Aliases, 2)

	for _, a := range entity1.Aliases {
		alias, err := identityStore.MemDBAliasByID(a.ID, true, false)
		require.NoError(t, err)
		require.NotNil(t, alias)
		require.Equal(t, entity1.ID, a.CanonicalID)
	}
}

func TestUnit_EntityMerge2(t *testing.T) {
	seedval := rand.Int63()
	seed := rand.New(rand.NewSource(seedval)) // Seed for deterministic test
	defer t.Logf("Test generated with seed: %d", seedval)

	logger := corehelpers.NewTestLogger(t)
	ims, err := inmem.NewTransactionalInmemHA(nil, logger)
	require.NoError(t, err)

	cfg := &CoreConfig{
		Physical:        ims,
		HAPhysical:      ims.(physical.HABackend),
		Logger:          logger,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}

	c, _, _ := TestCoreUnsealedWithConfig(t, cfg)

	upme, err := TestUserpassMount(c, false)
	require.NoError(t, err)

	ctx := context.Background()

	identityStore := c.identityStore

	name := fmt.Sprintf("entity-%d", 0)
	alias := fmt.Sprintf("alias-%d", 0)
	alias2 := fmt.Sprintf("alias-%d", 1)
	e1 := makeEntityForPacker(t, name, identityStore.entityPacker, seed)
	attachAlias(t, e1, alias, upme, seed)
	attachAlias(t, e1, alias2, upme, seed)
	err = identityStore.persistEntity(ctx, e1)
	require.NoError(t, err)

	name2 := fmt.Sprintf("entity-%d", 1)
	e2 := makeEntityForPacker(t, name2, identityStore.entityPacker, seed)
	attachAlias(t, e2, alias, upme, seed)
	attachAlias(t, e2, alias2, upme, seed)
	err = identityStore.persistEntity(ctx, e2)
	require.NoError(t, err)

	name3 := fmt.Sprintf("entity-%d", 2)
	e3 := makeEntityForPacker(t, name3, identityStore.entityPacker, seed)
	attachAlias(t, e3, alias, upme, seed)
	err = identityStore.persistEntity(ctx, e3)
	require.NoError(t, err)

	txn := identityStore.db.Txn(true)
	defer txn.Abort()
	err = identityStore.MemDBUpsertEntityInTxn(txn, e1)
	require.NoError(t, err)
	err = identityStore.MemDBUpsertEntityInTxn(txn, e2)
	require.NoError(t, err)
	err = identityStore.MemDBUpsertEntityInTxn(txn, e3)
	require.NoError(t, err)
	err = identityStore.MemDBUpsertAliasInTxn(txn, e1.Aliases[0], false)
	require.NoError(t, err)
	err = identityStore.MemDBUpsertAliasInTxn(txn, e2.Aliases[1], false)
	require.NoError(t, err)
	txn.Commit()

	ctx = namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	persist := true

	// Property 1:
	// After merging, ToEntity should retain all its original aliases
	// and all aliases from FromEntity, ensuring no alias is lost. Additionally, FromEntity should be deleted.
	// Merging toEntity and fromEntity should be idempotent

	txn = identityStore.db.Txn(true)
	defer txn.Abort()
	err1, err2 := identityStore.mergeEntityAsPartOfUpsert(ctx, txn, e1, e2.ID, persist)
	require.NoError(t, err1)
	require.NoError(t, err2)
	txn.Commit()

	txn = identityStore.db.Txn(true)
	defer txn.Abort()
	err1, err2 = identityStore.mergeEntityAsPartOfUpsert(ctx, txn, e1, e3.ID, persist)
	require.NoError(t, err1)
	require.NoError(t, err2)
	txn.Commit()

	entity1, err := identityStore.MemDBEntityByID(e1.ID, true)
	require.NoError(t, err)
	require.NotNil(t, entity1)
	require.Contains(t, entity1.MergedEntityIDs, e2.ID)

	entity2, err := identityStore.MemDBEntityByID(e2.ID, true)
	require.NoError(t, err)
	require.Nil(t, entity2)

	entity3, err := identityStore.MemDBEntityByID(e3.ID, true)
	require.NoError(t, err)
	require.Nil(t, entity3)

	require.Len(t, entity1.Aliases, 2)

	for _, a := range entity1.Aliases {
		alias, err := identityStore.MemDBAliasByID(a.ID, true, false)
		require.NoError(t, err)
		require.NotNil(t, alias)
		require.Equal(t, entity1.ID, a.CanonicalID)
	}
}
