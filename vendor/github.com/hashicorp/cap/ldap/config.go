// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
)

var derefAliasMap = map[string]int{
	"never":     ldap.NeverDerefAliases,
	"finding":   ldap.DerefFindingBaseObj,
	"searching": ldap.DerefInSearching,
	"always":    ldap.DerefAlways,
}

func validateDerefAlias(deref string) (string, error) {
	const op = "ldap.validateDerefAlias"
	lowerDeref := strings.ToLower(deref)
	_, found := derefAliasMap[lowerDeref]
	switch {
	case found:
		return lowerDeref, nil
	case deref == "":
		return DefaultDerefAliases, nil
	default:
		return "", fmt.Errorf("%s: invalid dereference_aliases %q: %w", op, deref, ErrInvalidParameter)
	}
}

const (
	// DefaultTimeout is the timeout value used for both dialing and requests to
	// the LDAP server
	DefaultTimeout = 60

	// DefaultURL for the ClientConfig.URLs
	DefaultURL = "ldaps://127.0.0.1:686"

	// DefaultUserAttr is the "username" attribute of the entry's DN and is
	// typically either the cn in ActiveDirectory or uid in openLDAP  (default:
	// cn)
	DefaultUserAttr = "cn"

	// DefaultGroupFilter for the ClientConfig.GroupFilter
	DefaultGroupFilter = `(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`

	// DefaultGroupAttr for the ClientConfig.GroupAttr
	DefaultGroupAttr = "cn"

	// DefaultTLSMinVersion for the ClientConfig.TLSMinVersion
	DefaultTLSMinVersion = "tls12"

	// DefaultTLSMaxVersion for the ClientConfig.TLSMaxVersion
	DefaultTLSMaxVersion = "tls13"

	// DefaultOpenLDAPUserPasswordAttribute defines the attribute name for the
	// openLDAP default password attribute which will always be excluded
	DefaultOpenLDAPUserPasswordAttribute = "userPassword"

	// DefaultADUserPasswordAttribute defines the attribute name for the
	// AD default password attribute which will always be excluded
	DefaultADUserPasswordAttribute = "unicodePwd"

	// DefaultDerefAliases defines the default for dereferencing aliases
	DefaultDerefAliases = "never"
)

