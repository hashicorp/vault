// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	gh "github.com/google/go-github/v68/github"
)

// PerPageMax is the maximum number of entities to request for enpoints that
// support pagination. 100 is usually the limit and we use it everywhere.
// We always request the maximum number of entities so that we we use the fewest
// possible API requests. There is a per-hour token limit after all.
const PerPageMax = 100

// ListWorkflowRunsReq is a request to list workflows runs. The fields represent
// various criteria we can use to filter.
type ListWorkflowRunsReq struct {
	Actor        string
	Branch       string
	CheckSuiteID int64
	Compact      bool
	DateQuery    string
	Event        string
	IncludePRs   bool
	Owner        string
	Repo         string
	Sha          string
	Status       string
	WorkflowName string
}

// ListWorkflowRunsRes is a list workflows response.
type ListWorkflowRunsRes struct {
	Workflow *gh.Workflow   `json:"workflow,omitempty"`
	Runs     []*WorkflowRun `json:"runs,omitempty"`
}

// WorkflowRun represents a Github actions workflow run. We include the raw
// Github API run response, the workflows jobs, and the associated check suite.
type WorkflowRun struct {
	Run       *gh.WorkflowRun `json:"run,omitempty"`
	Jobs      []*WorkflowJob  `json:"jobs,omitempty"`
	CheckRuns []*CheckRun     `json:"check_runs,omitempty"`
	summary   string
}

// CheckRun represents the run of a check suite for a workflow. We include the
// check suite annotations.
type CheckRun struct {
	Run         *gh.CheckRun             `json:"run,omitempty"`
	Annotations []*gh.CheckRunAnnotation `json:"annotations,omitempty"`
}

// WorkflowJob represents a singular job of a workflow. We include the raw
// Job response from Github along with log entries for any failed steps.
type WorkflowJob struct {
	Job        *gh.WorkflowJob `json:"job,omitempty"`
	LogEntries []*LogEntry     `json:"log_entries,omitempty"`
}

// Run runs the request to gather all instances of the workflow that match
// our filter criteria.
func (r *ListWorkflowRunsReq) Run(ctx context.Context, client *gh.Client) (*ListWorkflowRunsRes, error) {
	var err error
	res := &ListWorkflowRunsRes{}

	if err = r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	res.Workflow, err = r.getWorkflow(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("getting workflow: %w", err)
	}

	res.Runs, err = r.getWorkflowRuns(ctx, client, res.Workflow.GetID())
	if err != nil {
		return nil, fmt.Errorf("getting workflow runs: %w", err)
	}

	if len(res.Runs) < 1 {
		return nil, fmt.Errorf("fetching workflow runs: no workflow runs match the given filters")
	}

	err = r.getWorkflowCheckRuns(ctx, client, res.Runs)
	if err != nil {
		return nil, fmt.Errorf("fetching workflow check runs: %w", err)
	}

	err = r.getWorkflowCheckRunAnnotations(ctx, client, res.Runs)
	if err != nil {
		return nil, fmt.Errorf("fetching workflow check run annotations: %w", err)
	}

	err = r.getWorkflowJobs(ctx, client, res.Runs)
	if err != nil {
		return nil, fmt.Errorf("fetching workflow run jobs: %w", err)
	}

	// Logs have to be downloaded and parsed to get the relevant bits. Such
	// an expensive operation is limited to only failures.
	err = r.getUnsuccessfulWorkflowJobsLogs(ctx, client, res.Runs)
	if err != nil {
		return nil, fmt.Errorf("fetching failed workflow run job logs: %w", err)
	}

	return res, r.summarizeWorkflowRuns(res.Runs)
}

// validate ensures that we've been given the minimum filter arguments necessary to complete a
// request. It is always recommended that additional fitlers be given to reduce the response size
// and not exhaust API limits.
func (r *ListWorkflowRunsReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.DateQuery == "" {
		return errors.New("no date range query has been provided")
	}

	if r.WorkflowName == "" {
		return errors.New("no github actions workflow name has been provided")
	}

	return nil
}

// getWorkflow attempts to locate the workflow associated with our workflow name.
func (r *ListWorkflowRunsReq) getWorkflow(ctx context.Context, client *gh.Client) (*gh.Workflow, error) {
	opts := &gh.ListOptions{PerPage: PerPageMax}
	for {
		wfs, res, err := client.Actions.ListWorkflows(ctx, r.Owner, r.Repo, opts)
		if err != nil {
			return nil, err
		}

		for _, wf := range wfs.Workflows {
			if wf.GetName() == r.WorkflowName {
				return wf, nil
			}
		}

		if res.NextPage == 0 {
			return nil, fmt.Errorf("no workflow matching %s could be found", r.WorkflowName)
		}

		opts.Page = res.NextPage
	}
}

