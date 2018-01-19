package nomad

import (
	"context"
	"errors"

	"github.com/hashicorp/errwrap"
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

			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma-separated string or list of policies as previously created in Nomad. Required for 'client' token.",
			},

			"global": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Boolean value describing if the token should be global or not. Defaults to false.",
			},

			"type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "client",
				Description: `Which type of token to create: 'client'
or 'management'. If a 'management' token,
the "policies" parameter is not required.
Defaults to 'client'.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRolesRead,
			logical.CreateOperation: b.pathRolesWrite,
			logical.UpdateOperation: b.pathRolesWrite,
			logical.DeleteOperation: b.pathRolesDelete,
		},

		ExistenceCheck: b.rolesExistenceCheck,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) rolesExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)
	entry, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) Role(ctx context.Context, storage logical.Storage, name string) (*roleConfig, error) {
	if name == "" {
		return nil, errors.New("invalid role name")
	}

	entry, err := storage.Get(ctx, "role/"+name)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving role: {{err}}", err)
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

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"type":     role.TokenType,
			"global":   role.Global,
			"policies": role.Policies,
		},
	}
	return resp, nil
}

func (b *backend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		role = new(roleConfig)
	}

	policies, ok := d.GetOk("policies")
	if ok {
		role.Policies = policies.([]string)
	}

	role.TokenType = d.Get("type").(string)
	switch role.TokenType {
	case "client":
		if len(role.Policies) == 0 {
			return logical.ErrorResponse(
				"policies cannot be empty when using client tokens"), nil
		}
	case "management":
		if len(role.Policies) != 0 {
			return logical.ErrorResponse(
				"policies should be empty when using management tokens"), nil
		}
	default:
		return logical.ErrorResponse(
			`type must be "client" or "management"`), nil
	}

	global, ok := d.GetOk("global")
	if ok {
		role.Global = global.(bool)
	}

	entry, err := logical.StorageEntryJSON("role/"+name, role)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if err := req.Storage.Delete(ctx, "role/"+name); err != nil {
		return nil, err
	}
	return nil, nil
}

type roleConfig struct {
	Policies  []string `json:"policies"`
	TokenType string   `json:"type"`
	Global    bool     `json:"global"`
}
