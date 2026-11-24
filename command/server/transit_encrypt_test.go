// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestOpenInput_Stdin verifies that openInput reads the input "-" from the
// command's stdin, which RunCustom callers set with RunOptions.Stdin, rather
// than from the process's stdin.
func TestOpenInput_Stdin(t *testing.T) {
	t.Parallel()

	stdin := strings.NewReader("plaintext")
	defers, in, length, err := openInput("-", base.NewBaseCommand(nil, &base.RunOptions{Stdin: stdin}))
	require.NoError(t, err)
	require.Nil(t, defers)
	require.Nil(t, length)
	require.Same(t, stdin, in)
}
