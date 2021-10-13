package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const maxCustomMetadataKeys = 64
const maxCustomMetadataKeyLength = 128
const maxCustomMetadataValueLength = 512
const customMetadataValidationErrorPrefix = "custom_metadata validation failed"

// aliasPaths returns the API endpoints to operate on aliases.
// Following are the paths supported:
// entity-alias - To register/modify an alias
// entity-alias/id - To read, modify, delete and list aliases based on their ID
func aliasPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "entity-alias$",
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the entity alias. If set, updates the corresponding entity alias.",
				},
				// entity_id is deprecated in favor of canonical_id
				"entity_id": {
					Type: framework.TypeString,
					Description: `Entity ID to which this alias belongs.
This field is deprecated, use canonical_id.`,
				},
				"canonical_id": {
					Type:        framework.TypeString,
					Description: "Entity ID to which this alias belongs",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Mount accessor to which this alias belongs to; unused for a modify",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the alias; unused for a modify",
				},
				"custom_metadata": {
					Type:        framework.TypeKVPairs,
					Description: "User provided key-value pairs",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleAliasCreateUpdate(),
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias"][1]),
		},
		{
			Pattern: "entity-alias/id/" + framework.GenericNameRegex("id"),
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the alias",
				},
				// entity_id is deprecated
				"entity_id": {
					Type: framework.TypeString,
					Description: `Entity ID to which this alias belongs to.
This field is deprecated, use canonical_id.`,
				},
				"canonical_id": {
					Type:        framework.TypeString,
					Description: "Entity ID to which this alias should be tied to",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "(Unused)",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "(Unused)",
				},
				"custom_metadata": {
					Type:        framework.TypeKVPairs,
					Description: "User provided key-value pairs",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleAliasCreateUpdate(),
				logical.ReadOperation:   i.pathAliasIDRead(),
				logical.DeleteOperation: i.pathAliasIDDelete(),
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias-id"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias-id"][1]),
		},
		{
			Pattern: "entity-alias/id/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathAliasIDList(),
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias-id-list"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias-id-list"][1]),
		},
	}
}

// handleAliasCreateUpdate is used to create or update an alias
func (i *IdentityStore) handleAliasCreateUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		var err error

		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		// Get alias name, if any
		name := d.Get("name").(string)

		// Get mount accessor, if any
		mountAccessor := d.Get("mount_accessor").(string)

		// Get ID, if any
		id := d.Get("id").(string)

		// Get custom metadata, if any
		customMetadata := make(map[string]string)
		data, customMetadataExists := d.GetOk("custom_metadata")
		if customMetadataExists {
			err = mapstructure.Decode(data, &customMetadata)
			if err != nil {
				return nil, err
			}
		}

		// Get entity id
		canonicalID := d.Get("canonical_id").(string)
		if canonicalID == "" {
			// For backwards compatibility
			canonicalID = d.Get("entity_id").(string)
		}

		//validate customMetadata if provided
		if len(customMetadata) != 0 {

			err := validateCustomMetadata(customMetadata)
			if err != nil {
				return nil, err
			}
		}

		i.lock.Lock()
		defer i.lock.Unlock()

		// This block is run if they provided an ID
		{
			// If they provide an ID it must be an update. Find the alias, perform
			// due diligence, call the update function.
			if id != "" {
				alias, err := i.MemDBAliasByID(id, true, false)
				if err != nil {
					return nil, err
				}
				if alias == nil {
					return logical.ErrorResponse("invalid alias ID provided"), nil
				}
				if alias.NamespaceID != ns.ID {
					return logical.ErrorResponse("cannot modify aliases across namespaces"), logical.ErrPermissionDenied
				}
				switch {
				case mountAccessor == "" && name == "" && len(customMetadata) == 0:
					// Just a canonical ID update, maybe
					if canonicalID == "" {
						// Nothing to do, so be idempotent
						return nil, nil
					}
					name = alias.Name
					mountAccessor = alias.MountAccessor
					customMetadata = alias.CustomMetadata

				case mountAccessor == "":
					// No change to mount accessor
					mountAccessor = alias.MountAccessor
				case name == "":
					// No change to mount name
					name = alias.Name
				case len(customMetadata) == 0:
					// No change to custom metadata
					customMetadata = alias.CustomMetadata

				default:
					// mountAccessor, name and customMetadata  provided
				}
				return i.handleAliasUpdate(ctx, req, canonicalID, name, mountAccessor, alias, customMetadata)
			}
		}

		// If they didn't provide an ID, we must have both accessor and name provided
		if mountAccessor == "" || name == "" {
			return logical.ErrorResponse("'id' or 'mount_accessor' and 'name' must be provided"), nil
		}

		// Look up the alias by factors; if it's found it's an update
		mountEntry := i.router.MatchingMountByAccessor(mountAccessor)
		if mountEntry == nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid mount accessor %q", mountAccessor)), nil
		}
		if mountEntry.Local {
			return logical.ErrorResponse(fmt.Sprintf("mount accessor %q is of a local mount", mountAccessor)), nil
		}
		if mountEntry.NamespaceID != ns.ID {
			return logical.ErrorResponse("matching mount is in a different namespace than request"), logical.ErrPermissionDenied
		}
		alias, err := i.MemDBAliasByFactors(mountAccessor, name, false, false)
		if err != nil {
			return nil, err
		}
		if alias != nil {
			if alias.NamespaceID != ns.ID {
				return logical.ErrorResponse("cannot modify aliases across namespaces"), logical.ErrPermissionDenied
			}

			return i.handleAliasUpdate(ctx, req, alias.CanonicalID, name, mountAccessor, alias, customMetadata)
		}
		// At this point we know it's a new creation request
		return i.handleAliasCreate(ctx, req, canonicalID, name, mountAccessor, customMetadata)
	}
}

