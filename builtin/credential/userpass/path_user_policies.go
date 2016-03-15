package userpass

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUserPolicies(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/policies/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username for this user.",
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-separated list of policies",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathUserPoliciesUpdate,
		},

		HelpSynopsis:    pathUserPoliciesHelpSyn,
		HelpDescription: pathUserPoliciesHelpDesc,
	}
}
func (b *backend) pathUserPoliciesUpdate(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	policies := strings.Split(d.Get("policies").(string), ",")
	for i, p := range policies {
		policies[i] = strings.TrimSpace(p)
	}

	userEntry, err := b.User(req.Storage, strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	if userEntry == nil {
		return nil, nil
	}

	userEntry.Policies = policies

	// Store the UserEntry
	err = b.SetUser(req.Storage, username, userEntry)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

const pathUserPoliciesHelpSyn = `
Update the policies associated with the username.
`

const pathUserPoliciesHelpDesc = `
This endpoint allows updating the policies associated with the username.
`
