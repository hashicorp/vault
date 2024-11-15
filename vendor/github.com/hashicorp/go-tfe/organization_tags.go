// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

var _ OrganizationTags = (*organizationTags)(nil)

// OrganizationMemberships describes all the list of tags used with all resources across the organization.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/organization-tags
type OrganizationTags interface {
	// List all tags within an organization
	List(ctx context.Context, organization string, options *OrganizationTagsListOptions) (*OrganizationTagsList, error)

	// Delete tags from an organization
	Delete(ctx context.Context, organization string, options OrganizationTagsDeleteOptions) error

	// Associate an organization's workspace with a tag
	AddWorkspaces(ctx context.Context, tag string, options AddWorkspacesToTagOptions) error
}

// organizationTags implements OrganizationTags.
type organizationTags struct {
	client *Client
}

// OrganizationTagsList represents a list of organization tags
type OrganizationTagsList struct {
	*Pagination
	Items []*OrganizationTag
}

// OrganizationTag represents a Terraform Enterprise Organization tag
type OrganizationTag struct {
	ID string `jsonapi:"primary,tags"`
	// Optional:
	Name string `jsonapi:"attr,name,omitempty"`

	// Optional: Number of workspaces that have this tag
	InstanceCount int `jsonapi:"attr,instance-count,omitempty"`

	// The org this tag belongs to
	Organization *Organization `jsonapi:"relation,organization"`
}

// OrganizationTagsListOptions represents the options for listing organization tags
type OrganizationTagsListOptions struct {
	ListOptions
	// Optional:
	Filter string `url:"filter[exclude][taggable][id],omitempty"`

	// Optional: A search query string. Organization tags are searchable by name likeness.
	Query string `url:"q,omitempty"`
}

// OrganizationTagsDeleteOptions represents the request body for deleting a tag in an organization
type OrganizationTagsDeleteOptions struct {
	IDs []string // Required
}

// AddWorkspacesToTagOptions represents the request body to add a workspace to a tag
type AddWorkspacesToTagOptions struct {
	WorkspaceIDs []string // Required
}

// this represents a single tag ID
type tagID struct {
	ID string `jsonapi:"primary,tags"`
}

// this represents a single workspace ID
type workspaceID struct {
	ID string `jsonapi:"primary,workspaces"`
}

// List all the tags in an organization. You can provide query params through OrganizationTagsListOptions
func (s *organizationTags) List(ctx context.Context, organization string, options *OrganizationTagsListOptions) (*OrganizationTagsList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/tags", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	tags := &OrganizationTagsList{}
	err = req.Do(ctx, tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Delete tags from a Terraform Enterprise organization
func (s *organizationTags) Delete(ctx context.Context, organization string, options OrganizationTagsDeleteOptions) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}

	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("organizations/%s/tags", url.PathEscape(organization))
	var tagsToRemove []*tagID
	for _, id := range options.IDs {
		tagsToRemove = append(tagsToRemove, &tagID{ID: id})
	}

	req, err := s.client.NewRequest("DELETE", u, tagsToRemove)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Add workspaces to a tag
func (s *organizationTags) AddWorkspaces(ctx context.Context, tag string, options AddWorkspacesToTagOptions) error {
	if !validStringID(&tag) {
		return ErrInvalidTag
	}

	if err := options.valid(); err != nil {
		return err
	}

	var workspaces []*workspaceID
	for _, id := range options.WorkspaceIDs {
		workspaces = append(workspaces, &workspaceID{ID: id})
	}

	u := fmt.Sprintf("tags/%s/relationships/workspaces", url.PathEscape(tag))
	req, err := s.client.NewRequest("POST", u, workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (opts *OrganizationTagsDeleteOptions) valid() error {
	if opts.IDs == nil || len(opts.IDs) == 0 {
		return ErrRequiredTagID
	}

	for _, id := range opts.IDs {
		if !validStringID(&id) {
			errorMsg := fmt.Sprintf("%s is not a valid id value", id)
			return errors.New(errorMsg)
		}
	}

	return nil
}

func (w *AddWorkspacesToTagOptions) valid() error {
	if w.WorkspaceIDs == nil || len(w.WorkspaceIDs) == 0 {
		return ErrRequiredTagWorkspaceID
	}

	for _, id := range w.WorkspaceIDs {
		if !validStringID(&id) {
			errorMsg := fmt.Sprintf("%s is not a valid id value", id)
			return errors.New(errorMsg)
		}
	}

	return nil
}
