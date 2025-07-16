// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldaputil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"net/url"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/go-ldap/ldap/v3"
	hclog "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
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
			retErr = multierror.Append(retErr, fmt.Errorf(fmt.Sprintf("error parsing url %q: {{err}}", uut), err))
			continue
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			host = u.Host
		}

		var tlsConfig *tls.Config
		dialer := net.Dialer{
			Timeout: time.Duration(cfg.ConnectionTimeout) * time.Second,
		}

		switch u.Scheme {
		case "ldap":
			if port == "" {
				port = "389"
			}

			fullAddr := fmt.Sprintf("%s://%s", u.Scheme, net.JoinHostPort(host, port))
			opt := ldap.DialWithDialer(&dialer)

			conn, err = c.LDAP.DialURL(fullAddr, opt)
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

			fullAddr := fmt.Sprintf("%s://%s", u.Scheme, net.JoinHostPort(host, port))
			opt := ldap.DialWithDialer(&dialer)
			tls := ldap.DialWithTLSConfig(tlsConfig)

			conn, err = c.LDAP.DialURL(fullAddr, opt, tls)
			if err != nil {
				break
			}
		default:
			retErr = multierror.Append(retErr, fmt.Errorf("invalid LDAP scheme in url %q", net.JoinHostPort(host, port)))
			continue
		}
		if err == nil {
			if retErr != nil {
				if c.Logger.IsDebug() {
					c.Logger.Debug("errors connecting to some hosts", "error", retErr.Error())
				}
			}
			retErr = nil
			break
		}
		retErr = multierror.Append(retErr, fmt.Errorf(fmt.Sprintf("error connecting to host %q: {{err}}", uut), err))
	}
	if retErr != nil {
		return nil, retErr
	}
	if timeout := cfg.RequestTimeout; timeout > 0 {
		conn.SetTimeout(time.Duration(timeout) * time.Second)
	}
	return conn, nil
}

/*
 * Searches for a username in the ldap server, returning a minimal subset of the
 * user's attributes (if found)
 */
