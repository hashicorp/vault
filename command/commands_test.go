// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	"github.com/stretchr/testify/require"
)

func Test_Commands_HCPInit(t *testing.T) {
	tests := map[string]struct {
		expectError      bool
		expectedErrorMsg string
	}{
		"initialize with success": {
			expectError: false,
		},
		"initialize with error: existing commands conflict with init commands": {
			expectError:      true,
			expectedErrorMsg: "Failed to initialize HCP commands.",
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			mockUi := cli.NewMockUi()
			commands := initCommands(mockUi, nil, nil)
			if tst.expectError {
				initHCPCommands(mockUi, commands)
				errMsg := mockUi.ErrorWriter.String()
				require.NotEmpty(t, errMsg)
				require.Contains(t, errMsg, tst.expectedErrorMsg)
			} else {
				errMsg := mockUi.ErrorWriter.String()
				require.Empty(t, errMsg)
				require.NotEmpty(t, commands)
			}
		})
	}
}
