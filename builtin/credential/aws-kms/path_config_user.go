package awsKms

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigUser() *framework.Path {
	return &framework.Path{
		Pattern: "config/user",
		Fields: map[string]*framework.FieldSchema{
			"access_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Access key with kms:Decrypt permission.",
			},

			"secret_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Secret key with kms:Decrypt permission.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathConfigRootWrite,
		},

		HelpSynopsis:    pathConfigRootHelpSyn,
		HelpDescription: pathConfigRootHelpDesc,
	}
}

func pathConfigRootWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	entry, err := logical.StorageEntryJSON("config/user", rootConfig{
		AccessKey: data.Get("access_key").(string),
		SecretKey: data.Get("secret_key").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type rootConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

const pathConfigRootHelpSyn = `
Configure the credentials that are used to decrypt tokens.
`

const pathConfigRootHelpDesc = `
Before doing anything, the AWS KMS backend needs credentials that can be
used to decrypt the tokens. This endpoint is used to configure those
credentials. They don't necessarilly need to be root keys as long as they
have kms:Decrypt permission on keys that are used with Vault.
`
