package vault

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/golang/protobuf/ptypes"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
)

const (
	groupBucketsPrefix        = "packer/group/buckets/"
	localAliasesBucketsPrefix = "packer/local-aliases/buckets/"
)

var (
	caseSensitivityKey           = "casesensitivity"
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
		view:          config.StorageView,
		logger:        logger,
		router:        core.router,
		redirectAddr:  core.redirectAddr,
		localNode:     core,
		namespacer:    core,
		metrics:       core.MetricSink(),
		totpPersister: core,
		groupUpdater:  core,
		tokenStorer:   core,
		entityCreator: core,
		mfaBackend:    core.loginMFABackend,
	}

	// Create a memdb instance, which by default, operates on lower cased
	// identity names
	err := iStore.resetDB(ctx)
	if err != nil {
		return nil, err
	}

	entitiesPackerLogger := iStore.logger.Named("storagepacker").Named("entities")
	core.AddLogger(entitiesPackerLogger)
	localAliasesPackerLogger := iStore.logger.Named("storagepacker").Named("local-aliases")
	core.AddLogger(localAliasesPackerLogger)
	groupsPackerLogger := iStore.logger.Named("storagepacker").Named("groups")
	core.AddLogger(groupsPackerLogger)

	iStore.entityPacker, err = storagepacker.NewStoragePacker(iStore.view, entitiesPackerLogger, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create entity packer: %w", err)
	}

	iStore.localAliasPacker, err = storagepacker.NewStoragePacker(iStore.view, localAliasesPackerLogger, localAliasesBucketsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create local alias packer: %w", err)
	}

	iStore.groupPacker, err = storagepacker.NewStoragePacker(iStore.view, groupsPackerLogger, groupBucketsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create group packer: %w", err)
	}

	iStore.Backend = &framework.Backend{
		BackendType:    logical.TypeLogical,
		Paths:          iStore.paths(),
		Invalidate:     iStore.Invalidate,
		InitializeFunc: iStore.initialize,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"oidc/.well-known/*",
				"oidc/provider/+/.well-known/*",
				"oidc/provider/+/token",
			},
			LocalStorage: []string{
				localAliasesBucketsPrefix,
			},
		},
		PeriodicFunc: func(ctx context.Context, req *logical.Request) error {
			iStore.oidcPeriodicFunc(ctx)

			return nil
		},
	}

	iStore.oidcCache = newOIDCCache(cache.NoExpiration, cache.NoExpiration)
	iStore.oidcAuthCodeCache = newOIDCCache(5*time.Minute, 5*time.Minute)

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
		oidcPaths(i),
		oidcProviderPaths(i),
		mfaPaths(i),
	)
}

func mfaPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "mfa/method" + genericOptionalUUIDRegex("method_id"),
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: `The unique identifier for this MFA method.`,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodReadGlobal,
					Summary:  "Read the current configuration for the given ID regardless of the MFA method type",
				},
			},
		},
		{
			Pattern: "mfa/method/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodListGlobal,
					Summary:  "List MFA method configurations for all MFA methods",
				},
			},
		},
		{
			Pattern: "mfa/method/totp" + genericOptionalUUIDRegex("method_id"),
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: `The unique identifier for this MFA method.`,
				},
				"max_validation_attempts": {
					Type:        framework.TypeInt,
					Description: `Max number of allowed validation attempts.`,
				},
				"issuer": {
					Type:        framework.TypeString,
					Description: `The name of the key's issuing organization.`,
				},
				"period": {
					Type:        framework.TypeDurationSecond,
					Default:     30,
					Description: `The length of time used to generate a counter for the TOTP token calculation.`,
				},
				"key_size": {
					Type:        framework.TypeInt,
					Default:     20,
					Description: "Determines the size in bytes of the generated key.",
				},
				"qr_size": {
					Type:        framework.TypeInt,
					Default:     200,
					Description: `The pixel size of the generated square QR code.`,
				},
				"algorithm": {
					Type:        framework.TypeString,
					Default:     "SHA1",
					Description: `The hashing algorithm used to generate the TOTP token. Options include SHA1, SHA256 and SHA512.`,
				},
				"digits": {
					Type:        framework.TypeInt,
					Default:     6,
					Description: `The number of digits in the generated TOTP token. This value can either be 6 or 8.`,
				},
				"skew": {
					Type:        framework.TypeInt,
					Default:     1,
					Description: `The number of delay periods that are allowed when validating a TOTP token. This value can either be 0 or 1.`,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodTOTPRead,
					Summary:  "Read the current configuration for the given MFA method",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodTOTPUpdate,
					Summary:  "Update or create a configuration for the given MFA method",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodDelete,
					Summary:  "Delete a configuration for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/totp/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodListTOTP,
					Summary:  "List MFA method configurations for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/totp/generate$",
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: `The unique identifier for this MFA method.`,
					Required:    true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleLoginMFAGenerateUpdate,
					Summary:  "Update or create TOTP secret for the given method ID on the given entity.",
				},
			},
		},
		{
			Pattern: "mfa/method/totp/admin-generate$",
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: `The unique identifier for this MFA method.`,
					Required:    true,
				},
				"entity_id": {
					Type:        framework.TypeString,
					Description: "Entity ID on which the generated secret needs to get stored.",
					Required:    true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleLoginMFAAdminGenerateUpdate,
					Summary:  "Update or create TOTP secret for the given method ID on the given entity.",
				},
			},
		},
		{
			Pattern: "mfa/method/totp/admin-destroy$",
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: "The unique identifier for this MFA method.",
					Required:    true,
				},
				"entity_id": {
					Type:        framework.TypeString,
					Description: "Identifier of the entity from which the MFA method secret needs to be removed.",
					Required:    true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleLoginMFAAdminDestroyUpdate,
					Summary:  "Destroys a TOTP secret for the given MFA method ID on the given entity",
				},
			},
		},
		{
			Pattern: "mfa/method/okta" + genericOptionalUUIDRegex("method_id"),
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: `The unique identifier for this MFA method.`,
				},
				"username_format": {
					Type:        framework.TypeString,
					Description: `A template string for mapping Identity names to MFA method names. Values to substitute should be placed in {{}}. For example, "{{entity.name}}@example.com". If blank, the Entity's name field will be used as-is.`,
				},
				"org_name": {
					Type:        framework.TypeString,
					Description: "Name of the organization to be used in the Okta API.",
				},
				"api_token": {
					Type:        framework.TypeString,
					Description: "Okta API key.",
				},
				"base_url": {
					Type:        framework.TypeString,
					Description: `The base domain to use for the Okta API. When not specified in the configuration, "okta.com" is used.`,
				},
				"primary_email": {
					Type:        framework.TypeBool,
					Description: `If true, the username will only match the primary email for the account. Defaults to false.`,
				},
				"production": {
					Type:        framework.TypeBool,
					Description: "(DEPRECATED) Use base_url instead.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodOKTARead,
					Summary:  "Read the current configuration for the given MFA method",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodOKTAUpdate,
					Summary:  "Update or create a configuration for the given MFA method",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodDelete,
					Summary:  "Delete a configuration for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/okta/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodListOkta,
					Summary:  "List MFA method configurations for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/duo" + genericOptionalUUIDRegex("method_id"),
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: `The unique identifier for this MFA method.`,
				},
				"username_format": {
					Type:        framework.TypeString,
					Description: `A template string for mapping Identity names to MFA method names. Values to subtitute should be placed in {{}}. For example, "{{alias.name}}@example.com". Currently-supported mappings: alias.name: The name returned by the mount configured via the mount_accessor parameter If blank, the Alias's name field will be used as-is. `,
				},
				"secret_key": {
					Type:        framework.TypeString,
					Description: "Secret key for Duo.",
				},
				"integration_key": {
					Type:        framework.TypeString,
					Description: "Integration key for Duo.",
				},
				"api_hostname": {
					Type:        framework.TypeString,
					Description: "API host name for Duo.",
				},
				"push_info": {
					Type:        framework.TypeString,
					Description: "Push information for Duo.",
				},
				"use_passcode": {
					Type:        framework.TypeBool,
					Description: `If true, the user is reminded to use the passcode upon MFA validation. This option does not enforce using the passcode. Defaults to false.`,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodDuoRead,
					Summary:  "Read the current configuration for the given MFA method",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodDuoUpdate,
					Summary:  "Update or create a configuration for the given MFA method",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodDelete,
					Summary:  "Delete a configuration for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/duo/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodListDuo,
					Summary:  "List MFA method configurations for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/pingid" + genericOptionalUUIDRegex("method_id"),
			Fields: map[string]*framework.FieldSchema{
				"method_id": {
					Type:        framework.TypeString,
					Description: `The unique identifier for this MFA method.`,
				},
				"username_format": {
					Type:        framework.TypeString,
					Description: `A template string for mapping Identity names to MFA method names. Values to subtitute should be placed in {{}}. For example, "{{alias.name}}@example.com". Currently-supported mappings: alias.name: The name returned by the mount configured via the mount_accessor parameter If blank, the Alias's name field will be used as-is. `,
				},
				"settings_file_base64": {
					Type:        framework.TypeString,
					Description: "The settings file provided by Ping, Base64-encoded. This must be a settings file suitable for third-party clients, not the PingID SDK or PingFederate.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodPingIDRead,
					Summary:  "Read the current configuration for the given MFA method",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodPingIDUpdate,
					Summary:  "Update or create a configuration for the given MFA method",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodDelete,
					Summary:  "Delete a configuration for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/pingid/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodListPingID,
					Summary:  "List MFA method configurations for the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/login-enforcement/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name for this login enforcement configuration",
					Required:    true,
				},
				"mfa_method_ids": {
					Type:        framework.TypeStringSlice,
					Description: "Array of Method IDs that determine what methods will be enforced",
					Required:    true,
				},
				"auth_method_accessors": {
					Type:        framework.TypeStringSlice,
					Description: "Array of auth mount accessor IDs",
				},
				"auth_method_types": {
					Type:        framework.TypeStringSlice,
					Description: "Array of auth mount types",
				},
				"identity_group_ids": {
					Type:        framework.TypeStringSlice,
					Description: "Array of identity group IDs",
				},
				"identity_entity_ids": {
					Type:        framework.TypeStringSlice,
					Description: "Array of identity entity IDs",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.handleMFALoginEnforcementRead,
					Summary:  "Read the current login enforcement",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleMFALoginEnforcementUpdate,
					Summary:  "Create or update a login enforcement",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.handleMFALoginEnforcementDelete,
					Summary:  "Delete a login enforcement",
				},
			},
		},
		{
			Pattern: "mfa/login-enforcement/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.handleMFALoginEnforcementList,
					Summary:  "List login enforcements",
				},
			},
		},
	}
}

