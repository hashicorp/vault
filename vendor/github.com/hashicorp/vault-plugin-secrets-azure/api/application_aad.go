package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-uuid"
)

type ActiveDirectoryApplicationClient struct {
	Client    *graphrbac.ApplicationsClient
	Passwords Passwords
}

func (a *ActiveDirectoryApplicationClient) GetApplication(ctx context.Context, applicationObjectID string) (ApplicationResult, error) {
	app, err := a.Client.Get(ctx, applicationObjectID)
	if err != nil {
		return ApplicationResult{}, err
	}

	return ApplicationResult{
		AppID: app.AppID,
		ID:    app.ObjectID,
	}, nil
}

func (a *ActiveDirectoryApplicationClient) ListApplications(ctx context.Context, filter string) ([]ApplicationResult, error) {
	resp, err := a.Client.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := []ApplicationResult{}
	for resp.NotDone() {
		for _, app := range resp.Values() {
			passCreds := []*PasswordCredential{}
			for _, rawPC := range *app.PasswordCredentials {
				pc := &PasswordCredential{
					StartDate: rawPC.StartDate,
					EndDate:   rawPC.EndDate,
					KeyID:     rawPC.KeyID,
				}
				passCreds = append(passCreds, pc)
			}
			appResult := ApplicationResult{
				AppID:               app.AppID,
				ID:                  app.ObjectID,
				PasswordCredentials: passCreds,
			}
			results = append(results, appResult)
		}
		err = resp.NextWithContext(ctx)
		if err != nil {
			return results, fmt.Errorf("failed to get all results: %w", err)
		}
	}
	return results, nil
}

func (a *ActiveDirectoryApplicationClient) CreateApplication(ctx context.Context, displayName string) (ApplicationResult, error) {
	appURL := fmt.Sprintf("https://%s", displayName)

	app, err := a.Client.Create(ctx, graphrbac.ApplicationCreateParameters{
		AvailableToOtherTenants: to.BoolPtr(false),
		DisplayName:             to.StringPtr(displayName),
		Homepage:                to.StringPtr(appURL),
		IdentifierUris:          to.StringSlicePtr([]string{appURL}),
	})
	if err != nil {
		return ApplicationResult{}, err
	}

	return ApplicationResult{
		AppID: app.AppID,
		ID:    app.ObjectID,
	}, nil
}

func (a *ActiveDirectoryApplicationClient) DeleteApplication(ctx context.Context, applicationObjectID string) error {
	resp, err := a.Client.Delete(ctx, applicationObjectID)
	if resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if err != nil {
		return err
	}
	return nil
}

func (a *ActiveDirectoryApplicationClient) AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (PasswordCredentialResult, error) {
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		return PasswordCredentialResult{}, err
	}

	// Key IDs are not secret, and they're a convenient way for an operator to identify Vault-generated
	// passwords. These must be UUIDs, so the three leading bytes will be used as an indicator.
	keyID = "ffffff" + keyID[6:]

	password, err := a.Passwords.Generate(ctx)
	if err != nil {
		return PasswordCredentialResult{}, err
	}

	now := date.Time{Time: time.Now().UTC()}
	cred := graphrbac.PasswordCredential{
		StartDate: &now,
		EndDate:   &date.Time{endDateTime},
		KeyID:     to.StringPtr(keyID),
		Value:     to.StringPtr(password),
	}

	// Load current credentials
	resp, err := a.Client.ListPasswordCredentials(ctx, applicationObjectID)
	if err != nil {
		return PasswordCredentialResult{}, fmt.Errorf("error fetching credentials: %w", err)
	}
	curCreds := *resp.Value

	// Add and save credentials
	curCreds = append(curCreds, cred)

	if _, err := a.Client.UpdatePasswordCredentials(ctx, applicationObjectID,
		graphrbac.PasswordCredentialsUpdateParameters{
			Value: &curCreds,
		},
	); err != nil {
		if strings.Contains(err.Error(), "size of the object has exceeded its limit") {
			err = errors.New("maximum number of Application passwords reached")
		}
		return PasswordCredentialResult{}, fmt.Errorf("error updating credentials: %w", err)
	}

	result := PasswordCredentialResult{
		PasswordCredential: PasswordCredential{
			DisplayName: to.StringPtr(displayName),
			StartDate:   &now,
			EndDate:     &date.Time{endDateTime},
			KeyID:       to.StringPtr(keyID),
			SecretText:  to.StringPtr(password),
		},
	}
	return result, nil
}

func (a *ActiveDirectoryApplicationClient) RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID string) error {
	// Load current credentials
	resp, err := a.Client.ListPasswordCredentials(ctx, applicationObjectID)
	if err != nil {
		return fmt.Errorf("error fetching credentials: %w", err)
	}
	curCreds := *resp.Value

	// Remove credential
	found := false
	for i := range curCreds {
		if to.String(curCreds[i].KeyID) == keyID {
			curCreds[i] = curCreds[len(curCreds)-1]
			curCreds = curCreds[:len(curCreds)-1]
			found = true
			break
		}
	}

	// KeyID is not present, so nothing to do
	if !found {
		return nil
	}

	// Save new credentials list
	_, err = a.Client.UpdatePasswordCredentials(ctx, applicationObjectID,
		graphrbac.PasswordCredentialsUpdateParameters{
			Value: &curCreds,
		},
	)
	if err != nil {
		return fmt.Errorf("error updating credentials: %w", err)
	}

	return nil
}