func (i *IdentityStore) handleAliasCreate(ctx context.Context, req *logical.Request, canonicalID, name, mountAccessor string, customMetadata map[string]string) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	alias := &identity.Alias{
		MountAccessor:  mountAccessor,
		Name:           name,
		CustomMetadata: customMetadata,
	}

	entity := &identity.Entity{}

	// If a canonical ID is provided pull up the entity and make sure we're in
	// the right NS
	if canonicalID != "" {
		entity, err = i.MemDBEntityByID(canonicalID, true)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			return logical.ErrorResponse("invalid canonical ID"), nil
		}
		if entity.NamespaceID != ns.ID {
			return logical.ErrorResponse("entity found with 'canonical_id' not in request namespace"), logical.ErrPermissionDenied
		}
	}

	for _, currentAlias := range entity.Aliases {
		if currentAlias.MountAccessor == mountAccessor {
			return logical.ErrorResponse("Alias already exists for requested entity and mount accessor"), nil
		}
	}

	entity.Aliases = append(entity.Aliases, alias)

	// ID creation and other validations; This is more useful for new entities
	// and may not perform anything for the existing entities. Placing the
	// check here to make the flow common for both new and existing entities.
	err = i.sanitizeEntity(ctx, entity)
	if err != nil {
		return nil, err
	}

	// Set the canonical ID in the alias index. This should be done after
	// sanitizing entity in case it's a new entity that didn't have an ID.
	alias.CanonicalID = entity.ID

	// ID creation and other validations
	err = i.sanitizeAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	// Index entity and its aliases in MemDB and persist entity along with
	// aliases in storage.
	if err := i.upsertEntity(ctx, entity, nil, true); err != nil {
		return nil, err
	}

	// Return ID of both alias and entity
	return &logical.Response{
		Data: map[string]interface{}{
			"id":           alias.ID,
			"canonical_id": entity.ID,
		},
	}, nil
}

