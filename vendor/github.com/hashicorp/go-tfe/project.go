// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Projects = (*projects)(nil)

// Projects describes all the project related methods that the Terraform
// Enterprise API supports
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/projects
type Projects interface {
	// List all projects in the given organization
	List(ctx context.Context, organization string, options *ProjectListOptions) (*ProjectList, error)

	// Create a new project.
	Create(ctx context.Context, organization string, options ProjectCreateOptions) (*Project, error)

	// Read a project by its ID.
	Read(ctx context.Context, projectID string) (*Project, error)

	// Update a project.
	Update(ctx context.Context, projectID string, options ProjectUpdateOptions) (*Project, error)

	// Delete a project.
	Delete(ctx context.Context, projectID string) error

	// ListTagBindings lists all tag bindings associated with the project.
	ListTagBindings(ctx context.Context, projectID string) ([]*TagBinding, error)

	// AddTagBindings adds or modifies the value of existing tag binding keys for a project.
	AddTagBindings(ctx context.Context, projectID string, options ProjectAddTagBindingsOptions) ([]*TagBinding, error)
}

// projects implements Projects
type projects struct {
	client *Client
}

// ProjectList represents a list of projects
type ProjectList struct {
	*Pagination
	Items []*Project
}

// Project represents a Terraform Enterprise project
type Project struct {
	ID        string `jsonapi:"primary,projects"`
	IsUnified bool   `jsonapi:"attr,is-unified"`
	Name      string `jsonapi:"attr,name"`

	Description string `jsonapi:"attr,description"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
}

// ProjectListOptions represents the options for listing projects
type ProjectListOptions struct {
	ListOptions

	// Optional: String (complete project name) used to filter the results.
	// If multiple, comma separated values are specified, projects matching
	// any of the names are returned.
	Name string `url:"filter[names],omitempty"`

	// Optional: A query string to search projects by names.
	Query string `url:"q,omitempty"`

	// Optional: A filter string to list projects filtered by key/value tags.
	// These are not annotated and therefore not encoded by go-querystring
	TagBindings []*TagBinding
}

// ProjectCreateOptions represents the options for creating a project
type ProjectCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,projects"`

	// Required: A name to identify the project.
	Name string `jsonapi:"attr,name"`

	// Optional: A description for the project.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Associated TagBindings of the project.
	TagBindings []*TagBinding `jsonapi:"relation,tag-bindings,omitempty"`
}

// ProjectUpdateOptions represents the options for updating a project
type ProjectUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,projects"`

	// Optional: A name to identify the project
	Name *string `jsonapi:"attr,name,omitempty"`

	// Optional: A description for the project.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Associated TagBindings of the project. Note that this will replace
	// all existing tag bindings.
	TagBindings []*TagBinding `jsonapi:"relation,tag-bindings,omitempty"`
}

// ProjectAddTagBindingsOptions represents the options for adding tag bindings
// to a project.
type ProjectAddTagBindingsOptions struct {
	TagBindings []*TagBinding
}

// List all projects.
func (s *projects) List(ctx context.Context, organization string, options *ProjectListOptions) (*ProjectList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	var tagFilters map[string][]string
	if options != nil {
		tagFilters = encodeTagFiltersAsParams(options.TagBindings)
	}

	u := fmt.Sprintf("organizations/%s/projects", url.PathEscape(organization))
	req, err := s.client.NewRequestWithAdditionalQueryParams("GET", u, options, tagFilters)
	if err != nil {
		return nil, err
	}

	p := &ProjectList{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Create a project with the given options
func (s *projects) Create(ctx context.Context, organization string, options ProjectCreateOptions) (*Project, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/projects", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	p := &Project{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Read a single project by its ID.
func (s *projects) Read(ctx context.Context, projectID string) (*Project, error) {
	if !validStringID(&projectID) {
		return nil, ErrInvalidProjectID
	}

	u := fmt.Sprintf("projects/%s", url.PathEscape(projectID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	p := &Project{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *projects) ListTagBindings(ctx context.Context, projectID string) ([]*TagBinding, error) {
	if !validStringID(&projectID) {
		return nil, ErrInvalidProjectID
	}

	u := fmt.Sprintf("projects/%s/tag-bindings", url.PathEscape(projectID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var list struct {
		*Pagination
		Items []*TagBinding
	}

	err = req.Do(ctx, &list)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

// AddTagBindings adds or modifies the value of existing tag binding keys for a project
func (s *projects) AddTagBindings(ctx context.Context, projectID string, options ProjectAddTagBindingsOptions) ([]*TagBinding, error) {
	if !validStringID(&projectID) {
		return nil, ErrInvalidProjectID
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("projects/%s/tag-bindings", url.PathEscape(projectID))
	req, err := s.client.NewRequest("PATCH", u, options.TagBindings)
	if err != nil {
		return nil, err
	}

	var response = struct {
		*Pagination
		Items []*TagBinding
	}{}
	err = req.Do(ctx, &response)

	return response.Items, err
}

// Update a project by its ID
func (s *projects) Update(ctx context.Context, projectID string, options ProjectUpdateOptions) (*Project, error) {
	if !validStringID(&projectID) {
		return nil, ErrInvalidProjectID
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("projects/%s", url.PathEscape(projectID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	p := &Project{}
	err = req.Do(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Delete a project by its ID
func (s *projects) Delete(ctx context.Context, projectID string) error {
	if !validStringID(&projectID) {
		return ErrInvalidProjectID
	}

	u := fmt.Sprintf("projects/%s", url.PathEscape(projectID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o ProjectCreateOptions) valid() error {
	if !validString(&o.Name) {
		return ErrRequiredName
	}
	return nil
}

func (o ProjectUpdateOptions) valid() error {
	return nil
}

func (o ProjectAddTagBindingsOptions) valid() error {
	if len(o.TagBindings) == 0 {
		return ErrRequiredTagBindings
	}

	return nil
}
