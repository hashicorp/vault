package ldap

import (
	"net/url"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "ldap URL to connect to (default: ldap://127.0.0.1)",
			},
			"domain": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "LDAP domain to use (eg: dc=example,dc=org)",
			},
			"userattr": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Attribute used for users (default: cn)",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:  b.pathConfigRead,
			logical.WriteOperation: b.pathConfigWrite,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

func (b *backend) Config(req *logical.Request) (*ConfigEntry, error) {
	entry, err := req.Storage.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var result ConfigEntry
	result.SetDefaults()
	if err := entry.DecodeJSON(&result); err != nil {
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

	return &logical.Response{
		Data: map[string]interface{}{
			"url":      cfg.Url,
			"domain":   cfg.Domain,
			"userattr": cfg.UserAttr,
		},
	}, nil
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg := &ConfigEntry{}
	url := d.Get("url").(string)
	if url != "" {
		cfg.Url = strings.ToLower(url)
	}
	userattr := d.Get("userattr").(string)
	if url != "" {
		cfg.UserAttr = strings.ToLower(userattr)
	}
	domain := d.Get("domain").(string)
	if url != "" {
		cfg.Domain = domain
	}

	if !cfg.ValidateURL() {
		return logical.ErrorResponse("LDAP URL is malformed"), nil
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
	Url      string
	Domain   string
	UserAttr string
}

func (c *ConfigEntry) ValidateURL() bool {
	u, err := url.Parse(c.Url)
	if err != nil {
		return false
	}
	if u.Scheme != "ldap" && u.Scheme != "ldaps" {
		return false
	}
	if u.Path != "" {
		return false
	}
	return true
}

func (c *ConfigEntry) SetDefaults() {
	c.Url = "ldap://127.0.0.1"
	c.UserAttr = "cn"
}

const pathConfigHelpSyn = `
Configure the LDAP server to connect to.
`

const pathConfigHelpDesc = `
This endpoint allows you to configure the LDAP server to connect to, and give
basic information of the schema of that server.

The LDAP URL can use either the "ldap://" or "ldaps://" schema. In the former
case, an unencrypted connection will be done, with default port 389; in the latter
case, a SSL connection will be done, with default port 636.
`
