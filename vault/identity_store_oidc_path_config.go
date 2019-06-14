package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	oidcConfigFieldIssuer = "issuer"
	oidcConfigStorageKey  = "oidc/config/issuer"
)

type oidcConfig struct {
	Issuer string `json:"issuer"`
}

// oidcPathConfig returns the API endpoint for operations on OIDC configuration
func oidcPathConfig(i *IdentityStore) *framework.Path {
	return &framework.Path{
		Pattern: "oidc/config",
		Fields: map[string]*framework.FieldSchema{
			oidcConfigFieldIssuer: &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Issuer URL to be used in the iss claim of the token. If not set, Vault's app_addr will be used.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: writeOIDCConfig(),
			logical.ReadOperation:   readOIDCConfig(),
			logical.UpdateOperation: writeOIDCConfig(),
		},
		HelpSynopsis:    "HelpSynopsis here",
		HelpDescription: "HelpDecription here",
	}
}

// unsure if the following methods should be on the identityStore or not

// readOIDCConfig returns a framework.OperationFunc for reading OIDC configuration
func readOIDCConfig() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		var c = &oidcConfig{}

		entry, err := req.Storage.Get(ctx, oidcConfigStorageKey)
		if err != nil {
			return nil, err
		}

		if err := entry.DecodeJSON(c); err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"issuer": c.Issuer,
			},
		}, nil
	}
}

// writeOIDCConfig returns a framework.OperationFunc for creating and updating OIDC configuration
func writeOIDCConfig() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		var c = &oidcConfig{}

		value, ok := d.GetOk(oidcConfigFieldIssuer)
		if !ok {
			return nil, fmt.Errorf("OIDC config field not found for %s", oidcConfigFieldIssuer)
		}

		c.Issuer = value.(string)

		entry, err := logical.StorageEntryJSON(oidcConfigStorageKey, c)
		if err != nil {
			return nil, err
		}

		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		return &logical.Response{}, nil
	}
}