func (i *IdentityStore) handleAliasUpdate(ctx context.Context, req *logical.Request, canonicalID, name, mountAccessor string, alias *identity.Alias, customMetadata map[string]string) (*logical.Response, error) {
	if name == alias.Name &&
		mountAccessor == alias.MountAccessor &&
		(canonicalID == alias.CanonicalID || canonicalID == "") && (strutil.EqualStringMaps(customMetadata, alias.CustomMetadata)) {
		// Nothing to do; return nil to be idempotent
		return nil, nil
	}

	alias.LastUpdateTime = ptypes.TimestampNow()

	// Get our current entity, which may be the same as the new one if the
	// canonical ID hasn't changed
	currentEntity, err := i.MemDBEntityByAliasID(alias.ID, true)
	if err != nil {
		return nil, err
	}
	if currentEntity == nil {
		return logical.ErrorResponse("given alias is not associated with an entity"), nil
	}

	if currentEntity.NamespaceID != alias.NamespaceID {
		return logical.ErrorResponse("alias and entity do not belong to the same namespace"), logical.ErrPermissionDenied
	}

	// If the accessor is being changed but the entity is not, check if the entity
	// already has an alias corresponding to the new accessor
	if mountAccessor != alias.MountAccessor && (canonicalID == "" || canonicalID == alias.CanonicalID) {
		for _, currentAlias := range currentEntity.Aliases {
			if currentAlias.MountAccessor == mountAccessor {
				return logical.ErrorResponse("Alias cannot be updated as the entity already has an alias for the given 'mount_accessor' "), nil
			}
		}
	}
	// If we're changing one or the other or both of these, make sure that
	// there isn't a matching alias already, and make sure it's in the same
	// namespace.
	if name != alias.Name || mountAccessor != alias.MountAccessor || !strutil.EqualStringMaps(customMetadata, alias.CustomMetadata) {
		// Check here to see if such an alias already exists, if so bail
		mountEntry := i.router.MatchingMountByAccessor(mountAccessor)
		if mountEntry == nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid mount accessor %q", mountAccessor)), nil
		}
		if mountEntry.Local {
			return logical.ErrorResponse(fmt.Sprintf("mount_accessor %q is of a local mount", mountAccessor)), nil
		}
		if mountEntry.NamespaceID != alias.NamespaceID {
			return logical.ErrorResponse("given mount accessor is not in the same namespace as the existing alias"), logical.ErrPermissionDenied
		}

		existingAlias, err := i.MemDBAliasByFactors(mountAccessor, name, false, false)
		if err != nil {
			return nil, err
		}

		// Bail unless it's just a case change
		if existingAlias != nil && existingAlias.ID != alias.ID {
			return logical.ErrorResponse("alias with combination of mount accessor and name already exists"), nil
		}

		// Update the values in the alias
		alias.Name = name
		alias.MountAccessor = mountAccessor
		alias.CustomMetadata = customMetadata
	}

	newEntity := currentEntity
	if canonicalID != "" && canonicalID != alias.CanonicalID {
		newEntity, err = i.MemDBEntityByID(canonicalID, true)
		if err != nil {
			return nil, err
		}
		if newEntity == nil {
			return logical.ErrorResponse("given 'canonical_id' is not associated with an entity"), nil
		}
		if newEntity.NamespaceID != alias.NamespaceID {
			return logical.ErrorResponse("given 'canonical_id' associated with entity in a different namespace from the alias"), logical.ErrPermissionDenied
		}

		// Check if the entity the alias is being updated to, already has an alias for the mount
		for _, alias := range newEntity.Aliases {
			if alias.MountAccessor == mountAccessor {
				return logical.ErrorResponse("Alias cannot be updated as the given entity already has an alias for this mount "), nil
			}
		}

		// Update the canonical ID value and move it from the current entity to the new one
		alias.CanonicalID = newEntity.ID
		newEntity.Aliases = append(newEntity.Aliases, alias)
		for aliasIndex, item := range currentEntity.Aliases {
			if item.ID == alias.ID {
				currentEntity.Aliases = append(currentEntity.Aliases[:aliasIndex], currentEntity.Aliases[aliasIndex+1:]...)
				break
			}
		}
	} else {
		// If it's not moving we still need to update it in the existing
		// entity's aliases
		for aliasIndex, item := range currentEntity.Aliases {
			if item.ID == alias.ID {
				currentEntity.Aliases[aliasIndex] = alias
				break
			}
		}
		// newEntity will be pointing to the same entity; set currentEntity nil
		// so the upsertCall gets nil for the previous entity as we're only
		// changing one.
		currentEntity = nil
	}

	// Index entity and its aliases in MemDB and persist entity along with
	// aliases in storage. If the alias is being transferred over from
	// one entity to another, previous entity needs to get refreshed in MemDB
	// and persisted in storage as well.
	if err := i.upsertEntity(ctx, newEntity, currentEntity, true); err != nil {
		return nil, err
	}

	// Return ID of both alias and entity
	return &logical.Response{
		Data: map[string]interface{}{
			"id":           alias.ID,
			"canonical_id": newEntity.ID,
		},
	}, nil
}

func validateCustomMetadata(customMetadata map[string]string) error {
	var errs *multierror.Error

	if keyCount := len(customMetadata); keyCount > maxCustomMetadataKeys {
		errs = multierror.Append(errs, fmt.Errorf("%s: payload must contain at most %d keys, provided %d",
			customMetadataValidationErrorPrefix,
			maxCustomMetadataKeys,
			keyCount))

		return errs.ErrorOrNil()
	}

	// Perform validation on each key and value and return ALL errors
	for key, value := range customMetadata {
		if keyLen := len(key); 0 == keyLen || keyLen > maxCustomMetadataKeyLength {
			errs = multierror.Append(errs, fmt.Errorf("%s: length of key %q is %d but must be 0 < len(key) <= %d",
				customMetadataValidationErrorPrefix,
				key,
				keyLen,
				maxCustomMetadataKeyLength))
		}

		if valueLen := len(value); 0 == valueLen || valueLen > maxCustomMetadataValueLength {
			errs = multierror.Append(errs, fmt.Errorf("%s: length of value for key %q is %d but must be 0 < len(value) <= %d",
				customMetadataValidationErrorPrefix,
				key,
				valueLen,
				maxCustomMetadataValueLength))
		}

		if !strutil.Printable(key) {
			// Include unquoted format (%s) to also include the string without the unprintable
			//  characters visible to allow for easier debug and key identification
			errs = multierror.Append(errs, fmt.Errorf("%s: key %q (%s) contains unprintable characters",
				customMetadataValidationErrorPrefix,
				key,
				key))
		}

		if !strutil.Printable(value) {
			errs = multierror.Append(errs, fmt.Errorf("%s: value for key %q contains unprintable characters",
				customMetadataValidationErrorPrefix,
				key))
		}
	}

	return errs.ErrorOrNil()
}

