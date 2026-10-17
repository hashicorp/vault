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

// TestMondayGetItemIDCmd_FlagParsing verifies that required flags are enforced,
// that --board-id and --item-name are correctly populated when provided, and
// that --github-output is correctly parsed when set.
func TestMondayGetItemIDCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args             []string
		wantBoardID      string
		wantItemName     string
		wantGithubOutput bool
		wantErr          bool
	}{
		"required --board-id missing returns error": {
			args:    []string{"--item-name", "Release Checklist"},
			wantErr: true,
		},
		"required --item-name missing returns error": {
			args:    []string{"--board-id", "1010222"},
			wantErr: true,
		},
		"all required flags set": {
			args:         []string{"--board-id", "1010222", "--item-name", "Release Checklist"},
			wantBoardID:  "1010222",
			wantItemName: "Release Checklist",
		},
		"--github-output flag is set": {
			args:             []string{"--board-id", "1010222", "--item-name", "Release Checklist", "--github-output"},
			wantBoardID:      "1010222",
			wantItemName:     "Release Checklist",
			wantGithubOutput: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayGetItemIDCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			boardID, err := cmd.Flags().GetString("board-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantBoardID, boardID)

			itemName, err := cmd.Flags().GetString("item-name")
			require.NoError(t, err)
			require.Equal(t, tc.wantItemName, itemName)

			githubOutput, err := cmd.Flags().GetBool("github-output")
			require.NoError(t, err)
			require.Equal(t, tc.wantGithubOutput, githubOutput)
		})
	}
}

// TestMondayGetItemIDRes_ToTable verifies that ToTable renders the item ID
// in the output.
func TestMondayGetItemIDRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.GetItemIDRes{
		ItemID: "9876543",
	}

	output := res.ToTable().Render()
	require.Contains(t, output, "9876543")
}

// TestMondayGetItemIDRes_ToJSON verifies that ToJSON produces valid JSON
// containing the expected item_id field.
func TestMondayGetItemIDRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.GetItemIDRes{
		ItemID: "9876543",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, "9876543", got["item_id"])
}
