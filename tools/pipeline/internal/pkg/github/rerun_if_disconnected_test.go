// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync/atomic"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRerunIfDisconnected_ObservedJob checks a sanitized example of a disconnected
// job. Multiple matching jobs should still result in just one rerun request.
func TestRerunIfDisconnected_ObservedJob(t *testing.T) {
	t.Parallel()

	disconnectedExample, err := os.ReadFile("testfixtures/disconnected_job.json")
	require.NoError(t, err)
	var job gh.WorkflowJob
	require.NoError(t, json.Unmarshal(disconnectedExample, &job))

	var posts atomic.Int32
	var gets atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method + " " + r.URL.Path {
		case "GET /repos/example/project/actions/runs/2001/jobs":
			gets.Add(1)
			assert.Equal(t, "latest", r.URL.Query().Get("filter"))
			assert.Equal(t, "100", r.URL.Query().Get("per_page"))
			assert.NoError(t, json.NewEncoder(w).Encode(
				gh.Jobs{
					Jobs: []*gh.WorkflowJob{
						&job,
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
					},
				}))
		case "POST /repos/example/project/actions/runs/2001/rerun-failed-jobs":
			posts.Add(1)
			w.WriteHeader(http.StatusCreated)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			http.Error(w, "unexpected request", http.StatusNotFound)
		}
	}))

	t.Cleanup(srv.Close)

	client := gh.NewClient(srv.Client())
	client.BaseURL, err = url.Parse(srv.URL + "/")
	require.NoError(t, err)

	req := &github.RerunIfDisconnectedReq{Owner: "example", Repo: "project", RunID: 2001}
	res, err := req.Run(context.Background(), client)
	require.NoError(t, err)
	require.True(t, res.RerunRequested)
	require.EqualValues(t, 1, posts.Load())
	require.EqualValues(t, 1, gets.Load())
	require.Equal(t, 4, res.JobsInspected)
	require.Len(t, res.FailedJobs, 3)
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
	require.False(t, res.FailedJobs[2].SuspectedDisconnect)

	b, err := res.ToJSON()
	require.NoError(t, err)
	require.Contains(t, string(b), `"suspected_disconnect":false,"unfinished_steps":[]`)
}

// TestRerunIfDisconnected_Classification checks that we only suspect a disconnect
// when a completed, failed job still has unfinished steps.
func TestRerunIfDisconnected_Classification(t *testing.T) {
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

			var posts atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					posts.Add(1)
					w.WriteHeader(http.StatusCreated)
					return
				}
				_, err := fmt.Fprintf(w, `{"jobs":[null,{"status":%q,"conclusion":%q,"html_url":"https://ghe.example/job/1","steps":%s}]}`, test.status, test.conclusion, test.steps)
				assert.NoError(t, err)
			}))

			t.Cleanup(srv.Close)

			client := gh.NewClient(srv.Client())
			var err error
			client.BaseURL, err = url.Parse(srv.URL + "/")
			require.NoError(t, err)

			res, err := (&github.RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: 1}).Run(context.Background(), client)
			require.NoError(t, err)
			require.Equal(t, test.suspected, res.RerunRequested)
			if test.suspected {
				require.EqualValues(t, 1, posts.Load())
			} else {
				require.Zero(t, posts.Load())
			}
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

// TestRerunIfDisconnected_PaginationAndErrors checks that we fetch every page
// before requesting a rerun and return API errors without retrying.
func TestRerunIfDisconnected_PaginationAndErrors(t *testing.T) {
	t.Parallel()
	for name, test := range map[string]struct {
		firstStatus   int
		secondStatus  int
		postStatus    int
		matchFirst    bool
		matchSecond   bool
		shouldBeEmpty bool
		expectGetCnt  int32
		expectWantCnt int32
		shouldFail    string
	}{
		"match on later page": {200, 200, 201, false, true, false, 2, 1, ""},
		"later page failure":  {200, 500, 201, true, false, false, 2, 0, "inspecting o/r run 1"},
		"first page failure":  {403, 200, 201, false, false, false, 1, 0, "inspecting o/r run 1"},
		"rerun rejected":      {200, 200, 403, true, false, false, 2, 1, "rerunning failed jobs for o/r run 1"},
		"empty run":           {200, 200, 201, false, false, true, 1, 0, ""},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var gets, posts atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					assert.Equal(t, test.expectGetCnt, gets.Load())
					posts.Add(1)
					w.WriteHeader(test.postStatus)
					return
				}
				gets.Add(1)
				status, match := test.firstStatus, test.matchFirst
				if r.URL.Query().Get("page") == "2" {
					status, match = test.secondStatus, test.matchSecond
				} else if !test.shouldBeEmpty {
					w.Header().Set("Link", fmt.Sprintf(`<http://%s%s?page=2>; rel="next"`, r.Host, r.URL.Path))
				}
				w.WriteHeader(status)
				body := `{"jobs":[]}`
				if match {
					body = `{"jobs":[{"status":"completed","conclusion":"failure","steps":[{"status":"pending"}]}]}`
				}
				_, err := w.Write([]byte(body))
				assert.NoError(t, err)
			}))

			t.Cleanup(srv.Close)

			client := gh.NewClient(srv.Client())
			var err error
			client.BaseURL, err = url.Parse(srv.URL + "/")
			require.NoError(t, err)

			res, err := (&github.RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: 1}).Run(t.Context(), client)
			if test.shouldFail != "" {
				require.ErrorContains(t, err, test.shouldFail)
				require.Nil(t, res)
				var apiErr *gh.ErrorResponse
				require.ErrorAs(t, err, &apiErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectWantCnt > 0, res.RerunRequested)
				if test.shouldBeEmpty {
					b, err := res.ToJSON()
					require.NoError(t, err)
					require.Contains(t, string(b), `"rerun_requested":false,"failed_jobs":[]`)
				}
			}
			require.Equal(t, test.expectGetCnt, gets.Load())
			require.Equal(t, test.expectWantCnt, posts.Load())
		})
	}
}

