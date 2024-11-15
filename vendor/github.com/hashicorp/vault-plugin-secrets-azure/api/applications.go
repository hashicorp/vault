// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/google/uuid"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	auth "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type ApplicationsClient interface {
	GetApplication(ctx context.Context, applicationObjectID string) (Application, error)
	CreateApplication(ctx context.Context, displayName string, signInAudience string, tags []string) (Application, error)
	DeleteApplication(ctx context.Context, applicationObjectID string, permanentlyDelete bool) error
	ListApplications(ctx context.Context, filter string) ([]Application, error)
	AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (PasswordCredential, error)
	RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID string) error
}

var _ ApplicationsClient = (*MSGraphClient)(nil)
var _ GroupsClient = (*MSGraphClient)(nil)
var _ ServicePrincipalClient = (*MSGraphClient)(nil)

type MSGraphClient struct {
	client *msgraphsdkgo.GraphServiceClient
}

type Application struct {
	AppID               string
	AppObjectID         string
	PasswordCredentials []PasswordCredential
}

type PasswordCredential struct {
	EndDate    time.Time
	KeyID      string
	SecretText string
}

// NewMSGraphClient returns a new MSGraphClient configured to interact with
// the Microsoft Graph API. It can be configured to target alternative national cloud
// deployments via graphURI. For details on the client configuration see
// https://learn.microsoft.com/en-us/graph/sdks/national-clouds
func NewMSGraphClient(graphURI string, creds azcore.TokenCredential) (*MSGraphClient, error) {
	scopes := []string{
		fmt.Sprintf("%s/.default", graphURI),
	}

	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(creds, scopes)
	if err != nil {
		return nil, err
	}

	adapter, err := msgraphsdkgo.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return nil, err
	}

	adapter.SetBaseUrl(fmt.Sprintf("%s/v1.0", graphURI))
	client := msgraphsdkgo.NewGraphServiceClient(adapter)

	ac := &MSGraphClient{
		client: client,
	}
	return ac, nil
}

func (c *MSGraphClient) GetApplication(ctx context.Context, applicationObjectID string) (Application, error) {
	app, err := c.client.Applications().ByApplicationId(applicationObjectID).Get(ctx, nil)
	if err != nil {
		return Application{}, err
	}

	if app == nil {
		return Application{}, fmt.Errorf("no application found")
	}

	return Application{
		AppID:               *app.GetAppId(),
		AppObjectID:         *app.GetId(),
		PasswordCredentials: getPasswordCredentialsForApplication(app),
	}, nil
}

func (c *MSGraphClient) ListApplications(ctx context.Context, filter string) ([]Application, error) {

	req := &applications.ApplicationsRequestBuilderGetQueryParameters{
		Filter: &filter,
	}
	configuration := &applications.ApplicationsRequestBuilderGetRequestConfiguration{
		QueryParameters: req,
	}
	resp, err := c.client.Applications().Get(ctx, configuration)
	if err != nil {
		return nil, err
	}

	var apps []Application
	for _, app := range resp.GetValue() {
		apps = append(apps, getApplicationResponse(app))
	}

	return apps, nil
}

// CreateApplication create a new Azure application object.
func (c *MSGraphClient) CreateApplication(ctx context.Context, displayName string, signInAudience string, tags []string) (Application, error) {
	requestBody := models.NewApplication()
	requestBody.SetDisplayName(&displayName)
	requestBody.SetTags(tags)

	// only set signInAudience if it's non-empty
	if signInAudience != "" {
		requestBody.SetSignInAudience(&signInAudience)
	}

	resp, err := c.client.Applications().Post(ctx, requestBody, nil)
	if err != nil {
		return Application{}, err
	}

	return getApplicationResponse(resp), nil
}

// DeleteApplication deletes an Azure application object.
// This will in turn remove the service principal (but not the role assignments).
func (c *MSGraphClient) DeleteApplication(ctx context.Context, applicationObjectID string, permanentlyDelete bool) error {
	err := c.client.Applications().ByApplicationId(applicationObjectID).Delete(ctx, nil)
	if err != nil {
		return err
	}

	if permanentlyDelete {
		err = c.client.Directory().DeletedItems().ByDirectoryObjectId(applicationObjectID).Delete(ctx, nil)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *MSGraphClient) AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (PasswordCredential, error) {
	requestBody := applications.NewItemAddPasswordPostRequestBody()
	passwordCredential := models.NewPasswordCredential()
	passwordCredential.SetDisplayName(&displayName)
	passwordCredential.SetEndDateTime(&endDateTime)
	requestBody.SetPasswordCredential(passwordCredential)

	resp, err := c.client.Applications().ByApplicationId(applicationObjectID).AddPassword().Post(ctx, requestBody, nil)
	if err != nil {
		return PasswordCredential{}, err
	}

	return getPasswordCredentialResponse(resp), nil
}

func (c *MSGraphClient) RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID string) error {
	requestBody := applications.NewItemRemovePasswordPostRequestBody()
	kid, err := uuid.Parse(keyID)
	if err != nil {
		return err
	}

	requestBody.SetKeyId(&kid)

	return c.client.Applications().ByApplicationId(applicationObjectID).RemovePassword().Post(ctx, requestBody, nil)
}

func getPasswordCredentialsForApplication(app models.Applicationable) []PasswordCredential {
	var appCredentials []PasswordCredential
	creds := app.GetPasswordCredentials()
	if creds != nil {
		for _, cred := range creds {
			appCredentials = append(appCredentials, getPasswordCredentialResponse(cred))
		}
	}

	return appCredentials
}

func ptrToString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func getApplicationResponse(app models.Applicationable) Application {
	if app != nil {
		return Application{
			AppID:               ptrToString(app.GetAppId()),
			AppObjectID:         ptrToString(app.GetId()),
			PasswordCredentials: getPasswordCredentialsForApplication(app),
		}

	}

	// return zero-value result if app in nil
	// or fields can't be dereferenced
	return Application{
		AppID:               "",
		AppObjectID:         "",
		PasswordCredentials: []PasswordCredential{},
	}
}

func getPasswordCredentialResponse(cred models.PasswordCredentialable) PasswordCredential {
	if cred != nil {
		return PasswordCredential{
			SecretText: ptrToString(cred.GetSecretText()),
			EndDate:    *cred.GetEndDateTime(),
			KeyID:      cred.GetKeyId().String(),
		}
	}
	return PasswordCredential{
		SecretText: "",
		EndDate:    time.Time{},
		KeyID:      "",
	}
}
