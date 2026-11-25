// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github/githubtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRerunIfDisconnected_SkipReason checks that we only inspect a run that failed
// on the attempt we expected. Someone else might have already rerun it.
func TestRerunIfDisconnected_SkipReason(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		attempt    int
		runAttempt int
		conclusion string
		expect     string
	}{
		"failed":        {attempt: 1, runAttempt: 1, conclusion: "failure"},
		"any attempt":   {runAttempt: 2, conclusion: "failure"},
		"cancelled":     {runAttempt: 1, conclusion: "cancelled", expect: "the run didn't fail, it concluded with cancelled"},
		"succeeded":     {runAttempt: 1, conclusion: "success", expect: "the run didn't fail, it concluded with success"},
		"already rerun": {attempt: 1, runAttempt: 2, expect: "we expected attempt 1 but the run is on attempt 2"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &RerunIfDisconnectedReq{Attempt: test.attempt}
			run := &gh.WorkflowRun{RunAttempt: gh.Ptr(test.runAttempt), Conclusion: gh.Ptr(test.conclusion)}
			require.Equal(t, test.expect, req.skipReason(run))
		})
	}
}

// TestRerunIfDisconnected_NewerRun checks that we only treat a run as superseded
// by a different run that was created after it. Rerunning a superseded run might
// cancel the newer one.
func TestRerunIfDisconnected_NewerRun(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 11, 10, 18, 0, 0, 0, time.UTC)
	run := &gh.WorkflowRun{ID: gh.Ptr(int64(1)), CreatedAt: &gh.Timestamp{Time: created}}
	older := &gh.WorkflowRun{ID: gh.Ptr(int64(2)), CreatedAt: &gh.Timestamp{Time: created.Add(-time.Hour)}}
	newer := &gh.WorkflowRun{ID: gh.Ptr(int64(3)), CreatedAt: &gh.Timestamp{Time: created.Add(time.Hour)}}

	req := &RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: 1}
	require.Nil(t, req.newerWorkflowRun(run, nil))
	require.Nil(t, req.newerWorkflowRun(run, []*gh.WorkflowRun{run}))
	require.Nil(t, req.newerWorkflowRun(run, []*gh.WorkflowRun{older}))
	require.Equal(t, newer, req.newerWorkflowRun(run, []*gh.WorkflowRun{older, newer}))
}

// TestRerunIfDisconnected_APIFailures checks that we return Github's errors to the
// caller. If we can't tell whether a newer run exists we ought not rerun. If Github
// rejects the rerun, e.g. because the run is already running again, the caller
// needs to know. It also checks that we inspect the jobs of the attempt that we
// checked, since only that route is mocked.
func TestRerunIfDisconnected_APIFailures(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		latestStatus int
		expectPosts  int32
		expectErr    string
	}{
		"latest run lookup failure": {latestStatus: http.StatusForbidden, expectErr: "checking for runs newer than o/r run 1"},
		"rerun rejected":            {expectPosts: 1, expectErr: "rerunning failed jobs for o/r run 1"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, mux, teardown := githubtest.SetupTestClient(t)
			defer teardown()

			var posts atomic.Int32
			mux.HandleFunc("GET /repos/o/r/actions/runs/1", func(w http.ResponseWriter, r *http.Request) {
				_, err := io.WriteString(w, `{"id":1,"workflow_id":7,"status":"completed","conclusion":"failure","run_attempt":1}`)
				assert.NoError(t, err)
			})
			mux.HandleFunc("GET /repos/o/r/actions/runs/1/attempts/1/jobs", func(w http.ResponseWriter, r *http.Request) {
				_, err := io.WriteString(w, `{"jobs":[{"status":"completed","conclusion":"failure","steps":[{"status":"pending"}]}]}`)
				assert.NoError(t, err)
			})
			mux.HandleFunc("GET /repos/o/r/actions/workflows/7/runs", func(w http.ResponseWriter, r *http.Request) {
				if test.latestStatus != 0 {
					w.WriteHeader(test.latestStatus)
					return
				}
				_, err := io.WriteString(w, `{"workflow_runs":[{"id":1}]}`)
				assert.NoError(t, err)
			})
			mux.HandleFunc("POST /repos/o/r/actions/runs/1/rerun-failed-jobs", func(w http.ResponseWriter, r *http.Request) {
				posts.Add(1)
				w.WriteHeader(http.StatusForbidden)
				_, err := io.WriteString(w, `{"message":"This workflow is already running"}`)
				assert.NoError(t, err)
			})

			res, err := (&RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: 1}).Run(t.Context(), client)
			require.ErrorContains(t, err, test.expectErr)
			require.Nil(t, res)
			require.Equal(t, test.expectPosts, posts.Load())
		})
	}
}

