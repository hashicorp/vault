// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

// This file contains all "large"/expensive tests. These are running requests against a running backend

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestPlugin_lifecycle(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	env := []string{fmt.Sprintf("%s=%s", pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)}
	vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "mock-v4-database-plugin", consts.PluginTypeDatabase, "", "TestBackend_PluginMain_MockV4", env)
	vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "mock-v5-database-plugin", consts.PluginTypeDatabase, "", "TestBackend_PluginMain_MockV5", env)
	vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "mock-v6-database-plugin-muxed", consts.PluginTypeDatabase, "", "TestBackend_PluginMain_MockV6Multiplexed", env)

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys
	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to database backend")
	}
	defer b.Cleanup(context.Background())

	type testCase struct {
		dbName                string
		dbType                string
		configData            map[string]interface{}
		assertDynamicUsername stringAssertion
		assertDynamicPassword stringAssertion
	}

	tests := map[string]testCase{
		"v4": {
			dbName: "mockv4",
			dbType: "mock-v4-database-plugin",
			configData: map[string]interface{}{
				"name":              "mockv4",
				"plugin_name":       "mock-v4-database-plugin",
				"connection_url":    "sample_connection_url",
				"verify_connection": true,
				"allowed_roles":     []string{"*"},
				"username":          "mockv4-user",
				"password":          "mysecurepassword",
			},
			assertDynamicUsername: assertStringPrefix("mockv4_user_"),
			assertDynamicPassword: assertStringPrefix("mockv4_"),
		},
		"v5": {
			dbName: "mockv5",
			dbType: "mock-v5-database-plugin",
			configData: map[string]interface{}{
				"connection_url":    "sample_connection_url",
				"plugin_name":       "mock-v5-database-plugin",
				"verify_connection": true,
				"allowed_roles":     []string{"*"},
				"name":              "mockv5",
				"username":          "mockv5-user",
				"password":          "mysecurepassword",
			},
			assertDynamicUsername: assertStringPrefix("mockv5_user_"),
			assertDynamicPassword: assertStringRegex("^[a-zA-Z0-9-]{20}"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var cleanupReqs []*logical.Request
			defer func() {
				// Do not defer cleanup directly so that we can populate the
				// slice before the function gets executed.
				cleanup(t, b, cleanupReqs)
			}()

			// /////////////////////////////////////////////////////////////////
			// Configure
			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      fmt.Sprintf("config/%s", test.dbName),
				Storage:   config.StorageView,
				Data:      test.configData,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			assertNoRespData(t, resp)

			cleanupReqs = append(cleanupReqs, &logical.Request{
				Operation: logical.DeleteOperation,
				Path:      fmt.Sprintf("config/%s", test.dbName),
				Storage:   config.StorageView,
			})

			// /////////////////////////////////////////////////////////////////
			// Rotate root credentials
			req = &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("rotate-root/%s", test.dbName),
				Storage:   config.StorageView,
			}
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err = b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			assertNoRespData(t, resp)

			// /////////////////////////////////////////////////////////////////
			// Dynamic credentials

			// Create role
			dynamicRoleName := "dynamic-role"
			req = &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("roles/%s", dynamicRoleName),
				Storage:   config.StorageView,
				Data: map[string]interface{}{
					"db_name":     test.dbName,
					"default_ttl": "5s",
					"max_ttl":     "1m",
				},
			}
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err = b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			assertNoRespData(t, resp)

			cleanupReqs = append(cleanupReqs, &logical.Request{
				Operation: logical.DeleteOperation,
				Path:      fmt.Sprintf("roles/%s", dynamicRoleName),
				Storage:   config.StorageView,
			})

			// Generate credentials
			req = &logical.Request{
				Operation: logical.ReadOperation,
				Path:      fmt.Sprintf("creds/%s", dynamicRoleName),
				Storage:   config.StorageView,
			}
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err = b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			assertRespHasData(t, resp)

			// TODO: Figure out how to make a call to the cluster that gives back a lease ID
			// And also rotates the secret out after its TTL

			// /////////////////////////////////////////////////////////////////
			// Static credentials

			// Create static role
			staticRoleName := "static-role"
			req = &logical.Request{
				Operation: logical.CreateOperation,
				Path:      fmt.Sprintf("static-roles/%s", staticRoleName),
				Storage:   config.StorageView,
				Data: map[string]interface{}{
					"db_name":         test.dbName,
					"username":        "static-username",
					"rotation_period": "5",
				},
			}
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err = b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			assertNoRespData(t, resp)

			cleanupReqs = append(cleanupReqs, &logical.Request{
				Operation: logical.DeleteOperation,
				Path:      fmt.Sprintf("static-roles/%s", staticRoleName),
				Storage:   config.StorageView,
			})

			// Get credentials
			req = &logical.Request{
				Operation: logical.ReadOperation,
				Path:      fmt.Sprintf("static-creds/%s", staticRoleName),
				Storage:   config.StorageView,
			}
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err = b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			assertRespHasData(t, resp)
		})
	}
}

