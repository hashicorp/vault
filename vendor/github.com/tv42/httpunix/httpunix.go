// Package httpunix provides a HTTP transport (net/http.RoundTripper)
// that uses Unix domain sockets instead of HTTP.
//
// This is useful for non-browser connections within the same host, as
// it allows using the file system for credentials of both client
// and server, and guaranteeing unique names.
//
// The URLs look like this:
//
//     http+unix://LOCATION/PATH_ETC
//
// where LOCATION is translated to a file system path with
// Transport.RegisterLocation, and PATH_ETC follow normal http: scheme
// conventions.
package httpunix

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"
)

// Scheme is the URL scheme used for HTTP over UNIX domain sockets.
const Scheme = "http+unix"

// Transport is a http.RoundTripper that connects to Unix domain
// sockets.
type Transport struct {
	// DialTimeout is deprecated. Use context instead.
	DialTimeout time.Duration
	// RequestTimeout is deprecated and has no effect.
	RequestTimeout time.Duration
	// ResponseHeaderTimeout is deprecated. Use context instead.
	ResponseHeaderTimeout time.Duration

	onceInit  sync.Once
	transport http.Transport

	mu sync.Mutex
	// map a URL "hostname" to a UNIX domain socket path
	loc map[string]string
}

func (t *Transport) initTransport() {
	t.transport.DialContext = t.dialContext
	t.transport.DialTLS = t.dialTLS
	t.transport.DisableCompression = true
	t.transport.ResponseHeaderTimeout = t.ResponseHeaderTimeout
}

func (t *Transport) getTransport() *http.Transport {
	t.onceInit.Do(t.initTransport)
	return &t.transport
}

func (t *Transport) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if network != "tcp" {
		return nil, errors.New("httpunix internals are confused: network=" + network)
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if port != "80" {
		return nil, errors.New("httpunix internals are confused: port=" + port)
	}
	t.mu.Lock()
	path, ok := t.loc[host]
	t.mu.Unlock()
	if !ok {
		return nil, errors.New("unknown location: " + host)
	}
	d := net.Dialer{
		Timeout: t.DialTimeout,
	}
	return d.DialContext(ctx, "unix", path)
}

func (t *Transport) dialTLS(network, addr string) (net.Conn, error) {
	return nil, errors.New("httpunix: TLS over UNIX domain sockets is not supported")
}

// RegisterLocation registers an URL location and maps it to the given
// file system path.
//
// Calling RegisterLocation twice for the same location is a
// programmer error, and causes a panic.
func (t *Transport) RegisterLocation(loc string, path string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.loc == nil {
		t.loc = make(map[string]string)
	}
	if _, exists := t.loc[loc]; exists {
		panic("location " + loc + " already registered")
	}
	t.loc[loc] = path
}

var _ http.RoundTripper = (*Transport)(nil)

// RoundTrip executes a single HTTP transaction. See
// net/http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil {
		return nil, errors.New("http+unix: nil Request.URL")
	}
	if req.URL.Scheme != Scheme {
		return nil, errors.New("unsupported protocol scheme: " + req.URL.Scheme)
	}
	if req.URL.Host == "" {
		return nil, errors.New("http+unix: no Host in request URL")
	}

	tt := t.getTransport()
	req = req.Clone(req.Context())
	// get http.Transport to cooperate
	req.URL.Scheme = "http"
	return tt.RoundTrip(req)
}
