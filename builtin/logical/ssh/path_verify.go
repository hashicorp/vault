package ssh

import (
	"context"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathVerify(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "verify",
		Fields: map[string]*framework.FieldSchema{
			"otp": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "[Required] One-Time-Key that needs to be validated",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathVerifyWrite,
		},
		HelpSynopsis:    pathVerifyHelpSyn,
		HelpDescription: pathVerifyHelpDesc,
	}
}

func (b *backend) getOTP(ctx context.Context, s logical.Storage, n string) (*sshOTP, error) {
	entry, err := s.Get(ctx, "otp/"+n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result sshOTP
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathVerifyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	otp := d.Get("otp").(string)

	// If OTP is not a UUID and a string matching VerifyEchoRequest, then the
	// response will be VerifyEchoResponse. This is used by agent to check if
	// connection to Vault server is proper.
	if otp == api.VerifyEchoRequest {
		return &logical.Response{
			Data: map[string]interface{}{
				"message": api.VerifyEchoResponse,
			},
		}, nil
	}

	// Create the salt of OTP because entry would have been create with the
	// salt and not directly of the OTP. Salt will yield the same value which
	// because the seed is the same, the backend salt.
	salt, err := b.Salt(ctx)
	if err != nil {
		return nil, err
	}
	otpSalted := salt.SaltID(otp)

	// Return nil if there is no entry found for the OTP
	otpEntry, err := b.getOTP(ctx, req.Storage, otpSalted)
	if err != nil {
		return nil, err
	}
	if otpEntry == nil {
		return logical.ErrorResponse("OTP not found"), nil
	}

	// Delete the OTP if found. This is what makes the key an OTP.
	err = req.Storage.Delete(ctx, "otp/"+otpSalted)
	if err != nil {
		return nil, err
	}

	// Return username and IP only if there were no problems uptill this point.
	return &logical.Response{
		Data: map[string]interface{}{
			"username":  otpEntry.Username,
			"ip":        otpEntry.IP,
			"role_name": otpEntry.RoleName,
		},
	}, nil
}

const pathVerifyHelpSyn = `
Validate the OTP provided by Vault SSH Agent.
`

const pathVerifyHelpDesc = `
This path will be used by Vault SSH Agent running in the remote hosts. The OTP
provided by the client is sent to Vault for validation by the agent. If Vault
finds an entry for the OTP, it responds with the username and IP it is associated
with. Agent uses this information to authenticate the client. Vault deletes the
OTP after validating it once.
`
