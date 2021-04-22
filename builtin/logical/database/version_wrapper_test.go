package database

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInitDatabase_missingDB(t *testing.T) {
	dbw := databaseVersionWrapper{}

	req := v5.InitializeRequest{}
	resp, err := dbw.Initialize(context.Background(), req)
	if err == nil {
		t.Fatalf("err expected, got nil")
	}

	expectedResp := v5.InitializeResponse{}
	if !reflect.DeepEqual(resp, expectedResp) {
		t.Fatalf("Actual resp: %#v\nExpected resp: %#v", resp, expectedResp)
	}
}

func TestInitDatabase_newDB(t *testing.T) {
	type testCase struct {
		req v5.InitializeRequest

		newInitResp  v5.InitializeResponse
		newInitErr   error
		newInitCalls int

		expectedResp v5.InitializeResponse
		expectErr    bool
	}

	tests := map[string]testCase{
		"success": {
			req: v5.InitializeRequest{
				Config: map[string]interface{}{
					"foo": "bar",
				},
				VerifyConnection: true,
			},
			newInitResp: v5.InitializeResponse{
				Config: map[string]interface{}{
					"foo": "bar",
				},
			},
			newInitCalls: 1,
			expectedResp: v5.InitializeResponse{
				Config: map[string]interface{}{
					"foo": "bar",
				},
			},
			expectErr: false,
		},
		"error": {
			req: v5.InitializeRequest{
				Config: map[string]interface{}{
					"foo": "bar",
				},
				VerifyConnection: true,
			},
			newInitResp:  v5.InitializeResponse{},
			newInitErr:   fmt.Errorf("test error"),
			newInitCalls: 1,
			expectedResp: v5.InitializeResponse{},
			expectErr:    true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			newDB := new(mockNewDatabase)
			newDB.On("Initialize", mock.Anything, mock.Anything).
				Return(test.newInitResp, test.newInitErr)
			defer newDB.AssertNumberOfCalls(t, "Initialize", test.newInitCalls)

			dbw := databaseVersionWrapper{
				v5: newDB,
			}

			resp, err := dbw.Initialize(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual resp: %#v\nExpected resp: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestInitDatabase_legacyDB(t *testing.T) {
	type testCase struct {
		req v5.InitializeRequest

		initConfig map[string]interface{}
		initErr    error
		initCalls  int

		expectedResp v5.InitializeResponse
		expectErr    bool
	}

	tests := map[string]testCase{
		"success": {
			req: v5.InitializeRequest{
				Config: map[string]interface{}{
					"foo": "bar",
				},
				VerifyConnection: true,
			},
			initConfig: map[string]interface{}{
				"foo": "bar",
			},
			initCalls: 1,
			expectedResp: v5.InitializeResponse{
				Config: map[string]interface{}{
					"foo": "bar",
				},
			},
			expectErr: false,
		},
		"error": {
			req: v5.InitializeRequest{
				Config: map[string]interface{}{
					"foo": "bar",
				},
				VerifyConnection: true,
			},
			initErr:      fmt.Errorf("test error"),
			initCalls:    1,
			expectedResp: v5.InitializeResponse{},
			expectErr:    true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			legacyDB := new(mockLegacyDatabase)
			legacyDB.On("Init", mock.Anything, mock.Anything, mock.Anything).
				Return(test.initConfig, test.initErr)
			defer legacyDB.AssertNumberOfCalls(t, "Init", test.initCalls)

			dbw := databaseVersionWrapper{
				v4: legacyDB,
			}

			resp, err := dbw.Initialize(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual resp: %#v\nExpected resp: %#v", resp, test.expectedResp)
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

func TestGeneratePassword_missingDB(t *testing.T) {
	dbw := databaseVersionWrapper{}

	gen := fakePasswordGenerator{
		err: fmt.Errorf("this shouldn't be called"),
	}
	pass, err := dbw.GeneratePassword(context.Background(), gen, "policy")
	if err == nil {
		t.Fatalf("err expected, got nil")
	}

	if pass != "" {
		t.Fatalf("Password should be empty but was: %s", pass)
	}
}

func TestGeneratePassword_legacy(t *testing.T) {
	type testCase struct {
		legacyPassword string
		legacyErr      error
		legacyCalls    int

		expectedPassword string
		expectErr        bool
	}

	tests := map[string]testCase{
		"legacy password generation": {
			legacyPassword: "legacy_password",
			legacyErr:      nil,
			legacyCalls:    1,

			expectedPassword: "legacy_password",
			expectErr:        false,
		},
		"legacy password failure": {
			legacyPassword: "",
			legacyErr:      fmt.Errorf("failed :("),
			legacyCalls:    1,

			expectedPassword: "",
			expectErr:        true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			legacyDB := new(mockLegacyDatabase)
			legacyDB.On("GenerateCredentials", mock.Anything).
				Return(test.legacyPassword, test.legacyErr)
			defer legacyDB.AssertNumberOfCalls(t, "GenerateCredentials", test.legacyCalls)

			dbw := databaseVersionWrapper{
				v4: legacyDB,
			}

			passGen := fakePasswordGenerator{
				err: fmt.Errorf("this should not be called"),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			password, err := dbw.GeneratePassword(ctx, passGen, "test_policy")
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			if password != test.expectedPassword {
				t.Fatalf("Actual password: %s Expected password: %s", password, test.expectedPassword)
			}
		})
	}
}

func TestGeneratePassword_policies(t *testing.T) {
	type testCase struct {
		passwordPolicyPassword string
		passwordPolicyErr      error

		expectedPassword string
		expectErr        bool
	}

	tests := map[string]testCase{
		"password policy generation": {
			passwordPolicyPassword: "new_password",

			expectedPassword: "new_password",
			expectErr:        false,
		},
		"password policy error": {
			passwordPolicyPassword: "",
			passwordPolicyErr:      fmt.Errorf("test error"),

			expectedPassword: "",
			expectErr:        true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			newDB := new(mockNewDatabase)
			defer newDB.AssertExpectations(t)

			dbw := databaseVersionWrapper{
				v5: newDB,
			}

			passGen := fakePasswordGenerator{
				password: test.passwordPolicyPassword,
				err:      test.passwordPolicyErr,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			password, err := dbw.GeneratePassword(ctx, passGen, "test_policy")
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			if password != test.expectedPassword {
				t.Fatalf("Actual password: %s Expected password: %s", password, test.expectedPassword)
			}
		})
	}
}

func TestGeneratePassword_no_policy(t *testing.T) {
	newDB := new(mockNewDatabase)
	defer newDB.AssertExpectations(t)

	dbw := databaseVersionWrapper{
		v5: newDB,
	}

	passGen := fakePasswordGenerator{
		password: "",
		err:      fmt.Errorf("should not be called"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	password, err := dbw.GeneratePassword(ctx, passGen, "")
	if err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}
	if password == "" {
		t.Fatalf("missing password")
	}

	rawRegex := "^[a-zA-Z0-9-]{20}$"
	re := regexp.MustCompile(rawRegex)
	if !re.MatchString(password) {
		t.Fatalf("password %q did not match regex: %q", password, rawRegex)
	}
}

func TestNewUser_missingDB(t *testing.T) {
	dbw := databaseVersionWrapper{}

	req := v5.NewUserRequest{}
	resp, pass, err := dbw.NewUser(context.Background(), req)
	if err == nil {
		t.Fatalf("err expected, got nil")
	}

	expectedResp := v5.NewUserResponse{}
	if !reflect.DeepEqual(resp, expectedResp) {
		t.Fatalf("Actual resp: %#v\nExpected resp: %#v", resp, expectedResp)
	}

	if pass != "" {
		t.Fatalf("Password should be empty but was: %s", pass)
	}
}

func TestNewUser_newDB(t *testing.T) {
	type testCase struct {
		req v5.NewUserRequest

		newUserResp  v5.NewUserResponse
		newUserErr   error
		newUserCalls int

		expectedResp v5.NewUserResponse
		expectErr    bool
	}

	tests := map[string]testCase{
		"success": {
			req: v5.NewUserRequest{
				Password: "new_password",
			},

			newUserResp: v5.NewUserResponse{
				Username: "newuser",
			},
			newUserCalls: 1,

			expectedResp: v5.NewUserResponse{
				Username: "newuser",
			},
			expectErr: false,
		},
		"error": {
			req: v5.NewUserRequest{
				Password: "new_password",
			},

			newUserErr:   fmt.Errorf("test error"),
			newUserCalls: 1,

			expectedResp: v5.NewUserResponse{},
			expectErr:    true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			newDB := new(mockNewDatabase)
			newDB.On("NewUser", mock.Anything, mock.Anything).
				Return(test.newUserResp, test.newUserErr)
			defer newDB.AssertNumberOfCalls(t, "NewUser", test.newUserCalls)

			dbw := databaseVersionWrapper{
				v5: newDB,
			}

			resp, password, err := dbw.NewUser(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual resp: %#v\nExpected resp: %#v", resp, test.expectedResp)
			}

			if password != test.req.Password {
				t.Fatalf("Actual password: %s Expected password: %s", password, test.req.Password)
			}
		})
	}
}

func TestNewUser_legacyDB(t *testing.T) {
	type testCase struct {
		req v5.NewUserRequest

		createUserUsername string
		createUserPassword string
		createUserErr      error
		createUserCalls    int

		expectedResp     v5.NewUserResponse
		expectedPassword string
		expectErr        bool
	}

	tests := map[string]testCase{
		"success": {
			req: v5.NewUserRequest{
				Password: "new_password",
			},

			createUserUsername: "newuser",
			createUserPassword: "securepassword",
			createUserCalls:    1,

			expectedResp: v5.NewUserResponse{
				Username: "newuser",
			},
			expectedPassword: "securepassword",
			expectErr:        false,
		},
		"error": {
			req: v5.NewUserRequest{
				Password: "new_password",
			},

			createUserErr:   fmt.Errorf("test error"),
			createUserCalls: 1,

			expectedResp: v5.NewUserResponse{},
			expectErr:    true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			legacyDB := new(mockLegacyDatabase)
			legacyDB.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(test.createUserUsername, test.createUserPassword, test.createUserErr)
			defer legacyDB.AssertNumberOfCalls(t, "CreateUser", test.createUserCalls)

			dbw := databaseVersionWrapper{
				v4: legacyDB,
			}

			resp, password, err := dbw.NewUser(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual resp: %#v\nExpected resp: %#v", resp, test.expectedResp)
			}

			if password != test.expectedPassword {
				t.Fatalf("Actual password: %s Expected password: %s", password, test.req.Password)
			}
		})
	}
}

func TestUpdateUser_missingDB(t *testing.T) {
	dbw := databaseVersionWrapper{}

	req := v5.UpdateUserRequest{}
	resp, err := dbw.UpdateUser(context.Background(), req, false)
	if err == nil {
		t.Fatalf("err expected, got nil")
	}

	expectedConfig := map[string]interface{}(nil)
	if !reflect.DeepEqual(resp, expectedConfig) {
		t.Fatalf("Actual config: %#v\nExpected config: %#v", resp, expectedConfig)
	}
}

func TestUpdateUser_newDB(t *testing.T) {
	type testCase struct {
		req v5.UpdateUserRequest

		updateUserErr   error
		updateUserCalls int

		expectedResp v5.UpdateUserResponse
		expectErr    bool
	}

	tests := map[string]testCase{
		"success": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
			},
			updateUserCalls: 1,
			expectErr:       false,
		},
		"error": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
			},
			updateUserErr:   fmt.Errorf("test error"),
			updateUserCalls: 1,
			expectErr:       true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			newDB := new(mockNewDatabase)
			newDB.On("UpdateUser", mock.Anything, mock.Anything).
				Return(v5.UpdateUserResponse{}, test.updateUserErr)
			defer newDB.AssertNumberOfCalls(t, "UpdateUser", test.updateUserCalls)

			dbw := databaseVersionWrapper{
				v5: newDB,
			}

			_, err := dbw.UpdateUser(context.Background(), test.req, false)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func TestUpdateUser_legacyDB(t *testing.T) {
	type testCase struct {
		req        v5.UpdateUserRequest
		isRootUser bool

		setCredentialsErr   error
		setCredentialsCalls int

		rotateRootConfig map[string]interface{}
		rotateRootErr    error
		rotateRootCalls  int

		renewUserErr   error
		renewUserCalls int

		expectedConfig map[string]interface{}
		expectErr      bool
	}

	tests := map[string]testCase{
		"missing changes": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
			},
			isRootUser: false,

			setCredentialsCalls: 0,
			rotateRootCalls:     0,
			renewUserCalls:      0,

			expectErr: true,
		},
		"both password and expiration changes": {
			req: v5.UpdateUserRequest{
				Username:   "existing_user",
				Password:   &v5.ChangePassword{},
				Expiration: &v5.ChangeExpiration{},
			},
			isRootUser: false,

			setCredentialsCalls: 0,
			rotateRootCalls:     0,
			renewUserCalls:      0,

			expectErr: true,
		},
		"change password - SetCredentials": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
				Password: &v5.ChangePassword{
					NewPassword: "newpassowrd",
				},
			},
			isRootUser: false,

			setCredentialsErr:   nil,
			setCredentialsCalls: 1,
			rotateRootCalls:     0,
			renewUserCalls:      0,

			expectedConfig: nil,
			expectErr:      false,
		},
		"change password - SetCredentials failed": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
				Password: &v5.ChangePassword{
					NewPassword: "newpassowrd",
				},
			},
			isRootUser: false,

			setCredentialsErr:   fmt.Errorf("set credentials failed"),
			setCredentialsCalls: 1,
			rotateRootCalls:     0,
			renewUserCalls:      0,

			expectedConfig: nil,
			expectErr:      true,
		},
		"change password - SetCredentials unimplemented but not a root user": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
				Password: &v5.ChangePassword{
					NewPassword: "newpassowrd",
				},
			},
			isRootUser: false,

			setCredentialsErr:   status.Error(codes.Unimplemented, "SetCredentials is not implemented"),
			setCredentialsCalls: 1,

			rotateRootCalls: 0,
			renewUserCalls:  0,

			expectedConfig: nil,
			expectErr:      true,
		},
		"change password - RotateRootCredentials": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
				Password: &v5.ChangePassword{
					NewPassword: "newpassowrd",
				},
			},
			isRootUser: true,

			setCredentialsErr:   status.Error(codes.Unimplemented, "SetCredentials is not implemented"),
			setCredentialsCalls: 1,

			rotateRootConfig: map[string]interface{}{
				"foo": "bar",
			},
			rotateRootCalls: 1,

			renewUserCalls: 0,

			expectedConfig: map[string]interface{}{
				"foo": "bar",
			},
			expectErr: false,
		},
		"change password - RotateRootCredentials failed": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
				Password: &v5.ChangePassword{
					NewPassword: "newpassowrd",
				},
			},
			isRootUser: true,

			setCredentialsErr:   status.Error(codes.Unimplemented, "SetCredentials is not implemented"),
			setCredentialsCalls: 1,

			rotateRootErr:   fmt.Errorf("rotate root failed"),
			rotateRootCalls: 1,
			renewUserCalls:  0,

			expectedConfig: nil,
			expectErr:      true,
		},

		"change expiration": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
				Expiration: &v5.ChangeExpiration{
					NewExpiration: time.Now(),
				},
			},
			isRootUser: false,

			setCredentialsCalls: 0,
			rotateRootCalls:     0,

			renewUserErr:   nil,
			renewUserCalls: 1,

			expectedConfig: nil,
			expectErr:      false,
		},
		"change expiration failed": {
			req: v5.UpdateUserRequest{
				Username: "existing_user",
				Expiration: &v5.ChangeExpiration{
					NewExpiration: time.Now(),
				},
			},
			isRootUser: false,

			setCredentialsCalls: 0,
			rotateRootCalls:     0,

			renewUserErr:   fmt.Errorf("test error"),
			renewUserCalls: 1,

			expectedConfig: nil,
			expectErr:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			legacyDB := new(mockLegacyDatabase)
			legacyDB.On("SetCredentials", mock.Anything, mock.Anything, mock.Anything).
				Return("", "", test.setCredentialsErr)
			defer legacyDB.AssertNumberOfCalls(t, "SetCredentials", test.setCredentialsCalls)

			legacyDB.On("RotateRootCredentials", mock.Anything, mock.Anything).
				Return(test.rotateRootConfig, test.rotateRootErr)
			defer legacyDB.AssertNumberOfCalls(t, "RotateRootCredentials", test.rotateRootCalls)

			legacyDB.On("RenewUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(test.renewUserErr)
			defer legacyDB.AssertNumberOfCalls(t, "RenewUser", test.renewUserCalls)

			dbw := databaseVersionWrapper{
				v4: legacyDB,
			}

			newConfig, err := dbw.UpdateUser(context.Background(), test.req, test.isRootUser)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(newConfig, test.expectedConfig) {
				t.Fatalf("Actual config: %#v\nExpected config: %#v", newConfig, test.expectedConfig)
			}
		})
	}
}

