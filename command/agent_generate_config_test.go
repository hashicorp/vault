// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

// TestConstructTemplates tests the construcTemplates helper function
func TestConstructTemplates(t *testing.T) {
	// test setup
	client, closer := testVaultServer(t)
	defer closer()

	// enable kv-v1 backend
	if err := client.Sys().Mount("kv-v1/", &api.MountInput{
		Type: "kv-v1",
	}); err != nil {
		t.Fatal(err)
	}

	// enable kv-v2 backend
	if err := client.Sys().Mount("kv-v2/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	ctx, cancelContextFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelContextFunc()

	// populate secrets
	for _, path := range []string{
		"foo",
		"app-1/foo",
		"app-1/bar",
		"app-1/nested/baz",
	} {
		if err := client.KVv1("kv-v1").Put(ctx, path, map[string]interface{}{
			"user":     "test",
			"password": "Hashi123",
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := client.KVv2("kv-v2").Put(ctx, path, map[string]interface{}{
			"user":     "test",
			"password": "Hashi123",
		}); err != nil {
			t.Fatal(err)
		}
	}

	// tests
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
