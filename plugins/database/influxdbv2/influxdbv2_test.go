package influxdbv2

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stretchr/testify/require"
)

const createUserStatements = `influx user create --name {{.username}} --password {{.password}} --org {{.organization}}`

type Config struct {
	docker.ServiceURL
	Token string
}

var _ docker.ServiceConfig = &Config{}

func (c *Config) connectionParams() map[string]interface{} {
	pieces := strings.Split(c.Address(), ":")
	port, _ := strconv.Atoi(pieces[1])
	return map[string]interface{}{
		"host":         pieces[0],
		"port":         port,
		"token":        c.Token,
		"organization": "vault",
	}
}

func prepareInfluxdbTestContainer(t *testing.T) (func(), *Config, context.Context) {
	c := &Config{
		Token: "influx_token",
	}
	if host := os.Getenv("INFLUXDB_HOST"); host != "" {
		c.ServiceURL = *docker.NewServiceURL(url.URL{Scheme: "http", Host: host})
		return func() {}, c, context.Background()
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo: "influxdb",
		ImageTag:  "2.1-alpine",
		Env: []string{
			"INFLUXDB_DB=vault",
			"DOCKER_INFLUXDB_INIT_MODE=setup",
			"DOCKER_INFLUXDB_INIT_USERNAME=username",
			"DOCKER_INFLUXDB_INIT_PASSWORD=password",
			"DOCKER_INFLUXDB_INIT_ORG=vault",
			"DOCKER_INFLUXDB_INIT_BUCKET=vault",
			"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=" + c.Token,
			"INFLUXDB_HTTP_AUTH_ENABLED=true",
		},
		Ports: []string{"8086/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker InfluxDB: %s", err)
	}
	var ctxOuter context.Context
	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		c.ServiceURL = *docker.NewServiceURL(url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		})
		cli := influx.NewClient(c.URL().String(), c.Token)
		if err != nil {
			return nil, fmt.Errorf("error creating InfluxDB client: %w", err)
		}
		defer cli.Close()
		time.Sleep(10 * time.Second)
		_, err = cli.Ping(ctx)
		if err != nil {
			return nil, fmt.Errorf("error checking cluster status: %w", err)
		}
		ctxOuter = ctx
		return c, nil
	})
	if err != nil {
		t.Fatalf("Could not start docker InfluxDB: %s", err)
	}

	return svc.Cleanup, svc.Config.(*Config), ctxOuter
}

func TestInfluxdb_Initialize(t *testing.T) {
	cleanup, config, _ := prepareInfluxdbTestContainer(t)
	defer cleanup()

	type testCase struct {
		req               dbplugin.InitializeRequest
		expectedResponse  dbplugin.InitializeResponse
		expectErr         bool
		expectInitialized bool
	}

	tests := map[string]testCase{
		"port is an int": {
			req: dbplugin.InitializeRequest{
				Config:           makeConfig(config.connectionParams()),
				VerifyConnection: true,
			},
			expectedResponse: dbplugin.InitializeResponse{
				Config: config.connectionParams(),
			},
			expectErr:         false,
			expectInitialized: true,
		},
		"port is a string": {
			req: dbplugin.InitializeRequest{
				Config:           makeConfig(config.connectionParams(), "port", strconv.Itoa(config.connectionParams()["port"].(int))),
				VerifyConnection: true,
			},
			expectedResponse: dbplugin.InitializeResponse{
				Config: makeConfig(config.connectionParams(), "port", strconv.Itoa(config.connectionParams()["port"].(int))),
			},
			expectErr:         false,
			expectInitialized: true,
		},
		"missing config": {
			req: dbplugin.InitializeRequest{
				Config:           nil,
				VerifyConnection: true,
			},
			expectedResponse:  dbplugin.InitializeResponse{},
			expectErr:         true,
			expectInitialized: false,
		},
		"missing host": {
			req: dbplugin.InitializeRequest{
				Config:           makeConfig(config.connectionParams(), "host", ""),
				VerifyConnection: true,
			},
			expectedResponse:  dbplugin.InitializeResponse{},
			expectErr:         true,
			expectInitialized: false,
		},
		"missing token": {
			req: dbplugin.InitializeRequest{
				Config:           makeConfig(config.connectionParams(), "token", ""),
				VerifyConnection: true,
			},
			expectedResponse:  dbplugin.InitializeResponse{},
			expectErr:         true,
			expectInitialized: false,
		},
		"failed to validate connection": {
			req: dbplugin.InitializeRequest{
				// Host exists, but isn't a running instance
				Config:           makeConfig(config.connectionParams(), "host", "foobar://bad_connection"),
				VerifyConnection: true,
			},
			expectedResponse:  dbplugin.InitializeResponse{},
			expectErr:         true,
			expectInitialized: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := new()
			defer dbtesting.AssertClose(t, db)

			resp, err := db.Initialize(context.Background(), test.req)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(resp, test.expectedResponse) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResponse)
			}

			if test.expectInitialized && !db.Initialized {
				t.Fatalf("Database should be initialized but wasn't")
			} else if !test.expectInitialized && db.Initialized {
				t.Fatalf("Database was initiailized when it shouldn't")
			}
		})
	}
}

