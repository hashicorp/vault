// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql/driver"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/SAP/go-hdb/driver/dial"
	"github.com/SAP/go-hdb/internal/container/vermap"
	p "github.com/SAP/go-hdb/internal/protocol"
)

// Data Format Version values.
// Driver does currently support DfvLevel1, DfvLevel4, DfvLevel6 and DfvLevel8.
const (
	DfvLevel0 = 0 // base data format
	DfvLevel1 = 1 // eval types support all data types
	DfvLevel2 = 2 // reserved, broken, do not use
	DfvLevel3 = 3 // additional types Longdate, Secondate, Daydate, Secondtime supported for NGAP
	DfvLevel4 = 4 // generic support for new date/time types
	DfvLevel5 = 5 // spatial types in ODBC on request
	DfvLevel6 = 6 // BINTEXT
	DfvLevel7 = 7 // with boolean support
	DfvLevel8 = 8 // with FIXED8/12/16 support
)

var supportedDfvs = map[int]bool{DfvLevel1: true, DfvLevel4: true, DfvLevel6: true, DfvLevel8: true}

// Connector default values.
const (
	DefaultDfv          = DfvLevel8         // Default data version format level.
	DefaultTimeout      = 300 * time.Second // Default value connection timeout (300 seconds = 5 minutes).
	DefaultTCPKeepAlive = 15 * time.Second  // Default TCP keep-alive value (copied from net.dial.go)
	DefaultBufferSize   = 16276             // Default value bufferSize.
	DefaultFetchSize    = 128               // Default value fetchSize.
	DefaultBulkSize     = 10000             // Default value bulkSize.
	DefaultLobChunkSize = 4096              // Default value lobChunkSize.
	DefaultLegacy       = true              // Default value legacy.
)

// Connector minimal / maximal values.
const (
	minTimeout      = 0 * time.Second // Minimal timeout value.
	minFetchSize    = 1               // Minimal fetchSize value.
	minBulkSize     = 1               // Minimal bulkSize value.
	MaxBulkSize     = p.MaxNumArg     // Maximum bulk size.
	minLobChunkSize = 128             // Minimal lobChunkSize
	// TODO check maxLobChunkSize
	maxLobChunkSize = 1 << 14 // Maximal lobChunkSize
)

/*
SessionVariables maps session variables to their values.
All defined session variables will be set once after a database connection is opened.
*/
type SessionVariables map[string]string

/*
A Connector represents a hdb driver in a fixed configuration.
A Connector can be passed to sql.OpenDB (starting from go 1.10) allowing users to bypass a string based data source name.
*/
type Connector struct {
	mu                                            sync.RWMutex
	host, username, password                      string
	locale                                        string
	applicationName                               string
	bufferSize, fetchSize, bulkSize, lobChunkSize int
	timeout                                       time.Duration
	dfv                                           int
	pingInterval                                  time.Duration
	tcpKeepAlive                                  time.Duration // see net.Dialer
	tlsConfig                                     *tls.Config
	sessionVariables                              *vermap.VerMap
	defaultSchema                                 string
	legacy                                        bool
	dialer                                        dial.Dialer
}

// newConnector returns a new Connector instance with default values.
func newConnector() *Connector {
	return &Connector{
		applicationName:  defaultApplicationName,
		bufferSize:       DefaultBufferSize,
		fetchSize:        DefaultFetchSize,
		bulkSize:         DefaultBulkSize,
		lobChunkSize:     DefaultLobChunkSize,
		timeout:          DefaultTimeout,
		dfv:              DefaultDfv,
		tcpKeepAlive:     DefaultTCPKeepAlive,
		sessionVariables: vermap.NewVerMap(),
		legacy:           DefaultLegacy,
		dialer:           dial.DefaultDialer,
	}
}

// Connector attributes.
const (
	caDSN           = "dsn"
	caDefaultSchema = "defaultSchema"
	caPingInterval  = "pingInterval"
	caBufferSize    = "bufferSize"
	caBulkSize      = "bulkSize"
)

func stringAttr(attr string, value interface{}) (string, error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		return v.String(), nil
	}
	return "", fmt.Errorf("attribute %s: invalid parameter value %v", attr, value)
}

