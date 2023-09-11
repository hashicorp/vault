// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

// TestAddPrefixToKVPath tests the addPrefixToKVPath helper function
func TestAddPrefixToKVPath(t *testing.T) {
	cases := map[string]struct {
		path         string
		mountPath    string
		apiPrefix    string
		skipIfExists bool
		expected     string
	}{
		"simple": {
			path:         "kv-v2/foo",
			mountPath:    "kv-v2/",
			apiPrefix:    "data",
			skipIfExists: false,
			expected:     "kv-v2/data/foo",
		},

		"multi-part": {
			path:         "my/kv-v2/mount/path/foo/bar/baz",
			mountPath:    "my/kv-v2/mount/path",
			apiPrefix:    "metadata",
			skipIfExists: false,
			expected:     "my/kv-v2/mount/path/metadata/foo/bar/baz",
		},

		"with-namespace": {
			path:         "my/kv-v2/mount/path/foo/bar/baz",
			mountPath:    "my/ns1/my/kv-v2/mount/path",
			apiPrefix:    "metadata",
			skipIfExists: false,
			expected:     "my/kv-v2/mount/path/metadata/foo/bar/baz",
		},

		"skip-if-exists-true": {
			path:         "kv-v2/data/foo",
			mountPath:    "kv-v2/",
			apiPrefix:    "data",
			skipIfExists: true,
			expected:     "kv-v2/data/foo",
		},

		"skip-if-exists-false": {
			path:         "kv-v2/data/foo",
			mountPath:    "kv-v2",
			apiPrefix:    "data",
			skipIfExists: false,
			expected:     "kv-v2/data/data/foo",
		},

		"skip-if-exists-with-namespace": {
			path:         "my/kv-v2/mount/path/metadata/foo/bar/baz",
			mountPath:    "my/ns1/my/kv-v2/mount/path",
			apiPrefix:    "metadata",
			skipIfExists: true,
			expected:     "my/kv-v2/mount/path/metadata/foo/bar/baz",
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := addPrefixToKVPath(
				tc.path,
				tc.mountPath,
				tc.apiPrefix,
				tc.skipIfExists,
			)

			if tc.expected != actual {
				t.Fatalf("unexpected output; want: %v, got: %v", tc.expected, actual)
			}
		})
	}
}

// TestWalkSecretsTree tests the walkSecretsTree helper function
func TestWalkSecretsTree(t *testing.T) {
	// test setup
	client, closer := testVaultServer(t)
	defer closer()

	// enable kv-v1 backend
	if err := client.Sys().Mount("kv-v1/", &api.MountInput{
		Type: "kv-v1",
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

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
		"app-1/nested/x/y/z",
		"app-1/nested/x/y",
		"app-1/nested/bar",
	} {
		if err := client.KVv1("kv-v1").Put(ctx, path, map[string]interface{}{
			"password": "Hashi123",
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := client.KVv2("kv-v2").Put(ctx, path, map[string]interface{}{
			"password": "Hashi123",
		}); err != nil {
			t.Fatal(err)
		}
	}

	type treePath struct {
		path      string
		directory bool
	}

	cases := map[string]struct {
		path          string
		expected      []treePath
		expectedError bool
	}{
		"kv-v1-simple": {
			path: "kv-v1/app-1/nested/x/y",
			expected: []treePath{
				{path: "kv-v1/app-1/nested/x/y/z", directory: false},
			},
			expectedError: false,
		},

		"kv-v2-simple": {
			path: "kv-v2/metadata/app-1/nested/x/y",
			expected: []treePath{
				{path: "kv-v2/metadata/app-1/nested/x/y/z", directory: false},
			},
			expectedError: false,
		},

		"kv-v1-nested": {
			path: "kv-v1/app-1/nested/",
			expected: []treePath{
				{path: "kv-v1/app-1/nested/bar", directory: false},
				{path: "kv-v1/app-1/nested/x", directory: true},
				{path: "kv-v1/app-1/nested/x/y", directory: false},
				{path: "kv-v1/app-1/nested/x/y", directory: true},
				{path: "kv-v1/app-1/nested/x/y/z", directory: false},
			},
			expectedError: false,
		},

		"kv-v2-nested": {
			path: "kv-v2/metadata/app-1/nested/",
			expected: []treePath{
				{path: "kv-v2/metadata/app-1/nested/bar", directory: false},
				{path: "kv-v2/metadata/app-1/nested/x", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y", directory: false},
				{path: "kv-v2/metadata/app-1/nested/x/y", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y/z", directory: false},
			},
			expectedError: false,
		},

		"kv-v1-all": {
			path: "kv-v1",
			expected: []treePath{
				{path: "kv-v1/app-1", directory: true},
				{path: "kv-v1/app-1/bar", directory: false},
				{path: "kv-v1/app-1/foo", directory: false},
				{path: "kv-v1/app-1/nested", directory: true},
				{path: "kv-v1/app-1/nested/bar", directory: false},
				{path: "kv-v1/app-1/nested/x", directory: true},
				{path: "kv-v1/app-1/nested/x/y", directory: false},
				{path: "kv-v1/app-1/nested/x/y", directory: true},
				{path: "kv-v1/app-1/nested/x/y/z", directory: false},
				{path: "kv-v1/foo", directory: false},
			},
			expectedError: false,
		},

		"kv-v2-all": {
			path: "kv-v2/metadata",
			expected: []treePath{
				{path: "kv-v2/metadata/app-1", directory: true},
				{path: "kv-v2/metadata/app-1/bar", directory: false},
				{path: "kv-v2/metadata/app-1/foo", directory: false},
				{path: "kv-v2/metadata/app-1/nested", directory: true},
				{path: "kv-v2/metadata/app-1/nested/bar", directory: false},
				{path: "kv-v2/metadata/app-1/nested/x", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y", directory: false},
				{path: "kv-v2/metadata/app-1/nested/x/y", directory: true},
				{path: "kv-v2/metadata/app-1/nested/x/y/z", directory: false},
				{path: "kv-v2/metadata/foo", directory: false},
			},
			expectedError: false,
		},

		"kv-v1-not-found": {
			path:          "kv-v1/does/not/exist",
			expected:      nil,
			expectedError: true,
		},

		"kv-v2-not-found": {
			path:          "kv-v2/metadata/does/not/exist",
			expected:      nil,
			expectedError: true,
		},

		"kv-v1-not-listable-leaf-node": {
			path:          "kv-v1/foo",
			expected:      nil,
			expectedError: true,
		},

		"kv-v2-not-listable-leaf-node": {
			path:          "kv-v2/metadata/foo",
			expected:      nil,
			expectedError: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var descendants []treePath

			err := walkSecretsTree(ctx, client, tc.path, func(path string, directory bool) error {
				descendants = append(descendants, treePath{
					path:      path,
					directory: directory,
				})
				return nil
			})

			if tc.expectedError {
				if err == nil {
					t.Fatal("an error was expected but the test succeeded")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(tc.expected, descendants) {
					t.Fatalf("unexpected list output; want: %v, got: %v", tc.expected, descendants)
				}
			}
		})
	}
}
