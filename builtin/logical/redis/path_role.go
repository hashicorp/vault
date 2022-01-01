package redis

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of this role",
			},

			"rules": {
				Type:        framework.TypeCommaStringSlice,
				Description: "The rules to set for this role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.roleWrite,
			logical.ReadOperation:   b.roleRead,
			logical.UpdateOperation: b.roleWrite,
			logical.DeleteOperation: b.roleDelete,
		},

		ExistenceCheck: b.roleExists,
	}
}

type Role struct {
	Rules []string
}

func getRole(ctx context.Context, s logical.Storage, name string) (*Role, error) {
	entry, err := s.Get(ctx, "role/"+name)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	if entry == nil {
		return nil, nil
	}

	var role Role
	if err := entry.DecodeJSON(&role); err != nil {
		return nil, fmt.Errorf("failed to decode role: %w", err)
	}

	return &role, nil
}

func (r *Role) Response() *logical.Response {
	if r == nil {
		return logical.ErrorResponse("No role found")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"rules": r.Rules,
		},
	}
}

func (b *backend) roleExists(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := req.Storage.Get(ctx, "role/"+data.Get("name").(string))
	return entry != nil, err
}

func (b *backend) roleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role := Role{
		Rules: data.Get("rules").([]string),
	}

	entry, err := logical.StorageEntryJSON("role/"+data.Get("name").(string), role)
	if err != nil {
		return logical.ErrorResponse("failed to marshal role: %s", err), nil
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return logical.ErrorResponse("failed to save role: %s", err), nil
	}

	return role.Response(), nil
}

func (b *backend) roleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := getRole(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	return role.Response(), nil
}

func (b *backend) roleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "role/"+data.Get("name").(string))
	if err != nil {
		return logical.ErrorResponse("failed to delete role: %s", err), nil
	}
	return nil, nil
}
