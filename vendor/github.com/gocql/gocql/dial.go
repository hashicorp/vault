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
 * Copyright (c) 2016, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
)

// HostDialer allows customizing connection to cluster nodes.
type HostDialer interface {
	// DialHost establishes a connection to the host.
	// The returned connection must be directly usable for CQL protocol,
	// specifically DialHost is responsible also for setting up the TLS session if needed.
	// DialHost should disable write coalescing if the returned net.Conn does not support writev.
	// As of Go 1.18, only plain TCP connections support writev, TLS sessions should disable coalescing.
	// You can use WrapTLS helper function if you don't need to override the TLS setup.
	DialHost(ctx context.Context, host *HostInfo) (*DialedHost, error)
}

// DialedHost contains information about established connection to a host.
type DialedHost struct {
	// Conn used to communicate with the server.
	Conn net.Conn

	// DisableCoalesce disables write coalescing for the Conn.
	// If true, the effect is the same as if WriteCoalesceWaitTime was configured to 0.
	DisableCoalesce bool
}

// defaultHostDialer dials host in a default way.
type defaultHostDialer struct {
	dialer    Dialer
	tlsConfig *tls.Config
}

func (hd *defaultHostDialer) DialHost(ctx context.Context, host *HostInfo) (*DialedHost, error) {
	ip := host.ConnectAddress()
	port := host.Port()

	if !validIpAddr(ip) {
		return nil, fmt.Errorf("host missing connect ip address: %v", ip)
	} else if port == 0 {
		return nil, fmt.Errorf("host missing port: %v", port)
	}

	connAddr := host.ConnectAddressAndPort()
	conn, err := hd.dialer.DialContext(ctx, "tcp", connAddr)
	if err != nil {
		return nil, err
	}
	addr := host.HostnameAndPort()
	return WrapTLS(ctx, conn, addr, hd.tlsConfig)
}

func tlsConfigForAddr(tlsConfig *tls.Config, addr string) *tls.Config {
	// the TLS config is safe to be reused by connections but it must not
	// be modified after being used.
	if !tlsConfig.InsecureSkipVerify && tlsConfig.ServerName == "" {
		colonPos := strings.LastIndex(addr, ":")
		if colonPos == -1 {
			colonPos = len(addr)
		}
		hostname := addr[:colonPos]
		// clone config to avoid modifying the shared one.
		tlsConfig = tlsConfig.Clone()
		tlsConfig.ServerName = hostname
	}
	return tlsConfig
}

// WrapTLS optionally wraps a net.Conn connected to addr with the given tlsConfig.
// If the tlsConfig is nil, conn is not wrapped into a TLS session, so is insecure.
// If the tlsConfig does not have server name set, it is updated based on the default gocql rules.
func WrapTLS(ctx context.Context, conn net.Conn, addr string, tlsConfig *tls.Config) (*DialedHost, error) {
	if tlsConfig != nil {
		tlsConfig := tlsConfigForAddr(tlsConfig, addr)
		tconn := tls.Client(conn, tlsConfig)
		if err := tconn.HandshakeContext(ctx); err != nil {
			conn.Close()
			return nil, err
		}
		conn = tconn
	}

	return &DialedHost{
		Conn:            conn,
		DisableCoalesce: tlsConfig != nil, // write coalescing can't use writev when the connection is wrapped.
	}, nil
}
