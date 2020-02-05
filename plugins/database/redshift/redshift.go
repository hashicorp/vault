package redshift

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/lib/pq"
)

const (
	sqlTypeName     = "postgres"
	defaultRenewSQL = `
ALTER USER "{{name}}" VALID UNTIL '{{expiration}}';
`
	defaultRotateRootCredentialsSQL = `
ALTER USER "{{username}}" WITH PASSWORD '{{password}}';
`

	defaultRotateCredentialsSQL = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
`
)

var _ dbplugin.Database = &RedShift{}

// New implements builtinplugins.BuiltinFactory
func New(lowercaseUsername bool) func() (interface{}, error) {
	return func() (interface{}, error) {
		db := new(lowercaseUsername)
		// Wrap the plugin with middleware to sanitize errors
		dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.SecretValues)
		return dbType, nil
	}
}

func new(lowercaseUsername bool) *RedShift {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = sqlTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen:    8,
		RoleNameLen:       8,
		UsernameLen:       63,
		Separator:         "-",
		LowercaseUsername: lowercaseUsername,
	}

	db := &RedShift{
		SQLConnectionProducer: connProducer,
		CredentialsProducer:   credsProducer,
	}

	return db
}

// Run instantiates a RedShift object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New(true)()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

type RedShift struct {
	*connutil.SQLConnectionProducer
	credsutil.CredentialsProducer
}

func (p *RedShift) Type() (string, error) {
	return sqlTypeName, nil
}

func (p *RedShift) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := p.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (p *RedShift) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	if len(statements.Rotation) == 0 {
		return "", "", errors.New("empty rotation statements")
	}

	username = staticUser.Username
	password = staticUser.Password
	if username == "" || password == "" {
		return "", "", errors.New("must provide both username and password")
	}

	// Grab the lock
	p.Lock()
	defer p.Unlock()

	// Get the connection
	db, err := p.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	// Check if the role exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT exists (SELECT usename FROM pg_user WHERE usename=$1);", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return "", "", err
	}

	// Vault requires the database user already exist, and that the credentials
	// used to execute the rotation statements has sufficient privileges.
	stmts := statements.Rotation

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Execute each query
	for _, stmt := range stmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":     staticUser.Username,
				"password": password,
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

func (p *RedShift) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
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

func (p *RedShift) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	p.Lock()
	defer p.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	renewStmts := statements.Renewal
	if len(renewStmts) == 0 {
		renewStmts = []string{defaultRenewSQL}
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

func (p *RedShift) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the lock
	p.Lock()
	defer p.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Revocation) == 0 {
		return p.defaultRevokeUser(ctx, username)
	}

	return p.customRevokeUser(ctx, username, statements.Revocation)
}

func (p *RedShift) customRevokeUser(ctx context.Context, username string, revocationStmts []string) error {
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

func (p *RedShift) defaultRevokeUser(ctx context.Context, username string) error {
	db, err := p.getConnection(ctx)
	if err != nil {
		return err
	}

	// Check if the role exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT exists (SELECT usename FROM pg_user WHERE usename=$1);", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if !exists {
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
		"REVOKE USAGE ON SCHEMA public FROM %s;",
		pq.QuoteIdentifier(username)))

	// get the current database name so we can issue a REVOKE CONNECT for
	// this username
	var dbname sql.NullString
	if err := db.QueryRowContext(ctx, "SELECT current_database();").Scan(&dbname); err != nil {
		return err
	}

	if dbname.Valid {
		/*
			- SP means we can't have parallel sessions creating and destroying the SP
			- SP means we may accidentially leave a SP on the database if there's a network partition
			- Since we've already revoked all the client privs, and those revocations are instant, we should
			rely on the client killing its own connection instead of us doing it for them
			- For large redshift clusters, upper bound on connections is 5000
			- For small redshift clusters, upper bound on connections is less than 1000 (need to double check)
			- Sessions are memory bound, so this may be an issue if clients fail to exit their session
			- Scenario: TTL of 1 minute with hundreds of users means we may have hundreds of zombie connections
			  Reality: TTLs are longer lived since apps can't handle reconnecting this often anyways, so it's probably ok?
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
		//revocationStmts = append(revocationStmts, `DROP PROCEDURE terminateloop(varchar);`)
	}

	fmt.Printf("%+v\n", revocationStmts)

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
		`DROP USER IF EXISTS %s;`, pq.QuoteIdentifier(username)))
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (p *RedShift) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	p.Lock()
	defer p.Unlock()

	if len(p.Username) == 0 || len(p.Password) == 0 {
		return nil, errors.New("username and password are required to rotate")
	}

	rotateStatents := statements
	if len(rotateStatents) == 0 {
		rotateStatents = []string{defaultRotateRootCredentialsSQL}
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
