// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/require"
)

func Test_Commands_HCPInit(t *testing.T) {
	hcpCommands := hcpvlib.InitHCPCommand(cli.NewMockUi())

	tests := map[string]struct {
		expectError      bool
		expectedErrorMsg string
		commands         map[string]cli.CommandFactory
	}{
		"initialize with success": {
			expectError: false,
			commands:    map[string]cli.CommandFactory{},
		},
		"initialize with error: existing commands conflict with init commands": {
			expectError:      true,
			expectedErrorMsg: "Failed to initialize HCP commands.",
			commands:         hcpCommands,
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			mockUi := cli.NewMockUi()
			initHCPCommands(mockUi, tst.commands)

			errMsg := mockUi.ErrorWriter.String()
			if tst.expectError {
				require.NotEmpty(t, errMsg)
				require.Contains(t, errMsg, tst.expectedErrorMsg)
			} else {
				require.Empty(t, errMsg)
				require.NotEmpty(t, tst.commands)
			}
		})
	}
}