func (c *Client) makeLdapSearchRequest(cfg *ConfigEntry, conn Connection, username string) (*ldap.SearchResult, error) {
	// Note: The logic below drives the logic in ConfigEntry.Validate().
	// If updated, please update there as well.
	var err error
	if cfg.BindPassword != "" {
		err = conn.Bind(cfg.BindDN, cfg.BindPassword)
	} else {
		err = conn.UnauthenticatedBind(cfg.BindDN)
	}
	if err != nil {
		return nil, fmt.Errorf("LDAP bind (service) failed: %w", err)
	}

	renderedFilter, err := c.RenderUserSearchFilter(cfg, username)
	if err != nil {
		return nil, err
	}

	if c.Logger.IsDebug() {
		c.Logger.Debug("discovering user", "userdn", cfg.UserDN, "filter", renderedFilter)
	}
	ldapRequest := &ldap.SearchRequest{
		BaseDN:       cfg.UserDN,
		DerefAliases: ldapDerefAliasMap[cfg.DerefAliases],
		Scope:        ldap.ScopeWholeSubtree,
		Filter:       renderedFilter,
		SizeLimit:    2, // Should be only 1 result. Any number larger (2 or more) means access denied.
		Attributes: []string{
			cfg.UserAttr, // Return only needed attributes
		},
	}

	result, err := conn.Search(ldapRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/*
 * Discover and return the bind string for the user attempting to authenticate, as well as the
 * value to use for the identity alias.
 * This is handled in one of several ways:
 *
 * 1. If DiscoverDN is set, the user object will be searched for using userdn (base search path)
 *    and userattr (the attribute that maps to the provided username) or user search filter.
 *    The bind will either be anonymous or use binddn and bindpassword if they were provided.
 * 2. If upndomain is set, the user dn and alias attribte are constructed as 'username@upndomain'.
 *    See https://msdn.microsoft.com/en-us/library/cc223499.aspx
 *
 */
func (c *Client) GetUserBindDN(cfg *ConfigEntry, conn Connection, username string) (string, error) {
	bindDN := ""

	// Note: The logic below drives the logic in ConfigEntry.Validate().
	// If updated, please update there as well.
	if cfg.DiscoverDN || (cfg.BindDN != "" && cfg.BindPassword != "") {

		result, err := c.makeLdapSearchRequest(cfg, conn, username)
		if err != nil {
			return bindDN, fmt.Errorf("LDAP search for binddn failed %w", err)
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

func (c *Client) RenderUserSearchFilter(cfg *ConfigEntry, username string) (string, error) {
	// The UserFilter can be blank if not set, or running this version of the code
	// on an existing ldap configuration
	if cfg.UserFilter == "" {
		cfg.UserFilter = "({{.UserAttr}}={{.Username}})"
	}

	// If userfilter was defined, resolve it as a Go template and use the query to
	// find the login user
	if c.Logger.IsDebug() {
		c.Logger.Debug("compiling search filter", "search_filter", cfg.UserFilter)
	}

	// Parse the configuration as a template.
	// Example template "({{.UserAttr}}={{.Username}})"
	t, err := template.New("queryTemplate").Parse(cfg.UserFilter)
	if err != nil {
		return "", fmt.Errorf("LDAP search failed due to template compilation error: %w", err)
	}

	// Build context to pass to template - we will be exposing UserDn and Username.
	context := struct {
		UserAttr string
		Username string
	}{
		ldap.EscapeFilter(cfg.UserAttr),
		ldap.EscapeFilter(username),
	}
	if cfg.UPNDomain != "" {
		context.UserAttr = "userPrincipalName"
		// Intentionally, calling EscapeFilter(...) (vs EscapeValue) since the
		// username is being injected into a search filter.
		// As an untrusted string, the username must be escaped according to RFC
		// 4515, in order to prevent attackers from injecting characters that could modify the filter
		context.Username = fmt.Sprintf("%s@%s", ldap.EscapeFilter(username), cfg.UPNDomain)
	}

	// Execute the template. Note that the template context contains escaped input and does
	// not provide behavior via functions. Additionally, no function map has been provided
	// during template initialization. The only template functions available during execution
	// are the predefined global functions: https://pkg.go.dev/text/template#hdr-Functions
	var renderedFilter bytes.Buffer
	if err := t.Execute(&renderedFilter, context); err != nil {
		return "", fmt.Errorf("LDAP search failed due to template parsing error: %w", err)
	}

	return renderedFilter.String(), nil
}

/*
 * Returns the value to be used for the entity alias of this user
 * This is handled in one of several ways:
 *
 * 1. If DiscoverDN is set, the user will be searched for using userdn (base search path)
 *    and userattr (the attribute that maps to the provided username) or user search filter.
 *    The bind will either be anonymous or use binddn and bindpassword if they were provided.
 * 2. If upndomain is set, the alias attribte is constructed as 'username@upndomain'.
 *
 */
func (c *Client) GetUserAliasAttributeValue(cfg *ConfigEntry, conn Connection, username string) (string, error) {
	aliasAttributeValue := ""

	// Note: The logic below drives the logic in ConfigEntry.Validate().
	// If updated, please update there as well.
	if cfg.DiscoverDN || (cfg.BindDN != "" && cfg.BindPassword != "") {

		result, err := c.makeLdapSearchRequest(cfg, conn, username)
		if err != nil {
			return aliasAttributeValue, fmt.Errorf("LDAP search for entity alias attribute failed: %w", err)
		}
		if len(result.Entries) != 1 {
			return aliasAttributeValue, fmt.Errorf("LDAP search for entity alias attribute 0 or not unique")
		}

		if len(result.Entries[0].Attributes) != 1 {
			return aliasAttributeValue, fmt.Errorf("LDAP attribute missing for entity alias mapping")
		}

		if len(result.Entries[0].Attributes[0].Values) != 1 {
			return aliasAttributeValue, fmt.Errorf("LDAP entity alias attribute %s empty or not unique for entity alias mapping", cfg.UserAttr)
		}

		aliasAttributeValue = result.Entries[0].Attributes[0].Values[0]
	} else {
		if cfg.UPNDomain != "" {
			aliasAttributeValue = fmt.Sprintf("%s@%s", EscapeLDAPValue(username), cfg.UPNDomain)
		} else {
			aliasAttributeValue = fmt.Sprintf("%s=%s,%s", cfg.UserAttr, EscapeLDAPValue(username), cfg.UserDN)
		}
	}

	return aliasAttributeValue, nil
}

/*
 * Returns the DN of the object representing the authenticated user.
 */
func (c *Client) GetUserDN(cfg *ConfigEntry, conn Connection, bindDN, username string) (string, error) {
	userDN := ""
	if cfg.UPNDomain != "" {
		// Find the distinguished name for the user if userPrincipalName used for login
		filter := fmt.Sprintf("(userPrincipalName=%s@%s)", EscapeLDAPValue(username), cfg.UPNDomain)
		if c.Logger.IsDebug() {
			c.Logger.Debug("searching upn", "userdn", cfg.UserDN, "filter", filter)
		}
		result, err := conn.Search(&ldap.SearchRequest{
			BaseDN:       cfg.UserDN,
			Scope:        ldap.ScopeWholeSubtree,
			DerefAliases: ldapDerefAliasMap[cfg.DerefAliases],
			Filter:       filter,
			SizeLimit:    math.MaxInt32,
		})
		if err != nil {
			return userDN, fmt.Errorf("LDAP search failed for detecting user: %w", err)
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
		return nil, fmt.Errorf("LDAP search failed due to template compilation error: %w", err)
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
	if err := t.Execute(&renderedQuery, context); err != nil {
		return nil, fmt.Errorf("LDAP search failed due to template parsing error: %w", err)
	}

	if c.Logger.IsDebug() {
		c.Logger.Debug("searching", "groupdn", cfg.GroupDN, "rendered_query", renderedQuery.String())
	}

	result, err := conn.Search(&ldap.SearchRequest{
		BaseDN:       cfg.GroupDN,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldapDerefAliasMap[cfg.DerefAliases],
		Filter:       renderedQuery.String(),
		Attributes: []string{
			cfg.GroupAttr,
		},
		SizeLimit: math.MaxInt32,
	})
	if err != nil {
		return nil, fmt.Errorf("LDAP search failed: %w", err)
	}

	return result.Entries, nil
}

func (c *Client) performLdapFilterGroupsSearchPaging(cfg *ConfigEntry, conn PagingConnection, userDN string, username string) ([]*ldap.Entry, error) {
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
		return nil, fmt.Errorf("LDAP search failed due to template compilation error: %w", err)
	}

	// Build context to pass to template - we will be exposing UserDn and Username.
	context := struct {
		UserDN   string
		Username string
	}{
		ldap.EscapeFilter(userDN),
		ldap.EscapeFilter(username),
	}

	// Execute the template. Note that the template context contains escaped input and does
	// not provide behavior via functions. Additionally, no function map has been provided
	// during template initialization. The only template functions available during execution
	// are the predefined global functions: https://pkg.go.dev/text/template#hdr-Functions
	var renderedQuery bytes.Buffer
	if err := t.Execute(&renderedQuery, context); err != nil {
		return nil, fmt.Errorf("LDAP search failed due to template parsing error: %w", err)
	}

	if c.Logger.IsDebug() {
		c.Logger.Debug("searching", "groupdn", cfg.GroupDN, "rendered_query", renderedQuery.String())
	}

	result, err := conn.SearchWithPaging(&ldap.SearchRequest{
		BaseDN:       cfg.GroupDN,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldapDerefAliasMap[cfg.DerefAliases],
		Filter:       renderedQuery.String(),
		Attributes: []string{
			cfg.GroupAttr,
		},
		SizeLimit: math.MaxInt32,
	}, uint32(cfg.MaximumPageSize))
	if err != nil {
		return nil, fmt.Errorf("LDAP search failed: %w", err)
	}

	return result.Entries, nil
}

func sidBytesToString(b []byte) (string, error) {
	reader := bytes.NewReader(b)

	var revision, subAuthorityCount uint8
	var identifierAuthorityParts [3]uint16

	if err := binary.Read(reader, binary.LittleEndian, &revision); err != nil {
		return "", fmt.Errorf(fmt.Sprintf("SID %#v convert failed reading Revision: {{err}}", b), err)
	}

	if err := binary.Read(reader, binary.LittleEndian, &subAuthorityCount); err != nil {
		return "", fmt.Errorf(fmt.Sprintf("SID %#v convert failed reading SubAuthorityCount: {{err}}", b), err)
	}

	if err := binary.Read(reader, binary.BigEndian, &identifierAuthorityParts); err != nil {
		return "", fmt.Errorf(fmt.Sprintf("SID %#v convert failed reading IdentifierAuthority: {{err}}", b), err)
	}
	identifierAuthority := (uint64(identifierAuthorityParts[0]) << 32) + (uint64(identifierAuthorityParts[1]) << 16) + uint64(identifierAuthorityParts[2])

	subAuthority := make([]uint32, subAuthorityCount)
	if err := binary.Read(reader, binary.LittleEndian, &subAuthority); err != nil {
		return "", fmt.Errorf(fmt.Sprintf("SID %#v convert failed reading SubAuthority: {{err}}", b), err)
	}

	result := fmt.Sprintf("S-%d-%d", revision, identifierAuthority)
	for _, subAuthorityPart := range subAuthority {
		result += fmt.Sprintf("-%d", subAuthorityPart)
	}

	return result, nil
}

func (c *Client) performLdapTokenGroupsSearch(cfg *ConfigEntry, conn Connection, userDN string) ([]*ldap.Entry, error) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	taskChan := make(chan string)
	maxWorkers := 10

	result, err := conn.Search(&ldap.SearchRequest{
		BaseDN:       userDN,
		Scope:        ldap.ScopeBaseObject,
		DerefAliases: ldapDerefAliasMap[cfg.DerefAliases],
		Filter:       "(objectClass=*)",
		Attributes: []string{
			"tokenGroups",
		},
		SizeLimit: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("LDAP search failed: %w", err)
	}
	if len(result.Entries) == 0 {
		c.Logger.Warn("unable to read object for group attributes", "userdn", userDN, "groupattr", cfg.GroupAttr)
		return make([]*ldap.Entry, 0), nil
	}

	userEntry := result.Entries[0]
	groupAttrValues := userEntry.GetRawAttributeValues("tokenGroups")
	groupEntries := make([]*ldap.Entry, 0, len(groupAttrValues))

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for sid := range taskChan {
				groupResult, err := conn.Search(&ldap.SearchRequest{
					BaseDN:       fmt.Sprintf("<SID=%s>", sid),
					Scope:        ldap.ScopeBaseObject,
					DerefAliases: ldapDerefAliasMap[cfg.DerefAliases],
					Filter:       "(objectClass=*)",
					Attributes: []string{
						"1.1", // RFC no attributes
					},
					SizeLimit: 1,
				})
				if err != nil {
					c.Logger.Warn("unable to read the group sid", "sid", sid)
					continue
				}

				if len(groupResult.Entries) == 0 {
					c.Logger.Warn("unable to find the group", "sid", sid)
					continue
				}

				lock.Lock()
				groupEntries = append(groupEntries, groupResult.Entries[0])
				lock.Unlock()
			}
		}()
	}

	for _, sidBytes := range groupAttrValues {
		sidString, err := sidBytesToString(sidBytes)
		if err != nil {
			c.Logger.Warn("unable to read sid", "err", err)
			continue
		}
		taskChan <- sidString
	}

	close(taskChan)
	wg.Wait()

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
		if paging, ok := conn.(PagingConnection); ok && cfg.MaximumPageSize > 0 {
			entries, err = c.performLdapFilterGroupsSearchPaging(cfg, paging, userDN, username)
		} else {
			entries, err = c.performLdapFilterGroupsSearch(cfg, conn, userDN, username)
		}
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
				groupCN := getCN(cfg, val)
				ldapMap[groupCN] = true
			}
		} else {
			// If groupattr didn't resolve, use self (enumerating group objects)
			groupCN := getCN(cfg, e.DN)
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
// EscapeLDAPValue will properly escape the input string as an ldap value
// rfc4514 states the following must be escaped:
// - leading space or hash
// - trailing space
// - special characters '"', '+', ',', ';', '<', '>', '\\'
// - hex
func EscapeLDAPValue(input string) string {
	if input == "" {
		return ""
	}

	buf := bytes.Buffer{}

	escFn := func(c byte) {
		buf.WriteByte('\\')
		buf.WriteByte(c)
	}

	inputLen := len(input)
	for i := 0; i < inputLen; i++ {
		char := input[i]
		switch {
		case i == 0 && char == ' ' || char == '#':
			// leading space or hash.
			escFn(char)
			continue
		case i == inputLen-1 && char == ' ':
			// trailing space.
			escFn(char)
			continue
		case specialChar(char):
			escFn(char)
			continue
		case char < ' ' || char > '~':
			// anything that's not between the ascii space and tilde must be hex
			buf.WriteByte('\\')
			buf.WriteString(hex.EncodeToString([]byte{char}))
			continue
		default:
			// everything remaining, doesn't need to be escaped
			buf.WriteByte(char)
		}
	}
	return buf.String()
}

func specialChar(char byte) bool {
	switch char {
	case '"', '+', ',', ';', '<', '>', '\\':
		return true
	default:
		return false
	}
}

/*
 * Parses a distinguished name and returns the CN portion.
 * Given a non-conforming string (such as an already-extracted CN),
 * it will be returned as-is.
 */
func getCN(cfg *ConfigEntry, dn string) string {
	parsedDN, err := ldap.ParseDN(dn)
	if err != nil || len(parsedDN.RDNs) == 0 {
		// It was already a CN, return as-is
		return dn
	}

	for _, rdn := range parsedDN.RDNs {
		for _, rdnAttr := range rdn.Attributes {
			if cfg.UsePre111GroupCNBehavior == nil || *cfg.UsePre111GroupCNBehavior {
				if rdnAttr.Type == "CN" {
					return rdnAttr.Value
				}
			} else {
				if strings.EqualFold(rdnAttr.Type, "CN") {
					return rdnAttr.Value
				}
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
	if cfg.ClientTLSCert != "" && cfg.ClientTLSKey != "" {
		certificate, err := tls.X509KeyPair([]byte(cfg.ClientTLSCert), []byte(cfg.ClientTLSKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse client X509 key pair: %w", err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
	} else if cfg.ClientTLSCert != "" || cfg.ClientTLSKey != "" {
		return nil, fmt.Errorf("both client_tls_cert and client_tls_key must be set")
	}
	return tlsConfig, nil
}
