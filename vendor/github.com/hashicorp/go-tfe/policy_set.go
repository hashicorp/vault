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
var _ PolicySets = (*policySets)(nil)

// PolicyKind is an indicator of the underlying technology that the policy or policy set supports.
// There are two kinds documented in the enum.
type PolicyKind string

const (
	OPA      PolicyKind = "opa"
	Sentinel PolicyKind = "sentinel"
)

// PolicySets describes all the policy set related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-sets
type PolicySets interface {
	// List all the policy sets for a given organization.
	List(ctx context.Context, organization string, options *PolicySetListOptions) (*PolicySetList, error)

	// Create a policy set and associate it with an organization.
	Create(ctx context.Context, organization string, options PolicySetCreateOptions) (*PolicySet, error)

	// Read a policy set by its ID.
	Read(ctx context.Context, policySetID string) (*PolicySet, error)

	// ReadWithOptions reads a policy set by its ID using the options supplied.
	ReadWithOptions(ctx context.Context, policySetID string, options *PolicySetReadOptions) (*PolicySet, error)

	// Update an existing policy set.
	Update(ctx context.Context, policySetID string, options PolicySetUpdateOptions) (*PolicySet, error)

	// Add policies to a policy set. This function can only be used when
	// there is no VCS repository associated with the policy set.
	AddPolicies(ctx context.Context, policySetID string, options PolicySetAddPoliciesOptions) error

	// Remove policies from a policy set. This function can only be used
	// when there is no VCS repository associated with the policy set.
	RemovePolicies(ctx context.Context, policySetID string, options PolicySetRemovePoliciesOptions) error

	// Add workspaces to a policy set.
	AddWorkspaces(ctx context.Context, policySetID string, options PolicySetAddWorkspacesOptions) error

	// Remove workspaces from a policy set.
	RemoveWorkspaces(ctx context.Context, policySetID string, options PolicySetRemoveWorkspacesOptions) error

	// Add workspace exclusions to a policy set.
	AddWorkspaceExclusions(ctx context.Context, policySetID string, options PolicySetAddWorkspaceExclusionsOptions) error

	// Remove workspace exclusions from a policy set.
	RemoveWorkspaceExclusions(ctx context.Context, policySetID string, options PolicySetRemoveWorkspaceExclusionsOptions) error

	// Add projects to a policy set.
	AddProjects(ctx context.Context, policySetID string, options PolicySetAddProjectsOptions) error

	// Remove projects from a policy set.
	RemoveProjects(ctx context.Context, policySetID string, options PolicySetRemoveProjectsOptions) error

	// Delete a policy set by its ID.
	Delete(ctx context.Context, policyID string) error
}

// policySets implements PolicySets.
type policySets struct {
	client *Client
}

// PolicySetList represents a list of policy sets.
type PolicySetList struct {
	*Pagination
	Items []*PolicySet
}

