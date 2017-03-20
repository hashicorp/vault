package totp

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
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

	// Translate digits and algorithm to a format the totp library understands
	var digits otplib.Digits
	switch role.Digits {
	case 6:
		digits = otplib.DigitsSix
	case 8:
		digits = otplib.DigitsEight
	default:
		return logical.ErrorResponse("The digit value can only be 6 or 8."), nil
	}

	var algorithm otplib.Algorithm
	switch role.Algorithm {
	case "SHA1":
		algorithm = otplib.AlgorithmSHA1
	case "SHA256":
		algorithm = otplib.AlgorithmSHA256
	case "SHA512":
		algorithm = otplib.AlgorithmSHA512
	default:
		return logical.ErrorResponse("The algorithm value is not valid."), nil
	}

	period := uint(role.Period)

	// Generate password using totp library
	totpToken, err := totplib.GenerateCodeCustom(role.Key, time.Now().UTC(), totplib.ValidateOpts{
		Period:    period,
		Digits:    digits,
		Algorithm: algorithm,
	})

	if err != nil {
		return nil, err
	}

	// Return the secret
	b.logger.Trace("totp/pathRoleCreateRead: generating secret")

	resp, err := &logical.Response{
		Data: map[string]interface{}{
			"code": totpToken,
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
