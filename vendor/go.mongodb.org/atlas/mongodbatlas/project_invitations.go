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

const projectInvitationBasePath = projectBasePath + "/%s/invites"

// Invitations gets all unaccepted invitations to the specified Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/project-get-invitations/
func (s *ProjectsServiceOp) Invitations(ctx context.Context, groupID string, opts *InvitationOptions) ([]*Invitation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(projectInvitationBasePath, groupID)
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

// Invitation gets details for one unaccepted invitation to the specified Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/project-get-one-invitation/
func (s *ProjectsServiceOp) Invitation(ctx context.Context, groupID, invitationID string) (*Invitation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if invitationID == "" {
		return nil, nil, NewArgError("invitationID", "must be set")
	}

	basePath := fmt.Sprintf(projectInvitationBasePath, groupID)
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

// InviteUser invites one user to the Atlas project that you specify.
func (s *ProjectsServiceOp) InviteUser(ctx context.Context, groupID string, invitation *Invitation) (*Invitation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(projectInvitationBasePath, groupID)

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

// UpdateInvitation updates one pending invitation to the Atlas project that you specify.
//
// See more: https://docs.atlas.mongodb.com/reference/api/project-update-one-invitation/
func (s *ProjectsServiceOp) UpdateInvitation(ctx context.Context, groupID string, invitation *Invitation) (*Invitation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	return s.updateInvitation(ctx, groupID, "", invitation)
}

// UpdateInvitationByID updates one invitation to the Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/project-update-one-invitation-by-id/
func (s *ProjectsServiceOp) UpdateInvitationByID(ctx context.Context, groupID, invitationID string, invitation *Invitation) (*Invitation, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if invitationID == "" {
		return nil, nil, NewArgError("invitationID", "must be set")
	}

	return s.updateInvitation(ctx, groupID, invitationID, invitation)
}

// DeleteInvitation deletes one unaccepted invitation to the specified Atlas project. You can't delete an invitation that a user has accepted.
//
// See more: https://docs.atlas.mongodb.com/reference/api/project-delete-invitation/
func (s *ProjectsServiceOp) DeleteInvitation(ctx context.Context, groupID, invitationID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	if invitationID == "" {
		return nil, NewArgError("invitationID", "must be set")
	}

	basePath := fmt.Sprintf(projectInvitationBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, invitationID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

func (s *ProjectsServiceOp) updateInvitation(ctx context.Context, groupID, invitationID string, invitation *Invitation) (*Invitation, *Response, error) {
	path := fmt.Sprintf(projectInvitationBasePath, groupID)

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
