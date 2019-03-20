package radius

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/ory/dockertest"
)

const (
	testSysTTL    = time.Hour * 10
	testSysMaxTTL = time.Hour * 20

	envRadiusRadiusHost = "RADIUS_HOST"
	envRadiusPort       = "RADIUS_PORT"
	envRadiusSecret     = "RADIUS_SECRET"
	envRadiusUsername   = "RADIUS_USERNAME"
	envRadiusUserPass   = "RADIUS_USERPASS"
)

func prepareRadiusTestContainer(t *testing.T) (func(), string, int) {
	if os.Getenv(envRadiusRadiusHost) != "" {
		port, _ := strconv.Atoi(os.Getenv(envRadiusPort))
		return func() {}, os.Getenv(envRadiusRadiusHost), port
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	runOpts := &dockertest.RunOptions{
		Repository:   "jumanjiman/radiusd",
		Cmd:          []string{"-f", "-l", "stdout"},
		ExposedPorts: []string{"1812/udp"},
		Tag:          "latest",
	}
	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		t.Fatalf("Could not start local radius docker container: %s", err)
	}

	cleanup := func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	port, _ := strconv.Atoi(resource.GetPort("1812/udp"))
	address := fmt.Sprintf("127.0.0.1")

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		// There's no straightfoward way to check the state, but the server starts
		// up quick so a 2 second sleep should be enough.
		time.Sleep(2 * time.Second)
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to radius docker container: %s", err)
	}
	return cleanup, address, port
}

func TestBackend_Config(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	configDataBasic := map[string]interface{}{
		"host":   "test.radius.hostname.com",
		"secret": "test-secret",
	}

	configDataMissingRequired := map[string]interface{}{
		"host": "test.radius.hostname.com",
	}

	configDataEmptyPort := map[string]interface{}{
		"host":   "test.radius.hostname.com",
		"port":   "",
		"secret": "test-secret",
	}

	configDataInvalidPort := map[string]interface{}{
		"host":   "test.radius.hostname.com",
		"port":   "notnumeric",
		"secret": "test-secret",
	}

	configDataInvalidBool := map[string]interface{}{
		"host":                       "test.radius.hostname.com",
		"secret":                     "test-secret",
		"unregistered_user_policies": "test",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		// PreCheck:       func() { testAccPreCheck(t) },
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, configDataBasic, false),
			testConfigWrite(t, configDataMissingRequired, true),
			testConfigWrite(t, configDataEmptyPort, true),
			testConfigWrite(t, configDataInvalidPort, true),
			testConfigWrite(t, configDataInvalidBool, true),
		},
	})
}

func TestBackend_users(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testStepUpdateUser(t, "web", "foo"),
			testStepUpdateUser(t, "web2", "foo"),
			testStepUpdateUser(t, "web3", "foo"),
			testStepUserList(t, []string{"web", "web2", "web3"}),
		},
	})
}

func TestBackend_acceptance(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	cleanup, host, port := prepareRadiusTestContainer(t)
	defer cleanup()

	// These defaults are specific to the jumanjiman/radiusd docker image
	username := os.Getenv(envRadiusUsername)
	if username == "" {
		username = "test"
	}

	password := os.Getenv(envRadiusUserPass)
	if password == "" {
		password = "test"
	}

	secret := os.Getenv(envRadiusSecret)
	if len(secret) == 0 {
		secret = "testing123"
	}

	configDataAcceptanceAllowUnreg := map[string]interface{}{
		"host":                       host,
		"port":                       strconv.Itoa(port),
		"secret":                     secret,
		"unregistered_user_policies": "policy1,policy2",
	}
	if configDataAcceptanceAllowUnreg["port"] == "" {
		configDataAcceptanceAllowUnreg["port"] = "1812"
	}

	configDataAcceptanceNoAllowUnreg := map[string]interface{}{
		"host":                       host,
		"port":                       strconv.Itoa(port),
		"secret":                     secret,
		"unregistered_user_policies": "",
	}
	if configDataAcceptanceNoAllowUnreg["port"] == "" {
		configDataAcceptanceNoAllowUnreg["port"] = "1812"
	}

	dataRealpassword := map[string]interface{}{
		"password": password,
	}

	dataWrongpassword := map[string]interface{}{
		"password": "wrongpassword",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		PreCheck:          testAccPreCheck(t, host, port),
		AcceptanceTest:    true,
		Steps: []logicaltest.TestStep{
			// Login with valid but unknown user will fail because unregistered_user_policies is empty
			testConfigWrite(t, configDataAcceptanceNoAllowUnreg, false),
			testAccUserLogin(t, username, dataRealpassword, true),
			// Once the user is registered auth will succeed
			testStepUpdateUser(t, username, ""),
			testAccUserLoginPolicy(t, username, dataRealpassword, []string{"default"}, false),

			testStepUpdateUser(t, username, "foopolicy"),
			testAccUserLoginPolicy(t, username, dataRealpassword, []string{"default", "foopolicy"}, false),
			testAccStepDeleteUser(t, username),

			// When unregistered_user_policies is specified, an unknown user will be granted access and granted the listed policies
			testConfigWrite(t, configDataAcceptanceAllowUnreg, false),
			testAccUserLoginPolicy(t, username, dataRealpassword, []string{"default", "policy1", "policy2"}, false),

			// More tests
			testAccUserLogin(t, "nonexistinguser", dataRealpassword, true),
			testAccUserLogin(t, username, dataWrongpassword, true),
			testStepUpdateUser(t, username, "foopolicy"),
			testAccUserLoginPolicy(t, username, dataRealpassword, []string{"default", "foopolicy"}, false),
			testStepUpdateUser(t, username, "foopolicy, secondpolicy"),
			testAccUserLoginPolicy(t, username, dataRealpassword, []string{"default", "foopolicy", "secondpolicy"}, false),
			testAccUserLoginPolicy(t, username, dataRealpassword, []string{"default", "foopolicy", "secondpolicy", "thirdpolicy"}, true),
		},
	})
}

func testAccPreCheck(t *testing.T, host string, port int) func() {
	return func() {
		if host == "" {
			t.Fatal("Host must be set for acceptance tests")
		}

		if port == 0 {
			t.Fatal("Port must be non-zero for acceptance tests")
		}
	}
}

func testConfigWrite(t *testing.T, d map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
		ErrorOk:   expectError,
	}
}

func testAccStepDeleteUser(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "users/" + n,
	}
}

func testStepUserList(t *testing.T, users []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "users",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("got error response: %#v", *resp)
			}

			if !reflect.DeepEqual(users, resp.Data["keys"].([]string)) {
				return fmt.Errorf("expected:\n%#v\ngot:\n%#v\n", users, resp.Data["keys"])
			}
			return nil
		},
	}
}

func testStepUpdateUser(
	t *testing.T, name string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testAccUserLogin(t *testing.T, user string, data map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login/" + user,
		Data:            data,
		ErrorOk:         expectError,
		Unauthenticated: true,
	}
}

func testAccUserLoginPolicy(t *testing.T, user string, data map[string]interface{}, policies []string, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login/" + user,
		Data:            data,
		ErrorOk:         expectError,
		Unauthenticated: true,
		//Check:           logicaltest.TestCheckAuth(policies),
		Check: func(resp *logical.Response) error {
			res := logicaltest.TestCheckAuth(policies)(resp)
			if res != nil && expectError {
				return nil
			}
			return res
		},
	}
}
