// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/identity/mfa"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/protobuf/types/known/anypb"
)

func entityPathFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"id": {
			Type:        framework.TypeString,
			Description: "ID of the entity. If set, updates the corresponding existing entity.",
		},
		"name": {
			Type:        framework.TypeString,
			Description: "Name of the entity",
		},
		"metadata": {
			Type: framework.TypeKVPairs,
			Description: `Metadata to be associated with the entity.
In CLI, this parameter can be repeated multiple times, and it all gets merged together.
For example:
vault <command> <path> metadata=key1=value1 metadata=key2=value2
					`,
		},
		"policies": {
			Type:        framework.TypeCommaStringSlice,
			Description: "Policies to be tied to the entity.",
		},
		"disabled": {
			Type:        framework.TypeBool,
			Description: "If set true, tokens tied to this identity will not be able to be used (but will not be revoked).",
		},
	}
}

// entityPaths returns the API endpoints supported to operate on entities.
// Following are the paths supported:
// entity - To register a new entity
// entity/id - To lookup, modify, delete and list entities based on ID
// entity/merge - To merge entities based on ID
func entityPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "entity$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationVerb:   "create",
			},

			Fields: entityPathFields(),
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleEntityUpdateCommon(),
				},
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity"][1]),
		},
		{
			Pattern: "entity/name/(?P<name>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationSuffix: "by-name",
			},

			Fields: entityPathFields(),

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleEntityUpdateCommon(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathEntityNameRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathEntityNameDelete(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-name"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-name"][1]),
		},
		{
			Pattern: "entity/id/" + framework.GenericNameRegex("id"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationSuffix: "by-id",
			},

			Fields: entityPathFields(),

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleEntityUpdateCommon(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathEntityIDRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathEntityIDDelete(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-id"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-id"][1]),
		},
		{
			Pattern: "entity/batch-delete",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationVerb:   "batch-delete",
			},

			Fields: map[string]*framework.FieldSchema{
				"entity_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Entity IDs to delete",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleEntityBatchDelete(),
				},
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["batch-delete"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["batch-delete"][1]),
		},
		{
			Pattern: "entity/name/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationSuffix: "by-name",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.pathEntityNameList(),
				},
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-name-list"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-name-list"][1]),
		},
		{
			Pattern: "entity/id/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationSuffix: "by-id",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.pathEntityIDList(),
				},
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-id-list"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-id-list"][1]),
		},
		{
			Pattern: "entity/merge/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationVerb:   "merge",
			},

			Fields: map[string]*framework.FieldSchema{
				"from_entity_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Entity IDs which need to get merged",
				},
				"to_entity_id": {
					Type:        framework.TypeString,
					Description: "Entity ID into which all the other entities need to get merged",
				},
				"conflicting_alias_ids_to_keep": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Alias IDs to keep in case of conflicting aliases. Ignored if no conflicting aliases found",
				},
				"force": {
					Type:        framework.TypeBool,
					Description: "Setting this will follow the 'mine' strategy for merging MFA secrets. If there are secrets of the same type both in entities that are merged from and in entity into which all others are getting merged, secrets in the destination will be unaltered. If not set, this API will throw an error containing all the conflicts.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                  i.pathEntityMergeID(),
					ForwardPerformanceStandby: true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-merge-id"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-merge-id"][1]),
		},
	}
}

