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
			"provider_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Oauth2 API endpoint to use to authenticate users.",
			},
			"userinfo_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The URL to query after authenticating a user to get any group memberships",
			},
			"client_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Client ID Vault should use to authenticate with the oauth2 provider.",
			},
			"client_secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Client Secret Vault should use to authenticate with the oauth2 provider.",
			},
			"scope": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Oauth Scope that will provide access to group membership when making a request to the UserInfoURL.",
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
			"ProviderURL": cfg.ProviderURL,
			"UserInfoURL": cfg.UserInfoURL,
			"ClientID":    cfg.ClientID,
			"Scope":       cfg.Scope,
		},
	}

	return resp, nil
}

func (b *backend) pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// URL of provider
	providerURL := data.Get("provider_url").(string)
	if len(providerURL) == 0 {
		return logical.ErrorResponse("A provider_url must be specified."), nil
	}
	_, err := url.Parse(providerURL)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("Error parsing given provider_url: %s", err)), nil
	}

	// URL of userinfo endpoint to query groups.  If blank, only local
	// groups will be used.
	userinfoURL := data.Get("userinfo_url").(string)
	if len(userinfoURL) != 0 {
		_, err = url.Parse(userinfoURL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing given userinfo_url: %s", err)), nil
		}
	}

	// Client ID, Secret, & Scope Vault should use to authenticate with oauth2 provider
	clientID := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)
	scope := data.Get("scope").(string)

	entry, err := logical.StorageEntryJSON("config", ConfigEntry{
		ProviderURL:  providerURL,
		UserInfoURL:  userinfoURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	})
	if err != nil {
		return nil, err
	}
	if err = req.Storage.Put(entry); err != nil {
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
		Scopes: []string{c.Scope},
	}
	return config
}

// Vault ConfigEntry for oauth2
type ConfigEntry struct {
	ProviderURL  string `json:"provider_url"`
	UserInfoURL  string `json:"userinfo_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
}

const pathConfigHelp = `
This endpoint allows you to configure the oauth2 backend.
`
