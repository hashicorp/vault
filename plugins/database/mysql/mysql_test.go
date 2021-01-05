package mysql

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	stdmysql "github.com/go-sql-driver/mysql"
	mysqlhelper "github.com/hashicorp/vault/helper/testhelpers/mysql"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
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

	initReq := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := newMySQL(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Initialize(context.Background(), initReq)
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

	initReq = dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db = newMySQL(MetadataLen, MetadataLen, UsernameLen)
	_, err = db.Initialize(context.Background(), initReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMySQL_CreateUser(t *testing.T) {
	t.Run("missing creation statements", func(t *testing.T) {
		db := newMySQL(MetadataLen, MetadataLen, UsernameLen)

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
				Commands: []string{},
			},
			Password:   password,
			Expiration: time.Now().Add(time.Minute),
		}

		userResp, err := db.NewUser(context.Background(), createReq)
		if err == nil {
			t.Fatalf("expected err, got nil")
		}
		if userResp.Username != "" {
			t.Fatalf("expected empty username, got [%s]", userResp.Username)
		}
	})

	t.Run("non-legacy", func(t *testing.T) {
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

		db := newMySQL(MetadataLen, MetadataLen, UsernameLen)
		_, err := db.Initialize(context.Background(), initReq)
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

		initReq := dbplugin.InitializeRequest{
			Config:           connectionDetails,
			VerifyConnection: true,
		}

		db := newMySQL(credsutil.NoneLength, LegacyMetadataLen, LegacyUsernameLen)
		_, err := db.Initialize(context.Background(), initReq)
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
			createStmts: []string{
				`
				CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
				GRANT SELECT ON *.* TO '{{name}}'@'%';`,
			},
		},
		"create username": {
			createStmts: []string{
				`
				CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
				GRANT SELECT ON *.* TO '{{username}}'@'%';`,
			},
		},
		"prepared statement name": {
			createStmts: []string{
				`
				CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
				set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{name}}'@'%'");
				PREPARE grantStmt from @grants;
				EXECUTE grantStmt;
				DEALLOCATE PREPARE grantStmt;
				`,
			},
		},
		"prepared statement username": {
			createStmts: []string{
				`
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
					Commands: test.createStmts,
				},
				Password:   password,
				Expiration: time.Now().Add(time.Minute),
			}

			userResp, err := db.NewUser(context.Background(), createReq)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, userResp.Username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			// Test a second time to make sure usernames don't collide
			userResp, err = db.NewUser(context.Background(), createReq)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := mysqlhelper.TestCredsExist(t, connURL, userResp.Username, password); err != nil {
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

			db := newMySQL(MetadataLen, MetadataLen, UsernameLen)
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

	db := newMySQL(MetadataLen, MetadataLen, UsernameLen)
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

			db := newMySQL(MetadataLen, MetadataLen, UsernameLen)
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

	db := newMySQL(MetadataLen, MetadataLen, UsernameLen)
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
