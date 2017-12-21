package api

import "fmt"

// SentinelPolicies is used to query the Sentinel Policy endpoints.
type SentinelPolicies struct {
	client *Client
}

// SentinelPolicies returns a new handle on the Sentinel policies.
func (c *Client) SentinelPolicies() *SentinelPolicies {
	return &SentinelPolicies{client: c}
}

// List is used to dump all of the policies.
func (a *SentinelPolicies) List(q *QueryOptions) ([]*SentinelPolicyListStub, *QueryMeta, error) {
	var resp []*SentinelPolicyListStub
	qm, err := a.client.query("/v1/sentinel/policies", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Upsert is used to create or update a policy
func (a *SentinelPolicies) Upsert(policy *SentinelPolicy, q *WriteOptions) (*WriteMeta, error) {
	if policy == nil || policy.Name == "" {
		return nil, fmt.Errorf("missing policy name")
	}
	wm, err := a.client.write("/v1/sentinel/policy/"+policy.Name, policy, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Delete is used to delete a policy
func (a *SentinelPolicies) Delete(policyName string, q *WriteOptions) (*WriteMeta, error) {
	if policyName == "" {
		return nil, fmt.Errorf("missing policy name")
	}
	wm, err := a.client.delete("/v1/sentinel/policy/"+policyName, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Info is used to query a specific policy
func (a *SentinelPolicies) Info(policyName string, q *QueryOptions) (*SentinelPolicy, *QueryMeta, error) {
	if policyName == "" {
		return nil, nil, fmt.Errorf("missing policy name")
	}
	var resp SentinelPolicy
	wm, err := a.client.query("/v1/sentinel/policy/"+policyName, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

type SentinelPolicy struct {
	Name             string
	Description      string
	Scope            string
	EnforcementLevel string
	Policy           string
	CreateIndex      uint64
	ModifyIndex      uint64
}

type SentinelPolicyListStub struct {
	Name             string
	Description      string
	Scope            string
	EnforcementLevel string
	CreateIndex      uint64
	ModifyIndex      uint64
}
