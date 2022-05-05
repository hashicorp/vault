package ysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/lib/pq"
)

const (
	yugabyteDBType             = "yugabyte"
	defaultExpirationStatement = `
ALTER ROLE "{{name}}" VALID UNTIL '{{expiration}}';
`
	defaultChangePasswordStatement = `
ALTER ROLE "{{username}}" WITH PASSWORD '{{password}}';
`
	expirationFormat = "2006-01-02T15:04:05Z07:00" // "2006-01-02 15:04:05-0700"

	defaultUserNameTemplate = `{{ printf "v-%s-%s-%s-%s" (.DisplayName | truncate 8) (.RoleName | truncate 8) (random 20) (unix_time) | truncate 63 }}`
)

var (
	_ dbplugin.Database = &ysql{}

	// postgresEndStatement is basically the word "END" but
	// surrounded by a word boundary to differentiate it from
	// other words like "APPEND".
	postgresEndStatement = regexp.MustCompile(`\bEND\b`)

	// doubleQuotedPhrases finds substrings like "hello"
	// and pulls them out with the quotes included.
	doubleQuotedPhrases = regexp.MustCompile(`(".*?")`)

	// singleQuotedPhrases finds substrings like 'hello'
	// and pulls them out with the quotes included.
	singleQuotedPhrases = regexp.MustCompile(`('.*?')`)
)

type ysql struct {
	YugabyteConnectionProducer
	usernameProducer template.StringTemplate
}

func New() (interface{}, error) {
	db := new()

	// This middleware isn't strictly required, but highly recommended to prevent accidentally exposing
	// values such as passwords in error messages. An example of this is included below
	// DatabaseErrorSanitizerMiddleware wraps an implementation of Databases and
	// sanitizes returned error messages
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

var _ dbplugin.Database = (*ysql)(nil)

func new() *ysql {
	connProducer := YugabyteConnectionProducer{}
	connProducer.Type = yugabyteDBType

	yugabyte := &ysql{
		YugabyteConnectionProducer: connProducer,
		usernameProducer:           template.StringTemplate{},
	}
	return yugabyte
}

func (db *ysql) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	usernameTemplate, err := strutil.GetString(req.Config, "username_template")
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("failed to retrieve username_template: %w", err)
	}

	log.Println("initializing --> ", usernameTemplate)

	if usernameTemplate == "" {
		usernameTemplate = defaultUserNameTemplate
	}

	up, err := template.NewTemplate(template.Template(usernameTemplate))
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("unable to initialize username template: %w", err)
	}
	db.usernameProducer = up

	_, err = db.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	err = db.YugabyteConnectionProducer.Initialize(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return dbplugin.InitializeResponse{}, err
	}
	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

func (ydb *ysql) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	if len(req.Statements.Commands) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	ydb.Lock()
	defer ydb.Unlock()

	username, err := ydb.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	expirationStr := req.Expiration.Format(expirationFormat)

	db, err := ydb.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, fmt.Errorf("unable to get connection: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.NewUserResponse{}, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, stmt := range req.Statements.Commands {
		if containsMultilineStatement(stmt) {
			// Execute it as-is.
			m := map[string]string{
				"name":       username,
				"username":   username,
				"password":   req.Password,
				"expiration": expirationStr,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, stmt); err != nil {
				return dbplugin.NewUserResponse{}, fmt.Errorf("failed to execute query: %w", err)
			}
			continue
		}
		// Otherwise, it's fine to split the statements on the semicolon.
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       username,
				"username":   username,
				"password":   req.Password,
				"expiration": expirationStr,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return dbplugin.NewUserResponse{}, fmt.Errorf("failed to execute query: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	resp := dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

func (ydb *ysql) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Username == "" {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("missing username")
	}
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("no changes requested")
	}

	merr := &multierror.Error{}
	if req.Password != nil {
		err := ydb.changeUserPassword(ctx, req.Username, req.Password)
		merr = multierror.Append(merr, err)
	}
	if req.Expiration != nil {
		err := ydb.changeUserExpiration(ctx, req.Username, req.Expiration)
		merr = multierror.Append(merr, err)
	}
	return dbplugin.UpdateUserResponse{}, merr.ErrorOrNil()
}

