// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build go1.17
// +build go1.17

package topology

import (
	"context"
	"crypto/tls"
	"net"
)

type tlsConn interface {
	net.Conn

	// Require HandshakeContext on the interface for Go 1.17 and higher.
	HandshakeContext(ctx context.Context) error
	ConnectionState() tls.ConnectionState
}

var _ tlsConn = (*tls.Conn)(nil)

type tlsConnectionSource interface {
	Client(net.Conn, *tls.Config) tlsConn
}

type tlsConnectionSourceFn func(net.Conn, *tls.Config) tlsConn

var _ tlsConnectionSource = (tlsConnectionSourceFn)(nil)

func (t tlsConnectionSourceFn) Client(nc net.Conn, cfg *tls.Config) tlsConn {
	return t(nc, cfg)
}

var defaultTLSConnectionSource tlsConnectionSourceFn = func(nc net.Conn, cfg *tls.Config) tlsConn {
	return tls.Client(nc, cfg)
}

// clientHandshake will perform a handshake on Go 1.17 and higher with HandshakeContext.
func clientHandshake(ctx context.Context, client tlsConn) error {
	return client.HandshakeContext(ctx)
}
