package transit

import (
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"allow_upsert": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `Whether to allow upserting keys on first use,
rather than requiring them to be manually
specified through the keys endpoint`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.UpdateOperation: b.pathConfigWrite,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

type transitConfig struct {
	AllowUpsert bool `json:"allow_upsert" structs:"allow_upsert" mapstructure:"allow_upsert"`
}

func (b *backend) getConfig(s logical.Storage) (*transitConfig, error) {
	entry, err := s.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result transitConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.getConfig(req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &transitConfig{}
	}

	resp := &logical.Response{
		Data: structs.New(config).Map(),
	}

	return resp, nil
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	allowUpsertInt, ok := d.GetOk("allow_upsert")
	if !ok {
		return logical.ErrorResponse("no known configuration parameters supplied"), nil
	}

	config := &transitConfig{
		AllowUpsert: allowUpsertInt.(bool),
	}

	jsonEntry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(jsonEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

const pathConfigHelpSyn = `Configure the backend`

const pathConfigHelpDesc = `
This path is used to configure the backend. Currently, this allows configuring
whether or not keys can be created via upsert from the encryption endpoint.`
