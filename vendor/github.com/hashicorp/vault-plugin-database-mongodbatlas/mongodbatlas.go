package mongodbatlas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/atlas/mongodbatlas"
)

const mongoDBAtlasTypeName = "mongodbatlas"

// Verify interface is implemented
var _ dbplugin.Database = (*MongoDBAtlas)(nil)

type MongoDBAtlas struct {
	*mongoDBAtlasConnectionProducer
}

func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func new() *MongoDBAtlas {
	connProducer := &mongoDBAtlasConnectionProducer{
		Type: mongoDBAtlasTypeName,
	}

	return &MongoDBAtlas{
		mongoDBAtlasConnectionProducer: connProducer,
	}
}

func (m *MongoDBAtlas) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	m.Lock()
	defer m.Unlock()

	m.RawConfig = req.Config

	err := mapstructure.WeakDecode(req.Config, m.mongoDBAtlasConnectionProducer)
	if err != nil {
		return dbplugin.InitializeResponse{}, err
	}

	if len(m.PublicKey) == 0 {
		return dbplugin.InitializeResponse{}, errors.New("public Key is not set")
	}

	if len(m.PrivateKey) == 0 {
		return dbplugin.InitializeResponse{}, errors.New("private Key is not set")
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	m.Initialized = true

	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}

	return resp, nil
}

func (m *MongoDBAtlas) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	// Statement length checks
	if len(req.Statements.Commands) == 0 {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}
	if len(req.Statements.Commands) > 1 {
		return dbplugin.NewUserResponse{}, fmt.Errorf("only 1 creation statement supported for creation")
	}

	client, err := m.getConnection(ctx)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	username, err := credsutil.GenerateUsername(
		credsutil.DisplayName("", credsutil.NoneLength),
		credsutil.RoleName(req.UsernameConfig.RoleName, 15),
		credsutil.MaxLength(20),
		credsutil.Separator("-"),
	)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	// Unmarshal creation statements into mongodb roles
	var databaseUser mongoDBAtlasStatement
	err = json.Unmarshal([]byte(req.Statements.Commands[0]), &databaseUser)
	if err != nil {
		return dbplugin.NewUserResponse{}, fmt.Errorf("error unmarshalling statement %s", err)
	}

	// Default to "admin" if no db provided
	if databaseUser.DatabaseName == "" {
		databaseUser.DatabaseName = "admin"
	}

	if len(databaseUser.Roles) == 0 {
		return dbplugin.NewUserResponse{}, fmt.Errorf("roles array is required in creation statement")
	}

	databaseUserRequest := &mongodbatlas.DatabaseUser{
		Username:     username,
		Password:     req.Password,
		DatabaseName: databaseUser.DatabaseName,
		Roles:        databaseUser.Roles,
		Scopes:       databaseUser.Scopes,
	}

	_, _, err = client.DatabaseUsers.Create(ctx, m.ProjectID, databaseUserRequest)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	resp := dbplugin.NewUserResponse{
		Username: username,
	}

	return resp, nil
}

func (m *MongoDBAtlas) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	if req.Password != nil {
		err := m.changePassword(ctx, req.Username, req.Password.NewPassword)
		return dbplugin.UpdateUserResponse{}, err
	}

	// This also results in an no-op if the expiration is updated due to renewal.
	return dbplugin.UpdateUserResponse{}, nil
}

func (m *MongoDBAtlas) changePassword(ctx context.Context, username, password string) error {
	m.Lock()
	defer m.Unlock()

	client, err := m.getConnection(ctx)
	if err != nil {
		return err
	}

	databaseUserRequest := &mongodbatlas.DatabaseUser{
		Password: password,
	}

	_, _, err = client.DatabaseUsers.Update(context.Background(), m.ProjectID, username, databaseUserRequest)
	return err
}

func (m *MongoDBAtlas) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	m.Lock()
	defer m.Unlock()

	client, err := m.getConnection(ctx)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, err
	}

	var databaseUser mongoDBAtlasStatement
	if len(req.Statements.Commands) > 0 {
		err = json.Unmarshal([]byte(req.Statements.Commands[0]), &databaseUser)
		if err != nil {
			return dbplugin.DeleteUserResponse{}, fmt.Errorf("error unmarshalling statement %w", err)
		}
	}

	// Default to "admin" if no db provided
	if databaseUser.DatabaseName == "" {
		databaseUser.DatabaseName = "admin"
	}

	_, err = client.DatabaseUsers.Delete(ctx, databaseUser.DatabaseName, m.ProjectID, req.Username)
	return dbplugin.DeleteUserResponse{}, err
}

func (m *MongoDBAtlas) getConnection(ctx context.Context) (*mongodbatlas.Client, error) {
	client, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return client.(*mongodbatlas.Client), nil
}

// Type returns the TypeName for this backend
func (m *MongoDBAtlas) Type() (string, error) {
	return mongoDBAtlasTypeName, nil
}

type mongoDBAtlasStatement struct {
	DatabaseName string               `json:"database_name"`
	Roles        []mongodbatlas.Role  `json:"roles,omitempty"`
	Scopes       []mongodbatlas.Scope `json:"scopes,omitempty"`
}
