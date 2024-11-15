// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// NOTE WELL: This is a beta feature and is subject to change until noted otherwise in the
// release notes.
type StackPlanOperations interface {
	// Read returns a stack plan operation by its ID.
	Read(ctx context.Context, stackPlanOperationID string) (*StackPlanOperation, error)

	// Get Stack Plans from Configuration Version
	DownloadEventStream(ctx context.Context, stackPlanOperationID string) ([]byte, error)
}

type stackPlanOperations struct {
	client *Client
}

var _ StackPlanOperations = &stackPlanOperations{}

type StackPlanOperation struct {
	ID             string             `jsonapi:"primary,stack-plan-operations"`
	Type           string             `jsonapi:"attr,operation-type"`
	Status         string             `jsonapi:"attr,status"`
	EventStreamURL string             `jsonapi:"attr,event-stream-url"`
	Diagnostics    []*StackDiagnostic `jsonapi:"attr,diags"`

	// Relations
	StackPlan *StackPlan `jsonapi:"relation,stack-plan"`
}

func (s stackPlanOperations) Read(ctx context.Context, stackPlanOperationID string) (*StackPlanOperation, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("stack-plan-operations/%s", url.PathEscape(stackPlanOperationID)), nil)
	if err != nil {
		return nil, err
	}

	spo := &StackPlanOperation{}
	err = req.Do(ctx, spo)
	if err != nil {
		return nil, err
	}

	return spo, nil
}

func (s stackPlanOperations) DownloadEventStream(ctx context.Context, eventStreamURL string) ([]byte, error) {
	// Create a new request.
	req, err := http.NewRequest("GET", eventStreamURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Attach the default headers.
	for k, v := range s.client.headers {
		req.Header[k] = v
	}

	// Retrieve the next chunk.
	resp, err := s.client.http.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Basic response checking.
	if err := checkResponseCode(resp); err != nil {
		return nil, err
	}

	// Read the retrieved chunk.
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
