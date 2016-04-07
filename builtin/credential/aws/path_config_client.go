package aws

import (
	"fmt"

	"github.com/fatih/structs"
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
				Description: "Region for API calls. Defaults to the value of the AWS_REGION env var. Required.",
			},
		},

		ExistenceCheck: b.pathConfigClientExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigClientCreateUpdate,
			logical.DeleteOperation: b.pathConfigClientDelete,
			logical.ReadOperation:   b.pathConfigClientRead,
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
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

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

func (b *backend) pathConfigClientRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	clientConfig, err := clientConfigEntry(req.Storage)
	if err != nil {
		return nil, err
	}

	if clientConfig == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(clientConfig).Map(),
	}, nil
}

func (b *backend) pathConfigClientDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()

	err := req.Storage.Delete("config/client")
	if err != nil {
		b.configMutex.Unlock()
		return nil, err
	}

	b.configMutex.Unlock()

	_, err = b.clientEC2(req.Storage, true)
	if err != nil {
		return nil, fmt.Errorf("error creating client with updated credentials: %s", err)
	}

	return nil, nil
}

// pathConfigClientCreateUpdate is used to register the 'aws_secret_key' and 'aws_access_key'
// that can be used to interact with AWS EC2 API.
func (b *backend) pathConfigClientCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

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

	changedCreds := false

	accessKeyStr, ok := data.GetOk("access_key")
	if ok {
		if configEntry.AccessKey != accessKeyStr.(string) {
			changedCreds = true
		}
		configEntry.AccessKey = accessKeyStr.(string)
	} else if req.Operation == logical.CreateOperation {
		// Use the default
		configEntry.AccessKey = data.Get("access_key").(string)
	}

	secretKeyStr, ok := data.GetOk("secret_key")
	if ok {
		if configEntry.SecretKey != secretKeyStr.(string) {
			changedCreds = true
		}
		configEntry.SecretKey = secretKeyStr.(string)
	} else if req.Operation == logical.CreateOperation {
		configEntry.SecretKey = data.Get("secret_key").(string)
	}

	entry, err := logical.StorageEntryJSON("config/client", configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	if changedCreds {
		// We have to be careful here to re-lock as we have a deferred unlock
		// queued up and unlocking an unlocked mutex leads to a panic
		b.configMutex.Unlock()
		_, err = b.clientEC2(req.Storage, true)
		b.configMutex.Lock()
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"error creating client with updated credentials: %s", err),
			), nil
		}
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
