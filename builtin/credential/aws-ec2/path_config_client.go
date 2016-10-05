package awsec2

import (
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigClient(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/client$",
		Fields: map[string]*framework.FieldSchema{
			"access_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "AWS Access Key ID for the account used to make AWS API requests.",
			},

			"secret_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "AWS Secret Access Key for the account used to make AWS API requests.",
			},

			"endpoint": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "URL to override the default generated endpoint for making AWS EC2 API calls.",
			},
		},

		ExistenceCheck: b.pathConfigClientExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigClientCreateUpdate,
			logical.UpdateOperation: b.pathConfigClientCreateUpdate,
			logical.DeleteOperation: b.pathConfigClientDelete,
			logical.ReadOperation:   b.pathConfigClientRead,
		},

		HelpSynopsis:    pathConfigClientHelpSyn,
		HelpDescription: pathConfigClientHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigClientExistenceCheck(
	req *logical.Request, data *framework.FieldData) (bool, error) {

	entry, err := b.lockedClientConfigEntry(req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// Fetch the client configuration required to access the AWS API, after acquiring an exclusive lock.
func (b *backend) lockedClientConfigEntry(s logical.Storage) (*clientConfig, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedClientConfigEntry(s)
}

// Fetch the client configuration required to access the AWS API.
func (b *backend) nonLockedClientConfigEntry(s logical.Storage) (*clientConfig, error) {
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
	clientConfig, err := b.lockedClientConfigEntry(req.Storage)
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
	defer b.configMutex.Unlock()

	if err := req.Storage.Delete("config/client"); err != nil {
		return nil, err
	}

	// Remove all the cached EC2 client objects in the backend.
	b.flushCachedEC2Clients()

	// Remove all the cached EC2 client objects in the backend.
	b.flushCachedIAMClients()

	return nil, nil
}

// pathConfigClientCreateUpdate is used to register the 'aws_secret_key' and 'aws_access_key'
// that can be used to interact with AWS EC2 API.
func (b *backend) pathConfigClientCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedClientConfigEntry(req.Storage)
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		configEntry = &clientConfig{}
	}

	changedCreds := false

	accessKeyStr, ok := data.GetOk("access_key")
	if ok {
		if configEntry.AccessKey != accessKeyStr.(string) {
			changedCreds = true
			configEntry.AccessKey = accessKeyStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		// Use the default
		configEntry.AccessKey = data.Get("access_key").(string)
	}

	secretKeyStr, ok := data.GetOk("secret_key")
	if ok {
		if configEntry.SecretKey != secretKeyStr.(string) {
			changedCreds = true
			configEntry.SecretKey = secretKeyStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.SecretKey = data.Get("secret_key").(string)
	}

	endpointStr, ok := data.GetOk("endpoint")
	if ok {
		if configEntry.Endpoint != endpointStr.(string) {
			changedCreds = true
			configEntry.Endpoint = endpointStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.Endpoint = data.Get("endpoint").(string)
	}

	// Since this endpoint supports both create operation and update operation,
	// the error checks for access_key and secret_key not being set are not present.
	// This allows calling this endpoint multiple times to provide the values.
	// Hence, the readers of this endpoint should do the validation on
	// the validation of keys before using them.
	entry, err := logical.StorageEntryJSON("config/client", configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	if changedCreds {
		b.flushCachedEC2Clients()
		b.flushCachedIAMClients()
	}

	return nil, nil
}

// Struct to hold 'aws_access_key' and 'aws_secret_key' that are required to
// interact with the AWS EC2 API.
type clientConfig struct {
	AccessKey string `json:"access_key" structs:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" structs:"secret_key" mapstructure:"secret_key"`
	Endpoint  string `json:"endpoint" structs:"endpoint" mapstructure:"endpoint"`
}

const pathConfigClientHelpSyn = `
Configure AWS IAM credentials that are used to query instance and role details from the AWS API.
`

const pathConfigClientHelpDesc = `
The aws-ec2 auth backend makes AWS API queries to retrieve information
regarding EC2 instances that perform login operations. The 'aws_secret_key' and
'aws_access_key' parameters configured here should map to an AWS IAM user that
has permission to make the following API queries:

* ec2:DescribeInstances
* iam:GetInstanceProfile (if IAM Role binding is used)
`
