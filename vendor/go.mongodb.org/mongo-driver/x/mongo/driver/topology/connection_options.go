// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/httputil"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
)

// Dialer is used to make network connections.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// DialerFunc is a type implemented by functions that can be used as a Dialer.
type DialerFunc func(ctx context.Context, network, address string) (net.Conn, error)

// DialContext implements the Dialer interface.
func (df DialerFunc) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return df(ctx, network, address)
}

// DefaultDialer is the Dialer implementation that is used by this package. Changing this
// will also change the Dialer used for this package. This should only be changed why all
// of the connections being made need to use a different Dialer. Most of the time, using a
// WithDialer option is more appropriate than changing this variable.
var DefaultDialer Dialer = &net.Dialer{}

// Handshaker is the interface implemented by types that can perform a MongoDB
// handshake over a provided driver.Connection. This is used during connection
// initialization. Implementations must be goroutine safe.
type Handshaker = driver.Handshaker

// generationNumberFn is a callback type used by a connection to fetch its generation number given its service ID.
type generationNumberFn func(serviceID *primitive.ObjectID) uint64

type connectionConfig struct {
	connectTimeout           time.Duration
	dialer                   Dialer
	handshaker               Handshaker
	idleTimeout              time.Duration
	cmdMonitor               *event.CommandMonitor
	readTimeout              time.Duration
	writeTimeout             time.Duration
	tlsConfig                *tls.Config
	httpClient               *http.Client
	compressors              []string
	zlibLevel                *int
	zstdLevel                *int
	ocspCache                ocsp.Cache
	disableOCSPEndpointCheck bool
	tlsConnectionSource      tlsConnectionSource
	loadBalanced             bool
	getGenerationFn          generationNumberFn
}

func newConnectionConfig(opts ...ConnectionOption) *connectionConfig {
	cfg := &connectionConfig{
		connectTimeout:      30 * time.Second,
		dialer:              nil,
		tlsConnectionSource: defaultTLSConnectionSource,
		httpClient:          httputil.DefaultHTTPClient,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cfg)
	}

	if cfg.dialer == nil {
		// Use a zero value of net.Dialer when nothing is specified, so the Go driver applies default default behaviors
		// such as Timeout, KeepAlive, DNS resolving, etc. See https://golang.org/pkg/net/#Dialer for more information.
		cfg.dialer = &net.Dialer{}
	}

	return cfg
}

// ConnectionOption is used to configure a connection.
type ConnectionOption func(*connectionConfig)

func withTLSConnectionSource(fn func(tlsConnectionSource) tlsConnectionSource) ConnectionOption {
	return func(c *connectionConfig) {
		c.tlsConnectionSource = fn(c.tlsConnectionSource)
	}
}

// WithCompressors sets the compressors that can be used for communication.
func WithCompressors(fn func([]string) []string) ConnectionOption {
	return func(c *connectionConfig) {
		c.compressors = fn(c.compressors)
	}
}

// WithConnectTimeout configures the maximum amount of time a dial will wait for a
// Connect to complete. The default is 30 seconds.
func WithConnectTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.connectTimeout = fn(c.connectTimeout)
	}
}

// WithDialer configures the Dialer to use when making a new connection to MongoDB.
func WithDialer(fn func(Dialer) Dialer) ConnectionOption {
	return func(c *connectionConfig) {
		c.dialer = fn(c.dialer)
	}
}

// WithHandshaker configures the Handshaker that wll be used to initialize newly
// dialed connections.
func WithHandshaker(fn func(Handshaker) Handshaker) ConnectionOption {
	return func(c *connectionConfig) {
		c.handshaker = fn(c.handshaker)
	}
}

// WithIdleTimeout configures the maximum idle time to allow for a connection.
func WithIdleTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.idleTimeout = fn(c.idleTimeout)
	}
}

// WithReadTimeout configures the maximum read time for a connection.
func WithReadTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.readTimeout = fn(c.readTimeout)
	}
}

// WithWriteTimeout configures the maximum write time for a connection.
func WithWriteTimeout(fn func(time.Duration) time.Duration) ConnectionOption {
	return func(c *connectionConfig) {
		c.writeTimeout = fn(c.writeTimeout)
	}
}

// WithTLSConfig configures the TLS options for a connection.
func WithTLSConfig(fn func(*tls.Config) *tls.Config) ConnectionOption {
	return func(c *connectionConfig) {
		c.tlsConfig = fn(c.tlsConfig)
	}
}

// WithHTTPClient configures the HTTP client for a connection.
func WithHTTPClient(fn func(*http.Client) *http.Client) ConnectionOption {
	return func(c *connectionConfig) {
		c.httpClient = fn(c.httpClient)
	}
}

// WithMonitor configures a event for command monitoring.
func WithMonitor(fn func(*event.CommandMonitor) *event.CommandMonitor) ConnectionOption {
	return func(c *connectionConfig) {
		c.cmdMonitor = fn(c.cmdMonitor)
	}
}

// WithZlibLevel sets the zLib compression level.
func WithZlibLevel(fn func(*int) *int) ConnectionOption {
	return func(c *connectionConfig) {
		c.zlibLevel = fn(c.zlibLevel)
	}
}

// WithZstdLevel sets the zstd compression level.
func WithZstdLevel(fn func(*int) *int) ConnectionOption {
	return func(c *connectionConfig) {
		c.zstdLevel = fn(c.zstdLevel)
	}
}

// WithOCSPCache specifies a cache to use for OCSP verification.
func WithOCSPCache(fn func(ocsp.Cache) ocsp.Cache) ConnectionOption {
	return func(c *connectionConfig) {
		c.ocspCache = fn(c.ocspCache)
	}
}

// WithDisableOCSPEndpointCheck specifies whether or the driver should perform non-stapled OCSP verification. If set
// to true, the driver will only check stapled responses and will continue the connection without reaching out to
// OCSP responders.
func WithDisableOCSPEndpointCheck(fn func(bool) bool) ConnectionOption {
	return func(c *connectionConfig) {
		c.disableOCSPEndpointCheck = fn(c.disableOCSPEndpointCheck)
	}
}

// WithConnectionLoadBalanced specifies whether or not the connection is to a server behind a load balancer.
func WithConnectionLoadBalanced(fn func(bool) bool) ConnectionOption {
	return func(c *connectionConfig) {
		c.loadBalanced = fn(c.loadBalanced)
	}
}

func withGenerationNumberFn(fn func(generationNumberFn) generationNumberFn) ConnectionOption {
	return func(c *connectionConfig) {
		c.getGenerationFn = fn(c.getGenerationFn)
	}
}
