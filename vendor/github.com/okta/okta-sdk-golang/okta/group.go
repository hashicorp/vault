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

type GroupResource resource

type Group struct {
	Embedded              interface{}   `json:"_embedded,omitempty"`
	Links                 interface{}   `json:"_links,omitempty"`
	Created               *time.Time    `json:"created,omitempty"`
	Id                    string        `json:"id,omitempty"`
	LastMembershipUpdated *time.Time    `json:"lastMembershipUpdated,omitempty"`
	LastUpdated           *time.Time    `json:"lastUpdated,omitempty"`
	ObjectClass           []string      `json:"objectClass,omitempty"`
	Profile               *GroupProfile `json:"profile,omitempty"`
	Type                  string        `json:"type,omitempty"`
}

func (m *GroupResource) UpdateGroup(groupId string, body Group) (*Group, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%v", groupId)
	req, err := m.client.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var group *Group
	resp, err := m.client.requestExecutor.Do(req, &group)
	if err != nil {
		return nil, resp, err
	}
	return group, resp, nil
}
func (m *GroupResource) DeleteGroup(groupId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%v", groupId)
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
func (m *GroupResource) ListGroups(qp *query.Params) ([]*Group, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups")
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var group []*Group
	resp, err := m.client.requestExecutor.Do(req, &group)
	if err != nil {
		return nil, resp, err
	}
	return group, resp, nil
}
func (m *GroupResource) CreateGroup(body Group) (*Group, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups")
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var group *Group
	resp, err := m.client.requestExecutor.Do(req, &group)
	if err != nil {
		return nil, resp, err
	}
	return group, resp, nil
}
func (m *GroupResource) ListRules(qp *query.Params) ([]*GroupRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules")
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var groupRule []*GroupRule
	resp, err := m.client.requestExecutor.Do(req, &groupRule)
	if err != nil {
		return nil, resp, err
	}
	return groupRule, resp, nil
}
func (m *GroupResource) CreateRule(body GroupRule) (*GroupRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules")
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
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
func (m *GroupResource) DeleteRule(ruleId string, qp *query.Params) (*Response, error) {
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
func (m *GroupResource) GetRule(ruleId string, qp *query.Params) (*GroupRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v", ruleId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
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
func (m *GroupResource) UpdateRule(ruleId string, body GroupRule) (*GroupRule, *Response, error) {
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
func (m *GroupResource) ActivateRule(ruleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v/lifecycle/activate", ruleId)
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
func (m *GroupResource) DeactivateRule(ruleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v/lifecycle/deactivate", ruleId)
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
func (m *GroupResource) GetGroup(groupId string, qp *query.Params) (*Group, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%v", groupId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var group *Group
	resp, err := m.client.requestExecutor.Do(req, &group)
	if err != nil {
		return nil, resp, err
	}
	return group, resp, nil
}
func (m *GroupResource) ListGroupUsers(groupId string, qp *query.Params) ([]*User, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%v/users", groupId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var user []*User
	resp, err := m.client.requestExecutor.Do(req, &user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}
func (m *GroupResource) RemoveGroupUser(groupId string, userId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%v/users/%v", groupId, userId)
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
func (m *GroupResource) AddUserToGroup(groupId string, userId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%v/users/%v", groupId, userId)
	req, err := m.client.requestExecutor.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
