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

// TestMondayDuplicateBoardCmd_FlagParsing verifies that required flags are
// enforced, that --template-id and --board-name are populated, and that
// --duplicate-type defaults to duplicate_board_with_structure when omitted.
func TestMondayDuplicateBoardCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args              []string
		wantTemplateID    string
		wantBoardName     string
		wantDuplicateType monday.DuplicateBoardType
		wantErr           bool
	}{
		"required --template-id missing returns error": {
			args:    []string{"--board-name", "Release 1.21"},
			wantErr: true,
		},
		"required --board-name missing returns error": {
			args:    []string{"--template-id", "99999"},
			wantErr: true,
		},
		"both required flags set, --duplicate-type defaults to with_structure": {
			args:              []string{"--template-id", "99999", "--board-name", "Release 1.21"},
			wantTemplateID:    "99999",
			wantBoardName:     "Release 1.21",
			wantDuplicateType: monday.DuplicateBoardWithStructure,
		},
		"explicit --duplicate-type with_pulses is respected": {
			args:              []string{"--template-id", "99999", "--board-name", "Release 1.21", "--duplicate-type", "duplicate_board_with_pulses"},
			wantTemplateID:    "99999",
			wantBoardName:     "Release 1.21",
			wantDuplicateType: monday.DuplicateBoardWithPulses,
		},
		"explicit --duplicate-type with_pulses_and_updates is respected": {
			args:              []string{"--template-id", "99999", "--board-name", "Release 1.21", "--duplicate-type", "duplicate_board_with_pulses_and_updates"},
			wantTemplateID:    "99999",
			wantBoardName:     "Release 1.21",
			wantDuplicateType: monday.DuplicateBoardWithPulsesAndUpdates,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayDuplicateBoardCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			templateID, err := cmd.Flags().GetString("template-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantTemplateID, templateID)

			boardName, err := cmd.Flags().GetString("board-name")
			require.NoError(t, err)
			require.Equal(t, tc.wantBoardName, boardName)

			dupType, err := cmd.Flags().GetString("duplicate-type")
			require.NoError(t, err)
			require.Equal(t, string(tc.wantDuplicateType), dupType)
		})
	}
}

// TestMondayDuplicateBoardRes_ToTable verifies that ToTable renders the board
// ID and board name in the output.
func TestMondayDuplicateBoardRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.DuplicateBoardRes{
		BoardID:   "87654321",
		BoardName: "Release 1.21 (from template)",
	}

	output := res.ToTable().Render()
	require.Contains(t, output, "87654321")
	require.Contains(t, output, "Release 1.21 (from template)")
}

// TestMondayDuplicateBoardRes_ToJSON verifies that ToJSON produces valid JSON
// containing the expected board_id and board_name fields.
func TestMondayDuplicateBoardRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.DuplicateBoardRes{
		BoardID:   "87654321",
		BoardName: "Release 1.21 (from template)",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, "87654321", got["board_id"])
	require.Equal(t, "Release 1.21 (from template)", got["board_name"])
}
