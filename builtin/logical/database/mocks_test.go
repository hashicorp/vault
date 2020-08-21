package database

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	"github.com/stretchr/testify/mock"
)

var _ newdbplugin.Database = &mockNewDatabase{}

type mockNewDatabase struct {
	mock.Mock
}

func (m mockNewDatabase) Initialize(ctx context.Context, req newdbplugin.InitializeRequest) (newdbplugin.InitializeResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(newdbplugin.InitializeResponse), args.Error(1)
}

func (m mockNewDatabase) NewUser(ctx context.Context, req newdbplugin.NewUserRequest) (newdbplugin.NewUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(newdbplugin.NewUserResponse), args.Error(1)
}

func (m mockNewDatabase) UpdateUser(ctx context.Context, req newdbplugin.UpdateUserRequest) (newdbplugin.UpdateUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(newdbplugin.UpdateUserResponse), args.Error(1)
}

func (m mockNewDatabase) DeleteUser(ctx context.Context, req newdbplugin.DeleteUserRequest) (newdbplugin.DeleteUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(newdbplugin.DeleteUserResponse), args.Error(1)
}

func (m mockNewDatabase) Type() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m mockNewDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

var _ dbplugin.Database = &mockLegacyDatabase{}

type mockLegacyDatabase struct {
	mock.Mock
}

func (m *mockLegacyDatabase) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	args := m.Called(ctx, statements, usernameConfig, expiration)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockLegacyDatabase) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	args := m.Called(ctx, statements, username, expiration)
	return args.Error(0)
}

func (m *mockLegacyDatabase) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	args := m.Called(ctx, statements, username)
	return args.Error(0)
}

func (m *mockLegacyDatabase) RotateRootCredentials(ctx context.Context, statements []string) (config map[string]interface{}, err error) {
	args := m.Called(ctx, statements)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockLegacyDatabase) GenerateCredentials(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *mockLegacyDatabase) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticConfig dbplugin.StaticUserConfig) (username string, password string, err error) {
	args := m.Called(ctx, statements, staticConfig)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockLegacyDatabase) Init(ctx context.Context, config map[string]interface{}, verifyConnection bool) (saveConfig map[string]interface{}, err error) {
	args := m.Called(ctx, config, verifyConnection)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockLegacyDatabase) Type() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockLegacyDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockLegacyDatabase) Initialize(ctx context.Context, config map[string]interface{}, verifyConnection bool) (err error) {
	panic("Initialize should not be called")
}
