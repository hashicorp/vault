// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ RunEvents = (*runEvents)(nil)

// RunEvents describes all the run events that the Terraform Enterprise
// API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run
type RunEvents interface {
	// List all the runs events of the given run.
	List(ctx context.Context, runID string, options *RunEventListOptions) (*RunEventList, error)

	// Read a run event by its ID.
	Read(ctx context.Context, runEventID string) (*RunEvent, error)

	// ReadWithOptions reads a run event by its ID using the options supplied
	ReadWithOptions(ctx context.Context, runEventID string, options *RunEventReadOptions) (*RunEvent, error)
}

// runEvents implements RunEvents.
type runEvents struct {
	client *Client
}

// RunEventList represents a list of run events.
type RunEventList struct {
	// Pagination is not supported by the API
	*Pagination
	Items []*RunEvent
}

// RunEvent represents a Terraform Enterprise run event.
type RunEvent struct {
	ID          string    `jsonapi:"primary,run-events"`
	Action      string    `jsonapi:"attr,action"`
	CreatedAt   time.Time `jsonapi:"attr,created-at,iso8601"`
	Description string    `jsonapi:"attr,description"`

	// Relations - Note that `target` is not supported yet
	Actor   *User    `jsonapi:"relation,actor"`
	Comment *Comment `jsonapi:"relation,comment"`
}

// RunEventIncludeOpt represents the available options for include query params.
type RunEventIncludeOpt string

const (
	RunEventComment RunEventIncludeOpt = "comment"
	RunEventActor   RunEventIncludeOpt = "actor"
)

// RunEventListOptions represents the options for listing run events.
type RunEventListOptions struct {
	// Optional: A list of relations to include. See available resources:
	Include []RunEventIncludeOpt `url:"include,omitempty"`
}

// RunEventReadOptions represents the options for reading a run event.
type RunEventReadOptions struct {
	// Optional: A list of relations to include. See available resources:
	Include []RunEventIncludeOpt `url:"include,omitempty"`
}

// List all the run events of the given run.
func (s *runEvents) List(ctx context.Context, runID string, options *RunEventListOptions) (*RunEventList, error) {
	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("runs/%s/run-events", url.PathEscape(runID))

	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	rl := &RunEventList{}
	err = req.Do(ctx, rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

// Read a run by its ID.
func (s *runEvents) Read(ctx context.Context, runEventID string) (*RunEvent, error) {
	return s.ReadWithOptions(ctx, runEventID, nil)
}

// ReadWithOptions reads a run by its ID with the given options.
func (s *runEvents) ReadWithOptions(ctx context.Context, runEventID string, options *RunEventReadOptions) (*RunEvent, error) {
	if !validStringID(&runEventID) {
		return nil, ErrInvalidRunEventID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("run-events/%s", url.PathEscape(runEventID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	r := &RunEvent{}
	err = req.Do(ctx, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (o *RunEventReadOptions) valid() error {
	return nil
}

func (o *RunEventListOptions) valid() error {
	return nil
}
