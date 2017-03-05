package totp

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/pquerna/otp/totp"
)

func pathRoleCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathRoleCreateRead,
		},

		HelpSynopsis:    pathRoleCreateReadHelpSyn,
		HelpDescription: pathRoleCreateReadHelpDesc,
	}
}

func (b *backend) pathRoleCreateRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.logger.Trace("totp/pathRoleCreateRead: enter")
	defer b.logger.Trace("totp/pathRoleCreateRead: exit")

	name := data.Get("name").(string)

	// Get the role
	b.logger.Trace("totp/pathRoleCreateRead: getting role")
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Generate password using totp library
	totpToken, err := totp.GenerateCodeCustom(role.Key, time.Now().UTC(), ValdidateOpts{
		Period:    role.Period,
		Digits:    role.Digits,
		Algorithm: role.Algorithm,
	})

	if err != nil {
		return nil, err
	}

	// Return the secret
	b.logger.Trace("totp/pathRoleCreateRead: generating secret")

	resp := &logical.Response{
		Data: map[string]interface{}{
			"token": totpToken,
		},
	}, nil

	return resp, nil
}

const pathRoleCreateReadHelpSyn = `
Request time-based one-time use password for a certain role.
`
const pathRoleCreateReadHelpDesc = `
This path generates a time-based one-time use password for a certain role. 
`
