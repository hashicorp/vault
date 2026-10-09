// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import "testing"

// TestPredict_HasPathArg verifies that prediction stops after a command has
// received more than one positional argument.
func TestPredict_HasPathArg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		exp  bool
	}{
		{name: "nil", args: nil, exp: false},
		{name: "empty", args: []string{}, exp: false},
		{name: "empty_string", args: []string{""}, exp: false},
		{name: "single", args: []string{"foo"}, exp: false},
		{name: "multiple", args: []string{"foo", "bar", "baz"}, exp: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if actual := NewPredict().hasPathArg(test.args); actual != test.exp {
				t.Errorf("expected %t to be %t", actual, test.exp)
			}
		})
	}
}
