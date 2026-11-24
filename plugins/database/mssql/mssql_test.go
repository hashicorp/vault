// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	mssqldb "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/go-uuid"
	mssqlhelper "github.com/hashicorp/vault/helper/testhelpers/mssql"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMSSQLInitialize(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()

	type testCase struct {
		req dbplugin.InitializeRequest
	}

	tests := map[string]testCase{
		"happy path": {
			req: dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
				},
				VerifyConnection: true,
			},
		},
		"max_open_connections set": {
			dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":       connURL,
					"max_open_connections": "5",
				},
				VerifyConnection: true,
			},
		},
		"contained_db set": {
			dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
					"contained_db":   true,
				},
				VerifyConnection: true,
			},
		},
		"contained_db set string": {
			dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
					"contained_db":   "true",
				},
				VerifyConnection: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := new()
			dbtesting.AssertInitializeCircleCiTest(t, db, test.req)
			defer dbtesting.AssertClose(t, db)

			if !db.Initialized {
				t.Fatal("Database should be initialized")
			}
		})
	}
}

func TestMSSQLNewUser(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()

	type testCase struct {
		usernameTemplate string
		req              dbplugin.NewUserRequest
		usernameRegex    string
		expectErr        bool
		assertUser       func(t testing.TB, connURL, username, password string)
	}

	tests := map[string]testCase{
		"no creation statements": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{},
				Password:   "AG4qagho-dsvZ",
				Expiration: time.Now().Add(1 * time.Second),
			},
			usernameRegex: "^$",
			expectErr:     true,
			assertUser:    assertCredsDoNotExist,
		},
		"with creation statements": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test",
					RoleName:    "test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{testMSSQLRole},
				},
				Password:   "AG4qagho-dsvZ",
				Expiration: time.Now().Add(1 * time.Second),
			},
			usernameRegex: "^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$",
			expectErr:     false,
			assertUser:    assertCredsExist,
		},
		"custom username template": {
			usernameTemplate: "{{random 10}}_{{.RoleName}}.{{.DisplayName | sha256}}",
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "tokenwithlotsofextracharactershere",
					RoleName:    "myrolenamewithlotsofextracharacters",
				},
				Statements: dbplugin.Statements{
					Commands: []string{testMSSQLRole},
				},
				Password:   "AG4qagho-dsvZ",
				Expiration: time.Now().Add(1 * time.Second),
			},
			usernameRegex: "^[a-zA-Z0-9]{10}_myrolenamewithlotsofextracharacters.80d15d22dba29ddbd4994f8009b5ff7b17922c267eb49fb805a9488bd55d11f9$",
			expectErr:     false,
			assertUser:    assertCredsExist,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			usernameRe, err := regexp.Compile(test.usernameRegex)
			if err != nil {
				t.Fatalf("failed to compile username regex %q: %s", test.usernameRegex, err)
			}

			initReq := dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":    connURL,
					"username_template": test.usernameTemplate,
				},
				VerifyConnection: true,
			}

			db := new()
			dbtesting.AssertInitializeCircleCiTest(t, db, initReq)
			defer dbtesting.AssertClose(t, db)

			createResp, err := db.NewUser(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !usernameRe.MatchString(createResp.Username) {
				t.Fatalf("Generated username %q did not match regex %q", createResp.Username, test.usernameRegex)
			}

			// Protect against future fields that aren't specified
			expectedResp := dbplugin.NewUserResponse{
				Username: createResp.Username,
			}
			if !reflect.DeepEqual(createResp, expectedResp) {
				t.Fatalf("Fields missing from expected response: Actual: %#v", createResp)
			}

			test.assertUser(t, connURL, createResp.Username, test.req.Password)
		})
	}
}

