package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Organizations = (*organizations)(nil)

// Organizations describes all the organization related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/cloud/api/organizations.html
type Organizations interface {
	// List all the organizations visible to the current user.
	List(ctx context.Context, options OrganizationListOptions) (*OrganizationList, error)

	// Create a new organization with the given options.
	Create(ctx context.Context, options OrganizationCreateOptions) (*Organization, error)

	// Read an organization by its name.
	Read(ctx context.Context, organization string) (*Organization, error)

	// Update attributes of an existing organization.
	Update(ctx context.Context, organization string, options OrganizationUpdateOptions) (*Organization, error)

	// Delete an organization by its name.
	Delete(ctx context.Context, organization string) error

	// Capacity shows the current run capacity of an organization.
	Capacity(ctx context.Context, organization string) (*Capacity, error)

	// Entitlements shows the entitlements of an organization.
	Entitlements(ctx context.Context, organization string) (*Entitlements, error)

	// RunQueue shows the current run queue of an organization.
	RunQueue(ctx context.Context, organization string, options RunQueueOptions) (*RunQueue, error)
}

// organizations implements Organizations.
type organizations struct {
	client *Client
}

// AuthPolicyType represents an authentication policy type.
type AuthPolicyType string

// List of available authentication policies.
const (
	AuthPolicyPassword  AuthPolicyType = "password"
	AuthPolicyTwoFactor AuthPolicyType = "two_factor_mandatory"
)

// OrganizationList represents a list of organizations.
type OrganizationList struct {
	*Pagination
	Items []*Organization
}

// Organization represents a Terraform Enterprise organization.
type Organization struct {
	Name                                              string                   `jsonapi:"primary,organizations"`
	CollaboratorAuthPolicy                            AuthPolicyType           `jsonapi:"attr,collaborator-auth-policy"`
	CostEstimationEnabled                             bool                     `jsonapi:"attr,cost-estimation-enabled"`
	CreatedAt                                         time.Time                `jsonapi:"attr,created-at,iso8601"`
	Email                                             string                   `jsonapi:"attr,email"`
	ExternalID                                        string                   `jsonapi:"attr,external-id"`
	OwnersTeamSAMLRoleID                              string                   `jsonapi:"attr,owners-team-saml-role-id"`
	Permissions                                       *OrganizationPermissions `jsonapi:"attr,permissions"`
	SAMLEnabled                                       bool                     `jsonapi:"attr,saml-enabled"`
	SessionRemember                                   int                      `jsonapi:"attr,session-remember"`
	SessionTimeout                                    int                      `jsonapi:"attr,session-timeout"`
	TrialExpiresAt                                    time.Time                `jsonapi:"attr,trial-expires-at,iso8601"`
	TwoFactorConformant                               bool                     `jsonapi:"attr,two-factor-conformant"`
	SendPassingStatusesForUntriggeredSpeculativePlans bool                     `jsonapi:"attr,send-passing-statuses-for-untriggered-speculative-plans"`
}

// Capacity represents the current run capacity of an organization.
type Capacity struct {
	Organization string `jsonapi:"primary,organization-capacity"`
	Pending      int    `jsonapi:"attr,pending"`
	Running      int    `jsonapi:"attr,running"`
}

// Entitlements represents the entitlements of an organization.
type Entitlements struct {
	ID                    string `jsonapi:"primary,entitlement-sets"`
	Agents                bool   `jsonapi:"attr,agents"`
	AuditLogging          bool   `jsonapi:"attr,audit-logging"`
	CostEstimation        bool   `jsonapi:"attr,cost-estimation"`
	Operations            bool   `jsonapi:"attr,operations"`
	PrivateModuleRegistry bool   `jsonapi:"attr,private-module-registry"`
	SSO                   bool   `jsonapi:"attr,sso"`
	Sentinel              bool   `jsonapi:"attr,sentinel"`
	StateStorage          bool   `jsonapi:"attr,state-storage"`
	Teams                 bool   `jsonapi:"attr,teams"`
	VCSIntegrations       bool   `jsonapi:"attr,vcs-integrations"`
}

// RunQueue represents the current run queue of an organization.
type RunQueue struct {
	*Pagination
	Items []*Run
}

// OrganizationPermissions represents the organization permissions.
type OrganizationPermissions struct {
	CanCreateTeam               bool `jsonapi:"attr,can-create-team"`
	CanCreateWorkspace          bool `jsonapi:"attr,can-create-workspace"`
	CanCreateWorkspaceMigration bool `jsonapi:"attr,can-create-workspace-migration"`
	CanDestroy                  bool `jsonapi:"attr,can-destroy"`
	CanTraverse                 bool `jsonapi:"attr,can-traverse"`
	CanUpdate                   bool `jsonapi:"attr,can-update"`
	CanUpdateAPIToken           bool `jsonapi:"attr,can-update-api-token"`
	CanUpdateOAuth              bool `jsonapi:"attr,can-update-oauth"`
	CanUpdateSentinel           bool `jsonapi:"attr,can-update-sentinel"`
}

