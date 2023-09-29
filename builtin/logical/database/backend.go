// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/logical/database/schedule"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/syncmap"
	"github.com/hashicorp/vault/internalshared/configutil"
	v4 "github.com/hashicorp/vault/sdk/database/dbplugin"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	operationPrefixDatabase = "database"
	databaseConfigPath      = "config/"
	databaseRolePath        = "role/"
	databaseStaticRolePath  = "static-role/"
	minRootCredRollbackAge  = 1 * time.Minute
)

type dbPluginInstance struct {
	sync.RWMutex
	database databaseVersionWrapper

	id     string
	name   string
	closed bool
}

func (dbi *dbPluginInstance) ID() string {
	return dbi.id
}

func (dbi *dbPluginInstance) Close() error {
	dbi.Lock()
	defer dbi.Unlock()

	if dbi.closed {
		return nil
	}
	dbi.closed = true

	return dbi.database.Close()
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	b.credRotationQueue = queue.New()
	// Load queue and kickoff new periodic ticker
	go b.initQueue(b.queueCtx, conf)

	// collect metrics on number of plugin instances
	var err error
	b.gaugeCollectionProcess, err = metricsutil.NewGaugeCollectionProcess(
		[]string{"secrets", "database", "backend", "pluginInstances", "count"},
		[]metricsutil.Label{},
		b.collectPluginInstanceGaugeValues,
		metrics.Default(),
		configutil.UsageGaugeDefaultPeriod, // TODO: add config settings for these, or add plumbing to the main config settings
		configutil.MaximumGaugeCardinalityDefault,
		b.logger)
	if err != nil {
		return nil, err
	}
	go b.gaugeCollectionProcess.Run()
	return b, nil
}

