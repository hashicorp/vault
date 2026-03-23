// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFileGroups_Add tests the Add method which adds groups to a FileGroups set
// while maintaining uniqueness and sorted order. Tests cover nil/empty receivers,
// single/multiple additions, duplicates, and sorting behavior.
func TestFileGroups_Add(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		groups   FileGroups
		add      []FileGroup
		expected FileGroups
	}{
		{
			name:     "add to nil",
			groups:   nil,
			add:      []FileGroup{FileGroup("docs")},
			expected: FileGroups{FileGroup("docs")},
		},
		{
			name:     "add to empty",
			groups:   FileGroups{},
			add:      []FileGroup{FileGroup("docs")},
			expected: FileGroups{FileGroup("docs")},
		},
		{
			name:     "add single group",
			groups:   FileGroups{FileGroup("docs")},
			add:      []FileGroup{FileGroup("enos")},
			expected: FileGroups{FileGroup("docs"), FileGroup("enos")},
		},
		{
			name:     "add multiple groups",
			groups:   FileGroups{FileGroup("docs")},
			add:      []FileGroup{FileGroup("enos"), FileGroup("app"), FileGroup("ui")},
			expected: FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos"), FileGroup("ui")},
		},
		{
			name:     "add duplicate group",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			add:      []FileGroup{FileGroup("docs")},
			expected: FileGroups{FileGroup("docs"), FileGroup("enos")},
		},
		{
			name:     "add multiple with duplicates",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			add:      []FileGroup{FileGroup("app"), FileGroup("docs"), FileGroup("ui"), FileGroup("enos")},
			expected: FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos"), FileGroup("ui")},
		},
		{
			name:     "add maintains sorted order",
			groups:   FileGroups{FileGroup("ui")},
			add:      []FileGroup{FileGroup("app"), FileGroup("docs")},
			expected: FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("ui")},
		},
		{
			name:     "add nothing",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			add:      []FileGroup{},
			expected: FileGroups{FileGroup("docs"), FileGroup("enos")},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.groups.Add(test.add...)
			require.Equal(t, test.expected, result, "Add result should match expected")
			// Verify result is sorted
			for i := 1; i < len(result); i++ {
				require.Less(t, string(result[i-1]), string(result[i]), "result should be sorted")
			}
		})
	}
}

// TestFileGroups_In tests the In method which performs binary search to find a group
// in the sorted FileGroups set. Returns the index where the group is (or should be)
// and a boolean indicating if it was found. Tests cover nil/empty receivers, found/not found
// cases at various positions, and proper index calculation for insertion.
func TestFileGroups_In(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name          string
		groups        FileGroups
		search        FileGroup
		expectedIndex int
		expectedFound bool
	}{
		{
			name:          "search in nil",
			groups:        nil,
			search:        FileGroup("docs"),
			expectedIndex: 0,
			expectedFound: false,
		},
		{
			name:          "search in empty",
			groups:        FileGroups{},
			search:        FileGroup("docs"),
			expectedIndex: 0,
			expectedFound: false,
		},
		{
			name:          "found in single element",
			groups:        FileGroups{FileGroup("docs")},
			search:        FileGroup("docs"),
			expectedIndex: 0,
			expectedFound: true,
		},
		{
			name:          "not found in single element",
			groups:        FileGroups{FileGroup("docs")},
			search:        FileGroup("enos"),
			expectedIndex: 1,
			expectedFound: false,
		},
		{
			name:          "found at beginning",
			groups:        FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			search:        FileGroup("app"),
			expectedIndex: 0,
			expectedFound: true,
		},
		{
			name:          "found in middle",
			groups:        FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			search:        FileGroup("docs"),
			expectedIndex: 1,
			expectedFound: true,
		},
		{
			name:          "found at end",
			groups:        FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			search:        FileGroup("enos"),
			expectedIndex: 2,
			expectedFound: true,
		},
		{
			name:          "not found - would be at beginning",
			groups:        FileGroups{FileGroup("docs"), FileGroup("enos")},
			search:        FileGroup("app"),
			expectedIndex: 0,
			expectedFound: false,
		},
		{
			name:          "not found - would be in middle",
			groups:        FileGroups{FileGroup("app"), FileGroup("enos")},
			search:        FileGroup("docs"),
			expectedIndex: 1,
			expectedFound: false,
		},
		{
			name:          "not found - would be at end",
			groups:        FileGroups{FileGroup("app"), FileGroup("docs")},
			search:        FileGroup("ui"),
			expectedIndex: 2,
			expectedFound: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			idx, found := test.groups.In(test.search)
			require.Equal(t, test.expectedIndex, idx, "index should match expected")
			require.Equal(t, test.expectedFound, found, "found status should match expected")
		})
	}
}

