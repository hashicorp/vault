package aws

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigRoot() *framework.Path {
	return &framework.Path{
		Pattern: "config/root",
		Fields: map[string]*framework.FieldSchema{
			"access_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Access key with permission to create new keys.",
			},

			"secret_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Secret key with permission to create new keys.",
			},

			"region": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Region for API calls.",
			},
			"iamendpoint": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Endpoint to custom IAM server URL",
			},
			"stsendpoint": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Endpoint to custom STS server URL",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: pathConfigRootWrite,
		},

		HelpSynopsis:    pathConfigRootHelpSyn,
		HelpDescription: pathConfigRootHelpDesc,
	}
}

func pathConfigRootWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	region := data.Get("region").(string)
	iamendpoint := data.Get("iamendpoint").(string)
	stsendpoint := data.Get("stsendpoint").(string)

	entry, err := logical.StorageEntryJSON("config/root", rootConfig{
		AccessKey:   data.Get("access_key").(string),
		SecretKey:   data.Get("secret_key").(string),
		IAMEndpoint: iamendpoint,
		STSEndpoint: stsendpoint,
		Region:      region,
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
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	IAMEndpoint string `json:"iamendpoint"`
	STSEndpoint string `json:"stsendpoint"`
	Region      string `json:"region"`
}

const pathConfigRootHelpSyn = `
Configure the root credentials that are used to manage IAM.
`

const pathConfigRootHelpDesc = `
Before doing anything, the AWS backend needs credentials that are able
to manage IAM policies, users, access keys, etc. This endpoint is used
to configure those credentials. They don't necessarilly need to be root
keys as long as they have permission to manage IAM.
`
