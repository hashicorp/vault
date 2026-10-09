// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRunGo verifies that go runs with GOWORK=off so modules build with their
// own go.mod, and that extra environment entries override the process
// environment.
func TestRunGo(t *testing.T) {
	t.Parallel()

	res, err := runGo(t.Context(), t.TempDir(), []string{"GOFLAGS=-tags=fixture"}, "env", "GOWORK", "GOFLAGS")
	require.NoError(t, err)
	require.Equal(t, []string{"off", "-tags=fixture"}, strings.Fields(string(res.Stdout)))
}

// TestRunGo_Error verifies that a failing go command returns an error along
// with its stderr, which callers report to the user.
func TestRunGo_Error(t *testing.T) {
	t.Parallel()

	res, err := runGo(t.Context(), t.TempDir(), nil, "not-a-go-command")
	require.Error(t, err)
	require.Contains(t, string(res.Stderr), "unknown command")
}
