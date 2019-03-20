package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/errwrap"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
			Fields:  entityPathFields(),
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleEntityUpdateCommon(),
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity"][1]),
		},
		{
			Pattern: "entity/name/" + framework.GenericNameRegex("name"),
			Fields:  entityPathFields(),
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleEntityUpdateCommon(),
				logical.ReadOperation:   i.pathEntityNameRead(),
				logical.DeleteOperation: i.pathEntityNameDelete(),
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-name"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-name"][1]),
		},
		{
			Pattern: "entity/id/" + framework.GenericNameRegex("id"),
			Fields:  entityPathFields(),
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleEntityUpdateCommon(),
				logical.ReadOperation:   i.pathEntityIDRead(),
				logical.DeleteOperation: i.pathEntityIDDelete(),
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-id"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-id"][1]),
		},
		{
			Pattern: "entity/name/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathEntityNameList(),
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-name-list"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-name-list"][1]),
		},
		{
			Pattern: "entity/id/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathEntityIDList(),
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-id-list"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-id-list"][1]),
		},
		{
			Pattern: "entity/merge/?$",
			Fields: map[string]*framework.FieldSchema{
				"from_entity_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Entity IDs which needs to get merged",
				},
				"to_entity_id": {
					Type:        framework.TypeString,
					Description: "Entity ID into which all the other entities need to get merged",
				},
				"force": {
					Type:        framework.TypeBool,
					Description: "Setting this will follow the 'mine' strategy for merging MFA secrets. If there are secrets of the same type both in entities that are merged from and in entity into which all others are getting merged, secrets in the destination will be unaltered. If not set, this API will throw an error containing all the conflicts.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathEntityMergeID(),
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["entity-merge-id"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["entity-merge-id"][1]),
		},
	}
}

// pathEntityMergeID merges two or more entities into a single entity
func (i *IdentityStore) pathEntityMergeID() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		toEntityID := d.Get("to_entity_id").(string)
		if toEntityID == "" {
			return logical.ErrorResponse("missing entity id to merge to"), nil
		}

		fromEntityIDs := d.Get("from_entity_ids").([]string)
		if len(fromEntityIDs) == 0 {
			return logical.ErrorResponse("missing entity ids to merge from"), nil
		}

		force := d.Get("force").(bool)

		// Create a MemDB transaction to merge entities
		txn := i.db.Txn(true)
		defer txn.Abort()

		toEntity, err := i.MemDBEntityByID(toEntityID, true)
		if err != nil {
			return nil, err
		}

		userErr, intErr := i.mergeEntity(ctx, txn, toEntity, fromEntityIDs, force, true, false, true)
		if userErr != nil {
			return logical.ErrorResponse(userErr.Error()), nil
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
			entityByName, err := i.MemDBEntityByName(ctx, entityName, false)
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
			entity.Policies = entityPoliciesRaw.([]string)
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
			"id": entity.ID,
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
	respData["policies"] = entity.Policies
	respData["disabled"] = entity.Disabled

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

		if mountValidationResp := i.core.router.validateMountByAccessor(alias.MountAccessor); mountValidationResp != nil {
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

		err = i.handleEntityDeleteCommon(ctx, txn, entity)
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

		err = i.handleEntityDeleteCommon(ctx, txn, entity)
		if err != nil {
			return nil, err
		}

		txn.Commit()

		return nil, nil
	}
}

