package azuresecrets

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hashicorp/vault/sdk/helper/useragent"
)

// AzureProvider is an interface to access underlying Azure client objects and supporting services.
// Where practical the original function signature is preserved. client provides higher
// level operations atop AzureProvider.
type AzureProvider interface {
	ApplicationsClient
	ServicePrincipalsClient
	RoleAssignmentsClient
	RoleDefinitionsClient
}

type ApplicationsClient interface {
	CreateApplication(ctx context.Context, parameters graphrbac.ApplicationCreateParameters) (graphrbac.Application, error)
	DeleteApplication(ctx context.Context, applicationObjectID string) (autorest.Response, error)
	GetApplication(ctx context.Context, applicationObjectID string) (graphrbac.Application, error)
	UpdateApplicationPasswordCredentials(
		ctx context.Context,
		applicationObjectID string,
		parameters graphrbac.PasswordCredentialsUpdateParameters) (result autorest.Response, err error)
	ListApplicationPasswordCredentials(ctx context.Context, applicationObjectID string) (result graphrbac.PasswordCredentialListResult, err error)
}

type ServicePrincipalsClient interface {
	CreateServicePrincipal(ctx context.Context, parameters graphrbac.ServicePrincipalCreateParameters) (graphrbac.ServicePrincipal, error)
}

type RoleAssignmentsClient interface {
	CreateRoleAssignment(
		ctx context.Context,
		scope string,
		roleAssignmentName string,
		parameters authorization.RoleAssignmentCreateParameters) (authorization.RoleAssignment, error)
	DeleteRoleAssignmentByID(ctx context.Context, roleID string) (authorization.RoleAssignment, error)
}

type RoleDefinitionsClient interface {
	ListRoles(ctx context.Context, scope string, filter string) ([]authorization.RoleDefinition, error)
	GetRoleByID(ctx context.Context, roleID string) (result authorization.RoleDefinition, err error)
}

// provider is a concrete implementation of AzureProvider. In most cases it is a simple passthrough
// to the appropriate client object. But if the response requires processing that is more practical
// at this layer, the response signature may different from the Azure signature.
type provider struct {
	settings *clientSettings

	appClient *graphrbac.ApplicationsClient
	spClient  *graphrbac.ServicePrincipalsClient
	raClient  *authorization.RoleAssignmentsClient
	rdClient  *authorization.RoleDefinitionsClient
}

// newAzureProvider creates an azureProvider, backed by Azure client objects for underlying services.
func newAzureProvider(settings *clientSettings) (AzureProvider, error) {
	// build clients that use the GraphRBAC endpoint
	authorizer, err := getAuthorizer(settings, settings.Environment.GraphEndpoint)
	if err != nil {
		return nil, err
	}

	var userAgent string
	if settings.PluginEnv != nil {
		userAgent = useragent.PluginString(settings.PluginEnv, "azure-secrets")
	} else {
		userAgent = useragent.String()
	}

	appClient := graphrbac.NewApplicationsClient(settings.TenantID)
	appClient.Authorizer = authorizer
	appClient.AddToUserAgent(userAgent)

	spClient := graphrbac.NewServicePrincipalsClient(settings.TenantID)
	spClient.Authorizer = authorizer
	spClient.AddToUserAgent(userAgent)

	// build clients that use the Resource Manager endpoint
	authorizer, err = getAuthorizer(settings, settings.Environment.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}

	raClient := authorization.NewRoleAssignmentsClient(settings.SubscriptionID)
	raClient.Authorizer = authorizer
	raClient.AddToUserAgent(userAgent)

	rdClient := authorization.NewRoleDefinitionsClient(settings.SubscriptionID)
	rdClient.Authorizer = authorizer
	rdClient.AddToUserAgent(userAgent)

	p := &provider{
		settings: settings,

		appClient: &appClient,
		spClient:  &spClient,
		raClient:  &raClient,
		rdClient:  &rdClient,
	}

	return p, nil
}

