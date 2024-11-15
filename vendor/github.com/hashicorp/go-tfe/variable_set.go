// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ VariableSets = (*variableSets)(nil)

// VariableSets describes all the Variable Set related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/variable-sets
type VariableSets interface {
	// List all the variable sets within an organization.
	List(ctx context.Context, organization string, options *VariableSetListOptions) (*VariableSetList, error)

	// ListForWorkspace gets the associated variable sets for a workspace.
	ListForWorkspace(ctx context.Context, workspaceID string, options *VariableSetListOptions) (*VariableSetList, error)

	// ListForProject gets the associated variable sets for a project.
	ListForProject(ctx context.Context, projectID string, options *VariableSetListOptions) (*VariableSetList, error)

	// Create is used to create a new variable set.
	Create(ctx context.Context, organization string, options *VariableSetCreateOptions) (*VariableSet, error)

	// Read a variable set by its ID.
	Read(ctx context.Context, variableSetID string, options *VariableSetReadOptions) (*VariableSet, error)

	// Update an existing variable set.
	Update(ctx context.Context, variableSetID string, options *VariableSetUpdateOptions) (*VariableSet, error)

	// Delete a variable set by ID.
	Delete(ctx context.Context, variableSetID string) error

	// Apply variable set to workspaces in the supplied list.
	ApplyToWorkspaces(ctx context.Context, variableSetID string, options *VariableSetApplyToWorkspacesOptions) error

	// Remove variable set from workspaces in the supplied list.
	RemoveFromWorkspaces(ctx context.Context, variableSetID string, options *VariableSetRemoveFromWorkspacesOptions) error

	// Apply variable set to projects in the supplied list.
	ApplyToProjects(ctx context.Context, variableSetID string, options VariableSetApplyToProjectsOptions) error

	// Remove variable set from projects in the supplied list.
	RemoveFromProjects(ctx context.Context, variableSetID string, options VariableSetRemoveFromProjectsOptions) error

	// Update list of workspaces to which the variable set is applied to match the supplied list.
	UpdateWorkspaces(ctx context.Context, variableSetID string, options *VariableSetUpdateWorkspacesOptions) (*VariableSet, error)
}

// variableSets implements VariableSets.
type variableSets struct {
	client *Client
}

// VariableSetList represents a list of variable sets.
type VariableSetList struct {
	*Pagination
	Items []*VariableSet
}

// Parent represents the variable set's parent (currently only organizations and projects are supported).
// This relation is considered BETA, SUBJECT TO CHANGE, and likely unavailable to most users.
type Parent struct {
	Organization *Organization
	Project      *Project
}

// VariableSet represents a Terraform Enterprise variable set.
type VariableSet struct {
	ID          string `jsonapi:"primary,varsets"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
	Global      bool   `jsonapi:"attr,global"`
	Priority    bool   `jsonapi:"attr,priority"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
	// Optional: Parent represents the variable set's parent (currently only organizations and projects are supported).
	// This relation is considered BETA, SUBJECT TO CHANGE, and likely unavailable to most users.
	Parent     *Parent                `jsonapi:"polyrelation,parent"`
	Workspaces []*Workspace           `jsonapi:"relation,workspaces,omitempty"`
	Projects   []*Project             `jsonapi:"relation,projects,omitempty"`
	Variables  []*VariableSetVariable `jsonapi:"relation,vars,omitempty"`
}

// A list of relations to include. See available resources
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/organizations#available-related-resources
type VariableSetIncludeOpt string

const (
	VariableSetWorkspaces VariableSetIncludeOpt = "workspaces"
	VariableSetProjects   VariableSetIncludeOpt = "projects"
	VariableSetVars       VariableSetIncludeOpt = "vars"
)

// VariableSetListOptions represents the options for listing variable sets.
type VariableSetListOptions struct {
	ListOptions
	Include string `url:"include"`

	// Optional: A query string used to filter variable sets.
	// Any variable sets with a name partially matching this value will be returned.
	Query string `url:"q,omitempty"`
}

// VariableSetCreateOptions represents the options for creating a new variable set within in a organization.
type VariableSetCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,varsets"`

	// The name of the variable set.
	// Affects variable precedence when there are conflicts between Variable Sets
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/variable-sets#apply-variable-set-to-workspaces
	Name *string `jsonapi:"attr,name"`

	// A description to provide context for the variable set.
	Description *string `jsonapi:"attr,description,omitempty"`

	// If true the variable set is considered in all runs in the organization.
	Global *bool `jsonapi:"attr,global,omitempty"`

	// If true the variables in the set override any other variable values set
	// in a more specific scope including values set on the command line.
	Priority *bool `jsonapi:"attr,priority,omitempty"`

	// Optional: Parent represents the variable set's parent (currently only organizations and projects are supported).
	// This relation is considered BETA, SUBJECT TO CHANGE, and likely unavailable to most users.
	Parent *Parent `jsonapi:"polyrelation,parent"`
}