// TestFileGroups_All tests the All method which checks if all groups from the argument
// are present in the receiver FileGroups. Returns true only if every group in the check
// set exists in the receiver. Tests cover nil/empty cases, partial matches, full matches,
// and superset scenarios.
func TestFileGroups_All(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		groups   FileGroups
		check    FileGroups
		expected bool
	}{
		{
			name:     "nil groups, nil check",
			groups:   nil,
			check:    nil,
			expected: true,
		},
		{
			name:     "nil groups, empty check",
			groups:   nil,
			check:    FileGroups{},
			expected: true,
		},
		{
			name:     "nil groups, non-empty check",
			groups:   nil,
			check:    FileGroups{FileGroup("docs")},
			expected: false,
		},
		{
			name:     "empty groups, nil check",
			groups:   FileGroups{},
			check:    nil,
			expected: true,
		},
		{
			name:     "empty groups, empty check",
			groups:   FileGroups{},
			check:    FileGroups{},
			expected: true,
		},
		{
			name:     "empty groups, non-empty check",
			groups:   FileGroups{},
			check:    FileGroups{FileGroup("docs")},
			expected: false,
		},
		{
			name:     "non-empty groups, nil check",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    nil,
			expected: true,
		},
		{
			name:     "non-empty groups, empty check",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{},
			expected: true,
		},
		{
			name:     "all present - single element",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("docs")},
			expected: true,
		},
		{
			name:     "all present - multiple elements",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos"), FileGroup("ui")},
			check:    FileGroups{FileGroup("docs"), FileGroup("enos")},
			expected: true,
		},
		{
			name:     "all present - identical",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("docs"), FileGroup("enos")},
			expected: true,
		},
		{
			name:     "not all present - one missing",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("docs"), FileGroup("app")},
			expected: false,
		},
		{
			name:     "not all present - all missing",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("app"), FileGroup("ui")},
			expected: false,
		},
		{
			name:     "check is superset",
			groups:   FileGroups{FileGroup("docs")},
			check:    FileGroups{FileGroup("docs"), FileGroup("enos")},
			expected: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.groups.All(test.check)
			require.Equal(t, test.expected, result, "All result should match expected")
		})
	}
}

