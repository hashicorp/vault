package physical

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-metrics"
	"github.com/lib/pq"
)

// CockroachDBBackend Backend is a physical backend that stores data
// within a CockroachDB database.
type CockroachDBBackend struct {
	table        string
	client       *sql.DB
	put_query    string
	get_query    string
	delete_query string
	list_query   string
	logger       log.Logger
}

// newCockroachDBBackend constructs a CockroachDB backend using the given
// API client, server address, credentials, and database.
func newCockroachDBBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	// Get the CockroachDB credentials to perform read/write operations.
	connURL, ok := conf["connection_url"]
	if !ok || connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	unquotedTable, ok := conf["table"]
	if !ok {
		unquotedTable = "vault_kv_store"
	}
	quotedTable := pq.QuoteIdentifier(unquotedTable)

	// Create CockroachDB handle for the database.
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cockroachdb: %v", err)
	}

	// Create the required table if it doesn't exists.
	create_query := "CREATE TABLE IF NOT EXISTS " + unquotedTable +
		" (path STRING, value BYTES, PRIMARY KEY (path))"
	if _, err := db.Exec(create_query); err != nil {
		return nil, fmt.Errorf("failed to create mysql table: %v", err)
	}

	// Setup the backend.
	m := &CockroachDBBackend{
		table:  quotedTable,
		client: db,
		put_query: "INSERT INTO " + unquotedTable + " VALUES($1, $2)" +
			" ON CONFLICT (path) DO " +
			" UPDATE SET (path, value) = ($1, $2)",
		get_query:    "SELECT value FROM " + unquotedTable + " WHERE path = $1",
		delete_query: "DELETE FROM " + unquotedTable + " WHERE path = $1",
		list_query:   "SELECT path FROM " + unquotedTable + " WHERE path LIKE concat($1, '%')",
		logger:       logger,
	}

	return m, nil
}

// Put is used to insert or update an entry.
func (m *CockroachDBBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "put"}, time.Now())

	_, err := m.client.Exec(m.put_query, entry.Key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *CockroachDBBackend) Get(fullPath string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"cockroachdb", "get"}, time.Now())

	var result []byte
	err := m.client.QueryRow(m.get_query, fullPath).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ent := &Entry{
		Key:   fullPath,
		Value: result,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *CockroachDBBackend) Delete(fullPath string) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "delete"}, time.Now())

	_, err := m.client.Exec(m.delete_query, fullPath)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *CockroachDBBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"cockroachdb", "list"}, time.Now())

	rows, err := m.client.Query(m.list_query, prefix)
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

		key = strings.TrimPrefix(key, prefix)
		if i := strings.Index(key, "/"); i == -1 {
			// Add objects only from the current 'folder'
			keys = append(keys, key)
		} else if i != -1 {
			// Add truncated 'folder' paths
			keys = appendIfMissing(keys, string(key[:i+1]))
		}
	}

	sort.Strings(keys)
	return keys, nil
}
