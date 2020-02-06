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
	"time"
)

type PolicyRuleResource resource

type PolicyRule struct {
	Created     *time.Time `json:"created,omitempty"`
	Id          string     `json:"id,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Priority    int64      `json:"priority,omitempty"`
	Status      string     `json:"status,omitempty"`
	System      *bool      `json:"system,omitempty"`
	Type        string     `json:"type,omitempty"`
}

func (m *PolicyRuleResource) UpdatePolicyRule(policyId string, ruleId string, body PolicyRule) (*PolicyRule, *Response, error) {
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
func (m *PolicyRuleResource) DeletePolicyRule(policyId string, ruleId string) (*Response, error) {
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
