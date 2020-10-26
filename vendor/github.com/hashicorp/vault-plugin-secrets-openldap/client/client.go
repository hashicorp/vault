package client

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
)

type Config struct {
	*ldaputil.ConfigEntry
	LastBindPassword         string    `json:"last_bind_password"`
	LastBindPasswordRotation time.Time `json:"last_bind_password_rotation"`
	Schema                   string    `json:"schema"`
}

func New() Client {
	return Client{
		ldap: &ldaputil.Client{
			LDAP: ldaputil.NewLDAP(),
		},
	}
}

type Client struct {
	ldap *ldaputil.Client
}

func (c *Client) Search(cfg *Config, baseDN string, filters map[*Field][]string) ([]*Entry, error) {
	req := &ldap.SearchRequest{
		BaseDN:    baseDN,
		Scope:     ldap.ScopeBaseObject,
		Filter:    toString(filters),
		SizeLimit: math.MaxInt32,
	}

	conn, err := c.ldap.DialLDAP(cfg.ConfigEntry)
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

func (c *Client) UpdateEntry(cfg *Config, baseDN string, filters map[*Field][]string, newValues map[*Field][]string) error {
	entries, err := c.Search(cfg, baseDN, filters)
	if err != nil {
		return err
	}
	if len(entries) != 1 {
		return fmt.Errorf("expected one matching entry, but received %d", len(entries))
	}

	modifyReq := &ldap.ModifyRequest{
		DN: entries[0].DN,
	}

	for field, vals := range newValues {
		modifyReq.Replace(field.String(), vals)
	}

	conn, err := c.ldap.DialLDAP(cfg.ConfigEntry)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := bind(cfg, conn); err != nil {
		return err
	}

	return conn.Modify(modifyReq)
}

// UpdatePassword uses a Modify call under the hood instead of LDAP change password function.
// This allows AD and OpenLDAP secret engines to use the same api without changes to
// the interface.
func (c *Client) UpdatePassword(cfg *Config, baseDN string, newValues map[*Field][]string, filters map[*Field][]string) error {
	return c.UpdateEntry(cfg, baseDN, filters, newValues)
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

func bind(cfg *Config, conn ldaputil.Connection) error {
	if cfg.BindPassword == "" {
		return errors.New("unable to bind due to lack of configured password")
	}

	if cfg.BindDN == "" {
		return errors.New("must provide binddn")
	}

	origErr := conn.Bind(cfg.BindDN, cfg.BindPassword)
	if origErr == nil {
		return nil
	}
	if !shouldTryLastPwd(cfg.LastBindPassword, cfg.LastBindPasswordRotation) {
		return origErr
	}

	return conn.Bind(cfg.BindDN, cfg.LastBindPassword)
}

// shouldTryLastPwd determines if we should try a previous password.
// LDAP can return a variety of errors when a password is invalid.
// Rather than attempting to catalogue these errors across multiple versions of
// OpenLDAP, we simply try the last password if it's been less than a set amount of
// time since a rotation occurred.
func shouldTryLastPwd(lastPwd string, lastBindPasswordRotation time.Time) bool {
	if lastPwd == "" {
		return false
	}
	if lastBindPasswordRotation.IsZero() {
		return false
	}
	return lastBindPasswordRotation.Add(10 * time.Minute).After(time.Now())
}
