package hana

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/SAP/go-hdb/driver"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const (
	hanaTypeName = "hdb"
)

// HANA is an implementation of Database interface
type HANA struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

// New implements builtinplugins.BuiltinFactory
func New() (interface{}, error) {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = hanaTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 32,
		RoleNameLen:    20,
		UsernameLen:    128,
		Separator:      "_",
	}

	dbType := &HANA{
		ConnectionProducer:  connProducer,
		CredentialsProducer: credsProducer,
	}

	return dbType, nil
}

// Run instantiates a HANA object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(*HANA), apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (h *HANA) Type() (string, error) {
	return hanaTypeName, nil
}

func (h *HANA) getConnection() (*sql.DB, error) {
	db, err := h.Connection()
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

// CreateUser generates the username/password on the underlying HANA secret backend
// as instructed by the CreationStatement provided.
func (h *HANA) CreateUser(statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	h.Lock()
	defer h.Unlock()

	// Get the connection
	db, err := h.getConnection()
	if err != nil {
		return "", "", err
	}

	if statements.CreationStatements == "" {
		return "", "", dbutil.ErrEmptyCreationStatement
	}

	// Generate username
	username, err = h.GenerateUsername(usernameConfig)
	if err != nil {
		return "", "", err
	}

	// HANA does not allow hyphens in usernames, and highly prefers capital letters
	username = strings.Replace(username, "-", "_", -1)
	username = strings.ToUpper(username)

	// Generate password
	password, err = h.GeneratePassword()
	if err != nil {
		return "", "", err
	}
	// Most HANA configurations have password constraints
	// Prefix with A1a to satisfy these constraints. User will be forced to change upon login
	password = strings.Replace(password, "-", "_", -1)
	password = "A1a" + password

	// If expiration is in the role SQL, HANA will deactivate the user when time is up,
	// regardless of whether vault is alive to revoke lease
	expirationStr, err := h.GenerateExpiration(expiration)
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

// Renewing hana user just means altering user's valid until property
func (h *HANA) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	// Get connection
	db, err := h.getConnection()
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// If expiration is in the role SQL, HANA will deactivate the user when time is up,
	// regardless of whether vault is alive to revoke lease
	expirationStr, err := h.GenerateExpiration(expiration)
	if err != nil {
		return err
	}

	// Renew user's valid until property field
	stmt, err := tx.Prepare("ALTER USER " + username + " VALID UNTIL " + "'" + expirationStr + "'")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Revoking hana user will deactivate user and try to perform a soft drop
func (h *HANA) RevokeUser(statements dbplugin.Statements, username string) error {
	// default revoke will be a soft drop on user
	if statements.RevocationStatements == "" {
		return h.revokeUserDefault(username)
	}

	// Get connection
	db, err := h.getConnection()
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

func (h *HANA) revokeUserDefault(username string) error {
	// Get connection
	db, err := h.getConnection()
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Disable server login for user
	disableStmt, err := tx.Prepare(fmt.Sprintf("ALTER USER %s DEACTIVATE USER NOW", username))
	if err != nil {
		return err
	}
	defer disableStmt.Close()
	if _, err := disableStmt.Exec(); err != nil {
		return err
	}

	// Invalidates current sessions and performs soft drop (drop if no dependencies)
	// if hard drop is desired, custom revoke statements should be written for role
	dropStmt, err := tx.Prepare(fmt.Sprintf("DROP USER %s RESTRICT", username))
	if err != nil {
		return err
	}
	defer dropStmt.Close()
	if _, err := dropStmt.Exec(); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
