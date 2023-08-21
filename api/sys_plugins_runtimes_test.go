package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegisterPluginRuntime(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerRegister))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().RegisterPluginRuntimeWithContext(context.Background(), &RegisterPluginRuntimeInput{
		Name:         "gvisor",
		Type:         PluginRuntimeTypeContainer,
		OCIRuntime:   "runsc",
		ParentCGroup: "/cpulimit-cgroup/",
		CPU:          1,
		Memory:       10000,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPluginRuntime(t *testing.T) {
	for name, tc := range map[string]struct {
		body     string
		expected GetPluginRuntimeResponse
	}{
		"gvisor": {
			body: getPluginRuntimeResponse,
			expected: GetPluginRuntimeResponse{
				Name:         "gvisor",
				Type:         PluginRuntimeTypeContainer.String(),
				OCIRuntime:   "runsc",
				ParentCGroup: "/cpulimit-cgroup/",
				CPU:          1,
				Memory:       10000,
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

			input := GetPluginRuntimeInput{
				Name: "gvisor",
				Type: PluginRuntimeTypeContainer,
			}

			info, err := client.Sys().GetPluginRuntimeWithContext(context.Background(), &input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expected, *info) {
				t.Errorf("expected: %#v\ngot: %#v", tc.expected, info)
			}
		})
	}
}

const getPluginRuntimeResponse = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c123456",
    "data": {
        "name": "gvisor",
        "type": "container",
        "oci_runtime": "runsc",
        "parent_cgroup": "/cpulimit-cgroup/",
        "cpu": 1,
        "memory": 10000
    },
    "warnings": null,
    "auth": null
}`
