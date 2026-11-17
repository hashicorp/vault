// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestGithubCheckBackportLabelsCmd_FlagParsing verifies that flags are parsed,
// and the fix and fail flags default to off.
func TestGithubCheckBackportLabelsCmd_FlagParsing(t *testing.T) {
	for name, tc := range map[string]struct {
		args       []string
		wantCreate bool
		wantDelete bool
		wantFail   bool
	}{
		"defaults are report-only": {args: []string{}},
		"fix flags": {
			args:       []string{"--create-missing", "--delete-stale"},
			wantCreate: true,
			wantDelete: true,
		},
		"fail on missing": {args: []string{"--fail-on-missing"}, wantFail: true},
	} {
		t.Run(name, func(t *testing.T) {
			*checkBackportLabelsReq = github.CheckBackportLabelsReq{}
			checkBackportLabelsFailOnMissing = false

			cmd := newGithubCheckBackportLabelsCmd()
			cmd.RunE = func(c *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)
			require.NoError(t, cmd.Execute())

			require.Equal(t, tc.wantCreate, checkBackportLabelsReq.CreateMissing)
			require.Equal(t, tc.wantDelete, checkBackportLabelsReq.DeleteStale)
			require.Equal(t, tc.wantFail, checkBackportLabelsFailOnMissing)
		})
	}
}

// TestGithubCheckBackportLabelsCmd_Run verifies that drift exits zero by
// default, exits non-zero with --fail-on-missing, and that --github-output
// writes the result.
func TestGithubCheckBackportLabelsCmd_Run(t *testing.T) {
	for name, tc := range map[string]struct {
		args      []string
		labels    string
		wantErr   string
		wantInOut string
	}{
		"drift exits zero by default": {
			labels:    `[{"name":"backport/1.17.x"}]`,
			wantInOut: "in_sync=false",
		},
		"drift fails with flag": {
			args:      []string{"--fail-on-missing"},
			labels:    `[{"name":"backport/1.17.x"}]`,
			wantErr:   "out of sync",
			wantInOut: "in_sync=false",
		},
		"in sync passes with flag": {
			args:      []string{"--fail-on-missing"},
			labels:    `[{"name":"backport/1.19.x"}]`,
			wantInOut: "in_sync=true",
		},
	} {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tc.labels))
			}))
			t.Cleanup(srv.Close)

			client := gh.NewClient(nil)
			base, err := url.Parse(srv.URL + "/")
			require.NoError(t, err)
			client.BaseURL = base

			oldState, oldDecodeRes, oldFormat := *githubCmdState, rootCfg.versionsDecodeRes, rootCfg.format
			githubCmdState.GithubV3 = client
			githubCmdState.Owner, githubCmdState.Repo = "hashicorp", "vault-enterprise"
			rootCfg.format = "table"
			rootCfg.versionsDecodeRes = &releases.DecodeRes{Config: &releases.VersionsConfig{
				ActiveVersion: &releases.ActiveVersion{Versions: map[string]*releases.Version{"1.19.x": {}}},
			}}
			t.Cleanup(func() {
				*githubCmdState, rootCfg.versionsDecodeRes, rootCfg.format = oldState, oldDecodeRes, oldFormat
			})

			githubOutput := filepath.Join(t.TempDir(), "github_output")
			t.Setenv("GITHUB_OUTPUT", githubOutput)

			*checkBackportLabelsReq = github.CheckBackportLabelsReq{}
			checkBackportLabelsFailOnMissing = false
			cmd := newGithubCheckBackportLabelsCmd()
			suppressOutput(t, cmd)
			cmd.SetArgs(append([]string{"--github-output"}, tc.args...))

			err = cmd.Execute()
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			out, err := os.ReadFile(githubOutput)
			require.NoError(t, err)
			require.Contains(t, string(out), tc.wantInOut)
		})
	}
}
