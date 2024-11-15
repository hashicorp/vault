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
var _ OAuthClients = (*oAuthClients)(nil)

// OAuthClients describes all the OAuth client related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/oauth-clients
type OAuthClients interface {
	// List all the OAuth clients for a given organization.
	List(ctx context.Context, organization string, options *OAuthClientListOptions) (*OAuthClientList, error)

	// Create an OAuth client to connect an organization and a VCS provider.
	Create(ctx context.Context, organization string, options OAuthClientCreateOptions) (*OAuthClient, error)

	// Read an OAuth client by its ID.
	Read(ctx context.Context, oAuthClientID string) (*OAuthClient, error)

	// ReadWithOptions reads an oauth client by its ID using the options supplied.
	ReadWithOptions(ctx context.Context, oAuthClientID string, options *OAuthClientReadOptions) (*OAuthClient, error)

	// Update an existing OAuth client by its ID.
	Update(ctx context.Context, oAuthClientID string, options OAuthClientUpdateOptions) (*OAuthClient, error)

	// Delete an OAuth client by its ID.
	Delete(ctx context.Context, oAuthClientID string) error

	// AddProjects add projects to an oauth client.
	AddProjects(ctx context.Context, oAuthClientID string, options OAuthClientAddProjectsOptions) error

	// RemoveProjects remove projects from an oauth client.
	RemoveProjects(ctx context.Context, oAuthClientID string, options OAuthClientRemoveProjectsOptions) error
}

// oAuthClients implements OAuthClients.
type oAuthClients struct {
	client *Client
}

// ServiceProviderType represents a VCS type.
type ServiceProviderType string

// List of available VCS types.
const (
	ServiceProviderAzureDevOpsServer   ServiceProviderType = "ado_server"
	ServiceProviderAzureDevOpsServices ServiceProviderType = "ado_services"
	ServiceProviderBitbucketDataCenter ServiceProviderType = "bitbucket_data_center"
	ServiceProviderBitbucket           ServiceProviderType = "bitbucket_hosted"
	// Bitbucket Server v5.4.0 and above
	ServiceProviderBitbucketServer ServiceProviderType = "bitbucket_server"
	// Bitbucket Server v5.3.0 and below
	ServiceProviderBitbucketServerLegacy ServiceProviderType = "bitbucket_server_legacy"
	ServiceProviderGithub                ServiceProviderType = "github"
	ServiceProviderGithubEE              ServiceProviderType = "github_enterprise"
	ServiceProviderGitlab                ServiceProviderType = "gitlab_hosted"
	ServiceProviderGitlabCE              ServiceProviderType = "gitlab_community_edition"
	ServiceProviderGitlabEE              ServiceProviderType = "gitlab_enterprise_edition"
)

// OAuthClientList represents a list of OAuth clients.
type OAuthClientList struct {
	*Pagination
	Items []*OAuthClient
}

