package mongodb

import (
	"io"
	"strings"
	"time"

	"encoding/json"

	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/hashicorp/vault/plugins/helper/database/dbutil"
	"gopkg.in/mgo.v2"
)

const mongoDBTypeName = "mongodb"

// MongoDB is an implementation of Database interface
type MongoDB struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

// New returns a new MongoDB instance
func New() (interface{}, error) {
	connProducer := &mongoDBConnectionProducer{}
	connProducer.Type = mongoDBTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 15,
		RoleNameLen:    15,
		UsernameLen:    100,
		Separator:      "-",
	}

	dbType := &MongoDB{
		ConnectionProducer:  connProducer,
		CredentialsProducer: credsProducer,
	}
	return dbType, nil
}

// Run instantiates a MongoDB object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(*MongoDB), apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (m *MongoDB) Type() (string, error) {
	return mongoDBTypeName, nil
}

func (m *MongoDB) getConnection() (*mgo.Session, error) {
	session, err := m.Connection()
	if err != nil {
		return nil, err
	}

	return session.(*mgo.Session), nil
}

// CreateUser generates the username/password on the underlying secret backend as instructed by
// the CreationStatement provided. The creation statement is a JSON blob that has a db value,
// and an array of roles that accepts a role, and an optional db value pair. This array will
// be normalized the format specified in the mongoDB docs:
// https://docs.mongodb.com/manual/reference/command/createUser/#dbcmd.createUser
//
// JSON Example:
//  { "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "foo"}] }
func (m *MongoDB) CreateUser(statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	if statements.CreationStatements == "" {
		return "", "", dbutil.ErrEmptyCreationStatement
	}

	session, err := m.getConnection()
	if err != nil {
		return "", "", err
	}

	username, err = m.GenerateUsername(usernameConfig)
	if err != nil {
		return "", "", err
	}

	password, err = m.GeneratePassword()
	if err != nil {
		return "", "", err
	}

	// Unmarshal statements.CreationStatements into mongodbRoles
	var mongoCS mongoDBStatement
	err = json.Unmarshal([]byte(statements.CreationStatements), &mongoCS)
	if err != nil {
		return "", "", err
	}

	// Default to "admin" if no db provided
	if mongoCS.DB == "" {
		mongoCS.DB = "admin"
	}

	if len(mongoCS.Roles) == 0 {
		return "", "", fmt.Errorf("roles array is required in creation statement")
	}

	createUserCmd := createUserCommand{
		Username: username,
		Password: password,
		Roles:    mongoCS.Roles.toStandardRolesArray(),
	}

	err = session.DB(mongoCS.DB).Run(createUserCmd, nil)
	switch {
	case err == nil:
	case err == io.EOF, strings.Contains(err.Error(), "EOF"):
		if err := m.ConnectionProducer.Close(); err != nil {
			return "", "", errwrap.Wrapf("error closing EOF'd mongo connection: {{err}}", err)
		}
		session, err := m.getConnection()
		if err != nil {
			return "", "", err
		}
		err = session.DB(mongoCS.DB).Run(createUserCmd, nil)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", err
	}

	return username, password, nil
}

// RenewUser is not supported on MongoDB, so this is a no-op.
func (m *MongoDB) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser drops the specified user from the authentication databse. If none is provided
// in the revocation statement, the default "admin" authentication database will be assumed.
func (m *MongoDB) RevokeUser(statements dbplugin.Statements, username string) error {
	session, err := m.getConnection()
	if err != nil {
		return err
	}

	// If no revocation statements provided, pass in empty JSON
	revocationStatement := statements.RevocationStatements
	if revocationStatement == "" {
		revocationStatement = `{}`
	}

	// Unmarshal revocation statements into mongodbRoles
	var mongoCS mongoDBStatement
	err = json.Unmarshal([]byte(revocationStatement), &mongoCS)
	if err != nil {
		return err
	}

	db := mongoCS.DB
	// If db is not specified, use the default authenticationDatabase "admin"
	if db == "" {
		db = "admin"
	}

	err = session.DB(db).RemoveUser(username)
	switch {
	case err == nil, err == mgo.ErrNotFound:
	case err == io.EOF, strings.Contains(err.Error(), "EOF"):
		if err := m.ConnectionProducer.Close(); err != nil {
			return errwrap.Wrapf("error closing EOF'd mongo connection: {{err}}", err)
		}
		session, err := m.getConnection()
		if err != nil {
			return err
		}
		err = session.DB(db).RemoveUser(username)
		if err != nil {
			return err
		}
	default:
		return err
	}

	return nil
}
