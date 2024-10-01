// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	stdmysql "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	mysqlhelper "github.com/hashicorp/vault/helper/testhelpers/mysql"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/stretchr/testify/require"
)

var _ dbplugin.Database = (*MySQL)(nil)

func TestMySQL_Initialize(t *testing.T) {
	type testCase struct {
		rootPassword string
	}

	tests := map[string]testCase{
		"non-special characters in root password": {
			rootPassword: "B44a30c4C04D0aAaE140",
		},
		"special characters in root password": {
			rootPassword: "#secret!%25#{@}",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			testInitialize(t, test.rootPassword)
		})
	}
}

// TestMySQL_Initialize_CloudGCP validates the proper initialization of a MySQL backend pointing
// to a GCP CloudSQL MySQL instance. This expects some external setup (exact TBD)
func TestMySQL_Initialize_CloudGCP(t *testing.T) {
	envConnURL := "CONNECTION_URL"
	connURL := os.Getenv(envConnURL)
	if connURL == "" {
		t.Skipf("env var %s not set, skipping test", envConnURL)
	}

	credStr := dbtesting.GetGCPTestCredentials(t)

	tests := map[string]struct {
		req           dbplugin.InitializeRequest
		wantErr       bool
		expectedError string
	}{
		"empty auth type": {
			req: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
					"auth_type":      "",
				},
			},
		},
		"invalid auth type": {
			req: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
					"auth_type":      "invalid",
				},
			},
			wantErr:       true,
			expectedError: "invalid auth_type",
		},
		"JSON credentials": {
			req: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":       connURL,
					"auth_type":            connutil.AuthTypeGCPIAM,
					"service_account_json": credStr,
				},
				VerifyConnection: true,
			},
		},
	}

	for n, tc := range tests {
		t.Run(n, func(t *testing.T) {
			db := newMySQL(DefaultUserNameTemplate)
			defer dbtesting.AssertClose(t, db)
			_, err := db.Initialize(context.Background(), tc.req)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but received nil")
				}

				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("expected error %s, got %s", tc.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, received %s", err)
				}

				if !db.Initialized {
					t.Fatal("Database should be initialized")
				}
			}
		})
	}
}

func testInitialize(t *testing.T, rootPassword string) {
	cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, rootPassword)
	defer cleanup()

	mySQLConfig, err := stdmysql.ParseDSN(connURL)
	if err != nil {
		panic(fmt.Sprintf("Test failure: connection URL is invalid: %s", err))
	}
	rootUser := mySQLConfig.User
	mySQLConfig.User = "{{username}}"
	mySQLConfig.Passwd = "{{password}}"
	tmplConnURL := mySQLConfig.FormatDSN()

	type testCase struct {
		initRequest  dbplugin.InitializeRequest
		expectedResp dbplugin.InitializeResponse

		expectErr         bool
		expectInitialized bool
	}

	tests := map[string]testCase{
		"missing connection_url": {
			initRequest: dbplugin.InitializeRequest{
				Config:           map[string]interface{}{},
				VerifyConnection: true,
			},
			expectedResp:      dbplugin.InitializeResponse{},
			expectErr:         true,
			expectInitialized: false,
		},
		"basic config": {
			initRequest: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
				},
				VerifyConnection: true,
			},
			expectedResp: dbplugin.InitializeResponse{
				Config: map[string]interface{}{
					"connection_url": connURL,
				},
			},
			expectErr:         false,
			expectInitialized: true,
		},
		"username and password replacement in connection_url": {
			initRequest: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": tmplConnURL,
					"username":       rootUser,
					"password":       rootPassword,
				},
				VerifyConnection: true,
			},
			expectedResp: dbplugin.InitializeResponse{
				Config: map[string]interface{}{
					"connection_url": tmplConnURL,
					"username":       rootUser,
					"password":       rootPassword,
				},
			},
			expectErr:         false,
			expectInitialized: true,
		},
		"invalid username template": {
			initRequest: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":    connURL,
					"username_template": "{{.FieldThatDoesNotExist}}",
				},
				VerifyConnection: true,
			},
			expectedResp:      dbplugin.InitializeResponse{},
			expectErr:         true,
			expectInitialized: false,
		},
		"bad username template": {
			initRequest: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":    connURL,
					"username_template": "{{ .DisplayName", // Explicitly bad template
				},
				VerifyConnection: true,
			},
			expectedResp:      dbplugin.InitializeResponse{},
			expectErr:         true,
			expectInitialized: false,
		},
		"custom username template": {
			initRequest: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":    connURL,
					"username_template": "foo-{{random 10}}-{{.DisplayName}}",
				},
				VerifyConnection: true,
			},
			expectedResp: dbplugin.InitializeResponse{
				Config: map[string]interface{}{
					"connection_url":    connURL,
					"username_template": "foo-{{random 10}}-{{.DisplayName}}",
				},
			},
			expectErr:         false,
			expectInitialized: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := newMySQL(DefaultUserNameTemplate)
			defer dbtesting.AssertClose(t, db)
			initResp, err := db.Initialize(context.Background(), test.initRequest)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			require.Equal(t, test.expectedResp, initResp)
			require.Equal(t, test.expectInitialized, db.Initialized, "Initialized variable not set correctly")
		})
	}
}