func makeConfig(rootConfig map[string]interface{}, keyValues ...interface{}) map[string]interface{} {
	if len(keyValues)%2 != 0 {
		panic("makeConfig must be provided with key and value pairs")
	}

	// Make a copy of the map so there isn't a chance of test bleedover between maps
	config := make(map[string]interface{}, len(rootConfig)+(len(keyValues)/2))
	for k, v := range rootConfig {
		config[k] = v
	}
	for i := 0; i < len(keyValues); i += 2 {
		k := keyValues[i].(string) // Will panic if the key field isn't a string and that's fine in a test
		v := keyValues[i+1]
		config[k] = v
	}
	return config
}

func TestInfluxdb_CreateUser_DefaultUsernameTemplate(t *testing.T) {
	cleanup, config, ctx := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	req := dbplugin.InitializeRequest{
		Config:           config.connectionParams(),
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, req)

	password := "nuozxby98523u89bdfnkjl"
	newUserReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "token",
			RoleName:    "mylongrolenamewithmanycharacters",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}
	resp := dbtesting.AssertNewUser(t, db, newUserReq)

	if resp.Username == "" {
		t.Fatalf("Missing username")
	}

	assertCredsExist(t, resp.Username, password, *config, ctx)

	require.Regexp(t, `^v_token_mylongrolenamew_[a-z0-9]{20}_[0-9]{10}$`, resp.Username)
}

func TestInfluxdb_CreateUser_CustomUsernameTemplate(t *testing.T) {
	cleanup, config, ctx := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()

	conf := config.connectionParams()
	conf["username_template"] = "{{.DisplayName}}_{{random 10}}"

	req := dbplugin.InitializeRequest{
		Config:           conf,
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, req)

	password := "nuozxby98523u89bdfnkjl"
	newUserReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "token",
			RoleName:    "mylongrolenamewithmanycharacters",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}
	resp := dbtesting.AssertNewUser(t, db, newUserReq)

	if resp.Username == "" {
		t.Fatalf("Missing username")
	}

	assertCredsExist(t, resp.Username, password, *config, ctx)

	require.Regexp(t, `^token_[a-zA-Z0-9]{10}$`, resp.Username)
}

func TestInfluxdb_CreateUser_AddedToOrganization(t *testing.T) {
	cleanup, config, ctx := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	req := dbplugin.InitializeRequest{
		Config:           config.connectionParams(),
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, req)

	password := "nuozxby98523u89bdfnkjl"
	newUserReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "token",
			RoleName:    "mylongrolenamewithmanycharacters",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}
	resp := dbtesting.AssertNewUser(t, db, newUserReq)

	if resp.Username == "" {
		t.Fatalf("Missing username")
	}

	assertUserHasOrganization(t, resp.Username, config.connectionParams()["organization"].(string), *config, ctx)
}

func TestUpdateUser_expiration(t *testing.T) {
	// This test should end up with a no-op since the expiration doesn't do anything in Influx

	cleanup, config, ctx := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	req := dbplugin.InitializeRequest{
		Config:           config.connectionParams(),
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, req)

	password := "nuozxby98523u89bdfnkjl"
	newUserReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}
	newUserResp := dbtesting.AssertNewUser(t, db, newUserReq)

	assertCredsExist(t, newUserResp.Username, password, *config, ctx)

	renewReq := dbplugin.UpdateUserRequest{
		Username: newUserResp.Username,
		Expiration: &dbplugin.ChangeExpiration{
			NewExpiration: time.Now().Add(5 * time.Minute),
		},
	}
	dbtesting.AssertUpdateUser(t, db, renewReq)

	// Make sure the user hasn't changed
	assertCredsExist(t, newUserResp.Username, password, *config, ctx)
}

