// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package rabbitmq

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixRabbitMQ,
			OperationSuffix: "roles",
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixRabbitMQ,
			OperationSuffix: "role",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"tags": {
				Type:        framework.TypeString,
				Description: "Comma-separated list of tags for this role.",
			},
			"vhosts": {
				Type:        framework.TypeString,
				Description: "A map of virtual hosts to permissions.",
			},
			"vhost_topics": {
				Type:        framework.TypeString,
				Description: "A nested map of virtual hosts and exchanges to topic permissions.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleUpdate,
			logical.DeleteOperation: b.pathRoleDelete,
		},
		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

// Reads the role configuration from the storage
func (b *backend) Role(ctx context.Context, s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get(ctx, "role/"+n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Deletes an existing role
func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	return nil, req.Storage.Delete(ctx, "role/"+name)
}

// Reads an existing role
func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(role).Map(),
	}, nil
}

// Lists all the roles registered with the backend
func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(roles), nil
}

// Registers a new role with the backend
func (b *backend) pathRoleUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	tags := d.Get("tags").(string)
	rawVHosts := d.Get("vhosts").(string)
	rawVHostTopics := d.Get("vhost_topics").(string)

	// Either tags or VHost permissions are always required, but topic permissions are always optional.
	if tags == "" && rawVHosts == "" {
		return logical.ErrorResponse("both tags and vhosts not specified"), nil
	}

	var vhosts map[string]vhostPermission
	if len(rawVHosts) > 0 {
		if err := jsonutil.DecodeJSON([]byte(rawVHosts), &vhosts); err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to unmarshal vhosts: %s", err)), nil
		}
	}

	var vhostTopics map[string]map[string]vhostTopicPermission
	if len(rawVHostTopics) > 0 {
		if err := jsonutil.DecodeJSON([]byte(rawVHostTopics), &vhostTopics); err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to unmarshal vhost_topics: %s", err)), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("role/"+name, &roleEntry{
		Tags:        tags,
		VHosts:      vhosts,
		VHostTopics: vhostTopics,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// Role that defines the capabilities of the credentials issued against it.
// Maps are used because the names of vhosts and exchanges will vary widely.
// VHosts is a map with a vhost name as key and the permissions as value.
// VHostTopics is a nested map with vhost name and exchange name as keys and
// the topic permissions as value.
type roleEntry struct {
	Tags        string                                     `json:"tags" structs:"tags" mapstructure:"tags"`
	VHosts      map[string]vhostPermission                 `json:"vhosts" structs:"vhosts" mapstructure:"vhosts"`
	VHostTopics map[string]map[string]vhostTopicPermission `json:"vhost_topics" structs:"vhost_topics" mapstructure:"vhost_topics"`
}

// Structure representing the permissions of a vhost
type vhostPermission struct {
	Configure string `json:"configure" structs:"configure" mapstructure:"configure"`
	Write     string `json:"write" structs:"write" mapstructure:"write"`
	Read      string `json:"read" structs:"read" mapstructure:"read"`
}

// Structure representing the topic permissions of an exchange
type vhostTopicPermission struct {
	Write string `json:"write" structs:"write" mapstructure:"write"`
	Read  string `json:"read" structs:"read" mapstructure:"read"`
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.

The "tags" parameter customizes the tags used to create the role.
This is a comma separated list of strings. The "vhosts" parameter customizes
the virtual hosts that this user will be associated with. This is a JSON object
passed as a string in the form:
{
	"vhostOne": {
		"configure": ".*",
		"write": ".*",
		"read": ".*"
	},
	"vhostTwo": {
		"configure": ".*",
		"write": ".*",
		"read": ".*"
	}
}
The "vhost_topics" parameter customizes the topic permissions that this user
will be granted. This is a JSON object passed as a string in the form:
{
	"vhostOne": {
		"exchangeOneOne": {
			"write": ".*",
			"read": ".*"
		},
		"exchangeOneTwo": {
			"write": ".*",
			"read": ".*"
		}
	},
	"vhostTwo": {
		"exchangeTwoOne": {
			"write": ".*",
			"read": ".*"
		}
	}
}
`
