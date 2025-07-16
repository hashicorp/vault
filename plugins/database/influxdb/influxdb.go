// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package influxdb

import (
	"context"
	"fmt"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	influx "github.com/influxdata/influxdb1-client/v2"
)

const (
	defaultUserCreationIFQL           = `CREATE USER "{{username}}" WITH PASSWORD '{{password}}';`
	defaultUserDeletionIFQL           = `DROP USER "{{username}}";`
	defaultRootCredentialRotationIFQL = `SET PASSWORD FOR "{{username}}" = '{{password}}';`
	influxdbTypeName                  = "influxdb"

	defaultUserNameTemplate = `{{ printf "v_%s_%s_%s_%s" (.DisplayName | truncate 15) (.RoleName | truncate 15) (random 20) (unix_time) | truncate 100 | replace "-" "_" | lowercase }}`
)

var _ dbplugin.Database = &Influxdb{}

// Influxdb is an implementation of Database interface
type Influxdb struct {
	*influxdbConnectionProducer

	usernameProducer template.StringTemplate
}

// New returns a new Cassandra instance
func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)

	return dbType, nil
}

func new() *Influxdb {
	connProducer := &influxdbConnectionProducer{}
	connProducer.Type = influxdbTypeName

	return &Influxdb{
		influxdbConnectionProducer: connProducer,
	}
}

// Type returns the TypeName for this backend
func (i *Influxdb) Type() (string, error) {
	return influxdbTypeName, nil
}

func (i *Influxdb) getConnection(ctx context.Context) (influx.Client, error) {
	cli, err := i.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return cli.(influx.Client), nil
}

func (i *Influxdb) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (resp dbplugin.InitializeResponse, err error) {
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
	i.usernameProducer = up

	_, err = i.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	return i.influxdbConnectionProducer.Initialize(ctx, req)
}

// NewUser generates the username/password on the underlying Influxdb secret backend as instructed by
// the statements provided.
func (i *Influxdb) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (resp dbplugin.NewUserResponse, err error) {
	i.Lock()
	defer i.Unlock()

	cli, err := i.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, fmt.Errorf("unable to get connection: %w", err)
	}

	creationIFQL := req.Statements.Commands
	if len(creationIFQL) == 0 {
		creationIFQL = []string{defaultUserCreationIFQL}
	}

	rollbackIFQL := req.RollbackStatements.Commands
	if len(rollbackIFQL) == 0 {
		rollbackIFQL = []string{defaultUserDeletionIFQL}
	}

	username, err := i.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	for _, stmt := range creationIFQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			m := map[string]string{
				"username": username,
				"password": req.Password,
			}
			qry := influx.NewQuery(dbutil.QueryHelper(query, m), "", "")
			response, err := cli.Query(qry)
			// err can be nil with response.Error() being not nil, so both need to be handled
			merr := multierror.Append(err, response.Error())
			if merr.ErrorOrNil() != nil {
				// Attempt rollback only when the response has an error
				if response != nil && response.Error() != nil {
					attemptRollback(cli, username, rollbackIFQL)
				}

				return dbplugin.NewUserResponse{}, fmt.Errorf("failed to run query in InfluxDB: %w", merr)
			}
		}
	}
	resp = dbplugin.NewUserResponse{
		Username: username,
	}
	return resp, nil
}

// attemptRollback will attempt to roll back user creation if an error occurs in
// CreateUser
func attemptRollback(cli influx.Client, username string, rollbackStatements []string) error {
	for _, stmt := range rollbackStatements {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)

			if len(query) == 0 {
				continue
			}
			q := influx.NewQuery(dbutil.QueryHelper(query, map[string]string{
				"username": username,
			}), "", "")

			response, err := cli.Query(q)
			// err can be nil with response.Error() being not nil, so both need to be handled
			merr := multierror.Append(err, response.Error())
			if merr.ErrorOrNil() != nil {
				return merr
			}
		}
	}
	return nil
}

func (i *Influxdb) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	i.Lock()
	defer i.Unlock()

	cli, err := i.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("unable to get connection: %w", err)
	}

	revocationIFQL := req.Statements.Commands
	if len(revocationIFQL) == 0 {
		revocationIFQL = []string{defaultUserDeletionIFQL}
	}

	var result *multierror.Error
	for _, stmt := range revocationIFQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}
			m := map[string]string{
				"username": req.Username,
			}
			q := influx.NewQuery(dbutil.QueryHelper(query, m), "", "")
			response, err := cli.Query(q)
			result = multierror.Append(result, err)
			if response != nil {
				result = multierror.Append(result, response.Error())
			}
		}
	}
	if result.ErrorOrNil() != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("failed to delete user cleanly: %w", result.ErrorOrNil())
	}
	return dbplugin.DeleteUserResponse{}, nil
}

func (i *Influxdb) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password == nil && req.Expiration == nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("no changes requested")
	}

	i.Lock()
	defer i.Unlock()

	if req.Password != nil {
		err := i.changeUserPassword(ctx, req.Username, req.Password)
		if err != nil {
			return dbplugin.UpdateUserResponse{}, fmt.Errorf("failed to change %q password: %w", req.Username, err)
		}
	}
	// Expiration is a no-op
	return dbplugin.UpdateUserResponse{}, nil
}

func (i *Influxdb) changeUserPassword(ctx context.Context, username string, changePassword *dbplugin.ChangePassword) error {
	cli, err := i.getConnection(ctx)
	if err != nil {
		return fmt.Errorf("unable to get connection: %w", err)
	}

	rotateIFQL := changePassword.Statements.Commands
	if len(rotateIFQL) == 0 {
		rotateIFQL = []string{defaultRootCredentialRotationIFQL}
	}

	var result *multierror.Error
	for _, stmt := range rotateIFQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}
			m := map[string]string{
				"username": username,
				"password": changePassword.NewPassword,
			}
			q := influx.NewQuery(dbutil.QueryHelper(query, m), "", "")
			response, err := cli.Query(q)
			result = multierror.Append(result, err)
			if response != nil {
				result = multierror.Append(result, response.Error())
			}
		}
	}

	err = result.ErrorOrNil()
	if err != nil {
		return fmt.Errorf("failed to execute rotation queries: %w", err)
	}

	return nil
}