// getAuthorizer attempts to create an authorizer, preferring ClientID/Secret if present,
// and falling back to MSI if not.
func getAuthorizer(settings *clientSettings, resource string) (authorizer autorest.Authorizer, err error) {

	if settings.ClientID != "" && settings.ClientSecret != "" && settings.TenantID != "" {
		config := auth.NewClientCredentialsConfig(settings.ClientID, settings.ClientSecret, settings.TenantID)
		config.AADEndpoint = settings.Environment.ActiveDirectoryEndpoint
		config.Resource = resource
		authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	} else {
		config := auth.NewMSIConfig()
		config.Resource = resource
		authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	}

	return authorizer, nil
}

// CreateApplication create a new Azure application object.
func (p *provider) CreateApplication(ctx context.Context, parameters graphrbac.ApplicationCreateParameters) (graphrbac.Application, error) {
	return p.appClient.Create(ctx, parameters)
}

func (p *provider) GetApplication(ctx context.Context, applicationObjectID string) (graphrbac.Application, error) {
	return p.appClient.Get(ctx, applicationObjectID)
}

// DeleteApplication deletes an Azure application object.
// This will in turn remove the service principal (but not the role assignments).
func (p *provider) DeleteApplication(ctx context.Context, applicationObjectID string) (autorest.Response, error) {
	return p.appClient.Delete(ctx, applicationObjectID)
}

func (p *provider) UpdateApplicationPasswordCredentials(ctx context.Context, applicationObjectID string, parameters graphrbac.PasswordCredentialsUpdateParameters) (result autorest.Response, err error) {
	return p.appClient.UpdatePasswordCredentials(ctx, applicationObjectID, parameters)
}

func (p *provider) ListApplicationPasswordCredentials(ctx context.Context, applicationObjectID string) (result graphrbac.PasswordCredentialListResult, err error) {
	return p.appClient.ListPasswordCredentials(ctx, applicationObjectID)
}

// CreateServicePrincipal creates a new Azure service principal.
// An Application must be created prior to calling this and pass in parameters.
func (p *provider) CreateServicePrincipal(ctx context.Context, parameters graphrbac.ServicePrincipalCreateParameters) (graphrbac.ServicePrincipal, error) {
	return p.spClient.Create(ctx, parameters)
}

// ListRoles like all Azure roles with a scope (often subscription).
func (p *provider) ListRoles(ctx context.Context, scope string, filter string) (result []authorization.RoleDefinition, err error) {
	page, err := p.rdClient.List(ctx, scope, filter)

	if err != nil {
		return nil, err
	}

	return page.Values(), nil
}

// GetRoleByID fetches the full role definition given a roleID.
func (p *provider) GetRoleByID(ctx context.Context, roleID string) (result authorization.RoleDefinition, err error) {
	return p.rdClient.GetByID(ctx, roleID)
}

// CreateRoleAssignment assigns a role to a service principal.
func (p *provider) CreateRoleAssignment(ctx context.Context, scope string, roleAssignmentName string, parameters authorization.RoleAssignmentCreateParameters) (authorization.RoleAssignment, error) {
	return p.raClient.Create(ctx, scope, roleAssignmentName, parameters)
}

// GetRoleAssignmentByID fetches the full role assignment info given a roleAssignmentID.
func (p *provider) GetRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (result authorization.RoleAssignment, err error) {
	return p.raClient.GetByID(ctx, roleAssignmentID)
}

// DeleteRoleAssignmentByID deletes a role assignment.
func (p *provider) DeleteRoleAssignmentByID(ctx context.Context, roleAssignmentID string) (result authorization.RoleAssignment, err error) {
	return p.raClient.DeleteByID(ctx, roleAssignmentID)
}

// ListRoleAssignments lists all role assignments.
// There is no need for paging; the caller only cares about the the first match and whether
// there are 0, 1 or >1 items. Unpacking here is a simpler interface.
func (p *provider) ListRoleAssignments(ctx context.Context, filter string) ([]authorization.RoleAssignment, error) {
	page, err := p.raClient.List(ctx, filter)

	if err != nil {
		return nil, err
	}

	return page.Values(), nil
}
