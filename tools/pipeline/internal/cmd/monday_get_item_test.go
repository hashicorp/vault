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

// TestMondayGetItemCmd_FlagParsing verifies that --board-id and --item-id are
// both required and are populated when provided.
func TestMondayGetItemCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args        []string
		wantBoardID string
		wantItemID  string
		wantErr     bool
	}{
		"required --board-id missing returns error": {
			args:    []string{"--item-id", "2222222222"},
			wantErr: true,
		},
		"required --item-id missing returns error": {
			args:    []string{"--board-id", "1111111111"},
			wantErr: true,
		},
		"--board-id and --item-id set": {
			args:        []string{"--board-id", "1111111111", "--item-id", "2222222222"},
			wantBoardID: "1111111111",
			wantItemID:  "2222222222",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayGetItemCmd()
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
		})
	}
}

// TestMondayGetItemRes_ToTable verifies that ToTable renders the column text
// (not the raw value), the parent id and the subitem ids.
func TestMondayGetItemRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.GetItemRes{
		ItemID:       "2222222222",
		ItemName:     "Test item",
		BoardID:      "1111111111",
		GroupID:      "topics",
		GroupTitle:   "Group Title",
		State:        "active",
		ParentItemID: "3333333333",
		SubitemCount: 2,
		SubitemIDs:   []string{"5555555555", "5555555556"},
		UpdatedAt:    "2026-05-01T10:00:00Z",
		ColumnValues: []monday.ColumnValue{
			{ID: "status", Title: "Status", Type: "status", Text: "Done", Value: `{"index":1}`},
		},
	}

	output := res.ToTable().Render()
	for _, want := range []string{
		"2222222222", "Test item", "1111111111", "topics", "Group Title", "active",
		"3333333333", "5555555555,5555555556", "2026-05-01T10:00:00Z", "Status", "Done",
	} {
		require.Contains(t, output, want)
	}
	require.NotContains(t, output, `{"index":1}`)
}

// TestMondayGetItemRes_ToJSON verifies the JSON field names, that each column
// value carries both text and the raw value, that parent_item_id is omitted for
// a top-level item, and that subitem_ids marshals as [] rather than null.
func TestMondayGetItemRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.GetItemRes{
		ItemID:     "2222222222",
		ItemName:   "Test item",
		BoardID:    "1111111111",
		GroupID:    "topics",
		GroupTitle: "Group Title",
		State:      "active",
		SubitemIDs: []string{},
		UpdatedAt:  "2026-05-01T10:00:00Z",
		ColumnValues: []monday.ColumnValue{
			{ID: "status", Title: "Status", Type: "status", Text: "Done", Value: `{"index":1}`},
		},
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
		"state":         "active",
		"subitem_count": float64(0),
		"subitem_ids":   []any{},
		"updated_at":    "2026-05-01T10:00:00Z",
		"column_values": []any{
			map[string]any{"id": "status", "title": "Status", "type": "status", "text": "Done", "value": `{"index":1}`},
		},
	}, got)
}
