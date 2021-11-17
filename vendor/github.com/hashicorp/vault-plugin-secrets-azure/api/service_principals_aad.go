package api

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-uuid"
)

var _ ServicePrincipalClient = (*AADServicePrincipalsClient)(nil)

type AADServicePrincipalsClient struct {
	Client    graphrbac.ServicePrincipalsClient
	Passwords Passwords
}

func (c AADServicePrincipalsClient) CreateServicePrincipal(ctx context.Context, appID string, startDate time.Time, endDate time.Time) (string, string, error) {
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		return "", "", err
	}

	password, err := c.Passwords.Generate(ctx)
	if err != nil {
		return "", "", err
	}

	clientParams := graphrbac.ServicePrincipalCreateParameters{
		AppID:          to.StringPtr(appID),
		AccountEnabled: to.BoolPtr(true),
		PasswordCredentials: &[]graphrbac.PasswordCredential{
			graphrbac.PasswordCredential{
				StartDate: &date.Time{startDate},
				EndDate:   &date.Time{endDate},
				KeyID:     &keyID,
				Value:     &password,
			},
		},
	}
	sp, err := c.Client.Create(ctx, clientParams)
	if err != nil {
		return "", "", err
	}
	return *sp.ObjectID, password, nil
}
