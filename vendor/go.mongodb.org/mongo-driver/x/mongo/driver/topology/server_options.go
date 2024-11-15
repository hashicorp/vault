// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var defaultRegistry = bson.NewRegistryBuilder().Build()

type serverConfig struct {
	clock                *session.ClusterClock
	compressionOpts      []string
	connectionOpts       []ConnectionOption
	appname              string
	heartbeatInterval    time.Duration
	heartbeatTimeout     time.Duration
	serverMonitoringMode string
	serverMonitor        *event.ServerMonitor
	registry             *bsoncodec.Registry
	monitoringDisabled   bool
	serverAPI            *driver.ServerAPIOptions
	loadBalanced         bool

	// Connection pool options.
	maxConns             uint64
	minConns             uint64
	maxConnecting        uint64
	poolMonitor          *event.PoolMonitor
	logger               *logger.Logger
	poolMaxIdleTime      time.Duration
	poolMaintainInterval time.Duration
}

func newServerConfig(opts ...ServerOption) *serverConfig {
	cfg := &serverConfig{
		heartbeatInterval: 10 * time.Second,
		heartbeatTimeout:  10 * time.Second,
		registry:          defaultRegistry,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cfg)
	}

	return cfg
}

// ServerOption configures a server.
type ServerOption func(*serverConfig)

// ServerAPIFromServerOptions will return the server API options if they have been functionally set on the ServerOption
// slice.
func ServerAPIFromServerOptions(opts []ServerOption) *driver.ServerAPIOptions {
	return newServerConfig(opts...).serverAPI
}

func withMonitoringDisabled(fn func(bool) bool) ServerOption {
	return func(cfg *serverConfig) {
		cfg.monitoringDisabled = fn(cfg.monitoringDisabled)
	}
}

// WithConnectionOptions configures the server's connections.
func WithConnectionOptions(fn func(...ConnectionOption) []ConnectionOption) ServerOption {
	return func(cfg *serverConfig) {
		cfg.connectionOpts = fn(cfg.connectionOpts...)
	}
}

// WithCompressionOptions configures the server's compressors.
func WithCompressionOptions(fn func(...string) []string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.compressionOpts = fn(cfg.compressionOpts...)
	}
}

// WithServerAppName configures the server's application name.
func WithServerAppName(fn func(string) string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.appname = fn(cfg.appname)
	}
}

// WithHeartbeatInterval configures a server's heartbeat interval.
func WithHeartbeatInterval(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.heartbeatInterval = fn(cfg.heartbeatInterval)
	}
}

// WithHeartbeatTimeout configures how long to wait for a heartbeat socket to
// connection.
func WithHeartbeatTimeout(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.heartbeatTimeout = fn(cfg.heartbeatTimeout)
	}
}

// WithMaxConnections configures the maximum number of connections to allow for
// a given server. If max is 0, then maximum connection pool size is not limited.
func WithMaxConnections(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) {
		cfg.maxConns = fn(cfg.maxConns)
	}
}

// WithMinConnections configures the minimum number of connections to allow for
// a given server. If min is 0, then there is no lower limit to the number of
// connections.
func WithMinConnections(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) {
		cfg.minConns = fn(cfg.minConns)
	}
}

// WithMaxConnecting configures the maximum number of connections a connection
// pool may establish simultaneously. If maxConnecting is 0, the default value
// of 2 is used.
func WithMaxConnecting(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) {
		cfg.maxConnecting = fn(cfg.maxConnecting)
	}
}

// WithConnectionPoolMaxIdleTime configures the maximum time that a connection can remain idle in the connection pool
// before being removed. If connectionPoolMaxIdleTime is 0, then no idle time is set and connections will not be removed
// because of their age
func WithConnectionPoolMaxIdleTime(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.poolMaxIdleTime = fn(cfg.poolMaxIdleTime)
	}
}

// WithConnectionPoolMaintainInterval configures the interval that the background connection pool
// maintenance goroutine runs.
func WithConnectionPoolMaintainInterval(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.poolMaintainInterval = fn(cfg.poolMaintainInterval)
	}
}

// WithConnectionPoolMonitor configures the monitor for all connection pool actions
func WithConnectionPoolMonitor(fn func(*event.PoolMonitor) *event.PoolMonitor) ServerOption {
	return func(cfg *serverConfig) {
		cfg.poolMonitor = fn(cfg.poolMonitor)
	}
}

// WithServerMonitor configures the monitor for all SDAM events for a server
func WithServerMonitor(fn func(*event.ServerMonitor) *event.ServerMonitor) ServerOption {
	return func(cfg *serverConfig) {
		cfg.serverMonitor = fn(cfg.serverMonitor)
	}
}

// WithClock configures the ClusterClock for the server to use.
func WithClock(fn func(clock *session.ClusterClock) *session.ClusterClock) ServerOption {
	return func(cfg *serverConfig) {
		cfg.clock = fn(cfg.clock)
	}
}

// WithRegistry configures the registry for the server to use when creating
// cursors.
func WithRegistry(fn func(*bsoncodec.Registry) *bsoncodec.Registry) ServerOption {
	return func(cfg *serverConfig) {
		cfg.registry = fn(cfg.registry)
	}
}

// WithServerAPI configures the server API options for the server to use.
func WithServerAPI(fn func(serverAPI *driver.ServerAPIOptions) *driver.ServerAPIOptions) ServerOption {
	return func(cfg *serverConfig) {
		cfg.serverAPI = fn(cfg.serverAPI)
	}
}

// WithServerLoadBalanced specifies whether or not the server is behind a load balancer.
func WithServerLoadBalanced(fn func(bool) bool) ServerOption {
	return func(cfg *serverConfig) {
		cfg.loadBalanced = fn(cfg.loadBalanced)
	}
}

// withLogger configures the logger for the server to use.
func withLogger(fn func() *logger.Logger) ServerOption {
	return func(cfg *serverConfig) {
		cfg.logger = fn()
	}
}

// withServerMonitoringMode configures the mode (stream, poll, or auto) to use
// for monitoring.
func withServerMonitoringMode(mode *string) ServerOption {
	return func(cfg *serverConfig) {
		if mode != nil {
			cfg.serverMonitoringMode = *mode

			return
		}

		cfg.serverMonitoringMode = connstring.ServerMonitoringModeAuto
	}
}
