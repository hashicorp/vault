package vaultTotp

import (
	strings "strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	totp "github.com/ptollot/totp"
)

func pathTotpConfig() *framework.Path {
	return &framework.Path{

		Pattern: `totp/config`,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "User identifier",
			},
			"seed": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "string associated with the user for the creation of totp codes",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathTotpConfigWrite,
			logical.ReadOperation: pathTotpConfigRead,
			logical.DeleteOperation: pathTotpConfigDelete,
		},

		HelpSynopsis:    pathTotpConfigHelpSyn,
		HelpDescription: pathTotpConfigHelpDesc,
	}
}

func pathTotpConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	seed := d.Get("seed").(string)
	totp.TotpReference.Register(username, seed)
	return &logical.Response{
		Data: map[string]interface{}{
			"username": username,
		},
	}, nil
}

func pathTotpConfigDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	totp.TotpReference.Delete(username)
	return &logical.Response{
		Data: map[string]interface{}{
			"username": username,
		},
	}, nil
}

func pathTotpConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	userNames := totp.TotpReference.GetUsernames()

	return &logical.Response{
		Data: map[string]interface{}{
			"username": strings.Join(userNames, ", "),
		},
	}, nil
}

const pathTotpConfigHelpSyn = `
Configure totp second factor behavior.
`

const pathTotpConfigHelpDesc = `
This endpoint allows you to configure how the original auth backend username maps to
the totp username by providing a template format string.
`
