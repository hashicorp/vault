package ssh

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	VerifyEchoRequest  = "verify-echo-request"
	VerifyEchoResponse = "verify-echo-response"
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

func (b *backend) pathVerifyWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	otp := d.Get("otp").(string)

	if otp == VerifyEchoRequest {
		return &logical.Response{
			Data: map[string]interface{}{
				"message": VerifyEchoResponse,
			},
		}, nil
	}

	otpSalted := b.salt.SaltID(otp)
	entry, err := req.Storage.Get("otp/" + otpSalted)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var otpEntry sshOTP
	if err := entry.DecodeJSON(&otpEntry); err != nil {
		return nil, nil
	}

	err = req.Storage.Delete("otp/" + otpSalted)
	if err != nil {
		return nil, err
	}

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
that the key is valid only once (hence one-time).
`