// TestFileGroups_Any tests the Any method which checks if at least one group from the
// argument is present in the receiver FileGroups. Returns true if any group in the check
// set exists in the receiver. Tests cover nil/empty cases, single/multiple matches, and
// no match scenarios.
func TestFileGroups_Any(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		groups   FileGroups
		check    FileGroups
		expected bool
	}{
		{
			name:     "nil groups, nil check",
			groups:   nil,
			check:    nil,
			expected: false,
		},
		{
			name:     "nil groups, empty check",
			groups:   nil,
			check:    FileGroups{},
			expected: false,
		},
		{
			name:     "nil groups, non-empty check",
			groups:   nil,
			check:    FileGroups{FileGroup("docs")},
			expected: false,
		},
		{
			name:     "empty groups, nil check",
			groups:   FileGroups{},
			check:    nil,
			expected: false,
		},
		{
			name:     "empty groups, empty check",
			groups:   FileGroups{},
			check:    FileGroups{},
			expected: false,
		},
		{
			name:     "empty groups, non-empty check",
			groups:   FileGroups{},
			check:    FileGroups{FileGroup("docs")},
			expected: false,
		},
		{
			name:     "non-empty groups, nil check",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    nil,
			expected: false,
		},
		{
			name:     "non-empty groups, empty check",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{},
			expected: false,
		},
		{
			name:     "one match - single element check",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("docs")},
			expected: true,
		},
		{
			name:     "one match - multiple element check",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("app"), FileGroup("docs")},
			expected: true,
		},
		{
			name:     "multiple matches",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("docs"), FileGroup("enos")},
			expected: true,
		},
		{
			name:     "all match",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("docs"), FileGroup("enos")},
			expected: true,
		},
		{
			name:     "no match - single element",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("app")},
			expected: false,
		},
		{
			name:     "no match - multiple elements",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			check:    FileGroups{FileGroup("app"), FileGroup("ui")},
			expected: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.groups.Any(test.check)
			require.Equal(t, test.expected, result, "Any result should match expected")
		})
	}
}

// TestFileGroups_Groups tests the Groups method which converts FileGroups to a slice
// of strings. Each FileGroup is converted to its string representation. Tests cover
// nil/empty receivers and various group counts.
func TestFileGroups_Groups(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		groups   FileGroups
		expected []string
	}{
		{
			name:     "nil groups",
			groups:   nil,
			expected: []string{},
		},
		{
			name:     "empty groups",
			groups:   FileGroups{},
			expected: []string{},
		},
		{
			name:     "single group",
			groups:   FileGroups{FileGroup("docs")},
			expected: []string{"docs"},
		},
		{
			name:     "multiple groups",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			expected: []string{"app", "docs", "enos"},
		},
		{
			name:     "all file groups",
			groups:   FileGroups{FileGroup("autopilot"), FileGroup("changelog"), FileGroup("community"), FileGroup("docs"), FileGroup("enos")},
			expected: []string{"autopilot", "changelog", "community", "docs", "enos"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.groups.Groups()
			require.Equal(t, test.expected, result, "Groups result should match expected")
		})
	}
}

// TestFileGroups_String tests the String method which returns a comma-separated string
// representation of all groups in the FileGroups set. Tests cover nil/empty receivers
// and various group counts.
func TestFileGroups_String(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name     string
		groups   FileGroups
		expected string
	}{
		{
			name:     "nil groups",
			groups:   nil,
			expected: "",
		},
		{
			name:     "empty groups",
			groups:   FileGroups{},
			expected: "",
		},
		{
			name:     "single group",
			groups:   FileGroups{FileGroup("docs")},
			expected: "docs",
		},
		{
			name:     "multiple groups",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			expected: "app, docs, enos",
		},
		{
			name:     "all file groups",
			groups:   FileGroups{FileGroup("autopilot"), FileGroup("changelog"), FileGroup("community"), FileGroup("docs"), FileGroup("enos")},
			expected: "autopilot, changelog, community, docs, enos",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.groups.String()
			require.Equal(t, test.expected, result, "String result should match expected")
		})
	}
}

