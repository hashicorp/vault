// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package redshift

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/template"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	// This is how this plugin will be reflected in middleware
	// such as metrics.
	middlewareTypeName = "redshift"

	// This allows us to use the postgres database driver.
	sqlTypeName = "pgx"

	defaultRenewSQL = `
ALTER USER "{{name}}" VALID UNTIL '{{expiration}}';
`
	defaultRotateRootCredentialsSQL = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
`
	defaultUserNameTemplate = `{{ printf "v-%s-%s-%s-%s" (.DisplayName | truncate 8) (.RoleName | truncate 8) (random 20) (unix_time) | truncate 63 | lowercase }}`
)

var _ dbplugin.Database = (*RedShift)(nil)

// New implements builtinplugins.BuiltinFactory
// Redshift implements (mostly) a postgres 8 interface, and part of that is
// under the hood, it's lower-casing the usernames.
func New() (interface{}, error) {
	db := newRedshift()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func newRedshift() *RedShift {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = sqlTypeName

	db := &RedShift{
		SQLConnectionProducer: connProducer,
	}

	return db
}

type RedShift struct {
	*connutil.SQLConnectionProducer

	usernameProducer template.StringTemplate
}

func (r *RedShift) secretValues() map[string]string {
	return map[string]string{
		r.Password: "[password]",
	}
}

func (r *RedShift) Type() (string, error) {
	return middlewareTypeName, nil
}

// Initialize must be called on each new RedShift struct before use.
// It uses the connutil.SQLConnectionProducer's Init function to do all the lifting.
func (r *RedShift) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	conf, err := r.Init(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("error initializing db: %w", err)
	}

	usernameTemplate, err := strutil.GetString(req.Config, "username_template")
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("failed to retrieve username_template: %w", err)
	}
	if usernameTemplate == "" {
		usernameTemplate = defaultUserNameTemplate
	}

	up, err := template.NewTemplate(template.Template(usernameTemplate))
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("unable to initialize username template: %w", err)
	}
	r.usernameProducer = up

	_, err = r.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	return dbplugin.InitializeResponse{
		Config: conf,
	}, nil
}

// getConnection accepts a context and returns a new pointer to a sql.DB object.
// It's up to the caller to close the connection or handle reuse logic.
func (r *RedShift) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := r.Connection(ctx)
	if err != nil {
		return nil, err
	}
	return db.(*sql.DB), nil
}

// NewUser creates a new user in the database. There is no default statement for
// creating users, so one must be specified in the plugin config.
// Generated usernames are of the form v-{display-name}-{role-name}-{UUID}-{timestamp}
func (r *RedShift) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	if len(req.Statements.Commands) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	// Grab the lock
	r.Lock()
	defer r.Unlock()

	username, err := r.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	password := req.Password
	expirationStr := req.Expiration.Format("2006-01-02 15:04:05-0700")

	// Get the connection
	db, err := r.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	defer db.Close()

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	defer func() {
		tx.Rollback()
	}()

	// Execute each query
	for _, stmt := range req.Statements.Commands {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       username,
				"username":   username,
				"password":   password,
				"expiration": expirationStr,
			}
			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return dbplugin.NewUserResponse{}, err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	return dbplugin.NewUserResponse{
		Username: username,
	}, nil
}

// UpdateUser can update the expiration or the password of a user, or both.
// The updates all happen in a single transaction, so they will either all
// succeed or all fail.
// Both updates support both default and custom statements.
func (r *RedShift) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, errors.New("no changes requested")
	}

	r.Lock()
	defer r.Unlock()

	db, err := r.getConnection(ctx)
	if err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}
	defer func() {
		tx.Rollback()
	}()

	if req.Expiration != nil {
		err = updateUserExpiration(ctx, req, tx)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, err
		}
	}

	if req.Password != nil {
		err = updateUserPassword(ctx, req, tx)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, err
		}
	}

	err = tx.Commit()
	return dbplugin.UpdateUserResponse{}, err
}

func updateUserExpiration(ctx context.Context, req dbplugin.UpdateUserRequest, tx *sql.Tx) error {
	if req.Username == "" {
		return errors.New("must provide a username to update user expiration")
	}
	renewStmts := req.Expiration.Statements
	if len(renewStmts.Commands) == 0 {
		renewStmts.Commands = []string{defaultRenewSQL}
	}

	expirationStr := req.Expiration.NewExpiration.Format("2006-01-02 15:04:05-0700")

	for _, stmt := range renewStmts.Commands {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       req.Username,
				"username":   req.Username,
				"expiration": expirationStr,
			}
			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return err
			}
		}
	}

	return nil
}

func updateUserPassword(ctx context.Context, req dbplugin.UpdateUserRequest, tx *sql.Tx) error {
	username := req.Username
	password := req.Password.NewPassword
	if username == "" || password == "" {
		return errors.New("must provide both username and a new password to update user password")
	}

	// Check if the role exists
	var exists bool
	err := tx.QueryRowContext(ctx, "SELECT exists (SELECT usename FROM pg_user WHERE usename=$1);", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		// Server error
		return err
	}
	if err == sql.ErrNoRows || !exists {
		// Most likely a user error
		return fmt.Errorf("cannot update password for username %q because it does not exist", username)
	}

	// Vault requires the database user already exist, and that the credentials
	// used to execute the rotation statements has sufficient privileges.
	statements := req.Password.Statements.Commands
	if len(statements) == 0 {
		statements = []string{defaultRotateRootCredentialsSQL}
	}
	// Execute each query
	for _, stmt := range statements {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":     username,
				"username": username,
				"password": password,
			}
			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteUser supports both default and custom statements to delete a user.
func (r *RedShift) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	// Grab the lock
	r.Lock()
	defer r.Unlock()

	if len(req.Statements.Commands) == 0 {
		return r.defaultDeleteUser(ctx, req)
	}

	return r.customDeleteUser(ctx, req)
}

func (r *RedShift) customDeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	db, err := r.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer func() {
		tx.Rollback()
	}()

	for _, stmt := range req.Statements.Commands {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":     req.Username,
				"username": req.Username,
			}
			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return dbplugin.DeleteUserResponse{}, err
			}
		}
	}

	return dbplugin.DeleteUserResponse{}, tx.Commit()
}

func (r *RedShift) defaultDeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	db, err := r.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer db.Close()

	username := req.Username

	// Check if the role exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT exists (SELECT usename FROM pg_user WHERE usename=$1);", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return dbplugin.DeleteUserResponse{}, err
	}

	if !exists {
		// No error as Redshift may have deleted the user via TTL before we got to it.
		return dbplugin.DeleteUserResponse{}, nil
	}

	// Query for permissions; we need to revoke permissions before we can drop
	// the role
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	stmt, err := db.PrepareContext(ctx, "SELECT DISTINCT table_schema FROM information_schema.role_column_grants WHERE grantee=$1;")
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer rows.Close()

	const initialNumRevocations = 16
	revocationStmts := make([]string, 0, initialNumRevocations)
	for rows.Next() {
		var schema string
		err = rows.Scan(&schema)
		if err != nil {
			// keep going; remove as many permissions as possible right now
			continue
		}
		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA %s FROM %s;`,
			dbutil.QuoteIdentifier(schema),
			dbutil.QuoteIdentifier(username)))

		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE USAGE ON SCHEMA %s FROM %s;`,
			dbutil.QuoteIdentifier(schema),
			dbutil.QuoteIdentifier(username)))
	}

	// for good measure, revoke all privileges and usage on schema public
	revocationStmts = append(revocationStmts, fmt.Sprintf(
		`REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM %s;`,
		dbutil.QuoteIdentifier(username)))

	revocationStmts = append(revocationStmts, fmt.Sprintf(
		"REVOKE USAGE ON SCHEMA public FROM %s;",
		dbutil.QuoteIdentifier(username)))

	// get the current database name so we can issue a REVOKE CONNECT for
	// this username
	var dbname sql.NullString
	if err := db.QueryRowContext(ctx, "SELECT current_database();").Scan(&dbname); err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	if dbname.Valid {
		/*
			We create this stored procedure to ensure we can durably revoke users on Redshift. We do not
			clean up since that can cause race conditions with other instances of Vault attempting to use
			this SP at the same time.
		*/
		revocationStmts = append(revocationStmts, `CREATE OR REPLACE PROCEDURE terminateloop(dbusername varchar(100))
