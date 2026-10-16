// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	mssqlhelper "github.com/hashicorp/vault/helper/testhelpers/mssql"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMSSQLUsernameLifecycle checks exact generated, stored and authenticated
// names through plugin creation and default revocation, without harming a neighbor.
func TestMSSQLUsernameLifecycle(t *testing.T) {
	t.Parallel()
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	t.Cleanup(cleanup)

	tests := map[string]string{
		"plain":              "_reader",
		"bracket":            "]",
		"repeated brackets":  "]]",
		"bracket-shaped":     "[reader]",
		"apostrophe":         "'",
		"combined":           "]'",
		"SQL-looking":        "]; SELECT 1--",
		"space":              " reader",
		"double quote":       `"reader`,
		"dot":                ".reader",
		"connection syntax":  "@reader:/%",
		"Unicode":            "\u7528\u6237",
		"Unicode bracket":    "\u7528\u6237]",
		"Unicode apostrophe": "\u7528\u6237'",
		"127 characters":     strings.Repeat("]", 89),
		"128 characters":     strings.Repeat("]", 90),
	}
	// Subtests share one SQL Server and run sequentially: concurrent login and
	// user DDL can deadlock (error 1205) and fail setup nondeterministically.
	for name, suffix := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			id, err := uuid.GenerateUUID()
			require.NoError(t, err, "generate unique fixture")
			sentinel := "vault_" + strings.ReplaceAll(id, "-", "")
			username := sentinel + suffix
			dbName := sentinel

			admin, err := sql.Open("sqlserver", connURL)
			require.NoError(t, err, "open provisioning connection")
			t.Cleanup(func() { require.NoError(t, admin.Close(), "close provisioning connection") })
			_, err = admin.ExecContext(ctx, "CREATE DATABASE ["+dbName+"] COLLATE Latin1_General_100_CI_AS")
			require.NoError(t, err, "create isolated database")
			t.Cleanup(func() {
				cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Minute)
				defer cleanupCancel()
				_, err := admin.ExecContext(cleanupCtx, "ALTER DATABASE ["+dbName+"] SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE ["+dbName+"]")
				assert.NoError(t, err, "drop isolated database")
				for _, principal := range []string{username, sentinel} {
					_, err := admin.ExecContext(cleanupCtx, `
IF EXISTS (SELECT 1 FROM sys.server_principals WHERE name = @name)
BEGIN
	DECLARE @stmt nvarchar(max) = N'DROP LOGIN ' + QUOTENAME(@name);
	EXEC (@stmt);
END`, sql.Named("name", principal))
					assert.NoError(t, err, "remove fixture login")
				}
			})

			targetURL, err := url.Parse(connURL)
			require.True(t, err == nil, "parse fixture URL (error not printed: it would include the password)")
			params := targetURL.Query()
			params.Set("database", dbName)
			targetURL.RawQuery = params.Encode()
			observer, err := sql.Open("sqlserver", targetURL.String())
			require.NoError(t, err, "open catalog observer")
			t.Cleanup(func() { require.NoError(t, observer.Close(), "close catalog observer") })

			const password = "Vlt9-UsernameFixture!"
			_, err = observer.ExecContext(ctx, "CREATE LOGIN ["+sentinel+"] WITH PASSWORD = '"+password+"'; CREATE USER ["+sentinel+"] FOR LOGIN ["+sentinel+"]")
			require.NoError(t, err, "create similar-named sentinel")
			var sentinelID int
			err = observer.QueryRowContext(ctx, "SELECT principal_id FROM sys.server_principals WHERE name = @p1", sentinel).Scan(&sentinelID)
			require.NoError(t, err, "record sentinel identity")

			plugin := new()
			t.Cleanup(func() { require.NoError(t, plugin.Close(), "close plugin") })
			_, err = plugin.Initialize(ctx, dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":    targetURL.String(),
					"username_template": sentinel + fmt.Sprintf("{{%q}}", suffix),
				},
				VerifyConnection: true,
			})
			require.NoError(t, err, "initialize plugin")

			// Keep the controlled creation batch in a procedure: NewUser splits
			// administrator SQL on semicolons before substituting template values.
			_, err = observer.ExecContext(ctx, `
CREATE PROCEDURE dbo.create_fixture_user @password nvarchar(128) AS
BEGIN
	DECLARE @name nvarchar(128) = N'`+strings.ReplaceAll(username, "'", "''")+`';
	DECLARE @stmt nvarchar(max) = N'CREATE LOGIN ' + QUOTENAME(@name) +
		N' WITH PASSWORD = ' + QUOTENAME(@password, '''') +
		N'; CREATE USER ' + QUOTENAME(@name) + N' FOR LOGIN ' + QUOTENAME(@name);
	EXEC (@stmt);
END`)
			require.NoError(t, err, "prepare safely quoted creation fixture")
			created, err := plugin.NewUser(ctx, dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{DisplayName: "test", RoleName: "reader"},
				Statements: dbplugin.Statements{Commands: []string{
					"EXEC dbo.create_fixture_user @password = '{{password}}'",
				}},
				Password:   password,
				Expiration: time.Now().Add(time.Hour),
			})
			require.NoError(t, err, "create credentials through plugin")
			require.Equal(t, username, created.Username, "return the exact generated name")

			var userID int
			var storedUser, storedLogin string
			err = observer.QueryRowContext(ctx, "SELECT principal_id, name FROM sys.database_principals WHERE name = @p1", username).Scan(&userID, &storedUser)
			require.NoError(t, err, "observe created user")
			require.Equal(t, username, storedUser, "preserve database username")
			err = observer.QueryRowContext(ctx, "SELECT name FROM sys.server_principals WHERE name = @p1", username).Scan(&storedLogin)
			require.NoError(t, err, "observe created login")
			require.Equal(t, username, storedLogin, "preserve login name")

			credentialURL := *targetURL
			credentialURL.User = url.UserPassword(created.Username, password)
			application, err := sql.Open("sqlserver", credentialURL.String())
			require.NoError(t, err, "open application connection")
			t.Cleanup(func() { require.NoError(t, application.Close(), "close application connection") })
			var authenticatedLogin, authenticatedUser string
			err = application.QueryRowContext(ctx, "SELECT SUSER_SNAME(), USER_NAME()").Scan(&authenticatedLogin, &authenticatedUser)
			require.NoError(t, err, "authenticate with returned credentials")
			require.Equal(t, username, authenticatedLogin, "authenticate as the exact login")
			require.Equal(t, username, authenticatedUser, "map to the exact user")
			require.NoError(t, application.Close(), "release application session before revocation")

			_, err = plugin.DeleteUser(ctx, dbplugin.DeleteUserRequest{Username: created.Username})
			require.NoError(t, err, "use default revocation")
			var count int
			err = observer.QueryRowContext(ctx, "SELECT COUNT(*) FROM sys.database_principals WHERE principal_id = @p1", userID).Scan(&count)
			require.NoError(t, err, "query revoked user")
			require.Zero(t, count, "remove the original user")
			err = observer.QueryRowContext(ctx, "SELECT COUNT(*) FROM sys.server_principals WHERE name = @p1", username).Scan(&count)
			require.NoError(t, err, "query revoked login")
			require.Zero(t, count, "remove the original login, not merely disable it")
			err = observer.QueryRowContext(ctx, "SELECT COUNT(*) FROM sys.server_principals WHERE principal_id = @p1 AND is_disabled = 0", sentinelID).Scan(&count)
			require.NoError(t, err, "query sentinel login")
			require.Equal(t, 1, count, "preserve the enabled sentinel login")

			credentialURL.User = url.UserPassword(sentinel, password)
			survivor, err := sql.Open("sqlserver", credentialURL.String())
			require.NoError(t, err, "open sentinel connection")
			t.Cleanup(func() { require.NoError(t, survivor.Close(), "close sentinel connection") })
			err = survivor.QueryRowContext(ctx, "SELECT SUSER_SNAME(), USER_NAME()").Scan(&authenticatedLogin, &authenticatedUser)
			require.NoError(t, err, "sentinel must still authenticate")
			require.Equal(t, sentinel, authenticatedLogin, "preserve sentinel login mapping")
			require.Equal(t, sentinel, authenticatedUser, "preserve sentinel database user")
		})
	}
}