// pathAliasIDRead returns the properties of an alias for a given
// alias ID
func (i *IdentityStore) pathAliasIDRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		aliasID := d.Get("id").(string)
		if aliasID == "" {
			return logical.ErrorResponse("missing alias id"), nil
		}

		alias, err := i.MemDBAliasByID(aliasID, false, false)
		if err != nil {
			return nil, err
		}

		return i.handleAliasReadCommon(ctx, alias)
	}
}

func (i *IdentityStore) handleAliasReadCommon(ctx context.Context, alias *identity.Alias) (*logical.Response, error) {
	if alias == nil {
		return nil, nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns.ID != alias.NamespaceID {
		return logical.ErrorResponse("alias and request are in different namespaces"), logical.ErrPermissionDenied
	}

	respData := map[string]interface{}{}
	respData["id"] = alias.ID
	respData["canonical_id"] = alias.CanonicalID
	respData["mount_accessor"] = alias.MountAccessor
	respData["metadata"] = alias.Metadata
	respData["custom_metadata"] = alias.CustomMetadata
	respData["name"] = alias.Name
	respData["merged_from_canonical_ids"] = alias.MergedFromCanonicalIDs
	respData["namespace_id"] = alias.NamespaceID

	if mountValidationResp := i.router.ValidateMountByAccessor(alias.MountAccessor); mountValidationResp != nil {
		respData["mount_path"] = mountValidationResp.MountPath
		respData["mount_type"] = mountValidationResp.MountType
	}

	// Convert protobuf timestamp into RFC3339 format
	respData["creation_time"] = ptypes.TimestampString(alias.CreationTime)
	respData["last_update_time"] = ptypes.TimestampString(alias.LastUpdateTime)

	return &logical.Response{
		Data: respData,
	}, nil
}

// pathAliasIDDelete deletes the alias for a given alias ID
func (i *IdentityStore) pathAliasIDDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		aliasID := d.Get("id").(string)
		if aliasID == "" {
			return logical.ErrorResponse("missing alias ID"), nil
		}

		i.lock.Lock()
		defer i.lock.Unlock()

		// Create a MemDB transaction to delete entity
		txn := i.db.Txn(true)
		defer txn.Abort()

		// Fetch the alias
		alias, err := i.MemDBAliasByIDInTxn(txn, aliasID, false, false)
		if err != nil {
			return nil, err
		}

		// If there is no alias for the ID, do nothing
		if alias == nil {
			return nil, nil
		}

		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		if ns.ID != alias.NamespaceID {
			return logical.ErrorResponse("request and alias are in different namespaces"), logical.ErrPermissionDenied
		}

		// Fetch the associated entity
		entity, err := i.MemDBEntityByAliasIDInTxn(txn, alias.ID, true)
		if err != nil {
			return nil, err
		}

		// If there is no entity tied to a valid alias, something is wrong
		if entity == nil {
			return nil, fmt.Errorf("alias not associated to an entity")
		}

		aliases := []*identity.Alias{
			alias,
		}

		// Delete alias from the entity object
		err = i.deleteAliasesInEntityInTxn(txn, entity, aliases)
		if err != nil {
			return nil, err
		}

		// Update the entity index in the entities table
		err = i.MemDBUpsertEntityInTxn(txn, entity)
		if err != nil {
			return nil, err
		}

		// Persist the entity object
		entityAsAny, err := ptypes.MarshalAny(entity)
		if err != nil {
			return nil, err
		}
		item := &storagepacker.Item{
			ID:      entity.ID,
			Message: entityAsAny,
		}

		err = i.entityPacker.PutItem(ctx, item)
		if err != nil {
			return nil, err
		}

		// Committing the transaction *after* successfully updating entity in
		// storage
		txn.Commit()

		return nil, nil
	}
}

// pathAliasIDList lists the IDs of all the valid aliases in the identity
// store
func (i *IdentityStore) pathAliasIDList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		return i.handleAliasListCommon(ctx, false)
	}
}

var aliasHelp = map[string][2]string{
	"alias": {
		"Create a new alias.",
		"",
	},
	"alias-id": {
		"Update, read or delete an alias ID.",
		"",
	},
	"alias-id-list": {
		"List all the alias IDs.",
		"",
	},
}
