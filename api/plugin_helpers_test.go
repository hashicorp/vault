package api

import "testing"

func TestIsSudoPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path     string
		expected bool
	}{
		{
			"/not/in/sudo/paths/list",
			false,
		},
		{
			"/sys/raw/single-node-path",
			true,
		},
		{
			"/sys/raw/multiple/nodes/path",
			true,
		},
		{
			"/sys/raw/WEIRD(but_still_valid!)p4Th?🗿笑",
			true,
		},
		{
			"/sys/auth/path/in/middle/tune",
			true,
		},
		{
			"/sys/plugins/catalog/some-type",
			true,
		},
		{
			"/sys/plugins/catalog/some/type/or/name/with/slashes",
			false,
		},
		{
			"/sys/plugins/catalog/some-type/some-name",
			true,
		},
		{
			"/sys/plugins/catalog/some-type/some/name/with/slashes",
			false,
		},
	}

	for _, tc := range testCases {
		result := IsSudoPath(tc.path)
		if result != tc.expected {
			t.Fatalf("expected api.IsSudoPath to return %v for path %s but it returned %v", tc.expected, tc.path, result)
		}
	}
}
