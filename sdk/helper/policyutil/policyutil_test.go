// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policyutil

import "testing"

func TestSanitizePolicies(t *testing.T) {
	expected := []string{"foo", "bar"}
	actual := SanitizePolicies([]string{"foo", "bar"}, false)
	if !EquivalentPolicies(expected, actual) {
		t.Fatalf("bad: expected:%s\ngot:%s\n", expected, actual)
	}

	// If 'default' is already added, do not remove it.
	expected = []string{"foo", "bar", "default"}
	actual = SanitizePolicies([]string{"foo", "bar", "default"}, false)
	if !EquivalentPolicies(expected, actual) {
		t.Fatalf("bad: expected:%s\ngot:%s\n", expected, actual)
	}
}

func TestParsePolicies(t *testing.T) {
	expected := []string{"foo", "bar", "default"}
	actual := ParsePolicies("foo,bar")
	// add default if not present.
	if !EquivalentPolicies(expected, actual) {
		t.Fatalf("bad: expected:%s\ngot:%s\n", expected, actual)
	}

	// do not add default more than once.
	actual = ParsePolicies("foo,bar,default")
	if !EquivalentPolicies(expected, actual) {
		t.Fatalf("bad: expected:%s\ngot:%s\n", expected, actual)
	}

	// handle spaces and tabs.
	actual = ParsePolicies(" foo ,	bar	,   default")
	if !EquivalentPolicies(expected, actual) {
		t.Fatalf("bad: expected:%s\ngot:%s\n", expected, actual)
	}

	// ignore all others if root is present.
	expected = []string{"root"}
	actual = ParsePolicies("foo,bar,root")
	if !EquivalentPolicies(expected, actual) {
		t.Fatalf("bad: expected:%s\ngot:%s\n", expected, actual)
	}

	// with spaces and tabs.
	expected = []string{"root"}
	actual = ParsePolicies("foo ,bar, root		")
	if !EquivalentPolicies(expected, actual) {
		t.Fatalf("bad: expected:%s\ngot:%s\n", expected, actual)
	}
}

func TestEquivalentPolicies(t *testing.T) {
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
			if EquivalentPolicies(tc.A, tc.B) != tc.Expected {
				t.Fatal("bad")
			}
		})
	}
}
