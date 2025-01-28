// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/errwrap"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/activationflags"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/identity/mfa"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	errCycleDetectedPrefix = "cyclic relationship detected for member group ID"
	tmpSuffix              = ".tmp"
	entityLoadingTxMaxSize = 1024
)

// loadIdentityStoreArtifacts is responsible for loading entities, groups, and aliases from storage into MemDB.
func (c *Core) loadIdentityStoreArtifacts(ctx context.Context) error {
	if c.identityStore == nil {
		c.logger.Warn("identity store is not setup, skipping loading")
		return nil
	}

	loadFunc := func(context.Context) error {
		if err := c.identityStore.loadEntities(ctx); err != nil {
			return fmt.Errorf("failed to load entities: %w", err)
		}
		if err := c.identityStore.loadGroups(ctx); err != nil {
			return fmt.Errorf("failed to load groups: %w", err)
		}
		if err := c.identityStore.loadOIDCClients(ctx); err != nil {
			return fmt.Errorf("failed to load OIDC clients: %w", err)
		}
		if err := c.identityStore.loadCachedEntitiesOfLocalAliases(ctx); err != nil {
			return fmt.Errorf("failed to load cached local alias entities: %w", err)
		}

		return nil
	}

	// Resolve all conflicts by logging a warning and returning
	// errDuplicateIdentityName by default. The error will flip the
	// identityStore into case-sensitive mode by switching the underlying
	// schema to one with a relaxed lowerCase constraint and reload all
	// artifacts into MemDB.
	c.identityStore.conflictResolver = &errorResolver{c.identityStore.logger}

	// If the identity deduplication cleanup flag is activated, instead
	// deal with duplicate entities and groups by renaming with a -UUID
	// suffix. N.B. *entity alias* duplicates will still be merged as before.
	if c.FeatureActivationFlags.IsActivationFlagEnabled(activationflags.IdentityDeduplication) {
		c.identityStore.conflictResolver = &renameResolver{c.identityStore.logger}
	}

	// Load everything when MemDB is set to operate on lower cased names.
	// errDuplicateIdentityName below should only happen if we're using the
	// errorResolver (i.e. identity deduplication is not activated) and we
	// encounter non-alias duplicates.
	err := loadFunc(ctx)
	switch {
	case err == nil:
		// No error implies we've loaded the artifacts successfully
		// with no duplicates detected. This means there were no
		// unmerged duplicates detected while loading with the
		// errorResolver, or the renameResolver was activated and
		// resolved all duplicates. In either case, we can return
		// early, since there's nothing left to do.
		return nil
	case !errwrap.Contains(err, errDuplicateIdentityName.Error()):
		// All other errors are unexpected and should be returned.
		return err
	}

	// If we're here, it means we've encountered duplicates while loading
	// with the errorResolver.
	c.identityStore.logger.Warn("enabling case sensitive identity names")

	// Set identity store to operate on case sensitive identity names
	c.identityStore.disableLowerCasedNames = true

	// Swap out the MemDB instance and reload artifacts with the
	// new schema.
	if err := c.identityStore.resetDB(); err != nil {
		return err
	}

	// Also reset the conflict resolver so that we report potential duplicates to
	// be resolved before it's safe to return to case-insensitive mode.
	reporterResolver := newDuplicateReportingErrorResolver(c.identityStore.logger)
	c.identityStore.conflictResolver = reporterResolver

	// Attempt to load identity artifacts once more after memdb is reset to
	// accept case sensitive names
	err = loadFunc(ctx)

	// Log reported duplicates if any found whether or not we end up erroring.
	reporterResolver.LogReport(c.identityStore.logger)

	return err
}

func (i *IdentityStore) sanitizeName(name string) string {
	if i.disableLowerCasedNames {
		return name
	}
	return strings.ToLower(name)
}

func (i *IdentityStore) loadGroups(ctx context.Context) error {
	i.logger.Debug("identity loading groups")
	existing, err := i.groupPacker.View().List(ctx, groupBucketsPrefix)
	if err != nil {
		return fmt.Errorf("failed to scan for groups: %w", err)
	}
	i.logger.Debug("groups collected", "num_existing", len(existing))

	for _, key := range existing {
		bucket, err := i.groupPacker.GetBucket(ctx, groupBucketsPrefix+key)
		if err != nil {
			return err
		}

		if bucket == nil {
			continue
		}

		for _, item := range bucket.Items {
			group, err := i.parseGroupFromBucketItem(item)
			if err != nil {
				return err
			}
			if group == nil {
				continue
			}

			ns, err := i.namespacer.NamespaceByID(ctx, group.NamespaceID)
			if err != nil {
				return err
			}
			if ns == nil {
				// Remove dangling groups
				if !(i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) || i.localNode.HAState() == consts.PerfStandby) {
					// Group's namespace doesn't exist anymore but the group
					// from the namespace still exists.
					i.logger.Warn("deleting group and its any existing aliases", "name", group.Name, "namespace_id", group.NamespaceID)
					err = i.groupPacker.DeleteItem(ctx, group.ID)
					if err != nil {
						return err
					}
				}
				continue
			}
			nsCtx := namespace.ContextWithNamespace(ctx, ns)

			// Ensure that there are no groups with duplicate names
			groupByName, err := i.MemDBGroupByName(nsCtx, group.Name, false)
			if err != nil {
				return err
			}
			if err := i.conflictResolver.ResolveGroups(ctx, groupByName, group); err != nil && !i.disableLowerCasedNames {
				return err
			}

			if i.logger.IsDebug() {
				i.logger.Debug("loading group", "namespace", ns.ID, "name", group.Name, "id", group.ID)
			}

			txn := i.db.Txn(true)

			// Before pull#5786, entity memberships in groups were not getting
			// updated when respective entities were deleted. This is here to
			// check that the entity IDs in the group are indeed valid, and if
			// not remove them.
			persist := false
			for _, memberEntityID := range group.MemberEntityIDs {
				entity, err := i.MemDBEntityByID(memberEntityID, false)
				if err != nil {
					txn.Abort()
					return err
				}
				if entity == nil {
					persist = true
					group.MemberEntityIDs = strutil.StrListDelete(group.MemberEntityIDs, memberEntityID)
				}
			}

			err = i.UpsertGroupInTxn(ctx, txn, group, persist)

			if errors.Is(err, logical.ErrReadOnly) {
				// This is an imperfect solution to unblock customers who are running into
				// a readonly error during a DR failover (jira #28191). More specifically, if there
				// are duplicate aliases in storage then they are merged during loadEntities. Vault
				// attempts to remove the deleted duplicate entities from their groups to clean up.
				// If the node is a PR secondary though it will fail because the RPC client
				// is not yet initialized and the storage is read-only. This prevents the cluster from
				// unsealing entirely and can potentially block a DR failover from succeeding.
				i.logger.Warn("received a read only error while trying to upsert group to storage")
				err = nil
			}

			if err != nil {
				txn.Abort()
				return fmt.Errorf("failed to update group in memdb: %w", err)
			}

			txn.Commit()
		}
	}

	if i.logger.IsInfo() {
		i.logger.Info("groups restored")
	}

	return nil
}

func (i *IdentityStore) loadCachedEntitiesOfLocalAliases(ctx context.Context) error {
	// If we are performance secondary, load from temporary location those
	// entities that were created by the secondary via RPCs to the primary, and
	// also happen to have not yet been shipped to the secondary through
	// performance replication.
	if !i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) {
		return nil
	}

	i.logger.Debug("loading cached entities of local aliases")
	existing, err := i.localAliasPacker.View().List(ctx, localAliasesBucketsPrefix)
	if err != nil {
		return fmt.Errorf("failed to scan for cached entities of local alias: %w", err)
	}

	i.logger.Debug("cached entities of local alias entries", "num_buckets", len(existing))

	// Make the channels used for the worker pool
	broker := make(chan int)
	quit := make(chan bool)

	// We want to process the buckets in deterministic order so that duplicate
	// merging is deterministic. We still want to load in parallel though so
	// create a slice of result channels, one for each bucket. We need each result
	// and err chan to be 1 buffered so we can leave a result there even if the
	// processing loop is blocking on an earlier bucket still.
	results := make([]chan *storagepacker.Bucket, len(existing))
	errs := make([]chan error, len(existing))
	for j := range existing {
		results[j] = make(chan *storagepacker.Bucket, 1)
		errs[j] = make(chan error, 1)
	}

	// Use a wait group
	wg := &sync.WaitGroup{}

	// Create 64 workers to distribute work to
	for j := 0; j < consts.ExpirationRestoreWorkerCount; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case idx, ok := <-broker:
					// broker has been closed, we are done
					if !ok {
						return
					}
					key := existing[idx]

					bucket, err := i.localAliasPacker.GetBucket(ctx, localAliasesBucketsPrefix+key)
					if err != nil {
						errs[idx] <- err
						continue
					}

					// Write results out to the result channel
					results[idx] <- bucket

				// quit early
				case <-quit:
					return
				}
			}
		}()
	}

	// Distribute the collected keys to the workers in a go routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := range existing {
			if j%500 == 0 {
				i.logger.Debug("cached entities of local aliases loading", "progress", j)
			}

			select {
			case <-quit:
				return

			default:
				broker <- j
			}
		}

		// Close the broker, causing worker routines to exit
		close(broker)
	}()

	defer func() {
		// Let all go routines finish
		wg.Wait()

		i.logger.Info("cached entities of local aliases restored")
	}()

	// Restore each key by pulling from the slice of result chans
	for j := range existing {
		select {
		case err := <-errs[j]:
			// Close all go routines
			close(quit)

			return err

		case bucket := <-results[j]:
			// If there is no entry, nothing to restore
			if bucket == nil {
				continue
			}

			for _, item := range bucket.Items {
				if !strings.HasSuffix(item.ID, tmpSuffix) {
					continue
				}
				entity, err := i.parseCachedEntity(item)
				if err != nil {
					return err
				}
				ns, err := i.namespacer.NamespaceByID(ctx, entity.NamespaceID)
				if err != nil {
					return err
				}
				nsCtx := namespace.ContextWithNamespace(ctx, ns)
				err = i.upsertEntity(nsCtx, entity, nil, false)
				if err != nil {
					return fmt.Errorf("failed to update entity in MemDB: %w", err)
				}
			}
		}
	}

	return nil
}

