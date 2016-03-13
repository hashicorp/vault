package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"domain": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The domain users must be part of",
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Duration after which authentication will be expired`,
			},
			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Maximum duration after which authentication will be expired`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
		},
	}
}

func (b *backend) pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	domain := data.Get("domain").(string)

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
		Domain:     domain,
		TTL:     ttl,
		MaxTTL:  maxTTL,
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

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

type config struct {
	Domain string        `json:"domain"`
	TTL    time.Duration `json:"ttl"`
	MaxTTL time.Duration `json:"max_ttl"`
}
