package azuresecrets

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
)

const (
	appNamePrefix  = "vault-"
	retryTimeout   = 80 * time.Second
	clientLifetime = 30 * time.Minute
)

// client offers higher level Azure operations that provide a simpler interface
// for handlers. It in turn relies on a Provider interface to access the lower level
// Azure Client SDK methods.
type client struct {
	provider   AzureProvider
	settings   *clientSettings
	expiration time.Time
}

// Valid returns whether the client defined and not expired.
func (c *client) Valid() bool {
	return c != nil && time.Now().Before(c.expiration)
}

func (b *azureSecretBackend) getClient(ctx context.Context, s logical.Storage) (*client, error) {
	b.lock.RLock()
	unlockFunc := b.lock.RUnlock
	defer func() { unlockFunc() }()

	if b.client.Valid() {
		return b.client, nil
	}

	b.lock.RUnlock()
	b.lock.Lock()
	unlockFunc = b.lock.Unlock

	if b.client.Valid() {
		return b.client, nil
	}

	if b.settings == nil {
		config, err := b.getConfig(ctx, s)
		if err != nil {
			return nil, err
		}
		if config == nil {
			config = new(azureConfig)
		}

		settings, err := b.getClientSettings(ctx, config)
		if err != nil {
			return nil, err
		}
		b.settings = settings
	}

	p, err := b.getProvider(b.settings)
	if err != nil {
		return nil, err
	}

	c := &client{
		provider:   p,
		settings:   b.settings,
		expiration: time.Now().Add(clientLifetime),
	}
	b.client = c

	return c, nil
}

// createApp creates a new Azure application.
// An Application is a needed to create service principals used by
// the caller for authentication.
func (c *client) createApp(ctx context.Context) (app *graphrbac.Application, err error) {
	name, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	name = appNamePrefix + name

	appURL := fmt.Sprintf("https://%s", name)

	result, err := c.provider.CreateApplication(ctx, graphrbac.ApplicationCreateParameters{
		AvailableToOtherTenants: to.BoolPtr(false),
		DisplayName:             to.StringPtr(name),
		Homepage:                to.StringPtr(appURL),
		IdentifierUris:          to.StringSlicePtr([]string{appURL}),
	})

	return &result, err
}

// createSP creates a new service principal.
func (c *client) createSP(
	ctx context.Context,
	app *graphrbac.Application,
	duration time.Duration) (*graphrbac.ServicePrincipal, string, error) {

	// Generate a random key (which must be a UUID) and password
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, "", err
	}

	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, "", err
	}

	resultRaw, err := retry(ctx, func() (interface{}, bool, error) {
		now := time.Now().UTC()
		result, err := c.provider.CreateServicePrincipal(ctx, graphrbac.ServicePrincipalCreateParameters{
			AppID:          app.AppID,
			AccountEnabled: to.BoolPtr(true),
			PasswordCredentials: &[]graphrbac.PasswordCredential{
				graphrbac.PasswordCredential{
					StartDate: &date.Time{Time: now},
					EndDate:   &date.Time{Time: now.Add(duration)},
					KeyID:     to.StringPtr(keyID),
					Value:     to.StringPtr(password),
				},
			},
		})

		// Propagation delays within Azure can cause this error occasionally, so don't quit on it.
		if err != nil && strings.Contains(err.Error(), "does not reference a valid application object") {
			return nil, false, nil
		}

		return result, true, err
	})

	result := resultRaw.(graphrbac.ServicePrincipal)

	return &result, password, err
}

// addAppPassword adds a new password to an App's credentials list.
func (c *client) addAppPassword(ctx context.Context, appObjID string, duration time.Duration) (string, string, error) {
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		return "", "", err
	}

	// Key IDs are not secret, and they're a convenient way for an operator to identify Vault-generated
	// passwords. These must be UUIDs, so the three leading bytes will be used as an indicator.
	keyID = "ffffff" + keyID[6:]

	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", "", err
	}

	now := time.Now().UTC()
	cred := graphrbac.PasswordCredential{
		StartDate: &date.Time{Time: now},
		EndDate:   &date.Time{Time: now.Add(duration)},
		KeyID:     to.StringPtr(keyID),
		Value:     to.StringPtr(password),
	}

	// Load current credentials
	resp, err := c.provider.ListApplicationPasswordCredentials(ctx, appObjID)
	if err != nil {
		return "", "", errwrap.Wrapf("error fetching credentials: {{err}}", err)
	}
	curCreds := *resp.Value

	// Add and save credentials
	curCreds = append(curCreds, cred)

	if _, err := c.provider.UpdateApplicationPasswordCredentials(ctx, appObjID,
		graphrbac.PasswordCredentialsUpdateParameters{
			Value: &curCreds,
		},
	); err != nil {
		if strings.Contains(err.Error(), "size of the object has exceeded its limit") {
			err = errors.New("maximum number of Application passwords reached")
		}
		return "", "", errwrap.Wrapf("error updating credentials: {{err}}", err)
	}

	return keyID, password, nil
}

