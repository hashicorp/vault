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

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/testhelpers/certhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/postgresql"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getPostgreSQL(t *testing.T, options map[string]interface{}) (*PostgreSQL, func()) {
	cleanup, connURL := postgresql.PrepareTestContainer(t)

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

// TestPostgreSQL_InitializeMultiHost tests the functionality of Postgres's
// multi-host connection strings.
func TestPostgreSQL_InitializeMultiHost(t *testing.T) {
	cleanup, connURL := postgresql.PrepareTestContainerMultiHost(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": 5,
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
}

// TestPostgreSQL_InitializeSSLInlineFeatureFlag tests that the VAULT_PLUGIN_USE_POSTGRES_SSLINLINE
// flag guards against unwanted usage of the deprecated SSL client authentication path.
// TODO: remove this when we remove the underlying feature in a future SDK version
func TestPostgreSQL_InitializeSSLInlineFeatureFlag(t *testing.T) {
	// set the flag to true so we can call PrepareTestContainerWithSSL
	// which does a validation check on the connection
	t.Setenv(pluginutil.PluginUsePostgresSSLInline, "true")

	// Create certificates for postgres authentication
	caCert := certhelpers.NewCert(t, certhelpers.CommonName("ca"), certhelpers.IsCA(true), certhelpers.SelfSign())
	clientCert := certhelpers.NewCert(t, certhelpers.CommonName("postgres"), certhelpers.DNS("localhost"), certhelpers.Parent(caCert))
	cleanup, connURL := postgresql.PrepareTestContainerWithSSL(t, "verify-ca", caCert, clientCert, false)
	t.Cleanup(cleanup)

	type testCase struct {
		env           string
		wantErr       bool
		expectedError string
	}

	tests := map[string]testCase{
		"feature flag is true": {
			env:           "true",
			wantErr:       false,
			expectedError: "",
		},
		"feature flag is unset or empty": {
			env:     "",
			wantErr: true,
			// this error is expected because the env var unset means we are
			// using pgx's native connection string parsing which does not
			// support inlining of the certificate material in the sslrootcert,
			// sslcert, and sslkey fields
			expectedError: "error verifying connection",
		},
		"feature flag is false": {
			env:           "false",
			wantErr:       true,
			expectedError: "failed to open postgres connection with deprecated funtion",
		},
		"feature flag is invalid": {
			env:           "foo",
			wantErr:       true,
			expectedError: "failed to open postgres connection with deprecated funtion",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// update the env var with the value we are testing
			t.Setenv(pluginutil.PluginUsePostgresSSLInline, test.env)

			connectionDetails := map[string]interface{}{
				"connection_url":       connURL,
				"max_open_connections": 5,
			}

			req := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := new()
			_, err := dbtesting.VerifyInitialize(t, db, req)
			if test.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			} else if test.wantErr && !strings.Contains(err.Error(), test.expectedError) {
				t.Fatalf("got: %s, want: %s", err.Error(), test.expectedError)
			}

			if !test.wantErr && !db.Initialized {
				t.Fatal("Database should be initialized")
			}

			if err := db.Close(); err != nil {
				t.Fatalf("err: %s", err)
			}
			// unset for the next test case
			os.Unsetenv(pluginutil.PluginUsePostgresSSLInline)
		})
	}
}

