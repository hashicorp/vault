// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathUsersList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixOkta,
			OperationSuffix: "users",
			Navigation:      true,
			ItemType:        "User",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathUserList,
				Summary:  "List users configured in the Okta auth method.",
			},
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func pathUsers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `users/(?P<name>.+)`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixOkta,
			OperationSuffix: "user",
			Action:          "Create",
			ItemType:        "User",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the user.",
			},

			"groups": {
				Type:        framework.TypeCommaStringSlice,
				Description: "List of groups associated with the user.",
			},

			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: "List of policies associated with the user.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathUserDelete,
				Summary:  "Delete an Okta auth method user.",
				Responses: map[int][]framework.Response{
					204: {{Description: "No Content"}},
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathUserRead,
				Summary:  "Return the properties of an Okta auth method user.",
				Responses: map[int][]framework.Response{
					200: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"groups": {
								Type:        framework.TypeCommaStringSlice,
								Description: "List of groups associated with the user.",
							},
							"policies": {
								Type:        framework.TypeCommaStringSlice,
								Description: "List of policies associated with the user.",
							},
						},
					}},
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathUserWrite,
				Summary:  "Create or update an Okta auth method user.",
				Responses: map[int][]framework.Response{
					204: {{Description: "No Content"}},
				},
			},
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
