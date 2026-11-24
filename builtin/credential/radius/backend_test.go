// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package radius

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/logical"
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
	if strings.Contains(runtime.GOARCH, "arm") {
		t.Skip("Skipping, as this image is not supported on ARM architectures")
	}

	if os.Getenv(envRadiusRadiusHost) != "" {
		port, _ := strconv.Atoi(os.Getenv(envRadiusPort))
		return func() {}, os.Getenv(envRadiusRadiusHost), port
	}

	// Now allow any client to connect to this radiusd instance by writing our
	// own clients.conf file.
	//
	// This is necessary because we lack control over the container's network
	// IPs. We might be running in Circle CI (with variable IPs per new
	// network) or in Podman (which uses an entirely different set of default
	// ranges than Docker).
	//
	// See also: https://freeradius.org/radiusd/man/clients.conf.html
	ctx := context.Background()
	clientsConfig := `
client 0.0.0.0/1 {
 ipaddr = 0.0.0.0/1
 secret = testing123
 shortname = all-clients-first
}

client 128.0.0.0/1 {
 ipaddr = 128.0.0.0/1
 secret = testing123
 shortname = all-clients-second
}
`

	containerfile := `
FROM docker.mirror.hashicorp.services/jumanjiman/radiusd:latest

COPY clients.conf /etc/raddb/clients.conf
`

	bCtx := docker.NewBuildContext()
	bCtx["clients.conf"] = docker.PathContentsFromBytes([]byte(clientsConfig))

	imageName := "vault_radiusd_any_client"
	imageTag := "latest"

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     imageName,
		ImageTag:      imageTag,
		ContainerName: "radiusd",
		Cmd:           []string{"-f", "-l", "stdout", "-X"},
		Ports:         []string{"1812/udp"},
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	})
	if err != nil {
		t.Fatalf("Could not provision docker service runner: %s", err)
	}

	output, err := runner.BuildImage(ctx, containerfile, bCtx,
		docker.BuildRemove(true), docker.BuildForceRemove(true),
		docker.BuildPullParent(true),
		docker.BuildTags([]string{imageName + ":" + imageTag}))
	if err != nil {
		t.Fatalf("Could not build new image: %v", err)
	}

	t.Logf("Image build output: %v", string(output))

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		time.Sleep(2 * time.Second)
		return docker.NewServiceHostPort(host, port), nil
	})
	if err != nil {
		t.Fatalf("Could not start docker radiusd: %s", err)
	}

	pieces := strings.Split(svc.Config.Address(), ":")
	port, _ := strconv.Atoi(pieces[1])
	return svc.Cleanup, pieces[0], port
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

// TestBackend_Config_CaseInsensitiveWarnings verifies mixed-case and case-collision warnings are returned when enabling case-insensitive names.
func TestBackend_Config_CaseInsensitiveWarnings(t *testing.T) {
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

	baseConfig := map[string]interface{}{
		"host":                   "test.radius.hostname.com",
		"secret":                 "test-secret",
		"case_insensitive_names": false,
	}

	enableCaseInsensitive := map[string]interface{}{
		"case_insensitive_names": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, baseConfig, false),
			testStepUpdateUser(t, "Admin", "admin"),
			testStepUpdateUser(t, "admin", "attacker"),
			testConfigWriteWithCheck(t, enableCaseInsensitive, false, func(resp *logical.Response) error {
				if resp == nil {
					return fmt.Errorf("expected warning response, got nil")
				}

				if len(resp.Warnings) != 2 {
					return fmt.Errorf("expected 2 warnings, got: %v", resp.Warnings)
				}

				all := strings.Join(resp.Warnings, " ")
				if !strings.Contains(all, "uppercase characters") {
					return fmt.Errorf("warning did not mention mixed-case usernames: %v", resp.Warnings)
				}

				if !strings.Contains(all, "differ only by case") || !strings.Contains(all, "ExampleUser") || !strings.Contains(all, "exampleuser") {
					return fmt.Errorf("warning did not include expected collision guidance: %v", resp.Warnings)
				}

				return nil
			}),
		},
	})
}

