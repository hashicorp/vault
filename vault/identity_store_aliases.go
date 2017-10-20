package vault

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// aliasPaths returns the API endpoints to operate on aliases.
// Following are the paths supported:
// alias - To register/modify a alias
// alias/id - To lookup, delete and list aliases based on ID
func aliasPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "alias$",
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the alias",
				},
				"entity_id": {
					Type:        framework.TypeString,
					Description: "Entity ID to which this alias belongs to",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Mount accessor to which this alias belongs to",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the alias",
				},
				"metadata": {
					Type:        framework.TypeStringSlice,
					Description: "Metadata to be associated with the alias. Format should be a list of `key=value` pairs.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathAliasRegister,
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias"][1]),
		},
		{
			Pattern: "alias/id/" + framework.GenericNameRegex("id"),
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the alias",
				},
				"entity_id": {
					Type:        framework.TypeString,
					Description: "Entity ID to which this alias should be tied to",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Mount accessor to which this alias belongs to",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the alias",
				},
				"metadata": {
					Type:        framework.TypeStringSlice,
					Description: "Metadata to be associated with the alias. Format should be a comma separated list of `key=value` pairs.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathAliasIDUpdate,
				logical.ReadOperation:   i.pathAliasIDRead,
				logical.DeleteOperation: i.pathAliasIDDelete,
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias-id"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias-id"][1]),
		},
		{
			Pattern: "alias/id/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathAliasIDList,
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias-id-list"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias-id-list"][1]),
		},
	}
}

// pathAliasRegister is used to register new alias
func (i *IdentityStore) pathAliasRegister(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	_, ok := d.GetOk("id")
	if ok {
		return i.pathAliasIDUpdate(req, d)
	}

	return i.handleAliasUpdateCommon(req, d, nil)
}

// pathAliasIDUpdate is used to update a alias based on the given
// alias ID
func (i *IdentityStore) pathAliasIDUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get alias id
	aliasID := d.Get("id").(string)

	if aliasID == "" {
		return logical.ErrorResponse("missing alias ID"), nil
	}

	alias, err := i.memDBAliasByID(aliasID, true)
	if err != nil {
		return nil, err
	}
	if alias == nil {
		return logical.ErrorResponse("invalid alias id"), nil
	}

	return i.handleAliasUpdateCommon(req, d, alias)
}

