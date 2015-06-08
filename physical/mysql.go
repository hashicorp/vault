package physical

import (
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	_ "github.com/go-sql-driver/mysql"
)

var (
	MySQLDBNameMissing          = errors.New("database name is missing in the configuration")
	MySQLTableNameMissing       = errors.New("table name is missing in the configuration")
	MySQLHandlerCreationFailure = errors.New("failed to open handler with database")
	MySQLPrepareStmtFailure     = errors.New("failed to prepare statement")
	MySQLExecuteStmtFailure     = errors.New("failed to execute statement")
	MySQLGetColumnsFailure      = errors.New("failed to get columns")
	MySQLScanRowsFailure        = errors.New("failed to scan rows")
)

// MySQLBackend is a physical backend that stores data
// within MySQL database.
type MySQLBackend struct {
	table    string
	database string
	client   *sql.DB
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
		return nil, MySQLDBNameMissing
	}
	table, ok := conf["table"]
	if !ok {
		return nil, MySQLTableNameMissing
	}

	// Create MySQL handle for the database.
	dsn := username + ":" + password + "@tcp(" + address + ")/" + database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, MySQLHandlerCreationFailure
	}
	defer db.Close()

	// Create the required table.
	create_stmt := "CREATE TABLE IF NOT EXISTS " + database + "." + table + " (vault_key varchar(255), vault_value varchar(255), PRIMARY KEY (vault_key))"
	stmt, err := db.Prepare(create_stmt)
	if err != nil {
		return nil, MySQLPrepareStmtFailure
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil, MySQLExecuteStmtFailure
	}

	// Setup the backend.
	m := &MySQLBackend{
		client:   db,
		table:    table,
		database: database,
	}

	return m, nil
}

// Put is used to insert or update an entry.
func (m *MySQLBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"mysql", "put"}, time.Now())

	insert_stmt := "INSERT INTO " + m.database + "." + m.table + " VALUES( ?, ? ) ON DUPLICATE KEY UPDATE vault_value=VALUES(vault_value)"
	stmt, err := m.client.Prepare(insert_stmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}

	return nil
}

// Get is used to fetch and entry.
func (m *MySQLBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"mysql", "get"}, time.Now())

	select_stmt := "SELECT vault_value FROM " + m.database + "." + m.table + " WHERE vault_key = ?"
	stmt, err := m.client.Prepare(select_stmt)
	if err != nil {
		return nil, MySQLPrepareStmtFailure
	}
	defer stmt.Close()

	var result []byte

	err = stmt.QueryRow(key).Scan(&result)
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

	delete_stmt := "DELETE FROM " + m.database + "." + m.table + "WHERE vault_key = ?"
	stmt, err := m.client.Prepare(delete_stmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}

	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *MySQLBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"mysql", "list"}, time.Now())

	list_stmt := "SELECT vault_key FROM " + m.database + "." + m.table
	rows, err := m.client.Query(list_stmt)
	if err != nil {
		return nil, MySQLExecuteStmtFailure
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, MySQLGetColumnsFailure
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
			return nil, MySQLScanRowsFailure
		}

		for _, col := range values {
			if strings.HasPrefix(string(col), prefix) {
				keys = append(keys, string(col))
			}
		}
	}

	sort.Strings(keys)

	return keys, nil
}
