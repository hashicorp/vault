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
			"iam_endpoint": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Endpoint to custom IAM server URL",
			},
			"sts_endpoint": &framework.FieldSchema{
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
	iamendpoint := data.Get("iam_endpoint").(string)
	stsendpoint := data.Get("sts_endpoint").(string)

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
	IAMEndpoint string `json:"iam_endpoint"`
	STSEndpoint string `json:"sts_endpoint"`
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