// handleAliasUpdateCommon is used to update a alias
func (i *IdentityStore) handleAliasUpdateCommon(req *logical.Request, d *framework.FieldData, alias *identity.Alias) (*logical.Response, error) {
	var err error
	var newAlias bool
	var entity *identity.Entity
	var previousEntity *identity.Entity

	// Alias will be nil when a new alias is being registered; create a
	// new struct in that case.
	if alias == nil {
		alias = &identity.Alias{}
		newAlias = true
	}

	// Get entity id
	entityID := d.Get("entity_id").(string)
	if entityID != "" {
		entity, err = i.memDBEntityByID(entityID, true)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			return logical.ErrorResponse("invalid entity ID"), nil
		}
	}

	// Get alias name
	aliasName := d.Get("name").(string)
	if aliasName == "" {
		return logical.ErrorResponse("missing alias name"), nil
	}

	mountAccessor := d.Get("mount_accessor").(string)
	if mountAccessor == "" {
		return logical.ErrorResponse("missing mount_accessor"), nil
	}

	mountValidationResp := i.validateMountAccessorFunc(mountAccessor)
	if mountValidationResp == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid mount accessor %q", mountAccessor)), nil
	}

	// Get alias metadata

	// Accept metadata in the form of map[string]string to be able to index on
	// it
	var aliasMetadata map[string]string
	aliasMetadataRaw, ok := d.GetOk("metadata")
	if ok {
		aliasMetadata, err = parseMetadata(aliasMetadataRaw.([]string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse alias metadata: %v", err)), nil
		}
	}

	aliasByFactors, err := i.memDBAliasByFactors(mountValidationResp.MountAccessor, aliasName, false)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{}

	if newAlias {
		if aliasByFactors != nil {
			return logical.ErrorResponse("combination of mount and alias name is already in use"), nil
		}

		// If this is a alias being tied to a non-existent entity, create
		// a new entity for it.
		if entity == nil {
			entity = &identity.Entity{
				Aliases: []*identity.Alias{
					alias,
				},
			}
		} else {
			entity.Aliases = append(entity.Aliases, alias)
		}
	} else {
		// Verify that the combination of alias name and mount is not
		// already tied to a different alias
		if aliasByFactors != nil && aliasByFactors.ID != alias.ID {
			return logical.ErrorResponse("combination of mount and alias name is already in use"), nil
		}

		// Fetch the entity to which the alias is tied to
		existingEntity, err := i.memDBEntityByAliasID(alias.ID, true)
		if err != nil {
			return nil, err
		}

		if existingEntity == nil {
			return nil, fmt.Errorf("alias is not associated with an entity")
		}

		if entity != nil && entity.ID != existingEntity.ID {
			// Alias should be transferred from 'existingEntity' to 'entity'
			err = i.deleteAliasFromEntity(existingEntity, alias)
			if err != nil {
				return nil, err
			}
			previousEntity = existingEntity
			entity.Aliases = append(entity.Aliases, alias)
			resp.AddWarning(fmt.Sprintf("alias is being transferred from entity %q to %q", existingEntity.ID, entity.ID))
		} else {
			// Update entity with modified alias
			err = i.updateAliasInEntity(existingEntity, alias)
			if err != nil {
				return nil, err
			}
			entity = existingEntity
		}
	}

	// ID creation and other validations; This is more useful for new entities
	// and may not perform anything for the existing entities. Placing the
	// check here to make the flow common for both new and existing entities.
	err = i.sanitizeEntity(entity)
	if err != nil {
		return nil, err
	}

	// Update the fields
	alias.Name = aliasName
	alias.Metadata = aliasMetadata
	alias.MountType = mountValidationResp.MountType
	alias.MountAccessor = mountValidationResp.MountAccessor
	alias.MountPath = mountValidationResp.MountPath

	// Set the entity ID in the alias index. This should be done after
	// sanitizing entity.
	alias.EntityID = entity.ID

	// ID creation and other validations
	err = i.sanitizeAlias(alias)
	if err != nil {
		return nil, err
	}

	// Index entity and its aliases in MemDB and persist entity along with
	// aliases in storage. If the alias is being transferred over from
	// one entity to another, previous entity needs to get refreshed in MemDB
	// and persisted in storage as well.
	err = i.upsertEntity(entity, previousEntity, true)
	if err != nil {
		return nil, err
	}

	// Return ID of both alias and entity
	resp.Data = map[string]interface{}{
		"id":        alias.ID,
		"entity_id": entity.ID,
	}

	return resp, nil
}

// pathAliasIDRead returns the properties of a alias for a given
// alias ID
func (i *IdentityStore) pathAliasIDRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	aliasID := d.Get("id").(string)
	if aliasID == "" {
		return logical.ErrorResponse("missing alias id"), nil
	}

	alias, err := i.memDBAliasByID(aliasID, false)
	if err != nil {
		return nil, err
	}

	if alias == nil {
		return nil, nil
	}

	respData := map[string]interface{}{}
	respData["id"] = alias.ID
	respData["entity_id"] = alias.EntityID
	respData["mount_type"] = alias.MountType
	respData["mount_accessor"] = alias.MountAccessor
	respData["mount_path"] = alias.MountPath
	respData["metadata"] = alias.Metadata
	respData["name"] = alias.Name
	respData["merged_from_entity_ids"] = alias.MergedFromEntityIDs

	// Convert protobuf timestamp into RFC3339 format
	respData["creation_time"] = ptypes.TimestampString(alias.CreationTime)
	respData["last_update_time"] = ptypes.TimestampString(alias.LastUpdateTime)

	return &logical.Response{
		Data: respData,
	}, nil
}

// pathAliasIDDelete deleted the alias for a given alias ID
func (i *IdentityStore) pathAliasIDDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	aliasID := d.Get("id").(string)
	if aliasID == "" {
		return logical.ErrorResponse("missing alias ID"), nil
	}

	return nil, i.deleteAlias(aliasID)
}

// pathAliasIDList lists the IDs of all the valid aliases in the identity
// store
func (i *IdentityStore) pathAliasIDList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ws := memdb.NewWatchSet()
	iter, err := i.memDBAliases(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch iterator for aliases in memdb: %v", err)
	}

	var aliasIDs []string
	for {
		raw := iter.Next()
		if raw == nil {
			break
		}
		aliasIDs = append(aliasIDs, raw.(*identity.Alias).ID)
	}

	return logical.ListResponse(aliasIDs), nil
}

var aliasHelp = map[string][2]string{
	"alias": {
		"Create a new alias",
		"",
	},
	"alias-id": {
		"Update, read or delete an entity using alias ID",
		"",
	},
	"alias-id-list": {
		"List all the entity IDs",
		"",
	},
}
