package database

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/mock"
)

func TestInitDatabase_legacy(t *testing.T) {
	type testCase struct {
		legacyInitConfig map[string]interface{}
		legacyInitErr    error

		expectedConfig map[string]interface{}
		expectErr      bool
	}

	tests := map[string]testCase{
		"legacy database error": {
			legacyInitErr: fmt.Errorf("test error"),

			expectedConfig: nil,
			expectErr:      true,
		},
		"legacy database success": {
			legacyInitConfig: map[string]interface{}{
				"foo": "bar",
			},

			expectedConfig: map[string]interface{}{
				"foo": "bar",
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			legacyDB := new(mockLegacyDatabase)
			legacyDB.On("Init", mock.Anything, mock.Anything, mock.Anything).
				Return(test.legacyInitConfig, test.legacyInitErr).
				Once()
			defer legacyDB.AssertExpectations(t)

			dbw := databaseVersionWrapper{
				legacyDatabase: legacyDB,
			}
			config, err := initDatabase(context.Background(), dbw, map[string]interface{}{}, true)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(config, test.expectedConfig) {
				t.Fatalf("Config mismatch: Actual: %#v\nExpected: %#v", config, test.expectedConfig)
			}
		})
	}
}

func TestInitDatabase_newDB(t *testing.T) {
	type testCase struct {
		initResp newdbplugin.InitializeResponse
		initErr  error

		expectedConfig map[string]interface{}
		expectErr      bool
	}

	tests := map[string]testCase{
		"new database error": {
			initErr: fmt.Errorf("test error"),

			expectedConfig: nil,
			expectErr:      true,
		},
		"legacy database success": {
			initResp: newdbplugin.InitializeResponse{
				Config: map[string]interface{}{
					"foo": "bar",
				},
			},

			expectedConfig: map[string]interface{}{
				"foo": "bar",
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			newDB := new(mockNewDatabase)
			newDB.On("Initialize", mock.Anything, mock.Anything).
				Return(test.initResp, test.initErr).
				Once()
			defer newDB.AssertExpectations(t)

			dbw := databaseVersionWrapper{
				database: newDB,
			}
			config, err := initDatabase(context.Background(), dbw, map[string]interface{}{}, true)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(config, test.expectedConfig) {
				t.Fatalf("Config mismatch: Actual: %#v\nExpected: %#v", config, test.expectedConfig)
			}
		})
	}
}

type fakePasswordGenerator struct {
	password string
	err      error
}

func (pg fakePasswordGenerator) GeneratePasswordFromPolicy(ctx context.Context, policy string) (string, error) {
	return pg.password, pg.err
}

func TestGeneratePassword(t *testing.T) {
	t.Run("no policy", func(t *testing.T) {
		pg := fakePasswordGenerator{
			password: "",
			err:      fmt.Errorf("the password generator shouldn't be called"),
		}

		password, err := generatePassword(context.Background(), pg, "")
		assertErrIsNil(t, err)
		// Technically this is checking the number of bytes, not the number of runes
		// But since the default should be ASCII characters, this is simplified
		if len(password) != defaultPasswordGenerator.Length {
			t.Fatalf("Password should be %d characters, but was %d", defaultPasswordGenerator.Length, len(password))
		}
	})

	t.Run("with policy", func(t *testing.T) {
		expected := "foobarbaz"
		pg := fakePasswordGenerator{
			password: expected,
			err:      nil,
		}

		actual, err := generatePassword(context.Background(), pg, "testpolicy")
		assertErrIsNil(t, err)
		if actual != expected {
			t.Fatalf("Actual password: %s\nExpected password: %s", actual, expected)
		}
	})
}

func TestCreateUser_legacy(t *testing.T) {
	statements := dbplugin.Statements{
		Creation: []string{
			"foo",
			"bar",
		},
	}
	displayName := "disp_name"
	roleName := "role_name"

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: displayName,
		RoleName:    roleName,
	}

	expiration := time.Now()

	expectedUser := "username"
	expectedPassword := "myreallysecurepassword"
	legacyDB := new(mockLegacyDatabase)
	legacyDB.On("CreateUser", mock.Anything, statements, usernameConfig, expiration).
		Return(expectedUser, expectedPassword, error(nil)).
		Once()
	defer legacyDB.AssertExpectations(t)

	dbw := databaseVersionWrapper{
		legacyDatabase: legacyDB,
	}

	pg := fakePasswordGenerator{
		password: "",
		err:      fmt.Errorf("this should not be called"),
	}

	actualUser, actualPass, err := createUser(context.Background(),
		dbw,
		pg,
		statements,
		displayName,
		roleName,
		expiration,
		"testpolicy")
	assertErrIsNil(t, err)
	if actualUser != expectedUser {
		t.Fatalf("Actual User: %q Expected: %q", actualUser, expectedUser)
	}
	if actualPass != expectedPassword {
		t.Fatalf("Actual Password: %q Expected: %q", actualPass, expectedPassword)
	}
}