// deleteAppPassword removes a password, if present, from an App's credentials list.
func (c *client) deleteAppPassword(ctx context.Context, appObjID, keyID string) error {
	// Load current credentials
	resp, err := c.provider.ListApplicationPasswordCredentials(ctx, appObjID)
	if err != nil {
		return errwrap.Wrapf("error fetching credentials: {{err}}", err)
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
	if _, err := c.provider.UpdateApplicationPasswordCredentials(ctx, appObjID,
		graphrbac.PasswordCredentialsUpdateParameters{
			Value: &curCreds,
		},
	); err != nil {
		return errwrap.Wrapf("error updating credentials: {{err}}", err)
	}

	return nil
}

// deleteApp deletes an Azure application.
func (c *client) deleteApp(ctx context.Context, appObjectID string) error {
	resp, err := c.provider.DeleteApplication(ctx, appObjectID)

	// Don't consider it an error if the object wasn't present
	if err != nil && resp.Response != nil && resp.StatusCode == 404 {
		return nil
	}

	return err
}

// assignRoles assigns Azure roles to a service principal.
func (c *client) assignRoles(ctx context.Context, sp *graphrbac.ServicePrincipal, roles []*azureRole) ([]string, error) {
	var ids []string

	for _, role := range roles {
		assignmentID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		resultRaw, err := retry(ctx, func() (interface{}, bool, error) {
			ra, err := c.provider.CreateRoleAssignment(ctx, role.Scope, assignmentID,
				authorization.RoleAssignmentCreateParameters{
					RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
						RoleDefinitionID: to.StringPtr(role.RoleID),
						PrincipalID:      sp.ObjectID,
					},
				})

			// Propagation delays within Azure can cause this error occasionally, so don't quit on it.
			if err != nil && strings.Contains(err.Error(), "PrincipalNotFound") {
				return nil, false, nil
			}

			return to.String(ra.ID), true, err
		})

		if err != nil {
			return nil, errwrap.Wrapf("error while assigning roles: {{err}}", err)
		}

		ids = append(ids, resultRaw.(string))
	}

	return ids, nil
}

// unassignRoles deletes role assignments, if they existed.
// This is a clean-up operation that isn't essential to revocation. As such, an
// attempt is made to remove all assignments, and not return immediately if there
// is an error.
func (c *client) unassignRoles(ctx context.Context, roleIDs []string) error {
	var merr *multierror.Error

	for _, id := range roleIDs {
		if _, err := c.provider.DeleteRoleAssignmentByID(ctx, id); err != nil {
			merr = multierror.Append(merr, errwrap.Wrapf("error unassigning role: {{err}}", err))
		}
	}

	return merr.ErrorOrNil()
}

// search for roles by name
func (c *client) findRoles(ctx context.Context, roleName string) ([]authorization.RoleDefinition, error) {
	return c.provider.ListRoles(ctx, fmt.Sprintf("subscriptions/%s", c.settings.SubscriptionID), fmt.Sprintf("roleName eq '%s'", roleName))
}

// clientSettings is used by a client to configure the connections to Azure.
// It is created from a combination of Vault config settings and environment variables.
type clientSettings struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
	Environment    azure.Environment
	PluginEnv      *logical.PluginEnvironment
}

// getClientSettings creates a new clientSettings object.
// Environment variables have higher precedence than stored configuration.
func (b *azureSecretBackend) getClientSettings(ctx context.Context, config *azureConfig) (*clientSettings, error) {
	firstAvailable := func(opts ...string) string {
		for _, s := range opts {
			if s != "" {
				return s
			}
		}
		return ""
	}

	settings := new(clientSettings)

	settings.ClientID = firstAvailable(os.Getenv("AZURE_CLIENT_ID"), config.ClientID)
	settings.ClientSecret = firstAvailable(os.Getenv("AZURE_CLIENT_SECRET"), config.ClientSecret)
	settings.SubscriptionID = firstAvailable(os.Getenv("AZURE_SUBSCRIPTION_ID"))

	settings.SubscriptionID = firstAvailable(os.Getenv("AZURE_SUBSCRIPTION_ID"), config.SubscriptionID)
	if settings.SubscriptionID == "" {
		return nil, errors.New("subscription_id is required")
	}

	settings.TenantID = firstAvailable(os.Getenv("AZURE_TENANT_ID"), config.TenantID)
	if settings.TenantID == "" {
		return nil, errors.New("tenant_id is required")
	}

	envName := firstAvailable(os.Getenv("AZURE_ENVIRONMENT"), config.Environment, "AZUREPUBLICCLOUD")
	env, err := azure.EnvironmentFromName(envName)
	if err != nil {
		return nil, err
	}
	settings.Environment = env

	pluginEnv, err := b.System().PluginEnv(ctx)
	if err != nil {
		return nil, errwrap.Wrapf("error loading plugin environment: {{err}}", err)
	}
	settings.PluginEnv = pluginEnv

	return settings, nil
}

// retry will repeatedly call f until one of:
//
//   * f returns true
//   * the context is cancelled
//   * 80 seconds elapses. Vault's default request timeout is 90s; we want to expire before then.
//
// Delays are random but will average 5 seconds.
func retry(ctx context.Context, f func() (interface{}, bool, error)) (interface{}, error) {
	delayTimer := time.NewTimer(0)
	endCh := time.NewTimer(retryTimeout).C

	for {
		if result, done, err := f(); done {
			return result, err
		}

		delay := time.Duration(2000+rand.Intn(6000)) * time.Millisecond
		delayTimer.Reset(delay)

		select {
		case <-delayTimer.C:
		case <-endCh:
			return nil, errors.New("retry: timeout")
		case <-ctx.Done():
			return nil, errors.New("retry: cancelled")
		}
	}
}