func TestMSSQLUpdateUser_password(t *testing.T) {
	type testCase struct {
		req              dbplugin.UpdateUserRequest
		expectErr        bool
		expectedPassword string
	}

	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()
	initPassword := "p4$sw0rd"

	tests := map[string]testCase{
		"missing password": {
			req: dbplugin.UpdateUserRequest{
				Password: &dbplugin.ChangePassword{
					NewPassword: "",
					Statements:  dbplugin.Statements{},
				},
			},
			expectErr:        true,
			expectedPassword: initPassword,
		},
		"empty rotation statements": {
			req: dbplugin.UpdateUserRequest{
				Password: &dbplugin.ChangePassword{
					NewPassword: "N90gkKLy8$angf",
					Statements:  dbplugin.Statements{},
				},
			},
			expectErr:        false,
			expectedPassword: "N90gkKLy8$angf",
		},
		"username rotation": {
			req: dbplugin.UpdateUserRequest{
				Password: &dbplugin.ChangePassword{
					NewPassword: "N90gkKLy8$angf",
					Statements: dbplugin.Statements{
						Commands: []string{
							"ALTER LOGIN [{{username}}] WITH PASSWORD = '{{password}}'",
						},
					},
				},
			},
			expectErr:        false,
			expectedPassword: "N90gkKLy8$angf",
		},
		"bad statements": {
			req: dbplugin.UpdateUserRequest{
				Password: &dbplugin.ChangePassword{
					NewPassword: "N90gkKLy8$angf",
					Statements: dbplugin.Statements{
						Commands: []string{
							"ahosh98asjdffs",
						},
					},
				},
			},
			expectErr:        true,
			expectedPassword: initPassword,
		},
	}

	i := 0
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			initReq := dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url": connURL,
				},
				VerifyConnection: true,
			}

			db := new()
			dbtesting.AssertInitializeCircleCiTest(t, db, initReq)
			defer dbtesting.AssertClose(t, db)

			dbUser := fmt.Sprintf("vaultuser%d", i)
			test.req.Username = dbUser
			i++
			err := createTestMSSQLUser(connURL, dbUser, initPassword, testMSSQLLogin)
			if err != nil {
				t.Fatalf("Failed to create user: %s", err)
			}

			assertCredsExist(t, connURL, dbUser, initPassword)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			updateResp, err := db.UpdateUser(ctx, test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			// Protect against future fields that aren't specified
			expectedResp := dbplugin.UpdateUserResponse{}
			if !reflect.DeepEqual(updateResp, expectedResp) {
				t.Fatalf("Fields missing from expected response: Actual: %#v", updateResp)
			}

			assertCredsExist(t, connURL, dbUser, test.expectedPassword)

			// Delete user at the end of each test
			deleteReq := dbplugin.DeleteUserRequest{
				Username: dbUser,
			}

			ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			deleteResp, err := db.DeleteUser(ctx, deleteReq)
			if err != nil {
				t.Fatalf("Failed to delete user: %s", err)
			}

			// Protect against future fields that aren't specified
			expectedDeleteResp := dbplugin.DeleteUserResponse{}
			if !reflect.DeepEqual(deleteResp, expectedDeleteResp) {
				t.Fatalf("Fields missing from expected response: Actual: %#v", deleteResp)
			}

			assertCredsDoNotExist(t, connURL, dbUser, initPassword)
		})
	}
}

func TestMSSQLDeleteUser(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()

	dbUser := "vaultuser"
	initPassword := "p4$sw0rd"

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}

	db := new()

	dbtesting.AssertInitializeCircleCiTest(t, db, initReq)
	defer dbtesting.AssertClose(t, db)

	err := createTestMSSQLUser(connURL, dbUser, initPassword, testMSSQLLogin)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	assertCredsExist(t, connURL, dbUser, initPassword)

	deleteReq := dbplugin.DeleteUserRequest{
		Username: dbUser,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	deleteResp, err := db.DeleteUser(ctx, deleteReq)
	if err != nil {
		t.Fatalf("Failed to delete user: %s", err)
	}

	// Protect against future fields that aren't specified
	expectedResp := dbplugin.DeleteUserResponse{}
	if !reflect.DeepEqual(deleteResp, expectedResp) {
		t.Fatalf("Fields missing from expected response: Actual: %#v", deleteResp)
	}

	assertCredsDoNotExist(t, connURL, dbUser, initPassword)
}

