package userpass

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"strings"
)

func pathUserPassword(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("username") + "/password$",
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username for this user.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},

			"password_hash": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Pre-hashed password in bcrypt format for this user.",
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

	username := d.Get("username").(string)

	userEntry, err := b.user(req.Storage, username)
	if err != nil {
		return nil, err
	}
	if userEntry == nil {
		return nil, fmt.Errorf("username does not exist")
	}

	userErr, intErr := b.updateUserPassword(req, d, userEntry)
	if intErr != nil {
		return nil, err
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), logical.ErrInvalidRequest
	}

	return nil, b.setUser(req.Storage, username, userEntry)
}

func (b *backend) updateUserPassword(req *logical.Request, d *framework.FieldData, userEntry *UserEntry) (error, error) {
	password := d.Get("password").(string)
	prehashedPassword := d.Get("password_hash").(string)

	if password != "" && prehashedPassword != "" {
		return fmt.Errorf("can't provide both password and password_hash"), nil
	}

	if password == "" && prehashedPassword == "" {
		return fmt.Errorf("missing password"), nil
	}

	var hash []byte

	// If a hash was provided, use it.
	if prehashedPassword != "" {
		if strings.HasPrefix(prehashedPassword, "$2a$") {
			hash = []byte(prehashedPassword)
		} else {
			return nil, fmt.Errorf("password_hash doesn't appear to be a valid bcrypt hash")
		}
	} else {
		// Otherwise, generate a hash of the password
		var err error
		hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
	}

	userEntry.PasswordHash = hash
	return nil, nil
}

const pathUserPasswordHelpSyn = `
Reset user's password.
`

const pathUserPasswordHelpDesc = `
This endpoint allows resetting the user's password.
`
