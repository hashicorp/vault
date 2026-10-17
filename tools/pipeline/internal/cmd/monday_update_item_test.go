// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestMondayUpdateItemCmd_FlagParsing verifies that required flags are enforced,
// that --board-id, --item-id, --column-id, and --value are all populated, and
// that --github-output is correctly parsed when set.
func TestMondayUpdateItemCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args             []string
		wantBoardID      string
		wantItemID       string
		wantColumnID     string
		wantValue        string
		wantGithubOutput bool
		wantErr          bool
	}{
		"required --board-id missing returns error": {
			args:    []string{"--item-id", "1210003", "--column-id", "2020333", "--value", "3333"},
			wantErr: true,
		},
		"required --item-id missing returns error": {
			args:    []string{"--board-id", "1210003", "--column-id", "2020333", "--value", "3333"},
			wantErr: true,
		},
		"required --column-id missing returns error": {
			args:    []string{"--board-id", "1210003", "--item-id", "2020333", "--value", "3333"},
			wantErr: true,
		},
		"required --value missing returns error": {
			args:    []string{"--board-id", "1210003", "--item-id", "2020333", "--column-id", "3333"},
			wantErr: true,
		},
		"explicit all required flags set": {
			args:         []string{"--board-id", "1010222", "--item-id", "2020333", "--column-id", "3030444", "--value", "4040555"},
			wantBoardID:  "1010222",
			wantItemID:   "2020333",
			wantColumnID: "3030444",
			wantValue:    "4040555",
		},
		"--github-output flag is set": {
			args:             []string{"--board-id", "1010222", "--item-id", "2020333", "--column-id", "3030444", "--value", "4040555", "--github-output"},
			wantBoardID:      "1010222",
			wantItemID:       "2020333",
			wantColumnID:     "3030444",
			wantValue:        "4040555",
			wantGithubOutput: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayUpdateItemCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Inspect the flags directly from the command rather than a
			// package-level req variable - update-item uses a local req.
			boardID, err := cmd.Flags().GetString("board-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantBoardID, boardID)

			itemID, err := cmd.Flags().GetString("item-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantItemID, itemID)

			columnID, err := cmd.Flags().GetString("column-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantColumnID, columnID)

			value, err := cmd.Flags().GetString("value")
			require.NoError(t, err)
			require.Equal(t, tc.wantValue, value)

			githubOutput, err := cmd.Flags().GetBool("github-output")
			require.NoError(t, err)
			require.Equal(t, tc.wantGithubOutput, githubOutput)
		})
	}
}

// TestMondayUpdateItemRes_ToTable verifies that ToTable renders the
// item ID, item name and new value in the output.
func TestMondayUpdateItemRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.UpdateItemRes{
		ItemID:   "2020333",
		ItemName: "Monday Item",
		NewValue: "Completed",
	}

	output := res.ToTable().Render()
	require.Contains(t, output, "2020333")
	require.Contains(t, output, "Monday Item")
	require.Contains(t, output, "Completed")
}

// TestMondayUpdateItemRes_ToJSON verifies that ToJSON produces valid JSON
// containing the expected item_id, item_name and new_value
func TestMondayUpdateItemRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.UpdateItemRes{
		ItemID:   "2020333",
		ItemName: "Monday Item",
		NewValue: "Completed",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, "2020333", got["item_id"])
	require.Equal(t, "Monday Item", got["item_name"])
	require.Equal(t, "Completed", got["new_value"])
}