func TestPlugin_VersionSelection(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	for _, version := range []string{"v11.0.0", "v11.0.1-rc1", "v2.0.0"} {
		vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "mock-v5-database-plugin", consts.PluginTypeDatabase, version, "TestBackend_PluginMain_MockV5", []string{})
	}

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys
	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to database backend")
	}
	defer b.Cleanup(context.Background())

	test := func(t *testing.T, selectVersion, expectedVersion string) func(t *testing.T) {
		return func(t *testing.T) {
			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "config/db",
				Storage:   config.StorageView,
				Data: map[string]interface{}{
					"connection_url":    "sample_connection_url",
					"plugin_name":       "mock-v5-database-plugin",
					"plugin_version":    selectVersion,
					"verify_connection": true,
					"allowed_roles":     []string{"*"},
					"name":              "mockv5",
					"username":          "mockv5-user",
					"password":          "mysecurepassword",
				},
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			assertNoRespData(t, resp)

			defer func() {
				_, err := b.HandleRequest(context.Background(), &logical.Request{
					Operation: logical.DeleteOperation,
					Path:      "config/db",
					Storage:   config.StorageView,
				})
				if err != nil {
					t.Fatal(err)
				}
			}()

			req = &logical.Request{
				Operation: logical.ReadOperation,
				Path:      "config/db",
				Storage:   config.StorageView,
			}
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err = b.HandleRequest(ctx, req)
			assertErrIsNil(t, err)
			assertRespHasNoErr(t, resp)
			if resp.Data["plugin_version"].(string) != expectedVersion {
				t.Fatalf("Expected version %q but got %q", expectedVersion, resp.Data["plugin_version"].(string))
			}
		}
	}

	for name, tc := range map[string]struct {
		selectVersion   string
		expectedVersion string
	}{
		"no version specified, selects latest in the absence of unversioned plugins": {
			selectVersion:   "",
			expectedVersion: "v11.0.1-rc1",
		},
		"specific version selected": {
			selectVersion:   "11.0.0",
			expectedVersion: "v11.0.0",
		},
	} {
		t.Run(name, test(t, tc.selectVersion, tc.expectedVersion))
	}

	// Register a newer version of the plugin, and ensure that's the new default version selected.
	vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "mock-v5-database-plugin", consts.PluginTypeDatabase, "v11.0.1", "TestBackend_PluginMain_MockV5", []string{})
	t.Run("no version specified, new latest version selected", test(t, "", "v11.0.1"))

	// Register an unversioned plugin and ensure that is now selected when no version is specified.
	vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "mock-v5-database-plugin", consts.PluginTypeDatabase, "", "TestBackend_PluginMain_MockV5", []string{})
	for name, tc := range map[string]struct {
		selectVersion   string
		expectedVersion string
	}{
		"no version specified, selects unversioned": {
			selectVersion:   "",
			expectedVersion: "",
		},
		"specific version selected": {
			selectVersion:   "v2.0.0",
			expectedVersion: "v2.0.0",
		},
	} {
		t.Run(name, test(t, tc.selectVersion, tc.expectedVersion))
	}
}

