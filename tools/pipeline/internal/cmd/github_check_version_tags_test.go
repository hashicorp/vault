// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestGithubCheckVersionTagsCmd_FlagParsing verifies that required flags are
// enforced, versions are split correctly, and that the +ent suffix is derived
// automatically for vault-enterprise but not for other repos.
func TestGithubCheckVersionTagsCmd_FlagParsing(t *testing.T) {
	for name, tc := range map[string]struct {
		args          []string
		wantVersions  []string
		wantTagSuffix string
		wantErr       bool
	}{
		"required --repo missing returns error": {
			args:    []string{"--versions", "1.21.0"},
			wantErr: true,
		},
		"required --versions missing returns error": {
			args:    []string{"--repo", "vault-enterprise"},
			wantErr: true,
		},
		"vault-enterprise gets +ent suffix automatically": {
			args:          []string{"--repo", "vault-enterprise", "--versions", "1.21.0 1.20.4"},
			wantVersions:  []string{"1.21.0", "1.20.4"},
			wantTagSuffix: "+ent",
		},
		"non-enterprise repo gets no suffix": {
			args:          []string{"--repo", "some-oss-project", "--versions", "1.21.0"},
			wantVersions:  []string{"1.21.0"},
			wantTagSuffix: "",
		},
		"explicit --tag-suffix overrides automatic derivation": {
			args:          []string{"--repo", "vault-enterprise", "--versions", "1.21.0", "--tag-suffix", ""},
			wantVersions:  []string{"1.21.0"},
			wantTagSuffix: "",
		},
		"explicit non-empty --tag-suffix is preserved": {
			args:          []string{"--repo", "some-project", "--versions", "1.21.0", "--tag-suffix", "-custom"},
			wantVersions:  []string{"1.21.0"},
			wantTagSuffix: "-custom",
		},
	} {
		t.Run(name, func(t *testing.T) {
			checkVersionTagsReq = &github.CheckVersionTagsReq{}
			cmd := newGithubCheckVersionTagsCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tc.wantVersions != nil {
				require.Equal(t, tc.wantVersions, checkVersionTagsReq.Versions)
			}
			require.Equal(t, tc.wantTagSuffix, checkVersionTagsReq.TagSuffix)
		})
	}
}
