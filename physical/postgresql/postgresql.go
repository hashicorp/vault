package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/physical"
	//log "github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"

	"github.com/armon/go-metrics"
	"github.com/lib/pq"
)

const (

	// (seconds) The lock TTL matches the default that Consul API uses, 15 seconds.
	// used in as part of SQL commands to set/extend lock expiryt time relative to database clock
	PostgreSQLLockTTL = 15

	// The amount of time to wait between the lock renewals
	PostgreSQLLockRenewInterval = 5 * time.Second

	// PostgreSQLLockRetryInterval is the amount of time to wait
	// if a lock fails before trying again.
	PostgreSQLLockRetryInterval = time.Second
	// PostgreSQLWatchRetryMax is the number of times to re-try a
	// failed watch before signaling that leadership is lost.
	PostgreSQLWatchRetryMax = 5
	// PostgreSQLWatchRetryInterval is the amount of time to wait
	// if a watch fails before trying again.
	PostgreSQLWatchRetryInterval = 5 * time.Second
)

// Verify PostgreSQLBackend satisfies the correct interfaces
var _ physical.Backend = (*PostgreSQLBackend)(nil)

//
// HA backend was implemented based on the DynamoDB backend pattern
// With distinction using central postgres clock, hereby avoiding
// possible issues with multiple clocks
//
var _ physical.HABackend = (*PostgreSQLBackend)(nil)
var _ physical.Lock = (*PostgreSQLLock)(nil)

// PostgreSQL Backend is a physical backend that stores data
// within a PostgreSQL database.
type PostgreSQLBackend struct {
	table                    string
	client                   *sql.DB
	put_query                string
	get_query                string
	delete_query             string
	list_query               string
	haGetLockIdentityQuery   string
	haGetLockValueQuery      string
	haUpsertLockIdentityExec string
	haDeleteLockExec         string

	haEnabled  bool
	logger     log.Logger
	permitPool *physical.PermitPool
}

// PostgreSQLLock implements a lock using an PostgreSQL client.
type PostgreSQLLock struct {
	backend    *PostgreSQLBackend
	value, key string
	identity   string
	held       bool
	lock       sync.Mutex
}

