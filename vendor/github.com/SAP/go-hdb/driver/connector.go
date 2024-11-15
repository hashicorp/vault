package driver

import (
	"context"
	"database/sql/driver"
	"sync"

	"github.com/SAP/go-hdb/driver/internal/protocol/auth"
)

type redirectCacheKey struct {
	host, databaseName string
}

var redirectCache sync.Map

/*
A Connector represents a hdb driver in a fixed configuration.
A Connector can be passed to sql.OpenDB allowing users to bypass a string based data source name.
*/
type Connector struct {
	_host         string
	_databaseName string

	*connAttrs
	*authAttrs

	metrics *metrics
}

// NewConnector returns a new Connector instance with default values.
func NewConnector() *Connector {
	return &Connector{
		connAttrs: newConnAttrs(),
		authAttrs: &authAttrs{},
		metrics:   stdHdbDriver.metrics, // use default stdHdbDriver metrics
	}
}

// NewBasicAuthConnector creates a connector for basic authentication.
func NewBasicAuthConnector(host, username, password string) *Connector {
	c := NewConnector()
	c._host = host
	c._username = username
	c._password = password
	return c
}

// NewX509AuthConnector creates a connector for X509 (client certificate) authentication.
// Parameters clientCert and clientKey in PEM format, clientKey not password encryped.
func NewX509AuthConnector(host string, clientCert, clientKey []byte) (*Connector, error) {
	c := NewConnector()
	c._host = host
	var err error
	if c._certKey, err = auth.NewCertKey(clientCert, clientKey); err != nil {
		return nil, err
	}
	return c, nil
}

// NewX509AuthConnectorByFiles creates a connector for X509 (client certificate) authentication
// based on client certificate and client key files.
// Parameters clientCertFile and clientKeyFile in PEM format, clientKeyFile not password encryped.
func NewX509AuthConnectorByFiles(host, clientCertFile, clientKeyFile string) (*Connector, error) {
	c := NewConnector()
	c._host = host

	c._certKeyFiles = newCertKeyFiles(clientCertFile, clientKeyFile)
	clientCert, clientKey, err := c._certKeyFiles.read()
	if err != nil {
		return nil, err
	}
	if c._certKey, err = auth.NewCertKey(clientCert, clientKey); err != nil {
		return nil, err
	}
	return c, nil
}

// NewJWTAuthConnector creates a connector for token (JWT) based authentication.
func NewJWTAuthConnector(host, token string) *Connector {
	c := NewConnector()
	c._host = host
	c._token = token
	return c
}

func newDSNConnector(dsn *DSN) (*Connector, error) {
	c := NewConnector()
	c._host = dsn.host
	c._databaseName = dsn.databaseName
	c._pingInterval = dsn.pingInterval
	c._defaultSchema = dsn.defaultSchema
	c.setTimeout(dsn.timeout)
	if dsn.tls != nil {
		if err := c.connAttrs.setTLS(dsn.tls.ServerName, dsn.tls.InsecureSkipVerify, dsn.tls.RootCAFiles); err != nil {
			return nil, err
		}
	}
	c._username = dsn.username
	c._password = dsn.password
	return c, nil
}

// NewDSNConnector creates a connector from a data source name.
func NewDSNConnector(dsnStr string) (*Connector, error) {
	dsn, err := ParseDSN(dsnStr)
	if err != nil {
		return nil, err
	}
	return newDSNConnector(dsn)
}

// NativeDriver returns the concrete underlying Driver of the Connector.
func (c *Connector) NativeDriver() Driver { return stdHdbDriver }

// Host returns the host of the connector.
func (c *Connector) Host() string { return c._host }

// DatabaseName returns the tenant database name of the connector.
func (c *Connector) DatabaseName() string { return c._databaseName }

func (c *Connector) redirect(ctx context.Context) (driver.Conn, error) {
	connAttrs := c.connAttrs.clone()

	if redirectHost, found := redirectCache.Load(redirectCacheKey{host: c._host, databaseName: c._databaseName}); found {
		if conn, err := connect(ctx, redirectHost.(string), c.metrics, connAttrs, c.authAttrs); err == nil {
			return conn, nil
		}
	}

	redirectHost, err := fetchRedirectHost(ctx, c._host, c._databaseName, c.metrics, connAttrs)
	if err != nil {
		return nil, err
	}
	conn, err := connect(ctx, redirectHost, c.metrics, connAttrs, c.authAttrs)
	if err != nil {
		return nil, err
	}

	redirectCache.Store(redirectCacheKey{host: c._host, databaseName: c._databaseName}, redirectHost)

	return conn, err
}

// Connect implements the database/sql/driver/Connector interface.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	if c._databaseName != "" {
		return c.redirect(ctx)
	}
	return connect(ctx, c._host, c.metrics, c.connAttrs.clone(), c.authAttrs)
}

// Driver implements the database/sql/driver/Connector interface.
func (c *Connector) Driver() driver.Driver { return stdHdbDriver }

func (c *Connector) clone() *Connector {
	return &Connector{
		_host:         c._host,
		_databaseName: c._databaseName,
		connAttrs:     c.connAttrs.clone(),
		authAttrs:     c.authAttrs.clone(),
		metrics:       c.metrics,
	}
}

// WithDatabase returns a new Connector supporting tenant database connections via database name.
func (c *Connector) WithDatabase(databaseName string) *Connector {
	nc := c.clone()
	nc._databaseName = databaseName
	return nc
}