// getWorkflowRuns gets teh workflow runs associated with a workflow ID.
func (r *ListWorkflowRunsReq) getWorkflowRuns(ctx context.Context, client *gh.Client, id int64) ([]*WorkflowRun, error) {
	var runs []*WorkflowRun
	opts := &gh.ListWorkflowRunsOptions{
		Actor:               r.Actor,
		Branch:              r.Branch,
		CheckSuiteID:        r.CheckSuiteID,
		Created:             r.DateQuery,
		ExcludePullRequests: !r.IncludePRs,
		Event:               r.Event,
		HeadSHA:             r.Sha,
		ListOptions:         gh.ListOptions{PerPage: PerPageMax},
		Status:              r.Status,
	}
	if r.CheckSuiteID > 0 {
		opts.CheckSuiteID = r.CheckSuiteID
	}

	for {
		wfrs, res, err := client.Actions.ListWorkflowRunsByID(ctx, r.Owner, r.Repo, id, opts)
		if err != nil {
			return nil, err
		}

		for _, r := range wfrs.WorkflowRuns {
			runs = append(runs, &WorkflowRun{Run: r})
		}

		if res.NextPage == 0 {
			return runs, nil
		}

		opts.ListOptions.Page = res.NextPage
	}
}

// getWorkflowCheckRuns gets the check suite runs associated with the workflow runs.
func (r *ListWorkflowRunsReq) getWorkflowCheckRuns(ctx context.Context, client *gh.Client, wfrs []*WorkflowRun) error {
	filter := "latest"
	opts := &gh.ListCheckRunsOptions{
		Filter:      &filter, // "all" for all attemps
		ListOptions: gh.ListOptions{PerPage: PerPageMax},
	}

	if len(wfrs) < 1 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(wfrs))
	errC, resC, cancel := r.startErrorCollector(ctx)

	for _, wfr := range wfrs {
		go func() {
			defer wg.Done()

			for {
				chrs, res, err := client.Checks.ListCheckRunsCheckSuite(ctx, r.Owner, r.Repo, *wfr.Run.CheckSuiteID, opts)
				if err != nil {
					errC <- err
					return
				}

				if r.Status == "" || !r.Compact {
					for _, cr := range chrs.CheckRuns {
						wfr.CheckRuns = append(wfr.CheckRuns, &CheckRun{Run: cr})
					}
				} else {
					for _, cr := range chrs.CheckRuns {
						if cr.GetConclusion() == r.Status || cr.GetStatus() == r.Status {
							wfr.CheckRuns = append(wfr.CheckRuns, &CheckRun{Run: cr})
						}
					}
				}

				if res.NextPage == 0 {
					break
				}

				opts.ListOptions.Page = res.NextPage
			}
		}()
	}

	wg.Wait()

	return r.stopErrorCollector(errC, resC, cancel)
}

// getWorkflowCheckRunAnnotations gets the check suite annotations associated with the check suites
// that are associated with the runs.
func (r *ListWorkflowRunsReq) getWorkflowCheckRunAnnotations(ctx context.Context, client *gh.Client, wfrs []*WorkflowRun) error {
	opts := &gh.ListOptions{PerPage: PerPageMax}

	if len(wfrs) < 1 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(wfrs))
	errC, resC, cancel := r.startErrorCollector(ctx)

	for _, wfr := range wfrs {
		go func() {
			defer wg.Done()

			for _, cr := range wfr.CheckRuns {
				for {
					ans, res, err := client.Checks.ListCheckRunAnnotations(ctx, r.Owner, r.Repo, cr.Run.GetID(), opts)
					if err != nil {
						errC <- err
						return
					}
					cr.Annotations = append(cr.Annotations, ans...)

					if res.NextPage == 0 {
						break
					}

					opts.Page = res.NextPage
				}
			}
		}()
	}

	wg.Wait()

	return r.stopErrorCollector(errC, resC, cancel)
}

// getWorkflowJobs gets the jobs associated with the workflow runs.
func (r *ListWorkflowRunsReq) getWorkflowJobs(ctx context.Context, client *gh.Client, wfrs []*WorkflowRun) error {
	opts := &gh.ListWorkflowJobsOptions{
		Filter:      "latest", // "all" to include all attempts
		ListOptions: gh.ListOptions{PerPage: PerPageMax},
	}

	if len(wfrs) < 1 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(wfrs))
	errC, resC, cancel := r.startErrorCollector(ctx)

	for _, run := range wfrs {
		go func() {
			defer wg.Done()

			for {
				jobs, res, err := client.Actions.ListWorkflowJobs(ctx, r.Owner, r.Repo, *run.Run.ID, opts)
				if err != nil {
					errC <- err
					return
				}

				if r.Status == "" || !r.Compact {
					for _, job := range jobs.Jobs {
						run.Jobs = append(run.Jobs, &WorkflowJob{Job: job})
					}
				} else {
					for _, job := range jobs.Jobs {
						if job.GetConclusion() == r.Status || job.GetStatus() == r.Status {
							run.Jobs = append(run.Jobs, &WorkflowJob{Job: job})
						}
					}
				}

				if res.NextPage == 0 {
					break
				}

				opts.ListOptions.Page = res.NextPage
			}
		}()
	}

	wg.Wait()

	return r.stopErrorCollector(errC, resC, cancel)
}

