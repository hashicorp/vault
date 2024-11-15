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
var _ OrganizationTokens = (*organizationTokens)(nil)

type TokenType string

const (
	// A token which can only access the Audit Trails of an HCP Terraform Organization.
	// See https://developer.hashicorp.com/terraform/cloud-docs/api-docs/audit-trails-tokens
	AuditTrailToken TokenType = "audit-trails"
)

// OrganizationTokens describes all the organization token related methods
// that the Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/organization-tokens
type OrganizationTokens interface {
	// Create a new organization token, replacing any existing token.
	Create(ctx context.Context, organization string) (*OrganizationToken, error)

	// CreateWithOptions a new organization token with options, replacing any existing token.
	CreateWithOptions(ctx context.Context, organization string, options OrganizationTokenCreateOptions) (*OrganizationToken, error)

	// Read an organization token.
	Read(ctx context.Context, organization string) (*OrganizationToken, error)

	// Read an organization token with options.
	ReadWithOptions(ctx context.Context, organization string, options OrganizationTokenReadOptions) (*OrganizationToken, error)

	// Delete an organization token.
	Delete(ctx context.Context, organization string) error

	// Delete an organization token with options.
	DeleteWithOptions(ctx context.Context, organization string, options OrganizationTokenDeleteOptions) error
}

// organizationTokens implements OrganizationTokens.
type organizationTokens struct {
	client *Client
}

// OrganizationToken represents a Terraform Enterprise organization token.
type OrganizationToken struct {
	ID          string           `jsonapi:"primary,authentication-tokens"`
	CreatedAt   time.Time        `jsonapi:"attr,created-at,iso8601"`
	Description string           `jsonapi:"attr,description"`
	LastUsedAt  time.Time        `jsonapi:"attr,last-used-at,iso8601"`
	Token       string           `jsonapi:"attr,token"`
	ExpiredAt   time.Time        `jsonapi:"attr,expired-at,iso8601"`
	CreatedBy   *CreatedByChoice `jsonapi:"polyrelation,created-by"`
}

// OrganizationTokenCreateOptions contains the options for creating an organization token.
type OrganizationTokenCreateOptions struct {
	// Optional: The token's expiration date.
	// This feature is available in TFE release v202305-1 and later
	ExpiredAt *time.Time `jsonapi:"attr,expired-at,iso8601,omitempty" url:"-"`
	// Optional: What type of token to create
	// This option is only applicable to HCP Terraform and is ignored by TFE.
	TokenType *TokenType `url:"token,omitempty"`
}

// OrganizationTokenReadOptions contains the options for reading an organization token.
type OrganizationTokenReadOptions struct {
	// Optional: What type of token to read
	// This option is only applicable to HCP Terraform and is ignored by TFE.
	TokenType *TokenType `url:"token,omitempty"`
}

// OrganizationTokenDeleteOptions contains the options for deleting an organization token.
type OrganizationTokenDeleteOptions struct {
	// Optional: What type of token to delete
	// This option is only applicable to HCP Terraform and is ignored by TFE.
	TokenType *TokenType `url:"token,omitempty"`
}

// Create a new organization token, replacing any existing token.
func (s *organizationTokens) Create(ctx context.Context, organization string) (*OrganizationToken, error) {
	return s.CreateWithOptions(ctx, organization, OrganizationTokenCreateOptions{})
}

// CreateWithOptions a new organization token with options, replacing any existing token.
func (s *organizationTokens) CreateWithOptions(ctx context.Context, organization string, options OrganizationTokenCreateOptions) (*OrganizationToken, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/authentication-token", url.PathEscape(organization))
	qp, err := decodeQueryParams(options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequestWithAdditionalQueryParams("POST", u, &options, qp)
	if err != nil {
		return nil, err
	}

	ot := &OrganizationToken{}
	err = req.Do(ctx, ot)
	if err != nil {
		return nil, err
	}

	return ot, err
}

// Read an organization token.
func (s *organizationTokens) Read(ctx context.Context, organization string) (*OrganizationToken, error) {
	return s.ReadWithOptions(ctx, organization, OrganizationTokenReadOptions{})
}

// Read an organization token with options.
func (s *organizationTokens) ReadWithOptions(ctx context.Context, organization string, options OrganizationTokenReadOptions) (*OrganizationToken, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/authentication-token", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	ot := &OrganizationToken{}
	err = req.Do(ctx, ot)
	if err != nil {
		return nil, err
	}

	return ot, err
}

// Delete an organization token.
func (s *organizationTokens) Delete(ctx context.Context, organization string) error {
	return s.DeleteWithOptions(ctx, organization, OrganizationTokenDeleteOptions{})
}

// Delete an organization token with options
func (s *organizationTokens) DeleteWithOptions(ctx context.Context, organization string, options OrganizationTokenDeleteOptions) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/authentication-token", url.PathEscape(organization))
	qp, err := decodeQueryParams(options)
	if err != nil {
		return err
	}

	req, err := s.client.NewRequestWithAdditionalQueryParams("DELETE", u, nil, qp)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}
