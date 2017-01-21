package chefnode

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"base_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `The url to the chef server api endpoint`,
			},
			"client_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Name of the client to connect to chef server with`,
			},
			"client_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `PEM encoded client key to use for authenticating to chef server`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.UpdateOperation: b.pathConfigWrite,
		},
	}
}

func (b *backend) pathConfigRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(cfg).Map(),
	}
	resp.AddWarning("Read access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")
	return resp, nil
}

func (b *backend) pathConfigWrite(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	baseURL := data.Get("base_url").(string)
	clientName := data.Get("client_name").(string)
	clientKey := data.Get("client_key").(string)

	entry, err := logical.StorageEntryJSON("config", config{
		BaseURL:    baseURL,
		ClientName: clientName,
		ClientKey:  clientKey,
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) Config(s logical.Storage) (*config, error) {
	entry, err := s.Get("config")
	if err != nil {
		return nil, err
	}

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

type config struct {
	BaseURL    string `json:"base_url"`
	ClientKey  string `json:"client_key"`
	ClientName string `json:"client_name"`
}
