package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/connection",
			},
		},

		Paths: []*framework.Path{
			pathConfigConnection(&b),
			pathConfigLease(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Invalidate:  b.invalidate,
		Clean:       b.ResetDB,
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend

	db        *sql.DB
	defaultDb string
	lock      sync.Mutex
}

// DB returns the default database connection.
func (b *backend) DB(ctx context.Context, s logical.Storage) (*sql.DB, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// If we already have a DB, we got it!
	if b.db != nil {
		if err := b.db.Ping(); err == nil {
			return b.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		b.db.Close()
	}

	// Otherwise, attempt to make connection
	entry, err := s.Get(ctx, "config/connection")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("configure the DB connection with config/connection first")
	}

	var connConfig connectionConfig
	if err := entry.DecodeJSON(&connConfig); err != nil {
		return nil, err
	}
	connString := connConfig.ConnectionString

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	db.SetMaxOpenConns(connConfig.MaxOpenConnections)

	stmt, err := db.Prepare("SELECT db_name();")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&b.defaultDb)
	if err != nil {
		return nil, err
	}

	b.db = db
	return b.db, nil
}

// ResetDB forces a connection next time DB() is called.
func (b *backend) ResetDB(_ context.Context) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.db != nil {
		b.db.Close()
	}

	b.db = nil
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config/connection":
		b.ResetDB(ctx)
	}
}

// LeaseConfig returns the lease configuration
func (b *backend) LeaseConfig(ctx context.Context, s logical.Storage) (*configLease, error) {
	entry, err := s.Get(ctx, "config/lease")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result configLease
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

const backendHelp = `
The MSSQL backend dynamically generates database users.

After mounting this backend, configure it using the endpoints within
the "config/" path.

This backend does not support Azure SQL Databases.
`
