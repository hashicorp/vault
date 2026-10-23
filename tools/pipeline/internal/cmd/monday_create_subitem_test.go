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

// TestMondayCreateSubitemCmd_FlagParsing verifies that required flags are
// enforced, that --parent-item-id and --item-name are correctly populated when
// provided, and that --github-output is correctly parsed when set.
func TestMondayCreateSubitemCmd_FlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args             []string
		wantParentItemID string
		wantItemName     string
		wantGithubOutput bool
		wantErr          bool
	}{
		"required --parent-item-id missing returns error": {
			args:    []string{"--item-name", "Release Checklist"},
			wantErr: true,
		},
		"required --item-name missing returns error": {
			args:    []string{"--parent-item-id", "1210003"},
			wantErr: true,
		},
		"all required flags set": {
			args:             []string{"--parent-item-id", "1210003", "--item-name", "Release Checklist"},
			wantParentItemID: "1210003",
			wantItemName:     "Release Checklist",
		},
		"--github-output flag is set": {
			args:             []string{"--parent-item-id", "1210003", "--item-name", "Release Checklist", "--github-output"},
			wantParentItemID: "1210003",
			wantItemName:     "Release Checklist",
			wantGithubOutput: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newMondayCreateSubitemCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			parentItemID, err := cmd.Flags().GetString("parent-item-id")
			require.NoError(t, err)
			require.Equal(t, tc.wantParentItemID, parentItemID)

			itemName, err := cmd.Flags().GetString("item-name")
			require.NoError(t, err)
			require.Equal(t, tc.wantItemName, itemName)

			githubOutput, err := cmd.Flags().GetBool("github-output")
			require.NoError(t, err)
			require.Equal(t, tc.wantGithubOutput, githubOutput)
		})
	}
}

// TestMondayCreateSubitemRes_ToTable verifies that ToTable renders the
// sub-item ID in the output.
func TestMondayCreateSubitemRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &monday.CreateSubitemRes{
		SubitemID: "2020333",
	}

	output := res.ToTable().Render()
	require.Contains(t, output, "2020333")
}

// TestMondayCreateSubitemRes_ToJSON verifies that ToJSON produces valid JSON
// containing the expected subitem_id field.
func TestMondayCreateSubitemRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &monday.CreateSubitemRes{
		SubitemID: "2020333",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, "2020333", got["subitem_id"])
}
