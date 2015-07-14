package ldap

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUsers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `users/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the LDAP user.",
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-separated list of policies associated to the user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathUserDelete,
			logical.ReadOperation:   b.pathUserRead,
			logical.WriteOperation:  b.pathUserWrite,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func (b *backend) User(s logical.Storage, n string) (*UserEntry, error) {
	entry, err := s.Get("user/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result UserEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathUserDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("user/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathUserRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	user, err := b.User(req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": strings.Join(user.Policies, ","),
		},
	}, nil
}

func (b *backend) pathUserWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	policies := strings.Split(d.Get("policies").(string), ",")
	for i, p := range policies {
		policies[i] = strings.TrimSpace(p)
	}

	// Store it
	entry, err := logical.StorageEntryJSON("user/"+name, &UserEntry{
		Policies: policies,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type UserEntry struct {
	Policies []string
}

const pathUserHelpSyn = `
Manage users allowed to authenticate.
`

const pathUserHelpDesc = `
This endpoint allows you to create, read, update, and delete configuration
for LDAP users that are allowed to authenticate, and associate policies to
them.

Deleting a user will not revoke auth for prior authenticated users in that
user. To do this, do a revoke on "login/<username>" for
the usernames you want revoked.
`
