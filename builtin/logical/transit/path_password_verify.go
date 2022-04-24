package transit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/bcrypt"
)

func (b *backend) pathPasswordVerify() *framework.Path {
	return &framework.Path{
		Pattern: "password/verify",
		Fields: map[string]*framework.FieldSchema{
			"input": {
				Type:        framework.TypeString,
				Description: "Input data",
			},

			"hash": {
				Type:        framework.TypeString,
				Description: "Bcrypt hash to verify against input",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathPasswordVerifyWrite,
		},

		HelpSynopsis:    pathPasswordVerifyHelpSyn,
		HelpDescription: pathPasswordVerifyHelpDesc,
	}
}

func (b *backend) pathPasswordVerifyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rawInput, ok, err := d.GetOkErr("input")
	if err != nil {
		return nil, err
	}
	if !ok {
		return logical.ErrorResponse("input missing"), logical.ErrInvalidRequest
	}

	rawHash, ok, err := d.GetOkErr("hash")
	if err != nil {
		return nil, err
	}
	if !ok {
		return logical.ErrorResponse("hash missing"), logical.ErrInvalidRequest
	}

	input := rawInput.(string)
	hash := rawHash.(string)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("password verification failed: %s", err)), err
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"password": hash,
		},
	}
	return resp, nil
}

const pathPasswordVerifyHelpSyn = `Verify bcrypt hash against input data`

const pathPasswordVerifyHelpDesc = `
Verifies bcrypt hash against the given input data at specified cost.
`
