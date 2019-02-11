package github

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
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
				DisplayName: "Base URL",
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Duration after which authentication will be expired`,
				DisplayName: "TTL",
			},
			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Maximum duration after which authentication will be expired`,
				DisplayName: "Max TTL",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
	}
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

	var ttl time.Duration
	var err error
	ttlRaw, ok := data.GetOk("ttl")
	if !ok || len(ttlRaw.(string)) == 0 {
		ttl = 0
	} else {
		ttl, err = time.ParseDuration(ttlRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid 'ttl':%s", err)), nil
		}
	}

	var maxTTL time.Duration
	maxTTLRaw, ok := data.GetOk("max_ttl")
	if !ok || len(maxTTLRaw.(string)) == 0 {
		maxTTL = 0
	} else {
		maxTTL, err = time.ParseDuration(maxTTLRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid 'max_ttl':%s", err)), nil
		}
	}

	entry, err := logical.StorageEntryJSON("config", config{
		Organization: organization,
		BaseURL:      baseURL,
		TTL:          ttl,
		MaxTTL:       maxTTL,
	})

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

	config.TTL /= time.Second
	config.MaxTTL /= time.Second

	resp := &logical.Response{
		Data: map[string]interface{}{
			"organization": config.Organization,
			"base_url":     config.BaseURL,
			"ttl":          config.TTL,
			"max_ttl":      config.MaxTTL,
		},
	}
	return resp, nil
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
	Organization string        `json:"organization" structs:"organization" mapstructure:"organization"`
	BaseURL      string        `json:"base_url" structs:"base_url" mapstructure:"base_url"`
	TTL          time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	MaxTTL       time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
}
