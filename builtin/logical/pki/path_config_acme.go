package pki

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	storageAcmeConfig      = "config/acme"
	pathConfigAcmeHelpSyn  = "Configuration of ACME Endpoints"
	pathConfigAcmeHelpDesc = "Here we configure:\n\nenabled=false, whether ACME is enabled, defaults to false meaning that clusters will by default not get ACME support,\nallowed_issuers=\"default\", which issuers are allowed for use with ACME; by default, this will only be the primary (default) issuer,\nallowed_roles=\"*\", which roles are allowed for use with ACME; by default these will be all roles matching our selection criteria,\ndefault_role=\"\", if not empty, the role to be used for non-role-qualified ACME requests; by default this will be empty, meaning ACME issuance will be equivalent to sign-verbatim,\nallow_no_allowed_domains=false, whether ACME will allow the use of roles without any explicit list of allowed domains, and\nallow_any_domain=false, whether ACME will allow the use of roles with allow_any_name=true set."
)

type acmeConfigEntry struct {
	Enabled               bool     `json:"enabled"`
	AllowedIssuers        []string `json:"allowed_issuers="`
	AllowedRoles          []string `json:"allowed_roles"`
	DefaultRole           string   `json:"default_role"`
	AllowNoAllowedDomains bool     `json:"allow_no_allowed_domains"`
	AllowAnyDomain        bool     `json:"allow_any_domain"`
	DNSResolver           string   `json:"dns_resolver"`
}

func (sc *storageContext) getAcmeConfig() (*acmeConfigEntry, error) {
	entry, err := sc.Storage.Get(sc.Context, storageAcmeConfig)
	if err != nil {
		return nil, err
	}

	mapping := &acmeConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(mapping); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode ACME configuration: %v", err)}
		}
	}

	return mapping, nil
}

func (sc *storageContext) setAcmeConfig(entry *acmeConfigEntry) error {
	json, err := logical.StorageEntryJSON(storageAcmeConfig, entry)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, json)
}

func pathAcmeConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/acme",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

		Fields: map[string]*framework.FieldSchema{
			"enabled": {
				Type:        framework.TypeBool,
				Description: `whether ACME is enabled, defaults to false meaning that clusters will by default not get ACME support`,
				Default:     false,
			},
			"allowed_issuers": {
				Type:        framework.TypeCommaStringSlice,
				Description: `which issuers are allowed for use with ACME; by default, this will only be the primary (default) issuer`,
				Default:     "default",
			},
			"allowed_roles": {
				Type:        framework.TypeCommaStringSlice,
				Description: `which roles are allowed for use with ACME; by default these will be all roles matching our selection criteria`,
				Default:     "*",
			},
			"default_role": {
				Type:        framework.TypeString,
				Description: `if not empty, the role to be used for non-role-qualified ACME requests; by default this will be empty, meaning ACME issuance will be equivalent to sign-verbatim,`,
				Default:     "",
			},
			"allow_no_allowed_domains": {
				Type:        framework.TypeBool,
				Description: `whether ACME will allow the use of roles without any explicit list of allowed domains`,
				Default:     false,
			},
			"allow_any_domain": {
				Type:        framework.TypeBool,
				Description: `whether ACME will allow the use of roles with allow_any_name=true set.`,
			},
			"dns_resolver": {
				Type:        framework.TypeString,
				Description: `DNS resolver to use for domain resolution on this mount. Defaults to using the default system resolver. Must be in the format <host>:<port>, with both parts mandatory.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "acme-configuration",
				},
				Callback: b.pathAcmeRead,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathAcmeWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "acme",
				},
				// Read more about why these flags are set in backend.go.
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathConfigAcmeHelpSyn,
		HelpDescription: pathConfigAcmeHelpDesc,
	}
}

func (b *backend) pathAcmeRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := sc.getAcmeConfig()
	if err != nil {
		return nil, err
	}

	return genResponseFromAcmeConfig(config), nil
}

func genResponseFromAcmeConfig(config *acmeConfigEntry) *logical.Response {
	response := &logical.Response{
		Data: map[string]interface{}{
			"allow_any_domain":         config.AllowAnyDomain,
			"allowed_roles":            config.AllowedRoles,
			"allow_no_allowed_domains": config.AllowNoAllowedDomains,
			"allowed_issuers":          config.AllowedIssuers,
			"default_role":             config.DefaultRole,
			"enabled":                  config.Enabled,
			"dns_resolver":             config.DNSResolver,
		},
	}

	// TODO: Add some nice warning if we are on a replication cluster and path isn't set

	return response
}

func (b *backend) pathAcmeWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)

	config, err := sc.getAcmeConfig()
	if err != nil {
		return nil, err
	}

	if enabledRaw, ok := d.GetOk("enabled"); ok {
		config.Enabled = enabledRaw.(bool)
	}

	if allowAnyDomainRaw, ok := d.GetOk("allow_any_domain"); ok {
		config.AllowAnyDomain = allowAnyDomainRaw.(bool)
	}

	if allowedRolesRaw, ok := d.GetOk("allowed_roles"); ok {
		config.AllowedRoles = allowedRolesRaw.([]string)
	}

	if defaultRoleRaw, ok := d.GetOk("default_role"); ok {
		config.DefaultRole = defaultRoleRaw.(string)
	}

	if allowNoAllowedDomainsRaw, ok := d.GetOk("allow_no_allowed_domains"); ok {
		config.AllowNoAllowedDomains = allowNoAllowedDomainsRaw.(bool)
	}

	if allowedIssuersRaw, ok := d.GetOk("allowed_issuers"); ok {
		config.AllowedIssuers = allowedIssuersRaw.([]string)
	}

	if dnsResolverRaw, ok := d.GetOk("dns_resolver"); ok {
		config.DNSResolver = dnsResolverRaw.(string)
		if config.DNSResolver != "" {
			addr, _, err := net.SplitHostPort(config.DNSResolver)
			if err != nil {
				return nil, fmt.Errorf("failed to parse DNS resolver address: %w", err)
			}
			if addr != "" {
				return nil, fmt.Errorf("failed to parse DNS resolver address: got empty address")
			}
			if net.ParseIP(addr) != nil {
				return nil, fmt.Errorf("failed to parse DNS resolver address: expected IPv4/IPv6 address, likely got hostname")
			}
		}
	}

	err = sc.setAcmeConfig(config)
	if err != nil {
		return nil, err
	}

	return genResponseFromAcmeConfig(config), nil
}
