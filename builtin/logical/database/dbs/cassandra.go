package dbs

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/tlsutil"
)

type Cassandra struct {
	// Session is goroutine safe, however, since we reinitialize
	// it when connection info changes, we want to make sure we
	// can close it and use a new connection; hence the lock
	session *gocql.Session
	config  ConnectionConfig

	sync.RWMutex
}

func (c *Cassandra) Type() string {
	return cassandraTypeName
}

func (c *Cassandra) Connection() (*gocql.Session, error) {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	// If we already have a DB, we got it!
	if c.session != nil {
		return c.session, nil
	}

	session, err := createSession(c.config)
	if err != nil {
		return nil, err
	}

	//  Store the session in backend for reuse
	c.session = session

	return session, nil
}

func (p *Cassandra) Close() {
	// Grab the write lock
	p.Lock()
	defer p.Unlock()

	if p.session != nil {
		p.session.Close()
	}

	p.session = nil
}

func (p *Cassandra) Reset(config ConnectionConfig) (*sql.DB, error) {
	// Grab the write lock
	p.Lock()
	p.config = config
	p.Unlock()

	p.Close()
	return p.Connection()
}

func (p *Cassandra) CreateUser(createStmt, username, password, expiration string) error {
	// Get the connection
	db, err := p.Connection()
	if err != nil {
		return err
	}

	// TODO: This is racey
	// Grab a read lock
	p.RLock()
	defer p.RUnlock()

	return nil
}

func (p *Cassandra) RenewUser(username, expiration string) error {
	db, err := p.Connection()
	if err != nil {
		return err
	}
	// TODO: This is Racey
	// Grab the read lock
	p.RLock()
	defer p.RUnlock()

	return nil
}

func (p *Cassandra) CustomRevokeUser(username, revocationSQL string) error {
	db, err := p.Connection()
	if err != nil {
		return err
	}
	// TODO: this is Racey
	p.RLock()
	defer p.RUnlock()

	return nil
}

func (p *Cassandra) DefaultRevokeUser(username string) error {
	// Grab the read lock
	p.RLock()
	defer p.RUnlock()

	db, err := p.Connection()

	return nil
}

func createSession(cfg *ConnectionConfig) (*gocql.Session, error) {
	clusterConfig := gocql.NewCluster(strings.Split(cfg.Hosts, ",")...)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	clusterConfig.ProtoVersion = cfg.ProtocolVersion
	if clusterConfig.ProtoVersion == 0 {
		clusterConfig.ProtoVersion = 2
	}

	clusterConfig.Timeout = time.Duration(cfg.ConnectTimeout) * time.Second

	if cfg.TLS {
		var tlsConfig *tls.Config
		if len(cfg.Certificate) > 0 || len(cfg.IssuingCA) > 0 {
			if len(cfg.Certificate) > 0 && len(cfg.PrivateKey) == 0 {
				return nil, fmt.Errorf("Found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(cfg.Certificate) > 0 {
				certBundle.Certificate = cfg.Certificate
				certBundle.PrivateKey = cfg.PrivateKey
			}
			if len(cfg.IssuingCA) > 0 {
				certBundle.IssuingCA = cfg.IssuingCA
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate bundle: %s", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil || tlsConfig == nil {
				return nil, fmt.Errorf("failed to get TLS configuration: tlsConfig:%#v err:%v", tlsConfig, err)
			}
			tlsConfig.InsecureSkipVerify = cfg.InsecureTLS

			if cfg.TLSMinVersion != "" {
				var ok bool
				tlsConfig.MinVersion, ok = tlsutil.TLSLookup[cfg.TLSMinVersion]
				if !ok {
					return nil, fmt.Errorf("invalid 'tls_min_version' in config")
				}
			} else {
				// MinVersion was not being set earlier. Reset it to
				// zero to gracefully handle upgrades.
				tlsConfig.MinVersion = 0
			}
		}

		clusterConfig.SslOpts = &gocql.SslOptions{
			Config: *tlsConfig,
		}
	}

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("Error creating session: %s", err)
	}

	// Verify the info
	err = session.Query(`LIST USERS`).Exec()
	if err != nil {
		return nil, fmt.Errorf("Error validating connection info: %s", err)
	}

	return session, nil
}
