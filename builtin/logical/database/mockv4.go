// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"fmt"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	v4 "github.com/hashicorp/vault/sdk/database/dbplugin"
)

const mockV4Type = "mockv4"

// MockDatabaseV4 is an implementation of Database interface
type MockDatabaseV4 struct {
	config map[string]interface{}
}

var _ v4.Database = &MockDatabaseV4{}

// New returns a new in-memory instance
func NewV4() (interface{}, error) {
	return MockDatabaseV4{}, nil
}

// RunV4 instantiates a MongoDB object, and runs the RPC server for the plugin
func RunV4(apiTLSConfig *api.TLSConfig) error {
	dbType, err := NewV4()
	if err != nil {
		return err
	}

	v4.Serve(dbType.(v4.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

func (m MockDatabaseV4) Init(ctx context.Context, config map[string]interface{}, verifyConnection bool) (saveConfig map[string]interface{}, err error) {
	log.Default().Info("Init called",
		"config", config,
		"verifyConnection", verifyConnection)

	return config, nil
}

func (m MockDatabaseV4) Initialize(ctx context.Context, config map[string]interface{}, verifyConnection bool) (err error) {
	_, err = m.Init(ctx, config, verifyConnection)
	return err
}

func (m MockDatabaseV4) CreateUser(ctx context.Context, statements v4.Statements, usernameConfig v4.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	log.Default().Info("CreateUser called",
		"statements", statements,
		"usernameConfig", usernameConfig,
		"expiration", expiration)

	now := time.Now()
	user := fmt.Sprintf("mockv4_user_%s", now.Format(time.RFC3339))
	pass, err := m.GenerateCredentials(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate credentials: %w", err)
	}
	return user, pass, nil
}

func (m MockDatabaseV4) RenewUser(ctx context.Context, statements v4.Statements, username string, expiration time.Time) error {
	log.Default().Info("RenewUser called",
		"statements", statements,
		"username", username,
		"expiration", expiration)

	return nil
}

func (m MockDatabaseV4) RevokeUser(ctx context.Context, statements v4.Statements, username string) error {
	log.Default().Info("RevokeUser called",
		"statements", statements,
		"username", username)

	return nil
}

func (m MockDatabaseV4) RotateRootCredentials(ctx context.Context, statements []string) (config map[string]interface{}, err error) {
	log.Default().Info("RotateRootCredentials called",
		"statements", statements)

	newPassword, err := m.GenerateCredentials(ctx)
	if err != nil {
		return config, fmt.Errorf("failed to generate credentials: %w", err)
	}
	config["password"] = newPassword

	return m.config, nil
}

func (m MockDatabaseV4) SetCredentials(ctx context.Context, statements v4.Statements, staticConfig v4.StaticUserConfig) (username string, password string, err error) {
	log.Default().Info("SetCredentials called",
		"statements", statements,
		"staticConfig", staticConfig)
	return "", "", nil
}

func (m MockDatabaseV4) GenerateCredentials(ctx context.Context) (password string, err error) {
	now := time.Now()
	pass := fmt.Sprintf("mockv4_password_%s", now.Format(time.RFC3339))
	return pass, nil
}

func (m MockDatabaseV4) Type() (string, error) {
	log.Default().Info("Type called")
	return mockV4Type, nil
}

func (m MockDatabaseV4) Close() error {
	log.Default().Info("Close called")
	return nil
}
