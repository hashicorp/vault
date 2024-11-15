// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Policies = (*policies)(nil)

// Policies describes all the policy related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policies
type Policies interface {
	// List all the policies for a given organization
	List(ctx context.Context, organization string, options *PolicyListOptions) (*PolicyList, error)

	// Create a policy and associate it with an organization.
	Create(ctx context.Context, organization string, options PolicyCreateOptions) (*Policy, error)

	// Read a policy by its ID.
	Read(ctx context.Context, policyID string) (*Policy, error)

	// Update an existing policy.
	Update(ctx context.Context, policyID string, options PolicyUpdateOptions) (*Policy, error)

	// Delete a policy by its ID.
	Delete(ctx context.Context, policyID string) error

	// Upload the policy content of the policy.
	Upload(ctx context.Context, policyID string, content []byte) error

	// Download the policy content of the policy.
	Download(ctx context.Context, policyID string) ([]byte, error)
}

// policies implements Policies.
type policies struct {
	client *Client
}

// EnforcementLevel represents an enforcement level.
type EnforcementLevel string

// List the available enforcement types.
const (
	EnforcementAdvisory  EnforcementLevel = "advisory"
	EnforcementHard      EnforcementLevel = "hard-mandatory"
	EnforcementSoft      EnforcementLevel = "soft-mandatory"
	EnforcementMandatory EnforcementLevel = "mandatory"
)

// PolicyList represents a list of policies..
type PolicyList struct {
	*Pagination
	Items []*Policy
}

// Policy represents a Terraform Enterprise policy.
type Policy struct {
	ID          string     `jsonapi:"primary,policies"`
	Name        string     `jsonapi:"attr,name"`
	Kind        PolicyKind `jsonapi:"attr,kind"`
	Query       *string    `jsonapi:"attr,query"`
	Description string     `jsonapi:"attr,description"`
	// Deprecated: Use EnforcementLevel instead.
	Enforce          []*Enforcement   `jsonapi:"attr,enforce"`
	EnforcementLevel EnforcementLevel `jsonapi:"attr,enforcement-level"`
	PolicySetCount   int              `jsonapi:"attr,policy-set-count"`
	UpdatedAt        time.Time        `jsonapi:"attr,updated-at,iso8601"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
}

// Enforcement describes a enforcement.
type Enforcement struct {
	Path string           `jsonapi:"attr,path"`
	Mode EnforcementLevel `jsonapi:"attr,mode"`
}

// EnforcementOptions represents the enforcement options of a policy.
type EnforcementOptions struct {
	Path *string           `json:"path"`
	Mode *EnforcementLevel `json:"mode"`
}

// PolicyListOptions represents the options for listing policies.
type PolicyListOptions struct {
	ListOptions

	// Optional: A search string (partial policy name) used to filter the results.
	Search string `url:"search[name],omitempty"`

	// Optional: A kind string used to filter the results by the policy kind.
	Kind PolicyKind `url:"filter[kind],omitempty"`
}

// PolicyCreateOptions represents the options for creating a new policy.
type PolicyCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,policies"`

	// Required: The name of the policy.
	Name *string `jsonapi:"attr,name"`

	// Optional: The underlying technology that the policy supports. Defaults to Sentinel if not specified for PolicyCreate.
	Kind PolicyKind `jsonapi:"attr,kind,omitempty"`

	// Optional: The query passed to policy evaluation to determine the result of the policy. Only valid for OPA.
	Query *string `jsonapi:"attr,query,omitempty"`

	// Optional: A description of the policy's purpose.
	Description *string `jsonapi:"attr,description,omitempty"`

	// The enforcements of the policy.
	//
	// Deprecated: Use EnforcementLevel instead.
	Enforce []*EnforcementOptions `jsonapi:"attr,enforce,omitempty"`

	// Required: The enforcement level of the policy.
	// Either EnforcementLevel or Enforce must be set.
	EnforcementLevel *EnforcementLevel `jsonapi:"attr,enforcement-level,omitempty"`
}

// PolicyUpdateOptions represents the options for updating a policy.
type PolicyUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,policies"`

	// Optional: A description of the policy's purpose.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Optional: The query passed to policy evaluation to determine the result of the policy. Only valid for OPA.
	Query *string `jsonapi:"attr,query,omitempty"`

	// Optional: The enforcements of the policy.
	//
	// Deprecated: Use EnforcementLevel instead.
	Enforce []*EnforcementOptions `jsonapi:"attr,enforce,omitempty"`

	// Optional: The enforcement level of the policy.
	EnforcementLevel *EnforcementLevel `jsonapi:"attr,enforcement-level,omitempty"`
}

// List all the policies for a given organization
func (s *policies) List(ctx context.Context, organization string, options *PolicyListOptions) (*PolicyList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/policies", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	pl := &PolicyList{}
	err = req.Do(ctx, pl)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

// Create a policy and associate it with an organization.
func (s *policies) Create(ctx context.Context, organization string, options PolicyCreateOptions) (*Policy, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/policies", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	p := &Policy{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, err
}

// Read a policy by its ID.
func (s *policies) Read(ctx context.Context, policyID string) (*Policy, error) {
	if !validStringID(&policyID) {
		return nil, ErrInvalidPolicyID
	}

	u := fmt.Sprintf("policies/%s", url.PathEscape(policyID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	p := &Policy{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, err
}

// Update an existing policy.
func (s *policies) Update(ctx context.Context, policyID string, options PolicyUpdateOptions) (*Policy, error) {
	if !validStringID(&policyID) {
		return nil, ErrInvalidPolicyID
	}

	u := fmt.Sprintf("policies/%s", url.PathEscape(policyID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	p := &Policy{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, err
}

// Delete a policy by its ID.
func (s *policies) Delete(ctx context.Context, policyID string) error {
	if !validStringID(&policyID) {
		return ErrInvalidPolicyID
	}

	u := fmt.Sprintf("policies/%s", url.PathEscape(policyID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Upload the policy content of the policy.
func (s *policies) Upload(ctx context.Context, policyID string, content []byte) error {
	if !validStringID(&policyID) {
		return ErrInvalidPolicyID
	}

	u := fmt.Sprintf("policies/%s/upload", url.PathEscape(policyID))
	req, err := s.client.NewRequest("PUT", u, content)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Download the policy content of the policy.
func (s *policies) Download(ctx context.Context, policyID string) ([]byte, error) {
	if !validStringID(&policyID) {
		return nil, ErrInvalidPolicyID
	}

	u := fmt.Sprintf("policies/%s/download", url.PathEscape(policyID))
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

func (o PolicyCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validStringID(o.Name) {
		return ErrInvalidName
	}
	if o.Kind == OPA && !validString(o.Query) {
		return ErrRequiredQuery
	}
	if o.Enforce == nil && o.EnforcementLevel == nil {
		return ErrRequiredEnforce
	}
	if o.Enforce != nil && o.EnforcementLevel != nil {
		return ErrConflictingEnforceEnforcementLevel
	}
	if o.Enforce != nil {
		for _, e := range o.Enforce {
			if !validString(e.Path) {
				return ErrRequiredEnforcementPath
			}
			if e.Mode == nil {
				return ErrRequiredEnforcementMode
			}
		}
	}
	return nil
}
