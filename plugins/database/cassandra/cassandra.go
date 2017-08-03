package cassandra

import (
	"strings"
	"time"

	"github.com/gocql/gocql"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const (
	defaultUserCreationCQL = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;`
	defaultUserDeletionCQL = `DROP USER '{{username}}';`
	cassandraTypeName      = "cassandra"
)

// Cassandra is an implementation of Database interface
type Cassandra struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

// New returns a new Cassandra instance
func New() (interface{}, error) {
	connProducer := &cassandraConnectionProducer{}
	connProducer.Type = cassandraTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 15,
		RoleNameLen:    15,
		UsernameLen:    100,
		Separator:      "_",
	}

	dbType := &Cassandra{
		ConnectionProducer:  connProducer,
		CredentialsProducer: credsProducer,
	}

	return dbType, nil
}

// Run instantiates a Cassandra object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(*Cassandra), apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (c *Cassandra) Type() (string, error) {
	return cassandraTypeName, nil
}

func (c *Cassandra) getConnection() (*gocql.Session, error) {
	session, err := c.Connection()
	if err != nil {
		return nil, err
	}

	return session.(*gocql.Session), nil
}

// CreateUser generates the username/password on the underlying Cassandra secret backend as instructed by
// the CreationStatement provided.
func (c *Cassandra) CreateUser(statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	// Get the connection
	session, err := c.getConnection()
	if err != nil {
		return "", "", err
	}

	creationCQL := statements.CreationStatements
	if creationCQL == "" {
		creationCQL = defaultUserCreationCQL
	}
	rollbackCQL := statements.RollbackStatements
	if rollbackCQL == "" {
		rollbackCQL = defaultUserDeletionCQL
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
	for _, query := range strutil.ParseArbitraryStringSlice(creationCQL, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		err = session.Query(dbutil.QueryHelper(query, map[string]string{
			"username": username,
			"password": password,
		})).Exec()
		if err != nil {
			for _, query := range strutil.ParseArbitraryStringSlice(rollbackCQL, ";") {
				query = strings.TrimSpace(query)
				if len(query) == 0 {
					continue
				}

				session.Query(dbutil.QueryHelper(query, map[string]string{
					"username": username,
				})).Exec()
			}
			return "", "", err
		}
	}

	return username, password, nil
}

// RenewUser is not supported on Cassandra, so this is a no-op.
func (c *Cassandra) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser attempts to drop the specified user.
func (c *Cassandra) RevokeUser(statements dbplugin.Statements, username string) error {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	session, err := c.getConnection()
	if err != nil {
		return err
	}

	revocationCQL := statements.RevocationStatements
	if revocationCQL == "" {
		revocationCQL = defaultUserDeletionCQL
	}

	var result *multierror.Error
	for _, query := range strutil.ParseArbitraryStringSlice(revocationCQL, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		err := session.Query(dbutil.QueryHelper(query, map[string]string{
			"username": username,
		})).Exec()

		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}
