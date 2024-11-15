//go:build !go1.23

// Package dial provides types to implement go-hdb custom dialers.
package dial

import (
	"context"
	"net"
	"time"
)

// DialerOptions contains optional parameters that might be used by a Dialer.
type DialerOptions struct {
	Timeout, TCPKeepAlive time.Duration
}

// The Dialer interface needs to be implemented by custom Dialers. A Dialer for providing a custom driver connection
// to the database can be set in the driver.Connector object.
type Dialer interface {
	DialContext(ctx context.Context, address string, options DialerOptions) (net.Conn, error)
}

// DefaultDialer is the default driver Dialer implementation.
var DefaultDialer Dialer = &dialer{}

// default dialer implementation.
type dialer struct{}

func (d *dialer) DialContext(ctx context.Context, address string, options DialerOptions) (net.Conn, error) {
	dialer := net.Dialer{Timeout: options.Timeout, KeepAlive: options.TCPKeepAlive}
	return dialer.DialContext(ctx, "tcp", address)
}
