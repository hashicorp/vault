package mysql

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const (
	defaultMysqlRevocationStmts = `
		REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'%'; 
		DROP USER '{{name}}'@'%'
	`
	mySQLTypeName = "mysql"
)

var (
	MetadataLen       int = 10
	LegacyMetadataLen int = 4
	UsernameLen       int = 32
	LegacyUsernameLen int = 16
)

type MySQL struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

// New implements builtinplugins.BuiltinFactory
func New(displayNameLen, roleNameLen, usernameLen int) func() (interface{}, error) {
	return func() (interface{}, error) {
		connProducer := &connutil.SQLConnectionProducer{}
		connProducer.Type = mySQLTypeName

		credsProducer := &credsutil.SQLCredentialsProducer{
			DisplayNameLen: displayNameLen,
			RoleNameLen:    roleNameLen,
			UsernameLen:    usernameLen,
			Separator:      "-",
		}

		dbType := &MySQL{
			ConnectionProducer:  connProducer,
			CredentialsProducer: credsProducer,
		}

		return dbType, nil
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

	plugins.Serve(dbType.(*MySQL), apiTLSConfig)

	return nil
}

func (m *MySQL) Type() (string, error) {
	return mySQLTypeName, nil
}

func (m *MySQL) getConnection() (*sql.DB, error) {
	db, err := m.Connection()
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

func (m *MySQL) CreateUser(statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
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

// NOOP
func (m *MySQL) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	return nil
}

func (m *MySQL) RevokeUser(statements dbplugin.Statements, username string) error {
	// Grab the read lock
	m.Lock()
	defer m.Unlock()

	// Get the connection
	db, err := m.getConnection()
	if err != nil {
		return err
	}

	revocationStmts := statements.RevocationStatements
	// Use a default SQL statement for revocation if one cannot be fetched from the role
	if revocationStmts == "" {
		revocationStmts = defaultMysqlRevocationStmts
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
