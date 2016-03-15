package userpass

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUserPassword(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/password/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username for this user.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathUserPasswordUpdate,
		},

		HelpSynopsis:    pathUserPasswordHelpSyn,
		HelpDescription: pathUserPasswordHelpDesc,
	}
}

func (b *backend) pathUserPasswordUpdate(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	password := d.Get("password").(string)
	if password == "" {
		return nil, fmt.Errorf("missing password")
	}

	userEntry, err := b.User(req.Storage, strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	if userEntry == nil {
		return nil, nil
	}

	// Generate a hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set the new password hash
	userEntry.PasswordHash = hash

	// Store the UserEntry
	err = b.SetUser(req.Storage, username, userEntry)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

const pathUserPasswordHelpSyn = `
Reset user's password.
`

const pathUserPasswordHelpDesc = `
This endpoint allows resetting the user's password.
`
