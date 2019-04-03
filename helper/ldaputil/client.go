package ldaputil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"net/url"
	"strings"
	"text/template"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/tlsutil"
)

type Client struct {
	Logger hclog.Logger
	LDAP   LDAP
}

func (c *Client) DialLDAP(cfg *ConfigEntry) (Connection, error) {
	var retErr *multierror.Error
	var conn Connection
	urls := strings.Split(cfg.Url, ",")
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
			conn, err = c.LDAP.Dial("tcp", net.JoinHostPort(host, port))
			if err != nil {
				break
			}
			if conn == nil {
				err = fmt.Errorf("empty connection after dialing")
				break
			}
			if cfg.StartTLS {
				tlsConfig, err = getTLSConfig(cfg, host)
				if err != nil {
					break
				}
				err = conn.StartTLS(tlsConfig)
			}
		case "ldaps":
			if port == "" {
				port = "636"
			}
			tlsConfig, err = getTLSConfig(cfg, host)
			if err != nil {
				break
			}
			conn, err = c.LDAP.DialTLS("tcp", net.JoinHostPort(host, port), tlsConfig)
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
 * Discover and return the bind string for the user attempting to authenticate.
 * This is handled in one of several ways:
 *
 * 1. If DiscoverDN is set, the user object will be searched for using userdn (base search path)
 *    and userattr (the attribute that maps to the provided username).
 *    The bind will either be anonymous or use binddn and bindpassword if they were provided.
 * 2. If upndomain is set, the user dn is constructed as 'username@upndomain'. See https://msdn.microsoft.com/en-us/library/cc223499.aspx
 *
 */
func (c *Client) GetUserBindDN(cfg *ConfigEntry, conn Connection, username string) (string, error) {
	bindDN := ""
	// Note: The logic below drives the logic in ConfigEntry.Validate().
	// If updated, please update there as well.
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
			Scope:     ldap.ScopeWholeSubtree,
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
			bindDN = fmt.Sprintf("%s@%s", EscapeLDAPValue(username), cfg.UPNDomain)
		} else {
			bindDN = fmt.Sprintf("%s=%s,%s", cfg.UserAttr, EscapeLDAPValue(username), cfg.UserDN)
		}
	}

	return bindDN, nil
}

/*
 * Returns the DN of the object representing the authenticated user.
 */
func (c *Client) GetUserDN(cfg *ConfigEntry, conn Connection, bindDN string) (string, error) {
	userDN := ""
	if cfg.UPNDomain != "" {
		// Find the distinguished name for the user if userPrincipalName used for login
		filter := fmt.Sprintf("(userPrincipalName=%s)", ldap.EscapeFilter(bindDN))
		if c.Logger.IsDebug() {
			c.Logger.Debug("searching upn", "userdn", cfg.UserDN, "filter", filter)
		}
		result, err := conn.Search(&ldap.SearchRequest{
			BaseDN:    cfg.UserDN,
			Scope:     ldap.ScopeWholeSubtree,
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

func (c *Client) performLdapFilterGroupsSearch(cfg *ConfigEntry, conn Connection, userDN string, username string) ([]*ldap.Entry, error) {
	if cfg.GroupFilter == "" {
		c.Logger.Warn("groupfilter is empty, will not query server")
		return make([]*ldap.Entry, 0), nil
	}

	if cfg.GroupDN == "" {
		c.Logger.Warn("groupdn is empty, will not query server")
		return make([]*ldap.Entry, 0), nil
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
		Scope:  ldap.ScopeWholeSubtree,
		Filter: renderedQuery.String(),
		Attributes: []string{
			cfg.GroupAttr,
		},
		SizeLimit: math.MaxInt32,
	})
	if err != nil {
		return nil, errwrap.Wrapf("LDAP search failed: {{err}}", err)
	}

	return result.Entries, nil
}

func sidBytesToString(b []byte) (string, error) {
	reader := bytes.NewReader(b)

	var revision, subAuthorityCount uint8
	var identifierAuthorityParts [3]uint16

	if err := binary.Read(reader, binary.LittleEndian, &revision); err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("SID %#v convert failed reading Revision: {{err}}", b), err)
	}

	if err := binary.Read(reader, binary.LittleEndian, &subAuthorityCount); err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("SID %#v convert failed reading SubAuthorityCount: {{err}}", b), err)
	}

	if err := binary.Read(reader, binary.BigEndian, &identifierAuthorityParts); err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("SID %#v convert failed reading IdentifierAuthority: {{err}}", b), err)
	}
	identifierAuthority := (uint64(identifierAuthorityParts[0]) << 32) + (uint64(identifierAuthorityParts[1]) << 16) + uint64(identifierAuthorityParts[2])

	subAuthority := make([]uint32, subAuthorityCount)
	if err := binary.Read(reader, binary.LittleEndian, &subAuthority); err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("SID %#v convert failed reading SubAuthority: {{err}}", b), err)
	}

	result := fmt.Sprintf("S-%d-%d", revision, identifierAuthority)
	for _, subAuthorityPart := range subAuthority {
		result += fmt.Sprintf("-%d", subAuthorityPart)
	}

	return result, nil
}

