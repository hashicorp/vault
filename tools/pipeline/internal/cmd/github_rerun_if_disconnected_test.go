// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGithubRerunIfDisconnectedCmd checks that the command accepts valid arguments,
// uses inherited flags, writes the requested format, and returns errors to the caller.
func TestGithubRerunIfDisconnectedCmd(t *testing.T) {
	t.Parallel()
	for name, test := range map[string]struct {
		args         []string
		status       int
		expect       string
		expectErr    string
		reqCount     int32
		writeErr     bool
		disconnected bool
	}{
		"table":             {args: []string{"1"}, expect: "Summary:", reqCount: 1},
		"rerun accepted":    {args: []string{"1", "-f", "json"}, expect: `"rerun_requested":true`, reqCount: 2, disconnected: true},
		"markdown":          {args: []string{"1", "-f", "markdown"}, expect: "## Summary", reqCount: 1},
		"json":              {args: []string{"1", "--format", "json"}, expect: `"rerun_requested":false,"failed_jobs":[]`, reqCount: 1},
		"custom repository": {args: []string{"1", "-o", "custom", "-r", "repo", "--github-host", "ghe.example"}, expect: "custom/repo", reqCount: 1},
		"no ID":             {expectErr: "accepts 1 arg"},
		"extra ID":          {args: []string{"1", "2"}, expectErr: "accepts 1 arg"},
		"invalid ID":        {args: []string{"abc"}, expectErr: "positive int64"},
		"zero ID":           {args: []string{"0"}, expectErr: "positive int64"},
		"negative ID":       {args: []string{"--", "-1"}, expectErr: "positive int64"},
		"overflow ID":       {args: []string{"9223372036854775808"}, expectErr: "positive int64"},
		"csv":               {args: []string{"1", "-f", "csv"}, expectErr: "unsupported format"},
		"unknown format":    {args: []string{"1", "-f", "xml"}, expectErr: "unsupported format"},
		"API failure":       {args: []string{"1"}, status: 403, expectErr: "rerun if disconnected: inspecting", reqCount: 1},
		"write failure":     {args: []string{"1"}, expectErr: "output unavailable", reqCount: 1, writeErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			state := &githubCommandState{Owner: "hashicorp", Repo: "vault-enterprise", GithubHost: "github.com"}
			var requests atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if r.Method == http.MethodPost {
					assert.True(t, test.disconnected)
					assert.Equal(t, "/repos/"+state.Owner+"/"+state.Repo+"/actions/runs/1/rerun-failed-jobs", r.URL.Path)
					w.WriteHeader(http.StatusCreated)
					return
				}
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/repos/"+state.Owner+"/"+state.Repo+"/actions/runs/1/jobs", r.URL.Path)
				if test.status != 0 {
					w.WriteHeader(test.status)
					return
				}
				body := `{"jobs":[]}`
				if test.disconnected {
					body = `{"jobs":[{"status":"completed","conclusion":"failure","steps":[{"status":"pending"}]}]}`
				}
				_, err := io.WriteString(w, body)
				assert.NoError(t, err)
			}))
			t.Cleanup(srv.Close)
			state.GithubV3 = gh.NewClient(srv.Client())
			var err error
			state.GithubV3.BaseURL, err = url.Parse(srv.URL + "/")
			require.NoError(t, err)
			parent := &cobra.Command{Use: "github"}
			parent.PersistentFlags().StringVarP(&state.Owner, "owner", "o", "hashicorp", "")
			parent.PersistentFlags().StringVarP(&state.Repo, "repo", "r", "vault-enterprise", "")
			parent.PersistentFlags().StringVar(&state.GithubHost, "github-host", "github.com", "")
			parent.PersistentFlags().StringP("format", "f", "table", "")
			parent.AddCommand(newGithubRerunIfDisconnectedCmd(state))
			var output bytes.Buffer
			parent.SetOut(&output)
			parent.SetErr(io.Discard)
			if test.writeErr {
				parent.SetOut(rerunErrorWriter{})
			}
			parent.SetArgs(append([]string{"rerun-if-disconnected"}, test.args...))
			err = parent.Execute()
			if test.expectErr != "" {
				require.ErrorContains(t, err, test.expectErr)
			} else {
				require.NoError(t, err)
				require.Contains(t, output.String(), test.expect)
			}
			require.Equal(t, test.reqCount, requests.Load())
			if name == "custom repository" {
				require.Equal(t, "ghe.example", state.GithubHost)
			}
		})
	}
}

type rerunErrorWriter struct{}

func (rerunErrorWriter) Write([]byte) (int, error) {
	return 0, errors.New("output unavailable")
}

// TestGithubRerunIfDisconnectedCmd_Registration checks that inherited flags don't
// conflict with existing commands or change their defaults. We run it sequentially
// because the older commands bind flags to globals.
func TestGithubRerunIfDisconnectedCmd_Registration(t *testing.T) {
	parent := newGithubCmd()
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		require.NotPanics(t, func() { cmd.InheritedFlags() }, cmd.CommandPath())
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(parent)
	cmd, _, err := parent.Find([]string{"rerun-if-disconnected"})
	require.NoError(t, err)
	require.Equal(t, "rerun-if-disconnected", cmd.Name())
	require.Equal(t, "hashicorp", cmd.InheritedFlags().Lookup("owner").DefValue)
	require.Equal(t, "vault-enterprise", cmd.InheritedFlags().Lookup("repo").DefValue)
	for _, tc := range []struct {
		path []string
		repo string
	}{
		{[]string{"check", "workflow-run"}, "vault"},
		{[]string{"create", "backport"}, "vault-enterprise"},
	} {
		child, _, err := parent.Find(tc.path)
		require.NoError(t, err)
		require.Equal(t, tc.repo, child.PersistentFlags().Lookup("repo").DefValue)
	}
}
