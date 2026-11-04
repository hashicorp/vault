// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestMondayCreateItemCmd_FlagParsing verifies that required flags are enforced,
// that --board-id and --item-name are populated when provided, that the optional
// --group-id is accepted and forwarded, and that --github-output is correctly
// parsed when set.
//
// This test covers the CLI layer that CreateItemReq.Run tests cannot reach:
// flag registration, MarkFlagRequired enforcement, and optional flag defaults.
// It mirrors the isolated command tests in monday_create_board_test.go and
// monday_update_item_test.go to catch CLI regressions independently of the
// request-level behaviour.
func TestMondayCreateItemCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args             []string
		wantBoardID      string
		wantItemName     string
		wantGroupID      string
		wantGithubOutput bool
		wantErr          bool
	}{
		// Required flags — missing either must be rejected by cobra before RunE runs.
		"required --board-id missing returns error": {
			args:    []string{"--item-name", "Release Checklist"},
			wantErr: true,
		},
		"required --item-name missing returns error": {
			args:    []string{"--board-id", "12345"},
			wantErr: true,
		},

		// Both required flags present; optional --group-id is absent.
		// This is the common production call-site: place item in default top group.
		"both required flags set, --group-id defaults to empty string": {
			args:         []string{"--board-id", "12345", "--item-name", "Release Checklist"},
			wantBoardID:  "12345",
			wantItemName: "Release Checklist",
			wantGroupID:  "",
		},

		// Explicit --group-id must be forwarded so the item lands in the right group.
		"explicit --group-id is respected": {
			args:         []string{"--board-id", "12345", "--item-name", "Release Checklist", "--group-id", "topics"},
			wantBoardID:  "12345",
			wantItemName: "Release Checklist",
			wantGroupID:  "topics",
		},

		// --github-output writes the new item-id to $GITHUB_OUTPUT; must be off by default.
		"--github-output flag is set": {
			args:             []string{"--board-id", "12345", "--item-name", "Release Checklist", "--github-output"},
			wantBoardID:      "12345",
			wantItemName:     "Release Checklist",
			wantGithubOutput: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayCreateItemCmd()
			// Replace RunE with a no-op so the test exercises only flag parsing
			// without requiring a live Monday.com client.  This matches the
			// pattern used by monday_create_board_test.go and
			// monday_update_item_test.go.
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Inspect flags directly from the cobra.Command; the create-item
			// command stores state in a local req struct, so reading the flags
			// is the only portable way to assert what was parsed.
			boardID, err := cmd.Flags().GetString("board-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantBoardID, boardID)

			itemName, err := cmd.Flags().GetString("item-name")
			require.NoError(t, err)
			require.Equal(t, tc.wantItemName, itemName)

			groupID, err := cmd.Flags().GetString("group-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantGroupID, groupID)

			githubOutput, err := cmd.Flags().GetBool("github-output")
			require.NoError(t, err)
			require.Equal(t, tc.wantGithubOutput, githubOutput)
		})
	}
}
