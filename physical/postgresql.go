package physical

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-uuid"
	"github.com/lib/pq"
	log "github.com/mgutz/logxi/v1"
)

const (
	// PostgreSQLLockRetryInterval is the interval of time to wait between new
	// locking attempts
	PostgreSQLLockRetryInterval = time.Second
	// PostgreSQLLockErrorRetryMax is the number of retries to make after an
	// error fetching the status of a lock
	PostgreSQLLockErrorRetryMax = 5
	// PostgreSQLLockTTL is the maximum length of time of a lock, in seconds
	PostgreSQLLockTTL = 10
	// PostgreSQLLockRenewInterval is the interval of time locks are renewed
	PostgreSQLLockRenewInterval = 5 * time.Second
)

// PostgreSQL Backend is a physical backend that stores data
// within a PostgreSQL database.
type PostgreSQLBackend struct {
	table        string
	lock_table   string
	client       *sql.DB
	put_query    string
	get_query    string
	delete_query string
	list_query   string
	logger       log.Logger
}

// newPostgreSQLBackend constructs a PostgreSQL backend using the given
// API client, server address, credentials, and database.
func newPostgreSQLBackend(conf map[string]string, logger log.Logger) (Backend, error) {
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

	unquoted_lock_table, ok := conf["lock_table"]
	if !ok {
		unquoted_lock_table = "vault_lock"
	}
	quoted_lock_table := pq.QuoteIdentifier(unquoted_lock_table)

	// Create PostgreSQL handle for the database.
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	// Determine if we should use an upsert function (versions < 9.5)
	var upsert_required bool
	upsert_required_query := "SELECT string_to_array(setting, '.')::int[] < '{9,5}' FROM pg_settings WHERE name = 'server_version'"
	if err := db.QueryRow(upsert_required_query).Scan(&upsert_required); err != nil {
		return nil, fmt.Errorf("failed to check for native upsert: %v", err)
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
		lock_table:   quoted_lock_table,
		client:       db,
		put_query:    put_query,
		get_query:    "SELECT value FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		delete_query: "DELETE FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		list_query: "SELECT key FROM " + quoted_table + " WHERE path = $1" +
			"UNION SELECT DISTINCT substring(substr(path, length($1)+1) from '^.*?/') FROM " +
			quoted_table + " WHERE parent_path LIKE concat($1, '%')",
		logger: logger,
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
func (m *PostgreSQLBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"postgres", "put"}, time.Now())

	parentPath, path, key := m.splitKey(entry.Key)

	_, err := m.client.Exec(m.put_query, parentPath, path, key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *PostgreSQLBackend) Get(fullPath string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"postgres", "get"}, time.Now())

	_, path, key := m.splitKey(fullPath)

	var result []byte
	err := m.client.QueryRow(m.get_query, path, key).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ent := &Entry{
		Key:   key,
		Value: result,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *PostgreSQLBackend) Delete(fullPath string) error {
	defer metrics.MeasureSince([]string{"postgres", "delete"}, time.Now())

	_, path, key := m.splitKey(fullPath)

	_, err := m.client.Exec(m.delete_query, path, key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *PostgreSQLBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"postgres", "list"}, time.Now())

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
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (m *PostgreSQLBackend) HAEnabled() bool {
	return true
}

// PostgreSQLLock implements the Lock interface for PostgreSQL
type PostgreSQLLock struct {
	client *sql.DB
	table  string
	leader chan struct{}

	key     string
	value   string
	vaultID string
}

// LockWith initializes a Postgres backend lock
func (m *PostgreSQLBackend) LockWith(key, value string) (Lock, error) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &PostgreSQLLock{
		client:  m.client,
		table:   m.lock_table,
		leader:  make(chan struct{}),
		key:     key,
		value:   value,
		vaultID: id,
	}, nil
}

// Lock grabs a lock, or waits until it is available
func (m *PostgreSQLLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stopCh
		cancel()
	}()
	ticker := time.NewTicker(PostgreSQLLockRetryInterval)
	defer ticker.Stop()
	for {
		_, err := m.client.ExecContext(
			ctx,
			"DELETE FROM "+m.table+" WHERE expiration < now()",
		)
		if err != nil {
			return nil, err
		}
		lockTTL := strconv.Itoa(PostgreSQLLockTTL)
		_, err = m.client.ExecContext(
			ctx,
			`INSERT INTO `+m.table+` (key, value, vault_id, expiration)
				VALUES ($1, $2, $3, now() + interval '`+lockTTL+` seconds')`,
			m.key,
			m.value,
			m.vaultID,
		)
		if err == nil {
			break
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return nil, nil
		}
	}
	go m.watch()
	return m.leader, nil
}

// watch periodically queries the vault_lock table and closes the m.leader
// channel if the lock is lost
func (m *PostgreSQLLock) watch() {
	retries := PostgreSQLLockErrorRetryMax
	ticker := time.NewTicker(PostgreSQLLockRenewInterval)
	defer ticker.Stop()
	defer close(m.leader)
	for {
		select {
		case <-ticker.C:
			lockTTL := strconv.Itoa(PostgreSQLLockTTL)
			r, err := m.client.Exec(
				`UPDATE `+m.table+`
				SET expiration = now() + interval '`+lockTTL+` seconds'
				WHERE key = $1 AND vault_id = $2 AND expiration >= now()`,
				m.key,
				m.vaultID,
			)
			if err != nil || r == nil {
				retries--
				if retries == 0 {
					return
				}
				continue
			}
			if rows, _ := r.RowsAffected(); rows == 0 {
				// Lock lost!
				return
			}
			retries = PostgreSQLLockErrorRetryMax
		}
	}
}

// Unlock unlocks a lock. It returns an error if the lock was not in use.
func (m *PostgreSQLLock) Unlock() error {
	r, err := m.client.Exec(
		`DELETE FROM `+m.table+` WHERE key = $1 AND vault_id = $2`,
		m.key,
		m.vaultID,
	)
	if err != nil {
		return err
	}
	if r == nil {
		return errors.New("Unknown error reading a query's result")
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Lock not currently held")
	}
	return nil
}

// Value returns whether the lock is held and the value associated with it
func (m *PostgreSQLLock) Value() (held bool, value string, err error) {
	err = m.client.QueryRow(`
		SELECT expiration > now(), value
		FROM `+m.table+`
		WHERE key = $1`,
		m.key,
	).Scan(&held, &value)
	return held, value, err
}
