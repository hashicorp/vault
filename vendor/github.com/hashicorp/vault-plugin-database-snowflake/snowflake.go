// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package snowflake

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/template"
	_ "github.com/snowflakedb/gosnowflake"
)

const (
	snowflakeSQLTypeName     = "snowflake"
	defaultSnowflakeRenewSQL = `
alter user {{name}} set DAYS_TO_EXPIRY = {{expiration}};
`
	defaultSnowflakeRotatePasswordSQL = `
alter user {{name}} set PASSWORD = '{{password}}';
`
	defaultSnowflakeRotateRSAPublicKeySQL = `
alter user {{name}} set RSA_PUBLIC_KEY = '{{public_key}}';
`
	defaultSnowflakeDeleteSQL = `
drop user if exists {{name}};
`
	defaultUserNameTemplate = `{{ printf "v_%s_%s_%s_%s" (.DisplayName | truncate 32) (.RoleName | truncate 32) (random 20) (unix_time) | truncate 255 | replace "-" "_" }}`
)

var _ dbplugin.Database = (*SnowflakeSQL)(nil)

func New() (interface{}, error) {
	db := new()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func (s *SnowflakeSQL) secretValues() map[string]string {
	return map[string]string{
		s.Password: "[password]",
	}
}

func new() *SnowflakeSQL {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = snowflakeSQLTypeName

	db := &SnowflakeSQL{
		SQLConnectionProducer: connProducer,
	}

	return db
}

type SnowflakeSQL struct {
	*connutil.SQLConnectionProducer
	sync.RWMutex

	usernameProducer template.StringTemplate
}

func (s *SnowflakeSQL) Type() (string, error) {
	return snowflakeSQLTypeName, nil
}

func (s *SnowflakeSQL) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := s.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

func (s *SnowflakeSQL) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	err := s.SQLConnectionProducer.Initialize(ctx, req.Config, req.VerifyConnection)
	if err != nil {
		return dbplugin.InitializeResponse{}, err
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
	s.usernameProducer = up

	_, err = s.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}
	resp.SetSupportedCredentialTypes([]dbplugin.CredentialType{
		dbplugin.CredentialTypePassword,
		dbplugin.CredentialTypeRSAPrivateKey,
	})

	return resp, nil
}

func (s *SnowflakeSQL) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	s.RLock()
	defer s.RUnlock()

	statements := req.Statements.Commands
	if len(statements) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	username, err := s.generateUsername(req)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	expirationStr, err := calculateExpirationString(req.Expiration)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	// Get the connection
	db, err := s.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	defer tx.Rollback()

	m := map[string]string{
		"name":       username,
		"username":   username,
		"expiration": expirationStr,
	}

	switch req.CredentialType {
	case dbplugin.CredentialTypePassword:
		m["password"] = req.Password
	case dbplugin.CredentialTypeRSAPrivateKey:
		m["public_key"] = preparePublicKey(string(req.PublicKey))
	default:
		return dbplugin.NewUserResponse{}, fmt.Errorf("unsupported credential type %q",
			req.CredentialType.String())
	}

	// Execute each query
	for _, stmt := range statements {
		// it's fine to split the statements on the semicolon.
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return dbplugin.NewUserResponse{}, err
			}
		}
	}

	err = tx.Commit()
	resp := dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, err
}

func (s *SnowflakeSQL) generateUsername(req dbplugin.NewUserRequest) (string, error) {
	username, err := s.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return "", errwrap.Wrapf("error generating username: {{err}}", err)
	}
	return username, nil
}

func (s *SnowflakeSQL) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	s.RLock()
	defer s.RUnlock()

	if req.Username == "" {
		err := fmt.Errorf("a username must be provided to update a user")
		return dbplugin.UpdateUserResponse{}, err
	}

	db, err := s.getConnection(ctx)
	if err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}
	defer tx.Rollback()

	if req.Password != nil || req.PublicKey != nil {
		err = s.updateUserCredential(ctx, tx, req)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, err
		}
	}

	if req.Expiration != nil {
		err = s.updateUserExpiration(ctx, tx, req.Username, req.Expiration)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}

	return dbplugin.UpdateUserResponse{}, nil
}