type ClientConfig struct {
	// URLs are the URLs to use when connecting to a directory (default:
	// ldap://127.0.0.1).  When multiple URLs are specified; they are tried
	// in the order specified.
	URLs []string `json:"urls"`

	// UserDN is the base distinguished name to use when searching for users
	// (eg: ou=People,dc=example,dc=org)
	UserDN string `json:"userdn"`

	// AnonymousGroupSearch specifies that an anonymous bind should be used when
	// searching for groups (if true, the bind credentials will still be used
	// for the initial connection test).
	AnonymousGroupSearch bool `json:"anonymous_group_search"`

	// AllowEmptyAnonymousGroupSearches: if true it removes the userDN from
	// unauthenticated group searches (optional).
	AllowEmptyAnonymousGroupSearch bool `json:"allow_empty_anonymous_group_search"`

	// GroupDN is the distinguished name to use as base when searching for group
	// membership (eg: ou=Groups,dc=example,dc=org)
	GroupDN string `json:"groupdn"`

	// GroupFilter is a Go template for querying the group membership of user
	// (optional).  The template can access the following context variables:
	// UserDN, Username
	//
	// Example:
	// (&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))
	// Default: (|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`
	GroupFilter string `json:"groupfilter"`

	// GroupAttr is the attribute which identifies group members in entries
	// returned from GroupFilter queries.  Examples: for groupattr queries
	// returning group objects, use: cn. For queries returning user objects,
	// use: memberOf.
	// Default: cn
	GroupAttr string `json:"groupattr"`

	// UPNDomain is the userPrincipalName domain, which enables a
	// userPrincipalDomain login with [username]@UPNDomain (optional)
	UPNDomain string `json:"upndomain"`

	// UserFilter (optional) is a Go template used to construct a ldap user
	// search filter. The template can access the following context variables:
	// [UserAttr, Username]. The default userfilter is
	// ({{.UserAttr}}={{.Username}}) or
	// (userPrincipalName={{.Username}}@UPNDomain) if the upndomain parameter
	// is set. The user search filter can be used to  restrict what user can
	// attempt to log in. For example, to limit login to users that are not
	// contractors, you could write
	// (&(objectClass=user)({{.UserAttr}}={{.Username}})(!(employeeType=Contractor)))
	UserFilter string `json:"userfilter"`

	// UserAttr is the "username" attribute of the entry's DN and is typically
	// either the cn in ActiveDirectory or uid in openLDAP  (default: cn)
	UserAttr string `json:"userattr"`

	// Certificates to use verify the identity of the directory service and is a
	// set of PEM encoded x509 (optional)
	Certificates []string `json:"certificates"`

	// ClientTLSCert is the client certificate used with the ClientTLSKey to
	// authenticate the client to the directory service.  It must be PEM encoded
	// x509 (optional)
	ClientTLSCert string `json:"client_tls_cert"`

	// ClientTLSKey is the client certificate key used with the ClientTLSCert to
	// authenticate the client to the directory service.  It must be a PEM
	// encoded x509 (optional)
	ClientTLSKey string `json:"client_tls_key"`

	// InsecureTLS will skip the verification of the directory service's
	// certificate when making a client connection (optional).
	// Warning: this is insecure
	InsecureTLS bool `json:"insecure_tls"`

	// StartTLS will issue the StartTLS command after establishing an initial
	// non-TLS connection (optional)
	StartTLS bool `json:"starttls"`

	// BindDN is the distinguished name used when the client binds
	// (authenticates) to a directory service
	BindDN string `json:"binddn"`

	// BindPassword is the password used with the BindDN when the client binds
	// (authenticates) to a directory service (optional)
	BindPassword string `json:"bindpass"`

	// AllowEmptyPasswordBinds: if true it allows binds even if the user's
	// password is empty (zero length) (optional).
	AllowEmptyPasswordBinds bool `json:"allow_empty_passwd_bind"`

	// DiscoverDN: if true, it will use an anonymous bind with a search
	// to discover the bind DN of a user (optional)
	DiscoverDN bool `json:"discoverdn"`

	// TLSMinVersion version to use. Accepted values are
	// 'tls10', 'tls11', 'tls12' or 'tls13'. Defaults to 'tls12'
	TLSMinVersion string `json:"tls_min_version"`

	// TLSMaxVersion version to use. Accepted values are 'tls10', 'tls11',
	// 'tls12' or 'tls13'. Defaults to 'tls12'
	TLSMaxVersion string `json:"tls_max_version"`

	// UseTokenGroups: if true, use the Active Directory tokenGroups constructed
	// attribute of the user to find the group memberships. This will find all
	// security groups including nested ones.",
	UseTokenGroups bool `json:"use_token_groups"`

	// RequestTimeout in seconds is used when dialing to establish the
	// connection and when making requests against the server via a connection
	// before returning back an error. If not set, then the DefaultTimeout is
	// used.
	RequestTimeout int `json:"request_timeout"`

	// IncludeUserAttributes optionally specifies that the authenticating user's
	// DN and attributes be included an authentication AuthResult.
	//
	// Note: the default password attribute for both openLDAP (userPassword) and
	// AD (unicodePwd) will always be excluded.
	IncludeUserAttributes bool

	// ExcludedUserAttributes optionally defines a set of user attributes to be
	// excluded when an authenticating user's attributes are included in an
	// AuthResult (see: Config.IncludeUserAttributes or the WithUserAttributes()
	// option).
	//
	// Note: the default password attribute for both openLDAP (userPassword) and
	// AD (unicodePwd) will always be excluded.
	ExcludedUserAttributes []string

	// LowerUserAttributeKeys optionally specifies that the authenticating user's
	// DN and attributes be included in AuthResult use lowercase key names rather
	// than the default camel case.
	LowerUserAttributeKeys bool

	// IncludeUserGroups optionally specifies that the authenticating user's
	// group membership be included an authentication AuthResult.
	IncludeUserGroups bool

	// MaximumPageSize optionally specifies a maximum ldap search result size to
	// use when retrieving the authenticated user's group memberships. This can
	// be used to avoid reaching the LDAP server's max result size.
	MaximumPageSize int `json:"max_page_size"`

	// DerefAliases will control how aliases are dereferenced when
	// performing the search. Possible values are: never, finding, searching,
	// and always. If unset, a default of "never" is used. When set to
	// "finding", it will only dereference aliases during name resolution of the
	// base. When set to "searching", it will dereference aliases after name
	// resolution.
	DerefAliases string `json:"dereference_aliases"`

	// DeprecatedVaultPre111GroupCNBehavior: if true, group searching reverts to
	// the pre 1.1.1 Vault behavior.
	// see: https://www.vaultproject.io/docs/upgrading/upgrade-to-1.1.1
	DeprecatedVaultPre111GroupCNBehavior *bool `json:"use_pre111_group_cn_behavior"`
}

