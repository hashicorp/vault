// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// StackPlans describes all the stacks plans-related methods that the HCP Terraform API supports.
// NOTE WELL: This is a beta feature and is subject to change until noted otherwise in the
// release notes.
type StackPlans interface {
	// Read returns a stack plan by its ID.
	Read(ctx context.Context, stackPlanID string) (*StackPlan, error)

	// ListByConfiguration returns a list of stack plans for a given stack configuration.
	ListByConfiguration(ctx context.Context, stackConfigurationID string, options *StackPlansListOptions) (*StackPlanList, error)

	// Approve approves a stack plan.
	Approve(ctx context.Context, stackPlanID string) error

	// Cancel cancels a stack plan.
	Cancel(ctx context.Context, stackPlanID string) error

	// Discard discards a stack plan.
	Discard(ctx context.Context, stackPlanID string) error

	// PlanDescription returns the plan description for a stack plan.
	PlanDescription(ctx context.Context, stackPlanID string) (*JSONChangeDesc, error)

	// AwaitTerminal generates a channel that will receive the status of the stack plan as it progresses.
	// See WaitForStatusResult for more information.
	AwaitTerminal(ctx context.Context, stackPlanID string) <-chan WaitForStatusResult

	// AwaitRunning generates a channel that will receive the status of the stack plan as it progresses.
	// See WaitForStatusResult for more information.
	AwaitRunning(ctx context.Context, stackPlanID string) <-chan WaitForStatusResult
}

type StackPlansStatusFilter string

const (
	StackPlansStatusFilterCreated   StackPlansStatusFilter = "created"
	StackPlansStatusFilterRunning   StackPlansStatusFilter = "running"
	StackPlansStatusFilterPaused    StackPlansStatusFilter = "paused"
	StackPlansStatusFilterFinished  StackPlansStatusFilter = "finished"
	StackPlansStatusFilterDiscarded StackPlansStatusFilter = "discarded"
	StackPlansStatusFilterErrored   StackPlansStatusFilter = "errored"
	StackPlansStatusFilterCanceled  StackPlansStatusFilter = "canceled"
)

type StackPlanStatus string

const (
	StackPlanStatusCreated           StackPlanStatus = "created"
	StackPlanStatusRunning           StackPlanStatus = "running"
	StackPlanStatusRunningQueued     StackPlanStatus = "running_queued"
	StackPlanStatusRunningPlanning   StackPlanStatus = "running_planning"
	StackPlanStatusRunningApplying   StackPlanStatus = "running_applying"
	StackPlanStatusFinished          StackPlanStatus = "finished"
	StackPlanStatusFinishedNoChanges StackPlanStatus = "finished_no_changes"
	StackPlanStatusFinishedPlanned   StackPlanStatus = "finished_planned"
	StackPlanStatusFinishedApplied   StackPlanStatus = "finished_applied"
	StackPlanStatusDiscarded         StackPlanStatus = "discarded"
	StackPlanStatusErrored           StackPlanStatus = "errored"
	StackPlanStatusCanceled          StackPlanStatus = "canceled"
)

type StackPlansIncludeOpt string

func (s StackPlanStatus) String() string {
	return string(s)
}

const (
	StackPlansIncludeOperations StackPlansIncludeOpt = "stack_plan_operations"
)

type StackPlansListOptions struct {
	ListOptions

	// Optional: A query string to filter plans by status.
	Status StackPlansStatusFilter `url:"filter[status],omitempty"`

	// Optional: A query string to filter plans by deployment.
	Deployment string `url:"filter[deployment],omitempty"`

	Include []StackPlansIncludeOpt `url:"include,omitempty"`
}

type StackPlanList struct {
	*Pagination
	Items []*StackPlan
}

// stackPlans implements StackPlans.
type stackPlans struct {
	client *Client
}

var _ StackPlans = &stackPlans{}

// StackPlanStatusTimestamps are the timestamps of the status changes for a stack
type StackPlanStatusTimestamps struct {
	CreatedAt  time.Time `jsonapi:"attr,created-at,rfc3339"`
	RunningAt  time.Time `jsonapi:"attr,running-at,rfc3339"`
	PausedAt   time.Time `jsonapi:"attr,paused-at,rfc3339"`
	FinishedAt time.Time `jsonapi:"attr,finished-at,rfc3339"`
}

// PlanChanges is the summary of the planned changes
type PlanChanges struct {
	Add    int `jsonapi:"attr,add"`
	Total  int `jsonapi:"attr,total"`
	Change int `jsonapi:"attr,change"`
	Import int `jsonapi:"attr,import"`
	Remove int `jsonapi:"attr,remove"`
}

