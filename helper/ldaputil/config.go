package ldaputil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"strings"
	"text/template"

	"bytes"
	"math"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/tlsutil"
)

type ConfigEntry struct {
	Logger        log.Logger
	Url           string `json:"url"`
	UserDN        string `json:"userdn"`
	GroupDN       string `json:"groupdn"`
	GroupFilter   string `json:"groupfilter"`
	GroupAttr     string `json:"groupattr"`
	UPNDomain     string `json:"upndomain"`
	UserAttr      string `json:"userattr"`
	Certificate   string `json:"certificate"`
	InsecureTLS   bool   `json:"insecure_tls"`
	StartTLS      bool   `json:"starttls"`
	BindDN        string `json:"binddn"`
	BindPassword  string `json:"bindpass"`
	DenyNullBind  bool   `json:"deny_null_bind"`
	DiscoverDN    bool   `json:"discoverdn"`
	TLSMinVersion string `json:"tls_min_version"`
	TLSMaxVersion string `json:"tls_max_version"`

	// This json tag deviates from snake case because there was a past issue
	// where the tag was being ignored, causing it to be jsonified as "CaseSensitiveNames".
	// To continue reading in users' previously stored values,
	// we chose to carry that forward.
	CaseSensitiveNames *bool `json:"CaseSensitiveNames,omitempty"`
}

func (c *ConfigEntry) GetTLSConfig(host string) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		ServerName: host,
	}

	if c.TLSMinVersion != "" {
		tlsMinVersion, ok := tlsutil.TLSLookup[c.TLSMinVersion]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_min_version' in config")
		}
		tlsConfig.MinVersion = tlsMinVersion
	}

	if c.TLSMaxVersion != "" {
		tlsMaxVersion, ok := tlsutil.TLSLookup[c.TLSMaxVersion]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_max_version' in config")
		}
		tlsConfig.MaxVersion = tlsMaxVersion
	}

	if c.InsecureTLS {
		tlsConfig.InsecureSkipVerify = true
	}
	if c.Certificate != "" {
		caPool := x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM([]byte(c.Certificate))
		if !ok {
			return nil, fmt.Errorf("could not append CA certificate")
		}
		tlsConfig.RootCAs = caPool
	}
	return tlsConfig, nil
}

func (c *ConfigEntry) DialLDAP() (*ldap.Conn, error) {
	var retErr *multierror.Error
	var conn *ldap.Conn
	urls := strings.Split(c.Url, ",")
	for _, uut := range urls {
		u, err := url.Parse(uut)
		if err != nil {
			retErr = multierror.Append(retErr, errwrap.Wrapf(fmt.Sprintf("error parsing url %q: {{err}}", uut), err))
			continue
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			host = u.Host
		}

		var tlsConfig *tls.Config
		switch u.Scheme {
		case "ldap":
			if port == "" {
				port = "389"
			}
			conn, err = ldap.Dial("tcp", net.JoinHostPort(host, port))
			if err != nil {
				break
			}
			if conn == nil {
				err = fmt.Errorf("empty connection after dialing")
				break
			}
			if c.StartTLS {
				tlsConfig, err = c.GetTLSConfig(host)
				if err != nil {
					break
				}
				err = conn.StartTLS(tlsConfig)
			}
		case "ldaps":
			if port == "" {
				port = "636"
			}
			tlsConfig, err = c.GetTLSConfig(host)
			if err != nil {
				break
			}
			conn, err = ldap.DialTLS("tcp", net.JoinHostPort(host, port), tlsConfig)
		default:
			retErr = multierror.Append(retErr, fmt.Errorf("invalid LDAP scheme in url %q", net.JoinHostPort(host, port)))
			continue
		}
		if err == nil {
			if retErr != nil {
				if c.Logger.IsDebug() {
					c.Logger.Debug("errors connecting to some hosts: %s", retErr.Error())
				}
			}
			retErr = nil
			break
		}
		retErr = multierror.Append(retErr, errwrap.Wrapf(fmt.Sprintf("error connecting to host %q: {{err}}", uut), err))
	}

	return conn, retErr.ErrorOrNil()
}

