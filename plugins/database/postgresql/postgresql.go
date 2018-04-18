package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/dbtxn"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
	"github.com/lib/pq"
)

const (
	postgreSQLTypeName      = "postgres"
	defaultPostgresRenewSQL = `
ALTER ROLE "{{name}}" VALID UNTIL '{{expiration}}';
`
	defaultPostgresRotateRootCredentialsSQL = `
ALTER ROLE "{{username}}" WITH PASSWORD '{{password}}';
`
)

var _ dbplugin.Database = &PostgreSQL{}

// New implements builtinplugins.BuiltinFactory
func New() (interface{}, error) {
	db := new()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.SecretValues)
	return dbType, nil
}

func new() *PostgreSQL {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = postgreSQLTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 8,
		RoleNameLen:    8,
		UsernameLen:    63,
		Separator:      "-",
	}

	db := &PostgreSQL{
		SQLConnectionProducer: connProducer,
		CredentialsProducer:   credsProducer,
	}

	return db
}

// Run instantiates a PostgreSQL object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(dbplugin.Database), apiTLSConfig)

	return nil
}

type PostgreSQL struct {
	*connutil.SQLConnectionProducer
	credsutil.CredentialsProducer
}

func (p *PostgreSQL) Type() (string, error) {
	return postgreSQLTypeName, nil
}

func (p *PostgreSQL) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := p.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

func (p *PostgreSQL) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Creation) == 0 {
		return "", "", dbutil.ErrEmptyCreationStatement
	}

	// Grab the lock
	p.Lock()
	defer p.Unlock()

	username, err = p.GenerateUsername(usernameConfig)
	if err != nil {
		return "", "", err
	}

	password, err = p.GeneratePassword()
	if err != nil {
		return "", "", err
	}

	expirationStr, err := p.GenerateExpiration(expiration)
	if err != nil {
		return "", "", err
	}

	// Get the connection
	db, err := p.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", err

	}
	defer func() {
		tx.Rollback()
	}()
	// Return the secret

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

func (p *PostgreSQL) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	p.Lock()
	defer p.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	renewStmts := statements.Renewal
	if len(renewStmts) == 0 {
		renewStmts = []string{defaultPostgresRenewSQL}
	}

	db, err := p.getConnection(ctx)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	expirationStr, err := p.GenerateExpiration(expiration)
	if err != nil {
		return err
	}

	for _, stmt := range renewStmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       username,
				"expiration": expirationStr,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (p *PostgreSQL) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the lock
	p.Lock()
	defer p.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Revocation) == 0 {
		return p.defaultRevokeUser(ctx, username)
	}

	return p.customRevokeUser(ctx, username, statements.Revocation)
}

func (p *PostgreSQL) customRevokeUser(ctx context.Context, username string, revocationStmts []string) error {
	db, err := p.getConnection(ctx)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	for _, stmt := range revocationStmts {
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

	return tx.Commit()
}

func (p *PostgreSQL) defaultRevokeUser(ctx context.Context, username string) error {
	db, err := p.getConnection(ctx)
	if err != nil {
		return err
	}

	// Check if the role exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT exists (SELECT rolname FROM pg_roles WHERE rolname=$1);", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if exists == false {
		return nil
	}

	// Query for permissions; we need to revoke permissions before we can drop
	// the role
	// This isn't done in a transaction because even if we fail along the way,
	// we want to remove as much access as possible
	stmt, err := db.PrepareContext(ctx, "SELECT DISTINCT table_schema FROM information_schema.role_column_grants WHERE grantee=$1;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return err
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
			pq.QuoteIdentifier(schema),
			pq.QuoteIdentifier(username)))

		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE USAGE ON SCHEMA %s FROM %s;`,
			pq.QuoteIdentifier(schema),
			pq.QuoteIdentifier(username)))
	}

	// for good measure, revoke all privileges and usage on schema public
	revocationStmts = append(revocationStmts, fmt.Sprintf(
		`REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM %s;`,
		pq.QuoteIdentifier(username)))

	revocationStmts = append(revocationStmts, fmt.Sprintf(
		"REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM %s;",
		pq.QuoteIdentifier(username)))

	revocationStmts = append(revocationStmts, fmt.Sprintf(
		"REVOKE USAGE ON SCHEMA public FROM %s;",
		pq.QuoteIdentifier(username)))

	// get the current database name so we can issue a REVOKE CONNECT for
	// this username
	var dbname sql.NullString
	if err := db.QueryRowContext(ctx, "SELECT current_database();").Scan(&dbname); err != nil {
		return err
	}

	if dbname.Valid {
		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE CONNECT ON DATABASE %s FROM %s;`,
			pq.QuoteIdentifier(dbname.String),
			pq.QuoteIdentifier(username)))
	}

	// again, here, we do not stop on error, as we want to remove as
	// many permissions as possible right now
	var lastStmtError error
	for _, query := range revocationStmts {
		if err := dbtxn.ExecuteDBQuery(ctx, db, nil, query); err != nil {
			lastStmtError = err
		}
	}

	// can't drop if not all privileges are revoked
	if rows.Err() != nil {
		return errwrap.Wrapf("could not generate revocation statements for all rows: {{err}}", rows.Err())
	}
	if lastStmtError != nil {
		return errwrap.Wrapf("could not perform all revocation statements: {{err}}", lastStmtError)
	}

	// Drop this user
	stmt, err = db.PrepareContext(ctx, fmt.Sprintf(
		`DROP ROLE IF EXISTS %s;`, pq.QuoteIdentifier(username)))
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	p.Lock()
	defer p.Unlock()

	if len(p.Username) == 0 || len(p.Password) == 0 {
		return nil, errors.New("username and password are required to rotate")
	}

	rotateStatents := statements
	if len(rotateStatents) == 0 {
		rotateStatents = []string{defaultPostgresRotateRootCredentialsSQL}
	}

	db, err := p.getConnection(ctx)
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

	password, err := p.GeneratePassword()
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
				"username": p.Username,
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

	// Close the database connection to ensure no new connections come in
	if err := db.Close(); err != nil {
		return nil, err
	}

	p.RawConfig["password"] = password
	return p.RawConfig, nil
}
