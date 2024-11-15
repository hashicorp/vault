// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ PolicySetParameters = (*policySetParameters)(nil)

// PolicySetParameters describes all the parameter related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-set-params
type PolicySetParameters interface {
	// List all the parameters associated with the given policy-set.
	List(ctx context.Context, policySetID string, options *PolicySetParameterListOptions) (*PolicySetParameterList, error)

	// Create is used to create a new parameter.
	Create(ctx context.Context, policySetID string, options PolicySetParameterCreateOptions) (*PolicySetParameter, error)

	// Read a parameter by its ID.
	Read(ctx context.Context, policySetID string, parameterID string) (*PolicySetParameter, error)

	// Update values of an existing parameter.
	Update(ctx context.Context, policySetID string, parameterID string, options PolicySetParameterUpdateOptions) (*PolicySetParameter, error)

	// Delete a parameter by its ID.
	Delete(ctx context.Context, policySetID string, parameterID string) error
}

// policySetParameters implements Parameters.
type policySetParameters struct {
	client *Client
}

// PolicySetParameterList represents a list of parameters.
type PolicySetParameterList struct {
	*Pagination
	Items []*PolicySetParameter
}

// PolicySetParameter represents a Policy Set parameter
type PolicySetParameter struct {
	ID        string       `jsonapi:"primary,vars"`
	Key       string       `jsonapi:"attr,key"`
	Value     string       `jsonapi:"attr,value"`
	Category  CategoryType `jsonapi:"attr,category"`
	Sensitive bool         `jsonapi:"attr,sensitive"`

	// Relations
	PolicySet *PolicySet `jsonapi:"relation,configurable"`
}

// PolicySetParameterListOptions represents the options for listing parameters.
type PolicySetParameterListOptions struct {
	ListOptions
}

// PolicySetParameterCreateOptions represents the options for creating a new parameter.
type PolicySetParameterCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,vars"`

	// Required: The name of the parameter.
	Key *string `jsonapi:"attr,key"`

	// Optional: The value of the parameter.
	Value *string `jsonapi:"attr,value,omitempty"`

	// Required: The Category of the parameter, should always be "policy-set"
	Category *CategoryType `jsonapi:"attr,category"`

	// Optional: Whether the value is sensitive.
	Sensitive *bool `jsonapi:"attr,sensitive,omitempty"`
}

// PolicySetParameterUpdateOptions represents the options for updating a parameter.
type PolicySetParameterUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,vars"`

	// Optional: The name of the parameter.
	Key *string `jsonapi:"attr,key,omitempty"`

	// Optional: The value of the parameter.
	Value *string `jsonapi:"attr,value,omitempty"`

	// Optional: Whether the value is sensitive.
	Sensitive *bool `jsonapi:"attr,sensitive,omitempty"`
}

// List all the parameters associated with the given policy-set.
func (s *policySetParameters) List(ctx context.Context, policySetID string, options *PolicySetParameterListOptions) (*PolicySetParameterList, error) {
	if !validStringID(&policySetID) {
		return nil, ErrInvalidPolicySetID
	}

	u := fmt.Sprintf("policy-sets/%s/parameters", policySetID)
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	vl := &PolicySetParameterList{}
	err = req.Do(ctx, vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// Create is used to create a new parameter.
func (s *policySetParameters) Create(ctx context.Context, policySetID string, options PolicySetParameterCreateOptions) (*PolicySetParameter, error) {
	if !validStringID(&policySetID) {
		return nil, ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("policy-sets/%s/parameters", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	p := &PolicySetParameter{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Read a parameter by its ID.
func (s *policySetParameters) Read(ctx context.Context, policySetID, parameterID string) (*PolicySetParameter, error) {
	if !validStringID(&policySetID) {
		return nil, ErrInvalidPolicySetID
	}
	if !validStringID(&parameterID) {
		return nil, ErrInvalidParamID
	}

	u := fmt.Sprintf("policy-sets/%s/parameters/%s", url.PathEscape(policySetID), url.PathEscape(parameterID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	p := &PolicySetParameter{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, err
}

// Update values of an existing parameter.
func (s *policySetParameters) Update(ctx context.Context, policySetID, parameterID string, options PolicySetParameterUpdateOptions) (*PolicySetParameter, error) {
	if !validStringID(&policySetID) {
		return nil, ErrInvalidPolicySetID
	}
	if !validStringID(&parameterID) {
		return nil, ErrInvalidParamID
	}

	u := fmt.Sprintf("policy-sets/%s/parameters/%s", url.PathEscape(policySetID), url.PathEscape(parameterID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	p := &PolicySetParameter{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Delete a parameter by its ID.
func (s *policySetParameters) Delete(ctx context.Context, policySetID, parameterID string) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if !validStringID(&parameterID) {
		return ErrInvalidParamID
	}

	u := fmt.Sprintf("policy-sets/%s/parameters/%s", url.PathEscape(policySetID), url.PathEscape(parameterID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o PolicySetParameterCreateOptions) valid() error {
	if !validString(o.Key) {
		return ErrRequiredKey
	}
	if o.Category == nil {
		return ErrRequiredCategory
	}
	if *o.Category != CategoryPolicySet {
		return ErrInvalidCategory
	}
	return nil
}