func int64Attr(attr string, value interface{}) (int64, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := v.Uint()
		if u64 >= 1<<63 {
			return 0, fmt.Errorf("attribute %s: integer out of range %d", attr, value)
		}
		return int64(u64), nil
	}
	return 0, fmt.Errorf("attribute %s: invalid parameter value %v", attr, value)
}

// NewConnector returns a new connector instance setting connector attributes to
// values defined in attrs.
// Example:
//	dsn := "hdb://SYSTEM:MyPassword@localhost:39013"
//	schema:= "MySchema"
//	connector := NewConnector(map[string]interface{}{"dsn": dsn, "defaultSchema": schema}
func NewConnector(attrs map[string]interface{}) (*Connector, error) {
	c := newConnector()

	for attr, value := range attrs {
		switch attr {

		default:
			return nil, fmt.Errorf("connector: invalid attribute: %s", attr)

		case caDSN:
			dsn, err := stringAttr(attr, value)
			if err != nil {
				return nil, err
			}
			if err := c.setDSN(dsn); err != nil {
				return nil, err
			}

		case caDefaultSchema:
			defaultSchema, err := stringAttr(attr, value)
			if err != nil {
				return nil, err
			}
			c.defaultSchema = defaultSchema

		case caPingInterval:
			pingInterval, err := int64Attr(attr, value)
			if err != nil {
				return nil, err
			}
			c.pingInterval = time.Duration(pingInterval)

		case caBufferSize:
			bufferSize, err := int64Attr(attr, value)
			if err != nil {
				return nil, err
			}
			if err := c.setBufferSize(int(bufferSize)); err != nil {
				return nil, err
			}

		case caBulkSize:
			bulkSize, err := int64Attr(attr, value)
			if err != nil {
				return nil, err
			}
			if err := c.setBulkSize(int(bulkSize)); err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

// NewBasicAuthConnector creates a connector for basic authentication.
func NewBasicAuthConnector(host, username, password string) *Connector {
	c := newConnector()
	c.host = host
	c.username = username
	c.password = password
	return c
}

const parseDSNErrorText = "parse dsn error"

// ParseDSNError is the error returned in case DSN is invalid.
type ParseDSNError struct{ err error }

func (e ParseDSNError) Error() string {
	if err := errors.Unwrap(e.err); err != nil {
		return fmt.Sprintf("%s: %s", parseDSNErrorText, err.Error())
	}
	return parseDSNErrorText
}

// Unwrap returns the nested error.
func (e ParseDSNError) Unwrap() error { return e.err }

// NewDSNConnector creates a connector from a data source name.
func NewDSNConnector(dsn string) (*Connector, error) {
	c := newConnector()
	return c, c.setDSN(dsn)
}

func (c *Connector) setDSN(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("invalid DSN parameter error - DSN is empty")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return &ParseDSNError{err}
	}

	c.host = u.Host

	if u.User != nil {
		c.username = u.User.Username()
		c.password, _ = u.User.Password()
	}

	var certPool *x509.CertPool

	for k, v := range u.Query() {
		switch k {

		default:
			return fmt.Errorf("URL parameter %s is not supported", k)

		case DSNFetchSize:
			if len(v) == 0 {
				continue
			}
			fetchSize, err := strconv.Atoi(v[0])
			if err != nil {
				return fmt.Errorf("failed to parse fetchSize: %s", v[0])
			}
			if fetchSize < minFetchSize {
				c.fetchSize = minFetchSize
			} else {
				c.fetchSize = fetchSize
			}

		case DSNTimeout:
			if len(v) == 0 {
				continue
			}
			t, err := strconv.Atoi(v[0])
			if err != nil {
				return fmt.Errorf("failed to parse timeout: %s", v[0])
			}
			timeout := time.Duration(t) * time.Second
			if timeout < minTimeout {
				c.timeout = minTimeout
			} else {
				c.timeout = timeout
			}

		case DSNLocale:
			if len(v) == 0 {
				continue
			}
			c.locale = v[0]

		case DSNTLSServerName:
			if len(v) == 0 {
				continue
			}
			if c.tlsConfig == nil {
				c.tlsConfig = &tls.Config{}
			}
			c.tlsConfig.ServerName = v[0]

		case DSNTLSInsecureSkipVerify:
			if len(v) == 0 {
				continue
			}
			var err error
			b := true
			if v[0] != "" {
				b, err = strconv.ParseBool(v[0])
				if err != nil {
					return fmt.Errorf("failed to parse InsecureSkipVerify (bool): %s", v[0])
				}
			}
			if c.tlsConfig == nil {
				c.tlsConfig = &tls.Config{}
			}
			c.tlsConfig.InsecureSkipVerify = b

		case DSNTLSRootCAFile:
			for _, fn := range v {
				rootPEM, err := ioutil.ReadFile(fn)
				if err != nil {
					return err
				}
				if certPool == nil {
					certPool = x509.NewCertPool()
				}
				if ok := certPool.AppendCertsFromPEM(rootPEM); !ok {
					return fmt.Errorf("failed to parse root certificate - filename: %s", fn)
				}
			}
			if certPool != nil {
				if c.tlsConfig == nil {
					c.tlsConfig = &tls.Config{}
				}
				c.tlsConfig.RootCAs = certPool
			}
		}
	}
	return nil
}

// BasicAuthDSN return the connector DSN for basic authentication.
func (c *Connector) BasicAuthDSN() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	values := url.Values{}
	if c.locale != "" {
		values.Set(DSNLocale, c.locale)
	}
	if c.fetchSize != 0 {
		values.Set(DSNFetchSize, fmt.Sprintf("%d", c.fetchSize))
	}
	if c.timeout != 0 {
		values.Set(DSNTimeout, fmt.Sprintf("%d", c.timeout))
	}
	return (&url.URL{
		Scheme:   DriverName,
		User:     url.UserPassword(c.username, c.password),
		Host:     c.host,
		RawQuery: values.Encode(),
	}).String()
}

// Host returns the host of the connector.
func (c *Connector) Host() string { c.mu.RLock(); defer c.mu.RUnlock(); return c.host }

// Username returns the username of the connector.
func (c *Connector) Username() string { c.mu.RLock(); defer c.mu.RUnlock(); return c.username }

// Password returns the password of the connector.
func (c *Connector) Password() string { c.mu.RLock(); defer c.mu.RUnlock(); return c.password }

// Locale returns the locale of the connector.
func (c *Connector) Locale() string { c.mu.RLock(); defer c.mu.RUnlock(); return c.locale }

/*
SetLocale sets the locale of the connector.

For more information please see DSNLocale.
*/
func (c *Connector) SetLocale(locale string) { c.mu.Lock(); c.locale = locale; c.mu.Unlock() }

// DriverVersion returns the driver version of the connector.
func (c *Connector) DriverVersion() string { return DriverVersion }

// DriverName returns the driver name of the connector.
func (c *Connector) DriverName() string { return DriverName }

// ApplicationName returns the locale of the connector.
func (c *Connector) ApplicationName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.applicationName
}

