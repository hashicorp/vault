package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Runs = (*runs)(nil)

// Runs describes all the run related methods that the Terraform Enterprise
// API supports.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/run.html
type Runs interface {
	// List all the runs of the given workspace.
	List(ctx context.Context, workspaceID string, options RunListOptions) (*RunList, error)

	// Create a new run with the given options.
	Create(ctx context.Context, options RunCreateOptions) (*Run, error)

	// Read a run by its ID.
	Read(ctx context.Context, runID string) (*Run, error)

	// ReadWithOptions reads a run by its ID using the options supplied
	ReadWithOptions(ctx context.Context, runID string, options *RunReadOptions) (*Run, error)

	// Apply a run by its ID.
	Apply(ctx context.Context, runID string, options RunApplyOptions) error

	// Cancel a run by its ID.
	Cancel(ctx context.Context, runID string, options RunCancelOptions) error

	// Force-cancel a run by its ID.
	ForceCancel(ctx context.Context, runID string, options RunForceCancelOptions) error

	// Discard a run by its ID.
	Discard(ctx context.Context, runID string, options RunDiscardOptions) error
}

// runs implements Runs.
type runs struct {
	client *Client
}

// RunStatus represents a run state.
type RunStatus string

//List all available run statuses.
const (
	RunApplied            RunStatus = "applied"
	RunApplyQueued        RunStatus = "apply_queued"
	RunApplying           RunStatus = "applying"
	RunCanceled           RunStatus = "canceled"
	RunConfirmed          RunStatus = "confirmed"
	RunCostEstimated      RunStatus = "cost_estimated"
	RunCostEstimating     RunStatus = "cost_estimating"
	RunDiscarded          RunStatus = "discarded"
	RunErrored            RunStatus = "errored"
	RunPending            RunStatus = "pending"
	RunPlanQueued         RunStatus = "plan_queued"
	RunPlanned            RunStatus = "planned"
	RunPlannedAndFinished RunStatus = "planned_and_finished"
	RunPlanning           RunStatus = "planning"
	RunPolicyChecked      RunStatus = "policy_checked"
	RunPolicyChecking     RunStatus = "policy_checking"
	RunPolicyOverride     RunStatus = "policy_override"
	RunPolicySoftFailed   RunStatus = "policy_soft_failed"
)

// RunSource represents a source type of a run.
type RunSource string

// List all available run sources.
const (
	RunSourceAPI                  RunSource = "tfe-api"
	RunSourceConfigurationVersion RunSource = "tfe-configuration-version"
	RunSourceUI                   RunSource = "tfe-ui"
)

// RunList represents a list of runs.
type RunList struct {
	*Pagination
	Items []*Run
}