func (i *IdentityStore) loadEntities(ctx context.Context) error {
	// Accumulate existing entities
	i.logger.Debug("loading entities")
	existing, err := i.entityPacker.View().List(ctx, storagepacker.StoragePackerBucketsPrefix)
	if err != nil {
		return fmt.Errorf("failed to scan for entities: %w", err)
	}
	i.logger.Debug("entities collected", "num_existing", len(existing))

	duplicatedAccessors := make(map[string]struct{})
	// Make the channels used for the worker pool. We send the index into existing
	// so that we can keep results in the same order as inputs. Note that this is
	// goroutine safe as long as we never mutate existing again in this method
	// which we don't.
	broker := make(chan int)
	quit := make(chan bool)

	// We want to process the buckets in deterministic order so that duplicate
	// merging is deterministic. We still want to load in parallel though so
	// create a slice of result channels, one for each bucket. We need each result
	// and err chan to be 1 buffered so we can leave a result there even if the
	// processing loop is blocking on an earlier bucket still.
	results := make([]chan []*identity.Entity, len(existing))
	errs := make([]chan error, len(existing))
	for j := range existing {
		results[j] = make(chan []*identity.Entity, 1)
		errs[j] = make(chan error, 1)
	}

	// Use a wait group
	wg := &sync.WaitGroup{}

	// Create 64 workers to distribute work to
	for j := 0; j < consts.ExpirationRestoreWorkerCount; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case idx, ok := <-broker:
					// broker has been closed, we are done
					if !ok {
						return
					}
					key := existing[idx]

					bucket, err := i.entityPacker.GetBucket(ctx, storagepacker.StoragePackerBucketsPrefix+key)
					if err != nil {
						errs[idx] <- err
						continue
					}

					items := make([]*identity.Entity, len(bucket.Items))
					for j, item := range bucket.Items {
						entity, err := i.parseEntityFromBucketItem(ctx, item)
						if err != nil {
							errs[idx] <- err
							continue
						}
						items[j] = entity
					}

					// Write results out to the result channel
					results[idx] <- items

				// quit early
				case <-quit:
					return
				}
			}
		}()
	}

	// Distribute the collected keys to the workers in a go routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := range existing {
			select {
			case <-quit:
				return

			default:
				broker <- j
			}
		}

		// Close the broker, causing worker routines to exit
		close(broker)
	}()

	localAliasBuckets := make(map[string]*storagepacker.Bucket)

	// Restore each key by pulling from the result chan
LOOP:
	for j := range existing {
		select {
		case err = <-errs[j]:
			// Close all go routines
			close(quit)
			break LOOP

		case entities := <-results[j]:
			// If there is no entry, nothing to restore
			if entities == nil {
				continue
			}
			load := func(entities []*identity.Entity) error {
				tx := i.db.Txn(true)
				defer tx.Abort()
				upsertedItems := 0
				for _, entity := range entities {
					if entity == nil {
						continue
					}

					ns, err := i.namespacer.NamespaceByID(ctx, entity.NamespaceID)
					if err != nil {
						return err
					}
					if ns == nil {
						// Remove dangling entities
						if !(i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) || i.localNode.HAState() == consts.PerfStandby) {
							// Entity's namespace doesn't exist anymore but the
							// entity from the namespace still exists.
							i.logger.Warn("deleting entity and its any existing aliases", "name", entity.Name, "namespace_id", entity.NamespaceID)
							err = i.entityPacker.DeleteItem(ctx, entity.ID)
							if err != nil {
								return err
							}
						}
						continue
					}
					nsCtx := namespace.ContextWithNamespace(ctx, ns)

					// Ensure that there are no entities with duplicate names
					entityByName, err := i.MemDBEntityByName(nsCtx, entity.Name, false)
					if err != nil {
						return nil
					}
					if err := i.conflictResolver.ResolveEntities(ctx, entityByName, entity); err != nil && !i.disableLowerCasedNames {
						return err
					}

					mountAccessors := getAccessorsOnDuplicateAliases(entity.Aliases)

					if len(mountAccessors) > 0 {
						i.logger.Warn("Entity has multiple aliases on the same mount(s)", "entity_id", entity.ID, "mount_accessors", mountAccessors)
					}

					for _, accessor := range mountAccessors {
						if _, ok := duplicatedAccessors[accessor]; !ok {
							duplicatedAccessors[accessor] = struct{}{}
						}
					}

					err = i.loadLocalAliasesForEntity(ctx, entity, localAliasBuckets)
					if err != nil {
						return fmt.Errorf("failed to load local aliases from storage: %v", err)
					}

					toBeUpserted := 1 + len(entity.Aliases)
					if upsertedItems+toBeUpserted > entityLoadingTxMaxSize {
						tx.Commit()
						upsertedItems = 0
						tx = i.db.Txn(true)
						defer tx.Abort()
					}
					// Only update MemDB and don't hit the storage again
					err = i.upsertEntityInTxn(nsCtx, tx, entity, nil, false)
					if err != nil {
						return fmt.Errorf("failed to update entity in MemDB: %w", err)
					}
					upsertedItems += toBeUpserted
				}
				if upsertedItems > 0 {
					tx.Commit()
				}
				return nil
			}
			err := load(entities)
			if err != nil {
				return err
			}

		}
	}

	// Let all go routines finish
	wg.Wait()
	if err != nil {
		return err
	}

	// Flatten the map into a list of keys, in order to log them
	duplicatedAccessorsList := make([]string, len(duplicatedAccessors))
	accessorCounter := 0
	for accessor := range duplicatedAccessors {
		duplicatedAccessorsList[accessorCounter] = accessor
		accessorCounter++
	}

	if i.logger.IsInfo() {
		i.logger.Info("entities restored")
	}

	return nil
}

// loadLocalAliasesForEntity upserts local aliases into the entity by retrieving
// the local aliases from the cache (if present) or storage
func (i *IdentityStore) loadLocalAliasesForEntity(ctx context.Context, entity *identity.Entity, localAliasCache map[string]*storagepacker.Bucket) error {
	bucketKey := i.localAliasPacker.BucketKey(entity.ID)
	if len(bucketKey) == 0 {
		return fmt.Errorf("no bucket key for ID %s", entity.ID)
	}
	bucket, ok := localAliasCache[bucketKey]
	if !ok {
		var err error
		bucket, err = i.localAliasPacker.GetBucket(ctx, bucketKey)
		if err != nil {
			return fmt.Errorf("failed to load local alias bucket from storage: %v", err)
		}
		localAliasCache[bucketKey] = bucket
	}
	if bucket == nil {
		return nil
	}
	for _, item := range bucket.Items {
		if item.ID == entity.ID {
			var localAliases identity.LocalAliases
			err := ptypes.UnmarshalAny(item.Message, &localAliases)
			if err != nil {
				return fmt.Errorf("failed to unmarshal local alias: %v", err)
			}
			for _, alias := range localAliases.Aliases {
				entity.UpsertAlias(alias)
			}
		}
	}
	return nil
}

// getAccessorsOnDuplicateAliases returns a list of accessors by checking aliases in
// the passed in list which belong to the same accessor(s)
func getAccessorsOnDuplicateAliases(aliases []*identity.Alias) []string {
	accessorCounts := make(map[string]int)
	var mountAccessors []string

	for _, alias := range aliases {
		accessorCounts[alias.MountAccessor] += 1
	}

	for accessor, accessorCount := range accessorCounts {
		if accessorCount > 1 {
			mountAccessors = append(mountAccessors, accessor)
		}
	}

	return mountAccessors
}

