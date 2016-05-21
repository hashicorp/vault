package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func rolesFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"name": &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: "Name of the role.",
		},

		"tags": &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: "Comma-separated list of tags for this role.",
		},

		"vhosts": &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: "A map of virtual hosts to permissions.",
		},
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields:  rolesFields(),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.CreateOperation: b.pathRoleCreate,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *backend) Role(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
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

func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name, err := validateName(data)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Delete("role/" + name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name, err := validateName(data)
	if err != nil {
		return nil, err
	}

	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"tags":   role.Tags,
			"vhosts": role.VHosts,
		},
	}, nil
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name, err := validateName(data)
	if err != nil {
		return nil, err
	}
	tags := data.Get("tags").(string)
	rawVHosts := data.Get("vhosts").(string)

	var vhosts map[string]vhostPermission
	if len(rawVHosts) > 0 {
		err := json.Unmarshal([]byte(rawVHosts), &vhosts)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to unmarshal vhosts: %s", err)), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("role/"+name, &roleEntry{
		Tags:   tags,
		VHosts: vhosts,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleEntry struct {
	Tags   string                     `json:"tags"`
	VHosts map[string]vhostPermission `json:"vhosts"`
}

type vhostPermission struct {
	Configure string `json:"configure"`
	Write     string `json:"write"`
	Read      string `json:"read"`
}

func validateName(data *framework.FieldData) (string, error) {
	name := data.Get("name").(string)
	if len(name) == 0 {
		return "", errors.New("name is required")
	}

	return name, nil
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
`