// SetApplicationName sets the application name of the connector.
func (c *Connector) SetApplicationName(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.applicationName = name
	return nil
}

// BufferSize returns the bufferSize of the connector.
func (c *Connector) BufferSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c.bufferSize }

/*
SetBufferSize sets the bufferSize of the connector.
*/
func (c *Connector) SetBufferSize(bufferSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.setBufferSize(bufferSize)
}

func (c *Connector) setBufferSize(bufferSize int) error {
	c.bufferSize = bufferSize
	return nil
}

// FetchSize returns the fetchSize of the connector.
func (c *Connector) FetchSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c.fetchSize }

/*
SetFetchSize sets the fetchSize of the connector.

For more information please see DSNFetchSize.
*/
func (c *Connector) SetFetchSize(fetchSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if fetchSize < minFetchSize {
		fetchSize = minFetchSize
	}
	c.fetchSize = fetchSize
	return nil
}

// BulkSize returns the bulkSize of the connector.
func (c *Connector) BulkSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c.bulkSize }

/*
SetBulkSize sets the bulkSize of the connector.
*/
func (c *Connector) SetBulkSize(bulkSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.setBulkSize(bulkSize)
}

func (c *Connector) setBulkSize(bulkSize int) error {
	switch {
	case bulkSize < minBulkSize:
		bulkSize = minBulkSize
	case bulkSize > MaxBulkSize:
		bulkSize = MaxBulkSize
	}
	c.bulkSize = bulkSize
	return nil
}

