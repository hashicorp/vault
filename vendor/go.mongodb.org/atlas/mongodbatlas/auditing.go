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
	auditingsPath = "api/atlas/v1.0/groups/%s/auditLog"
)

// AuditingsService is an interface for interfacing with the Auditing
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/auditing/
type AuditingsService interface {
	Get(context.Context, string) (*Auditing, *Response, error)
	Configure(context.Context, string, *Auditing) (*Auditing, *Response, error)
}

// AuditingsServiceOp handles communication with the Auditings related methods
// of the MongoDB Atlas API.
type AuditingsServiceOp service

var _ AuditingsService = &AuditingsServiceOp{}

// Auditing represents MongoDB Maintenance Windows.
type Auditing struct {
	AuditAuthorizationSuccess *bool  `json:"auditAuthorizationSuccess,omitempty"` // Indicates whether the auditing system captures successful authentication attempts for audit filters using the "atype" : "authCheck" auditing event. For more information, see auditAuthorizationSuccess
	AuditFilter               string `json:"auditFilter,omitempty"`               // JSON-formatted audit filter used by the project
	ConfigurationType         string `json:"configurationType,omitempty"`         // Denotes the configuration method for the audit filter. Possible values are: NONE - auditing not configured for the project.m FILTER_BUILDER - auditing configured via Atlas UI filter builderm FILTER_JSON - auditing configured via Atlas custom filter or API
	Enabled                   *bool  `json:"enabled,omitempty"`                   // Denotes whether or not the project associated with the {GROUP-ID} has database auditing enabled.
}

// Get audit configuration for the project associated with {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/auditing-get-auditLog/
func (s *AuditingsServiceOp) Get(ctx context.Context, groupID string) (*Auditing, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(auditingsPath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Auditing)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Configure the audit configuration for the project associated with {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/auditing-set-auditLog/
func (s *AuditingsServiceOp) Configure(ctx context.Context, groupID string, configRequest *Auditing) (*Auditing, *Response, error) {
	if configRequest == nil {
		return nil, nil, NewArgError("configRequest", "cannot be nil")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "cannot be nil")
	}

	path := fmt.Sprintf(auditingsPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, configRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Auditing)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
