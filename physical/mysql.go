package physical

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	mysql "github.com/go-sql-driver/mysql"
)

// Unreserved tls key
// Reserved values are "true", "false", "skip-verify"
const mysqlTLSKey = "default"

// MySQLBackend is a physical backend that stores data
// within MySQL database.
type MySQLBackend struct {
	dbTable    string
	client     *sql.DB
	statements map[string]*sql.Stmt
	logger     *log.Logger
}

// newMySQLBackend constructs a MySQL backend using the given API client and
// server address and credential for accessing mysql database.
func newMySQLBackend(conf map[string]string, logger *log.Logger) (Backend, error) {
	// Get the MySQL credentials to perform read/write operations.
	username, ok := conf["username"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing username")
	}
	password, ok := conf["password"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get or set MySQL server address. Defaults to localhost and default port(3306)
	address, ok := conf["address"]
	if !ok {
		address = "127.0.0.1:3306"
	}

	// Get the MySQL database and table details.
	database, ok := conf["database"]
	if !ok {
		database = "vault"
	}
	table, ok := conf["table"]
	if !ok {
		table = "vault"
	}
	dbTable := database + "." + table

	dsnParams := url.Values{}
	tlsCaFile, ok := conf["tls_ca_file"]
	if ok {
		if err := setupMySQLTLSConfig(tlsCaFile); err != nil {
			return nil, fmt.Errorf("failed register TLS config: %v", err)
		}

		dsnParams.Add("tls", mysqlTLSKey)
	}

	// Create MySQL handle for the database.
	dsn := username + ":" + password + "@tcp(" + address + ")/?" + dsnParams.Encode()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %v", err)
	}

	// Create the required database if it doesn't exists.
	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		return nil, fmt.Errorf("failed to create mysql database: %v", err)
	}

	// Create the required table if it doesn't exists.
	create_query := "CREATE TABLE IF NOT EXISTS " + dbTable +
		" (vault_key varbinary(512), vault_value mediumblob, PRIMARY KEY (vault_key))"
	if _, err := db.Exec(create_query); err != nil {
		return nil, fmt.Errorf("failed to create mysql table: %v", err)
	}

	// Setup the backend.
	m := &MySQLBackend{
		dbTable:    dbTable,
		client:     db,
		statements: make(map[string]*sql.Stmt),
		logger:     logger,
	}

	// Prepare all the statements required
	statements := map[string]string{
		"put": "INSERT INTO " + dbTable +
			" VALUES( ?, ? ) ON DUPLICATE KEY UPDATE vault_value=VALUES(vault_value)",
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
func (m *MySQLBackend) prepare(name, query string) error {
	stmt, err := m.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare '%s': %v", name, err)
	}
	m.statements[name] = stmt
	return nil
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

// Establish a TLS connection with a given CA certificate
// Register a tsl.Config associted with the same key as the dns param from sql.Open
// foo:bar@tcp(127.0.0.1:3306)/dbname?tls=default
func setupMySQLTLSConfig(tlsCaFile string) error {
	rootCertPool := x509.NewCertPool()

	pem, err := ioutil.ReadFile(tlsCaFile)
	if err != nil {
		return err
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return err
	}

	err = mysql.RegisterTLSConfig(mysqlTLSKey, &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		return err
	}

	return nil
}