func TestMSSQLDeleteUserContainedDB(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()

	dbUser := "vaultuser"
	initPassword := "p4$sw0rd"

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
			"contained_db":   true,
		},
		VerifyConnection: true,
	}

	db := new()

	dbtesting.AssertInitializeCircleCiTest(t, db, initReq)
	defer dbtesting.AssertClose(t, db)

	err := createTestMSSQLUser(connURL, dbUser, initPassword, testMSSQLContainedLogin)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	assertCredsExist(t, connURL, dbUser, initPassword)

	deleteReq := dbplugin.DeleteUserRequest{
		Username: dbUser,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	deleteResp, err := db.DeleteUser(ctx, deleteReq)
	if err != nil {
		t.Fatalf("Failed to delete user: %s", err)
	}

	// Protect against future fields that aren't specified
	expectedResp := dbplugin.DeleteUserResponse{}
	if !reflect.DeepEqual(deleteResp, expectedResp) {
		t.Fatalf("Fields missing from expected response: Actual: %#v", deleteResp)
	}

	assertContainedDBCredsDoNotExist(t, connURL, dbUser)
}

// TestMSSQLDeleteUserWithBatch verifies that USE database context switches work correctly when statements are executed as a single batch.
func TestMSSQLDeleteUserWithBatch(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()

	dbUser := "vaultuser"
	initPassword := "p4$sw0rd"

	// Create a test database
	setupDB, err := sql.Open("mssql", connURL)
	if err != nil {
		t.Fatalf("Failed to open connection for setup: %s", err)
	}
	defer setupDB.Close()

	_, err = setupDB.ExecContext(context.Background(), "CREATE DATABASE testdb")
	if err != nil {
		t.Fatalf("Failed to create test database: %s", err)
	}

	// Create login in master and user in testdb
	createSQL := `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
USE [testdb];
CREATE USER [{{name}}] FOR LOGIN [{{name}}];`

	err = createTestMSSQLUser(connURL, dbUser, initPassword, createSQL)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	assertCredsExist(t, connURL, dbUser, initPassword)

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}

	db := new()
	dbtesting.AssertInitializeCircleCiTest(t, db, initReq)
	defer dbtesting.AssertClose(t, db)

	deleteReq := dbplugin.DeleteUserRequest{
		Username: dbUser,
		Statements: dbplugin.Statements{
			Commands: []string{
				`USE [testdb]; DROP USER [{{name}}]; USE [master]; DROP LOGIN [{{name}}];`,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	deleteResp, err := db.DeleteUser(ctx, deleteReq)
	if err != nil {
		t.Fatalf("Failed to delete user with batch statement: %s", err)
	}

	expectedResp := dbplugin.DeleteUserResponse{}
	if !reflect.DeepEqual(deleteResp, expectedResp) {
		t.Fatalf("Fields missing from expected response: Actual: %#v", deleteResp)
	}

	assertCredsDoNotExist(t, connURL, dbUser, initPassword)
}

func TestMSSQLContainedDBSQLSanitization(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()

	injectionString := "vaultuser]"
	dbUser := "vaultuser"
	initPassword := "p4$sw0rd"

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}

	db := new()

	dbtesting.AssertInitializeCircleCiTest(t, db, initReq)
	defer dbtesting.AssertClose(t, db)

	err := createTestMSSQLUser(connURL, dbUser, initPassword, testMSSQLContainedLogin)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	assertCredsExist(t, connURL, dbUser, initPassword)

	deleteReq := dbplugin.DeleteUserRequest{
		Username: injectionString,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.DeleteUser(ctx, deleteReq)

	assert.EqualError(t, err, "mssql: Cannot alter the login 'vaultuser]', because it does not exist or you do not have permission.")
}

func TestMSSQLSanitization(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	defer cleanup()

	injectionString := "vaultuser]"
	dbUser := "vaultuser"
	initPassword := "p4$sw0rd"

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}

	db := new()

	dbtesting.AssertInitializeCircleCiTest(t, db, initReq)
	defer dbtesting.AssertClose(t, db)

	err := createTestMSSQLUser(connURL, dbUser, initPassword, testMSSQLLogin)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	assertCredsExist(t, connURL, dbUser, initPassword)

	deleteReq := dbplugin.DeleteUserRequest{
		Username: injectionString,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.DeleteUser(ctx, deleteReq)

	assert.EqualError(t, err, "mssql: Cannot alter the login 'vaultuser]', because it does not exist or you do not have permission.")
}

// TestMSSQLDropUserSQL verifies that database and user names remain literal
// identifiers, including Unicode names on a non-Unicode code page, and that
// name matching follows the target database's case and accent sensitivity.
func TestMSSQLDropUserSQL(t *testing.T) {
	t.Parallel()
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	t.Cleanup(cleanup)

	setupDB, err := sql.Open("sqlserver", connURL)
	require.NoError(t, err, "open setup connection")
	t.Cleanup(func() { require.NoError(t, setupDB.Close(), "close setup connection") })

	tests := map[string]struct {
		dbSuffix  string
		collation string
		username  string
		lookup    string
		neighbor  string
	}{
		"plain":               {username: "vaultuser"},
		"closing bracket":     {dbSuffix: "]", username: "vault]user", neighbor: "vault]]user"},
		"repeated brackets":   {dbSuffix: "]]", username: "vault]]user", neighbor: "vault]user"},
		"bracket-shaped name": {username: "[vaultuser]"},
		"apostrophe":          {dbSuffix: "'", username: "vault'user"},
		"combined delimiters": {dbSuffix: "]'", username: "vault]u'ser"},
		"SQL-looking name":    {dbSuffix: "]; SELECT 1--", username: "vault]; SELECT 1--"},
		"double quote":        {username: `vault"user`},
		"space":               {username: "vault user"},
		"leading space":       {username: " vaultuser"},
		"trailing space":      {username: "vaultuser "},
		"literal dot":         {username: "vault.user"},
		"reserved word":       {username: "SELECT"},
		"connection syntax":   {username: "vault@user:/%"},
		"Unicode":             {dbSuffix: "\u6570\u636e", username: "\u7528\u6237"},
		"Unicode bracket":     {username: "\u7528\u6237]"},
		"Unicode apostrophe":  {username: "\u7528\u6237'"},
		"minimum username":    {username: "v"},
		"127 characters":      {username: strings.Repeat("]", 127)},
		"maximum username":    {username: strings.Repeat("]", 128)},
		"maximum Unicode":     {username: strings.Repeat("\u7528", 128)},
		"case sensitive":      {collation: "Latin1_General_100_CS_AS", username: "VaultUser", neighbor: "vaultuser"},
		"accent sensitive":    {collation: "Latin1_General_100_CS_AS", username: "vault\u00e9", neighbor: "vaulte"},
		"case insensitive":    {collation: "Latin1_General_100_CI_AS", username: "VaultUser", lookup: "vaultuser"},
		"accent insensitive":  {collation: "Latin1_General_100_CI_AI", username: "vault\u00e9", lookup: "vaulte"},
	}
	// Subtests share one SQL Server and run sequentially: concurrent login and
	// user DDL can deadlock (error 1205) and fail setup nondeterministically.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			neighbor := tc.neighbor
			if neighbor == "" {
				neighbor = "vault_sentinel"
			}
			collation := tc.collation
			if collation == "" {
				collation = "Latin1_General_100_CI_AS"
			}
			lookup := tc.lookup
			if lookup == "" {
				lookup = tc.username
			}
			id, err := uuid.GenerateUUID()
			require.NoError(t, err, "generate isolated database name")
			dbName := "vault_quote_" + strings.ReplaceAll(id, "-", "") + tc.dbSuffix
			quotedDB := "[" + strings.ReplaceAll(dbName, "]", "]]") + "]"
			_, err = setupDB.Exec("CREATE DATABASE " + quotedDB + " COLLATE " + collation)
			require.NoError(t, err, "create isolated database")
			t.Cleanup(func() {
				_, err := setupDB.Exec("ALTER DATABASE " + quotedDB + " SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE " + quotedDB)
				require.NoError(t, err, "drop isolated database")
			})

			targetURL, err := url.Parse(connURL)
			require.True(t, err == nil, "parse connection URL (error not printed: it would include the password)")
			params := targetURL.Query()
			params.Set("database", dbName)
			targetURL.RawQuery = params.Encode()
			targetDB, err := sql.Open("sqlserver", targetURL.String())
			require.NoError(t, err, "open target database")
			t.Cleanup(func() { require.NoError(t, targetDB.Close(), "close target database") })

			_, err = targetDB.Exec(`
DECLARE @stmt nvarchar(max) = N'CREATE USER ' + QUOTENAME(@username) + N' WITHOUT LOGIN';
EXEC (@stmt);
SET @stmt = N'CREATE USER ' + QUOTENAME(@neighbor) + N' WITHOUT LOGIN';
EXEC (@stmt);`, sql.Named("username", tc.username), sql.Named("neighbor", neighbor))
			require.NoError(t, err, "create target and unrelated principal")

			var principalID, neighborID int
			var storedName string
			err = targetDB.QueryRow("SELECT principal_id, name FROM sys.database_principals WHERE name = @p1", tc.username).Scan(&principalID, &storedName)
			require.NoError(t, err, "query target principal before deletion")
			require.Equal(t, tc.username, storedName, "SQL Server must preserve the exact identifier")
			err = targetDB.QueryRow("SELECT principal_id FROM sys.database_principals WHERE name = @p1", neighbor).Scan(&neighborID)
			require.NoError(t, err, "query unrelated principal before deletion")
			require.NotEqual(t, principalID, neighborID, "fixture must contain distinct principals")

			_, err = setupDB.Exec(dropUserSQL, sql.Named("dbname", dbName), sql.Named("username", lookup))
			require.NoError(t, err, "revoke the exact target principal")
			var count int
			err = targetDB.QueryRow("SELECT COUNT(*) FROM sys.database_principals WHERE principal_id = @p1", principalID).Scan(&count)
			require.NoError(t, err, "query target principal after deletion")
			require.Zero(t, count, "target principal must actually be deleted")
			err = targetDB.QueryRow("SELECT COUNT(*) FROM sys.database_principals WHERE principal_id = @p1", neighborID).Scan(&count)
			require.NoError(t, err, "query unrelated principal")
			require.Equal(t, 1, count, "unrelated principal must remain")

			_, err = setupDB.Exec(dropUserSQL, sql.Named("dbname", dbName), sql.Named("username", lookup))
			require.NoError(t, err, "revoking an already absent database user must succeed")
		})
	}
}

