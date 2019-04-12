// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2013 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

// +build appengine

package mysql

import (
	"context"
	"net"

	"google.golang.org/appengine/cloudsql"
)

func init() {
	RegisterDialContext("cloudsql", func(_ context.Context, instance string) (net.Conn, error) {
		// XXX: the cloudsql driver still does not export a Context-aware dialer.
		return cloudsql.Dial(instance)
	})
}
