package totp

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},

			"key": {
				Type:        framework.TypeString,
				Description: "The shared master key used to generate a TOTP token.",
			},

			"issuer": {
				Type:        framework.TypeString,
				Description: `The name of the key's issuing organization.`,
			},

			"account_name": {
				Type:        framework.TypeString,
				Description: `The name of the account associated with the key.`,
			},

			"period": {
				Type:        framework.TypeInt,
				Description: `The length of time used to generate a counter for the TOTP token calculation.`,
			},

			"algorithm": {
				Type:        framework.TypeString,
				Description: `The hashing algorithm used to generate the TOTP token.`,
			},

			"digits": {
				Type:        framework.TypeInt,
				Description: `The number of digits in the generated TOTP token.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleCreate,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *backend) Role(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("role/" + data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := b.Role(req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Return values of role
	return &logical.Response{
		Data: map[string]interface{}{
			"issuer":       role.Issuer,
			"account_name": role.Account_Name,
			"period":       role.Period,
			"algorithm":    role.Algorithm,
			"digits":       role.Digits,
		},
	}, nil
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	key := data.Get("key").(string)
	issuer := data.Get("issuer").(string)
	account_name := data.Get("account_name").(string)
	period := data.Get("period").(int)
	algorithm := data.Get("algorithm").(string)
	digits := data.Get("digits").(int)

	// Set optional parameters if neccessary
	if period == 0 {
		period = 30
	}

	switch algorithm {
	case "SHA1", "SHA256", "SHA512", "MD5":
	default:
		algorithm = "SHA1"
	}

	switch digits {
	case 6, 8:
	default:
		digits = 6
	}

	// Store it
	entry, err := logical.StorageEntryJSON("role/"+name, &roleEntry{
		Key:          key,
		Issuer:       issuer,
		Account_Name: account_name,
		Period:       period,
		Algorithm:    algorithm,
		Digits:       digits,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleEntry struct {
	Key          string `json:"key" mapstructure:"key" structs:"key"`
	Issuer       string `json:"issuer" mapstructure:"issuer" structs:"issuer"`
	Account_Name string `json:"account_name" mapstructure:"account_name" structs:"account_name"`
	Period       int    `json:"period" mapstructure:"period" structs:"period"`
	Algorithm    string `json:"algorithm" mapstructure:"algorithm" structs:"algorithm"`
	Digits       int    `json:"digits" mapstructure:"digits" structs:"digits"`
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.

Role Parameters:

  * "key" - required - The shared master key used to generate a TOTP token.

  * "issuer" - required - The name of the key's issuing organization.

  * "account_name" - required - The name of the account associated with the key.

  * "period" - optional - The length of time used to generate a counter for the TOTP token calculation. Default value is 30 seconds.

  * "algorithm" - optional - The hashing algorithm used to generate the TOTP token. Default value is "SHA1". Other options include "SHA256", "SHA512" and "MD5". 

  * "digits" - optional - The number of digits in the generated TOTP token. Default value is 6. Options include 6 or 8.
`