// NewPostgreSQLBackend constructs a PostgreSQL backend using the given
// API client, server address, credentials, and database.
func NewPostgreSQLBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the PostgreSQL credentials to perform read/write operations.
	connURL, ok := conf["connection_url"]
	if !ok || connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	unquoted_table, ok := conf["table"]
	if !ok {
		unquoted_table = "vault_kv_store"
	}
	quoted_table := pq.QuoteIdentifier(unquoted_table)

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	var err error
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	} else {
		maxParInt = physical.DefaultParallelOperations
	}

	var hae bool = false
	haestr, ok := conf["ha_enabled"]
	if ok && haestr == "true" {
		hae = true
	}

	// Create PostgreSQL handle for the database.
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, errwrap.Wrapf("failed to connect to postgres: {{err}}", err)
	}
	db.SetMaxOpenConns(maxParInt)

	// Determine if we should use an upsert function (versions < 9.5)
	var upsert_required bool
	upsert_required_query := "SELECT current_setting('server_version_num')::int < 90500"
	if err := db.QueryRow(upsert_required_query).Scan(&upsert_required); err != nil {
		return nil, errwrap.Wrapf("failed to check for native upsert: {{err}}", err)
	}

	// Setup our put strategy based on the presence or absence of a native
	// upsert.
	var put_query string
	if upsert_required {
		put_query = "SELECT vault_kv_put($1, $2, $3, $4)"
	} else {
		put_query = "INSERT INTO " + quoted_table + " VALUES($1, $2, $3, $4)" +
			" ON CONFLICT (path, key) DO " +
			" UPDATE SET (parent_path, path, key, value) = ($1, $2, $3, $4)"
	}

	// Setup the backend.
	m := &PostgreSQLBackend{
		table:        quoted_table,
		client:       db,
		put_query:    put_query,
		get_query:    "SELECT value FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		delete_query: "DELETE FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		list_query: "SELECT key FROM " + quoted_table + " WHERE path = $1" +
			" UNION SELECT DISTINCT substring(substr(path, length($1)+1) from '^.*?/') FROM " + quoted_table +
			" WHERE parent_path LIKE $1 || '%'",
		haGetLockIdentityQuery:
		//only read non expired data
		" SELECT ha_identity FROM vault_ha_store WHERE NOW() <= valid_until AND ha_key = $1 ",
		haGetLockValueQuery:
		//only read non expired data
		" SELECT ha_value FROM vault_ha_store WHERE NOW() <= valid_until AND ha_key = $1 ",
		haUpsertLockIdentityExec:
		// $1=identity $2=ha_key $3=ha_value $4=TTL in seconds
		//update either steal expired lock OR update expiry for lock owned by me
		" INSERT INTO vault_ha_store (ha_identity, ha_key, ha_value, valid_until) VALUES ($1, $2, $3, NOW() + $4 * INTERVAL '1 seconds'  ) " +
			" ON CONFLICT (ha_key) DO " +
			" UPDATE SET (ha_identity, ha_key, ha_value, valid_until) = ($1, $2, $3, NOW() + $4 * INTERVAL '1 seconds') " +
			" WHERE vault_ha_store.valid_until < NOW() AND vault_ha_store.ha_key = $2 OR " +
			" vault_ha_store.ha_identity = $1 AND vault_ha_store.ha_key = $2  ",
		haDeleteLockExec:
		//$1=ha_identity $2=ha_key
		" DELETE FROM vault_ha_store WHERE ha_identity=$1 AND ha_key=$2 ",
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
		haEnabled:  hae,
	}

	return m, nil
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

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	parentPath, path, key := m.splitKey(entry.Key)

	_, err := m.client.Exec(m.put_query, parentPath, path, key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *PostgreSQLBackend) Get(ctx context.Context, fullPath string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"postgres", "get"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	_, path, key := m.splitKey(fullPath)

	var result []byte
	err := m.client.QueryRow(m.get_query, path, key).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ent := &physical.Entry{
		Key:   key,
		Value: result,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *PostgreSQLBackend) Delete(ctx context.Context, fullPath string) error {
	defer metrics.MeasureSince([]string{"postgres", "delete"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	_, path, key := m.splitKey(fullPath)

	_, err := m.client.Exec(m.delete_query, path, key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *PostgreSQLBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"postgres", "list"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	rows, err := m.client.Query(m.list_query, "/"+prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, errwrap.Wrapf("failed to scan rows: {{err}}", err)
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
		backend:  p,
		key:      key,
		value:    value,
		identity: identity,
	}, nil
}

func (p *PostgreSQLBackend) HAEnabled() bool {
	return p.haEnabled
}

// Lock tries to acquire the lock by repeatedly trying to create
// a record in the PostgreSQL table. It will block until either the
// stop channel is closed or the lock could be acquired successfully.
// The returned channel will be closed once the lock is deleted or
// changed in the PostgreSQL table.
func (l *PostgreSQLLock) Lock(stopCh <-chan struct{}) (doneCh <-chan struct{}, retErr error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.held {
		return nil, fmt.Errorf("lock already held")
	}

	done := make(chan struct{})
	// close done channel even in case of error
	defer func() {
		if retErr != nil {
			close(done)
		}
	}()

	var (
		stop    = make(chan struct{})
		success = make(chan struct{})
		errors  = make(chan error)
		leader  = make(chan struct{})
	)
	// try to acquire the lock asynchronously
	go l.tryToLock(stop, success, errors)

	select {
	case <-success:
		l.held = true
		// after acquiring it successfully, we must renew the lock periodically,
		// and watch the lock in order to close the leader channel
		// once it is lost.
		go l.periodicallyRenewLock(leader)
		go l.watch(leader)
	case retErr = <-errors:
		close(stop)
		return nil, retErr
	case <-stopCh:
		close(stop)
		return nil, nil
	}

	return leader, retErr
}

// Unlock releases the lock by deleting the lock record from the
// PostgreSQL table.
func (l *PostgreSQLLock) Unlock() error {
	pg := l.backend
	pg.permitPool.Acquire()
	defer pg.permitPool.Release()

	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.held {
		return nil
	}
	l.held = false
	//Delete lock owned by me
	_, err := pg.client.Exec(pg.haDeleteLockExec, l.identity, l.key)
	return err
}

// Value checks whether or not the lock is held by any instance of PostgreSQLLock,
// including this one, and returns the current value.
func (l *PostgreSQLLock) Value() (bool, string, error) {

	pg := l.backend
	pg.permitPool.Acquire()
	defer pg.permitPool.Release()
	var result string

	err := pg.client.QueryRow(pg.haGetLockValueQuery, l.key).Scan(&result)

	if err != nil {
		return false, "", err
	}

	return l.held, result, nil
}

// tryToLock tries to create a new item in PostgreSQL
// every `PostgreSQLLockRetryInterval`. As long as the item
// cannot be created (because it already exists), it will
// be retried. If the operation fails due to an error, it
// is sent to the errors channel.
// When the lock could be acquired successfully, the success
// channel is closed.
func (l *PostgreSQLLock) tryToLock(stop, success chan struct{}, errors chan error) {
	ticker := time.NewTicker(PostgreSQLLockRetryInterval)

	for {
		select {
		case <-stop:
			ticker.Stop()
		case <-ticker.C:
			err := l.writeItem()
			if err != nil {
				errors <- err
				return
			} else {
				ticker.Stop()
				close(success)
				return
			}
		}
	}
}

func (l *PostgreSQLLock) periodicallyRenewLock(done chan struct{}) {
	ticker := time.NewTicker(PostgreSQLLockRenewInterval)
	for {
		select {
		case <-ticker.C:
			l.writeItem()
		case <-done:
			ticker.Stop()
			return
		}
	}
}

// watch checks whether the lock has changed in the
// PostgreSQL table and closes the leader channel if so.
// The interval is set by `PostgreSQLWatchRetryInterval`.
// If an error occurs during the check, watch will retry
// the operation for `PostgreSQLWatchRetryMax` times and
// close the leader channel if it can't succeed.
func (l *PostgreSQLLock) watch(lost chan struct{}) {
	ticker := time.NewTicker(PostgreSQLLockRetryInterval)

	pb := l.backend
	pb.permitPool.Acquire()
	defer pb.permitPool.Release()

WatchLoop:
	for {
		select {
		case <-ticker.C:
			var result string
			err := pb.client.QueryRow(pb.haGetLockIdentityQuery, l.key).Scan(&result)
			//no row or different indenty, we have lost the lock !
			if err != nil || result != l.identity {
				break WatchLoop
			}
		}
	}
	//signal lost lock
	close(lost)
}

// Attempts to put/update the PostgreSQL item using condition expressions to
// evaluate the TTL.
func (l *PostgreSQLLock) writeItem() error {

	pg := l.backend
	pg.permitPool.Acquire()
	defer pg.permitPool.Release()

	//Try steal lock or update expiry on my lock

	sqlResult, err := pg.client.Exec(pg.haUpsertLockIdentityExec, l.identity, l.key, l.value, PostgreSQLLockTTL)
	if err != nil {
		return err
	}
	if sqlResult != nil {
		ar, err := sqlResult.RowsAffected()
		if err == nil && ar == 1 {
			return nil
		}
	}
	return fmt.Errorf("Could not obtain lock.")
}