func (i *IdentityStore) initialize(ctx context.Context, req *logical.InitializationRequest) error {
	// Only primary should write the status
	if i.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary | consts.ReplicationPerformanceStandby | consts.ReplicationDRSecondary) {
		return nil
	}

	if err := i.storeOIDCDefaultResources(ctx, req.Storage); err != nil {
		return err
	}

	entry, err := logical.StorageEntryJSON(caseSensitivityKey, &casesensitivity{
		DisableLowerCasedNames: i.disableLowerCasedNames,
	})
	if err != nil {
		return err
	}

	return i.view.Put(ctx, entry)
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
	case key == caseSensitivityKey:
		entry, err := i.view.Get(ctx, caseSensitivityKey)
		if err != nil {
			i.logger.Error("failed to read case sensitivity setting during invalidation", "error", err)
			return
		}
		if entry == nil {
			return
		}

		var setting casesensitivity
		if err := entry.DecodeJSON(&setting); err != nil {
			i.logger.Error("failed to decode case sensitivity setting during invalidation", "error", err)
			return
		}

		// Fast return if the setting is the same
		if i.disableLowerCasedNames == setting.DisableLowerCasedNames {
			return
		}

		// If the setting is different, reset memdb and reload all the artifacts
		i.disableLowerCasedNames = setting.DisableLowerCasedNames
		if err := i.resetDB(ctx); err != nil {
			i.logger.Error("failed to reset memdb during invalidation", "error", err)
			return
		}
		if err := i.loadEntities(ctx); err != nil {
			i.logger.Error("failed to load entities during invalidation", "error", err)
			return
		}
		if err := i.loadGroups(ctx); err != nil {
			i.logger.Error("failed to load groups during invalidation", "error", err)
			return
		}
		if err := i.loadOIDCClients(ctx); err != nil {
			i.logger.Error("failed to load OIDC clients during invalidation", "error", err)
			return
		}
	// Check if the key is a storage entry key for an entity bucket
	case strings.HasPrefix(key, storagepacker.StoragePackerBucketsPrefix):
		// Create a MemDB transaction
		txn := i.db.Txn(true)
		defer txn.Abort()

		// Each entity object in MemDB holds the MD5 hash of the storage
		// entry key of the entity bucket. Fetch all the entities that
		// belong to this bucket using the hash value. Remove these entities
		// from MemDB along with all the aliases of each entity.
		entitiesFetched, err := i.MemDBEntitiesByBucketKeyInTxn(txn, key)
		if err != nil {
			i.logger.Error("failed to fetch entities using the bucket key", "key", key)
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
		bucket, err := i.entityPacker.GetBucket(ctx, key)
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
		var entityIDs []string
		if bucket != nil {
			entityIDs = make([]string, 0, len(bucket.Items))
			for _, item := range bucket.Items {
				entity, err := i.parseEntityFromBucketItem(ctx, item)
				if err != nil {
					i.logger.Error("failed to parse entity from bucket entry item", "error", err)
					return
				}

				localAliases, err := i.parseLocalAliases(entity.ID)
				if err != nil {
					i.logger.Error("failed to load local aliases from storage", "error", err)
					return
				}
				if localAliases != nil {
					for _, alias := range localAliases.Aliases {
						entity.UpsertAlias(alias)
					}
				}

				// Only update MemDB and don't touch the storage
				err = i.upsertEntityInTxn(ctx, txn, entity, nil, false)
				if err != nil {
					i.logger.Error("failed to update entity in MemDB", "error", err)
					return
				}

				// If we are a secondary, the entity created by the secondary
				// via the CreateEntity RPC would have been cached. Now that the
				// invalidation of the same has hit, there is no need of the
				// cache. Clearing the cache. Writing to storage can't be
				// performed by perf standbys. So only doing this in the active
				// node of the secondary.
				if i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) && i.localNode.HAState() != consts.PerfStandby {
					if err := i.localAliasPacker.DeleteItem(ctx, entity.ID+tmpSuffix); err != nil {
						i.logger.Error("failed to clear local alias entity cache", "error", err, "entity_id", entity.ID)
						return
					}
				}

				entityIDs = append(entityIDs, entity.ID)
			}
		}

		// entitiesFetched are the entities before invalidation. entityIDs
		// represent entities that are valid after invalidation. Clear the
		// storage entries of local aliases for those entities that are
		// indicated deleted by this invalidation.
		if i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) && i.localNode.HAState() != consts.PerfStandby {
			for _, entity := range entitiesFetched {
				if !strutil.StrListContains(entityIDs, entity.ID) {
					if err := i.localAliasPacker.DeleteItem(ctx, entity.ID); err != nil {
						i.logger.Error("failed to clear local alias for entity", "error", err, "entity_id", entity.ID)
						return
					}
				}
			}
		}

		txn.Commit()
		return

	// Check if the key is a storage entry key for an group bucket
	// For those entities that are deleted, clear up the local alias entries
	case strings.HasPrefix(key, groupBucketsPrefix):
		// Create a MemDB transaction
		txn := i.db.Txn(true)
		defer txn.Abort()

		groupsFetched, err := i.MemDBGroupsByBucketKeyInTxn(txn, key)
		if err != nil {
			i.logger.Error("failed to fetch groups using the bucket key", "key", key)
			return
		}

		for _, group := range groupsFetched {
			// Delete the group using the same transaction
			err = i.MemDBDeleteGroupByIDInTxn(txn, group.ID)
			if err != nil {
				i.logger.Error("failed to delete group from MemDB", "group_id", group.ID, "error", err)
				return
			}

			if group.Alias != nil {
				err := i.MemDBDeleteAliasByIDInTxn(txn, group.Alias.ID, true)
				if err != nil {
					i.logger.Error("failed to delete group alias from MemDB", "error", err)
					return
				}
			}
		}

		// Get the storage bucket entry
		bucket, err := i.groupPacker.GetBucket(ctx, key)
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
				err = i.UpsertGroupInTxn(ctx, txn, group, false)
				if err != nil {
					i.logger.Error("failed to update group in MemDB", "error", err)
					return
				}
			}
		}

		txn.Commit()
		return

	case strings.HasPrefix(key, oidcTokensPrefix):
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			i.logger.Error("error retrieving namespace", "error", err)
			return
		}

		// Wipe the cache for the requested namespace. This will also clear
		// the shared namespace as well.
		if err := i.oidcCache.Flush(ns); err != nil {
			i.logger.Error("error flushing oidc cache", "error", err)
		}
	case strings.HasPrefix(key, clientPath):
		name := strings.TrimPrefix(key, clientPath)

		// Invalidate the cached client in memdb
		if err := i.memDBDeleteClientByName(ctx, name); err != nil {
			i.logger.Error("error invalidating client", "error", err, "key", key)
			return
		}
	case strings.HasPrefix(key, localAliasesBucketsPrefix):
		//
		// This invalidation only happens on perf standbys
		//

		txn := i.db.Txn(true)
		defer txn.Abort()

		// Find all the local aliases belonging to this bucket and remove it
		// both from aliases table and entities table. We will add the local
		// aliases back by parsing the storage key. This way the deletion
		// invalidation gets handled.
		aliases, err := i.MemDBLocalAliasesByBucketKeyInTxn(txn, key)
		if err != nil {
			i.logger.Error("failed to fetch entities using the bucket key", "key", key)
			return
		}

		for _, alias := range aliases {
			entity, err := i.MemDBEntityByIDInTxn(txn, alias.CanonicalID, true)
			if err != nil {
				i.logger.Error("failed to fetch entity during local alias invalidation", "entity_id", alias.CanonicalID, "error", err)
				return
			}
			if entity == nil {
				i.logger.Error("failed to fetch entity during local alias invalidation, missing entity", "entity_id", alias.CanonicalID, "error", err)
				continue
			}

			// Delete local aliases from the entity.
			err = i.deleteAliasesInEntityInTxn(txn, entity, []*identity.Alias{alias})
			if err != nil {
				i.logger.Error("failed to delete aliases in entity", "entity_id", entity.ID, "error", err)
				return
			}

			// Update the entity with removed alias.
			if err := i.MemDBUpsertEntityInTxn(txn, entity); err != nil {
				i.logger.Error("failed to delete entity from MemDB", "entity_id", entity.ID, "error", err)
				return
			}
		}

		// Now read the invalidated storage key
		bucket, err := i.localAliasPacker.GetBucket(ctx, key)
		if err != nil {
			i.logger.Error("failed to refresh local aliases", "key", key, "error", err)
			return
		}
		if bucket != nil {
			for _, item := range bucket.Items {
				if strings.HasSuffix(item.ID, tmpSuffix) {
					continue
				}

				var localAliases identity.LocalAliases
				err = ptypes.UnmarshalAny(item.Message, &localAliases)
				if err != nil {
					i.logger.Error("failed to parse local aliases during invalidation", "error", err)
					return
				}
				for _, alias := range localAliases.Aliases {
					// Add to the aliases table
					if err := i.MemDBUpsertAliasInTxn(txn, alias, false); err != nil {
						i.logger.Error("failed to insert local alias to memdb during invalidation", "error", err)
						return
					}

					// Fetch the associated entity and add the alias to that too.
					entity, err := i.MemDBEntityByIDInTxn(txn, alias.CanonicalID, false)
					if err != nil {
						i.logger.Error("failed to fetch entity during local alias invalidation", "error", err)
						return
					}
					if entity == nil {
						cachedEntityItem, err := i.localAliasPacker.GetItem(alias.CanonicalID + tmpSuffix)
						if err != nil {
							i.logger.Error("failed to fetch cached entity", "key", key, "error", err)
							return
						}
						if cachedEntityItem != nil {
							entity, err = i.parseCachedEntity(cachedEntityItem)
							if err != nil {
								i.logger.Error("failed to parse cached entity", "key", key, "error", err)
								return
							}
						}
					}
					if entity == nil {
						i.logger.Error("received local alias invalidation for an invalid entity", "item.ID", item.ID)
						return
					}
					entity.UpsertAlias(alias)

					// Update the entities table
					if err := i.MemDBUpsertEntityInTxn(txn, entity); err != nil {
						i.logger.Error("failed to upsert entity during local alias invalidation", "error", err)
						return
					}
				}
			}
		}
		txn.Commit()
		return
	}
}

