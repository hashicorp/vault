// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	storageAcmeConfig      = "config/acme"
	pathConfigAcmeHelpSyn  = "Configuration of ACME Endpoints"
	pathConfigAcmeHelpDesc = "Here we configure:\n\nenabled=false, whether ACME is enabled, defaults to false meaning that clusters will by default not get ACME support,\nallowed_issuers=\"default\", which issuers are allowed for use with ACME; by default, this will only be the primary (default) issuer,\nallowed_roles=\"*\", which roles are allowed for use with ACME; by default these will be all roles matching our selection criteria,\ndefault_directory_policy=\"\", either \"forbid\", preventing the default directory from being used at all, \"role:<role_name>\" which is the role to be used for non-role-qualified ACME requests; or \"sign-verbatim\", the default meaning ACME issuance will be equivalent to sign-verbatim.,\ndns_resolver=\"\", which specifies a custom DNS resolver to use for all ACME-related DNS lookups"
	disableAcmeEnvVar      = "VAULT_DISABLE_PUBLIC_ACME"
	defaultAcmeMaxTTL      = 90 * (24 * time.Hour)
)

type acmeConfigEntry struct {
	Enabled                bool          `json:"enabled"`
	AllowedIssuers         []string      `json:"allowed_issuers="`
	AllowedRoles           []string      `json:"allowed_roles"`
	AllowRoleExtKeyUsage   bool          `json:"allow_role_ext_key_usage"`
	DefaultDirectoryPolicy string        `json:"default_directory_policy"`
	DNSResolver            string        `json:"dns_resolver"`
	EabPolicyName          EabPolicyName `json:"eab_policy_name"`
	MaxTTL                 time.Duration `json:"max_ttl"`
}

var defaultAcmeConfig = acmeConfigEntry{
	Enabled:                false,
	AllowedIssuers:         []string{"*"},
	AllowedRoles:           []string{"*"},
	AllowRoleExtKeyUsage:   false,
	DefaultDirectoryPolicy: "sign-verbatim",
	DNSResolver:            "",
	EabPolicyName:          eabPolicyNotRequired,
	MaxTTL:                 defaultAcmeMaxTTL,
}

var (
	extPolicyPrefix       = "external-policy"
	extPolicyPrefixLength = len(extPolicyPrefix)
	extPolicyRegex        = regexp.MustCompile(framework.GenericNameRegex("policy"))
	rolePrefix            = "role:"
	rolePrefixLength      = len(rolePrefix)
)