func TestDeleteUser_missingDB(t *testing.T) {
	dbw := databaseVersionWrapper{}

	req := v5.DeleteUserRequest{}
	_, err := dbw.DeleteUser(context.Background(), req)
	if err == nil {
		t.Fatalf("err expected, got nil")
	}
}

func TestDeleteUser_newDB(t *testing.T) {
	type testCase struct {
		req v5.DeleteUserRequest

		deleteUserErr   error
		deleteUserCalls int

		expectErr bool
	}

	tests := map[string]testCase{
		"success": {
			req: v5.DeleteUserRequest{
				Username: "existing_user",
			},

			deleteUserErr:   nil,
			deleteUserCalls: 1,

			expectErr: false,
		},
		"error": {
			req: v5.DeleteUserRequest{
				Username: "existing_user",
			},

			deleteUserErr:   fmt.Errorf("test error"),
			deleteUserCalls: 1,

			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			newDB := new(mockNewDatabase)
			newDB.On("DeleteUser", mock.Anything, mock.Anything).
				Return(v5.DeleteUserResponse{}, test.deleteUserErr)
			defer newDB.AssertNumberOfCalls(t, "DeleteUser", test.deleteUserCalls)

			dbw := databaseVersionWrapper{
				v5: newDB,
			}

			_, err := dbw.DeleteUser(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func TestDeleteUser_legacyDB(t *testing.T) {
	type testCase struct {
		req v5.DeleteUserRequest

		revokeUserErr   error
		revokeUserCalls int

		expectErr bool
	}

	tests := map[string]testCase{
		"success": {
			req: v5.DeleteUserRequest{
				Username: "existing_user",
			},

			revokeUserErr:   nil,
			revokeUserCalls: 1,

			expectErr: false,
		},
		"error": {
			req: v5.DeleteUserRequest{
				Username: "existing_user",
			},

			revokeUserErr:   fmt.Errorf("test error"),
			revokeUserCalls: 1,

			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			legacyDB := new(mockLegacyDatabase)
			legacyDB.On("RevokeUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(test.revokeUserErr)
			defer legacyDB.AssertNumberOfCalls(t, "RevokeUser", test.revokeUserCalls)

			dbw := databaseVersionWrapper{
				v4: legacyDB,
			}

			_, err := dbw.DeleteUser(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
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