// OAuthClient represents a connection between an organization and a VCS
// provider.
type OAuthClient struct {
	ID                  string              `jsonapi:"primary,oauth-clients"`
	APIURL              string              `jsonapi:"attr,api-url"`
	CallbackURL         string              `jsonapi:"attr,callback-url"`
	ConnectPath         string              `jsonapi:"attr,connect-path"`
	CreatedAt           time.Time           `jsonapi:"attr,created-at,iso8601"`
	HTTPURL             string              `jsonapi:"attr,http-url"`
	Key                 string              `jsonapi:"attr,key"`
	RSAPublicKey        string              `jsonapi:"attr,rsa-public-key"`
	Name                *string             `jsonapi:"attr,name"`
	Secret              string              `jsonapi:"attr,secret"`
	ServiceProvider     ServiceProviderType `jsonapi:"attr,service-provider"`
	ServiceProviderName string              `jsonapi:"attr,service-provider-display-name"`
	OrganizationScoped  *bool               `jsonapi:"attr,organization-scoped"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
	OAuthTokens  []*OAuthToken `jsonapi:"relation,oauth-tokens"`
	AgentPool    *AgentPool    `jsonapi:"relation,agent-pool"`
	// The projects to which the oauth client applies.
	Projects []*Project `jsonapi:"relation,projects"`
}

// A list of relations to include
type OAuthClientIncludeOpt string

const (
	OauthClientOauthTokens OAuthClientIncludeOpt = "oauth_tokens"
	OauthClientProjects    OAuthClientIncludeOpt = "projects"
)

// OAuthClientListOptions represents the options for listing
// OAuth clients.
type OAuthClientListOptions struct {
	ListOptions

	Include []OAuthClientIncludeOpt `url:"include,omitempty"`
}

// OAuthClientReadOptions are read options.
// For a full list of relations, please see:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/oauth-clients#relationships
type OAuthClientReadOptions struct {
	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/oauth-clients#available-related-resources
	Include []OAuthClientIncludeOpt `url:"include,omitempty"`
}

// OAuthClientCreateOptions represents the options for creating an OAuth client.
type OAuthClientCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,oauth-clients"`

	// A display name for the OAuth Client.
	Name *string `jsonapi:"attr,name"`

	// Required: The base URL of your VCS provider's API.
	APIURL *string `jsonapi:"attr,api-url"`

	// Required: The homepage of your VCS provider.
	HTTPURL *string `jsonapi:"attr,http-url"`

	// Optional: The OAuth Client key.
	Key *string `jsonapi:"attr,key,omitempty"`

	// Optional: The token string you were given by your VCS provider.
	OAuthToken *string `jsonapi:"attr,oauth-token-string,omitempty"`

	// Optional: The initial list of projects for which the oauth client should be associated with.
	Projects []*Project `jsonapi:"relation,projects,omitempty"`

	// Optional: Private key associated with this vcs provider - only available for ado_server
	PrivateKey *string `jsonapi:"attr,private-key,omitempty"`

	// Optional: Secret key associated with this vcs provider - only available for ado_server
	Secret *string `jsonapi:"attr,secret,omitempty"`

	// Optional: RSAPublicKey the text of the SSH public key associated with your
	// BitBucket Data Center Application Link.
	RSAPublicKey *string `jsonapi:"attr,rsa-public-key,omitempty"`

	// Required: The VCS provider being connected with.
	ServiceProvider *ServiceProviderType `jsonapi:"attr,service-provider"`

	// Optional: AgentPool to associate the VCS Provider with, for PrivateVCS support
	AgentPool *AgentPool `jsonapi:"relation,agent-pool,omitempty"`

	// Optional: Whether the OAuthClient is available to all workspaces in the organization.
	// True if the oauth client is organization scoped, false otherwise.
	OrganizationScoped *bool `jsonapi:"attr,organization-scoped,omitempty"`
}

// OAuthClientUpdateOptions represents the options for updating an OAuth client.
type OAuthClientUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,oauth-clients"`

	// Optional: A display name for the OAuth Client.
	Name *string `jsonapi:"attr,name,omitempty"`

	// Optional: The OAuth Client key.
	Key *string `jsonapi:"attr,key,omitempty"`

	// Optional: Secret key associated with this vcs provider - only available for ado_server
	Secret *string `jsonapi:"attr,secret,omitempty"`

	// Optional: RSAPublicKey the text of the SSH public key associated with your BitBucket
	// Server Application Link.
	RSAPublicKey *string `jsonapi:"attr,rsa-public-key,omitempty"`

	// Optional: The token string you were given by your VCS provider.
	OAuthToken *string `jsonapi:"attr,oauth-token-string,omitempty"`

	// Optional: Whether the OAuthClient is available to all workspaces in the organization.
	// True if the oauth client is organization scoped, false otherwise.
	OrganizationScoped *bool `jsonapi:"attr,organization-scoped,omitempty"`
}

// OAuthClientAddProjectsOptions represents the options for adding projects
// to an oauth client.
type OAuthClientAddProjectsOptions struct {
	// The projects to add to an oauth client.
	Projects []*Project
}

// OAuthClientRemoveProjectsOptions represents the options for removing
// projects from an oauth client.
type OAuthClientRemoveProjectsOptions struct {
	// The projects to remove from an oauth client.
	Projects []*Project
}

// List all the OAuth clients for a given organization.
func (s *oAuthClients) List(ctx context.Context, organization string, options *OAuthClientListOptions) (*OAuthClientList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/oauth-clients", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	ocl := &OAuthClientList{}
	err = req.Do(ctx, ocl)
	if err != nil {
		return nil, err
	}

	return ocl, nil
}