// Run represents a Terraform Enterprise run.
type Run struct {
	ID                     string               `jsonapi:"primary,runs"`
	Actions                *RunActions          `jsonapi:"attr,actions"`
	AutoApply              bool                 `jsonapi:"attr,auto-apply,omitempty"`
	CreatedAt              time.Time            `jsonapi:"attr,created-at,iso8601"`
	ForceCancelAvailableAt time.Time            `jsonapi:"attr,force-cancel-available-at,iso8601"`
	HasChanges             bool                 `jsonapi:"attr,has-changes"`
	IsDestroy              bool                 `jsonapi:"attr,is-destroy"`
	Message                string               `jsonapi:"attr,message"`
	Permissions            *RunPermissions      `jsonapi:"attr,permissions"`
	PositionInQueue        int                  `jsonapi:"attr,position-in-queue"`
	Refresh                bool                 `jsonapi:"attr,refresh"`
	RefreshOnly            bool                 `jsonapi:"attr,refresh-only"`
	ReplaceAddrs           []string             `jsonapi:"attr,replace-addrs,omitempty"`
	Source                 RunSource            `jsonapi:"attr,source"`
	Status                 RunStatus            `jsonapi:"attr,status"`
	StatusTimestamps       *RunStatusTimestamps `jsonapi:"attr,status-timestamps"`
	TargetAddrs            []string             `jsonapi:"attr,target-addrs,omitempty"`

	// Relations
	Apply                *Apply                `jsonapi:"relation,apply"`
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`
	CostEstimate         *CostEstimate         `jsonapi:"relation,cost-estimate"`
	CreatedBy            *User                 `jsonapi:"relation,created-by"`
	Plan                 *Plan                 `jsonapi:"relation,plan"`
	PolicyChecks         []*PolicyCheck        `jsonapi:"relation,policy-checks"`
	Workspace            *Workspace            `jsonapi:"relation,workspace"`
}

// RunActions represents the run actions.
type RunActions struct {
	IsCancelable      bool `jsonapi:"attr,is-cancelable"`
	IsConfirmable     bool `jsonapi:"attr,is-confirmable"`
	IsDiscardable     bool `jsonapi:"attr,is-discardable"`
	IsForceCancelable bool `jsonapi:"attr,is-force-cancelable"`
}

// RunPermissions represents the run permissions.
type RunPermissions struct {
	CanApply        bool `jsonapi:"attr,can-apply"`
	CanCancel       bool `jsonapi:"attr,can-cancel"`
	CanDiscard      bool `jsonapi:"attr,can-discard"`
	CanForceCancel  bool `jsonapi:"attr,can-force-cancel"`
	CanForceExecute bool `jsonapi:"attr,can-force-execute"`
}

// RunStatusTimestamps holds the timestamps for individual run statuses.
type RunStatusTimestamps struct {
	AppliedAt            time.Time `jsonapi:"attr,applied-at,rfc3339"`
	ApplyQueuedAt        time.Time `jsonapi:"attr,apply-queued-at,rfc3339"`
	ApplyingAt           time.Time `jsonapi:"attr,applying-at,rfc3339"`
	CanceledAt           time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	ConfirmedAt          time.Time `jsonapi:"attr,confirmed-at,rfc3339"`
	CostEstimatedAt      time.Time `jsonapi:"attr,cost-estimated-at,rfc3339"`
	CostEstimatingAt     time.Time `jsonapi:"attr,cost-estimating-at,rfc3339"`
	DiscardedAt          time.Time `jsonapi:"attr,discarded-at,rfc3339"`
	ErroredAt            time.Time `jsonapi:"attr,errored-at,rfc3339"`
	ForceCanceledAt      time.Time `jsonapi:"attr,force-canceled-at,rfc3339"`
	PlanQueueableAt      time.Time `jsonapi:"attr,plan-queueable-at,rfc3339"`
	PlanQueuedAt         time.Time `jsonapi:"attr,plan-queued-at,rfc3339"`
	PlannedAndFinishedAt time.Time `jsonapi:"attr,planned-and-finished-at,rfc3339"`
	PlannedAt            time.Time `jsonapi:"attr,planned-at,rfc3339"`
	PlanningAt           time.Time `jsonapi:"attr,planning-at,rfc3339"`
	PolicyCheckedAt      time.Time `jsonapi:"attr,policy-checked-at,rfc3339"`
	PolicySoftFailedAt   time.Time `jsonapi:"attr,policy-soft-failed-at,rfc3339"`
}

// RunListOptions represents the options for listing runs.
type RunListOptions struct {
	ListOptions

	// A list of relations to include. See available resources:
	// https://www.terraform.io/docs/cloud/api/run.html#available-related-resources
	Include *string `url:"include"`
}

// List all the runs of the given workspace.
func (s *runs) List(ctx context.Context, workspaceID string, options RunListOptions) (*RunList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/runs", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	rl := &RunList{}
	err = s.client.do(ctx, req, rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

// RunCreateOptions represents the options for creating a new run.
type RunCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,runs"`

	// Specifies if this plan is a destroy plan, which will destroy all
	// provisioned resources.
	IsDestroy *bool `jsonapi:"attr,is-destroy,omitempty"`

	// Refresh determines if the run should
	// update the state prior to checking for differences
	Refresh *bool `jsonapi:"attr,refresh,omitempty"`

	// RefreshOnly determines whether the run should ignore config changes
	// and refresh the state only
	RefreshOnly *bool `jsonapi:"attr,refresh-only,omitempty"`

	// Specifies the message to be associated with this run.
	Message *string `jsonapi:"attr,message,omitempty"`

	// Specifies the configuration version to use for this run. If the
	// configuration version object is omitted, the run will be created using the
	// workspace's latest configuration version.
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`

	// Specifies the workspace where the run will be executed.
	Workspace *Workspace `jsonapi:"relation,workspace"`

	// If non-empty, requests that Terraform should create a plan including
	// actions only for the given objects (specified using resource address
	// syntax) and the objects they depend on.
	//
	// This capability is provided for exceptional circumstances only, such as
	// recovering from mistakes or working around existing Terraform
	// limitations. Terraform will generally mention the -target command line
	// option in its error messages describing situations where setting this
	// argument may be appropriate. This argument should not be used as part
	// of routine workflow and Terraform will emit warnings reminding about
	// this whenever this property is set.
	TargetAddrs []string `jsonapi:"attr,target-addrs,omitempty"`

	// If non-empty, requests that Terraform create a plan that replaces
	// (destroys and then re-creates) the objects specified by the given
	// resource addresses.
	ReplaceAddrs []string `jsonapi:"attr,replace-addrs,omitempty"`

	// AutoApply determines if the run should be applied automatically without
	// user confirmation. It defaults to the Workspace.AutoApply setting.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`
}

