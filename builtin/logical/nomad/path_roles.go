package nomad

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"policy": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma separated list of policies as previously created in Nomad. Required",
			},

			"global": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "Boolean value describing if the token should be global or not. Defaults to false",
			},

			"type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "client",
				Description: `Which type of token to create: 'client'
or 'management'. If a 'management' token,
the "policy" parameter is not required.
Defaults to 'client'.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRolesRead,
			logical.UpdateOperation: b.pathRolesWrite,
			logical.DeleteOperation: b.pathRolesDelete,
		},
	}
}

func (b *backend) Role(storage logical.Storage, name string) (*roleConfig, error) {
	entry, err := storage.Get("role/" + name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %s", err)
	}
	if entry == nil {
		return nil, nil
	}

	var result roleConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRolesRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"type":   role.TokenType,
			"global": role.Global,
			"policy": role.Policy,
		},
	}
	return resp, nil
}

func (b *backend) pathRolesWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	tokenType := d.Get("type").(string)
	name := d.Get("name").(string)
	global := d.Get("global").(bool)
	policy := d.Get("policy").([]string)

	switch tokenType {
	case "client":
		if len(policy) == 0 {
			return logical.ErrorResponse(
				"policy cannot be empty when using client tokens"), nil
		}
	case "management":
		if len(policy) != 0 {
			return logical.ErrorResponse(
				"policy should be empty when using management tokens"), nil
		}
	default:
		return logical.ErrorResponse(
			"type must be \"client\" or \"management\""), nil
	}

	entry, err := logical.StorageEntryJSON("role/"+name, roleConfig{
		Policy:    policy,
		TokenType: tokenType,
		Global:    global,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRolesDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if err := req.Storage.Delete("role/" + name); err != nil {
		return nil, err
	}
	return nil, nil
}

type roleConfig struct {
	Policy    []string `json:"policy"`
	TokenType string   `json:"type"`
	Global    bool     `json:"global"`
}