// pathEntityMergeID merges two or more entities into a single entity
func (i *IdentityStore) pathEntityMergeID() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		toEntityIDInterface, ok := d.GetOk("to_entity_id")
		if !ok || toEntityIDInterface == "" {
			return logical.ErrorResponse("missing entity id to merge to"), nil
		}
		toEntityID := toEntityIDInterface.(string)

		fromEntityIDsInterface, ok := d.GetOk("from_entity_ids")
		if !ok || len(fromEntityIDsInterface.([]string)) == 0 {
			return logical.ErrorResponse("missing entity ids to merge from"), nil
		}
		fromEntityIDs := fromEntityIDsInterface.([]string)

		var conflictingAliasIDsToKeep []string
		if conflictingAliasIDsToKeepInterface, ok := d.GetOk("conflicting_alias_ids_to_keep"); ok {
			conflictingAliasIDsToKeep = conflictingAliasIDsToKeepInterface.([]string)
		}

		var force bool
		if forceInterface, ok := d.GetOk("force"); ok {
			force = forceInterface.(bool)
		}

		// Create a MemDB transaction to merge entities
		i.lock.Lock()
		defer i.lock.Unlock()

		txn := i.db.Txn(true)
		defer txn.Abort()

		toEntity, err := i.MemDBEntityByID(toEntityID, true)
		if err != nil {
			return nil, err
		}

		userErr, intErr, aliases := i.mergeEntity(ctx, txn, toEntity, fromEntityIDs, conflictingAliasIDsToKeep, force, false, false, true, false)
		if userErr != nil {
			// Not an error due to alias clash, return like normal
			if len(aliases) == 0 {
				return logical.ErrorResponse(userErr.Error()), nil
			}
			// Alias clash error, so include additional details
			return logical.ErrorResponseWithData(aliases, userErr.Error()), nil
		}
		if intErr != nil {
			return nil, intErr
		}

		// Committing the transaction *after* successfully performing storage
		// persistence
		txn.Commit()

		return nil, nil
	}
}

// handleEntityUpdateCommon is used to update an entity
func (i *IdentityStore) handleEntityUpdateCommon() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		i.lock.Lock()
		defer i.lock.Unlock()

		entity := new(identity.Entity)
		var err error

		entityID := d.Get("id").(string)
		if entityID != "" {
			entity, err = i.MemDBEntityByID(entityID, true)
			if err != nil {
				return nil, err
			}
			if entity == nil {
				return logical.ErrorResponse("entity not found from id"), nil
			}
		}

		// Get the name
		entityName := d.Get("name").(string)
		if entityName != "" {
			entityByName, err := i.MemDBEntityByName(ctx, entityName, true)
			if err != nil {
				return nil, err
			}
			switch {
			case entityByName == nil:
				// Not found, safe to use this name with an existing or new entity
			case entity.ID == "":
				// Entity by ID was not found, but and entity for the supplied
				// name was found. Continue updating the entity.
				entity = entityByName
			case entity.ID == entityByName.ID:
				// Same exact entity, carry on (this is basically a noop then)
			default:
				return logical.ErrorResponse("entity name is already in use"), nil
			}
		}

		if entityName != "" {
			entity.Name = entityName
		}

		// Update the policies if supplied
		entityPoliciesRaw, ok := d.GetOk("policies")
		if ok {
			entity.Policies = strutil.RemoveDuplicates(entityPoliciesRaw.([]string), false)
		}

		if strutil.StrListContains(entity.Policies, "root") {
			return logical.ErrorResponse("policies cannot contain root"), nil
		}

		disabledRaw, ok := d.GetOk("disabled")
		if ok {
			entity.Disabled = disabledRaw.(bool)
		}

		// Get entity metadata
		metadata, ok, err := d.GetOkErr("metadata")
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse metadata: %v", err)), nil
		}
		if ok {
			entity.Metadata = metadata.(map[string]string)
		}

		// At this point, if entity.ID is empty, it indicates that a new entity
		// is being created. Using this to respond data in the response.
		newEntity := entity.ID == ""

		// ID creation and some validations
		err = i.sanitizeEntity(ctx, entity)
		if err != nil {
			return nil, err
		}

		if err := i.upsertEntity(ctx, entity, nil, true); err != nil {
			return nil, err
		}

		// If this operation was an update to an existing entity, return 204
		if !newEntity {
			return nil, nil
		}

		// Prepare the response
		respData := map[string]interface{}{
			"id":   entity.ID,
			"name": entity.Name,
		}

		var aliasIDs []string
		for _, alias := range entity.Aliases {
			aliasIDs = append(aliasIDs, alias.ID)
		}

		respData["aliases"] = aliasIDs

		// Return ID of the entity that was either created or updated along with
		// its aliases
		return &logical.Response{
			Data: respData,
		}, nil
	}
}

// pathEntityNameRead returns the properties of an entity for a given entity ID
func (i *IdentityStore) pathEntityNameRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		entityName := d.Get("name").(string)
		if entityName == "" {
			return logical.ErrorResponse("missing entity name"), nil
		}

		entity, err := i.MemDBEntityByName(ctx, entityName, false)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			return nil, nil
		}

		return i.handleEntityReadCommon(ctx, entity)
	}
}

