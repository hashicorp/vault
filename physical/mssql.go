package physical

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	_ "github.com/denisenkom/go-mssqldb"
)

// MySQLBackend is a physical backend that stores data
// within MySQL database.
type MsSQLBackend struct {
	dbTable    string
	client     *sql.DB
	statements map[string]*sql.Stmt
}

// newMsSQLBackend constructs a MsSQL backend using the given API client and
// server address and credential for accessing mysql database.
func newMsSQLBackend(conf map[string]string) (Backend, error) {
	// Get the MySQL credentials to perform read/write operations.
	username, ok := conf["username"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing username")
	}
	password, ok := conf["password"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get or set MsSQL server address. Defaults to localhost and default port(1433)
	var port string
	address, ok := conf["address"]
	if ok {
		s := strings.Split(address, ",")

		address = s[0]
		if len(s) > 1 {
			port = s[1]
		} else {
			port = "1433"
		}
	}	else {
		address = "127.0.0.1"
		port = "1433"
	}

	// Get the MsSQL database and table details.
	database, ok := conf["database"]
	if !ok {
		database = "vault"
	}
	table, ok := conf["table"]
	if !ok {
		table = "vault"
	}
	dbTable := database + ".." + table

	// Create MySQL handle for the database.
	//dsn := username + ":" + password + "@tcp(" + address + ")/"
	dsn := "port=" + port + ";server=" + address + ";user id=" + username + ";" + "password=" + password
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mssql: %v", err)
	}

	// Create the required database if it doesn't exists.
	if _, err := db.Exec("if not exists(select * from sys.databases where name = '" + database + "') create database " + database); err != nil {
		return nil, fmt.Errorf("failed to create mssql database: %v", err)
	}

	// Create the required table if it doesn't exists.
	create_query := "IF NOT EXISTS(SELECT 1 FROM " + database + ".INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_NAME='" + table + "') CREATE TABLE " + dbTable +
		" (vault_key varchar(512) PRIMARY KEY, vault_value varbinary(max))"
	if _, err := db.Exec(create_query); err != nil {
		return nil, fmt.Errorf("failed to create mssql table: %v", err)
	}

	// Setup the backend.
	m := &MsSQLBackend{
		dbTable:    dbTable,
		client:     db,
		statements: make(map[string]*sql.Stmt),
	}

	// Prepare all the statements required
	statements := map[string]string{
		"put": "IF EXISTS(SELECT 1 FROM " + dbTable + " WHERE vault_key = ?) UPDATE " + dbTable + " SET vault_value = ? WHERE vault_key = ? " +
		"ELSE INSERT INTO " + dbTable + " VALUES( ?, ? )",
		"get":    "SELECT vault_value FROM " + dbTable + " WHERE vault_key = ?",
		"delete": "DELETE FROM " + dbTable + " WHERE vault_key = ?",
		"list":   "SELECT vault_key FROM " + dbTable + " WHERE vault_key LIKE ?",
	}
	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// prepare is a helper to prepare a query for future execution
func (m *MsSQLBackend) prepare(name, query string) error {
	stmt, err := m.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare '%s': %v", name, err)
	}
	m.statements[name] = stmt
	return nil
}

// Put is used to insert or update an entry.
func (m *MsSQLBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"mssql", "put"}, time.Now())

	_, err := m.statements["put"].Exec(entry.Key, entry.Value, entry.Key, entry.Key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *MsSQLBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"mssql", "get"}, time.Now())

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
func (m *MsSQLBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"mssql", "delete"}, time.Now())

	_, err := m.statements["delete"].Exec(key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *MsSQLBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"mssql", "list"}, time.Now())

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
