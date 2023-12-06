// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldaputil

import (
	"github.com/go-ldap/ldap/v3"
)

func NewLDAP() LDAP {
	return &ldapIfc{}
}

// LDAP provides ldap functionality, but through an interface
// rather than statically. This allows faking it for tests.
type LDAP interface {
	DialURL(addr string, opts ...ldap.DialOpt) (Connection, error)
}

type ldapIfc struct{}

func (l *ldapIfc) DialURL(addr string, opts ...ldap.DialOpt) (Connection, error) {
	return ldap.DialURL(addr, opts...)
}