// TestRerunIfDisconnected_Wait checks that we poll a run until it has completed
// before we inspect it, and that we stop waiting if someone else reruns it.
func TestRerunIfDisconnected_Wait(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		polls []string
		rerun bool
		skip  string
	}{
		"completes": {
			polls: []string{
				`"status":"queued","run_attempt":1`,
				`"status":"in_progress","run_attempt":1`,
				`"status":"completed","conclusion":"failure","run_attempt":1`,
			},
			rerun: true,
		},
		"rerun by someone else": {
			polls: []string{
				`"status":"in_progress","run_attempt":1`,
				`"status":"queued","run_attempt":2`,
			},
			skip: "we expected attempt 1 but the run is on attempt 2",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, mux, teardown := githubtest.SetupTestClient(t)
			defer teardown()

			var runGets, posts atomic.Int32
			mux.HandleFunc("GET /repos/o/r/actions/runs/1", func(w http.ResponseWriter, r *http.Request) {
				poll := test.polls[min(int(runGets.Add(1)), len(test.polls))-1]
				_, err := fmt.Fprintf(w, `{"id":1,"workflow_id":7,%s}`, poll)
				assert.NoError(t, err)
			})
			mux.HandleFunc("GET /repos/o/r/actions/runs/1/attempts/1/jobs", func(w http.ResponseWriter, r *http.Request) {
				_, err := io.WriteString(w, `{"jobs":[{"status":"completed","conclusion":"failure","steps":[{"status":"pending"}]}]}`)
				assert.NoError(t, err)
			})
			mux.HandleFunc("GET /repos/o/r/actions/workflows/7/runs", func(w http.ResponseWriter, r *http.Request) {
				_, err := io.WriteString(w, `{"workflow_runs":[{"id":1}]}`)
				assert.NoError(t, err)
			})
			mux.HandleFunc("POST /repos/o/r/actions/runs/1/rerun-failed-jobs", func(w http.ResponseWriter, r *http.Request) {
				posts.Add(1)
				w.WriteHeader(http.StatusCreated)
			})

			req := &RerunIfDisconnectedReq{
				Owner: "o", Repo: "r", RunID: 1, Attempt: 1, Wait: true, PollInterval: time.Millisecond,
			}
			res, err := req.Run(t.Context(), client)
			require.NoError(t, err)
			require.Equal(t, test.rerun, res.RerunRequested)
			require.Equal(t, test.skip, res.SkipReason)
			require.EqualValues(t, len(test.polls), runGets.Load())
			if test.rerun {
				require.EqualValues(t, 1, posts.Load())
			} else {
				require.Zero(t, posts.Load())
			}
		})
	}
}

