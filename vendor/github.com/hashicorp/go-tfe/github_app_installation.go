// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ GHAInstallations = (*gHAInstallations)(nil)

// GHAInstallations describes all the GitHub App Installation related methods that the
// Terraform Enterprise API supports. The APIs require the user token for the user who
// already has the GitHub App Installation set up via the UI.
// (https://developer.hashicorp.com/terraform/enterprise/admin/application/github-app-integration)
type GHAInstallations interface {
	// List all the GitHub App Installations for the user.
	List(ctx context.Context, options *GHAInstallationListOptions) (*GHAInstallationList, error)

	// Read a GitHub App Installations by its external id.
	Read(ctx context.Context, GHAInstallationID string) (*GHAInstallation, error)
}

// gHAInstallations implements GHAInstallations.
type gHAInstallations struct {
	client *Client
}

// GHAInstallationList represents a list of github installations.
type GHAInstallationList struct {
	*Pagination
	Items []*GHAInstallation
}

// GHAInstallation represents a github app installation
type GHAInstallation struct {
	ID             *string `jsonapi:"primary,github-app-installations"`
	InstallationID *int    `jsonapi:"attr,installation-id"`
	Name           *string `jsonapi:"attr,name"`
}

// GHAInstallationListOptions represents the options for listing.
type GHAInstallationListOptions struct {
	ListOptions
}

// List all the github app installations.
func (s *gHAInstallations) List(ctx context.Context, options *GHAInstallationListOptions) (*GHAInstallationList, error) {
	u := "github-app/installations"
	req, err := s.client.NewRequest("GET", u, options)
	fmt.Println(u)
	if err != nil {
		return nil, err
	}

	ghil := &GHAInstallationList{}

	err = req.Do(ctx, ghil)
	if err != nil {
		return nil, err
	}

	return ghil, nil
}

// Read a GitHub App Installations by its ID.
func (s *gHAInstallations) Read(ctx context.Context, id string) (*GHAInstallation, error) {
	if !validStringID(&id) {
		return nil, ErrInvalidOauthClientID
	}

	u := fmt.Sprintf("github-app/installation/%s", url.PathEscape(id))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	ghi := &GHAInstallation{}
	err = req.Do(ctx, ghi)
	if err != nil {
		return nil, err
	}

	return ghi, err
}