// pathEntityIDRead returns the properties of an entity for a given entity ID
func (i *IdentityStore) pathEntityIDRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		entityID := d.Get("id").(string)
		if entityID == "" {
			return logical.ErrorResponse("missing entity id"), nil
		}

		entity, err := i.MemDBEntityByID(entityID, false)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			return nil, nil
		}

		return i.handleEntityReadCommon(ctx, entity)
	}
}

func (i *IdentityStore) handleEntityReadCommon(ctx context.Context, entity *identity.Entity) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns.ID != entity.NamespaceID {
		return nil, nil
	}

	respData := map[string]interface{}{}
	respData["id"] = entity.ID
	respData["name"] = entity.Name
	respData["metadata"] = entity.Metadata
	respData["merged_entity_ids"] = entity.MergedEntityIDs
	respData["policies"] = strutil.RemoveDuplicates(entity.Policies, false)
	respData["disabled"] = entity.Disabled
	respData["namespace_id"] = entity.NamespaceID

	// Convert protobuf timestamp into RFC3339 format
	respData["creation_time"] = ptypes.TimestampString(entity.CreationTime)
	respData["last_update_time"] = ptypes.TimestampString(entity.LastUpdateTime)

	// Convert each alias into a map and replace the time format in each
	aliasesToReturn := make([]interface{}, len(entity.Aliases))
	for aliasIdx, alias := range entity.Aliases {
		aliasMap := map[string]interface{}{}
		aliasMap["id"] = alias.ID
		aliasMap["canonical_id"] = alias.CanonicalID
		aliasMap["mount_accessor"] = alias.MountAccessor
		aliasMap["metadata"] = alias.Metadata
		aliasMap["name"] = alias.Name
		aliasMap["merged_from_canonical_ids"] = alias.MergedFromCanonicalIDs
		aliasMap["creation_time"] = ptypes.TimestampString(alias.CreationTime)
		aliasMap["last_update_time"] = ptypes.TimestampString(alias.LastUpdateTime)
		aliasMap["local"] = alias.Local
		aliasMap["custom_metadata"] = alias.CustomMetadata

		if mountValidationResp := i.router.ValidateMountByAccessor(alias.MountAccessor); mountValidationResp != nil {
			aliasMap["mount_type"] = mountValidationResp.MountType
			aliasMap["mount_path"] = mountValidationResp.MountPath
		}

		aliasesToReturn[aliasIdx] = aliasMap
	}

	// Add the aliases information to the response which has the correct time
	// formats
	respData["aliases"] = aliasesToReturn

	addExtraEntityDataToResponse(entity, respData)

	// Fetch the groups this entity belongs to and return their identifiers
	groups, inheritedGroups, err := i.groupsByEntityID(entity.ID)
	if err != nil {
		return nil, err
	}

	groupIDs := make([]string, len(groups))
	for i, group := range groups {
		groupIDs[i] = group.ID
	}
	respData["direct_group_ids"] = groupIDs

	inheritedGroupIDs := make([]string, len(inheritedGroups))
	for i, group := range inheritedGroups {
		inheritedGroupIDs[i] = group.ID
	}
	respData["inherited_group_ids"] = inheritedGroupIDs

	respData["group_ids"] = append(groupIDs, inheritedGroupIDs...)

	return &logical.Response{
		Data: respData,
	}, nil
}

// pathEntityIDDelete deletes the entity for a given entity ID
func (i *IdentityStore) pathEntityIDDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		entityID := d.Get("id").(string)
		if entityID == "" {
			return logical.ErrorResponse("missing entity id"), nil
		}

		i.lock.Lock()
		defer i.lock.Unlock()

		// Create a MemDB transaction to delete entity
		txn := i.db.Txn(true)
		defer txn.Abort()

		// Fetch the entity using its ID
		entity, err := i.MemDBEntityByIDInTxn(txn, entityID, true)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			return nil, nil
		}

		err = i.handleEntityDeleteCommon(ctx, txn, entity, true)
		if err != nil {
			return nil, err
		}

		txn.Commit()

		return nil, nil
	}
}

