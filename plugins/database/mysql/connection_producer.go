// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

	cloudmysql "cloud.google.com/go/cloudsqlconn/mysql/mysql"

	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-uuid"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
)

const (
	authTypeIAM   = "iam"
	cloudSQLMySQL = "cloudsql-mysql"
	driverMySQL   = "mysql"
)

// ensure cloud driver registration only happens once.
// the complication here is that registration happens at point X, and is not allowed to happen again, ever.
// This cannot be stored as state within a connection producer, because those producers are config-specific,
// so if we configure, say, two databases that are both cloud-mysql, we must only register once.
// we might be able to cleverly do this with init().
var onceler sync.Once

// mySQLConnectionProducer implements ConnectionProducer and provides a generic producer for most sql databases
type mySQLConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url"          mapstructure:"connection_url"          structs:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections"    mapstructure:"max_open_connections"    structs:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections"    mapstructure:"max_idle_connections"    structs:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`
	Username                 string      `json:"username" mapstructure:"username" structs:"username"`
	Password                 string      `json:"password" mapstructure:"password" structs:"password"`

	TLSCertificateKeyData []byte `json:"tls_certificate_key" mapstructure:"tls_certificate_key" structs:"-"`
	TLSCAData             []byte `json:"tls_ca"              mapstructure:"tls_ca"              structs:"-"`
	TLSServerName         string `json:"tls_server_name" mapstructure:"tls_server_name" structs:"tls_server_name"`
	TLSSkipVerify         bool   `json:"tls_skip_verify" mapstructure:"tls_skip_verify" structs:"tls_skip_verify"`

	// tlsConfigName is a globally unique name that references the TLS config for this instance in the mysql driver
	tlsConfigName string

	// cloudDriverName is a globally unique name that references the cloud dialer config for this instance of the driver
	// I would like to reuse the tlsConfigName value, but there are parts of the code that expect its emptyness to equal "no-tls"
	cloudDriverName string
	isCloud         bool

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
		if _, err = c.Connection(ctx); err != nil {
			return nil, fmt.Errorf("error verifying - connection: %w", err)
		}

		if err := c.db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("error verifying - ping: %w", err)
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

	driverName := driverMySQL
	if c.RawConfig["auth_type"] == authTypeIAM {
		var err error
		c.cloudDriverName, err = uuid.GenerateUUID()
		c.isCloud = true
		if err != nil {
			return nil, fmt.Errorf("unable to generate UUID for connection producer: %w", err)
		}
		driverName = c.cloudDriverName

		filename := c.RawConfig["filename"]
		credentials := c.RawConfig["credentials"]

		// for _most_ sql databases, the driver itself contains no state. In the case of google's cloudsql drivers,
		// however, the driver might store a credentials file, in which case the state stored by the driver is in
		// fact critical to the proper function of the connection.
		//
		// This poses one big obvious problem - each configured cloud database might/will need its own driver registration,
		// the name of which we have to track, and even worse in the MySQL case, that name is also the name of the
		// dialer that is registered by MySQL in the dsn: protocol in the DSN acts double duty if a custom dialer
		// is registered. This means that either we can't hid the name of the dialer from the user OR we have to rewrite
		// the DSN after the user provides it. We already do this, KIND OF for the TLS config, but this modification
		// is much more dramatic.
		_, err = registerDriverMySQL(driverName, filename, credentials)

		if err != nil {
			return nil, err
		}

		//@TODO store driver cleanup
		//c.cleanupDrivers = append(c.cleanupDrivers, cleanupFunc)
	}

	connURL, err := c.addTLStoDSN()
	if err != nil {
		return nil, err
	}

	cloudURL, err := c.rewriteProtocolForGCP(connURL)
	if err != nil {
		return nil, err
	}

	c.db, err = sql.Open(driverName, cloudURL)
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
		// if auth_type is IAM, ensure cleanup
		// of cloudSQL resources
		if c.RawConfig["auth_type"] == authTypeIAM {
			// @TODO implement cloudSQL Driver cleanup from cache
		} else {
			c.db.Close()
		}
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
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		ServerName:         c.TLSServerName,
		InsecureSkipVerify: c.TLSSkipVerify,
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

// rewriteProtocolForGCP rewrites the protocl in the DSN to contain the protocol name associated
// with the dialer and therefore driver associated with the provided cloudsqlconn.DialerOpts
func (c *mySQLConnectionProducer) rewriteProtocolForGCP(inDSN string) (string, error) {
	config, err := mysql.ParseDSN(inDSN)
	if err != nil {
		return "", fmt.Errorf("unable to parseeeee connectionURL: %s", err)
	}

	if config.Net != "cloudsql-mysql" {
		return "", fmt.Errorf("didn't update net name because it wasn't what we expected as a placeholder: %s", config.Net)
	}

	if c.isCloud {
		config.Net = c.cloudDriverName // todo format to something reasonable
	}

	return config.FormatDSN(), nil
}

func registerDriverMySQL(driverName string, filename, credentials interface{}) (func() error, error) {
	// @TODO implement driver cleanup cache
	// if driver is already registered, return

	opts, err := connutil.GetCloudSQLAuthOptions(filename, credentials)
	if err != nil {
		return nil, err
	}
	return cloudmysql.RegisterDriver(driverName, opts...)
}
