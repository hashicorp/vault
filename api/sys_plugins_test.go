// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/go-secure-stdlib/strutil"
)

func TestRegisterPlugin(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerRegister))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().RegisterPluginWithContext(context.Background(), &RegisterPluginInput{
		Version: "v1.0.0",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListPlugins(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerList))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for name, tc := range map[string]struct {
		input           ListPluginsInput
		expectedPlugins map[PluginType][]string
	}{
		"no type specified": {
			input: ListPluginsInput{},
			expectedPlugins: map[PluginType][]string{
				PluginTypeCredential: {"alicloud"},
				PluginTypeDatabase:   {"cassandra-database-plugin"},
				PluginTypeSecrets:    {"ad", "alicloud"},
			},
		},
		"only auth plugins": {
			input: ListPluginsInput{Type: PluginTypeCredential},
			expectedPlugins: map[PluginType][]string{
				PluginTypeCredential: {"alicloud"},
			},
		},
		"only database plugins": {
			input: ListPluginsInput{Type: PluginTypeDatabase},
			expectedPlugins: map[PluginType][]string{
				PluginTypeDatabase: {"cassandra-database-plugin"},
			},
		},
		"only secret plugins": {
			input: ListPluginsInput{Type: PluginTypeSecrets},
			expectedPlugins: map[PluginType][]string{
				PluginTypeSecrets: {"ad", "alicloud"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			resp, err := client.Sys().ListPluginsWithContext(context.Background(), &tc.input)
			if err != nil {
				t.Fatal(err)
			}

			for pluginType, expected := range tc.expectedPlugins {
				actualPlugins := resp.PluginsByType[pluginType]
				if len(expected) != len(actualPlugins) {
					t.Fatal("Wrong number of plugins", expected, actualPlugins)
				}
				for i := range actualPlugins {
					if expected[i] != actualPlugins[i] {
						t.Fatalf("Expected %q but got %q", expected[i], actualPlugins[i])
					}
				}

				for _, expectedPlugin := range expected {
					found := false
					for _, plugin := range resp.Details {
						if plugin.Type == pluginType.String() && plugin.Name == expectedPlugin {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected to find %s plugin %s but not found in details: %#v", pluginType.String(), expectedPlugin, resp.Details)
					}
				}
			}

			for _, actual := range resp.Details {
				pluginType, err := ParsePluginType(actual.Type)
				if err != nil {
					t.Fatal(err)
				}
				if !strutil.StrListContains(tc.expectedPlugins[pluginType], actual.Name) {
					t.Errorf("Did not expect to find %s in details", actual.Name)
				}
			}
		})
	}
}

func TestGetPlugin(t *testing.T) {
	for name, tc := range map[string]struct {
		version  string
		body     string
		expected GetPluginResponse
	}{
		"builtin": {
			body: getResponse,
			expected: GetPluginResponse{
				Args:              nil,
				Builtin:           true,
				Command:           "",
				Name:              "azure",
				SHA256:            "",
				DeprecationStatus: "supported",
				Version:           "v0.14.0+builtin",
			},
		},
		"external": {
			version: "v1.0.0",
			body:    getResponseExternal,
			expected: GetPluginResponse{
				Args:              []string{},
				Builtin:           false,
				Command:           "azure-plugin",
				Name:              "azure",
				SHA256:            "8ba442dba253803685b05e35ad29dcdebc48dec16774614aa7a4ebe53c1e90e1",
				DeprecationStatus: "",
				Version:           "v1.0.0",
			},
		},
		"old server": {
			body: getResponseOldServerVersion,
			expected: GetPluginResponse{
				Args:              nil,
				Builtin:           true,
				Command:           "",
				Name:              "azure",
				SHA256:            "",
				DeprecationStatus: "",
				Version:           "",
			},
		},
		"oci image": {
			version: "v0.16.0",
			body:    getResponseOCIImageVersion,
			expected: GetPluginResponse{
				Args:              []string{},
				Builtin:           false,
				Command:           "",
				Name:              "jwt",
				OCIImage:          "hashicorp/vault-plugin-auth-jwt",
				Runtime:           "gvisor",
				SHA256:            "8ba442dba253803685b05e35ad29dcdebc48dec16774614aa7a4ebe53c1e90e1",
				DeprecationStatus: "",
				Version:           "v0.16.0",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerInfo(tc.body)))
			defer mockVaultServer.Close()

			cfg := DefaultConfig()
			cfg.Address = mockVaultServer.URL
			client, err := NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}

			input := GetPluginInput{
				Name: "azure",
				Type: PluginTypeSecrets,
			}
			if tc.version != "" {
				input.Version = tc.version
			}

			info, err := client.Sys().GetPluginWithContext(context.Background(), &input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expected, *info) {
				t.Errorf("expected: %#v\ngot: %#v", tc.expected, info)
			}
		})
	}
}

func mockVaultHandlerInfo(body string) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	}
}

const getResponse = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c495241",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "args": null,
        "builtin": true,
        "command": "",
        "deprecation_status": "supported",
        "name": "azure",
        "sha256": "",
        "version": "v0.14.0+builtin"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`

const getResponseExternal = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c495241",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "args": [],
        "builtin": false,
        "command": "azure-plugin",
        "name": "azure",
        "sha256": "8ba442dba253803685b05e35ad29dcdebc48dec16774614aa7a4ebe53c1e90e1",
        "version": "v1.0.0"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`

const getResponseOldServerVersion = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c495241",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "args": null,
        "builtin": true,
        "command": "",
        "name": "azure",
        "sha256": ""
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`

const getResponseOCIImageVersion = `{
	"request_id": "e93d3f93-8e4f-8443-a803-f1c97c495241",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "args": [],
        "builtin": false,
        "name": "jwt",
		"oci_image" : "hashicorp/vault-plugin-auth-jwt",
		"runtime" : "gvisor",
        "sha256": "8ba442dba253803685b05e35ad29dcdebc48dec16774614aa7a4ebe53c1e90e1",
        "version": "v0.16.0"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`

func mockVaultHandlerList(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(listUntypedResponse))
}

const listUntypedResponse = `{
  "request_id": "82601a91-cd7a-718f-feca-f573449cc1bb",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "auth": [
      "alicloud"
    ],
    "database": [
      "cassandra-database-plugin"
    ],
    "secret": [
      "ad",
      "alicloud"
    ],
    "some_other_unexpected_key": [
      {
        "objectKey": "objectValue"
      },
      {
        "arbitraryData": 7
      }
    ],
    "detailed": [
      {
        "type": "auth",
        "name": "alicloud",
        "version": "v0.13.0+builtin",
        "builtin": true,
        "deprecation_status": "supported"
      },
      {
        "type": "database",
        "name": "cassandra-database-plugin",
        "version": "v1.13.0+builtin.vault",
        "builtin": true,
        "deprecation_status": "supported"
      },
      {
        "type": "secret",
        "name": "ad",
        "version": "v0.14.0+builtin",
        "builtin": true,
        "deprecation_status": "supported"
      },
      {
        "type": "secret",
        "name": "alicloud",
        "version": "v0.13.0+builtin",
        "builtin": true,
        "deprecation_status": "supported"
      }
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`

func mockVaultHandlerRegister(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(registerResponse))
}

const registerResponse = `{}`