func TestMySQL_NewUser_nonLegacy(t *testing.T) {
	displayName := "token"
	roleName := "testrole"

	type testCase struct {
		usernameTemplate string

		newUserReq dbplugin.NewUserRequest

		expectedUsernameRegex string
		expectErr             bool
	}

	tests := map[string]testCase{
		"name statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
						GRANT SELECT ON *.* TO '{{name}}'@'%';`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-token-testrole-[a-zA-Z0-9]{15}$`,
			expectErr:             false,
		},
		"username statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
						GRANT SELECT ON *.* TO '{{username}}'@'%';`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-token-testrole-[a-zA-Z0-9]{15}$`,
			expectErr:             false,
		},
		"prepared name statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
						set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{name}}'@'%'");
						PREPARE grantStmt from @grants;
						EXECUTE grantStmt;
						DEALLOCATE PREPARE grantStmt;`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-token-testrole-[a-zA-Z0-9]{15}$`,
			expectErr:             false,
		},
		"prepared username statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
						set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{username}}'@'%'");
						PREPARE grantStmt from @grants;
						EXECUTE grantStmt;
						DEALLOCATE PREPARE grantStmt;`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-token-testrole-[a-zA-Z0-9]{15}$`,
			expectErr:             false,
		},
		"custom username template": {
			usernameTemplate: "foo-{{random 10}}-{{.RoleName | uppercase}}",

			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
						GRANT SELECT ON *.* TO '{{username}}'@'%';`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^foo-[a-zA-Z0-9]{10}-TESTROLE$`,
			expectErr:             false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
			defer cleanup()

			connectionDetails := map[string]interface{}{
				"connection_url":    connURL,
				"username_template": test.usernameTemplate,
			}

			initReq := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := newMySQL(DefaultUserNameTemplate)
			defer db.Close()
			_, err := db.Initialize(context.Background(), initReq)
			require.NoError(t, err)

			userResp, err := db.NewUser(context.Background(), test.newUserReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			require.Regexp(t, test.expectedUsernameRegex, userResp.Username)

			err = mysqlhelper.TestCredsExist(t, connURL, userResp.Username, test.newUserReq.Password)
			require.NoError(t, err, "Failed to connect with credentials")
		})
	}
}

func TestMySQL_NewUser_legacy(t *testing.T) {
	displayName := "token"
	roleName := "testrole"

	type testCase struct {
		usernameTemplate string

		newUserReq dbplugin.NewUserRequest

		expectedUsernameRegex string
		expectErr             bool
	}

	tests := map[string]testCase{
		"name statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
						GRANT SELECT ON *.* TO '{{name}}'@'%';`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-test-[a-zA-Z0-9]{9}$`,
			expectErr:             false,
		},
		"username statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
						GRANT SELECT ON *.* TO '{{username}}'@'%';`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-test-[a-zA-Z0-9]{9}$`,
			expectErr:             false,
		},
		"prepared name statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
						set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{name}}'@'%'");
						PREPARE grantStmt from @grants;
						EXECUTE grantStmt;
						DEALLOCATE PREPARE grantStmt;`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-test-[a-zA-Z0-9]{9}$`,
			expectErr:             false,
		},
		"prepared username statements": {
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
						set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{username}}'@'%'");
						PREPARE grantStmt from @grants;
						EXECUTE grantStmt;
						DEALLOCATE PREPARE grantStmt;`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^v-test-[a-zA-Z0-9]{9}$`,
			expectErr:             false,
		},
		"custom username template": {
			usernameTemplate: `{{printf "foo-%s-%s" (random 5) (.RoleName | uppercase) | truncate 16}}`,

			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: displayName,
					RoleName:    roleName,
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
						GRANT SELECT ON *.* TO '{{username}}'@'%';`,
					},
				},
				Password:   "09g8hanbdfkVSM",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: `^foo-[a-zA-Z0-9]{5}-TESTRO$`,
			expectErr:             false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
			defer cleanup()

			connectionDetails := map[string]interface{}{
				"connection_url":    connURL,
				"username_template": test.usernameTemplate,
			}

			initReq := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := newMySQL(DefaultLegacyUserNameTemplate)
			defer db.Close()
			_, err := db.Initialize(context.Background(), initReq)
			require.NoError(t, err)

			userResp, err := db.NewUser(context.Background(), test.newUserReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			require.Regexp(t, test.expectedUsernameRegex, userResp.Username)

			err = mysqlhelper.TestCredsExist(t, connURL, userResp.Username, test.newUserReq.Password)
			require.NoError(t, err, "Failed to connect with credentials")
		})
	}
}

func TestMySQL_RotateRootCredentials(t *testing.T) {
	type testCase struct {
		statements []string
	}

	tests := map[string]testCase{
		"empty statements": {
			statements: nil,
		},
		"default username": {
			statements: []string{defaultMySQLRotateCredentialsSQL},
		},
		"default name": {
			statements: []string{
				`
				ALTER USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
			defer cleanup()

			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
				"username":       "root",
				"password":       "secret",
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			initReq := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := newMySQL(DefaultUserNameTemplate)
			defer db.Close()
			_, err := db.Initialize(context.Background(), initReq)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !db.Initialized {
				t.Fatal("Database should be initialized")
			}

			updateReq := dbplugin.UpdateUserRequest{
				Username: "root",
				Password: &dbplugin.ChangePassword{
					NewPassword: "different_sercret",
					Statements: dbplugin.Statements{
						Commands: test.statements,
					},
				},
			}

			_, err = db.UpdateUser(ctx, updateReq)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			err = mysqlhelper.TestCredsExist(t, connURL, updateReq.Username, updateReq.Password.NewPassword)
			if err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			// verify old password doesn't work
			if err := mysqlhelper.TestCredsExist(t, connURL, updateReq.Username, "secret"); err == nil {
				t.Fatalf("Should not be able to connect with initial credentials")
			}

			err = db.Close()
			if err != nil {
				t.Fatalf("err: %s", err)
			}
		})
	}
}

