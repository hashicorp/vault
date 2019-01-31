package github

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/tokenhelper"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	ret := &framework.Path{
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
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
	}
	tokenhelper.AddTokenFields(ret.Fields)

	return ret
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	organization := data.Get("organization").(string)
	baseURL := data.Get("base_url").(string)
	if len(baseURL) != 0 {
		_, err := url.Parse(baseURL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing given base_url: %s", err)), nil
		}
	}

	cfg := &config{
		Organization: organization,
		BaseURL:      baseURL,
	}

	if err := cfg.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	entry, err := logical.StorageEntryJSON("config", cfg)

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, fmt.Errorf("configuration object not found")
	}

	respData := map[string]interface{}{
		"organization": config.Organization,
		"base_url":     config.BaseURL,
	}
	config.PopulateTokenData(respData)

	return &logical.Response{
		Data: respData,
	}, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, errwrap.Wrapf("error reading configuration: {{err}}", err)
		}
	}

	return &result, nil
}

type config struct {
	tokenhelper.TokenParams

	Organization string `json:"organization" structs:"organization" mapstructure:"organization"`
	BaseURL      string `json:"base_url" structs:"base_url" mapstructure:"base_url"`
}