// pathEntityNameDelete deletes the entity for a given entity ID
func (i *IdentityStore) pathEntityNameDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		entityName := d.Get("name").(string)
		if entityName == "" {
			return logical.ErrorResponse("missing entity name"), nil
		}

		i.lock.Lock()
		defer i.lock.Unlock()

		// Create a MemDB transaction to delete entity
		txn := i.db.Txn(true)
		defer txn.Abort()

		// Fetch the entity using its name
		entity, err := i.MemDBEntityByNameInTxn(ctx, txn, entityName, true)
		if err != nil {
			return nil, err
		}
		// If there is no entity for the ID, do nothing
		if entity == nil {
			return nil, nil
		}

		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		if entity.NamespaceID != ns.ID {
			return nil, nil
		}

		err = i.handleEntityDeleteCommon(ctx, txn, entity, true)
		if err != nil {
			return nil, err
		}

		txn.Commit()

		return nil, nil
	}
}

// pathEntityIDDelete deletes the entity for a given entity ID
func (i *IdentityStore) handleEntityBatchDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		entityIDs := d.Get("entity_ids").([]string)
		if len(entityIDs) == 0 {
			return logical.ErrorResponse("missing entity ids to delete"), nil
		}

		// Sort the ids by the bucket they will be deleted from
		byBucket := make(map[string]map[string]struct{})
		for _, id := range entityIDs {
			bucketKey := i.entityPacker.BucketKey(id)

			bucket, ok := byBucket[bucketKey]
			if !ok {
				bucket = make(map[string]struct{})
				byBucket[bucketKey] = bucket
			}

			bucket[id] = struct{}{}
		}

		deleteIdsForBucket := func(entityIDs []string) error {
			i.lock.Lock()
			defer i.lock.Unlock()

			// Create a MemDB transaction to delete entities from the inmem database
			// without altering storage. Batch deletion on storage bucket items is
			// performed directly through entityPacker.
			txn := i.db.Txn(true)
			defer txn.Abort()

			for _, entityID := range entityIDs {
				// Fetch the entity using its ID
				entity, err := i.MemDBEntityByIDInTxn(txn, entityID, true)
				if err != nil {
					return err
				}
				if entity == nil {
					continue
				}

				err = i.handleEntityDeleteCommon(ctx, txn, entity, false)
				if err != nil {
					return err
				}
			}

			// Write all updates for this bucket.
			err := i.entityPacker.DeleteMultipleItems(ctx, i.logger, entityIDs)
			if err != nil {
				return err
			}

			txn.Commit()
			return nil
		}

		for _, bucket := range byBucket {
			ids := make([]string, len(bucket))
			i := 0
			for id := range bucket {
				ids[i] = id
				i++
			}

			err := deleteIdsForBucket(ids)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}
}

// handleEntityDeleteCommon deletes an entity by removing it from groups of
// which it's a member and then, if update is true, deleting the entity itself.
func (i *IdentityStore) handleEntityDeleteCommon(ctx context.Context, txn *memdb.Txn, entity *identity.Entity, update bool) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	if entity.NamespaceID != ns.ID {
		return nil
	}

	// Remove entity ID as a member from all the groups it belongs, both
	// internal and external
	groups, err := i.MemDBGroupsByMemberEntityIDInTxn(txn, entity.ID, true, false)
	if err != nil {
		return err
	}

	for _, group := range groups {
		group.MemberEntityIDs = strutil.StrListDelete(group.MemberEntityIDs, entity.ID)
		err = i.UpsertGroupInTxn(ctx, txn, group, true)
		if err != nil {
			return err
		}
	}

	// Delete all the aliases in the entity and the respective indexes
	err = i.deleteAliasesInEntityInTxn(txn, entity, entity.Aliases)
	if err != nil {
		return err
	}

	// Delete the entity using the same transaction
	err = i.MemDBDeleteEntityByIDInTxn(txn, entity.ID)
	if err != nil {
		return err
	}

	if update {
		// Delete the entity from storage
		err = i.entityPacker.DeleteItem(ctx, entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *IdentityStore) pathEntityIDList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		return i.handlePathEntityListCommon(ctx, req, d, true)
	}
}

func (i *IdentityStore) pathEntityNameList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		return i.handlePathEntityListCommon(ctx, req, d, false)
	}
}

