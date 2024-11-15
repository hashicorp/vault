// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresecrets

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault-plugin-secrets-azure/api"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
)

// AzureProvider is an interface to access underlying Azure Client objects and supporting services.
// Where practical the original function signature is preserved. Client provides higher
// level operations atop AzureProvider.
type AzureProvider interface {
	api.ApplicationsClient
	api.GroupsClient
	api.ServicePrincipalClient

	CreateRoleAssignment(
		ctx context.Context,
		scope string,
		roleAssignmentName string,
		parameters armauthorization.RoleAssignmentCreateParameters) (armauthorization.RoleAssignmentsClientCreateResponse, error)
	DeleteRoleAssignmentByID(ctx context.Context, roleID string) (armauthorization.RoleAssignmentsClientDeleteByIDResponse, error)
	ListRoleDefinitions(ctx context.Context, scope string, filter string) (result []*armauthorization.RoleDefinition, err error)
	GetRoleDefinitionByID(ctx context.Context, roleID string) (result armauthorization.RoleDefinitionsClientGetByIDResponse, err error)
}

var _ AzureProvider = (*provider)(nil)

// provider is a concrete implementation of AzureProvider. In most cases it is a simple passthrough
// to the appropriate client object. But if the response requires processing that is more practical
// at this layer, the response signature may different from the Azure signature.
type provider struct {
	settings *clientSettings

	appClient    api.ApplicationsClient
	spClient     api.ServicePrincipalClient
	groupsClient api.GroupsClient
	raClient     *armauthorization.RoleAssignmentsClient
	rdClient     *armauthorization.RoleDefinitionsClient
}

// newAzureProvider creates an azureProvider, backed by Azure client objects for underlying services.
func newAzureProvider(ctx context.Context, logger hclog.Logger, sys logical.SystemView, settings *clientSettings) (AzureProvider, error) {
	httpClient := cleanhttp.DefaultClient()

	cred, err := getTokenCredential(ctx, logger, sys, settings)
	if err != nil {
		return nil, err
	}

	msGraphAppClient, err := api.NewMSGraphClient(settings.GraphURI, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to create MS graph client: %w", err)
	}

	opts := getClientOptions(settings, httpClient)

	raClient, err := armauthorization.NewRoleAssignmentsClient(settings.SubscriptionID, cred, opts)
	if err != nil {
		return nil, err
	}

	rdClient, err := armauthorization.NewRoleDefinitionsClient(cred, opts)
	if err != nil {
		return nil, err
	}

	p := &provider{
		appClient:    msGraphAppClient,
		spClient:     msGraphAppClient,
		groupsClient: msGraphAppClient,
		raClient:     raClient,
		rdClient:     rdClient,
	}

	return p, nil
}

func getTokenCredential(ctx context.Context, logger hclog.Logger, sys logical.SystemView, s *clientSettings) (azcore.TokenCredential, error) {
	clientCloudOpts := azcore.ClientOptions{Cloud: s.CloudConfig}

	if s.ClientSecret != "" {
		options := &azidentity.ClientSecretCredentialOptions{
			ClientOptions: clientCloudOpts,
		}

		cred, err := azidentity.NewClientSecretCredential(s.TenantID, s.ClientID,
			s.ClientSecret, options)
		if err != nil {
			return nil, fmt.Errorf("failed to create client secret token credential: %w", err)
		}

		return cred, nil
	}

	if s.IdentityTokenAudience != "" {
		options := &azidentity.ClientAssertionCredentialOptions{
			ClientOptions: clientCloudOpts,
		}
		getAssertion := getAssertionFunc(ctx, logger, sys, s)
		cred, err := azidentity.NewClientAssertionCredential(s.TenantID, s.ClientID,
			getAssertion, options)
		if err != nil {
			return nil, fmt.Errorf("failed to create client assertion credential: %w", err)
		}

		return cred, nil
	}

	// Fall back to using managed service identity
	options := &azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: clientCloudOpts,
	}
	cred, err := azidentity.NewManagedIdentityCredential(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create managed identity token credential: %w", err)
	}

	return cred, nil
}

type getAssertion func(context.Context) (string, error)

func getAssertionFunc(ctx context.Context, logger hclog.Logger, sys logical.SystemView, s *clientSettings) getAssertion {
	return func(ctx context.Context) (string, error) {
		req := &pluginutil.IdentityTokenRequest{
			Audience: s.IdentityTokenAudience,
			TTL:      s.IdentityTokenTTL * time.Second,
		}
		resp, err := sys.GenerateIdentityToken(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to generate plugin identity token: %w", err)
		}
		logger.Info("fetched new plugin identity token")

		if resp.TTL < req.TTL {
			logger.Debug("generated plugin identity token has shorter TTL than requested",
				"requested", req.TTL, "actual", resp.TTL)
		}

		return resp.Token.Token(), nil
	}
}

// transporter implements the azure exported.Transporter interface to send HTTP
// requests. This allows us to set our custom http client and user agent.
type transporter struct {
	pluginEnv *logical.PluginEnvironment
	sender    *http.Client
}

