package aws

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigClient(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/client",
		Fields: map[string]*framework.FieldSchema{
			"access_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Access key with permission to query instance metadata.",
			},

			"secret_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Secret key with permission to query instance metadata.",
			},

			"region": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "us-east-1",
				Description: "Region for API calls.",
			},
		},

		ExistenceCheck: b.pathConfigClientExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigClientCreateUpdate,
			logical.UpdateOperation: b.pathConfigClientCreateUpdate,
		},

		HelpSynopsis:    pathConfigClientHelpSyn,
		HelpDescription: pathConfigClientHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigClientExistenceCheck(
	req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := clientConfigEntry(req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// Fetch the client configuration required to access the AWS API.
func clientConfigEntry(s logical.Storage) (*clientConfig, error) {
	entry, err := s.Get("config/client")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result clientConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// pathConfigClientCreateUpdate is used to register the 'aws_secret_key' and 'aws_access_key'
// that can be used to interact with AWS EC2 API.
func (b *backend) pathConfigClientCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	configEntry, err := clientConfigEntry(req.Storage)
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		configEntry = &clientConfig{}
	}

	regionStr, ok := data.GetOk("region")
	if ok {
		configEntry.Region = regionStr.(string)
	} else if req.Operation == logical.CreateOperation {
		configEntry.Region = data.Get("region").(string)
	}

	// Either a valid region needs to be provided or it should be left empty
	// so a default value could take over.
	if configEntry.Region == "" {
		return nil, fmt.Errorf("invalid region")

	}

	accessKeyStr, ok := data.GetOk("access_key")
	if ok {
		configEntry.AccessKey = accessKeyStr.(string)
	} else if req.Operation == logical.CreateOperation {
		if configEntry.AccessKey = data.Get("access_key").(string); configEntry.AccessKey == "" {
			return nil, fmt.Errorf("missing access_key")
		}
	}

	secretKeyStr, ok := data.GetOk("secret_key")
	if ok {
		configEntry.SecretKey = secretKeyStr.(string)
	} else if req.Operation == logical.CreateOperation {
		if configEntry.SecretKey = data.Get("secret_key").(string); configEntry.SecretKey == "" {
			return nil, fmt.Errorf("missing secret_key")
		}
	}

	entry, err := logical.StorageEntryJSON("config/client", configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// Struct to hold 'aws_access_key' and 'aws_secret_key' that are required to
// interact with the AWS EC2 API.
type clientConfig struct {
	AccessKey string `json:"access_key" structs:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" structs:"secret_key" mapstructure:"secret_key"`
	Region    string `json:"region" structs:"region" mapstructure:"region"`
}

const pathConfigClientHelpSyn = `
Configure the client credentials that are used to query instance details from AWS EC2 API.
`

const pathConfigClientHelpDesc = `
AWS auth backend makes API calls to retrieve EC2 instance metadata.
The aws_secret_key and aws_access_key registered with Vault should have the
permissions to make these API calls.
`
