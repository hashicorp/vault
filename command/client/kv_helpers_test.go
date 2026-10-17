// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// TestAddPrefixToKVPath tests the addPrefixToKVPath helper function
func TestAddPrefixToKVPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path         string
		mountPath    string
		apiPrefix    string
		skipIfExists bool
		expected     string
	}{
		"simple": {
			path: "kv-v2/foo", mountPath: "kv-v2/", apiPrefix: "data", expected: "kv-v2/data/foo",
		},
		"multi-part": {
			path: "my/kv-v2/mount/path/foo/bar/baz", mountPath: "my/kv-v2/mount/path", apiPrefix: "metadata", expected: "my/kv-v2/mount/path/metadata/foo/bar/baz",
		},
		"with-namespace": {
			path: "my/kv-v2/mount/path/foo/bar/baz", mountPath: "my/ns1/my/kv-v2/mount/path", apiPrefix: "metadata", expected: "my/kv-v2/mount/path/metadata/foo/bar/baz",
		},
		"skip-if-exists-true": {
			path: "kv-v2/data/foo", mountPath: "kv-v2/", apiPrefix: "data", skipIfExists: true, expected: "kv-v2/data/foo",
		},
		"skip-if-exists-false": {
			path: "kv-v2/data/foo", mountPath: "kv-v2", apiPrefix: "data", expected: "kv-v2/data/data/foo",
		},
		"skip-if-exists-with-namespace": {
			path: "my/kv-v2/mount/path/metadata/foo/bar/baz", mountPath: "my/ns1/my/kv-v2/mount/path", apiPrefix: "metadata", skipIfExists: true, expected: "my/kv-v2/mount/path/metadata/foo/bar/baz",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(
				t,
				test.expected,
				addPrefixToKVPath(test.path, test.mountPath, test.apiPrefix, test.skipIfExists),
			)
		})
	}
}

// TestWalkSecretsTree tests the walkSecretsTree helper function
func TestWalkSecretsTree(t *testing.T) {
	t.Parallel()

	responses := map[string][]string{
		"kv-v1":                           {"app-1/", "foo"},
		"kv-v1/app-1":                     {"bar", "foo", "nested/"},
		"kv-v1/app-1/nested":              {"bar", "x/"},
		"kv-v1/app-1/nested/x":            {"y", "y/"},
		"kv-v1/app-1/nested/x/y":          {"z"},
		"kv-v2/metadata":                  {"app-1/", "foo"},
		"kv-v2/metadata/app-1":            {"bar", "foo", "nested/"},
		"kv-v2/metadata/app-1/nested":     {"bar", "x/"},
		"kv-v2/metadata/app-1/nested/x":   {"y", "y/"},
		"kv-v2/metadata/app-1/nested/x/y": {"z"},
	}
	type treePath struct {
		path      string
		directory bool
	}
	tests := map[string]struct {
		path          string
		expected      []treePath
		expectedError bool
	}{
		"kv-v1-simple": {
			path: "kv-v1/app-1/nested/x/y", expected: []treePath{{path: "kv-v1/app-1/nested/x/y/z"}},
		},
		"kv-v2-simple": {
			path: "kv-v2/metadata/app-1/nested/x/y", expected: []treePath{{path: "kv-v2/metadata/app-1/nested/x/y/z"}},
		},
		"kv-v1-nested": {
			path: "kv-v1/app-1/nested/", expected: []treePath{
				{path: "kv-v1/app-1/nested/bar"},
				{path: "kv-v1/app-1/nested/x", directory: true},
				{path: "kv-v1/app-1/nested/x/y"},
				{path: "kv-v1/app-1/nested/x/y", directory: true},
				{path: "kv-v1/app-1/nested/x/y/z"},
			},
		},
		"kv-v2-nested": {
			path: "kv-v2/metadata/app-1/nested/", expected: []treePath{
				{path: "kv-v2/metadata/app-1/nested/bar"},
				{path: "kv-v2/metadata/app-1/nested/x", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y"},
				{path: "kv-v2/metadata/app-1/nested/x/y", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y/z"},
			},
		},
		"kv-v1-all": {
			path: "kv-v1", expected: []treePath{
				{path: "kv-v1/app-1", directory: true},
				{path: "kv-v1/app-1/bar"},
				{path: "kv-v1/app-1/foo"},
				{path: "kv-v1/app-1/nested", directory: true},
				{path: "kv-v1/app-1/nested/bar"},
				{path: "kv-v1/app-1/nested/x", directory: true},
				{path: "kv-v1/app-1/nested/x/y"},
				{path: "kv-v1/app-1/nested/x/y", directory: true},
				{path: "kv-v1/app-1/nested/x/y/z"},
				{path: "kv-v1/foo"},
			},
		},
		"kv-v2-all": {
			path: "kv-v2/metadata", expected: []treePath{
				{path: "kv-v2/metadata/app-1", directory: true},
				{path: "kv-v2/metadata/app-1/bar"},
				{path: "kv-v2/metadata/app-1/foo"},
				{path: "kv-v2/metadata/app-1/nested", directory: true},
				{path: "kv-v2/metadata/app-1/nested/bar"},
				{path: "kv-v2/metadata/app-1/nested/x", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y"},
				{path: "kv-v2/metadata/app-1/nested/x/y", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y/z"},
				{path: "kv-v2/metadata/foo"},
			},
		},
		"kv-v1-not-found":              {path: "kv-v1/does/not/exist", expectedError: true},
		"kv-v2-not-found":              {path: "kv-v2/metadata/does/not/exist", expectedError: true},
		"kv-v1-not-listable-leaf-node": {path: "kv-v1/foo", expectedError: true},
		"kv-v2-not-listable-leaf-node": {path: "kv-v2/metadata/foo", expectedError: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					path := strings.TrimPrefix(r.URL.Path, "/v1/")
					keys, ok := responses[path]
					if !ok {
						http.NotFound(w, r)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					require.NoError(
						t,
						json.NewEncoder(w).
							Encode(map[string]interface{}{"data": map[string]interface{}{"keys": keys}}),
					)
				}),
			)
			t.Cleanup(server.Close)

			config := api.DefaultConfig()
			config.Address = server.URL
			client, err := api.NewClient(config)
			require.NoError(t, err)

			var descendants []treePath
			err = walkSecretsTree(
				t.Context(),
				client,
				test.path,
				func(path string, directory bool) error {
					descendants = append(descendants, treePath{path: path, directory: directory})
					return nil
				},
			)
			if test.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(
				t,
				reflect.DeepEqual(test.expected, descendants),
				"unexpected descendants: %#v",
				descendants,
			)
		})
	}
}

// TestPadEqualSigns verifies that table section headings receive balanced
// padding for even, odd, and shorter target widths.
func TestPadEqualSigns(t *testing.T) {
	t.Parallel()

	header := "Test Header"
	tests := []struct {
		name          string
		totalPathLen  int
		expectedCount int
	}{
		{name: "path with even length", totalPathLen: 20, expectedCount: 4},
		{name: "path with odd length", totalPathLen: 19, expectedCount: 3},
		{name: "smallest possible path", totalPathLen: 8, expectedCount: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			signs := strings.Split(padEqualSigns(header, test.totalPathLen), " "+header+" ")
			require.Len(t, signs, 2)
			require.Len(t, signs[0], len(signs[1]))
			for _, sign := range signs {
				require.Equal(t, test.expectedCount, strings.Count(sign, "="))
			}
		})
	}
}
