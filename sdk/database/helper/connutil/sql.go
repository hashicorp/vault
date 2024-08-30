// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/cacheutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/mitchellh/mapstructure"
)

const (
	AuthTypeGCPIAM           = "gcp_iam"
	AuthTypeCert             = "cert"
	AuthTypeUsernamePassword = ""
)

const (
	dbTypePostgres   = "pgx"
	cloudSQLPostgres = "cloudsql-postgres"

	// controls the size of the static account cache
	// as part of the self-managed workflow
	defaultStaticCacheSize     = 4
	defaultSelfManagedUsername = "self-managed-user"
	defaultSelfManagedPassword = "self-managed-password"
)

var _ ConnectionProducer = &SQLConnectionProducer{}

// SQLConnectionProducer implements ConnectionProducer and provides a generic producer for most sql databases
type SQLConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url" mapstructure:"connection_url" structs:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections" mapstructure:"max_open_connections" structs:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections" mapstructure:"max_idle_connections" structs:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`
	DisableEscaping          bool        `json:"disable_escaping" mapstructure:"disable_escaping" structs:"disable_escaping"`
	usePrivateIP             bool        `json:"use_private_ip" mapstructure:"use_private_ip" structs:"use_private_ip"`
	SelfManaged              bool        `json:"self_managed" mapstructure:"self_managed" structs:"self_managed"`

	// Username/Password is the default auth type when AuthType is not set
	Username string `json:"username" mapstructure:"username" structs:"username"`
	Password string `json:"password" mapstructure:"password" structs:"password"`

	// AuthType defines the type of client authenticate used for this connection
	AuthType           string `json:"auth_type" mapstructure:"auth_type" structs:"auth_type"`
	ServiceAccountJSON string `json:"service_account_json" mapstructure:"service_account_json" structs:"service_account_json"`
	TLSConfig          *tls.Config

	// cloudDriverName is globally unique, but only needs to be retained for the lifetime
	// of driver registration, not across plugin restarts.
	cloudDriverName    string
	cloudDialerCleanup func() error

	Type                  string
	RawConfig             map[string]interface{}
	maxConnectionLifetime time.Duration
	Initialized           bool
	db                    *sql.DB
	staticAccountsCache   *cacheutil.Cache
	sync.Mutex
}

func (c *SQLConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

func (c *SQLConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
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

	isTemplatedURL := true
	if !strings.Contains(c.ConnectionURL, "{{username}}") || !strings.Contains(c.ConnectionURL, "{{password}}") {
		isTemplatedURL = false
	}

	// Do not allow the username or password template pattern to be used as
	// part of the user-supplied username or password
	if strings.Contains(c.Username, "{{username}}") ||
		strings.Contains(c.Username, "{{password}}") ||
		strings.Contains(c.Password, "{{username}}") ||
		strings.Contains(c.Password, "{{password}}") {

		return nil, fmt.Errorf("username and/or password cannot contain the template variables")
	}

	// validate that at least one of username/password / self_managed is set
	if !c.SelfManaged && (c.Username == "" && c.Password == "") && isTemplatedURL {
		return nil, fmt.Errorf("must either provide username/password or set self-managed to 'true'")
	}

	// validate that self-managed and username/password are mutually exclusive
	if c.SelfManaged {
		if (c.Username != "" || c.Password != "") || !isTemplatedURL {
			return nil, fmt.Errorf("cannot use both self-managed and vault-managed workflows")
		}
	}

	var username string
	var password string
	if !c.SelfManaged {
		// Default behavior
		username = c.Username
		password = c.Password

		// Don't escape special characters for MySQL password
		// Also don't escape special characters for the username and password if
		// the disable_escaping parameter is set to true
		if !c.DisableEscaping {
			username = url.PathEscape(c.Username)
		}
		if (c.Type != "mysql") && !c.DisableEscaping {
			password = url.PathEscape(c.Password)
		}

	} else {
		// this is added to make middleware happy
		// these placeholders are replaced when we make the actual static connection
		username = defaultSelfManagedUsername
		password = defaultSelfManagedPassword
	}

	// QueryHelper doesn't do any SQL escaping, but if it starts to do so
	// then maybe we won't be able to use it to do URL substitution any more.
	c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": username,
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
		return nil, errwrap.Wrapf("invalid max_connection_lifetime: {{err}}", err)
	}

	if ok := ValidateAuthType(c.AuthType); !ok {
		return nil, fmt.Errorf("invalid auth_type: %s", c.AuthType)
	}

	if c.AuthType == AuthTypeGCPIAM {
		c.cloudDriverName, err = uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("unable to generate UUID for IAM configuration: %w", err)
		}

		// for _most_ sql databases, the driver itself contains no state. In the case of google's cloudsql drivers,
		// however, the driver might store a credentials file, in which case the state stored by the driver is in
		// fact critical to the proper function of the connection. So it needs to be registered here inside the
		// ConnectionProducer init.
		dialerCleanup, err := c.registerDrivers(c.cloudDriverName, c.ServiceAccountJSON, c.usePrivateIP)
		if err != nil {
			return nil, err
		}

		c.cloudDialerCleanup = dialerCleanup
	}

	if c.SelfManaged && c.staticAccountsCache == nil {
		logger := log.New(&log.LoggerOptions{
			Level: log.Trace,
		})

		closer := func(key interface{}, value interface{}) {
			logger.Trace(fmt.Sprintf("Evicting key %s from static LRU cache", key))
			conn, ok := value.(*sql.DB)
			if !ok {
				logger.Error(fmt.Sprintf("error retrieving connection %s from static LRU cache, err=%s", key, err))
			}

			if err := conn.Close(); err != nil {
				logger.Error(fmt.Sprintf("error closing connection for %s, err=%s", key, err))
			}
			logger.Trace(fmt.Sprintf("closed DB connection for %s", key))
		}
		c.staticAccountsCache, err = cacheutil.NewCache(defaultStaticCacheSize, closer)
		if err != nil {
			return nil, fmt.Errorf("error initializing static account cache: %s", err)
		}
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	// only verify if not self-managed
	if verifyConnection && !c.SelfManaged {
		if _, err := c.Connection(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}

		if err := c.db.PingContext(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: ping failed: {{err}}", err)
		}
	}

	return c.RawConfig, nil
}

