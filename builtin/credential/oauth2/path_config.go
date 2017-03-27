package oauth2

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	goauth2 "golang.org/x/oauth2"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"client_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Client ID Vault should use to authenticate with the oauth2 provider.",
			},
			"client_secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Client Secret Vault should use to authenticate with the oauth2 provider.",
			},
			"provider_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Oauth2 API endpoint to use to authenticate users.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.CreateOperation: b.pathConfigWrite,
			logical.UpdateOperation: b.pathConfigWrite,
		},

		ExistenceCheck: b.pathConfigExistenceCheck,

		HelpSynopsis: pathConfigHelp,
	}
}

// Config returns the configuration for this backend.
func (b *backend) Config(s logical.Storage) (*ConfigEntry, error) {
	entry, err := s.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result ConfigEntry
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	// Don't reveal the client secret
	resp := &logical.Response{
		Data: map[string]interface{}{
			"ClientID":    cfg.ClientID,
			"ProviderURL": cfg.ProviderURL,
		},
	}

	return resp, nil
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	// Due to the existence check, entry will only be nil if it's a create
	// operation, so just create a new one
	if cfg == nil {
		cfg = &ConfigEntry{}
	}

	// Client ID & Secret Vault should use to authenticate with oauth2 provider
	cfg.ClientID = d.Get("client_id").(string)
	cfg.ClientSecret = d.Get("client_secret").(string)

	// URL of provider
	providerURL, ok := d.GetOk("provider_url")
	if ok {
		providerURLString := providerURL.(string)
		if len(providerURLString) != 0 {
			_, err = url.Parse(providerURLString)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("Error parsing given provider_url: %s", err)), nil
			}
			cfg.ProviderURL = providerURLString
		}
	} else if req.Operation == logical.CreateOperation {
		cfg.ProviderURL = d.Get("provider_url").(string)
	}

	jsonCfg, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(jsonCfg); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigExistenceCheck(
	req *logical.Request, d *framework.FieldData) (bool, error) {
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return false, err
	}

	return cfg != nil, nil
}

// OauthConfig creates a basic oauth2 client Config
func (c *ConfigEntry) OauthConfig() *goauth2.Config {
	config := &goauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint: goauth2.Endpoint{
			TokenURL: c.ProviderURL,
		},
		//Scopes: []string{},
	}
	return config
}

// Vault ConfigEntry for oauth2
type ConfigEntry struct {
	ProviderURL  string `json:"provider_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

const pathConfigHelp = `
This endpoint allows you to configure the oauth2 backend.
`
