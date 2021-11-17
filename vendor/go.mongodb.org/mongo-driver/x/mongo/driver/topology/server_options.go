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
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var defaultRegistry = bson.NewRegistryBuilder().Build()

type serverConfig struct {
	clock                     *session.ClusterClock
	compressionOpts           []string
	connectionOpts            []ConnectionOption
	appname                   string
	heartbeatInterval         time.Duration
	heartbeatTimeout          time.Duration
	maxConns                  uint64
	minConns                  uint64
	poolMonitor               *event.PoolMonitor
	serverMonitor             *event.ServerMonitor
	connectionPoolMaxIdleTime time.Duration
	registry                  *bsoncodec.Registry
	monitoringDisabled        bool
	serverAPI                 *driver.ServerAPIOptions
	loadBalanced              bool
}

func newServerConfig(opts ...ServerOption) (*serverConfig, error) {
	cfg := &serverConfig{
		heartbeatInterval: 10 * time.Second,
		heartbeatTimeout:  10 * time.Second,
		maxConns:          100,
		registry:          defaultRegistry,
	}

	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// ServerOption configures a server.
type ServerOption func(*serverConfig) error

func withMonitoringDisabled(fn func(bool) bool) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.monitoringDisabled = fn(cfg.monitoringDisabled)
		return nil
	}
}

// WithConnectionOptions configures the server's connections.
func WithConnectionOptions(fn func(...ConnectionOption) []ConnectionOption) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.connectionOpts = fn(cfg.connectionOpts...)
		return nil
	}
}

// WithCompressionOptions configures the server's compressors.
func WithCompressionOptions(fn func(...string) []string) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.compressionOpts = fn(cfg.compressionOpts...)
		return nil
	}
}

// WithServerAppName configures the server's application name.
func WithServerAppName(fn func(string) string) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.appname = fn(cfg.appname)
		return nil
	}
}

// WithHeartbeatInterval configures a server's heartbeat interval.
func WithHeartbeatInterval(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.heartbeatInterval = fn(cfg.heartbeatInterval)
		return nil
	}
}

// WithHeartbeatTimeout configures how long to wait for a heartbeat socket to
// connection.
func WithHeartbeatTimeout(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.heartbeatTimeout = fn(cfg.heartbeatTimeout)
		return nil
	}
}

// WithMaxConnections configures the maximum number of connections to allow for
// a given server. If max is 0, then the default will be math.MaxInt64.
func WithMaxConnections(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.maxConns = fn(cfg.maxConns)
		return nil
	}
}

// WithMinConnections configures the minimum number of connections to allow for
// a given server. If min is 0, then there is no lower limit to the number of
// connections.
func WithMinConnections(fn func(uint64) uint64) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.minConns = fn(cfg.minConns)
		return nil
	}
}

// WithConnectionPoolMaxIdleTime configures the maximum time that a connection can remain idle in the connection pool
// before being removed. If connectionPoolMaxIdleTime is 0, then no idle time is set and connections will not be removed
// because of their age
func WithConnectionPoolMaxIdleTime(fn func(time.Duration) time.Duration) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.connectionPoolMaxIdleTime = fn(cfg.connectionPoolMaxIdleTime)
		return nil
	}
}

// WithConnectionPoolMonitor configures the monitor for all connection pool actions
func WithConnectionPoolMonitor(fn func(*event.PoolMonitor) *event.PoolMonitor) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.poolMonitor = fn(cfg.poolMonitor)
		return nil
	}
}

// WithServerMonitor configures the monitor for all SDAM events for a server
func WithServerMonitor(fn func(*event.ServerMonitor) *event.ServerMonitor) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.serverMonitor = fn(cfg.serverMonitor)
		return nil
	}
}

// WithClock configures the ClusterClock for the server to use.
func WithClock(fn func(clock *session.ClusterClock) *session.ClusterClock) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.clock = fn(cfg.clock)
		return nil
	}
}

// WithRegistry configures the registry for the server to use when creating
// cursors.
func WithRegistry(fn func(*bsoncodec.Registry) *bsoncodec.Registry) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.registry = fn(cfg.registry)
		return nil
	}
}

// WithServerAPI configures the server API options for the server to use.
func WithServerAPI(fn func(serverAPI *driver.ServerAPIOptions) *driver.ServerAPIOptions) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.serverAPI = fn(cfg.serverAPI)
		return nil
	}
}

// WithServerLoadBalanced specifies whether or not the server is behind a load balancer.
func WithServerLoadBalanced(fn func(bool) bool) ServerOption {
	return func(cfg *serverConfig) error {
		cfg.loadBalanced = fn(cfg.loadBalanced)
		return nil
	}
}
