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

// TestMondayCreateBoardCmd_FlagParsing verifies that required flags are
// enforced, that --board-name is populated, and that --board-kind defaults to
// "public" when omitted.
func TestMondayCreateBoardCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args          []string
		wantBoardName string
		wantBoardKind monday.BoardKind
		wantErr       bool
	}{
		"required --board-name missing returns error": {
			args:    []string{},
			wantErr: true,
		},
		"--board-name is set, --board-kind defaults to public": {
			args:          []string{"--board-name", "Release 1.21"},
			wantBoardName: "Release 1.21",
			wantBoardKind: monday.BoardKindPublic,
		},
		"explicit --board-kind private is respected": {
			args:          []string{"--board-name", "Release 1.21", "--board-kind", "private"},
			wantBoardName: "Release 1.21",
			wantBoardKind: monday.BoardKindPrivate,
		},
		"explicit --board-kind share is respected": {
			args:          []string{"--board-name", "Release 1.21", "--board-kind", "share"},
			wantBoardName: "Release 1.21",
			wantBoardKind: monday.BoardKindShare,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayCreateBoardCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Inspect the flags directly from the command rather than a
			// package-level req variable — create-board uses a local req.
			boardName, err := cmd.Flags().GetString("board-name")
			require.NoError(t, err)
			require.Equal(t, tc.wantBoardName, boardName)

			boardKind, err := cmd.Flags().GetString("board-kind")
			require.NoError(t, err)
			require.Equal(t, string(tc.wantBoardKind), boardKind)
		})
	}
}

// TestMondayCreateBoardRes_ToTable verifies that ToTable renders the board ID
// and board name in the output.
func TestMondayCreateBoardRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.CreateBoardRes{
		BoardID:   "12345678",
		BoardName: "Release 1.21",
	}

	output := res.ToTable().Render()
	require.Contains(t, output, "12345678")
	require.Contains(t, output, "Release 1.21")
}

// TestMondayCreateBoardRes_ToJSON verifies that ToJSON produces valid JSON
// containing the expected board_id and board_name fields.
func TestMondayCreateBoardRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.CreateBoardRes{
		BoardID:   "12345678",
		BoardName: "Release 1.21",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, "12345678", got["board_id"])
	require.Equal(t, "Release 1.21", got["board_name"])
}
