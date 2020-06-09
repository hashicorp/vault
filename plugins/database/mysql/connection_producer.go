package mysql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	stdmysql "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/mitchellh/mapstructure"
)

// SQLConnectionProducer implements ConnectionProducer and provides a generic producer for most sql databases
type mySQLConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url" mapstructure:"connection_url" structs:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections" mapstructure:"max_open_connections" structs:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections" mapstructure:"max_idle_connections" structs:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`

	Username                 string      `json:"username" mapstructure:"username" structs:"username"`
	Password                 string      `json:"password" mapstructure:"password" structs:"password"`

	TLSCertificateKeyData []byte `json:"tls_certificate_key" structs:"-" mapstructure:"tls_certificate_key"`
	TLSCertificateData		[]byte `json:"tls_client_cert" structs:"-" mapstructure:"tls_client_cert"`
	TLSKeyData						[]byte `json:"tls_client_key" structs:"-" mapstructure:"tls_client_key"`
	TLSCAData							[]byte `json:"tls_ca"							 structs:"-" mapstructure:"tls_ca"`
	TLSConfigName					string

	Type                  string
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

	password := c.Password

	// QueryHelper doesn't do any SQL escaping, but if it starts to do so
	// then maybe we won't be able to use it to do URL substitution any more.
	c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": url.PathEscape(c.Username),
		"password": password,
	})

	c.Type = "mysql"

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
		return nil, errwrap.Wrapf("invalid max_connection_lifetime: {{err}}", err)
	}

	err = c.getTLSAuth()
	if err != nil {
		return nil, err
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}

		if err := c.db.PingContext(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
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

	dbType := c.Type

	// Otherwise, attempt to make connection
	conn := c.ConnectionURL

	var err error
	c.db, err = sql.Open(dbType, conn)
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

func (c *mySQLConnectionProducer) SecretValues() map[string]interface{} {
	return map[string]interface{}{
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

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (c *mySQLConnectionProducer) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	return "", "", dbutil.Unimplemented()
}

func (c *mySQLConnectionProducer) getTLSAuth() (err error) {
	if len(c.TLSCAData) == 0 && 
	   len(c.TLSCertificateKeyData) == 0 &&
		 (len(c.TLSCertificateData) == 0 || len(c.TLSKeyData) == 0) {
		return nil
	}

	tlsConfig := &tls.Config{}

	if len(c.TLSCAData) > 0 {
		tlsConfig.RootCAs = x509.NewCertPool()

		ok := tlsConfig.RootCAs.AppendCertsFromPEM(c.TLSCAData)
		if !ok {
			return fmt.Errorf("failed to append CA to client options")
		}
	}

	if len(c.TLSCertificateKeyData) > 0 {
		certificate, err := tls.X509KeyPair(c.TLSCertificateKeyData, c.TLSCertificateKeyData)
		if err != nil {
			return fmt.Errorf("unable to load tls_certificate_key_data: %w", err)
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
	} else if len(c.TLSCertificateData) > 0 && len(c.TLSKeyData) > 0 {
		certificate, err := tls.X509KeyPair(c.TLSCertificateData, c.TLSKeyData)
		if err != nil {
			return fmt.Errorf("unable to load tls_certificate_data or tls_key_data: %w", err)
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
	}

	if c.TLSConfigName == "" {
		c.TLSConfigName = "custom"
	}

	stdmysql.RegisterTLSConfig(c.TLSConfigName, tlsConfig)

	if !strings.Contains(c.ConnectionURL, "tls=") {
		if !strings.Contains(c.ConnectionURL, "?") {
			c.ConnectionURL += "?"
		} else {
			c.ConnectionURL += "&"
		}
		c.ConnectionURL += "tls=" + c.TLSConfigName
	}
	return nil
}