// TestFileGroups_Intersection tests the Intersection method which returns a new FileGroups
// containing only the groups present in both the receiver and the argument. Returns an empty
// FileGroups if there's no intersection. Tests cover nil/empty cases, no intersection, partial
// intersections, full intersections, and various subset scenarios.
func TestFileGroups_Intersection(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		groups   FileGroups
		arg      FileGroups
		expected FileGroups
	}{
		{
			name:     "both nil",
			groups:   nil,
			arg:      nil,
			expected: FileGroups{},
		},
		{
			name:     "groups nil, arg has values",
			groups:   nil,
			arg:      FileGroups{FileGroup("docs"), FileGroup("enos")},
			expected: FileGroups{},
		},
		{
			name:     "groups has values, arg nil",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			arg:      nil,
			expected: FileGroups{},
		},
		{
			name:     "both empty",
			groups:   FileGroups{},
			arg:      FileGroups{},
			expected: FileGroups{},
		},
		{
			name:     "groups empty, arg has values",
			groups:   FileGroups{},
			arg:      FileGroups{FileGroup("docs"), FileGroup("enos")},
			expected: FileGroups{},
		},
		{
			name:     "groups has values, arg empty",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			arg:      FileGroups{},
			expected: FileGroups{},
		},
		{
			name:     "no intersection",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos")},
			arg:      FileGroups{FileGroup("app"), FileGroup("ui")},
			expected: FileGroups{},
		},
		{
			name:     "partial intersection - one common element",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			arg:      FileGroups{FileGroup("app"), FileGroup("ui")},
			expected: FileGroups{FileGroup("app")},
		},
		{
			name:     "partial intersection - multiple common elements",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos"), FileGroup("ui")},
			arg:      FileGroups{FileGroup("app"), FileGroup("enos"), FileGroup("tools"), FileGroup("ui")},
			expected: FileGroups{FileGroup("app"), FileGroup("enos"), FileGroup("ui")},
		},
		{
			name:     "full intersection - identical groups",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			arg:      FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
			expected: FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos")},
		},
		{
			name:     "full intersection - arg is subset of groups",
			groups:   FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos"), FileGroup("ui")},
			arg:      FileGroups{FileGroup("app"), FileGroup("enos")},
			expected: FileGroups{FileGroup("app"), FileGroup("enos")},
		},
		{
			name:     "full intersection - groups is subset of arg",
			groups:   FileGroups{FileGroup("app"), FileGroup("enos")},
			arg:      FileGroups{FileGroup("app"), FileGroup("docs"), FileGroup("enos"), FileGroup("ui")},
			expected: FileGroups{FileGroup("app"), FileGroup("enos")},
		},
		{
			name:     "single element in both - match",
			groups:   FileGroups{FileGroup("docs")},
			arg:      FileGroups{FileGroup("docs")},
			expected: FileGroups{FileGroup("docs")},
		},
		{
			name:     "single element in both - no match",
			groups:   FileGroups{FileGroup("docs")},
			arg:      FileGroups{FileGroup("enos")},
			expected: FileGroups{},
		},
		{
			name:     "single element in groups, multiple in arg - match",
			groups:   FileGroups{FileGroup("docs")},
			arg:      FileGroups{FileGroup("docs"), FileGroup("enos"), FileGroup("app")},
			expected: FileGroups{FileGroup("docs")},
		},
		{
			name:     "single element in groups, multiple in arg - no match",
			groups:   FileGroups{FileGroup("ui")},
			arg:      FileGroups{FileGroup("docs"), FileGroup("enos"), FileGroup("app")},
			expected: FileGroups{},
		},
		{
			name:     "multiple in groups, single in arg - match",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos"), FileGroup("app")},
			arg:      FileGroups{FileGroup("enos")},
			expected: FileGroups{FileGroup("enos")},
		},
		{
			name:     "multiple in groups, single in arg - no match",
			groups:   FileGroups{FileGroup("docs"), FileGroup("enos"), FileGroup("app")},
			arg:      FileGroups{FileGroup("ui")},
			expected: FileGroups{},
		},
		{
			name:     "all file groups intersection",
			groups:   FileGroups{FileGroup("autopilot"), FileGroup("changelog"), FileGroup("community"), FileGroup("docs"), FileGroup("enos"), FileGroup("enterprise")},
			arg:      FileGroups{FileGroup("docs"), FileGroup("enos"), FileGroup("enterprise"), FileGroup("github"), FileGroup("app"), FileGroup("gotoolchain")},
			expected: FileGroups{FileGroup("docs"), FileGroup("enos"), FileGroup("enterprise")},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.groups.Intersection(test.arg)
			require.Equal(t, test.expected, result, "intersection result should match expected")
		})
	}
}
