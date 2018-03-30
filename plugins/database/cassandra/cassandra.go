package cassandra

import (
	"context"
	"strings"
	"time"

	"github.com/gocql/gocql"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const (
	defaultUserCreationCQL           = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;`
	defaultUserDeletionCQL           = `DROP USER '{{username}}';`
	defaultRootCredentialRotationCQL = `ALTER USER {{username}} WITH PASSWORD '{{password}}';`
	cassandraTypeName                = "cassandra"
)

var _ dbplugin.Database = &Cassandra{}

// Cassandra is an implementation of Database interface
type Cassandra struct {
	*cassandraConnectionProducer
	credsutil.CredentialsProducer
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

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 15,
		RoleNameLen:    15,
		UsernameLen:    100,
		Separator:      "_",
	}

	return &Cassandra{
		cassandraConnectionProducer: connProducer,
		CredentialsProducer:         credsProducer,
	}
}

// Run instantiates a Cassandra object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(dbplugin.Database), apiTLSConfig)

	return nil
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

// CreateUser generates the username/password on the underlying Cassandra secret backend as instructed by
// the CreationStatement provided.
func (c *Cassandra) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	// Get the connection
	session, err := c.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	creationCQL := statements.Creation
	if len(creationCQL) == 0 {
		creationCQL = []string{defaultUserCreationCQL}
	}

	rollbackCQL := statements.Rollback
	if len(rollbackCQL) == 0 {
		rollbackCQL = []string{defaultUserDeletionCQL}
	}

	username, err = c.GenerateUsername(usernameConfig)
	username = strings.Replace(username, "-", "_", -1)
	if err != nil {
		return "", "", err
	}
	// Cassandra doesn't like the uppercase usernames
	username = strings.ToLower(username)

	password, err = c.GeneratePassword()
	if err != nil {
		return "", "", err
	}

	// Execute each query
	for _, stmt := range creationCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			err = session.Query(dbutil.QueryHelper(query, map[string]string{
				"username": username,
				"password": password,
			})).Exec()
			if err != nil {
				for _, stmt := range rollbackCQL {
					for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
						query = strings.TrimSpace(query)
						if len(query) == 0 {
							continue
						}

						session.Query(dbutil.QueryHelper(query, map[string]string{
							"username": username,
						})).Exec()
					}
				}
				return "", "", err
			}
		}
	}

	return username, password, nil
}

// RenewUser is not supported on Cassandra, so this is a no-op.
func (c *Cassandra) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser attempts to drop the specified user.
func (c *Cassandra) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	session, err := c.getConnection(ctx)
	if err != nil {
		return err
	}

	revocationCQL := statements.Revocation
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

			err := session.Query(dbutil.QueryHelper(query, map[string]string{
				"username": username,
			})).Exec()

			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func (c *Cassandra) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	session, err := c.getConnection(ctx)
	if err != nil {
		return nil, err
	}

	rotateCQL := statements
	if len(rotateCQL) == 0 {
		rotateCQL = []string{defaultRootCredentialRotationCQL}
	}

	password, err := c.GeneratePassword()
	if err != nil {
		return nil, err
	}

	var result *multierror.Error
	for _, stmt := range rotateCQL {
		for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
			query = strings.TrimSpace(query)
			if len(query) == 0 {
				continue
			}

			err := session.Query(dbutil.QueryHelper(query, map[string]string{
				"username": c.Username,
				"password": password,
			})).Exec()

			result = multierror.Append(result, err)
		}
	}

	err = result.ErrorOrNil()
	if err != nil {
		return nil, err
	}

	c.rawConfig["password"] = password
	return c.rawConfig, nil
}
