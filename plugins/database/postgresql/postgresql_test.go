// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/postgresql"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getPostgreSQL(t *testing.T, options map[string]interface{}) (*PostgreSQL, func()) {
	cleanup, connURL := postgresql.PrepareTestContainer(t, "13.4-buster")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}
	for k, v := range options {
		connectionDetails[k] = v
	}

	req := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	dbtesting.AssertInitialize(t, db, req)

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}
	return db, cleanup
}

func TestPostgreSQL_Initialize(t *testing.T) {
	db, cleanup := getPostgreSQL(t, map[string]interface{}{
		"max_open_connections": 5,
	})
	defer cleanup()

	if err := db.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostgreSQL_InitializeWithStringVals(t *testing.T) {
	db, cleanup := getPostgreSQL(t, map[string]interface{}{
		"max_open_connections": "5",
	})
	defer cleanup()

	if err := db.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostgreSQL_Initialize_ConnURLWithDSNFormat(t *testing.T) {
	cleanup, connURL := postgresql.PrepareTestContainer(t, "13.4-buster")
	defer cleanup()

	dsnConnURL, err := dbutil.ParseURL(connURL)
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": dsnConnURL,
	}

	req := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	dbtesting.AssertInitialize(t, db, req)

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}
}

// Ensures we can successfully initialize and connect to a CloudSQL database
// Requires the following:
// - GOOGLE_APPLICATION_CREDENTIALS either JSON or path to file
// - CONNECTION_URL to a valid Postgres instance on Google CloudSQL
func TestPostgreSQL_Initialize_CloudGCP(t *testing.T) {
	envConnURL := "CONNECTION_URL"
	connURL := os.Getenv(envConnURL)
	if connURL == "" {
		t.Skipf("env var %s not set, skipping test", envConnURL)
	}

	credStr := dbtesting.GetGCPTestCredentials(t)

	type testCase struct {
		req           dbplugin.InitializeRequest
		wantErr       bool
		expectedError string
	}

	tests := map[string]testCase{
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
		"default credentials": {
			req: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
					"auth_type":      connutil.AuthTypeGCPIAM,
				},
				VerifyConnection: true,
			},
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

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := new()
			defer dbtesting.AssertClose(t, db)

			_, err := dbtesting.VerifyInitialize(t, db, test.req)

			if test.wantErr {
				if err == nil {
					t.Fatalf("expected error but received nil")
				}

				if !strings.Contains(err.Error(), test.expectedError) {
					t.Fatalf("expected error %s, got %s", test.expectedError, err.Error())
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

// TestPostgreSQL_PasswordAuthentication tests that the default "password_authentication" is "none", and that
// an error is returned if an invalid "password_authentication" is provided.
func TestPostgreSQL_PasswordAuthentication(t *testing.T) {
	cleanup, connURL := postgresql.PrepareTestContainer(t, "13.4-buster")
	defer cleanup()

	dsnConnURL, err := dbutil.ParseURL(connURL)
	assert.NoError(t, err)
	db := new()

	ctx := context.Background()

	t.Run("invalid-password-authentication", func(t *testing.T) {
		connectionDetails := map[string]interface{}{
			"connection_url":          dsnConnURL,
			"password_authentication": "invalid-password-authentication",
		}

		req := dbplugin.InitializeRequest{
			Config:           connectionDetails,
			VerifyConnection: true,
		}

		_, err := db.Initialize(ctx, req)
		assert.EqualError(t, err, "'invalid-password-authentication' is not a valid password authentication type")
	})

	t.Run("default-is-none", func(t *testing.T) {
		connectionDetails := map[string]interface{}{
			"connection_url": dsnConnURL,
		}

		req := dbplugin.InitializeRequest{
			Config:           connectionDetails,
			VerifyConnection: true,
		}

		_ = dbtesting.AssertInitialize(t, db, req)
		assert.Equal(t, passwordAuthenticationPassword, db.passwordAuthentication)
	})
}

// TestPostgreSQL_PasswordAuthentication_SCRAMSHA256 tests that password_authentication works when set to scram-sha-256.
// When sending an encrypted password, the raw password should still successfully authenticate the user.
func TestPostgreSQL_PasswordAuthentication_SCRAMSHA256(t *testing.T) {
	cleanup, connURL := postgresql.PrepareTestContainer(t, "13.4-buster")
	defer cleanup()

	dsnConnURL, err := dbutil.ParseURL(connURL)
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url":          dsnConnURL,
		"password_authentication": string(passwordAuthenticationSCRAMSHA256),
	}

	req := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	resp := dbtesting.AssertInitialize(t, db, req)
	assert.Equal(t, string(passwordAuthenticationSCRAMSHA256), resp.Config["password_authentication"])

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	ctx := context.Background()
	newUserRequest := dbplugin.NewUserRequest{
		Statements: dbplugin.Statements{
			Commands: []string{
				`
						CREATE ROLE "{{name}}" WITH
						  LOGIN
						  PASSWORD '{{password}}'
						  VALID UNTIL '{{expiration}}';
						GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";`,
			},
		},
		Password:   "somesecurepassword",
		Expiration: time.Now().Add(1 * time.Minute),
	}
	newUserResponse, err := db.NewUser(ctx, newUserRequest)

	assertCredsExist(t, db.ConnectionURL, newUserResponse.Username, newUserRequest.Password)
}

func TestPostgreSQL_NewUser(t *testing.T) {
	type testCase struct {
		req            dbplugin.NewUserRequest
		expectErr      bool
		credsAssertion credsAssertion
	}

	tests := map[string]testCase{
		"no creation statements": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				// No statements
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: true,
			credsAssertion: assertCreds(
				assertUsernameRegex("^$"),
				assertCredsDoNotExist,
			),
		},
		"admin name": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`
						CREATE ROLE "{{name}}" WITH
						  LOGIN
						  PASSWORD '{{password}}'
						  VALID UNTIL '{{expiration}}';
						GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";`,
					},
				},
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: false,
			credsAssertion: assertCreds(
				assertUsernameRegex("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$"),
				assertCredsExist,
			),
		},
		"admin username": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`
						CREATE ROLE "{{username}}" WITH
						  LOGIN
						  PASSWORD '{{password}}'
						  VALID UNTIL '{{expiration}}';
						GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{username}}";`,
					},
				},
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: false,
			credsAssertion: assertCreds(
				assertUsernameRegex("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$"),
				assertCredsExist,
			),
		},
		"read only name": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`
						CREATE ROLE "{{name}}" WITH
						  LOGIN
						  PASSWORD '{{password}}'
						  VALID UNTIL '{{expiration}}';
						GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
						GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{name}}";`,
					},
				},
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: false,
			credsAssertion: assertCreds(
				assertUsernameRegex("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$"),
				assertCredsExist,
			),
		},
		"read only username": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`
						CREATE ROLE "{{username}}" WITH
						  LOGIN
						  PASSWORD '{{password}}'
						  VALID UNTIL '{{expiration}}';
						GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{username}}";
						GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{username}}";`,
					},
				},
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: false,
			credsAssertion: assertCreds(
				assertUsernameRegex("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$"),
				assertCredsExist,
			),
		},
		// https://github.com/hashicorp/vault/issues/6098
		"reproduce GH-6098": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						// NOTE: "rolname" in the following line is not a typo.
						"DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname='my_role') THEN CREATE ROLE my_role; END IF; END $$",
					},
				},
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: false,
			credsAssertion: assertCreds(
				assertUsernameRegex("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$"),
				assertCredsDoNotExist,
			),
		},
		"reproduce issue with template": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{
						`DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname='my_role') THEN CREATE ROLE "{{username}}"; END IF; END $$`,
					},
				},
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: false,
			credsAssertion: assertCreds(
				assertUsernameRegex("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$"),
				assertCredsDoNotExist,
			),
		},
		"large block statements": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: newUserLargeBlockStatements,
				},
				Password:   "somesecurepassword",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr: false,
			credsAssertion: assertCreds(
				assertUsernameRegex("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$"),
				assertCredsExist,
			),
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	db, cleanup := getPostgreSQL(t, nil)
	defer cleanup()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			resp, err := db.NewUser(ctx, test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			test.credsAssertion(t, db.ConnectionURL, resp.Username, test.req.Password)

			// Ensure that the role doesn't expire immediately
			time.Sleep(2 * time.Second)

			test.credsAssertion(t, db.ConnectionURL, resp.Username, test.req.Password)
		})
	}
}

