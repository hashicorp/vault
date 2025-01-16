// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"reflect"
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
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (i *IdentityStore) resetDB() error {
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
		mountLister:   core,
		mfaBackend:    core.loginMFABackend,
		aliasLocks:    locksutil.CreateLocks(),
	}

	// Create a memdb instance, which by default, operates on lower cased
	// identity names
	err := iStore.resetDB()
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
				"oidc/+/.well-known/*",
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
		RunningVersion: versions.DefaultBuiltinVersion,
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
		mfaCommonPaths(i),
		mfaTOTPPaths(i),
		mfaTOTPExtraPaths(i),
		mfaOktaPaths(i),
		mfaDuoPaths(i),
		mfaPingIDPaths(i),
		mfaLoginEnforcementPaths(i),
	)
}

func mfaCommonPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "mfa/method/" + uuidRegex("method_id"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationSuffix: "method",
			},
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
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationVerb:   "list",
				OperationSuffix: "methods",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.handleMFAMethodListGlobal,
					Summary:  "List MFA method configurations for all MFA methods",
				},
			},
		},
	}
}

func makeMFAMethodPaths(
	methodType string,
	methodTypeForOpenAPIOperationID string,
	methodFields map[string]*framework.FieldSchema,
	i *IdentityStore,
) []*framework.Path {
	methodFieldsIncludingMethodID := make(map[string]*framework.FieldSchema)
	for k, v := range methodFields {
		methodFieldsIncludingMethodID[k] = v
	}
	methodFieldsIncludingMethodID["method_id"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: `The unique identifier for this MFA method.`,
	}

	return []*framework.Path{
		{
			Pattern: "mfa/method/" + methodType + "/" + uuidRegex("method_id"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationSuffix: methodTypeForOpenAPIOperationID + "-method",
			},
			Fields: methodFieldsIncludingMethodID,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
						return i.handleMFAMethodReadCommon(ctx, req, d, methodType)
					},
					Summary: "Read the current configuration for the given MFA method",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
						return i.handleMFAMethodWriteCommon(ctx, req, d, methodType)
					},
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
					Summary: "Update the configuration for the given MFA method",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
						return i.handleMFAMethodDeleteCommon(ctx, req, d, methodType)
					},
					Summary: "Delete the given MFA method",
				},
			},
		},
		{
			Pattern: "mfa/method/" + methodType + "/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
			},
			Fields: methodFields,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
						return i.handleMFAMethodListCommon(ctx, req, d, methodType)
					},
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: methodTypeForOpenAPIOperationID + "-methods",
					},
					Summary: "List MFA method configurations for the given MFA method",
				},
				// Conceptually, it would make more sense to treat this as a CreateOperation, but we have to leave it
				// as an UpdateOperation, because the API was originally released that way, and we don't want to change
				// the meaning of ACL policies that users have written.
				logical.UpdateOperation: &framework.PathOperation{
					Callback: func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
						return i.handleMFAMethodWriteCommon(ctx, req, d, methodType)
					},
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "create",
						OperationSuffix: methodTypeForOpenAPIOperationID + "-method",
					},
					Summary: "Create the given MFA method",
				},
			},
		},
	}
}

func mfaTOTPPaths(i *IdentityStore) []*framework.Path {
	return makeMFAMethodPaths(
		mfaMethodTypeTOTP,
		mfaMethodTypeTOTP,
		map[string]*framework.FieldSchema{
			"method_name": {
				Type:        framework.TypeString,
				Description: `The unique name identifier for this MFA method.`,
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
		i,
	)
}

func mfaTOTPExtraPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "mfa/method/totp/generate$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationVerb:   "generate",
				OperationSuffix: "totp-secret",
			},
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
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationVerb:   "admin-generate",
				OperationSuffix: "totp-secret",
			},
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
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationVerb:   "admin-destroy",
				OperationSuffix: "totp-secret",
			},
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
	}
}

