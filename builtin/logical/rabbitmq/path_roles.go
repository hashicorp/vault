// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func pathListRoles(b *backend) []*framework.Path {
	return []*framework.Path{
		{
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
		},
		{
			Pattern: "static-roles/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixRabbitMQ,
				OperationVerb:   "list",
				OperationSuffix: "static-roles",
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList,
			},

			HelpSynopsis:    pathStaticRoleHelpSyn,
			HelpDescription: pathStaticRoleHelpDesc,
		},
	}
}

func pathRoles(b *backend) []*framework.Path {
	return []*framework.Path{
		{
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
				logical.CreateOperation: b.pathRoleUpdate,
				logical.UpdateOperation: b.pathRoleUpdate,
				logical.DeleteOperation: b.pathRoleDelete,
			},
			HelpSynopsis:    pathRoleHelpSyn,
			HelpDescription: pathRoleHelpDesc,
		},
		{
			Pattern: "static-roles/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixRabbitMQ,
				OperationSuffix: "static-roles",
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
				"username": {
					Type:        framework.TypeString,
					Description: "A username to create this static roles as.",
				},
				"revoke_user_on_delete": {
					Type:        framework.TypeBool,
					Description: "Whether to revoke the user associated to this role when the role is deleted.",
				},
				"rotation_period": {
					Type:        framework.TypeDurationSecond,
					Description: "Period for automatic credential rotation of the given username.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathStaticRoleRead,
				logical.CreateOperation: b.pathStaticRoleCreateUpdate,
				logical.UpdateOperation: b.pathStaticRoleCreateUpdate,
				logical.DeleteOperation: b.pathStaticRoleDelete,
			},

			HelpSynopsis:    pathStaticRoleHelpSyn,
			HelpDescription: pathStaticRoleHelpDesc,
		},
	}
}

