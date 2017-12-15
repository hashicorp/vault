package nomad

import (
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",
		Fields: map[string]*framework.FieldSchema{
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nomad server address",
			},

			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigAccessRead,
			logical.UpdateOperation: b.pathConfigAccessWrite,
		},
	}
}

func (b *backend) readConfigAccess(storage logical.Storage) (*accessConfig, error) {
	entry, err := storage.Get("config/access")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	conf := &accessConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, errwrap.Wrapf("error reading nomad access configuration: {{err}}", err)
	}

	return conf, nil
}

func (b *backend) pathConfigAccessRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readConfigAccess(req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": conf.Address,
		},
	}, nil
}

func (b *backend) pathConfigAccessWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	address := data.Get("address").(string)
	if address == "" {
		return logical.ErrorResponse("missing nomad server address"), nil
	}
	token := data.Get("token").(string)
	if token == "" {
		return logical.ErrorResponse("missing nomad management token"), nil
	}
	entry, err := logical.StorageEntryJSON("config/access", accessConfig{
		Address: address,
		Token:   token,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type accessConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}
