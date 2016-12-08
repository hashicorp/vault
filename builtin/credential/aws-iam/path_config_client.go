package awsiam

import (
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigClient(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/client$",
		Fields: map[string]*framework.FieldSchema{
			"endpoint": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "API endpoint to submit GetCallerIdentity requests to",
			},
			"vault_header_value": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "Value to require in the X-Vault-AWSIAM-Server-ID request header",
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

func (b *backend) lockedClientConfigEntry(s logical.Storage) (*clientConfig, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedClientConfigEntry(s)
}

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

func (b *backend) pathConfigClientExistenceCheck(
	req *logical.Request, data *framework.FieldData) (bool, error) {

	entry, err := b.lockedClientConfigEntry(req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
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

	endpointStr, ok := data.GetOk("endpoint")
	if ok {
		configEntry.Endpoint = endpointStr.(string)
	} else if req.Operation == logical.CreateOperation {
		configEntry.Endpoint = data.Get("endpoint").(string)
	}

	headerValueStr, ok := data.GetOk("vault_header_value")
	if ok {
		configEntry.HeaderValue = headerValueStr.(string)
	} else if req.Operation == logical.CreateOperation {
		configEntry.HeaderValue = data.Get("vault_header_value").(string)
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

func (b *backend) pathConfigClientDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return nil, req.Storage.Delete("config/client")
}

type clientConfig struct {
	Endpoint    string `json:"endpoint" structs:"endpoint" mapstructure:"endpoint"`
	HeaderValue string `json:"vault_header_value" structs:"vault_header_value" mapstructure:"vault_header_value"`
}

const pathConfigClientHelpSyn = `
Configures the STS client used to validate requests.
`

const pathConfigClientHelpDesc = `
aws-iam auth backend validates signed GetCallerIdentity requests to determine
the AWS IAM principal making the requests. The endpoint allows you to specify
which STS endpoint to forward the signed query along to, while the
vault_header_value allows you to specify a required value to be included in the
signed headers to mitigate certain types of attacks.
`