// StackPlan represents a plan for a stack.
type StackPlan struct {
	ID               string                     `jsonapi:"primary,stack-plans"`
	PlanMode         string                     `jsonapi:"attr,plan-mode"`
	PlanNumber       string                     `jsonapi:"attr,plan-number"`
	Status           StackPlanStatus            `jsonapi:"attr,status"`
	StatusTimestamps *StackPlanStatusTimestamps `jsonapi:"attr,status-timestamps"`
	IsPlanned        bool                       `jsonapi:"attr,is-planned"`
	Changes          *PlanChanges               `jsonapi:"attr,changes"`
	Deployment       string                     `jsonapi:"attr,deployment"`

	// Relationships
	StackConfiguration  *StackConfiguration   `jsonapi:"relation,stack-configuration"`
	Stack               *Stack                `jsonapi:"relation,stack"`
	StackPlanOperations []*StackPlanOperation `jsonapi:"relation,stack-plan-operations"`
}

// JSONChangeDesc represents a change description of a stack plan / apply operation.
type JSONChangeDesc struct {
	FormatVersion             uint64                         `json:"terraform_stack_change_description"`
	Interim                   bool                           `json:"interim,omitempty"`
	Applyable                 bool                           `json:"applyable"`
	PlanMode                  string                         `json:"plan_mode"`
	Components                []JSONComponent                `json:"components"`
	ResourceInstances         []JSONResourceInstance         `json:"resource_instances"`
	DeferredResourceInstances []JSONResourceInstanceDeferral `json:"deferred_resource_instances"`
	Outputs                   map[string]JSONOutput          `json:"outputs"`
}

// JSONComponent represents a change description of a single component in a plan.
type JSONComponent struct {
	Address             string         `json:"address"`
	ComponentAddress    string         `json:"component_address"`
	InstanceCorrelator  string         `json:"instance_correlator"`
	ComponentCorrelator string         `json:"component_correlator"`
	Actions             []ChangeAction `json:"actions"`
	Complete            bool           `json:"complete"`
}

// ChangeAction are the actions a change can have: no-op, create, read, update, delte, forget.
type ChangeAction string

// JSONResourceInstance is the change description of a single resource instance in a plan.
type JSONResourceInstance struct {
	ComponentInstanceCorrelator      string          `json:"component_instance_correlator"`
	ComponentInstanceAddress         string          `json:"component_instance_address"`
	Address                          string          `json:"address"`
	PreviousComponentInstanceAddress string          `json:"previous_component_instance_address,omitempty"`
	PreviousAddress                  string          `json:"previous_address,omitempty"`
	DeposedKey                       string          `json:"deposed,omitempty"`
	ResourceMode                     string          `json:"mode,omitempty"`
	ResourceType                     string          `json:"type"`
	ProviderAddr                     string          `json:"provider_name"`
	Change                           Change          `json:"change"`
	ResourceName                     string          `json:"resource_name"`
	Index                            json.RawMessage `json:"index"`
	IndexUnknown                     bool            `json:"index_unknown"`
	ModuleAddr                       string          `json:"module_address"`
	ActionReason                     string          `json:"action_reason,omitempty"`
}

// JSONResourceInstanceDeferral is the change description of a single resource instance that is deferred.
type JSONResourceInstanceDeferral struct {
	ResourceInstance JSONResourceInstance `json:"resource_instance"`
	Deferred         JSONDeferred         `json:"deferred"`
}

// JSONDeferred contains the reason why a resource instance is deferred: instance_count_unknown, resource_config_unknown, provider_config_unknown, provider_config_unknown, or deferred_prereq.
type JSONDeferred struct {
	Reason string `json:"reason"`
}

// JSONOutput is the value of a single output in a plan.
type JSONOutput struct {
	Change json.RawMessage `json:"change"`
}

// Change represents the change of a resource instance in a plan.
type Change struct {
	Actions         []ChangeAction  `json:"actions"`
	After           json.RawMessage `json:"after"`
	Before          json.RawMessage `json:"before"`
	AfterUnknown    json.RawMessage `json:"after_unknown"`
	BeforeSensitive json.RawMessage `json:"before_sensitive"`
	AfterSensitive  json.RawMessage `json:"after_sensitive"`
	Importing       *JSONImporting  `json:"importing,omitempty"`
	ReplacePaths    json.RawMessage `json:"replace_paths,omitempty"`
}

