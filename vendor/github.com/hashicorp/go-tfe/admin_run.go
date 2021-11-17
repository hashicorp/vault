package tfe

import (
	"context"
	"errors"
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
// TFE API docs: https://www.terraform.io/docs/cloud/api/admin/runs.html
type AdminRuns interface {
	// List all the runs of the given installation.
	List(ctx context.Context, options AdminRunsListOptions) (*AdminRunsList, error)

	// Force-cancel a run by its ID.
	ForceCancel(ctx context.Context, runID string, options AdminRunForceCancelOptions) error
}

// adminRuns implements the AdminRuns interface.
type adminRuns struct {
	client *Client
}

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

// AdminRunsListOptions represents the options for listing runs.
// https://www.terraform.io/docs/cloud/api/admin/runs.html#query-parameters
type AdminRunsListOptions struct {
	ListOptions

	RunStatus *string `url:"filter[status],omitempty"`
	Query     *string `url:"q,omitempty"`
	Include   *string `url:"include,omitempty"`
}

// List all the runs of the terraform enterprise installation.
// https://www.terraform.io/docs/cloud/api/admin/runs.html#list-all-runs
func (s *adminRuns) List(ctx context.Context, options AdminRunsListOptions) (*AdminRunsList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := "admin/runs"
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	rl := &AdminRunsList{}
	err = s.client.do(ctx, req, rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

// AdminRunForceCancelOptions represents the options for force-canceling a run.
type AdminRunForceCancelOptions struct {
	// An optional comment explaining the reason for the force-cancel.
	// https://www.terraform.io/docs/cloud/api/admin/runs.html#request-body
	Comment *string `json:"comment,omitempty"`
}

// ForceCancel is used to forcefully cancel a run by its ID.
// https://www.terraform.io/docs/cloud/api/admin/runs.html#force-a-run-into-the-quot-cancelled-quot-state
func (s *adminRuns) ForceCancel(ctx context.Context, runID string, options AdminRunForceCancelOptions) error {
	if !validStringID(&runID) {
		return errors.New("invalid value for run ID")
	}

	u := fmt.Sprintf("admin/runs/%s/actions/force-cancel", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

func (o AdminRunsListOptions) valid() error {
	if validString(o.RunStatus) {
		validRunStatus := map[string]int{
			string(RunApplied):            1,
			string(RunApplyQueued):        1,
			string(RunApplying):           1,
			string(RunCanceled):           1,
			string(RunConfirmed):          1,
			string(RunCostEstimated):      1,
			string(RunCostEstimating):     1,
			string(RunDiscarded):          1,
			string(RunErrored):            1,
			string(RunPending):            1,
			string(RunPlanQueued):         1,
			string(RunPlanned):            1,
			string(RunPlannedAndFinished): 1,
			string(RunPlanning):           1,
			string(RunPolicyChecked):      1,
			string(RunPolicyChecking):     1,
			string(RunPolicyOverride):     1,
			string(RunPolicySoftFailed):   1,
		}
		runStatus := strings.Split(*o.RunStatus, ",")

		// iterate over our statuses, and ensure it is valid.
		for _, status := range runStatus {
			if _, present := validRunStatus[status]; !present {
				return fmt.Errorf("invalid value %s for run status", status)
			}
		}
	}
	return nil
}
