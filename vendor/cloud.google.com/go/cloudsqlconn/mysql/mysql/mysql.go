// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mysql provides a Cloud SQL MySQL driver that uses go-sql-driver/mysql
// and works with database/sql
package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net"
	"syscall"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/go-sql-driver/mysql"
)

// RegisterDriver registers a MySQL driver that uses the cloudsqlconn.Dialer
// configured with the provided options. The choice of name is entirely up to
// the caller and may be used to distinguish between multiple registrations of
// differently configured Dialers.
func RegisterDriver(name string, opts ...cloudsqlconn.Option) (func() error, error) {
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return func() error { return nil }, err
	}
	mysql.RegisterDialContext(name, mysql.DialContextFunc(func(ctx context.Context, addr string) (net.Conn, error) {
		conn, err := d.Dial(ctx, addr)
		if err != nil {
			return nil, err
		}
		return LivenessCheckConn{Conn: conn}, nil
	}))
	sql.Register(name, &mysqlDriver{
		d: &mysql.MySQLDriver{},
	})
	return func() error { return d.Close() }, nil
}

// LivenessCheckConn wraps the underlying connection with support for a
// liveness check.
//
// See https://github.com/go-sql-driver/mysql/pull/934 for details.
type LivenessCheckConn struct {
	net.Conn
}

// SyscallConn supports a connection check in the MySQL driver by delegating to
// the underlying non-TLS net.Conn.
func (c *LivenessCheckConn) SyscallConn() (syscall.RawConn, error) {
	sconn, ok := c.Conn.(syscall.Conn)
	if !ok {
		return nil, errors.New("connection is not a syscall.Conn")
	}
	return sconn.SyscallConn()
}

type mysqlDriver struct {
	d *mysql.MySQLDriver
}

// Open accepts a DSN using the go-sql-driver/mysql format. See
// https://github.com/go-sql-driver/mysql#dsn-data-source-name for details.
// Note the protocol should match the name used when registering a driver. For
// example, a connection string might look like this where "cloudsql-mysql" is
// the named used when registering the driver:
//
//	my-user:mypass@cloudsql-mysql(my-proj:us-central1:my-inst)/my-db
func (d *mysqlDriver) Open(name string) (driver.Conn, error) {
	return d.d.Open(name)
}