func TestUpdateUser_Password(t *testing.T) {
	type testCase struct {
		statements     []string
		expectErr      bool
		credsAssertion credsAssertion
	}

	tests := map[string]testCase{
		"default statements": {
			statements:     nil,
			expectErr:      false,
			credsAssertion: assertCredsExist,
		},
		"explicit default statements": {
			statements:     []string{defaultChangePasswordStatement},
			expectErr:      false,
			credsAssertion: assertCredsExist,
		},
		"name instead of username": {
			statements:     []string{`ALTER ROLE "{{name}}" WITH PASSWORD '{{password}}';`},
			expectErr:      false,
			credsAssertion: assertCredsExist,
		},
		"bad statements": {
			statements:     []string{`asdofyas8uf77asoiajv`},
			expectErr:      true,
			credsAssertion: assertCredsDoNotExist,
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	db, cleanup := getPostgreSQL(t, nil)
	defer cleanup()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			initialPass := "myreallysecurepassword"
			createReq := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{createAdminUser},
				},
				Password:   initialPass,
				Expiration: time.Now().Add(2 * time.Second),
			}
			createResp := dbtesting.AssertNewUser(t, db, createReq)

			assertCredsExist(t, db.ConnectionURL, createResp.Username, initialPass)

			newPass := "somenewpassword"
			updateReq := dbplugin.UpdateUserRequest{
				Username: createResp.Username,
				Password: &dbplugin.ChangePassword{
					NewPassword: newPass,
					Statements: dbplugin.Statements{
						Commands: test.statements,
					},
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := db.UpdateUser(ctx, updateReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			test.credsAssertion(t, db.ConnectionURL, createResp.Username, newPass)
		})
	}

	t.Run("user does not exist", func(t *testing.T) {
		newPass := "somenewpassword"
		updateReq := dbplugin.UpdateUserRequest{
			Username: "missing-user",
			Password: &dbplugin.ChangePassword{
				NewPassword: newPass,
				Statements:  dbplugin.Statements{},
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := db.UpdateUser(ctx, updateReq)
		if err == nil {
			t.Fatalf("err expected, got nil")
		}

		assertCredsDoNotExist(t, db.ConnectionURL, updateReq.Username, newPass)
	})
}

func TestUpdateUser_Expiration(t *testing.T) {
	type testCase struct {
		initialExpiration  time.Time
		newExpiration      time.Time
		expectedExpiration time.Time
		statements         []string
		expectErr          bool
	}

	now := time.Now()
	tests := map[string]testCase{
		"no statements": {
			initialExpiration:  now.Add(1 * time.Minute),
			newExpiration:      now.Add(5 * time.Minute),
			expectedExpiration: now.Add(5 * time.Minute),
			statements:         nil,
			expectErr:          false,
		},
		"default statements with name": {
			initialExpiration:  now.Add(1 * time.Minute),
			newExpiration:      now.Add(5 * time.Minute),
			expectedExpiration: now.Add(5 * time.Minute),
			statements:         []string{defaultExpirationStatement},
			expectErr:          false,
		},
		"default statements with username": {
			initialExpiration:  now.Add(1 * time.Minute),
			newExpiration:      now.Add(5 * time.Minute),
			expectedExpiration: now.Add(5 * time.Minute),
			statements:         []string{`ALTER ROLE "{{username}}" VALID UNTIL '{{expiration}}';`},
			expectErr:          false,
		},
		"bad statements": {
			initialExpiration:  now.Add(1 * time.Minute),
			newExpiration:      now.Add(5 * time.Minute),
			expectedExpiration: now.Add(1 * time.Minute),
			statements:         []string{"ladshfouay09sgj"},
			expectErr:          true,
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	db, cleanup := getPostgreSQL(t, nil)
	defer cleanup()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			password := "myreallysecurepassword"
			initialExpiration := test.initialExpiration.Truncate(time.Second)
			createReq := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{createAdminUser},
				},
				Password:   password,
				Expiration: initialExpiration,
			}
			createResp := dbtesting.AssertNewUser(t, db, createReq)

			assertCredsExist(t, db.ConnectionURL, createResp.Username, password)

			actualExpiration := getExpiration(t, db, createResp.Username)
			if actualExpiration.IsZero() {
				t.Fatalf("Initial expiration is zero but should be set")
			}
			if !actualExpiration.Equal(initialExpiration) {
				t.Fatalf("Actual expiration: %s Expected expiration: %s", actualExpiration, initialExpiration)
			}

			newExpiration := test.newExpiration.Truncate(time.Second)
			updateReq := dbplugin.UpdateUserRequest{
				Username: createResp.Username,
				Expiration: &dbplugin.ChangeExpiration{
					NewExpiration: newExpiration,
					Statements: dbplugin.Statements{
						Commands: test.statements,
					},
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := db.UpdateUser(ctx, updateReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			expectedExpiration := test.expectedExpiration.Truncate(time.Second)
			actualExpiration = getExpiration(t, db, createResp.Username)
			if !actualExpiration.Equal(expectedExpiration) {
				t.Fatalf("Actual expiration: %s Expected expiration: %s", actualExpiration, expectedExpiration)
			}
		})
	}
}

func getExpiration(t testing.TB, db *PostgreSQL, username string) time.Time {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("select valuntil from pg_catalog.pg_user where usename = '%s'", username)
	conn, err := db.getConnection(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection to database: %s", err)
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		t.Fatalf("Failed to prepare statement: %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		t.Fatalf("Failed to execute query to get expiration: %s", err)
	}

	if !rows.Next() {
		return time.Time{} // No expiration
	}
	rawExp := ""
	err = rows.Scan(&rawExp)
	if err != nil {
		t.Fatalf("Unable to get raw expiration: %s", err)
	}
	if rawExp == "" {
		return time.Time{} // No expiration
	}
	exp, err := time.Parse(time.RFC3339, rawExp)
	if err != nil {
		t.Fatalf("Failed to parse expiration %q: %s", rawExp, err)
	}
	return exp
}

func TestDeleteUser(t *testing.T) {
	type testCase struct {
		revokeStmts    []string
		expectErr      bool
		credsAssertion credsAssertion
	}

	tests := map[string]testCase{
		"no statements": {
			revokeStmts: nil,
			expectErr:   false,
			// Wait for a short time before failing because postgres takes a moment to finish deleting the user
			credsAssertion: waitUntilCredsDoNotExist(2 * time.Second),
		},
		"statements with name": {
			revokeStmts: []string{`
				REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{name}}";
				REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{name}}";
				REVOKE USAGE ON SCHEMA public FROM "{{name}}";
		
				DROP ROLE IF EXISTS "{{name}}";`},
			expectErr: false,
			// Wait for a short time before failing because postgres takes a moment to finish deleting the user
			credsAssertion: waitUntilCredsDoNotExist(2 * time.Second),
		},
		"statements with username": {
			revokeStmts: []string{`
				REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{username}}";
				REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{username}}";
				REVOKE USAGE ON SCHEMA public FROM "{{username}}";
		
				DROP ROLE IF EXISTS "{{username}}";`},
			expectErr: false,
			// Wait for a short time before failing because postgres takes a moment to finish deleting the user
			credsAssertion: waitUntilCredsDoNotExist(2 * time.Second),
		},
		"bad statements": {
			revokeStmts: []string{`8a9yhfoiasjff`},
			expectErr:   true,
			// Wait for a short time before checking because postgres takes a moment to finish deleting the user
			credsAssertion: assertCredsExistAfter(100 * time.Millisecond),
		},
		"multiline": {
			revokeStmts: []string{`
				DO $$ BEGIN
					REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{username}}";
					REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{username}}";
					REVOKE USAGE ON SCHEMA public FROM "{{username}}";
					DROP ROLE IF EXISTS "{{username}}";
				END $$;
				`},
			expectErr: false,
			// Wait for a short time before checking because postgres takes a moment to finish deleting the user
			credsAssertion: waitUntilCredsDoNotExist(2 * time.Second),
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	db, cleanup := getPostgreSQL(t, nil)
	defer cleanup()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			password := "myreallysecurepassword"
			createReq := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{createAdminUser},
				},
				Password:   password,
				Expiration: time.Now().Add(2 * time.Second),
			}
			createResp := dbtesting.AssertNewUser(t, db, createReq)

			assertCredsExist(t, db.ConnectionURL, createResp.Username, password)

			deleteReq := dbplugin.DeleteUserRequest{
				Username: createResp.Username,
				Statements: dbplugin.Statements{
					Commands: test.revokeStmts,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := db.DeleteUser(ctx, deleteReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			test.credsAssertion(t, db.ConnectionURL, createResp.Username, password)
		})
	}
}

type credsAssertion func(t testing.TB, connURL, username, password string)

func assertCreds(assertions ...credsAssertion) credsAssertion {
	return func(t testing.TB, connURL, username, password string) {
		t.Helper()
		for _, assertion := range assertions {
			assertion(t, connURL, username, password)
		}
	}
}

func assertUsernameRegex(rawRegex string) credsAssertion {
	return func(t testing.TB, _, username, _ string) {
		t.Helper()
		require.Regexp(t, rawRegex, username)
	}
}

func assertCredsExist(t testing.TB, connURL, username, password string) {
	t.Helper()
	err := testCredsExist(t, connURL, username, password)
	if err != nil {
		t.Fatalf("user does not exist: %s", err)
	}
}

func assertCredsDoNotExist(t testing.TB, connURL, username, password string) {
	t.Helper()
	err := testCredsExist(t, connURL, username, password)
	if err == nil {
		t.Fatalf("user should not exist but does")
	}
}

func waitUntilCredsDoNotExist(timeout time.Duration) credsAssertion {
	return func(t testing.TB, connURL, username, password string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("Timed out waiting for user %s to be deleted", username)
			case <-ticker.C:
				err := testCredsExist(t, connURL, username, password)
				if err != nil {
					// Happy path
					return
				}
			}
		}
	}
}

func assertCredsExistAfter(timeout time.Duration) credsAssertion {
	return func(t testing.TB, connURL, username, password string) {
		t.Helper()
		time.Sleep(timeout)
		assertCredsExist(t, connURL, username, password)
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	t.Helper()
	// Log in with the new creds
	connURL = strings.Replace(connURL, "postgres:secret", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const createAdminUser = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

var newUserLargeBlockStatements = []string{
	`
DO $$
BEGIN
   IF NOT EXISTS (SELECT * FROM pg_catalog.pg_roles WHERE rolname='foo-role') THEN
      CREATE ROLE "foo-role";
      CREATE SCHEMA IF NOT EXISTS foo AUTHORIZATION "foo-role";
      ALTER ROLE "foo-role" SET search_path = foo;
      GRANT TEMPORARY ON DATABASE "postgres" TO "foo-role";
      GRANT ALL PRIVILEGES ON SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA foo TO "foo-role";
   END IF;
END
$$
`,
	`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';`,
	`GRANT "foo-role" TO "{{name}}";`,
	`ALTER ROLE "{{name}}" SET search_path = foo;`,
	`GRANT CONNECT ON DATABASE "postgres" TO "{{name}}";`,
}

func TestContainsMultilineStatement(t *testing.T) {
	type testCase struct {
		Input    string
		Expected bool
	}

	testCases := map[string]*testCase{
		"issue 6098 repro": {
			Input:    `DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname='my_role') THEN CREATE ROLE my_role; END IF; END $$`,
			Expected: true,
		},
		"multiline with template fields": {
			Input:    `DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname="{{name}}") THEN CREATE ROLE {{name}}; END IF; END $$`,
			Expected: true,
		},
		"docs example": {
			Input: `CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";`,
			Expected: false,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			if containsMultilineStatement(tCase.Input) != tCase.Expected {
				t.Fatalf("%q should be %t for multiline input", tCase.Input, tCase.Expected)
			}
		})
	}
}

func TestExtractQuotedStrings(t *testing.T) {
	type testCase struct {
		Input    string
		Expected []string
	}

	testCases := map[string]*testCase{
		"no quotes": {
			Input:    `Five little monkeys jumping on the bed`,
			Expected: []string{},
		},
		"two of both quote types": {
			Input:    `"Five" little 'monkeys' "jumping on" the' 'bed`,
			Expected: []string{`"Five"`, `"jumping on"`, `'monkeys'`, `' '`},
		},
		"one single quote": {
			Input:    `Five little monkeys 'jumping on the bed`,
			Expected: []string{},
		},
		"empty string": {
			Input:    ``,
			Expected: []string{},
		},
		"templated field": {
			Input:    `DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname="{{name}}") THEN CREATE ROLE {{name}}; END IF; END $$`,
			Expected: []string{`"{{name}}"`},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			results, err := extractQuotedStrings(tCase.Input)
			if err != nil {
				t.Fatal(err)
			}
			if len(results) != len(tCase.Expected) {
				t.Fatalf("%s isn't equal to %s", results, tCase.Expected)
			}
			for i := range results {
				if results[i] != tCase.Expected[i] {
					t.Fatalf(`expected %q but received %q`, tCase.Expected, results[i])
				}
			}
		})
	}
}

func TestUsernameGeneration(t *testing.T) {
	type testCase struct {
		data          dbplugin.UsernameMetadata
		expectedRegex string
	}

	tests := map[string]testCase{
		"simple display and role names": {
			data: dbplugin.UsernameMetadata{
				DisplayName: "token",
				RoleName:    "myrole",
			},
			expectedRegex: `v-token-myrole-[a-zA-Z0-9]{20}-[0-9]{10}`,
		},
		"display name has dash": {
			data: dbplugin.UsernameMetadata{
				DisplayName: "token-foo",
				RoleName:    "myrole",
			},
			expectedRegex: `v-token-fo-myrole-[a-zA-Z0-9]{20}-[0-9]{10}`,
		},
		"display name has underscore": {
			data: dbplugin.UsernameMetadata{
				DisplayName: "token_foo",
				RoleName:    "myrole",
			},
			expectedRegex: `v-token_fo-myrole-[a-zA-Z0-9]{20}-[0-9]{10}`,
		},
		"display name has period": {
			data: dbplugin.UsernameMetadata{
				DisplayName: "token.foo",
				RoleName:    "myrole",
			},
			expectedRegex: `v-token.fo-myrole-[a-zA-Z0-9]{20}-[0-9]{10}`,
		},
		"role name has dash": {
			data: dbplugin.UsernameMetadata{
				DisplayName: "token",
				RoleName:    "myrole-foo",
			},
			expectedRegex: `v-token-myrole-f-[a-zA-Z0-9]{20}-[0-9]{10}`,
		},
		"role name has underscore": {
			data: dbplugin.UsernameMetadata{
				DisplayName: "token",
				RoleName:    "myrole_foo",
			},
			expectedRegex: `v-token-myrole_f-[a-zA-Z0-9]{20}-[0-9]{10}`,
		},
		"role name has period": {
			data: dbplugin.UsernameMetadata{
				DisplayName: "token",
				RoleName:    "myrole.foo",
			},
			expectedRegex: `v-token-myrole.f-[a-zA-Z0-9]{20}-[0-9]{10}`,
		},
	}

	for name, test := range tests {
		t.Run(fmt.Sprintf("new-%s", name), func(t *testing.T) {
			up, err := template.NewTemplate(
				template.Template(defaultUserNameTemplate),
			)
			require.NoError(t, err)

			for i := 0; i < 1000; i++ {
				username, err := up.Generate(test.data)
				require.NoError(t, err)
				require.Regexp(t, test.expectedRegex, username)
			}
		})
	}
}

func TestNewUser_CustomUsername(t *testing.T) {
	cleanup, connURL := postgresql.PrepareTestContainer(t, "13.4-buster")
	defer cleanup()

	type testCase struct {
		usernameTemplate string
		newUserData      dbplugin.UsernameMetadata
		expectedRegex    string
	}

	tests := map[string]testCase{
		"default template": {
			usernameTemplate: "",
			newUserData: dbplugin.UsernameMetadata{
				DisplayName: "displayname",
				RoleName:    "longrolename",
			},
			expectedRegex: "^v-displayn-longrole-[a-zA-Z0-9]{20}-[0-9]{10}$",
		},
		"explicit default template": {
			usernameTemplate: defaultUserNameTemplate,
			newUserData: dbplugin.UsernameMetadata{
				DisplayName: "displayname",
				RoleName:    "longrolename",
			},
			expectedRegex: "^v-displayn-longrole-[a-zA-Z0-9]{20}-[0-9]{10}$",
		},
		"unique template": {
			usernameTemplate: "foo-bar",
			newUserData: dbplugin.UsernameMetadata{
				DisplayName: "displayname",
				RoleName:    "longrolename",
			},
			expectedRegex: "^foo-bar$",
		},
		"custom prefix": {
			usernameTemplate: "foobar-{{.DisplayName | truncate 8}}-{{.RoleName | truncate 8}}-{{random 20}}-{{unix_time}}",
			newUserData: dbplugin.UsernameMetadata{
				DisplayName: "displayname",
				RoleName:    "longrolename",
			},
			expectedRegex: "^foobar-displayn-longrole-[a-zA-Z0-9]{20}-[0-9]{10}$",
		},
		"totally custom template": {
			usernameTemplate: "foobar_{{random 10}}-{{.RoleName | uppercase}}.{{unix_time}}x{{.DisplayName | truncate 5}}",
			newUserData: dbplugin.UsernameMetadata{
				DisplayName: "displayname",
				RoleName:    "longrolename",
			},
			expectedRegex: `^foobar_[a-zA-Z0-9]{10}-LONGROLENAME\.[0-9]{10}xdispl$`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			initReq := dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":    connURL,
					"username_template": test.usernameTemplate,
				},
				VerifyConnection: true,
			}

			db := new()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_, err := db.Initialize(ctx, initReq)
			require.NoError(t, err)

			newUserReq := dbplugin.NewUserRequest{
				UsernameConfig: test.newUserData,
				Statements: dbplugin.Statements{
					Commands: []string{
						`
						CREATE ROLE "{{name}}" WITH
						  LOGIN
						  PASSWORD '{{password}}'
						  VALID UNTIL '{{expiration}}';
						GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";`,
					},
				},
				Password:   "myReally-S3curePassword",
				Expiration: time.Now().Add(1 * time.Hour),
			}
			ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			newUserResp, err := db.NewUser(ctx, newUserReq)
			require.NoError(t, err)

			require.Regexp(t, test.expectedRegex, newUserResp.Username)
		})
	}
}