func mfaOktaPaths(i *IdentityStore) []*framework.Path {
	return makeMFAMethodPaths(
		mfaMethodTypeOkta,
		mfaMethodTypeOkta,
		map[string]*framework.FieldSchema{
			"method_name": {
				Type:        framework.TypeString,
				Description: `The unique name identifier for this MFA method.`,
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
		i,
	)
}

func mfaDuoPaths(i *IdentityStore) []*framework.Path {
	return makeMFAMethodPaths(
		mfaMethodTypeDuo,
		mfaMethodTypeDuo,
		map[string]*framework.FieldSchema{
			"method_name": {
				Type:        framework.TypeString,
				Description: `The unique name identifier for this MFA method.`,
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
		i,
	)
}

func mfaPingIDPaths(i *IdentityStore) []*framework.Path {
	return makeMFAMethodPaths(
		mfaMethodTypePingID,
		// This overridden name helps code generation using the OpenAPI spec choose better method names, that avoid
		// treating "Pingid" as a single word:
		"ping-id",
		map[string]*framework.FieldSchema{
			"method_name": {
				Type:        framework.TypeString,
				Description: `The unique name identifier for this MFA method.`,
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
		i,
	)
}

func mfaLoginEnforcementPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "mfa/login-enforcement/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationSuffix: "login-enforcement",
			},
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
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationSuffix: "login-enforcements",
			},
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
		i.logger.Error("failed to write OIDC default resources to storage", "error", err)
		return err
	}

	// if the storage entry for caseSensitivityKey exists, remove it
	storageEntry, err := i.view.Get(ctx, caseSensitivityKey)
	if err != nil {
		i.logger.Error("could not get storage entry for case sensitivity key", "error", err)
		return nil
	}

	if storageEntry != nil {
		var setting casesensitivity
		err := storageEntry.DecodeJSON(&setting)
		switch err {
		case nil:
			i.logger.Debug("removing storage entry for case sensitivity key", "value", setting.DisableLowerCasedNames)
		default:
			i.logger.Error("failed to decode case sensitivity key, removing its storage entry anyway", "error", err)
		}

		err = i.view.Delete(ctx, caseSensitivityKey)
		if err != nil {
			i.logger.Error("could not delete storage entry for case sensitivity key", "error", err)
			return nil
		}
	}

	return nil
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
	case strings.HasPrefix(key, storagepacker.StoragePackerBucketsPrefix):
		// key is for a entity bucket in storage.
		i.invalidateEntityBucket(ctx, key)
	case strings.HasPrefix(key, groupBucketsPrefix):
		// key is for a group bucket in storage.
		i.invalidateGroupBucket(ctx, key)
	case strings.HasPrefix(key, oidcTokensPrefix):
		// key is for oidc tokens in storage.
		i.invalidateOIDCToken(ctx)
	case strings.HasPrefix(key, clientPath):
		// key is for a client in storage.
		i.invalidateClientPath(ctx, key)
	case strings.HasPrefix(key, localAliasesBucketsPrefix):
		// key is for a local alias bucket in storage.
		i.invalidateLocalAliasesBucket(ctx, key)
	}
}

func (i *IdentityStore) invalidateEntityBucket(ctx context.Context, key string) {
	txn := i.db.Txn(true)
	defer txn.Abort()

	// The handling of entities has the added quirk of dealing with a temporary
	// copy of the entity written in storage on the active node of performance
	// secondary clusters. These temporary entity entries in storage must be
	// removed once the actual entity appears in the storage bucket (as
	// replicated from the primary cluster).
	//
	// This function retrieves all entities from MemDB that have a corresponding
	// storage key that matches the provided key to invalidate. This is the set
	// of entities that need to be updated, removed, or left alone in MemDB.
	//
	// The logic iterates over every entity stored in the invalidated storage
	// bucket. For each entity read from the storage bucket, the set of entities
	// read from MemDB is searched for the same entity. If it can't be found,
	// it means that it needs to be inserted into MemDB. On the other hand, if
	// the entity is found, it the storage bucket entity is compared to the
	// MemDB entity. If they do not match, then the storage entity state needs
	// to be used to update the MemDB entity; if they did match, then it means
	// that the MemDB entity can be left alone. As each MemDB entity is
	// processed in the loop, it is removed from the set of MemDB entities.
	//
	// Once all entities from the storage bucket have been compared to those
	// retrieved from MemDB, the remaining entities from the set retrieved from
	// MemDB are those that have been deleted from storage and must be removed
	// from MemDB (because as MemDB entities that matches a storage bucket
	// entity were processed, they were removed from the set).
	memDBEntities, err := i.MemDBEntitiesByBucketKeyInTxn(txn, key)
	if err != nil {
		i.logger.Error("failed to fetch entities using the bucket key", "key", key)
		return
	}

	bucket, err := i.entityPacker.GetBucket(ctx, key)
	if err != nil {
		i.logger.Error("failed to refresh entities", "key", key, "error", err)
		return
	}

	if bucket != nil {
		// The storage entry for the entity bucket exists, so we need to compare
		// the entities in that bucket with those in MemDB and only update those
		// that are different. The entities in the bucket storage entry are the
		// source of truth.

		// Iterate over each entity item from the bucket
		for _, item := range bucket.Items {
			bucketEntity, err := i.parseEntityFromBucketItem(ctx, item)
			if err != nil {
				i.logger.Error("failed to parse entity from bucket entry item", "error", err)
				return
			}

			localAliases, err := i.parseLocalAliases(bucketEntity.ID)
			if err != nil {
				i.logger.Error("failed to load local aliases from storage", "error", err)
				return
			}

			if localAliases != nil {
				for _, alias := range localAliases.Aliases {
					bucketEntity.UpsertAlias(alias)
				}
			}

			var memDBEntity *identity.Entity
			for i, entity := range memDBEntities {
				if entity.ID == bucketEntity.ID {
					memDBEntity = entity

					// Remove this processed entity from the slice, so that
					// all tht will be left are unprocessed entities.
					copy(memDBEntities[i:], memDBEntities[i+1:])
					memDBEntities = memDBEntities[:len(memDBEntities)-1]
					break
				}
			}

			// We've considered the use of github.com/google/go-cmp here,
			// but opted for sticking with reflect.DeepEqual because go-cmp
			// is intended for testing and is able to panic in some
			// situations.
			if memDBEntity != nil && reflect.DeepEqual(memDBEntity, bucketEntity) {
				// No changes on this entity, move on to the next one.
				continue
			}

			// If the entity exists in MemDB it must differ from the entity in
			// the storage bucket because of above test. Blindly delete the
			// current aliases associated with the MemDB entity. The correct set
			// of aliases will be created in MemDB by the upsertEntityInTxn
			// function. We need to do this because the upsertEntityInTxn
			// function does not delete those aliases, it only creates missing
			// ones.
			if memDBEntity != nil {
				if err := i.deleteAliasesInEntityInTxn(txn, memDBEntity, memDBEntity.Aliases); err != nil {
					i.logger.Error("failed to remove entity aliases from changed entity", "entity_id", memDBEntity.ID, "error", err)
					return
				}

				if err := i.MemDBDeleteEntityByIDInTxn(txn, memDBEntity.ID); err != nil {
					i.logger.Error("failed to delete changed entity", "entity_id", memDBEntity.ID, "error", err)
					return
				}
			}

			err = i.upsertEntityInTxn(ctx, txn, bucketEntity, nil, false)
			if err != nil {
				i.logger.Error("failed to update entity in MemDB", "entity_id", bucketEntity.ID, "error", err)
				return
			}

			// If this is a performance secondary, the entity created on
			// this node would have been cached in a local cache based on
			// the result of the CreateEntity RPC call to the primary
			// cluster. Since this invalidation is signaling that the
			// entity is now in the primary cluster's storage, the locally
			// cached entry can be removed.
			if i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) && i.localNode.HAState() == consts.Active {
				if err := i.localAliasPacker.DeleteItem(ctx, bucketEntity.ID+tmpSuffix); err != nil {
					i.logger.Error("failed to clear local alias entity cache", "error", err, "entity_id", bucketEntity.ID)
					return
				}
			}

		}
	}

	// Any entities that are still in the memDBEntities slice are ones that do
	// not exist in the bucket storage entry. These entities have to be removed
	// from MemDB.
	for _, memDBEntity := range memDBEntities {
		err = i.deleteAliasesInEntityInTxn(txn, memDBEntity, memDBEntity.Aliases)
		if err != nil {
			i.logger.Error("failed to delete aliases in entity", "entity_id", memDBEntity.ID, "error", err)
			return
		}

		err = i.MemDBDeleteEntityByIDInTxn(txn, memDBEntity.ID)
		if err != nil {
			i.logger.Error("failed to delete entity from MemDB", "entity_id", memDBEntity.ID, "error", err)
			return
		}

		// In addition, if this is an active node of a performance secondary
		// cluster, remove the local alias storage entry for this deleted entity.
		if i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) && i.localNode.HAState() == consts.Active {
			if err := i.localAliasPacker.DeleteItem(ctx, memDBEntity.ID); err != nil {
				i.logger.Error("failed to clear local alias for entity", "error", err, "entity_id", memDBEntity.ID)
				return
			}
		}
	}

	txn.Commit()
}