// VariableSetReadOptions represents the options for reading variable sets.
type VariableSetReadOptions struct {
	Include *[]VariableSetIncludeOpt `url:"include,omitempty"`
}

// VariableSetUpdateOptions represents the options for updating a variable set.
type VariableSetUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,varsets"`

	// The name of the variable set.
	// Affects variable precedence when there are conflicts between Variable Sets
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/variable-sets#apply-variable-set-to-workspaces
	Name *string `jsonapi:"attr,name,omitempty"`

	// A description to provide context for the variable set.
	Description *string `jsonapi:"attr,description,omitempty"`

	// If true the variable set is considered in all runs in the organization.
	Global *bool `jsonapi:"attr,global,omitempty"`

	// If true the variables in the set override any other variable values set
	// in a more specific scope including values set on the command line.
	Priority *bool `jsonapi:"attr,priority,omitempty"`
}

// VariableSetApplyToWorkspacesOptions represents the options for applying variable sets to workspaces.
type VariableSetApplyToWorkspacesOptions struct {
	// The workspaces to apply the variable set to (additive).
	Workspaces []*Workspace
}

// VariableSetRemoveFromWorkspacesOptions represents the options for removing variable sets from workspaces.
type VariableSetRemoveFromWorkspacesOptions struct {
	// The workspaces to remove the variable set from.
	Workspaces []*Workspace
}

// VariableSetApplyToProjectsOptions represents the options for applying variable sets to projects.
type VariableSetApplyToProjectsOptions struct {
	// The projects to apply the variable set to (additive).
	Projects []*Project
}

// VariableSetRemoveFromProjectsOptions represents the options for removing variable sets from projects.
type VariableSetRemoveFromProjectsOptions struct {
	// The projects to remove the variable set from.
	Projects []*Project
}

// VariableSetUpdateWorkspacesOptions represents a subset of update options specifically for applying variable sets to workspaces
type VariableSetUpdateWorkspacesOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,varsets"`

	// The workspaces to be applied to. An empty set means remove all applied
	Workspaces []*Workspace `jsonapi:"relation,workspaces"`
}

type privateVariableSetUpdateWorkspacesOptions struct {
	Type       string       `jsonapi:"primary,varsets"`
	Global     bool         `jsonapi:"attr,global"`
	Workspaces []*Workspace `jsonapi:"relation,workspaces"`
}

