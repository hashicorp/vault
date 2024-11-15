// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
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

func (c *Client) Search(cfg *ADConf, baseDN string, filters map[*Field][]string) ([]*Entry, error) {
	req := &ldap.SearchRequest{
		BaseDN:    baseDN,
		Scope:     ldap.ScopeWholeSubtree,
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

func (c *Client) UpdateEntry(cfg *ADConf, baseDN string, filters map[*Field][]string, newValues map[*Field][]string) error {
	entries, err := c.Search(cfg, baseDN, filters)
	if err != nil {
		return err
	}
	if len(entries) != 1 {
		return fmt.Errorf("filter of %s doesn't match just one entry: %+v", filters, entries)
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

// UpdatePassword uses a Modify call under the hood because
// Active Directory doesn't recognize the passwordModify method.
// See https://github.com/go-ldap/ldap/issues/106
// for more.
func (c *Client) UpdatePassword(cfg *ADConf, baseDN string, filters map[*Field][]string, newPassword string) error {
	pwdEncoded, err := formatPassword(newPassword)
	if err != nil {
		return err
	}

	newValues := map[*Field][]string{
		FieldRegistry.UnicodePassword: {pwdEncoded},
	}

	return c.UpdateEntry(cfg, baseDN, filters, newValues)
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

func bind(cfg *ADConf, conn ldaputil.Connection) error {
	if cfg.BindPassword == "" {
		return errors.New("unable to bind due to lack of configured password")
	}

	if cfg.UPNDomain != "" {
		origErr := conn.Bind(fmt.Sprintf("%s@%s", ldaputil.EscapeLDAPValue(cfg.BindDN), cfg.UPNDomain), cfg.BindPassword)
		if origErr == nil {
			return nil
		}
		if !shouldTryLastPwd(cfg.LastBindPassword, cfg.LastBindPasswordRotation) {
			return origErr
		}
		if err := conn.Bind(fmt.Sprintf("%s@%s", ldaputil.EscapeLDAPValue(cfg.BindDN), cfg.UPNDomain), cfg.LastBindPassword); err != nil {
			// Return the original error because it'll be more helpful for debugging.
			return origErr
		}
		return nil
	}

	if cfg.BindDN != "" {
		origErr := conn.Bind(cfg.BindDN, cfg.BindPassword)
		if origErr == nil {
			return nil
		}
		if !shouldTryLastPwd(cfg.LastBindPassword, cfg.LastBindPasswordRotation) {
			return origErr
		}
		if err := conn.Bind(cfg.BindDN, cfg.LastBindPassword); err != nil {
			// Return the original error because it'll be more helpful for debugging.
			return origErr
		}
	}
	return errors.New("must provide binddn or upndomain")
}

// shouldTryLastPwd determines if we should try a previous password.
// Active Directory can return a variety of errors when a password is invalid.
// Rather than attempting to catalogue these errors across multiple versions of
// AD, we simply try the last password if it's been less than a set amount of
// time since a rotation occurred.
func shouldTryLastPwd(lastPwd string, lastBindPasswordRotation time.Time) bool {
	if lastPwd == "" {
		return false
	}
	if lastBindPasswordRotation.Equal(time.Time{}) {
		return false
	}
	return lastBindPasswordRotation.Add(10 * time.Minute).After(time.Now())
}