func (i *IdentityStore) invalidateGroupBucket(ctx context.Context, key string) {
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
}

// invalidateOIDCToken is called by the Invalidate function to handle the
// invalidation of an OIDC token storage entry.
func (i *IdentityStore) invalidateOIDCToken(ctx context.Context) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		i.logger.Error("error retrieving namespace", "error", err)
		return
	}

	// Wipe the cache for the requested namespace. This will also clear
	// the shared namespace as well.
	if err := i.oidcCache.Flush(ns); err != nil {
		i.logger.Error("error flushing oidc cache", "error", err)
		return
	}
}

// invalidateClientPath is called by the Invalidate function to handle the
// invalidation of a client path storage entry.
func (i *IdentityStore) invalidateClientPath(ctx context.Context, key string) {
	name := strings.TrimPrefix(key, clientPath)

	// Invalidate the cached client in memdb
	if err := i.memDBDeleteClientByName(ctx, name); err != nil {
		i.logger.Error("error invalidating client", "error", err, "key", key)
		return
	}
}

// invalidateLocalAliasBucket is called by the Invalidate function to handle the
// invalidation of a local alias bucket storage entry.
func (i *IdentityStore) invalidateLocalAliasesBucket(ctx context.Context, key string) {
	// This invalidation only happens on performance standby servers

	// Create a MemDB transaction and abort it once this function returns
	txn := i.db.Txn(true)
	defer txn.Abort()

	// Local aliases have the added complexity of being associated with
	// entities. Whenever a local alias is updated or inserted into MemDB, its
	// associated MemDB-stored entity must also be updated.
	//
	// This function retrieves all local aliases that have a corresponding
	// storage key that matches the provided key to invalidate. This is the
	// set of local aliases that need to be updated, removed, or left
	// alone in MemDB. Each of these operations is done as its own MemDB
	// operation, but the corresponding changes that need to be made to the
	// associated entities can be batched together to cut down on the number of
	// MemDB operations.
	//
	// The logic iterates over every local alias stored at the invalidated key.
	// For each local alias read from the storage entry, the set of local
	// aliases read from MemDB is searched for the same local alias. If it can't
	// be found, it means that it needs to be inserted into MemDB. However, if
	// it's found, it must be compared with the local alias from the storage. If
	// they don't match, it means that the local alias in MemDB needs to be
	// updated. If they did match, it means that this particular local alias did
	// not change in storage, so nothing further needs to be done. Each local
	// alias processed in this loop is removed from the set of retrieved local
	// aliases. The local alias is also added to the map tracking local aliases
	// that need to be upserted in their associated entities in MemDB.
	//
	// Once the code is done iterating over all of the local aliases from
	// storage, any local aliases still in the set retrieved from MemDB
	// corresponds to a local alias that is no longer in storage and must be
	// removed from MemDB. These local aliases are added to the map tracking
	// local aliases to remove from their entities in MemDB. The actual removal
	// of the local aliases themselves is done as part of the tidying up of the
	// associated entities, described below.
	//
	// In order to batch the changes to the associated entities, a map of entity
	// to local aliases (slice of local alias) is built up in the loop that
	// iterates over the local aliases from storage. Similarly, the code that
	// detects which local aliases to remove from MemDB also builds a separate
	// map of entity to local aliases (slice of local alias). Each element in
	// the map of local aliases to update in their entity is processed as
	// follows: the mapped slice of local aliases is iterated over and each
	// local alias is upserted into the entity and then the entity itself is
	// upserted. Then, each element in the map of local aliases to remove from
	// their entity is processed as follows: the

	// Get all cached local aliases to compare with invalidated bucket
	memDBLocalAliases, err := i.MemDBLocalAliasesByBucketKeyInTxn(txn, key)
	if err != nil {
		i.logger.Error("failed to fetch local aliases using the bucket key", "key", key, "error", err)
		return
	}

	// Get local aliases from the invalidated bucket
	bucket, err := i.localAliasPacker.GetBucket(ctx, key)
	if err != nil {
		i.logger.Error("failed to refresh local aliases", "key", key, "error", err)
		return
	}

	// This map tracks the set of local aliases that need to be updated in each
	// affected entity in MemDB.
	entityLocalAliasesToUpsert := map[*identity.Entity][]*identity.Alias{}

	// This map tracks the set of local aliases that need to be removed from
	// their affected entity in MemDB, as well as removing the local alias
	// themselves.
	entityLocalAliasesToRemove := map[*identity.Entity][]*identity.Alias{}

	if bucket != nil {
		// The storage entry for the local alias bucket exists, so we need to
		// compare the local aliases in that bucket with those in MemDB and only
		// update those that are different. The local aliases in the bucket are
		// the source of truth.

		// Iterate over each local alias item from the bucket
		for _, item := range bucket.Items {
			if strings.HasSuffix(item.ID, tmpSuffix) {
				continue
			}

			var bucketLocalAliases identity.LocalAliases

			err = anypb.UnmarshalTo(item.Message, &bucketLocalAliases, proto.UnmarshalOptions{})
			if err != nil {
				i.logger.Error("failed to parse local aliases during invalidation", "item_id", item.ID, "error", err)
				return
			}

			for _, bucketLocalAlias := range bucketLocalAliases.Aliases {
				// Find the entity related to bucketLocalAlias in MemDB in order
				// to track any local aliases modifications that must be made in
				// this entity.
				memDBEntity := i.FetchEntityForLocalAliasInTxn(txn, bucketLocalAlias)
				if memDBEntity == nil {
					// FetchEntityForLocalAliasInTxn already logs any error
					return
				}

				// memDBLocalAlias starts off nil but gets set to the local
				// alias from memDBLocalAliases whose ID matches the ID of
				// bucketLocalAlias.
				var memDBLocalAlias *identity.Alias
				for i, localAlias := range memDBLocalAliases {
					if localAlias.ID == bucketLocalAlias.ID {
						memDBLocalAlias = localAlias

						// Remove this processed local alias from the
						// memDBLocalAliases slice, so that all that
						// will be left are unprocessed local aliases.
						copy(memDBLocalAliases[i:], memDBLocalAliases[i+1:])
						memDBLocalAliases = memDBLocalAliases[:len(memDBLocalAliases)-1]

						break
					}
				}

				// We've considered the use of github.com/google/go-cmp here,
				// but opted for sticking with reflect.DeepEqual because go-cmp
				// is intended for testing and is able to panic in some
				// situations.
				if memDBLocalAlias == nil || !reflect.DeepEqual(memDBLocalAlias, bucketLocalAlias) {
					// The bucketLocalAlias is not in MemDB or it has changed in
					// storage.
					err = i.MemDBUpsertAliasInTxn(txn, bucketLocalAlias, false)
					if err != nil {
						i.logger.Error("failed to update local alias in MemDB", "alias_id", bucketLocalAlias.ID, "error", err)
						return
					}

					// Add this local alias to the set of local aliases that
					// need to be updated for memDBEntity.
					entityLocalAliasesToUpsert[memDBEntity] = append(entityLocalAliasesToUpsert[memDBEntity], bucketLocalAlias)
				}
			}
		}
	}

	// Any local aliases still remaining in memDBLocalAliases do not exist in
	// storage and should be removed from MemDB.
	for _, memDBLocalAlias := range memDBLocalAliases {
		memDBEntity := i.FetchEntityForLocalAliasInTxn(txn, memDBLocalAlias)
		if memDBEntity == nil {
			// FetchEntityForLocalAliasInTxn already logs any error
			return
		}

		entityLocalAliasesToRemove[memDBEntity] = append(entityLocalAliasesToRemove[memDBEntity], memDBLocalAlias)
	}

	// Now process the entityLocalAliasesToUpsert map.
	for entity, localAliases := range entityLocalAliasesToUpsert {
		for _, localAlias := range localAliases {
			entity.UpsertAlias(localAlias)
		}

		err = i.MemDBUpsertEntityInTxn(txn, entity)
		if err != nil {
			i.logger.Error("failed to update entity in MemDB", "entity_id", entity.ID, "error", err)
			return
		}
	}

	// Finally process the entityLocalAliasesToRemove map.
	for entity, localAliases := range entityLocalAliasesToRemove {
		// The deleteAliasesInEntityInTxn removes the provided aliases from
		// the entity, but it also removes the aliases themselves from MemDB.
		err := i.deleteAliasesInEntityInTxn(txn, entity, localAliases)
		if err != nil {
			i.logger.Error("failed to delete aliases in entity", "entity_id", entity.ID, "error", err)
			return
		}

		err = i.MemDBUpsertEntityInTxn(txn, entity)
		if err != nil {
			i.logger.Error("failed to update entity in MemDB", "entity_id", entity.ID, "error", err)
			return
		}
	}

	txn.Commit()
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
		entityAsAny, err := anypb.New(&entity)
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

	return entity.Clone()
}

