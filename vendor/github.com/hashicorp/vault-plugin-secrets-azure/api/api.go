package api

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/authorization/mgmt/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
)

// AzureProvider is an interface to access underlying Azure Client objects and supporting services.
// Where practical the original function signature is preserved. Client provides higher
// level operations atop AzureProvider.
type AzureProvider interface {
	ApplicationsClient
	GroupsClient
	ServicePrincipalClient

	CreateRoleAssignment(
		ctx context.Context,
		scope string,
		roleAssignmentName string,
		parameters authorization.RoleAssignmentCreateParameters) (authorization.RoleAssignment, error)
	DeleteRoleAssignmentByID(ctx context.Context, roleID string) (authorization.RoleAssignment, error)

	ListRoleDefinitions(ctx context.Context, scope string, filter string) ([]authorization.RoleDefinition, error)
	GetRoleDefinitionByID(ctx context.Context, roleID string) (authorization.RoleDefinition, error)
}

type ApplicationsClient interface {
	GetApplication(ctx context.Context, applicationObjectID string) (ApplicationResult, error)
	CreateApplication(ctx context.Context, displayName string) (ApplicationResult, error)
	DeleteApplication(ctx context.Context, applicationObjectID string) error
	ListApplications(ctx context.Context, filter string) ([]ApplicationResult, error)
	AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (PasswordCredentialResult, error)
	RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID string) error
}

type PasswordCredential struct {
	DisplayName *string    `json:"displayName"`
	StartDate   *date.Time `json:"startDateTime,omitempty"`
	EndDate     *date.Time `json:"endDateTime,omitempty"`
	KeyID       *string    `json:"keyId,omitempty"`
	SecretText  *string    `json:"secretText,omitempty"`
}

type PasswordCredentialResult struct {
	autorest.Response `json:"-"`

	PasswordCredential
}

type ApplicationResult struct {
	autorest.Response `json:"-"`

	AppID               *string               `json:"appId,omitempty"`
	ID                  *string               `json:"id,omitempty"`
	PasswordCredentials []*PasswordCredential `json:"passwordCredentials,omitempty"`
}
