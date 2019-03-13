package ldap

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
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
				Description: "Name of the LDAP user.",
			},

			"groups": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-separated list of additional groups associated with the user.",
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma-separated list of policies associated with the user.",
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
	err := req.Storage.Delete(ctx, "user/"+d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathUserRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("name").(string)

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("ldap backend not configured"), nil
	}
	if !*cfg.CaseSensitiveNames {
		username = strings.ToLower(username)
	}

	user, err := b.User(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"groups":   strings.Join(user.Groups, ","),
			"policies": user.Policies,
		},
	}, nil
}

func (b *backend) pathUserWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lowercaseGroups := false
	username := d.Get("name").(string)

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("ldap backend not configured"), nil
	}
	if !*cfg.CaseSensitiveNames {
		username = strings.ToLower(username)
		lowercaseGroups = true
	}

	groups := strutil.RemoveDuplicates(strutil.ParseStringSlice(d.Get("groups").(string), ","), lowercaseGroups)
	policies := policyutil.ParsePolicies(d.Get("policies"))
	for i, g := range groups {
		groups[i] = strings.TrimSpace(g)
	}

	// Store it
	entry, err := logical.StorageEntryJSON("user/"+username, &UserEntry{
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
	keys, err := logical.CollectKeys(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	retKeys := make([]string, 0)
	for _, key := range keys {
		if strings.HasPrefix(key, "user/") && !strings.HasPrefix(key, "/") {
			retKeys = append(retKeys, strings.TrimPrefix(key, "user/"))
		}
	}
	return logical.ListResponse(retKeys), nil

}

type UserEntry struct {
	Groups   []string
	Policies []string
}

const pathUserHelpSyn = `
Manage users allowed to authenticate.
`

const pathUserHelpDesc = `
This endpoint allows you to create, read, update, and delete configuration
for LDAP users that are allowed to authenticate, in particular associating
additional groups to them.

Deleting a user will not revoke their auth. To do this, do a revoke on "login/<username>" for
the usernames you want revoked.
`