/*
 * GetLdapGroups queries LDAP and returns a slice describing the set of groups the authenticated user is a member of.
 *
 * The search query is constructed according to cfg.GroupFilter, and run in context of cfg.GroupDN.
 * Groups will be resolved from the query results by following the attribute defined in cfg.GroupAttr.
 *
 * cfg.GroupFilter is a go template and is compiled with the following context: [UserDN, Username]
 *    UserDN - The DN of the authenticated user
 *    Username - The Username of the authenticated user
 *
 * Example:
 *   cfg.GroupFilter = "(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))"
 *   cfg.GroupDN     = "OU=Groups,DC=myorg,DC=com"
 *   cfg.GroupAttr   = "cn"
 *
 * NOTE - If cfg.GroupFilter is empty, no query is performed and an empty result slice is returned.
 *
 */
func (c *ConfigEntry) GetLdapGroups(cfg *ConfigEntry, conn *ldap.Conn, userDN string, username string) ([]string, error) {
	// retrieve the groups in a string/bool map as a structure to avoid duplicates inside
	ldapMap := make(map[string]bool)

	if cfg.GroupFilter == "" {
		c.Logger.Warn("groupfilter is empty, will not query server")
		return make([]string, 0), nil
	}

	if cfg.GroupDN == "" {
		c.Logger.Warn("groupdn is empty, will not query server")
		return make([]string, 0), nil
	}

	// If groupfilter was defined, resolve it as a Go template and use the query for
	// returning the user's groups
	if c.Logger.IsDebug() {
		c.Logger.Debug("compiling group filter", "group_filter", cfg.GroupFilter)
	}

	// Parse the configuration as a template.
	// Example template "(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))"
	t, err := template.New("queryTemplate").Parse(cfg.GroupFilter)
	if err != nil {
		return nil, errwrap.Wrapf("LDAP search failed due to template compilation error: {{err}}", err)
	}

	// Build context to pass to template - we will be exposing UserDn and Username.
	context := struct {
		UserDN   string
		Username string
	}{
		ldap.EscapeFilter(userDN),
		ldap.EscapeFilter(username),
	}

	var renderedQuery bytes.Buffer
	t.Execute(&renderedQuery, context)

	if c.Logger.IsDebug() {
		c.Logger.Debug("searching", "groupdn", cfg.GroupDN, "rendered_query", renderedQuery.String())
	}

	result, err := conn.Search(&ldap.SearchRequest{
		BaseDN: cfg.GroupDN,
		Scope:  2, // subtree
		Filter: renderedQuery.String(),
		Attributes: []string{
			cfg.GroupAttr,
		},
		SizeLimit: math.MaxInt32,
	})
	if err != nil {
		return nil, errwrap.Wrapf("LDAP search failed: {{err}}", err)
	}

	for _, e := range result.Entries {
		dn, err := ldap.ParseDN(e.DN)
		if err != nil || len(dn.RDNs) == 0 {
			continue
		}

		// Enumerate attributes of each result, parse out CN and add as group
		values := e.GetAttributeValues(cfg.GroupAttr)
		if len(values) > 0 {
			for _, val := range values {
				groupCN := c.getCN(val)
				ldapMap[groupCN] = true
			}
		} else {
			// If groupattr didn't resolve, use self (enumerating group objects)
			groupCN := c.getCN(e.DN)
			ldapMap[groupCN] = true
		}
	}

	ldapGroups := make([]string, 0, len(ldapMap))
	for key, _ := range ldapMap {
		ldapGroups = append(ldapGroups, key)
	}

	return ldapGroups, nil
}

/*
 * Parses a distinguished name and returns the CN portion.
 * Given a non-conforming string (such as an already-extracted CN),
 * it will be returned as-is.
 */
func (c *ConfigEntry) getCN(dn string) string {
	parsedDN, err := ldap.ParseDN(dn)
	if err != nil || len(parsedDN.RDNs) == 0 {
		// It was already a CN, return as-is
		return dn
	}

	for _, rdn := range parsedDN.RDNs {
		for _, rdnAttr := range rdn.Attributes {
			if rdnAttr.Type == "CN" {
				return rdnAttr.Value
			}
		}
	}

	// Default, return self
	return dn
}

