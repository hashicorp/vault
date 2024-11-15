// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ RegistryNoCodeModules = (*registryNoCodeModules)(nil)

// RegistryNoCodeModules describes all the registry no-code module related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: (TODO: Add link to API docs)
type RegistryNoCodeModules interface {

	// Create a registry no-code module
	// **Note: This API is still in BETA and subject to change.**
	Create(ctx context.Context, organization string, options RegistryNoCodeModuleCreateOptions) (*RegistryNoCodeModule, error)

	// Read a registry no-code  module
	// **Note: This API is still in BETA and subject to change.**
	Read(ctx context.Context, noCodeModuleID string, options *RegistryNoCodeModuleReadOptions) (*RegistryNoCodeModule, error)

	// ReadVariables returns the variables for a version of a no-code module
	// **Note: This API is still in BETA and subject to change.**
	ReadVariables(ctx context.Context, noCodeModuleID, noCodeModuleVersion string, options *RegistryNoCodeModuleReadVariablesOptions) (*RegistryModuleVariableList, error)

	// Update a registry no-code module
	// **Note: This API is still in BETA and subject to change.**
	Update(ctx context.Context, noCodeModuleID string, options RegistryNoCodeModuleUpdateOptions) (*RegistryNoCodeModule, error)

	// Delete a registry no-code module
	// **Note: This API is still in BETA and subject to change.**
	Delete(ctx context.Context, ID string) error

	// CreateWorkspace creates a workspace using a no-code module.
	CreateWorkspace(ctx context.Context, noCodeModuleID string, options *RegistryNoCodeModuleCreateWorkspaceOptions) (*Workspace, error)

	// UpgradeWorkspace initiates an upgrade of an existing no-code module workspace.
	UpgradeWorkspace(ctx context.Context, noCodeModuleID string, workspaceID string, options *RegistryNoCodeModuleUpgradeWorkspaceOptions) (*WorkspaceUpgrade, error)
}

// RegistryModuleVariableList is a list of registry module variables.
// **Note: This API is still in BETA and subject to change.**
type RegistryModuleVariableList struct {
	Items []*RegistryModuleVariable

	// NOTE: At the time of authoring this comment, the API endpoint to fetch
	// registry module variables does not support pagination. This field is
	// included to satisfy jsonapi unmarshaler implementation here:
	// https://github.com/hashicorp/go-tfe/blob/3d29602707fa4b10469d1a02685644bd159d3ccc/tfe.go#L859
	*Pagination
}

// RegistryModuleVariable represents a registry module variable.
type RegistryModuleVariable struct {
	// ID is the ID of the variable.
	ID string `jsonapi:"primary,registry-module-variables"`

	// Name is the name of the variable.
	Name string `jsonapi:"attr,name"`

	// VariableType is the type of the variable.
	VariableType string `jsonapi:"attr,type"`

	// Description is the description of the variable.
	Description string `jsonapi:"attr,description"`

	// Required is a boolean indicating if the variable is required.
	Required bool `jsonapi:"attr,required"`

	// Sensitive is a boolean indicating if the variable is sensitive.
	Sensitive bool `jsonapi:"attr,sensitive"`

	// Options is a slice of strings representing the options for the variable.
	Options []string `jsonapi:"attr,options"`

	// HasGlobal is a boolean indicating if the variable is global.
	HasGlobal bool `jsonapi:"attr,has-global"`
}

type RegistryNoCodeModuleCreateWorkspaceOptions struct {
	Type string `jsonapi:"primary,no-code-module-workspace"`

	// Name is the name of the workspace, which can only include letters,
	// numbers, and _. This will be used as an identifier and must be unique in
	// the organization.
	Name string `jsonapi:"attr,name"`

	// Description is a description for the workspace.
	Description *string `jsonapi:"attr,description,omitempty"`

	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// Project is the associated project with the workspace. If not provided,
	// default project of the organization will be assigned to the workspace.
	Project *Project `jsonapi:"relation,project,omitempty"`

	// Variables is the slice of variables to be configured for the no-code
	// workspace.
	Variables []*Variable `jsonapi:"relation,vars,omitempty"`

	// SourceName is the name of the source of the workspace.
	SourceName *string `jsonapi:"attr,source-name,omitempty"`

	// SourceUrl is the URL of the source of the workspace.
	SourceURL *string `jsonapi:"attr,source-url,omitempty"`

	// ExecutionMode is the execution mode of the workspace.
	ExecutionMode *string `jsonapi:"attr,execution-mode,omitempty"`

	// AgentPoolId is the ID of the agent pool to use for the workspace.
	// This is required when execution mode is set to "agent".
	// This must not be specified when execution mode is set to "remote".
	AgentPoolID *string `jsonapi:"attr,agent-pool-id,omitempty"`
}