func (c *Client) performLdapTokenGroupsSearch(cfg *ConfigEntry, conn Connection, userDN string) ([]*ldap.Entry, error) {
	result, err := conn.Search(&ldap.SearchRequest{
		BaseDN: userDN,
		Scope:  ldap.ScopeBaseObject,
		Filter: "(objectClass=*)",
		Attributes: []string{
			"tokenGroups",
		},
		SizeLimit: 1,
	})
	if err != nil {
		return nil, errwrap.Wrapf("LDAP search failed: {{err}}", err)
	}
	if len(result.Entries) == 0 {
		c.Logger.Warn("unable to read object for group attributes", "userdn", userDN, "groupattr", cfg.GroupAttr)
		return make([]*ldap.Entry, 0), nil
	}

	userEntry := result.Entries[0]
	groupAttrValues := userEntry.GetRawAttributeValues("tokenGroups")

	groupEntries := make([]*ldap.Entry, 0, len(groupAttrValues))
	for _, sidBytes := range groupAttrValues {
		sidString, err := sidBytesToString(sidBytes)
		if err != nil {
			c.Logger.Warn("unable to read sid", "err", err)
			continue
		}

		groupResult, err := conn.Search(&ldap.SearchRequest{
			BaseDN: fmt.Sprintf("<SID=%s>", sidString),
			Scope:  ldap.ScopeBaseObject,
			Filter: "(objectClass=*)",
			Attributes: []string{
				"1.1", // RFC no attributes
			},
			SizeLimit: 1,
		})
		if err != nil {
			c.Logger.Warn("unable to read the group sid", "sid", sidString)
			continue
		}
		if len(groupResult.Entries) == 0 {
			c.Logger.Warn("unable to find the group", "sid", sidString)
			continue
		}

		groupEntries = append(groupEntries, groupResult.Entries[0])
	}

	return groupEntries, nil
}

/*
 * getLdapGroups queries LDAP and returns a slice describing the set of groups the authenticated user is a member of.
 *
 * If cfg.UseTokenGroups is true then the search is performed directly on the userDN.
 * The values of those attributes are converted to string SIDs, and then looked up to get ldap.Entry objects.
 * Otherwise, the search query is constructed according to cfg.GroupFilter, and run in context of cfg.GroupDN.
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
func (c *Client) GetLdapGroups(cfg *ConfigEntry, conn Connection, userDN string, username string) ([]string, error) {
	var entries []*ldap.Entry
	var err error
	if cfg.UseTokenGroups {
		entries, err = c.performLdapTokenGroupsSearch(cfg, conn, userDN)
	} else {
		entries, err = c.performLdapFilterGroupsSearch(cfg, conn, userDN, username)
	}
	if err != nil {
		return nil, err
	}

	// retrieve the groups in a string/bool map as a structure to avoid duplicates inside
	ldapMap := make(map[string]bool)

	for _, e := range entries {
		dn, err := ldap.ParseDN(e.DN)
		if err != nil || len(dn.RDNs) == 0 {
			continue
		}

		// Enumerate attributes of each result, parse out CN and add as group
		values := e.GetAttributeValues(cfg.GroupAttr)
		if len(values) > 0 {
			for _, val := range values {
				groupCN := getCN(val)
				ldapMap[groupCN] = true
			}
		} else {
			// If groupattr didn't resolve, use self (enumerating group objects)
			groupCN := getCN(e.DN)
			ldapMap[groupCN] = true
		}
	}

	ldapGroups := make([]string, 0, len(ldapMap))
	for key := range ldapMap {
		ldapGroups = append(ldapGroups, key)
	}

	return ldapGroups, nil
}

// EscapeLDAPValue is exported because a plugin uses it outside this package.
func EscapeLDAPValue(input string) string {
	if input == "" {
		return ""
	}

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

/*
 * Parses a distinguished name and returns the CN portion.
 * Given a non-conforming string (such as an already-extracted CN),
 * it will be returned as-is.
 */
func getCN(dn string) string {
	parsedDN, err := ldap.ParseDN(dn)
	if err != nil || len(parsedDN.RDNs) == 0 {
		// It was already a CN, return as-is
		return dn
	}

	for _, rdn := range parsedDN.RDNs {
		for _, rdnAttr := range rdn.Attributes {
			if strings.EqualFold(rdnAttr.Type, "CN") {
				return rdnAttr.Value
			}
		}
	}

	// Default, return self
	return dn
}

func getTLSConfig(cfg *ConfigEntry, host string) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		ServerName: host,
	}

	if cfg.TLSMinVersion != "" {
		tlsMinVersion, ok := tlsutil.TLSLookup[cfg.TLSMinVersion]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_min_version' in config")
		}
		tlsConfig.MinVersion = tlsMinVersion
	}

	if cfg.TLSMaxVersion != "" {
		tlsMaxVersion, ok := tlsutil.TLSLookup[cfg.TLSMaxVersion]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_max_version' in config")
		}
		tlsConfig.MaxVersion = tlsMaxVersion
	}

	if cfg.InsecureTLS {
		tlsConfig.InsecureSkipVerify = true
	}
	if cfg.Certificate != "" {
		caPool := x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM([]byte(cfg.Certificate))
		if !ok {
			return nil, fmt.Errorf("could not append CA certificate")
		}
		tlsConfig.RootCAs = caPool
	}
	return tlsConfig, nil
}
