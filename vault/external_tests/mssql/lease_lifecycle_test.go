// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package mssql_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	mssqlhelper "github.com/hashicorp/vault/helper/testhelpers/mssql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sqlTimeout = 10 * time.Second
	leaseTTL   = 60
	renewedTTL = 180
	// Short enough to keep expiry cases fast, long enough to outlast the
	// login and permission checks that run before the lease expires.
	expiryTTL    = 15
	revokeWait   = 2 * time.Minute
	pollInterval = 250 * time.Millisecond
)

// TestMSSQL_DefaultRevocationLeaseLifecycle verifies that API-issued SQL credentials
// are read-only, renew correctly, and disappear on explicit revocation or expiry
// without affecting an unrelated login, including when names contain SQL delimiters.
func TestMSSQL_DefaultRevocationLeaseLifecycle(t *testing.T) {
	t.Parallel()
	cleanup, connectionURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	t.Cleanup(cleanup)

	// Subtests share one SQL Server and run sequentially: concurrent login and
	// user DDL can deadlock (error 1205) and fail setup nondeterministically.
	for _, tc := range []struct {
		name     string
		suffix   string
		explicit bool
		ttl      int
	}{
		{name: "ordinary/renew_and_revoke", suffix: "reader", explicit: true, ttl: leaseTTL},
		{name: "bracket_apostrophe/renew_and_revoke", suffix: "reader]O'Brien", explicit: true, ttl: leaseTTL},
		{name: "ordinary/natural_expiry", suffix: "reader", ttl: expiryTTL},
		{name: "bracket_apostrophe/natural_expiry", suffix: "reader]O'Brien", ttl: expiryTTL},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id, err := uuid.GenerateUUID()
			require.NoError(t, err, "generate isolated SQL resource names")
			id = strings.ReplaceAll(id, "-", "")
			databaseName := "vault_lease_" + id
			username := "vault_" + id + "_" + tc.suffix
			sentinel := "sentinel_" + id
			sentinelSecret, err := uuid.GenerateUUID()
			require.NoError(t, err, "generate sentinel password")
			sentinelPassword := "Az9!" + sentinelSecret

			adminURL, err := url.Parse(connectionURL)
			require.True(t, err == nil, "parse SQL connection URL (error not printed: it would include the password)")
			require.Equal(t, "sqlserver", adminURL.Scheme, "MSSQL_URL must use the sqlserver URL format")
			query := adminURL.Query()
			query.Set("database", "master")
			adminURL.RawQuery = query.Encode()
			admin := openSQL(t, adminURL.String())

			// Register before creation so partially provisioned credentials are also removed.
			// Each cleanup has its own deadline and runs even if another cleanup fails.
			for _, name := range []string{username, sentinel} {
				t.Cleanup(func() {
					execSQL(t, admin, `
IF EXISTS (SELECT 1 FROM sys.server_principals WHERE name = @name)
BEGIN
	DECLARE @stmt nvarchar(max) = N'DROP LOGIN ' + QUOTENAME(@name);
	EXEC(@stmt);
END`, sql.Named("name", name))
				})
			}
			t.Cleanup(func() {
				execSQL(t, admin, `
IF DB_ID(@name) IS NOT NULL
BEGIN
	DECLARE @stmt nvarchar(max) = N'ALTER DATABASE ' + QUOTENAME(@name) +
		N' SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE ' + QUOTENAME(@name);
	EXEC(@stmt);
END`, sql.Named("name", databaseName))
			})
			execSQL(t, admin, `
DECLARE @stmt nvarchar(max) = N'CREATE DATABASE ' + QUOTENAME(@name);
EXEC(@stmt);`, sql.Named("name", databaseName))

			databaseURL := *adminURL
			query.Set("database", databaseName)
			databaseURL.RawQuery = query.Encode()
			database := openSQL(t, databaseURL.String())
			execSQL(t, database, "CREATE TABLE dbo.lease_probe (value int NOT NULL); INSERT INTO dbo.lease_probe VALUES (1);")
			execSQL(t, admin, `
DECLARE @stmt nvarchar(max) = N'CREATE LOGIN ' + QUOTENAME(@name) +
	N' WITH PASSWORD = ' + QUOTENAME(@password, '''');
EXEC(@stmt);`, sql.Named("name", sentinel), sql.Named("password", sentinelPassword))
			execSQL(t, database, `
DECLARE @stmt nvarchar(max) = N'CREATE USER ' + QUOTENAME(@name) + N' FOR LOGIN ' + QUOTENAME(@name) +
	N'; GRANT SELECT ON dbo.lease_probe TO ' + QUOTENAME(@name);
EXEC(@stmt);`, sql.Named("name", sentinel))
			assertReadOnly(t, &databaseURL, sentinel, sentinelPassword)
			sqlCtx, sqlCancel := context.WithTimeout(context.Background(), sqlTimeout)
			var sentinelLoginID, sentinelUserID int
			err = admin.QueryRowContext(sqlCtx, "SELECT principal_id FROM sys.server_principals WHERE name = @name",
				sql.Named("name", sentinel)).Scan(&sentinelLoginID)
			require.NoError(t, err, "record sentinel login identity")
			err = database.QueryRowContext(sqlCtx, "SELECT principal_id FROM sys.database_principals WHERE name = @name",
				sql.Named("name", sentinel)).Scan(&sentinelUserID)
			sqlCancel()
			require.NoError(t, err, "record sentinel database-user identity")

			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
			defer cancel()
			err = client.Sys().MountWithContext(ctx, "database", &api.MountInput{Type: "database"})
			require.NoError(t, err, "mount the database secrets engine")

			_, err = client.Logical().WriteWithContext(ctx, "sys/policies/password/sql-test", map[string]interface{}{
				"policy": `length = 24
rule "charset" {
	charset = "abcdefghijklmnopqrstuvwxyz"
	min-chars = 1
}
rule "charset" {
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	min-chars = 1
}
rule "charset" {
	charset = "0123456789"
	min-chars = 1
}`,
			})
			require.NoError(t, err, "restrict generated passwords to SQL literal-safe characters")
			_, err = client.Logical().WriteWithContext(ctx, "database/config/sql", map[string]interface{}{
				"plugin_name":       "mssql-database-plugin",
				"connection_url":    databaseURL.String(),
				"allowed_roles":     []string{"reader"},
				"username_template": username,
				"password_policy":   "sql-test",
				"verify_connection": true,
			})
			require.NoError(t, err, "configure the SQL plugin")

			// Creation uses literal substitution, not bound parameters. These controlled
			// names contain no double quotes or semicolons; passwords are alphanumeric.
			_, err = client.Logical().WriteWithContext(ctx, "database/roles/reader", map[string]interface{}{
				"db_name": "sql",
				"creation_statements": `SET QUOTED_IDENTIFIER ON;
CREATE LOGIN "{{name}}" WITH PASSWORD = '{{password}}';
CREATE USER "{{name}}" FOR LOGIN "{{name}}";
GRANT SELECT ON dbo.lease_probe TO "{{name}}";`,
				"default_ttl": tc.ttl,
				"max_ttl":     300,
			})
			require.NoError(t, err, "create a role using default revocation, with no revocation_statements")

			secret, err := client.Logical().ReadWithContext(ctx, "database/creds/reader")
			require.NoError(t, err, "issue SQL credentials")
			require.NotNil(t, secret, "credential issuance returns a secret")
			require.NotEmpty(t, secret.LeaseID, "credentials have a Vault-managed lease")
			require.True(t, secret.Renewable, "the credential lease supports renewal")
			require.Positive(t, secret.LeaseDuration, "the lease has a positive lifetime")
			actualUsername, ok := secret.Data["username"].(string)
			require.True(t, ok, "credential response includes a username string")
			require.Equal(t, username, actualUsername, "preserve the configured username exactly")
			password, ok := secret.Data["password"].(string)
			require.True(t, ok, "credential response includes a password string")
			require.NotEmpty(t, password, "return a non-empty SQL password")
			assertPrincipal(t, admin, database, username)
			assertReadOnly(t, &databaseURL, actualUsername, password)

			before, err := client.Sys().LookupWithContext(ctx, secret.LeaseID)
			require.NoError(t, err, "lookup the issued lease")
			require.NotNil(t, before, "lease lookup returns metadata")
			require.Equal(t, secret.LeaseID, before.Data["id"], "lookup identifies the issued lease")
			require.Equal(t, true, before.Data["renewable"], "stored lease is renewable")
			require.Positive(t, leaseSeconds(t, before), "stored lease has time remaining")

			if tc.explicit {
				renewed, renewErr := client.Sys().RenewWithContext(ctx, secret.LeaseID, renewedTTL)
				require.NoError(t, renewErr, "renew the credential lease through the API")
				require.NotNil(t, renewed, "renewal returns a lease")
				require.Equal(t, secret.LeaseID, renewed.LeaseID, "renewal preserves lease identity")
				after, lookupErr := client.Sys().LookupWithContext(ctx, secret.LeaseID)
				require.NoError(t, lookupErr, "read persisted renewal metadata")
				require.NotNil(t, after, "renewed lease remains registered")
				require.Equal(t, before.Data["issue_time"], after.Data["issue_time"], "renewal does not reissue credentials")
				require.True(t, leaseExpiry(t, after).After(leaseExpiry(t, before)), "renewal advances persisted expiry")
				require.Greater(t, leaseSeconds(t, after), leaseSeconds(t, before), "renewal increases persisted TTL")
				assertReadOnly(t, &databaseURL, actualUsername, password)

				err = client.Sys().RevokeWithContext(ctx, secret.LeaseID)
				require.NoError(t, err, "explicitly revoke with the plugin's default revocation")
			}

			// Expiry cases deliberately make no revoke call; both paths must remove
			// the SQL principals and the lease rather than only disabling access.
			require.EventuallyWithT(t, func(c *assert.CollectT) {
				sqlCtx, sqlCancel := context.WithTimeout(context.Background(), sqlTimeout)
				defer sqlCancel()
				var logins, users int
				catalogErr := admin.QueryRowContext(sqlCtx,
					"SELECT COUNT(*) FROM sys.server_principals WHERE name = @name",
					sql.Named("name", username)).Scan(&logins)
				require.NoError(c, catalogErr, "inspect target login after revocation")
				require.Zero(c, logins, "default revocation drops the target login")
				catalogErr = database.QueryRowContext(sqlCtx,
					"SELECT COUNT(*) FROM sys.database_principals WHERE name = @name",
					sql.Named("name", username)).Scan(&users)
				require.NoError(c, catalogErr, "inspect target database user after revocation")
				require.Zero(c, users, "default revocation drops the target database user")
				_, lookupErr := client.Sys().LookupWithContext(sqlCtx, secret.LeaseID)
				var responseErr *api.ResponseError
				require.ErrorAs(c, lookupErr, &responseErr, "revoked lease lookup must fail at the API")
				require.Equal(c, 400, responseErr.StatusCode, "lease lookup reports an invalid lease, not a transport failure")
				require.Contains(c, responseErr.Errors, "invalid lease", "Vault no longer stores the credential lease")
			}, revokeWait, pollInterval, "wait for default revocation, SQL cleanup, and lease removal")

			assertPrincipal(t, admin, database, sentinel)
			assertReadOnly(t, &databaseURL, sentinel, sentinelPassword)
			sqlCtx, sqlCancel = context.WithTimeout(context.Background(), sqlTimeout)
			var remainingLoginID, remainingUserID int
			err = admin.QueryRowContext(sqlCtx, "SELECT principal_id FROM sys.server_principals WHERE name = @name",
				sql.Named("name", sentinel)).Scan(&remainingLoginID)
			require.NoError(t, err, "inspect surviving login identity")
			require.Equal(t, sentinelLoginID, remainingLoginID, "preserve the original sentinel login")
			err = database.QueryRowContext(sqlCtx, "SELECT principal_id FROM sys.database_principals WHERE name = @name",
				sql.Named("name", sentinel)).Scan(&remainingUserID)
			sqlCancel()
			require.NoError(t, err, "inspect surviving user identity")
			require.Equal(t, sentinelUserID, remainingUserID, "preserve the original sentinel user")
			deniedURL := databaseURL
			deniedURL.User = url.UserPassword(actualUsername, password)
			denied := openSQL(t, deniedURL.String())
			sqlCtx, sqlCancel = context.WithTimeout(context.Background(), sqlTimeout)
			defer sqlCancel()
			loginErr := denied.PingContext(sqlCtx)
			var sqlErr mssql.Error
			require.True(t, errors.As(loginErr, &sqlErr), "revoked credentials must fail with a SQL Server error")
			require.Equal(t, int32(18456), sqlErr.Number, "fresh authentication fails because the login is revoked")
		})
	}
}

