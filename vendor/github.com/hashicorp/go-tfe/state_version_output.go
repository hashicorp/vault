// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ StateVersionOutputs = (*stateVersionOutputs)(nil)

// State version outputs are the output values from a Terraform state file.
// They include the name and value of the output, as well as a sensitive boolean
// if the value should be hidden by default in UIs.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/state-version-outputs
type StateVersionOutputs interface {
	Read(ctx context.Context, outputID string) (*StateVersionOutput, error)
	ReadCurrent(ctx context.Context, workspaceID string) (*StateVersionOutputsList, error)
}

// stateVersionOutputs implements StateVersionOutputs.
type stateVersionOutputs struct {
	client *Client
}

// StateVersionOutput represents a State Version Outputs
type StateVersionOutput struct {
	ID        string      `jsonapi:"primary,state-version-outputs"`
	Name      string      `jsonapi:"attr,name"`
	Sensitive bool        `jsonapi:"attr,sensitive"`
	Type      string      `jsonapi:"attr,type"`
	Value     interface{} `jsonapi:"attr,value"`
	// BETA: This field is experimental and not universally present in all versions of TFE/Terraform
	DetailedType interface{} `jsonapi:"attr,detailed-type"`
}

// ReadCurrent reads the current state version outputs for the specified workspace
func (s *stateVersionOutputs) ReadCurrent(ctx context.Context, workspaceID string) (*StateVersionOutputsList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/current-state-version-outputs", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	so := &StateVersionOutputsList{}
	err = req.Do(ctx, so)
	if err != nil {
		return nil, err
	}

	return so, nil
}

// Read a State Version Output
func (s *stateVersionOutputs) Read(ctx context.Context, outputID string) (*StateVersionOutput, error) {
	if !validStringID(&outputID) {
		return nil, ErrInvalidOutputID
	}

	u := fmt.Sprintf("state-version-outputs/%s", url.PathEscape(outputID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	so := &StateVersionOutput{}
	err = req.Do(ctx, so)
	if err != nil {
		return nil, err
	}

	return so, nil
}