func TestNewUser_CloudGCP(t *testing.T) {
	envConnURL := "CONNECTION_URL"
	connURL := os.Getenv(envConnURL)
	if connURL == "" {
		t.Skipf("env var %s not set, skipping test", envConnURL)
	}

	credStr := dbtesting.GetGCPTestCredentials(t)

	type testCase struct {
		usernameTemplate string
		newUserData      dbplugin.UsernameMetadata
		expectedRegex    string
	}

	tests := map[string]testCase{
		"default template": {
			usernameTemplate: "",
			newUserData: dbplugin.UsernameMetadata{
				DisplayName: "displayname",
				RoleName:    "longrolename",
			},
			expectedRegex: "^v-displayn-longrole-[a-zA-Z0-9]{20}-[0-9]{10}$",
		},
		"unique template": {
			usernameTemplate: "foo-bar",
			newUserData: dbplugin.UsernameMetadata{
				DisplayName: "displayname",
				RoleName:    "longrolename",
			},
			expectedRegex: "^foo-bar$",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			initReq := dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":       connURL,
					"username_template":    test.usernameTemplate,
					"auth_type":            connutil.AuthTypeGCPIAM,
					"service_account_json": credStr,
				},
				VerifyConnection: true,
			}

			db := new()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_, err := db.Initialize(ctx, initReq)
			require.NoError(t, err)

			newUserReq := dbplugin.NewUserRequest{
				UsernameConfig: test.newUserData,
				Statements: dbplugin.Statements{
					Commands: []string{
						`
						CREATE ROLE "{{name}}" WITH
						  LOGIN
						  PASSWORD '{{password}}'
						  VALID UNTIL '{{expiration}}';
						GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";`,
					},
				},
				Password:   "myReally-S3curePassword",
				Expiration: time.Now().Add(1 * time.Hour),
			}
			ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			newUserResp, err := db.NewUser(ctx, newUserReq)
			require.NoError(t, err)

			require.Regexp(t, test.expectedRegex, newUserResp.Username)
		})
	}
}

