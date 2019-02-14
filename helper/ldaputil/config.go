package ldaputil

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/logical/framework"

	"github.com/hashicorp/errwrap"
)

// ConfigFields returns all the config fields that can potentially be used by the LDAP client.
// Not all fields will be used by every integration.
func ConfigFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"url": {
			Type:        framework.TypeString,
			Default:     "ldap://127.0.0.1",
			Description: "LDAP URL to connect to (default: ldap://127.0.0.1). Multiple URLs can be specified by concatenating them with commas; they will be tried in-order.",
			DisplayName: "URL",
		},

		"userdn": {
			Type:        framework.TypeString,
			Description: "LDAP domain to use for users (eg: ou=People,dc=example,dc=org)",
			DisplayName: "User DN",
		},

		"binddn": {
			Type:        framework.TypeString,
			Description: "LDAP DN for searching for the user DN (optional)",
			DisplayName: "Name of Object to bind (binddn)",
		},

		"bindpass": {
			Type:             framework.TypeString,
			Description:      "LDAP password for searching for the user DN (optional)",
			DisplaySensitive: true,
		},

		"groupdn": {
			Type:        framework.TypeString,
			Description: "LDAP search base to use for group membership search (eg: ou=Groups,dc=example,dc=org)",
			DisplayName: "Group DN",
		},

		"groupfilter": {
			Type:    framework.TypeString,
			Default: "(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))",
			Description: `Go template for querying group membership of user (optional)
The template can access the following context variables: UserDN, Username
Example: (&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))
Default: (|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`,
		},

		"groupattr": {
			Type:    framework.TypeString,
			Default: "cn",
			Description: `LDAP attribute to follow on objects returned by <groupfilter>
in order to enumerate user group membership.
Examples: "cn" or "memberOf", etc.
Default: cn`,
			DisplayName: "Group Attribute",
		},

		"upndomain": {
			Type:        framework.TypeString,
			Description: "Enables userPrincipalDomain login with [username]@UPNDomain (optional)",
			DisplayName: "User Principal (UPN) Domain",
		},

		"userattr": {
			Type:        framework.TypeString,
			Default:     "cn",
			Description: "Attribute used for users (default: cn)",
			DisplayName: "User Attribute",
		},

		"certificate": {
			Type:        framework.TypeString,
			Description: "CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded (optional)",
		},

		"discoverdn": {
			Type:        framework.TypeBool,
			Description: "Use anonymous bind to discover the bind DN of a user (optional)",
			DisplayName: "Discover DN",
		},

		"insecure_tls": {
			Type:        framework.TypeBool,
			Description: "Skip LDAP server SSL Certificate verification - VERY insecure (optional)",
			DisplayName: "Insecure TLS",
		},

		"starttls": {
			Type:        framework.TypeBool,
			Description: "Issue a StartTLS command after establishing unencrypted connection (optional)",
			DisplayName: "Issue StartTLS command after establishing an unencrypted connection",
		},

		"tls_min_version": {
			Type:          framework.TypeString,
			Default:       "tls12",
			Description:   "Minimum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'",
			DisplayName:   "Minimum TLS Version",
			AllowedValues: []interface{}{"tls10", "tls11", "tls12"},
		},

		"tls_max_version": {
			Type:          framework.TypeString,
			Default:       "tls12",
			Description:   "Maximum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'",
			DisplayName:   "Maxumum TLS Version",
			AllowedValues: []interface{}{"tls10", "tls11", "tls12"},
		},

		"deny_null_bind": {
			Type:        framework.TypeBool,
			Default:     true,
			Description: "Denies an unauthenticated LDAP bind request if the user's password is empty; defaults to true",
		},

		"case_sensitive_names": {
			Type:        framework.TypeBool,
			Description: "If true, case sensitivity will be used when comparing usernames and groups for matching policies.",
		},

		"use_token_groups": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: "If true, use the Active Directory tokenGroups constructed attribute of the user to find the group memberships. This will find all security groups including nested ones.",
		},
	}
}

/*
 * Creates and initializes a ConfigEntry object with its default values,
 * as specified by the passed schema.
 */
