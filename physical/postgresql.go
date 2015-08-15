package physical

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// PostgresqlBackend is a physical backend that stores data on postgres database
type PostgresqlBackend struct {
	dbTable    string
	client     *sql.DB
	statements map[string]*sql.Stmt
}

// newPostgresqlBackend constructs a new postgresql backend, opens a pool of connection and prepares sql statements, creates table is not already there
func newPostgresqlBackend(conf map[string]string) (Backend, error) {

	url, ok := conf["url"]
	if !ok {
		return nil, fmt.Errorf("'url' must be set")
	}

	tableName, ok := conf["table_name"]
	if !ok {
		tableName = "vault"
	}

	fmt.Println("ERROR::::::::::::- new backend")

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("Error opening connection to postgresql database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error running ping to postgresql databases: %v", err)
	}

	sqlStmt := "CREATE TABLE if not exists " + tableName + " ( key TEXT not null PRIMARY KEY, value bytea, created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL default (now() at time zone 'utc'), updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL default (now() at time zone 'utc'));"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("Error creating database table: %v", err)
	}

	b := &PostgresqlBackend{
		dbTable:    tableName,
		client:     db,
		statements: make(map[string]*sql.Stmt),
	}

	statements := map[string]string{
		"put_update": "update " + tableName + " set value =  $1, updated_at = $2 where key = $3",
		"put_insert": "insert into  " + tableName + " (key, value, created_at, updated_at) values ($1, $2, $3, $4)",
		"list":       "SELECT key FROM " + tableName + " WHERE key like $1",
		"delete":     "DELETE FROM " + tableName + " WHERE key = $1",
		"get":        "SELECT value FROM " + tableName + " WHERE key = $1",
	}
	for name, query := range statements {
		if err := b.prepare(name, query); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (b *PostgresqlBackend) prepare(name, query string) error {
	stmt, err := b.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("Prepare: failed to prepare '%s': %v", name, err)
	}
	b.statements[name] = stmt
	return nil
}

// Delete - Delete rows from table which match key
func (b *PostgresqlBackend) Delete(k string) error {
	txn, err := b.client.Begin()
	if err != nil {
		return fmt.Errorf("DELETE: Error starting transaction: %v", err)
	}

	_, err = b.statements["delete"].Exec(k)
	if err != nil {
		return fmt.Errorf("DELETE: Error executing delete satement: %v", err)
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("DELETE: Error committing transaction: %v", err)
	}

	return err
}

// Get - fetch data from tables that match key
func (b *PostgresqlBackend) Get(k string) (*Entry, error) {
	var result []byte

	err := b.statements["get"].QueryRow(k).Scan(&result)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	ent := &Entry{
		Key:   k,
		Value: result,
	}

	return ent, nil
}

// Put - Update a row in database
func (b *PostgresqlBackend) Put(entry *Entry) error {

	// TODO: fix me
	rows, err := b.statements["get"].Query(entry.Key)
	if err != nil {
		return fmt.Errorf("PUT: Error running query: %v", err)
	}
	defer rows.Close()

	txn, err := b.client.Begin()
	if err != nil {
		return fmt.Errorf("PUT: Error starting transaction: %v", err)
	}
	next := rows.Next()

	// need to update if already there
	time := time.Now().UTC()
	if next {
		_, err = b.statements["put_update"].Exec(entry.Value, time, entry.Key)
		if err != nil {
			return fmt.Errorf("PUT: Error executing update: %v", err)
		}

	} else {
		_, err = b.statements["put_insert"].Exec(entry.Key, entry.Value, time, time)
		if err != nil {
			return fmt.Errorf("PUT: Error executing insert on put: %v", err)
		}

	}
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("PUT: Error committing transaction: %v", err)
	}

	return err
}

// List - query database for matches
func (b *PostgresqlBackend) List(prefix string) ([]string, error) {

	buffer := bytes.NewBufferString(prefix)
	buffer.WriteString("%")
	query := buffer.String()

	rows, err := b.statements["list"].Query(query)
	if err != nil {
		return nil, fmt.Errorf("List: Error querying during list: %v", err)
	}
	defer rows.Close()

	result := make([]string, 0)

	for rows.Next() {
		var message string
		err = rows.Scan(&message)
		if err != nil {
			return nil, fmt.Errorf("List: failed to scan rows: %v", err)
		}

		message = strings.TrimPrefix(message, prefix)
		if i := strings.Index(message, "/"); i == -1 {
			// Add objects only from the current 'folder'
			result = append(result, message)
		} else if i != -1 {
			// Add truncated 'folder' paths
			result = appendIfMissing(result, string(message[:i+1]))
		}
	}
	sort.Strings(result)

	return result, nil
}