type RegistryNoCodeModuleUpgradeWorkspaceOptions struct {
	Type string `jsonapi:"primary,no-code-module-workspace"`

	// Variables is the slice of variables to be configured for the no-code
	// workspace.
	Variables []*Variable `jsonapi:"relation,vars,omitempty"`
}

// registryNoCodeModules implements RegistryNoCodeModules.
type registryNoCodeModules struct {
	client *Client
}

// RegistryNoCodeModule represents a registry no-code module
type RegistryNoCodeModule struct {
	ID         string `jsonapi:"primary,no-code-modules"`
	VersionPin string `jsonapi:"attr,version-pin"`
	Enabled    bool   `jsonapi:"attr,enabled"`

	// Relations
	Organization    *Organization           `jsonapi:"relation,organization"`
	RegistryModule  *RegistryModule         `jsonapi:"relation,registry-module"`
	VariableOptions []*NoCodeVariableOption `jsonapi:"relation,variable-options"`
}

// NoCodeVariableOption represents a registry no-code module variable and its
// options.
type NoCodeVariableOption struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	Type string `jsonapi:"primary,variable-options"`

	// Required: The variable name
	VariableName string `jsonapi:"attr,variable-name"`

	// Required: The variable type
	VariableType string `jsonapi:"attr,variable-type"`

	// Optional: The options for the variable
	Options []string `jsonapi:"attr,options"`
}

// RegistryNoCodeModuleCreateOptions is used when creating a registry no-code module
type RegistryNoCodeModuleCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,no-code-modules"`

	// Required: the registry module to use for the no-code module (only the ID is used)
	RegistryModule *RegistryModule `jsonapi:"relation,registry-module"`

	// Optional: whether no-code is enabled for the module
	Enabled *bool `jsonapi:"attr,enabled,omitempty"`

	// Optional: the version pin for the module. valid values are "latest" or a semver string
	VersionPin string `jsonapi:"attr,version-pin,omitempty"`

	// Optional: the variable options for the registry module
	VariableOptions []*NoCodeVariableOption `jsonapi:"relation,variable-options,omitempty"`
}

// RegistryNoCodeModuleIncludeOpt represents the available options for include query params.
type RegistryNoCodeModuleIncludeOpt string

var (
	// RegistryNoCodeIncludeVariableOptions is used to include variable options in the response
	RegistryNoCodeIncludeVariableOptions RegistryNoCodeModuleIncludeOpt = "variable-options"
)

// RegistryNoCodeModuleReadOptions is used when reading a registry no-code module
type RegistryNoCodeModuleReadOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-updating
	Type string `jsonapi:"primary,no-code-modules"`

	// Optional: Include is used to specify the related resources to include in the response.
	Include []RegistryNoCodeModuleIncludeOpt `url:"include,omitempty"`
}

// RegistryNoCodeModuleReadVariablesOptions is used when reading the variables
// for a no-code module.
type RegistryNoCodeModuleReadVariablesOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-updating
	Type string `jsonapi:"primary,no-code-modules"`
}

// RegistryNoCodeModuleUpdateOptions is used when updating a registry no-code module
type RegistryNoCodeModuleUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-updating
	Type string `jsonapi:"primary,no-code-modules"`

	// Required: the registry module to use for the no-code module (only the ID is used)
	RegistryModule *RegistryModule `jsonapi:"relation,registry-module"`

	// Optional: the version pin for the module. valid values are "latest" or a semver string
	VersionPin string `jsonapi:"attr,version-pin,omitempty"`

	// Optional: whether no-code is enabled for the module
	Enabled *bool `jsonapi:"attr,enabled,omitempty"`

	// Optional: are the variable options for the module
	VariableOptions []*NoCodeVariableOption `jsonapi:"relation,variable-options,omitempty"`
}

// WorkspaceUpgrade contains the data returned by the no-code workspace upgrade
// API endpoint.
type WorkspaceUpgrade struct {
	// Status is the status of the run of the upgrade
	Status string `jsonapi:"attr,status"`

	// PlanURL is the URL to the plan of the upgrade
	PlanURL string `jsonapi:"attr,plan-url"`

	// Message is the message returned by the API when an upgrade is not available.
	Message string `jsonapi:"attr,message"`
}

