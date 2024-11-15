package driver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"maps"
	"math"
	"os"
	"path"
	"time"

	"github.com/SAP/go-hdb/driver/dial"
	p "github.com/SAP/go-hdb/driver/internal/protocol"
	"github.com/SAP/go-hdb/driver/unicode/cesu8"
	"golang.org/x/text/transform"
)

/*
SessionVariables maps session variables to their values.
All defined session variables will be set once after a database connection is opened.
*/
type SessionVariables map[string]string

// conn attributes default values.
const (
	defaultBufferSize   = 16276             // default value bufferSize.
	defaultBulkSize     = 10000             // default value bulkSize.
	defaultTimeout      = 300 * time.Second // default value connection timeout (300 seconds = 5 minutes).
	defaultTCPKeepAlive = 15 * time.Second  // default TCP keep-alive value (copied from net.dial.go)
)

// minimal / maximal values.
const (
	minTimeout  = 0 * time.Second // minimal timeout value.
	minBulkSize = 1               // minimal bulkSize value.
	maxBulkSize = p.MaxNumArg     // maximum bulk size.
)

const (
	defaultFetchSize    = 128         // Default value fetchSize.
	defaultLobChunkSize = 1 << 16     // Default value lobChunkSize.
	defaultDfv          = p.DfvLevel8 // Default data version format level.
)

const (
	minFetchSize    = 1             // Minimal fetchSize value.
	minLobChunkSize = 128           // Minimal lobChunkSize
	maxLobChunkSize = math.MaxInt32 // Maximal lobChunkSize
)

func (c *connAttrs) setTimeout(timeout time.Duration) {
	if timeout < minTimeout {
		timeout = minTimeout
	}
	c._timeout = timeout
}
func (c *connAttrs) setBulkSize(bulkSize int) {
	switch {
	case bulkSize < minBulkSize:
		bulkSize = minBulkSize
	case bulkSize > maxBulkSize:
		bulkSize = maxBulkSize
	}
	c._bulkSize = bulkSize
}
func (c *connAttrs) setTLS(serverName string, insecureSkipVerify bool, rootCAFiles []string) error {
	c._tlsConfig = &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: insecureSkipVerify, //nolint:gosec
	}
	var certPool *x509.CertPool
	for _, fn := range rootCAFiles {
		rootPEM, err := os.ReadFile(path.Clean(fn))
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
		c._tlsConfig.RootCAs = certPool
	}
	return nil
}
func (c *connAttrs) setDialer(dialer dial.Dialer) {
	if dialer == nil {
		dialer = dial.DefaultDialer
	}
	c._dialer = dialer
}
func (c *connAttrs) setFetchSize(fetchSize int) {
	if fetchSize < minFetchSize {
		fetchSize = minFetchSize
	}
	c._fetchSize = fetchSize
}
func (c *connAttrs) setLobChunkSize(lobChunkSize int) {
	switch {
	case lobChunkSize < minLobChunkSize:
		lobChunkSize = minLobChunkSize
	case lobChunkSize > maxLobChunkSize:
		lobChunkSize = maxLobChunkSize
	}
	c._lobChunkSize = lobChunkSize
}
func (c *connAttrs) setDfv(dfv int) {
	if !p.IsSupportedDfv(dfv) {
		dfv = defaultDfv
	}
	c._dfv = dfv
}

// Timeout returns the timeout of the connector.
func (c *connAttrs) Timeout() time.Duration { c.mu.RLock(); defer c.mu.RUnlock(); return c._timeout }

/*
SetTimeout sets the timeout of the connector.

For more information please see DSNTimeout.
*/
func (c *connAttrs) SetTimeout(timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setTimeout(timeout)
}

// PingInterval returns the connection ping interval of the connector.
func (c *connAttrs) PingInterval() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._pingInterval
}

/*
SetPingInterval sets the connection ping interval value of the connector.

Using a ping interval supports detecting broken connections. In case the ping
is not successful a new or another connection out of the connection pool would
be used automatically instead of retuning an error.

Parameter d defines the time between the pings as duration.
If d is zero no ping is executed. If d is not zero a database ping is executed if
an idle connection out of the connection pool is reused and the time since the
last connection access is greater or equal than d.
*/
func (c *connAttrs) SetPingInterval(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._pingInterval = d
}

// BufferSize returns the bufferSize of the connector.
func (c *connAttrs) BufferSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c._bufferSize }

/*
SetBufferSize sets the bufferSize of the connector.
*/
func (c *connAttrs) SetBufferSize(bufferSize int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._bufferSize = bufferSize
}

// BulkSize returns the bulkSize of the connector.
func (c *connAttrs) BulkSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c._bulkSize }

// SetBulkSize sets the bulkSize of the connector.
func (c *connAttrs) SetBulkSize(bulkSize int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setBulkSize(bulkSize)
}

// TCPKeepAlive returns the tcp keep-alive value of the connector.
func (c *connAttrs) TCPKeepAlive() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._tcpKeepAlive
}

/*
SetTCPKeepAlive sets the tcp keep-alive value of the connector.

For more information please see net.Dialer structure.
*/
func (c *connAttrs) SetTCPKeepAlive(tcpKeepAlive time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._tcpKeepAlive = tcpKeepAlive
}

// DefaultSchema returns the database default schema of the connector.
func (c *connAttrs) DefaultSchema() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._defaultSchema
}

