package centrify

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"client_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "OAuth2 Client ID",
			},
			"client_secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "OAuth2 Client Secret",
			},
			"service_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Service URL (https://<tenant>.my.centrify.com)",
			},
			"app_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "OAuth2 App ID",
				Default:     "vault_io_integration",
			},
			"scope": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "OAuth2 App Scope",
				Default:     "vault_io_integration",
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
				Deprecated:  true,
			},
		},

		ExistenceCheck: b.pathConfigExistCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigCreateOrUpdate,
			logical.CreateOperation: b.pathConfigCreateOrUpdate,
			logical.ReadOperation:   b.pathConfigRead,
		},

		HelpSynopsis: pathSyn,
	}

	tokenutil.AddTokenFieldsWithAllowList(p.Fields, []string{
		"token_bound_cidrs",
		"token_no_default_policy",
		"token_policies",
		"token_type",
		"token_ttl",
		"token_num_uses",
	})
	return p
}

func (b *backend) pathConfigExistCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	if config == nil {
		return false, nil
	}

	return true, nil
}

func (b *backend) pathConfigCreateOrUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = &config{}
	}

	val, ok := data.GetOk("client_id")
	if ok {
		cfg.ClientID = val.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.ClientID = data.Get("client_id").(string)
	}
	if cfg.ClientID == "" {
		return logical.ErrorResponse("config parameter `client_id` cannot be empty"), nil
	}

	val, ok = data.GetOk("client_secret")
	if ok {
		cfg.ClientSecret = val.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.ClientSecret = data.Get("client_secret").(string)
	}
	if cfg.ClientSecret == "" {
		return logical.ErrorResponse("config parameter `client_secret` cannot be empty"), nil
	}

	val, ok = data.GetOk("service_url")
	if ok {
		cfg.ServiceURL = val.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.ServiceURL = data.Get("service_url").(string)
	}
	if cfg.ServiceURL == "" {
		return logical.ErrorResponse("config parameter `service_url` cannot be empty"), nil
	}

	val, ok = data.GetOk("app_id")
	if ok {
		cfg.AppID = val.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.AppID = data.Get("app_id").(string)
	}

	val, ok = data.GetOk("scope")
	if ok {
		cfg.Scope = val.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.Scope = data.Get("scope").(string)
	}

	if err := cfg.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	{
		if err := tokenutil.UpgradeValue(data, "policies", "token_policies", &cfg.Policies, &cfg.TokenPolicies); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	// We want to normalize the service url to https://
	url, err := url.Parse(cfg.ServiceURL)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("config parameter 'service_url' is not a valid url: %s", err)), nil
	}

	// Its a proper url, just force the scheme to https, and strip any paths
	url.Scheme = "https"
	url.Path = ""
	cfg.ServiceURL = url.String()

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
		return nil, nil
	}

	d := map[string]interface{}{
		"client_id":   config.ClientID,
		"service_url": config.ServiceURL,
		"app_id":      config.AppID,
		"scope":       config.Scope,
		"policies":    config.Policies,
	}
	config.PopulateTokenData(d)
	delete(d, "token_explicit_max_ttl")
	delete(d, "token_max_ttl")
	delete(d, "token_period")

	if len(config.Policies) > 0 {
		d["policies"] = d["token_policies"]
	}

	return &logical.Response{
		Data: d,
	}, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result config
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("error reading configuration: %s", err)
	}

	if len(result.TokenPolicies) == 0 && len(result.Policies) > 0 {
		result.TokenPolicies = result.Policies
	}

	return &result, nil
}

type config struct {
	tokenutil.TokenParams

	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	ServiceURL   string   `json:"service_url"`
	AppID        string   `json:"app_id"`
	Scope        string   `json:"scope"`
	Policies     []string `json:"policies"`
}

const pathSyn = `
This path allows you to configure the centrify auth provider to interact with the Centrify Identity Services Platform
for authenticating users.  
`