func openSQL(t *testing.T, connectionURL string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlserver", connectionURL)
	require.NoError(t, err, "open SQL connection")
	t.Cleanup(func() {
		require.NoError(t, db.Close(), "close SQL connection pool")
	})
	return db
}

func execSQL(t *testing.T, db *sql.DB, statement string, args ...interface{}) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), sqlTimeout)
	defer cancel()
	_, err := db.ExecContext(ctx, statement, args...)
	require.NoError(t, err, "execute SQL fixture statement")
}

func assertReadOnly(t *testing.T, connectionURL *url.URL, username, password string) {
	t.Helper()
	loginURL := *connectionURL
	loginURL.User = url.UserPassword(username, password)
	db := openSQL(t, loginURL.String())
	ctx, cancel := context.WithTimeout(context.Background(), sqlTimeout)
	defer cancel()
	var value int
	err := db.QueryRowContext(ctx, "SELECT value FROM dbo.lease_probe").Scan(&value)
	require.NoError(t, err, "authenticate and select using issued SQL credentials")
	require.Equal(t, 1, value, "read the isolated table's known value")
	_, writeErr := db.ExecContext(ctx, "INSERT INTO dbo.lease_probe VALUES (2)")
	var sqlErr mssql.Error
	require.True(t, errors.As(writeErr, &sqlErr), "write must fail with a SQL Server error")
	require.Equal(t, int32(229), sqlErr.Number, "INSERT is denied by permissions, not by a broken connection")
	require.NoError(t, db.Close(), "release application session before lease cleanup")
}

