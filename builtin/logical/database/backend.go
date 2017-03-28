package database

import (
	"fmt"
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

		PathsSpecial: &logical.Paths{
			Root: []string{
				"dbs/plugin/*",
			},
		},

		Paths: []*framework.Path{
			pathConfigureBuiltinConnection(&b),
			pathConfigurePluginConnection(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathRoleCreate(&b),
			pathResetConnection(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean: b.closeAllDBs,
	}

	b.logger = conf.Logger
	b.connections = make(map[string]dbs.DatabaseType)
	return &b
}

type databaseBackend struct {
	connections map[string]dbs.DatabaseType
	logger      log.Logger

	*framework.Backend
	sync.Mutex
}

// resetAllDBs closes all connections from all database types
func (b *databaseBackend) closeAllDBs() {
	b.Lock()
	defer b.Unlock()

	for _, db := range b.connections {
		db.Close()
	}
}

// This function is used to retrieve a database object either from the cached
// connection map or by using the database config in storage. The caller of this
// function needs to hold the backend's lock.
func (b *databaseBackend) getOrCreateDBObj(s logical.Storage, name string) (dbs.DatabaseType, error) {
	// if the object already is built and cached, return it
	db, ok := b.connections[name]
	if ok {
		return db, nil
	}

	entry, err := s.Get(fmt.Sprintf("dbs/%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to read connection configuration with name: %s", name)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to find entry for connection with name: %s", name)
	}

	var config dbs.DatabaseConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}

	factory := config.GetFactory()

	db, err = factory(&config, b.System(), b.logger)
	if err != nil {
		return nil, err
	}

	err = db.Initialize(config.ConnectionDetails)
	if err != nil {
		return nil, err
	}

	b.connections[name] = db

	return db, nil
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