// TestMSSQLDeleteUserQuotedUsername verifies that the default revocation flow
// removes both the database user and server login, not merely disables login.
func TestMSSQLDeleteUserQuotedUsername(t *testing.T) {
	t.Parallel()
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	t.Cleanup(cleanup)

	tests := map[string]string{
		"bracket":            "]",
		"repeated brackets":  "]]",
		"bracket-shaped":     "[user]",
		"apostrophe":         "'",
		"combined":           "]'",
		"SQL-looking":        "]; SELECT 1--",
		"Unicode":            "\u7528\u6237",
		"Unicode bracket":    "\u7528\u6237]",
		"Unicode apostrophe": "\u7528\u6237'",
		"space":              " reader",
		"dot":                ".reader",
		"double quote":       `"reader`,
		"connection syntax":  "@reader:/%",
		"127 characters":     strings.Repeat("]", 89),
		"128 characters":     strings.Repeat("]", 90),
	}
	// Subtests share one SQL Server and run sequentially: concurrent login and
	// user DDL can deadlock (error 1205) and fail setup nondeterministically.
	for name, suffix := range tests {
		t.Run(name, func(t *testing.T) {
			id, err := uuid.GenerateUUID()
			require.NoError(t, err, "generate unique principal")
			sentinel := "vault_" + strings.ReplaceAll(id, "-", "")
			username := sentinel + suffix
			quotedUser := "[" + strings.ReplaceAll(username, "]", "]]") + "]"

			setupDB, err := sql.Open("sqlserver", connURL)
			require.NoError(t, err, "open setup connection")
			t.Cleanup(func() { require.NoError(t, setupDB.Close(), "close setup connection") })
			t.Cleanup(func() {
				for _, principal := range []string{username, sentinel} {
					quoted := "[" + strings.ReplaceAll(principal, "]", "]]") + "]"
					_, err := setupDB.Exec("DROP USER IF EXISTS "+quoted+"; IF EXISTS (SELECT 1 FROM sys.server_principals WHERE name = @p1) DROP LOGIN "+quoted, principal)
					assert.NoError(t, err, "clean up test principals")
				}
			})
			_, err = setupDB.Exec("CREATE LOGIN " + quotedUser + " WITH PASSWORD = 'p4$sw0rd'; CREATE USER " + quotedUser + " FOR LOGIN " + quotedUser)
			require.NoError(t, err, "create target login and database user")
			_, err = setupDB.Exec("CREATE LOGIN [" + sentinel + "] WITH PASSWORD = 'p4$sw0rd'; CREATE USER [" + sentinel + "] FOR LOGIN [" + sentinel + "]")
			require.NoError(t, err, "create similar-named sentinel")

			var sentinelID int
			var storedUser, storedLogin string
			err = setupDB.QueryRow("SELECT name FROM sys.database_principals WHERE name = @p1", username).Scan(&storedUser)
			require.NoError(t, err, "find target user before revocation")
			require.Equal(t, username, storedUser, "preserve exact database username")
			err = setupDB.QueryRow("SELECT name FROM sys.server_principals WHERE name = @p1", username).Scan(&storedLogin)
			require.NoError(t, err, "find target login before revocation")
			require.Equal(t, username, storedLogin, "preserve exact login name")
			err = setupDB.QueryRow("SELECT principal_id FROM sys.server_principals WHERE name = @p1", sentinel).Scan(&sentinelID)
			require.NoError(t, err, "find sentinel login before revocation")

			db := new()
			dbtesting.AssertInitializeCircleCiTest(t, db, dbplugin.InitializeRequest{
				Config:           map[string]interface{}{"connection_url": connURL},
				VerifyConnection: true,
			})
			t.Cleanup(func() { dbtesting.AssertClose(t, db) })
			_, err = db.DeleteUser(context.Background(), dbplugin.DeleteUserRequest{Username: username})
			require.NoError(t, err, "revoke the exact principal")

			var count int
			err = setupDB.QueryRow("SELECT COUNT(*) FROM sys.database_principals WHERE name = @p1", username).Scan(&count)
			require.NoError(t, err, "query database user")
			require.Zero(t, count, "database user must be deleted")
			err = setupDB.QueryRow("SELECT COUNT(*) FROM sys.server_principals WHERE name = @p1", username).Scan(&count)
			require.NoError(t, err, "query server login")
			require.Zero(t, count, "server login must be deleted")
			err = setupDB.QueryRow("SELECT COUNT(*) FROM sys.database_principals WHERE name = @p1", sentinel).Scan(&count)
			require.NoError(t, err, "query sentinel user")
			require.Equal(t, 1, count, "similar-named user must remain")
			err = setupDB.QueryRow("SELECT COUNT(*) FROM sys.server_principals WHERE principal_id = @p1 AND is_disabled = 0", sentinelID).Scan(&count)
			require.NoError(t, err, "query sentinel login")
			require.Equal(t, 1, count, "same sentinel login must remain enabled")

			sentinelURL, err := url.Parse(connURL)
			require.True(t, err == nil, "parse sentinel connection URL (error not printed: it would include the password)")
			sentinelURL.User = url.UserPassword(sentinel, "p4$sw0rd")
			sentinelDB, err := sql.Open("sqlserver", sentinelURL.String())
			require.NoError(t, err, "open sentinel connection")
			t.Cleanup(func() { require.NoError(t, sentinelDB.Close(), "close sentinel connection") })
			require.NoError(t, sentinelDB.PingContext(context.Background()), "sentinel must still authenticate")
		})
	}
}

