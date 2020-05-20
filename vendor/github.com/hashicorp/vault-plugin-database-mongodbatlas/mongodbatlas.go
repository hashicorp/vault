package mongodbatlas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const mongoDBAtlasTypeName = "mongodbatlas"

// Verify interface is implemented
var _ dbplugin.Database = &MongoDBAtlas{}

type MongoDBAtlas struct {
	*mongoDBAtlasConnectionProducer
	credsutil.CredentialsProducer
}

func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func new() *MongoDBAtlas {
	connProducer := &mongoDBAtlasConnectionProducer{}
	connProducer.Type = mongoDBAtlasTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: credsutil.NoneLength,
		RoleNameLen:    15,
		UsernameLen:    20,
		Separator:      "-",
	}

	return &MongoDBAtlas{
		mongoDBAtlasConnectionProducer: connProducer,
		CredentialsProducer:            credsProducer,
	}
}

// Run instantiates a MongoDBAtlas object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

func (m *MongoDBAtlas) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
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
	var databaseUser mongoDBAtlasStatement
	err = json.Unmarshal([]byte(statements.Creation[0]), &databaseUser)
	if err != nil {
		return "", "", fmt.Errorf("Error unmarshalling statement %s", err)
	}

	// Default to "admin" if no db provided
	if databaseUser.DatabaseName == "" {
		databaseUser.DatabaseName = "admin"
	}

	if len(databaseUser.Roles) == 0 {
		return "", "", fmt.Errorf("roles array is required in creation statement")
	}

	databaseUserRequest := &mongodbatlas.DatabaseUser{
		Username:     username,
		Password:     password,
		DatabaseName: databaseUser.DatabaseName,
		Roles:        databaseUser.Roles,
	}

	_, _, err = client.DatabaseUsers.Create(ctx, m.ProjectID, databaseUserRequest)
	if err != nil {
		return "", "", err
	}
	return username, password, nil
}

// RenewUser is not supported on MongoDB, so this is a no-op.
func (m *MongoDBAtlas) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser drops the specified user from the authentication database. If none is provided
// in the revocation statement, the default "admin" authentication database will be assumed.
func (m *MongoDBAtlas) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	m.Lock()
	defer m.Unlock()

	statements = dbutil.StatementCompatibilityHelper(statements)

	client, err := m.getConnection(ctx)
	if err != nil {
		return err
	}

	_, err = client.DatabaseUsers.Delete(ctx, m.ProjectID, username)
	return err
}

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (m *MongoDBAtlas) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
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

	username = staticUser.Username
	password = staticUser.Password

	databaseUserRequest := &mongodbatlas.DatabaseUser{
		Password: password,
	}

	_, _, err = client.DatabaseUsers.Update(context.Background(), m.ProjectID, username, databaseUserRequest)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func (m *MongoDBAtlas) getConnection(ctx context.Context) (*mongodbatlas.Client, error) {
	client, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return client.(*mongodbatlas.Client), nil
}

// RotateRootCredentials is not currently supported on MongoDB
func (m *MongoDBAtlas) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	return nil, errors.New("root credential rotation is not currently implemented in this database secrets engine")
}

// Type returns the TypeName for this backend
func (m *MongoDBAtlas) Type() (string, error) {
	return mongoDBAtlasTypeName, nil
}

type mongoDBAtlasStatement struct {
	DatabaseName string              `json:"database_name"`
	Roles        []mongodbatlas.Role `json:"roles,omitempty"`
}