// This is a long-running integration test which tests the functionality of Postgres's multi-host
// connection strings. It uses two Postgres containers preconfigured with Replication Manager
// provided by Bitnami. This test currently does not run in CI and must be run manually. This is
// due to the test length, as it requires multiple sleep calls to ensure cluster setup and
// primary node failover occurs before the test steps continue.
//
// To run the test, set the environment variable POSTGRES_MULTIHOST_NET to the value of
// a docker network you've preconfigured, e.g.
// 'docker network create -d bridge postgres-repmgr'
// 'export POSTGRES_MULTIHOST_NET=postgres-repmgr'
func TestPostgreSQL_Repmgr(t *testing.T) {
	_, exists := os.LookupEnv("POSTGRES_MULTIHOST_NET")
	if !exists {
		t.Skipf("POSTGRES_MULTIHOST_NET not set, skipping test")
	}

	// Run two postgres-repmgr containers in a replication cluster
	db0, runner0, url0, container0 := testPostgreSQL_Repmgr_Container(t, "psql-repl-node-0")
	_, _, url1, _ := testPostgreSQL_Repmgr_Container(t, "psql-repl-node-1")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	time.Sleep(10 * time.Second)

	// Write a read role to the cluster
	_, err := db0.NewUser(ctx, dbplugin.NewUserRequest{
		Statements: dbplugin.Statements{
			Commands: []string{
				`CREATE ROLE "ro" NOINHERIT;
				GRANT SELECT ON ALL TABLES IN SCHEMA public TO "ro";`,
			},
		},
	})
	if err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}

	// Open a connection to both databases using the multihost connection string
	connectionDetails := map[string]interface{}{
		"connection_url": fmt.Sprintf("postgresql://{{username}}:{{password}}@%s,%s/postgres?target_session_attrs=read-write", getHost(url0), getHost(url1)),
		"username":       "postgres",
		"password":       "secret",
	}
	req := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	dbtesting.AssertInitialize(t, db, req)
	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}
	defer db.Close()

	// Add a user to the cluster, then stop the primary container
	if err = testPostgreSQL_Repmgr_AddUser(ctx, db); err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}
	postgresql.StopContainer(t, ctx, runner0, container0)

	// Try adding a new user immediately - expect failure as the database
	// cluster is still switching primaries
	err = testPostgreSQL_Repmgr_AddUser(ctx, db)
	if !strings.HasSuffix(err.Error(), "ValidateConnect failed (read only connection)") {
		t.Fatalf("expected error was not received, got: %s", err)
	}

	time.Sleep(20 * time.Second)

	// Try adding a new user again which should succeed after the sleep
	// as the primary failover should have finished. Then, restart
	// the first container which should become a secondary DB.
	if err = testPostgreSQL_Repmgr_AddUser(ctx, db); err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}
	postgresql.RestartContainer(t, ctx, runner0, container0)

	time.Sleep(10 * time.Second)

	// A final new user to add, which should succeed after the secondary joins.
	if err = testPostgreSQL_Repmgr_AddUser(ctx, db); err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testPostgreSQL_Repmgr_Container(t *testing.T, name string) (*PostgreSQL, *docker.Runner, string, string) {
	envVars := []string{
		"REPMGR_NODE_NAME=" + name,
		"REPMGR_NODE_NETWORK_NAME=" + name,
	}

	runner, cleanup, connURL, containerID := postgresql.PrepareTestContainerRepmgr(t, name, "13.4.0", envVars)
	t.Cleanup(cleanup)

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}
	req := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}
	db := new()
	dbtesting.AssertInitialize(t, db, req)
	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	if err := db.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}

	return db, runner, connURL, containerID
}

func testPostgreSQL_Repmgr_AddUser(ctx context.Context, db *PostgreSQL) error {
	_, err := db.NewUser(ctx, dbplugin.NewUserRequest{
		Statements: dbplugin.Statements{
			Commands: []string{
				`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}' INHERIT;
				GRANT ro TO "{{name}}";`,
			},
		},
	})

	return err
}

func getHost(url string) string {
	splitCreds := strings.Split(url, "@")[1]

	return strings.Split(splitCreds, "/")[0]
}
