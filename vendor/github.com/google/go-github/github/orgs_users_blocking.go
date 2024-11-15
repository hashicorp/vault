// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ListBlockedUsers lists all the users blocked by an organization.
//
// GitHub API docs: https://developer.github.com/v3/orgs/blocking/#list-blocked-users
func (s *OrganizationsService) ListBlockedUsers(ctx context.Context, org string, opt *ListOptions) ([]*User, *Response, error) {
	u := fmt.Sprintf("orgs/%v/blocks", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeBlockUsersPreview)

	var blockedUsers []*User
	resp, err := s.client.Do(ctx, req, &blockedUsers)
	if err != nil {
		return nil, resp, err
	}

	return blockedUsers, resp, nil
}

// IsBlocked reports whether specified user is blocked from an organization.
//
// GitHub API docs: https://developer.github.com/v3/orgs/blocking/#check-whether-a-user-is-blocked-from-an-organization
func (s *OrganizationsService) IsBlocked(ctx context.Context, org string, user string) (bool, *Response, error) {
	u := fmt.Sprintf("orgs/%v/blocks/%v", org, user)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeBlockUsersPreview)

	resp, err := s.client.Do(ctx, req, nil)
	isBlocked, err := parseBoolResponse(err)
	return isBlocked, resp, err
}

// BlockUser blocks specified user from an organization.
//
// GitHub API docs: https://developer.github.com/v3/orgs/blocking/#block-a-user
func (s *OrganizationsService) BlockUser(ctx context.Context, org string, user string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/blocks/%v", org, user)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeBlockUsersPreview)

	return s.client.Do(ctx, req, nil)
}

// UnblockUser unblocks specified user from an organization.
//
// GitHub API docs: https://developer.github.com/v3/orgs/blocking/#unblock-a-user
func (s *OrganizationsService) UnblockUser(ctx context.Context, org string, user string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/blocks/%v", org, user)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeBlockUsersPreview)

	return s.client.Do(ctx, req, nil)
}
