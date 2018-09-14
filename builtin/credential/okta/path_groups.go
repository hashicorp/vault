package okta

import (
	"context"
	"strings"

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
				Description: "Name of the Okta group.",
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

// We look up groups in a case-insensitive manner since Okta is case-preserving
// but case-insensitive for comparisons
func (b *backend) Group(ctx context.Context, s logical.Storage, n string) (*GroupEntry, string, error) {
	canonicalName := n
	entry, err := s.Get(ctx, "group/"+n)
	if err != nil {
		return nil, "", err
	}
	if entry == nil {
		entries, err := s.List(ctx, "group/")
		if err != nil {
			return nil, "", err
		}
		for _, groupName := range entries {
			if strings.EqualFold(groupName, n) {
				entry, err = s.Get(ctx, "group/"+groupName)
				if err != nil {
					return nil, "", err
				}
				canonicalName = groupName
				break
			}
		}
	}
	if entry == nil {
		return nil, "", nil
	}

	var result GroupEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, "", err
	}

	return &result, canonicalName, nil
}

func (b *backend) pathGroupDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if len(name) == 0 {
		return logical.ErrorResponse("'name' must be supplied"), nil
	}

	entry, canonicalName, err := b.Group(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		err := req.Storage.Delete(ctx, "group/"+canonicalName)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (b *backend) pathGroupRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if len(name) == 0 {
		return logical.ErrorResponse("'name' must be supplied"), nil
	}

	group, _, err := b.Group(ctx, req.Storage, name)
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

func (b *backend) pathGroupWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if len(name) == 0 {
		return logical.ErrorResponse("'name' must be supplied"), nil
	}

	// Check for an existing group, possibly lowercased so that we keep using
	// existing user set values
	_, canonicalName, err := b.Group(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if canonicalName != "" {
		name = canonicalName
	} else {
		name = strings.ToLower(name)
	}

	entry, err := logical.StorageEntryJSON("group/"+name, &GroupEntry{
		Policies: policyutil.ParsePolicies(d.Get("policies")),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathGroupList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	groups, err := req.Storage.List(ctx, "group/")
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
for Okta groups that are allowed to authenticate, and associate policies to
them.

Deleting a group will not revoke auth for prior authenticated users in that
group. To do this, do a revoke on "login/<username>" for
the usernames you want revoked.
`
