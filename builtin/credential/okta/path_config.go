package okta

import (
	"fmt"
	"net/url"
	"strings"

	"time"

	"github.com/chrismalek/oktasdk-go/okta"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"organization": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Okta organization to authenticate against (DEPRECATED)",
			},
			"org_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the organization to be used in the Okta API.",
			},
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Okta admin API token (DEPRECATED)",
			},
			"api_token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Okta API key.",
			},
			"base_url": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The API endpoint to use. Useful if you
are using Okta development accounts. (DEPRECATED)`,
			},
			"production": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: `If set, production API URL prefix will be used to communicate with Okta and if not set, a preview production API URL prefix will be used. Defaults to true.`,
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: `Duration after which authentication will be expired`,
			},
			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: `Maximum duration after which authentication will be expired`,
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

	resp := &logical.Response{
		Data: map[string]interface{}{
			"organization": cfg.Org,
			"org_name":     cfg.Org,
			"production":   *cfg.Production,
			"ttl":          cfg.TTL,
			"max_ttl":      cfg.MaxTTL,
		},
	}
	if cfg.BaseURL != "" {
		resp.Data["base_url"] = cfg.BaseURL
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

	org, ok := d.GetOk("org_name")
	if ok {
		cfg.Org = org.(string)
	}
	if cfg.Org == "" {
		org, ok = d.GetOk("organization")
		if ok {
			cfg.Org = org.(string)
		}
	}
	if cfg.Org == "" && req.Operation == logical.CreateOperation {
		return logical.ErrorResponse("org_name is missing"), nil
	}

	token, ok := d.GetOk("api_token")
	if ok {
		cfg.Token = token.(string)
	}
	if cfg.Token == "" {
		token, ok = d.GetOk("token")
		if ok {
			cfg.Token = token.(string)
		}
	}
	if cfg.Token == "" && req.Operation == logical.CreateOperation {
		return logical.ErrorResponse("api_token is missing"), nil
	}

	baseURL, ok := d.GetOk("base_url")
	if ok {
		baseURLString := baseURL.(string)
		if len(baseURLString) != 0 {
			_, err = url.Parse(baseURLString)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("Error parsing given base_url: %s", err)), nil
			}
			cfg.BaseURL = baseURLString
		}
	} else if req.Operation == logical.CreateOperation {
		cfg.BaseURL = d.Get("base_url").(string)
	}

	productionRaw := d.Get("production").(bool)
	cfg.Production = &productionRaw

	ttl, ok := d.GetOk("ttl")
	if ok {
		cfg.TTL = time.Duration(ttl.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		cfg.TTL = time.Duration(d.Get("ttl").(int)) * time.Second
	}

	maxTTL, ok := d.GetOk("max_ttl")
	if ok {
		cfg.MaxTTL = time.Duration(maxTTL.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		cfg.MaxTTL = time.Duration(d.Get("max_ttl").(int)) * time.Second
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

// OktaClient creates a basic okta client connection
func (c *ConfigEntry) OktaClient() *okta.Client {
	production := true
	if c.Production != nil {
		production = *c.Production
	}
	if c.BaseURL != "" {
		if strings.Contains(c.BaseURL, "oktapreview.com") {
			production = false
		}
	}
	return okta.NewClient(cleanhttp.DefaultClient(), c.Org, c.Token, production)
}

// ConfigEntry for Okta
type ConfigEntry struct {
	Org        string        `json:"organization"`
	Token      string        `json:"token"`
	BaseURL    string        `json:"base_url"`
	Production *bool         `json:"is_production,omitempty"`
	TTL        time.Duration `json:"ttl"`
	MaxTTL     time.Duration `json:"max_ttl"`
}

const pathConfigHelp = `
This endpoint allows you to configure the Okta and its
configuration options.

The Okta organization are the characters at the front of the URL for Okta.
Example https://ORG.okta.com
`
