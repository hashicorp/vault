package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	groupBucketsPrefix = "packer/group/buckets/"
)

var (
	sendGroupUpgrade             = func(*IdentityStore, *identity.Group) (bool, error) { return false, nil }
	parseExtraEntityFromBucket   = func(context.Context, *IdentityStore, *identity.Entity) (bool, error) { return false, nil }
	addExtraEntityDataToResponse = func(*identity.Entity, map[string]interface{}) {}
)

func (c *Core) IdentityStore() *IdentityStore {
	return c.identityStore
}

func (i *IdentityStore) resetDB(ctx context.Context) error {
	var err error

	i.db, err = memdb.NewMemDB(identityStoreSchema(!i.disableLowerCasedNames))
	if err != nil {
		return err
	}

	return nil
}

func NewIdentityStore(ctx context.Context, core *Core, config *logical.BackendConfig, logger log.Logger) (*IdentityStore, error) {
	iStore := &IdentityStore{
		view:   config.StorageView,
		logger: logger,
		core:   core,
	}

	// Create a memdb instance, which by default, operates on lower cased
	// identity names
	err := iStore.resetDB(ctx)
	if err != nil {
		return nil, err
	}

	entitiesPackerLogger := iStore.logger.Named("storagepacker").Named("entities")
	core.AddLogger(entitiesPackerLogger)
	groupsPackerLogger := iStore.logger.Named("storagepacker").Named("groups")
	core.AddLogger(groupsPackerLogger)
	iStore.entityPacker, err = storagepacker.NewStoragePacker(iStore.view, entitiesPackerLogger, "")
	if err != nil {
		return nil, errwrap.Wrapf("failed to create entity packer: {{err}}", err)
	}

	iStore.groupPacker, err = storagepacker.NewStoragePacker(iStore.view, groupsPackerLogger, groupBucketsPrefix)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create group packer: {{err}}", err)
	}

	iStore.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Paths:       iStore.paths(),
		Invalidate:  iStore.Invalidate,
	}

	err = iStore.Setup(ctx, config)
	if err != nil {
		return nil, err
	}

	return iStore, nil
}

func (i *IdentityStore) paths() []*framework.Path {
	return framework.PathAppend(
		entityPaths(i),
		aliasPaths(i),
		groupAliasPaths(i),
		groupPaths(i),
		lookupPaths(i),
		upgradePaths(i),
	)
}

