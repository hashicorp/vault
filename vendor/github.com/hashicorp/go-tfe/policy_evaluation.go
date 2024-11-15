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
var _ PolicyEvaluations = (*policyEvaluation)(nil)

// PolicyEvaluationStatus is an enum that represents all possible statuses for a policy evaluation
type PolicyEvaluationStatus string

const (
	PolicyEvaluationPassed      PolicyEvaluationStatus = "passed"
	PolicyEvaluationFailed      PolicyEvaluationStatus = "failed"
	PolicyEvaluationPending     PolicyEvaluationStatus = "pending"
	PolicyEvaluationRunning     PolicyEvaluationStatus = "running"
	PolicyEvaluationUnreachable PolicyEvaluationStatus = "unreachable"
	PolicyEvaluationOverridden  PolicyEvaluationStatus = "overridden"
	PolicyEvaluationCanceled    PolicyEvaluationStatus = "canceled"
	PolicyEvaluationErrored     PolicyEvaluationStatus = "errored"
)

// PolicyResultCount represents the count of the policy results
type PolicyResultCount struct {
	AdvisoryFailed  int `jsonapi:"attr,advisory-failed"`
	MandatoryFailed int `jsonapi:"attr,mandatory-failed"`
	Passed          int `jsonapi:"attr,passed"`
	Errored         int `jsonapi:"attr,errored"`
}

// The task stage the policy evaluation belongs to
type PolicyAttachable struct {
	ID   string `jsonapi:"attr,id"`
	Type string `jsonapi:"attr,type"`
}

// PolicyEvaluationStatusTimestamps represents the set of timestamps recorded for a policy evaluation
type PolicyEvaluationStatusTimestamps struct {
	ErroredAt  time.Time `jsonapi:"attr,errored-at,rfc3339"`
	RunningAt  time.Time `jsonapi:"attr,running-at,rfc3339"`
	CanceledAt time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	FailedAt   time.Time `jsonapi:"attr,failed-at,rfc3339"`
	PassedAt   time.Time `jsonapi:"attr,passed-at,rfc3339"`
}