// List all Variable Sets in the organization
func (s *variableSets) List(ctx context.Context, organization string, options *VariableSetListOptions) (*VariableSetList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if options != nil {
		if err := options.valid(); err != nil {
			return nil, err
		}
	}

	u := fmt.Sprintf("organizations/%s/varsets", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	vl := &VariableSetList{}
	err = req.Do(ctx, vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// ListForWorkspace gets the associated variable sets for a workspace.
func (s *variableSets) ListForWorkspace(ctx context.Context, workspaceID string, options *VariableSetListOptions) (*VariableSetList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if options != nil {
		if err := options.valid(); err != nil {
			return nil, err
		}
	}

	u := fmt.Sprintf("workspaces/%s/varsets", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	vl := &VariableSetList{}
	err = req.Do(ctx, vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// ListForProject gets the associated variable sets for a project.
func (s *variableSets) ListForProject(ctx context.Context, projectID string, options *VariableSetListOptions) (*VariableSetList, error) {
	if !validStringID(&projectID) {
		return nil, ErrInvalidProjectID
	}
	if options != nil {
		if err := options.valid(); err != nil {
			return nil, err
		}
	}

	u := fmt.Sprintf("projects/%s/varsets", url.PathEscape(projectID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	vl := &VariableSetList{}
	err = req.Do(ctx, vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// Create is used to create a new variable set.
func (s *variableSets) Create(ctx context.Context, organization string, options *VariableSetCreateOptions) (*VariableSet, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/varsets", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, options)
	if err != nil {
		return nil, err
	}

	vl := &VariableSet{}
	err = req.Do(ctx, vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}

// Read is used to inspect a given variable set based on ID
func (s *variableSets) Read(ctx context.Context, variableSetID string, options *VariableSetReadOptions) (*VariableSet, error) {
	if !validStringID(&variableSetID) {
		return nil, ErrInvalidVariableSetID
	}

	u := fmt.Sprintf("varsets/%s", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	vs := &VariableSet{}
	err = req.Do(ctx, vs)
	if err != nil {
		return nil, err
	}

	return vs, err
}

// Update an existing variable set.
func (s *variableSets) Update(ctx context.Context, variableSetID string, options *VariableSetUpdateOptions) (*VariableSet, error) {
	if !validStringID(&variableSetID) {
		return nil, ErrInvalidVariableSetID
	}

	u := fmt.Sprintf("varsets/%s", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("PATCH", u, options)
	if err != nil {
		return nil, err
	}

	v := &VariableSet{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Delete a variable set by its ID.
func (s *variableSets) Delete(ctx context.Context, variableSetID string) error {
	if !validStringID(&variableSetID) {
		return ErrInvalidVariableSetID
	}

	u := fmt.Sprintf("varsets/%s", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Apply variable set to workspaces in the supplied list.
// Note: this method will return an error if the variable set has global = true.
func (s *variableSets) ApplyToWorkspaces(ctx context.Context, variableSetID string, options *VariableSetApplyToWorkspacesOptions) error {
	if !validStringID(&variableSetID) {
		return ErrInvalidVariableSetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("varsets/%s/relationships/workspaces", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("POST", u, options.Workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Remove variable set from workspaces in the supplied list.
// Note: this method will return an error if the variable set has global = true.
func (s *variableSets) RemoveFromWorkspaces(ctx context.Context, variableSetID string, options *VariableSetRemoveFromWorkspacesOptions) error {
	if !validStringID(&variableSetID) {
		return ErrInvalidVariableSetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("varsets/%s/relationships/workspaces", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("DELETE", u, options.Workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// ApplyToProjects applies the variable set to projects in the supplied list.
// This method will return an error if the variable set has global = true.
func (s variableSets) ApplyToProjects(ctx context.Context, variableSetID string, options VariableSetApplyToProjectsOptions) error {
	if !validStringID(&variableSetID) {
		return ErrInvalidVariableSetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("varsets/%s/relationships/projects", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("POST", u, options.Projects)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveFromProjects removes the variable set from projects in the supplied list.
// This method will return an error if the variable set has global = true.
func (s variableSets) RemoveFromProjects(ctx context.Context, variableSetID string, options VariableSetRemoveFromProjectsOptions) error {
	if !validStringID(&variableSetID) {
		return ErrInvalidVariableSetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("varsets/%s/relationships/projects", url.PathEscape(variableSetID))
	req, err := s.client.NewRequest("DELETE", u, options.Projects)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Update variable set to be applied to only the workspaces in the supplied list.
func (s *variableSets) UpdateWorkspaces(ctx context.Context, variableSetID string, options *VariableSetUpdateWorkspacesOptions) (*VariableSet, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Use private struct to ensure global is set to false when applying to workspaces
	o := privateVariableSetUpdateWorkspacesOptions{
		Global:     bool(false),
		Workspaces: options.Workspaces,
	}

	// We force inclusion of workspaces as that is the primary data for which we are concerned with confirming changes.
	u := fmt.Sprintf("varsets/%s?include=%s", url.PathEscape(variableSetID), VariableSetWorkspaces)
	req, err := s.client.NewRequest("PATCH", u, &o)
	if err != nil {
		return nil, err
	}

	v := &VariableSet{}
	err = req.Do(ctx, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (o *VariableSetListOptions) valid() error {
	return nil
}

func (o *VariableSetCreateOptions) valid() error {
	if o == nil {
		return nil
	}
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if o.Global == nil {
		return ErrRequiredGlobalFlag
	}
	return nil
}

func (o *VariableSetApplyToWorkspacesOptions) valid() error {
	for _, s := range o.Workspaces {
		if !validStringID(&s.ID) {
			return ErrRequiredWorkspaceID
		}
	}
	return nil
}

func (o *VariableSetRemoveFromWorkspacesOptions) valid() error {
	for _, s := range o.Workspaces {
		if !validStringID(&s.ID) {
			return ErrRequiredWorkspaceID
		}
	}
	return nil
}

func (o *VariableSetApplyToProjectsOptions) valid() error {
	for _, s := range o.Projects {
		if !validStringID(&s.ID) {
			return ErrRequiredProjectID
		}
	}
	return nil
}

func (o VariableSetRemoveFromProjectsOptions) valid() error {
	for _, s := range o.Projects {
		if !validStringID(&s.ID) {
			return ErrRequiredProjectID
		}
	}
	return nil
}

func (o *VariableSetUpdateWorkspacesOptions) valid() error {
	if o == nil || o.Workspaces == nil {
		return ErrRequiredWorkspacesList
	}
	return nil
}
