// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Applies = (*applies)(nil)

// Applies describes all the apply related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/applies
type Applies interface {
	// Read an apply by its ID.
	Read(ctx context.Context, applyID string) (*Apply, error)

	// Logs retrieves the logs of an apply.
	Logs(ctx context.Context, applyID string) (io.Reader, error)
}

// applies implements Applies interface.
type applies struct {
	client *Client
}

// ApplyStatus represents an apply state.
type ApplyStatus string

// List all available apply statuses.
const (
	ApplyCanceled    ApplyStatus = "canceled"
	ApplyCreated     ApplyStatus = "created"
	ApplyErrored     ApplyStatus = "errored"
	ApplyFinished    ApplyStatus = "finished"
	ApplyMFAWaiting  ApplyStatus = "mfa_waiting"
	ApplyPending     ApplyStatus = "pending"
	ApplyQueued      ApplyStatus = "queued"
	ApplyRunning     ApplyStatus = "running"
	ApplyUnreachable ApplyStatus = "unreachable"
)

// Apply represents a Terraform Enterprise apply.
type Apply struct {
	ID                   string                 `jsonapi:"primary,applies"`
	LogReadURL           string                 `jsonapi:"attr,log-read-url"`
	ResourceAdditions    int                    `jsonapi:"attr,resource-additions"`
	ResourceChanges      int                    `jsonapi:"attr,resource-changes"`
	ResourceDestructions int                    `jsonapi:"attr,resource-destructions"`
	ResourceImports      int                    `jsonapi:"attr,resource-imports"`
	Status               ApplyStatus            `jsonapi:"attr,status"`
	StatusTimestamps     *ApplyStatusTimestamps `jsonapi:"attr,status-timestamps"`
}

// ApplyStatusTimestamps holds the timestamps for individual apply statuses.
type ApplyStatusTimestamps struct {
	CanceledAt      time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	ErroredAt       time.Time `jsonapi:"attr,errored-at,rfc3339"`
	FinishedAt      time.Time `jsonapi:"attr,finished-at,rfc3339"`
	ForceCanceledAt time.Time `jsonapi:"attr,force-canceled-at,rfc3339"`
	QueuedAt        time.Time `jsonapi:"attr,queued-at,rfc3339"`
	StartedAt       time.Time `jsonapi:"attr,started-at,rfc3339"`
}

// Read an apply by its ID.
func (s *applies) Read(ctx context.Context, applyID string) (*Apply, error) {
	if !validStringID(&applyID) {
		return nil, ErrInvalidApplyID
	}

	u := fmt.Sprintf("applies/%s", url.PathEscape(applyID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	a := &Apply{}
	err = req.Do(ctx, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Logs retrieves the logs of an apply.
func (s *applies) Logs(ctx context.Context, applyID string) (io.Reader, error) {
	if !validStringID(&applyID) {
		return nil, ErrInvalidApplyID
	}

	// Get the apply to make sure it exists.
	a, err := s.Read(ctx, applyID)
	if err != nil {
		return nil, err
	}

	// Return an error if the log URL is empty.
	if a.LogReadURL == "" {
		return nil, fmt.Errorf("apply %s does not have a log URL", applyID)
	}

	u, err := url.Parse(a.LogReadURL)
	if err != nil {
		return nil, fmt.Errorf("invalid log URL: %w", err)
	}

	done := func() (bool, error) {
		a, err := s.Read(ctx, a.ID)
		if err != nil {
			return false, err
		}

		switch a.Status {
		case ApplyCanceled, ApplyErrored, ApplyFinished, ApplyUnreachable:
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
