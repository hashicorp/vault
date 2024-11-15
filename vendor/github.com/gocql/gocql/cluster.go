/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2012, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"context"
	"errors"
	"net"
	"time"
)

// PoolConfig configures the connection pool used by the driver, it defaults to
// using a round-robin host selection policy and a round-robin connection selection
// policy for each host.
type PoolConfig struct {
	// HostSelectionPolicy sets the policy for selecting which host to use for a
	// given query (default: RoundRobinHostPolicy())
	// It is not supported to use a single HostSelectionPolicy in multiple sessions
	// (even if you close the old session before using in a new session).
	HostSelectionPolicy HostSelectionPolicy
}

func (p PoolConfig) buildPool(session *Session) *policyConnPool {
	return newPolicyConnPool(session)
}

// ClusterConfig is a struct to configure the default cluster implementation
// of gocql. It has a variety of attributes that can be used to modify the
// behavior to fit the most common use cases. Applications that require a
// different setup must implement their own cluster.
type ClusterConfig struct {
	// addresses for the initial connections. It is recommended to use the value set in
	// the Cassandra config for broadcast_address or listen_address, an IP address not
	// a domain name. This is because events from Cassandra will use the configured IP
	// address, which is used to index connected hosts. If the domain name specified
	// resolves to more than 1 IP address then the driver may connect multiple times to
	// the same host, and will not mark the node being down or up from events.
	Hosts []string

	// CQL version (default: 3.0.0)
	CQLVersion string

	// ProtoVersion sets the version of the native protocol to use, this will
	// enable features in the driver for specific protocol versions, generally this
	// should be set to a known version (2,3,4) for the cluster being connected to.
	//
	// If it is 0 or unset (the default) then the driver will attempt to discover the
	// highest supported protocol for the cluster. In clusters with nodes of different
	// versions the protocol selected is not defined (ie, it can be any of the supported in the cluster)
	ProtoVersion int

	// Timeout limits the time spent on the client side while executing a query.
	// Specifically, query or batch execution will return an error if the client does not receive a response
	// from the server within the Timeout period.
	// Timeout is also used to configure the read timeout on the underlying network connection.
	// Client Timeout should always be higher than the request timeouts configured on the server,
	// so that retries don't overload the server.
	// Timeout has a default value of 11 seconds, which is higher than default server timeout for most query types.
	// Timeout is not applied to requests during initial connection setup, see ConnectTimeout.
	Timeout time.Duration

	// ConnectTimeout limits the time spent during connection setup.
	// During initial connection setup, internal queries, AUTH requests will return an error if the client
	// does not receive a response within the ConnectTimeout period.
	// ConnectTimeout is applied to the connection setup queries independently.
	// ConnectTimeout also limits the duration of dialing a new TCP connection
	// in case there is no Dialer nor HostDialer configured.
	// ConnectTimeout has a default value of 11 seconds.
	ConnectTimeout time.Duration

	// WriteTimeout limits the time the driver waits to write a request to a network connection.
	// WriteTimeout should be lower than or equal to Timeout.
	// WriteTimeout defaults to the value of Timeout.
	WriteTimeout time.Duration

	// Port used when dialing.
	// Default: 9042
	Port int

	// Initial keyspace. Optional.
	Keyspace string

	// Number of connections per host.
	// Default: 2
	NumConns int

	// Default consistency level.
	// Default: Quorum
	Consistency Consistency

	// Compression algorithm.
	// Default: nil
	Compressor Compressor

	// Default: nil
	Authenticator Authenticator

	// An Authenticator factory. Can be used to create alternative authenticators.
	// Default: nil
	AuthProvider func(h *HostInfo) (Authenticator, error)

	// Default retry policy to use for queries.
	// Default: no retries.
	RetryPolicy RetryPolicy

	// ConvictionPolicy decides whether to mark host as down based on the error and host info.
	// Default: SimpleConvictionPolicy
	ConvictionPolicy ConvictionPolicy

	// Default reconnection policy to use for reconnecting before trying to mark host as down.
	ReconnectionPolicy ReconnectionPolicy

	// The keepalive period to use, enabled if > 0 (default: 0)
	// SocketKeepalive is used to set up the default dialer and is ignored if Dialer or HostDialer is provided.
	SocketKeepalive time.Duration

	// Maximum cache size for prepared statements globally for gocql.
	// Default: 1000
	MaxPreparedStmts int

	// Maximum cache size for query info about statements for each session.
	// Default: 1000
	MaxRoutingKeyInfo int

	// Default page size to use for created sessions.
	// Default: 5000
	PageSize int

	// Consistency for the serial part of queries, values can be either SERIAL or LOCAL_SERIAL.
	// Default: unset
	SerialConsistency SerialConsistency

	// SslOpts configures TLS use when HostDialer is not set.
	// SslOpts is ignored if HostDialer is set.
	SslOpts *SslOptions

	// Sends a client side timestamp for all requests which overrides the timestamp at which it arrives at the server.
	// Default: true, only enabled for protocol 3 and above.
	DefaultTimestamp bool

	// PoolConfig configures the underlying connection pool, allowing the
	// configuration of host selection and connection selection policies.
	PoolConfig PoolConfig

	// If not zero, gocql attempt to reconnect known DOWN nodes in every ReconnectInterval.
	ReconnectInterval time.Duration

	// The maximum amount of time to wait for schema agreement in a cluster after
	// receiving a schema change frame. (default: 60s)
	MaxWaitSchemaAgreement time.Duration

	// HostFilter will filter all incoming events for host, any which don't pass
	// the filter will be ignored. If set will take precedence over any options set
	// via Discovery
	HostFilter HostFilter

	// AddressTranslator will translate addresses found on peer discovery and/or
	// node change events.
	AddressTranslator AddressTranslator

	// If IgnorePeerAddr is true and the address in system.peers does not match
	// the supplied host by either initial hosts or discovered via events then the
	// host will be replaced with the supplied address.
	//
	// For example if an event comes in with host=10.0.0.1 but when looking up that
	// address in system.local or system.peers returns 127.0.0.1, the peer will be
	// set to 10.0.0.1 which is what will be used to connect to.
	IgnorePeerAddr bool

	// If DisableInitialHostLookup then the driver will not attempt to get host info
	// from the system.peers table, this will mean that the driver will connect to
	// hosts supplied and will not attempt to lookup the hosts information, this will
	// mean that data_centre, rack and token information will not be available and as
	// such host filtering and token aware query routing will not be available.
	DisableInitialHostLookup bool

	// Configure events the driver will register for
	Events struct {
		// disable registering for status events (node up/down)
		DisableNodeStatusEvents bool
		// disable registering for topology events (node added/removed/moved)
		DisableTopologyEvents bool
		// disable registering for schema events (keyspace/table/function removed/created/updated)
		DisableSchemaEvents bool
	}

	// DisableSkipMetadata will override the internal result metadata cache so that the driver does not
	// send skip_metadata for queries, this means that the result will always contain
	// the metadata to parse the rows and will not reuse the metadata from the prepared
	// statement.
	//
	// See https://issues.apache.org/jira/browse/CASSANDRA-10786
	DisableSkipMetadata bool

	// QueryObserver will set the provided query observer on all queries created from this session.
	// Use it to collect metrics / stats from queries by providing an implementation of QueryObserver.
	QueryObserver QueryObserver

	// BatchObserver will set the provided batch observer on all queries created from this session.
	// Use it to collect metrics / stats from batch queries by providing an implementation of BatchObserver.
	BatchObserver BatchObserver

	// ConnectObserver will set the provided connect observer on all queries
	// created from this session.
	ConnectObserver ConnectObserver

	// FrameHeaderObserver will set the provided frame header observer on all frames' headers created from this session.
	// Use it to collect metrics / stats from frames by providing an implementation of FrameHeaderObserver.
	FrameHeaderObserver FrameHeaderObserver

	// StreamObserver will be notified of stream state changes.
	// This can be used to track in-flight protocol requests and responses.
	StreamObserver StreamObserver

	// Default idempotence for queries
	DefaultIdempotence bool

	// The time to wait for frames before flushing the frames connection to Cassandra.
	// Can help reduce syscall overhead by making less calls to write. Set to 0 to
	// disable.
	//
	// (default: 200 microseconds)
	WriteCoalesceWaitTime time.Duration

	// Dialer will be used to establish all connections created for this Cluster.
	// If not provided, a default dialer configured with ConnectTimeout will be used.
	// Dialer is ignored if HostDialer is provided.
	Dialer Dialer

	// HostDialer will be used to establish all connections for this Cluster.
	// If not provided, Dialer will be used instead.
	HostDialer HostDialer

	// Logger for this ClusterConfig.
	// If not specified, defaults to the global gocql.Logger.
	Logger StdLogger

	// internal config for testing
	disableControlConn bool
}

type Dialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

// NewCluster generates a new config for the default cluster implementation.
//
// The supplied hosts are used to initially connect to the cluster then the rest of
// the ring will be automatically discovered. It is recommended to use the value set in
// the Cassandra config for broadcast_address or listen_address, an IP address not
// a domain name. This is because events from Cassandra will use the configured IP
// address, which is used to index connected hosts. If the domain name specified
// resolves to more than 1 IP address then the driver may connect multiple times to
// the same host, and will not mark the node being down or up from events.
func NewCluster(hosts ...string) *ClusterConfig {
	cfg := &ClusterConfig{
		Hosts:                  hosts,
		CQLVersion:             "3.0.0",
		Timeout:                11 * time.Second,
		ConnectTimeout:         11 * time.Second,
		Port:                   9042,
		NumConns:               2,
		Consistency:            Quorum,
		MaxPreparedStmts:       defaultMaxPreparedStmts,
		MaxRoutingKeyInfo:      1000,
		PageSize:               5000,
		DefaultTimestamp:       true,
		MaxWaitSchemaAgreement: 60 * time.Second,
		ReconnectInterval:      60 * time.Second,
		ConvictionPolicy:       &SimpleConvictionPolicy{},
		ReconnectionPolicy:     &ConstantReconnectionPolicy{MaxRetries: 3, Interval: 1 * time.Second},
		WriteCoalesceWaitTime:  200 * time.Microsecond,
	}
	return cfg
}

