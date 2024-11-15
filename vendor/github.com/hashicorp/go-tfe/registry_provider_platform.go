// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation
var _ RegistryProviderPlatforms = (*registryProviderPlatforms)(nil)

// RegistryProviderPlatforms describes the registry provider platform methods supported by the Terraform Enterprise API.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/private-registry/provider-versions-platforms#private-provider-versions-and-platforms-api
type RegistryProviderPlatforms interface {
	// Create a provider platform for an organization
	Create(ctx context.Context, versionID RegistryProviderVersionID, options RegistryProviderPlatformCreateOptions) (*RegistryProviderPlatform, error)

	// List all provider platforms for a single version
	List(ctx context.Context, versionID RegistryProviderVersionID, options *RegistryProviderPlatformListOptions) (*RegistryProviderPlatformList, error)

	// Read a provider platform by ID
	Read(ctx context.Context, platformID RegistryProviderPlatformID) (*RegistryProviderPlatform, error)

	// Delete a provider platform
	Delete(ctx context.Context, platformID RegistryProviderPlatformID) error
}

// registryProviders implements RegistryProviders
type registryProviderPlatforms struct {
	client *Client
}

// RegistryProviderPlatform represents a registry provider platform
type RegistryProviderPlatform struct {
	ID                     string `jsonapi:"primary,registry-provider-platforms"`
	OS                     string `jsonapi:"attr,os"`
	Arch                   string `jsonapi:"attr,arch"`
	Filename               string `jsonapi:"attr,filename"`
	Shasum                 string `jsonapi:"attr,shasum"`
	ProviderBinaryUploaded bool   `jsonapi:"attr,provider-binary-uploaded"`

	// Relations
	RegistryProviderVersion *RegistryProviderVersion `jsonapi:"relation,registry-provider-version"`

	// Links
	Links map[string]interface{} `jsonapi:"links,omitempty"`
}

// RegistryProviderPlatformID is the multi key ID for identifying a provider platform
type RegistryProviderPlatformID struct {
	RegistryProviderVersionID
	OS   string
	Arch string
}

// RegistryProviderPlatformCreateOptions represents the set of options for creating a registry provider platform
type RegistryProviderPlatformCreateOptions struct {
	// Required: A valid operating system string
	OS string `jsonapi:"attr,os"`

	// Required: A valid architecture string
	Arch string `jsonapi:"attr,arch"`

	// Required: A valid shasum string
	Shasum string `jsonapi:"attr,shasum"`

	// Required: A valid filename string
	Filename string `jsonapi:"attr,filename"`
}

type RegistryProviderPlatformList struct {
	*Pagination
	Items []*RegistryProviderPlatform
}

type RegistryProviderPlatformListOptions struct {
	ListOptions
}

// Create a new registry provider platform
func (r *registryProviderPlatforms) Create(ctx context.Context, versionID RegistryProviderVersionID, options RegistryProviderPlatformCreateOptions) (*RegistryProviderPlatform, error) {
	if err := versionID.valid(); err != nil {
		return nil, err
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	// POST /organizations/:organization_name/registry-providers/:registry_name/:namespace/:name/versions/:version/platforms
	u := fmt.Sprintf(
		"organizations/%s/registry-providers/%s/%s/%s/versions/%s/platforms",
		url.PathEscape(versionID.OrganizationName),
		url.PathEscape(string(versionID.RegistryName)),
		url.PathEscape(versionID.Namespace),
		url.PathEscape(versionID.Name),
		url.PathEscape(versionID.Version),
	)
	req, err := r.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rpp := &RegistryProviderPlatform{}
	err = req.Do(ctx, rpp)
	if err != nil {
		return nil, err
	}

	return rpp, nil
}

// List all provider platforms for a single version
func (r *registryProviderPlatforms) List(ctx context.Context, versionID RegistryProviderVersionID, options *RegistryProviderPlatformListOptions) (*RegistryProviderPlatformList, error) {
	if err := versionID.valid(); err != nil {
		return nil, err
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	// GET /organizations/:organization_name/registry-providers/:registry_name/:namespace/:name/versions/:version/platforms
	u := fmt.Sprintf(
		"organizations/%s/registry-providers/%s/%s/%s/versions/%s/platforms",
		url.PathEscape(versionID.RegistryProviderID.OrganizationName),
		url.PathEscape(string(versionID.RegistryProviderID.RegistryName)),
		url.PathEscape(versionID.RegistryProviderID.Namespace),
		url.PathEscape(versionID.RegistryProviderID.Name),
		url.PathEscape(versionID.Version),
	)
	req, err := r.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	ppl := &RegistryProviderPlatformList{}
	err = req.Do(ctx, ppl)
	if err != nil {
		return nil, err
	}

	return ppl, nil
}

// Read is used to read an organization's example by ID
func (r *registryProviderPlatforms) Read(ctx context.Context, platformID RegistryProviderPlatformID) (*RegistryProviderPlatform, error) {
	if err := platformID.valid(); err != nil {
		return nil, err
	}

	// GET /organizations/:organization_name/registry-providers/:registry_name/:namespace/:name/versions/:version/platforms/:os/:arch
	u := fmt.Sprintf(
		"organizations/%s/registry-providers/%s/%s/%s/versions/%s/platforms/%s/%s",
		url.PathEscape(platformID.RegistryProviderID.OrganizationName),
		url.PathEscape(string(platformID.RegistryProviderID.RegistryName)),
		url.PathEscape(platformID.RegistryProviderID.Namespace),
		url.PathEscape(platformID.RegistryProviderID.Name),
		url.PathEscape(platformID.RegistryProviderVersionID.Version),
		url.PathEscape(platformID.OS),
		url.PathEscape(platformID.Arch),
	)
	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	rpp := &RegistryProviderPlatform{}
	err = req.Do(ctx, rpp)

	if err != nil {
		return nil, err
	}

	return rpp, nil
}

// Delete a registry provider platform
func (r *registryProviderPlatforms) Delete(ctx context.Context, platformID RegistryProviderPlatformID) error {
	if err := platformID.valid(); err != nil {
		return err
	}

	// DELETE /organizations/:organization_name/registry-providers/:registry_name/:namespace/:name/versions/:version/platforms/:os/:arch
	u := fmt.Sprintf(
		"organizations/%s/registry-providers/%s/%s/%s/versions/%s/platforms/%s/%s",
		url.PathEscape(platformID.OrganizationName),
		url.PathEscape(string(platformID.RegistryName)),
		url.PathEscape(platformID.Namespace),
		url.PathEscape(platformID.Name),
		url.PathEscape(platformID.Version),
		url.PathEscape(platformID.OS),
		url.PathEscape(platformID.Arch),
	)
	req, err := r.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (id RegistryProviderPlatformID) valid() error {
	if err := id.RegistryProviderID.valid(); err != nil {
		return err
	}
	if !validString(&id.OS) {
		return ErrInvalidOS
	}
	if !validString(&id.Arch) {
		return ErrInvalidArch
	}
	return nil
}

func (o RegistryProviderPlatformCreateOptions) valid() error {
	if !validString(&o.OS) {
		return ErrRequiredOS
	}
	if !validString(&o.Arch) {
		return ErrRequiredArch
	}
	if !validStringID(&o.Shasum) {
		return ErrRequiredShasum
	}
	if !validStringID(&o.Filename) {
		return ErrRequiredFilename
	}
	return nil
}

func (o *RegistryProviderPlatformListOptions) valid() error {
	return nil
}
