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

// Package pgxv4 provides a Cloud SQL Postgres driver that uses pgx v4 and works
// with the database/sql package.
package pgxv4

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"net"
	"sync"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

// RegisterDriver registers a Postgres driver that uses the cloudsqlconn.Dialer
// configured with the provided options. The choice of name is entirely up to
// the caller and may be used to distinguish between multiple registrations of
// differently configured Dialers. The driver uses pgx/v4 internally.
// RegisterDriver returns a cleanup function that should be called one the
// database connection is no longer needed.
func RegisterDriver(name string, opts ...cloudsqlconn.Option) (func() error, error) {
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return func() error { return nil }, err
	}
	sql.Register(name, &pgDriver{
		d:      d,
		dbURIs: make(map[string]string),
	})
	return func() error { return d.Close() }, nil
}

type pgDriver struct {
	d  *cloudsqlconn.Dialer
	mu sync.RWMutex
	// dbURIs is a map of DSN to DB URI for registered connection names.
	dbURIs map[string]string
}

// Open accepts a keyword/value formatted connection string and returns a
// connection to the database using cloudsqlconn.Dialer. The Cloud SQL instance
// connection name should be specified in the host field. For example:
//
// "host=my-project:us-central1:my-db-instance user=myuser password=mypass"
func (p *pgDriver) Open(name string) (driver.Conn, error) {
	dbURI, err := p.dbURI(name)
	if err != nil {
		return nil, err
	}
	return stdlib.GetDefaultDriver().Open(dbURI)

}

// dbURI registers a driver using the provided DSN. If the name has already
// been registered, dbURI returns the existing registration.
func (p *pgDriver) dbURI(name string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	dbURI, ok := p.dbURIs[name]
	if ok {
		return dbURI, nil
	}

	config, err := pgx.ParseConfig(name)
	if err != nil {
		return "", err
	}
	instConnName := config.Config.Host // Extract instance connection name
	config.Config.Host = "localhost"   // Replace it with a default value
	config.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
		return p.d.Dial(ctx, instConnName)
	}

	dbURI = stdlib.RegisterConnConfig(config)
	p.dbURIs[name] = dbURI

	return dbURI, nil
}