func (i *IdentityStore) handleEntityDeleteCommon(ctx context.Context, txn *memdb.Txn, entity *identity.Entity) error {
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
		return nil
	}

	for _, group := range groups {
		group.MemberEntityIDs = strutil.StrListDelete(group.MemberEntityIDs, entity.ID)
		err = i.UpsertGroupInTxn(txn, group, true)
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

	// Delete the entity from storage
	err = i.entityPacker.DeleteItem(entity.ID)
	if err != nil {
		return err
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
		return nil, errwrap.Wrapf("failed to fetch iterator for entities in memdb: {{err}}", err)
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
					if mountValidationResp := i.core.router.validateMountByAccessor(alias.MountAccessor); mountValidationResp != nil {
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

func (i *IdentityStore) mergeEntity(ctx context.Context, txn *memdb.Txn, toEntity *identity.Entity, fromEntityIDs []string, force, grabLock, mergePolicies, persist bool) (error, error) {
	if grabLock {
		i.lock.Lock()
		defer i.lock.Unlock()
	}

	if toEntity == nil {
		return errors.New("entity id to merge to is invalid"), nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if toEntity.NamespaceID != ns.ID {
		return errors.New("entity id to merge into does not belong to the request's namespace"), nil
	}

	// Merge the MFA secrets
	for _, fromEntityID := range fromEntityIDs {
		if fromEntityID == toEntity.ID {
			return errors.New("to_entity_id should not be present in from_entity_ids"), nil
		}

		fromEntity, err := i.MemDBEntityByID(fromEntityID, false)
		if err != nil {
			return nil, err
		}

		if fromEntity == nil {
			return errors.New("entity id to merge from is invalid"), nil
		}

		if fromEntity.NamespaceID != toEntity.NamespaceID {
			return errors.New("entity id to merge from does not belong to this namespace"), nil
		}

		for configID, configSecret := range fromEntity.MFASecrets {
			_, ok := toEntity.MFASecrets[configID]
			if ok && !force {
				return nil, fmt.Errorf("conflicting MFA config ID %q in entity ID %q", configID, fromEntity.ID)
			} else {
				toEntity.MFASecrets[configID] = configSecret
			}
		}
	}

	isPerfSecondaryOrStandby := i.core.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) || i.core.perfStandby
	for _, fromEntityID := range fromEntityIDs {
		if fromEntityID == toEntity.ID {
			return errors.New("to_entity_id should not be present in from_entity_ids"), nil
		}

		fromEntity, err := i.MemDBEntityByID(fromEntityID, false)
		if err != nil {
			return nil, err
		}

		if fromEntity == nil {
			return errors.New("entity id to merge from is invalid"), nil
		}

		if fromEntity.NamespaceID != toEntity.NamespaceID {
			return errors.New("entity id to merge from does not belong to this namespace"), nil
		}

		for _, alias := range fromEntity.Aliases {
			// Set the desired canonical ID
			alias.CanonicalID = toEntity.ID

			alias.MergedFromCanonicalIDs = append(alias.MergedFromCanonicalIDs, fromEntity.ID)

			err = i.MemDBUpsertAliasInTxn(txn, alias, false)
			if err != nil {
				return nil, errwrap.Wrapf("failed to update alias during merge: {{err}}", err)
			}

			// Add the alias to the desired entity
			toEntity.Aliases = append(toEntity.Aliases, alias)
		}

		// If told to, merge policies
		if mergePolicies {
			toEntity.Policies = strutil.MergeSlices(toEntity.Policies, fromEntity.Policies)
		}

		// If the entity from which we are merging from was already a merged
		// entity, transfer over the Merged set to the entity we are
		// merging into.
		toEntity.MergedEntityIDs = append(toEntity.MergedEntityIDs, fromEntity.MergedEntityIDs...)

		// Add the entity from which we are merging from to the list of entities
		// the entity we are merging into is composed of.
		toEntity.MergedEntityIDs = append(toEntity.MergedEntityIDs, fromEntity.ID)

		// Delete the entity which we are merging from in MemDB using the same transaction
		err = i.MemDBDeleteEntityByIDInTxn(txn, fromEntity.ID)
		if err != nil {
			return nil, err
		}

		if persist && !isPerfSecondaryOrStandby {
			// Delete the entity which we are merging from in storage
			err = i.entityPacker.DeleteItem(fromEntity.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	// Update MemDB with changes to the entity we are merging to
	err = i.MemDBUpsertEntityInTxn(txn, toEntity)
	if err != nil {
		return nil, err
	}

	if persist && !isPerfSecondaryOrStandby {
		// Persist the entity which we are merging to
		toEntityAsAny, err := ptypes.MarshalAny(toEntity)
		if err != nil {
			return nil, err
		}
		item := &storagepacker.Item{
			ID:      toEntity.ID,
			Message: toEntityAsAny,
		}

		err = i.entityPacker.PutItem(item)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
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
}