// Reads the role configuration from the storage
func (b *backend) Role(ctx context.Context, s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get(ctx, rabbitMQRolePath+n)
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

// Reads the static role configuration from the storage
func (b *backend) StaticRole(ctx context.Context, s logical.Storage, n string) (*staticRoleEntry, error) {
	entry, err := s.Get(ctx, rabbitMQStaticRolePath+n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result staticRoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Deletes an existing role
func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return b.pathRoleDeleteWithPrefix(ctx, rabbitMQRolePath, req, d)
}

func (b *backend) pathRoleDeleteWithPrefix(ctx context.Context, prefix string, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	return nil, req.Storage.Delete(ctx, prefix+name)
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
	path := rabbitMQRolePath
	if strings.HasPrefix(req.Path, "static-roles") {
		path = rabbitMQStaticRolePath
	}
	entries, err := req.Storage.List(ctx, path)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

// Registers a new role with the backend
func (b *backend) pathRoleUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.parseRoleEntryFromRequest(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	entry, err := logical.StorageEntryJSON(rabbitMQRolePath+name, role)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) parseRoleEntryFromRequest(d *framework.FieldData) (*roleEntry, error) {
	tags := d.Get("tags").(string)
	rawVHosts := d.Get("vhosts").(string)
	rawVHostTopics := d.Get("vhost_topics").(string)

	// Either tags or VHost permissions are always required, but topic permissions are always optional.
	if tags == "" && rawVHosts == "" {
		return nil, fmt.Errorf("both tags and vhosts not specified")
	}

	var vhosts map[string]vhostPermission
	if len(rawVHosts) > 0 {
		if err := jsonutil.DecodeJSON([]byte(rawVHosts), &vhosts); err != nil {
			return nil, fmt.Errorf("failed to unmarshal vhosts: %s", err)
		}
	}

	var vhostTopics map[string]map[string]vhostTopicPermission
	if len(rawVHostTopics) > 0 {
		if err := jsonutil.DecodeJSON([]byte(rawVHostTopics), &vhostTopics); err != nil {
			return nil, fmt.Errorf("failed to unmarshal vhost_topics: %s", err)
		}
	}

	entry := roleEntry{
		Tags:        tags,
		VHosts:      vhosts,
		VHostTopics: vhostTopics,
	}
	return &entry, nil
}

func (b *backend) pathStaticRoleRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.StaticRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	data := map[string]interface{}{
		"rotation_period":       role.RotationPeriod.Seconds(),
		"revoke_user_on_delete": role.RevokeUserOnDelete,
		"username":              role.Username,
	}
	roleInfo := structs.New(role.RoleEntry).Map()
	for k, v := range roleInfo {
		data[k] = v
	}
	if !role.LastVaultRotation.IsZero() {
		data["last_vault_rotation"] = role.LastVaultRotation
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) pathStaticRoleCreateUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.parseRoleEntryFromRequest(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	username := d.Get("username").(string)
	if username == "" {
		return logical.ErrorResponse("username is a required field to create a static account"), nil
	}
	revokeUserOnDelete := d.Get("revoke_user_on_delete").(bool)
	rotationPeriodRaw, ok := d.GetOk("rotation_period")
	if !ok {
		// TODO allow to set default in backend
		return logical.ErrorResponse("missing %q parameter", "rotation_period"), nil
	}
	rotationPeriod := rotationPeriodRaw.(int)
	if err := b.validateRotationPeriod(time.Duration(rotationPeriod) * time.Second); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	staticRole := staticRoleEntry{
		RoleEntry:          *role,
		Username:           username,
		RotationPeriod:     time.Duration(rotationPeriod) * time.Second,
		RevokeUserOnDelete: revokeUserOnDelete,
	}

	_, err = b.credRotationQueue.PopByKey(name)
	if err != nil {
		return nil, fmt.Errorf("failed to pop item from the rotation queue for role %q: %w", name, err)
	}
	// regenerate either way as username or permissions might have changed
	if err := b.createStaticCredential(ctx, req.Storage, &staticRole, name); err != nil {
		return nil, err
	}
	priority := time.Now().Add(staticRole.RotationPeriod).Unix()
	err = b.credRotationQueue.Push(&queue.Item{
		Key:      name,
		Value:    staticRole,
		Priority: priority,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add item into the rotation queue for role %q: %w", name, err)
	}

	return nil, nil
}

func (b *backend) pathStaticRoleDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}
	item, err := b.credRotationQueue.PopByKey(name)
	if err != nil {
		return nil, fmt.Errorf("failed to pop item from the rotation queue for role %q: %w", name, err)
	}
	if item == nil {
		return nil, fmt.Errorf("failed to find item on the rotation queue for role %q: %w", name, err)
	}
	role := item.Value.(staticRoleEntry)
	if err := b.deleteStaticCredential(ctx, req.Storage, role, false); err != nil {
		return nil, err
	}
	return b.pathRoleDeleteWithPrefix(ctx, rabbitMQStaticRolePath, req, d)
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

type staticRoleEntry struct {
	Username           string        `json:"username" structs:"username" mapstructure:"username"`
	Password           string        `json:"password" structs:"password" mapstructure:"password"`
	LastVaultRotation  time.Time     `json:"last_vault_rotation" structs:"last_vault_rotation" mapstructure:"last_vault_rotation"`
	RotationPeriod     time.Duration `json:"rotation_period" structs:"rotation_period" mapstructure:"rotation_period"`
	RevokeUserOnDelete bool          `json:"revoke_user_on_delete" structs:"revoke_user_on_delete" mapstructure:"revoke_user_on_delete"`
	RoleEntry          roleEntry     `json:"role_entry" structs:"role_entry" mapstructure:"role_entry"`
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

const pathStaticRoleHelpSyn = `
Manage the static roles that can be created with this backend.
`

const pathStaticRoleHelpDesc = `
This path lets you manage the static roles that can be created with this
backend. Static Roles are associated with a single RabbitMQ user, and manage the
credential based on a rotation period, automatically rotating the credential.

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