// PolicyEvaluation represents the policy evaluations that are part of the task stage.
type PolicyEvaluation struct {
	ID               string                           `jsonapi:"primary,policy-evaluations"`
	Status           PolicyEvaluationStatus           `jsonapi:"attr,status"`
	PolicyKind       PolicyKind                       `jsonapi:"attr,policy-kind"`
	StatusTimestamps PolicyEvaluationStatusTimestamps `jsonapi:"attr,status-timestamps"`
	ResultCount      *PolicyResultCount               `jsonapi:"attr,result-count"`
	CreatedAt        time.Time                        `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt        time.Time                        `jsonapi:"attr,updated-at,iso8601"`

	// The task stage this evaluation belongs to
	TaskStage *PolicyAttachable `jsonapi:"relation,policy-attachable"`
}

// PolicyEvalutations describes all the policy evaluation related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-checks
type PolicyEvaluations interface {
	// **Note: This method is still in BETA and subject to change.**
	// List all policy evaluations in the task stage. Only available for OPA policies.
	List(ctx context.Context, taskStageID string, options *PolicyEvaluationListOptions) (*PolicyEvaluationList, error)
}

// policyEvaluation implements PolicyEvaluations.
type policyEvaluation struct {
	client *Client
}

// PolicyEvaluationListOptions represents the options for listing policy evaluations.
type PolicyEvaluationListOptions struct {
	ListOptions
}

// PolicyEvaluationList represents a list of policy evaluation.
type PolicyEvaluationList struct {
	*Pagination
	Items []*PolicyEvaluation
}

// List all policy evaluations in a task stage.
func (s *policyEvaluation) List(ctx context.Context, taskStageID string, options *PolicyEvaluationListOptions) (*PolicyEvaluationList, error) {
	if !validStringID(&taskStageID) {
		return nil, ErrInvalidTaskStageID
	}

	u := fmt.Sprintf("task-stages/%s/policy-evaluations", url.PathEscape(taskStageID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	pcl := &PolicyEvaluationList{}
	err = req.Do(ctx, pcl)
	if err != nil {
		return nil, err
	}

	return pcl, nil
}

// Compile-time proof of interface implementation.
var _ PolicySetOutcomes = (*policySetOutcome)(nil)

// PolicySetOutcomes describes all the policy set outcome related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-checks
type PolicySetOutcomes interface {
	// **Note: This method is still in BETA and subject to change.**
	// List all policy set outcomes in the policy evaluation. Only available for OPA policies.
	List(ctx context.Context, policyEvaluationID string, options *PolicySetOutcomeListOptions) (*PolicySetOutcomeList, error)

	// **Note: This method is still in BETA and subject to change.**
	// Read a policy set outcome by its ID. Only available for OPA policies.
	Read(ctx context.Context, policySetOutcomeID string) (*PolicySetOutcome, error)
}

// policySetOutcome implements PolicySetOutcomes.
type policySetOutcome struct {
	client *Client
}

// PolicySetOutcomeListFilter represents the filters that are supported while listing a policy set outcome
type PolicySetOutcomeListFilter struct {
	// Optional: A status string used to filter the results.
	// Must be either "passed", "failed", or "errored".
	Status string

	// Optional: The enforcement level used to filter the results.
	// Must be either "advisory" or "mandatory".
	EnforcementLevel string
}

// PolicySetOutcomeListOptions represents the options for listing policy set outcomes.
type PolicySetOutcomeListOptions struct {
	*ListOptions

	// Optional: A filter map used to filter the results of the policy outcome.
	// You can use filter[n] to combine combinations of statuses and enforcement levels filters
	Filter map[string]PolicySetOutcomeListFilter
}

// PolicySetOutcomeList represents a list of policy set outcomes.
type PolicySetOutcomeList struct {
	*Pagination
	Items []*PolicySetOutcome
}

// Outcome represents the outcome of the individual policy
type Outcome struct {
	EnforcementLevel EnforcementLevel `jsonapi:"attr,enforcement_level"`
	Query            string           `jsonapi:"attr,query"`
	Status           string           `jsonapi:"attr,status"`
	PolicyName       string           `jsonapi:"attr,policy_name"`
	Description      string           `jsonapi:"attr,description"`
}

// PolicySetOutcome represents outcome of the policy set that are part of the policy evaluation
type PolicySetOutcome struct {
	ID                   string            `jsonapi:"primary,policy-set-outcomes"`
	Outcomes             []Outcome         `jsonapi:"attr,outcomes"`
	Error                string            `jsonapi:"attr,error"`
	Overridable          *bool             `jsonapi:"attr,overridable"`
	PolicySetName        string            `jsonapi:"attr,policy-set-name"`
	PolicySetDescription string            `jsonapi:"attr,policy-set-description"`
	ResultCount          PolicyResultCount `jsonapi:"attr,result_count"`

	// The policy evaluation that this outcome belongs to
	PolicyEvaluation *PolicyEvaluation `jsonapi:"relation,policy-evaluation"`
}

// List all policy set outcomes in a policy evaluation.
func (s *policySetOutcome) List(ctx context.Context, policyEvaluationID string, options *PolicySetOutcomeListOptions) (*PolicySetOutcomeList, error) {
	if !validStringID(&policyEvaluationID) {
		return nil, ErrInvalidPolicyEvaluationID
	}

	additionalQueryParams := options.buildQueryString()

	u := fmt.Sprintf("policy-evaluations/%s/policy-set-outcomes", url.QueryEscape(policyEvaluationID))

	var opts *ListOptions
	if options != nil && options.ListOptions != nil {
		opts = options.ListOptions
	}

	req, err := s.client.NewRequestWithAdditionalQueryParams("GET", u, opts, additionalQueryParams)
	if err != nil {
		return nil, err
	}

	psol := &PolicySetOutcomeList{}
	err = req.Do(ctx, psol)
	if err != nil {
		return nil, err
	}

	return psol, nil
}

// buildQueryString takes the PolicySetOutcomeListOptions and returns a filters map.
// This function is required due to the limitations of the current library,
// we cannot encode map of objects using the current library that is used by go-tfe: https://github.com/google/go-querystring/issues/7
func (opts *PolicySetOutcomeListOptions) buildQueryString() map[string][]string {
	result := make(map[string][]string)
	if opts == nil || opts.Filter == nil {
		return nil
	}
	for k, v := range opts.Filter {
		if v.Status != "" {
			newKey := fmt.Sprintf("filter[%s][status]", k)
			result[newKey] = append(result[newKey], v.Status)
		}
		if v.EnforcementLevel != "" {
			newKey := fmt.Sprintf("filter[%s][enforcement_level]", k)
			result[newKey] = append(result[newKey], v.EnforcementLevel)
		}
	}
	return result
}

// Read reads a policy set outcome by its ID
func (s *policySetOutcome) Read(ctx context.Context, policySetOutcomeID string) (*PolicySetOutcome, error) {
	if !validStringID(&policySetOutcomeID) {
		return nil, ErrInvalidPolicySetOutcomeID
	}

	u := fmt.Sprintf("policy-set-outcomes/%s", url.PathEscape(policySetOutcomeID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	pso := &PolicySetOutcome{}
	err = req.Do(ctx, pso)
	if err != nil {
		return nil, err
	}

	return pso, err
}
