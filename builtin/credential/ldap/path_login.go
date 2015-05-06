package ldap

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login`,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "DN (distinguished name) to be used for login.",
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
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("ldap backend not configured"), nil
	}

	c, err := cfg.DialLDAP()
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Try to authenticate to the server using the provided credentials
	binddn := fmt.Sprintf("%s=%s,%s", cfg.UserAttr, username, cfg.Domain)
	cerr := c.Bind(binddn, password)
	if cerr != nil {
		return logical.ErrorResponse(fmt.Sprintf("LDAP bind failed: %s", cerr.Error())), nil
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: []string{"root"},
			Metadata: map[string]string{
				"username": username,
			},
			DisplayName: username,
		},
	}, nil
}

// func (b *backend) pathLoginRenew(
// 	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
// 	// Get the user and validate auth
// 	user, err := b.User(req.Storage, req.Auth.Metadata["username"])
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil {
// 		// User no longer exists, do not renew
// 		return nil, nil
// 	}

// 	return framework.LeaseExtend(1*time.Hour, 0)(req, d)
// }

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