// upsertEntityInTxn either creates or updates an existing entity. The
// operations will be updated in both MemDB and storage. If 'persist' is set to
// false, then storage will not be updated. When an alias is transferred from
// one entity to another, both the source and destination entities should get
// updated, in which case, callers should send in both entity and
// previousEntity.
func (i *IdentityStore) upsertEntityInTxn(ctx context.Context, txn *memdb.Txn, entity *identity.Entity, previousEntity *identity.Entity, persist bool) error {
	defer metrics.MeasureSince([]string{"identity", "upsert_entity_txn"}, time.Now())
	var err error

	if txn == nil {
		return errors.New("txn is nil")
	}

	if entity == nil {
		return errors.New("entity is nil")
	}

	if entity.NamespaceID == "" {
		entity.NamespaceID = namespace.RootNamespaceID
	}

	if previousEntity != nil && previousEntity.NamespaceID != entity.NamespaceID {
		return errors.New("entity and previous entity are not in the same namespace")
	}

	aliasFactors := make([]string, len(entity.Aliases))

	for index, alias := range entity.Aliases {
		// Verify that alias is not associated to a different one already
		aliasByFactors, err := i.MemDBAliasByFactors(alias.MountAccessor, alias.Name, false, false)
		if err != nil {
			return err
		}

		if alias.NamespaceID == "" {
			alias.NamespaceID = namespace.RootNamespaceID
		}

		switch {
		case aliasByFactors == nil:
			// Not found, no merging needed, just check namespace
			if alias.NamespaceID != entity.NamespaceID {
				return errors.New("alias and entity are not in the same namespace")
			}

		case aliasByFactors.CanonicalID == entity.ID:
			// Lookup found the same entity, so it's already attached to the
			// right place
			if aliasByFactors.NamespaceID != entity.NamespaceID {
				return errors.New("alias from factors and entity are not in the same namespace")
			}

		case previousEntity != nil && aliasByFactors.CanonicalID == previousEntity.ID:
			// previousEntity isn't upserted yet so may still contain the old
			// alias reference in memdb if it was just changed; validate
			// whether or not it's _actually_ still tied to the entity
			var found bool
			for _, prevEntAlias := range previousEntity.Aliases {
				if prevEntAlias.ID == alias.ID {
					found = true
					break
				}
			}
			// If we didn't find the alias still tied to previousEntity, we
			// shouldn't use the merging logic and should bail
			if !found {
				break
			}

			// Otherwise it's still tied to previousEntity and fall through
			// into merging. We don't need a namespace check here as existing
			// checks when creating the aliases should ensure that all line up.
			fallthrough

		default:
			// Though this is technically a conflict that should be resolved by the
			// ConflictResolver implementation, the behavior here is a bit nuanced.
			// Rather than introduce a behavior change, we handle this case directly
			// as before by merging.
			i.logger.Warn("alias is already tied to a different entity; these entities are being merged",
				"alias_id", alias.ID,
				"other_entity_id", aliasByFactors.CanonicalID,
				"entity_aliases", entity.Aliases,
				"alias_by_factors", aliasByFactors)

			respErr, intErr := i.mergeEntityAsPartOfUpsert(ctx, txn, entity, aliasByFactors.CanonicalID, persist)
			switch {
			case respErr != nil:
				return respErr
			case intErr != nil:
				return intErr
			}

			// The entity and aliases will be loaded into memdb and persisted
			// as a result of the merge, so we are done here
			return nil
		}

		// This is subtle. We want to call `ResolveAliases` so that the resolver can
		// get full insight into all the aliases loaded and generate useful reports
		// about duplicates. However, we don't want to actually change the error
		// handling behavior from before which would only return an error in a very
		// specific case (when the alias being added is a duplicate of one for the
		// same entity and we are not in case-sensitive mode). So we call the method
		// here unconditionally, but then only handle the resultant error in the
		// specific case we care about. Note that we choose not to call it `err` to
		// avoid it being left non-nil in some cases and tripping up later error
		// handling code, and to signal something different is happening here. Note
		// that we explicitly _want_ this to be here an not before we merge
		// duplicates above, because duplicates that have always merged are not a
		// problem to the user and are already logged. We care about different-case
		// duplicates that are not being considered duplicates right now because we
		// are in case-sensitive mode so we can report these to the operator ahead
		// of them disabling case-sensitive mode.
		conflictErr := i.conflictResolver.ResolveAliases(ctx, entity, aliasByFactors, alias)

		// This appears to be accounting for any duplicate aliases for the same
		// Entity. In that case we would have skipped over the merge above in the
		// `aliasByFactors.CanonicalID == entity.ID` case and made it here. Now we
		// are here, duplicates are reported and may cause an insert error but only
		// if we are in default case-insensitive mode. Once we are in case-sensitive
		// mode we'll happily ignore duplicates of any case! This doesn't seem
		// especially desirable to me, but we'd rather not change behavior for now.
		if strutil.StrListContains(aliasFactors, i.sanitizeName(alias.Name)+alias.MountAccessor) &&
			conflictErr != nil && !i.disableLowerCasedNames {
			return conflictErr
		}

		// Insert or update alias in MemDB using the transaction created above
		err = i.MemDBUpsertAliasInTxn(txn, alias, false)
		if err != nil {
			return err
		}

		aliasFactors[index] = i.sanitizeName(alias.Name) + alias.MountAccessor
	}

	// If previous entity is set, update it in MemDB and persist it
	if previousEntity != nil {
		err = i.MemDBUpsertEntityInTxn(txn, previousEntity)
		if err != nil {
			return err
		}

		if persist {
			// Persist the previous entity object
			if err := i.persistEntity(ctx, previousEntity); err != nil {
				return err
			}
		}
	}

	// Insert or update entity in MemDB using the transaction created above
	err = i.MemDBUpsertEntityInTxn(txn, entity)
	if err != nil {
		return err
	}

	if persist {
		if err := i.persistEntity(ctx, entity); err != nil {
			return err
		}
	}

	return nil
}

func (i *IdentityStore) processLocalAlias(ctx context.Context, lAlias *logical.Alias, entity *identity.Entity, updateDb bool) (*identity.Alias, error) {
	if !lAlias.Local {
		return nil, fmt.Errorf("alias is not local")
	}

	mountValidationResp := i.router.ValidateMountByAccessor(lAlias.MountAccessor)
	if mountValidationResp == nil {
		return nil, fmt.Errorf("invalid mount accessor %q", lAlias.MountAccessor)
	}

	if !mountValidationResp.MountLocal {
		return nil, fmt.Errorf("mount accessor %q is not local", lAlias.MountAccessor)
	}

	alias, err := i.MemDBAliasByFactors(lAlias.MountAccessor, lAlias.Name, true, false)
	if err != nil {
		return nil, err
	}

	if alias == nil {
		alias = &identity.Alias{}
	}

	alias.CanonicalID = entity.ID
	alias.Name = lAlias.Name
	alias.MountAccessor = lAlias.MountAccessor
	alias.Metadata = lAlias.Metadata
	alias.MountPath = mountValidationResp.MountPath
	alias.MountType = mountValidationResp.MountType
	alias.Local = lAlias.Local
	alias.CustomMetadata = lAlias.CustomMetadata

	if err := i.sanitizeAlias(ctx, alias); err != nil {
		return nil, err
	}

	entity.UpsertAlias(alias)

	localAliases, err := i.parseLocalAliases(entity.ID)
	if err != nil {
		return nil, err
	}
	if localAliases == nil {
		localAliases = &identity.LocalAliases{}
	}

	updated := false
	for i, item := range localAliases.Aliases {
		if item.ID == alias.ID {
			localAliases.Aliases[i] = alias
			updated = true
			break
		}
	}

	if !updated {
		localAliases.Aliases = append(localAliases.Aliases, alias)
	}

	marshaledAliases, err := anypb.New(localAliases)
	if err != nil {
		return nil, err
	}
	if err := i.localAliasPacker.PutItem(ctx, &storagepacker.Item{
		ID:      entity.ID,
		Message: marshaledAliases,
	}); err != nil {
		return nil, err
	}

	if updateDb {
		txn := i.db.Txn(true)
		defer txn.Abort()
		if err := i.MemDBUpsertAliasInTxn(txn, alias, false); err != nil {
			return nil, err
		}
		if err := i.upsertEntityInTxn(ctx, txn, entity, nil, false); err != nil {
			return nil, err
		}
		txn.Commit()
	}

	return alias, nil
}