// PolicySet represents a Terraform Enterprise policy set.
type PolicySet struct {
	ID           string     `jsonapi:"primary,policy-sets"`
	Name         string     `jsonapi:"attr,name"`
	Description  string     `jsonapi:"attr,description"`
	Kind         PolicyKind `jsonapi:"attr,kind"`
	Overridable  *bool      `jsonapi:"attr,overridable"`
	Global       bool       `jsonapi:"attr,global"`
	PoliciesPath string     `jsonapi:"attr,policies-path"`
	// **Note: This field is still in BETA and subject to change.**
	PolicyCount       int       `jsonapi:"attr,policy-count"`
	VCSRepo           *VCSRepo  `jsonapi:"attr,vcs-repo"`
	WorkspaceCount    int       `jsonapi:"attr,workspace-count"`
	ProjectCount      int       `jsonapi:"attr,project-count"`
	CreatedAt         time.Time `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt         time.Time `jsonapi:"attr,updated-at,iso8601"`
	AgentEnabled      bool      `jsonapi:"attr,agent-enabled"`
	PolicyToolVersion string    `jsonapi:"attr,policy-tool-version"`

	// Relations
	// The organization to which the policy set belongs to.
	Organization *Organization `jsonapi:"relation,organization"`
	// The workspaces to which the policy set applies.
	Workspaces []*Workspace `jsonapi:"relation,workspaces"`
	// Individually managed policies which are associated with the policy set.
	Policies []*Policy `jsonapi:"relation,policies"`
	// The most recently created policy set version, regardless of status.
	// Note that this relationship may include an errored and unusable version,
	// and is intended to allow checking for errors.
	NewestVersion *PolicySetVersion `jsonapi:"relation,newest-version"`
	// The most recent successful policy set version.
	CurrentVersion *PolicySetVersion `jsonapi:"relation,current-version"`
	// The workspace exclusions to which the policy set applies.
	WorkspaceExclusions []*Workspace `jsonapi:"relation,workspace-exclusions"`
	// The projects to which the policy set applies.
	Projects []*Project `jsonapi:"relation,projects"`
}

// PolicySetIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-sets#available-related-resources
type PolicySetIncludeOpt string

const (
	PolicySetPolicies            PolicySetIncludeOpt = "policies"
	PolicySetWorkspaces          PolicySetIncludeOpt = "workspaces"
	PolicySetCurrentVersion      PolicySetIncludeOpt = "current_version"
	PolicySetNewestVersion       PolicySetIncludeOpt = "newest_version"
	PolicySetProjects            PolicySetIncludeOpt = "projects"
	PolicySetWorkspaceExclusions PolicySetIncludeOpt = "workspace_exclusions"
)

// PolicySetListOptions represents the options for listing policy sets.
type PolicySetListOptions struct {
	ListOptions

	// Optional: A search string (partial policy set name) used to filter the results.
	Search string `url:"search[name],omitempty"`

	// Optional: A kind string used to filter the results by the policy set kind.
	Kind PolicyKind `url:"filter[kind],omitempty"`

	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-sets#available-related-resources
	Include []PolicySetIncludeOpt `url:"include,omitempty"`
}

// PolicySetReadOptions are read options.
// For a full list of relations, please see:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-sets#relationships
type PolicySetReadOptions struct {
	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-sets#available-related-resources
	Include []PolicySetIncludeOpt `url:"include,omitempty"`
}

// PolicySetCreateOptions represents the options for creating a new policy set.
type PolicySetCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,policy-sets"`

	// Required: The name of the policy set.
	Name *string `jsonapi:"attr,name"`

	// Optional: The description of the policy set.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Optional: Whether or not the policy set is global.
	Global *bool `jsonapi:"attr,global,omitempty"`

	// Optional: The underlying technology that the policy set supports
	Kind PolicyKind `jsonapi:"attr,kind,omitempty"`

	// Optional: Whether or not users can override this policy when it fails during a run. Only valid for policy evaluations.
	// https://developer.hashicorp.com/terraform/cloud-docs/policy-enforcement/manage-policy-sets#policy-checks-versus-policy-evaluations
	Overridable *bool `jsonapi:"attr,overridable,omitempty"`

	// Optional: Whether or not the policy is run as an evaluation inside the agent.
	AgentEnabled *bool `jsonapi:"attr,agent-enabled,omitempty"`

	// Optional: The policy tool version to run the evaluation against.
	PolicyToolVersion *string `jsonapi:"attr,policy-tool-version,omitempty"`

	// Optional: The sub-path within the attached VCS repository to ingress. All
	// files and directories outside of this sub-path will be ignored.
	// This option may only be specified when a VCS repo is present.
	PoliciesPath *string `jsonapi:"attr,policies-path,omitempty"`

	// Optional: The initial members of the policy set.
	Policies []*Policy `jsonapi:"relation,policies,omitempty"`

	// Optional: VCS repository information. When present, the policies and
	// configuration will be sourced from the specified VCS repository
	// instead of being defined within the policy set itself. Note that
	// this option is mutually exclusive with the Policies option and
	// both cannot be used at the same time.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`

	// Optional: The initial list of workspaces for which the policy set should be enforced.
	Workspaces []*Workspace `jsonapi:"relation,workspaces,omitempty"`

	// Optional: The initial list of workspace exclusions for which the policy set should be enforced.
	WorkspaceExclusions []*Workspace `jsonapi:"relation,workspace-exclusions,omitempty"`

	// Optional: The initial list of projects for which the policy set should be enforced.
	Projects []*Project `jsonapi:"relation,projects,omitempty"`
}

