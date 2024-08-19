// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"
)

const (
	cloudSQLMySQL = "cloudsql-mysql"
	driverMySQL   = "mysql"
)

// mySQLConnectionProducer implements ConnectionProducer and provides a generic producer for most sql databases
type mySQLConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url"          mapstructure:"connection_url"          structs:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections"    mapstructure:"max_open_connections"    structs:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections"    mapstructure:"max_idle_connections"    structs:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`
	Username                 string      `json:"username" mapstructure:"username" structs:"username"`
	Password                 string      `json:"password" mapstructure:"password" structs:"password"`
	AuthType                 string      `json:"auth_type" mapstructure:"auth_type" structs:"auth_type"`
	ServiceAccountJSON       string      `json:"service_account_json" mapstructure:"service_account_json" structs:"service_account_json"`

	TLSCertificateKeyData []byte `json:"tls_certificate_key" mapstructure:"tls_certificate_key" structs:"-"`
	TLSCAData             []byte `json:"tls_ca"              mapstructure:"tls_ca"              structs:"-"`
	TLSServerName         string `json:"tls_server_name" mapstructure:"tls_server_name" structs:"tls_server_name"`
	TLSSkipVerify         bool   `json:"tls_skip_verify" mapstructure:"tls_skip_verify" structs:"tls_skip_verify"`

	// tlsConfigName is a globally unique name that references the TLS config for this instance in the mysql driver
	tlsConfigName string

	// cloudDriverName is a globally unique name that references the cloud dialer config for this instance of the driver
	cloudDriverName    string
	cloudDialerCleanup func() error

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

	// validate auth_type if provided
	if ok := connutil.ValidateAuthType(c.AuthType); !ok {
		return nil, fmt.Errorf("invalid auth_type: %s", c.AuthType)
	}

	if c.AuthType == connutil.AuthTypeGCPIAM {
		c.cloudDriverName, err = uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("unable to generate UUID for IAM configuration: %w", err)
		}

		// for _most_ sql databases, the driver itself contains no state. In the case of google's cloudsql drivers,
		// however, the driver might store a credentials file, in which case the state stored by the driver is in
		// fact critical to the proper function of the connection. So it needs to be registered here inside the
		// ConnectionProducer init.
		dialerCleanup, err := registerDriverMySQL(c.cloudDriverName, c.ServiceAccountJSON)
		if err != nil {
			return nil, err
		}

		c.cloudDialerCleanup = dialerCleanup
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

		// if IAM authentication was enabled
		// ensure open dialer is also closed
		if c.AuthType == connutil.AuthTypeGCPIAM {
			if c.cloudDialerCleanup != nil {
				c.cloudDialerCleanup()
			}
		}

	}

	driverName := driverMySQL
	if c.cloudDriverName != "" {
		driverName = c.cloudDriverName
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
		if c.AuthType == connutil.AuthTypeGCPIAM {
			if c.cloudDialerCleanup != nil {
				c.cloudDialerCleanup()
			}
		}
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

// rewriteProtocolForGCP rewrites the protocol in the DSN to contain the protocol name associated
// with the dialer and therefore driver associated with the provided cloudsqlconn.DialerOpts.
// As a safety/sanity check, it will only do this for protocol "cloudsql-mysql", the name GCP uses in its documentation.
//
// For example, it will rewrite the dsn "user@cloudsql-mysql(zone:region:instance)/ to
// "user@the-uuid-generated(zone:region:instance)/
func (c *mySQLConnectionProducer) rewriteProtocolForGCP(inDSN string) (string, error) {
	if c.cloudDriverName == "" {
		// unchanged if not cloud
		return inDSN, nil
	}

	config, err := mysql.ParseDSN(inDSN)
	if err != nil {
		return "", fmt.Errorf("unable to parse connectionURL: %s", err)
	}

	if config.Net != cloudSQLMySQL {
		return "", fmt.Errorf("didn't update net name because it wasn't what we expected as a placeholder: %s", config.Net)
	}

	config.Net = c.cloudDriverName

	return config.FormatDSN(), nil
}

func registerDriverMySQL(driverName, credentials string) (cleanup func() error, err error) {
	opts, err := connutil.GetCloudSQLAuthOptions(credentials, false)
	if err != nil {
		return nil, err
	}

	return cloudmysql.RegisterDriver(driverName, opts...)
}
