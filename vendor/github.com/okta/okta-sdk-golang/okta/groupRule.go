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

type GroupRuleResource resource

type GroupRule struct {
	Embedded       interface{}          `json:"_embedded,omitempty"`
	Actions        *GroupRuleAction     `json:"actions,omitempty"`
	AllGroupsValid *bool                `json:"allGroupsValid,omitempty"`
	Conditions     *GroupRuleConditions `json:"conditions,omitempty"`
	Created        *time.Time           `json:"created,omitempty"`
	Id             string               `json:"id,omitempty"`
	LastUpdated    *time.Time           `json:"lastUpdated,omitempty"`
	Name           string               `json:"name,omitempty"`
	Status         string               `json:"status,omitempty"`
	Type           string               `json:"type,omitempty"`
}

func (m *GroupRuleResource) UpdateRule(ruleId string, body GroupRule) (*GroupRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v", ruleId)
	req, err := m.client.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var groupRule *GroupRule
	resp, err := m.client.requestExecutor.Do(req, &groupRule)
	if err != nil {
		return nil, resp, err
	}
	return groupRule, resp, nil
}
func (m *GroupRuleResource) DeleteRule(ruleId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v", ruleId)
	if qp != nil {
		url = url + qp.String()
	}
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
