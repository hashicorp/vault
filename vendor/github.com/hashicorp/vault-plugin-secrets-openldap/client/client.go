// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/go-ldap/ldif"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
)

type Config struct {
	*ldaputil.ConfigEntry
	LastBindPassword         string    `json:"last_bind_password"`
	LastBindPasswordRotation time.Time `json:"last_bind_password_rotation"`
	Schema                   string    `json:"schema"`
}

func New(logger hclog.Logger) Client {
	if logger == nil {
		logger = hclog.NewNullLogger()
	}

	return Client{
		ldap: &ldaputil.Client{
			LDAP:   ldaputil.NewLDAP(),
			Logger: logger,
		},
	}
}

func NewWithClient(logger hclog.Logger, ldap ldaputil.LDAP) Client {
	return Client{
		ldap: &ldaputil.Client{
			Logger: logger,
			LDAP:   ldap,
		},
	}
}

type Client struct {
	ldap *ldaputil.Client
}

func (c *Client) Search(cfg *Config, baseDN string, scope int, filters map[*Field][]string) ([]*Entry, error) {
	req := &ldap.SearchRequest{
		BaseDN:    baseDN,
		Scope:     scope,
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
		return nil, fmt.Errorf("failed to search ldap server: %w", err)
	}

	entries := make([]*Entry, len(result.Entries))
	for i, rawEntry := range result.Entries {
		entries[i] = NewEntry(rawEntry)
	}
	return entries, nil
}

func (c *Client) UpdateEntry(cfg *Config, baseDN string, scope int, filters map[*Field][]string, newValues map[*Field][]string) error {
	entries, err := c.Search(cfg, baseDN, scope, filters)
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

// UpdatePassword uses a Modify call under the hood instead of LDAP change
// password function. This allows AD and OpenLDAP schemas to use the same
// api without changes to the interface.
func (c *Client) UpdatePassword(cfg *Config, baseDN string, scope int, newValues map[*Field][]string, filters map[*Field][]string) error {
	return c.UpdateEntry(cfg, baseDN, scope, filters, newValues)
}

// toString turns the following map of filters into LDAP search filter strings
// For example: "(cn=Ellen Jones)"
// when multiple filters are applied, they get AND'ed together.
// example: (&(x=1)(y=2))
// for test assertions, this sorts the filters alphabetically.
func toString(filters map[*Field][]string) string {
	var fieldEquals []string
	sortedFilters := make([]*Field, 0, len(filters))
	for filter := range filters {
		sortedFilters = append(sortedFilters, filter)
	}
	sort.Slice(sortedFilters, func(i, j int) bool {
		return sortedFilters[i].String() < sortedFilters[j].String()
	})
	for _, filter := range sortedFilters {
		values := filters[filter]
		// make values deterministic
		sort.Strings(values)
		for _, v := range values {
			fieldEquals = append(fieldEquals, fmt.Sprintf("(%s=%s)", filter, v))
		}
	}
	if len(fieldEquals) <= 1 {
		return strings.Join(fieldEquals, "")
	}

	return "(&" + strings.Join(fieldEquals, "") + ")"
}

func bind(cfg *Config, conn ldaputil.Connection) error {
	if cfg.BindPassword == "" {
		return errors.New("unable to bind due to lack of configured password")
	}

	// Determine the user to bind with
	var bindUser string
	switch {
	case cfg.UPNDomain != "":
		bindUser = fmt.Sprintf("%s@%s", ldaputil.EscapeLDAPValue(cfg.BindDN), cfg.UPNDomain)
	case cfg.BindDN != "":
		bindUser = cfg.BindDN
	default:
		return errors.New("must provide binddn or upndomain")
	}

	merr := new(multierror.Error)

	// Bind using the bind password. If this fails, attempt to bind with the prior
	// bind password for at most 10 minutes. We do this to allow continued operation
	// after a root credential rotation where we may not be able to bind with the new
	// password immediately.
	err := conn.Bind(bindUser, cfg.BindPassword)
	if err == nil {
		return nil
	}
	merr = multierror.Append(merr, err)

	if !shouldTryLastPwd(cfg.LastBindPassword, cfg.LastBindPasswordRotation) {
		return fmt.Errorf("failed to bind with current password: %w", merr.ErrorOrNil())
	}

	err = conn.Bind(bindUser, cfg.LastBindPassword)
	if err == nil {
		return nil
	}
	merr = multierror.Append(merr, err)

	return fmt.Errorf("failed to bind with current and prior password: %w", merr.ErrorOrNil())
}

// shouldTryLastPwd determines if we should try a previous password.
// LDAP can return a variety of errors when a password is invalid.
// Rather than attempting to catalogue these errors across multiple implementations of
// LDAP, we simply try the last password if it's been less than a set amount of
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

func (c *Client) Execute(cfg *Config, entries []*ldif.Entry, continueOnFailure bool) (err error) {
	if len(entries) == 0 {
		return nil
	}

	conn, err := c.ldap.DialLDAP(cfg.ConfigEntry)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := bind(cfg, conn); err != nil {
		return err
	}

	merr := new(multierror.Error)
	for _, entry := range entries {
		if entry == nil {
			// Skip entries that are nil since they don't indicate an error in execution. Since these entries
			// are usually coming from an ldif parse, this should generally not happen so it's mainly to
			// protect against developers from screwing up and creating a panic due to a nil reference
			continue
		}
		var err error
		switch {
		case entry.Entry != nil:
			addReq := coerceToAddRequest(entry.Entry)
			err = errorf("failed to run AddRequest: %w", conn.Add(addReq))
		case entry.Add != nil:
			err = errorf("failed to run AddRequest: %w", conn.Add(entry.Add))
		case entry.Modify != nil:
			err = errorf("failed to run ModifyRequest: %w", conn.Modify(entry.Modify))
		case entry.Del != nil:
			err = errorf("failed to run DelRequest: %w", conn.Del(entry.Del))
		default:
			err = fmt.Errorf("unrecognized or missing LDIF command")
		}

		if err != nil {
			if continueOnFailure {
				merr = multierror.Append(merr, err)
			} else {
				return err
			}
		}
	}
	return merr.ErrorOrNil()
}

func coerceToAddRequest(entry *ldap.Entry) *ldap.AddRequest {
	attributes := make([]ldap.Attribute, 0, len(entry.Attributes))
	for _, entryAttribute := range entry.Attributes {
		attribute := ldap.Attribute{
			Type: entryAttribute.Name,
			Vals: entryAttribute.Values,
		}
		attributes = append(attributes, attribute)
	}
	addReq := &ldap.AddRequest{
		DN:         entry.DN,
		Attributes: attributes,
		Controls:   nil,
	}
	return addReq
}

func errorf(format string, wrappedErr error) error {
	if wrappedErr == nil {
		return nil
	}

	return fmt.Errorf(format, wrappedErr)
}
