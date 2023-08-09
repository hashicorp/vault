// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hana

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/stretchr/testify/require"
)

func TestHANA_Initialize(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	expectedConfig := copyConfig(connectionDetails)

	initReq := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	initResp := dbtesting.AssertInitialize(t, db, initReq)
	defer dbtesting.AssertClose(t, db)

	if !reflect.DeepEqual(initResp.Config, expectedConfig) {
		t.Fatalf("Actual config: %#v\nExpected config: %#v", initResp.Config, expectedConfig)
	}
}

// this test will leave a lingering user on the system
func TestHANA_NewUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}

	connURL := os.Getenv("HANA_URL")

	type testCase struct {
		commands   []string
		expectErr  bool
		assertUser func(t testing.TB, connURL, username, password string)
	}

	tests := map[string]testCase{
		"no creation statements": {
			commands:   []string{},
			expectErr:  true,
			assertUser: assertCredsDoNotExist,
		},
		"with creation statements": {
			commands:   []string{testHANARole},
			expectErr:  false,
			assertUser: assertCredsExist,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
			}

			initReq := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := new()
			dbtesting.AssertInitialize(t, db, initReq)
			defer dbtesting.AssertClose(t, db)

			req := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test-test",
					RoleName:    "test-test",
				},
				Statements: dbplugin.Statements{
					Commands: test.commands,
				},
				Password:   "AG4qagho_dsvZ",
				Expiration: time.Now().Add(1 * time.Second),
			}

			createResp, err := db.NewUser(context.Background(), req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, received nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			test.assertUser(t, connURL, createResp.Username, req.Password)
		})
	}
}

func TestHANA_UpdateUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	type testCase struct {
		commands         []string
		expectErrOnLogin bool
		expectedErrMsg   string
	}

	tests := map[string]testCase{
		"no update statements": {
			commands:         []string{},
			expectErrOnLogin: true,
			expectedErrMsg:   "user is forced to change password",
		},
		"with custom update statements": {
			commands:         []string{testHANAUpdate},
			expectErrOnLogin: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
			}

			initReq := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := new()
			dbtesting.AssertInitialize(t, db, initReq)
			defer dbtesting.AssertClose(t, db)

			password := "this_is_Thirty_2_characters_wow_"
			newReq := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test-test",
					RoleName:    "test-test",
				},
				Password: password,
				Statements: dbplugin.Statements{
					Commands: []string{testHANARole},
				},
				Expiration: time.Now().Add(time.Hour),
			}

			userResp := dbtesting.AssertNewUser(t, db, newReq)
			assertCredsExist(t, connURL, userResp.Username, password)

			req := dbplugin.UpdateUserRequest{
				Username: userResp.Username,
				Password: &dbplugin.ChangePassword{
					NewPassword: "this_is_ALSO_Thirty_2_characters_",
					Statements: dbplugin.Statements{
						Commands: test.commands,
					},
				},
			}

			dbtesting.AssertUpdateUser(t, db, req)
			err := testCredsExist(t, connURL, userResp.Username, req.Password.NewPassword)
			if test.expectErrOnLogin {
				if err == nil {
					t.Fatalf("Able to login with new creds when expecting an issue")
				} else if test.expectedErrMsg != "" && !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Fatalf("Expected error message to contain %q, received: %s", test.expectedErrMsg, err)
				}
			}
			if !test.expectErrOnLogin && err != nil {
				t.Fatalf("Unable to login: %s", err)
			}
		})
	}
}

func TestHANA_DeleteUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	type testCase struct {
		commands []string
	}

	tests := map[string]testCase{
		"no update statements": {
			commands: []string{},
		},
		"with custom update statements": {
			commands: []string{testHANADrop},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
			}

			initReq := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := new()
			dbtesting.AssertInitialize(t, db, initReq)
			defer dbtesting.AssertClose(t, db)

			password := "this_is_Thirty_2_characters_wow_"

			newReq := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test-test",
					RoleName:    "test-test",
				},
				Password: password,
				Statements: dbplugin.Statements{
					Commands: []string{testHANARole},
				},
				Expiration: time.Now().Add(time.Hour),
			}

			userResp := dbtesting.AssertNewUser(t, db, newReq)
			assertCredsExist(t, connURL, userResp.Username, password)

			req := dbplugin.DeleteUserRequest{
				Username: userResp.Username,
				Statements: dbplugin.Statements{
					Commands: test.commands,
				},
			}

			dbtesting.AssertDeleteUser(t, db, req)
			assertCredsDoNotExist(t, connURL, userResp.Username, password)
		})
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	// Log in with the new creds
	parts := strings.Split(connURL, "@")
	connURL = fmt.Sprintf("hdb://%s:%s@%s", username, password, parts[1])
	db, err := sql.Open("hdb", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

func assertCredsExist(t testing.TB, connURL, username, password string) {
	t.Helper()
	err := testCredsExist(t, connURL, username, password)
	if err != nil {
		t.Fatalf("Unable to log in as %q: %s", username, err)
	}
}

func assertCredsDoNotExist(t testing.TB, connURL, username, password string) {
	t.Helper()
	err := testCredsExist(t, connURL, username, password)
	if err == nil {
		t.Fatalf("Able to log in when we should not be able to")
	}
}

func copyConfig(config map[string]interface{}) map[string]interface{} {
	newConfig := map[string]interface{}{}
	for k, v := range config {
		newConfig[k] = v
	}
	return newConfig
}

func TestHANA_DefaultUsernameTemplate(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	initReq := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	dbtesting.AssertInitialize(t, db, initReq)

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	resp := dbtesting.AssertNewUser(t, db, dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements: dbplugin.Statements{
			Commands: []string{testHANARole},
		},
		Expiration: time.Now().Add(5 * time.Minute),
	})
	username := resp.Username

	if resp.Username == "" {
		t.Fatalf("Missing username")
	}

	testCredsExist(t, connURL, username, password)

	require.Regexp(t, `^V_TEST_TEST_[A-Z0-9]{20}_[0-9]{10}$`, resp.Username)

	defer dbtesting.AssertClose(t, db)
}

func TestHANA_CustomUsernameTemplate(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url":    connURL,
		"username_template": "{{.DisplayName}}_{{random 10}}",
	}

	initReq := dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	dbtesting.AssertInitialize(t, db, initReq)

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	resp := dbtesting.AssertNewUser(t, db, dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements: dbplugin.Statements{
			Commands: []string{testHANARole},
		},
		Expiration: time.Now().Add(5 * time.Minute),
	})
	username := resp.Username

	if resp.Username == "" {
		t.Fatalf("Missing username")
	}

	testCredsExist(t, connURL, username, password)

	require.Regexp(t, `^TEST_[A-Z0-9]{10}$`, resp.Username)

	defer dbtesting.AssertClose(t, db)
}

const testHANARole = `
CREATE USER {{name}} PASSWORD "{{password}}" NO FORCE_FIRST_PASSWORD_CHANGE VALID UNTIL '{{expiration}}';`

const testHANADrop = `
DROP USER {{name}} CASCADE;`

const testHANAUpdate = `
ALTER USER {{name}} PASSWORD "{{password}}" NO FORCE_FIRST_PASSWORD_CHANGE;`
