// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/physical"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (

	// The lock TTL matches the default that Consul API uses, 15 seconds.
	// Used as part of SQL commands to set/extend lock expiry time relative to
	// database clock.
	PostgreSQLLockTTLSeconds = 15

	// The amount of time to wait between the lock renewals
	PostgreSQLLockRenewInterval = 5 * time.Second

	// PostgreSQLLockRetryInterval is the amount of time to wait
	// if a lock fails before trying again.
	PostgreSQLLockRetryInterval = time.Second
)

// Verify PostgreSQLBackend satisfies the correct interfaces
var _ physical.Backend = (*PostgreSQLBackend)(nil)

// HA backend was implemented based on the DynamoDB backend pattern
// With distinction using central postgres clock, hereby avoiding
// possible issues with multiple clocks
var (
	_ physical.HABackend = (*PostgreSQLBackend)(nil)
	_ physical.Lock      = (*PostgreSQLLock)(nil)
)

// PostgreSQL Backend is a physical backend that stores data
// within a PostgreSQL database.
type PostgreSQLBackend struct {
	table        string
	client       *sql.DB
	put_query    string
	get_query    string
	delete_query string
	list_query   string

	ha_table                 string
	haGetLockValueQuery      string
	haUpsertLockIdentityExec string
	haDeleteLockExec         string

	haEnabled  bool
	logger     log.Logger
	permitPool *permitpool.Pool
}

// PostgreSQLLock implements a lock using an PostgreSQL client.
type PostgreSQLLock struct {
	backend    *PostgreSQLBackend
	value, key string
	identity   string
	lock       sync.Mutex

	renewTicker *time.Ticker

	// ttlSeconds is how long a lock is valid for
	ttlSeconds int

	// renewInterval is how much time to wait between lock renewals.  must be << ttl
	renewInterval time.Duration

	// retryInterval is how much time to wait between attempts to grab the lock
	retryInterval time.Duration
}

// NewPostgreSQLBackend constructs a PostgreSQL backend using the given
// API client, server address, credentials, and database.
func NewPostgreSQLBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the PostgreSQL credentials to perform read/write operations.
	connURL := connectionURL(conf)
	if connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	unquoted_table, ok := conf["table"]
	if !ok {
		unquoted_table = "vault_kv_store"
	}
	quoted_table := dbutil.QuoteIdentifier(unquoted_table)

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	var err error
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	} else {
		maxParInt = physical.DefaultParallelOperations
	}

	maxIdleConnsStr, maxIdleConnsIsSet := conf["max_idle_connections"]
	var maxIdleConns int
	if maxIdleConnsIsSet {
		maxIdleConns, err = strconv.Atoi(maxIdleConnsStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_idle_connections parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_idle_connections set", "max_idle_connections", maxIdleConnsStr)
		}
	}

	// Create PostgreSQL handle for the database.
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	db.SetMaxOpenConns(maxParInt)

	if maxIdleConnsIsSet {
		db.SetMaxIdleConns(maxIdleConns)
	}

	// Determine if we should use a function to work around lack of upsert (versions < 9.5)
	var upsertAvailable bool
	upsertAvailableQuery := "SELECT current_setting('server_version_num')::int >= 90500"
	if err := db.QueryRow(upsertAvailableQuery).Scan(&upsertAvailable); err != nil {
		return nil, fmt.Errorf("failed to check for native upsert: %w", err)
	}

	if !upsertAvailable && conf["ha_enabled"] == "true" {
		return nil, fmt.Errorf("ha_enabled=true in config but PG version doesn't support HA, must be at least 9.5")
	}

	// Setup our put strategy based on the presence or absence of a native
	// upsert.
	var put_query string
	if !upsertAvailable {
		put_query = "SELECT vault_kv_put($1, $2, $3, $4)"
	} else {
		put_query = "INSERT INTO " + quoted_table + " VALUES($1, $2, $3, $4)" +
			" ON CONFLICT (path, key) DO " +
			" UPDATE SET (parent_path, path, key, value) = ($1, $2, $3, $4)"
	}

	unquoted_ha_table, ok := conf["ha_table"]
	if !ok {
		unquoted_ha_table = "vault_ha_locks"
	}
	quoted_ha_table := dbutil.QuoteIdentifier(unquoted_ha_table)

	// Setup the backend.
	m := &PostgreSQLBackend{
		table:        quoted_table,
		client:       db,
		put_query:    put_query,
		get_query:    "SELECT value FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		delete_query: "DELETE FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		list_query: "SELECT key FROM " + quoted_table + " WHERE path = $1" +
			" UNION ALL SELECT DISTINCT substring(substr(path, length($1)+1) from '^.*?/') FROM " + quoted_table +
			" WHERE parent_path LIKE $1 || '%'",
		haGetLockValueQuery:
		// only read non expired data
		" SELECT ha_value FROM " + quoted_ha_table + " WHERE NOW() <= valid_until AND ha_key = $1 ",
		haUpsertLockIdentityExec:
		// $1=identity $2=ha_key $3=ha_value $4=TTL in seconds
		// update either steal expired lock OR update expiry for lock owned by me
		" INSERT INTO " + quoted_ha_table + " as t (ha_identity, ha_key, ha_value, valid_until) VALUES ($1, $2, $3, NOW() + $4 * INTERVAL '1 seconds'  ) " +
			" ON CONFLICT (ha_key) DO " +
			" UPDATE SET (ha_identity, ha_key, ha_value, valid_until) = ($1, $2, $3, NOW() + $4 * INTERVAL '1 seconds') " +
			" WHERE (t.valid_until < NOW() AND t.ha_key = $2) OR " +
			" (t.ha_identity = $1 AND t.ha_key = $2)  ",
		haDeleteLockExec:
		// $1=ha_identity $2=ha_key
		" DELETE FROM " + quoted_ha_table + " WHERE ha_identity=$1 AND ha_key=$2 ",
		logger:     logger,
		permitPool: permitpool.New(maxParInt),
		haEnabled:  conf["ha_enabled"] == "true",
	}

	return m, nil
}

