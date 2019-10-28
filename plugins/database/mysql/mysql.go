package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	stdmysql "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

const (
	defaultMysqlRevocationStmts = `
		REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'%'; 
		DROP USER '{{name}}'@'%'
	`

	defaultMySQLRotateCredentialsSQL = `
		ALTER USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
	`

	mySQLTypeName = "mysql"
)

var (
	MetadataLen       int = 10
	LegacyMetadataLen int = 4
	UsernameLen       int = 32
	LegacyUsernameLen int = 16
)

var _ dbplugin.Database = (*MySQL)(nil)

type MySQL struct {
	*connutil.SQLConnectionProducer
	credsutil.CredentialsProducer
}

// New implements builtinplugins.BuiltinFactory
func New(displayNameLen, roleNameLen, usernameLen int) func() (interface{}, error) {
	return func() (interface{}, error) {
		db := new(displayNameLen, roleNameLen, usernameLen)
		// Wrap the plugin with middleware to sanitize errors
		dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.SecretValues)

		return dbType, nil
	}
}

func new(displayNameLen, roleNameLen, usernameLen int) *MySQL {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = mySQLTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: displayNameLen,
		RoleNameLen:    roleNameLen,
		UsernameLen:    usernameLen,
		Separator:      "-",
	}

	return &MySQL{
		SQLConnectionProducer: connProducer,
		CredentialsProducer:   credsProducer,
	}
}

// Run instantiates a MySQL object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	return runCommon(false, apiTLSConfig)
}

// Run instantiates a MySQL object, and runs the RPC server for the plugin
func RunLegacy(apiTLSConfig *api.TLSConfig) error {
	return runCommon(true, apiTLSConfig)
}

func runCommon(legacy bool, apiTLSConfig *api.TLSConfig) error {
	var f func() (interface{}, error)
	if legacy {
		f = New(credsutil.NoneLength, LegacyMetadataLen, LegacyUsernameLen)
	} else {
		f = New(MetadataLen, MetadataLen, UsernameLen)
	}
	dbType, err := f()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

func (m *MySQL) Type() (string, error) {
	return mySQLTypeName, nil
}

func (m *MySQL) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

func (m *MySQL) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	statements = dbutil.StatementCompatibilityHelper(statements)

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

	queryMap := map[string]string{
		"name":       username,
		"password":   password,
		"expiration": expirationStr,
	}

	if err := m.executePreparedStatmentsWithMap(ctx, statements.Creation, queryMap); err != nil {
		return "", "", err
	}
	return username, password, nil
}

// NOOP
func (m *MySQL) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	return nil
}

func (m *MySQL) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the read lock
	m.Lock()
	defer m.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	// Get the connection
	db, err := m.getConnection(ctx)
	if err != nil {
		return err
	}

	revocationStmts := statements.Revocation
	// Use a default SQL statement for revocation if one cannot be fetched from the role
	if len(revocationStmts) == 0 {
		revocationStmts = []string{defaultMysqlRevocationStmts}
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range revocationStmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			// This is not a prepared statement because not all commands are supported
			// 1295: This command is not supported in the prepared statement protocol yet
			// Reference https://mariadb.com/kb/en/mariadb/prepare-statement/
			query = strings.Replace(query, "{{name}}", username, -1)
			_, err = tx.ExecContext(ctx, query)
			if err != nil {
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

func (m *MySQL) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	m.Lock()
	defer m.Unlock()

	if len(m.Username) == 0 || len(m.Password) == 0 {
		return nil, errors.New("username and password are required to rotate")
	}

	rotateStatements := statements
	if len(rotateStatements) == 0 {
		rotateStatements = []string{defaultMySQLRotateCredentialsSQL}
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

	for _, stmt := range rotateStatements {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			// This is not a prepared statement because not all commands are supported
			// 1295: This command is not supported in the prepared statement protocol yet
			// Reference https://mariadb.com/kb/en/mariadb/prepare-statement/
			query = strings.Replace(query, "{{username}}", m.Username, -1)
			query = strings.Replace(query, "{{password}}", password, -1)

			if _, err := tx.ExecContext(ctx, query); err != nil {
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

// SetCredentials uses provided information to set the password to a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for setting
// the password of static accounts, as well as rolling back passwords in the
// database in the event an updated database fails to save in Vault's storage.
func (m *MySQL) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	rotateStatements := statements.Rotation
	if len(rotateStatements) == 0 {
		rotateStatements = []string{defaultMySQLRotateCredentialsSQL}
	}

	username = staticUser.Username
	password = staticUser.Password
	if username == "" || password == "" {
		return "", "", errors.New("must provide both username and password")
	}

	queryMap := map[string]string{
		"name":     username,
		"password": password,
	}

	if err := m.executePreparedStatmentsWithMap(ctx, statements.Rotation, queryMap); err != nil {
		return "", "", err
	}
	return username, password, nil
}

// executePreparedStatmentsWithMap loops through the given templated SQL statements and
// applies the a map to them, interpolating values into the templates,returning
// tthe resulting username and password
func (m *MySQL) executePreparedStatmentsWithMap(ctx context.Context, statements []string, queryMap map[string]string) error {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	// Get the connection
	db, err := m.getConnection(ctx)
	if err != nil {
		return err
	}
	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Execute each query
	for _, stmt := range statements {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			query = dbutil.QueryHelper(query, queryMap)

			stmt, err := tx.PrepareContext(ctx, query)
			if err != nil {
				// If the error code we get back is Error 1295: This command is not
				// supported in the prepared statement protocol yet, we will execute
				// the statement without preparing it. This allows the caller to
				// manually prepare statements, as well as run other not yet
				// prepare supported commands. If there is no error when running we
				// will continue to the next statement.
				if e, ok := err.(*stdmysql.MySQLError); ok && e.Number == 1295 {
					_, err = tx.ExecContext(ctx, query)
					if err != nil {
						stmt.Close()
						return err
					}
					continue
				}

				return err
			}
			if _, err := stmt.ExecContext(ctx); err != nil {
				stmt.Close()
				return err
			}
			stmt.Close()
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
