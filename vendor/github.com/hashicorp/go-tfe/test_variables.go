// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ TestVariables = (*testVariables)(nil)

// Variables describes all the variable related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/private-registry/tests
type TestVariables interface {
	// List all the test variables associated with the given module.
	List(ctx context.Context, moduleID RegistryModuleID, options *VariableListOptions) (*VariableList, error)

	// Read a test variable by its ID.
	Read(ctx context.Context, moduleID RegistryModuleID, variableID string) (*Variable, error)

	// Create is used to create a new variable.
	Create(ctx context.Context, moduleID RegistryModuleID, options VariableCreateOptions) (*Variable, error)

	// Update values of an existing variable.
	Update(ctx context.Context, moduleID RegistryModuleID, variableID string, options VariableUpdateOptions) (*Variable, error)

	// Delete a variable by its ID.
	Delete(ctx context.Context, moduleID RegistryModuleID, variableID string) error
}

// variables implements Variables.
type testVariables struct {
	client *Client
}

// List all the variables associated with the given module.
func (s *testVariables) List(ctx context.Context, moduleID RegistryModuleID, options *VariableListOptions) (*VariableList, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", testVarsPath(moduleID), options)
	if err != nil {
		return nil, err
	}

	vl := &VariableList{}
	err = req.Do(ctx, vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// Read a variable by its ID.
func (s *testVariables) Read(ctx context.Context, moduleID RegistryModuleID, variableID string) (*Variable, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if !validStringID(&variableID) {
		return nil, ErrInvalidVariableID
	}

	u := fmt.Sprintf("%s/%s", testVarsPath(moduleID), url.PathEscape(variableID))

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	v := &Variable{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

// Create is used to create a new variable.
func (s *testVariables) Create(ctx context.Context, moduleID RegistryModuleID, options VariableCreateOptions) (*Variable, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", testVarsPath(moduleID), &options)
	if err != nil {
		return nil, err
	}

	v := &Variable{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Update values of an existing variable.
func (s *testVariables) Update(ctx context.Context, moduleID RegistryModuleID, variableID string, options VariableUpdateOptions) (*Variable, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}
	if !validStringID(&variableID) {
		return nil, ErrInvalidVariableID
	}

	u := fmt.Sprintf("%s/%s", testVarsPath(moduleID), url.PathEscape(variableID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	v := &Variable{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Delete a variable by its ID.
func (s *testVariables) Delete(ctx context.Context, moduleID RegistryModuleID, variableID string) error {
	if err := moduleID.valid(); err != nil {
		return err
	}
	if !validStringID(&variableID) {
		return ErrInvalidVariableID
	}

	u := fmt.Sprintf("%s/%s", testVarsPath(moduleID), url.PathEscape(variableID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func testVarsPath(moduleID RegistryModuleID) string {
	return fmt.Sprintf("organizations/%s/tests/registry-modules/%s/%s/%s/%s/vars",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(string(moduleID.RegistryName)),
		url.PathEscape(moduleID.Namespace),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider))
}
