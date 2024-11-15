// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"time"
)

// RunTaskRequest is the payload object that TFC/E sends to the Run Task's URL.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/run-tasks/run-tasks-integration#common-properties
type RunTaskRequest struct {
	AccessToken                     string                      `json:"access_token"`
	Capabilitites                   RunTaskRequestCapabilitites `json:"capabilitites,omitempty"`
	ConfigurationVersionDownloadURL string                      `json:"configuration_version_download_url,omitempty"`
	ConfigurationVersionID          string                      `json:"configuration_version_id,omitempty"`
	IsSpeculative                   bool                        `json:"is_speculative"`
	OrganizationName                string                      `json:"organization_name"`
	PayloadVersion                  int                         `json:"payload_version"`
	PlanJSONAPIURL                  string                      `json:"plan_json_api_url,omitempty"` // Specific to post_plan, pre_apply or post_apply stage
	RunAppURL                       string                      `json:"run_app_url"`
	RunCreatedAt                    time.Time                   `json:"run_created_at"`
	RunCreatedBy                    string                      `json:"run_created_by"`
	RunID                           string                      `json:"run_id"`
	RunMessage                      string                      `json:"run_message"`
	Stage                           string                      `json:"stage"`
	TaskResultCallbackURL           string                      `json:"task_result_callback_url"`
	TaskResultEnforcementLevel      string                      `json:"task_result_enforcement_level"`
	TaskResultID                    string                      `json:"task_result_id"`
	VcsBranch                       string                      `json:"vcs_branch,omitempty"`
	VcsCommitURL                    string                      `json:"vcs_commit_url,omitempty"`
	VcsPullRequestURL               string                      `json:"vcs_pull_request_url,omitempty"`
	VcsRepoURL                      string                      `json:"vcs_repo_url,omitempty"`
	WorkspaceAppURL                 string                      `json:"workspace_app_url"`
	WorkspaceID                     string                      `json:"workspace_id"`
	WorkspaceName                   string                      `json:"workspace_name"`
	WorkspaceWorkingDirectory       string                      `json:"workspace_working_directory,omitempty"`
}

// RunTaskRequestCapabilitites defines the capabilities that the caller supports.
type RunTaskRequestCapabilitites struct {
	Outcomes bool `json:"outcomes"`
}
