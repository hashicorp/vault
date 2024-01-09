// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package totp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

func pathCode(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "code/" + framework.GenericNameWithAtRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTOTP,
			OperationSuffix: "code",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key.",
			},
			"code": {
				Type:        framework.TypeString,
				Description: "TOTP code to be validated.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathReadCode,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "generate",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathValidateCode,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "validate",
				},
			},
		},

		HelpSynopsis:    pathCodeHelpSyn,
		HelpDescription: pathCodeHelpDesc,
	}
}

func (b *backend) pathReadCode(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the key
	key, err := b.Key(ctx, req.Storage, name)
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

func (b *backend) pathValidateCode(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	code := data.Get("code").(string)

	// Enforce input value requirements
	if code == "" {
		return logical.ErrorResponse("the code value is required"), nil
	}

	// Get the key's stored values
	key, err := b.Key(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown key: %s", name)), nil
	}

	usedName := fmt.Sprintf("%s_%s", name, code)

	_, ok := b.usedCodes.Get(usedName)
	if ok {
		return logical.ErrorResponse("code already used; wait until the next time period"), nil
	}

	valid, err := totplib.ValidateCustom(code, key.Key, time.Now(), totplib.ValidateOpts{
		Period:    key.Period,
		Skew:      key.Skew,
		Digits:    key.Digits,
		Algorithm: key.Algorithm,
	})
	if err != nil && err != otplib.ErrValidateInputInvalidLength {
		return logical.ErrorResponse("an error occurred while validating the code"), err
	}

	// Take the key skew, add two for behind and in front, and multiple that by
	// the period to cover the full possibility of the validity of the key
	err = b.usedCodes.Add(usedName, nil, time.Duration(
		int64(time.Second)*
			int64(key.Period)*
			int64((2+key.Skew))))
	if err != nil {
		return nil, fmt.Errorf("error adding code to used cache: %w", err)
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
