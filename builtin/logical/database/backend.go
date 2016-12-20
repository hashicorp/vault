package database

import (
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/builtin/logical/database/dbs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend(conf).Setup(conf)
}

func Backend(conf *logical.BackendConfig) *databaseBackend {
	var b databaseBackend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		Paths: []*framework.Path{
			pathConfigConnection(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathRoleCreate(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean: b.resetAllDBs,
	}

	b.logger = conf.Logger
	b.connections = make(map[string]dbs.DatabaseType)
	return &b
}

type databaseBackend struct {
	connections map[string]dbs.DatabaseType
	logger      log.Logger

	*framework.Backend
	sync.RWMutex
}

// resetAllDBs closes all connections from all database types
func (b *databaseBackend) resetAllDBs() {
	b.logger.Trace("postgres/resetdb: enter")
	defer b.logger.Trace("postgres/resetdb: exit")

	b.Lock()
	defer b.Unlock()

	for _, db := range b.connections {
		db.Close()
	}
}

func (b *databaseBackend) Role(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
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