// JSONImporting represents the import status of a resource instance in a plan.
type JSONImporting struct {
	// True within a deferred instance
	Unknown         bool   `json:"unknown"`
	ID              string `json:"id"`
	GeneratedConfig string `json:"generated_config"`
}

func (s stackPlans) Read(ctx context.Context, stackPlanID string) (*StackPlan, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("stack-plans/%s", url.PathEscape(stackPlanID)), nil)
	if err != nil {
		return nil, err
	}

	sp := &StackPlan{}
	err = req.Do(ctx, sp)
	if err != nil {
		return nil, err
	}

	return sp, nil
}

func (s stackPlans) ListByConfiguration(ctx context.Context, stackConfigurationID string, options *StackPlansListOptions) (*StackPlanList, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("stack-configurations/%s/stack-plans", url.PathEscape(stackConfigurationID)), options)
	if err != nil {
		return nil, err
	}

	sl := &StackPlanList{}
	err = req.Do(ctx, sl)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (s stackPlans) Approve(ctx context.Context, stackPlanID string) error {
	req, err := s.client.NewRequest("POST", fmt.Sprintf("stack-plans/%s/approve", url.PathEscape(stackPlanID)), nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (s stackPlans) Discard(ctx context.Context, stackPlanID string) error {
	req, err := s.client.NewRequest("POST", fmt.Sprintf("stack-plans/%s/discard", url.PathEscape(stackPlanID)), nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (s stackPlans) Cancel(ctx context.Context, stackPlanID string) error {
	req, err := s.client.NewRequest("POST", fmt.Sprintf("stack-plans/%s/cancel", url.PathEscape(stackPlanID)), nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (s stackPlans) PlanDescription(ctx context.Context, stackPlanID string) (*JSONChangeDesc, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("stack-plans/%s/plan-description", url.PathEscape(stackPlanID)), nil)
	if err != nil {
		return nil, err
	}

	jd := &JSONChangeDesc{}
	err = req.DoJSON(ctx, jd)
	if err != nil {
		return nil, err
	}

	return jd, nil
}

// AwaitTerminal generates a channel that will receive the status of the stack plan as it progresses.
// The channel will be closed when the stack plan reaches a final status or an error occurs. The
// read will be retried dependending on the configuration of the client. When the channel is closed,
// the last value will either be a terminal status (finished, finished_no_changes, finished_applied,
// finished_planned, discarded, canceled, errorer), or an error. The status check will continue even
// if the stack plan is waiting for approval. Check the status within the the channel to determine
// if the stack plan needs approval.
func (s stackPlans) AwaitTerminal(ctx context.Context, stackPlanID string) <-chan WaitForStatusResult {
	return awaitPoll(ctx, stackPlanID, func(ctx context.Context) (string, error) {
		stackPlan, err := s.Read(ctx, stackPlanID)
		if err != nil {
			return "", err
		}

		return stackPlan.Status.String(), nil
	}, []string{
		StackPlanStatusFinished.String(),
		StackPlanStatusFinishedNoChanges.String(),
		StackPlanStatusFinishedApplied.String(),
		StackPlanStatusFinishedPlanned.String(),
		StackPlanStatusDiscarded.String(),
		StackPlanStatusErrored.String(),
		StackPlanStatusCanceled.String(),
	})
}

// AwaitRunning generates a channel that will receive the status of the stack plan as it progresses.
// The channel will be closed when the stack plan reaches a running status (running, running_queued,
// running_planning, running_applying), a terminal status (finished, finished_no_changes, finished_applied,
// finished_planned, discarded, canceled, errorer), or an error occurs. The read will be retried
// dependending on the configuration of the client.
func (s stackPlans) AwaitRunning(ctx context.Context, stackPlanID string) <-chan WaitForStatusResult {
	return awaitPoll(ctx, stackPlanID, func(ctx context.Context) (string, error) {
		stackPlan, err := s.Read(ctx, stackPlanID)
		if err != nil {
			return "", err
		}

		return stackPlan.Status.String(), nil
	}, []string{
		StackPlanStatusRunning.String(),
		StackPlanStatusRunningPlanning.String(),
		StackPlanStatusRunningApplying.String(),
		StackPlanStatusFinished.String(),
		StackPlanStatusFinishedNoChanges.String(),
		StackPlanStatusFinishedApplied.String(),
		StackPlanStatusFinishedPlanned.String(),
		StackPlanStatusDiscarded.String(),
		StackPlanStatusErrored.String(),
		StackPlanStatusCanceled.String(),
	})
}