func (c *SQLConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	if !c.Initialized {
		return nil, ErrNotInitialized
	}

	// If we already have a DB, test it and return
	if c.db != nil {
		if err := c.db.PingContext(ctx); err == nil {
			return c.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		c.db.Close()

		// if IAM authentication is enabled
		// ensure open dialer is also closed
		if c.AuthType == AuthTypeGCPIAM {
			if c.cloudDialerCleanup != nil {
				c.cloudDialerCleanup()
			}
		}
	}

	// default non-IAM behavior
	driverName := c.Type

	if c.AuthType == AuthTypeGCPIAM {
		driverName = c.cloudDriverName
	} else if c.Type == "mssql" {
		// For mssql backend, switch to sqlserver instead
		driverName = "sqlserver"
	}

	// Otherwise, attempt to make connection
	// Apply PostgreSQL specific settings if needed
	conn := applyPostgresSettings(c.ConnectionURL)

	if driverName == dbTypePostgres && c.TLSConfig != nil {
		config, err := pgx.ParseConfig(conn)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
		if config.TLSConfig == nil {
			// handle sslmode=disable
			config.TLSConfig = &tls.Config{}
		}

		config.TLSConfig.RootCAs = c.TLSConfig.RootCAs
		config.TLSConfig.ClientCAs = c.TLSConfig.ClientCAs
		config.TLSConfig.Certificates = c.TLSConfig.Certificates

		// Ensure there are no stale fallbacks when manually setting TLSConfig
		for _, fallback := range config.Fallbacks {
			fallback.TLSConfig = config.TLSConfig
		}

		c.db = stdlib.OpenDB(*config)
		if err != nil {
			return nil, fmt.Errorf("failed to open connection: %w", err)
		}
	} else if driverName == dbTypePostgres && os.Getenv(pluginutil.PluginUsePostgresSSLInline) != "" {
		var err error
		// TODO: remove this deprecated function call in a future SDK version
		c.db, err = openPostgres(driverName, conn)
		if err != nil {
			return nil, fmt.Errorf("failed to open connection: %w", err)
		}
	} else {
		var err error
		c.db, err = sql.Open(driverName, conn)
		if err != nil {
			return nil, fmt.Errorf("failed to open connection: %w", err)
		}
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	c.db.SetMaxOpenConns(c.MaxOpenConnections)
	c.db.SetMaxIdleConns(c.MaxIdleConnections)
	c.db.SetConnMaxLifetime(c.maxConnectionLifetime)

	return c.db, nil
}

func (c *SQLConnectionProducer) SecretValues() map[string]interface{} {
	return map[string]interface{}{
		c.Password: "[password]",
	}
}

// Close attempts to close the connection
func (c *SQLConnectionProducer) Close() error {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		c.db.Close()

		// cleanup IAM dialer if it exists
		if c.AuthType == AuthTypeGCPIAM {
			if c.cloudDialerCleanup != nil {
				c.cloudDialerCleanup()
			}
		}
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
func (c *SQLConnectionProducer) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	return "", "", dbutil.Unimplemented()
}

func applyPostgresSettings(connURL string) string {
	res := connURL
	if strings.HasPrefix(res, "postgres://") || strings.HasPrefix(res, "postgresql://") {
		// Ensure timezone is set to UTC for all the connections
		if strings.Contains(res, "?") {
			res += "&timezone=UTC"
		} else {
			res += "?timezone=UTC"
		}

		// Ensure a reasonable application_name is set
		if !strings.Contains(res, "application_name") {
			res += "&application_name=vault"
		}
	}

	return res
}

var configurableAuthTypes = map[string]bool{
	AuthTypeUsernamePassword: true,
	AuthTypeCert:             true,
	AuthTypeGCPIAM:           true,
}

func ValidateAuthType(authType string) bool {
	return configurableAuthTypes[authType]
}
