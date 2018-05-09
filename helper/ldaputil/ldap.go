package ldaputil

import (
	"crypto/tls"

	"github.com/go-ldap/ldap"
)

func NewLDAP() LDAP {
	return &ldapIfc{}
}

// LDAP provides ldap functionality, but through an interface
// rather than statically. This allows faking it for tests.
type LDAP interface {
	Dial(network, addr string) (Connection, error)
	DialTLS(network, addr string, config *tls.Config) (Connection, error)
}

type ldapIfc struct{}

func (l *ldapIfc) Dial(network, addr string) (Connection, error) {
	return ldap.Dial(network, addr)
}

func (l *ldapIfc) DialTLS(network, addr string, config *tls.Config) (Connection, error) {
	return ldap.DialTLS(network, addr, config)
}
