// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Compile-time proof of interface implementation.
var _ AdminRuns = (*adminRuns)(nil)

// AdminRuns describes all the admin run related methods that the Terraform
// Enterprise  API supports.
// It contains endpoints to help site administrators manage their runs.
//
// TFE API docs: https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/runs
type AdminRuns interface {
	// List all the runs of the given installation.
	List(ctx context.Context, options *AdminRunsListOptions) (*AdminRunsList, error)

	// Force-cancel a run by its ID.
	ForceCancel(ctx context.Context, runID string, options AdminRunForceCancelOptions) error
}

// AdminRun represents AdminRuns interface.
type AdminRun struct {
	ID               string               `jsonapi:"primary,runs"`
	CreatedAt        time.Time            `jsonapi:"attr,created-at,iso8601"`
	HasChanges       bool                 `jsonapi:"attr,has-changes"`
	Status           RunStatus            `jsonapi:"attr,status"`
	StatusTimestamps *RunStatusTimestamps `jsonapi:"attr,status-timestamps"`

	// Relations
	Workspace    *AdminWorkspace    `jsonapi:"relation,workspace"`
	Organization *AdminOrganization `jsonapi:"relation,workspace.organization"`
}

// AdminRunsList represents a list of runs.
type AdminRunsList struct {
	*Pagination
	Items []*AdminRun
}

// AdminRunIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/runs#available-related-resources
type AdminRunIncludeOpt string

const (
	AdminRunWorkspace          AdminRunIncludeOpt = "workspace"
	AdminRunWorkspaceOrg       AdminRunIncludeOpt = "workspace.organization"
	AdminRunWorkspaceOrgOwners AdminRunIncludeOpt = "workspace.organization.owners"
)

// AdminRunsListOptions represents the options for listing runs.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/runs#query-parameters
type AdminRunsListOptions struct {
	ListOptions

	RunStatus     string `url:"filter[status],omitempty"`
	CreatedBefore string `url:"filter[to],omitempty"`
	CreatedAfter  string `url:"filter[from],omitempty"`
	Query         string `url:"q,omitempty"`
	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/runs#available-related-resources
	Include []AdminRunIncludeOpt `url:"include,omitempty"`
}

// adminRuns implements the AdminRuns interface.
type adminRuns struct {
	client *Client
}

// List all the runs of the terraform enterprise installation.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/runs#list-all-runs
func (s *adminRuns) List(ctx context.Context, options *AdminRunsListOptions) (*AdminRunsList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := "admin/runs"
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	rl := &AdminRunsList{}
	err = req.Do(ctx, rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

// AdminRunForceCancelOptions represents the options for force-canceling a run.
type AdminRunForceCancelOptions struct {
	// An optional comment explaining the reason for the force-cancel.
	// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/runs#request-body
	Comment *string `json:"comment,omitempty"`
}

// ForceCancel is used to forcefully cancel a run by its ID.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/runs#force-a-run-into-the-quot-cancelled-quot-state
func (s *adminRuns) ForceCancel(ctx context.Context, runID string, options AdminRunForceCancelOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("admin/runs/%s/actions/force-cancel", url.PathEscape(runID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o *AdminRunsListOptions) valid() error {
	if o == nil { // nothing to validate
		return nil
	}

	if err := validateAdminRunDateRanges(o.CreatedBefore, o.CreatedAfter); err != nil {
		return err
	}

	if err := validateAdminRunFilterParams(o.RunStatus); err != nil {
		return err
	}

	return nil
}

func validateAdminRunDateRanges(before, after string) error {
	if validString(&before) {
		_, err := time.Parse(time.RFC3339, before)
		if err != nil {
			return fmt.Errorf("invalid date format for CreatedBefore: '%s', must be in RFC3339 format", before)
		}
	}

	if validString(&after) {
		_, err := time.Parse(time.RFC3339, after)
		if err != nil {
			return fmt.Errorf("invalid date format for CreatedAfter: '%s', must be in RFC3339 format", after)
		}
	}

	return nil
}

func validateAdminRunFilterParams(runStatus string) error {
	// For the platform, an invalid filter value is a semantically understood query that returns an empty set, no error, no warning. But for go-tfe, an invalid value is good enough reason to error prior to a network call to the platform:
	if validString(&runStatus) {
		sanitizedRunstatus := strings.TrimSpace(runStatus)
		runStatuses := strings.Split(sanitizedRunstatus, ",")
		// iterate over our statuses, and ensure it is valid.
		for _, status := range runStatuses {
			switch status {
			case string(RunApplied),
				string(RunApplyQueued),
				string(RunApplying),
				string(RunCanceled),
				string(RunConfirmed),
				string(RunCostEstimate),
				string(RunCostEstimating),
				string(RunDiscarded),
				string(RunErrored),
				string(RunPending),
				string(RunPlanQueued),
				string(RunPlanned),
				string(RunPlannedAndFinished),
				string(RunPlanning),
				string(RunPolicyChecked),
				string(RunPolicyChecking),
				string(RunPolicyOverride),
				string(RunPolicySoftFailed),
				"":
				// do nothing
			default:
				return fmt.Errorf(`invalid value "%s" for run status`, status)
			}
		}
	}

	return nil
}
