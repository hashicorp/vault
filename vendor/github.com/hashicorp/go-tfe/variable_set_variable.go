// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ VariableSetVariables = (*variableSetVariables)(nil)

// VariableSetVariables describes all variable variable related methods within the scope of
// Variable Sets that the Terraform Enterprise API supports
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/variable-sets#variable-relationships
type VariableSetVariables interface {
	// List all variables in the variable set.
	List(ctx context.Context, variableSetID string, options *VariableSetVariableListOptions) (*VariableSetVariableList, error)

	// Create is used to create a new variable within a given variable set
	Create(ctx context.Context, variableSetID string, options *VariableSetVariableCreateOptions) (*VariableSetVariable, error)

	// Read a variable by its ID
	Read(ctx context.Context, variableSetID string, variableID string) (*VariableSetVariable, error)

	// Update valuse of an existing variable
	Update(ctx context.Context, variableSetID string, variableID string, options *VariableSetVariableUpdateOptions) (*VariableSetVariable, error)

	// Delete a variable by its ID
	Delete(ctx context.Context, variableSetID string, variableID string) error
}

type variableSetVariables struct {
	client *Client
}

type VariableSetVariableList struct {
	*Pagination
	Items []*VariableSetVariable
}

type VariableSetVariable struct {
	ID          string       `jsonapi:"primary,vars"`
	Key         string       `jsonapi:"attr,key"`
	Value       string       `jsonapi:"attr,value"`
	Description string       `jsonapi:"attr,description"`
	Category    CategoryType `jsonapi:"attr,category"`
	HCL         bool         `jsonapi:"attr,hcl"`
	Sensitive   bool         `jsonapi:"attr,sensitive"`
	VersionID   string       `jsonapi:"attr,version-id"`

	// Relations
	VariableSet *VariableSet `jsonapi:"relation,varset"`
}

type VariableSetVariableListOptions struct {
	ListOptions
}

func (o VariableSetVariableListOptions) valid() error {
	return nil
}

// List all variables associated with the given variable set.
func (s *variableSetVariables) List(ctx context.Context, variableSetID string, options *VariableSetVariableListOptions) (*VariableSetVariableList, error) {
	if !validStringID(&variableSetID) {
		return nil, ErrInvalidVariableSetID
	}
	if options != nil {
		if err := options.valid(); err != nil {
			return nil, err
		}
	}

	u := fmt.Sprintf("varsets/%s/relationships/vars", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	vl := &VariableSetVariableList{}
	err = req.Do(ctx, vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// VariableSetVariableCreatOptions represents the options for creating a new variable within a variable set
type VariableSetVariableCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,vars"`

	// The name of the variable.
	Key *string `jsonapi:"attr,key"`

	// The value of the variable.
	Value *string `jsonapi:"attr,value,omitempty"`

	// The description of the variable.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Whether this is a Terraform or environment variable.
	Category *CategoryType `jsonapi:"attr,category"`

	// Whether to evaluate the value of the variable as a string of HCL code.
	HCL *bool `jsonapi:"attr,hcl,omitempty"`

	// Whether the value is sensitive.
	Sensitive *bool `jsonapi:"attr,sensitive,omitempty"`
}

func (o VariableSetVariableCreateOptions) valid() error {
	if !validString(o.Key) {
		return ErrRequiredKey
	}
	if o.Category == nil {
		return ErrRequiredCategory
	}
	return nil
}

// Create is used to create a new variable.
func (s *variableSetVariables) Create(ctx context.Context, variableSetID string, options *VariableSetVariableCreateOptions) (*VariableSetVariable, error) {
	if !validStringID(&variableSetID) {
		return nil, ErrInvalidVariableSetID
	}
	if options != nil {
		if err := options.valid(); err != nil {
			return nil, err
		}
	}

	u := fmt.Sprintf("varsets/%s/relationships/vars", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("POST", u, options)
	if err != nil {
		return nil, err
	}

	v := &VariableSetVariable{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Read a variable by its ID.
func (s *variableSetVariables) Read(ctx context.Context, variableSetID, variableID string) (*VariableSetVariable, error) {
	if !validStringID(&variableSetID) {
		return nil, ErrInvalidVariableSetID
	}
	if !validStringID(&variableID) {
		return nil, ErrInvalidVariableID
	}

	u := fmt.Sprintf("varsets/%s/relationships/vars/%s", url.PathEscape(variableSetID), url.PathEscape(variableID))
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, err
	}

	v := &VariableSetVariable{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

// VariableSetVariableUpdateOptions represents the options for updating a variable.
type VariableSetVariableUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,vars"`

	// The name of the variable.
	Key *string `jsonapi:"attr,key,omitempty"`

	// The value of the variable.
	Value *string `jsonapi:"attr,value,omitempty"`

	// The description of the variable.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Whether to evaluate the value of the variable as a string of HCL code.
	HCL *bool `jsonapi:"attr,hcl,omitempty"`

	// Whether the value is sensitive.
	Sensitive *bool `jsonapi:"attr,sensitive,omitempty"`
}

// Update values of an existing variable.
func (s *variableSetVariables) Update(ctx context.Context, variableSetID, variableID string, options *VariableSetVariableUpdateOptions) (*VariableSetVariable, error) {
	if !validStringID(&variableSetID) {
		return nil, ErrInvalidVariableSetID
	}
	if !validStringID(&variableID) {
		return nil, ErrInvalidVariableID
	}

	u := fmt.Sprintf("varsets/%s/relationships/vars/%s", url.PathEscape(variableSetID), url.PathEscape(variableID))
	req, err := s.client.NewRequest("PATCH", u, options)
	if err != nil {
		return nil, err
	}

	v := &VariableSetVariable{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Delete a variable by its ID.
func (s *variableSetVariables) Delete(ctx context.Context, variableSetID, variableID string) error {
	if !validStringID(&variableSetID) {
		return ErrInvalidVariableSetID
	}
	if !validStringID(&variableID) {
		return ErrInvalidVariableID
	}

	u := fmt.Sprintf("varsets/%s/relationships/vars/%s", url.PathEscape(variableSetID), url.PathEscape(variableID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}
