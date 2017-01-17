package chefnode

import (
	"strings"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathEnvironmentsList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "environments/?$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathEnvironmentList,
		},
	}
}

func pathEnvironments(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `environment/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the Chef environment",
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-seperated list of policies associated to this Chef environment",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathEnvironmentDelete,
			logical.ReadOperation:   b.pathEnvironmentRead,
			logical.UpdateOperation: b.pathEnvironmentWrite,
		},
	}
}

func (b *backend) pathEnvironmentList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	envs, err := req.Storage.List("environment/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(envs), nil
}

func (b *backend) Environment(s logical.Storage, n string) (*EnvironmentEntry, error) {
	entry, err := s.Get("environment/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result EnvironmentEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathEnvironmentDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("environment/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathEnvironmentRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	env, err := b.Environment(req.Storage, d.Get("name").(string))
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

func (b *backend) pathEnvironmentWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("environment/"+d.Get("name").(string), &EnvironmentEntry{
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

type EnvironmentEntry struct {
	Policies []string
}
