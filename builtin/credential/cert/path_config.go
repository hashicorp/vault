package cert

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"disable_binding": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: `If set, during renewal, skips the matching of presented client identity with the client identity used during login. Defaults to false.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
		},
	}
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	disableBinding := data.Get("disable_binding").(bool)

	entry, err := logical.StorageEntryJSON("config", config{
		DisableBinding: disableBinding,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	// Returning a default configuration if an entry is not found
	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, errwrap.Wrapf("error reading configuration: {{err}}", err)
		}
	}
	return &result, nil
}

type config struct {
	DisableBinding bool `json:"disable_binding"`
}
