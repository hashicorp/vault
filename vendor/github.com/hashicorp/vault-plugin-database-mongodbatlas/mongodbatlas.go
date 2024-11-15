// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	mongoDBAtlasTypeName    = "mongodbatlas"
	defaultUserNameTemplate = `{{ printf "v-%s-%s" (.RoleName | truncate 15) (random 20) | truncate 20 }}`
)

// Verify interface is implemented
var _ dbplugin.Database = (*MongoDBAtlas)(nil)

type MongoDBAtlas struct {
	*mongoDBAtlasConnectionProducer

	usernameProducer template.StringTemplate
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
	m.usernameProducer = up

	_, err = m.usernameProducer.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid username template: %w", err)
	}

	err = m.mongoDBAtlasConnectionProducer.Initialize(ctx, req)
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("failed to initialize: %w", err)
	}

	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}
	resp.SetSupportedCredentialTypes([]dbplugin.CredentialType{
		dbplugin.CredentialTypePassword,
		dbplugin.CredentialTypeClientCertificate,
	})
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

	var username string
	switch req.CredentialType {
	case dbplugin.CredentialTypePassword:
		username, err = m.usernameProducer.Generate(req.UsernameConfig)
		if err != nil {
			return dbplugin.NewUserResponse{}, err
		}
	case dbplugin.CredentialTypeClientCertificate:
		// MongoDb Atlas expects the username to equal the client certificate subject
		// https://www.mongodb.com/docs/manual/tutorial/configure-x509-client-authentication/
		username = req.Subject
	default:
		return dbplugin.NewUserResponse{}, fmt.Errorf("unsupported credential type %q",
			req.CredentialType)
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
		X509Type:     databaseUser.X509Type,
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

	// If the user is an X.509 user, delete the user from the $external database
	if isX509User(req.Username) {
		if databaseUser.DatabaseName == "" {
			databaseUser.DatabaseName = "$external"
		}
	} else {
		// If the user is not an X.509 user, delete the user from the MongoDB Atlas project
		if databaseUser.DatabaseName == "" {
			databaseUser.DatabaseName = "admin"
		}
	}

	_, err = client.DatabaseUsers.Delete(ctx, databaseUser.DatabaseName, m.ProjectID, req.Username)
	if err != nil {
		return dbplugin.DeleteUserResponse{}, fmt.Errorf("error deleting user from project: %w", err)
	}
	return dbplugin.DeleteUserResponse{}, nil
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

// Check to see if the user is a X509 user
func isX509User(username string) bool {
	return strings.HasPrefix(username, "CN=")
}

type mongoDBAtlasStatement struct {
	DatabaseName string               `json:"database_name"`
	Roles        []mongodbatlas.Role  `json:"roles,omitempty"`
	Scopes       []mongodbatlas.Scope `json:"scopes,omitempty"`
	X509Type     string               `json:"x509Type,omitempty"`
}
