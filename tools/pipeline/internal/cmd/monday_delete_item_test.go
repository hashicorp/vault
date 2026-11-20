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

// TestMondayDeleteItemCmd_FlagParsing verifies that --board-id and --item-id
// are required and that --expect-name and --dry-run are parsed. The rule that
// --expect-name is required unless --dry-run is set lives in
// monday.DeleteItemReq.Run, so it applies to every caller, not only the CLI.
func TestMondayDeleteItemCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args           []string
		wantBoardID    string
		wantItemID     string
		wantExpectName string
		wantDryRun     bool
		wantErr        bool
	}{
		"required --board-id missing returns error": {
			args:    []string{"--item-id", "2222222222", "--expect-name", "Test item"},
			wantErr: true,
		},
		"required --item-id missing returns error": {
			args:    []string{"--board-id", "1111111111", "--expect-name", "Test item"},
			wantErr: true,
		},
		"real run with --expect-name": {
			args:           []string{"--board-id", "1111111111", "--item-id", "2222222222", "--expect-name", "Test item"},
			wantBoardID:    "1111111111",
			wantItemID:     "2222222222",
			wantExpectName: "Test item",
		},
		"dry run without --expect-name": {
			args:        []string{"--board-id", "1111111111", "--item-id", "2222222222", "--dry-run"},
			wantBoardID: "1111111111",
			wantItemID:  "2222222222",
			wantDryRun:  true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayDeleteItemCmd()
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

			itemID, err := cmd.Flags().GetString("item-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantItemID, itemID)

			expectName, err := cmd.Flags().GetString("expect-name")
			require.NoError(t, err)
			require.Equal(t, tc.wantExpectName, expectName)

			dryRun, err := cmd.Flags().GetBool("dry-run")
			require.NoError(t, err)
			require.Equal(t, tc.wantDryRun, dryRun)
		})
	}
}

// TestMondayDeleteItemRes_ToTable verifies that ToTable renders the details a
// human needs to confirm the target, including the parent and subitem count.
func TestMondayDeleteItemRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.DeleteItemRes{
		ItemID:       "6666666666",
		ItemName:     "Test subitem",
		BoardID:      "4444444444",
		GroupID:      "topics",
		GroupTitle:   "Subitems",
		ParentItemID: "2222222222",
		SubitemCount: 3,
	}

	output := res.ToTable().Render()
	for _, want := range []string{"6666666666", "Test subitem", "4444444444", "topics", "Subitems", "2222222222", "3", "false"} {
		require.Contains(t, output, want)
	}
}

// TestMondayDeleteItemRes_ToJSON verifies the JSON field names, that
// parent_item_id is omitted for a top-level item, and that deleted is reported.
func TestMondayDeleteItemRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.DeleteItemRes{
		ItemID:       "2222222222",
		ItemName:     "Test item",
		BoardID:      "1111111111",
		GroupID:      "topics",
		GroupTitle:   "Group Title",
		SubitemCount: 1,
		Deleted:      true,
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, map[string]any{
		"item_id":       "2222222222",
		"item_name":     "Test item",
		"board_id":      "1111111111",
		"group_id":      "topics",
		"group_title":   "Group Title",
		"subitem_count": float64(1),
		"deleted":       true,
	}, got)
}
