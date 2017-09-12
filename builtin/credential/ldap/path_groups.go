package ldap

import (
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathGroupsList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "groups/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathGroupList,
		},

		HelpSynopsis:    pathGroupHelpSyn,
		HelpDescription: pathGroupHelpDesc,
	}
}

func pathGroups(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `groups/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the LDAP group.",
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma-separated list of policies associated to the group.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathGroupDelete,
			logical.ReadOperation:   b.pathGroupRead,
			logical.UpdateOperation: b.pathGroupWrite,
		},

		HelpSynopsis:    pathGroupHelpSyn,
		HelpDescription: pathGroupHelpDesc,
	}
}

func (b *backend) Group(s logical.Storage, n string) (*GroupEntry, error) {
	entry, err := s.Get("group/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result GroupEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathGroupDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("group/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathGroupRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	group, err := b.Group(req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": group.Policies,
		},
	}, nil
}

func (b *backend) pathGroupWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Store it
	entry, err := logical.StorageEntryJSON("group/"+d.Get("name").(string), &GroupEntry{
		Policies: policyutil.ParsePolicies(d.Get("policies")),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathGroupList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	groups, err := req.Storage.List("group/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(groups), nil
}

type GroupEntry struct {
	Policies []string
}

const pathGroupHelpSyn = `
Manage users allowed to authenticate.
`

const pathGroupHelpDesc = `
This endpoint allows you to create, read, update, and delete configuration
for LDAP groups that are allowed to authenticate, and associate policies to
them.

Deleting a group will not revoke auth for prior authenticated users in that
group. To do this, do a revoke on "login/<username>" for
the usernames you want revoked.
`
