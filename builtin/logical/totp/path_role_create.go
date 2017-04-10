package totp

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	totplib "github.com/pquerna/otp/totp"
)

func pathRoleCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "code/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"code": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "TOTP code to be validated.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathReadCode,
			logical.UpdateOperation: b.pathValidateCode,
		},

		HelpSynopsis:    pathRoleCreateReadHelpSyn,
		HelpDescription: pathRoleCreateReadHelpDesc,
	}
}

func (b *backend) pathReadCode(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.logger.Trace("totp/pathReadCode: enter")
	defer b.logger.Trace("totp/pathReadCode: exit")

	name := data.Get("name").(string)

	// Get the key
	b.logger.Trace("totp/pathReadCode: getting key")
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown key: %s", name)), nil
	}

	// Generate password using totp library
	totpToken, err := totplib.GenerateCodeCustom(role.Key, time.Now(), totplib.ValidateOpts{
		Period:    role.Period,
		Digits:    role.Digits,
		Algorithm: role.Algorithm,
	})

	if err != nil {
		return nil, err
	}

	// Return the secret
	b.logger.Trace("totp/pathReadCode: generating secret")

	resp, err := &logical.Response{
		Data: map[string]interface{}{
			"code": totpToken,
		},
	}, nil

	return resp, nil
}

func (b *backend) pathValidateCode(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	code := data.Get("code").(string)

	// Enforce input value requirements
	if code == "" {
		return logical.ErrorResponse("The code value is required."), nil
	}

	// Get the key's stored values
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown key: %s", name)), nil
	}

	valid, err := totplib.ValidateCustom(code, role.Key, time.Now(), totplib.ValidateOpts{
		Period:    role.Period,
		Skew:      role.Skew,
		Digits:    role.Digits,
		Algorithm: role.Algorithm,
	})

	resp, err := &logical.Response{
		Data: map[string]interface{}{
			"valid": valid,
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