func (o RunCreateOptions) valid() error {
	if o.Workspace == nil {
		return errors.New("workspace is required")
	}
	return nil
}

// Create a new run with the given options.
func (s *runs) Create(ctx context.Context, options RunCreateOptions) (*Run, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("POST", "runs", &options)
	if err != nil {
		return nil, err
	}

	r := &Run{}
	err = s.client.do(ctx, req, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Read a run by its ID.
func (s *runs) Read(ctx context.Context, runID string) (*Run, error) {
	return s.ReadWithOptions(ctx, runID, nil)
}

// RunReadOptions represents the options for reading a run.
type RunReadOptions struct {
	Include string `url:"include"`
}

// Read a run by its ID with the given options.
func (s *runs) ReadWithOptions(ctx context.Context, runID string, options *RunReadOptions) (*Run, error) {
	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s", url.QueryEscape(runID))
	req, err := s.client.newRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	r := &Run{}
	err = s.client.do(ctx, req, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// RunApplyOptions represents the options for applying a run.
type RunApplyOptions struct {
	// An optional comment about the run.
	Comment *string `jsonapi:"attr,comment,omitempty"`
}

// Apply a run by its ID.
func (s *runs) Apply(ctx context.Context, runID string, options RunApplyOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/apply", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// RunCancelOptions represents the options for canceling a run.
type RunCancelOptions struct {
	// An optional explanation for why the run was canceled.
	Comment *string `jsonapi:"attr,comment,omitempty"`
}

// Cancel a run by its ID.
func (s *runs) Cancel(ctx context.Context, runID string, options RunCancelOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/cancel", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// RunForceCancelOptions represents the options for force-canceling a run.
type RunForceCancelOptions struct {
	// An optional comment explaining the reason for the force-cancel.
	Comment *string `jsonapi:"attr,comment,omitempty"`
}

// ForceCancel is used to forcefully cancel a run by its ID.
func (s *runs) ForceCancel(ctx context.Context, runID string, options RunForceCancelOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/force-cancel", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// RunDiscardOptions represents the options for discarding a run.
type RunDiscardOptions struct {
	// An optional explanation for why the run was discarded.
	Comment *string `jsonapi:"attr,comment,omitempty"`
}

// Discard a run by its ID.
func (s *runs) Discard(ctx context.Context, runID string, options RunDiscardOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/discard", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
