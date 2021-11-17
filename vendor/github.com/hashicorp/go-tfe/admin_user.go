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
// TFE API docs: https://www.terraform.io/docs/cloud/api/admin/users.html
type AdminUsers interface {
	// List all the users of the given installation.
	List(ctx context.Context, options AdminUserListOptions) (*AdminUserList, error)

	// Delete a user by its ID.
	Delete(ctx context.Context, userID string) error

	// Suspend a user by its ID.
	Suspend(ctx context.Context, userID string) (*AdminUser, error)

	// Unsuspend a user by its ID.
	Unsuspend(ctx context.Context, userID string) (*AdminUser, error)

	// GrantAdmin grants admin privilages to a user by its ID.
	GrantAdmin(ctx context.Context, userID string) (*AdminUser, error)

	// RevokeAdmin revokees admin privilages to a user by its ID.
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

// AdminUserListOptions represents the options for listing users.
// https://www.terraform.io/docs/cloud/api/admin/users.html#query-parameters
type AdminUserListOptions struct {
	ListOptions

	// A search query string. Users are searchable by username and email address.
	Query *string `url:"q,omitempty"`

	// Can be "true" or "false" to show only administrators or non-administrators.
	Administrators *string `url:"filter[admin]"`

	// Can be "true" or "false" to show only suspended users or users who are not suspended.
	SuspendedUsers *string `url:"filter[suspended]"`

	// A list of relations to include. See available resources
	// https://www.terraform.io/docs/cloud/api/admin/users.html#available-related-resources
	Include *string `url:"include"`
}

// List all user accounts in the Terraform Enterprise installation
func (a *adminUsers) List(ctx context.Context, options AdminUserListOptions) (*AdminUserList, error) {
	u := "admin/users"
	req, err := a.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	aul := &AdminUserList{}
	err = a.client.do(ctx, req, aul)
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

	u := fmt.Sprintf("admin/users/%s", url.QueryEscape(userID))
	req, err := a.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return a.client.do(ctx, req, nil)
}

// Suspend a user by its ID.
func (a *adminUsers) Suspend(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/suspend", url.QueryEscape(userID))
	req, err := a.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = a.client.do(ctx, req, au)
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

	u := fmt.Sprintf("admin/users/%s/actions/unsuspend", url.QueryEscape(userID))
	req, err := a.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = a.client.do(ctx, req, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}

// GrantAdmin grants admin privilages to a user by its ID.
func (a *adminUsers) GrantAdmin(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/grant_admin", url.QueryEscape(userID))
	req, err := a.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = a.client.do(ctx, req, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}

// RevokeAdmin revokes admin privilages to a user by its ID.
func (a *adminUsers) RevokeAdmin(ctx context.Context, userID string) (*AdminUser, error) {
	if !validStringID(&userID) {
		return nil, ErrInvalidUserValue
	}

	u := fmt.Sprintf("admin/users/%s/actions/revoke_admin", url.QueryEscape(userID))
	req, err := a.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = a.client.do(ctx, req, au)
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

	u := fmt.Sprintf("admin/users/%s/actions/disable_two_factor", url.QueryEscape(userID))
	req, err := a.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	au := &AdminUser{}
	err = a.client.do(ctx, req, au)
	if err != nil {
		return nil, err
	}

	return au, nil
}
