package cert

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"verify_cert": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: `If set, during renewal, checks if the presented client certificate matches the client certificate used during login. Defaults to true.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
		},
	}
}

func (b *backend) pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	verifyCert := data.Get("verify_cert").(bool)

	entry, err := logical.StorageEntryJSON("config", config{
		VerifyCert: verifyCert,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(s logical.Storage) (*config, error) {
	entry, err := s.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result config
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("error reading configuration: %s", err)
	}
	return &result, nil
}

type config struct {
	VerifyCert bool `json:"verify_cert"`
}
