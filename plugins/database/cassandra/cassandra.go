package cassandra

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const (
	defaultCreationCQL = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;`
	defaultRollbackCQL = `DROP USER '{{username}}';`
	cassandraTypeName  = "cassandra"
)

type Cassandra struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

func New() *Cassandra {
	connProducer := &connutil.CassandraConnectionProducer{}
	connProducer.Type = cassandraTypeName

	credsProducer := &credsutil.CassandraCredentialsProducer{}

	dbType := &Cassandra{
		ConnectionProducer:  connProducer,
		CredentialsProducer: credsProducer,
	}

	return dbType
}

// Run instantiates a MySQL object, and runs the RPC server for the plugin
func Run() error {
	dbType := New()

	dbplugin.NewPluginServer(dbType)

	return nil
}

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

func (c *Cassandra) CreateUser(statements dbplugin.Statements, usernamePrefix string, expiration time.Time) (username string, password string, err error) {
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
		creationCQL = defaultCreationCQL
	}
	rollbackCQL := statements.RollbackStatements
	if rollbackCQL == "" {
		rollbackCQL = defaultRollbackCQL
	}

	username, err = c.GenerateUsername(usernamePrefix)
	if err != nil {
		return "", "", err
	}

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
					"password": password,
				})).Exec()
			}
			return "", "", err
		}
	}

	return username, password, nil
}

func (c *Cassandra) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

func (c *Cassandra) RevokeUser(statements dbplugin.Statements, username string) error {
	// Grab the lock
	c.Lock()
	defer c.Unlock()

	session, err := c.getConnection()
	if err != nil {
		return err
	}

	err = session.Query(fmt.Sprintf("DROP USER '%s'", username)).Exec()
	if err != nil {
		return fmt.Errorf("error removing user '%s': %s", username, err)
	}

	return nil
}