// TestPostgreSQL_InitializeSSLInline tests that we can successfully authenticate
// with a postgres server via ssl with a URL connection string or DSN (key/value)
// for each ssl mode.
// TODO: remove this when we remove the underlying feature in a future SDK version
func TestPostgreSQL_InitializeSSLInline(t *testing.T) {
	// required to enable the sslinline custom parsing
	t.Setenv(pluginutil.PluginUsePostgresSSLInline, "true")

	type testCase struct {
		sslMode       string
		useDSN        bool
		useFallback   bool
		wantErr       bool
		expectedError string
	}

	tests := map[string]testCase{
		"disable sslmode": {
			sslMode:       "disable",
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode": {
			sslMode: "allow",
			wantErr: false,
		},
		"prefer sslmode": {
			sslMode: "prefer",
			wantErr: false,
		},
		"require sslmode": {
			sslMode: "require",
			wantErr: false,
		},
		"verify-ca sslmode": {
			sslMode: "verify-ca",
			wantErr: false,
		},
		"disable sslmode with DSN": {
			sslMode:       "disable",
			useDSN:        true,
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode with DSN": {
			sslMode: "allow",
			useDSN:  true,
			wantErr: false,
		},
		"prefer sslmode with DSN": {
			sslMode: "prefer",
			useDSN:  true,
			wantErr: false,
		},
		"require sslmode with DSN": {
			sslMode: "require",
			useDSN:  true,
			wantErr: false,
		},
		"verify-ca sslmode with DSN": {
			sslMode: "verify-ca",
			useDSN:  true,
			wantErr: false,
		},
		"disable sslmode with fallback": {
			sslMode:       "disable",
			useFallback:   true,
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode with fallback": {
			sslMode:     "allow",
			useFallback: true,
		},
		"prefer sslmode with fallback": {
			sslMode:     "prefer",
			useFallback: true,
		},
		"require sslmode with fallback": {
			sslMode:     "require",
			useFallback: true,
		},
		"verify-ca sslmode with fallback": {
			sslMode:     "verify-ca",
			useFallback: true,
		},
		"disable sslmode with DSN with fallback": {
			sslMode:       "disable",
			useDSN:        true,
			useFallback:   true,
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode with DSN with fallback": {
			sslMode:     "allow",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
		"prefer sslmode with DSN with fallback": {
			sslMode:     "prefer",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
		"require sslmode with DSN with fallback": {
			sslMode:     "require",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
		"verify-ca sslmode with DSN with fallback": {
			sslMode:     "verify-ca",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create certificates for postgres authentication
			caCert := certhelpers.NewCert(t, certhelpers.CommonName("ca"), certhelpers.IsCA(true), certhelpers.SelfSign())
			clientCert := certhelpers.NewCert(t, certhelpers.CommonName("postgres"), certhelpers.DNS("localhost"), certhelpers.Parent(caCert))
			cleanup, connURL := postgresql.PrepareTestContainerWithSSL(t, test.sslMode, caCert, clientCert, test.useFallback)
			t.Cleanup(cleanup)

			if test.useDSN {
				var err error
				connURL, err = dbutil.ParseURL(connURL)
				if err != nil {
					t.Fatal(err)
				}
			}
			connectionDetails := map[string]interface{}{
				"connection_url":       connURL,
				"max_open_connections": 5,
			}

			req := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := new()
			_, err := dbtesting.VerifyInitialize(t, db, req)
			if test.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			} else if test.wantErr && !strings.Contains(err.Error(), test.expectedError) {
				t.Fatalf("got: %s, want: %s", err.Error(), test.expectedError)
			}

			if !test.wantErr && !db.Initialized {
				t.Fatal("Database should be initialized")
			}

			if err := db.Close(); err != nil {
				t.Fatalf("err: %s", err)
			}
		})
	}
}

// TestPostgreSQL_InitializeSSL tests that we can successfully authenticate
// with a postgres server via ssl with a URL connection string or DSN (key/value)
// for each ssl mode.
func TestPostgreSQL_InitializeSSL(t *testing.T) {
	type testCase struct {
		sslMode       string
		useDSN        bool
		useFallback   bool
		wantErr       bool
		expectedError string
	}

	tests := map[string]testCase{
		"disable sslmode": {
			sslMode:       "disable",
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode": {
			sslMode: "allow",
			wantErr: false,
		},
		"prefer sslmode": {
			sslMode: "prefer",
			wantErr: false,
		},
		"require sslmode": {
			sslMode: "require",
			wantErr: false,
		},
		"verify-ca sslmode": {
			sslMode: "verify-ca",
			wantErr: false,
		},
		"verify-full sslmode": {
			sslMode: "verify-full",
			wantErr: false,
		},
		"disable sslmode with DSN": {
			sslMode:       "disable",
			useDSN:        true,
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode with DSN": {
			sslMode: "allow",
			useDSN:  true,
			wantErr: false,
		},
		"prefer sslmode with DSN": {
			sslMode: "prefer",
			useDSN:  true,
			wantErr: false,
		},
		"require sslmode with DSN": {
			sslMode: "require",
			useDSN:  true,
			wantErr: false,
		},
		"verify-ca sslmode with DSN": {
			sslMode: "verify-ca",
			useDSN:  true,
			wantErr: false,
		},
		"verify-full sslmode with DSN": {
			sslMode: "verify-full",
			useDSN:  true,
			wantErr: false,
		},
		"disable sslmode with fallback": {
			sslMode:       "disable",
			useFallback:   true,
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode with fallback": {
			sslMode:     "allow",
			useFallback: true,
		},
		"prefer sslmode with fallback": {
			sslMode:     "prefer",
			useFallback: true,
		},
		"require sslmode with fallback": {
			sslMode:     "require",
			useFallback: true,
		},
		"verify-ca sslmode with fallback": {
			sslMode:     "verify-ca",
			useFallback: true,
		},
		"verify-full sslmode with fallback": {
			sslMode:     "verify-full",
			useFallback: true,
		},
		"disable sslmode with DSN with fallback": {
			sslMode:       "disable",
			useDSN:        true,
			useFallback:   true,
			wantErr:       true,
			expectedError: "error verifying connection",
		},
		"allow sslmode with DSN with fallback": {
			sslMode:     "allow",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
		"prefer sslmode with DSN with fallback": {
			sslMode:     "prefer",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
		"require sslmode with DSN with fallback": {
			sslMode:     "require",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
		"verify-ca sslmode with DSN with fallback": {
			sslMode:     "verify-ca",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
		"verify-full sslmode with DSN with fallback": {
			sslMode:     "verify-full",
			useDSN:      true,
			useFallback: true,
			wantErr:     false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create certificates for postgres authentication
			caCert := certhelpers.NewCert(t, certhelpers.CommonName("ca"), certhelpers.IsCA(true), certhelpers.SelfSign())
			clientCert := certhelpers.NewCert(t, certhelpers.CommonName("postgres"), certhelpers.DNS("localhost"), certhelpers.Parent(caCert))
			cleanup, connURL := postgresql.PrepareTestContainerWithSSL(t, test.sslMode, caCert, clientCert, test.useFallback)
			t.Cleanup(cleanup)

			if test.useDSN {
				var err error
				connURL, err = dbutil.ParseURL(connURL)
				if err != nil {
					t.Fatal(err)
				}
			}
			connectionDetails := map[string]interface{}{
				"connection_url":       connURL,
				"max_open_connections": 5,
				"tls_certificate":      string(clientCert.CombinedPEM()),
				"private_key":          string(clientCert.PrivateKeyPEM()),
				"tls_ca":               string(caCert.CombinedPEM()),
			}

			req := dbplugin.InitializeRequest{
				Config:           connectionDetails,
				VerifyConnection: true,
			}

			db := new()
			_, err := dbtesting.VerifyInitialize(t, db, req)
			if test.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			} else if test.wantErr && !strings.Contains(err.Error(), test.expectedError) {
				t.Fatalf("got: %s, want: %s", err.Error(), test.expectedError)
			}

			if !test.wantErr && !db.Initialized {
				t.Fatal("Database should be initialized")
			}

			if err := db.Close(); err != nil {
				t.Fatalf("err: %s", err)
			}
		})
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
	cleanup, connURL := postgresql.PrepareTestContainer(t)
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

// TestPostgreSQL_Initialize_SelfManaged_OSS tests the initialization of
// the self-managed flow and ensures an error is returned on OSS.
func TestPostgreSQL_Initialize_SelfManaged_OSS(t *testing.T) {
	if constants.IsEnterprise {
		t.Skip("this test is only valid on OSS")
	}

	cleanup, url := postgresql.PrepareTestContainerSelfManaged(t)
	defer cleanup()

	connURL := fmt.Sprintf("postgresql://{{username}}:{{password}}@%s/postgres?sslmode=disable", url.Host)

	testCases := []struct {
		name              string
		connectionDetails map[string]interface{}
		wantErr           bool
		errContains       string
	}{
		{
			name: "no parameters set",
			connectionDetails: map[string]interface{}{
				"connection_url": connURL,
				"self_managed":   false,
				"username":       "",
				"password":       "",
			},
			wantErr:     true,
			errContains: "must either provide username/password or set self-managed to 'true'",
		},
		{
			name: "both sets of parameters set",
			connectionDetails: map[string]interface{}{
				"connection_url": connURL,
				"self_managed":   true,
				"username":       "test",
				"password":       "test",
			},
			wantErr:     true,
			errContains: "cannot use both self-managed and vault-managed workflows",
		},
		{
			name: "either username/password with self-managed",
			connectionDetails: map[string]interface{}{
				"connection_url": connURL,
				"self_managed":   true,
				"username":       "test",
				"password":       "",
			},
			wantErr:     true,
			errContains: "cannot use both self-managed and vault-managed workflows",
		},
		{
			name: "cache not implemented",
			connectionDetails: map[string]interface{}{
				"connection_url": connURL,
				"self_managed":   true,
				"username":       "",
				"password":       "",
			},
			wantErr:     true,
			errContains: "self-managed static roles only available in Vault Enterprise",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := dbplugin.InitializeRequest{
				Config:           tc.connectionDetails,
				VerifyConnection: true,
			}

			db := new()
			_, err := dbtesting.VerifyInitialize(t, db, req)
			if err == nil && tc.wantErr {
				t.Fatalf("got: %s, wantErr: %t", err, tc.wantErr)
			}

			if err != nil && !strings.Contains(err.Error(), tc.errContains) {
				t.Fatalf("expected error: %s, received error: %s", tc.errContains, err)
			}

			if !tc.wantErr && !db.Initialized {
				t.Fatal("Database should be initialized")
			}

			if err := db.Close(); err != nil {
				t.Fatalf("err closing DB: %s", err)
			}
		})
	}
}

// TestPostgreSQL_PasswordAuthentication tests that the default "password_authentication" is "none", and that
// an error is returned if an invalid "password_authentication" is provided.
func TestPostgreSQL_PasswordAuthentication(t *testing.T) {
	cleanup, connURL := postgresql.PrepareTestContainer(t)
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
	cleanup, connURL := postgresql.PrepareTestContainer(t)
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

// TestUpdateUser_SelfManaged_OSS checks basic validation
// for self-managed fields and confirms an error is returned on OSS
func TestUpdateUser_SelfManaged_OSS(t *testing.T) {
	if constants.IsEnterprise {
		t.Skip("this test is only valid on OSS")
	}

	// Shared test container for speed - there should not be any overlap between the tests
	db, cleanup := getPostgreSQL(t, nil)
	defer cleanup()

	updateReq := dbplugin.UpdateUserRequest{
		Username: "static",
		Password: &dbplugin.ChangePassword{
			NewPassword: "somenewpassword",
			Statements: dbplugin.Statements{
				Commands: nil,
			},
		},
		SelfManagedPassword: "test",
	}

	expectedErr := "self-managed static roles only available in Vault Enterprise"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.UpdateUser(ctx, updateReq)
	if err == nil {
		t.Fatalf("err expected, got nil")
	}
	if !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("err expected: %s, got: %s", expectedErr, err)
	}
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
	cleanup, connURL := postgresql.PrepareTestContainer(t)
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

func getHost(url string) string {
	splitCreds := strings.Split(url, "@")[1]

	return strings.Split(splitCreds, "/")[0]
}
