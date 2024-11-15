//go:build !go1.23

package driver

import (
	"context"
	"crypto/tls"
	"log/slog"
	"maps"
	"net"
	"sync"
	"time"

	"github.com/SAP/go-hdb/driver/dial"
	"github.com/SAP/go-hdb/driver/unicode/cesu8"
	"golang.org/x/text/transform"
)

// connAttrs is holding connection relevant attributes.
type connAttrs struct {
	mu                sync.RWMutex
	_timeout          time.Duration
	_pingInterval     time.Duration
	_bufferSize       int
	_bulkSize         int
	_tcpKeepAlive     time.Duration // see net.Dialer
	_tlsConfig        *tls.Config
	_defaultSchema    string
	_dialer           dial.Dialer
	_applicationName  string
	_sessionVariables map[string]string
	_locale           string
	_fetchSize        int
	_lobChunkSize     int
	_dfv              int
	_cesu8Decoder     func() transform.Transformer
	_cesu8Encoder     func() transform.Transformer
	_emptyDateAsNull  bool
	_logger           *slog.Logger
}

func newConnAttrs() *connAttrs {
	return &connAttrs{
		_timeout:         defaultTimeout,
		_bufferSize:      defaultBufferSize,
		_bulkSize:        defaultBulkSize,
		_tcpKeepAlive:    defaultTCPKeepAlive,
		_dialer:          dial.DefaultDialer,
		_applicationName: defaultApplicationName,
		_fetchSize:       defaultFetchSize,
		_lobChunkSize:    defaultLobChunkSize,
		_dfv:             defaultDfv,
		_cesu8Decoder:    cesu8.DefaultDecoder,
		_cesu8Encoder:    cesu8.DefaultEncoder,
		_logger:          slog.Default(),
	}
}

/*
keep c as the instance name, so that the generated help does have the same variable name when object is
included in connector
*/

func (c *connAttrs) clone() *connAttrs {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &connAttrs{
		_timeout:          c._timeout,
		_pingInterval:     c._pingInterval,
		_bufferSize:       c._bufferSize,
		_bulkSize:         c._bulkSize,
		_tcpKeepAlive:     c._tcpKeepAlive,
		_tlsConfig:        c._tlsConfig.Clone(),
		_defaultSchema:    c._defaultSchema,
		_dialer:           c._dialer,
		_applicationName:  c._applicationName,
		_sessionVariables: maps.Clone(c._sessionVariables),
		_locale:           c._locale,
		_fetchSize:        c._fetchSize,
		_lobChunkSize:     c._lobChunkSize,
		_dfv:              c._dfv,
		_cesu8Decoder:     c._cesu8Decoder,
		_cesu8Encoder:     c._cesu8Encoder,
		_emptyDateAsNull:  c._emptyDateAsNull,
		_logger:           c._logger,
	}
}

func (c *connAttrs) dialContext(ctx context.Context, host string) (net.Conn, error) {
	return c._dialer.DialContext(ctx, host, dial.DialerOptions{Timeout: c._timeout, TCPKeepAlive: c._tcpKeepAlive})
}
