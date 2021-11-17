package tfe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ PolicyChecks = (*policyChecks)(nil)

// PolicyChecks describes all the policy check related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/cloud/api/policy-checks.html
type PolicyChecks interface {
	// List all policy checks of the given run.
	List(ctx context.Context, runID string, options PolicyCheckListOptions) (*PolicyCheckList, error)

	// Read a policy check by its ID.
	Read(ctx context.Context, policyCheckID string) (*PolicyCheck, error)

	// Override a soft-mandatory or warning policy.
	Override(ctx context.Context, policyCheckID string) (*PolicyCheck, error)

	// Logs retrieves the logs of a policy check.
	Logs(ctx context.Context, policyCheckID string) (io.Reader, error)
}

// policyChecks implements PolicyChecks.
type policyChecks struct {
	client *Client
}

// PolicyScope represents a policy scope.
type PolicyScope string

// List all available policy scopes.
const (
	PolicyScopeOrganization PolicyScope = "organization"
	PolicyScopeWorkspace    PolicyScope = "workspace"
)

// PolicyStatus represents a policy check state.
type PolicyStatus string

//List all available policy check statuses.
const (
	PolicyCanceled    PolicyStatus = "canceled"
	PolicyErrored     PolicyStatus = "errored"
	PolicyHardFailed  PolicyStatus = "hard_failed"
	PolicyOverridden  PolicyStatus = "overridden"
	PolicyPasses      PolicyStatus = "passed"
	PolicyPending     PolicyStatus = "pending"
	PolicyQueued      PolicyStatus = "queued"
	PolicySoftFailed  PolicyStatus = "soft_failed"
	PolicyUnreachable PolicyStatus = "unreachable"
)

// PolicyCheckList represents a list of policy checks.
type PolicyCheckList struct {
	*Pagination
	Items []*PolicyCheck
}

// PolicyCheck represents a Terraform Enterprise policy check..
type PolicyCheck struct {
	ID               string                  `jsonapi:"primary,policy-checks"`
	Actions          *PolicyActions          `jsonapi:"attr,actions"`
	Permissions      *PolicyPermissions      `jsonapi:"attr,permissions"`
	Result           *PolicyResult           `jsonapi:"attr,result"`
	Scope            PolicyScope             `jsonapi:"attr,scope"`
	Status           PolicyStatus            `jsonapi:"attr,status"`
	StatusTimestamps *PolicyStatusTimestamps `jsonapi:"attr,status-timestamps"`
	Run              *Run                    `jsonapi:"relation,run"`
}

// PolicyActions represents the policy check actions.
type PolicyActions struct {
	IsOverridable bool `jsonapi:"attr,is-overridable"`
}

// PolicyPermissions represents the policy check permissions.
type PolicyPermissions struct {
	CanOverride bool `jsonapi:"attr,can-override"`
}

// PolicyResult represents the complete policy check result,
type PolicyResult struct {
	AdvisoryFailed int  `jsonapi:"attr,advisory-failed"`
	Duration       int  `jsonapi:"attr,duration"`
	HardFailed     int  `jsonapi:"attr,hard-failed"`
	Passed         int  `jsonapi:"attr,passed"`
	Result         bool `jsonapi:"attr,result"`
	SoftFailed     int  `jsonapi:"attr,soft-failed"`
	TotalFailed    int  `jsonapi:"attr,total-failed"`
}

// PolicyStatusTimestamps holds the timestamps for individual policy check
// statuses.
type PolicyStatusTimestamps struct {
	ErroredAt    time.Time `jsonapi:"attr,errored-at,rfc3339"`
	HardFailedAt time.Time `jsonapi:"attr,hard-failed-at,rfc3339"`
	PassedAt     time.Time `jsonapi:"attr,passed-at,rfc3339"`
	QueuedAt     time.Time `jsonapi:"attr,queued-at,rfc3339"`
	SoftFailedAt time.Time `jsonapi:"attr,soft-failed-at,rfc3339"`
}

// PolicyCheckListOptions represents the options for listing policy checks.
type PolicyCheckListOptions struct {
	ListOptions
}

// List all policy checks of the given run.
func (s *policyChecks) List(ctx context.Context, runID string, options PolicyCheckListOptions) (*PolicyCheckList, error) {
	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/policy-checks", url.QueryEscape(runID))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	pcl := &PolicyCheckList{}
	err = s.client.do(ctx, req, pcl)
	if err != nil {
		return nil, err
	}

	return pcl, nil
}

// Read a policy check by its ID.
func (s *policyChecks) Read(ctx context.Context, policyCheckID string) (*PolicyCheck, error) {
	if !validStringID(&policyCheckID) {
		return nil, errors.New("invalid value for policy check ID")
	}

	u := fmt.Sprintf("policy-checks/%s", url.QueryEscape(policyCheckID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	pc := &PolicyCheck{}
	err = s.client.do(ctx, req, pc)
	if err != nil {
		return nil, err
	}

	return pc, nil
}

// Override a soft-mandatory or warning policy.
func (s *policyChecks) Override(ctx context.Context, policyCheckID string) (*PolicyCheck, error) {
	if !validStringID(&policyCheckID) {
		return nil, errors.New("invalid value for policy check ID")
	}

	u := fmt.Sprintf("policy-checks/%s/actions/override", url.QueryEscape(policyCheckID))
	req, err := s.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	pc := &PolicyCheck{}
	err = s.client.do(ctx, req, pc)
	if err != nil {
		return nil, err
	}

	return pc, nil
}

// Logs retrieves the logs of a policy check.
func (s *policyChecks) Logs(ctx context.Context, policyCheckID string) (io.Reader, error) {
	if !validStringID(&policyCheckID) {
		return nil, errors.New("invalid value for policy check ID")
	}

	// Loop until the context is canceled or the policy check is finished
	// running. The policy check logs are not streamed and so only available
	// once the check is finished.
	for {
		pc, err := s.Read(ctx, policyCheckID)
		if err != nil {
			return nil, err
		}

		switch pc.Status {
		case PolicyPending, PolicyQueued:
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(500 * time.Millisecond):
				continue
			}
		}

		u := fmt.Sprintf("policy-checks/%s/output", url.QueryEscape(policyCheckID))
		req, err := s.client.newRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}

		logs := bytes.NewBuffer(nil)
		err = s.client.do(ctx, req, logs)
		if err != nil {
			return nil, err
		}

		return logs, nil
	}
}
