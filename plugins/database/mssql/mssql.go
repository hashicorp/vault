package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/dbtxn"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const msSQLTypeName = "mssql"

var _ dbplugin.Database = &MSSQL{}

// MSSQL is an implementation of Database interface
type MSSQL struct {
	*connutil.SQLConnectionProducer
	credsutil.CredentialsProducer
}

func New() (interface{}, error) {
	db := new()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.SecretValues)

	return dbType, nil
}

func new() *MSSQL {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = msSQLTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 20,
		RoleNameLen:    20,
		UsernameLen:    128,
		Separator:      "-",
	}

	return &MSSQL{
		SQLConnectionProducer: connProducer,
		CredentialsProducer:   credsProducer,
	}
}

// Run instantiates a MSSQL object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(dbplugin.Database), apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (m *MSSQL) Type() (string, error) {
	return msSQLTypeName, nil
}

func (m *MSSQL) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

// CreateUser generates the username/password on the underlying MSSQL secret backend as instructed by
// the CreationStatement provided.
func (m *MSSQL) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	// Get the connection
	db, err := m.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	if len(statements.Creation) == 0 {
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
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	// Execute each query
	for _, stmt := range statements.Creation {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       username,
				"password":   password,
				"expiration": expirationStr,
			}

			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return "", "", err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return "", "", err
	}

	return username, password, nil
}

// RenewUser is not supported on MSSQL, so this is a no-op.
func (m *MSSQL) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser attempts to drop the specified user. It will first attempt to disable login,
// then kill pending connections from that user, and finally drop the user and login from the
// database instance.
func (m *MSSQL) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Revocation) == 0 {
		return m.revokeUserDefault(ctx, username)
	}

	// Get connection
	db, err := m.getConnection(ctx)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute each query
	for _, stmt := range statements.Revocation {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name": username,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *MSSQL) revokeUserDefault(ctx context.Context, username string) error {
	// Get connection
	db, err := m.getConnection(ctx)
	if err != nil {
		return err
	}

	// First disable server login
	disableStmt, err := db.PrepareContext(ctx, fmt.Sprintf("ALTER LOGIN [%s] DISABLE;", username))
	if err != nil {
		return err
	}
	defer disableStmt.Close()
	if _, err := disableStmt.ExecContext(ctx); err != nil {
		return err
	}

	// Query for sessions for the login so that we can kill any outstanding
	// sessions.  There cannot be any active sessions before we drop the logins
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	sessionStmt, err := db.PrepareContext(ctx, fmt.Sprintf(
		"SELECT session_id FROM sys.dm_exec_sessions WHERE login_name = '%s';", username))
	if err != nil {
		return err
	}
	defer sessionStmt.Close()

	sessionRows, err := sessionStmt.QueryContext(ctx)
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
	stmt, err := db.PrepareContext(ctx, fmt.Sprintf("EXEC master.dbo.sp_msloginmappings '%s';", username))
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
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
		if err := dbtxn.ExecuteDBQuery(ctx, db, nil, query); err != nil {
			lastStmtError = err
		}
	}

	// can't drop if not all database users are dropped
	if rows.Err() != nil {
		return errwrap.Wrapf("could not generate sql statements for all rows: {{err}}", rows.Err())
	}
	if lastStmtError != nil {
		return errwrap.Wrapf("could not perform all sql statements: {{err}}", lastStmtError)
	}

	// Drop this login
	stmt, err = db.PrepareContext(ctx, fmt.Sprintf(dropLoginSQL, username, username))
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (m *MSSQL) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	m.Lock()
	defer m.Unlock()

	if len(m.Username) == 0 || len(m.Password) == 0 {
		return nil, errors.New("username and password are required to rotate")
	}

	rotateStatents := statements
	if len(rotateStatents) == 0 {
		rotateStatents = []string{rotateRootCredentialsSQL}
	}

	db, err := m.getConnection(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		tx.Rollback()
	}()

	password, err := m.GeneratePassword()
	if err != nil {
		return nil, err
	}

	for _, stmt := range rotateStatents {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": m.Username,
				"password": password,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	m.RawConfig["password"] = password
	return m.RawConfig, nil
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

const rotateRootCredentialsSQL = `
ALTER LOGIN [{{username}}] WITH PASSWORD = '{{password}}' 
`
