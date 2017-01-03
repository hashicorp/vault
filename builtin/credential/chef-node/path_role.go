package chefNode

import (
	"strings"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRolesList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `role/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the Chef role",
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-seperated list of policies associated to this Chef role",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathRoleDelete,
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleWrite,
		},
	}
}

func (b *backend) pathRoleList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	envs, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(envs), nil
}

func (b *backend) Role(s logical.Storage, n string) (*RoleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result RoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("role/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	env, err := b.Role(req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if env == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": strings.Join(env.Policies, ","),
		},
	}, nil
}

func (b *backend) pathRoleWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("role/"+d.Get("name").(string), &RoleEntry{
		Policies: policyutil.ParsePolicies(d.Get("policies").(string)),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type RoleEntry struct {
	Policies []string
}
