package client

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-errors/errors"
	"github.com/go-ldap/ldap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/ldaputil"
	"golang.org/x/text/encoding/unicode"
)

func NewClient(logger hclog.Logger) *Client {
	return &Client{
		ldap: &ldaputil.Client{
			Logger: logger,
			LDAP:   ldaputil.NewLDAP(),
		},
	}
}

type Client struct {
	ldap *ldaputil.Client
}

func (c *Client) Search(cfg *ldaputil.ConfigEntry, filters map[*Field][]string) ([]*Entry, error) {
	req := &ldap.SearchRequest{
		BaseDN:    cfg.UserDN,
		Scope:     ldap.ScopeWholeSubtree,
		Filter:    toString(filters),
		SizeLimit: math.MaxInt32,
	}

	conn, err := c.ldap.DialLDAP(cfg)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := bind(cfg, conn); err != nil {
		return nil, err
	}

	result, err := conn.Search(req)
	if err != nil {
		return nil, err
	}

	entries := make([]*Entry, len(result.Entries))
	for i, rawEntry := range result.Entries {
		entries[i] = NewEntry(rawEntry)
	}
	return entries, nil
}

func (c *Client) UpdateEntry(cfg *ldaputil.ConfigEntry, filters map[*Field][]string, newValues map[*Field][]string) error {
	entries, err := c.Search(cfg, filters)
	if err != nil {
		return err
	}
	if len(entries) != 1 {
		return fmt.Errorf("filter of %s doesn't match just one entry: %s", filters, entries)
	}

	modifyReq := &ldap.ModifyRequest{
		DN: entries[0].DN,
	}

	for field, vals := range newValues {
		modifyReq.Replace(field.String(), vals)
	}

	conn, err := c.ldap.DialLDAP(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := bind(cfg, conn); err != nil {
		return err
	}
	return conn.Modify(modifyReq)
}

// UpdatePassword uses a Modify call under the hood because
// Active Directory doesn't recognize the passwordModify method.
// See https://github.com/go-ldap/ldap/issues/106
// for more.
func (c *Client) UpdatePassword(cfg *ldaputil.ConfigEntry, filters map[*Field][]string, newPassword string) error {
	pwdEncoded, err := formatPassword(newPassword)
	if err != nil {
		return err
	}

	newValues := map[*Field][]string{
		FieldRegistry.UnicodePassword: {pwdEncoded},
	}

	return c.UpdateEntry(cfg, filters, newValues)
}

// According to the MS docs, the password needs to be utf16 and enclosed in quotes.
func formatPassword(original string) (string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return utf16.NewEncoder().String("\"" + original + "\"")
}

// Ex. "(cn=Ellen Jones)"
func toString(filters map[*Field][]string) string {
	var fieldEquals []string
	for f, values := range filters {
		for _, v := range values {
			fieldEquals = append(fieldEquals, fmt.Sprintf("%s=%s", f, v))
		}
	}
	result := strings.Join(fieldEquals, ",")
	return "(" + result + ")"
}

func bind(cfg *ldaputil.ConfigEntry, conn ldaputil.Connection) error {
	if cfg.BindPassword == "" {
		return errors.New("unable to bind due to lack of configured password")
	}
	if cfg.UPNDomain != "" {
		return conn.Bind(fmt.Sprintf("%s@%s", ldaputil.EscapeLDAPValue(cfg.BindDN), cfg.UPNDomain), cfg.BindPassword)
	}
	if cfg.BindDN != "" {
		return conn.Bind(cfg.BindDN, cfg.BindPassword)
	}
	return errors.New("must provide binddn or upndomain")
}
