// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	"github.com/stretchr/testify/require"
)

// TestInitClientCommands_AuthMethodHelp verifies that the registered userpass and
// AWS handlers expose credential-specific usage without contacting either service.
func TestInitClientCommands_AuthMethodHelp(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		method string
		hint   string
	}{
		{method: "userpass", hint: "username=<string>"},
		{method: "aws", hint: "aws_access_key_id=<string>"},
	} {
		t.Run(tc.method, func(t *testing.T) {
			t.Parallel()

			ui := cli.NewMockUi()
			commands := InitClientCommands(ui, ui, nil)
			cmd, err := commands["auth help"]()
			require.NoError(t, err, "auth help factory should construct a command")
			handler := cmd.(*AuthHelpCommand).Handlers[tc.method]
			require.NotNil(t, handler, "auth help must register the reported method")
			help := handler.Help()
			require.Contains(t, help, "login -method="+tc.method, "usage must identify the selected method")
			require.Contains(t, help, tc.hint, "usage must describe the method's credentials")
		})
	}
}
