package command

import "testing"

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
