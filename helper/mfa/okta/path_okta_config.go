package okta

import (
	"fmt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"net/url"
)

func pathOktaConfig() *framework.Path {
	return &framework.Path{
		Pattern: `okta/config`,
		Fields: map[string]*framework.FieldSchema{
			"api_token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "API token to connect to Okta (Required)",
			},
			"base_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Base URL for Okta (Required)",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: pathOktaConfigWrite,
			logical.ReadOperation:   pathOktaConfigRead,
		},

		HelpSynopsis:    pathOktaConfigHelpSyn,
		HelpDescription: pathOktaConfigHelpDesc,
	}
}

func GetOktaConfig(req *logical.Request) (*OktaConfig, error) {
	var result OktaConfig
	// all config parameters are optional, so path need not exist
	entry, err := req.Storage.Get("okta/config")
	if err == nil && entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, err
		}
	}
	return &result, nil
}

func pathOktaConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	api_token := d.Get("api_token").(string)
	base_url := d.Get("base_url").(string)

	if base_url == "" {
		return nil, fmt.Errorf("missing base_url")
	}

	url_base_url, err := url.ParseRequestURI(base_url)

	if err != nil {
		return nil, fmt.Errorf("unable to parse base_url")
	}

	if api_token == "" {
		return nil, fmt.Errorf("missing api_token")
	}

	entry, err := logical.StorageEntryJSON("okta/config", OktaConfig{
		ApiToken: api_token,
		BaseURL:  url_base_url,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func pathOktaConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	config, err := GetOktaConfig(req)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"api_token": config.ApiToken,
			"base_url":  config.BaseURL,
		},
	}, nil
}

type OktaConfig struct {
	ApiToken string   `json:"api_token"`
	BaseURL  *url.URL `json:"base_url"`
}

const pathOktaConfigHelpSyn = `
Configure Okta second factor behavior.
`

const pathOktaConfigHelpDesc = `
This endpoint allows you to configure the Okta second factor authentication.
`
