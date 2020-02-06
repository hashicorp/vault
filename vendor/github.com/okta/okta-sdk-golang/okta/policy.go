/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package okta

import (
	"fmt"
	"github.com/okta/okta-sdk-golang/okta/query"
	"time"
)

type PolicyResource resource

type Policy struct {
	Embedded    interface{} `json:"_embedded,omitempty"`
	Links       interface{} `json:"_links,omitempty"`
	Created     *time.Time  `json:"created,omitempty"`
	Description string      `json:"description,omitempty"`
	Id          string      `json:"id,omitempty"`
	LastUpdated *time.Time  `json:"lastUpdated,omitempty"`
	Name        string      `json:"name,omitempty"`
	Priority    int64       `json:"priority,omitempty"`
	Status      string      `json:"status,omitempty"`
	System      *bool       `json:"system,omitempty"`
	Type        string      `json:"type,omitempty"`
}

func (m *PolicyResource) GetPolicy(policyId string, qp *query.Params) (*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v", policyId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policy *Policy
	resp, err := m.client.requestExecutor.Do(req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}
func (m *PolicyResource) UpdatePolicy(policyId string, body Policy) (*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v", policyId)
	req, err := m.client.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policy *Policy
	resp, err := m.client.requestExecutor.Do(req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}
func (m *PolicyResource) DeletePolicy(policyId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v", policyId)
	req, err := m.client.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (m *PolicyResource) ListPolicies(qp *query.Params) ([]*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies")
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policy []*Policy
	resp, err := m.client.requestExecutor.Do(req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}
func (m *PolicyResource) CreatePolicy(body Policy, qp *query.Params) (*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies")
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policy *Policy
	resp, err := m.client.requestExecutor.Do(req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}
func (m *PolicyResource) ActivatePolicy(policyId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/lifecycle/activate", policyId)
	req, err := m.client.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (m *PolicyResource) DeactivatePolicy(policyId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/lifecycle/deactivate", policyId)
	req, err := m.client.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (m *PolicyResource) ListPolicyRules(policyId string) ([]*PolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules", policyId)
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policyRule []*PolicyRule
	resp, err := m.client.requestExecutor.Do(req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}
func (m *PolicyResource) AddPolicyRule(policyId string, body PolicyRule, qp *query.Params) (*PolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules", policyId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policyRule *PolicyRule
	resp, err := m.client.requestExecutor.Do(req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}
func (m *PolicyResource) DeletePolicyRule(policyId string, ruleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyId, ruleId)
	req, err := m.client.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (m *PolicyResource) GetPolicyRule(policyId string, ruleId string) (*PolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyId, ruleId)
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policyRule *PolicyRule
	resp, err := m.client.requestExecutor.Do(req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}
func (m *PolicyResource) UpdatePolicyRule(policyId string, ruleId string, body PolicyRule) (*PolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyId, ruleId)
	req, err := m.client.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policyRule *PolicyRule
	resp, err := m.client.requestExecutor.Do(req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}
func (m *PolicyResource) ActivatePolicyRule(policyId string, ruleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v/lifecycle/activate", policyId, ruleId)
	req, err := m.client.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (m *PolicyResource) DeactivatePolicyRule(policyId string, ruleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v/lifecycle/deactivate", policyId, ruleId)
	req, err := m.client.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
