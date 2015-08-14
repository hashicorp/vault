package physical

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"
)

// PostGreSQLBackend is a physical backend that stores data on postgres database
type PostGreSQLBackend struct {
	dbTable    string
	db         *sql.DB
	statements map[string]*sql.Stmt
}

// newPostGreSQLBackend constructs a new postgresql backend, opens a pool of connection and prepares sql statements, creates table is not already there
func newPostGreSQLBackend(conf map[string]string) (Backend, error) {

	url, ok := conf["url"]
	if !ok {
		return nil, fmt.Errorf("'url' must be set")
	}

	table_name, ok := conf["table_name"]
	if !ok {
		table_name = "vault"
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("Error opening connection to postgresql database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error running ping to postgresql databases: %v", err)
	}

	sqlStmt = "CREATE TABLE if not exists " + table_name + " ( key TEXT not null, value bytea, created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL, updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("Error creating database table: %v", err)
	}

	b = PostGreSQLBackend{
		dbTable:    table_name,
		db:         db,
		statements: make(map[string]*sql.Stmt),
	}, nil

	statements := map[string]string{
		"put_update": "update " + table_name + " set value =  $1, updated_at = $2 where key = $3",
		"put_insert": "insert into  " + table_name + " (key, value, created_at, updated_at) values ($1, $2, $3, $4)",
		"list":       "SELECT key FROM " + table_name + " WHERE key like ?",
		"delete":     "DELETE FROM " + table_name + " WHERE key = ?",
		"get":        "SELECT value FROM " + table_name + " WHERE key = ?",
	}
	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
			return nil, err
		}
	}

	return &b
}

func (b *PostGreSQLBackend) prepare(name, query string) error {
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare '%s': %v", name, err)
	}
	b.statements[name] = stmt
	return nil
}

// Delete - Delete rows from table which match key
func (b *PostGreSQLBackend) Delete(k string) error {
	txn, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %v", err)
	}

	_, err = statements["delete"].Exec(k)
	if err != nil {
		return fmt.Errorf("Error executing delete satement: %v", err)
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("Error committing transaction: %v", err)
	}

	return err
}

// Get - fetch data from tables that match key
func (b *PostGreSQLBackend) Get(k string) (*Entry, error) {
	var value []byte
	row, err := statements["get"].Query(k)
	if err != nil {
		return nil, fmt.Errorf("Error committing transaction: %v", err)
	}

	if row.Next() {
		row.Scan(&value)
		entry := Entry{k, value}

		return &entry, nil
	}
	return nil, err

}

// Put - Update a row in database
func (b *PostGreSQLBackend) Put(entry *Entry) error {

	// TODO: fix me
	rows, err := statements["get"].Query(entry.Key)
	if err != nil {
		return fmt.Errorf("Error committing transaction: %v", err)
	}
	defer rows.Close()

	txn, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %v", err)
	}
	next := rows.Next()

	// need to update if already there
	time := time.Now()
	if next {
		_, err = statements["put_update"].Exec(entry.Value, time, entry.Key)
		if err != nil {
			return fmt.Errorf("Error executing update: %v", err)
		}

	} else {
		_, err = statements["put_insert"].Exec(entry.Key, entry.Value, time, time)
		if err != nil {
			return fmt.Errorf("Error executing insert on put: %v", err)
		}

	}
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("Error committing transaction: %v", err)
	}

	return err
}

// List - query database for matches
func (b *PostGreSQLBackend) List(prefix string) ([]string, error) {

	buffer := bytes.NewBufferString(prefix)
	buffer.WriteString("%")
	query := buffer.String()

	// TODO: fix me
	rows, err := statements["list"].Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error querying during list: %v", err)
	}
	defer rows.Close()

	result := make([]string, 0)

	for rows.Next() {
		var message string
		rows.Scan(&message)
		result = append(result, message)
	}

	return result, nil
}
