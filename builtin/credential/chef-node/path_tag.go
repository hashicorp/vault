package chefnode

import (
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathTagsList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tags/?$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathTagList,
		},
		HelpSynopsis:    pathTagHelpSyn,
		HelpDescription: pathTagHelpDesc,
	}
}

func pathTags(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `tag/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the Chef tag",
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-seperated list of policies associated to this Chef tag",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathTagDelete,
			logical.ReadOperation:   b.pathTagRead,
			logical.UpdateOperation: b.pathTagWrite,
		},
	}
}

func (b *backend) pathTagList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	tags, err := req.Storage.List("tag/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(tags), nil
}

func (b *backend) Tag(s logical.Storage, n string) (*TagEntry, error) {
	entry, err := s.Get("tag/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result TagEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathTagDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("tag/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathTagRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	tag, err := b.Tag(req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": tag.Policies,
		},
	}, nil
}

func (b *backend) pathTagWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("tag/"+d.Get("name").(string), &TagEntry{
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

type TagEntry struct {
	Policies []string
}

const pathTagHelpSyn = `
Manage Vault policies assigned to a Chef tag
`
const pathTagHelpDesc = `
This endpoint allows you to create, read, update and delete configurations for policies
associated with a Chef tag.
`
