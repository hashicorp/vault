// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ChecksService provides access to the Checks API in the
// GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/checks/
type ChecksService service

// CheckRun represents a GitHub check run on a repository associated with a GitHub app.
type CheckRun struct {
	ID           *int64          `json:"id,omitempty"`
	HeadSHA      *string         `json:"head_sha,omitempty"`
	ExternalID   *string         `json:"external_id,omitempty"`
	URL          *string         `json:"url,omitempty"`
	HTMLURL      *string         `json:"html_url,omitempty"`
	Status       *string         `json:"status,omitempty"`
	Conclusion   *string         `json:"conclusion,omitempty"`
	StartedAt    *Timestamp      `json:"started_at,omitempty"`
	CompletedAt  *Timestamp      `json:"completed_at,omitempty"`
	Output       *CheckRunOutput `json:"output,omitempty"`
	Name         *string         `json:"name,omitempty"`
	CheckSuite   *CheckSuite     `json:"check_suite,omitempty"`
	App          *App            `json:"app,omitempty"`
	PullRequests []*PullRequest  `json:"pull_requests,omitempty"`
}

// CheckRunOutput represents the output of a CheckRun.
type CheckRunOutput struct {
	Title            *string               `json:"title,omitempty"`
	Summary          *string               `json:"summary,omitempty"`
	Text             *string               `json:"text,omitempty"`
	AnnotationsCount *int                  `json:"annotations_count,omitempty"`
	AnnotationsURL   *string               `json:"annotations_url,omitempty"`
	Annotations      []*CheckRunAnnotation `json:"annotations,omitempty"`
	Images           []*CheckRunImage      `json:"images,omitempty"`
}

// CheckRunAnnotation represents an annotation object for a CheckRun output.
type CheckRunAnnotation struct {
	FileName     *string `json:"filename,omitempty"`
	BlobHRef     *string `json:"blob_href,omitempty"`
	StartLine    *int    `json:"start_line,omitempty"`
	EndLine      *int    `json:"end_line,omitempty"`
	WarningLevel *string `json:"warning_level,omitempty"`
	Message      *string `json:"message,omitempty"`
	Title        *string `json:"title,omitempty"`
	RawDetails   *string `json:"raw_details,omitempty"`
}

// CheckRunImage represents an image object for a CheckRun output.
type CheckRunImage struct {
	Alt      *string `json:"alt,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
	Caption  *string `json:"caption,omitempty"`
}

// CheckSuite represents a suite of check runs.
type CheckSuite struct {
	ID           *int64         `json:"id,omitempty"`
	HeadBranch   *string        `json:"head_branch,omitempty"`
	HeadSHA      *string        `json:"head_sha,omitempty"`
	URL          *string        `json:"url,omitempty"`
	BeforeSHA    *string        `json:"before,omitempty"`
	AfterSHA     *string        `json:"after,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Conclusion   *string        `json:"conclusion,omitempty"`
	App          *App           `json:"app,omitempty"`
	Repository   *Repository    `json:"repository,omitempty"`
	PullRequests []*PullRequest `json:"pull_requests,omitempty"`
}

func (c CheckRun) String() string {
	return Stringify(c)
}

func (c CheckSuite) String() string {
	return Stringify(c)
}

// GetCheckRun gets a check-run for a repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#get-a-single-check-run
func (s *ChecksService) GetCheckRun(ctx context.Context, owner, repo string, checkRunID int64) (*CheckRun, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs/%v", owner, repo, checkRunID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	checkRun := new(CheckRun)
	resp, err := s.client.Do(ctx, req, checkRun)
	if err != nil {
		return nil, resp, err
	}

	return checkRun, resp, nil
}

// GetCheckSuite gets a single check suite.
//
// GitHub API docs: https://developer.github.com/v3/checks/suites/#get-a-single-check-suite
func (s *ChecksService) GetCheckSuite(ctx context.Context, owner, repo string, checkSuiteID int64) (*CheckSuite, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-suites/%v", owner, repo, checkSuiteID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	checkSuite := new(CheckSuite)
	resp, err := s.client.Do(ctx, req, checkSuite)
	if err != nil {
		return nil, resp, err
	}

	return checkSuite, resp, nil
}

// CreateCheckRunOptions sets up parameters needed to create a CheckRun.
type CreateCheckRunOptions struct {
	Name        string          `json:"name"`                   // The name of the check (e.g., "code-coverage"). (Required.)
	HeadBranch  string          `json:"head_branch"`            // The name of the branch to perform a check against. (Required.)
	HeadSHA     string          `json:"head_sha"`               // The SHA of the commit. (Required.)
	DetailsURL  *string         `json:"details_url,omitempty"`  // The URL of the integrator's site that has the full details of the check. (Optional.)
	ExternalID  *string         `json:"external_id,omitempty"`  // A reference for the run on the integrator's system. (Optional.)
	Status      *string         `json:"status,omitempty"`       // The current status. Can be one of "queued", "in_progress", or "completed". Default: "queued". (Optional.)
	Conclusion  *string         `json:"conclusion,omitempty"`   // Can be one of "success", "failure", "neutral", "cancelled", "timed_out", or "action_required". (Optional. Required if you provide a status of "completed".)
	StartedAt   *Timestamp      `json:"started_at,omitempty"`   // The time that the check run began. (Optional.)
	CompletedAt *Timestamp      `json:"completed_at,omitempty"` // The time the check completed. (Optional. Required if you provide conclusion.)
	Output      *CheckRunOutput `json:"output,omitempty"`       // Provide descriptive details about the run. (Optional)
}