// SetDefaultSchema sets the database default schema of the connector.
func (c *connAttrs) SetDefaultSchema(schema string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._defaultSchema = schema
}

// TLSConfig returns the TLS configuration of the connector.
func (c *connAttrs) TLSConfig() *tls.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._tlsConfig.Clone()
}

// SetTLS sets the TLS configuration of the connector with given parameters. An existing connector TLS configuration is replaced.
func (c *connAttrs) SetTLS(serverName string, insecureSkipVerify bool, rootCAFiles ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.setTLS(serverName, insecureSkipVerify, rootCAFiles)
}

// SetTLSConfig sets the TLS configuration of the connector.
func (c *connAttrs) SetTLSConfig(tlsConfig *tls.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._tlsConfig = tlsConfig.Clone()
}

// Dialer returns the dialer object of the connector.
func (c *connAttrs) Dialer() dial.Dialer { c.mu.RLock(); defer c.mu.RUnlock(); return c._dialer }

// SetDialer sets the dialer object of the connector.
func (c *connAttrs) SetDialer(dialer dial.Dialer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setDialer(dialer)
}

// ApplicationName returns the application name of the connector.
func (c *connAttrs) ApplicationName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._applicationName
}

// SetApplicationName sets the application name of the connector.
func (c *connAttrs) SetApplicationName(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._applicationName = name
}

// SessionVariables returns the session variables stored in connector.
func (c *connAttrs) SessionVariables() SessionVariables {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return maps.Clone(c._sessionVariables)
}

// SetSessionVariables sets the session varibles of the connector.
func (c *connAttrs) SetSessionVariables(sessionVariables SessionVariables) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._sessionVariables = maps.Clone(sessionVariables)
}

// Locale returns the locale of the connector.
func (c *connAttrs) Locale() string { c.mu.RLock(); defer c.mu.RUnlock(); return c._locale }

/*
SetLocale sets the locale of the connector.

For more information please see http://help.sap.com/hana/SAP_HANA_SQL_Command_Network_Protocol_Reference_en.pdf.
*/
func (c *connAttrs) SetLocale(locale string) { c.mu.Lock(); defer c.mu.Unlock(); c._locale = locale }

// FetchSize returns the fetchSize of the connector.
func (c *connAttrs) FetchSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c._fetchSize }

/*
SetFetchSize sets the fetchSize of the connector.

For more information please see DSNFetchSize.
*/
func (c *connAttrs) SetFetchSize(fetchSize int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setFetchSize(fetchSize)
}

// LobChunkSize returns the lobChunkSize of the connector.
func (c *connAttrs) LobChunkSize() int { c.mu.RLock(); defer c.mu.RUnlock(); return c._lobChunkSize }

// SetLobChunkSize sets the lobChunkSize of the connector.
func (c *connAttrs) SetLobChunkSize(lobChunkSize int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setLobChunkSize(lobChunkSize)
}

// Dfv returns the client data format version of the connector.
func (c *connAttrs) Dfv() int { c.mu.RLock(); defer c.mu.RUnlock(); return c._dfv }

// SetDfv sets the client data format version of the connector.
func (c *connAttrs) SetDfv(dfv int) { c.mu.Lock(); defer c.mu.Unlock(); c.setDfv(dfv) }

// CESU8Decoder returns the CESU-8 decoder of the connector.
func (c *connAttrs) CESU8Decoder() func() transform.Transformer {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._cesu8Decoder
}

// SetCESU8Decoder sets the CESU-8 decoder of the connector.
func (c *connAttrs) SetCESU8Decoder(cesu8Decoder func() transform.Transformer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cesu8Decoder == nil {
		cesu8Decoder = cesu8.DefaultDecoder
	}
	c._cesu8Decoder = cesu8Decoder
}

// CESU8Encoder returns the CESU-8 encoder of the connector.
func (c *connAttrs) CESU8Encoder() func() transform.Transformer {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._cesu8Encoder
}

// SetCESU8Encoder sets the CESU-8 encoder of the connector.
func (c *connAttrs) SetCESU8Encoder(cesu8Encoder func() transform.Transformer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cesu8Encoder == nil {
		cesu8Encoder = cesu8.DefaultEncoder
	}
	c._cesu8Encoder = cesu8Encoder
}

/*
EmptyDateAsNull returns NULL for empty dates ('0000-00-00') if true, otherwise:

For data format version 1 the backend does return the NULL indicator for empty date fields.
For data format version non equal 1 (field type daydate) the NULL indicator is not set and the return value is 0.
As value 1 represents '0001-01-01' (the minimal valid date) without setting EmptyDateAsNull '0000-12-31' is returned,
so that NULL, empty and valid dates can be distinguished.

https://help.sap.com/docs/HANA_SERVICE_CF/7c78579ce9b14a669c1f3295b0d8ca16/3f81ccc7e35d44cbbc595c7d552c202a.html?locale=en-US
*/
func (c *connAttrs) EmptyDateAsNull() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._emptyDateAsNull
}

// SetEmptyDateAsNull sets the EmptyDateAsNull flag of the connector.
func (c *connAttrs) SetEmptyDateAsNull(emptyDateAsNull bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c._emptyDateAsNull = emptyDateAsNull
}

// Logger returns the Logger instance of the connector.
func (c *connAttrs) Logger() *slog.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c._logger
}

// SetLogger sets the Logger instance of the connector.
func (c *connAttrs) SetLogger(logger *slog.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if logger == nil {
		logger = slog.Default()
	}
	c._logger = logger
}
