package influxdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	defaultUserCreationIFQL           = `CREATE USER "{{username}}" WITH PASSWORD '{{password}}';`
	defaultUserDeletionIFQL           = `DROP USER "{{username}}";`
	defaultRootCredentialRotationIFQL = `SET PASSWORD FOR "{{username}}" = '{{password}}';`
	influxdbTypeName                  = "influxdb"
)

var _ dbplugin.Database = &Influxdb{}

// Influxdb is an implementation of Database interface
type Influxdb struct {
	*influxdbConnectionProducer
	credsutil.CredentialsProducer
}

// New returns a new Cassandra instance
func New() (interface{}, error) {
	db, err := new()
	if err != nil {
		return nil, err
	}
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)

	return dbType, nil
}

func new() (*Influxdb, error) {
	connProducer := &influxdbConnectionProducer{}
	connProducer.Type = influxdbTypeName

	credsProducer, err := credsutil.NewUsernamePasswordProducer(
		credsutil.UsernameOpts(
			credsutil.UsernameTemplate("v_{{.DisplayName | truncate 15}}_{{.RoleName | truncate 15}}_{{rand 20}}_{{now_seconds}}"),
			credsutil.UsernameMaxLength(100),
		),
		credsutil.PasswordOpts(
			credsutil.PasswordLength(20),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create username & password producer: %w", err)
	}

	idb := &Influxdb{
		influxdbConnectionProducer: connProducer,
		CredentialsProducer:        credsProducer,
	}
	return idb, nil
}

// Run instantiates a Influxdb object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

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

// CreateUser generates the username/password on the underlying Influxdb secret backend as instructed by
// the CreationStatement provided.
func (i *Influxdb) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	i.Lock()
	defer i.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	// Get the connection
	cli, err := i.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	creationIFQL := statements.Creation
	if len(creationIFQL) == 0 {
		creationIFQL = []string{defaultUserCreationIFQL}
	}

	rollbackIFQL := statements.Rollback
	if len(rollbackIFQL) == 0 {
		rollbackIFQL = []string{defaultUserDeletionIFQL}
	}

	username, err = i.GenerateUsername(usernameConfig)
	username = strings.Replace(username, "-", "_", -1)
	if err != nil {
		return "", "", err
	}
	username = strings.ToLower(username)
	password, err = i.GeneratePassword()
	if err != nil {
		return "", "", err
	}

	// Execute each query
	for _, stmt := range creationIFQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			q := influx.NewQuery(dbutil.QueryHelper(query, map[string]string{
				"username": username,
				"password": password,
			}), "", "")
			response, err := cli.Query(q)
			if err != nil {
				if response != nil && response.Error() != nil {
					attemptRollback(cli, username, rollbackIFQL)
				}
				return "", "", err
			}
		}
	}
	return username, password, nil
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

// RenewUser is not supported on Influxdb, so this is a no-op.
func (i *Influxdb) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser attempts to drop the specified user.
func (i *Influxdb) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the lock
	i.Lock()
	defer i.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	cli, err := i.getConnection(ctx)
	if err != nil {
		return err
	}

	revocationIFQL := statements.Revocation
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
			q := influx.NewQuery(dbutil.QueryHelper(query, map[string]string{
				"username": username,
			}), "", "")
			response, err := cli.Query(q)
			result = multierror.Append(result, err)
			if response != nil {
				result = multierror.Append(result, response.Error())
			}
		}
	}
	return result.ErrorOrNil()
}

// RotateRootCredentials is useful when we try to change root credential
func (i *Influxdb) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	// Grab the lock
	i.Lock()
	defer i.Unlock()

	cli, err := i.getConnection(ctx)
	if err != nil {
		return nil, err
	}

	rotateIFQL := statements
	if len(rotateIFQL) == 0 {
		rotateIFQL = []string{defaultRootCredentialRotationIFQL}
	}

	password, err := i.GeneratePassword()
	if err != nil {
		return nil, err
	}

	var result *multierror.Error
	for _, stmt := range rotateIFQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}
			q := influx.NewQuery(dbutil.QueryHelper(query, map[string]string{
				"username": i.Username,
				"password": password,
			}), "", "")
			response, err := cli.Query(q)
			result = multierror.Append(result, err)
			if response != nil {
				result = multierror.Append(result, response.Error())
			}
		}
	}

	err = result.ErrorOrNil()
	if err != nil {
		return nil, err
	}

	i.rawConfig["password"] = password
	return i.rawConfig, nil
}