// TestRerunIfDisconnected_WaitTimeout checks that we give up on a run that never
// completes and that we tell the caller why we stopped waiting.
func TestRerunIfDisconnected_WaitTimeout(t *testing.T) {
	t.Parallel()

	client, mux, teardown := githubtest.SetupTestClient(t)
	defer teardown()

	mux.HandleFunc("GET /repos/o/r/actions/runs/1", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"id":1,"status":"in_progress"}`)
		assert.NoError(t, err)
	})

	timedOut := errors.New("timed out")
	ctx, cancel := context.WithTimeoutCause(t.Context(), 50*time.Millisecond, timedOut)
	defer cancel()

	req := &RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: 1, Wait: true, PollInterval: 5 * time.Millisecond}
	res, err := req.Run(ctx, client)
	require.ErrorIs(t, err, timedOut)
	require.Nil(t, res)
}

// TestRerunIfDisconnected_Validation checks that invalid requests fail before we
// call the API.
func TestRerunIfDisconnected_Validation(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req    *RerunIfDisconnectedReq
		client *gh.Client
		expect string
	}{
		"nil request": {nil, gh.NewClient(nil), "uninitialized request"},
		"nil client":  {&RerunIfDisconnectedReq{}, nil, "client is required"},
		"owner":       {&RerunIfDisconnectedReq{Repo: "r", RunID: 1}, gh.NewClient(nil), "owner is required"},
		"repo":        {&RerunIfDisconnectedReq{Owner: "o", RunID: 1}, gh.NewClient(nil), "repo is required"},
		"zero ID":     {&RerunIfDisconnectedReq{Owner: "o", Repo: "r"}, gh.NewClient(nil), "run ID must be positive"},
		"negative attempt": {
			&RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: 1, Attempt: -1},
			gh.NewClient(nil), "attempt can't be negative",
		},
		"no poll interval": {
			&RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: 1, Wait: true},
			gh.NewClient(nil), "poll interval must be positive",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, err := test.req.Run(t.Context(), test.client)
			require.ErrorContains(t, err, test.expect)
			require.Nil(t, res)
		})
	}
}

// TestRerunIfDisconnected_Formatting checks that we tell the reader whether or
// not we requested a rerun, and why not if we skipped one.
func TestRerunIfDisconnected_Formatting(t *testing.T) {
	t.Parallel()

	res := &RerunIfDisconnectedRes{
		ListDisconnectedJobsRes: &ListDisconnectedJobsRes{
			Owner: "o", Repo: "r", RunID: 1, JobsInspected: 1, SuspectedDisconnect: true,
			FailedJobs: []*DisconnectedJobResult{{
				ID: 2, Name: "build", SuspectedDisconnect: true,
				UnfinishedSteps: []*DisconnectedStepResult{{Number: 51, Name: "test", Status: "pending"}},
			}},
		},
		RerunRequested: true,
	}
	for _, output := range []string{res.ToTable(), res.ToMarkdown()} {
		require.Contains(t, output, "Summary")
		require.Contains(t, output, "Unfinished Steps")
		require.Contains(t, output, "Rerun of failed jobs requested")
	}
	b, err := res.ToJSON()
	require.NoError(t, err)
	require.Contains(t, string(b), `"run_id":1,`)
	require.Contains(t, string(b), `"rerun_requested":true}`)
	var decoded RerunIfDisconnectedRes
	require.NoError(t, json.Unmarshal(b, &decoded))
	require.Equal(t, *res, decoded)

	skipped := &RerunIfDisconnectedRes{
		ListDisconnectedJobsRes: &ListDisconnectedJobsRes{FailedJobs: []*DisconnectedJobResult{}},
		SkipReason:              "newer run 2 has superseded it",
	}
	for _, output := range []string{skipped.ToTable(), skipped.ToMarkdown()} {
		require.Contains(t, output, "No rerun requested: newer run 2 has superseded it")
	}
	b, err = skipped.ToJSON()
	require.NoError(t, err)
	require.Contains(t, string(b), `"rerun_requested":false,"skip_reason":"newer run 2 has superseded it"`)

	none := &RerunIfDisconnectedRes{ListDisconnectedJobsRes: &ListDisconnectedJobsRes{FailedJobs: []*DisconnectedJobResult{}}}
	for _, output := range []string{none.ToTable(), none.ToMarkdown()} {
		require.Contains(t, output, "No suspected disconnects found. No rerun requested.")
	}

	var absent *RerunIfDisconnectedRes
	require.Empty(t, absent.ToTable())
	require.Empty(t, absent.ToMarkdown())
	_, err = absent.ToJSON()
	require.Error(t, err)
}
