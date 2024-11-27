// Copyright (c) HashiCorp, Inc.
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
