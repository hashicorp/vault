package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/strutil"
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
		expectedPlugins map[consts.PluginType][]string
	}{
		"no type specified": {
			input: ListPluginsInput{},
			expectedPlugins: map[consts.PluginType][]string{
				consts.PluginTypeCredential: {"alicloud"},
				consts.PluginTypeDatabase:   {"cassandra-database-plugin"},
				consts.PluginTypeSecrets:    {"ad", "alicloud"},
			},
		},
		"only auth plugins": {
			input: ListPluginsInput{Type: consts.PluginTypeCredential},
			expectedPlugins: map[consts.PluginType][]string{
				consts.PluginTypeCredential: {"alicloud"},
			},
		},
		"only database plugins": {
			input: ListPluginsInput{Type: consts.PluginTypeDatabase},
			expectedPlugins: map[consts.PluginType][]string{
				consts.PluginTypeDatabase: {"cassandra-database-plugin"},
			},
		},
		"only secret plugins": {
			input: ListPluginsInput{Type: consts.PluginTypeSecrets},
			expectedPlugins: map[consts.PluginType][]string{
				consts.PluginTypeSecrets: {"ad", "alicloud"},
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
				pluginType, err := consts.ParsePluginType(actual.Type)
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
