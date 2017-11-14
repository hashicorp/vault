package vault

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func upgradePaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "persona$",
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
					Type:        framework.TypeKVPairs,
					Description: "Metadata to be associated with the alias.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathAliasRegister,
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias"][1]),
		},
		{
			Pattern: "persona/id/" + framework.GenericNameRegex("id"),
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
					Type:        framework.TypeKVPairs,
					Description: "Metadata to be associated with the alias.",
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
			Pattern: "persona/id/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathAliasIDList,
			},

			HelpSynopsis:    strings.TrimSpace(aliasHelp["alias-id-list"][0]),
			HelpDescription: strings.TrimSpace(aliasHelp["alias-id-list"][1]),
		},
		{
			Pattern: "alias$",
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the alias",
				},
				"entity_id": {
					Type:        framework.TypeString,
					Description: "Entity ID to which this alias belongs to. This field is deprecated in favor of 'canonical_id'.",
				},
				"canonical_id": {
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
					Type:        framework.TypeKVPairs,
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
					Description: "Entity ID to which this alias should be tied to. This field is deprecated in favor of 'canonical_id'.",
				},
				"canonical_id": {
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
					Type:        framework.TypeKVPairs,
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
