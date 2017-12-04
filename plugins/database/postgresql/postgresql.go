package postgresql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	postgreSQLTypeName      string = "postgres"
	defaultPostgresRenewSQL        = `
ALTER ROLE "{{name}}" VALID UNTIL '{{expiration}}';
`
)

// New implements builtinplugins.BuiltinFactory
func New() (interface{}, error) {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = postgreSQLTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 8,
		RoleNameLen:    8,
		UsernameLen:    63,
		Separator:      "-",
	}

	dbType := &PostgreSQL{
		ConnectionProducer:  connProducer,
		CredentialsProducer: credsProducer,
	}

	return dbType, nil
}

// Run instantiates a PostgreSQL object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(*PostgreSQL), apiTLSConfig)

	return nil
}

type PostgreSQL struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

func (p *PostgreSQL) Type() (string, error) {
	return postgreSQLTypeName, nil
}

func (p *PostgreSQL) getConnection() (*sql.DB, error) {
	db, err := p.Connection()
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

func (p *PostgreSQL) CreateUser(statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	if statements.CreationStatements == "" {
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
	db, err := p.getConnection()
	if err != nil {
		return "", "", err

	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return "", "", err

	}
	defer func() {
		tx.Rollback()
	}()
	// Return the secret

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

func (p *PostgreSQL) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	p.Lock()
	defer p.Unlock()

	renewStmts := statements.RenewStatements
	if renewStmts == "" {
		renewStmts = defaultPostgresRenewSQL
	}

	db, err := p.getConnection()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
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

	for _, query := range strutil.ParseArbitraryStringSlice(renewStmts, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}
		stmt, err := tx.Prepare(dbutil.QueryHelper(query, map[string]string{
			"name":       username,
			"expiration": expirationStr,
		}))
		if err != nil {
			return err
		}

		defer stmt.Close()
		if _, err := stmt.Exec(); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) RevokeUser(statements dbplugin.Statements, username string) error {
	// Grab the lock
	p.Lock()
	defer p.Unlock()

	if statements.RevocationStatements == "" {
		return p.defaultRevokeUser(username)
	}

	return p.customRevokeUser(username, statements.RevocationStatements)
}

func (p *PostgreSQL) customRevokeUser(username, revocationStmts string) error {
	db, err := p.getConnection()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	for _, query := range strutil.ParseArbitraryStringSlice(revocationStmts, ";") {
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

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) defaultRevokeUser(username string) error {
	db, err := p.getConnection()
	if err != nil {
		return err
	}

	// Check if the role exists
	var exists bool
	err = db.QueryRow("SELECT exists (SELECT rolname FROM pg_roles WHERE rolname=$1);", username).Scan(&exists)
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
	stmt, err := db.Prepare("SELECT DISTINCT table_schema FROM information_schema.role_column_grants WHERE grantee=$1;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
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
	if err := db.QueryRow("SELECT current_database();").Scan(&dbname); err != nil {
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

	// can't drop if not all privileges are revoked
	if rows.Err() != nil {
		return fmt.Errorf("could not generate revocation statements for all rows: %s", rows.Err())
	}
	if lastStmtError != nil {
		return fmt.Errorf("could not perform all revocation statements: %s", lastStmtError)
	}

	// Drop this user
	stmt, err = db.Prepare(fmt.Sprintf(
		`DROP ROLE IF EXISTS %s;`, pq.QuoteIdentifier(username)))
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}
