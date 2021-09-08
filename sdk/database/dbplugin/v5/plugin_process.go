package dbplugin

import (
	"fmt"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-plugin"
)

type pluginProcess struct {
	mu      *sync.RWMutex
	stopped bool

	logger     log.Logger
	process    process
	pluginType string

	dbs map[string]*multiplexedDatabase // ID -> Database
}

type process interface {
	Client() (plugin.ClientProtocol, error)
	Kill()
}

func newPluginProcess(logger log.Logger, process *plugin.Client) (*pluginProcess, error) {
	rpcClient, err := process.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	rawDB, err := rpcClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	db, ok := rawDB.(gRPCClient)
	if !ok {
		return nil, fmt.Errorf("unsupported client type")
	}

	pluginType, err := db.Type()
	if err != nil {
		return nil, err
	}

	p := &pluginProcess{
		mu:         new(sync.RWMutex),
		stopped:    false,
		logger:     logger,
		process:    process,
		pluginType: pluginType,
		dbs:        make(map[string]*multiplexedDatabase),
	}

	return p, nil
}

// GetDatabase retrieves an existing Database instance if that Database already exists
// If the Database does not exist, this will create a new instance and add it to its
// cache for fast retrieval later.
func (p *pluginProcess) getDatabase(id string) (*multiplexedDatabase, error) {
	p.mu.RLock()
	multiplexedDB, exists := p.dbs[id]
	if exists {
		p.mu.RUnlock()
		return multiplexedDB, nil
	}

	// Upgrade the lock to a write lock
	p.mu.RUnlock()
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check to see if another goroutine has added the DB to the cache
	multiplexedDB, exists = p.dbs[id]
	if exists {
		return multiplexedDB, nil
	}

	// DB doesn't exist in the cache - create one, add it to the cache and return
	rpcClient, err := p.process.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin client: %w", err)
	}

	rawDB, err := rpcClient.Dispense("database")
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin client database: %w", err)
	}

	grpcClient, ok := rawDB.(*gRPCClient)
	if !ok {
		return nil, fmt.Errorf("plugin is not a gRPCClient")
	}
	// Set the ID on the client since we can't do this as a part of the Dispense operation
	// because go-plugin is what creates the client instance
	grpcClient.id = id

	namedLogger := p.logger.Named(id)

	db := addMiddleware(grpcClient, p.pluginType, namedLogger)

	multiplexedDB = &multiplexedDatabase{
		Database: db,
		id:       id,
		closer: func() error {
			return p.stopDatabase(id)
		},
	}
	p.dbs[id] = multiplexedDB

	return multiplexedDB, nil
}

var _ Database = (*multiplexedDatabase)(nil)

type multiplexedDatabase struct {
	Database
	id     string
	closer func() error
}

func (m multiplexedDatabase) Close() error {
	dbErr := m.Database.Close()
	err := m.closer()
	merr := multierror.Append(dbErr, err)
	return merr.ErrorOrNil()
}

func (p *pluginProcess) stopDatabase(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stopped {
		return fmt.Errorf("plugin process has already been stopped")
	}

	db, exists := p.dbs[id]
	if !exists {
		return fmt.Errorf("database [%s] does not exist in this plugin process", id)
	}

	dbErr := db.Close()
	delete(p.dbs, id)

	if len(p.dbs) > 0 {
		// Other DBs running with this process - keep the plugin running
		return dbErr
	}

	// No other DBs are running on this process, kill the plugin
	p.process.Kill()
	return dbErr
}

// Close the entire process down
func (p *pluginProcess) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	merr := new(multierror.Error)
	for id := range p.dbs {
		err := p.stopDatabase(id)
		merr = multierror.Append(merr, err)
	}
	p.stopped = true
	return merr.ErrorOrNil()
}

func addMiddleware(db Database, pluginType string, logger log.Logger) Database {
	db = &databaseMetricsMiddleware{
		next:    db,
		typeStr: pluginType,
	}

	// Wrap with tracing middleware
	if logger.IsTrace() {
		db = &databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
	}

	return db
}
