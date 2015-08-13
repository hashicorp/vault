package ssh

import (
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
				Description: "One-time-key for SSH session",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathVerifyWrite,
		},
		HelpSynopsis:    pathVerifyHelpSyn,
		HelpDescription: pathVerifyHelpDesc,
	}
}

func (b *backend) getOTP(s logical.Storage, n string) (*sshOTP, error) {
	entry, err := s.Get("otp/" + n)
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

func (b *backend) pathVerifyWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

	otpSalted := b.salt.SaltID(otp)

	// Return nil if there is no entry found for the OTP
	otpEntry, err := b.getOTP(req.Storage, otpSalted)
	if err != nil {
		return nil, err
	}
	if otpEntry == nil {
		return nil, nil
	}

	// Delete the OTP if found. This is what makes the key an OTP.
	err = req.Storage.Delete("otp/" + otpSalted)
	if err != nil {
		return nil, err
	}

	// Return username and IP only if there were no problems uptill this point.
	return &logical.Response{
		Data: map[string]interface{}{
			"username": otpEntry.Username,
			"ip":       otpEntry.IP,
		},
	}, nil
}

const pathVerifyHelpSyn = `
Tells if the key provided by the client is valid or not.
`

const pathVerifyHelpDesc = `
This path will be used by the vault agent running in the
target machine to check if the key provided by the client
to establish the SSH connection is valid or not.

This key will be a one-time-password. The vault server responds
that the key is valid and then deletes it, hence the key is OTP. 
`
