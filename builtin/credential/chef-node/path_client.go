package chefnode

import (
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathClientsList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "clients/?$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathClientList,
		},
		HelpSynopsis:    pathClientHelpSyn,
		HelpDescription: pathClientHelpDesc,
	}
}

func pathClients(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `client/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the Chef client",
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-seperated list of policies associated to this Chef client",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathClientDelete,
			logical.ReadOperation:   b.pathClientRead,
			logical.UpdateOperation: b.pathClientWrite,
		},
	}
}

func (b *backend) pathClientList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	clients, err := req.Storage.List("client/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(clients), nil
}

func (b *backend) Client(s logical.Storage, n string) (*ClientEntry, error) {
	entry, err := s.Get("client/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result ClientEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathClientDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("client/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathClientRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, err := b.Client(req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": client.Policies,
		},
	}, nil
}

func (b *backend) pathClientWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("client/"+d.Get("name").(string), &ClientEntry{
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

type ClientEntry struct {
	Policies []string
}

const pathClientHelpSyn = `
Manage Vault policies assigned to a Chef client
`
const pathClientHelpDesc = `
This endpoint allows you to create, read, update, and delete configuration for policies
associated with Chef clients
`
