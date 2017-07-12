package physical

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-uuid"
	"github.com/lib/pq"
	"github.com/mgutz/logxi/v1"
)

const (
	// DefaultPostgreSQLPollInterval is the interval of time to wait between new
	// locking attempts
	DefaultPostgreSQLPollInterval = 1 * time.Second
	// DefaultPostgreSQLLockTTL is the default TTL for a HA lock
	DefaultPostgreSQLLockTTL        = 10 * time.Second
	DefaultPostgreSQLLockSchemaName = ""
	DefaultPostgreSQLLockTableName  = "vault_lock"

	PostgreSQLLockTTLConf            = "lock_ttl"
	PostgreSQLPollIntervalConf       = "poll_interval"
	PostgreSQLLockTableNameConf      = "lock_table"
	PostgreSQLLockSchemaNameConf     = "lock_schema"
	MinimumPostgreSQLPollInterval    = 1 * time.Second
	MinimumPostgreSQLLockGracePeriod = 1 * time.Second
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
	logger       log.Logger

	lockSchemaName string
	lockTableName  string
	pollInterval   time.Duration
	lockTTL        time.Duration
}

// newPostgreSQLBackend constructs a PostgreSQL backend using the given
// API client, server address, credentials, and database.
func newPostgreSQLBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	// Get the PostgreSQL credentials to perform read/write operations.
	connURL, ok := conf["connection_url"]
	if !ok || connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	unquoted_table, ok := conf["lockTableName"]
	if !ok {
		unquoted_table = "vault_kv_store"
	}
	quoted_table := pq.QuoteIdentifier(unquoted_table)

	var err error

	lockTTL := DefaultPostgreSQLLockTTL
	if rawLockTTL, found := conf[PostgreSQLLockTTLConf]; found {
		if lockTTL, err = time.ParseDuration(rawLockTTL); err != nil {
			return nil, fmt.Errorf("%s error: %v", PostgreSQLLockTTLConf, err)
		}
	}

	pollInterval := DefaultPostgreSQLPollInterval
	if rawPollInterval, found := conf[PostgreSQLPollIntervalConf]; found {
		if pollInterval, err = time.ParseDuration(rawPollInterval); err != nil {
			return nil, fmt.Errorf("%s error: %v", err, PostgreSQLPollIntervalConf)
		}
	}

	lockTableName := DefaultPostgreSQLLockTableName
	if rawLockTableName, found := conf[PostgreSQLLockTableNameConf]; found {
		lockTableName = strings.TrimSpace(rawLockTableName)
	}

	lockSchemaName := DefaultPostgreSQLLockSchemaName
	if rawLockSchemaName, found := conf[PostgreSQLLockSchemaNameConf]; found {
		lockSchemaName = strings.TrimSpace(rawLockSchemaName)
	}

	// Sanity check inputs

	if pollInterval < 0 {
		return nil, fmt.Errorf("%s (%q) must be a positive time duration",
			PostgreSQLPollIntervalConf, pollInterval)
	}

	if !(pollInterval < lockTTL) {
		return nil, fmt.Errorf("%s (%q) must be smaller than the %s (%q)",
			PostgreSQLPollIntervalConf, PostgreSQLLockTTLConf, pollInterval,
			lockTTL)
	}

	if pollInterval < MinimumPostgreSQLPollInterval {
		return nil, fmt.Errorf("%s (%q) can not be less than %q",
			PostgreSQLPollIntervalConf, pollInterval,
			MinimumPostgreSQLPollInterval)
	}

	if lockTTL-pollInterval < MinimumPostgreSQLLockGracePeriod {
		return nil, fmt.Errorf(
			"There must be at least %s between the %s (%q) and %s (%q)",
			MinimumPostgreSQLLockGracePeriod, PostgreSQLPollIntervalConf,
			pollInterval, PostgreSQLLockTTLConf, lockTTL)
	}

	if lockTableName == "" {
		return nil, fmt.Errorf("%s error: can not be an empty string",
			PostgreSQLLockTableNameConf)
	}

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
		client:       db,
		put_query:    put_query,
		get_query:    "SELECT value FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		delete_query: "DELETE FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		list_query: "SELECT key FROM " + quoted_table + " WHERE path = $1" +
			"UNION SELECT DISTINCT substring(substr(path, length($1)+1) from '^.*?/') FROM " +
			quoted_table + " WHERE parent_path LIKE concat($1, '%')",
		logger: logger,

		lockSchemaName: lockSchemaName,
		lockTableName:  lockTableName,
		pollInterval:   pollInterval,
		lockTTL:        lockTTL,
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
	client         *sql.DB
	hostname       string
	lockSchemaName string
	lockTableName  string
	pollInterval   time.Duration
	lockTTL        time.Duration

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
	// Record the hostname to give DBAs a chance to figure out which Vault
	// service has the lock
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "vault"
	}
	return &PostgreSQLLock{
		client:         m.client,
		hostname:       hostname,
		lockSchemaName: m.lockSchemaName,
		lockTableName:  m.lockTableName,
		pollInterval:   m.pollInterval,
		lockTTL:        m.lockTTL,
		leader:         make(chan struct{}),
		key:            key,
		value:          value,
		vaultID:        fmt.Sprintf("%s-%s", hostname, id),
	}, nil
}