func (c *ClientConfig) clone() (*ClientConfig, error) {
	clone := *c
	return &clone, nil
}

func (c *ClientConfig) validate() error {
	const op = "ldap.(ClientConfig).validate"
	if len(c.URLs) == 0 {
		return fmt.Errorf("%s: at least one url must be provided: %w", op, ErrInvalidParameter)
	}
	tlsMinVersion, ok := tlsutil.TLSLookup[c.TLSMinVersion]
	if !ok {
		return fmt.Errorf("%s: invalid 'tls_min_version' in config: %w", op, ErrInvalidParameter)
	}
	tlsMaxVersion, ok := tlsutil.TLSLookup[c.TLSMaxVersion]
	if !ok {
		return fmt.Errorf("%s: invalid 'tls_max_version' in config: %w", op, ErrInvalidParameter)
	}
	if tlsMaxVersion < tlsMinVersion {
		return fmt.Errorf("%s: 'tls_max_version' must be greater than or equal to 'tls_min_version': %w", op, ErrInvalidParameter)
	}
	if c.Certificates != nil {
		for _, cert := range c.Certificates {
			if err := validateCertificate([]byte(cert)); err != nil {
				return fmt.Errorf("%s: failed to parse server tls cert: %w", op, err)
			}
		}
	}
	if (c.ClientTLSCert != "" && c.ClientTLSKey == "") ||
		(c.ClientTLSCert == "" && c.ClientTLSKey != "") {
		return fmt.Errorf("%s: both client_tls_cert and client_tls_key must be set in configuration: %w", op, ErrInvalidParameter)
	}
	if c.ClientTLSCert != "" && c.ClientTLSKey != "" {
		if _, err := tls.X509KeyPair([]byte(c.ClientTLSCert), []byte(c.ClientTLSKey)); err != nil {
			return fmt.Errorf("%s: failed to parse client X509 key pair: %w", op, err)
		}
	}
	var err error
	c.DerefAliases, err = validateDerefAlias(c.DerefAliases)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func validateCertificate(pemBlock []byte) error {
	const op = "ldap.validateCertificate"
	if pemBlock == nil {
		return fmt.Errorf("%s: missing certificate pem block: %w", op, ErrInvalidParameter)
	}
	block, _ := pem.Decode([]byte(pemBlock))
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("%s: failed to decode PEM block in the certificate: %w", op, ErrInvalidParameter)
	}
	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("%s: failed to parse certificate %w", op, err)
	}
	return nil
}
