// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AdminUsers = (*adminUsers)(nil)

// AdminUsers describes all the admin user related methods that the Terraform
// Enterprise  API supports.
// It contains endpoints to help site administrators manage their users.
//
// TFE API docs: https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/users
type AdminUsers interface {
	// List all the users of the given installation.
	List(ctx context.Context, options *AdminUserListOptions) (*AdminUserList, error)

	// Delete a user by its ID.
	Delete(ctx context.Context, userID string) error

	// Suspend a user by its ID.
	Suspend(ctx context.Context, userID string) (*AdminUser, error)

	// Unsuspend a user by its ID.
	Unsuspend(ctx context.Context, userID string) (*AdminUser, error)

	// GrantAdmin grants admin privileges to a user by its ID.
	GrantAdmin(ctx context.Context, userID string) (*AdminUser, error)

	// RevokeAdmin revokees admin privileges to a user by its ID.
	RevokeAdmin(ctx context.Context, userID string) (*AdminUser, error)

	// Disable2FA disables a user's two-factor authentication in the situation
	// where they have lost access to their device and recovery codes.
	Disable2FA(ctx context.Context, userID string) (*AdminUser, error)
}

// adminUsers implements the AdminUsers interface.
type adminUsers struct {
	client *Client
}

// AdminUser represents a user as seen by an Admin.
type AdminUser struct {
	ID               string     `jsonapi:"primary,users"`
	Email            string     `jsonapi:"attr,email"`
	Username         string     `jsonapi:"attr,username"`
	AvatarURL        string     `jsonapi:"attr,avatar-url"`
	TwoFactor        *TwoFactor `jsonapi:"attr,two-factor"`
	IsAdmin          bool       `jsonapi:"attr,is-admin"`
	IsSuspended      bool       `jsonapi:"attr,is-suspended"`
	IsServiceAccount bool       `jsonapi:"attr,is-service-account"`

	// Relations
	Organizations []*Organization `jsonapi:"relation,organizations"`
}

// AdminUserList represents a list of users.
type AdminUserList struct {
	*Pagination
	Items []*AdminUser
}

// AdminUserIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/users#available-related-resources
type AdminUserIncludeOpt string

const AdminUserOrgs AdminUserIncludeOpt = "organizations"

// AdminUserListOptions represents the options for listing users.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/users#query-parameters
type AdminUserListOptions struct {
	ListOptions

	// Optional: A search query string. Users are searchable by username and email address.
	Query string `url:"q,omitempty"`

	// Optional: Can be "true" or "false" to show only administrators or non-administrators.
	Administrators string `url:"filter[admin],omitempty"`

	// Optional: Can be "true" or "false" to show only suspended users or users who are not suspended.
	SuspendedUsers string `url:"filter[suspended],omitempty"`

	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/users#available-related-resources
	Include []AdminUserIncludeOpt `url:"include,omitempty"`
}

// List all user accounts in the Terraform Enterprise installation
func (a *adminUsers) List(ctx context.Context, options *AdminUserListOptions) (*AdminUserList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := "admin/users"
	req, err := a.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	aul := &AdminUserList{}
	err = req.Do(ctx, aul)
	if err != nil {
		return nil, err
	}

	return aul, nil
}

// Delete a user by its ID.
func (a *adminUsers) Delete(ctx context.Context, userID string) error {
	if !validStringID(&userID) {
		return ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s", url.PathEscape(userID))
	req, err := a.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Suspend a user by its ID.
func (a *adminUsers) Suspend(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/suspend", url.PathEscape(userID))
	req, err := a.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = req.Do(ctx, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}

// Unsuspend a user by its ID.
func (a *adminUsers) Unsuspend(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/unsuspend", url.PathEscape(userID))
	req, err := a.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = req.Do(ctx, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}

// GrantAdmin grants admin privileges to a user by its ID.
func (a *adminUsers) GrantAdmin(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/grant_admin", url.PathEscape(userID))
	req, err := a.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = req.Do(ctx, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}

// RevokeAdmin revokes admin privileges to a user by its ID.
func (a *adminUsers) RevokeAdmin(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/revoke_admin", url.PathEscape(userID))
	req, err := a.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = req.Do(ctx, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}

// Disable2FA disables a user's two-factor authentication in the situation
// where they have lost access to their device and recovery codes.
func (a *adminUsers) Disable2FA(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/disable_two_factor", url.PathEscape(userID))
	req, err := a.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = req.Do(ctx, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}

func (o *AdminUserListOptions) valid() error {
	return nil
}