// LobChunkSize returns the lobChunkSize of the connector.
func (c *Connector) LobChunkSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c.lobChunkSize }

// Dialer returns the dialer object of the connector.
func (c *Connector) Dialer() dial.Dialer { c.mu.RLock(); defer c.mu.RUnlock(); return c.dialer }

// SetDialer sets the dialer object of the connector.
func (c *Connector) SetDialer(dialer dial.Dialer) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if dialer == nil {
		dialer = dial.DefaultDialer
	}
	c.dialer = dialer
	return nil
}

// Timeout returns the timeout of the connector.
func (c *Connector) Timeout() time.Duration { c.mu.RLock(); defer c.mu.RUnlock(); return c.timeout }

/*
SetTimeout sets the timeout of the connector.

For more information please see DSNTimeout.
*/
func (c *Connector) SetTimeout(timeout time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout < minTimeout {
		timeout = minTimeout
	}
	c.timeout = timeout
	return nil
}

// PingInterval returns the connection ping interval of the connector.
func (c *Connector) PingInterval() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pingInterval
}

/*
SetPingInterval sets the connection ping interval value of the connector.

If the ping interval is greater than zero, the driver pings all open
connections (active or idle in connection pool) periodically.
Parameter d defines the time between the pings.
*/
func (c *Connector) SetPingInterval(d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pingInterval = d
	return nil
}

// TCPKeepAlive returns the tcp keep-alive value of the connector.
func (c *Connector) TCPKeepAlive() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tcpKeepAlive
}

/*
SetTCPKeepAlive sets the tcp keep-alive value of the connector.

For more information please see net.Dialer structure.
*/
func (c *Connector) SetTCPKeepAlive(tcpKeepAlive time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tcpKeepAlive = tcpKeepAlive
	return nil
}

// Dfv returns the client data format version of the connector.
func (c *Connector) Dfv() int { c.mu.RLock(); defer c.mu.RUnlock(); return c.dfv }

// SetDfv sets the client data format version of the connector.
func (c *Connector) SetDfv(dfv int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := supportedDfvs[dfv]; ok {
		c.dfv = dfv
	} else {
		c.dfv = DefaultDfv
	}
	return nil
}

// TLSConfig returns the TLS configuration of the connector.
func (c *Connector) TLSConfig() *tls.Config { c.mu.RLock(); defer c.mu.RUnlock(); return c.tlsConfig }

// SetTLSConfig sets the TLS configuration of the connector.
func (c *Connector) SetTLSConfig(tlsConfig *tls.Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tlsConfig = tlsConfig
	return nil
}

// SessionVariables returns the session variables stored in connector.
func (c *Connector) SessionVariables() SessionVariables {
	return SessionVariables(c.sessionVariables.Load())
}

// SetSessionVariables sets the session varibles of the connector.
func (c *Connector) SetSessionVariables(sessionVariables SessionVariables) error {
	c.sessionVariables.Store(map[string]string(sessionVariables))
	return nil
}

// DefaultSchema returns the database default schema of the connector.
func (c *Connector) DefaultSchema() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.defaultSchema
}

// SetDefaultSchema sets the database default schema of the connector.
func (c *Connector) SetDefaultSchema(schema string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.defaultSchema = schema
	return nil
}

// Legacy returns the connector legacy flag.
func (c *Connector) Legacy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.legacy
}

// SetLegacy sets the connector legacy flag.
func (c *Connector) SetLegacy(b bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.legacy = b
	return nil
}

// Connect implements the database/sql/driver/Connector interface.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) { return newConn(ctx, c) }

// Driver implements the database/sql/driver/Connector interface.
func (c *Connector) Driver() driver.Driver { return hdbDriver }