func (cfg *ClusterConfig) logger() StdLogger {
	if cfg.Logger == nil {
		return Logger
	}
	return cfg.Logger
}

// CreateSession initializes the cluster based on this config and returns a
// session object that can be used to interact with the database.
func (cfg *ClusterConfig) CreateSession() (*Session, error) {
	return NewSession(*cfg)
}

// translateAddressPort is a helper method that will use the given AddressTranslator
// if defined, to translate the given address and port into a possibly new address
// and port, If no AddressTranslator or if an error occurs, the given address and
// port will be returned.
func (cfg *ClusterConfig) translateAddressPort(addr net.IP, port int) (net.IP, int) {
	if cfg.AddressTranslator == nil || len(addr) == 0 {
		return addr, port
	}
	newAddr, newPort := cfg.AddressTranslator.Translate(addr, port)
	if gocqlDebug {
		cfg.logger().Printf("gocql: translating address '%v:%d' to '%v:%d'", addr, port, newAddr, newPort)
	}
	return newAddr, newPort
}

func (cfg *ClusterConfig) filterHost(host *HostInfo) bool {
	return !(cfg.HostFilter == nil || cfg.HostFilter.Accept(host))
}

var (
	ErrNoHosts              = errors.New("no hosts provided")
	ErrNoConnectionsStarted = errors.New("no connections were made when creating the session")
	ErrHostQueryFailed      = errors.New("unable to populate Hosts")
)
