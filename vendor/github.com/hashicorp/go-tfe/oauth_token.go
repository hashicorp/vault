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
var _ OAuthTokens = (*oAuthTokens)(nil)

// OAuthTokens describes all the OAuth token related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/oauth-tokens
type OAuthTokens interface {
	// List all the OAuth tokens for a given organization.
	List(ctx context.Context, organization string, options *OAuthTokenListOptions) (*OAuthTokenList, error)
	// Read a OAuth token by its ID.
	Read(ctx context.Context, oAuthTokenID string) (*OAuthToken, error)

	// Update an existing OAuth token.
	Update(ctx context.Context, oAuthTokenID string, options OAuthTokenUpdateOptions) (*OAuthToken, error)

	// Delete a OAuth token by its ID.
	Delete(ctx context.Context, oAuthTokenID string) error
}

// oAuthTokens implements OAuthTokens.
type oAuthTokens struct {
	client *Client
}

// OAuthTokenList represents a list of OAuth tokens.
type OAuthTokenList struct {
	*Pagination
	Items []*OAuthToken
}

// OAuthToken represents a VCS configuration including the associated
// OAuth token
type OAuthToken struct {
	ID                  string    `jsonapi:"primary,oauth-tokens"`
	UID                 string    `jsonapi:"attr,uid"`
	CreatedAt           time.Time `jsonapi:"attr,created-at,iso8601"`
	HasSSHKey           bool      `jsonapi:"attr,has-ssh-key"`
	ServiceProviderUser string    `jsonapi:"attr,service-provider-user"`

	// Relations
	OAuthClient *OAuthClient `jsonapi:"relation,oauth-client"`
}

// OAuthTokenListOptions represents the options for listing
// OAuth tokens.
type OAuthTokenListOptions struct {
	ListOptions
}

// OAuthTokenUpdateOptions represents the options for updating an OAuth token.
type OAuthTokenUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,oauth-tokens"`

	// Optional: A private SSH key to be used for git clone operations.
	PrivateSSHKey *string `jsonapi:"attr,ssh-key,omitempty"`
}

// List all the OAuth tokens for a given organization.
func (s *oAuthTokens) List(ctx context.Context, organization string, options *OAuthTokenListOptions) (*OAuthTokenList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/oauth-tokens", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	otl := &OAuthTokenList{}
	err = req.Do(ctx, otl)
	if err != nil {
		return nil, err
	}

	return otl, nil
}

// Read an OAuth token by its ID.
func (s *oAuthTokens) Read(ctx context.Context, oAuthTokenID string) (*OAuthToken, error) {
	if !validStringID(&oAuthTokenID) {
		return nil, ErrInvalidOauthTokenID
	}

	u := fmt.Sprintf("oauth-tokens/%s", url.PathEscape(oAuthTokenID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	ot := &OAuthToken{}
	err = req.Do(ctx, ot)
	if err != nil {
		return nil, err
	}

	return ot, err
}

// Update an existing OAuth token.
func (s *oAuthTokens) Update(ctx context.Context, oAuthTokenID string, options OAuthTokenUpdateOptions) (*OAuthToken, error) {
	if !validStringID(&oAuthTokenID) {
		return nil, ErrInvalidOauthTokenID
	}

	u := fmt.Sprintf("oauth-tokens/%s", url.PathEscape(oAuthTokenID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	ot := &OAuthToken{}
	err = req.Do(ctx, ot)
	if err != nil {
		return nil, err
	}

	return ot, err
}

// Delete an OAuth token by its ID.
func (s *oAuthTokens) Delete(ctx context.Context, oAuthTokenID string) error {
	if !validStringID(&oAuthTokenID) {
		return ErrInvalidOauthTokenID
	}

	u := fmt.Sprintf("oauth-tokens/%s", url.PathEscape(oAuthTokenID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}
