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
var _ RunTriggers = (*runTriggers)(nil)

// RunTriggers describes all the Run Trigger
// related methods that the HCP Terraform API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers
type RunTriggers interface {
	// List all the run triggers within a workspace.
	List(ctx context.Context, workspaceID string, options *RunTriggerListOptions) (*RunTriggerList, error)

	// Create a new run trigger with the given options.
	Create(ctx context.Context, workspaceID string, options RunTriggerCreateOptions) (*RunTrigger, error)

	// Read a run trigger by its ID.
	Read(ctx context.Context, RunTriggerID string) (*RunTrigger, error)

	// Delete a run trigger by its ID.
	Delete(ctx context.Context, RunTriggerID string) error
}

// runTriggers implements RunTriggers.
type runTriggers struct {
	client *Client
}

// RunTriggerList represents a list of Run Triggers
type RunTriggerList struct {
	*Pagination
	Items []*RunTrigger
}

// SourceableChoice is a choice type struct that represents the possible values
// within a polymorphic relation. If a value is available, exactly one field
// will be non-nil.
type SourceableChoice struct {
	Workspace *Workspace
}

// RunTrigger represents a run trigger.
type RunTrigger struct {
	ID             string    `jsonapi:"primary,run-triggers"`
	CreatedAt      time.Time `jsonapi:"attr,created-at,iso8601"`
	SourceableName string    `jsonapi:"attr,sourceable-name"`
	WorkspaceName  string    `jsonapi:"attr,workspace-name"`
	// DEPRECATED. The sourceable field is polymorphic. Use SourceableChoice instead.
	Sourceable       *Workspace        `jsonapi:"relation,sourceable"`
	SourceableChoice *SourceableChoice `jsonapi:"polyrelation,sourceable"`
	Workspace        *Workspace        `jsonapi:"relation,workspace"`
}

// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#query-parameters
type RunTriggerFilterOp string

const (
	RunTriggerOutbound RunTriggerFilterOp = "outbound" // create runs in other workspaces.
	RunTriggerInbound  RunTriggerFilterOp = "inbound"  // create runs in the specified workspace
)

// A list of relations to include
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run-triggers#available-related-resources
type RunTriggerIncludeOpt string

const (
	RunTriggerWorkspace  RunTriggerIncludeOpt = "workspace"
	RunTriggerSourceable RunTriggerIncludeOpt = "sourceable"
)

// RunTriggerListOptions represents the options for listing
// run triggers.
type RunTriggerListOptions struct {
	ListOptions
	RunTriggerType RunTriggerFilterOp     `url:"filter[run-trigger][type]"` // Required
	Include        []RunTriggerIncludeOpt `url:"include,omitempty"`         // optional
}

// RunTriggerCreateOptions represents the options for
// creating a new run trigger.
type RunTriggerCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,run-triggers"`

	// The source workspace
	Sourceable *Workspace `jsonapi:"relation,sourceable"`
}

// List all the run triggers associated with a workspace.
func (s *runTriggers) List(ctx context.Context, workspaceID string, options *RunTriggerListOptions) (*RunTriggerList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/run-triggers", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	rtl := &RunTriggerList{}
	err = req.Do(ctx, rtl)
	if err != nil {
		return nil, err
	}

	return rtl, nil
}

// Create a run trigger with the given options.
func (s *runTriggers) Create(ctx context.Context, workspaceID string, options RunTriggerCreateOptions) (*RunTrigger, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/run-triggers", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rt := &RunTrigger{}
	err = req.Do(ctx, rt)
	if err != nil {
		return nil, err
	}

	return rt, nil
}

// Read a run trigger by its ID.
func (s *runTriggers) Read(ctx context.Context, runTriggerID string) (*RunTrigger, error) {
	if !validStringID(&runTriggerID) {
		return nil, ErrInvalidRunTriggerID
	}

	u := fmt.Sprintf("run-triggers/%s", url.PathEscape(runTriggerID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	rt := &RunTrigger{}
	err = req.Do(ctx, rt)
	if err != nil {
		return nil, err
	}

	return rt, nil
}

// Delete a run trigger by its ID.
func (s *runTriggers) Delete(ctx context.Context, runTriggerID string) error {
	if !validStringID(&runTriggerID) {
		return ErrInvalidRunTriggerID
	}

	u := fmt.Sprintf("run-triggers/%s", url.PathEscape(runTriggerID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o RunTriggerCreateOptions) valid() error {
	if o.Sourceable == nil {
		return ErrRequiredSourceable
	}
	return nil
}

func (o *RunTriggerListOptions) valid() error {
	if o == nil {
		return ErrRequiredRunTriggerListOps
	}

	if err := validateRunTriggerFilterParam(o.RunTriggerType, o.Include); err != nil {
		return err
	}

	return nil
}

func validateRunTriggerFilterParam(filterParam RunTriggerFilterOp, includeParams []RunTriggerIncludeOpt) error {
	switch filterParam {
	case RunTriggerOutbound, RunTriggerInbound:
		// Do nothing
	default:
		return ErrInvalidRunTriggerType // return an error even if string is empty because this a required field
	}

	if len(includeParams) > 0 {
		if filterParam != RunTriggerInbound {
			return ErrUnsupportedRunTriggerType // if user passes RunTriggerOutbound the platform will not return any "include" data
		}
	}

	return nil
}
