package userpass

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

	// Get the user and validate auth
	user, err := b.User(req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Password != d.Get("password").(string) {
		return logical.ErrorResponse("unknown username or password"), nil
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

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
