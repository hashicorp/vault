// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hana

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/SAP/go-hdb/driver"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/template"
)

const (
	hanaTypeName = "hdb"

	defaultUserNameTemplate = `{{ printf "v_%s_%s_%s_%s" (.DisplayName | truncate 32) (.RoleName | truncate 20) (random 20) (unix_time) | truncate 127 | replace "-" "_" | uppercase }}`
)

// HANA is an implementation of Database interface
type HANA struct {
	*connutil.SQLConnectionProducer

	usernameProducer template.StringTemplate
}

var _ dbplugin.Database = (*HANA)(nil)

// New implements builtinplugins.BuiltinFactory
func New() (interface{}, error) {
	db := new()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)

	return dbType, nil
}

func new() *HANA {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = hanaTypeName

	return &HANA{
		SQLConnectionProducer: connProducer,
	}
}

func (h *HANA) secretValues() map[string]string {
	return map[string]string{
		h.Password: "[password]",
	}
}

func (h *HANA) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	conf, err := h.Init(ctx, req.Config, req.VerifyConnection)
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
	h.usernameProducer = up

	_, err = h.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	return dbplugin.InitializeResponse{
		Config: conf,
	}, nil
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

// NewUser generates the username/password on the underlying HANA secret backend
// as instructed by the CreationStatement provided.
func (h *HANA) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (response dbplugin.NewUserResponse, err error) {
	// Grab the lock
	h.Lock()
	defer h.Unlock()

	// Get the connection
	db, err := h.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	if len(req.Statements.Commands) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	// Generate username
	username, err := h.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	// HANA does not allow hyphens in usernames, and highly prefers capital letters
	username = strings.ReplaceAll(username, "-", "_")
	username = strings.ToUpper(username)

	// If expiration is in the role SQL, HANA will deactivate the user when time is up,
	// regardless of whether vault is alive to revoke lease
	expirationStr := req.Expiration.UTC().Format("2006-01-02 15:04:05")

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	defer tx.Rollback()

	// Execute each query
	for _, stmt := range req.Statements.Commands {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       username,
				"password":   req.Password,
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

	resp := dbplugin.NewUserResponse{
		Username: username,
	}

	return resp, nil
}

// UpdateUser allows for updating the expiration or password of the user mentioned in
// the UpdateUserRequest
func (h *HANA) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	h.Lock()
	defer h.Unlock()

	// No change requested
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, nil
	}

	// Get connection
	db, err := h.getConnection(ctx)
	if err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}
	defer tx.Rollback()

	if req.Password != nil {
		err = h.updateUserPassword(ctx, tx, req.Username, req.Password)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, err
		}
	}

	if req.Expiration != nil {
		err = h.updateUserExpiration(ctx, tx, req.Username, req.Expiration)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}

	return dbplugin.UpdateUserResponse{}, nil
}

func (h *HANA) updateUserPassword(ctx context.Context, tx *sql.Tx, username string, req *dbplugin.ChangePassword) error {
	password := req.NewPassword

	if username == "" || password == "" {
		return fmt.Errorf("must provide both username and password")
	}

	stmts := req.Statements.Commands
	if len(stmts) == 0 {
		stmts = []string{"ALTER USER {{username}} PASSWORD \"{{password}}\""}
	}

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

			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
	}

	return nil
}

func (h *HANA) updateUserExpiration(ctx context.Context, tx *sql.Tx, username string, req *dbplugin.ChangeExpiration) error {
	// If expiration is in the role SQL, HANA will deactivate the user when time is up,
	// regardless of whether vault is alive to revoke lease
	expirationStr := req.NewExpiration.String()

	if username == "" || expirationStr == "" {
		return fmt.Errorf("must provide both username and expiration")
	}

	stmts := req.Statements.Commands
	if len(stmts) == 0 {
		stmts = []string{"ALTER USER {{username}} VALID UNTIL '{{expiration}}'"}
	}

	for _, stmt := range stmts {
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

			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
	}

	return nil
}

// Revoking hana user will deactivate user and try to perform a soft drop
func (h *HANA) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	h.Lock()
	defer h.Unlock()

	// default revoke will be a soft drop on user
	if len(req.Statements.Commands) == 0 {
		return h.revokeUserDefault(ctx, req)
	}

	// Get connection
	db, err := h.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer tx.Rollback()

	// Execute each query
	for _, stmt := range req.Statements.Commands {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name": req.Username,
			}
			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return dbplugin.DeleteUserResponse{}, err
			}
		}
	}

	return dbplugin.DeleteUserResponse{}, tx.Commit()
}

func (h *HANA) revokeUserDefault(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	// Get connection
	db, err := h.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer tx.Rollback()

	// Disable server login for user
	disableStmt, err := tx.PrepareContext(ctx, fmt.Sprintf("ALTER USER %s DEACTIVATE USER NOW", req.Username))
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer disableStmt.Close()
	if _, err := disableStmt.ExecContext(ctx); err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	// Invalidates current sessions and performs soft drop (drop if no dependencies)
	// if hard drop is desired, custom revoke statements should be written for role
	dropStmt, err := tx.PrepareContext(ctx, fmt.Sprintf("DROP USER %s RESTRICT", req.Username))
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer dropStmt.Close()
	if _, err := dropStmt.ExecContext(ctx); err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	return dbplugin.DeleteUserResponse{}, nil
}