func TestCreateUser_newDB(t *testing.T) {
	type testCase struct {
		reqPassword  string
		respUsername string

		newUserErr error

		expectedUsername string
		expectedPassword string
		expectErr        bool
	}

	tests := map[string]testCase{
		"errored": {
			reqPassword:      "mysecurepassword",
			respUsername:     "username",
			newUserErr:       fmt.Errorf("failed to create user because reasons"),
			expectedUsername: "",
			expectedPassword: "",
			expectErr:        true,
		},
		"happy path": {
			reqPassword:      "mysecurepassword",
			respUsername:     "username",
			newUserErr:       nil,
			expectedUsername: "username",
			expectedPassword: "mysecurepassword",
			expectErr:        false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			statements := dbplugin.Statements{
				Creation: []string{
					"foo",
					"bar",
				},
			}
			displayName := "disp_name"
			roleName := "role_name"
			expiration := time.Now()

			req := newdbplugin.NewUserRequest{
				UsernameConfig: newdbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Password:   test.reqPassword,
				Expiration: expiration,
				Statements: newdbplugin.Statements{
					Commands: statements.Creation,
				},
			}
			resp := newdbplugin.NewUserResponse{
				Username: test.respUsername,
			}

			newDB := new(mockNewDatabase)
			newDB.On("NewUser", mock.Anything, req).
				Return(resp, test.newUserErr).
				Once()
			defer newDB.AssertExpectations(t)

			dbw := databaseVersionWrapper{
				database: newDB,
			}

			pg := fakePasswordGenerator{
				password: test.reqPassword,
				err:      nil,
			}

			actualUser, actualPass, err := createUser(context.Background(),
				dbw,
				pg,
				statements,
				displayName,
				roleName,
				expiration,
				"testpolicy")
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if actualUser != test.expectedUsername {
				t.Fatalf("Actual User: %q Expected: %q", actualUser, test.expectedUsername)
			}
			if actualPass != test.expectedPassword {
				t.Fatalf("Actual Password: %q Expected: %q", actualPass, test.expectedPassword)
			}
		})
	}
}

