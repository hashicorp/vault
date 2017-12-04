package mssql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const msSQLTypeName = "mssql"

// MSSQL is an implementation of Database interface
type MSSQL struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

func New() (interface{}, error) {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = msSQLTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 20,
		RoleNameLen:    20,
		UsernameLen:    128,
		Separator:      "-",
	}

	dbType := &MSSQL{
		ConnectionProducer:  connProducer,
		CredentialsProducer: credsProducer,
	}

	return dbType, nil
}

// Run instantiates a MSSQL object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(*MSSQL), apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (m *MSSQL) Type() (string, error) {
	return msSQLTypeName, nil
}

func (m *MSSQL) getConnection() (*sql.DB, error) {
	db, err := m.Connection()
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

// CreateUser generates the username/password on the underlying MSSQL secret backend as instructed by
// the CreationStatement provided.
func (m *MSSQL) CreateUser(statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	// Get the connection
	db, err := m.getConnection()
	if err != nil {
		return "", "", err
	}

	if statements.CreationStatements == "" {
		return "", "", dbutil.ErrEmptyCreationStatement
	}

	username, err = m.GenerateUsername(usernameConfig)
	if err != nil {
		return "", "", err
	}

	password, err = m.GeneratePassword()
	if err != nil {
		return "", "", err
	}

	expirationStr, err := m.GenerateExpiration(expiration)
	if err != nil {
		return "", "", err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	// Execute each query
	for _, query := range strutil.ParseArbitraryStringSlice(statements.CreationStatements, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		stmt, err := tx.Prepare(dbutil.QueryHelper(query, map[string]string{
			"name":       username,
			"password":   password,
			"expiration": expirationStr,
		}))
		if err != nil {
			return "", "", err
		}
		defer stmt.Close()
		if _, err := stmt.Exec(); err != nil {
			return "", "", err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return "", "", err
	}

	return username, password, nil
}

// RenewUser is not supported on MSSQL, so this is a no-op.
func (m *MSSQL) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser attempts to drop the specified user. It will first attempt to disable login,
// then kill pending connections from that user, and finally drop the user and login from the
// database instance.
func (m *MSSQL) RevokeUser(statements dbplugin.Statements, username string) error {
	if statements.RevocationStatements == "" {
		return m.revokeUserDefault(username)
	}

	// Get connection
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
	for _, query := range strutil.ParseArbitraryStringSlice(statements.RevocationStatements, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		stmt, err := tx.Prepare(dbutil.QueryHelper(query, map[string]string{
			"name": username,
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

func (m *MSSQL) revokeUserDefault(username string) error {
	// Get connection
	db, err := m.getConnection()
	if err != nil {
		return err
	}

	// First disable server login
	disableStmt, err := db.Prepare(fmt.Sprintf("ALTER LOGIN [%s] DISABLE;", username))
	if err != nil {
		return err
	}
	defer disableStmt.Close()
	if _, err := disableStmt.Exec(); err != nil {
		return err
	}

	// Query for sessions for the login so that we can kill any outstanding
	// sessions.  There cannot be any active sessions before we drop the logins
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	sessionStmt, err := db.Prepare(fmt.Sprintf(
		"SELECT session_id FROM sys.dm_exec_sessions WHERE login_name = '%s';", username))
	if err != nil {
		return err
	}
	defer sessionStmt.Close()

	sessionRows, err := sessionStmt.Query()
	if err != nil {
		return err
	}
	defer sessionRows.Close()

	var revokeStmts []string
	for sessionRows.Next() {
		var sessionID int
		err = sessionRows.Scan(&sessionID)
		if err != nil {
			return err
		}
		revokeStmts = append(revokeStmts, fmt.Sprintf("KILL %d;", sessionID))
	}

	// Query for database users using undocumented stored procedure for now since
	// it is the easiest way to get this information;
	// we need to drop the database users before we can drop the login and the role
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	stmt, err := db.Prepare(fmt.Sprintf("EXEC master.dbo.sp_msloginmappings '%s';", username))
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var loginName, dbName, qUsername string
		var aliasName sql.NullString
		err = rows.Scan(&loginName, &dbName, &qUsername, &aliasName)
		if err != nil {
			return err
		}
		revokeStmts = append(revokeStmts, fmt.Sprintf(dropUserSQL, dbName, username, username))
	}

	// we do not stop on error, as we want to remove as
	// many permissions as possible right now
	var lastStmtError error
	for _, query := range revokeStmts {
		stmt, err := db.Prepare(query)
		if err != nil {
			lastStmtError = err
			continue
		}
		defer stmt.Close()
		_, err = stmt.Exec()
		if err != nil {
			lastStmtError = err
		}
	}

	// can't drop if not all database users are dropped
	if rows.Err() != nil {
		return fmt.Errorf("cound not generate sql statements for all rows: %s", rows.Err())
	}
	if lastStmtError != nil {
		return fmt.Errorf("could not perform all sql statements: %s", lastStmtError)
	}

	// Drop this login
	stmt, err = db.Prepare(fmt.Sprintf(dropLoginSQL, username, username))
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

const dropUserSQL = `
USE [%s]
IF EXISTS
  (SELECT name
   FROM sys.database_principals
   WHERE name = N'%s')
BEGIN
  DROP USER [%s]
END
`

const dropLoginSQL = `
IF EXISTS
  (SELECT name
   FROM master.sys.server_principals
   WHERE name = N'%s')
BEGIN
  DROP LOGIN [%s]
END
`
