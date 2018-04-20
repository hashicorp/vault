package activedirectory

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/go-errors/errors"
	"github.com/go-ldap/ldap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/ldapifc"
	"golang.org/x/text/encoding/unicode"
)

func NewClient(logger hclog.Logger, conf *Configuration) *Client {
	return &Client{logger, conf, ldapifc.NewClient()}
}

func NewClientWith(logger hclog.Logger, conf *Configuration, ldapClient ldapifc.Client) *Client {
	return &Client{logger, conf, ldapClient}
}

type Client struct {
	logger     hclog.Logger
	conf       *Configuration
	ldapClient ldapifc.Client
}

func (c *Client) Search(filters map[*Field][]string) ([]*Entry, error) {

	req := &ldap.SearchRequest{
		BaseDN: toDNString(strings.Split(c.conf.RootDomainName, ",")),
		Scope:  ldap.ScopeWholeSubtree,
		Filter: toFilterString(filters),
	}

	conn, err := c.getFirstSucceedingConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	result, err := conn.Search(req)
	if err != nil {
		return nil, err
	}

	entries := make([]*Entry, len(result.Entries))
	for i, rawEntry := range result.Entries {
		entries[i] = NewEntry(c.logger, rawEntry)
	}
	return entries, nil
}

func (c *Client) UpdateEntry(filters map[*Field][]string, newValues map[*Field][]string) error {

	entries, err := c.Search(filters)
	if err != nil {
		return err
	}
	if len(entries) != 1 {
		return fmt.Errorf("filter of %s doesn't match just one entry: %s", filters, entries)
	}

	replaceAttributes := make([]ldap.PartialAttribute, len(newValues))
	i := 0
	for field, vals := range newValues {
		replaceAttributes[i] = ldap.PartialAttribute{
			Type: field.String(),
			Vals: vals,
		}
		i++
	}

	modifyReq := &ldap.ModifyRequest{
		DN:                entries[0].DN,
		ReplaceAttributes: replaceAttributes,
	}

	conn, err := c.getFirstSucceedingConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Modify(modifyReq)
}

// UpdatePassword uses a Modify call under the hood because
// Active Directory doesn't recognize the passwordModify method.
// See https://github.com/go-ldap/ldap/issues/106
// for more.
func (c *Client) UpdatePassword(filters map[*Field][]string, newPassword string) error {

	if !c.conf.StartTLS {
		return errors.New("per Active Directory, a TLS session must be in progress to update passwords, please update your StartTLS setting")
	}

	pwdEncoded, err := formatPassword(newPassword)
	if err != nil {
		return err
	}

	newValues := map[*Field][]string{
		FieldRegistry.UnicodePassword: {pwdEncoded},
	}

	return c.UpdateEntry(filters, newValues)
}

func (c *Client) getFirstSucceedingConnection() (ldapifc.Connection, error) {

	var retErr *multierror.Error

	tlsConfigs, err := c.conf.GetTLSConfigs(c.logger)
	if err != nil {
		return nil, err
	}

	for u, tlsConfig := range tlsConfigs {
		conn, err := c.connect(u, tlsConfig)
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("error parsing url %v: %v", u, err.Error()))
			continue
		}

		if c.conf.Username != "" && c.conf.Password != "" {
			if err := conn.Bind(c.conf.Username, c.conf.Password); err != nil {
				retErr = multierror.Append(retErr, fmt.Errorf("error binding to url %s: %s", u, err.Error()))
				continue
			}
		}

		return conn, nil
	}

	if c.logger.IsDebug() {
		c.logger.Debug("ldap: errors connecting to some hosts: %s", retErr.Error())
	}

	return nil, retErr
}

func (c *Client) connect(u *url.URL, tlsConfig *tls.Config) (ldapifc.Connection, error) {

	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		// err intentionally ignored - we'll fall back to default ldap ports if we're unable to parse this
		port = ""
	}

	switch u.Scheme {

	case "ldap":

		if port == "" {
			port = "389"
		}

		conn, err := c.ldapClient.Dial("tcp", net.JoinHostPort(tlsConfig.ServerName, port))
		if err != nil {
			return nil, err
		}

		if c.conf.StartTLS {
			if err = conn.StartTLS(tlsConfig); err != nil {
				return nil, err
			}
		}
		return conn, nil

	case "ldaps":

		if port == "" {
			port = "636"
		}

		conn, err := c.ldapClient.DialTLS("tcp", net.JoinHostPort(tlsConfig.ServerName, port), tlsConfig)
		if err != nil {
			return nil, err
		}
		return conn, nil

	default:
		return nil, fmt.Errorf("invalid LDAP scheme in url %q", net.JoinHostPort(tlsConfig.ServerName, port))
	}
}

// According to the MS docs, the password needs to be utf16 and enclosed in quotes.
func formatPassword(original string) (string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return utf16.NewEncoder().String("\"" + original + "\"")
}

func toDNString(dnValues []string) string {
	m := map[*Field][]string{
		FieldRegistry.DomainComponent: dnValues,
	}
	return toJoinedFieldEqualsString(m)
}

// Ex. "dc=example,dc=com"
func toJoinedFieldEqualsString(fieldValues map[*Field][]string) string {
	var fieldEquals []string
	for f, values := range fieldValues {
		for _, v := range values {
			fieldEquals = append(fieldEquals, fmt.Sprintf("%s=%s", f, v))
		}
	}
	return strings.Join(fieldEquals, ",")
}

// Ex. "(cn=Ellen Jones)"
func toFilterString(filters map[*Field][]string) string {
	result := toJoinedFieldEqualsString(filters)
	return "(" + result + ")"
}
