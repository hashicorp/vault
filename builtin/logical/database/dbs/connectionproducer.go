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
	"github.com/mitchellh/mapstructure"
)

type ConnectionProducer interface {
	Connection() (interface{}, error)
	Close()
	// TODO: Should we make this immutable instead?
	Reset(*DatabaseConfig) error
}

// sqlConnectionProducer impliments ConnectionProducer and provides a generic producer for most sql databases
type sqlConnectionDetails struct {
	ConnectionURL string `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`
}

type sqlConnectionProducer struct {
	config *DatabaseConfig
	// TODO: Should we merge these two structures make it immutable?
	connDetails *sqlConnectionDetails

	db *sql.DB
	sync.Mutex
}

func (cp *sqlConnectionProducer) Connection() (interface{}, error) {
	// Grab the write lock
	cp.Lock()
	defer cp.Unlock()

	// If we already have a DB, we got it!
	if cp.db != nil {
		if err := cp.db.Ping(); err == nil {
			return cp.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		cp.db.Close()
	}

	// Otherwise, attempt to make connection
	conn := cp.connDetails.ConnectionURL

	// Ensure timezone is set to UTC for all the conenctions
	if strings.HasPrefix(conn, "postgres://") || strings.HasPrefix(conn, "postgresql://") {
		if strings.Contains(conn, "?") {
			conn += "&timezone=utc"
		} else {
			conn += "?timezone=utc"
		}
	}

	var err error
	cp.db, err = sql.Open(cp.config.DatabaseType, conn)
	if err != nil {
		return nil, err
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	cp.db.SetMaxOpenConns(cp.config.MaxOpenConnections)
	cp.db.SetMaxIdleConns(cp.config.MaxIdleConnections)
	cp.db.SetConnMaxLifetime(cp.config.MaxConnectionLifetime)

	return cp.db, nil
}

func (cp *sqlConnectionProducer) Close() {
	// Grab the write lock
	cp.Lock()
	defer cp.Unlock()

	if cp.db != nil {
		cp.db.Close()
	}

	cp.db = nil
}

func (cp *sqlConnectionProducer) Reset(config *DatabaseConfig) error {
	// Grab the write lock
	cp.Lock()

	var details *sqlConnectionDetails
	err := mapstructure.Decode(config.ConnectionDetails, &details)
	if err != nil {
		return err
	}

	cp.connDetails = details
	cp.config = config

	cp.Unlock()

	cp.Close()
	_, err = cp.Connection()
	return err
}

// cassandraConnectionProducer impliments ConnectionProducer and provides connections for cassandra
type cassandraConnectionDetails struct {
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
}

type cassandraConnectionProducer struct {
	config *DatabaseConfig
	// TODO: Should we merge these two structures make it immutable?
	connDetails *cassandraConnectionDetails

	session *gocql.Session
	sync.Mutex
}

func (cp *cassandraConnectionProducer) Connection() (interface{}, error) {
	// Grab the write lock
	cp.Lock()
	defer cp.Unlock()

	// If we already have a DB, we got it!
	if cp.session != nil {
		return cp.session, nil
	}

	session, err := cp.createSession(cp.connDetails)
	if err != nil {
		return nil, err
	}

	//  Store the session in backend for reuse
	cp.session = session

	return session, nil
}

func (cp *cassandraConnectionProducer) Close() {
	// Grab the write lock
	cp.Lock()
	defer cp.Unlock()

	if cp.session != nil {
		cp.session.Close()
	}

	cp.session = nil
}

func (cp *cassandraConnectionProducer) Reset(config *DatabaseConfig) error {
	// Grab the write lock
	cp.Lock()
	cp.config = config
	cp.Unlock()

	cp.Close()
	_, err := cp.Connection()

	return err
}

func (cp *cassandraConnectionProducer) createSession(cfg *cassandraConnectionDetails) (*gocql.Session, error) {
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

	// Set consistency
	if cfg.Consistency != "" {
		consistencyValue, err := gocql.ParseConsistencyWrapper(cfg.Consistency)
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
