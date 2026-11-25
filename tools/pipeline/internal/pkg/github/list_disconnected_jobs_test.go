// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github/githubtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListDisconnectedJobs_ObservedJob checks a sanitized example of a job whose
// runner disconnected. We ought to flag it, list its unfinished steps in order,
// and ignore any jobs that haven't failed.
func TestListDisconnectedJobs_ObservedJob(t *testing.T) {
	t.Parallel()

	b, err := os.ReadFile("testfixtures/disconnected_job.json")
	require.NoError(t, err)
	var job gh.WorkflowJob
	require.NoError(t, json.Unmarshal(b, &job))

	req := &ListDisconnectedJobsReq{Owner: "o", Repo: "r", RunID: 1}
	res := req.inspectJobs(t.Context(), []*gh.WorkflowJob{
		&job,
		{
			ID:         gh.Ptr(int64(2)),
			Name:       gh.Ptr("lint"),
			Status:     gh.Ptr("completed"),
			Conclusion: gh.Ptr("failure"),
		},
		{
			Status: gh.Ptr("in_progress"),
			Steps:  []*gh.TaskStep{{Status: gh.Ptr("pending")}},
		},
	})
	require.True(t, res.SuspectedDisconnect)
	require.Equal(t, 1, res.RunAttempt)
	require.Equal(t, 3, res.JobsInspected)
	require.Len(t, res.FailedJobs, 2)
	result := res.FailedJobs[0]
	require.EqualValues(t, 1001, result.ID)
	require.Equal(t, job.GetName(), result.Name)
	require.Equal(t, job.GetHTMLURL(), result.URL)
	require.True(t, result.SuspectedDisconnect)
	require.Len(t, result.UnfinishedSteps, 15)
	for i, step := range result.UnfinishedSteps {
		require.Equal(t, job.Steps[i+13].GetNumber(), step.Number)
		require.Equal(t, job.Steps[i+13].GetName(), step.Name)
		require.Equal(t, job.Steps[i+13].GetStatus(), step.Status)
	}
	require.False(t, res.FailedJobs[1].SuspectedDisconnect)

	b, err = res.ToJSON()
	require.NoError(t, err)
	require.Contains(t, string(b), `"suspected_disconnect":false,"unfinished_steps":[]`)
}

// TestListDisconnectedJobs_Classification checks that we only suspect a disconnect
// when a completed, failed job still has unfinished steps.
func TestListDisconnectedJobs_Classification(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		status     string
		conclusion string
		steps      string
		failed     bool
		suspected  bool
	}{
		"queued":           {"completed", "failure", `[{"status":"queued"}]`, true, true},
		"pending":          {"completed", "failure", `[{"status":"pending"}]`, true, true},
		"in progress":      {"completed", "failure", `[{"status":"in_progress"}]`, true, true},
		"ordinary failure": {"completed", "failure", `[{"status":"completed","conclusion":"failure"}]`, true, false},
		"missing evidence": {"completed", "failure", `[null,{}, {"status":"unknown"}]`, true, false},
		"no steps":         {"completed", "failure", `null`, true, false},
		"empty steps":      {"completed", "failure", `[]`, true, false},
		"success":          {"completed", "success", `[{"status":"pending"}]`, false, false},
		"cancelled":        {"completed", "cancelled", `[{"status":"pending"}]`, false, false},
		"timed out":        {"completed", "timed_out", `[{"status":"pending"}]`, false, false},
		"skipped":          {"completed", "skipped", `[{"status":"pending"}]`, false, false},
		"running":          {"in_progress", "failure", `[{"status":"pending"}]`, false, false},
		"queued job":       {"queued", "failure", `[{"status":"pending"}]`, false, false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var jobs gh.Jobs
			require.NoError(t, json.Unmarshal(fmt.Appendf(nil,
				`{"jobs":[null,{"status":%q,"conclusion":%q,"html_url":"https://ghe.example/job/1","steps":%s}]}`,
				test.status, test.conclusion, test.steps,
			), &jobs))

			res := (&ListDisconnectedJobsReq{Owner: "o", Repo: "r", RunID: 1}).inspectJobs(t.Context(), jobs.Jobs)
			require.Equal(t, test.suspected, res.SuspectedDisconnect)
			if test.failed {
				require.Len(t, res.FailedJobs, 1)
				require.Equal(t, test.suspected, res.FailedJobs[0].SuspectedDisconnect)
				require.Equal(t, "https://ghe.example/job/1", res.FailedJobs[0].URL)
			} else {
				require.Empty(t, res.FailedJobs)
			}
		})
	}
}

// TestListDisconnectedJobs_Attempt checks that we inspect the latest attempt of a
// run unless we're given one, and that we tell the caller which attempt we
// inspected. If someone has already rerun the run, the disconnect will only show
// up in an earlier attempt. We need a mock because Github serves each attempt's
// jobs from a different route.
func TestListDisconnectedJobs_Attempt(t *testing.T) {
	t.Parallel()

	client, mux, teardown := githubtest.SetupTestClient(t)
	defer teardown()

	mux.HandleFunc("GET /repos/o/r/actions/runs/1/jobs", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"jobs":[{"run_attempt":2,"status":"completed","conclusion":"success","steps":[{"status":"completed"}]}]}`)
		assert.NoError(t, err)
	})
	mux.HandleFunc("GET /repos/o/r/actions/runs/1/attempts/1/jobs", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"jobs":[{"run_attempt":1,"status":"completed","conclusion":"failure","steps":[{"status":"in_progress"}]}]}`)
		assert.NoError(t, err)
	})

	latest, err := (&ListDisconnectedJobsReq{Owner: "o", Repo: "r", RunID: 1}).Run(t.Context(), client)
	require.NoError(t, err)
	require.Equal(t, 2, latest.RunAttempt)
	require.False(t, latest.SuspectedDisconnect)
	require.Empty(t, latest.FailedJobs)

	first, err := (&ListDisconnectedJobsReq{Owner: "o", Repo: "r", RunID: 1, Attempt: 1}).Run(t.Context(), client)
	require.NoError(t, err)
	require.Equal(t, 1, first.RunAttempt)
	require.True(t, first.SuspectedDisconnect)
	require.Len(t, first.FailedJobs, 1)
}