// TestMSSQLUpdateUserDefaultQuoted verifies that default password rotation,
// used when no rotation statements are configured, treats the username and
// password as literal values for both server logins and contained database
// users, and leaves a similarly named principal's password unchanged.
func TestMSSQLUpdateUserDefaultQuoted(t *testing.T) {
	t.Parallel()
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	t.Cleanup(cleanup)

	setupDB, err := sql.Open("sqlserver", connURL)
	require.NoError(t, err, "open setup connection")
	t.Cleanup(func() { require.NoError(t, setupDB.Close(), "close setup connection") })
	// Contained database users need this server option. It only allows
	// contained databases to be created; existing databases are unaffected.
	_, err = setupDB.Exec("EXEC sp_configure 'contained database authentication', 1; RECONFIGURE;")
	require.NoError(t, err, "enable contained database authentication")

	const initPassword = "Az9-Initial-Password"
	tests := map[string]struct {
		contained bool
		suffix    string
		password  string
	}{
		"login bracket":           {suffix: "]", password: "Az9]Rotated]]"},
		"login apostrophe":        {suffix: "'", password: "Az9'Rotated''"},
		"login SQL-looking":       {suffix: "]; SELECT 1--", password: "Az9'; SELECT 1--"},
		"login Unicode":           {suffix: "\u7528\u6237", password: "Az9-\u5bc6\u7801'xyz"},
		"login 128 characters":    {suffix: strings.Repeat("]", 90), password: "Az9]Rotated"},
		"login 128-char password": {suffix: "_reader", password: "Az9'" + strings.Repeat("x", 124)},
		"contained bracket":       {contained: true, suffix: "]", password: "Az9]Rotated]]"},
		"contained SQL-looking":   {contained: true, suffix: "]; SELECT 1--", password: "Az9'; SELECT 1--"},
		"contained Unicode":       {contained: true, suffix: "\u7528\u6237", password: "Az9-\u5bc6\u7801'xyz"},
	}
	// Subtests share one SQL Server and run sequentially: concurrent login and
	// user DDL can deadlock (error 1205) and fail setup nondeterministically.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			id, err := uuid.GenerateUUID()
			require.NoError(t, err, "generate unique principal")
			id = strings.ReplaceAll(id, "-", "")
			sentinel := "vault_" + id
			username := sentinel + tc.suffix

			targetURL, err := url.Parse(connURL)
			require.True(t, err == nil, "parse connection URL (error not printed: it would include the password)")
			targetDB := setupDB
			if tc.contained {
				dbName := "vault_rotate_" + id
				_, err = setupDB.Exec("CREATE DATABASE [" + dbName + "]; ALTER DATABASE [" + dbName + "] SET CONTAINMENT = PARTIAL;")
				require.NoError(t, err, "create contained database")
				t.Cleanup(func() {
					_, err := setupDB.Exec("ALTER DATABASE [" + dbName + "] SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE [" + dbName + "]")
					require.NoError(t, err, "drop contained database")
				})
				params := targetURL.Query()
				params.Set("database", dbName)
				targetURL.RawQuery = params.Encode()
				targetDB, err = sql.Open("sqlserver", targetURL.String())
				require.NoError(t, err, "open contained database")
				t.Cleanup(func() { require.NoError(t, targetDB.Close(), "close contained database") })
			} else {
				t.Cleanup(func() {
					for _, principal := range []string{username, sentinel} {
						_, err := setupDB.Exec(`
IF EXISTS (SELECT 1 FROM sys.server_principals WHERE name = @name)
BEGIN
	DECLARE @stmt nvarchar(max) = N'DROP LOGIN ' + QUOTENAME(@name);
	EXEC (@stmt);
END`, sql.Named("name", principal))
						assert.NoError(t, err, "clean up test login")
					}
				})
			}

			createPrincipal := "CREATE LOGIN "
			if tc.contained {
				createPrincipal = "CREATE USER "
			}
			for _, principal := range []string{username, sentinel} {
				_, err = targetDB.Exec(`
DECLARE @stmt nvarchar(max) = @create + QUOTENAME(@name) + N' WITH PASSWORD = ' + QUOTENAME(@password, '''');
EXEC (@stmt);`, sql.Named("create", createPrincipal), sql.Named("name", principal), sql.Named("password", initPassword))
				require.NoError(t, err, "create target and similar-named principals")
			}

			db := new()
			dbtesting.AssertInitializeCircleCiTest(t, db, dbplugin.InitializeRequest{
				Config:           map[string]interface{}{"connection_url": targetURL.String(), "contained_db": tc.contained},
				VerifyConnection: true,
			})
			t.Cleanup(func() { dbtesting.AssertClose(t, db) })

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			_, err = db.UpdateUser(ctx, dbplugin.UpdateUserRequest{
				Username: username,
				Password: &dbplugin.ChangePassword{NewPassword: tc.password},
			})
			require.NoError(t, err, "rotate the password with the default statement")

			login := func(user, password string) error {
				loginURL := *targetURL
				loginURL.User = url.UserPassword(user, password)
				conn, err := sql.Open("sqlserver", loginURL.String())
				require.NoError(t, err, "open login connection")
				defer conn.Close()
				return conn.PingContext(ctx)
			}
			require.NoError(t, login(username, tc.password), "target authenticates with the rotated password")
			var sqlErr mssqldb.Error
			require.True(t, errors.As(login(username, initPassword), &sqlErr), "old password fails with a SQL Server error")
			require.Equal(t, int32(18456), sqlErr.Number, "old password is rejected by authentication")
			require.NoError(t, login(sentinel, initPassword), "similar-named principal keeps its password")
		})
	}

	t.Run("username over 128 characters", func(t *testing.T) {
		db := new()
		dbtesting.AssertInitializeCircleCiTest(t, db, dbplugin.InitializeRequest{
			Config:           map[string]interface{}{"connection_url": connURL},
			VerifyConnection: true,
		})
		t.Cleanup(func() { dbtesting.AssertClose(t, db) })

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		_, err := db.UpdateUser(ctx, dbplugin.UpdateUserRequest{
			Username: strings.Repeat("v", 129),
			Password: &dbplugin.ChangePassword{NewPassword: "Az9]Rotated"},
		})
		require.ErrorContains(t, err, "username must be at most 128 characters", "rotation must fail rather than silently skip")
	})
}