func (ydb *ysql) changeUserPassword(ctx context.Context, username string, changePass *dbplugin.ChangePassword) error {
	stmts := changePass.Statements.Commands
	if len(stmts) == 0 {
		stmts = []string{defaultChangePasswordStatement}
	}

	password := changePass.NewPassword
	if password == "" {
		return fmt.Errorf("missing password")
	}

	ydb.Lock()
	defer ydb.Unlock()

	db, err := ydb.getConnection(ctx)
	if err != nil {
		return fmt.Errorf("unable to get connection: %w", err)
	}

	// Check if the role exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT exists (SELECT rolname FROM pg_roles WHERE rolname=$1);", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("user does not appear to exist: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
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
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ydb *ysql) changeUserExpiration(ctx context.Context, username string, changeExp *dbplugin.ChangeExpiration) error {
	ydb.Lock()
	defer ydb.Unlock()

	renewStmts := changeExp.Statements.Commands
	if len(renewStmts) == 0 {
		renewStmts = []string{defaultExpirationStatement}
	}

	db, err := ydb.getConnection(ctx)
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

	expirationStr := changeExp.NewExpiration.Format(expirationFormat)

	for _, stmt := range renewStmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       username,
				"username":   username,
				"expiration": expirationStr,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (ydb *ysql) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	ydb.Lock()
	defer ydb.Unlock()
	if len(req.Statements.Commands) == 0 {
		return dbplugin.DeleteUserResponse{}, ydb.defaultDeleteUser(ctx, req.Username)
	}

	return dbplugin.DeleteUserResponse{}, ydb.customDeleteUser(ctx, req.Username, req.Statements.Commands)
}

func (ydb *ysql) customDeleteUser(ctx context.Context, username string, revocationStmts []string) error {
	db, err := ydb.getConnection(ctx)
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
				"name":     username,
				"username": username,
			}
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (ydb *ysql) defaultDeleteUser(ctx context.Context, username string) error {
	db, err := ydb.getConnection(ctx)
	if err != nil {
		return err
	}

	// Check if the role exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT exists (SELECT rolname FROM pg_roles WHERE rolname=$1);", username).Scan(&exists)
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
		return fmt.Errorf("unable to prepare context : %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return fmt.Errorf("unable to execute query: %w ", err)
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
			`REVOKE ALL PRIVILEGES  ON SCHEMA %s FROM %s;`,
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
		"REVOKE ALL PRIVILEGES  ON SCHEMA public FROM %s;",
		pq.QuoteIdentifier(username)))

	// get the current database name so we can issue a REVOKE CONNECT for
	// this username
	var dbname sql.NullString
	if err := db.QueryRowContext(ctx, "SELECT current_database();").Scan(&dbname); err != nil {
		return err
	}

	if dbname.Valid {
		revocationStmts = append(revocationStmts, fmt.Sprintf(
			`REVOKE ALL PRIVILEGES ON DATABASE %s FROM %s;`,
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
		return fmt.Errorf("could not generate revocation statements for all rows: %w", rows.Err())
	}
	if lastStmtError != nil {
		return fmt.Errorf("could not perform all revocation statements: %w", lastStmtError)
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

func (ydb *ysql) secretValues() map[string]string {
	return map[string]string{
		ydb.Password: "[password]",
	}
}

// containsMultilineStatement is a best effort to determine whether
// a particular statement is multiline, and therefore should not be
// split upon semicolons. If it's unsure, it defaults to false.
func containsMultilineStatement(stmt string) bool {
	// We're going to look for the word "END", but first let's ignore
	// anything the user provided within single or double quotes since
	// we're looking for an "END" within the Postgres syntax.
	literals, err := extractQuotedStrings(stmt)
	if err != nil {
		return false
	}
	stmtWithoutLiterals := stmt
	for _, literal := range literals {
		stmtWithoutLiterals = strings.Replace(stmt, literal, "", -1)
	}
	// Now look for the word "END" specifically. This will miss any
	// representations of END that aren't surrounded by spaces, but
	// it should be easy to change on the user's side.
	return postgresEndStatement.MatchString(stmtWithoutLiterals)
}

// extractQuotedStrings extracts 0 or many substrings
// that have been single- or double-quoted. Ex:
// `"Hello", silly 'elephant' from the "zoo".`
// returns [ `Hello`, `'elephant'`, `"zoo"` ]
func extractQuotedStrings(s string) ([]string, error) {
	var found []string
	toFind := []*regexp.Regexp{
		doubleQuotedPhrases,
		singleQuotedPhrases,
	}
	for _, typeOfPhrase := range toFind {
		found = append(found, typeOfPhrase.FindAllString(s, -1)...)
	}
	return found, nil
}

func (db *ysql) Type() (string, error) {
	return yugabyteDBType, nil
}

func (db *ysql) getConnection(ctx context.Context) (*sql.DB, error) {
	conn, err := db.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return conn.(*sql.DB), nil
}
