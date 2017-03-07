package dbs

import (
	"database/sql"
	"strings"

	"github.com/hashicorp/vault/helper/strutil"
)

const defaultRevocationStmts = `
	REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'%'; 
	DROP USER '{{name}}'@'%'
`

type MySQL struct {
	ConnectionProducer
	CredentialsProducer
}

func (m *MySQL) Type() string {
	return mySQLTypeName
}

func (m *MySQL) getConnection() (*sql.DB, error) {
	db, err := m.connection()
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

func (m *MySQL) CreateUser(createStmts, rollbackStmts, username, password, expiration string) error {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	// Get the connection
	db, err := m.getConnection()
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute each query
	for _, query := range strutil.ParseArbitraryStringSlice(createStmts, ";") {
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
func (m *MySQL) RenewUser(username, expiration string) error {
	return nil
}

func (m *MySQL) RevokeUser(username, revocationStmts string) error {
	// Grab the read lock
	m.Lock()
	defer m.Unlock()

	// Get the connection
	db, err := m.getConnection()
	if err != nil {
		return err
	}

	// Use a default SQL statement for revocation if one cannot be fetched from the role
	if revocationStmts == "" {
		revocationStmts = defaultRevocationStmts
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, query := range strutil.ParseArbitraryStringSlice(revocationStmts, ";") {
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
