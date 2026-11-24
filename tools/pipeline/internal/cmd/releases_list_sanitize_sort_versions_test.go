// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestReleasesSanitizeSortVersionsCmd_FlagParsing verifies that --versions
// is correctly split into the Versions slice and that --asc toggles sort
// direction without requiring a real Run execution.
func TestReleasesSanitizeSortVersionsCmd_FlagParsing(t *testing.T) {
	for name, tc := range map[string]struct {
		args           []string
		wantVersions   []string
		wantDescending bool
		wantErr        bool
	}{
		"versions flag is split on whitespace": {
			args:           []string{"--versions", "1.21.0 1.20.4 1.19.10"},
			wantVersions:   []string{"1.21.0", "1.20.4", "1.19.10"},
			wantDescending: true,
		},
		"default sort is descending": {
			args:           []string{"--versions", "1.21.0 1.20.4"},
			wantDescending: true,
		},
		"--asc sets descending to false": {
			args:           []string{"--versions", "1.21.0 1.20.4", "--asc"},
			wantDescending: false,
		},
		"missing --versions flag returns error": {
			args:    []string{},
			wantErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			sanitizeSortVersionsReq = &releases.SanitizeSortVersionsReq{Descending: true}
			cmd := newReleasesSanitizeSortVersionsCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tc.wantVersions != nil {
				require.Equal(t, tc.wantVersions, sanitizeSortVersionsReq.Versions)
			}
			require.Equal(t, tc.wantDescending, sanitizeSortVersionsReq.Descending)
		})
	}
}
