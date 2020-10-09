package influxdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/database/helper/credsutil"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	defaultUserCreationIFQL           = `CREATE USER "{{username}}" WITH PASSWORD '{{password}}';`
	defaultUserDeletionIFQL           = `DROP USER "{{username}}";`
	defaultRootCredentialRotationIFQL = `SET PASSWORD FOR "{{username}}" = '{{password}}';`
	influxdbTypeName                  = "influxdb"
)

var _ newdbplugin.Database = &Influxdb{}

// Influxdb is an implementation of Database interface
type Influxdb struct {
	*influxdbConnectionProducer
}

// New returns a new Cassandra instance
func New() (interface{}, error) {
	db := new()
	dbType := newdbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)

	return dbType, nil
}

func new() *Influxdb {
	connProducer := &influxdbConnectionProducer{}
	connProducer.Type = influxdbTypeName

	return &Influxdb{
		influxdbConnectionProducer: connProducer,
	}
}

// Run instantiates a Influxdb object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	newdbplugin.Serve(dbType.(newdbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
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

// NewUser generates the username/password on the underlying Influxdb secret backend as instructed by
// the statements provided.
func (i *Influxdb) NewUser(ctx context.Context, req newdbplugin.NewUserRequest) (resp newdbplugin.NewUserResponse, err error) {
	i.Lock()
	defer i.Unlock()

	cli, err := i.getConnection(ctx)
	if err != nil {
		return newdbplugin.NewUserResponse{}, fmt.Errorf("unable to get connection: %w", err)
	}

	creationIFQL := req.Statements.Commands
	if len(creationIFQL) == 0 {
		creationIFQL = []string{defaultUserCreationIFQL}
	}

	rollbackIFQL := req.RollbackStatements.Commands
	if len(rollbackIFQL) == 0 {
		rollbackIFQL = []string{defaultUserDeletionIFQL}
	}

	username, err := credsutil.GenerateUsername(
		credsutil.DisplayName(req.UsernameConfig.DisplayName, 15),
		credsutil.RoleName(req.UsernameConfig.RoleName, 15),
		credsutil.MaxLength(100),
		credsutil.Separator("_"),
		credsutil.ToLower(),
	)
	if err != nil {
		return newdbplugin.NewUserResponse{}, fmt.Errorf("failed to generate username: %w", err)
	}
	username = strings.Replace(username, "-", "_", -1)

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
			q := influx.NewQuery(dbutil.QueryHelper(query, m), "", "")
			response, err := cli.Query(q)
			if err != nil {
				if response != nil && response.Error() != nil {
					attemptRollback(cli, username, rollbackIFQL)
				}
				return newdbplugin.NewUserResponse{}, fmt.Errorf("failed to run query in InfluxDB: %w", err)
			}
		}
	}
	resp = newdbplugin.NewUserResponse{
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
			if err != nil {
				if response != nil && response.Error() != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (i *Influxdb) DeleteUser(ctx context.Context, req newdbplugin.DeleteUserRequest) (newdbplugin.DeleteUserResponse, error) {
	i.Lock()
	defer i.Unlock()

	cli, err := i.getConnection(ctx)
	if err != nil {
		return newdbplugin.DeleteUserResponse{}, fmt.Errorf("unable to get connection: %w", err)
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
		return newdbplugin.DeleteUserResponse{}, fmt.Errorf("failed to delete user cleanly: %w", result.ErrorOrNil())
	}
	return newdbplugin.DeleteUserResponse{}, nil
}

func (i *Influxdb) UpdateUser(ctx context.Context, req newdbplugin.UpdateUserRequest) (newdbplugin.UpdateUserResponse, error) {
	i.Lock()
	defer i.Unlock()

	if req.Password != nil {
		err := i.changeUserPassword(ctx, req.Username, req.Password)
		if err != nil {
			return newdbplugin.UpdateUserResponse{}, fmt.Errorf("failed to change %q password: %w", req.Username, err)
		}
	}
	// Expiration is a no-op
	return newdbplugin.UpdateUserResponse{}, nil
}

func (i *Influxdb) changeUserPassword(ctx context.Context, username string, changePassword *newdbplugin.ChangePassword) error {
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