func (i *IdentityStore) parseLocalAliases(entityID string) (*identity.LocalAliases, error) {
	item, err := i.localAliasPacker.GetItem(entityID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	var localAliases identity.LocalAliases
	err = ptypes.UnmarshalAny(item.Message, &localAliases)
	if err != nil {
		return nil, err
	}

	return &localAliases, nil
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
			return nil, fmt.Errorf("failed to decode entity from storage bucket item: %w", err)
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
		entity.BucketKey = oldEntity.BucketKeyHash
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
			entity.UpsertAlias(&newAlias)
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

	if persistNeeded && !i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) {
		entityAsAny, err := ptypes.MarshalAny(&entity)
		if err != nil {
			return nil, err
		}

		item := &storagepacker.Item{
			ID:      entity.ID,
			Message: entityAsAny,
		}

		// Store the entity with new format
		err = i.entityPacker.PutItem(ctx, item)
		if err != nil {
			return nil, err
		}
	}

	if entity.NamespaceID == "" {
		entity.NamespaceID = namespace.RootNamespaceID
	}

	return &entity, nil
}

func (i *IdentityStore) parseCachedEntity(item *storagepacker.Item) (*identity.Entity, error) {
	if item == nil {
		return nil, fmt.Errorf("nil item")
	}

	var entity identity.Entity
	err := ptypes.UnmarshalAny(item.Message, &entity)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cached entity from storage bucket item: %w", err)
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
		return nil, fmt.Errorf("failed to decode group from storage bucket item: %w", err)
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

// entityByAliasFactorsInTxn fetches the entity based on factors of alias, i.e
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

// CreateEntity creates a new entity.
func (i *IdentityStore) CreateEntity(ctx context.Context) (*identity.Entity, error) {
	defer metrics.MeasureSince([]string{"identity", "create_entity"}, time.Now())

	entity := new(identity.Entity)
	err := i.sanitizeEntity(ctx, entity)
	if err != nil {
		return nil, err
	}
	if err := i.upsertEntity(ctx, entity, nil, true); err != nil {
		return nil, err
	}

	// Emit a metric for the new entity
	ns, err := i.namespacer.NamespaceByID(ctx, entity.NamespaceID)
	var nsLabel metrics.Label
	if err != nil {
		nsLabel = metrics.Label{"namespace", "unknown"}
	} else {
		nsLabel = metricsutil.NamespaceLabel(ns)
	}
	i.metrics.IncrCounterWithLabels(
		[]string{"identity", "entity", "creation"},
		1,
		[]metrics.Label{
			nsLabel,
		})

	return entity, nil
}

// CreateOrFetchEntity creates a new entity. This is used by core to
// associate each login attempt by an alias to a unified entity in Vault.
func (i *IdentityStore) CreateOrFetchEntity(ctx context.Context, alias *logical.Alias) (*identity.Entity, error) {
	defer metrics.MeasureSince([]string{"identity", "create_or_fetch_entity"}, time.Now())

	var entity *identity.Entity
	var err error
	var update bool

	if alias == nil {
		return nil, fmt.Errorf("alias is nil")
	}

	if alias.Name == "" {
		return nil, fmt.Errorf("empty alias name")
	}

	mountValidationResp := i.router.ValidateMountByAccessor(alias.MountAccessor)
	if mountValidationResp == nil {
		return nil, fmt.Errorf("invalid mount accessor %q", alias.MountAccessor)
	}

	if mountValidationResp.MountType != alias.MountType {
		return nil, fmt.Errorf("mount accessor %q is not a mount of type %q", alias.MountAccessor, alias.MountType)
	}

	// Check if an entity already exists for the given alias
	entity, err = i.entityByAliasFactors(alias.MountAccessor, alias.Name, true)
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
			Local:         alias.Local,
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

		// Emit a metric for the new entity
		ns, err := i.namespacer.NamespaceByID(ctx, entity.NamespaceID)
		var nsLabel metrics.Label
		if err != nil {
			nsLabel = metrics.Label{"namespace", "unknown"}
		} else {
			nsLabel = metricsutil.NamespaceLabel(ns)
		}
		i.metrics.IncrCounterWithLabels(
			[]string{"identity", "entity", "creation"},
			1,
			[]metrics.Label{
				nsLabel,
				{"auth_method", newAlias.MountType},
				{"mount_point", newAlias.MountPath},
			})
	}

	// Update MemDB and persist entity object
	err = i.upsertEntityInTxn(ctx, txn, entity, nil, true)
	if err != nil {
		return nil, err
	}

	txn.Commit()
	return entity.Clone()
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
