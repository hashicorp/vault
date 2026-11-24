// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestConstructTemplates tests the constructTemplates helper function
func TestConstructTemplates(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/")
		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.HasPrefix(path, "sys/internal/ui/mounts/kv-v1/"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"path": "kv-v1/", "options": nil},
			})
		case strings.HasPrefix(path, "sys/internal/ui/mounts/kv-v2/"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"path": "kv-v2/", "options": map[string]interface{}{"version": "2"}},
			})
		case r.URL.Query().Get("list") == "true" && (path == "kv-v1/app-1" || path == "kv-v2/metadata/app-1"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"keys": []string{"foo", "bar", "nested/"}},
			})
		case r.URL.Query().Get("list") == "true" && (path == "kv-v1/app-1/nested" || path == "kv-v2/metadata/app-1/nested"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"keys": []string{"baz"}},
			})
		case strings.HasPrefix(path, "kv-v1/") && !strings.Contains(path, "does/not/exist"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"user": "test", "password": "Hashi123"},
			})
		case strings.HasPrefix(path, "kv-v2/data/") && !strings.Contains(path, "does/not/exist"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"data": map[string]interface{}{"user": "test", "password": "Hashi123"}},
			})
		default:
			writeAgentConfigTestResponse(t, w, http.StatusNotFound, map[string]interface{}{"errors": []string{}})
		}
	}))
	t.Cleanup(server.Close)

	config := api.DefaultConfig()
	config.Address = server.URL
	client, err := api.NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := map[string]struct {
		paths         []string
		expected      []generatedConfigEnvTemplate
		expectedError bool
	}{
		"kv-v1-simple": {
			paths: []string{"kv-v1/foo"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v1/foo" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/foo" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
			},
		},
		"kv-v2-simple": {
			paths: []string{"kv-v2/foo"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
			},
		},
		"kv-v2-data-in-path": {
			paths: []string{"kv-v2/data/foo"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
			},
		},
		"kv-v1-nested": {
			paths: []string{"kv-v1/app-1/*"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v1/app-1/bar" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/app-1/bar" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_USER"},
				{Contents: `{{ with secret "kv-v1/app-1/foo" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/app-1/foo" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
				{Contents: `{{ with secret "kv-v1/app-1/nested/baz" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAZ_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/app-1/nested/baz" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAZ_USER"},
			},
		},
		"kv-v2-nested": {
			paths: []string{"kv-v2/app-1/*"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v2/data/app-1/bar" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/app-1/bar" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_USER"},
				{Contents: `{{ with secret "kv-v2/data/app-1/foo" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/app-1/foo" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
				{Contents: `{{ with secret "kv-v2/data/app-1/nested/baz" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAZ_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/app-1/nested/baz" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAZ_USER"},
			},
		},
		"kv-v1-multi-path": {
			paths: []string{"kv-v1/foo", "kv-v1/app-1/bar"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v1/foo" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/foo" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
				{Contents: `{{ with secret "kv-v1/app-1/bar" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/app-1/bar" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_USER"},
			},
		},
		"kv-v2-multi-path": {
			paths: []string{"kv-v2/foo", "kv-v2/app-1/bar"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
				{Contents: `{{ with secret "kv-v2/data/app-1/bar" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/app-1/bar" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_USER"},
			},
		},
		"kv-v1-path-not-found": {paths: []string{"kv-v1/does/not/exist"}, expectedError: true},
		"kv-v2-path-not-found": {paths: []string{"kv-v2/does/not/exist"}, expectedError: true},
		"kv-v1-early-wildcard": {paths: []string{"kv-v1/*/foo"}, expectedError: true},
		"kv-v2-early-wildcard": {paths: []string{"kv-v2/*/foo"}, expectedError: true},
	}

	for name, test := range tests {
		templates, err := constructTemplates(ctx, client, test.paths)
		if test.expectedError {
			require.Errorf(t, err, "%s expected an error", name)
			continue
		}
		require.NoErrorf(t, err, "%s returned an error", name)
		require.Truef(t, reflect.DeepEqual(test.expected, templates), "%s: want %v, got %v", name, test.expected, templates)
	}
}

// TestAgentGenerateConfigCommand_Run verifies that the public command writes
// environment templates for KV v1 and KV v2 secrets to the requested file.
func TestAgentGenerateConfigCommand_Run(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/")
		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.HasPrefix(path, "sys/internal/ui/mounts/kv-v1/"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"path": "kv-v1/", "options": nil},
			})
		case strings.HasPrefix(path, "sys/internal/ui/mounts/kv-v2/"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"path": "kv-v2/", "options": map[string]interface{}{"version": "2"}},
			})
		case path == "kv-v1/foo":
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"user": "test", "password": "Hashi123"},
			})
		case path == "kv-v2/data/foo":
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"data": map[string]interface{}{"user": "test", "password": "Hashi123"}},
			})
		default:
			writeAgentConfigTestResponse(t, w, http.StatusNotFound, map[string]interface{}{"errors": []string{}})
		}
	}))
	t.Cleanup(server.Close)

	config := api.DefaultConfig()
	config.Address = server.URL
	client, err := api.NewClient(config)
	require.NoError(t, err)

	tests := map[string]struct {
		exec             string
		path             string
		expectedContents []string
	}{
		"kv-v1-simple": {
			exec: "./my-app arg1 arg2",
			path: "kv-v1/foo",
			expectedContents: []string{
				`env_template "FOO_PASSWORD"`,
				`secret \"kv-v1/foo\"`,
				`env_template "FOO_USER"`,
				`command                   = ["./my-app", "arg1", "arg2"]`,
			},
		},
		"kv-v2-default-exec": {
			path: "kv-v2/foo",
			expectedContents: []string{
				`env_template "FOO_PASSWORD"`,
				`secret \"kv-v2/data/foo\"`,
				`env_template "FOO_USER"`,
				`command                   = ["env"]`,
			},
		},
	}

	for name, test := range tests {
		ui := cli.NewMockUi()
		cmd := &AgentGenerateConfigCommand{BaseCommand: &base.BaseCommand{UI: ui}}
		cmd.SetClient(client)
		cmd.SetHCPTokenHelper(testHCPTokenHelper{})

		outputPath := filepath.Join(t.TempDir(), "agent.hcl")
		args := []string{"-type=env-template", "-path=" + test.path}
		if test.exec != "" {
			args = append(args, "-exec="+test.exec)
		}
		args = append(args, outputPath)

		require.Equalf(t, 0, cmd.Run(args), "%s command errors: %s", name, ui.ErrorWriter.String())
		generatedConfig, err := os.ReadFile(outputPath)
		require.NoErrorf(t, err, "%s failed to read generated config", name)
		for _, expected := range test.expectedContents {
			require.Containsf(t, string(generatedConfig), expected, "%s generated unexpected config", name)
		}
	}
}

