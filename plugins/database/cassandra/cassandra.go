// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/template"
)

const (
	defaultUserCreationCQL   = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;`
	defaultUserDeletionCQL   = `DROP USER '{{username}}';`
	defaultChangePasswordCQL = `ALTER USER '{{username}}' WITH PASSWORD '{{password}}';`
	cassandraTypeName        = "cassandra"

	defaultUserNameTemplate = `{{ printf "v_%s_%s_%s_%s" (.DisplayName | truncate 15) (.RoleName | truncate 15) (random 20) (unix_time) | truncate 100 | replace "-" "_" | lowercase }}`
)

var _ dbplugin.Database = &Cassandra{}

// Cassandra is an implementation of Database interface
type Cassandra struct {
	*cassandraConnectionProducer

	usernameProducer template.StringTemplate
}

// New returns a new Cassandra instance
func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)

	return dbType, nil
}

func new() *Cassandra {
	connProducer := &cassandraConnectionProducer{}
	connProducer.Type = cassandraTypeName

	return &Cassandra{
		cassandraConnectionProducer: connProducer,
	}
}

// Type returns the TypeName for this backend
func (c *Cassandra) Type() (string, error) {
	return cassandraTypeName, nil
}

func (c *Cassandra) getConnection(ctx context.Context) (*gocql.Session, error) {
	session, err := c.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return session.(*gocql.Session), nil
}

func (c *Cassandra) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
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
	c.usernameProducer = up

	_, err = c.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	err = c.cassandraConnectionProducer.Initialize(ctx, req)
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("failed to initialize: %w", err)
	}

	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

// NewUser generates the username/password on the underlying Cassandra secret backend as instructed by
// the statements provided.
func (c *Cassandra) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	c.Lock()
	defer c.Unlock()

	session, err := c.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	creationCQL := req.Statements.Commands
	if len(creationCQL) == 0 {
		creationCQL = []string{defaultUserCreationCQL}
	}

	rollbackCQL := req.RollbackStatements.Commands
	if len(rollbackCQL) == 0 {
		rollbackCQL = []string{defaultUserDeletionCQL}
	}

	username, err := c.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	for _, stmt := range creationCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": username,
				"password": req.Password,
			}
			err = session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()
			if err != nil {
				rollbackErr := rollbackUser(ctx, session, username, rollbackCQL)
				if rollbackErr != nil {
					err = multierror.Append(err, rollbackErr)
				}
				return dbplugin.NewUserResponse{}, err
			}
		}
	}

	resp := dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

func rollbackUser(ctx context.Context, session *gocql.Session, username string, rollbackCQL []string) error {
	for _, stmt := range rollbackCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": username,
			}
			err := session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()
			if err != nil {
				return fmt.Errorf("failed to roll back user %s: %w", username, err)
			}
		}
	}
	return nil
}

func (c *Cassandra) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("no changes requested")
	}

	if req.Password != nil {
		err := c.changeUserPassword(ctx, req.Username, req.Password)
		return dbplugin.UpdateUserResponse{}, err
	}
	// Expiration is no-op
	return dbplugin.UpdateUserResponse{}, nil
}

func (c *Cassandra) changeUserPassword(ctx context.Context, username string, changePass *dbplugin.ChangePassword) error {
	session, err := c.getConnection(ctx)
	if err != nil {
		return err
	}

	rotateCQL := changePass.Statements.Commands
	if len(rotateCQL) == 0 {
		rotateCQL = []string{defaultChangePasswordCQL}
	}

	var result *multierror.Error
	for _, stmt := range rotateCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": username,
				"password": changePass.NewPassword,
			}
			err := session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// DeleteUser attempts to drop the specified user.
func (c *Cassandra) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	c.Lock()
	defer c.Unlock()

	session, err := c.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	revocationCQL := req.Statements.Commands
	if len(revocationCQL) == 0 {
		revocationCQL = []string{defaultUserDeletionCQL}
	}

	var result *multierror.Error
	for _, stmt := range revocationCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": req.Username,
			}
			err := session.
				Query(dbutil.QueryHelper(query, m)).
				WithContext(ctx).
				Exec()

			result = multierror.Append(result, err)
		}
	}

	return dbplugin.DeleteUserResponse{}, result.ErrorOrNil()
}