func assertPrincipal(t *testing.T, admin, database *sql.DB, username string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), sqlTimeout)
	defer cancel()
	var disabled bool
	err := admin.QueryRowContext(ctx, "SELECT is_disabled FROM sys.server_principals WHERE name = @name",
		sql.Named("name", username)).Scan(&disabled)
	require.NoError(t, err, "the SQL login still exists")
	require.False(t, disabled, "the SQL login is enabled")
	var users int
	err = database.QueryRowContext(ctx, "SELECT COUNT(*) FROM sys.database_principals WHERE name = @name",
		sql.Named("name", username)).Scan(&users)
	require.NoError(t, err, "lookup the SQL database user")
	require.Equal(t, 1, users, "the database user still exists")
}

func leaseSeconds(t *testing.T, secret *api.Secret) int64 {
	t.Helper()
	ttl, ok := secret.Data["ttl"].(json.Number)
	require.True(t, ok, "lease lookup exposes a numeric TTL")
	seconds, err := ttl.Int64()
	require.NoError(t, err, "decode lease TTL")
	return seconds
}

func leaseExpiry(t *testing.T, secret *api.Secret) time.Time {
	t.Helper()
	value, ok := secret.Data["expire_time"].(string)
	require.True(t, ok, "lease lookup exposes an expiry timestamp")
	expiry, err := time.Parse(time.RFC3339Nano, value)
	require.NoError(t, err, "decode lease expiry")
	return expiry
}
