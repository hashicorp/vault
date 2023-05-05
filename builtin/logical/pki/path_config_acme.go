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
	pathConfigAcmeHelpDesc = "Here we configure:\n\nenabled=false, whether ACME is enabled, defaults to false meaning that clusters will by default not get ACME support,\nallowed_issuers=\"default\", which issuers are allowed for use with ACME; by default, this will only be the primary (default) issuer,\nallowed_roles=\"*\", which roles are allowed for use with ACME; by default these will be all roles matching our selection criteria,\ndefault_role=\"\", if not empty, the role to be used for non-role-qualified ACME requests; by default this will be empty, meaning ACME issuance will be equivalent to sign-verbatim.,\ndns_resolver=\"\", which specifies a custom DNS resolver to use for all ACME-related DNS lookups"
)

type acmeConfigEntry struct {
	Enabled        bool          `json:"enabled"`
	AllowedIssuers []string      `json:"allowed_issuers="`
	AllowedRoles   []string      `json:"allowed_roles"`
	DefaultRole    string        `json:"default_role"`
	DNSResolver    string        `json:"dns_resolver"`
	EabPolicyName  EabPolicyName `json:"eab_policy_name"`
}

var defaultAcmeConfig = acmeConfigEntry{
	Enabled:        false,
	AllowedIssuers: []string{"*"},
	AllowedRoles:   []string{"*"},
	DefaultRole:    "",
	DNSResolver:    "",
	EabPolicyName:  eabPolicyAlwaysRequired,
}

func (sc *storageContext) getAcmeConfig() (*acmeConfigEntry, error) {
	entry, err := sc.Storage.Get(sc.Context, storageAcmeConfig)
	if err != nil {
		return nil, err
	}

	var mapping acmeConfigEntry
	if entry == nil {
		mapping = defaultAcmeConfig
		return &mapping, nil
	}

	if err := entry.DecodeJSON(&mapping); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode ACME configuration: %v", err)}
	}

	return &mapping, nil
}

func (sc *storageContext) setAcmeConfig(entry *acmeConfigEntry) error {
	json, err := logical.StorageEntryJSON(storageAcmeConfig, entry)
	if err != nil {
		return fmt.Errorf("failed creating storage entry: %w", err)
	}

	if err := sc.Storage.Put(sc.Context, json); err != nil {
		return fmt.Errorf("failed writing storage entry: %w", err)
	}

	sc.Backend.acmeState.markConfigDirty()
	return nil
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
				Default:     "*",
			},
			"allowed_roles": {
				Type:        framework.TypeCommaStringSlice,
				Description: `which roles are allowed for use with ACME; by default via '*', these will be all roles including sign-verbatim; when concrete role names are specified, sign-verbatim is not allowed and a default_role must be specified in order to allow usage of the default acme directories under /pki/acme/directory and /pki/issuer/:issuer_id/acme/directory.`,
				Default:     "*",
			},
			"default_role": {
				Type:        framework.TypeString,
				Description: `if not empty, the role to be used for non-role-qualified ACME requests; by default this will be empty, meaning ACME issuance will be equivalent to sign-verbatim; must be specified in allowed_roles if non-empty`,
				Default:     "",
			},
			"dns_resolver": {
				Type:        framework.TypeString,
				Description: `DNS resolver to use for domain resolution on this mount. Defaults to using the default system resolver. Must be in the format <host>:<port>, with both parts mandatory.`,
				Default:     "",
			},
			"eab_policy": {
				Type:        framework.TypeString,
				Description: `Specify the policy to use for external account binding behaviour, 'not-required', 'new-account-required' or 'always-required'`,
				Default:     "always-required",
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
			"allowed_roles":   config.AllowedRoles,
			"allowed_issuers": config.AllowedIssuers,
			"default_role":    config.DefaultRole,
			"enabled":         config.Enabled,
			"dns_resolver":    config.DNSResolver,
			"eab_policy":      config.EabPolicyName,
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

	if allowedRolesRaw, ok := d.GetOk("allowed_roles"); ok {
		config.AllowedRoles = allowedRolesRaw.([]string)
		if len(config.AllowedRoles) == 0 {
			return nil, fmt.Errorf("allowed_roles must take a non-zero length value; specify '*' as the value to allow anything or specify enabled=false to disable ACME entirely")
		}
	}

	if defaultRoleRaw, ok := d.GetOk("default_role"); ok {
		config.DefaultRole = defaultRoleRaw.(string)
	}

	if allowedIssuersRaw, ok := d.GetOk("allowed_issuers"); ok {
		config.AllowedIssuers = allowedIssuersRaw.([]string)
		if len(config.AllowedIssuers) == 0 {
			return nil, fmt.Errorf("allowed_issuers must take a non-zero length value; specify '*' as the value to allow anything or specify enabled=false to disable ACME entirely")
		}
	}

	if dnsResolverRaw, ok := d.GetOk("dns_resolver"); ok {
		config.DNSResolver = dnsResolverRaw.(string)
		if config.DNSResolver != "" {
			addr, _, err := net.SplitHostPort(config.DNSResolver)
			if err != nil {
				return nil, fmt.Errorf("failed to parse DNS resolver address: %w", err)
			}
			if addr == "" {
				return nil, fmt.Errorf("failed to parse DNS resolver address: got empty address")
			}
			if net.ParseIP(addr) == nil {
				return nil, fmt.Errorf("failed to parse DNS resolver address: expected IPv4/IPv6 address, likely got hostname")
			}
		}
	}

	if eabPolicyRaw, ok := d.GetOk("eab_policy"); ok {
		eabPolicy, err := getEabPolicyByString(eabPolicyRaw.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid eab policy name provided")
		}
		config.EabPolicyName = eabPolicy.Name
	}

	allowAnyRole := len(config.AllowedRoles) == 1 && config.AllowedRoles[0] == "*"
	if !allowAnyRole {
		foundDefault := len(config.DefaultRole) == 0
		for index, name := range config.AllowedRoles {
			if name == "*" {
				return nil, fmt.Errorf("cannot use '*' as role name at index %d", index)
			}

			role, err := sc.Backend.getRole(sc.Context, sc.Storage, name)
			if err != nil {
				return nil, fmt.Errorf("failed validating allowed_roles: unable to fetch role: %v: %w", name, err)
			}

			if role == nil {
				return nil, fmt.Errorf("role %v specified in allowed_roles does not exist", name)
			}

			if name == config.DefaultRole {
				foundDefault = true
			}
		}

		if !foundDefault {
			return nil, fmt.Errorf("default role %v was not specified in allowed_roles: %v", config.DefaultRole, config.AllowedRoles)
		}
	}

	allowAnyIssuer := len(config.AllowedIssuers) == 1 && config.AllowedIssuers[0] == "*"
	if !allowAnyIssuer {
		for index, name := range config.AllowedIssuers {
			if name == "*" {
				return nil, fmt.Errorf("cannot use '*' as issuer name at index %d", index)
			}

			_, err := sc.resolveIssuerReference(name)
			if err != nil {
				return nil, fmt.Errorf("failed validating allowed_issuers: unable to fetch issuer: %v: %w", name, err)
			}
		}
	}

	err = sc.setAcmeConfig(config)
	if err != nil {
		return nil, err
	}

	return genResponseFromAcmeConfig(config), nil
}