func TestPlugin_VersionMustBeExplicitlyUpgraded(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys
	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to database backend")
	}
	defer b.Cleanup(context.Background())

	configData := func(extraData ...string) map[string]interface{} {
		data := map[string]interface{}{
			"connection_url":    "sample_connection_url",
			"plugin_name":       "mysql-database-plugin",
			"verify_connection": false,
			"allowed_roles":     []string{"*"},
			"username":          "mockv5-user",
			"password":          "mysecurepassword",
		}
		if len(extraData)%2 != 0 {
			t.Fatal("Expected an even number of args in extraData")
		}
		for i := 0; i < len(extraData); i += 2 {
			data[extraData[i]] = extraData[i+1]
		}
		return data
	}

	readVersion := func() string {
		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "config/db",
			Storage:   config.StorageView,
		})
		assertErrIsNil(t, err)
		assertRespHasNoErr(t, resp)
		return resp.Data["plugin_version"].(string)
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/db",
		Storage:   config.StorageView,
		Data:      configData(),
	})
	assertErrIsNil(t, err)
	assertRespHasNoErr(t, resp)
	assertNoRespData(t, resp)

	version := readVersion()
	expectedVersion := ""
	if version != expectedVersion {
		t.Fatalf("Expected version %q but got %q", expectedVersion, version)
	}

	// Register versioned plugin, and check that a new write to existing config doesn't upgrade the plugin implicitly.
	vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "mysql-database-plugin", consts.PluginTypeDatabase, "v1.0.0", "TestBackend_PluginMain_MockV5", []string{})
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/db",
		Storage:   config.StorageView,
		Data:      configData(),
	})
	assertErrIsNil(t, err)
	assertRespHasNoErr(t, resp)
	assertNoRespData(t, resp)

	version = readVersion()
	if version != expectedVersion {
		t.Fatalf("Expected version %q but got %q", expectedVersion, version)
	}

	// Now explicitly upgrade.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/db",
		Storage:   config.StorageView,
		Data:      configData("plugin_version", "1.0.0"),
	})
	assertErrIsNil(t, err)
	assertRespHasNoErr(t, resp)
	assertNoRespData(t, resp)

	version = readVersion()
	expectedVersion = "v1.0.0"
	if version != expectedVersion {
		t.Fatalf("Expected version %q but got %q", expectedVersion, version)
	}
}

func cleanup(t *testing.T, b *databaseBackend, reqs []*logical.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Go in stack order so it works similar to defer
	for i := len(reqs) - 1; i >= 0; i-- {
		req := reqs[i]
		resp, err := b.HandleRequest(ctx, req)
		if err != nil {
			t.Fatalf("Error cleaning up: %s", err)
		}
		if resp != nil && resp.IsError() {
			t.Fatalf("Error cleaning up: %s", resp.Error())
		}
	}
}

func TestBackend_PluginMain_MockV4(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args := []string{"--ca-cert=" + caPEM}

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)

	RunV4(apiClientMeta.GetTLSConfig())
}

func TestBackend_PluginMain_MockV5(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	RunV5()
}

func TestBackend_PluginMain_MockV6Multiplexed(t *testing.T) {
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	RunV6Multiplexed()
}

func assertNoRespData(t *testing.T, resp *logical.Response) {
	t.Helper()
	if resp != nil && len(resp.Data) > 0 {
		t.Fatalf("Response had data when none was expected: %#v", resp.Data)
	}
}

func assertRespHasData(t *testing.T, resp *logical.Response) {
	t.Helper()
	if resp == nil || len(resp.Data) == 0 {
		t.Fatalf("Response didn't have any data when some was expected")
	}
}

type stringAssertion func(t *testing.T, str string)

func assertStringPrefix(expectedPrefix string) stringAssertion {
	return func(t *testing.T, str string) {
		t.Helper()
		if !strings.HasPrefix(str, expectedPrefix) {
			t.Fatalf("Missing prefix %q: Actual: %q", expectedPrefix, str)
		}
	}
}

func assertStringRegex(expectedRegex string) stringAssertion {
	re := regexp.MustCompile(expectedRegex)
	return func(t *testing.T, str string) {
		if !re.MatchString(str) {
			t.Fatalf("Actual: %q did not match regexp %q", str, expectedRegex)
		}
	}
}

func assertRespHasNoErr(t *testing.T, resp *logical.Response) {
	t.Helper()
	if resp != nil && resp.IsError() {
		t.Fatalf("response is error: %#v\n", resp)
	}
}

func assertErrIsNil(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("No error expected, got: %s", err)
	}
}