LANGUAGE plpgsql
AS $$
DECLARE
  currentpid int;
  loopvar int;
  qtyconns int;
BEGIN
SELECT COUNT(process) INTO qtyconns FROM stv_sessions WHERE user_name=dbusername;
  FOR loopvar IN 1..qtyconns LOOP
    SELECT INTO currentpid process FROM stv_sessions WHERE user_name=dbusername ORDER BY process ASC LIMIT 1;
    SELECT pg_terminate_backend(currentpid);
  END LOOP;
END
$$;`)

		revocationStmts = append(revocationStmts, fmt.Sprintf(`call terminateloop('%s');`, username))
	}

	// again, here, we do not stop on error, as we want to remove as
	// many permissions as possible right now
	var lastStmtError *multierror.Error // error
	for _, query := range revocationStmts {
		if err := dbtxn.ExecuteDBQueryDirect(ctx, db, nil, query); err != nil {
			lastStmtError = multierror.Append(lastStmtError, err)
		}
	}

	// can't drop if not all privileges are revoked
	if rows.Err() != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("could not generate revocation statements for all rows: %w", rows.Err())
	}
	if lastStmtError != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("could not perform all revocation statements: %w", lastStmtError)
	}

	// Drop this user
	stmt, err = db.PrepareContext(ctx, fmt.Sprintf(
		`DROP USER IF EXISTS %s;`, dbutil.QuoteIdentifier(username)))
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx); err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	return dbplugin.DeleteUserResponse{}, nil
}
