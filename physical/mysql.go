package physical

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/armon/go-metrics"
	_ "github.com/go-sql-driver/mysql"
)

var (
	MySQLPrepareStmtFailure = errors.New("failed to prepare statement")
	MySQLExecuteStmtFailure = errors.New("failed to execute statement")
)

// MySQLBackend is a physical backend that stores data
// within MySQL database.
type MySQLBackend struct {
	table      string
	database   string
	client     *sql.DB
	statements map[string]*sql.Stmt
}

// newMySQLBackend constructs a MySQL backend using the given API client and
// server address and credential for accessing mysql database.
func newMySQLBackend(conf map[string]string) (Backend, error) {
	// Get or set MySQL server address. Defaults to localhost and default port(3306)
	address, ok := conf["address"]
	if !ok {
		address = "127.0.0.1:3306"
	}

	// Get the MySQL credentials to perform read/write operations.
	username, ok := conf["username"]
	password, ok := conf["password"]

	// Get the MySQL database and table details.
	database, ok := conf["database"]
	if !ok {
		return nil, fmt.Errorf("database name is missing in the configuration")
	}
	table, ok := conf["table"]
	if !ok {
		return nil, fmt.Errorf("table name is missing in the configuration")
	}

	// Create MySQL handle for the database.
	dsn := username + ":" + password + "@tcp(" + address + ")/" + database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open handler with database")
	}
	defer db.Close()

	// Create the required table if it doesn't exists.
	create_query := "CREATE TABLE IF NOT EXISTS " + database + "." + table + " (vault_key varchar(512), vault_value mediumblob, PRIMARY KEY (vault_key))"
	create_stmt, err := db.Prepare(create_query)
	if err != nil {
		return nil, MySQLPrepareStmtFailure
	}
	defer create_stmt.Close()

	_, err = create_stmt.Exec()
	if err != nil {
		return nil, MySQLExecuteStmtFailure
	}

	// Map of query type as key to prepared statement.
	statements := make(map[string]*sql.Stmt)

	// Prepare statement for put query.
	insert_query := "INSERT INTO " + database + "." + table + " VALUES( ?, ? ) ON DUPLICATE KEY UPDATE vault_value=VALUES(vault_value)"
	insert_stmt, err := db.Prepare(insert_query)
	if err != nil {
		return nil, MySQLPrepareStmtFailure
	}
	statements["put"] = insert_stmt
	defer insert_stmt.Close()

	// Prepare statement for select query.
	select_query := "SELECT vault_value FROM " + database + "." + table + " WHERE vault_key = ?"
	select_stmt, err := db.Prepare(select_query)
	if err != nil {
		return nil, MySQLPrepareStmtFailure
	}
	statements["get"] = select_stmt
	defer select_stmt.Close()

	// Prepare statement for delete query.
	delete_query := "DELETE FROM " + database + "." + table + " WHERE vault_key = ?"
	delete_stmt, err := db.Prepare(delete_query)
	if err != nil {
		return nil, MySQLPrepareStmtFailure
	}
	statements["delete"] = delete_stmt
	defer delete_stmt.Close()

	// Setup the backend.
	m := &MySQLBackend{
		client:     db,
		table:      table,
		database:   database,
		statements: statements,
	}

	return m, nil
}

// Put is used to insert or update an entry.
func (m *MySQLBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"mysql", "put"}, time.Now())

	_, err := m.statements["put"].Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}

	return nil
}

// Get is used to fetch and entry.
func (m *MySQLBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"mysql", "get"}, time.Now())

	var result []byte

	err := m.statements["get"].QueryRow(key).Scan(&result)
	if err != nil {
		return nil, MySQLExecuteStmtFailure
	}

	ent := &Entry{
		Key:   key,
		Value: result,
	}

	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *MySQLBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"mysql", "delete"}, time.Now())

	_, err := m.statements["delete"].Exec(key)
	if err != nil {
		return err
	}

	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *MySQLBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"mysql", "list"}, time.Now())

	// Query to get all keys matching a prefix.
	list_query := "SELECT vault_key FROM " + m.database + "." + m.table + " WHERE vault_key LIKE '" + prefix + "%'"
	rows, err := m.client.Query(list_query)
	if err != nil {
		return nil, MySQLExecuteStmtFailure
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns")
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	keys := []string{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows")
		}

		for _, col := range values {
			keys = append(keys, string(col))
		}
	}

	sort.Strings(keys)

	return keys, nil
}
