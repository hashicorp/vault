package okta

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUsersList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathUserList,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func pathUsers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `users/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the user.",
			},

			"groups": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "List of groups associated with the user.",
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "List of policies associated with the user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathUserDelete,
			logical.ReadOperation:   b.pathUserRead,
			logical.UpdateOperation: b.pathUserWrite,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func (b *backend) User(ctx context.Context, s logical.Storage, n string) (*UserEntry, error) {
	entry, err := s.Get(ctx, "user/"+n)
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

func (b *backend) pathUserDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if len(name) == 0 {
		return logical.ErrorResponse("Error empty name"), nil
	}

	err := req.Storage.Delete(ctx, "user/"+name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathUserRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if len(name) == 0 {
		return logical.ErrorResponse("Error empty name"), nil
	}

	user, err := b.User(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"groups":   user.Groups,
			"policies": user.Policies,
		},
	}, nil
}

func (b *backend) pathUserWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if len(name) == 0 {
		return logical.ErrorResponse("Error empty name"), nil
	}

	groups := d.Get("groups").([]string)
	policies := d.Get("policies").([]string)

	// Store it
	entry, err := logical.StorageEntryJSON("user/"+name, &UserEntry{
		Groups:   groups,
		Policies: policies,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathUserList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	users, err := req.Storage.List(ctx, "user/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(users), nil
}

type UserEntry struct {
	Groups   []string
	Policies []string
}

const pathUserHelpSyn = `
Manage additional groups for users allowed to authenticate.
`

const pathUserHelpDesc = `
This endpoint allows you to create, read, update, and delete configuration
for Okta users that are allowed to authenticate, in particular associating
additional groups to them.

Deleting a user will not revoke their auth. To do this, do a revoke on "login/<username>" for
the usernames you want revoked.
`