/*
 * Discover and return the bind string for the user attempting to authenticate.
 * This is handled in one of several ways:
 *
 * 1. If DiscoverDN is set, the user object will be searched for using userdn (base search path)
 *    and userattr (the attribute that maps to the provided username).
 *    The bind will either be anonymous or use binddn and bindpassword if they were provided.
 * 2. If upndomain is set, the user dn is constructed as 'username@upndomain'. See https://msdn.microsoft.com/en-us/library/cc223499.aspx
 *
 */
func (c *ConfigEntry) GetUserBindDN(cfg *ConfigEntry, conn *ldap.Conn, username string) (string, error) {
	bindDN := ""
	if cfg.DiscoverDN || (cfg.BindDN != "" && cfg.BindPassword != "") {
		var err error
		if cfg.BindPassword != "" {
			err = conn.Bind(cfg.BindDN, cfg.BindPassword)
		} else {
			err = conn.UnauthenticatedBind(cfg.BindDN)
		}
		if err != nil {
			return bindDN, errwrap.Wrapf("LDAP bind (service) failed: {{err}}", err)
		}

		filter := fmt.Sprintf("(%s=%s)", cfg.UserAttr, ldap.EscapeFilter(username))
		if c.Logger.IsDebug() {
			c.Logger.Debug("discovering user", "userdn", cfg.UserDN, "filter", filter)
		}
		result, err := conn.Search(&ldap.SearchRequest{
			BaseDN:    cfg.UserDN,
			Scope:     2, // subtree
			Filter:    filter,
			SizeLimit: math.MaxInt32,
		})
		if err != nil {
			return bindDN, errwrap.Wrapf("LDAP search for binddn failed: {{err}}", err)
		}
		if len(result.Entries) != 1 {
			return bindDN, fmt.Errorf("LDAP search for binddn 0 or not unique")
		}
		bindDN = result.Entries[0].DN
	} else {
		if cfg.UPNDomain != "" {
			bindDN = fmt.Sprintf("%s@%s", escapeLDAPValue(username), cfg.UPNDomain)
		} else {
			bindDN = fmt.Sprintf("%s=%s,%s", cfg.UserAttr, escapeLDAPValue(username), cfg.UserDN)
		}
	}

	return bindDN, nil
}

/*
 * Returns the DN of the object representing the authenticated user.
 */
func (c *ConfigEntry) GetUserDN(cfg *ConfigEntry, conn *ldap.Conn, bindDN string) (string, error) {
	userDN := ""
	if cfg.UPNDomain != "" {
		// Find the distinguished name for the user if userPrincipalName used for login
		filter := fmt.Sprintf("(userPrincipalName=%s)", ldap.EscapeFilter(bindDN))
		if c.Logger.IsDebug() {
			c.Logger.Debug("searching upn", "userdn", cfg.UserDN, "filter", filter)
		}
		result, err := conn.Search(&ldap.SearchRequest{
			BaseDN:    cfg.UserDN,
			Scope:     2, // subtree
			Filter:    filter,
			SizeLimit: math.MaxInt32,
		})
		if err != nil {
			return userDN, errwrap.Wrapf("LDAP search failed for detecting user: {{err}}", err)
		}
		for _, e := range result.Entries {
			userDN = e.DN
		}
	} else {
		userDN = bindDN
	}

	return userDN, nil
}

func escapeLDAPValue(input string) string {
	// RFC4514 forbids un-escaped:
	// - leading space or hash
	// - trailing space
	// - special characters '"', '+', ',', ';', '<', '>', '\\'
	// - null
	for i := 0; i < len(input); i++ {
		escaped := false
		if input[i] == '\\' {
			i++
			escaped = true
		}
		switch input[i] {
		case '"', '+', ',', ';', '<', '>', '\\':
			if !escaped {
				input = input[0:i] + "\\" + input[i:]
				i++
			}
			continue
		}
		if escaped {
			input = input[0:i] + "\\" + input[i:]
			i++
		}
	}
	if input[0] == ' ' || input[0] == '#' {
		input = "\\" + input
	}
	if input[len(input)-1] == ' ' {
		input = input[0:len(input)-1] + "\\ "
	}
	return input
}
