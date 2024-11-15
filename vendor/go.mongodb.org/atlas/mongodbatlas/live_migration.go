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

const linkTokensBasePath = "api/atlas/v1.0/orgs/%s/liveMigrations/linkTokens" //nolint:gosec //false positive
const liveMigrationBasePath = "api/atlas/v1.0/groups/%s/liveMigrations"

// LiveMigrationService is an interface for interfacing with the Live Migration
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/
type LiveMigrationService interface {
	CreateLinkToken(context.Context, string, *TokenCreateRequest) (*LinkToken, *Response, error)
	DeleteLinkToken(context.Context, string) (*Response, error)
	CreateValidation(context.Context, string, *LiveMigration) (*Validation, *Response, error)
	GetValidationStatus(context.Context, string, string) (*Validation, *Response, error)
	Create(context.Context, string, *LiveMigration) (*LiveMigration, *Response, error)
	Get(context.Context, string, string) (*LiveMigration, *Response, error)
	Start(context.Context, string, string) (*Validation, *Response, error)
}

// LiveMigrationServiceOp provides an implementation of AlertsService.
type LiveMigrationServiceOp service

var _ LiveMigrationService = &LiveMigrationServiceOp{}

// TokenCreateRequest represents the Request Body Parameters of LiveMigrationService.CreateLinkToken.
type TokenCreateRequest struct {
	AccessListIPs []string `json:"accessListIps"` // One IP address access list entry associated with the API key.
}

type LinkToken struct {
	LinkToken string `json:"linkToken,omitempty"` // Atlas-generated token that links the source (Cloud Manager or Ops Manager) and destination (Atlas) clusters for migration.
}

// LiveMigration represents a live migration.
type LiveMigration struct {
	Source          *Source      `json:"source,omitempty"`          // Source describes the Cloud Manager or Ops Manager source of the migrating cluster.
	Destination     *Destination `json:"destination,omitempty"`     // Destination describes the Atlas destination of the migrating Cloud Manager or Ops Manager cluster.
	MigrationHosts  []string     `json:"migrationHosts,omitempty"`  // MigrationHosts is a list of hosts running the MongoDB Agent that can transfer your MongoDB data from the source (Cloud Manager or Ops Manager) to destination (Atlas) deployments.
	DropEnabled     *bool        `json:"dropEnabled,omitempty"`     // DropEnabled indicates whether this process should drop existing collections from the destination (Atlas) cluster given in destination.clusterName before starting the migration of data from the source cluster.
	ID              string       `json:"_id,omitempty"`             // ID Unique 24-hexadecimal digit string that identifies the migration.
	Status          string       `json:"status,omitempty"`          // Status of the migration when you submitted this request.
	ReadyForCutover *bool        `json:"readyForCutover,omitempty"` // ReadyForCutover indicates whether the live migration process is ready to start the cutover process.
}

// Source represents the Cloud Manager or Ops Manager source of the migrating cluster.
type Source struct {
	ClusterName           string `json:"clusterName,omitempty"`           // Human-readable label that identifies the source Cloud Manager or Ops Manager cluster.
	GroupID               string `json:"groupId,omitempty"`               // Unique 24-hexadecimal digit string that identifies the source project.
	Username              string `json:"username,omitempty"`              // Human-readable label that identifies the SCRAM-SHA user that connects to the source Cloud Manager or Ops Manager cluster.
	Password              string `json:"password,omitempty"`              // Password that authenticates the username to the source Cloud Manager or Ops Manager cluster.
	SSL                   *bool  `json:"ssl,omitempty"`                   // Flag that indicates whether you have TLS enabled.
	CACertificatePath     string `json:"caCertificatePath,omitempty"`     // Path to the CA certificate that signed TLS certificates use to authenticate to the source Cloud Manager or Ops Manager cluster.
	ManagedAuthentication *bool  `json:"managedAuthentication,omitempty"` // Flag that indicates whether MongoDB Automation manages authentication to the source Cloud Manager or Ops Manager cluster.
}

// Destination represents settings of the Atlas destination.
type Destination struct {
	ClusterName string `json:"clusterName,omitempty"` // Human-readable label that identifies the Atlas destination cluster.
	GroupID     string `json:"groupId,omitempty"`     // Unique 24-hexadecimal digit string that identifies the Atlas destination project.
}

type Validation struct {
	ID            string `json:"_id,omitempty"`           // Unique 24-hexadecimal digit string that identifies this process validating the live migration.
	GroupID       string `json:"groupId,omitempty"`       // Unique 24-hexadecimal digit string that identifies the Atlas project to validate.
	Status        string `json:"status,omitempty"`        // State of the validation job when you submitted this request.
	SourceGroupID string `json:"sourceGroupId,omitempty"` // Unique 24-hexadecimal digit string that identifies the source (Cloud Manager or Ops Manager) project.
	ErrorMessage  string `json:"errorMessage,omitempty"`  // Reason why the validation job failed.
}

// CreateLinkToken create one new link-token.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/create-one-link-token/
func (s *LiveMigrationServiceOp) CreateLinkToken(ctx context.Context, orgID string, body *TokenCreateRequest) (*LinkToken, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(linkTokensBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	root := new(LinkToken)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeleteLinkToken deletes one link-token.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/delete-one-link-token/
func (s *LiveMigrationServiceOp) DeleteLinkToken(ctx context.Context, orgID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(linkTokensBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// CreateValidation creates one new validation request.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/validate-one-migration-request/
func (s *LiveMigrationServiceOp) CreateValidation(ctx context.Context, groupID string, body *LiveMigration) (*Validation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(liveMigrationBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "validate")

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	root := new(Validation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetValidationStatus returns one validation job.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/return-one-specific-validation-job/
func (s *LiveMigrationServiceOp) GetValidationStatus(ctx context.Context, groupID, id string) (*Validation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if id == "" {
		return nil, nil, NewArgError("id", "must be set")
	}

	basePath := fmt.Sprintf(liveMigrationBasePath, groupID)
	path := fmt.Sprintf("%s/validate/%s", basePath, id)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Validation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates one new live migration.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/create-one-migration/
func (s *LiveMigrationServiceOp) Create(ctx context.Context, groupID string, body *LiveMigration) (*LiveMigration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(liveMigrationBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	root := new(LiveMigration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Get returns one migration job.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/return-one-specific-migration/
func (s *LiveMigrationServiceOp) Get(ctx context.Context, groupID, id string) (*LiveMigration, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if id == "" {
		return nil, nil, NewArgError("id", "must be set")
	}

	basePath := fmt.Sprintf(liveMigrationBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, id)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(LiveMigration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Start starts the migration of one deployment.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/start-the-migration-cutover/
func (s *LiveMigrationServiceOp) Start(ctx context.Context, groupID, id string) (*Validation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if id == "" {
		return nil, nil, NewArgError("id", "must be set")
	}

	basePath := fmt.Sprintf(liveMigrationBasePath, groupID)
	path := fmt.Sprintf("%s/%s/cutover", basePath, id)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Validation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
