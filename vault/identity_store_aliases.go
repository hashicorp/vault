package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

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
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleAliasUpdateCommon(),
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
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleAliasUpdateCommon(),
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

// handleAliasUpdateCommon is used to update an alias
func (i *IdentityStore) handleAliasUpdateCommon() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		var err error
		var alias *identity.Alias
		var entity *identity.Entity
		var previousEntity *identity.Entity

		i.lock.Lock()
		defer i.lock.Unlock()

		// Check for update or create
		aliasID := d.Get("id").(string)
		if aliasID != "" {
			alias, err = i.MemDBAliasByID(aliasID, true, false)
			if err != nil {
				return nil, err
			}
			if alias == nil {
				return logical.ErrorResponse("invalid alias id"), nil
			}
		} else {
			alias = &identity.Alias{}
		}

		// Get entity id
		canonicalID := d.Get("canonical_id").(string)
		if canonicalID == "" {
			// For backwards compatibility
			canonicalID = d.Get("entity_id").(string)
		}

		// Get alias name
		if aliasName := d.Get("name").(string); aliasName == "" {
			if alias.Name == "" {
				return logical.ErrorResponse("missing alias name"), nil
			}
		} else {
			alias.Name = aliasName
		}

		// Get mount accessor
		if mountAccessor := d.Get("mount_accessor").(string); mountAccessor == "" {
			if alias.MountAccessor == "" {
				return logical.ErrorResponse("missing mount_accessor"), nil
			}
		} else {
			alias.MountAccessor = mountAccessor
		}

		mountValidationResp := i.core.router.validateMountByAccessor(alias.MountAccessor)
		if mountValidationResp == nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid mount accessor %q", alias.MountAccessor)), nil
		}
		if mountValidationResp.MountLocal {
			return logical.ErrorResponse(fmt.Sprintf("mount_accessor %q is of a local mount", alias.MountAccessor)), nil
		}

		// Verify that the combination of alias name and mount is not
		// already tied to a different alias
		aliasByFactors, err := i.MemDBAliasByFactors(mountValidationResp.MountAccessor, alias.Name, false, false)
		if err != nil {
			return nil, err
		}
		if aliasByFactors != nil {
			// If it's a create we won't have an alias ID so this will correctly
			// bail. If it's an update alias will be the same as aliasbyfactors so
			// we don't need to transfer any info over
			if aliasByFactors.ID != alias.ID {
				return logical.ErrorResponse("combination of mount and alias name is already in use"), nil
			}

			// Fetch the entity to which the alias is tied. We don't need to append
			// here, so the only further checking is whether the canonical ID is
			// different
			entity, err = i.MemDBEntityByAliasID(alias.ID, true)
			if err != nil {
				return nil, err
			}
			if entity == nil {
				return nil, fmt.Errorf("existing alias is not associated with an entity")
			}
		} else if alias.ID != "" {
			// This is an update, not a create; if we have an associated entity
			// already, load it
			entity, err = i.MemDBEntityByAliasID(alias.ID, true)
			if err != nil {
				return nil, err
			}
		}

		resp := &logical.Response{}

		// If we found an existing alias we won't hit this condition because
		// canonicalID being empty will result in nil being returned in the block
		// above, so in this case we know that creating a new entity is the right
		// thing.
		if canonicalID == "" {
			entity = &identity.Entity{
				Aliases: []*identity.Alias{
					alias,
				},
			}
		} else {
			// If we can look up by the given canonical ID, see if this is a
			// transfer; otherwise if we found no previous entity but we find one
			// here, use it.
			canonicalEntity, err := i.MemDBEntityByID(canonicalID, true)
			if err != nil {
				return nil, err
			}
			if canonicalEntity == nil {
				return logical.ErrorResponse("invalid canonical ID"), nil
			}
			if entity == nil {
				// If entity is nil, we didn't find a previous alias from factors,
				// so append to this entity
				entity = canonicalEntity
				entity.Aliases = append(entity.Aliases, alias)
			} else if entity.ID != canonicalEntity.ID {
				// In this case we found an entity from alias factors or given
				// alias ID but it's not the same, so it's a migration
				previousEntity = entity
				entity = canonicalEntity

				for aliasIndex, item := range previousEntity.Aliases {
					if item.ID == alias.ID {
						previousEntity.Aliases = append(previousEntity.Aliases[:aliasIndex], previousEntity.Aliases[aliasIndex+1:]...)
						break
					}
				}

				entity.Aliases = append(entity.Aliases, alias)
				resp.AddWarning(fmt.Sprintf("alias is being transferred from entity %q to %q", previousEntity.ID, entity.ID))
			}
		}

		// ID creation and other validations; This is more useful for new entities
		// and may not perform anything for the existing entities. Placing the
		// check here to make the flow common for both new and existing entities.
		err = i.sanitizeEntity(ctx, entity)
		if err != nil {
			return nil, err
		}

		// Explicitly set to empty as in the past we incorrectly saved it
		alias.MountPath = ""
		alias.MountType = ""

		// Set the canonical ID in the alias index. This should be done after
		// sanitizing entity.
		alias.CanonicalID = entity.ID

		// ID creation and other validations
		err = i.sanitizeAlias(ctx, alias)
		if err != nil {
			return nil, err
		}

		for index, item := range entity.Aliases {
			if item.ID == alias.ID {
				entity.Aliases[index] = alias
			}
		}

		// Index entity and its aliases in MemDB and persist entity along with
		// aliases in storage. If the alias is being transferred over from
		// one entity to another, previous entity needs to get refreshed in MemDB
		// and persisted in storage as well.
		if err := i.upsertEntity(ctx, entity, previousEntity, true); err != nil {
			return nil, err
		}

		// Return ID of both alias and entity
		resp.Data = map[string]interface{}{
			"id":           alias.ID,
			"canonical_id": entity.ID,
		}

		return resp, nil
	}
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
		return nil, nil
	}

	respData := map[string]interface{}{}
	respData["id"] = alias.ID
	respData["canonical_id"] = alias.CanonicalID
	respData["mount_accessor"] = alias.MountAccessor
	respData["metadata"] = alias.Metadata
	respData["name"] = alias.Name
	respData["merged_from_canonical_ids"] = alias.MergedFromCanonicalIDs

	if mountValidationResp := i.core.router.validateMountByAccessor(alias.MountAccessor); mountValidationResp != nil {
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
			return nil, logical.ErrUnsupportedPath
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

		err = i.entityPacker.PutItem(item)
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