// TestBackend_Config_CaseInsensitiveMixedCaseWarnings verifies a mixed-case-only dataset returns only the mixed-case warning.
func TestBackend_Config_CaseInsensitiveMixedCaseWarnings(t *testing.T) {
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

	baseConfig := map[string]interface{}{
		"host":                   "test.radius.hostname.com",
		"secret":                 "test-secret",
		"case_insensitive_names": false,
	}

	enableCaseInsensitive := map[string]interface{}{
		"case_insensitive_names": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, baseConfig, false),
			testStepUpdateUser(t, "Admin", "admin"),
			testConfigWriteWithCheck(t, enableCaseInsensitive, false, func(resp *logical.Response) error {
				if resp == nil {
					return fmt.Errorf("expected warning response, got nil")
				}

				if len(resp.Warnings) != 1 {
					return fmt.Errorf("expected 1 warning, got: %v", resp.Warnings)
				}

				all := strings.Join(resp.Warnings, " ")
				if !strings.Contains(all, "uppercase characters") {
					return fmt.Errorf("warning did not mention mixed-case usernames: %v", resp.Warnings)
				}

				if strings.Contains(all, "differ only by case") {
					return fmt.Errorf("did not expect collision warning for a single mixed-case username: %v", resp.Warnings)
				}

				return nil
			}),
		},
	})
}

// TestBackend_Config_CaseInsensitiveNoWarnings verifies no warnings are returned when existing usernames are already lowercase.
func TestBackend_Config_CaseInsensitiveNoWarnings(t *testing.T) {
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

	baseConfig := map[string]interface{}{
		"host":                   "test.radius.hostname.com",
		"secret":                 "test-secret",
		"case_insensitive_names": false,
	}

	enableCaseInsensitive := map[string]interface{}{
		"case_insensitive_names": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, baseConfig, false),
			testStepUpdateUser(t, "admin", "attacker"),
			testConfigWriteWithCheck(t, enableCaseInsensitive, false, func(resp *logical.Response) error {
				if resp == nil {
					return nil
				}

				if len(resp.Warnings) != 0 {
					return fmt.Errorf("expected no warnings, got: %v", resp.Warnings)
				}

				return nil
			}),
		},
	})
}

func testConfigWriteWithCheck(t *testing.T, d map[string]interface{}, expectError bool, check func(*logical.Response) error) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
		ErrorOk:   expectError,
		Check:     check,
	}
}

func testAccUserLoginPolicyAndMetaUsername(
	t *testing.T,
	user string,
	data map[string]interface{},
	policies []string,
	expectError bool,
	expectedUsername string,
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login/" + user,
		Data:            data,
		ErrorOk:         expectError,
		Unauthenticated: true,
		Check: func(resp *logical.Response) error {
			res := logicaltest.TestCheckAuth(policies)(resp)
			if res != nil {
				if expectError {
					return nil
				}

				return res
			}

			if resp == nil || resp.Auth == nil || resp.Auth.Metadata == nil {
				return fmt.Errorf("expected auth metadata, got none")
			}

			got := resp.Auth.Metadata["username"]
			if got != expectedUsername {
				return fmt.Errorf("expected metadata username %q, got %q", expectedUsername, got)
			}

			return nil
		},
	}
}

func policiesFromData(v interface{}) []string {
	switch policies := v.(type) {
	case []string:
		return policies
	case []interface{}:
		out := make([]string, 0, len(policies))
		for _, p := range policies {
			if s, ok := p.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
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
			testStepUpdateUser(t, "Web4", "foo"),
			testStepUserList(t, []string{"Web4", "web", "web2", "web3"}),
		},
	})
}

func TestBackend_acceptance(t *testing.T) {
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
	t *testing.T, name string, policies string,
) logicaltest.TestStep {
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
		// Check:           logicaltest.TestCheckAuth(policies),
		Check: func(resp *logical.Response) error {
			res := logicaltest.TestCheckAuth(policies)(resp)
			if res != nil && expectError {
				return nil
			}
			return res
		},
	}
}

// TestBackend_Users_CaseInsensitiveCollisionBlocked verifies case-variant writes are rejected when case-insensitive names are enabled.
func TestBackend_Users_CaseInsensitiveCollisionBlocked(t *testing.T) {
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

	configTrue := map[string]interface{}{
		"host":                   "test.radius.hostname.com",
		"secret":                 "test-secret",
		"case_insensitive_names": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, configTrue, false),
			testStepUpdateUser(t, "admin", "attacker"),
			{
				Operation: logical.UpdateOperation,
				Path:      "users/Admin",
				Data: map[string]interface{}{
					"policies": "admin",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp == nil || !resp.IsError() {
						return fmt.Errorf("expected error response, got: %#v", resp)
					}

					msg, _ := resp.Data["error"].(string)
					if !strings.Contains(msg, "collides with existing username") {
						return fmt.Errorf("unexpected error message: %q", msg)
					}

					return nil
				},
			},
		},
	})
}

