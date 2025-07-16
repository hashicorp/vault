// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import "testing"

func TestIsSudoPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path     string
		expected bool
	}{
		// Testing: Not a real endpoint
		{
			"/not/in/sudo/paths/list",
			false,
		},
		// Testing: sys/raw/{path}
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
		// Testing: sys/auth/{path}/tune
		{
			"/sys/auth/path/in/middle/tune",
			true,
		},
		// Testing: sys/plugins/catalog/{type} and sys/plugins/catalog/{name} (regexes overlap)
		{
			"/sys/plugins/catalog/some-type",
			true,
		},
		// Testing: Not a real endpoint
		{
			"/sys/plugins/catalog/some/type/or/name/with/slashes",
			false,
		},
		// Testing: sys/plugins/catalog/{type}/{name}
		{
			"/sys/plugins/catalog/some-type/some-name",
			true,
		},
		// Testing: Not a real endpoint
		{
			"/sys/plugins/catalog/some-type/some/name/with/slashes",
			false,
		},
		// Testing: sys/plugins/runtimes/catalog/{type}/{name}
		{
			"/sys/plugins/runtimes/catalog/some-type/some-name",
			true,
		},
		// Testing: auth/token/accessors (an example of a sudo path that only accepts list operations)
		// It is matched as sudo without the trailing slash...
		{
			"/auth/token/accessors",
			true,
		},
		// ...and also with it.
		// (Although at the time of writing, the only caller of IsSudoPath always removes trailing slashes.)
		{
			"/auth/token/accessors/",
			true,
		},
	}

	for _, tc := range testCases {
		result := IsSudoPath(tc.path)
		if result != tc.expected {
			t.Fatalf("expected api.IsSudoPath to return %v for path %s but it returned %v", tc.expected, tc.path, result)
		}
	}
}
