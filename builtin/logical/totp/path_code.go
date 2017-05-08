package totp

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

func pathCode(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "code/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key.",
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

		HelpSynopsis:    pathCodeHelpSyn,
		HelpDescription: pathCodeHelpDesc,
	}
}

func (b *backend) pathReadCode(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the key
	key, err := b.Key(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown key: %s", name)), nil
	}

	// Generate password using totp library
	totpToken, err := totplib.GenerateCodeCustom(key.Key, time.Now(), totplib.ValidateOpts{
		Period:    key.Period,
		Digits:    key.Digits,
		Algorithm: key.Algorithm,
	})
	if err != nil {
		return nil, err
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"code": totpToken,
		},
	}, nil
}

func (b *backend) pathValidateCode(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	code := data.Get("code").(string)

	// Enforce input value requirements
	if code == "" {
		return logical.ErrorResponse("the code value is required"), nil
	}

	// Get the key's stored values
	key, err := b.Key(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown key: %s", name)), nil
	}

	valid, err := totplib.ValidateCustom(code, key.Key, time.Now(), totplib.ValidateOpts{
		Period:    key.Period,
		Skew:      key.Skew,
		Digits:    key.Digits,
		Algorithm: key.Algorithm,
	})
	if err != nil && err != otplib.ErrValidateInputInvalidLength {
		return logical.ErrorResponse("an error occured while validating the code"), err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"valid": valid,
		},
	}, nil
}

const pathCodeHelpSyn = `
Request time-based one-time use password or validate a password for a certain key .
`
const pathCodeHelpDesc = `
This path generates and validates time-based one-time use passwords for a certain key. 

`
