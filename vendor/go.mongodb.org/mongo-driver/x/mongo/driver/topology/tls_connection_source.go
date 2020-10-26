// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"crypto/tls"
	"net"
)

type tlsConnectionSource interface {
	Client(net.Conn, *tls.Config) *tls.Conn
}

type tlsConnectionSourceFn func(net.Conn, *tls.Config) *tls.Conn

func (t tlsConnectionSourceFn) Client(nc net.Conn, cfg *tls.Config) *tls.Conn {
	return t(nc, cfg)
}

var defaultTLSConnectionSource tlsConnectionSourceFn = func(nc net.Conn, cfg *tls.Config) *tls.Conn {
	return tls.Client(nc, cfg)
}