// Invalidate is a callback wherein the backend is informed that the value at
// the given key is updated. In identity store's case, it would be the entity
// storage entries that get updated. The value needs to be read and MemDB needs
// to be updated accordingly.
func (i *IdentityStore) Invalidate(ctx context.Context, key string) {
	i.logger.Debug("invalidate notification received", "key", key)

	i.lock.Lock()
	defer i.lock.Unlock()

	switch {
	// Check if the key is a storage entry key for an entity bucket
	case strings.HasPrefix(key, storagepacker.StoragePackerBucketsPrefix):
		// Get the hash value of the storage bucket entry key
		bucketKeyHash := i.entityPacker.BucketKeyHashByKey(key)
		if len(bucketKeyHash) == 0 {
			i.logger.Error("failed to get the bucket entry key hash")
			return
		}

		// Create a MemDB transaction
		txn := i.db.Txn(true)
		defer txn.Abort()

		// Each entity object in MemDB holds the MD5 hash of the storage
		// entry key of the entity bucket. Fetch all the entities that
		// belong to this bucket using the hash value. Remove these entities
		// from MemDB along with all the aliases of each entity.
		entitiesFetched, err := i.MemDBEntitiesByBucketEntryKeyHashInTxn(txn, string(bucketKeyHash))
		if err != nil {
			i.logger.Error("failed to fetch entities using the bucket entry key hash", "bucket_entry_key_hash", bucketKeyHash)
			return
		}

		for _, entity := range entitiesFetched {
			// Delete all the aliases in the entity. This function will also remove
			// the corresponding alias indexes too.
			err = i.deleteAliasesInEntityInTxn(txn, entity, entity.Aliases)
			if err != nil {
				i.logger.Error("failed to delete aliases in entity", "entity_id", entity.ID, "error", err)
				return
			}

			// Delete the entity using the same transaction
			err = i.MemDBDeleteEntityByIDInTxn(txn, entity.ID)
			if err != nil {
				i.logger.Error("failed to delete entity from MemDB", "entity_id", entity.ID, "error", err)
				return
			}
		}

		// Get the storage bucket entry
		bucket, err := i.entityPacker.GetBucket(key)
		if err != nil {
			i.logger.Error("failed to refresh entities", "key", key, "error", err)
			return
		}

		// If the underlying entry is nil, it means that this invalidation
		// notification is for the deletion of the underlying storage entry. At
		// this point, since all the entities belonging to this bucket are
		// already removed, there is nothing else to be done. But, if the
		// storage entry is non-nil, its an indication of an update. In this
		// case, entities in the updated bucket needs to be reinserted into
		// MemDB.
		if bucket != nil {
			for _, item := range bucket.Items {
				entity, err := i.parseEntityFromBucketItem(ctx, item)
				if err != nil {
					i.logger.Error("failed to parse entity from bucket entry item", "error", err)
					return
				}

				// Only update MemDB and don't touch the storage
				err = i.upsertEntityInTxn(ctx, txn, entity, nil, false)
				if err != nil {
					i.logger.Error("failed to update entity in MemDB", "error", err)
					return
				}
			}
		}

		txn.Commit()
		return

	// Check if the key is a storage entry key for an group bucket
	case strings.HasPrefix(key, groupBucketsPrefix):
		// Get the hash value of the storage bucket entry key
		bucketKeyHash := i.groupPacker.BucketKeyHashByKey(key)
		if len(bucketKeyHash) == 0 {
			i.logger.Error("failed to get the bucket entry key hash")
			return
		}

		// Create a MemDB transaction
		txn := i.db.Txn(true)
		defer txn.Abort()

		groupsFetched, err := i.MemDBGroupsByBucketEntryKeyHashInTxn(txn, string(bucketKeyHash))
		if err != nil {
			i.logger.Error("failed to fetch groups using the bucket entry key hash", "bucket_entry_key_hash", bucketKeyHash)
			return
		}

		for _, group := range groupsFetched {
			// Delete the group using the same transaction
			err = i.MemDBDeleteGroupByIDInTxn(txn, group.ID)
			if err != nil {
				i.logger.Error("failed to delete group from MemDB", "group_id", group.ID, "error", err)
				return
			}
		}

		// Get the storage bucket entry
		bucket, err := i.groupPacker.GetBucket(key)
		if err != nil {
			i.logger.Error("failed to refresh group", "key", key, "error", err)
			return
		}

		if bucket != nil {
			for _, item := range bucket.Items {
				group, err := i.parseGroupFromBucketItem(item)
				if err != nil {
					i.logger.Error("failed to parse group from bucket entry item", "error", err)
					return
				}

				// Before updating the group, check if the group exists. If it
				// does, then delete the group alias from memdb, for the
				// invalidation would have sent an update.
				groupFetched, err := i.MemDBGroupByIDInTxn(txn, group.ID, true)
				if err != nil {
					i.logger.Error("failed to fetch group from MemDB", "error", err)
					return
				}

				// If the group has an alias remove it from memdb
				if groupFetched != nil && groupFetched.Alias != nil {
					err := i.MemDBDeleteAliasByIDInTxn(txn, groupFetched.Alias.ID, true)
					if err != nil {
						i.logger.Error("failed to delete old group alias from MemDB", "error", err)
						return
					}
				}

				// Only update MemDB and don't touch the storage
				err = i.UpsertGroupInTxn(txn, group, false)
				if err != nil {
					i.logger.Error("failed to update group in MemDB", "error", err)
					return
				}
			}
		}

		txn.Commit()
		return
	}
}

