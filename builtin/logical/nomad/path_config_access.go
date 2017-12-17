package nomad

import (
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const configAccessKey = "config/access"

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
			logical.CreateOperation: b.pathConfigAccessWrite,
			logical.UpdateOperation: b.pathConfigAccessWrite,
			logical.DeleteOperation: b.pathConfigAccessDelete,
		},

		ExistenceCheck: b.configExistenceCheck,
	}
}

func (b *backend) configExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.readConfigAccess(req.Storage)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

func (b *backend) readConfigAccess(storage logical.Storage) (*accessConfig, error) {
	entry, err := storage.Get(configAccessKey)
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
	conf, err := b.readConfigAccess(req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		conf = &accessConfig{}
	}

	address, ok := data.GetOk("address")
	if ok {
		conf.Address = address.(string)
	}
	token, ok := data.GetOk("token")
	if ok {
		conf.Token = token.(string)
	}

	entry, err := logical.StorageEntryJSON("config/access", conf)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigAccessDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(configAccessKey); err != nil {
		return nil, err
	}
	return nil, nil
}

type accessConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}
