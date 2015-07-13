package userpass

import (
	"crypto/subtle"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/bcrypt"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username of the user.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLogin(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	password := d.Get("password").(string)

	// Get the user and validate auth
	user, err := b.User(req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return logical.ErrorResponse("unknown username or password"), nil
	}

	// Check for a password match. Check for a hash collision for Vault 0.2+,
	// but handle the older legacy passwords with a constant time comparison.
	passwordBytes := []byte(password)
	if user.PasswordHash != nil {
		if err := bcrypt.CompareHashAndPassword(user.PasswordHash, passwordBytes); err != nil {
			return logical.ErrorResponse("unknown username or password"), nil
		}
	} else {
		if subtle.ConstantTimeCompare([]byte(user.Password), passwordBytes) != 1 {
			return logical.ErrorResponse("unknown username or password"), nil
		}
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: user.Policies,
			Metadata: map[string]string{
				"username": username,
			},
			DisplayName: username,
		},
	}, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the user and validate auth
	user, err := b.User(req.Storage, req.Auth.Metadata["username"])
	if err != nil {
		return nil, err
	}
	if user == nil {
		// User no longer exists, do not renew
		return nil, nil
	}

	return framework.LeaseExtend(1*time.Hour, 0, false)(req, d)
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