func (i *IdentityStore) parseEntityFromBucketItem(ctx context.Context, item *storagepacker.Item) (*identity.Entity, error) {
	if item == nil {
		return nil, fmt.Errorf("nil item")
	}

	persistNeeded := false

	var entity identity.Entity
	err := ptypes.UnmarshalAny(item.Message, &entity)
	if err != nil {
		// If we encounter an error, it would mean that the format of the
		// entity is an older one. Try decoding using the older format and if
		// successful, upgrage the storage with the newer format.
		var oldEntity identity.EntityStorageEntry
		oldEntityErr := ptypes.UnmarshalAny(item.Message, &oldEntity)
		if oldEntityErr != nil {
			return nil, errwrap.Wrapf("failed to decode entity from storage bucket item: {{err}}", err)
		}

		i.logger.Debug("upgrading the entity using patch introduced with vault 0.8.2.1", "entity_id", oldEntity.ID)

		// Successfully decoded entity using older format. Entity is stored
		// with older format. Upgrade it.
		entity.ID = oldEntity.ID
		entity.Name = oldEntity.Name
		entity.Metadata = oldEntity.Metadata
		entity.CreationTime = oldEntity.CreationTime
		entity.LastUpdateTime = oldEntity.LastUpdateTime
		entity.MergedEntityIDs = oldEntity.MergedEntityIDs
		entity.Policies = oldEntity.Policies
		entity.BucketKeyHash = oldEntity.BucketKeyHash
		entity.MFASecrets = oldEntity.MFASecrets
		// Copy each alias individually since the format of aliases were
		// also different
		for _, oldAlias := range oldEntity.Personas {
			var newAlias identity.Alias
			newAlias.ID = oldAlias.ID
			newAlias.Name = oldAlias.Name
			newAlias.CanonicalID = oldAlias.EntityID
			newAlias.MountType = oldAlias.MountType
			newAlias.MountAccessor = oldAlias.MountAccessor
			newAlias.MountPath = oldAlias.MountPath
			newAlias.Metadata = oldAlias.Metadata
			newAlias.CreationTime = oldAlias.CreationTime
			newAlias.LastUpdateTime = oldAlias.LastUpdateTime
			newAlias.MergedFromCanonicalIDs = oldAlias.MergedFromEntityIDs
			entity.Aliases = append(entity.Aliases, &newAlias)
		}

		persistNeeded = true
	}

	pN, err := parseExtraEntityFromBucket(ctx, i, &entity)
	if err != nil {
		return nil, err
	}
	if pN {
		persistNeeded = true
	}

	if persistNeeded && !i.core.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) {
		entityAsAny, err := ptypes.MarshalAny(&entity)
		if err != nil {
			return nil, err
		}

		item := &storagepacker.Item{
			ID:      entity.ID,
			Message: entityAsAny,
		}

		// Store the entity with new format
		err = i.entityPacker.PutItem(item)
		if err != nil {
			return nil, err
		}
	}

	if entity.NamespaceID == "" {
		entity.NamespaceID = namespace.RootNamespaceID
	}

	return &entity, nil
}

func (i *IdentityStore) parseGroupFromBucketItem(item *storagepacker.Item) (*identity.Group, error) {
	if item == nil {
		return nil, fmt.Errorf("nil item")
	}

	var group identity.Group
	err := ptypes.UnmarshalAny(item.Message, &group)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode group from storage bucket item: {{err}}", err)
	}

	if group.NamespaceID == "" {
		group.NamespaceID = namespace.RootNamespaceID
	}

	return &group, nil
}

// entityByAliasFactors fetches the entity based on factors of alias, i.e mount
// accessor and the alias name.
func (i *IdentityStore) entityByAliasFactors(mountAccessor, aliasName string, clone bool) (*identity.Entity, error) {
	if mountAccessor == "" {
		return nil, fmt.Errorf("missing mount accessor")
	}

	if aliasName == "" {
		return nil, fmt.Errorf("missing alias name")
	}

	txn := i.db.Txn(false)

	return i.entityByAliasFactorsInTxn(txn, mountAccessor, aliasName, clone)
}

