package postgresql

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(conf *logical.BackendConfig) *backend {
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
			pathRoleCreate(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean:       b.ResetDB,
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	b.logger = conf.Logger
	return &b
}

type backend struct {
	*framework.Backend

	db   *sql.DB
	lock sync.Mutex

	logger log.Logger
}

// DB returns the database connection.
func (b *backend) DB(s logical.Storage) (*sql.DB, error) {
	b.logger.Trace("postgres/db: enter")
	defer b.logger.Trace("postgres/db: exit")

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
	entry, err := s.Get("config/connection")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil,
			fmt.Errorf("configure the DB connection with config/connection first")
	}

	var connConfig connectionConfig
	if err := entry.DecodeJSON(&connConfig); err != nil {
		return nil, err
	}

	conn := connConfig.ConnectionURL
	if len(conn) == 0 {
		conn = connConfig.ConnectionString
	}

	// Ensure timezone is set to UTC for all the conenctions
	if strings.HasPrefix(conn, "postgres://") || strings.HasPrefix(conn, "postgresql://") {
		if strings.Contains(conn, "?") {
			conn += "&timezone=utc"
		} else {
			conn += "?timezone=utc"
		}
	} else {
		conn += " timezone=utc"
	}

	b.db, err = sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	b.db.SetMaxOpenConns(connConfig.MaxOpenConnections)
	b.db.SetMaxIdleConns(connConfig.MaxIdleConnections)

	return b.db, nil
}

// ResetDB forces a connection next time DB() is called.
func (b *backend) ResetDB() {
	b.logger.Trace("postgres/resetdb: enter")
	defer b.logger.Trace("postgres/resetdb: exit")

	b.lock.Lock()
	defer b.lock.Unlock()

	if b.db != nil {
		b.db.Close()
	}

	b.db = nil
}

func (b *backend) invalidate(key string) {
	switch key {
	case "config/connection":
		b.ResetDB()
	}
}

// Lease returns the lease information
func (b *backend) Lease(s logical.Storage) (*configLease, error) {
	entry, err := s.Get("config/lease")
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
The PostgreSQL backend dynamically generates database users.

After mounting this backend, configure it using the endpoints within
the "config/" path.
`
