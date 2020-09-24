package ldaputil

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

func NewLDAP() LDAP {
	return &ldapIfc{}
}

// LDAP provides ldap functionality, but through an interface
// rather than statically. This allows faking it for tests.
type LDAP interface {
	Dial(network string, addr string) (Connection, error)
	DialTLS(network string, addr string, config *tls.Config) (Connection, error)
}

type ldapIfc struct{}

func (l *ldapIfc) Dial(_network string, addr string) (Connection, error) {
	dialer := ldap.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second})
	return ldap.DialURL(addr, dialer)
}

func (l *ldapIfc) DialTLS(_network string, addr string, config *tls.Config) (Connection, error) {
	tlsOpts := ldap.DialWithTLSConfig(config)
	dialer := ldap.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second})
	return ldap.DialURL(addr, dialer, tlsOpts)
}