func TestUpdateUser_password(t *testing.T) {
	cleanup, config, ctx := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	req := dbplugin.InitializeRequest{
		Config:           config.connectionParams(),
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, req)

	initialPassword := "nuozxby98523u89bdfnkjl"
	newUserReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   initialPassword,
		Expiration: time.Now().Add(1 * time.Minute),
	}
	newUserResp := dbtesting.AssertNewUser(t, db, newUserReq)

	assertCredsExist(t, newUserResp.Username, initialPassword, *config, ctx)

	newPassword := "y89qgmbzadiygry8uazodijnb"
	newPasswordReq := dbplugin.UpdateUserRequest{
		Username: newUserResp.Username,
		Password: &dbplugin.ChangePassword{
			NewPassword: newPassword,
		},
	}
	dbtesting.AssertUpdateUser(t, db, newPasswordReq)

	assertCredsDoNotExist(t, newUserResp.Username, initialPassword, *config, ctx)
	assertCredsExist(t, newUserResp.Username, newPassword, *config, ctx)
}

// TestInfluxdb_RevokeDeletedUser tests attempting to revoke a user that was
// deleted externally. Guards against a panic, see
// https://github.com/hashicorp/vault/issues/6734
// Updated to attempt to delete a user that never existed to replicate a similar scenario since
// the cleanup function from `prepareInfluxdbTestContainer` does not do anything if using an
// external InfluxDB instance rather than spinning one up for the test.
func TestInfluxdb_RevokeDeletedUser(t *testing.T) {
	cleanup, config, _ := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	req := dbplugin.InitializeRequest{
		Config:           config.connectionParams(),
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, req)

	// attempt to revoke a user that does not exist
	delReq := dbplugin.DeleteUserRequest{
		Username: "someuser",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.DeleteUser(ctx, delReq)
	if err == nil {
		t.Fatalf("Expected err, got nil")
	}
}

func TestInfluxdb_RevokeUser(t *testing.T) {
	cleanup, config, ctx := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	req := dbplugin.InitializeRequest{
		Config:           config.connectionParams(),
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, req)

	initialPassword := "nuozxby98523u89bdfnkjl"
	newUserReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   initialPassword,
		Expiration: time.Now().Add(1 * time.Minute),
	}
	newUserResp := dbtesting.AssertNewUser(t, db, newUserReq)

	assertCredsExist(t, newUserResp.Username, initialPassword, *config, ctx)

	delReq := dbplugin.DeleteUserRequest{
		Username: newUserResp.Username,
	}
	dbtesting.AssertDeleteUser(t, db, delReq)
	assertCredsDoNotExist(t, newUserResp.Username, initialPassword, *config, ctx)
}

func assertCredsExist(t testing.TB, username, password string, c Config, ctx context.Context) {
	t.Helper()
	err := testCredsExist(username, password, c, ctx)
	if err != nil {
		t.Fatalf("Could not log in as %q %s", username, err)
	}
}

func assertCredsDoNotExist(t testing.TB, username, password string, c Config, ctx context.Context) {
	t.Helper()
	err := testCredsExist(username, password, c, ctx)
	if err == nil {
		t.Fatalf("Able to log in as %q when it shouldn't", username)
	}
}

func assertUserHasOrganization(t testing.TB, username, organizationName string, c Config, ctx context.Context) {
	t.Helper()
	err := testUserHasOrganization(username, organizationName, c, ctx)
	if err != nil {
		t.Fatalf("User %q does not belong to organization %q when it should", username, organizationName)
	}
}

func testCredsExist(username, password string, c Config, ctx context.Context) error {
	cli := influx.NewClient(c.URL().String(), c.Token)
	defer cli.Close()
	_, err := cli.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error checking server ping: %w", err)
	}
	err = cli.UsersAPI().SignIn(ctx, username, password)
	if err != nil {
		return fmt.Errorf("error querying influxdb server: %w", err)
	}
	return nil
}

func testUserHasOrganization(username, organizationName string, c Config, ctx context.Context) error {
	cli := influx.NewClient(c.URL().String(), c.Token)
	defer cli.Close()
	_, err := cli.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error checking server ping: %w", err)
	}
	res, err := cli.UsersAPI().FindUserByName(ctx, username)
	if err != nil {
		return fmt.Errorf("error querying influxdb server: %w", err)
	}
	organizations, err := cli.OrganizationsAPI().FindOrganizationsByUserID(ctx, *res.Id)
	for _, organization := range *organizations {
		if organization.Name == organizationName {
			return nil
		}
	}
	return fmt.Errorf("user %s does not belong to organization %s", username, organizationName)
}