func TestMySQL_DeleteUser(t *testing.T) {
	type testCase struct {
		revokeStmts []string
	}

	tests := map[string]testCase{
		"empty statements": {
			revokeStmts: nil,
		},
		"default name": {
			revokeStmts: []string{defaultMysqlRevocationStmts},
		},
		"default username": {
			revokeStmts: []string{
				`
				REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{username}}'@'%'; 
				DROP USER '{{username}}'@'%'`,
			},
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	initReq := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := newMySQL(DefaultUserNameTemplate)
	defer db.Close()
	_, err := db.Initialize(context.Background(), initReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			password, err := credsutil.RandomAlphaNumeric(32, false)
			if err != nil {
				t.Fatalf("unable to generate password: %s", err)
			}

			createReq := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`
						CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
						GRANT SELECT ON *.* TO '{{name}}'@'%';`,
					},
				},
				Password:   password,
				Expiration: time.Now().Add(time.Minute),
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			userResp, err := db.NewUser(ctx, createReq)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, userResp.Username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			deleteReq := dbplugin.DeleteUserRequest{
				Username: userResp.Username,
				Statements: dbplugin.Statements{
					Commands: test.revokeStmts,
				},
			}
			_, err = db.DeleteUser(context.Background(), deleteReq)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, userResp.Username, password); err == nil {
				t.Fatalf("Credentials were not revoked!")
			}
		})
	}
}

func TestMySQL_UpdateUser(t *testing.T) {
	type testCase struct {
		rotateStmts []string
	}

	tests := map[string]testCase{
		"empty statements": {
			rotateStmts: nil,
		},
		"custom statement name": {
			rotateStmts: []string{`
				ALTER USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';`},
		},
		"custom statement username": {
			rotateStmts: []string{`
				ALTER USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';`},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
			defer cleanup()

			// create the database user and verify we can access
			dbUser := "vaultstatictest"
			initPassword := "password"

			createStatements := `
				CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
				GRANT SELECT ON *.* TO '{{name}}'@'%';`

			createTestMySQLUser(t, connURL, dbUser, initPassword, createStatements)
			if err := mysqlhelper.TestCredsExist(t, connURL, dbUser, initPassword); err != nil {
				t.Fatalf("Could not connect with credentials: %s", err)
			}

			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
			}

			initReq := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			db := newMySQL(DefaultUserNameTemplate)
			defer db.Close()
			_, err := db.Initialize(context.Background(), initReq)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			newPassword, err := credsutil.RandomAlphaNumeric(32, false)
			if err != nil {
				t.Fatalf("unable to generate password: %s", err)
			}

			updateReq := dbplugin.UpdateUserRequest{
				Username: dbUser,
				Password: &dbplugin.ChangePassword{
					NewPassword: newPassword,
					Statements: dbplugin.Statements{
						Commands: test.rotateStmts,
					},
				},
			}

			_, err = db.UpdateUser(ctx, updateReq)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			// verify new password works
			if err := mysqlhelper.TestCredsExist(t, connURL, dbUser, newPassword); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			// verify old password doesn't work
			if err := mysqlhelper.TestCredsExist(t, connURL, dbUser, initPassword); err == nil {
				t.Fatalf("Should not be able to connect with initial credentials")
			}
		})
	}
}

func createTestMySQLUser(t *testing.T, connURL, username, password, query string) {
	t.Helper()
	db, err := sql.Open("mysql", connURL)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Start a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// copied from mysql.go
	for _, query := range strutil.ParseArbitraryStringSlice(query, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}
		query = dbutil.QueryHelper(query, map[string]string{
			"name":     username,
			"password": password,
		})

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			if e, ok := err.(*stdmysql.MySQLError); ok && e.Number == 1295 {
				_, err = tx.ExecContext(ctx, query)
				if err != nil {
					t.Fatal(err)
				}
				stmt.Close()
				continue
			}

			t.Fatal(err)
		}
		if _, err := stmt.ExecContext(ctx); err != nil {
			stmt.Close()
			t.Fatal(err)
		}
		stmt.Close()
	}
}