// PolicySetUpdateOptions represents the options for updating a policy set.
type PolicySetUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,policy-sets"`

	// Optional: The name of the policy set.
	Name *string `jsonapi:"attr,name,omitempty"`

	// Optional: The description of the policy set.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Optional: Whether or not the policy set is global.
	Global *bool `jsonapi:"attr,global,omitempty"`

	// Optional: Whether or not users can override this policy when it fails during a run. Only valid for policy evaluations.
	// https://developer.hashicorp.com/terraform/cloud-docs/policy-enforcement/manage-policy-sets#policy-checks-versus-policy-evaluations
	Overridable *bool `jsonapi:"attr,overridable,omitempty"`

	// Optional: Whether or not the policy is run as an evaluation inside the agent.
	AgentEnabled *bool `jsonapi:"attr,agent-enabled,omitempty"`

	// Optional: The policy tool version to run the evaluation against.
	PolicyToolVersion *string `jsonapi:"attr,policy-tool-version,omitempty"`

	// Optional: The sub-path within the attached VCS repository to ingress. All
	// files and directories outside of this sub-path will be ignored.
	// This option may only be specified when a VCS repo is present.
	PoliciesPath *string `jsonapi:"attr,policies-path,omitempty"`

	// Optional: VCS repository information. When present, the policies and
	// configuration will be sourced from the specified VCS repository
	// instead of being defined within the policy set itself. Note that
	// specifying this option may only be used on policy sets with no
	// directly-attached policies (*PolicySet.Policies). Specifying this
	// option when policies are already present will result in an error.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`
}

// PolicySetAddPoliciesOptions represents the options for adding policies
// to a policy set.
type PolicySetAddPoliciesOptions struct {
	// The policies to add to the policy set.
	Policies []*Policy
}

// PolicySetRemovePoliciesOptions represents the options for removing
// policies from a policy set.
type PolicySetRemovePoliciesOptions struct {
	// The policies to remove from the policy set.
	Policies []*Policy
}

// PolicySetAddWorkspacesOptions represents the options for adding workspaces
// to a policy set.
type PolicySetAddWorkspacesOptions struct {
	// The workspaces to add to the policy set.
	Workspaces []*Workspace
}

// PolicySetRemoveWorkspacesOptions represents the options for removing
// workspaces from a policy set.
type PolicySetRemoveWorkspacesOptions struct {
	// The workspaces to remove from the policy set.
	Workspaces []*Workspace
}

// PolicySetAddWorkspaceExclusionsOptions represents the options for adding workspace exclusions to a policy set.
type PolicySetAddWorkspaceExclusionsOptions struct {
	// The workspaces to add to the policy set exclusion list.
	WorkspaceExclusions []*Workspace
}

// PolicySetRemoveWorkspaceExclusionsOptions represents the options for removing workspace exclusions from a policy set.
type PolicySetRemoveWorkspaceExclusionsOptions struct {
	// The workspaces to remove from the policy set exclusion list.
	WorkspaceExclusions []*Workspace
}

// PolicySetAddProjectsOptions represents the options for adding projects
// to a policy set.
type PolicySetAddProjectsOptions struct {
	// The projects to add to the policy set.
	Projects []*Project
}

// PolicySetRemoveProjectsOptions represents the options for removing
// projects from a policy set.
type PolicySetRemoveProjectsOptions struct {
	// The projects to remove from the policy set.
	Projects []*Project
}

// List all the policies for a given organization.
func (s *policySets) List(ctx context.Context, organization string, options *PolicySetListOptions) (*PolicySetList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/policy-sets", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	psl := &PolicySetList{}
	err = req.Do(ctx, psl)
	if err != nil {
		return nil, err
	}

	return psl, nil
}

// Create a policy set and associate it with an organization.
func (s *policySets) Create(ctx context.Context, organization string, options PolicySetCreateOptions) (*PolicySet, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/policy-sets", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	ps := &PolicySet{}
	err = req.Do(ctx, ps)
	if err != nil {
		return nil, err
	}

	return ps, err
}

// Read a policy set by its ID.
func (s *policySets) Read(ctx context.Context, policySetID string) (*PolicySet, error) {
	return s.ReadWithOptions(ctx, policySetID, nil)
}

