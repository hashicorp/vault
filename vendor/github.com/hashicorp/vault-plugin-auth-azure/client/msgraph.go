// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

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

type MSGraphClient interface {
	GetApplication(ctx context.Context, clientID string) (models.Applicationable, error)
	AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (models.PasswordCredentialable, error)
	RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID *uuid.UUID) error
}

var _ MSGraphClient = (*AppClient)(nil)

type AppClient struct {
	client *msgraphsdkgo.GraphServiceClient
}

// NewMSGraphApplicationClient returns a new AppClient configured to interact with
// the Microsoft Graph API. It can be configured to target alternative national cloud
// deployments via graphURI. For details on the client configuration see
// https://learn.microsoft.com/en-us/graph/sdks/national-clouds
func NewMSGraphApplicationClient(graphURI string, creds azcore.TokenCredential) (*AppClient, error) {
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

	ac := &AppClient{
		client: client,
	}
	return ac, nil
}

func (c *AppClient) GetApplication(ctx context.Context, clientID string) (models.Applicationable, error) {
	filter := fmt.Sprintf("appId eq '%s'", clientID)
	req := applications.ApplicationsRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	resp, err := c.client.Applications().Get(ctx, &req)
	if err != nil {
		return nil, err
	}

	apps := resp.GetValue()
	if len(apps) == 0 {
		return nil, fmt.Errorf("no application found")
	}
	if len(apps) > 1 {
		return nil, fmt.Errorf("multiple applications found - double check your client_id")
	}

	return apps[0], nil
}

func (c *AppClient) AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (models.PasswordCredentialable, error) {
	requestBody := applications.NewItemAddPasswordPostRequestBody()
	passwordCredential := models.NewPasswordCredential()
	passwordCredential.SetDisplayName(&displayName)
	passwordCredential.SetEndDateTime(&endDateTime)
	requestBody.SetPasswordCredential(passwordCredential)

	resp, err := c.client.Applications().ByApplicationId(applicationObjectID).AddPassword().Post(ctx, requestBody, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *AppClient) RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID *uuid.UUID) error {
	requestBody := applications.NewItemRemovePasswordPostRequestBody()
	requestBody.SetKeyId(keyID)

	return c.client.Applications().ByApplicationId(applicationObjectID).RemovePassword().Post(ctx, requestBody, nil)
}
