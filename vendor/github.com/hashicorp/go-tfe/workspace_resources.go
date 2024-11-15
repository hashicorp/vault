// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ WorkspaceResources = (*workspaceResources)(nil)

// WorkspaceResources describes all the workspace resources related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspace-resources
type WorkspaceResources interface {
	// List all the workspaces resources within a workspace
	List(ctx context.Context, workspaceID string, options *WorkspaceResourceListOptions) (*WorkspaceResourcesList, error)
}

// workspaceResources implements WorkspaceResources.
type workspaceResources struct {
	client *Client
}

// WorkspaceResourcesList represents a list of workspace resources.
type WorkspaceResourcesList struct {
	*Pagination
	Items []*WorkspaceResource
}

// WorkspaceResource represents a Terraform Enterprise workspace resource.
type WorkspaceResource struct {
	ID                       string  `jsonapi:"primary,resources"`
	Address                  string  `jsonapi:"attr,address"`
	Name                     string  `jsonapi:"attr,name"`
	CreatedAt                string  `jsonapi:"attr,created-at"`
	UpdatedAt                string  `jsonapi:"attr,updated-at"`
	Module                   string  `jsonapi:"attr,module"`
	Provider                 string  `jsonapi:"attr,provider"`
	ProviderType             string  `jsonapi:"attr,provider-type"`
	ModifiedByStateVersionID string  `jsonapi:"attr,modified-by-state-version-id"`
	NameIndex                *string `jsonapi:"attr,name-index"`
}

// WorkspaceResourceListOptions represents the options for listing workspace resources.
type WorkspaceResourceListOptions struct {
	ListOptions
}

// List all the workspaces resources within a workspace
func (s *workspaceResources) List(ctx context.Context, workspaceID string, options *WorkspaceResourceListOptions) (*WorkspaceResourcesList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/resources", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	wl := &WorkspaceResourcesList{}
	err = req.Do(ctx, wl)
	if err != nil {
		return nil, err
	}
	return wl, nil
}

func (o *WorkspaceResourceListOptions) valid() error {
	return nil
}
