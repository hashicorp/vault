// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"context"
	"fmt"
	"time"

	log "github.com/hashicorp/go-hclog"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

const mockV5Type = "mockv5"

// MockDatabaseV5 is an implementation of Database interface
type MockDatabaseV5 struct {
	config map[string]interface{}
}

var _ v5.Database = &MockDatabaseV5{}

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

	v5.Serve(dbType.(v5.Database))

	return nil
}

// Run instantiates a MongoDB object, and runs the RPC server for the plugin
func RunV6Multiplexed() error {
	v5.ServeMultiplex(New)

	return nil
}

func (m MockDatabaseV5) Initialize(ctx context.Context, req v5.InitializeRequest) (v5.InitializeResponse, error) {
	log.Default().Info("Initialize called",
		"req", req)

	config := req.Config
	config["from-plugin"] = "this value is from the plugin itself"

	resp := v5.InitializeResponse{
		Config: req.Config,
	}
	return resp, nil
}

func (m MockDatabaseV5) NewUser(ctx context.Context, req v5.NewUserRequest) (v5.NewUserResponse, error) {
	log.Default().Info("NewUser called",
		"req", req)

	now := time.Now()
	user := fmt.Sprintf("mockv5_user_%s", now.Format(time.RFC3339))
	resp := v5.NewUserResponse{
		Username: user,
	}
	return resp, nil
}

func (m MockDatabaseV5) UpdateUser(ctx context.Context, req v5.UpdateUserRequest) (v5.UpdateUserResponse, error) {
	log.Default().Info("UpdateUser called",
		"req", req)
	return v5.UpdateUserResponse{}, nil
}

func (m MockDatabaseV5) DeleteUser(ctx context.Context, req v5.DeleteUserRequest) (v5.DeleteUserResponse, error) {
	log.Default().Info("DeleteUser called",
		"req", req)
	return v5.DeleteUserResponse{}, nil
}

func (m MockDatabaseV5) Type() (string, error) {
	log.Default().Info("Type called")
	return mockV5Type, nil
}

func (m MockDatabaseV5) Close() error {
	log.Default().Info("Close called")
	return nil
}
