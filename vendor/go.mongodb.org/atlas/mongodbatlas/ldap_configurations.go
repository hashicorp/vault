// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ldapConfigurationPath                = "api/atlas/v1.0/groups/%s/userSecurity"
	ldapConfigurationPathuserToDNMapping = ldapConfigurationPath + "/ldap/userToDNMapping"
	ldapVerifyConfigurationPath          = ldapConfigurationPath + "/ldap/verify"
)

// LDAPConfigurationsService is an interface of the LDAP Configuration
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ldaps-configuration/
type LDAPConfigurationsService interface {
	Verify(context.Context, string, *LDAP) (*LDAPConfiguration, *Response, error)
	Get(context.Context, string) (*LDAPConfiguration, *Response, error)
	GetStatus(context.Context, string, string) (*LDAPConfiguration, *Response, error)
	Save(context.Context, string, *LDAPConfiguration) (*LDAPConfiguration, *Response, error)
	Delete(context.Context, string) (*LDAPConfiguration, *Response, error)
}

// LDAPConfigurationsServiceOp handles communication with the LDAP Configuration related methods of the MongoDB Atlas API.
type LDAPConfigurationsServiceOp service

var _ LDAPConfigurationsService = &LDAPConfigurationsServiceOp{}

// LDAPConfiguration represents MongoDB LDAP Configuration.
type LDAPConfiguration struct {
	RequestID   string            `json:"requestId,omitempty"`   // Identifier for the Atlas project associated with the request to verify an LDAP over TLS/SSL configuration.
	GroupID     string            `json:"groupId,omitempty"`     // Unique identifier of the project that owns this alert configuration.
	Request     *LDAPRequest      `json:"request,omitempty"`     // Contains the details of the request to verify an LDAP over TLS/SSL configuration.
	Status      string            `json:"status,omitempty"`      // The current status of the LDAP over TLS/SSL configuration.
	Validations []*LDAPValidation `json:"validations,omitempty"` // Array of validation messages related to the verification of the provided LDAP over TLS/SSL configuration details.
	Links       []*Link           `json:"links,omitempty"`
	LDAP        *LDAP             `json:"ldap,omitempty"` // Specifies the LDAP over TLS/SSL configuration details for an Atlas group.
}

// LDAP specifies an LDAP configuration for a Atlas project.
type LDAP struct {
	AuthenticationEnabled *bool              `json:"authenticationEnabled,omitempty"` // Specifies whether user authentication with LDAP is enabled.
	AuthorizationEnabled  *bool              `json:"authorizationEnabled,omitempty"`  // The current status of the LDAP over TLS/SSL configuration.
	Hostname              *string            `json:"hostname,omitempty"`              // The hostname or IP address of the LDAP server
	Port                  *int               `json:"port,omitempty"`                  // The port to which the LDAP server listens for client connections.
	BindUsername          *string            `json:"bindUsername,omitempty"`          // The user DN that Atlas uses to connect to the LDAP server.
	UserToDNMapping       []*UserToDNMapping `json:"userToDNMapping,omitempty"`       // Maps an LDAP username for authentication to an LDAP Distinguished Name (DN).
	BindPassword          *string            `json:"bindPassword,omitempty"`          // The password used to authenticate the bindUsername.
	CaCertificate         *string            `json:"caCertificate,omitempty"`         // CA certificate used to verify the identity of the LDAP server.
	AuthzQueryTemplate    *string            `json:"authzQueryTemplate,omitempty"`    // An LDAP query template that Atlas executes to obtain the LDAP groups to which the authenticated user belongs.
}

// UserToDNMapping maps an LDAP username for authentication to an LDAP Distinguished Name (DN). Each document contains a match regular expression and either a substitution or ldapQuery template used to transform the LDAP username extracted from the regular expression.
type UserToDNMapping struct {
	Match        string `json:"match,omitempty"`        // A regular expression to match against a provided LDAP username.
	Substitution string `json:"substitution,omitempty"` // An LDAP Distinguished Name (DN) formatting template that converts the LDAP name matched by the match regular expression into an LDAP Distinguished Name.
	LDAPQuery    string `json:"ldapQuery,omitempty"`    // An LDAP query formatting template that inserts the LDAP name matched by the match regular expression into an LDAP query URI as specified by RFC 4515 and RFC 4516.
}

// LDAPValidation contains an array of validation messages related to the verification of the provided LDAP over TLS/SSL configuration details.
type LDAPValidation struct {
	Status         string `json:"status,omitempty"`         // The status of the validation.
	ValidationType string `json:"validationType,omitempty"` // The type of the validation.
}

// LDAPRequest contains the details of the request to verify an LDAP over TLS/SSL configuration.
type LDAPRequest struct {
	Hostname     string `json:"hostname,omitempty"`     // The hostname or IP address of the LDAP server.
	Port         int    `json:"port,omitempty"`         // The port to which the LDAP server listens for client connections from Atlas.
	BindUsername string `json:"bindUsername,omitempty"` // The user DN that Atlas uses to connect to the LDAP server.
}

// Verify requests verification of an LDAP configuration. Use this endpoint to test your LDAP configuration details before saving them.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-request-verification/
func (s *LDAPConfigurationsServiceOp) Verify(ctx context.Context, groupID string, configuration *LDAP) (*LDAPConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if configuration == nil {
		return nil, nil, NewArgError("configuration", "must be set")
	}

	path := fmt.Sprintf(ldapVerifyConfigurationPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, configuration)
	if err != nil {
		return nil, nil, err
	}

	root := new(LDAPConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetStatus retrieves the status of a request for verification of an LDAP configuration.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-verification-status/
func (s *LDAPConfigurationsServiceOp) GetStatus(ctx context.Context, groupID, requestID string) (*LDAPConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if requestID == "" {
		return nil, nil, NewArgError("requestID", "must be set")
	}

	basePath := fmt.Sprintf(ldapVerifyConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, requestID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(LDAPConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Save saves an LDAP configuration for a Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-save/
func (s *LDAPConfigurationsServiceOp) Save(ctx context.Context, groupID string, configuration *LDAPConfiguration) (*LDAPConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if configuration == nil {
		return nil, nil, NewArgError("configuration", "must be set")
	}

	path := fmt.Sprintf(ldapConfigurationPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, configuration)
	if err != nil {
		return nil, nil, err
	}

	root := new(LDAPConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Get retrieves the current LDAP configuration for an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-get-current/
func (s *LDAPConfigurationsServiceOp) Get(ctx context.Context, groupID string) (*LDAPConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(ldapConfigurationPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(LDAPConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete removes the current userToDNMapping from the LDAP configuration for an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ldaps-configuration-remove-usertodnmapping/
func (s *LDAPConfigurationsServiceOp) Delete(ctx context.Context, groupID string) (*LDAPConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(ldapConfigurationPathuserToDNMapping, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(LDAPConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