func (tp transporter) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", useragent.PluginString(tp.pluginEnv,
		userAgentPluginName))

	client := tp.sender

	// don't attempt redirects so we aren't acting as an unintended network proxy
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getClientOptions(s *clientSettings, httpClient *http.Client) *arm.ClientOptions {
	return &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: s.CloudConfig,
			Transport: transporter{
				pluginEnv: s.PluginEnv,
				sender:    httpClient,
			},
		},
	}
}

// CreateApplication create a new Azure application object.
func (p *provider) CreateApplication(ctx context.Context, displayName string, signInAudience string, tags []string) (result api.Application, err error) {
	return p.appClient.CreateApplication(ctx, displayName, signInAudience, tags)
}

func (p *provider) GetApplication(ctx context.Context, applicationObjectID string) (result api.Application, err error) {
	return p.appClient.GetApplication(ctx, applicationObjectID)
}

func (p *provider) ListApplications(ctx context.Context, filter string) ([]api.Application, error) {
	return p.appClient.ListApplications(ctx, filter)
}

// DeleteApplication deletes an Azure application object.
// This will in turn remove the service principal (but not the role assignments).
func (p *provider) DeleteApplication(ctx context.Context, applicationObjectID string, permanentlyDelete bool) error {
	return p.appClient.DeleteApplication(ctx, applicationObjectID, permanentlyDelete)
}

func (p *provider) AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (result api.PasswordCredential, err error) {
	return p.appClient.AddApplicationPassword(ctx, applicationObjectID, displayName, endDateTime)
}

func (p *provider) RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID string) (err error) {
	return p.appClient.RemoveApplicationPassword(ctx, applicationObjectID, keyID)
}

// CreateServicePrincipal creates a new Azure service principal.
// An Application must be created prior to calling this and pass in parameters.
func (p *provider) CreateServicePrincipal(ctx context.Context, appID string, startDate time.Time, endDate time.Time) (id string, password string, err error) {
	return p.spClient.CreateServicePrincipal(ctx, appID, startDate, endDate)
}

func (p *provider) DeleteServicePrincipal(ctx context.Context, spObjectID string, permanentlyDelete bool) error {
	return p.spClient.DeleteServicePrincipal(ctx, spObjectID, permanentlyDelete)
}

// ListRoles like all Azure roles with a scope (often subscription).
func (p *provider) ListRoleDefinitions(ctx context.Context, scope string, filter string) (result []*armauthorization.RoleDefinition, err error) {
	options := armauthorization.RoleDefinitionsClientListOptions{
		Filter: &filter,
	}
	page := p.rdClient.NewListPager(scope, &options)
	listResp, err := page.NextPage(ctx)
	if err != nil {
		return nil, err
	}

	return listResp.Value, err
}

// GetRoleDefinitionByID fetches the full role definition given a roleID.
func (p *provider) GetRoleDefinitionByID(ctx context.Context, roleID string) (result armauthorization.RoleDefinitionsClientGetByIDResponse, err error) {
	return p.rdClient.GetByID(ctx, roleID, nil)
}

// CreateRoleAssignment assigns a role to a service principal.
func (p *provider) CreateRoleAssignment(ctx context.Context, scope string, roleAssignmentName string, parameters armauthorization.RoleAssignmentCreateParameters) (armauthorization.RoleAssignmentsClientCreateResponse, error) {
	return p.raClient.Create(ctx, scope, roleAssignmentName, parameters, nil)
}

// GetRoleAssignmentByID fetches the full role assignment info given a roleAssignmentID.
func (p *provider) GetRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (armauthorization.RoleAssignmentsClientGetByIDResponse, error) {
	return p.raClient.GetByID(ctx, roleAssignmentID, nil)
}

// DeleteRoleAssignmentByID deletes a role assignment.
func (p *provider) DeleteRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (armauthorization.RoleAssignmentsClientDeleteByIDResponse, error) {
	return p.raClient.DeleteByID(ctx, roleAssignmentID, nil)
}

// AddGroupMember adds a member to a Group.
func (p *provider) AddGroupMember(ctx context.Context, groupObjectID string, memberObjectID string) (err error) {
	return p.groupsClient.AddGroupMember(ctx, groupObjectID, memberObjectID)
}

// RemoveGroupMember removes a member from a Group.
func (p *provider) RemoveGroupMember(ctx context.Context, groupObjectID, memberObjectID string) (err error) {
	return p.groupsClient.RemoveGroupMember(ctx, groupObjectID, memberObjectID)
}

// GetGroup gets group information from the directory.
func (p *provider) GetGroup(ctx context.Context, objectID string) (result api.Group, err error) {
	return p.groupsClient.GetGroup(ctx, objectID)
}

// ListGroups gets list of groups for the current tenant.
func (p *provider) ListGroups(ctx context.Context, filter string) (result []api.Group, err error) {
	return p.groupsClient.ListGroups(ctx, filter)
}
