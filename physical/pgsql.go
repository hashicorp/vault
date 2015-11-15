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

// PostgreSQLBackend is a physical backend that stores data
// within PostgreSQL database.
type PostgreSQLBackend struct {
	dbTable    string
	client     *sql.DB
	statements map[string]*sql.Stmt
}

// newPostgreSQLBackend constructs a PostgreSQL backend using the given API client and
// server address and credential for accessing postgresql database.
func newPostgreSQLBackend(conf map[string]string) (Backend, error) {
	// Get the PostgreSQL credentials to perform read/write operations.
	username, ok := conf["username"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing username")
	}
	password, ok := conf["password"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get or set PostgreSQL server address. Defaults to localhost and default port(5432)
	address, ok := conf["address"]
	if !ok {
		address = "127.0.0.1:5432"
	}

	// Get the PostgreSQL database and table details.
	database, ok := conf["database"]
	if !ok {
		database = "vault"
	}
	table, ok := conf["table"]
	if !ok {
		table = "vault"
	}
	dbTable := database + "." + table

	// Create PostgreSQL handle for the database.
	dsn := "postgres://" + username + ":" + password + "@" + address + "/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %v", err)
	}

       // Check if the database exists.
       var exists []byte
       exists_query := "SELECT FROM pg_database WHERE datname='" + database + "'"
       err = db.QueryRow(exists_query).Scan(&exists)
       if err == sql.ErrNoRows {
               // Database doesn't exists, create the database first.
               if _, err := db.Exec("CREATE DATABASE " + database); err != nil {
                       return nil, fmt.Errorf("failed to create postgresql database: %v", err)
               }
       }

       // Create the required schema if it doesn't exists.
       if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + database); err != nil {
               return nil, fmt.Errorf("failed to create postgresql schema: %v", err)
       }

	// Create the required table if it doesn't exists.
	create_query := "CREATE TABLE IF NOT EXISTS " + dbTable +
		" (vault_key bytea, vault_value bytea, PRIMARY KEY (vault_key))"
	if _, err := db.Exec(create_query); err != nil {
		return nil, fmt.Errorf("failed to create postgresql table: %v", err)
	}

	// Setup the backend.
	m := &PostgreSQLBackend{
		dbTable:    dbTable,
		client:     db,
		statements: make(map[string]*sql.Stmt),
	}

	// Prepare all the statements required
	statements := map[string]string{
		"put": "INSERT INTO " + dbTable +
			" VALUES( $1, $2 ) ON CONFLICT (vault_key) DO UPDATE SET vault_value = $2",
		"get":    "SELECT vault_value FROM " + dbTable + " WHERE vault_key = $1",
		"delete": "DELETE FROM " + dbTable + " WHERE vault_key = $1",
		"list":   "SELECT vault_key FROM " + dbTable + " WHERE vault_key LIKE $1",
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
	defer metrics.MeasureSince([]string{"postgresql", "put"}, time.Now())

	_, err := m.statements["put"].Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *PostgreSQLBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"postgresql", "get"}, time.Now())

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
	defer metrics.MeasureSince([]string{"postgresql", "delete"}, time.Now())

	_, err := m.statements["delete"].Exec(key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *PostgreSQLBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"postgresql", "list"}, time.Now())

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
