// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Plans = (*plans)(nil)

// Plans describes all the plan related methods that the Terraform Enterprise
// API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/plans
type Plans interface {
	// Read a plan by its ID.
	Read(ctx context.Context, planID string) (*Plan, error)

	// Logs retrieves the logs of a plan.
	Logs(ctx context.Context, planID string) (io.Reader, error)

	// Retrieve the JSON execution plan
	ReadJSONOutput(ctx context.Context, planID string) ([]byte, error)
}

// plans implements Plans.
type plans struct {
	client *Client
}

// PlanStatus represents a plan state.
type PlanStatus string

// List all available plan statuses.
const (
	PlanCanceled    PlanStatus = "canceled"
	PlanCreated     PlanStatus = "created"
	PlanErrored     PlanStatus = "errored"
	PlanFinished    PlanStatus = "finished"
	PlanMFAWaiting  PlanStatus = "mfa_waiting"
	PlanPending     PlanStatus = "pending"
	PlanQueued      PlanStatus = "queued"
	PlanRunning     PlanStatus = "running"
	PlanUnreachable PlanStatus = "unreachable"
)

// Plan represents a Terraform Enterprise plan.
type Plan struct {
	ID                     string                `jsonapi:"primary,plans"`
	HasChanges             bool                  `jsonapi:"attr,has-changes"`
	GeneratedConfiguration bool                  `jsonapi:"attr,generated-configuration"`
	LogReadURL             string                `jsonapi:"attr,log-read-url"`
	ResourceAdditions      int                   `jsonapi:"attr,resource-additions"`
	ResourceChanges        int                   `jsonapi:"attr,resource-changes"`
	ResourceDestructions   int                   `jsonapi:"attr,resource-destructions"`
	ResourceImports        int                   `jsonapi:"attr,resource-imports"`
	Status                 PlanStatus            `jsonapi:"attr,status"`
	StatusTimestamps       *PlanStatusTimestamps `jsonapi:"attr,status-timestamps"`

	// Relations
	Exports []*PlanExport `jsonapi:"relation,exports"`
}

// PlanStatusTimestamps holds the timestamps for individual plan statuses.
type PlanStatusTimestamps struct {
	CanceledAt      time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	ErroredAt       time.Time `jsonapi:"attr,errored-at,rfc3339"`
	FinishedAt      time.Time `jsonapi:"attr,finished-at,rfc3339"`
	ForceCanceledAt time.Time `jsonapi:"attr,force-canceled-at,rfc3339"`
	QueuedAt        time.Time `jsonapi:"attr,queued-at,rfc3339"`
	StartedAt       time.Time `jsonapi:"attr,started-at,rfc3339"`
}

// Read a plan by its ID.
func (s *plans) Read(ctx context.Context, planID string) (*Plan, error) {
	if !validStringID(&planID) {
		return nil, ErrInvalidPlanID
	}

	u := fmt.Sprintf("plans/%s", url.PathEscape(planID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	p := &Plan{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Logs retrieves the logs of a plan.
func (s *plans) Logs(ctx context.Context, planID string) (io.Reader, error) {
	if !validStringID(&planID) {
		return nil, ErrInvalidPlanID
	}

	// Get the plan to make sure it exists.
	p, err := s.Read(ctx, planID)
	if err != nil {
		return nil, err
	}

	// Return an error if the log URL is empty.
	if p.LogReadURL == "" {
		return nil, fmt.Errorf("plan %s does not have a log URL", planID)
	}

	u, err := url.Parse(p.LogReadURL)
	if err != nil {
		return nil, fmt.Errorf("invalid log URL: %w", err)
	}

	done := func() (bool, error) {
		p, err := s.Read(ctx, p.ID)
		if err != nil {
			return false, err
		}

		switch p.Status {
		case PlanCanceled, PlanErrored, PlanFinished, PlanUnreachable:
			return true, nil
		default:
			return false, nil
		}
	}

	return &LogReader{
		client: s.client,
		ctx:    ctx,
		done:   done,
		logURL: u,
	}, nil
}

// Retrieve the JSON execution plan
func (s *plans) ReadJSONOutput(ctx context.Context, planID string) ([]byte, error) {
	if !validStringID(&planID) {
		return nil, ErrInvalidPlanID
	}

	u := fmt.Sprintf("plans/%s/json-output", url.PathEscape(planID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = req.Do(ctx, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
