package physical

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	_ "github.com/lib/pq"
)

// PostgreSQL Backend is a physical backend that stores data
// within a PostgreSQL database.
type PostgreSQLBackend struct {
	table      string
	client     *sql.DB
	statements map[string]*sql.Stmt
}

// newPostgreSQLBackend constructs a PostgreSQL backend using the given
// API client, server address, credentials, and database.
func newPostgreSQLBackend(conf map[string]string) (Backend, error) {
	// Get the PostgreSQL credentials to perform read/write operations.
	connURL, ok := conf["connection_url"]
	if !ok || connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	table, ok := conf["table"]
	if !ok {
		table = "vault"
	}

	// Create PostgreSQL handle for the database.
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	// Create the required table if it doesn't exists.
	create_query := "CREATE TABLE IF NOT EXISTS " + table +
		" (vault_key TEXT, vault_value BYTEA, PRIMARY KEY (vault_key))"
	if _, err := db.Exec(create_query); err != nil {
		return nil, fmt.Errorf("failed to create postgres table: %v", err)
	}

	// Determine if we need to create an upsert function (versions < 9.5)
	var upsert_missing bool
	upsert_missing_query := "SELECT string_to_array(setting, '.')::int[] < '{9,5}' FROM pg_settings WHERE name = 'server_version'"
	if err := db.QueryRow(upsert_missing_query).Scan(&upsert_missing); err != nil {
		return nil, fmt.Errorf("failed to check postgres server version: %v", err)
	}

	put_statement := "INSERT INTO " + table + " VALUES($1, $2) ON CONFLICT (vault_key) DO UPDATE SET vault_value = $2"
	if upsert_missing {
		// Create the upsert function
		// http://www.postgresql.org/docs/9.4/static/plpgsql-control-structures.html
		create_upsert := `
CREATE OR REPLACE FUNCTION upsert_vault(key TEXT, value BYTEA) RETURNS VOID AS
$$
BEGIN
    LOOP
        -- first try to update the key
        UPDATE ` + table + ` SET vault_value = value WHERE vault_key = key;
        IF found THEN
            RETURN;
        END IF;
        -- not there, so try to insert the key
        -- if someone else inserts the same key concurrently,
        -- we could get a unique-key failure
        BEGIN
            INSERT INTO ` + table + `(vault_key, vault_value) VALUES (key, value);
            RETURN;
        EXCEPTION WHEN unique_violation THEN
            -- Do nothing, and loop to try the UPDATE again.
        END;
    END LOOP;
END;
$$
LANGUAGE plpgsql;`

		if _, err := db.Exec(create_upsert); err != nil {
			return nil, fmt.Errorf("failed to create upsert function: %v", err)
		}

		put_statement = "SELECT upsert_vault($1, $2)"
	}

	// Setup the backend.
	m := &PostgreSQLBackend{
		table:      table,
		client:     db,
		statements: make(map[string]*sql.Stmt),
	}

	// Prepare all the statements required
	statements := map[string]string{
		"put":    put_statement,
		"get":    "SELECT vault_value FROM " + table + " WHERE vault_key = $1",
		"delete": "DELETE FROM " + table + " WHERE vault_key = $1",
		"list":   "SELECT vault_key FROM " + table + " WHERE vault_key LIKE $1",
	}
	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// prepare is a helper to prepare a query for future execution
func (m *PostgreSQLBackend) prepare(name, query string) error {
	stmt, err := m.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare '%s': %v", name, err)
	}
	m.statements[name] = stmt
	return nil
}

// Put is used to insert or update an entry.
func (m *PostgreSQLBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"postgres", "put"}, time.Now())

	_, err := m.statements["put"].Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *PostgreSQLBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"postgres", "get"}, time.Now())

	var result []byte
	err := m.statements["get"].QueryRow(key).Scan(&result)
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
func (m *PostgreSQLBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"postgres", "delete"}, time.Now())

	_, err := m.statements["delete"].Exec(key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *PostgreSQLBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"postgres", "list"}, time.Now())

	// Add the % wildcard to the prefix to do the prefix search
	likePrefix := prefix + "%"
	rows, err := m.statements["list"].Query(likePrefix)

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