// OrganizationListOptions represents the options for listing organizations.
type OrganizationListOptions struct {
	ListOptions
}

// List all the organizations visible to the current user.
func (s *organizations) List(ctx context.Context, options OrganizationListOptions) (*OrganizationList, error) {
	req, err := s.client.newRequest("GET", "organizations", &options)
	if err != nil {
		return nil, err
	}

	orgl := &OrganizationList{}
	err = s.client.do(ctx, req, orgl)
	if err != nil {
		return nil, err
	}

	return orgl, nil
}

// OrganizationCreateOptions represents the options for creating an organization.
type OrganizationCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,organizations"`

	// Name of the organization.
	Name *string `jsonapi:"attr,name"`

	// Admin email address.
	Email *string `jsonapi:"attr,email"`

	// Session expiration (minutes).
	SessionRemember *int `jsonapi:"attr,session-remember,omitempty"`

	// Session timeout after inactivity (minutes).
	SessionTimeout *int `jsonapi:"attr,session-timeout,omitempty"`

	// Authentication policy.
	CollaboratorAuthPolicy *AuthPolicyType `jsonapi:"attr,collaborator-auth-policy,omitempty"`

	// Enable Cost Estimation
	CostEstimationEnabled *bool `jsonapi:"attr,cost-estimation-enabled,omitempty"`

	// The name of the "owners" team
	OwnersTeamSAMLRoleID *string `jsonapi:"attr,owners-team-saml-role-id,omitempty"`

	// SendPassingStatusesForUntriggeredSpeculativePlans toggles behavior of untriggered speculative plans to send status updates to version control systems like GitHub.
	SendPassingStatusesForUntriggeredSpeculativePlans *bool `jsonapi:"attr,send-passing-statuses-for-untriggered-speculative-plans,omitempty"`
}

func (o OrganizationCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validStringID(o.Name) {
		return ErrInvalidName
	}
	if !validString(o.Email) {
		return errors.New("email is required")
	}
	return nil
}

// Create a new organization with the given options.
func (s *organizations) Create(ctx context.Context, options OrganizationCreateOptions) (*Organization, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("POST", "organizations", &options)
	if err != nil {
		return nil, err
	}

	org := &Organization{}
	err = s.client.do(ctx, req, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Read an organization by its name.
func (s *organizations) Read(ctx context.Context, organization string) (*Organization, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s", url.QueryEscape(organization))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	org := &Organization{}
	err = s.client.do(ctx, req, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// OrganizationUpdateOptions represents the options for updating an organization.
type OrganizationUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,organizations"`

	// New name for the organization.
	Name *string `jsonapi:"attr,name,omitempty"`

	// New admin email address.
	Email *string `jsonapi:"attr,email,omitempty"`

	// Session expiration (minutes).
	SessionRemember *int `jsonapi:"attr,session-remember,omitempty"`

	// Session timeout after inactivity (minutes).
	SessionTimeout *int `jsonapi:"attr,session-timeout,omitempty"`

	// Authentication policy.
	CollaboratorAuthPolicy *AuthPolicyType `jsonapi:"attr,collaborator-auth-policy,omitempty"`

	// Enable Cost Estimation
	CostEstimationEnabled *bool `jsonapi:"attr,cost-estimation-enabled,omitempty"`

	// The name of the "owners" team
	OwnersTeamSAMLRoleID *string `jsonapi:"attr,owners-team-saml-role-id,omitempty"`

	// SendPassingStatusesForUntriggeredSpeculativePlans toggles behavior of untriggered speculative plans to send status updates to version control systems like GitHub.
	SendPassingStatusesForUntriggeredSpeculativePlans *bool `jsonapi:"attr,send-passing-statuses-for-untriggered-speculative-plans,omitempty"`
}

// Update attributes of an existing organization.
func (s *organizations) Update(ctx context.Context, organization string, options OrganizationUpdateOptions) (*Organization, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s", url.QueryEscape(organization))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	org := &Organization{}
	err = s.client.do(ctx, req, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Delete an organization by its name.
func (s *organizations) Delete(ctx context.Context, organization string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s", url.QueryEscape(organization))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// Capacity shows the currently used capacity of an organization.
func (s *organizations) Capacity(ctx context.Context, organization string) (*Capacity, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/capacity", url.QueryEscape(organization))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	c := &Capacity{}
	err = s.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Entitlements shows the entitlements of an organization.
func (s *organizations) Entitlements(ctx context.Context, organization string) (*Entitlements, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/entitlement-set", url.QueryEscape(organization))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	e := &Entitlements{}
	err = s.client.do(ctx, req, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// RunQueueOptions represents the options for showing the queue.
type RunQueueOptions struct {
	ListOptions
}

// RunQueue shows the current run queue of an organization.
func (s *organizations) RunQueue(ctx context.Context, organization string, options RunQueueOptions) (*RunQueue, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/runs/queue", url.QueryEscape(organization))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	rq := &RunQueue{}
	err = s.client.do(ctx, req, rq)
	if err != nil {
		return nil, err
	}

	return rq, nil
}
