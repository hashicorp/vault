package neo4j

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/url"
	"strings"
	"sync"
	"time"
)

var _ connutil.ConnectionProducer = (*ConnectionProducer)(nil)

type ConnectionProducer struct {
	ConnectionURL                   string      `json:"connection_url" mapstructure:"connection_url" structs:"connection_url"`
	MaxConnectionPoolSize           int         `json:"max_connection_pool_size" mapstructure:"max_connection_pool_size" structs:"max_connection_pool_size"`
	MaxTransactionRetryTimeRaw      interface{} `json:"max_transaction_retry_time" mapstructure:"max_transaction_retry_time" structs:"max_transaction_retry_time"`
	MaxConnectionLifetimeRaw        interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`
	ConnectionAcquisitionTimeoutRaw interface{} `json:"connection_acquisition_timeout" mapstructure:"connection_acquisition_timeout" structs:"connection_acquisition_timeout"`
	SocketConnectTimeoutRaw         interface{} `json:"socket_connect_timeout" mapstructure:"socket_connect_timeout" structs:"socket_connect_timeout"`
	Username                        string      `json:"username" mapstructure:"username" structs:"username"`
	Password                        string      `json:"password" mapstructure:"password" structs:"password"`
	RootCAPemBundle                 string      `json:"root_ca_pem_bundle" structs:"root_ca_pem_bundle" mapstructure:"root_ca_pem_bundle"`
	TLSCertChainPEM                 string      `json:"tls_cert_chain_pem" structs:"tls_cert_chain_pem" mapstructure:"tls_cert_chain_pem"`
	TLSKeyPEM                       string      `json:"tls_key_pem" structs:"tls_key_pem" mapstructure:"tls_key_pem"`
	DatabaseName                    string      `json:"database_name" structs:"database_name" mapstructure:"database_name"`

	RawConfig                    map[string]interface{}
	maxConnectionLifetime        time.Duration
	maxTransactionRetryTime      time.Duration
	connectionAcquisitionTimeout time.Duration
	socketConnectTimeout         time.Duration
	Initialized                  bool
	db                           neo4j.DriverWithContext
	tlsConf                      *tls.Config
	sync.Mutex
}

func (c *ConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	// Pointless interface hangover?
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

func (c *ConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
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
	// Do not allow the username or password template pattern to be used as
	// part of the user-supplied username or password
	if strings.Contains(c.Username, "{{username}}") ||
		strings.Contains(c.Username, "{{password}}") ||
		strings.Contains(c.Password, "{{username}}") ||
		strings.Contains(c.Password, "{{password}}") {

		return nil, fmt.Errorf("username and/or password cannot contain the template variables")
	}

	c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": url.PathEscape(c.Username),
		"password": c.Password,
	})

	if c.MaxConnectionPoolSize == 0 {
		c.MaxConnectionPoolSize = 100
	}

	if c.MaxConnectionLifetimeRaw == nil {
		c.MaxConnectionLifetimeRaw = "1h"
	}

	if c.MaxTransactionRetryTimeRaw == nil {
		c.MaxTransactionRetryTimeRaw = "30s"
	}

	if c.ConnectionAcquisitionTimeoutRaw == nil {
		c.ConnectionAcquisitionTimeoutRaw = "1m"
	}

	if c.SocketConnectTimeoutRaw == nil {
		c.SocketConnectTimeoutRaw = "5s"
	}

	c.maxConnectionLifetime, err = parseutil.ParseDurationSecond(c.MaxConnectionLifetimeRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid max_connection_lifetime: %w", err)
	}

	c.maxTransactionRetryTime, err = parseutil.ParseDurationSecond(c.MaxTransactionRetryTimeRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid max_transaction_retry_time: %w", err)
	}

	c.connectionAcquisitionTimeout, err = parseutil.ParseDurationSecond(c.ConnectionAcquisitionTimeoutRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid connection_acquisition_timeout: %w", err)
	}

	c.socketConnectTimeout, err = parseutil.ParseDurationSecond(c.SocketConnectTimeoutRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid socket_connect_timeout: %w", err)
	}

	c.tlsConf = &tls.Config{}

	if c.TLSCertChainPEM != "" && c.TLSKeyPEM != "" {
		cert, err := tls.X509KeyPair([]byte(c.TLSCertChainPEM), []byte(c.TLSKeyPEM))
		if err != nil {
			return nil, fmt.Errorf("could not parse 'tls_cert_chain_pem' and 'tls_key_pem' as x509 key pair: %w", err)
		}
		c.tlsConf.Certificates = []tls.Certificate{cert}
	} else if c.TLSCertChainPEM+c.TLSKeyPEM != "" {
		return nil, errors.New("both 'tls_cert_chain_pem' and 'tls_key_pem' must be set")
	}

	if c.RootCAPemBundle != "" {
		rootCA := x509.NewCertPool()
		ok := rootCA.AppendCertsFromPEM([]byte(c.RootCAPemBundle))
		if !ok {
			return nil, errors.New("failed to parse 'pem_bundle' as x509 root CAs")
		}
		c.tlsConf.RootCAs = rootCA
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			return nil, fmt.Errorf("error verifying connection: %w", err)
		}

		if err := c.db.VerifyConnectivity(ctx); err != nil {
			return nil, fmt.Errorf("error verifying connection: %w", err)
		}
	}

	return c.RawConfig, nil
}

func (c *ConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	// If we already have a DB, test it and return
	if c.db != nil {
		if err := c.db.VerifyConnectivity(ctx); err == nil {
			return c.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		_ = c.db.Close(ctx)
	}

	config := &neo4j.Config{
		TlsConfig:                    c.tlsConf,
		MaxTransactionRetryTime:      c.maxTransactionRetryTime,
		MaxConnectionPoolSize:        c.MaxConnectionPoolSize,
		MaxConnectionLifetime:        c.maxConnectionLifetime,
		ConnectionAcquisitionTimeout: c.connectionAcquisitionTimeout,
		SocketConnectTimeout:         c.socketConnectTimeout,
		SocketKeepalive:              true,
		UserAgent:                    neo4j.UserAgent,
		FetchSize:                    neo4j.FetchDefault,
	}

	configFunc := func(c *neo4j.Config) {
		*c = *config
	}

	var err error
	c.db, err = neo4j.NewDriverWithContext(c.ConnectionURL, neo4j.BasicAuth(c.Username, c.Password, ""), configFunc)
	if err != nil {
		return nil, err
	}

	return c.db, nil
}

func (c *ConnectionProducer) getConnection(ctx context.Context) (neo4j.DriverWithContext, error) {
	db, err := c.Connection(ctx)
	if err != nil {
		return nil, err
	}
	return db.(neo4j.DriverWithContext), nil
}

func (c *ConnectionProducer) SecretValues() map[string]string {
	return map[string]string{
		c.Password: "[password]",
	}
}

// Close attempts to close the connection
func (c *ConnectionProducer) Close() error {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()
	var err error

	if c.db != nil {
		err = c.db.Close(context.Background())
	}

	c.db = nil

	return err
}
