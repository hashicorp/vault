package consul

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigAccess() *framework.Path {
	return &framework.Path{
		Pattern: "config/access",
		Fields: map[string]*framework.FieldSchema{
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Consul server address",
			},

			"scheme": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "URI scheme for the Consul address",

				// https would be a better default but Consul on its own
				// defaults to HTTP access, and when HTTPS is enabled it
				// disables HTTP, so there isn't really any harm done here.
				Default: "http",
			},

			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   pathConfigAccessRead,
			logical.UpdateOperation: pathConfigAccessWrite,
		},
	}
}

func readConfigAccess(storage logical.Storage) (*accessConfig, error, error) {
	entry, err := storage.Get("config/access")
	if err != nil {
		return nil, nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf(
				"Access credentials for the backend itself haven't been configured. Please configure them at the '/config/access' endpoint"),
			nil
	}

	conf := &accessConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, nil, fmt.Errorf("error reading consul access configuration: %s", err)
	}

	return conf, nil, nil
}

func pathConfigAccessRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, userErr, intErr := readConfigAccess(req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}
	if conf == nil {
		return nil, fmt.Errorf("no user error reported but consul access configuration not found")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": conf.Address,
			"scheme":  conf.Scheme,
		},
	}, nil
}

func pathConfigAccessWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/access", accessConfig{
		Address: data.Get("address").(string),
		Scheme:  data.Get("scheme").(string),
		Token:   data.Get("token").(string),
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
	Scheme  string `json:"scheme"`
	Token   string `json:"token"`
}
