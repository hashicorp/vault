// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ TeamTokens = (*teamTokens)(nil)

// TeamTokens describes all the team token related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/team-tokens
type TeamTokens interface {
	// Create a new team token, replacing any existing token.
	Create(ctx context.Context, teamID string) (*TeamToken, error)

	// CreateWithOptions a new team token, with options, replacing any existing token.
	CreateWithOptions(ctx context.Context, teamID string, options TeamTokenCreateOptions) (*TeamToken, error)

	// Read a team token by its ID.
	Read(ctx context.Context, teamID string) (*TeamToken, error)

	// Delete a team token by its ID.
	Delete(ctx context.Context, teamID string) error
}

// teamTokens implements TeamTokens.
type teamTokens struct {
	client *Client
}

// TeamToken represents a Terraform Enterprise team token.
type TeamToken struct {
	ID          string           `jsonapi:"primary,authentication-tokens"`
	CreatedAt   time.Time        `jsonapi:"attr,created-at,iso8601"`
	Description string           `jsonapi:"attr,description"`
	LastUsedAt  time.Time        `jsonapi:"attr,last-used-at,iso8601"`
	Token       string           `jsonapi:"attr,token"`
	ExpiredAt   time.Time        `jsonapi:"attr,expired-at,iso8601"`
	CreatedBy   *CreatedByChoice `jsonapi:"polyrelation,created-by"`
}

// TeamTokenCreateOptions contains the options for creating a team token.
type TeamTokenCreateOptions struct {
	// Optional: The token's expiration date.
	// This feature is available in TFE release v202305-1 and later
	ExpiredAt *time.Time `jsonapi:"attr,expired-at,iso8601,omitempty"`
}

// Create a new team token, replacing any existing token.
func (s *teamTokens) Create(ctx context.Context, teamID string) (*TeamToken, error) {
	return s.CreateWithOptions(ctx, teamID, TeamTokenCreateOptions{})
}

// CreateWithOptions a new team token, with options, replacing any existing token.
func (s *teamTokens) CreateWithOptions(ctx context.Context, teamID string, options TeamTokenCreateOptions) (*TeamToken, error) {
	if !validStringID(&teamID) {
		return nil, ErrInvalidTeamID
	}

	u := fmt.Sprintf("teams/%s/authentication-token", url.PathEscape(teamID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	tt := &TeamToken{}
	err = req.Do(ctx, tt)
	if err != nil {
		return nil, err
	}

	return tt, err
}

// Read a team token by its ID.
func (s *teamTokens) Read(ctx context.Context, teamID string) (*TeamToken, error) {
	if !validStringID(&teamID) {
		return nil, ErrInvalidTeamID
	}

	u := fmt.Sprintf("teams/%s/authentication-token", url.PathEscape(teamID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tt := &TeamToken{}
	err = req.Do(ctx, tt)
	if err != nil {
		return nil, err
	}

	return tt, err
}

// Delete a team token by its ID.
func (s *teamTokens) Delete(ctx context.Context, teamID string) error {
	if !validStringID(&teamID) {
		return ErrInvalidTeamID
	}

	u := fmt.Sprintf("teams/%s/authentication-token", url.PathEscape(teamID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}
