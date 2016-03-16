package userpass

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUserPolicies(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("username") + "/policies$",
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username for this user.",
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-separated list of policies",
			},
		},

		ExistenceCheck: b.userPoliciesExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathUserPoliciesUpdate,
		},

		HelpSynopsis:    pathUserPoliciesHelpSyn,
		HelpDescription: pathUserPoliciesHelpDesc,
	}
}

func (b *backend) userPoliciesExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	return true, nil
}

func (b *backend) pathUserPoliciesUpdate(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	username := d.Get("username").(string)

	userEntry, err := b.user(req.Storage, username)
	if err != nil {
		return nil, err
	}
	if userEntry == nil {
		return nil, fmt.Errorf("username does not exist")
	}

	err = b.updateUserPolicies(req, d, userEntry)
	if err != nil {
		return nil, err
	}

	return nil, b.setUser(req.Storage, username, userEntry)
}

func (b *backend) updateUserPolicies(req *logical.Request, d *framework.FieldData, userEntry *UserEntry) error {
	policies := strings.Split(d.Get("policies").(string), ",")
	for i, p := range policies {
		policies[i] = strings.TrimSpace(p)
	}
	userEntry.Policies = policies
	return nil
}

const pathUserPoliciesHelpSyn = `
Update the policies associated with the username.
`

const pathUserPoliciesHelpDesc = `
This endpoint allows updating the policies associated with the username.
`
