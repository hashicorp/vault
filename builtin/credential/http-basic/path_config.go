package httpBasic

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "HTTP authentication URL",
			},
			"unregistered_user_policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "Comma-separated list of policies to grant upon successful HTTP authentication of an unregisted user (default: empty)",
			},
			"timeout": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     10,
				Description: "Number of seconds before request times out (default: 10)",
			},
		},

		ExistenceCheck: b.configExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.CreateOperation: b.pathConfigCreateUpdate,
			logical.UpdateOperation: b.pathConfigCreateUpdate,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) configExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.Config(req)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

/*
 * Construct ConfigEntry struct using stored configuration.
 */
func (b *backend) Config(req *logical.Request) (*ConfigEntry, error) {

	storedConfig, err := req.Storage.Get("config")
	if err != nil {
		return nil, err
	}

	if storedConfig == nil {
		return nil, nil
	}

	var result ConfigEntry

	if err := storedConfig.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(cfg).Map(),
	}
	return resp, nil
}

func (b *backend) pathConfigCreateUpdate(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	// Build a ConfigEntry struct out of the supplied FieldData
	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &ConfigEntry{}
	}

	url, ok := d.GetOk("url")
	if ok {
		cfg.Url = url.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.Url = d.Get("url").(string)
	}
	if cfg.Url == "" {
		return logical.ErrorResponse("config parameter `url` cannot be empty"), nil
	}

	policies := make([]string, 0)
	unregisteredUserPoliciesRaw, ok := d.GetOk("unregistered_user_policies")
	if ok {
		unregisteredUserPoliciesStr := unregisteredUserPoliciesRaw.(string)
		if strings.TrimSpace(unregisteredUserPoliciesStr) != "" {
			policies = strings.Split(unregisteredUserPoliciesStr, ",")
			for _, policy := range policies {
				if policy == "root" {
					return logical.ErrorResponse("root policy cannot be granted by an authentication backend"), nil
				}
			}
		}
		cfg.UnregisteredUserPolicies = policies
	} else if req.Operation == logical.CreateOperation {
		cfg.UnregisteredUserPolicies = policies
	}

	timeout, ok := d.GetOk("timeout")
	if ok {
		cfg.Timeout = timeout.(int)
	} else if req.Operation == logical.CreateOperation {
		cfg.Timeout = d.Get("timeout").(int)
	}

	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type ConfigEntry struct {
	Url                      string   `json:"url" structs:"url" mapstructure:"url"`
	Timeout                  int      `json:"timeout" structs:"timeout" mapstructure:"timeout"`
	UnregisteredUserPolicies []string `json:"unregistered_user_policies" structs:"unregistered_user_policies" mapstructure:"unregistered_user_policies"`
}

const pathConfigHelpSyn = `
Configure the HTTP server to connect to, along with its options.
`

const pathConfigHelpDesc = `
This endpoint allows you to configure the HTTP server to connect to and its
configuration options.
`
