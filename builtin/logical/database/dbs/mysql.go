package dbs

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/hashicorp/vault/helper/strutil"
)

const defaultRevocationSQL = `
	REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'%'; 
	DROP USER '{{name}}'@'%'
`

type MySQL struct {
	db *sql.DB

	ConnectionProducer
	CredentialsProducer
	sync.RWMutex
}

func (p *MySQL) Type() string {
	return postgreSQLTypeName
}

func (p *MySQL) getConnection() (*sql.DB, error) {
	db, err := p.Connection()
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

func (p *MySQL) CreateUser(createStmt, rollbackStmt, username, password, expiration string) error {
	// Get the connection
	db, err := p.getConnection()
	if err != nil {
		return err
	}

	// TODO: This is racey
	// Grab a read lock
	p.RLock()
	defer p.RUnlock()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute each query
	for _, query := range strutil.ParseArbitraryStringSlice(createStmt, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		stmt, err := tx.Prepare(queryHelper(query, map[string]string{
			"name":     username,
			"password": password,
		}))
		if err != nil {
			return err
		}
		defer stmt.Close()
		if _, err := stmt.Exec(); err != nil {
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// NOOP
func (p *MySQL) RenewUser(username, expiration string) error {
	return nil
}

func (p *MySQL) RevokeUser(username, revocationStmt string) error {
	// Get the connection
	db, err := p.getConnection()
	if err != nil {
		return err
	}

	// Grab the read lock
	p.RLock()
	defer p.RUnlock()

	// Use a default SQL statement for revocation if one cannot be fetched from the role

	if revocationStmt == "" {
		revocationStmt = defaultRevocationSQL
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, query := range strutil.ParseArbitraryStringSlice(revocationStmt, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		// This is not a prepared statement because not all commands are supported
		// 1295: This command is not supported in the prepared statement protocol yet
		// Reference https://mariadb.com/kb/en/mariadb/prepare-statement/
		query = strings.Replace(query, "{{name}}", username, -1)
		_, err = tx.Exec(query)
		if err != nil {
			return err
		}

	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