// ReadWithOptions reads a policy by its ID using the options supplied.
func (s *policySets) ReadWithOptions(ctx context.Context, policySetID string, options *PolicySetReadOptions) (*PolicySet, error) {
	if !validStringID(&policySetID) {
		return nil, ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("policy-sets/%s", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	ps := &PolicySet{}
	err = req.Do(ctx, ps)
	if err != nil {
		return nil, err
	}

	return ps, err
}

// Update an existing policy set.
func (s *policySets) Update(ctx context.Context, policySetID string, options PolicySetUpdateOptions) (*PolicySet, error) {
	if !validStringID(&policySetID) {
		return nil, ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("policy-sets/%s", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	ps := &PolicySet{}
	err = req.Do(ctx, ps)
	if err != nil {
		return nil, err
	}

	return ps, err
}

// AddPolicies adds policies to a policy set
func (s *policySets) AddPolicies(ctx context.Context, policySetID string, options PolicySetAddPoliciesOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/policies", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("POST", u, options.Policies)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemovePolicies remove policies from a policy set
func (s *policySets) RemovePolicies(ctx context.Context, policySetID string, options PolicySetRemovePoliciesOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/policies", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("DELETE", u, options.Policies)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Addworkspaces adds workspaces to a policy set.
func (s *policySets) AddWorkspaces(ctx context.Context, policySetID string, options PolicySetAddWorkspacesOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/workspaces", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("POST", u, options.Workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveWorkspaces removes workspaces from a policy set.
func (s *policySets) RemoveWorkspaces(ctx context.Context, policySetID string, options PolicySetRemoveWorkspacesOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/workspaces", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("DELETE", u, options.Workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// AddWorkspaceExclusions adds workspace exclusions to a policy set.
func (s *policySets) AddWorkspaceExclusions(ctx context.Context, policySetID string, options PolicySetAddWorkspaceExclusionsOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/workspace-exclusions", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("POST", u, options.WorkspaceExclusions)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveWorkspaceExclusions removes workspace exclusions from a policy set.
func (s *policySets) RemoveWorkspaceExclusions(ctx context.Context, policySetID string, options PolicySetRemoveWorkspaceExclusionsOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/workspace-exclusions", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("DELETE", u, options.WorkspaceExclusions)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// AddProjects adds projects to a given policy set.
func (s *policySets) AddProjects(ctx context.Context, policySetID string, options PolicySetAddProjectsOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/projects", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("POST", u, options.Projects)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveProjects removes projects from a policy set.
func (s *policySets) RemoveProjects(ctx context.Context, policySetID string, options PolicySetRemoveProjectsOptions) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("policy-sets/%s/relationships/projects", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("DELETE", u, options.Projects)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Delete a policy set by its ID.
func (s *policySets) Delete(ctx context.Context, policySetID string) error {
	if !validStringID(&policySetID) {
		return ErrInvalidPolicySetID
	}

	u := fmt.Sprintf("policy-sets/%s", url.PathEscape(policySetID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o PolicySetCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validStringID(o.Name) {
		return ErrInvalidName
	}
	return nil
}

func (o PolicySetRemoveWorkspacesOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

func (o PolicySetRemoveWorkspaceExclusionsOptions) valid() error {
	if o.WorkspaceExclusions == nil {
		return ErrWorkspacesRequired
	}
	if len(o.WorkspaceExclusions) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

func (o PolicySetRemoveProjectsOptions) valid() error {
	if o.Projects == nil {
		return ErrRequiredProject
	}
	if len(o.Projects) == 0 {
		return ErrProjectMinLimit
	}
	return nil
}

func (o PolicySetUpdateOptions) valid() error {
	if o.Name != nil && !validStringID(o.Name) {
		return ErrInvalidName
	}
	return nil
}

func (o PolicySetAddPoliciesOptions) valid() error {
	if o.Policies == nil {
		return ErrRequiredPolicies
	}
	if len(o.Policies) == 0 {
		return ErrInvalidPolicies
	}
	return nil
}

func (o PolicySetRemovePoliciesOptions) valid() error {
	if o.Policies == nil {
		return ErrRequiredPolicies
	}
	if len(o.Policies) == 0 {
		return ErrInvalidPolicies
	}
	return nil
}

func (o PolicySetAddWorkspacesOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

func (o PolicySetAddWorkspaceExclusionsOptions) valid() error {
	if o.WorkspaceExclusions == nil {
		return ErrWorkspacesRequired
	}
	if len(o.WorkspaceExclusions) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

func (o PolicySetAddProjectsOptions) valid() error {
	if o.Projects == nil {
		return ErrRequiredProject
	}
	if len(o.Projects) == 0 {
		return ErrProjectMinLimit
	}
	return nil
}

func (o *PolicySetReadOptions) valid() error {
	return nil
}
