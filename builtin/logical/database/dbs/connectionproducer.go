package dbs

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	// Import sql drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/tlsutil"
)

type ConnectionProducer interface {
	Connection() (interface{}, error)
	Close()
}

// sqlConnectionProducer impliments ConnectionProducer and provides a generic producer for most sql databases
type sqlConnectionProducer struct {
	ConnectionURL string `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`

	config *DatabaseConfig

	db *sql.DB
	sync.Mutex
}

func (c *sqlConnectionProducer) Connection() (interface{}, error) {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	// If we already have a DB, we got it!
	if c.db != nil {
		if err := c.db.Ping(); err == nil {
			return c.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		c.db.Close()
	}

	// Otherwise, attempt to make connection
	conn := c.ConnectionURL

	// Ensure timezone is set to UTC for all the conenctions
	if strings.HasPrefix(conn, "postgres://") || strings.HasPrefix(conn, "postgresql://") {
		if strings.Contains(conn, "?") {
			conn += "&timezone=utc"
		} else {
			conn += "?timezone=utc"
		}
	}

	var err error
	c.db, err = sql.Open(c.config.DatabaseType, conn)
	if err != nil {
		return nil, err
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	c.db.SetMaxOpenConns(c.config.MaxOpenConnections)
	c.db.SetMaxIdleConns(c.config.MaxIdleConnections)
	c.db.SetConnMaxLifetime(c.config.MaxConnectionLifetime)

	return c.db, nil
}

func (c *sqlConnectionProducer) Close() {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		c.db.Close()
	}

	c.db = nil
}

type cassandraConnectionProducer struct {
	Hosts           string `json:"hosts" structs:"hosts" mapstructure:"hosts"`
	Username        string `json:"username" structs:"username" mapstructure:"username"`
	Password        string `json:"password" structs:"password" mapstructure:"password"`
	TLS             bool   `json:"tls" structs:"tls" mapstructure:"tls"`
	InsecureTLS     bool   `json:"insecure_tls" structs:"insecure_tls" mapstructure:"insecure_tls"`
	Certificate     string `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	PrivateKey      string `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	IssuingCA       string `json:"issuing_ca" structs:"issuing_ca" mapstructure:"issuing_ca"`
	ProtocolVersion int    `json:"protocol_version" structs:"protocol_version" mapstructure:"protocol_version"`
	ConnectTimeout  int    `json:"connect_timeout" structs:"connect_timeout" mapstructure:"connect_timeout"`
	TLSMinVersion   string `json:"tls_min_version" structs:"tls_min_version" mapstructure:"tls_min_version"`
	Consistency     string `json:"consistency" structs:"consistency" mapstructure:"consistency"`

	config *DatabaseConfig

	session *gocql.Session
	sync.Mutex
}

func (c *cassandraConnectionProducer) Connection() (interface{}, error) {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	// If we already have a DB, we got it!
	if c.session != nil {
		return c.session, nil
	}

	session, err := c.createSession()
	if err != nil {
		return nil, err
	}

	//  Store the session in backend for reuse
	c.session = session

	return session, nil
}

func (c *cassandraConnectionProducer) Close() {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	if c.session != nil {
		c.session.Close()
	}

	c.session = nil
}

func (c *cassandraConnectionProducer) createSession() (*gocql.Session, error) {
	clusterConfig := gocql.NewCluster(strings.Split(c.Hosts, ",")...)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	}

	clusterConfig.ProtoVersion = c.ProtocolVersion
	if clusterConfig.ProtoVersion == 0 {
		clusterConfig.ProtoVersion = 2
	}

	clusterConfig.Timeout = time.Duration(c.ConnectTimeout) * time.Second

	if c.TLS {
		var tlsConfig *tls.Config
		if len(c.Certificate) > 0 || len(c.IssuingCA) > 0 {
			if len(c.Certificate) > 0 && len(c.PrivateKey) == 0 {
				return nil, fmt.Errorf("Found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(c.Certificate) > 0 {
				certBundle.Certificate = c.Certificate
				certBundle.PrivateKey = c.PrivateKey
			}
			if len(c.IssuingCA) > 0 {
				certBundle.IssuingCA = c.IssuingCA
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate bundle: %s", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil || tlsConfig == nil {
				return nil, fmt.Errorf("failed to get TLS configuration: tlsConfig:%#v err:%v", tlsConfig, err)
			}
			tlsConfig.InsecureSkipVerify = c.InsecureTLS

			if c.TLSMinVersion != "" {
				var ok bool
				tlsConfig.MinVersion, ok = tlsutil.TLSLookup[c.TLSMinVersion]
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

	// Set consistency
	if c.Consistency != "" {
		consistencyValue, err := gocql.ParseConsistencyWrapper(c.Consistency)
		if err != nil {
			return nil, err
		}

		session.SetConsistency(consistencyValue)
	}

	// Verify the info
	err = session.Query(`LIST USERS`).Exec()
	if err != nil {
		return nil, fmt.Errorf("Error validating connection info: %s", err)
	}

	return session, nil
}
