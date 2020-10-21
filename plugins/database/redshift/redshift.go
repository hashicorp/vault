package redshift

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/lib/pq"
)

const (
	// This is how this plugin will be reflected in middleware
	// such as metrics.
	middlewareTypeName = "redshift"

	// This allows us to use the postgres database driver.
	sqlTypeName = "postgres"

	defaultRenewSQL = `
ALTER USER "{{name}}" VALID UNTIL '{{expiration}}';
`
	defaultRotateRootCredentialsSQL = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
`
)

var _ dbplugin.Database = (*RedShift)(nil)

// lowercaseUsername is the reason we wrote this plugin. Redshift implements (mostly)
// a postgres 8 interface, and part of that is under the hood, it's lowercasing the
// usernames.
func New(lowercaseUsername bool) func() (interface{}, error) {
	return func() (interface{}, error) {
		db := newRedshift(lowercaseUsername)
		// Wrap the plugin with middleware to sanitize errors
		dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
		return dbType, nil
	}
}

func newRedshift(lowercaseUsername bool) *RedShift {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = sqlTypeName

	db := &RedShift{
		SQLConnectionProducer: connProducer,
		lowerCaseUsername:     lowercaseUsername,
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
	lowerCaseUsername bool
}

func (r *RedShift) secretValues() map[string]string {
	return map[string]string{
		r.Password: "[password]",
	}
}

func (r *RedShift) Type() (string, error) {
	return middlewareTypeName, nil
}

func (r *RedShift) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	conf, err := r.Init(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("error initializing db: %w", err)
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

func (r *RedShift) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	if len(req.Statements.Commands) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	// Grab the lock
	r.Lock()
	defer r.Unlock()

	usernameOpts := []credsutil.UsernameOpt{
		credsutil.DisplayName(req.UsernameConfig.DisplayName, 8),
		credsutil.RoleName(req.UsernameConfig.RoleName, 8),
		credsutil.MaxLength(63),
		credsutil.Separator("-"),
	}
	if r.lowerCaseUsername {
		usernameOpts = append(usernameOpts, credsutil.ToLower())
	}

	username, err := credsutil.GenerateUsername(usernameOpts...)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	password := req.Password
	expirationStr := req.Expiration.UTC().Format("2006-01-02 15:04:05")

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
				"password":   password,
				"expiration": expirationStr,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
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

func (r *RedShift) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, nil
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
		renewStmts := req.Expiration.Statements
		if len(renewStmts.Commands) == 0 {
			renewStmts.Commands = []string{defaultRenewSQL}
		}

		expirationStr := req.Expiration.NewExpiration.UTC().Format("2006-01-02 15:04:05")

		for _, stmt := range renewStmts.Commands {
			for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
				query = strings.TrimSpace(query)
				if len(query) == 0 {
					continue
				}

				m := map[string]string{
					"name":       req.Username,
					"expiration": expirationStr,
				}
				if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
					return dbplugin.UpdateUserResponse{}, err
				}
			}
		}
	}

	if req.Password != nil {
		username := req.Username
		password := req.Password.NewPassword
		if username == "" || password == "" {
			return dbplugin.UpdateUserResponse{}, errors.New("must provide both username and password")
		}

		// Check if the role exists
		var exists bool
		err = db.QueryRowContext(ctx, "SELECT exists (SELECT usename FROM pg_user WHERE usename=$1);", username).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			return dbplugin.UpdateUserResponse{}, err
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
				if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
					return dbplugin.UpdateUserResponse{}, err
				}
			}
		}
	}

	return dbplugin.UpdateUserResponse{}, tx.Commit()
}

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
				"name": req.Username,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
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
	var lastStmtError *multierror.Error //error
	for _, query := range revocationStmts {
		if err := dbtxn.ExecuteDBQuery(ctx, db, nil, query); err != nil {
			lastStmtError = multierror.Append(lastStmtError, err)
		}
	}

	// can't drop if not all privileges are revoked
	if rows.Err() != nil {
		return dbplugin.DeleteUserResponse{}, errwrap.Wrapf("could not generate revocation statements for all rows: {{err}}", rows.Err())
	}
	if lastStmtError != nil {
		return dbplugin.DeleteUserResponse{}, errwrap.Wrapf("could not perform all revocation statements: {{err}}", lastStmtError)
	}

	// Drop this user
	stmt, err = db.PrepareContext(ctx, fmt.Sprintf(
		`DROP USER IF EXISTS %s;`, pq.QuoteIdentifier(username)))
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx); err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	return dbplugin.DeleteUserResponse{}, nil
}
