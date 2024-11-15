// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ Users = (*users)(nil)

// Users describes all the user related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/account
type Users interface {
	// ReadCurrent reads the details of the currently authenticated user.
	ReadCurrent(ctx context.Context) (*User, error)

	// UpdateCurrent updates attributes of the currently authenticated user.
	UpdateCurrent(ctx context.Context, options UserUpdateOptions) (*User, error)
}

// users implements Users.
type users struct {
	client *Client
}

// User represents a Terraform Enterprise user.
type User struct {
	ID               string     `jsonapi:"primary,users"`
	AvatarURL        string     `jsonapi:"attr,avatar-url"`
	Email            string     `jsonapi:"attr,email"`
	IsServiceAccount bool       `jsonapi:"attr,is-service-account"`
	TwoFactor        *TwoFactor `jsonapi:"attr,two-factor"`
	UnconfirmedEmail string     `jsonapi:"attr,unconfirmed-email"`
	Username         string     `jsonapi:"attr,username"`
	V2Only           bool       `jsonapi:"attr,v2-only"`
	// Deprecated: IsSiteAdmin was deprecated in v202406 and will be removed in a future version of Terraform Enterprise
	IsSiteAdmin *bool            `jsonapi:"attr,is-site-admin"`
	IsAdmin     *bool            `jsonapi:"attr,is-admin"`
	IsSsoLogin  *bool            `jsonapi:"attr,is-sso-login"`
	Permissions *UserPermissions `jsonapi:"attr,permissions"`

	// Relations
	// AuthenticationTokens *AuthenticationTokens `jsonapi:"relation,authentication-tokens"`
}

// UserPermissions represents the user permissions.
type UserPermissions struct {
	CanCreateOrganizations bool `jsonapi:"attr,can-create-organizations"`
	CanChangeEmail         bool `jsonapi:"attr,can-change-email"`
	CanChangeUsername      bool `jsonapi:"attr,can-change-username"`
	CanManageUserTokens    bool `jsonapi:"attr,can-manage-user-tokens"`
	CanView2FaSettings     bool `jsonapi:"attr,can-view2fa-settings"`
	CanManageHcpAccount    bool `jsonapi:"attr,can-manage-hcp-account"`
}

// TwoFactor represents the organization permissions.
type TwoFactor struct {
	Enabled  bool `jsonapi:"attr,enabled"`
	Verified bool `jsonapi:"attr,verified"`
}

// UserUpdateOptions represents the options for updating a user.
type UserUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,users"`

	// Optional: New username.
	Username *string `jsonapi:"attr,username,omitempty"`

	// Optional: New email address (must be consumed afterwards to take effect).
	Email *string `jsonapi:"attr,email,omitempty"`
}

// ReadCurrent reads the details of the currently authenticated user.
func (s *users) ReadCurrent(ctx context.Context) (*User, error) {
	req, err := s.client.NewRequest("GET", "account/details", nil)
	if err != nil {
		return nil, err
	}

	u := &User{}
	err = req.Do(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdateCurrent updates attributes of the currently authenticated user.
func (s *users) UpdateCurrent(ctx context.Context, options UserUpdateOptions) (*User, error) {
	req, err := s.client.NewRequest("PATCH", "account/update", &options)
	if err != nil {
		return nil, err
	}

	u := &User{}
	err = req.Do(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