// Lock grabs a lock, or waits until it is available
func (m *PostgreSQLLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stopCh
		cancel()
	}()
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()
	for {
		cleanupSQL := fmt.Sprintf(
			"DELETE FROM %s WHERE expiration < now()", m.relationName(),
		)
		_, err := m.client.ExecContext(ctx, cleanupSQL)
		if err != nil {
			return nil, err
		}
		lockSQL := fmt.Sprintf(`INSERT INTO %s (key, value, vault_id,
			expiration) VALUES ($1, $2, $3, now() + $4::INTERVAL)`,
			m.relationName())
		_, err = m.client.ExecContext(ctx, lockSQL, m.key, m.value, m.vaultID,
			m.lockTTL.String())
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

func (m *PostgreSQLLock) relationName() string {
	relationName := pq.QuoteIdentifier(m.lockTableName)
	if m.lockSchemaName != "" {
		relationName = fmt.Sprintf(
			"%s.%s",
			pq.QuoteIdentifier(m.lockSchemaName),
			pq.QuoteIdentifier(m.lockTableName),
		)
	}
	return relationName
}

// watch periodically queries the lock table and closes the m.leader channel if
// the lock is lost
func (m *PostgreSQLLock) watch() {
	// refresh the lock halfway through its expiration
	refreshTicker := time.NewTicker(m.lockTTL / 2)
	defer refreshTicker.Stop()
	pollTicker := time.NewTicker(m.pollInterval)
	defer pollTicker.Stop()
	ticker := refreshTicker

	defer close(m.leader)

	var lastRefresh time.Time
	for {
		select {
		case <-ticker.C:
			refreshLockSQL := fmt.Sprintf(
				`UPDATE %s SET expiration = now() + $1::INTERVAL WHERE key = $2
				AND vault_id = $3 AND expiration >= now()`,
				m.relationName(),
			)
			r, err := m.client.Exec(refreshLockSQL, m.lockTTL.String(), m.key,
				m.vaultID)
			if err != nil || r == nil {
				if lastRefresh.Add(m.lockTTL).Before(time.Now()) {
					// Lock is definitely expired by now
					return
				}
				// Refresh faster
				ticker = pollTicker
				continue
			}
			ticker = refreshTicker

			if rows, _ := r.RowsAffected(); rows == 0 {
				// Lock lost!
				return
			}
			lastRefresh = time.Now()
		}
	}
}

// Unlock unlocks a lock. It returns an error if the lock was not in use.
func (m *PostgreSQLLock) Unlock() error {
	unlockSQL := fmt.Sprintf(
		"DELETE FROM %s WHERE key = $1 AND vault_id = $2", m.relationName(),
	)
	r, err := m.client.Exec(unlockSQL, m.key, m.vaultID)
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
	valueSQL := fmt.Sprintf(
		"SELECT expiration > now(), value FROM %s WHERE key = $1",
		m.lockTableName,
	)
	err = m.client.QueryRow(valueSQL, m.key).Scan(&held, &value)
	return held, value, err
}