// entityByAlaisFactorsInTxn fetches the entity based on factors of alias, i.e
// mount accessor and the alias name.
func (i *IdentityStore) entityByAliasFactorsInTxn(txn *memdb.Txn, mountAccessor, aliasName string, clone bool) (*identity.Entity, error) {
	if txn == nil {
		return nil, fmt.Errorf("nil txn")
	}

	if mountAccessor == "" {
		return nil, fmt.Errorf("missing mount accessor")
	}

	if aliasName == "" {
		return nil, fmt.Errorf("missing alias name")
	}

	alias, err := i.MemDBAliasByFactorsInTxn(txn, mountAccessor, aliasName, false, false)
	if err != nil {
		return nil, err
	}

	if alias == nil {
		return nil, nil
	}

	return i.MemDBEntityByAliasIDInTxn(txn, alias.ID, clone)
}

// CreateOrFetchEntity creates a new entity. This is used by core to
// associate each login attempt by an alias to a unified entity in Vault.
func (i *IdentityStore) CreateOrFetchEntity(ctx context.Context, alias *logical.Alias) (*identity.Entity, error) {
	var entity *identity.Entity
	var err error
	var update bool

	if alias == nil {
		return nil, fmt.Errorf("alias is nil")
	}

	if alias.Name == "" {
		return nil, fmt.Errorf("empty alias name")
	}

	mountValidationResp := i.core.router.validateMountByAccessor(alias.MountAccessor)
	if mountValidationResp == nil {
		return nil, fmt.Errorf("invalid mount accessor %q", alias.MountAccessor)
	}

	if mountValidationResp.MountLocal {
		return nil, fmt.Errorf("mount_accessor %q is of a local mount", alias.MountAccessor)
	}

	if mountValidationResp.MountType != alias.MountType {
		return nil, fmt.Errorf("mount accessor %q is not a mount of type %q", alias.MountAccessor, alias.MountType)
	}

	// Check if an entity already exists for the given alias
	entity, err = i.entityByAliasFactors(alias.MountAccessor, alias.Name, false)
	if err != nil {
		return nil, err
	}
	if entity != nil && changedAliasIndex(entity, alias) == -1 {
		return entity, nil
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	// Create a MemDB transaction to update both alias and entity
	txn := i.db.Txn(true)
	defer txn.Abort()

	// Check if an entity was created before acquiring the lock
	entity, err = i.entityByAliasFactorsInTxn(txn, alias.MountAccessor, alias.Name, true)
	if err != nil {
		return nil, err
	}
	if entity != nil {
		idx := changedAliasIndex(entity, alias)
		if idx == -1 {
			return entity, nil
		}
		a := entity.Aliases[idx]
		a.Metadata = alias.Metadata
		a.LastUpdateTime = ptypes.TimestampNow()

		update = true
	}

	if !update {
		entity = new(identity.Entity)
		err = i.sanitizeEntity(ctx, entity)
		if err != nil {
			return nil, err
		}

		// Create a new alias
		newAlias := &identity.Alias{
			CanonicalID:   entity.ID,
			Name:          alias.Name,
			MountAccessor: alias.MountAccessor,
			Metadata:      alias.Metadata,
			MountPath:     mountValidationResp.MountPath,
			MountType:     mountValidationResp.MountType,
		}

		err = i.sanitizeAlias(ctx, newAlias)
		if err != nil {
			return nil, err
		}

		i.logger.Debug("creating a new entity", "alias", newAlias)

		// Append the new alias to the new entity
		entity.Aliases = []*identity.Alias{
			newAlias,
		}
	}

	// Update MemDB and persist entity object
	err = i.upsertEntityInTxn(ctx, txn, entity, nil, true)
	if err != nil {
		return nil, err
	}

	txn.Commit()

	return entity, nil
}

// changedAliasIndex searches an entity for changed alias metadata.
//
// If a match is found, the changed alias's index is returned. If no alias
// names match or no metadata is different, -1 is returned.
func changedAliasIndex(entity *identity.Entity, alias *logical.Alias) int {
	for i, a := range entity.Aliases {
		if a.Name == alias.Name && !strutil.EqualStringMaps(a.Metadata, alias.Metadata) {
			return i
		}
	}

	return -1
}
