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
var _ AdminOPAVersions = (*adminOPAVersions)(nil)

// AdminOPAVersions describes all the admin OPA versions related methods that
// the Terraform Enterprise API supports.
// Note that admin OPA versions are only available in Terraform Enterprise.
//
// TFE API docs: https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/opa-versions
type AdminOPAVersions interface {
	// List all the OPA versions.
	List(ctx context.Context, options *AdminOPAVersionsListOptions) (*AdminOPAVersionsList, error)

	// Read a OPA version by its ID.
	Read(ctx context.Context, id string) (*AdminOPAVersion, error)

	// Create a OPA version.
	Create(ctx context.Context, options AdminOPAVersionCreateOptions) (*AdminOPAVersion, error)

	// Update a OPA version.
	Update(ctx context.Context, id string, options AdminOPAVersionUpdateOptions) (*AdminOPAVersion, error)

	// Delete a OPA version
	Delete(ctx context.Context, id string) error
}

// adminOPAVersions implements AdminOPAVersions.
type adminOPAVersions struct {
	client *Client
}

// AdminOPAVersion represents a OPA Version
type AdminOPAVersion struct {
	ID               string    `jsonapi:"primary,opa-versions"`
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

// AdminOPAVersionsListOptions represents the options for listing
// OPA versions.
type AdminOPAVersionsListOptions struct {
	ListOptions

	// Optional: A query string to find an exact version
	Filter string `url:"filter[version],omitempty"`

	// Optional: A search query string to find all versions that match version substring
	Search string `url:"search[version],omitempty"`
}

// AdminOPAVersionCreateOptions for creating an OPA version.
type AdminOPAVersionCreateOptions struct {
	Type             string  `jsonapi:"primary,opa-versions"`
	Version          string  `jsonapi:"attr,version"` // Required
	URL              string  `jsonapi:"attr,url"`     // Required
	SHA              string  `jsonapi:"attr,sha"`     // Required
	Official         *bool   `jsonapi:"attr,official,omitempty"`
	Deprecated       *bool   `jsonapi:"attr,deprecated,omitempty"`
	DeprecatedReason *string `jsonapi:"attr,deprecated-reason,omitempty"`
	Enabled          *bool   `jsonapi:"attr,enabled,omitempty"`
	Beta             *bool   `jsonapi:"attr,beta,omitempty"`
}

// AdminOPAVersionUpdateOptions for updating OPA version.
type AdminOPAVersionUpdateOptions struct {
	Type             string  `jsonapi:"primary,opa-versions"`
	Version          *string `jsonapi:"attr,version,omitempty"`
	URL              *string `jsonapi:"attr,url,omitempty"`
	SHA              *string `jsonapi:"attr,sha,omitempty"`
	Official         *bool   `jsonapi:"attr,official,omitempty"`
	Deprecated       *bool   `jsonapi:"attr,deprecated,omitempty"`
	DeprecatedReason *string `jsonapi:"attr,deprecated-reason,omitempty"`
	Enabled          *bool   `jsonapi:"attr,enabled,omitempty"`
	Beta             *bool   `jsonapi:"attr,beta,omitempty"`
}

// AdminOPAVersionsList represents a list of OPA versions.
type AdminOPAVersionsList struct {
	*Pagination
	Items []*AdminOPAVersion
}

// List all the OPA versions.
func (a *adminOPAVersions) List(ctx context.Context, options *AdminOPAVersionsListOptions) (*AdminOPAVersionsList, error) {
	req, err := a.client.NewRequest("GET", "admin/opa-versions", options)
	if err != nil {
		return nil, err
	}

	ol := &AdminOPAVersionsList{}
	err = req.Do(ctx, ol)
	if err != nil {
		return nil, err
	}

	return ol, nil
}

// Read a OPA version by its ID.
func (a *adminOPAVersions) Read(ctx context.Context, id string) (*AdminOPAVersion, error) {
	if !validStringID(&id) {
		return nil, ErrInvalidOPAVersionID
	}

	u := fmt.Sprintf("admin/opa-versions/%s", url.PathEscape(id))
	req, err := a.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	ov := &AdminOPAVersion{}
	err = req.Do(ctx, ov)
	if err != nil {
		return nil, err
	}

	return ov, nil
}

// Create a new OPA version.
func (a *adminOPAVersions) Create(ctx context.Context, options AdminOPAVersionCreateOptions) (*AdminOPAVersion, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	req, err := a.client.NewRequest("POST", "admin/opa-versions", &options)
	if err != nil {
		return nil, err
	}

	ov := &AdminOPAVersion{}
	err = req.Do(ctx, ov)
	if err != nil {
		return nil, err
	}

	return ov, nil
}

// Update an existing OPA version.
func (a *adminOPAVersions) Update(ctx context.Context, id string, options AdminOPAVersionUpdateOptions) (*AdminOPAVersion, error) {
	if !validStringID(&id) {
		return nil, ErrInvalidOPAVersionID
	}

	u := fmt.Sprintf("admin/opa-versions/%s", url.PathEscape(id))
	req, err := a.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	ov := &AdminOPAVersion{}
	err = req.Do(ctx, ov)
	if err != nil {
		return nil, err
	}

	return ov, nil
}

// Delete a OPA version.
func (a *adminOPAVersions) Delete(ctx context.Context, id string) error {
	if !validStringID(&id) {
		return ErrInvalidOPAVersionID
	}

	u := fmt.Sprintf("admin/opa-versions/%s", url.PathEscape(id))
	req, err := a.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o AdminOPAVersionCreateOptions) valid() error {
	if (o == AdminOPAVersionCreateOptions{}) {
		return ErrRequiredOPAVerCreateOps
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
