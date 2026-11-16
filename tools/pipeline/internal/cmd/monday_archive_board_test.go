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

// TestMondayArchiveBoardCmd_FlagParsing verifies that --board-id is required
// and that --dry-run and --hard-delete are parsed.
func TestMondayArchiveBoardCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args           []string
		wantBoardID    string
		wantDryRun     bool
		wantHardDelete bool
		wantErr        bool
	}{
		"required --board-id missing returns error": {
			args:    []string{"--dry-run"},
			wantErr: true,
		},
		"archive": {
			args:        []string{"--board-id", "1111111111"},
			wantBoardID: "1111111111",
		},
		"hard delete": {
			args:           []string{"--board-id", "1111111111", "--hard-delete"},
			wantBoardID:    "1111111111",
			wantHardDelete: true,
		},
		"dry run": {
			args:        []string{"--board-id", "1111111111", "--dry-run"},
			wantBoardID: "1111111111",
			wantDryRun:  true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayArchiveBoardCmd()
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

			dryRun, err := cmd.Flags().GetBool("dry-run")
			require.NoError(t, err)
			require.Equal(t, tc.wantDryRun, dryRun)

			hardDelete, err := cmd.Flags().GetBool("hard-delete")
			require.NoError(t, err)
			require.Equal(t, tc.wantHardDelete, hardDelete)
		})
	}
}

// TestMondayArchiveBoardRes_ToTable verifies that ToTable renders the details a
// human needs to confirm the target, along with the action taken.
func TestMondayArchiveBoardRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.ArchiveBoardRes{
		BoardID:       "1111111111",
		BoardName:     "vault-52079 copy",
		State:         "active",
		WorkspaceID:   "3333333333",
		WorkspaceName: "Vault Releases",
		ItemsCount:    12,
		UpdatedAt:     "2026-11-10T15:04:05Z",
		Action:        monday.BoardActionNone,
	}

	output := res.ToTable().Render()
	for _, want := range []string{"1111111111", "vault-52079 copy", "active", "3333333333", "Vault Releases", "12", "2026-11-10T15:04:05Z", "none"} {
		require.Contains(t, output, want)
	}
}

// TestMondayArchiveBoardRes_ToJSON verifies the JSON field names, that
// workspace_id is omitted for the Main workspace, and that the action is
// reported.
func TestMondayArchiveBoardRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.ArchiveBoardRes{
		BoardID:       "2222222222",
		BoardName:     "vault-52079 copy",
		State:         "active",
		WorkspaceName: monday.MainWorkspaceName,
		ItemsCount:    0,
		UpdatedAt:     "2026-11-09T08:00:00Z",
		Action:        monday.BoardActionDeleted,
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, map[string]any{
		"board_id":       "2222222222",
		"board_name":     "vault-52079 copy",
		"state":          "active",
		"workspace_name": "Main workspace",
		"items_count":    float64(0),
		"updated_at":     "2026-11-09T08:00:00Z",
		"action":         "deleted",
	}, got)
}
