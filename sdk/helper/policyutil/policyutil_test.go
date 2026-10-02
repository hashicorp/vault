// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package policyutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestValidatePolicyName verifies policy-name canonicalization and segment-aware
// traversal validation.
func TestValidatePolicyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "accepts valid name",
			input: " Team-Read ",
			want:  "team-read",
		},
		{
			name:  "accepts embedded dot dot substring",
			input: "policy..backup",
			want:  "policy..backup",
		},
		{
			name:  "accepts embedded triple dots",
			input: "version...read",
			want:  "version...read",
		},
		{
			name:    "rejects empty canonical name",
			input:   " \t ",
			wantErr: true,
		},
		{
			name:  "accepts slash-separated name",
			input: "team/read",
			want:  "team/read",
		},
		{
			name:  "accepts backslash-containing name",
			input: `team\read`,
			want:  `team\read`,
		},
		{
			name:  "accepts dotted suffix",
			input: "read.v2",
			want:  "read.v2",
		},
		{
			name:  "accepts hidden-style segment",
			input: ".hidden",
			want:  ".hidden",
		},
		{
			name:  "accepts hidden subsegment",
			input: "team/.hidden",
			want:  "team/.hidden",
		},
		{
			name:  "accepts prefixed dotdot subsegment",
			input: "team/..hidden",
			want:  "team/..hidden",
		},
		{
			name:  "accepts suffixed dotdot subsegment",
			input: "team/hidden..",
			want:  "team/hidden..",
		},
		{
			name:    "rejects parent segment name",
			input:   "..",
			wantErr: true,
		},
		{
			name:    "rejects parent segment at beginning",
			input:   "../admin",
			wantErr: true,
		},
		{
			name:    "rejects parent segment in middle",
			input:   "team/../admin",
			wantErr: true,
		},
		{
			name:    "rejects parent segment at end",
			input:   "team/..",
			wantErr: true,
		},
		{
			name:    "rejects multi-level parent traversal",
			input:   "../root/admin",
			wantErr: true,
		},
		{
			name:    "rejects exact dot segment name",
			input:   ".",
			wantErr: true,
		},
		{
			name:    "rejects dot segment at beginning",
			input:   "./admin",
			wantErr: true,
		},
		{
			name:    "rejects dot segment in middle",
			input:   "team/./admin",
			wantErr: true,
		},
		{
			name:    "rejects dot segment at end",
			input:   "team/.",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ValidatePolicyName(tc.input)
			if tc.wantErr {
				require.Error(t, err, "expected validation error for %q", tc.input)
				return
			}

			require.NoError(t, err, "expected no validation error for %q", tc.input)
			require.Equal(t, tc.want, got, "expected canonical policy name")
		})
	}
}

// TestSanitizePolicies verifies policy sanitization preserves valid values and
// default-policy handling semantics.
func TestSanitizePolicies(t *testing.T) {
	t.Parallel()

	expected := []string{"foo", "bar"}
	actual := SanitizePolicies([]string{"foo", "bar"}, false)
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// If 'default' is already added, do not remove it.
	expected = []string{"foo", "bar", "default"}
	actual = SanitizePolicies([]string{"foo", "bar", "default"}, false)
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)
}

// TestParsePolicies verifies parsing normalizes policy lists and root-policy
// precedence semantics.
func TestParsePolicies(t *testing.T) {
	t.Parallel()

	expected := []string{"foo", "bar", "default"}
	actual := ParsePolicies("foo,bar")
	// add default if not present.
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// do not add default more than once.
	actual = ParsePolicies("foo,bar,default")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// handle spaces and tabs.
	actual = ParsePolicies(" foo ,	bar	,   default")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// ignore all others if root is present.
	expected = []string{"root"}
	actual = ParsePolicies("foo,bar,root")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// with spaces and tabs.
	expected = []string{"root"}
	actual = ParsePolicies("foo ,bar, root		")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)
}

// TestEquivalentPolicies verifies semantic policy-list equivalence, including
// case-insensitivity and implicit default handling.
func TestEquivalentPolicies(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		A        []string
		B        []string
		Expected bool
	}{
		"nil": {
			A:        nil,
			B:        nil,
			Expected: true,
		},
		"empty": {
			A:        []string{"foo", "bar"},
			B:        []string{},
			Expected: false,
		},
		"missing": {
			A:        []string{"foo", "bar"},
			B:        []string{"foo"},
			Expected: false,
		},
		"equal": {
			A:        []string{"bar", "foo"},
			B:        []string{"bar", "foo"},
			Expected: true,
		},
		"default": {
			A:        []string{"bar", "foo"},
			B:        []string{"foo", "default", "bar"},
			Expected: true,
		},
		"case-insensitive": {
			A:        []string{"test"},
			B:        []string{"Test"},
			Expected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.Expected, EquivalentPolicies(tc.A, tc.B))
		})
	}
}