// getUnsuccessfulWorkflowJobsLogs downloads the job log and parses out the
// out failed entries for any unsuccesful jobs.
func (r *ListWorkflowRunsReq) getUnsuccessfulWorkflowJobsLogs(ctx context.Context, client *gh.Client, wfrs []*WorkflowRun) error {
	if len(wfrs) < 1 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(wfrs))
	errC, resC, cancel := r.startErrorCollector(ctx)

	for _, run := range wfrs {
		go func() {
			defer wg.Done()

			for _, job := range run.Jobs {
				if job.Job == nil {
					continue
				}

				if job.Job.GetStatus() == "completed" && job.Job.GetConclusion() == "successful" {
					continue
				}

				url, _, err := client.Actions.GetWorkflowJobLogs(ctx, r.Owner, r.Repo, *job.Job.ID, 3) // last is max redirects
				if err != nil {
					errC <- err
					return
				}

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
				if err != nil {
					errC <- err
					return
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					errC <- err
					return
				}

				defer res.Body.Close()
				scanner := NewLogScaner(
					WithLogScannerTruncate(), // truncate the body
					// Our max size here is a magic number but seems to be a nice sweet
					// spot where most of failure diagnostic will show up but you don't
					// get walls of text.
					WithLogScannerMaxSize(4000),
					WithLogScannerOnlyUnsuccessful(), // only keep unsuccessful step logs
				)

				job.LogEntries, err = scanner.Scan(res.Body)
				if err != nil {
					errC <- err
					return
				}
			}
		}()
	}

	wg.Wait()

	return r.stopErrorCollector(errC, resC, cancel)
}

// startErrorCollector starts a helper go routine that listens and aggregrates
// errors passed to it via it's returned channel. When the caller has completed
// all work that may write errors it can call the stopErrorCollector() method
// with the channel and cancel func to safely shut down the go routine and
// receive the aggregrate error.
func (r *ListWorkflowRunsReq) startErrorCollector(ctx context.Context) (chan error, chan error, context.CancelFunc) {
	errC := make(chan error)
	resC := make(chan error)
	errCtx, cancelCollector := context.WithCancel(ctx)

	go func() {
		var err error

	LOOP:
		for {
			select {
			case err1 := <-errC:
				err = errors.Join(err, err1)
				continue
			default:
			}

			select {
			case err1 := <-errC:
				err = errors.Join(err, err1)
				continue
			case <-errCtx.Done():
				// Don't bubble up the close signal, only bubble up outer context errors
				if err1 := ctx.Err(); err1 != nil {
					err = errors.Join(err, err1)
				}
				break LOOP
			}
		}

	DRAIN:
		for {
			select {
			case err1 := <-errC:
				err = errors.Join(err, err1)
				continue
			default:
				break DRAIN
			}
		}

		resC <- err
	}()

	return errC, resC, cancelCollector
}

// stopErrorCollector stops the error collector and returns the aggregrated error. The given channel
// is closed. The caller must be sure that all work has concluded prior to calling.
func (r *ListWorkflowRunsReq) stopErrorCollector(errC chan error, resC chan error, cancel context.CancelFunc) error {
	cancel()
	err := <-resC
	close(errC)
	close(resC)

	return err
}

// summarizeWorkflowRuns creates a human readable summary for all workflow runs.
func (r *ListWorkflowRunsReq) summarizeWorkflowRuns(wrfs []*WorkflowRun) error {
	if len(wrfs) < 1 {
		return nil
	}

	var err error
	for _, run := range wrfs {
		_, err = run.Summary()
		if err != nil {
			return err
		}
	}

	return nil
}

// Summary returns the human readable summary of the workflow run.
func (r *WorkflowRun) Summary() (string, error) {
	if r == nil {
		return "", errors.New("uninitialized workflow run")
	}

	if r.summary != "" {
		return r.summary, nil
	}

	var err error
	r.summary, err = summarizeWorkflowRun(r)
	return r.summary, err
}

// UnsuccessfulSteps returns any unsuccesful steps in the workflow job.
func (j *WorkflowJob) UnsuccessfulSteps() []*gh.TaskStep {
	if j == nil || j.Job == nil || len(j.Job.Steps) < 1 {
		return nil
	}

	res := []*gh.TaskStep{}
	for _, step := range j.Job.Steps {
		if step.GetStatus() == "completed" && (step.GetConclusion() == "success" || step.GetConclusion() == "skipped") {
			continue
		}

		res = append(res, step)
	}

	return res
}

// UnsuccessfulSteps returns the names of any unsuccesful steps in the workflow job.
func (j *WorkflowJob) UnsuccessfulStepNames() []string {
	steps := j.UnsuccessfulSteps()
	if len(steps) < 1 {
		return nil
	}

	res := []string{}
	for _, step := range steps {
		res = append(res, step.GetName())
	}

	return res
}