// cacheTemporaryEntity stores in secondary's storage, the entity returned by
// the primary cluster via the CreateEntity RPC. This is so that the secondary
// cluster knows and retains information about the existence of these entities
// before the replication invalidation informs the secondary of the same. This
// also happens to cover the case where the secondary's replication is lagging
// behind the primary by hours and/or days which sometimes may happen. Even if
// the nodes of the secondary are restarted in the interim, the cluster would
// still be aware of the entities. This temporary cache will be cleared when the
// invalidation hits the secondary nodes.
func (i *IdentityStore) cacheTemporaryEntity(ctx context.Context, entity *identity.Entity) error {
	if i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) && i.localNode.HAState() == consts.Active {
		marshaledEntity, err := anypb.New(entity)
		if err != nil {
			return err
		}
		if err := i.localAliasPacker.PutItem(ctx, &storagepacker.Item{
			ID:      entity.ID + tmpSuffix,
			Message: marshaledEntity,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (i *IdentityStore) persistEntity(ctx context.Context, entity *identity.Entity) error {
	// If the entity that is passed into this function is resulting from a memdb
	// query without cloning, then modifying it will result in a direct DB edit,
	// bypassing the transaction. To avoid any surprises arising from this
	// effect, work on a replica of the entity struct.
	var err error
	entity, err = entity.Clone()
	if err != nil {
		return err
	}

	// Separate the local and non-local aliases.
	var localAliases []*identity.Alias
	var nonLocalAliases []*identity.Alias
	for _, alias := range entity.Aliases {
		switch alias.Local {
		case true:
			localAliases = append(localAliases, alias)
		default:
			nonLocalAliases = append(nonLocalAliases, alias)
		}
	}

	// Store the entity with non-local aliases.
	entity.Aliases = nonLocalAliases
	marshaledEntity, err := anypb.New(entity)
	if err != nil {
		return err
	}
	if err := i.entityPacker.PutItem(ctx, &storagepacker.Item{
		ID:      entity.ID,
		Message: marshaledEntity,
	}); err != nil {
		return err
	}

	if len(localAliases) == 0 {
		return nil
	}

	// Store the local aliases separately.
	aliases := &identity.LocalAliases{
		Aliases: localAliases,
	}

	marshaledAliases, err := anypb.New(aliases)
	if err != nil {
		return err
	}
	if err := i.localAliasPacker.PutItem(ctx, &storagepacker.Item{
		ID:      entity.ID,
		Message: marshaledAliases,
	}); err != nil {
		return err
	}

	return nil
}

// upsertEntity either creates or updates an existing entity. The operations
// will be updated in both MemDB and storage. If 'persist' is set to false,
// then storage will not be updated. When an alias is transferred from one
// entity to another, both the source and destination entities should get
// updated, in which case, callers should send in both entity and
// previousEntity.
func (i *IdentityStore) upsertEntity(ctx context.Context, entity *identity.Entity, previousEntity *identity.Entity, persist bool) error {
	defer metrics.MeasureSince([]string{"identity", "upsert_entity"}, time.Now())

	// Create a MemDB transaction to update both alias and entity
	txn := i.db.Txn(true)
	defer txn.Abort()

	err := i.upsertEntityInTxn(ctx, txn, entity, previousEntity, persist)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (i *IdentityStore) MemDBUpsertAliasInTxn(txn *memdb.Txn, alias *identity.Alias, groupAlias bool) error {
	if txn == nil {
		return fmt.Errorf("nil txn")
	}

	if alias == nil {
		return fmt.Errorf("alias is nil")
	}

	if alias.NamespaceID == "" {
		alias.NamespaceID = namespace.RootNamespaceID
	}

	tableName := entityAliasesTable
	if groupAlias {
		tableName = groupAliasesTable
	}

	aliasRaw, err := txn.First(tableName, "id", alias.ID)
	if err != nil {
		return fmt.Errorf("failed to lookup alias from memdb using alias ID: %w", err)
	}

	if aliasRaw != nil {
		err = txn.Delete(tableName, aliasRaw)
		if err != nil {
			return fmt.Errorf("failed to delete alias from memdb: %w", err)
		}
	}

	if err := txn.Insert(tableName, alias); err != nil {
		return fmt.Errorf("failed to update alias into memdb: %w", err)
	}

	return nil
}

func (i *IdentityStore) MemDBAliasByIDInTxn(txn *memdb.Txn, aliasID string, clone bool, groupAlias bool) (*identity.Alias, error) {
	if aliasID == "" {
		return nil, fmt.Errorf("missing alias ID")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	tableName := entityAliasesTable
	if groupAlias {
		tableName = groupAliasesTable
	}

	aliasRaw, err := txn.First(tableName, "id", aliasID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch alias from memdb using alias ID: %w", err)
	}

	if aliasRaw == nil {
		return nil, nil
	}

	alias, ok := aliasRaw.(*identity.Alias)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched alias")
	}

	if clone {
		return alias.Clone()
	}

	return alias, nil
}

func (i *IdentityStore) MemDBAliasByID(aliasID string, clone bool, groupAlias bool) (*identity.Alias, error) {
	if aliasID == "" {
		return nil, fmt.Errorf("missing alias ID")
	}

	txn := i.db.Txn(false)

	return i.MemDBAliasByIDInTxn(txn, aliasID, clone, groupAlias)
}

func (i *IdentityStore) MemDBAliasByFactors(mountAccessor, aliasName string, clone bool, groupAlias bool) (*identity.Alias, error) {
	if aliasName == "" {
		return nil, fmt.Errorf("missing alias name")
	}

	if mountAccessor == "" {
		return nil, fmt.Errorf("missing mount accessor")
	}

	txn := i.db.Txn(false)

	return i.MemDBAliasByFactorsInTxn(txn, mountAccessor, aliasName, clone, groupAlias)
}

func (i *IdentityStore) MemDBAliasByFactorsInTxn(txn *memdb.Txn, mountAccessor, aliasName string, clone bool, groupAlias bool) (*identity.Alias, error) {
	if txn == nil {
		return nil, fmt.Errorf("nil txn")
	}

	if aliasName == "" {
		return nil, fmt.Errorf("missing alias name")
	}

	if mountAccessor == "" {
		return nil, fmt.Errorf("missing mount accessor")
	}

	tableName := entityAliasesTable
	if groupAlias {
		tableName = groupAliasesTable
	}

	aliasRaw, err := txn.First(tableName, "factors", mountAccessor, aliasName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch alias from memdb using factors: %w", err)
	}

	if aliasRaw == nil {
		return nil, nil
	}

	alias, ok := aliasRaw.(*identity.Alias)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched alias")
	}

	if clone {
		return alias.Clone()
	}

	return alias, nil
}

func (i *IdentityStore) MemDBDeleteAliasByIDInTxn(txn *memdb.Txn, aliasID string, groupAlias bool) error {
	if aliasID == "" {
		return nil
	}

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	alias, err := i.MemDBAliasByIDInTxn(txn, aliasID, false, groupAlias)
	if err != nil {
		return err
	}

	if alias == nil {
		return nil
	}

	tableName := entityAliasesTable
	if groupAlias {
		tableName = groupAliasesTable
	}

	err = txn.Delete(tableName, alias)
	if err != nil {
		return fmt.Errorf("failed to delete alias from memdb: %w", err)
	}

	return nil
}

func (i *IdentityStore) MemDBAliases(ws memdb.WatchSet, groupAlias bool) (memdb.ResultIterator, error) {
	txn := i.db.Txn(false)

	tableName := entityAliasesTable
	if groupAlias {
		tableName = groupAliasesTable
	}

	iter, err := txn.Get(tableName, "id")
	if err != nil {
		return nil, err
	}

	ws.Add(iter.WatchCh())

	return iter, nil
}

func (i *IdentityStore) MemDBUpsertEntityInTxn(txn *memdb.Txn, entity *identity.Entity) error {
	if txn == nil {
		return fmt.Errorf("nil txn")
	}

	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	if entity.NamespaceID == "" {
		entity.NamespaceID = namespace.RootNamespaceID
	}

	entityRaw, err := txn.First(entitiesTable, "id", entity.ID)
	if err != nil {
		return fmt.Errorf("failed to lookup entity from memdb using entity id: %w", err)
	}

	if entityRaw != nil {
		err = txn.Delete(entitiesTable, entityRaw)
		if err != nil {
			return fmt.Errorf("failed to delete entity from memdb: %w", err)
		}
	}

	if err := txn.Insert(entitiesTable, entity); err != nil {
		return fmt.Errorf("failed to update entity into memdb: %w", err)
	}

	return nil
}

func (i *IdentityStore) MemDBEntityByIDInTxn(txn *memdb.Txn, entityID string, clone bool) (*identity.Entity, error) {
	if entityID == "" {
		return nil, fmt.Errorf("missing entity id")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	entityRaw, err := txn.First(entitiesTable, "id", entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entity from memdb using entity id: %w", err)
	}

	if entityRaw == nil {
		return nil, nil
	}

	entity, ok := entityRaw.(*identity.Entity)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched entity")
	}

	if clone {
		return entity.Clone()
	}

	return entity, nil
}

func (i *IdentityStore) MemDBEntityByID(entityID string, clone bool) (*identity.Entity, error) {
	if entityID == "" {
		return nil, fmt.Errorf("missing entity id")
	}

	txn := i.db.Txn(false)

	return i.MemDBEntityByIDInTxn(txn, entityID, clone)
}

func (i *IdentityStore) MemDBEntityByName(ctx context.Context, entityName string, clone bool) (*identity.Entity, error) {
	if entityName == "" {
		return nil, fmt.Errorf("missing entity name")
	}

	txn := i.db.Txn(false)

	return i.MemDBEntityByNameInTxn(ctx, txn, entityName, clone)
}

func (i *IdentityStore) MemDBEntityByNameInTxn(ctx context.Context, txn *memdb.Txn, entityName string, clone bool) (*identity.Entity, error) {
	if entityName == "" {
		return nil, fmt.Errorf("missing entity name")
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	entityRaw, err := txn.First(entitiesTable, "name", ns.ID, entityName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entity from memdb using entity name: %w", err)
	}

	if entityRaw == nil {
		return nil, nil
	}

	entity, ok := entityRaw.(*identity.Entity)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched entity")
	}

	if clone {
		return entity.Clone()
	}

	return entity, nil
}

func (i *IdentityStore) MemDBLocalAliasesByBucketKeyInTxn(txn *memdb.Txn, bucketKey string) ([]*identity.Alias, error) {
	if txn == nil {
		return nil, fmt.Errorf("nil txn")
	}

	if bucketKey == "" {
		return nil, fmt.Errorf("empty bucket key")
	}

	iter, err := txn.Get(entityAliasesTable, "local_bucket_key", bucketKey)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup aliases using local bucket entry key hash: %w", err)
	}

	var aliases []*identity.Alias
	for item := iter.Next(); item != nil; item = iter.Next() {
		alias := item.(*identity.Alias)
		if alias.Local {
			aliases = append(aliases, alias)
		}
	}

	return aliases, nil
}

func (i *IdentityStore) MemDBEntitiesByBucketKeyInTxn(txn *memdb.Txn, bucketKey string) ([]*identity.Entity, error) {
	if txn == nil {
		return nil, fmt.Errorf("nil txn")
	}

	if bucketKey == "" {
		return nil, fmt.Errorf("empty bucket key")
	}

	entitiesIter, err := txn.Get(entitiesTable, "bucket_key", bucketKey)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup entities using bucket entry key hash: %w", err)
	}

	var entities []*identity.Entity
	for item := entitiesIter.Next(); item != nil; item = entitiesIter.Next() {
		entity, err := item.(*identity.Entity).Clone()
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func (i *IdentityStore) MemDBEntityByMergedEntityID(mergedEntityID string, clone bool) (*identity.Entity, error) {
	if mergedEntityID == "" {
		return nil, fmt.Errorf("missing merged entity id")
	}

	txn := i.db.Txn(false)

	entityRaw, err := txn.First(entitiesTable, "merged_entity_ids", mergedEntityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entity from memdb using merged entity id: %w", err)
	}

	if entityRaw == nil {
		return nil, nil
	}

	entity, ok := entityRaw.(*identity.Entity)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched entity")
	}

	if clone {
		return entity.Clone()
	}

	return entity, nil
}

func (i *IdentityStore) MemDBEntityByAliasIDInTxn(txn *memdb.Txn, aliasID string, clone bool) (*identity.Entity, error) {
	if aliasID == "" {
		return nil, fmt.Errorf("missing alias ID")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	alias, err := i.MemDBAliasByIDInTxn(txn, aliasID, false, false)
	if err != nil {
		return nil, err
	}

	if alias == nil {
		return nil, nil
	}

	return i.MemDBEntityByIDInTxn(txn, alias.CanonicalID, clone)
}

func (i *IdentityStore) MemDBEntityByAliasID(aliasID string, clone bool) (*identity.Entity, error) {
	if aliasID == "" {
		return nil, fmt.Errorf("missing alias ID")
	}

	txn := i.db.Txn(false)

	return i.MemDBEntityByAliasIDInTxn(txn, aliasID, clone)
}

func (i *IdentityStore) MemDBDeleteEntityByID(entityID string) error {
	if entityID == "" {
		return nil
	}

	txn := i.db.Txn(true)
	defer txn.Abort()

	err := i.MemDBDeleteEntityByIDInTxn(txn, entityID)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

// FetchEntityForLocalAliasInTxn fetches the entity associated with the provided
// local identity.Alias. MemDB will first be searched for the entity. If it is
// not found there, the localAliasPacker storagepacker.StoragePacker will be
// used. If an error occurs, an appropriate error message is logged and nil is
// returned.
func (i *IdentityStore) FetchEntityForLocalAliasInTxn(txn *memdb.Txn, alias *identity.Alias) *identity.Entity {
	entity, err := i.MemDBEntityByIDInTxn(txn, alias.CanonicalID, false)
	if err != nil {
		i.logger.Error("failed to fetch entity from local alias", "entity_id", alias.CanonicalID, "error", err)
		return nil
	}

	if entity == nil {
		cachedEntityItem, err := i.localAliasPacker.GetItem(alias.CanonicalID + tmpSuffix)
		if err != nil {
			i.logger.Error("failed to fetch cached entity from local alias", "key", alias.CanonicalID+tmpSuffix, "error", err)
			return nil
		}
		if cachedEntityItem != nil {
			entity, err = i.parseCachedEntity(cachedEntityItem)
			if err != nil {
				i.logger.Error("failed to parse cached entity", "key", alias.CanonicalID+tmpSuffix, "error", err)
				return nil
			}
		}
	}

	return entity
}

func (i *IdentityStore) MemDBDeleteEntityByIDInTxn(txn *memdb.Txn, entityID string) error {
	if entityID == "" {
		return nil
	}

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	entity, err := i.MemDBEntityByIDInTxn(txn, entityID, false)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	err = txn.Delete(entitiesTable, entity)
	if err != nil {
		return fmt.Errorf("failed to delete entity from memdb: %w", err)
	}

	return nil
}

func (i *IdentityStore) sanitizeAlias(ctx context.Context, alias *identity.Alias) error {
	var err error

	if alias == nil {
		return fmt.Errorf("alias is nil")
	}

	// Alias must always be tied to a canonical object
	if alias.CanonicalID == "" {
		return fmt.Errorf("missing canonical ID")
	}

	// Alias must have a name
	if alias.Name == "" {
		return fmt.Errorf("missing alias name %q", alias.Name)
	}

	// Alias metadata should always be map[string]string
	err = validateMetadata(alias.Metadata)
	if err != nil {
		return fmt.Errorf("invalid alias metadata: %w", err)
	}

	// Create an ID if there isn't one already
	if alias.ID == "" {
		alias.ID, err = uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("failed to generate alias ID")
		}

		alias.LocalBucketKey = i.localAliasPacker.BucketKey(alias.CanonicalID)
	}

	if alias.NamespaceID == "" {
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return err
		}
		alias.NamespaceID = ns.ID
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	if ns.ID != alias.NamespaceID {
		return errors.New("alias belongs to a different namespace")
	}

	// Set the creation and last update times
	if alias.CreationTime == nil {
		alias.CreationTime = timestamppb.Now()
		alias.LastUpdateTime = alias.CreationTime
	} else {
		alias.LastUpdateTime = timestamppb.Now()
	}

	return nil
}

func (i *IdentityStore) sanitizeEntity(ctx context.Context, entity *identity.Entity) error {
	var err error

	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	// Create an ID if there isn't one already
	if entity.ID == "" {
		entity.ID, err = uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("failed to generate entity id")
		}

		// Set the storage bucket key in entity
		entity.BucketKey = i.entityPacker.BucketKey(entity.ID)
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	if entity.NamespaceID == "" {
		entity.NamespaceID = ns.ID
	}
	if ns.ID != entity.NamespaceID {
		return errors.New("entity does not belong to this namespace")
	}

	// Create a name if there isn't one already
	if entity.Name == "" {
		entity.Name, err = i.generateName(ctx, "entity")
		if err != nil {
			return fmt.Errorf("failed to generate entity name")
		}
	}

	// Entity metadata should always be map[string]string
	err = validateMetadata(entity.Metadata)
	if err != nil {
		return fmt.Errorf("invalid entity metadata: %w", err)
	}

	// Set the creation and last update times
	if entity.CreationTime == nil {
		entity.CreationTime = timestamppb.Now()
		entity.LastUpdateTime = entity.CreationTime
	} else {
		entity.LastUpdateTime = timestamppb.Now()
	}

	// Ensure that MFASecrets is non-nil at any time. This is useful when MFA
	// secret generation procedures try to append MFA info to entity.
	if entity.MFASecrets == nil {
		entity.MFASecrets = make(map[string]*mfa.Secret)
	}

	return nil
}

func (i *IdentityStore) sanitizeAndUpsertGroup(ctx context.Context, group *identity.Group, previousGroup *identity.Group, memberGroupIDs []string) error {
	var err error

	if group == nil {
		return fmt.Errorf("group is nil")
	}

	// Create an ID if there isn't one already
	if group.ID == "" {
		group.ID, err = uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("failed to generate group id")
		}

		// Set the hash value of the storage bucket key in group
		group.BucketKey = i.groupPacker.BucketKey(group.ID)
	}

	if group.NamespaceID == "" {
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return err
		}
		group.NamespaceID = ns.ID
	}
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	if ns.ID != group.NamespaceID {
		return errors.New("group does not belong to this namespace")
	}

	// Create a name if there isn't one already
	if group.Name == "" {
		group.Name, err = i.generateName(ctx, "group")
		if err != nil {
			return fmt.Errorf("failed to generate group name")
		}
	}

	// Entity metadata should always be map[string]string
	err = validateMetadata(group.Metadata)
	if err != nil {
		return fmt.Errorf("invalid group metadata: %w", err)
	}

	// Set the creation and last update times
	if group.CreationTime == nil {
		group.CreationTime = timestamppb.Now()
		group.LastUpdateTime = group.CreationTime
	} else {
		group.LastUpdateTime = timestamppb.Now()
	}

	// Remove duplicate entity IDs and check if all IDs are valid
	if group.MemberEntityIDs != nil {
		group.MemberEntityIDs = strutil.RemoveDuplicates(group.MemberEntityIDs, false)
		for _, entityID := range group.MemberEntityIDs {
			entity, err := i.MemDBEntityByID(entityID, false)
			if err != nil {
				return fmt.Errorf("failed to validate entity ID %q: %w", entityID, err)
			}
			if entity == nil {
				return fmt.Errorf("invalid entity ID %q", entityID)
			}
		}
	}

	// Remove duplicate policies
	if group.Policies != nil {
		group.Policies = strutil.RemoveDuplicates(group.Policies, false)
	}

	txn := i.db.Txn(true)
	defer txn.Abort()

	var currentMemberGroupIDs []string
	var currentMemberGroups []*identity.Group

	// If there are no member group IDs supplied, then it shouldn't be
	// processed. If an empty set of member group IDs are supplied, then it
	// should be processed. Hence the nil check instead of the length check.
	if memberGroupIDs == nil {
		goto ALIAS
	}

	memberGroupIDs = strutil.RemoveDuplicates(memberGroupIDs, false)

	// For those group member IDs that are removed from the list, remove current
	// group ID as their respective ParentGroupID.

	// Get the current MemberGroups IDs for this group
	currentMemberGroups, err = i.MemDBGroupsByParentGroupID(group.ID, false)
	if err != nil {
		return err
	}
	for _, currentMemberGroup := range currentMemberGroups {
		currentMemberGroupIDs = append(currentMemberGroupIDs, currentMemberGroup.ID)
	}

	// Update parent group IDs in the removed members
	for _, currentMemberGroupID := range currentMemberGroupIDs {
		if strutil.StrListContains(memberGroupIDs, currentMemberGroupID) {
			continue
		}

		currentMemberGroup, err := i.MemDBGroupByID(currentMemberGroupID, true)
		if err != nil {
			return err
		}
		if currentMemberGroup == nil {
			return fmt.Errorf("invalid member group ID %q", currentMemberGroupID)
		}

		// Remove group ID from the parent group IDs
		currentMemberGroup.ParentGroupIDs = strutil.StrListDelete(currentMemberGroup.ParentGroupIDs, group.ID)

		err = i.UpsertGroupInTxn(ctx, txn, currentMemberGroup, true)
		if err != nil {
			return err
		}
	}

	// After the group lock is held, make membership updates to all the
	// relevant groups
	for _, memberGroupID := range memberGroupIDs {
		memberGroup, err := i.MemDBGroupByID(memberGroupID, true)
		if err != nil {
			return err
		}
		if memberGroup == nil {
			return fmt.Errorf("invalid member group ID %q", memberGroupID)
		}

		// Skip if memberGroupID is already a member of group.ID
		if strutil.StrListContains(memberGroup.ParentGroupIDs, group.ID) {
			continue
		}

		// Ensure that adding memberGroupID does not lead to cyclic
		// relationships
		// Detect self loop
		if group.ID == memberGroupID {
			return fmt.Errorf("member group ID %q is same as the ID of the group", group.ID)
		}

		groupByID, err := i.MemDBGroupByID(group.ID, true)
		if err != nil {
			return err
		}

		// If group is nil, that means that a group doesn't already exist and its
		// okay to add any group as its member group.
		if groupByID != nil {
			// If adding the memberGroupID to groupID creates a cycle, then groupID must
			// be a hop in that loop. Start a DFS traversal from memberGroupID and see if
			// it reaches back to groupID. If it does, then it's a loop.

			// Created a visited set
			visited := make(map[string]bool)
			cycleDetected, err := i.detectCycleDFS(visited, groupByID.ID, memberGroupID)
			if err != nil {
				return fmt.Errorf("failed to perform cyclic relationship detection for member group ID %q", memberGroupID)
			}
			if cycleDetected {
				return fmt.Errorf("%s %q", errCycleDetectedPrefix, memberGroupID)
			}
		}

		memberGroup.ParentGroupIDs = append(memberGroup.ParentGroupIDs, group.ID)

		// This technically is not upsert. It is only update, only the method
		// name is upsert here.
		err = i.UpsertGroupInTxn(ctx, txn, memberGroup, true)
		if err != nil {
			// Ideally we would want to revert the whole operation in case of
			// errors while persisting in member groups. But there is no
			// storage transaction support yet. When we do have it, this will need
			// an update.
			return err
		}
	}

ALIAS:
	// Sanitize the group alias
	if group.Alias != nil {
		group.Alias.CanonicalID = group.ID
		err = i.sanitizeAlias(ctx, group.Alias)
		if err != nil {
			return err
		}
	}

	// If previousGroup is not nil, we are moving the alias from the previous
	// group to the new one. As a result we need to upsert both in the context
	// of this same transaction.
	if previousGroup != nil {
		err = i.UpsertGroupInTxn(ctx, txn, previousGroup, true)
		if err != nil {
			return err
		}
	}

	err = i.UpsertGroupInTxn(ctx, txn, group, true)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (i *IdentityStore) deleteAliasesInEntityInTxn(txn *memdb.Txn, entity *identity.Entity, aliases []*identity.Alias) error {
	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	var remainList []*identity.Alias
	var removeList []*identity.Alias
	for _, item := range entity.Aliases {
		remove := false
		for _, alias := range aliases {
			if alias.ID == item.ID {
				remove = true
			}
		}
		if remove {
			removeList = append(removeList, item)
		} else {
			remainList = append(remainList, item)
		}
	}

	// Remove identity indices from aliases table for those that needs to
	// be removed
	for _, alias := range removeList {
		err := i.MemDBDeleteAliasByIDInTxn(txn, alias.ID, false)
		if err != nil {
			return err
		}
	}

	// Update the entity with remaining items
	entity.Aliases = remainList

	return nil
}

// validateMeta validates a set of key/value pairs from the agent config
func validateMetadata(meta map[string]string) error {
	if len(meta) > metaMaxKeyPairs {
		return fmt.Errorf("metadata cannot contain more than %d key/value pairs", metaMaxKeyPairs)
	}

	for key, value := range meta {
		if err := validateMetaPair(key, value); err != nil {
			return fmt.Errorf("failed to load metadata pair (%q, %q): %w", key, value, err)
		}
	}

	return nil
}

// validateMetaPair checks that the given key/value pair is in a valid format
func validateMetaPair(key, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be blank")
	}
	if !metaKeyFormatRegEx(key) {
		return fmt.Errorf("key contains invalid characters")
	}
	if len(key) > metaKeyMaxLength {
		return fmt.Errorf("key is too long (limit: %d characters)", metaKeyMaxLength)
	}
	if strings.HasPrefix(key, metaKeyReservedPrefix) {
		return fmt.Errorf("key prefix %q is reserved for internal use", metaKeyReservedPrefix)
	}
	if len(value) > metaValueMaxLength {
		return fmt.Errorf("value is too long (limit: %d characters)", metaValueMaxLength)
	}
	return nil
}

func (i *IdentityStore) MemDBGroupByNameInTxn(ctx context.Context, txn *memdb.Txn, groupName string, clone bool) (*identity.Group, error) {
	if groupName == "" {
		return nil, fmt.Errorf("missing group name")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	groupRaw, err := txn.First(groupsTable, "name", ns.ID, groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group from memdb using group name: %w", err)
	}

	if groupRaw == nil {
		return nil, nil
	}

	group, ok := groupRaw.(*identity.Group)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched group")
	}

	if clone {
		return group.Clone()
	}

	return group, nil
}

func (i *IdentityStore) MemDBGroupByName(ctx context.Context, groupName string, clone bool) (*identity.Group, error) {
	if groupName == "" {
		return nil, fmt.Errorf("missing group name")
	}

	txn := i.db.Txn(false)

	return i.MemDBGroupByNameInTxn(ctx, txn, groupName, clone)
}

func (i *IdentityStore) UpsertGroup(ctx context.Context, group *identity.Group, persist bool) error {
	defer metrics.MeasureSince([]string{"identity", "upsert_group"}, time.Now())

	txn := i.db.Txn(true)
	defer txn.Abort()

	err := i.UpsertGroupInTxn(ctx, txn, group, persist)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (i *IdentityStore) UpsertGroupInTxn(ctx context.Context, txn *memdb.Txn, group *identity.Group, persist bool) error {
	defer metrics.MeasureSince([]string{"identity", "upsert_group_txn"}, time.Now())

	var err error

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	if group == nil {
		return fmt.Errorf("group is nil")
	}

	// Increment the modify index of the group
	group.ModifyIndex++

	// Clear the old alias from memdb
	groupClone, err := i.MemDBGroupByID(group.ID, true)
	if err != nil {
		return err
	}
	if groupClone != nil && groupClone.Alias != nil {
		err = i.MemDBDeleteAliasByIDInTxn(txn, groupClone.Alias.ID, true)
		if err != nil {
			return err
		}
	}

	// Add the new alias to memdb
	if group.Alias != nil {
		err = i.MemDBUpsertAliasInTxn(txn, group.Alias, true)
		if err != nil {
			return err
		}
	}

	// Insert or update group in MemDB using the transaction created above
	err = i.MemDBUpsertGroupInTxn(txn, group)
	if err != nil {
		return err
	}

	if persist {
		groupAsAny, err := anypb.New(group)
		if err != nil {
			return err
		}

		item := &storagepacker.Item{
			ID:      group.ID,
			Message: groupAsAny,
		}

		sent, err := i.groupUpdater.SendGroupUpdate(ctx, group)
		if err != nil {
			return err
		}
		if !sent {
			if err := i.groupPacker.PutItem(ctx, item); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *IdentityStore) MemDBUpsertGroupInTxn(txn *memdb.Txn, group *identity.Group) error {
	if txn == nil {
		return fmt.Errorf("nil txn")
	}

	if group == nil {
		return fmt.Errorf("group is nil")
	}

	if group.NamespaceID == "" {
		group.NamespaceID = namespace.RootNamespaceID
	}

	groupRaw, err := txn.First(groupsTable, "id", group.ID)
	if err != nil {
		return fmt.Errorf("failed to lookup group from memdb using group id: %w", err)
	}

	if groupRaw != nil {
		err = txn.Delete(groupsTable, groupRaw)
		if err != nil {
			return fmt.Errorf("failed to delete group from memdb: %w", err)
		}
	}

	if err := txn.Insert(groupsTable, group); err != nil {
		return fmt.Errorf("failed to update group into memdb: %w", err)
	}

	return nil
}

func (i *IdentityStore) MemDBDeleteGroupByIDInTxn(txn *memdb.Txn, groupID string) error {
	if groupID == "" {
		return nil
	}

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	group, err := i.MemDBGroupByIDInTxn(txn, groupID, false)
	if err != nil {
		return err
	}

	if group == nil {
		return nil
	}

	err = txn.Delete("groups", group)
	if err != nil {
		return fmt.Errorf("failed to delete group from memdb: %w", err)
	}

	return nil
}

func (i *IdentityStore) MemDBGroupByIDInTxn(txn *memdb.Txn, groupID string, clone bool) (*identity.Group, error) {
	if groupID == "" {
		return nil, fmt.Errorf("missing group ID")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	groupRaw, err := txn.First(groupsTable, "id", groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group from memdb using group ID: %w", err)
	}

	if groupRaw == nil {
		return nil, nil
	}

	group, ok := groupRaw.(*identity.Group)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched group")
	}

	if clone {
		return group.Clone()
	}

	return group, nil
}

func (i *IdentityStore) MemDBGroupByID(groupID string, clone bool) (*identity.Group, error) {
	if groupID == "" {
		return nil, fmt.Errorf("missing group ID")
	}

	txn := i.db.Txn(false)

	return i.MemDBGroupByIDInTxn(txn, groupID, clone)
}

func (i *IdentityStore) MemDBGroupsByParentGroupIDInTxn(txn *memdb.Txn, memberGroupID string, clone bool) ([]*identity.Group, error) {
	if memberGroupID == "" {
		return nil, fmt.Errorf("missing member group ID")
	}

	groupsIter, err := txn.Get(groupsTable, "parent_group_ids", memberGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup groups using member group ID: %w", err)
	}

	var groups []*identity.Group
	for group := groupsIter.Next(); group != nil; group = groupsIter.Next() {
		entry := group.(*identity.Group)
		if clone {
			entry, err = entry.Clone()
			if err != nil {
				return nil, err
			}
		}
		groups = append(groups, entry)
	}

	return groups, nil
}

func (i *IdentityStore) MemDBGroupsByParentGroupID(memberGroupID string, clone bool) ([]*identity.Group, error) {
	if memberGroupID == "" {
		return nil, fmt.Errorf("missing member group ID")
	}

	txn := i.db.Txn(false)

	return i.MemDBGroupsByParentGroupIDInTxn(txn, memberGroupID, clone)
}

func (i *IdentityStore) MemDBGroupsByMemberEntityID(entityID string, clone bool, externalOnly bool) ([]*identity.Group, error) {
	txn := i.db.Txn(false)
	defer txn.Abort()

	return i.MemDBGroupsByMemberEntityIDInTxn(txn, entityID, clone, externalOnly)
}

func (i *IdentityStore) MemDBGroupsByMemberEntityIDInTxn(txn *memdb.Txn, entityID string, clone bool, externalOnly bool) ([]*identity.Group, error) {
	if entityID == "" {
		return nil, fmt.Errorf("missing entity ID")
	}

	groupsIter, err := txn.Get(groupsTable, "member_entity_ids", entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup groups using entity ID: %w", err)
	}

	var groups []*identity.Group
	for group := groupsIter.Next(); group != nil; group = groupsIter.Next() {
		entry := group.(*identity.Group)
		if externalOnly && entry.Type == groupTypeInternal {
			continue
		}
		if clone {
			entry, err = entry.Clone()
			if err != nil {
				return nil, err
			}
		}
		groups = append(groups, entry)
	}

	return groups, nil
}

func (i *IdentityStore) groupPoliciesByEntityID(entityID string) (map[string][]string, error) {
	if entityID == "" {
		return nil, fmt.Errorf("empty entity ID")
	}

	groups, err := i.MemDBGroupsByMemberEntityID(entityID, false, false)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	policies := make(map[string][]string)
	for _, group := range groups {
		err := i.collectPoliciesReverseDFS(group, visited, policies)
		if err != nil {
			return nil, err
		}
	}

	return policies, nil
}

func (i *IdentityStore) groupsByEntityID(entityID string) ([]*identity.Group, []*identity.Group, error) {
	if entityID == "" {
		return nil, nil, fmt.Errorf("empty entity ID")
	}

	groups, err := i.MemDBGroupsByMemberEntityID(entityID, true, false)
	if err != nil {
		return nil, nil, err
	}

	visited := make(map[string]bool)
	var tGroups []*identity.Group
	for _, group := range groups {
		gGroups, err := i.collectGroupsReverseDFS(group, visited, nil)
		if err != nil {
			return nil, nil, err
		}
		tGroups = append(tGroups, gGroups...)
	}

	// Remove duplicates
	groupMap := make(map[string]*identity.Group)
	for _, group := range tGroups {
		groupMap[group.ID] = group
	}

	tGroups = make([]*identity.Group, 0, len(groupMap))
	for _, group := range groupMap {
		tGroups = append(tGroups, group)
	}

	diff := diffGroups(groups, tGroups)

	// For sanity
	// There should not be any group that gets deleted
	if len(diff.Deleted) != 0 {
		return nil, nil, fmt.Errorf("failed to diff group memberships")
	}

	return diff.Unmodified, diff.New, nil
}

func (i *IdentityStore) collectGroupsReverseDFS(group *identity.Group, visited map[string]bool, groups []*identity.Group) ([]*identity.Group, error) {
	if group == nil {
		return nil, fmt.Errorf("nil group")
	}

	// If traversal for a groupID is performed before, skip it
	if visited[group.ID] {
		return groups, nil
	}
	visited[group.ID] = true

	groups = append(groups, group)

	// Traverse all the parent groups
	for _, parentGroupID := range group.ParentGroupIDs {
		parentGroup, err := i.MemDBGroupByID(parentGroupID, false)
		if err != nil {
			return nil, err
		}
		if parentGroup == nil {
			continue
		}
		groups, err = i.collectGroupsReverseDFS(parentGroup, visited, groups)
		if err != nil {
			return nil, fmt.Errorf("failed to collect group at parent group ID %q", parentGroup.ID)
		}
	}

	return groups, nil
}

func (i *IdentityStore) collectPoliciesReverseDFS(group *identity.Group, visited map[string]bool, policies map[string][]string) error {
	if group == nil {
		return fmt.Errorf("nil group")
	}

	// If traversal for a groupID is performed before, skip it
	if visited[group.ID] {
		return nil
	}
	visited[group.ID] = true

	policies[group.NamespaceID] = append(policies[group.NamespaceID], group.Policies...)

	// Traverse all the parent groups
	for _, parentGroupID := range group.ParentGroupIDs {
		parentGroup, err := i.MemDBGroupByID(parentGroupID, false)
		if err != nil {
			return err
		}
		if parentGroup == nil {
			continue
		}
		err = i.collectPoliciesReverseDFS(parentGroup, visited, policies)
		if err != nil {
			return fmt.Errorf("failed to collect policies at parent group ID %q", parentGroup.ID)
		}
	}

	return nil
}

func (i *IdentityStore) detectCycleDFS(visited map[string]bool, startingGroupID, groupID string) (bool, error) {
	// If the traversal reaches the startingGroupID, a loop is detected
	if startingGroupID == groupID {
		return true, nil
	}

	// If traversal for a groupID is performed before, skip it
	if visited[groupID] {
		return false, nil
	}
	visited[groupID] = true

	group, err := i.MemDBGroupByID(groupID, true)
	if err != nil {
		return false, err
	}
	if group == nil {
		return false, nil
	}

	// Fetch all groups in which groupID is present as a ParentGroupID. In
	// other words, find all the subgroups of groupID.
	memberGroups, err := i.MemDBGroupsByParentGroupID(groupID, false)
	if err != nil {
		return false, err
	}

	// DFS traverse the member groups
	for _, memberGroup := range memberGroups {
		cycleDetected, err := i.detectCycleDFS(visited, startingGroupID, memberGroup.ID)
		if err != nil {
			return false, fmt.Errorf("failed to perform cycle detection at member group ID %q", memberGroup.ID)
		}
		if cycleDetected {
			return true, nil
		}
	}

	return false, nil
}

func (i *IdentityStore) memberGroupIDsByID(groupID string) ([]string, error) {
	var memberGroupIDs []string
	memberGroups, err := i.MemDBGroupsByParentGroupID(groupID, false)
	if err != nil {
		return nil, err
	}
	for _, memberGroup := range memberGroups {
		memberGroupIDs = append(memberGroupIDs, memberGroup.ID)
	}
	return memberGroupIDs, nil
}

func (i *IdentityStore) generateName(ctx context.Context, entryType string) (string, error) {
	var name string
OUTER:
	for {
		randBytes, err := uuid.GenerateRandomBytes(4)
		if err != nil {
			return "", err
		}
		name = fmt.Sprintf("%s_%s", entryType, fmt.Sprintf("%08x", randBytes[0:4]))

		switch entryType {
		case "entity":
			entity, err := i.MemDBEntityByName(ctx, name, false)
			if err != nil {
				return "", err
			}
			if entity == nil {
				break OUTER
			}
		case "group":
			group, err := i.MemDBGroupByName(ctx, name, false)
			if err != nil {
				return "", err
			}
			if group == nil {
				break OUTER
			}
		default:
			return "", fmt.Errorf("unrecognized type %q", entryType)
		}
	}

	return name, nil
}

func (i *IdentityStore) MemDBGroupsByBucketKeyInTxn(txn *memdb.Txn, bucketKey string) ([]*identity.Group, error) {
	if txn == nil {
		return nil, fmt.Errorf("nil txn")
	}

	if bucketKey == "" {
		return nil, fmt.Errorf("empty bucket key")
	}

	groupsIter, err := txn.Get(groupsTable, "bucket_key", bucketKey)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup groups using bucket entry key hash: %w", err)
	}

	var groups []*identity.Group
	for group := groupsIter.Next(); group != nil; group = groupsIter.Next() {
		groups = append(groups, group.(*identity.Group))
	}

	return groups, nil
}

func (i *IdentityStore) MemDBGroupByAliasIDInTxn(txn *memdb.Txn, aliasID string, clone bool) (*identity.Group, error) {
	if aliasID == "" {
		return nil, fmt.Errorf("missing alias ID")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	alias, err := i.MemDBAliasByIDInTxn(txn, aliasID, false, true)
	if err != nil {
		return nil, err
	}

	if alias == nil {
		return nil, nil
	}

	return i.MemDBGroupByIDInTxn(txn, alias.CanonicalID, clone)
}

func (i *IdentityStore) MemDBGroupByAliasID(aliasID string, clone bool) (*identity.Group, error) {
	if aliasID == "" {
		return nil, fmt.Errorf("missing alias ID")
	}

	txn := i.db.Txn(false)

	return i.MemDBGroupByAliasIDInTxn(txn, aliasID, clone)
}

func (i *IdentityStore) refreshExternalGroupMembershipsByEntityID(ctx context.Context, entityID string, groupAliases []*logical.Alias, mountAccessor string) ([]*logical.Alias, error) {
	defer metrics.MeasureSince([]string{"identity", "refresh_external_groups"}, time.Now())

	if entityID == "" {
		return nil, fmt.Errorf("empty entity ID")
	}

	refreshFunc := func(dryRun bool) (bool, []*logical.Alias, error) {
		if !dryRun {
			i.groupLock.Lock()
			defer i.groupLock.Unlock()
		}

		txn := i.db.Txn(!dryRun)
		defer txn.Abort()

		oldGroups, err := i.MemDBGroupsByMemberEntityIDInTxn(txn, entityID, true, true)
		if err != nil {
			return false, nil, err
		}

		var newGroups []*identity.Group
		var validAliases []*logical.Alias
		for _, alias := range groupAliases {
			aliasByFactors, err := i.MemDBAliasByFactorsInTxn(txn, alias.MountAccessor, alias.Name, true, true)
			if err != nil {
				return false, nil, err
			}
			if aliasByFactors == nil {
				continue
			}
			mappingGroup, err := i.MemDBGroupByAliasIDInTxn(txn, aliasByFactors.ID, true)
			if err != nil {
				return false, nil, err
			}
			if mappingGroup == nil {
				return false, nil, fmt.Errorf("group unavailable for a valid alias ID %q", aliasByFactors.ID)
			}

			newGroups = append(newGroups, mappingGroup)
			validAliases = append(validAliases, alias)
		}

		diff := diffGroups(oldGroups, newGroups)

		// Add the entity ID to all the new groups
		for _, group := range diff.New {
			if group.Type != groupTypeExternal {
				continue
			}

			// We need to update a group, if we are in a dry run we should
			// report back that a change needs to take place.
			if dryRun {
				return true, nil, nil
			}

			i.logger.Debug("adding member entity ID to external group", "member_entity_id", entityID, "group_id", group.ID)

			group.MemberEntityIDs = append(group.MemberEntityIDs, entityID)

			err = i.UpsertGroupInTxn(ctx, txn, group, true)
			if err != nil {
				return false, nil, err
			}
		}

		// Remove the entity ID from all the deleted groups
		for _, group := range diff.Deleted {
			if group.Type != groupTypeExternal {
				continue
			}

			// If the external group is from a different mount, don't remove the
			// entity ID from it.
			if mountAccessor != "" && group.Alias != nil && group.Alias.MountAccessor != mountAccessor {
				continue
			}

			// We need to update a group, if we are in a dry run we should
			// report back that a change needs to take place.
			if dryRun {
				return true, nil, nil
			}

			i.logger.Debug("removing member entity ID from external group", "member_entity_id", entityID, "group_id", group.ID)

			group.MemberEntityIDs = strutil.StrListDelete(group.MemberEntityIDs, entityID)

			err = i.UpsertGroupInTxn(ctx, txn, group, true)
			if err != nil {
				return false, nil, err
			}
		}

		txn.Commit()
		return false, validAliases, nil
	}

	// dryRun
	needsUpdate, validAliases, err := refreshFunc(true)
	if err != nil {
		return nil, err
	}

	if needsUpdate || len(groupAliases) > 0 {
		i.logger.Debug("refreshing external group memberships", "entity_id", entityID, "group_aliases", groupAliases)
	}

	if !needsUpdate {
		return validAliases, nil
	}

	// Run the update
	_, validAliases, err = refreshFunc(false)
	if err != nil {
		return nil, err
	}

	return validAliases, nil
}

// diffGroups is used to diff two sets of groups
func diffGroups(old, new []*identity.Group) *groupDiff {
	diff := &groupDiff{}

	existing := make(map[string]*identity.Group)
	for _, group := range old {
		existing[group.ID] = group
	}

	for _, group := range new {
		// Check if the entry in new is present in the old
		_, ok := existing[group.ID]

		// If its not present, then its a new entry
		if !ok {
			diff.New = append(diff.New, group)
			continue
		}

		// If its present, it means that its unmodified
		diff.Unmodified = append(diff.Unmodified, group)

		// By deleting the unmodified from the old set, we could determine the
		// ones that are stale by looking at the remaining ones.
		delete(existing, group.ID)
	}

	// Any remaining entries must have been deleted
	for _, me := range existing {
		diff.Deleted = append(diff.Deleted, me)
	}

	return diff
}

func (i *IdentityStore) handleAliasListCommon(ctx context.Context, groupAlias bool) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	tableName := entityAliasesTable
	if groupAlias {
		tableName = groupAliasesTable
	}

	ws := memdb.NewWatchSet()

	txn := i.db.Txn(false)

	iter, err := txn.Get(tableName, "namespace_id", ns.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch iterator for aliases in memdb: %w", err)
	}

	ws.Add(iter.WatchCh())

	var aliasIDs []string
	aliasInfo := map[string]interface{}{}

	type mountInfo struct {
		MountType string
		MountPath string
	}
	mountAccessorMap := map[string]mountInfo{}

	for {
		raw := iter.Next()
		if raw == nil {
			break
		}
		alias := raw.(*identity.Alias)
		aliasIDs = append(aliasIDs, alias.ID)
		aliasInfoEntry := map[string]interface{}{
			"name":            alias.Name,
			"canonical_id":    alias.CanonicalID,
			"mount_accessor":  alias.MountAccessor,
			"custom_metadata": alias.CustomMetadata,
			"metadata":        alias.Metadata,
			"local":           alias.Local,
		}

		mi, ok := mountAccessorMap[alias.MountAccessor]
		if ok {
			aliasInfoEntry["mount_type"] = mi.MountType
			aliasInfoEntry["mount_path"] = mi.MountPath
		} else {
			mi = mountInfo{}
			if mountValidationResp := i.router.ValidateMountByAccessor(alias.MountAccessor); mountValidationResp != nil {
				mi.MountType = mountValidationResp.MountType
				mi.MountPath = mountValidationResp.MountPath
				aliasInfoEntry["mount_type"] = mi.MountType
				aliasInfoEntry["mount_path"] = mi.MountPath
			}
			mountAccessorMap[alias.MountAccessor] = mi
		}

		aliasInfo[alias.ID] = aliasInfoEntry
	}

	return logical.ListResponseWithInfo(aliasIDs, aliasInfo), nil
}

func (i *IdentityStore) countEntities() (int, error) {
	txn := i.db.Txn(false)

	iter, err := txn.Get(entitiesTable, "id")
	if err != nil {
		return -1, err
	}

	count := 0
	val := iter.Next()
	for val != nil {
		count++
		val = iter.Next()
	}

	return count, nil
}

// Sum up the number of entities belonging to each namespace (keyed by ID)
func (i *IdentityStore) countEntitiesByNamespace(ctx context.Context) (map[string]int, error) {
	txn := i.db.Txn(false)
	iter, err := txn.Get(entitiesTable, "id")
	if err != nil {
		return nil, err
	}

	byNamespace := make(map[string]int)
	val := iter.Next()
	for val != nil {
		// Check if runtime exceeded.
		select {
		case <-ctx.Done():
			return byNamespace, errors.New("context cancelled")
		default:
			break
		}

		// Count in the namespace attached to the entity.
		entity := val.(*identity.Entity)
		byNamespace[entity.NamespaceID] = byNamespace[entity.NamespaceID] + 1
		val = iter.Next()
	}

	return byNamespace, nil
}

// Sum up the number of entities belonging to each mount point (keyed by accessor)
func (i *IdentityStore) countEntitiesByMountAccessor(ctx context.Context) (map[string]int, error) {
	txn := i.db.Txn(false)
	iter, err := txn.Get(entitiesTable, "id")
	if err != nil {
		return nil, err
	}

	byMountAccessor := make(map[string]int)
	val := iter.Next()
	for val != nil {
		// Check if runtime exceeded.
		select {
		case <-ctx.Done():
			return byMountAccessor, errors.New("context cancelled")
		default:
			break
		}

		// Count each alias separately; will translate to mount point and type
		// in the caller.
		entity := val.(*identity.Entity)
		for _, alias := range entity.Aliases {
			byMountAccessor[alias.MountAccessor] = byMountAccessor[alias.MountAccessor] + 1
		}
		val = iter.Next()
	}

	return byMountAccessor, nil
}

func makeEntityForPacker(t *testing.T, name string, p *storagepacker.StoragePacker) *identity.Entity {
	t.Helper()
	return makeEntityForPackerWithNamespace(t, namespace.RootNamespaceID, name, p)
}

func makeEntityForPackerWithNamespace(t *testing.T, namespaceID, name string, p *storagepacker.StoragePacker) *identity.Entity {
	t.Helper()
	id, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return &identity.Entity{
		ID:          id,
		Name:        name,
		NamespaceID: namespaceID,
		BucketKey:   p.BucketKey(id),
	}
}

func attachAlias(t *testing.T, e *identity.Entity, name string, me *MountEntry) *identity.Alias {
	t.Helper()
	id, err := uuid.GenerateUUID()
	require.NoError(t, err)
	if e.NamespaceID != me.NamespaceID {
		panic("mount and entity in different namespaces")
	}
	require.NoError(t, err)
	a := &identity.Alias{
		ID:            id,
		Name:          name,
		NamespaceID:   me.NamespaceID,
		CanonicalID:   e.ID,
		MountType:     me.Type,
		MountAccessor: me.Accessor,
		Local:         me.Local,
	}
	e.UpsertAlias(a)
	return a
}

func identityCreateCaseDuplicates(t *testing.T, ctx context.Context, c *Core, upme, localme *MountEntry) {
	t.Helper()

	if upme.NamespaceID != localme.NamespaceID {
		panic("both replicated and local auth mounts must be in the same namespace")
	}

	// Create entities with both case-sensitive and case-insensitive duplicate
	// suffixes.
	for i, suffix := range []string{"-case", "-case", "-cAsE"} {
		// Entity duplicated by name
		e := makeEntityForPackerWithNamespace(t, upme.NamespaceID, "entity"+suffix, c.identityStore.entityPacker)
		err := TestHelperWriteToStoragePacker(ctx, c.identityStore.entityPacker, e.ID, e)
		require.NoError(t, err)

		// Entity that isn't a dupe itself but has duplicated aliases
		e2 := makeEntityForPackerWithNamespace(t, upme.NamespaceID, fmt.Sprintf("entity-%d", i), c.identityStore.entityPacker)
		// Add local and non-local aliases for this entity (which will also be
		// duplicated)
		attachAlias(t, e2, "alias"+suffix, upme)
		attachAlias(t, e2, "local-alias"+suffix, localme)
		err = TestHelperWriteToStoragePacker(ctx, c.identityStore.entityPacker, e2.ID, e2)
		require.NoError(t, err)

		// Group duplicated by name
		g := makeGroupWithNameAndAlias(t, "group"+suffix, "", c.identityStore.groupPacker, upme)
		err = TestHelperWriteToStoragePacker(ctx, c.identityStore.groupPacker, g.ID, g)
		require.NoError(t, err)
	}
}

func makeGroupWithNameAndAlias(t *testing.T, name, alias string, p *storagepacker.StoragePacker, me *MountEntry) *identity.Group {
	t.Helper()
	id, err := uuid.GenerateUUID()
	require.NoError(t, err)
	id2, err := uuid.GenerateUUID()
	require.NoError(t, err)
	g := &identity.Group{
		ID:          id,
		Name:        name,
		NamespaceID: me.NamespaceID,
		BucketKey:   p.BucketKey(id),
	}
	if alias != "" {
		g.Alias = &identity.Alias{
			ID:            id2,
			Name:          alias,
			CanonicalID:   id,
			MountType:     me.Type,
			MountAccessor: me.Accessor,
			NamespaceID:   me.NamespaceID,
		}
	}
	return g
}

func makeLocalAliasWithName(t *testing.T, name, entityID string, bucketKey string, me *MountEntry) *identity.LocalAliases {
	t.Helper()
	id, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return &identity.LocalAliases{
		Aliases: []*identity.Alias{
			{
				ID:            id,
				Name:          name,
				CanonicalID:   entityID,
				MountType:     me.Type,
				MountAccessor: me.Accessor,
				NamespaceID:   me.NamespaceID,
			},
		},
	}
}