func getAcmeConfig(sc *storageContext) (*acmeConfigEntry, error) {
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

	// Update previous stored configurations to use the default max ttl we used to enforce
	if mapping.MaxTTL == 0 {
		mapping.MaxTTL = defaultAcmeMaxTTL
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
				Default:     []string{"*"},
			},
			"allowed_roles": {
				Type:        framework.TypeCommaStringSlice,
				Description: `which roles are allowed for use with ACME; by default via '*', these will be all roles including sign-verbatim; when concrete role names are specified, any default_directory_policy role must be included to allow usage of the default acme directories under /pki/acme/directory and /pki/issuer/:issuer_id/acme/directory.`,
				Default:     []string{"*"},
			},
			"allow_role_ext_key_usage": {
				Type:        framework.TypeBool,
				Description: `whether the ExtKeyUsage field from a role is used, defaults to false meaning that certificate will be signed with ServerAuth.`,
				Default:     false,
			},
			"default_directory_policy": {
				Type:        framework.TypeString,
				Description: `the policy to be used for non-role-qualified ACME requests; by default ACME issuance will be otherwise unrestricted, equivalent to the sign-verbatim endpoint; one may also specify a role to use as this policy, as "role:<role_name>", the specified role must be allowed by allowed_roles`,
				Default:     "sign-verbatim",
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
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: `Specify the maximum TTL for ACME certificates. Role TTL values will be limited to this value`,
				Default:     defaultAcmeMaxTTL.Seconds(),
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
	config, err := b.GetAcmeState().getConfigWithForcedUpdate(sc)
	if err != nil {
		return nil, err
	}

	var warnings []string
	if config.Enabled {
		_, err := getBasePathFromClusterConfig(sc)
		if err != nil {
			warnings = append(warnings, err.Error())
		}
	}

	return genResponseFromAcmeConfig(config, warnings), nil
}

func genResponseFromAcmeConfig(config *acmeConfigEntry, warnings []string) *logical.Response {
	response := &logical.Response{
		Data: map[string]interface{}{
			"allowed_roles":            config.AllowedRoles,
			"allow_role_ext_key_usage": config.AllowRoleExtKeyUsage,
			"allowed_issuers":          config.AllowedIssuers,
			"default_directory_policy": config.DefaultDirectoryPolicy,
			"enabled":                  config.Enabled,
			"dns_resolver":             config.DNSResolver,
			"eab_policy":               config.EabPolicyName,
			"max_ttl":                  config.MaxTTL.Seconds(),
		},
		Warnings: warnings,
	}

	// TODO: Add some nice warning if we are on a replication cluster and path isn't set

	return response
}

func (b *backend) pathAcmeWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)

	config, err := b.GetAcmeState().getConfigWithForcedUpdate(sc)
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

	if allowRoleExtKeyUsageRaw, ok := d.GetOk("allow_role_ext_key_usage"); ok {
		config.AllowRoleExtKeyUsage = allowRoleExtKeyUsageRaw.(bool)
	}

	if defaultDirectoryPolicyRaw, ok := d.GetOk("default_directory_policy"); ok {
		config.DefaultDirectoryPolicy = defaultDirectoryPolicyRaw.(string)
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
			return nil, fmt.Errorf("invalid eab policy name provided, valid values are '%s', '%s', '%s'",
				eabPolicyNotRequired, eabPolicyNewAccountRequired, eabPolicyAlwaysRequired)
		}
		config.EabPolicyName = eabPolicy.Name
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		maxTTL := time.Second * time.Duration(maxTTLRaw.(int))
		if maxTTL <= 0 {
			return nil, fmt.Errorf("invalid max_ttl value, must be greater than 0")
		}
		config.MaxTTL = maxTTL
	}

	// Validate Default Directory Behavior:
	defaultDirectoryPolicyType, extraInfo, err := getDefaultDirectoryPolicyType(config.DefaultDirectoryPolicy)
	if err != nil {
		return nil, fmt.Errorf("invalid default_directory_policy: %w", err)
	}
	defaultDirectoryRoleName := ""
	switch defaultDirectoryPolicyType {
	case Forbid:
	case SignVerbatim:
	case ExternalPolicy:
		if !constants.IsEnterprise {
			return nil, fmt.Errorf("external-policy is only available in enterprise versions of Vault")
		}
	case Role:
		defaultDirectoryRoleName = extraInfo

		_, err := getAndValidateAcmeRole(sc, defaultDirectoryRoleName)
		if err != nil {
			return nil, fmt.Errorf("default directory policy role %v is not a valid ACME role: %w", defaultDirectoryRoleName, err)
		}
	default:
		return nil, fmt.Errorf("validation for the type of policy defined by %v is undefined", config.DefaultDirectoryPolicy)
	}

	// Validate Allowed Roles
	allowAnyRole := len(config.AllowedRoles) == 1 && config.AllowedRoles[0] == "*"
	foundDefault := false
	if !allowAnyRole {
		for index, name := range config.AllowedRoles {
			if name == "*" {
				return nil, fmt.Errorf("cannot use '*' as role name at index %d", index)
			}

			_, err := getAndValidateAcmeRole(sc, name)
			if err != nil {
				return nil, fmt.Errorf("allowed_role %v is not a valid acme role: %w", name, err)
			}

			if defaultDirectoryPolicyType == Role && name == defaultDirectoryRoleName {
				foundDefault = true
			}
		}

		if !foundDefault && defaultDirectoryPolicyType == Role {
			return nil, fmt.Errorf("default directory policy %v was not specified in allowed_roles: %v", config.DefaultDirectoryPolicy, config.AllowedRoles)
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

	// Check to make sure that we have a proper value for the cluster path which ACME requires
	if config.Enabled {
		_, err = getBasePathFromClusterConfig(sc)
		if err != nil {
			return nil, err
		}
	}

	var warnings []string
	// Lastly lets verify that the configuration is honored/invalidated by the public ACME env var.
	isPublicAcmeDisabledByEnv, err := isPublicACMEDisabledByEnv()
	if err != nil {
		warnings = append(warnings, err.Error())
	}
	if isPublicAcmeDisabledByEnv && config.Enabled {
		eabPolicy := getEabPolicyByName(config.EabPolicyName)
		if !eabPolicy.OverrideEnvDisablingPublicAcme() {
			resp := logical.ErrorResponse("%s env var is enabled, ACME EAB policy needs to be '%s' with ACME enabled",
				disableAcmeEnvVar, eabPolicyAlwaysRequired)
			resp.Warnings = warnings
			return resp, nil
		}
	}

	if _, err := b.GetAcmeState().writeConfig(sc, config); err != nil {
		return nil, fmt.Errorf("failed persisting: %w", err)
	}

	return genResponseFromAcmeConfig(config, warnings), nil
}

func isPublicACMEDisabledByEnv() (bool, error) {
	disableAcmeRaw, ok := os.LookupEnv(disableAcmeEnvVar)
	if !ok {
		return false, nil
	}

	disableAcme, err := strconv.ParseBool(disableAcmeRaw)
	if err != nil {
		// So the environment variable was set but we couldn't parse the value as a string, assume
		// the operator wanted public ACME disabled.
		return true, fmt.Errorf("failed parsing environment variable %s: %w", disableAcmeEnvVar, err)
	}

	return disableAcme, nil
}

func getDefaultDirectoryPolicyType(defaultDirectoryPolicy string) (DefaultDirectoryPolicyType, string, error) {
	switch {
	case defaultDirectoryPolicy == "forbid":
		return Forbid, "", nil
	case defaultDirectoryPolicy == "sign-verbatim":
		return SignVerbatim, "", nil
	case strings.HasPrefix(defaultDirectoryPolicy, rolePrefix):
		if len(defaultDirectoryPolicy) == rolePrefixLength {
			return Forbid, "", fmt.Errorf("no role specified by policy %v", defaultDirectoryPolicy)
		}
		roleName := defaultDirectoryPolicy[rolePrefixLength:]
		return Role, roleName, nil
	case strings.HasPrefix(defaultDirectoryPolicy, extPolicyPrefix):
		if len(defaultDirectoryPolicy) == extPolicyPrefixLength {
			// default external-policy case without a specified policy
			return ExternalPolicy, "", nil
		}

		if strings.HasPrefix(defaultDirectoryPolicy, extPolicyPrefix+":") &&
			len(defaultDirectoryPolicy) == extPolicyPrefixLength+1 {
			// end user set 'external-policy:', so no policy which is acceptable
			return ExternalPolicy, "", nil
		}

		policyName := defaultDirectoryPolicy[extPolicyPrefixLength+1:]
		if ok := extPolicyRegex.MatchString(policyName); !ok {
			return Forbid, "", fmt.Errorf("invalid characters within external-policy name: %s", defaultDirectoryPolicy)
		}
		return ExternalPolicy, policyName, nil
	default:
		return Forbid, "", fmt.Errorf("string %v not a valid Default Directory Policy", defaultDirectoryPolicy)
	}
}

//go:generate enumer -type=DefaultDirectoryPolicyType
type DefaultDirectoryPolicyType int

const (
	Forbid DefaultDirectoryPolicyType = iota
	SignVerbatim
	Role
	ExternalPolicy
)
