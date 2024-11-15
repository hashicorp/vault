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

const invitationBasePath = orgsBasePath + "/%s/invites"

// InvitationOptions filtering options for invitations.
type InvitationOptions struct {
	Username string `url:"username,omitempty"`
}

// Invitation represents the structure of an Invitation.
type Invitation struct {
	ID              string   `json:"id,omitempty"`
	OrgID           string   `json:"orgId,omitempty"`
	OrgName         string   `json:"orgName,omitempty"`
	GroupID         string   `json:"groupId,omitempty"`
	GroupName       string   `json:"groupName,omitempty"`
	CreatedAt       string   `json:"createdAt,omitempty"`
	ExpiresAt       string   `json:"expiresAt,omitempty"`
	InviterUsername string   `json:"inviterUsername,omitempty"`
	Username        string   `json:"username,omitempty"`
	Roles           []string `json:"roles,omitempty"`
	TeamIDs         []string `json:"teamIds,omitempty"`
}

// Invitations gets all unaccepted invitations to the specified Atlas organization.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-get-invitations/
func (s *OrganizationsServiceOp) Invitations(ctx context.Context, orgID string, opts *InvitationOptions) ([]*Invitation, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(invitationBasePath, orgID)
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*Invitation
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Invitation gets details for one unaccepted invitation to the specified Atlas organization.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-get-one-invitation/
func (s *OrganizationsServiceOp) Invitation(ctx context.Context, orgID, invitationID string) (*Invitation, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	if invitationID == "" {
		return nil, nil, NewArgError("invitationID", "must be set")
	}

	basePath := fmt.Sprintf(invitationBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, invitationID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Invitation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// InviteUser invites one user to the Atlas organization that you specify.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-create-one-invitation/
func (s *OrganizationsServiceOp) InviteUser(ctx context.Context, orgID string, invitation *Invitation) (*Invitation, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(invitationBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, invitation)
	if err != nil {
		return nil, nil, err
	}

	root := new(Invitation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// UpdateInvitation updates one pending invitation to the Atlas organization that you specify.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-update-one-invitation/
func (s *OrganizationsServiceOp) UpdateInvitation(ctx context.Context, orgID string, invitation *Invitation) (*Invitation, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	return s.updateInvitation(ctx, orgID, "", invitation)
}

// UpdateInvitationByID updates one invitation to the Atlas organization.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-update-one-invitation-by-id/
func (s *OrganizationsServiceOp) UpdateInvitationByID(ctx context.Context, orgID, invitationID string, invitation *Invitation) (*Invitation, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	if invitationID == "" {
		return nil, nil, NewArgError("invitationID", "must be set")
	}

	return s.updateInvitation(ctx, orgID, invitationID, invitation)
}

// DeleteInvitation deletes one unaccepted invitation to the specified Atlas organization. You can't delete an invitation that a user has accepted.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-delete-invitation/
func (s *OrganizationsServiceOp) DeleteInvitation(ctx context.Context, orgID, invitationID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}

	if invitationID == "" {
		return nil, NewArgError("invitationID", "must be set")
	}

	basePath := fmt.Sprintf(invitationBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, invitationID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

func (s *OrganizationsServiceOp) updateInvitation(ctx context.Context, orgID, invitationID string, invitation *Invitation) (*Invitation, *Response, error) {
	path := fmt.Sprintf(invitationBasePath, orgID)

	if invitationID != "" {
		path = fmt.Sprintf("%s/%s", path, invitationID)
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, invitation)
	if err != nil {
		return nil, nil, err
	}

	root := new(Invitation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
