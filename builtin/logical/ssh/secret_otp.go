// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const SecretOTPType = "secret_otp_type"

func secretOTP(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretOTPType,
		Fields: map[string]*framework.FieldSchema{
			"otp": {
				Type:        framework.TypeString,
				Description: "One time password",
			},
		},

		Revoke: b.secretOTPRevoke,
	}
}

func (b *backend) secretOTPRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	otpRaw, ok := req.Secret.InternalData["otp"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	otp, ok := otpRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}

	salt, err := b.Salt(ctx)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Delete(ctx, "otp/"+salt.SaltID(otp))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