func TestChangeUserPassword_legacy(t *testing.T) {
	type testCase struct {
		setCredsErr error
		expectErr   bool
	}

	tests := map[string]testCase{
		"errored update": {
			setCredsErr: fmt.Errorf("failed to update user because reasons"),
			expectErr:   true,
		},
		"happy path": {
			setCredsErr: nil,
			expectErr:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			statements := dbplugin.Statements{
				Rotation: []string{
					"foo",
					"bar",
				},
			}

			username := "username"
			expectedPassword := "myreallysecurepassword"

			userConfig := dbplugin.StaticUserConfig{
				Username: username,
				Password: expectedPassword,
			}
			legacyDB := new(mockLegacyDatabase)
			legacyDB.On("SetCredentials", mock.Anything, statements, userConfig).
				Return(username, expectedPassword, test.setCredsErr).
				Once()
			defer legacyDB.AssertExpectations(t)

			dbw := databaseVersionWrapper{
				legacyDatabase: legacyDB,
			}

			err := changeUserPassword(context.Background(),
				dbw,
				username,
				expectedPassword,
				statements.Rotation)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func TestChangeUserPassword_newDB(t *testing.T) {
	type testCase struct {
		updateErr error
		expectErr bool
	}

	tests := map[string]testCase{
		"errored update": {
			updateErr: fmt.Errorf("failed to update user because reasons"),
			expectErr: true,
		},
		"happy path": {
			updateErr: nil,
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			statements := []string{
				"foo",
				"bar",
			}

			username := "username"
			expectedPassword := "myreallysecurepassword"

			req := newdbplugin.UpdateUserRequest{
				Username: username,
				Password: &newdbplugin.ChangePassword{
					NewPassword: expectedPassword,
					Statements: newdbplugin.Statements{
						Commands: statements,
					},
				},
			}

			resp := newdbplugin.UpdateUserResponse{}

			newDB := new(mockNewDatabase)
			newDB.On("UpdateUser", mock.Anything, req).
				Return(resp, test.updateErr).
				Once()
			defer newDB.AssertExpectations(t)

			dbw := databaseVersionWrapper{
				database: newDB,
			}

			err := changeUserPassword(context.Background(),
				dbw,
				username,
				expectedPassword,
				statements)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func TestRenewUser_legacy(t *testing.T) {
	statements := dbplugin.Statements{
		Renewal: []string{
			"foo",
			"bar",
		},
	}
	username := "username"
	expiration := time.Now()

	legacyDB := new(mockLegacyDatabase)
	legacyDB.On("RenewUser", mock.Anything, statements, username, expiration).
		Return(error(nil)).
		Once()
	defer legacyDB.AssertExpectations(t)

	dbw := databaseVersionWrapper{
		legacyDatabase: legacyDB,
	}

	err := renewUser(context.Background(), dbw, username, expiration, statements.Renewal)
	assertErrIsNil(t, err)
}

func TestRenewUser_newDB(t *testing.T) {
	statements := []string{
		"foo",
		"bar",
	}
	username := "username"
	expiration := time.Now()

	req := newdbplugin.UpdateUserRequest{
		Username: username,
		Expiration: &newdbplugin.ChangeExpiration{
			NewExpiration: expiration,
			Statements: newdbplugin.Statements{
				Commands: statements,
			},
		},
	}

	resp := newdbplugin.UpdateUserResponse{}

	newDB := new(mockNewDatabase)
	newDB.On("UpdateUser", mock.Anything, req).
		Return(resp, error(nil)).
		Once()
	defer newDB.AssertExpectations(t)

	dbw := databaseVersionWrapper{
		database: newDB,
	}

	err := renewUser(context.Background(), dbw, username, expiration, statements)
	assertErrIsNil(t, err)
}

func TestDeleteUser_legacy(t *testing.T) {
	statements := dbplugin.Statements{
		Revocation: []string{
			"foo",
			"bar",
		},
	}
	username := "username"

	legacyDB := new(mockLegacyDatabase)
	legacyDB.On("RevokeUser", mock.Anything, statements, username).
		Return(error(nil)).
		Once()
	defer legacyDB.AssertExpectations(t)

	dbw := databaseVersionWrapper{
		legacyDatabase: legacyDB,
	}

	err := deleteUser(context.Background(), dbw, username, statements.Revocation)
	assertErrIsNil(t, err)
}

func TestDeleteUser_newDB(t *testing.T) {
	statements := []string{
		"foo",
		"bar",
	}
	username := "username"

	req := newdbplugin.DeleteUserRequest{
		Username: username,
		Statements: newdbplugin.Statements{
			Commands: statements,
		},
	}

	resp := newdbplugin.DeleteUserResponse{}

	newDB := new(mockNewDatabase)
	newDB.On("DeleteUser", mock.Anything, req).
		Return(resp, error(nil)).
		Once()
	defer newDB.AssertExpectations(t)

	dbw := databaseVersionWrapper{
		database: newDB,
	}

	err := deleteUser(context.Background(), dbw, username, statements)
	assertErrIsNil(t, err)
}

type badValue struct{}

func (badValue) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("this value cannot be marshalled to JSON")
}

var _ logical.Storage = fakeStorage{}

type fakeStorage struct {
	putErr error
}

func (f fakeStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	return f.putErr
}

func (f fakeStorage) List(ctx context.Context, s string) ([]string, error) {
	panic("list not implemented")
}
func (f fakeStorage) Get(ctx context.Context, s string) (*logical.StorageEntry, error) {
	panic("get not implemented")
}
func (f fakeStorage) Delete(ctx context.Context, s string) error {
	panic("delete not implemented")
}

func TestStoreConfig(t *testing.T) {
	type testCase struct {
		config    *DatabaseConfig
		putErr    error
		expectErr bool
	}

	tests := map[string]testCase{
		"bad config": {
			config: &DatabaseConfig{
				PluginName: "testplugin",
				ConnectionDetails: map[string]interface{}{
					"bad value": badValue{},
				},
			},
			putErr:    nil,
			expectErr: true,
		},
		"storage error": {
			config: &DatabaseConfig{
				PluginName: "testplugin",
				ConnectionDetails: map[string]interface{}{
					"foo": "bar",
				},
			},
			putErr:    fmt.Errorf("failed to store config"),
			expectErr: true,
		},
		"happy path": {
			config: &DatabaseConfig{
				PluginName: "testplugin",
				ConnectionDetails: map[string]interface{}{
					"foo": "bar",
				},
			},
			putErr:    nil,
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage := fakeStorage{
				putErr: test.putErr,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			err := storeConfig(ctx, storage, "testconfig", test.config)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func assertErrIsNil(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("No error expected, got: %s", err)
	}
}

func assertErrIsNotNil(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error, but didn't get one")
	}
}