// handlePathEntityListCommon lists the IDs or names of all the valid entities
// in the identity store
func (i *IdentityStore) handlePathEntityListCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, byID bool) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	ws := memdb.NewWatchSet()

	txn := i.db.Txn(false)

	iter, err := txn.Get(entitiesTable, "namespace_id", ns.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch iterator for entities in memdb: %w", err)
	}

	ws.Add(iter.WatchCh())

	var keys []string
	entityInfo := map[string]interface{}{}

	type mountInfo struct {
		MountType string
		MountPath string
	}
	mountAccessorMap := map[string]mountInfo{}

	for {
		// Check for timeouts
		select {
		case <-ctx.Done():
			resp := logical.ListResponseWithInfo(keys, entityInfo)
			resp.AddWarning("partial response due to timeout")
			return resp, nil
		default:
			break
		}

		raw := iter.Next()
		if raw == nil {
			break
		}
		entity := raw.(*identity.Entity)
		if byID {
			keys = append(keys, entity.ID)
		} else {
			keys = append(keys, entity.Name)
		}
		entityInfoEntry := map[string]interface{}{
			"name": entity.Name,
		}
		if len(entity.Aliases) > 0 {
			aliasList := make([]interface{}, 0, len(entity.Aliases))
			for _, alias := range entity.Aliases {
				entry := map[string]interface{}{
					"id":             alias.ID,
					"name":           alias.Name,
					"mount_accessor": alias.MountAccessor,
				}

				mi, ok := mountAccessorMap[alias.MountAccessor]
				if ok {
					entry["mount_type"] = mi.MountType
					entry["mount_path"] = mi.MountPath
				} else {
					mi = mountInfo{}
					if mountValidationResp := i.router.ValidateMountByAccessor(alias.MountAccessor); mountValidationResp != nil {
						mi.MountType = mountValidationResp.MountType
						mi.MountPath = mountValidationResp.MountPath
						entry["mount_type"] = mi.MountType
						entry["mount_path"] = mi.MountPath
					}
					mountAccessorMap[alias.MountAccessor] = mi
				}

				aliasList = append(aliasList, entry)
			}
			entityInfoEntry["aliases"] = aliasList
		}
		entityInfo[entity.ID] = entityInfoEntry
	}

	return logical.ListResponseWithInfo(keys, entityInfo), nil
}

func (i *IdentityStore) mergeEntityAsPartOfUpsert(ctx context.Context, txn *memdb.Txn, toEntity *identity.Entity, fromEntityID string, persist bool) (error, error) {
	err1, err2, _ := i.mergeEntity(ctx, txn, toEntity, []string{fromEntityID}, []string{}, true, false, true, persist, true)
	return err1, err2
}

// A small type to return useful information to the UI after an entity clash
// Every alias involved in a clash will be returned.
type aliasClashInformation struct {
	Alias     string `json:"alias"`
	Entity    string `json:"entity"`
	EntityId  string `json:"entity_id"`
	Mount     string `json:"mount"`
	MountPath string `json:"mount_path"`
}