// CreateOrFetchEntity creates a new entity. This is used by core to
// associate each login attempt by an alias to a unified entity in Vault.
func (i *IdentityStore) CreateOrFetchEntity(ctx context.Context, alias *logical.Alias) (*identity.Entity, bool, error) {
	defer metrics.MeasureSince([]string{"identity", "create_or_fetch_entity"}, time.Now())

	var entity *identity.Entity
	var err error
	var update bool
	var entityCreated bool

	if alias == nil {
		return nil, false, fmt.Errorf("alias is nil")
	}

	if alias.Name == "" {
		return nil, false, fmt.Errorf("empty alias name")
	}

	mountValidationResp := i.router.ValidateMountByAccessor(alias.MountAccessor)
	if mountValidationResp == nil {
		return nil, false, fmt.Errorf("invalid mount accessor %q", alias.MountAccessor)
	}

	if mountValidationResp.MountType != alias.MountType {
		return nil, false, fmt.Errorf("mount accessor %q is not a mount of type %q", alias.MountAccessor, alias.MountType)
	}

	// Check if an entity already exists for the given alias
	entity, err = i.entityByAliasFactors(alias.MountAccessor, alias.Name, true)
	if err != nil {
		return nil, false, err
	}
	if entity != nil && changedAliasIndex(entity, alias) == -1 {
		return entity, false, nil
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	// Create a MemDB transaction to update both alias and entity
	txn := i.db.Txn(true)
	defer txn.Abort()

	// Check if an entity was created before acquiring the lock
	entity, err = i.entityByAliasFactorsInTxn(txn, alias.MountAccessor, alias.Name, true)
	if err != nil {
		return nil, false, err
	}
	if entity != nil {
		idx := changedAliasIndex(entity, alias)
		if idx == -1 {
			return entity, false, nil
		}
		a := entity.Aliases[idx]
		a.Metadata = alias.Metadata
		a.LastUpdateTime = timestamppb.Now()

		update = true
	}

	if !update {
		entity = new(identity.Entity)
		err = i.sanitizeEntity(ctx, entity)
		if err != nil {
			return nil, false, err
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
			return nil, false, err
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
		entityCreated = true
	}

	// Update MemDB and persist entity object
	err = i.upsertEntityInTxn(ctx, txn, entity, nil, true)
	if err != nil {
		return entity, entityCreated, err
	}

	txn.Commit()
	clonedEntity, err := entity.Clone()
	return clonedEntity, entityCreated, err
}

// changedAliasIndex searches an entity for changed alias metadata.
//
// If a match is found, the changed alias's index is returned. If no alias
// names match or no metadata is different, -1 is returned.
func changedAliasIndex(entity *identity.Entity, alias *logical.Alias) int {
	for i, a := range entity.Aliases {
		if a.Name == alias.Name && a.MountAccessor == alias.MountAccessor && !strutil.EqualStringMaps(a.Metadata, alias.Metadata) {
			return i
		}
	}

	return -1
}