// TestBackend_Users_CaseInsensitiveCollisionBlockedLegacyMixedCase verifies a legacy mixed-case key cannot be rewritten through the normalized CRUD path.
func TestBackend_Users_CaseInsensitiveCollisionBlockedLegacyMixedCase(t *testing.T) {
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

	baseConfig := map[string]interface{}{
		"host":                   "test.radius.hostname.com",
		"secret":                 "test-secret",
		"case_insensitive_names": false,
	}

	enableCaseInsensitive := map[string]interface{}{
		"case_insensitive_names": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, baseConfig, false),
			testStepUpdateUser(t, "Admin", "attacker"),
			testConfigWriteWithCheck(t, enableCaseInsensitive, false, func(resp *logical.Response) error {
				if resp == nil {
					return fmt.Errorf("expected warning response, got nil")
				}

				if len(resp.Warnings) == 0 {
					return fmt.Errorf("expected warnings when enabling case_insensitive_names with legacy mixed-case usernames")
				}

				return nil
			}),
			{
				Operation: logical.UpdateOperation,
				Path:      "users/Admin",
				Data: map[string]interface{}{
					"policies": "admin",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp == nil || !resp.IsError() {
						return fmt.Errorf("expected error response, got: %#v", resp)
					}

					msg, _ := resp.Data["error"].(string)
					if !strings.Contains(msg, "collides with existing username") {
						return fmt.Errorf("unexpected error message: %q", msg)
					}

					return nil
				},
			},
		},
	})
}

// TestBackend_Users_CaseInsensitiveReadDeleteNormalized verifies read and delete normalize mixed-case usernames to the canonical lowercase key.
func TestBackend_Users_CaseInsensitiveReadDeleteNormalized(t *testing.T) {
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

	configTrue := map[string]interface{}{
		"host":                   "test.radius.hostname.com",
		"secret":                 "test-secret",
		"case_insensitive_names": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, configTrue, false),
			testStepUpdateUser(t, "admin", "attacker"),
			{
				Operation: logical.ReadOperation,
				Path:      "users/AdMiN",
				Check: func(resp *logical.Response) error {
					if resp == nil || resp.IsError() {
						return fmt.Errorf("expected successful response, got: %#v", resp)
					}

					got := policiesFromData(resp.Data["policies"])
					if !reflect.DeepEqual(got, []string{"attacker"}) {
						return fmt.Errorf("expected policies [attacker], got: %#v", got)
					}

					return nil
				},
			},
			{
				Operation: logical.DeleteOperation,
				Path:      "users/AdMiN",
			},
			{
				Operation: logical.ReadOperation,
				Path:      "users/admin",
				Check: func(resp *logical.Response) error {
					if resp != nil {
						return fmt.Errorf("expected canonical user to be deleted, got %#v", resp.Data)
					}

					return nil
				},
			},
		},
	})
}

// TestBackend_Acceptance_CaseInsensitiveLoginNormalization verifies login metadata uses the normalized lowercase username when enabled.
func TestBackend_Acceptance_CaseInsensitiveLoginNormalization(t *testing.T) {
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

	username := os.Getenv(envRadiusUsername)
	if username == "" {
		username = "test"
	}

	password := os.Getenv(envRadiusUserPass)
	if password == "" {
		password = "test"
	}

	secret := os.Getenv(envRadiusSecret)
	if secret == "" {
		secret = "testing123"
	}

	cfg := map[string]interface{}{
		"host":                   host,
		"port":                   strconv.Itoa(port),
		"secret":                 secret,
		"case_insensitive_names": true,
	}

	if cfg["port"] == "" {
		cfg["port"] = "1812"
	}

	lower := strings.ToLower(username)
	upper := strings.ToUpper(username)

	loginData := map[string]interface{}{
		"password": password,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		PreCheck:          testAccPreCheck(t, host, port),
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, cfg, false),
			testStepUpdateUser(t, lower, "foopolicy"),
			testAccUserLoginPolicyAndMetaUsername(t, upper, loginData, []string{"default", "foopolicy"}, false, lower),
		},
	})
}
