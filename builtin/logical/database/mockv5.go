package database

import (
	"context"
	"fmt"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
)

const mockV5Type = "mockv5"

// MockDatabaseV5 is an implementation of Database interface
type MockDatabaseV5 struct {
	config map[string]interface{}
}

var _ newdbplugin.Database = &MockDatabaseV5{}

// New returns a new in-memory instance
func New() (interface{}, error) {
	db := MockDatabaseV5{}
	return db, nil
}

// Run instantiates a MongoDB object, and runs the RPC server for the plugin
func RunV5() error {
	dbType, err := New()
	if err != nil {
		return err
	}

	newdbplugin.Serve(dbType.(newdbplugin.Database))

	return nil
}

func (m MockDatabaseV5) Initialize(ctx context.Context, req newdbplugin.InitializeRequest) (newdbplugin.InitializeResponse, error) {
	log.Default().Info("Initialize called",
		"req", req)

	config := req.Config
	config["from-plugin"] = "this value is from the plugin itself"

	resp := newdbplugin.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

func (m MockDatabaseV5) NewUser(ctx context.Context, req newdbplugin.NewUserRequest) (newdbplugin.NewUserResponse, error) {
	log.Default().Info("NewUser called",
		"req", req)

	now := time.Now()
	user := fmt.Sprintf("mockv5_user_%s", now.Format(time.RFC3339))
	resp := newdbplugin.NewUserResponse{
		Username: user,
	}
	return resp, nil
}

func (m MockDatabaseV5) UpdateUser(ctx context.Context, req newdbplugin.UpdateUserRequest) (newdbplugin.UpdateUserResponse, error) {
	log.Default().Info("UpdateUser called",
		"req", req)
	return newdbplugin.UpdateUserResponse{}, nil
}

func (m MockDatabaseV5) DeleteUser(ctx context.Context, req newdbplugin.DeleteUserRequest) (newdbplugin.DeleteUserResponse, error) {
	log.Default().Info("DeleteUser called",
		"req", req)
	return newdbplugin.DeleteUserResponse{}, nil
}

func (m MockDatabaseV5) Type() (string, error) {
	log.Default().Info("Type called")
	return mockV5Type, nil
}

func (m MockDatabaseV5) Close() error {
	log.Default().Info("Close called")
	return nil
}