// TestGenerateConfiguration tests the generateConfiguration helper function
func TestGenerateConfiguration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/")
		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.HasPrefix(path, "sys/internal/ui/mounts/kv-v1/"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"path": "kv-v1/", "options": nil},
			})
		case strings.HasPrefix(path, "sys/internal/ui/mounts/kv-v2/"):
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"path": "kv-v2/", "options": map[string]interface{}{"version": "2"}},
			})
		case path == "kv-v1/foo":
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"user": "test", "password": "Hashi123"},
			})
		case path == "kv-v2/data/foo":
			writeAgentConfigTestResponse(t, w, http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{"data": map[string]interface{}{"user": "test", "password": "Hashi123"}},
			})
		default:
			writeAgentConfigTestResponse(t, w, http.StatusNotFound, map[string]interface{}{"errors": []string{}})
		}
	}))
	t.Cleanup(server.Close)

	config := api.DefaultConfig()
	config.Address = server.URL
	client, err := api.NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := map[string]struct {
		exec     string
		paths    []string
		expected *regexp.Regexp
	}{
		"kv-v1-simple": {
			exec:  "./my-app arg1 arg2",
			paths: []string{"kv-v1/foo"},
			expected: regexp.MustCompile(`
auto_auth \{

  method \{
    type = "token_file"

    config \{
      token_file_path = ".*/.vault-token"
    \}
  \}
\}

template_config \{
  static_secret_render_interval = "5m"
  exit_on_retry_failure         = true
  max_connections_per_host      = 10
\}

vault \{
  address = "http://127.0.0.1:[0-9]{5}"
\}

env_template "FOO_PASSWORD" \{
  contents             = "\{\{ with secret \\"kv-v1/foo\\" }}\{\{ .Data.password }}\{\{ end }}"
  error_on_missing_key = true
\}
env_template "FOO_USER" \{
  contents             = "\{\{ with secret \\"kv-v1/foo\\" }}\{\{ .Data.user }}\{\{ end }}"
  error_on_missing_key = true
\}

exec \{
  command                   = \["./my-app", "arg1", "arg2"\]
  restart_on_secret_changes = "always"
  restart_stop_signal       = "SIGTERM"
\}
`),
		},
		"kv-v2-default-exec": {
			paths: []string{"kv-v2/foo"},
			expected: regexp.MustCompile(`
auto_auth \{

  method \{
    type = "token_file"

    config \{
      token_file_path = ".*/.vault-token"
    \}
  \}
\}

template_config \{
  static_secret_render_interval = "5m"
  exit_on_retry_failure         = true
  max_connections_per_host      = 10
\}

vault \{
  address = "http://127.0.0.1:[0-9]{5}"
\}

env_template "FOO_PASSWORD" \{
  contents             = "\{\{ with secret \\"kv-v2/data/foo\\" }}\{\{ .Data.data.password }}\{\{ end }}"
  error_on_missing_key = true
\}
env_template "FOO_USER" \{
  contents             = "\{\{ with secret \\"kv-v2/data/foo\\" }}\{\{ .Data.data.user }}\{\{ end }}"
  error_on_missing_key = true
\}

exec \{
  command                   = \["env"\]
  restart_on_secret_changes = "always"
  restart_stop_signal       = "SIGTERM"
\}
`),
		},
	}

	for name, test := range tests {
		generated, err := generateConfiguration(ctx, client, test.exec, test.paths)
		require.NoErrorf(t, err, "%s returned an error", name)

		var output bytes.Buffer
		_, err = generated.WriteTo(&output)
		require.NoErrorf(t, err, "%s failed to write generated config", name)
		require.Regexpf(t, test.expected, output.String(), "%s generated unexpected config", name)
	}
}

func writeAgentConfigTestResponse(t *testing.T, w http.ResponseWriter, status int, body interface{}) {
	t.Helper()
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Errorf("failed to write test response: %v", err)
	}
}
