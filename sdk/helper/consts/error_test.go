// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package consts

import "testing"

func TestPathContainsParentReferences(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		// True cases: actual parent references
		{"bare double dot", "..", true},
		{"leading parent ref", "../foo", true},
		{"trailing parent ref", "foo/..", true},
		{"middle parent ref", "foo/../bar", true},
		{"multiple parent refs", "foo/../../bar", true},
		{"only slashes and parent", "/../", true},

		// False cases: dots that are NOT parent references
		{"triple dots", "foo/.../bar", false},
		{"triple dots trailing", "test_...", false},
		{"quadruple dots", "foo/..../bar", false},
		{"single dot", "foo/./bar", false},
		{"dots in name", "foo/bar..baz", false},
		{"double dots in name", "foo/bar..baz/qux", false},
		{"empty string", "", false},
		{"simple path", "foo/bar", false},
		{"single segment", "foo", false},
		{"trailing slash", "foo/bar/", false},
		{"dots prefix in segment", "..foo/bar", false},
		{"dots suffix in segment", "foo../bar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PathContainsParentReferences(tt.path)
			if got != tt.want {
				t.Errorf("PathContainsParentReferences(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
