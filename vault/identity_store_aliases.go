// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/custommetadata"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// aliasPaths returns the API endpoints to operate on aliases.
// Following are the paths supported:
// entity-alias - To register/modify an alias
// entity-alias/id - To read, modify, delete and list aliases based on their ID
func aliasPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "entity-alias$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationVerb:   "create",
				OperationSuffix: "alias",
			},

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

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationSuffix: "alias-by-id",
			},

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

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.handleAliasCreateUpdate(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathAliasIDRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathAliasIDDelete(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias-id"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias-id"][1]),
		},
		{
			Pattern: "entity-alias/id/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationVerb:   "list",
				OperationSuffix: "aliases-by-id",
			},

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
			customMetadata = data.(map[string]string)
		}

		// Get entity id
		canonicalID := d.Get("canonical_id").(string)
		if canonicalID == "" {
			// For backwards compatibility
			canonicalID = d.Get("entity_id").(string)
		}

		// validate customMetadata if provided
		if len(customMetadata) != 0 {
			if err := custommetadata.Validate(customMetadata); err != nil {
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
				if !customMetadataExists {
					customMetadata = alias.CustomMetadata
				}
				switch {
				case mountAccessor == "" && name == "":
					// Check if the canonicalID or the customMetadata are being
					// updated
					if canonicalID == "" && !customMetadataExists {
						// Nothing to do, so be idempotent
						return nil, nil
					}
					name = alias.Name
					mountAccessor = alias.MountAccessor
				case mountAccessor == "":
					// No change to mount accessor
					mountAccessor = alias.MountAccessor
				case name == "":
					// No change to mount name
					name = alias.Name
				default:
					// mountAccessor, name and customMetadata  provided
				}
				return i.handleAliasUpdate(ctx, canonicalID, name, mountAccessor, alias, customMetadata)
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
		if mountEntry.NamespaceID != ns.ID {
			return logical.ErrorResponse("matching mount is in a different namespace than request"), logical.ErrPermissionDenied
		}
		alias, err := i.MemDBAliasByFactors(mountAccessor, name, true, false)
		if err != nil {
			return nil, err
		}
		if alias != nil {
			if alias.NamespaceID != ns.ID {
				return logical.ErrorResponse("cannot modify aliases across namespaces"), logical.ErrPermissionDenied
			}
			return i.handleAliasUpdate(ctx, canonicalID, name, mountAccessor, alias, customMetadata)
		}
		// At this point we know it's a new creation request
		return i.handleAliasCreate(ctx, canonicalID, name, mountAccessor, mountEntry.Local, customMetadata)
	}
}

func (i *IdentityStore) handleAliasCreate(ctx context.Context, canonicalID, name, mountAccessor string, local bool, customMetadata map[string]string) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var entity *identity.Entity
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

	if entity == nil && local {
		// Check to see if the entity creation should be forwarded.
		entity, err = i.entityCreator.CreateEntity(ctx)
		if err != nil {
			return nil, err
		}
	}

	persist := false
	// If the request was not forwarded, then this is the active node of the
	// primary. Create the entity here itself.
	if entity == nil {
		persist = true
		entity = new(identity.Entity)
		err = i.sanitizeEntity(ctx, entity)
		if err != nil {
			return nil, err
		}
	}

	for _, currentAlias := range entity.Aliases {
		if currentAlias.MountAccessor == mountAccessor {
			return logical.ErrorResponse("Alias already exists for requested entity and mount accessor"), nil
		}
	}

	var alias *identity.Alias
	switch local {
	case true:
		alias, err = i.processLocalAlias(ctx, &logical.Alias{
			MountAccessor:  mountAccessor,
			Name:           name,
			Local:          local,
			CustomMetadata: customMetadata,
		}, entity, false)
		if err != nil {
			return nil, err
		}
	default:
		alias = &identity.Alias{
			MountAccessor:  mountAccessor,
			Name:           name,
			CustomMetadata: customMetadata,
			CanonicalID:    entity.ID,
		}
		err = i.sanitizeAlias(ctx, alias)
		if err != nil {
			return nil, err
		}
		entity.UpsertAlias(alias)
		persist = true
	}

	// Index entity and its aliases in MemDB and persist entity along with
	// aliases in storage.
	if err := i.upsertEntity(ctx, entity, nil, persist); err != nil {
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

func (i *IdentityStore) handleAliasUpdate(ctx context.Context, canonicalID, name, mountAccessor string, alias *identity.Alias, customMetadata map[string]string) (*logical.Response, error) {
	if name == alias.Name &&
		mountAccessor == alias.MountAccessor &&
		(canonicalID == alias.CanonicalID || canonicalID == "") && (strutil.EqualStringMaps(customMetadata, alias.CustomMetadata)) {
		// Nothing to do; return nil to be idempotent
		return nil, nil
	}

	alias.LastUpdateTime = timestamppb.Now()

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

	mountValidationResp := i.router.ValidateMountByAccessor(alias.MountAccessor)
	if mountValidationResp == nil {
		return nil, fmt.Errorf("invalid mount accessor %q", alias.MountAccessor)
	}

	newEntity := currentEntity
	if canonicalID != "" && canonicalID != alias.CanonicalID {
		// Don't allow moving local aliases between entities.
		if mountValidationResp.MountLocal {
			return logical.ErrorResponse("local aliases can't be moved between entities"), nil
		}

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

	if mountValidationResp.MountLocal {
		alias, err = i.processLocalAlias(ctx, &logical.Alias{
			MountAccessor:  mountAccessor,
			Name:           name,
			Local:          mountValidationResp.MountLocal,
			CustomMetadata: customMetadata,
		}, newEntity, true)
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"id":           alias.ID,
				"canonical_id": newEntity.ID,
			},
		}, nil
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
	respData["local"] = alias.Local

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

		switch alias.Local {
		case true:
			localAliases, err := i.parseLocalAliases(entity.ID)
			if err != nil {
				return nil, err
			}

			if localAliases == nil {
				return nil, nil
			}

			for i, item := range localAliases.Aliases {
				if item.ID == alias.ID {
					localAliases.Aliases = append(localAliases.Aliases[:i], localAliases.Aliases[i+1:]...)
					break
				}
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
		default:
			if err := i.persistEntity(ctx, entity); err != nil {
				return nil, err
			}
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
