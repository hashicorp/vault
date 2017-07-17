package awsauth

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

			"iam_endpoint": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "URL to override the default generated endpoint for making AWS IAM API calls.",
			},

			"sts_endpoint": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "URL to override the default generated endpoint for making AWS STS API calls.",
			},

			"iam_server_id_header_value": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "Value to require in the X-Vault-AWS-IAM-Server-ID request header",
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

	// unset the cached default AWS account ID
	b.defaultAWSAccountID = ""

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

	// changedCreds is whether we need to flush the cached AWS clients and store in the backend
	changedCreds := false
	// changedOtherConfig is whether other config has changed that requires storing in the backend
	// but does not require flushing the cached clients
	changedOtherConfig := false

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

	iamEndpointStr, ok := data.GetOk("iam_endpoint")
	if ok {
		if configEntry.IAMEndpoint != iamEndpointStr.(string) {
			changedCreds = true
			configEntry.IAMEndpoint = iamEndpointStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.IAMEndpoint = data.Get("iam_endpoint").(string)
	}

	stsEndpointStr, ok := data.GetOk("sts_endpoint")
	if ok {
		if configEntry.STSEndpoint != stsEndpointStr.(string) {
			// We don't directly cache STS clients as they are ever directly used.
			// However, they are potentially indirectly used as credential providers
			// for the EC2 and IAM clients, and thus we would be indirectly caching
			// them there. So, if we change the STS endpoint, we should flush those
			// cached clients.
			changedCreds = true
			configEntry.STSEndpoint = stsEndpointStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.STSEndpoint = data.Get("sts_endpoint").(string)
	}

	headerValStr, ok := data.GetOk("iam_server_id_header_value")
	if ok {
		if configEntry.IAMServerIdHeaderValue != headerValStr.(string) {
			// NOT setting changedCreds here, since this isn't really cached
			configEntry.IAMServerIdHeaderValue = headerValStr.(string)
			changedOtherConfig = true
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.IAMServerIdHeaderValue = data.Get("iam_server_id_header_value").(string)
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

	if changedCreds || changedOtherConfig || req.Operation == logical.CreateOperation {
		if err := req.Storage.Put(entry); err != nil {
			return nil, err
		}
	}

	if changedCreds {
		b.flushCachedEC2Clients()
		b.flushCachedIAMClients()
		b.defaultAWSAccountID = ""
	}

	return nil, nil
}

// Struct to hold 'aws_access_key' and 'aws_secret_key' that are required to
// interact with the AWS EC2 API.
type clientConfig struct {
	AccessKey              string `json:"access_key" structs:"access_key" mapstructure:"access_key"`
	SecretKey              string `json:"secret_key" structs:"secret_key" mapstructure:"secret_key"`
	Endpoint               string `json:"endpoint" structs:"endpoint" mapstructure:"endpoint"`
	IAMEndpoint            string `json:"iam_endpoint" structs:"iam_endpoint" mapstructure:"iam_endpoint"`
	STSEndpoint            string `json:"sts_endpoint" structs:"sts_endpoint" mapstructure:"sts_endpoint"`
	IAMServerIdHeaderValue string `json:"iam_server_id_header_value" structs:"iam_server_id_header_value" mapstructure:"iam_server_id_header_value"`
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
