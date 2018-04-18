package hana

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/SAP/go-hdb/driver"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/dbtxn"
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
	*connutil.SQLConnectionProducer
	credsutil.CredentialsProducer
}

var _ dbplugin.Database = &HANA{}

// New implements builtinplugins.BuiltinFactory
func New() (interface{}, error) {
	db := new()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.SecretValues)

	return dbType, nil
}

func new() *HANA {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = hanaTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 32,
		RoleNameLen:    20,
		UsernameLen:    128,
		Separator:      "_",
	}

	return &HANA{
		SQLConnectionProducer: connProducer,
		CredentialsProducer:   credsProducer,
	}
}

// Run instantiates a HANA object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(dbplugin.Database), apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (h *HANA) Type() (string, error) {
	return hanaTypeName, nil
}

func (h *HANA) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := h.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

// CreateUser generates the username/password on the underlying HANA secret backend
// as instructed by the CreationStatement provided.
func (h *HANA) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	h.Lock()
	defer h.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	// Get the connection
	db, err := h.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	if len(statements.Creation) == 0 {
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

// Renewing hana user just means altering user's valid until property
func (h *HANA) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	statements = dbutil.StatementCompatibilityHelper(statements)

	// Get connection
	db, err := h.getConnection(ctx)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
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
	stmt, err := tx.PrepareContext(ctx, "ALTER USER "+username+" VALID UNTIL "+"'"+expirationStr+"'")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Revoking hana user will deactivate user and try to perform a soft drop
func (h *HANA) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	statements = dbutil.StatementCompatibilityHelper(statements)

	// default revoke will be a soft drop on user
	if len(statements.Revocation) == 0 {
		return h.revokeUserDefault(ctx, username)
	}

	// Get connection
	db, err := h.getConnection(ctx)
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

	return tx.Commit()
}

func (h *HANA) revokeUserDefault(ctx context.Context, username string) error {
	// Get connection
	db, err := h.getConnection(ctx)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Disable server login for user
	disableStmt, err := tx.PrepareContext(ctx, fmt.Sprintf("ALTER USER %s DEACTIVATE USER NOW", username))
	if err != nil {
		return err
	}
	defer disableStmt.Close()
	if _, err := disableStmt.ExecContext(ctx); err != nil {
		return err
	}

	// Invalidates current sessions and performs soft drop (drop if no dependencies)
	// if hard drop is desired, custom revoke statements should be written for role
	dropStmt, err := tx.PrepareContext(ctx, fmt.Sprintf("DROP USER %s RESTRICT", username))
	if err != nil {
		return err
	}
	defer dropStmt.Close()
	if _, err := dropStmt.ExecContext(ctx); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// RotateRootCredentials is not currently supported on HANA
func (h *HANA) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	return nil, errors.New("root credentaion rotation is not currently implemented in this database secrets engine")
}
