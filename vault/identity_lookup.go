package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func lookupPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "lookup/entity$",
			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Description: "Type of lookup. Current supported values are 'id' and 'name'.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the entity.",
				},
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the entity.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathLookupEntityUpdate,
			},

			HelpSynopsis:    strings.TrimSpace(lookupHelp["lookup-entity"][0]),
			HelpDescription: strings.TrimSpace(lookupHelp["lookup-entity"][1]),
		},
		{
			Pattern: "lookup/group$",
			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Description: "Type of lookup. Current supported values are 'id' and 'name'",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the group.",
				},
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathLookupGroupUpdate,
			},

			HelpSynopsis:    strings.TrimSpace(lookupHelp["lookup-group"][0]),
			HelpDescription: strings.TrimSpace(lookupHelp["lookup-group"][1]),
		},
		{
			Pattern: "lookup/entity-alias$",
			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Description: "Type of lookup. Current supported values are 'id', 'canonical_id' and 'factors'.",
				},
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the entity.",
				},
				"canonical_id": {
					Type:        framework.TypeString,
					Description: "ID of the entity to which the alias belongs to.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the alias.",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Accessor of the mount to which the entity alias belongs to.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathLookupEntityAliasUpdate,
			},

			HelpSynopsis:    strings.TrimSpace(lookupHelp["lookup-entity-alias"][0]),
			HelpDescription: strings.TrimSpace(lookupHelp["lookup-entity-alias"][1]),
		},
		{
			Pattern: "lookup/group-alias$",
			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Description: "Type of lookup. Current supported values are 'id', 'canonical_id' and 'factors'.",
				},
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group.",
				},
				"canonical_id": {
					Type:        framework.TypeString,
					Description: "ID of the group to which the alias belongs to.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the alias.",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Accessor of the mount to which the group alias belongs to.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathLookupGroupAliasUpdate,
			},

			HelpSynopsis:    strings.TrimSpace(lookupHelp["lookup-group-alias"][0]),
			HelpDescription: strings.TrimSpace(lookupHelp["lookup-group-alias"][1]),
		},
	}
}

func (i *IdentityStore) pathLookupEntityUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lookupType := d.Get("type").(string)
	if lookupType == "" {
		return logical.ErrorResponse("empty type"), nil
	}

	var entity *identity.Entity
	var err error

	switch lookupType {
	case "id":
		entityID := d.Get("id").(string)
		if entityID == "" {
			return logical.ErrorResponse("empty id"), nil
		}
		entity, err = i.MemDBEntityByID(entityID, false)
		if err != nil {
			return nil, err
		}

	case "name":
		entityName := d.Get("name").(string)
		if entityName == "" {
			return logical.ErrorResponse("empty name"), nil
		}
		entity, err = i.MemDBEntityByName(entityName, false)
		if err != nil {
			return nil, err
		}

	default:
		return logical.ErrorResponse(fmt.Sprintf("unrecognized type %q", lookupType)), nil
	}

	if entity == nil {
		return nil, nil
	}

	return i.handleEntityReadCommon(entity)
}

func (i *IdentityStore) pathLookupGroupUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lookupType := d.Get("type").(string)
	if lookupType == "" {
		return logical.ErrorResponse("empty type"), nil
	}

	switch lookupType {
	case "id":
		groupID := d.Get("id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty ID"), nil
		}
		group, err := i.MemDBGroupByID(groupID, false)
		if err != nil {
			return nil, err
		}
		return i.handleGroupReadCommon(group)
	case "name":
		groupName := d.Get("name").(string)
		if groupName == "" {
			return logical.ErrorResponse("empty name"), nil
		}
		group, err := i.MemDBGroupByName(groupName, false)
		if err != nil {
			return nil, err
		}
		return i.handleGroupReadCommon(group)
	default:
		return logical.ErrorResponse(fmt.Sprintf("unrecognized type %q", lookupType)), nil
	}

	return nil, nil
}

func (i *IdentityStore) pathLookupEntityAliasUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return i.handleLookupAliasUpdateCommon(req, d, false)
}

func (i *IdentityStore) pathLookupGroupAliasUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return i.handleLookupAliasUpdateCommon(req, d, true)
}

func (i *IdentityStore) handleLookupAliasUpdateCommon(req *logical.Request, d *framework.FieldData, groupAlias bool) (*logical.Response, error) {
	lookupType := d.Get("type").(string)
	if lookupType == "" {
		return logical.ErrorResponse("empty type"), nil
	}

	switch lookupType {
	case "id":
		aliasID := d.Get("id").(string)
		if aliasID == "" {
			return logical.ErrorResponse("empty ID"), nil
		}

		alias, err := i.MemDBAliasByID(aliasID, false, groupAlias)
		if err != nil {
			return nil, err
		}

		return i.handleAliasReadCommon(alias)

	case "canonical_id":
		canonicalID := d.Get("canonical_id").(string)
		if canonicalID == "" {
			return logical.ErrorResponse("empty canonical_id"), nil
		}

		alias, err := i.MemDBAliasByCanonicalID(canonicalID, false, groupAlias)
		if err != nil {
			return nil, err
		}

		return i.handleAliasReadCommon(alias)

	case "factors":
		aliasName := d.Get("name").(string)
		if aliasName == "" {
			return logical.ErrorResponse("empty name"), nil
		}
		mountAccessor := d.Get("mount_accessor").(string)
		if mountAccessor == "" {
			return logical.ErrorResponse("empty 'mount_accessor'"), nil
		}

		alias, err := i.MemDBAliasByFactors(mountAccessor, aliasName, false, groupAlias)
		if err != nil {
			return nil, err
		}

		return i.handleAliasReadCommon(alias)

	default:
		return logical.ErrorResponse(fmt.Sprintf("unrecognized type %q", lookupType)), nil
	}
}

var lookupHelp = map[string][2]string{
	"lookup-entity": {
		"Query entities based on types.",
		`Supported types:
		- 'id'
		To query the entity by its ID. This requires 'id' parameter to be set.
		- 'name'
		To query the entity by its name. This requires 'name' parameter to be set.
		`,
	},
	"lookup-group": {
		"Query groups based on types.",
		`Supported types:
		- 'id'
		To query the group by its ID. This requires 'id' parameter to be set.
		- 'name'
		To query the group by its name. This requires 'name' parameter to be set.
		`,
	},
	"lookup-group-alias": {
		"Query group alias based on types.",
		`Supported types:
		- 'id'
		To query the group alias by its ID. This requires 'id' parameter to be set.
		- 'canonical_id'
		To query the group alias by the ID of the group it belongs to. This requires the 'canonical_id' parameter to be set.
		- 'factors'
		To query the group alias using the factors that uniquely identifies a group alias; its name and the mount accessor. This requires the 'name' and 'mount_accessor' parameters to be set.
		`,
	},
	"lookup-entity-alias": {
		"Query entity alias based on types.",
		`Supported types:
		- 'id'
		To query the entity alias by its ID. This requires 'id' parameter to be set.
		- 'canonical_id'
		To query the entity alias by the ID of the entity it belongs to. This requires the 'canonical_id' parameter to be set.
		- 'factors'
		To query the entity alias using the factors that uniquely identifies an entity alias; its name and the mount accessor. This requires the 'name' and 'mount_accessor' parameters to be set.
		`,
	},
}