func (i *IdentityStore) mergeEntity(ctx context.Context, txn *memdb.Txn, toEntity *identity.Entity, fromEntityIDs, conflictingAliasIDsToKeep []string, force, grabLock, mergePolicies, persist, forceMergeAliases bool) (error, error, []aliasClashInformation) {
	if grabLock {
		i.lock.Lock()
		defer i.lock.Unlock()
	}

	if toEntity == nil {
		return errors.New("entity id to merge to is invalid"), nil, nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err, nil
	}
	if toEntity.NamespaceID != ns.ID {
		return errors.New("entity id to merge into does not belong to the request's namespace"), nil, nil
	}

	if len(fromEntityIDs) > 1 && len(conflictingAliasIDsToKeep) > 1 {
		return errors.New("aliases conflicts cannot be resolved with multiple from entity ids - merge one entity at a time"), nil, nil
	}

	sanitizedFromEntityIDs := strutil.RemoveDuplicates(fromEntityIDs, false)

	// A map to check if there are any clashes between mount accessors for any of the sanitizedFromEntityIDs
	fromEntityAccessors := make(map[string]string)

	// A list detailing all aliases where a clash has occurred, so that the error
	// can be understood by the UI
	aliasesInvolvedInClashes := make([]aliasClashInformation, 0)

	// An error detailing if any alias clashes happen (shared mount accessor)
	var aliasClashError error

	for _, fromEntityID := range sanitizedFromEntityIDs {
		if fromEntityID == toEntity.ID {
			return errors.New("to_entity_id should not be present in from_entity_ids"), nil, nil
		}

		fromEntity, err := i.MemDBEntityByID(fromEntityID, false)
		if err != nil {
			return nil, err, nil
		}

		if fromEntity == nil {
			// If forceMergeAliases is true, and we didn't find a fromEntity, then something
			// is wrong with storage. This function was called as part of an automated
			// merge from CreateOrFetchEntity or Invalidate to repatriate an alias with its 'true'
			// entity. As a result, the entity _should_ exist, but we can't find it.
			// MemDb may be in a bad state, because fromEntity should be non-nil in the
			// automated merge case.
			if forceMergeAliases {
				return fmt.Errorf("fromEntity %s was not found in memdb as part of an automated entity merge into %s; storage/memdb may be in a bad state", fromEntityID, toEntity.ID), nil, nil
			}
			return errors.New("entity id to merge from is invalid"), nil, nil
		}

		if fromEntity.NamespaceID != toEntity.NamespaceID {
			return errors.New("entity id to merge from does not belong to this namespace"), nil, nil
		}

		// If we're not resolving a conflict, we check to see if
		// any aliases conflict between the toEntity and this fromEntity:
		if !forceMergeAliases && len(conflictingAliasIDsToKeep) == 0 {
			for _, toAlias := range toEntity.Aliases {
				for _, fromAlias := range fromEntity.Aliases {
					// First, check to see if this alias clashes with an alias from any of the other fromEntities:
					id, mountAccessorInAnotherFromEntity := fromEntityAccessors[fromAlias.MountAccessor]
					if mountAccessorInAnotherFromEntity && (id != fromEntityID) {
						return fmt.Errorf("mount accessor %s found in multiple fromEntities, merge should be done with one fromEntity at a time", fromAlias.MountAccessor), nil, nil
					}

					fromEntityAccessors[fromAlias.MountAccessor] = fromEntityID

					// If it doesn't, check if it clashes with the toEntities
					if toAlias.MountAccessor == fromAlias.MountAccessor {
						if aliasClashError == nil {
							aliasClashError = multierror.Append(aliasClashError, fmt.Errorf("toEntity and at least one fromEntity have aliases with the same mount accessor, repeat the merge request specifying exactly one fromEntity, clashes: "))
						}
						aliasClashError = multierror.Append(aliasClashError,
							fmt.Errorf("mountAccessor: %s, toEntity ID: %s, fromEntity ID: %s, conflicting toEntity alias ID: %s, conflicting fromEntity alias ID: %s",
								toAlias.MountAccessor, toEntity.ID, fromEntityID, toAlias.ID, fromAlias.ID))

						var toAliasMountType string
						var toAliasMountPath string
						mountValidationRespToAlias := i.router.ValidateMountByAccessor(toAlias.MountAccessor)
						if mountValidationRespToAlias != nil {
							toAliasMountType = mountValidationRespToAlias.MountType
							toAliasMountPath = mountValidationRespToAlias.MountPath
						}

						var fromAliasMountType string
						var fromAliasMountPath string
						mountValidationRespFromAlias := i.router.ValidateMountByAccessor(fromAlias.MountAccessor)
						if mountValidationRespFromAlias != nil {
							fromAliasMountType = mountValidationRespFromAlias.MountType
							fromAliasMountPath = mountValidationRespFromAlias.MountPath
						}

						// Also add both to our summary of all clashes:
						aliasesInvolvedInClashes = append(aliasesInvolvedInClashes, aliasClashInformation{
							Entity:    toEntity.Name,
							EntityId:  toEntity.ID,
							Alias:     toAlias.Name,
							Mount:     toAliasMountType,
							MountPath: toAliasMountPath,
						})
						aliasesInvolvedInClashes = append(aliasesInvolvedInClashes, aliasClashInformation{
							Entity:    fromEntity.Name,
							EntityId:  fromEntityID,
							Alias:     fromAlias.Name,
							Mount:     fromAliasMountType,
							MountPath: fromAliasMountPath,
						})
					}
				}
			}
		}

		for configID, configSecret := range fromEntity.MFASecrets {
			_, ok := toEntity.MFASecrets[configID]
			if ok && !force {
				return nil, fmt.Errorf("conflicting MFA config ID %q in entity ID %q", configID, fromEntity.ID), nil
			} else {
				if toEntity.MFASecrets == nil {
					toEntity.MFASecrets = make(map[string]*mfa.Secret)
				}
				toEntity.MFASecrets[configID] = configSecret
			}
		}
	}

	// Check alias clashes after validating every fromEntity, so that we have a full list of errors
	if aliasClashError != nil {
		return aliasClashError, nil, aliasesInvolvedInClashes
	}

	isPerfSecondaryOrStandby := i.localNode.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) ||
		i.localNode.HAState() == consts.PerfStandby
	var fromEntityGroups []*identity.Group

	toEntityAccessors := make(map[string][]string)

	for _, alias := range toEntity.Aliases {
		if accessors, ok := toEntityAccessors[alias.MountAccessor]; !ok {
			// While it is not supported to have multiple aliases with the same mount accessor in one entity
			// we do not strictly enforce the invariant. Thus, we account for multiple just to be safe
			if accessors == nil {
				toEntityAccessors[alias.MountAccessor] = []string{alias.ID}
			} else {
				toEntityAccessors[alias.MountAccessor] = append(accessors, alias.ID)
			}
		}
	}

	for _, fromEntityID := range sanitizedFromEntityIDs {
		if fromEntityID == toEntity.ID {
			return errors.New("to_entity_id should not be present in from_entity_ids"), nil, nil
		}

		fromEntity, err := i.MemDBEntityByID(fromEntityID, true)
		if err != nil {
			return nil, err, nil
		}

		if fromEntity == nil {
			// If forceMergeAliases is true, and we didn't find a fromEntity, then something
			// is wrong with storage. This function was called as part of an automated
			// merge from CreateOrFetchEntity or Invalidate to repatriate an alias with its 'true'
			// entity. As a result, the entity _should_ exist, but we can't find it.
			// MemDb may be in a bad state, because fromEntity should be non-nil in the
			// automated merge case.
			if forceMergeAliases {
				return fmt.Errorf("fromEntity %s was not found in memdb as part of an automated entity merge into %s; storage/memdb may be in a bad state", fromEntityID, toEntity.ID), nil, nil
			}
			return errors.New("entity id to merge from is invalid"), nil, nil
		}

		if fromEntity.NamespaceID != toEntity.NamespaceID {
			return errors.New("entity id to merge from does not belong to this namespace"), nil, nil
		}

		for _, fromAlias := range fromEntity.Aliases {
			// If true, we need to handle conflicts (conflict = both aliases share the same mount accessor)
			if toAliasIds, ok := toEntityAccessors[fromAlias.MountAccessor]; ok {
				for _, toAliasId := range toAliasIds {
					// When forceMergeAliases is true (as part of the merge-during-upsert case), we make the decision
					// for the user, and keep the to_entity alias, merging the from_entity
					// This case's code is the same as when the user selects to keep the from_entity alias
					// but is kept separate for clarity
					if forceMergeAliases {
						i.logger.Info("Deleting to_entity alias during entity merge", "to_entity", toEntity.ID, "deleted_alias", toAliasId)
						err := i.MemDBDeleteAliasByIDInTxn(txn, toAliasId, false)
						if err != nil {
							return nil, fmt.Errorf("aborting entity merge - failed to delete orphaned alias %q during merge into entity %q: %w", toAliasId, toEntity.ID, err), nil
						}
						// Remove the alias from the entity's list in memory too!
						toEntity.DeleteAliasByID(toAliasId)
					} else if strutil.StrListContains(conflictingAliasIDsToKeep, toAliasId) {
						i.logger.Info("Deleting from_entity alias during entity merge", "from_entity", fromEntityID, "deleted_alias", fromAlias.ID)
						err := i.MemDBDeleteAliasByIDInTxn(txn, fromAlias.ID, false)
						if err != nil {
							return nil, fmt.Errorf("aborting entity merge - failed to delete orphaned alias %q during merge into entity %q: %w", fromAlias.ID, toEntity.ID, err), nil
						}
						// Remove the alias from the entity's list in memory too!
						toEntity.DeleteAliasByID(toAliasId)

						// Continue to next alias, as there's no alias to merge left in the from_entity
						continue
					} else if strutil.StrListContains(conflictingAliasIDsToKeep, fromAlias.ID) {
						i.logger.Info("Deleting to_entity alias during entity merge", "to_entity", toEntity.ID, "deleted_alias", toAliasId)
						err := i.MemDBDeleteAliasByIDInTxn(txn, toAliasId, false)
						if err != nil {
							return nil, fmt.Errorf("aborting entity merge - failed to delete orphaned alias %q during merge into entity %q: %w", toAliasId, toEntity.ID, err), nil
						}
						// Remove the alias from the entity's list in memory too!
						toEntity.DeleteAliasByID(toAliasId)
					} else {
						return fmt.Errorf("conflicting mount accessors in following alias IDs and neither were present in conflicting_alias_ids_to_keep: %s, %s", fromAlias.ID, toAliasId), nil, nil
					}
				}
			}

			// Set the desired canonical ID
			fromAlias.CanonicalID = toEntity.ID

			fromAlias.MergedFromCanonicalIDs = append(fromAlias.MergedFromCanonicalIDs, fromEntity.ID)

			err = i.MemDBUpsertAliasInTxn(txn, fromAlias, false)
			if err != nil {
				return nil, fmt.Errorf("failed to update alias during merge: %w", err), nil
			}

			// Add the alias to the desired entity
			toEntity.Aliases = append(toEntity.Aliases, fromAlias)
		}

		// If told to, merge policies
		if mergePolicies {
			toEntity.Policies = strutil.RemoveDuplicates(strutil.MergeSlices(toEntity.Policies, fromEntity.Policies), false)
		}

		// If the entity from which we are merging from was already a merged
		// entity, transfer over the Merged set to the entity we are
		// merging into.
		toEntity.MergedEntityIDs = append(toEntity.MergedEntityIDs, fromEntity.MergedEntityIDs...)

		// Add the entity from which we are merging from to the list of entities
		// the entity we are merging into is composed of.
		toEntity.MergedEntityIDs = append(toEntity.MergedEntityIDs, fromEntity.ID)

		// Remove entity ID as a member from all the groups it belongs, both
		// internal and external
		groups, err := i.MemDBGroupsByMemberEntityIDInTxn(txn, fromEntity.ID, true, false)
		if err != nil {
			return nil, err, nil
		}
		for _, group := range groups {
			group.MemberEntityIDs = strutil.StrListDelete(group.MemberEntityIDs, fromEntity.ID)
			err = i.UpsertGroupInTxn(ctx, txn, group, persist && !isPerfSecondaryOrStandby)
			if err != nil {
				return nil, err, nil
			}

			fromEntityGroups = append(fromEntityGroups, group)
		}

		// Delete the entity which we are merging from in MemDB using the same transaction
		err = i.MemDBDeleteEntityByIDInTxn(txn, fromEntity.ID)
		if err != nil {
			return nil, err, nil
		}

		if persist && !isPerfSecondaryOrStandby {
			// Delete the entity which we are merging from in storage
			err = i.entityPacker.DeleteItem(ctx, fromEntity.ID)
			if err != nil {
				return nil, err, nil
			}
		}
	}

	// Update MemDB with changes to the entity we are merging to
	err = i.MemDBUpsertEntityInTxn(txn, toEntity)
	if err != nil {
		return nil, err, nil
	}

	for _, group := range fromEntityGroups {
		group.MemberEntityIDs = append(group.MemberEntityIDs, toEntity.ID)
		err = i.UpsertGroupInTxn(ctx, txn, group, persist && !isPerfSecondaryOrStandby)
		if err != nil {
			return nil, err, nil
		}
	}

	if persist && !isPerfSecondaryOrStandby {
		// Persist the entity which we are merging to
		toEntityAsAny, err := anypb.New(toEntity)
		if err != nil {
			return nil, err, nil
		}
		item := &storagepacker.Item{
			ID:      toEntity.ID,
			Message: toEntityAsAny,
		}

		err = i.entityPacker.PutItem(ctx, item)
		if err != nil {
			return nil, err, nil
		}
	}

	return nil, nil, nil
}

var entityHelp = map[string][2]string{
	"entity": {
		"Create a new entity",
		"",
	},
	"entity-id": {
		"Update, read or delete an entity using entity ID",
		"",
	},
	"entity-name": {
		"Update, read or delete an entity using entity name",
		"",
	},
	"entity-id-list": {
		"List all the entity IDs",
		"",
	},
	"entity-name-list": {
		"List all the entity names",
		"",
	},
	"entity-merge-id": {
		"Merge two or more entities together",
		"",
	},
	"batch-delete": {
		"Delete all of the entities provided",
		"",
	},
}
