package postgresql

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/pquerna/totp"
)

// Update with TOTP values
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

	// Generate TOTP token
	/*
		//Generate key using totp library
		totpKey, err := totp.GenerateCodeCustom(role.key, time.Now().UTC(), ValdidateOpts{
			Period: role.period,
			Skew: 1,
			Digits: otp.DigitsSix
			Algorithm: otp.AlgorithmSHA1
		});

		if err != nil {
			return nil, err
		}
	*/

	// Return the secret
	b.logger.Trace("totp/pathRoleCreateRead: generating secret")

	/*
			return &logical.Response{
			Data: map[string]interface{}{
				"token":            totpKey,
			},
		}, nil
	*/
	return resp, nil
}

// Update help strings
const pathRoleCreateReadHelpSyn = `
Request database credentials for a certain role.
`

const pathRoleCreateReadHelpDesc = `
This path reads database credentials for a certain role. The
database credentials will be generated on demand and will be automatically
revoked when the lease is up.
`