// Create a new registry no-code module
func (r *registryNoCodeModules) Create(ctx context.Context, organization string, options RegistryNoCodeModuleCreateOptions) (*RegistryNoCodeModule, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/no-code-modules", url.PathEscape(organization))
	req, err := r.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryNoCodeModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// Read a registry no-code module
func (r *registryNoCodeModules) Read(ctx context.Context, noCodeModuleID string, options *RegistryNoCodeModuleReadOptions) (*RegistryNoCodeModule, error) {
	if !validStringID(&noCodeModuleID) {
		return nil, ErrInvalidModuleID
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("no-code-modules/%s", url.PathEscape(noCodeModuleID))
	req, err := r.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryNoCodeModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// ReadVariables retrieves the no-code variable options for a version of a
// module.
func (r *registryNoCodeModules) ReadVariables(
	ctx context.Context,
	noCodeModuleID, noCodeModuleVersion string,
	options *RegistryNoCodeModuleReadVariablesOptions,
) (*RegistryModuleVariableList, error) {
	if !validStringID(&noCodeModuleID) {
		return nil, ErrInvalidModuleID
	}
	if !validVersion(noCodeModuleVersion) {
		return nil, ErrInvalidVersion
	}

	u := fmt.Sprintf(
		"no-code-modules/%s/versions/%s/module-variables",
		url.PathEscape(noCodeModuleID),
		url.PathEscape(noCodeModuleVersion),
	)
	req, err := r.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	resp := &RegistryModuleVariableList{}
	err = req.Do(ctx, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update a registry no-code module
func (r *registryNoCodeModules) Update(ctx context.Context, noCodeModuleID string, options RegistryNoCodeModuleUpdateOptions) (*RegistryNoCodeModule, error) {
	if !validString(&noCodeModuleID) {
		return nil, ErrInvalidModuleID
	}
	if !validStringID(&noCodeModuleID) {
		return nil, ErrInvalidModuleID
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("no-code-modules/%s", url.PathEscape(noCodeModuleID))
	req, err := r.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryNoCodeModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// Delete is used to delete the registry no-code module
func (r *registryNoCodeModules) Delete(ctx context.Context, noCodeModuleID string) error {
	if !validStringID(&noCodeModuleID) {
		return ErrInvalidModuleID
	}

	u := fmt.Sprintf("no-code-modules/%s", url.PathEscape(noCodeModuleID))
	req, err := r.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// CreateWorkspace creates a no-code workspace using a no-code module.
func (r *registryNoCodeModules) CreateWorkspace(
	ctx context.Context,
	noCodeModuleID string,
	options *RegistryNoCodeModuleCreateWorkspaceOptions,
) (*Workspace, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("no-code-modules/%s/workspaces", url.PathEscape(noCodeModuleID))
	req, err := r.client.NewRequest("POST", u, options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// UpgradeWorkspace initiates an upgrade of an existing no-code module workspace.
func (r *registryNoCodeModules) UpgradeWorkspace(
	ctx context.Context,
	noCodeModuleID string,
	workspaceID string,
	options *RegistryNoCodeModuleUpgradeWorkspaceOptions,
) (*WorkspaceUpgrade, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("no-code-modules/%s/workspaces/%s/upgrade",
		url.PathEscape(noCodeModuleID),
		workspaceID,
	)
	req, err := r.client.NewRequest("POST", u, options)
	if err != nil {
		return nil, err
	}

	wu := &WorkspaceUpgrade{}
	err = req.Do(ctx, wu)
	if err != nil {
		return nil, err
	}

	return wu, nil
}

func (o RegistryNoCodeModuleCreateOptions) valid() error {
	if o.RegistryModule == nil || o.RegistryModule.ID == "" {
		return ErrRequiredRegistryModule
	}

	return nil
}

func (o *RegistryNoCodeModuleUpdateOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if o.RegistryModule == nil || o.RegistryModule.ID == "" {
		return ErrRequiredRegistryModule
	}

	return nil
}

func (o *RegistryNoCodeModuleReadOptions) valid() error {
	return nil
}

func (o *RegistryNoCodeModuleCreateWorkspaceOptions) valid() error {
	if !validString(&o.Name) {
		return ErrRequiredName
	}

	return nil
}

func (o *RegistryNoCodeModuleUpgradeWorkspaceOptions) valid() error {
	return nil
}