// CreateCheckRun creates a check run for repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#create-a-check-run
func (s *ChecksService) CreateCheckRun(ctx context.Context, owner, repo string, opt CreateCheckRunOptions) (*CheckRun, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs", owner, repo)
	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	checkRun := new(CheckRun)
	resp, err := s.client.Do(ctx, req, checkRun)
	if err != nil {
		return nil, resp, err
	}

	return checkRun, resp, nil
}

// UpdateCheckRunOptions sets up parameters needed to update a CheckRun.
type UpdateCheckRunOptions struct {
	Name        string          `json:"name"`                   // The name of the check (e.g., "code-coverage"). (Required.)
	HeadBranch  *string         `json:"head_branch,omitempty"`  // The name of the branch to perform a check against. (Optional.)
	HeadSHA     *string         `json:"head_sha,omitempty"`     // The SHA of the commit. (Optional.)
	DetailsURL  *string         `json:"details_url,omitempty"`  // The URL of the integrator's site that has the full details of the check. (Optional.)
	ExternalID  *string         `json:"external_id,omitempty"`  // A reference for the run on the integrator's system. (Optional.)
	Status      *string         `json:"status,omitempty"`       // The current status. Can be one of "queued", "in_progress", or "completed". Default: "queued". (Optional.)
	Conclusion  *string         `json:"conclusion,omitempty"`   // Can be one of "success", "failure", "neutral", "cancelled", "timed_out", or "action_required". (Optional. Required if you provide a status of "completed".)
	CompletedAt *Timestamp      `json:"completed_at,omitempty"` // The time the check completed. (Optional. Required if you provide conclusion.)
	Output      *CheckRunOutput `json:"output,omitempty"`       // Provide descriptive details about the run. (Optional)
}

// UpdateCheckRun updates a check run for a specific commit in a repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#update-a-check-run
func (s *ChecksService) UpdateCheckRun(ctx context.Context, owner, repo string, checkRunID int64, opt UpdateCheckRunOptions) (*CheckRun, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs/%v", owner, repo, checkRunID)
	req, err := s.client.NewRequest("PATCH", u, opt)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	checkRun := new(CheckRun)
	resp, err := s.client.Do(ctx, req, checkRun)
	if err != nil {
		return nil, resp, err
	}

	return checkRun, resp, nil
}

// ListCheckRunAnnotations lists the annotations for a check run.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#list-annotations-for-a-check-run
func (s *ChecksService) ListCheckRunAnnotations(ctx context.Context, owner, repo string, checkRunID int64, opt *ListOptions) ([]*CheckRunAnnotation, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs/%v/annotations", owner, repo, checkRunID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	var checkRunAnnotations []*CheckRunAnnotation
	resp, err := s.client.Do(ctx, req, &checkRunAnnotations)
	if err != nil {
		return nil, resp, err
	}

	return checkRunAnnotations, resp, nil
}

// ListCheckRunsOptions represents parameters to list check runs.
type ListCheckRunsOptions struct {
	CheckName *string `url:"check_name,omitempty"` // Returns check runs with the specified name.
	Status    *string `url:"status,omitempty"`     // Returns check runs with the specified status. Can be one of "queued", "in_progress", or "completed".
	Filter    *string `url:"filter,omitempty"`     // Filters check runs by their completed_at timestamp. Can be one of "latest" (returning the most recent check runs) or "all". Default: "latest"

	ListOptions
}

// ListCheckRunsResults represents the result of a check run list.
type ListCheckRunsResults struct {
	Total     *int        `json:"total_count,omitempty"`
	CheckRuns []*CheckRun `json:"check_runs,omitempty"`
}

// ListCheckRunsForRef lists check runs for a specific ref.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#list-check-runs-for-a-specific-ref
func (s *ChecksService) ListCheckRunsForRef(ctx context.Context, owner, repo, ref string, opt *ListCheckRunsOptions) (*ListCheckRunsResults, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/commits/%v/check-runs", owner, repo, ref)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	var checkRunResults *ListCheckRunsResults
	resp, err := s.client.Do(ctx, req, &checkRunResults)
	if err != nil {
		return nil, resp, err
	}

	return checkRunResults, resp, nil
}

