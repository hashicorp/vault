package ldapifc

import (
	"crypto/tls"
	"github.com/go-ldap/ldap"
)

func NewClient() Client {
	return &client{}
}

// Client provides ldap functionality, but through an interface
// rather than statically.
type Client interface {
	Dial(network, addr string) (Connection, error)
	DialTLS(network, addr string, config *tls.Config) (Connection, error)
}

type client struct{}

func (c *client) Dial(network, addr string) (Connection, error) {
	return ldap.Dial(network, addr)
}

func (c *client) DialTLS(network, addr string, config *tls.Config) (Connection, error) {
	return ldap.DialTLS(network, addr, config)
}
