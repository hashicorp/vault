package hana

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
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
		req           dbplugin.NewUserRequest
		expectErr     bool
		assertUser    func(t testing.TB, connURL, username, password string)
	}

	tests := map[string]testCase{
		"no creation statements": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test-test",
					RoleName:    "test-test",
				},
				Statements: dbplugin.Statements{},
				Password:   "AG4qagho_dsvZ",
				Expiration: time.Now().Add(1 * time.Second),
			},
			expectErr:     true,
			assertUser:    assertCredsDoNotExist,
		},
		"with creation statements": {
			req: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test-test",
					RoleName:    "test-test",
				},
				Statements: dbplugin.Statements{
					Commands: []string{testHANARole},
				},
				Password:   "AG4qagho_dsvZ",
				Expiration: time.Now().Add(1 * time.Second),
			},
			expectErr:     false,
			assertUser:    assertCredsExist,
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

			createResp, err := db.NewUser(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, received nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			test.assertUser(t, connURL, createResp.Username, test.req.Password)
		})
	}
}

func TestHANA_UpdateUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	type testCase struct {
		req              dbplugin.UpdateUserRequest
		startingPassword string
		expectErrOnLogin bool
		expectedErrMsg   string
	}

	tests := map[string]testCase{
		"no update statements": {
			req: dbplugin.UpdateUserRequest{
				Password: &dbplugin.ChangePassword{
					NewPassword: "this_is_ALSO_Thirty_2_characters_",
				},
			},
			startingPassword: "this_is_Thirty_2_characters_wow_",
			expectErrOnLogin: true,
			expectedErrMsg:   "user is forced to change password",
		},
		"with custom update statements": {
			req: dbplugin.UpdateUserRequest{
				Password: &dbplugin.ChangePassword{
					NewPassword: "this_is_ALSO_Thirty_2_characters_",
					Statements: dbplugin.Statements{
						Commands: []string{testHANAUpdate},
					},
				},
			},
			startingPassword: "this_is_Thirty_2_characters_wow_",
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

			newReq := dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "test-test",
					RoleName:    "test-test",
				},
				Password: test.startingPassword,
				Statements: dbplugin.Statements{
					Commands: []string{testHANARole},
				},
				Expiration: time.Now().Add(time.Hour),
			}

			userResp := dbtesting.AssertNewUser(t, db, newReq)
			assertCredsExist(t, connURL, userResp.Username, test.startingPassword)

			test.req.Username = userResp.Username

			dbtesting.AssertUpdateUser(t, db, test.req)
			err := testCredsExist(t, connURL, userResp.Username, test.req.Password.NewPassword)
			if test.expectErrOnLogin {
				if err == nil {
					t.Fatalf("Able to login with new creds when expecting an issue")
				} else if test.expectedErrMsg != "" && !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Fatalf("Expected error message to contain \"%s\", received: %s", test.expectedErrMsg, err)
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
		req              dbplugin.DeleteUserRequest
	}

	tests := map[string]testCase{
		"no update statements": {
			req: dbplugin.DeleteUserRequest{},
		},
		"with custom update statements": {
			req: dbplugin.DeleteUserRequest{
				Statements: dbplugin.Statements{
					Commands: []string{testHANADrop},
				},
			},
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

			test.req.Username = userResp.Username

			dbtesting.AssertDeleteUser(t, db, test.req)
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

const testHANARole = `
CREATE USER {{name}} PASSWORD "{{password}}" NO FORCE_FIRST_PASSWORD_CHANGE VALID UNTIL '{{expiration}}';`

const testHANADrop = `
DROP USER {{name}} CASCADE;`

const testHANAUpdate = `
ALTER USER {{name}} PASSWORD "{{password}}" NO FORCE_FIRST_PASSWORD_CHANGE;`
