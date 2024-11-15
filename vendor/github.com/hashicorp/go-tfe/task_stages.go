// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"time"
)

// Compile-time proof of interface  implementation
var _ TaskStages = (*taskStages)(nil)

// TaskStages describes all the task stage related methods that the HCP Terraform and Terraform Enterprise API
// supports.
type TaskStages interface {
	// Read a task stage by ID
	Read(ctx context.Context, taskStageID string, options *TaskStageReadOptions) (*TaskStage, error)

	// List all task stages for a given run
	List(ctx context.Context, runID string, options *TaskStageListOptions) (*TaskStageList, error)

	// **Note: This function is still in BETA and subject to change.**
	// Override a task stage for a given run
	Override(ctx context.Context, taskStageID string, options TaskStageOverrideOptions) (*TaskStage, error)
}

// taskStages implements TaskStages
type taskStages struct {
	client *Client
}

// Stage is an enum that represents the possible run stages for run tasks
type Stage string

const (
	PrePlan   Stage = "pre_plan"
	PostPlan  Stage = "post_plan"
	PreApply  Stage = "pre_apply"
	PostApply Stage = "post_apply"
)

// TaskStageStatus is an enum that represents all possible statuses for a task stage
type TaskStageStatus string

const (
	TaskStagePending          TaskStageStatus = "pending"
	TaskStageRunning          TaskStageStatus = "running"
	TaskStagePassed           TaskStageStatus = "passed"
	TaskStageFailed           TaskStageStatus = "failed"
	TaskStageAwaitingOverride TaskStageStatus = "awaiting_override"
	TaskStageCanceled         TaskStageStatus = "canceled"
	TaskStageErrored          TaskStageStatus = "errored"
	TaskStageUnreachable      TaskStageStatus = "unreachable"
)

// Permissions represents the permission types for overridding a task stage
type Permissions struct {
	CanOverridePolicy *bool `jsonapi:"attr,can-override-policy"`
	CanOverrideTasks  *bool `jsonapi:"attr,can-override-tasks"`
	CanOverride       *bool `jsonapi:"attr,can-override"`
}

// Actions represents a task stage actions
type Actions struct {
	IsOverridable *bool `jsonapi:"attr,is-overridable"`
}

// TaskStage represents a HCP Terraform or Terraform Enterprise run's stage where run tasks can occur
type TaskStage struct {
	ID               string                    `jsonapi:"primary,task-stages"`
	Stage            Stage                     `jsonapi:"attr,stage"`
	Status           TaskStageStatus           `jsonapi:"attr,status"`
	StatusTimestamps TaskStageStatusTimestamps `jsonapi:"attr,status-timestamps"`
	CreatedAt        time.Time                 `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt        time.Time                 `jsonapi:"attr,updated-at,iso8601"`
	Permissions      *Permissions              `jsonapi:"attr,permissions"`
	Actions          *Actions                  `jsonapi:"attr,actions"`

	Run               *Run                `jsonapi:"relation,run"`
	TaskResults       []*TaskResult       `jsonapi:"relation,task-results"`
	PolicyEvaluations []*PolicyEvaluation `jsonapi:"relation,policy-evaluations"`
}

// TaskStageOverrideOptions represents the options for overriding a TaskStage.
type TaskStageOverrideOptions struct {
	// An optional explanation for why the stage was overridden
	Comment *string `json:"comment,omitempty"`
}

// TaskStageList represents a list of task stages
type TaskStageList struct {
	*Pagination
	Items []*TaskStage
}

// TaskStageStatusTimestamps represents the set of timestamps recorded for a task stage
type TaskStageStatusTimestamps struct {
	ErroredAt  time.Time `jsonapi:"attr,errored-at,rfc3339"`
	RunningAt  time.Time `jsonapi:"attr,running-at,rfc3339"`
	CanceledAt time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	FailedAt   time.Time `jsonapi:"attr,failed-at,rfc3339"`
	PassedAt   time.Time `jsonapi:"attr,passed-at,rfc3339"`
}

// TaskStageIncludeOpt represents the available options for include query params.
type TaskStageIncludeOpt string

const TaskStageTaskResults TaskStageIncludeOpt = "task_results"

// **Note: This field is still in BETA and subject to change.**
const PolicyEvaluationsTaskResults TaskStageIncludeOpt = "policy_evaluations"

// TaskStageReadOptions represents the set of options when reading a task stage
type TaskStageReadOptions struct {
	// Optional: A list of relations to include.
	Include []TaskStageIncludeOpt `url:"include,omitempty"`
}

// TaskStageListOptions represents the options for listing task stages for a run
type TaskStageListOptions struct {
	ListOptions
}

// Read a task stage by ID
func (s *taskStages) Read(ctx context.Context, taskStageID string, options *TaskStageReadOptions) (*TaskStage, error) {
	if !validStringID(&taskStageID) {
		return nil, ErrInvalidTaskStageID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("task-stages/%s", taskStageID)
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	t := &TaskStage{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// List task stages for a run
func (s *taskStages) List(ctx context.Context, runID string, options *TaskStageListOptions) (*TaskStageList, error) {
	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/task-stages", runID)
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	tlist := &TaskStageList{}

	err = req.Do(ctx, tlist)
	if err != nil {
		return nil, err
	}

	return tlist, nil
}

// **Note: This function is still in BETA and subject to change.**
// Override a task stages for a run
func (s *taskStages) Override(ctx context.Context, taskStageID string, options TaskStageOverrideOptions) (*TaskStage, error) {
	if !validStringID(&taskStageID) {
		return nil, ErrInvalidTaskStageID
	}

	u := fmt.Sprintf("task-stages/%s/actions/override", taskStageID)
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	t := &TaskStage{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (o *TaskStageReadOptions) valid() error {
	return nil
}
