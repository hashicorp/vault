package github

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"organization": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The organization users must be part of",
			},

			"base_url": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The API endpoint to use. Useful if you
are running GitHub Enterprise or an
API-compatible authentication server.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathConfigWrite,
		},
	}
}

func pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf := config{
		Org: data.Get("organization").(string),
	}
	baseURL := data.Get("base_url").(string)
	if len(baseURL) != 0 {
		_, err := url.Parse(baseURL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing given base_url: %s", err)), nil
		}
		conf.BaseURL = baseURL
	}

	entry, err := logical.StorageEntryJSON("config", conf)
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

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

type config struct {
	Org     string `json:"organization"`
	BaseURL string `json:"base_url"`
}
