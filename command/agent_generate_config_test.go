// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"context"
	"reflect"
	"regexp"
	"testing"
	"time"
)

// TestConstructTemplates tests the construcTemplates helper function
func TestConstructTemplates(t *testing.T) {
	ctx, cancelContextFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelContextFunc()

	client, closer := testVaultServerWithSecrets(ctx, t)
	defer closer()

	cases := map[string]struct {
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
			expectedError: false,
		},

		"kv-v2-simple": {
			paths: []string{"kv-v2/foo"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
			},
			expectedError: false,
		},

		"kv-v2-data-in-path": {
			paths: []string{"kv-v2/data/foo"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
			},
			expectedError: false,
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
			expectedError: false,
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
			expectedError: false,
		},

		"kv-v1-multi-path": {
			paths: []string{"kv-v1/foo", "kv-v1/app-1/bar"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v1/foo" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/foo" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
				{Contents: `{{ with secret "kv-v1/app-1/bar" }}{{ .Data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_PASSWORD"},
				{Contents: `{{ with secret "kv-v1/app-1/bar" }}{{ .Data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_USER"},
			},
			expectedError: false,
		},

		"kv-v2-multi-path": {
			paths: []string{"kv-v2/foo", "kv-v2/app-1/bar"},
			expected: []generatedConfigEnvTemplate{
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/foo" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "FOO_USER"},
				{Contents: `{{ with secret "kv-v2/data/app-1/bar" }}{{ .Data.data.password }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_PASSWORD"},
				{Contents: `{{ with secret "kv-v2/data/app-1/bar" }}{{ .Data.data.user }}{{ end }}`, ErrorOnMissingKey: true, Name: "BAR_USER"},
			},
			expectedError: false,
		},

		"kv-v1-path-not-found": {
			paths:         []string{"kv-v1/does/not/exist"},
			expected:      nil,
			expectedError: true,
		},

		"kv-v2-path-not-found": {
			paths:         []string{"kv-v2/does/not/exist"},
			expected:      nil,
			expectedError: true,
		},

		"kv-v1-early-wildcard": {
			paths:         []string{"kv-v1/*/foo"},
			expected:      nil,
			expectedError: true,
		},

		"kv-v2-early-wildcard": {
			paths:         []string{"kv-v2/*/foo"},
			expected:      nil,
			expectedError: true,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			templates, err := constructTemplates(ctx, client, tc.paths)

			if tc.expectedError {
				if err == nil {
					t.Fatal("an error was expected but the test succeeded")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(tc.expected, templates) {
					t.Fatalf("unexpected output; want: %v, got: %v", tc.expected, templates)
				}
			}
		})
	}
}

// TestGenerateConfiguration tests the generateConfiguration helper function
func TestGenerateConfiguration(t *testing.T) {
	ctx, cancelContextFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelContextFunc()

	client, closer := testVaultServerWithSecrets(ctx, t)
	defer closer()

	cases := map[string]struct {
		flagExec      string
		flagPaths     []string
		expected      *regexp.Regexp
		expectedError bool
	}{
		"kv-v1-simple": {
			flagExec:  "./my-app arg1 arg2",
			flagPaths: []string{"kv-v1/foo"},
			expected: regexp.MustCompile(`
auto_auth \{

  method \{
    type = "token_file"

    config \{
      token_file_path = ".*/.vault-token"
    }
  }
}

template_config \{
  static_secret_render_interval = "5m"
  exit_on_retry_failure         = true
  max_connections_per_host      = 10
}

vault \{
  address = "https://127.0.0.1:[0-9]{5}"
}

env_template "FOO_PASSWORD" \{
  contents             = "\{\{ with secret \\"kv-v1/foo\\" }}\{\{ .Data.password }}\{\{ end }}"
  error_on_missing_key = true
}
env_template "FOO_USER" \{
  contents             = "\{\{ with secret \\"kv-v1/foo\\" }}\{\{ .Data.user }}\{\{ end }}"
  error_on_missing_key = true
}

exec \{
  command                   = \["./my-app", "arg1", "arg2"\]
  restart_on_secret_changes = "always"
  restart_stop_signal       = "SIGTERM"
}
`),
			expectedError: false,
		},

		"kv-v2-default-exec": {
			flagExec:  "",
			flagPaths: []string{"kv-v2/foo"},
			expected: regexp.MustCompile(`
auto_auth \{

  method \{
    type = "token_file"

    config \{
      token_file_path = ".*/.vault-token"
    }
  }
}

template_config \{
  static_secret_render_interval = "5m"
  exit_on_retry_failure         = true
  max_connections_per_host      = 10
}

vault \{
  address = "https://127.0.0.1:[0-9]{5}"
}

env_template "FOO_PASSWORD" \{
  contents             = "\{\{ with secret \\"kv-v2/data/foo\\" }}\{\{ .Data.data.password }}\{\{ end }}"
  error_on_missing_key = true
}
env_template "FOO_USER" \{
  contents             = "\{\{ with secret \\"kv-v2/data/foo\\" }}\{\{ .Data.data.user }}\{\{ end }}"
  error_on_missing_key = true
}

exec \{
  command                   = \["env"\]
  restart_on_secret_changes = "always"
  restart_stop_signal       = "SIGTERM"
}
`),
			expectedError: false,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			var config bytes.Buffer

			c, err := generateConfiguration(ctx, client, tc.flagExec, tc.flagPaths)
			c.WriteTo(&config)

			if tc.expectedError {
				if err == nil {
					t.Fatal("an error was expected but the test succeeded")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}

				if !tc.expected.MatchString(config.String()) {
					t.Fatalf("unexpected output; want: %v, got: %v", tc.expected.String(), config.String())
				}
			}
		})
	}
}