// TestRerunIfDisconnected_Validation checks that invalid requests fail before we call the API.
func TestRerunIfDisconnected_Validation(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req    *github.RerunIfDisconnectedReq
		client *gh.Client
		expect string
	}{
		"nil request": {nil, gh.NewClient(nil), "uninitialized request"},
		"nil client":  {&github.RerunIfDisconnectedReq{}, nil, "client is required"},
		"owner":       {&github.RerunIfDisconnectedReq{Repo: "r", RunID: 1}, gh.NewClient(nil), "owner is required"},
		"repo":        {&github.RerunIfDisconnectedReq{Owner: "o", RunID: 1}, gh.NewClient(nil), "repo is required"},
		"zero ID":     {&github.RerunIfDisconnectedReq{Owner: "o", Repo: "r"}, gh.NewClient(nil), "run ID must be positive"},
		"negative ID": {&github.RerunIfDisconnectedReq{Owner: "o", Repo: "r", RunID: -1}, gh.NewClient(nil), "run ID must be positive"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res, err := test.req.Run(context.Background(), test.client)
			require.ErrorContains(t, err, test.expect)
			require.Nil(t, res)
		})
	}
}

// TestRerunIfDisconnected_Formatting checks that we escape Markdown cell text
// without changing the values returned as JSON.
func TestRerunIfDisconnected_Formatting(t *testing.T) {
	t.Parallel()

	res := &github.RerunIfDisconnectedRes{
		Owner: "o", Repo: "r", RunID: 1, JobsInspected: 2, RerunRequested: true,
		FailedJobs: []*github.DisconnectedJobResult{{
			ID: 2, Name: "build|test\n<script>", URL: "https://ghe.example/job/2", SuspectedDisconnect: true,
			UnfinishedSteps: []*github.DisconnectedStepResult{{Number: 51, Name: "test_*`[x]", Status: "pending"}},
		}},
	}

	text := res.ToTable()
	for _, want := range []string{"Summary:", "Failed Jobs:", "Unfinished Steps:", "build|test", "pending", "51", "Rerun of failed jobs requested"} {
		require.Contains(t, text, want)
	}

	md := res.ToMarkdown()
	for _, want := range []string{"## Summary", "## Failed Jobs", "## Unfinished Steps", `build\|test<br/>&lt;script&gt;`, "test&#95;&#42;&#96;&#91;x&#93;", "https://ghe.example/job/2"} {
		require.Contains(t, md, want)
	}

	b, err := res.ToJSON()
	require.NoError(t, err)
	var decoded github.RerunIfDisconnectedRes
	require.NoError(t, json.Unmarshal(b, &decoded))
	require.Equal(t, *res, decoded)

	noop := &github.RerunIfDisconnectedRes{FailedJobs: []*github.DisconnectedJobResult{}}
	for _, output := range []string{noop.ToTable(), noop.ToMarkdown()} {
		require.Contains(t, output, "No suspected disconnects found. No rerun requested.")
		require.NotContains(t, output, "Failed Jobs")
		require.NotContains(t, output, "Unfinished Steps")
	}

	var absent *github.RerunIfDisconnectedRes
	require.Empty(t, absent.ToTable())
	require.Empty(t, absent.ToMarkdown())
	_, err = absent.ToJSON()
	require.Error(t, err)
}