func NewConfigEntry(d *framework.FieldData) (*ConfigEntry, error) {
	cfg := new(ConfigEntry)

	url := d.Get("url").(string)
	if url != "" {
		cfg.Url = strings.ToLower(url)
	}
	userattr := d.Get("userattr").(string)
	if userattr != "" {
		cfg.UserAttr = strings.ToLower(userattr)
	}
	userdn := d.Get("userdn").(string)
	if userdn != "" {
		cfg.UserDN = userdn
	}
	groupdn := d.Get("groupdn").(string)
	if groupdn != "" {
		cfg.GroupDN = groupdn
	}
	groupfilter := d.Get("groupfilter").(string)
	if groupfilter != "" {
		// Validate the template before proceeding
		_, err := template.New("queryTemplate").Parse(groupfilter)
		if err != nil {
			return nil, errwrap.Wrapf("invalid groupfilter: {{err}}", err)
		}

		cfg.GroupFilter = groupfilter
	}
	groupattr := d.Get("groupattr").(string)
	if groupattr != "" {
		cfg.GroupAttr = groupattr
	}
	upndomain := d.Get("upndomain").(string)
	if upndomain != "" {
		cfg.UPNDomain = upndomain
	}
	certificate := d.Get("certificate").(string)
	if certificate != "" {
		block, _ := pem.Decode([]byte(certificate))

		if block == nil || block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("failed to decode PEM block in the certificate")
		}
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errwrap.Wrapf("failed to parse certificate: {{err}}", err)
		}
		cfg.Certificate = certificate
	}
	insecureTLS := d.Get("insecure_tls").(bool)
	if insecureTLS {
		cfg.InsecureTLS = insecureTLS
	}
	cfg.TLSMinVersion = d.Get("tls_min_version").(string)
	if cfg.TLSMinVersion == "" {
		return nil, fmt.Errorf("failed to get 'tls_min_version' value")
	}

	var ok bool
	_, ok = tlsutil.TLSLookup[cfg.TLSMinVersion]
	if !ok {
		return nil, fmt.Errorf("invalid 'tls_min_version'")
	}

	cfg.TLSMaxVersion = d.Get("tls_max_version").(string)
	if cfg.TLSMaxVersion == "" {
		return nil, fmt.Errorf("failed to get 'tls_max_version' value")
	}

	_, ok = tlsutil.TLSLookup[cfg.TLSMaxVersion]
	if !ok {
		return nil, fmt.Errorf("invalid 'tls_max_version'")
	}
	if cfg.TLSMaxVersion < cfg.TLSMinVersion {
		return nil, fmt.Errorf("'tls_max_version' must be greater than or equal to 'tls_min_version'")
	}

	startTLS := d.Get("starttls").(bool)
	if startTLS {
		cfg.StartTLS = startTLS
	}

	bindDN := d.Get("binddn").(string)
	if bindDN != "" {
		cfg.BindDN = bindDN
	}

	bindPass := d.Get("bindpass").(string)
	if bindPass != "" {
		cfg.BindPassword = bindPass
	}

	denyNullBind := d.Get("deny_null_bind").(bool)
	if denyNullBind {
		cfg.DenyNullBind = denyNullBind
	}

	discoverDN := d.Get("discoverdn").(bool)
	if discoverDN {
		cfg.DiscoverDN = discoverDN
	}

	caseSensitiveNames, ok := d.GetOk("case_sensitive_names")
	if ok {
		cfg.CaseSensitiveNames = new(bool)
		*cfg.CaseSensitiveNames = caseSensitiveNames.(bool)
	}

	useTokenGroups := d.Get("use_token_groups").(bool)
	if useTokenGroups {
		cfg.UseTokenGroups = useTokenGroups
	}

	return cfg, nil
}

type ConfigEntry struct {
	Url            string `json:"url"`
	UserDN         string `json:"userdn"`
	GroupDN        string `json:"groupdn"`
	GroupFilter    string `json:"groupfilter"`
	GroupAttr      string `json:"groupattr"`
	UPNDomain      string `json:"upndomain"`
	UserAttr       string `json:"userattr"`
	Certificate    string `json:"certificate"`
	InsecureTLS    bool   `json:"insecure_tls"`
	StartTLS       bool   `json:"starttls"`
	BindDN         string `json:"binddn"`
	BindPassword   string `json:"bindpass"`
	DenyNullBind   bool   `json:"deny_null_bind"`
	DiscoverDN     bool   `json:"discoverdn"`
	TLSMinVersion  string `json:"tls_min_version"`
	TLSMaxVersion  string `json:"tls_max_version"`
	UseTokenGroups bool   `json:"use_token_groups"`

	// This json tag deviates from snake case because there was a past issue
	// where the tag was being ignored, causing it to be jsonified as "CaseSensitiveNames".
	// To continue reading in users' previously stored values,
	// we chose to carry that forward.
	CaseSensitiveNames *bool `json:"CaseSensitiveNames,omitempty"`
}

func (c *ConfigEntry) Map() map[string]interface{} {
	m := c.PasswordlessMap()
	m["bindpass"] = c.BindPassword
	return m
}

func (c *ConfigEntry) PasswordlessMap() map[string]interface{} {
	m := map[string]interface{}{
		"url":              c.Url,
		"userdn":           c.UserDN,
		"groupdn":          c.GroupDN,
		"groupfilter":      c.GroupFilter,
		"groupattr":        c.GroupAttr,
		"upndomain":        c.UPNDomain,
		"userattr":         c.UserAttr,
		"certificate":      c.Certificate,
		"insecure_tls":     c.InsecureTLS,
		"starttls":         c.StartTLS,
		"binddn":           c.BindDN,
		"deny_null_bind":   c.DenyNullBind,
		"discoverdn":       c.DiscoverDN,
		"tls_min_version":  c.TLSMinVersion,
		"tls_max_version":  c.TLSMaxVersion,
		"use_token_groups": c.UseTokenGroups,
	}
	if c.CaseSensitiveNames != nil {
		m["case_sensitive_names"] = *c.CaseSensitiveNames
	}
	return m
}

func (c *ConfigEntry) Validate() error {
	if len(c.Url) == 0 {
		return errors.New("at least one url must be provided")
	}
	// Note: This logic is driven by the logic in GetUserBindDN.
	// If updating this, please also update the logic there.
	if !c.DiscoverDN && (c.BindDN == "" || c.BindPassword == "") && c.UPNDomain == "" && c.UserDN == "" {
		return errors.New("cannot derive UserBindDN")
	}
	tlsMinVersion, ok := tlsutil.TLSLookup[c.TLSMinVersion]
	if !ok {
		return errors.New("invalid 'tls_min_version' in config")
	}
	tlsMaxVersion, ok := tlsutil.TLSLookup[c.TLSMaxVersion]
	if !ok {
		return errors.New("invalid 'tls_max_version' in config")
	}
	if tlsMaxVersion < tlsMinVersion {
		return errors.New("'tls_max_version' must be greater than or equal to 'tls_min_version'")
	}
	if c.Certificate != "" {
		block, _ := pem.Decode([]byte(c.Certificate))
		if block == nil || block.Type != "CERTIFICATE" {
			return errors.New("failed to decode PEM block in the certificate")
		}
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate %s", err.Error())
		}
	}
	return nil
}
