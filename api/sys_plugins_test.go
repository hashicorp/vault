package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
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

	resp, err := client.Sys().ListPluginsWithContext(context.Background(), &ListPluginsInput{})
	if err != nil {
		t.Fatal(err)
	}

	expectedPlugins := map[consts.PluginType][]string{
		consts.PluginTypeCredential: {"alicloud"},
		consts.PluginTypeDatabase:   {"cassandra-database-plugin"},
		consts.PluginTypeSecrets:    {"ad", "alicloud"},
	}

	for pluginType, expected := range expectedPlugins {
		actualPlugins := resp.PluginsByType[pluginType]
		if len(expected) != len(actualPlugins) {
			t.Fatal("Wrong number of plugins", expected, actualPlugins)
		}
		for i := range actualPlugins {
			if expected[i] != actualPlugins[i] {
				t.Fatalf("Expected %q but got %q", expected[i], actualPlugins[i])
			}
		}
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
