package cassandra

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Factory creates a new backend
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

// Backend contains the base information for the backend's functionality
func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		Paths: []*framework.Path{
			pathConfigConnection(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},
	}

	return b.Backend
}

type backend struct {
	*framework.Backend

	// Session is goroutine safe, however, since we reinitialize
	// it when connection info changes, we want to make sure we
	// can close it and use a new connection; hence the lock
	session *gocql.Session
	lock    sync.Mutex
}

type sessionConfig struct {
	Hosts           string `json:"hosts" structs:"hosts"`
	Username        string `json:"username" structs:"username"`
	Password        string `json:"password" structs:"password"`
	TLS             bool   `json:"tls" structs:"tls"`
	InsecureTLS     bool   `json:"insecure_tls" structs:"insecure_tls"`
	Certificate     string `json:"certificate" structs:"certificate"`
	PrivateKey      string `json:"private_key" structs:"private_key"`
	IssuingCA       string `json:"issuing_ca" structs:"issuing_ca"`
	ProtocolVersion int    `json:"protocol_version" structs:"protocol_version"`
}

// DB returns the database connection.
func (b *backend) DB(s logical.Storage) (*gocql.Session, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// If we already have a DB, we got it!
	if b.session != nil {
		return b.session, nil
	}

	entry, err := s.Get("config/connection")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil,
			fmt.Errorf("Configure the DB connection with config/connection first")
	}

	config := &sessionConfig{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}

	return createSession(config, s)
}

// ResetDB forces a connection next time DB() is called.
func (b *backend) ResetDB(newSession *gocql.Session) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.session != nil {
		b.session.Close()
	}

	b.session = newSession
}

const backendHelp = `
The Cassandra backend dynamically generates database users.

After mounting this backend, configure it using the endpoints within
the "config/" path.
`