// TestListDisconnectedJobs_PageFailure checks that we return an error instead of
// partial results if Github fails to give us a later page of jobs. Otherwise we
// might miss a disconnect.
func TestListDisconnectedJobs_PageFailure(t *testing.T) {
	t.Parallel()

	for name, attempt := range map[string]int{"latest attempt": 0, "earlier attempt": 1} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, mux, teardown := githubtest.SetupTestClient(t)
			defer teardown()

			var gets atomic.Int32
			handler := func(w http.ResponseWriter, r *http.Request) {
				if gets.Add(1) > 1 {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Link", fmt.Sprintf(`<http://%s%s?page=2>; rel="next"`, r.Host, r.URL.Path))
				_, err := io.WriteString(w, `{"jobs":[]}`)
				assert.NoError(t, err)
			}
			mux.HandleFunc("GET /repos/o/r/actions/runs/1/jobs", handler)
			mux.HandleFunc("GET /repos/o/r/actions/runs/1/attempts/1/jobs", handler)

			req := &ListDisconnectedJobsReq{Owner: "o", Repo: "r", RunID: 1, Attempt: attempt}
			res, err := req.Run(t.Context(), client)
			require.ErrorContains(t, err, "inspecting o/r run 1")
			require.Nil(t, res)
			require.EqualValues(t, 2, gets.Load())
		})
	}
}

// TestListDisconnectedJobs_Validation checks that invalid requests fail before
// we call the API.
func TestListDisconnectedJobs_Validation(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req    *ListDisconnectedJobsReq
		client *gh.Client
		expect string
	}{
		"nil request": {nil, gh.NewClient(nil), "uninitialized request"},
		"nil client":  {&ListDisconnectedJobsReq{}, nil, "client is required"},
		"owner":       {&ListDisconnectedJobsReq{Repo: "r", RunID: 1}, gh.NewClient(nil), "owner is required"},
		"repo":        {&ListDisconnectedJobsReq{Owner: "o", RunID: 1}, gh.NewClient(nil), "repo is required"},
		"zero ID":     {&ListDisconnectedJobsReq{Owner: "o", Repo: "r"}, gh.NewClient(nil), "run ID must be positive"},
		"negative ID": {&ListDisconnectedJobsReq{Owner: "o", Repo: "r", RunID: -1}, gh.NewClient(nil), "run ID must be positive"},
		"negative attempt": {
			&ListDisconnectedJobsReq{Owner: "o", Repo: "r", RunID: 1, Attempt: -1},
			gh.NewClient(nil), "attempt can't be negative",
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

// TestListDisconnectedJobs_Formatting checks that we escape Markdown cell text
// without changing the values we return as JSON or write to $GITHUB_OUTPUT.
func TestListDisconnectedJobs_Formatting(t *testing.T) {
	t.Parallel()

	res := &ListDisconnectedJobsRes{
		Owner: "o", Repo: "r", RunID: 1, RunAttempt: 2, JobsInspected: 2, SuspectedDisconnect: true,
		FailedJobs: []*DisconnectedJobResult{{
			ID: 2, Name: "build|test\n<script>", URL: "https://ghe.example/job/2", SuspectedDisconnect: true,
			UnfinishedSteps: []*DisconnectedStepResult{{Number: 51, Name: "test_*`[x]", Status: "pending"}},
		}},
	}

	text := res.ToTable()
	for _, want := range []string{"Summary:", "Failed Jobs:", "Unfinished Steps:", "build|test", "pending", "51"} {
		require.Contains(t, text, want)
	}

	md := res.ToMarkdown()
	for _, want := range []string{
		"## Summary", "## Failed Jobs", "## Unfinished Steps", `build\|test<br/>&lt;script&gt;`,
		"test&#95;&#42;&#96;&#91;x&#93;", "https://ghe.example/job/2",
	} {
		require.Contains(t, md, want)
	}

	b, err := res.ToJSON()
	require.NoError(t, err)
	var decoded ListDisconnectedJobsRes
	require.NoError(t, json.Unmarshal(b, &decoded))
	require.Equal(t, *res, decoded)

	output, err := res.ToGithubOutput()
	require.NoError(t, err)
	require.Equal(t, b, output)

	none := &ListDisconnectedJobsRes{FailedJobs: []*DisconnectedJobResult{}}
	for _, output := range []string{none.ToTable(), none.ToMarkdown()} {
		require.Contains(t, output, "Summary")
		require.NotContains(t, output, "Failed Jobs")
		require.NotContains(t, output, "Unfinished Steps")
	}

	var absent *ListDisconnectedJobsRes
	require.Empty(t, absent.ToTable())
	require.Empty(t, absent.ToMarkdown())
	_, err = absent.ToJSON()
	require.Error(t, err)
	_, err = absent.ToGithubOutput()
	require.Error(t, err)
}