// ListCheckRunsCheckSuite lists check runs for a check suite.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#list-check-runs-in-a-check-suite
func (s *ChecksService) ListCheckRunsCheckSuite(ctx context.Context, owner, repo string, checkSuiteID int64, opt *ListCheckRunsOptions) (*ListCheckRunsResults, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-suites/%v/check-runs", owner, repo, checkSuiteID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	var checkRunResults *ListCheckRunsResults
	resp, err := s.client.Do(ctx, req, &checkRunResults)
	if err != nil {
		return nil, resp, err
	}

	return checkRunResults, resp, nil
}

// ListCheckSuiteOptions represents parameters to list check suites.
type ListCheckSuiteOptions struct {
	CheckName *string `url:"check_name,omitempty"` // Filters checks suites by the name of the check run.
	AppID     *int    `url:"app_id,omitempty"`     // Filters check suites by GitHub App id.

	ListOptions
}

// ListCheckSuiteResults represents the result of a check run list.
type ListCheckSuiteResults struct {
	Total       *int          `json:"total_count,omitempty"`
	CheckSuites []*CheckSuite `json:"check_suites,omitempty"`
}

// ListCheckSuitesForRef lists check suite for a specific ref.
//
// GitHub API docs: https://developer.github.com/v3/checks/suites/#list-check-suites-for-a-specific-ref
func (s *ChecksService) ListCheckSuitesForRef(ctx context.Context, owner, repo, ref string, opt *ListCheckSuiteOptions) (*ListCheckSuiteResults, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/commits/%v/check-suites", owner, repo, ref)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	var checkSuiteResults *ListCheckSuiteResults
	resp, err := s.client.Do(ctx, req, &checkSuiteResults)
	if err != nil {
		return nil, resp, err
	}

	return checkSuiteResults, resp, nil
}

// AutoTriggerCheck enables or disables automatic creation of CheckSuite events upon pushes to the repository.
type AutoTriggerCheck struct {
	AppID   *int64 `json:"app_id,omitempty"`  // The id of the GitHub App. (Required.)
	Setting *bool  `json:"setting,omitempty"` // Set to "true" to enable automatic creation of CheckSuite events upon pushes to the repository, or "false" to disable them. Default: "true" (Required.)
}

// CheckSuitePreferenceOptions set options for check suite preferences for a repository.
type CheckSuitePreferenceOptions struct {
	PreferenceList *PreferenceList `json:"auto_trigger_checks,omitempty"` // A list of auto trigger checks that can be set for a check suite in a repository.
}

// CheckSuitePreferenceResults represents the results of the preference set operation.
type CheckSuitePreferenceResults struct {
	Preferences *PreferenceList `json:"preferences,omitempty"`
	Repository  *Repository     `json:"repository,omitempty"`
}

// PreferenceList represents a list of auto trigger checks for  repository
type PreferenceList struct {
	AutoTriggerChecks []*AutoTriggerCheck `json:"auto_trigger_checks,omitempty"` // A slice of auto trigger checks that can be set for a check suite in a repository.
}

// SetCheckSuitePreferences changes the default automatic flow when creating check suites.
//
// GitHub API docs: https://developer.github.com/v3/checks/suites/#set-preferences-for-check-suites-on-a-repository
func (s *ChecksService) SetCheckSuitePreferences(ctx context.Context, owner, repo string, opt CheckSuitePreferenceOptions) (*CheckSuitePreferenceResults, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-suites/preferences", owner, repo)
	req, err := s.client.NewRequest("PATCH", u, opt)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	var checkSuitePrefResults *CheckSuitePreferenceResults
	resp, err := s.client.Do(ctx, req, &checkSuitePrefResults)
	if err != nil {
		return nil, resp, err
	}

	return checkSuitePrefResults, resp, nil
}

// CreateCheckSuiteOptions sets up parameters to manually create a check suites
type CreateCheckSuiteOptions struct {
	HeadSHA    string  `json:"head_sha"`              // The sha of the head commit. (Required.)
	HeadBranch *string `json:"head_branch,omitempty"` // The name of the head branch where the code changes are implemented.
}

// CreateCheckSuite manually creates a check suite for a repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/suites/#create-a-check-suite
func (s *ChecksService) CreateCheckSuite(ctx context.Context, owner, repo string, opt CreateCheckSuiteOptions) (*CheckSuite, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-suites", owner, repo)
	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	checkSuite := new(CheckSuite)
	resp, err := s.client.Do(ctx, req, checkSuite)
	if err != nil {
		return nil, resp, err
	}

	return checkSuite, resp, nil
}

// RequestCheckSuiteOptions sets up the parameters for a request check suite endpoint.
type RequestCheckSuiteOptions struct {
	HeadSHA string `json:"head_sha"` // The sha of the head commit. (Required.)
}

// RequestCheckSuite triggers GitHub to create a new check suite, without pushing new code to a repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/suites/#request-check-suites
func (s *ChecksService) RequestCheckSuite(ctx context.Context, owner, repo string, opt RequestCheckSuiteOptions) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-suite-requests", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
