// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
)

type ServicePrincipalClient interface {
	// CreateServicePrincipal in Azure. The password returned is the actual password that the appID was created with
	CreateServicePrincipal(ctx context.Context, appID string, startDate time.Time, endDate time.Time) (id string, password string, err error)
	DeleteServicePrincipal(ctx context.Context, spObjectID string, permanentlyDelete bool) error
}

type ServicePrincipal struct {
	ID    string
	AppID string
}

func (c *MSGraphClient) CreateServicePrincipal(ctx context.Context, appID string, startDate time.Time, endDate time.Time) (string, string, error) {
	spReq := models.NewServicePrincipal()
	spReq.SetAppId(&appID)

	sp, err := c.client.ServicePrincipals().Post(ctx, spReq, nil)
	if err != nil {
		return "", "", err
	}

	spID := sp.GetId()

	passwordReq := serviceprincipals.NewItemAddPasswordPostRequestBody()
	passwordCredential := models.NewPasswordCredential()
	passwordCredential.SetStartDateTime(&startDate)
	passwordCredential.SetEndDateTime(&endDate)

	passwordReq.SetPasswordCredential(passwordCredential)

	password, err := c.client.ServicePrincipals().ByServicePrincipalId(*spID).AddPassword().Post(ctx, passwordReq, nil)

	if err != nil {
		e := c.DeleteServicePrincipal(ctx, *spID, false)
		merr := multierror.Append(err, e)
		return "", "", merr.ErrorOrNil()
	}
	return *spID, *password.GetSecretText(), nil
}

func (c *MSGraphClient) DeleteServicePrincipal(ctx context.Context, spObjectID string, permanentlyDelete bool) error {
	err := c.client.ServicePrincipals().ByServicePrincipalId(spObjectID).Delete(ctx, nil)

	if permanentlyDelete {
		e := c.client.Directory().DeletedItems().ByDirectoryObjectId(spObjectID).Delete(ctx, nil)
		merr := multierror.Append(err, e)
		return merr.ErrorOrNil()
	}

	return err
}

func (c *MSGraphClient) ListServicePrincipals(ctx context.Context, spObjectID string) ([]ServicePrincipal, error) {
	filter := fmt.Sprintf("appId eq '%s'", spObjectID)
	requestParameters := &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
		Filter: &filter,
	}

	configuration := &serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	spList, err := c.client.ServicePrincipals().Get(ctx, configuration)
	if err != nil {
		return nil, err
	}

	var result []ServicePrincipal
	for _, sp := range spList.GetValue() {
		result = append(result, getServicePrincipalResponse(sp))
	}
	return result, nil
}

func (c *MSGraphClient) GetServicePrincipalByID(ctx context.Context, spObjectID string) (ServicePrincipal, error) {
	sp, err := c.client.ServicePrincipals().ByServicePrincipalId(spObjectID).Get(ctx, nil)
	if err != nil {
		return ServicePrincipal{}, err
	}

	return getServicePrincipalResponse(sp), nil
}

func getServicePrincipalResponse(sp models.ServicePrincipalable) ServicePrincipal {
	if sp != nil {
		return ServicePrincipal{
			ID:    ptrToString(sp.GetId()),
			AppID: ptrToString(sp.GetAppId()),
		}
	}
	return ServicePrincipal{
		ID:    "",
		AppID: "",
	}
}
