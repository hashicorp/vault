package azuresecrets

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/authorization/mgmt/authorization"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	multierror "github.com/hashicorp/go-multierror"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault-plugin-secrets-azure/api"
	"github.com/hashicorp/vault/sdk/logical"
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
	provider   api.AzureProvider
	settings   *clientSettings
	expiration time.Time
	passwords  api.Passwords
}

// Valid returns whether the client defined and not expired.
func (c *client) Valid() bool {
	return c != nil && time.Now().Before(c.expiration)
}

// createApp creates a new Azure application.
// An Application is a needed to create service principals used by
// the caller for authentication.
func (c *client) createApp(ctx context.Context) (app *api.ApplicationResult, err error) {
	// TODO: Make this name customizable with the same logic as username customization
	name, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	name = appNamePrefix + name

	result, err := c.provider.CreateApplication(ctx, name)

	return &result, err
}

// createSP creates a new service principal.
func (c *client) createSP(
	ctx context.Context,
	app *api.ApplicationResult,
	duration time.Duration) (spID string, password string, err error) {

	type idPass struct {
		ID       string
		Password string
	}

	resultRaw, err := retry(ctx, func() (interface{}, bool, error) {
		now := time.Now()
		spID, password, err := c.provider.CreateServicePrincipal(ctx, *app.AppID, now, now.Add(duration))

		// Propagation delays within Azure can cause this error occasionally, so don't quit on it.
		if err != nil && strings.Contains(err.Error(), "does not reference a valid application object") {
			return nil, false, nil
		}

		result := idPass{
			ID:       spID,
			Password: password,
		}

		return result, true, err
	})

	if err != nil {
		return "", "", fmt.Errorf("error creating service principal: %w", err)
	}

	result := resultRaw.(idPass)

	return result.ID, result.Password, nil
}

// addAppPassword adds a new password to an App's credentials list.
func (c *client) addAppPassword(ctx context.Context, appObjID string, expiresIn time.Duration) (string, string, error) {
	exp := time.Now().Add(expiresIn)
	resp, err := c.provider.AddApplicationPassword(ctx, appObjID, "vault-plugin-secrets-azure", exp)
	if err != nil {
		if strings.Contains(err.Error(), "size of the object has exceeded its limit") {
			err = errors.New("maximum number of Application passwords reached")
		}
		return "", "", fmt.Errorf("error updating credentials: %w", err)
	}

	return to.String(resp.KeyID), to.String(resp.SecretText), nil
}

// deleteAppPassword removes a password, if present, from an App's credentials list.
func (c *client) deleteAppPassword(ctx context.Context, appObjID, keyID string) error {
	err := c.provider.RemoveApplicationPassword(ctx, appObjID, keyID)
	if err != nil {
		if strings.Contains(err.Error(), "No password credential found with keyId") {
			return nil
		}
		return fmt.Errorf("error removing credentials: %w", err)
	}

	return nil
}

// deleteApp deletes an Azure application.
func (c *client) deleteApp(ctx context.Context, appObjectID string) error {
	return c.provider.DeleteApplication(ctx, appObjectID)
}

// assignRoles assigns Azure roles to a service principal.
func (c *client) assignRoles(ctx context.Context, spID string, roles []*AzureRole) ([]string, error) {
	var ids []string

	for _, role := range roles {
		assignmentID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		resultRaw, err := retry(ctx, func() (interface{}, bool, error) {
			ra, err := c.provider.CreateRoleAssignment(ctx, role.Scope, assignmentID,
				authorization.RoleAssignmentCreateParameters{
					Properties: &authorization.RoleAssignmentProperties{
						RoleDefinitionID: &role.RoleID,
						PrincipalID:      &spID,
					},
				})

			// Propagation delays within Azure can cause this error occasionally, so don't quit on it.
			if err != nil && strings.Contains(err.Error(), "PrincipalNotFound") {
				return nil, false, nil
			}

			return to.String(ra.ID), true, err
		})

		if err != nil {
			return nil, fmt.Errorf("error while assigning roles: %w", err)
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
			merr = multierror.Append(merr, fmt.Errorf("error unassigning role: %w", err))
		}
	}

	return merr.ErrorOrNil()
}

// addGroupMemberships adds the service principal to the Azure groups.
func (c *client) addGroupMemberships(ctx context.Context, spID string, groups []*AzureGroup) error {
	for _, group := range groups {
		_, err := retry(ctx, func() (interface{}, bool, error) {
			err := c.provider.AddGroupMember(ctx, group.ObjectID, spID)

			// Propagation delays within Azure can cause this error occasionally, so don't quit on it.
			if err != nil && strings.Contains(err.Error(), "Request_ResourceNotFound") {
				return nil, false, nil
			}

			return nil, true, err
		})

		if err != nil {
			return fmt.Errorf("error while adding group membership: %w", err)
		}
	}

	return nil
}

// removeGroupMemberships removes the passed service principal from the passed
// groups. This is a clean-up operation that isn't essential to revocation. As
// such, an attempt is made to remove all memberships, and not return
// immediately if there is an error.
func (c *client) removeGroupMemberships(ctx context.Context, servicePrincipalObjectID string, groupIDs []string) error {
	var merr *multierror.Error

	for _, id := range groupIDs {
		if err := c.provider.RemoveGroupMember(ctx, servicePrincipalObjectID, id); err != nil {
			merr = multierror.Append(merr, fmt.Errorf("error removing group membership: %w", err))
		}
	}

	return merr.ErrorOrNil()
}

// groupObjectIDs is a helper for converting a list of AzureGroup
// objects to a list of their object IDs.
func groupObjectIDs(groups []*AzureGroup) []string {
	groupIDs := make([]string, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ObjectID)

	}
	return groupIDs
}

// search for roles by name
func (c *client) findRoles(ctx context.Context, roleName string) ([]authorization.RoleDefinition, error) {
	return c.provider.ListRoleDefinitions(ctx, fmt.Sprintf("subscriptions/%s", c.settings.SubscriptionID), fmt.Sprintf("roleName eq '%s'", roleName))
}

// findGroups is used to find a group by name. It returns all groups matching
// the provided name.
func (c *client) findGroups(ctx context.Context, groupName string) ([]api.Group, error) {
	return c.provider.ListGroups(ctx, fmt.Sprintf("displayName eq '%s'", groupName))
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
		return nil, fmt.Errorf("error loading plugin environment: %w", err)
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
	if _, hasTimeout := ctx.Deadline(); !hasTimeout {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, retryTimeout)
		defer cancel()
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var lastErr error
	for {
		select {
		case <-delayTimer.C:
			result, done, err := f()
			if done {
				return result, err
			}
			lastErr = err

			delay := time.Duration(2+rng.Intn(6)) * time.Second
			delayTimer.Reset(delay)
		case <-ctx.Done():
			err := lastErr
			if err == nil {
				err = ctx.Err()
			}
			return nil, fmt.Errorf("retry failed: %w", err)
		}
	}
}
