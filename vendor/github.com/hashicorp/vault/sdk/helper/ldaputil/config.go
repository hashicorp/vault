package ldaputil

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"

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
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "URL",
			},
		},

		"userdn": {
			Type:        framework.TypeString,
			Description: "LDAP domain to use for users (eg: ou=People,dc=example,dc=org)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "User DN",
			},
		},

		"binddn": {
			Type:        framework.TypeString,
			Description: "LDAP DN for searching for the user DN (optional)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Name of Object to bind (binddn)",
			},
		},

		"bindpass": {
			Type:        framework.TypeString,
			Description: "LDAP password for searching for the user DN (optional)",
			DisplayAttrs: &framework.DisplayAttributes{
				Sensitive: true,
			},
		},

		"groupdn": {
			Type:        framework.TypeString,
			Description: "LDAP search base to use for group membership search (eg: ou=Groups,dc=example,dc=org)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Group DN",
			},
		},

		"groupfilter": {
			Type:    framework.TypeString,
			Default: "(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))",
			Description: `Go template for querying group membership of user (optional)
The template can access the following context variables: UserDN, Username
Example: (&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))
Default: (|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`,
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Group Filter",
			},
		},

		"groupattr": {
			Type:    framework.TypeString,
			Default: "cn",
			Description: `LDAP attribute to follow on objects returned by <groupfilter>
in order to enumerate user group membership.
Examples: "cn" or "memberOf", etc.
Default: cn`,
			DisplayAttrs: &framework.DisplayAttributes{
				Name:  "Group Attribute",
				Value: "cn",
			},
		},

		"upndomain": {
			Type:        framework.TypeString,
			Description: "Enables userPrincipalDomain login with [username]@UPNDomain (optional)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "User Principal (UPN) Domain",
			},
		},

		"userattr": {
			Type:        framework.TypeString,
			Default:     "cn",
			Description: "Attribute used for users (default: cn)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name:  "User Attribute",
				Value: "cn",
			},
		},

		"certificate": {
			Type:        framework.TypeString,
			Description: "CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded (optional)",
		},

		"discoverdn": {
			Type:        framework.TypeBool,
			Description: "Use anonymous bind to discover the bind DN of a user (optional)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Discover DN",
			},
		},

		"insecure_tls": {
			Type:        framework.TypeBool,
			Description: "Skip LDAP server SSL Certificate verification - VERY insecure (optional)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Insecure TLS",
			},
		},

		"starttls": {
			Type:        framework.TypeBool,
			Description: "Issue a StartTLS command after establishing unencrypted connection (optional)",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Issue StartTLS",
			},
		},

		"tls_min_version": {
			Type:        framework.TypeString,
			Default:     "tls12",
			Description: "Minimum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Minimum TLS Version",
			},
			AllowedValues: []interface{}{"tls10", "tls11", "tls12"},
		},

		"tls_max_version": {
			Type:        framework.TypeString,
			Default:     "tls12",
			Description: "Maximum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Maximum TLS Version",
			},
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
func NewConfigEntry(existing *ConfigEntry, d *framework.FieldData) (*ConfigEntry, error) {
	var hadExisting bool
	var cfg *ConfigEntry

	if existing != nil {
		cfg = existing
		hadExisting = true
	} else {
		cfg = new(ConfigEntry)
	}

	if _, ok := d.Raw["url"]; ok || !hadExisting {
		cfg.Url = strings.ToLower(d.Get("url").(string))
	}

	if _, ok := d.Raw["userattr"]; ok || !hadExisting {
		cfg.UserAttr = strings.ToLower(d.Get("userattr").(string))
	}

	if _, ok := d.Raw["userdn"]; ok || !hadExisting {
		cfg.UserDN = d.Get("userdn").(string)
	}

	if _, ok := d.Raw["groupdn"]; ok || !hadExisting {
		cfg.GroupDN = d.Get("groupdn").(string)
	}

	if _, ok := d.Raw["groupfilter"]; ok || !hadExisting {
		groupfilter := d.Get("groupfilter").(string)
		if groupfilter != "" {
			// Validate the template before proceeding
			_, err := template.New("queryTemplate").Parse(groupfilter)
			if err != nil {
				return nil, errwrap.Wrapf("invalid groupfilter: {{err}}", err)
			}
		}

		cfg.GroupFilter = groupfilter
	}

	if _, ok := d.Raw["groupattr"]; ok || !hadExisting {
		cfg.GroupAttr = d.Get("groupattr").(string)
	}

	if _, ok := d.Raw["upndomain"]; ok || !hadExisting {
		cfg.UPNDomain = d.Get("upndomain").(string)
	}

	if _, ok := d.Raw["certificate"]; ok || !hadExisting {
		certificate := d.Get("certificate").(string)
		if certificate != "" {
			block, _ := pem.Decode([]byte(certificate))

			if block == nil || block.Type != "CERTIFICATE" {
				return nil, errors.New("failed to decode PEM block in the certificate")
			}
			_, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, errwrap.Wrapf("failed to parse certificate: {{err}}", err)
			}
		}

		cfg.Certificate = certificate
	}

	if _, ok := d.Raw["insecure_tls"]; ok || !hadExisting {
		cfg.InsecureTLS = d.Get("insecure_tls").(bool)
	}

	if _, ok := d.Raw["tls_min_version"]; ok || !hadExisting {
		cfg.TLSMinVersion = d.Get("tls_min_version").(string)
		_, ok = tlsutil.TLSLookup[cfg.TLSMinVersion]
		if !ok {
			return nil, errors.New("invalid 'tls_min_version'")
		}
	}

	if _, ok := d.Raw["tls_max_version"]; ok || !hadExisting {
		cfg.TLSMaxVersion = d.Get("tls_max_version").(string)
		_, ok = tlsutil.TLSLookup[cfg.TLSMaxVersion]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_max_version'")
		}
	}
	if cfg.TLSMaxVersion < cfg.TLSMinVersion {
		return nil, fmt.Errorf("'tls_max_version' must be greater than or equal to 'tls_min_version'")
	}

	if _, ok := d.Raw["starttls"]; ok || !hadExisting {
		cfg.StartTLS = d.Get("starttls").(bool)
	}

	if _, ok := d.Raw["binddn"]; ok || !hadExisting {
		cfg.BindDN = d.Get("binddn").(string)
	}

	if _, ok := d.Raw["bindpass"]; ok || !hadExisting {
		cfg.BindPassword = d.Get("bindpass").(string)
	}

	if _, ok := d.Raw["deny_null_bind"]; ok || !hadExisting {
		cfg.DenyNullBind = d.Get("deny_null_bind").(bool)
	}

	if _, ok := d.Raw["discoverdn"]; ok || !hadExisting {
		cfg.DiscoverDN = d.Get("discoverdn").(bool)
	}

	if _, ok := d.Raw["case_sensitive_names"]; ok || !hadExisting {
		cfg.CaseSensitiveNames = new(bool)
		*cfg.CaseSensitiveNames = d.Get("case_sensitive_names").(bool)
	}

	if _, ok := d.Raw["use_token_groups"]; ok || !hadExisting {
		cfg.UseTokenGroups = d.Get("use_token_groups").(bool)
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
