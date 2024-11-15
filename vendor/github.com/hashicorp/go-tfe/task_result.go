// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"time"
)

// Compile-time proof of interface implementation
var _ TaskResults = (*taskResults)(nil)

// TaskResults describes all the task result related methods that the HCP Terraform or Terraform Enterprise API supports.
type TaskResults interface {
	// Read a task result by ID
	Read(ctx context.Context, taskResultID string) (*TaskResult, error)
}

// taskResults implements TaskResults
type taskResults struct {
	client *Client
}

// TaskResultStatus is an enum that represents all possible statuses for a task result
type TaskResultStatus string

const (
	TaskPassed      TaskResultStatus = "passed"
	TaskFailed      TaskResultStatus = "failed"
	TaskPending     TaskResultStatus = "pending"
	TaskRunning     TaskResultStatus = "running"
	TaskUnreachable TaskResultStatus = "unreachable"
	TaskErrored     TaskResultStatus = "errored"
)

// TaskEnforcementLevel is an enum that describes the enforcement levels for a run task
type TaskEnforcementLevel string

const (
	Advisory  TaskEnforcementLevel = "advisory"
	Mandatory TaskEnforcementLevel = "mandatory"
)

// TaskResultStatusTimestamps represents the set of timestamps recorded for a task result
type TaskResultStatusTimestamps struct {
	ErroredAt  time.Time `jsonapi:"attr,errored-at,rfc3339"`
	RunningAt  time.Time `jsonapi:"attr,running-at,rfc3339"`
	CanceledAt time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	FailedAt   time.Time `jsonapi:"attr,failed-at,rfc3339"`
	PassedAt   time.Time `jsonapi:"attr,passed-at,rfc3339"`
}

// TaskResult represents the result of a HCP Terraform or Terraform Enterprise run task
type TaskResult struct {
	ID                            string                     `jsonapi:"primary,task-results"`
	Status                        TaskResultStatus           `jsonapi:"attr,status"`
	Message                       string                     `jsonapi:"attr,message"`
	StatusTimestamps              TaskResultStatusTimestamps `jsonapi:"attr,status-timestamps"`
	URL                           string                     `jsonapi:"attr,url"`
	CreatedAt                     time.Time                  `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt                     time.Time                  `jsonapi:"attr,updated-at,iso8601"`
	TaskID                        string                     `jsonapi:"attr,task-id"`
	TaskName                      string                     `jsonapi:"attr,task-name"`
	TaskURL                       string                     `jsonapi:"attr,task-url"`
	WorkspaceTaskID               string                     `jsonapi:"attr,workspace-task-id"`
	WorkspaceTaskEnforcementLevel TaskEnforcementLevel       `jsonapi:"attr,workspace-task-enforcement-level"`

	// The task stage this result belongs to
	TaskStage *TaskStage `jsonapi:"relation,task_stage"`
}

// Read a task result by ID
func (t *taskResults) Read(ctx context.Context, taskResultID string) (*TaskResult, error) {
	if !validStringID(&taskResultID) {
		return nil, ErrInvalidTaskResultID
	}

	u := fmt.Sprintf("task-results/%s", taskResultID)
	req, err := t.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	r := &TaskResult{}
	err = req.Do(ctx, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
