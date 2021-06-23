package mysql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/mitchellh/mapstructure"
)

// mySQLConnectionProducer implements ConnectionProducer and provides a generic producer for most sql databases
type mySQLConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url"          mapstructure:"connection_url"          structs:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections"    mapstructure:"max_open_connections"    structs:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections"    mapstructure:"max_idle_connections"    structs:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`

	Username string `json:"username" mapstructure:"username" structs:"username"`
	Password string `json:"password" mapstructure:"password" structs:"password"`

	TLSCertificateKeyData []byte `json:"tls_certificate_key" mapstructure:"tls_certificate_key" structs:"-"`
	TLSCAData             []byte `json:"tls_ca"              mapstructure:"tls_ca"              structs:"-"`

	// tlsConfigName is a globally unique name that references the TLS config for this instance in the mysql driver
	tlsConfigName string

	RawConfig             map[string]interface{}
	maxConnectionLifetime time.Duration
	Initialized           bool
	db                    *sql.DB
	sync.Mutex
}

func (c *mySQLConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

func (c *mySQLConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	c.RawConfig = conf

	err := mapstructure.WeakDecode(conf, &c)
	if err != nil {
		return nil, err
	}

	if len(c.ConnectionURL) == 0 {
		return nil, fmt.Errorf("connection_url cannot be empty")
	}

	// Don't escape special characters for MySQL password
	password := c.Password

	// QueryHelper doesn't do any SQL escaping, but if it starts to do so
	// then maybe we won't be able to use it to do URL substitution any more.
	c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": url.PathEscape(c.Username),
		"password": password,
	})

	if c.MaxOpenConnections == 0 {
		c.MaxOpenConnections = 4
	}

	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = c.MaxOpenConnections
	}
	if c.MaxIdleConnections > c.MaxOpenConnections {
		c.MaxIdleConnections = c.MaxOpenConnections
	}
	if c.MaxConnectionLifetimeRaw == nil {
		c.MaxConnectionLifetimeRaw = "0s"
	}

	c.maxConnectionLifetime, err = parseutil.ParseDurationSecond(c.MaxConnectionLifetimeRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid max_connection_lifetime: %w", err)
	}

	tlsConfig, err := c.getTLSAuth()
	if err != nil {
		return nil, err
	}

	if tlsConfig != nil {
		if c.tlsConfigName == "" {
			c.tlsConfigName, err = uuid.GenerateUUID()
			if err != nil {
				return nil, fmt.Errorf("unable to generate UUID for TLS configuration: %w", err)
			}
		}

		mysql.RegisterTLSConfig(c.tlsConfigName, tlsConfig)
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			return nil, fmt.Errorf("error verifying connection: %w", err)
		}

		if err := c.db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("error verifying connection: %w", err)
		}
	}

	return c.RawConfig, nil
}

func (c *mySQLConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	// If we already have a DB, test it and return
	if c.db != nil {
		if err := c.db.PingContext(ctx); err == nil {
			return c.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		c.db.Close()
	}

	connURL, err := c.addTLStoDSN()
	if err != nil {
		return nil, err
	}

	c.db, err = sql.Open("mysql", connURL)
	if err != nil {
		return nil, err
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	c.db.SetMaxOpenConns(c.MaxOpenConnections)
	c.db.SetMaxIdleConns(c.MaxIdleConnections)
	c.db.SetConnMaxLifetime(c.maxConnectionLifetime)

	return c.db, nil
}

func (c *mySQLConnectionProducer) SecretValues() map[string]string {
	return map[string]string{
		c.Password: "[password]",
	}
}

// Close attempts to close the connection
func (c *mySQLConnectionProducer) Close() error {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		c.db.Close()
	}

	c.db = nil

	return nil
}

func (c *mySQLConnectionProducer) getTLSAuth() (tlsConfig *tls.Config, err error) {
	if len(c.TLSCAData) == 0 &&
		len(c.TLSCertificateKeyData) == 0 {
		return nil, nil
	}

	rootCertPool := x509.NewCertPool()
	if len(c.TLSCAData) > 0 {
		ok := rootCertPool.AppendCertsFromPEM(c.TLSCAData)
		if !ok {
			return nil, fmt.Errorf("failed to append CA to client options")
		}
	}

	clientCert := make([]tls.Certificate, 0, 1)

	if len(c.TLSCertificateKeyData) > 0 {
		certificate, err := tls.X509KeyPair(c.TLSCertificateKeyData, c.TLSCertificateKeyData)
		if err != nil {
			return nil, fmt.Errorf("unable to load tls_certificate_key_data: %w", err)
		}

		clientCert = append(clientCert, certificate)
	}

	tlsConfig = &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
	}

	return tlsConfig, nil
}

func (c *mySQLConnectionProducer) addTLStoDSN() (connURL string, err error) {
	config, err := mysql.ParseDSN(c.ConnectionURL)
	if err != nil {
		return "", fmt.Errorf("unable to parse connectionURL: %s", err)
	}

	if len(c.tlsConfigName) > 0 {
		config.TLSConfig = c.tlsConfigName
	}

	connURL = config.FormatDSN()

	return connURL, nil
}