func (s *SnowflakeSQL) updateUserCredential(ctx context.Context, tx *sql.Tx, req dbplugin.UpdateUserRequest) error {
	m := map[string]string{
		"name":     req.Username,
		"username": req.Username,
	}

	var stmts []string
	switch req.CredentialType {
	case dbplugin.CredentialTypePassword:
		if req.Password == nil || req.Password.NewPassword == "" {
			return fmt.Errorf("new password credential must not be empty")
		}

		stmts = req.Password.Statements.Commands
		if len(stmts) == 0 {
			stmts = []string{defaultSnowflakeRotatePasswordSQL}
		}

		m["password"] = req.Password.NewPassword

	case dbplugin.CredentialTypeRSAPrivateKey:
		if req.PublicKey == nil || len(req.PublicKey.NewPublicKey) == 0 {
			return fmt.Errorf("new public key credential must not be empty")
		}

		stmts = req.PublicKey.Statements.Commands
		if len(stmts) == 0 {
			stmts = []string{defaultSnowflakeRotateRSAPublicKeySQL}
		}

		m["public_key"] = preparePublicKey(string(req.PublicKey.NewPublicKey))

	default:
		return fmt.Errorf("unsupported credential type %q", req.CredentialType.String())
	}

	for _, stmt := range stmts {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
	}

	return nil
}

func (s *SnowflakeSQL) updateUserExpiration(ctx context.Context, tx *sql.Tx, username string, req *dbplugin.ChangeExpiration) error {
	expiration := req.NewExpiration

	if username == "" || expiration.IsZero() {
		return fmt.Errorf("must provide both username and valid expiration to modify expiration")
	}

	expirationStr, err := calculateExpirationString(expiration)
	if err != nil {
		return err
	}

	stmts := req.Statements.Commands
	if len(stmts) == 0 {
		stmts = []string{defaultSnowflakeRenewSQL}
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

func (s *SnowflakeSQL) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	s.RLock()
	defer s.RUnlock()

	username := req.Username
	statements := req.Statements.Commands
	if len(statements) == 0 {
		statements = []string{defaultSnowflakeDeleteSQL}
	}

	db, err := s.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}
	defer tx.Rollback()

	for _, stmt := range statements {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":     username,
				"username": username,
			}
			if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
				return dbplugin.DeleteUserResponse{}, err
			}
		}
	}

	err = tx.Commit()
	return dbplugin.DeleteUserResponse{}, err
}

// calculateExpirationString has a minimum expiration of 1 Day. This
// limitation is due to Snowflake requiring any expiration to be in
// terms of days, with 1 being the minimum.
func calculateExpirationString(expiration time.Time) (string, error) {
	currentTime := time.Now()

	if currentTime.Before(expiration) {
		timeDiff := expiration.Sub(currentTime)
		inSeconds := timeDiff.Seconds()
		inDays := math.Max(math.Floor(inSeconds/float64(60*60*24)), 1)

		expirationStr := fmt.Sprintf("%d", int(inDays))
		return expirationStr, nil
	} else {
		err := fmt.Errorf("expiration time earlier than current time")
		return "", err
	}
}

// preparePublicKey strips the BEGIN and END lines from the given PEM string.
// This is required by Snowflake when setting the RSA_PUBLIC_KEY credential per
// the statement "Exclude the public key delimiters in the SQL statement" in
// https://docs.snowflake.com/en/user-guide/key-pair-auth.html#step-4-assign-the-public-key-to-a-snowflake-user
func preparePublicKey(pub string) string {
	pub = strings.Replace(pub, "-----BEGIN PUBLIC KEY-----\n", "", 1)
	return strings.Replace(pub, "-----END PUBLIC KEY-----\n", "", 1)
}