// Create an OAuth client to connect an organization and a VCS provider.
func (s *oAuthClients) Create(ctx context.Context, organization string, options OAuthClientCreateOptions) (*OAuthClient, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/oauth-clients", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	oc := &OAuthClient{}
	err = req.Do(ctx, oc)
	if err != nil {
		return nil, err
	}

	return oc, nil
}

// Read an OAuth client by its ID.
func (s *oAuthClients) Read(ctx context.Context, oAuthClientID string) (*OAuthClient, error) {
	return s.ReadWithOptions(ctx, oAuthClientID, nil)
}

func (s *oAuthClients) ReadWithOptions(ctx context.Context, oAuthClientID string, options *OAuthClientReadOptions) (*OAuthClient, error) {
	if !validStringID(&oAuthClientID) {
		return nil, ErrInvalidOauthClientID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("oauth-clients/%s", url.PathEscape(oAuthClientID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	oc := &OAuthClient{}
	err = req.Do(ctx, oc)
	if err != nil {
		return nil, err
	}

	return oc, err
}

// Update an OAuth client by its ID.
func (s *oAuthClients) Update(ctx context.Context, oAuthClientID string, options OAuthClientUpdateOptions) (*OAuthClient, error) {
	if !validStringID(&oAuthClientID) {
		return nil, ErrInvalidOauthClientID
	}

	u := fmt.Sprintf("oauth-clients/%s", url.PathEscape(oAuthClientID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	oc := &OAuthClient{}
	err = req.Do(ctx, oc)
	if err != nil {
		return nil, err
	}

	return oc, err
}

// Delete an OAuth client by its ID.
func (s *oAuthClients) Delete(ctx context.Context, oAuthClientID string) error {
	if !validStringID(&oAuthClientID) {
		return ErrInvalidOauthClientID
	}

	u := fmt.Sprintf("oauth-clients/%s", url.PathEscape(oAuthClientID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o OAuthClientCreateOptions) valid() error {
	if !validString(o.APIURL) {
		return ErrRequiredAPIURL
	}
	if !validString(o.HTTPURL) {
		return ErrRequiredHTTPURL
	}
	if o.ServiceProvider == nil {
		return ErrRequiredServiceProvider
	}
	if !validString(o.OAuthToken) &&
		*o.ServiceProvider != *ServiceProvider(ServiceProviderBitbucketServer) &&
		*o.ServiceProvider != *ServiceProvider(ServiceProviderBitbucketDataCenter) {
		return ErrRequiredOauthToken
	}
	if validString(o.PrivateKey) && *o.ServiceProvider != *ServiceProvider(ServiceProviderAzureDevOpsServer) {
		return ErrUnsupportedPrivateKey
	}
	return nil
}

func (o *OAuthClientListOptions) valid() error {
	return nil
}

// AddProjects adds projects to a given oauth client.
func (s *oAuthClients) AddProjects(ctx context.Context, oAuthClientID string, options OAuthClientAddProjectsOptions) error {
	if !validStringID(&oAuthClientID) {
		return ErrInvalidOauthClientID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("oauth-clients/%s/relationships/projects", url.PathEscape(oAuthClientID))
	req, err := s.client.NewRequest("POST", u, options.Projects)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveProjects removes projects from an oauth client.
func (s *oAuthClients) RemoveProjects(ctx context.Context, oAuthClientID string, options OAuthClientRemoveProjectsOptions) error {
	if !validStringID(&oAuthClientID) {
		return ErrInvalidOauthClientID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("oauth-clients/%s/relationships/projects", url.PathEscape(oAuthClientID))
	req, err := s.client.NewRequest("DELETE", u, options.Projects)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o OAuthClientAddProjectsOptions) valid() error {
	if o.Projects == nil {
		return ErrRequiredProject
	}
	if len(o.Projects) == 0 {
		return ErrProjectMinLimit
	}
	return nil
}

func (o OAuthClientRemoveProjectsOptions) valid() error {
	if o.Projects == nil {
		return ErrRequiredProject
	}
	if len(o.Projects) == 0 {
		return ErrProjectMinLimit
	}
	return nil
}

func (o *OAuthClientReadOptions) valid() error {
	return nil
}
