package transit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/bcrypt"
)

func (b *backend) pathPasswordGenerate() *framework.Path {
	return &framework.Path{
		Pattern: "password/generate",
		Fields: map[string]*framework.FieldSchema{
			"input": {
				Type:        framework.TypeString,
				Description: "Input data",
			},

			"cost": {
				Type:        framework.TypeInt,
				Default:     bcrypt.DefaultCost,
				Description: fmt.Sprintf("Input cost to use, defaults to %d if lower than %d or unspecified", bcrypt.DefaultCost, bcrypt.MinCost),
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathPasswordGenerateWrite,
		},

		HelpSynopsis:    pathPasswordGenerateHelpSyn,
		HelpDescription: pathPasswordGenerateHelpDesc,
	}
}

func (b *backend) pathPasswordGenerateWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rawInput, ok, err := d.GetOkErr("input")
	if err != nil {
		return nil, err
	}
	if !ok {
		return logical.ErrorResponse("input missing"), logical.ErrInvalidRequest
	}

	input := rawInput.(string)
	cost := d.Get("cost").(int)

	retBytes, err := bcrypt.GenerateFromPassword([]byte(input), cost)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("password generation failed: %s", err)), logical.ErrInvalidRequest
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"password": string(retBytes),
		},
	}
	return resp, nil
}

const pathPasswordGenerateHelpSyn = `Generate a bcrypt hash for input data`

const pathPasswordGenerateHelpDesc = `
Generates a bcrypt hash against the given input data at default or specified cost.
`
