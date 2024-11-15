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
var _ AdminSentinelVersions = (*adminSentinelVersions)(nil)

// AdminSentinelVersions describes all the admin Sentinel versions related methods that
// the Terraform Enterprise API supports.
// Note that admin Sentinel versions are only available in Terraform Enterprise.
//
// TFE API docs: https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/sentinel-versions
type AdminSentinelVersions interface {
	// List all the Sentinel versions.
	List(ctx context.Context, options *AdminSentinelVersionsListOptions) (*AdminSentinelVersionsList, error)

	// Read a Sentinel version by its ID.
	Read(ctx context.Context, id string) (*AdminSentinelVersion, error)

	// Create a Sentinel version.
	Create(ctx context.Context, options AdminSentinelVersionCreateOptions) (*AdminSentinelVersion, error)

	// Update a Sentinel version.
	Update(ctx context.Context, id string, options AdminSentinelVersionUpdateOptions) (*AdminSentinelVersion, error)

	// Delete a Sentinel version
	Delete(ctx context.Context, id string) error
}

// adminSentinelVersions implements AdminSentinelVersions.
type adminSentinelVersions struct {
	client *Client
}

// AdminSentinelVersion represents a Sentinel Version
type AdminSentinelVersion struct {
	ID               string    `jsonapi:"primary,sentinel-versions"`
	Version          string    `jsonapi:"attr,version"`
	URL              string    `jsonapi:"attr,url"`
	SHA              string    `jsonapi:"attr,sha"`
	Deprecated       bool      `jsonapi:"attr,deprecated"`
	DeprecatedReason *string   `jsonapi:"attr,deprecated-reason,omitempty"`
	Official         bool      `jsonapi:"attr,official"`
	Enabled          bool      `jsonapi:"attr,enabled"`
	Beta             bool      `jsonapi:"attr,beta"`
	Usage            int       `jsonapi:"attr,usage"`
	CreatedAt        time.Time `jsonapi:"attr,created-at,iso8601"`
}

// AdminSentinelVersionsListOptions represents the options for listing
// Sentinel versions.
type AdminSentinelVersionsListOptions struct {
	ListOptions

	// Optional: A query string to find an exact version
	Filter string `url:"filter[version],omitempty"`

	// Optional: A search query string to find all versions that match version substring
	Search string `url:"search[version],omitempty"`
}

// AdminSentinelVersionCreateOptions for creating an Sentinel version.
type AdminSentinelVersionCreateOptions struct {
	Type             string  `jsonapi:"primary,sentinel-versions"`
	Version          string  `jsonapi:"attr,version"` // Required
	URL              string  `jsonapi:"attr,url"`     // Required
	SHA              string  `jsonapi:"attr,sha"`     // Required
	Official         *bool   `jsonapi:"attr,official,omitempty"`
	Deprecated       *bool   `jsonapi:"attr,deprecated,omitempty"`
	DeprecatedReason *string `jsonapi:"attr,deprecated-reason,omitempty"`
	Enabled          *bool   `jsonapi:"attr,enabled,omitempty"`
	Beta             *bool   `jsonapi:"attr,beta,omitempty"`
}

// AdminSentinelVersionUpdateOptions for updating Sentinel version.
type AdminSentinelVersionUpdateOptions struct {
	Type             string  `jsonapi:"primary,sentinel-versions"`
	Version          *string `jsonapi:"attr,version,omitempty"`
	URL              *string `jsonapi:"attr,url,omitempty"`
	SHA              *string `jsonapi:"attr,sha,omitempty"`
	Official         *bool   `jsonapi:"attr,official,omitempty"`
	Deprecated       *bool   `jsonapi:"attr,deprecated,omitempty"`
	DeprecatedReason *string `jsonapi:"attr,deprecated-reason,omitempty"`
	Enabled          *bool   `jsonapi:"attr,enabled,omitempty"`
	Beta             *bool   `jsonapi:"attr,beta,omitempty"`
}

// AdminSentinelVersionsList represents a list of Sentinel versions.
type AdminSentinelVersionsList struct {
	*Pagination
	Items []*AdminSentinelVersion
}

// List all the Sentinel versions.
func (a *adminSentinelVersions) List(ctx context.Context, options *AdminSentinelVersionsListOptions) (*AdminSentinelVersionsList, error) {
	req, err := a.client.NewRequest("GET", "admin/sentinel-versions", options)
	if err != nil {
		return nil, err
	}

	sl := &AdminSentinelVersionsList{}
	err = req.Do(ctx, sl)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

// Read a Sentinel version by its ID.
func (a *adminSentinelVersions) Read(ctx context.Context, id string) (*AdminSentinelVersion, error) {
	if !validStringID(&id) {
		return nil, ErrInvalidSentinelVersionID
	}

	u := fmt.Sprintf("admin/sentinel-versions/%s", url.PathEscape(id))
	req, err := a.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	sv := &AdminSentinelVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Create a new Sentinel version.
func (a *adminSentinelVersions) Create(ctx context.Context, options AdminSentinelVersionCreateOptions) (*AdminSentinelVersion, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	req, err := a.client.NewRequest("POST", "admin/sentinel-versions", &options)
	if err != nil {
		return nil, err
	}

	sv := &AdminSentinelVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Update an existing Sentinel version.
func (a *adminSentinelVersions) Update(ctx context.Context, id string, options AdminSentinelVersionUpdateOptions) (*AdminSentinelVersion, error) {
	if !validStringID(&id) {
		return nil, ErrInvalidSentinelVersionID
	}

	u := fmt.Sprintf("admin/sentinel-versions/%s", url.PathEscape(id))
	req, err := a.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	sv := &AdminSentinelVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Delete a Sentinel version.
func (a *adminSentinelVersions) Delete(ctx context.Context, id string) error {
	if !validStringID(&id) {
		return ErrInvalidSentinelVersionID
	}

	u := fmt.Sprintf("admin/sentinel-versions/%s", url.PathEscape(id))
	req, err := a.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o AdminSentinelVersionCreateOptions) valid() error {
	if (o == AdminSentinelVersionCreateOptions{}) {
		return ErrRequiredSentinelVerCreateOps
	}
	if o.Version == "" {
		return ErrRequiredVersion
	}
	if o.URL == "" {
		return ErrRequiredURL
	}
	if o.SHA == "" {
		return ErrRequiredSha
	}

	return nil
}
