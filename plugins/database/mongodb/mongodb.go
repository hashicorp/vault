package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const mongoDBTypeName = "mongodb"

// MongoDB is an implementation of Database interface
type MongoDB struct {
	*mongoDBConnectionProducer
	credsutil.CredentialsProducer
}

var _ dbplugin.Database = &MongoDB{}

// New returns a new MongoDB instance
func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func new() *MongoDB {
	connProducer := &mongoDBConnectionProducer{}
	connProducer.Type = mongoDBTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 15,
		RoleNameLen:    15,
		UsernameLen:    100,
		Separator:      "-",
	}

	return &MongoDB{
		mongoDBConnectionProducer: connProducer,
		CredentialsProducer:       credsProducer,
	}
}

// Run instantiates a MongoDB object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

// Type returns the TypeName for this backend
func (m *MongoDB) Type() (string, error) {
	return mongoDBTypeName, nil
}

func (m *MongoDB) getConnection(ctx context.Context) (*mongo.Client, error) {
	client, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return client.(*mongo.Client), nil
}

// CreateUser generates the username/password on the underlying secret backend as instructed by
// the CreationStatement provided. The creation statement is a JSON blob that has a db value,
// and an array of roles that accepts a role, and an optional db value pair. This array will
// be normalized the format specified in the mongoDB docs:
// https://docs.mongodb.com/manual/reference/command/createUser/#dbcmd.createUser
//
// JSON Example:
//  { "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "foo"}] }
func (m *MongoDB) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Creation) == 0 {
		return "", "", dbutil.ErrEmptyCreationStatement
	}

	client, err := m.getConnection(ctx)
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
	err = json.Unmarshal([]byte(statements.Creation[0]), &mongoCS)
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

	result := client.Database(mongoCS.DB).RunCommand(ctx, createUserCmd, nil)
	switch {
	case result.Err() == nil:
	case result.Err() == io.EOF, strings.Contains(result.Err().Error(), "EOF"):
		// Call getConnection to reset and retry query if we get an ErrNoDocuments error on first attempt.
		client, err := m.getConnection(ctx)
		if err != nil {
			return "", "", err
		}
		result = client.Database(mongoCS.DB).RunCommand(ctx, createUserCmd, nil)
		if result.Err() != nil {
			return "", "", err
		}
	default:
		return "", "", err
	}

	return username, password, nil
}

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (m *MongoDB) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	client, err := m.getConnection(ctx)
	if err != nil {
		return "", "", err
	}

	username = staticUser.Username
	password = staticUser.Password

	changeUserCmd := &updateUserCommand{
		Username: username,
		Password: password,
	}

	cs, err := connstring.Parse(m.ConnectionURL)
	if err != nil {
		return "", "", err
	}
	result := client.Database(cs.Database).RunCommand(ctx, changeUserCmd, nil)

	switch {
	case result.Err() == nil:
	case result.Err() == io.EOF, strings.Contains(result.Err().Error(), "EOF"):
		// Call getConnection to reset and retry query if we get an ErrNoDocuments error on first attempt.
		client, err := m.getConnection(ctx)
		if err != nil {
			return "", "", err
		}
		result = client.Database(cs.Database).RunCommand(ctx, changeUserCmd, nil)
		if result.Err() != nil {
			return "", "", err
		}
	default:
		return "", "", err
	}

	return username, password, nil
}

// RenewUser is not supported on MongoDB, so this is a no-op.
func (m *MongoDB) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser drops the specified user from the authentication database. If none is provided
// in the revocation statement, the default "admin" authentication database will be assumed.
func (m *MongoDB) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	m.Lock()
	defer m.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	client, err := m.getConnection(ctx)
	if err != nil {
		return err
	}

	// If no revocation statements provided, pass in empty JSON
	var revocationStatement string
	switch len(statements.Revocation) {
	case 0:
		revocationStatement = `{}`
	case 1:
		revocationStatement = statements.Revocation[0]
	default:
		return fmt.Errorf("expected 0 or 1 revocation statements, got %d", len(statements.Revocation))
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

	dropUserCmd := &dropUserCommand{
		Username:     username,
		WriteConcern: writeconcern.New(writeconcern.WMajority()),
	}

	result := client.Database(db).RunCommand(ctx, dropUserCmd, nil)
	switch {
	case result.Err() == nil:
	case result.Err() == io.EOF, strings.Contains(result.Err().Error(), "EOF"):
		if err := m.Close(); err != nil {
			return errwrap.Wrapf("error closing EOF'd mongo connection: {{err}}", err)
		}
		client, err := m.getConnection(ctx)
		if err != nil {
			return err
		}
		result = client.Database(db).RunCommand(ctx, dropUserCmd, nil)
		if result.Err() != nil {
			return err
		}
	default:
		return err
	}

	return nil
}

// RotateRootCredentials is not currently supported on MongoDB
func (m *MongoDB) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	return nil, errors.New("root credential rotation is not currently implemented in this database secrets engine")
}