func Backend(conf *logical.BackendConfig) *databaseBackend {
	var b databaseBackend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config/*",
				"static-role/*",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathListPluginConnection(&b),
				pathConfigurePluginConnection(&b),
				pathResetConnection(&b),
			},
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathRotateRootCredentials(&b),
		),

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},
		Clean:             b.clean,
		Invalidate:        b.invalidate,
		WALRollback:       b.walRollback,
		WALRollbackMinAge: minRootCredRollbackAge,
		BackendType:       logical.TypeLogical,
	}

	b.logger = conf.Logger
	b.connections = syncmap.NewSyncMap[string, *dbPluginInstance]()
	b.queueCtx, b.cancelQueueCtx = context.WithCancel(context.Background())
	b.roleLocks = locksutil.CreateLocks()
	b.schedule = &schedule.DefaultSchedule{}

	return &b
}

func (b *databaseBackend) collectPluginInstanceGaugeValues(context.Context) ([]metricsutil.GaugeLabelValues, error) {
	// copy the map so we can release the lock
	connectionsCopy := b.connections.Values()
	counts := map[string]int{}
	for _, v := range connectionsCopy {
		dbType, err := v.database.Type()
		if err != nil {
			// there's a chance this will already be closed since we don't hold the lock
			continue
		}
		if _, ok := counts[dbType]; !ok {
			counts[dbType] = 0
		}
		counts[dbType] += 1
	}
	var gauges []metricsutil.GaugeLabelValues
	for k, v := range counts {
		gauges = append(gauges, metricsutil.GaugeLabelValues{Labels: []metricsutil.Label{{Name: "dbType", Value: k}}, Value: float32(v)})
	}
	return gauges, nil
}

type databaseBackend struct {
	// connections holds configured database connections by config name
	connections *syncmap.SyncMap[string, *dbPluginInstance]
	logger      log.Logger

	*framework.Backend
	// credRotationQueue is an in-memory priority queue used to track Static Roles
	// that require periodic rotation. Backends will have a PriorityQueue
	// initialized on setup, but only backends that are mounted by a primary
	// server or mounted as a local mount will perform the rotations.
	credRotationQueue *queue.PriorityQueue
	// queueCtx is the context for the priority queue
	queueCtx context.Context
	// cancelQueueCtx is used to terminate the background ticker
	cancelQueueCtx context.CancelFunc

	// roleLocks is used to lock modifications to roles in the queue, to ensure
	// concurrent requests are not modifying the same role and possibly causing
	// issues with the priority queue.
	roleLocks []*locksutil.LockEntry

	// the running gauge collection process
	gaugeCollectionProcess     *metricsutil.GaugeCollectionProcess
	gaugeCollectionProcessStop sync.Once

	schedule schedule.Scheduler
}

func (b *databaseBackend) DatabaseConfig(ctx context.Context, s logical.Storage, name string) (*DatabaseConfig, error) {
	entry, err := s.Get(ctx, fmt.Sprintf("config/%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to read connection configuration: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to find entry for connection with name: %q", name)
	}

	var config DatabaseConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

type upgradeStatements struct {
	// This json tag has a typo in it, the new version does not. This
	// necessitates this upgrade logic.
	CreationStatements   string `json:"creation_statments"`
	RevocationStatements string `json:"revocation_statements"`
	RollbackStatements   string `json:"rollback_statements"`
	RenewStatements      string `json:"renew_statements"`
}

type upgradeCheck struct {
	// This json tag has a typo in it, the new version does not. This
	// necessitates this upgrade logic.
	Statements *upgradeStatements `json:"statments,omitempty"`
}

func (b *databaseBackend) Role(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	return b.roleAtPath(ctx, s, roleName, databaseRolePath)
}

func (b *databaseBackend) StaticRole(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	return b.roleAtPath(ctx, s, roleName, databaseStaticRolePath)
}

func (b *databaseBackend) roleAtPath(ctx context.Context, s logical.Storage, roleName string, pathPrefix string) (*roleEntry, error) {
	entry, err := s.Get(ctx, pathPrefix+roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var upgradeCh upgradeCheck
	if err := entry.DecodeJSON(&upgradeCh); err != nil {
		return nil, err
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	switch {
	case upgradeCh.Statements != nil:
		var stmts v4.Statements
		if upgradeCh.Statements.CreationStatements != "" {
			stmts.Creation = []string{upgradeCh.Statements.CreationStatements}
		}
		if upgradeCh.Statements.RevocationStatements != "" {
			stmts.Revocation = []string{upgradeCh.Statements.RevocationStatements}
		}
		if upgradeCh.Statements.RollbackStatements != "" {
			stmts.Rollback = []string{upgradeCh.Statements.RollbackStatements}
		}
		if upgradeCh.Statements.RenewStatements != "" {
			stmts.Renewal = []string{upgradeCh.Statements.RenewStatements}
		}
		result.Statements = stmts
	}

	result.Statements.Revocation = strutil.RemoveEmpty(result.Statements.Revocation)

	// For backwards compatibility, copy the values back into the string form
	// of the fields
	result.Statements = dbutil.StatementCompatibilityHelper(result.Statements)

	return &result, nil
}

func (b *databaseBackend) invalidate(ctx context.Context, key string) {
	switch {
	case strings.HasPrefix(key, databaseConfigPath):
		name := strings.TrimPrefix(key, databaseConfigPath)
		b.ClearConnection(name)
	}
}

func (b *databaseBackend) GetConnection(ctx context.Context, s logical.Storage, name string) (*dbPluginInstance, error) {
	config, err := b.DatabaseConfig(ctx, s, name)
	if err != nil {
		return nil, err
	}

	return b.GetConnectionWithConfig(ctx, name, config)
}

func (b *databaseBackend) GetConnectionWithConfig(ctx context.Context, name string, config *DatabaseConfig) (*dbPluginInstance, error) {
	dbi := b.connections.Get(name)
	if dbi != nil {
		return dbi, nil
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	dbw, err := newDatabaseWrapper(ctx, config.PluginName, config.PluginVersion, b.System(), b.logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create database instance: %w", err)
	}

	initReq := v5.InitializeRequest{
		Config:           config.ConnectionDetails,
		VerifyConnection: true,
	}
	_, err = dbw.Initialize(ctx, initReq)
	if err != nil {
		dbw.Close()
		return nil, err
	}

	dbi = &dbPluginInstance{
		database: dbw,
		id:       id,
		name:     name,
	}
	oldConn := b.connections.Put(name, dbi)
	if oldConn != nil {
		err := oldConn.Close()
		if err != nil {
			b.Logger().Warn("Error closing database connection", "error", err)
		}
	}
	return dbi, nil
}

// ClearConnection closes the database connection and
// removes it from the b.connections map.
func (b *databaseBackend) ClearConnection(name string) error {
	db := b.connections.Pop(name)
	if db != nil {
		// Ignore error here since the database client is always killed
		db.Close()
	}
	return nil
}

// ClearConnectionId closes the database connection with a specific id and
// removes it from the b.connections map.
func (b *databaseBackend) ClearConnectionId(name, id string) error {
	db := b.connections.PopIfEqual(name, id)
	if db != nil {
		// Ignore error here since the database client is always killed
		db.Close()
	}
	return nil
}

func (b *databaseBackend) CloseIfShutdown(db *dbPluginInstance, err error) {
	// Plugin has shutdown, close it so next call can reconnect.
	switch err {
	case rpc.ErrShutdown, v4.ErrPluginShutdown, v5.ErrPluginShutdown:
		// Put this in a goroutine so that requests can run with the read or write lock
		// and simply defer the unlock.  Since we are attaching the instance and matching
		// the id in the connection map, we can safely do this.
		go func() {
			db.Close()

			// Delete the connection if it is still active.
			b.connections.PopIfEqual(db.name, db.id)
		}()
	}
}

// clean closes all connections from all database types
// and cancels any rotation queue loading operation.
func (b *databaseBackend) clean(_ context.Context) {
	// kill the queue and terminate the background ticker
	if b.cancelQueueCtx != nil {
		b.cancelQueueCtx()
	}

	connections := b.connections.Clear()
	for _, db := range connections {
		go db.Close()
	}
	b.gaugeCollectionProcessStop.Do(func() {
		if b.gaugeCollectionProcess != nil {
			b.gaugeCollectionProcess.Stop()
		}
		b.gaugeCollectionProcess = nil
	})
}

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, mssql, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/config/" path.
`