// connectionURL first check the environment variables for a connection URL. If
// no connection URL exists in the environment variable, the Vault config file is
// checked. If neither the environment variables or the config file set the connection
// URL for the Postgres backend, because it is a required field, an error is returned.
func connectionURL(conf map[string]string) string {
	connURL := conf["connection_url"]
	if envURL := os.Getenv("VAULT_PG_CONNECTION_URL"); envURL != "" {
		connURL = envURL
	}

	return connURL
}

// splitKey is a helper to split a full path key into individual
// parts: parentPath, path, key
func (m *PostgreSQLBackend) splitKey(fullPath string) (string, string, string) {
	var parentPath string
	var path string

	pieces := strings.Split(fullPath, "/")
	depth := len(pieces)
	key := pieces[depth-1]

	if depth == 1 {
		parentPath = ""
		path = "/"
	} else if depth == 2 {
		parentPath = "/"
		path = "/" + pieces[0] + "/"
	} else {
		parentPath = "/" + strings.Join(pieces[:depth-2], "/") + "/"
		path = "/" + strings.Join(pieces[:depth-1], "/") + "/"
	}

	return parentPath, path, key
}

// Put is used to insert or update an entry.
func (m *PostgreSQLBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"postgres", "put"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer m.permitPool.Release()

	parentPath, path, key := m.splitKey(entry.Key)

	_, err := m.client.ExecContext(ctx, m.put_query, parentPath, path, key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *PostgreSQLBackend) Get(ctx context.Context, fullPath string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"postgres", "get"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer m.permitPool.Release()

	_, path, key := m.splitKey(fullPath)

	var result []byte
	err := m.client.QueryRowContext(ctx, m.get_query, path, key).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ent := &physical.Entry{
		Key:   fullPath,
		Value: result,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *PostgreSQLBackend) Delete(ctx context.Context, fullPath string) error {
	defer metrics.MeasureSince([]string{"postgres", "delete"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer m.permitPool.Release()

	_, path, key := m.splitKey(fullPath)

	_, err := m.client.ExecContext(ctx, m.delete_query, path, key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *PostgreSQLBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"postgres", "list"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer m.permitPool.Release()

	rows, err := m.client.QueryContext(ctx, m.list_query, "/"+prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// LockWith is used for mutual exclusion based on the given key.
func (p *PostgreSQLBackend) LockWith(key, value string) (physical.Lock, error) {
	identity, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &PostgreSQLLock{
		backend:       p,
		key:           key,
		value:         value,
		identity:      identity,
		ttlSeconds:    PostgreSQLLockTTLSeconds,
		renewInterval: PostgreSQLLockRenewInterval,
		retryInterval: PostgreSQLLockRetryInterval,
	}, nil
}

func (p *PostgreSQLBackend) HAEnabled() bool {
	return p.haEnabled
}

// Lock tries to acquire the lock by repeatedly trying to create a record in the
// PostgreSQL table. It will block until either the stop channel is closed or
// the lock could be acquired successfully. The returned channel will be closed
// once the lock in the PostgreSQL table cannot be renewed, either due to an
// error speaking to PostgreSQL or because someone else has taken it.
func (l *PostgreSQLLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	var (
		success = make(chan struct{})
		errors  = make(chan error)
		leader  = make(chan struct{})
	)
	// try to acquire the lock asynchronously
	go l.tryToLock(stopCh, success, errors)

	select {
	case <-success:
		// after acquiring it successfully, we must renew the lock periodically
		l.renewTicker = time.NewTicker(l.renewInterval)
		go l.periodicallyRenewLock(leader)
	case err := <-errors:
		return nil, err
	case <-stopCh:
		return nil, nil
	}

	return leader, nil
}

// Unlock releases the lock by deleting the lock record from the
// PostgreSQL table.
func (l *PostgreSQLLock) Unlock() error {
	pg := l.backend
	if err := pg.permitPool.Acquire(context.Background()); err != nil {
		return err
	}
	defer pg.permitPool.Release()

	if l.renewTicker != nil {
		l.renewTicker.Stop()
	}

	// Delete lock owned by me
	_, err := pg.client.Exec(pg.haDeleteLockExec, l.identity, l.key)
	return err
}

// Value checks whether or not the lock is held by any instance of PostgreSQLLock,
// including this one, and returns the current value.
func (l *PostgreSQLLock) Value() (bool, string, error) {
	pg := l.backend
	if err := pg.permitPool.Acquire(context.Background()); err != nil {
		return false, "", err
	}
	defer pg.permitPool.Release()
	var result string
	err := pg.client.QueryRow(pg.haGetLockValueQuery, l.key).Scan(&result)

	switch err {
	case nil:
		return true, result, nil
	case sql.ErrNoRows:
		return false, "", nil
	default:
		return false, "", err

	}
}

// tryToLock tries to create a new item in PostgreSQL every `retryInterval`.
// As long as the item cannot be created (because it already exists), it will
// be retried. If the operation fails due to an error, it is sent to the errors
// channel. When the lock could be acquired successfully, the success channel
// is closed.
func (l *PostgreSQLLock) tryToLock(stop <-chan struct{}, success chan struct{}, errors chan error) {
	ticker := time.NewTicker(l.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			gotlock, err := l.writeItem()
			switch {
			case err != nil:
				errors <- err
				return
			case gotlock:
				close(success)
				return
			}
		}
	}
}

func (l *PostgreSQLLock) periodicallyRenewLock(done chan struct{}) {
	for range l.renewTicker.C {
		gotlock, err := l.writeItem()
		if err != nil || !gotlock {
			close(done)
			l.renewTicker.Stop()
			return
		}
	}
}

// Attempts to put/update the PostgreSQL item using condition expressions to
// evaluate the TTL.  Returns true if the lock was obtained, false if not.
// If false error may be nil or non-nil: nil indicates simply that someone
// else has the lock, whereas non-nil means that something unexpected happened.
func (l *PostgreSQLLock) writeItem() (bool, error) {
	pg := l.backend
	if err := pg.permitPool.Acquire(context.Background()); err != nil {
		return false, err
	}
	defer pg.permitPool.Release()

	// Try steal lock or update expiry on my lock

	sqlResult, err := pg.client.Exec(pg.haUpsertLockIdentityExec, l.identity, l.key, l.value, l.ttlSeconds)
	if err != nil {
		return false, err
	}
	if sqlResult == nil {
		return false, fmt.Errorf("empty SQL response received")
	}

	ar, err := sqlResult.RowsAffected()
	if err != nil {
		return false, err
	}
	return ar == 1, nil
}
