package mysql

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	stdmysql "github.com/go-sql-driver/mysql"
	mysqlhelper "github.com/hashicorp/vault/helper/testhelpers/mysql"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

var _ dbplugin.Database = (*MySQL)(nil)

func TestMySQL_Initialize(t *testing.T) {
	cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test decoding a string value for max_open_connections
	connectionDetails = map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": "5",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMySQL_CreateUser(t *testing.T) {
	t.Run("missing creation statements", func(t *testing.T) {
		db := new(MetadataLen, MetadataLen, UsernameLen)

		usernameConfig := dbplugin.UsernameConfig{
			DisplayName: "test-long-displayname",
			RoleName:    "test-long-rolename",
		}

		username, password, err := db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
		if err == nil {
			t.Fatalf("expected err, got nil")
		}
		if username != "" {
			t.Fatalf("expected empty username, got [%s]", username)
		}
		if password != "" {
			t.Fatalf("expected empty password, got [%s]", password)
		}
	})

	t.Run("non-legacy", func(t *testing.T) {
		// Shared test container for speed - there should not be any overlap between the tests
		cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
		defer cleanup()

		connectionDetails := map[string]interface{}{
			"connection_url": connURL,
		}

		db := new(MetadataLen, MetadataLen, UsernameLen)
		_, err := db.Init(context.Background(), connectionDetails, true)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		testCreateUser(t, db, connURL)
	})

	t.Run("legacy", func(t *testing.T) {
		// Shared test container for speed - there should not be any overlap between the tests
		cleanup, connURL := mysqlhelper.PrepareTestContainer(t, true, "secret")
		defer cleanup()

		connectionDetails := map[string]interface{}{
			"connection_url": connURL,
		}

		db := new(credsutil.NoneLength, LegacyMetadataLen, LegacyUsernameLen)
		_, err := db.Init(context.Background(), connectionDetails, true)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		testCreateUser(t, db, connURL)
	})
}

func testCreateUser(t *testing.T, db *MySQL, connURL string) {
	type testCase struct {
		createStmts []string
	}

	tests := map[string]testCase{
		"create name": {
			createStmts: []string{`
				CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
				GRANT SELECT ON *.* TO '{{name}}'@'%';`,
			},
		},
		"create username": {
			createStmts: []string{`
				CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
				GRANT SELECT ON *.* TO '{{username}}'@'%';`,
			},
		},
		"prepared statement name": {
			createStmts: []string{`
				CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
				set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{name}}'@'%'");
				PREPARE grantStmt from @grants;
				EXECUTE grantStmt;
				DEALLOCATE PREPARE grantStmt;
				`,
			},
		},
		"prepared statement username": {
			createStmts: []string{`
				CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
				set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{username}}'@'%'");
				PREPARE grantStmt from @grants;
				EXECUTE grantStmt;
				DEALLOCATE PREPARE grantStmt;
				`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			usernameConfig := dbplugin.UsernameConfig{
				DisplayName: "test-long-displayname",
				RoleName:    "test-long-rolename",
			}

			statements := dbplugin.Statements{
				Creation: test.createStmts,
			}

			username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			// Test a second time to make sure usernames don't collide
			username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}
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
			statements: []string{`
				ALTER USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
			defer cleanup()

			connURL = strings.Replace(connURL, "root:secret", `{{username}}:{{password}}`, -1)

			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
				"username":       "root",
				"password":       "secret",
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			db := new(MetadataLen, MetadataLen, UsernameLen)
			_, err := db.Init(ctx, connectionDetails, true)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !db.Initialized {
				t.Fatal("Database should be initialized")
			}

			newConf, err := db.RotateRootCredentials(ctx, test.statements)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if newConf["password"] == "secret" {
				t.Fatal("password was not updated")
			}

			err = db.Close()
			if err != nil {
				t.Fatalf("err: %s", err)
			}
		})
	}
}

func TestMySQL_RevokeUser(t *testing.T) {
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
			revokeStmts: []string{`
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

	// Give a timeout just in case the test decides to be problematic
	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(initCtx, connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			statements := dbplugin.Statements{
				Creation: []string{`
					CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
					GRANT SELECT ON *.* TO '{{name}}'@'%';`,
				},
				Revocation: test.revokeStmts,
			}

			usernameConfig := dbplugin.UsernameConfig{
				DisplayName: "test",
				RoleName:    "test",
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			username, password, err := db.CreateUser(ctx, statements, usernameConfig, time.Now().Add(time.Minute))
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			err = db.RevokeUser(context.Background(), statements, username)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, username, password); err == nil {
				t.Fatal("Credentials were not revoked")
			}
		})
	}
}

func TestMySQL_SetCredentials(t *testing.T) {
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
			if err := mysqlhelper.TestCredsExist(t, connURL, dbUser, "password"); err != nil {
				t.Fatalf("Could not connect with credentials: %s", err)
			}

			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			db := new(MetadataLen, MetadataLen, UsernameLen)
			_, err := db.Init(ctx, connectionDetails, true)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			newPassword, err := db.GenerateCredentials(ctx)
			if err != nil {
				t.Fatalf("unable to generate password: %s", err)
			}

			userConfig := dbplugin.StaticUserConfig{
				Username: dbUser,
				Password: newPassword,
			}

			statements := dbplugin.Statements{
				Rotation: test.rotateStmts,
			}

			username, password, err := db.SetCredentials(ctx, statements, userConfig)
			if err != nil {
				t.Fatalf("err: %s", err)
			}
			if username != userConfig.Username {
				t.Fatalf("expected username [%s], got [%s]", userConfig.Username, username)
			}
			if password != userConfig.Password {
				t.Fatalf("expected password [%s] got [%s]", userConfig.Password, password)
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

func TestMySQL_Initialize_ReservedChars(t *testing.T) {
	pw := "#secret!%25#{@}"
	cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, pw)
	defer cleanup()

	// Revert password set to test replacement by db.Init
	connURL = strings.ReplaceAll(connURL, pw, "{{password}}")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
		"password":       pw,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
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
