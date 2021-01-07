package snowflake

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"

	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	_ "github.com/snowflakedb/gosnowflake"
)

const (
	snowflakeSQLTypeName     = "snowflake"
	maxIdentifierLength      = 255
	maxUsernameChunkLen      = 32 // arbitrarily chosen
	maxRolenameChunkLen      = 32 // arbitrarily chosen
	defaultSnowflakeRenewSQL = `
alter user {{name}} set DAYS_TO_EXPIRY = {{expiration}};
`
	defaultSnowflakeRotateCredsSQL = `
alter user {{name}} set PASSWORD = '{{password}}';
`
	defaultSnowflakeDeleteSQL = `
drop user {{name}};
`
)

var (
	_ dbplugin.Database = (*SnowflakeSQL)(nil)
)

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

	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}

	return resp, nil
}

func (s *SnowflakeSQL) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	s.Lock()
	defer s.Unlock()

	statements := req.Statements.Commands
	if len(statements) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	username, err := s.generateUsername(req)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	password := req.Password
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

	// Execute each query
	for _, stmt := range statements {
		// it's fine to split the statements on the semicolon.
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"name":       username,
				"username":   username,
				"password":   password,
				"expiration": expirationStr,
			}

			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
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
	username, err := credsutil.GenerateUsername(
		credsutil.DisplayName(req.UsernameConfig.DisplayName, maxUsernameChunkLen),
		credsutil.RoleName(req.UsernameConfig.RoleName, maxRolenameChunkLen),
		credsutil.MaxLength(maxIdentifierLength),
	)
	if err != nil {
		return "", errwrap.Wrapf("error generating username: {{err}}", err)
	}
	return username, nil
}

func (s *SnowflakeSQL) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	s.Lock()
	defer s.Unlock()

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

	if req.Password != nil {
		err = s.updateUserPassword(ctx, tx, req.Username, req.Password)
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

func (s *SnowflakeSQL) updateUserPassword(ctx context.Context, tx *sql.Tx, username string, req *dbplugin.ChangePassword) error {
	password := req.NewPassword

	if username == "" || password == "" {
		return fmt.Errorf("must provide both username and password to modify password")
	}

	stmts := req.Statements.Commands
	if len(stmts) == 0 {
		stmts = []string{defaultSnowflakeRotateCredsSQL}
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

			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
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

			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
	}

	return nil
}

func (s *SnowflakeSQL) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	s.Lock()
	defer s.Unlock()

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
			if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
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