func assertCredsExist(t testing.TB, connURL, username, password string) {
	t.Helper()
	err := testCredsExist(connURL, username, password)
	if err != nil {
		t.Fatalf("Unable to log in as %q: %s", username, err)
	}
}

func assertCredsDoNotExist(t testing.TB, connURL, username, password string) {
	t.Helper()
	err := testCredsExist(connURL, username, password)
	if err == nil {
		t.Fatalf("Able to log in when it shouldn't")
	}
}

func assertContainedDBCredsDoNotExist(t testing.TB, connURL, username string) {
	t.Helper()
	err := testContainedDBCredsExist(connURL, username)
	assert.EqualError(t, err, "mssql: Cannot drop the user 'vaultuser', because it does not exist or you do not have permission.")
}

func testContainedDBCredsExist(connURL, username string) error {
	ctx := context.Background()
	// Log in
	db, err := sql.Open("mssql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	userStmt, err := db.PrepareContext(ctx, fmt.Sprintf("DROP USER [%s]", username))
	if err != nil {
		return err
	}
	_, err = userStmt.ExecContext(ctx)
	defer userStmt.Close()
	return err
}

func testCredsExist(connURL, username, password string) error {
	// Log in with the new creds
	parts := strings.Split(connURL, "@")
	connURL = fmt.Sprintf("sqlserver://%s:%s@%s", username, password, parts[1])
	db, err := sql.Open("mssql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

func createTestMSSQLUser(connURL string, username, password, query string) error {
	db, err := sql.Open("mssql", connURL)
	defer db.Close()
	if err != nil {
		return err
	}

	// Start a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	m := map[string]string{
		"name":     username,
		"password": password,
	}
	if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
		return err
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

const testMSSQLRole = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
CREATE USER [{{name}}] FOR LOGIN [{{name}}];
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::dbo TO [{{name}}];`

const testMSSQLLogin = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
`

const testMSSQLContainedLogin = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
CREATE USER [{{name}}] FOR LOGIN [{{name}}];
`

const testMSSQLContainedLoginAdmin = `
CREATE USER [{{name}}] WITH PASSWORD = '{{password}}';

ALTER ROLE db_datareader ADD MEMBER [{{name}}];
ALTER ROLE db_datawriter ADD MEMBER [{{name}}];
ALTER ROLE db_owner ADD MEMBER [{{name}}];
`
