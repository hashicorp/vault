// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestGithubCheckGoVersionCmd_RequiredFlags verifies that --repo is enforced
// and the command errors without it.
func TestGithubCheckGoVersionCmd_RequiredFlags(t *testing.T) {
	for name, tc := range map[string]struct {
		args    []string
		wantErr bool
	}{
		"missing --repo": {args: []string{}, wantErr: true},
		"repo set":       {args: []string{"--repo", "vault-enterprise"}, wantErr: false},
	} {
		t.Run(name, func(t *testing.T) {
			checkGoVersionReq = &github.CheckGoVersionReq{}
			cmd := newGithubCheckGoVersionCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestBranchesToCheck verifies that positional branch arguments take precedence
// over the branches derived from the versions config, including when
// --include-main is set, and that the config is used when no args are given.
func TestBranchesToCheck(t *testing.T) {
	t.Parallel()

	dec := &releases.DecodeRes{
		Config: &releases.VersionsConfig{
			ActiveVersion: &releases.ActiveVersion{
				Versions: map[string]*releases.Version{"1.21.x": {}},
			},
		},
	}

	for name, tc := range map[string]struct {
		args        []string
		includeMain bool
		want        []string
	}{
		"positional args bypass versions config": {
			args: []string{"main", "release/1.20.x+ent"},
			want: []string{"main", "release/1.20.x+ent"},
		},
		"positional args ignore include-main": {
			args:        []string{"release/1.20.x+ent"},
			includeMain: true,
			want:        []string{"release/1.20.x+ent"},
		},
		"no args uses versions config": {
			want: []string{"release/1.21.x+ent"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, branchesToCheck(tc.args, dec, tc.includeMain))
		})
	}
}
