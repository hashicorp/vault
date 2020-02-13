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

type UserResource resource

type User struct {
	Embedded              interface{}      `json:"_embedded,omitempty"`
	Links                 interface{}      `json:"_links,omitempty"`
	Activated             *time.Time       `json:"activated,omitempty"`
	Created               *time.Time       `json:"created,omitempty"`
	Credentials           *UserCredentials `json:"credentials,omitempty"`
	Id                    string           `json:"id,omitempty"`
	LastLogin             *time.Time       `json:"lastLogin,omitempty"`
	LastUpdated           *time.Time       `json:"lastUpdated,omitempty"`
	PasswordChanged       *time.Time       `json:"passwordChanged,omitempty"`
	Profile               *UserProfile     `json:"profile,omitempty"`
	Status                string           `json:"status,omitempty"`
	StatusChanged         *time.Time       `json:"statusChanged,omitempty"`
	TransitioningToStatus string           `json:"transitioningToStatus,omitempty"`
}

func (m *UserResource) CreateUser(body User, qp *query.Params) (*User, *Response, error) {
	url := fmt.Sprintf("/api/v1/users")
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var user *User
	resp, err := m.client.requestExecutor.Do(req, &user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}
func (m *UserResource) GetUser(userId string) (*User, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v", userId)
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var user *User
	resp, err := m.client.requestExecutor.Do(req, &user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}
func (m *UserResource) UpdateUser(userId string, body User, qp *query.Params) (*User, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var user *User
	resp, err := m.client.requestExecutor.Do(req, &user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}
func (m *UserResource) DeactivateOrDeleteUser(userId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v", userId)
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
func (m *UserResource) ListUsers(qp *query.Params) ([]*User, *Response, error) {
	url := fmt.Sprintf("/api/v1/users")
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
func (m *UserResource) ListAppLinks(userId string, qp *query.Params) ([]*AppLink, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/appLinks", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var appLink []*AppLink
	resp, err := m.client.requestExecutor.Do(req, &appLink)
	if err != nil {
		return nil, resp, err
	}
	return appLink, resp, nil
}
func (m *UserResource) ChangePassword(userId string, body ChangePasswordRequest, qp *query.Params) (*UserCredentials, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/credentials/change_password", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var userCredentials *UserCredentials
	resp, err := m.client.requestExecutor.Do(req, &userCredentials)
	if err != nil {
		return nil, resp, err
	}
	return userCredentials, resp, nil
}
func (m *UserResource) ChangeRecoveryQuestion(userId string, body UserCredentials) (*UserCredentials, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/credentials/change_recovery_question", userId)
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var userCredentials *UserCredentials
	resp, err := m.client.requestExecutor.Do(req, &userCredentials)
	if err != nil {
		return nil, resp, err
	}
	return userCredentials, resp, nil
}
func (m *UserResource) ForgotPassword(userId string, body UserCredentials, qp *query.Params) (*ForgotPasswordResponse, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/credentials/forgot_password", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var forgotPasswordResponse *ForgotPasswordResponse
	resp, err := m.client.requestExecutor.Do(req, &forgotPasswordResponse)
	if err != nil {
		return nil, resp, err
	}
	return forgotPasswordResponse, resp, nil
}
func (m *UserResource) ListUserGroups(userId string, qp *query.Params) ([]*Group, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/groups", userId)
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
func (m *UserResource) ActivateUser(userId string, qp *query.Params) (*UserActivationToken, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/activate", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var userActivationToken *UserActivationToken
	resp, err := m.client.requestExecutor.Do(req, &userActivationToken)
	if err != nil {
		return nil, resp, err
	}
	return userActivationToken, resp, nil
}
func (m *UserResource) DeactivateUser(userId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/deactivate", userId)
	if qp != nil {
		url = url + qp.String()
	}
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
func (m *UserResource) ExpirePassword(userId string, qp *query.Params) (*TempPassword, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/expire_password", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var tempPassword *TempPassword
	resp, err := m.client.requestExecutor.Do(req, &tempPassword)
	if err != nil {
		return nil, resp, err
	}
	return tempPassword, resp, nil
}
func (m *UserResource) ResetAllFactors(userId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/reset_factors", userId)
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
func (m *UserResource) ResetPassword(userId string, qp *query.Params) (*ResetPasswordToken, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/reset_password", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var resetPasswordToken *ResetPasswordToken
	resp, err := m.client.requestExecutor.Do(req, &resetPasswordToken)
	if err != nil {
		return nil, resp, err
	}
	return resetPasswordToken, resp, nil
}
func (m *UserResource) SuspendUser(userId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/suspend", userId)
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
func (m *UserResource) UnlockUser(userId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/unlock", userId)
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
func (m *UserResource) UnsuspendUser(userId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/lifecycle/unsuspend", userId)
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
func (m *UserResource) ListAssignedRoles(userId string, qp *query.Params) ([]*Role, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/roles", userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var role []*Role
	resp, err := m.client.requestExecutor.Do(req, &role)
	if err != nil {
		return nil, resp, err
	}
	return role, resp, nil
}
func (m *UserResource) AddRoleToUser(userId string, body Role) (*Role, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/roles", userId)
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var role *Role
	resp, err := m.client.requestExecutor.Do(req, &role)
	if err != nil {
		return nil, resp, err
	}
	return role, resp, nil
}
func (m *UserResource) RemoveRoleFromUser(userId string, roleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/roles/%v", userId, roleId)
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
func (m *UserResource) ListGroupTargetsForRole(userId string, roleId string, qp *query.Params) ([]*Group, *Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/roles/%v/targets/groups", userId, roleId)
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
func (m *UserResource) RemoveGroupTargetFromRole(userId string, roleId string, groupId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/roles/%v/targets/groups/%v", userId, roleId, groupId)
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
func (m *UserResource) AddGroupTargetToRole(userId string, roleId string, groupId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/roles/%v/targets/groups/%v", userId, roleId, groupId)
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
func (m *UserResource) EndAllUserSessions(userId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/sessions", userId)
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
