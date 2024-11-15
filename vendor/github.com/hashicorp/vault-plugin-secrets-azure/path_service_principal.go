// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresecrets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	SecretTypeSP       = "service_principal"
	SecretTypeStaticSP = "static_service_principal"
)

// SPs will be created with a far-future expiration in Azure
// unless `explicit_max_ttl` is set in the role
var spExpiration = 10 * 365 * 24 * time.Hour

func secretServicePrincipal(b *azureSecretBackend) *framework.Secret {
	return &framework.Secret{
		Type:   SecretTypeSP,
		Renew:  b.spRenew,
		Revoke: b.spRevoke,
	}
}

func secretStaticServicePrincipal(b *azureSecretBackend) *framework.Secret {
	return &framework.Secret{
		Type:   SecretTypeStaticSP,
		Renew:  b.spRenew,
		Revoke: b.staticSPRevoke,
	}
}

func pathServicePrincipal(b *azureSecretBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("creds/%s", framework.GenericNameRegex("role")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAzure,
			OperationVerb:   "request",
			OperationSuffix: "service-principal-credentials",
		},
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the Vault role",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.pathSPRead,
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathServicePrincipalHelpSyn,
		HelpDescription: pathServicePrincipalHelpDesc,
	}
}

// pathSPRead generates Azure credentials based on the role credential type.
func (b *azureSecretBackend) pathSPRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	roleName := d.Get("role").(string)

	role, err := getRole(ctx, roleName, req.Storage)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role '%s' does not exist", roleName)), nil
	}

	var resp *logical.Response

	if role.ApplicationObjectID != "" {
		resp, err = b.createStaticSPSecret(ctx, client, roleName, role)
	} else {
		resp, err = b.createSPSecret(ctx, req.Storage, client, roleName, role)
	}

	if err != nil {
		return nil, err
	}

	resp.Secret.TTL = role.TTL
	resp.Secret.MaxTTL = role.MaxTTL
	if role.ExplicitMaxTTL != 0 && (role.ExplicitMaxTTL < role.MaxTTL || role.MaxTTL == 0) {
		resp.Secret.MaxTTL = role.ExplicitMaxTTL
	}
	return resp, nil
}

// createSPSecret generates a new App/Service Principal.
func (b *azureSecretBackend) createSPSecret(ctx context.Context, s logical.Storage, c *client, roleName string, role *roleEntry) (*logical.Response, error) {
	// Create the App, which is the top level object to be tracked in the secret
	// and deleted upon revocation. If any subsequent step fails, the App will be
	// deleted as part of WAL rollback.
	app, err := c.createApp(ctx, role.SignInAudience, role.Tags)
	if err != nil {
		return nil, err
	}
	appID := app.AppID
	appObjID := app.AppObjectID

	// Write a WAL entry in case the SP create process doesn't complete
	walID, err := framework.PutWAL(ctx, s, walAppKey, &walApp{
		AppID:      appID,
		AppObjID:   appObjID,
		Expiration: time.Now().Add(maxWALAge),
	})
	if err != nil {
		return nil, fmt.Errorf("error writing WAL: %w", err)
	}

	// Determine SP duration
	spDuration := spExpiration
	if role.ExplicitMaxTTL != 0 {
		spDuration = role.ExplicitMaxTTL
	}

	// Create a service principal associated with the new App
	spID, password, endDate, err := c.createSP(ctx, app, spDuration)
	if err != nil {
		return nil, err
	}

	assignmentIDs, err := c.generateUUIDs(len(role.AzureRoles))
	if err != nil {
		return nil, fmt.Errorf("error generating assignment IDs; err=%w", err)
	}

	// Write a second WAL entry in case the Role assignments don't complete
	rWALID, err := framework.PutWAL(ctx, s, walAppRoleAssignment, &walAppRoleAssign{
		SpID:          spID,
		AssignmentIDs: assignmentIDs,
		AzureRoles:    role.AzureRoles,
		Expiration:    time.Now().Add(maxWALAge),
	})
	if err != nil {
		return nil, fmt.Errorf("error writing WAL: %w", err)
	}

	// Assign Azure roles to the new SP
	raIDs, err := c.assignRoles(ctx, spID, role.AzureRoles, assignmentIDs)
	if err != nil {
		return nil, err
	}

	// Assign Azure group memberships to the new SP
	if err := c.addGroupMemberships(ctx, spID, role.AzureGroups); err != nil {
		return nil, err
	}

	// SP is fully created so delete the WALs
	if err := framework.DeleteWAL(ctx, s, walID); err != nil {
		return nil, fmt.Errorf("error deleting WAL: %w", err)
	}

	if err := framework.DeleteWAL(ctx, s, rWALID); err != nil {
		return nil, fmt.Errorf("error deleting role assignment WAL: %w", err)
	}

	data := map[string]interface{}{
		"client_id":     appID,
		"client_secret": password,
	}
	internalData := map[string]interface{}{
		"app_object_id":        appObjID,
		"sp_object_id":         spID,
		"role_assignment_ids":  raIDs,
		"group_membership_ids": groupObjectIDs(role.AzureGroups),
		"role":                 roleName,
		"permanently_delete":   role.PermanentlyDelete,
		"key_end_date":         endDate.Format(time.RFC3339Nano),
	}

	return b.Secret(SecretTypeSP).Response(data, internalData), nil
}

// createStaticSPSecret adds a new password to the App associated with the role.
func (b *azureSecretBackend) createStaticSPSecret(ctx context.Context, c *client, roleName string, role *roleEntry) (*logical.Response, error) {
	lock := locksutil.LockForKey(b.appLocks, role.ApplicationObjectID)
	lock.Lock()
	defer lock.Unlock()

	// Determine SP duration
	spDuration := spExpiration
	if role.ExplicitMaxTTL != 0 {
		spDuration = role.ExplicitMaxTTL
	}

	keyID, password, endDate, err := c.addAppPassword(ctx, role.ApplicationObjectID, spDuration)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"client_id":     role.ApplicationID,
		"client_secret": password,
	}
	internalData := map[string]interface{}{
		"app_object_id": role.ApplicationObjectID,
		"key_id":        keyID,
		"key_end_date":  endDate.Format(time.RFC3339Nano),
		"role":          roleName,
	}

	return b.Secret(SecretTypeStaticSP).Response(data, internalData), nil
}

func (b *azureSecretBackend) spRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return nil, errors.New("internal data 'role' not found")
	}

	role, err := getRole(ctx, roleRaw.(string), req.Storage)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, nil
	}

	// Determine remaining lifetime of SP secret in Azure
	keyEndDateRaw, ok := req.Secret.InternalData["key_end_date"]
	if !ok {
		return nil, errors.New("internal data 'key_end_date' not found")
	}
	keyEndDate, err := time.Parse(time.RFC3339Nano, keyEndDateRaw.(string))
	if err != nil {
		return nil, errors.New("cannot parse 'key_end_date' to timestamp")
	}
	keyLifetime := time.Until(keyEndDate)

	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = min(role.TTL, keyLifetime)
	resp.Secret.MaxTTL = min(role.MaxTTL, keyLifetime)
	resp.Secret.Renewable = role.TTL < keyLifetime // Lease cannot be renewed beyond service-side endDate

	return resp, nil
}

func (b *azureSecretBackend) spRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := new(logical.Response)

	appObjectIDRaw, ok := req.Secret.InternalData["app_object_id"]
	if !ok {
		return nil, errors.New("internal data 'app_object_id' not found")
	}

	appObjectID := appObjectIDRaw.(string)

	// Get the service principal object ID. Only set if using dynamic service
	// principals.
	var spObjectID string
	if spObjectIDRaw, ok := req.Secret.InternalData["sp_object_id"]; ok {
		spObjectID = spObjectIDRaw.(string)
	}

	// Get the permanently delete setting. Only set if using dynamic service
	// principals.
	var permanentlyDelete bool
	if permanentlyDeleteRaw, ok := req.Secret.InternalData["permanently_delete"]; ok {
		permanentlyDelete = permanentlyDeleteRaw.(bool)
	}

	var raIDs []string
	if req.Secret.InternalData["role_assignment_ids"] != nil {
		for _, v := range req.Secret.InternalData["role_assignment_ids"].([]interface{}) {
			raIDs = append(raIDs, v.(string))
		}
	}

	var gmIDs []string
	if req.Secret.InternalData["group_membership_ids"] != nil {
		for _, v := range req.Secret.InternalData["group_membership_ids"].([]interface{}) {
			gmIDs = append(gmIDs, v.(string))
		}
	}

	if len(gmIDs) != 0 && spObjectID == "" {
		return nil, errors.New("internal data 'sp_object_id' not found")
	}

	c, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error during revoke: %w", err)
	}

	// unassigning roles is effectively a garbage collection operation. Errors will be noted but won't fail the
	// revocation process. Deleting the app, however, *is* required to consider the secret revoked.
	if err := c.unassignRoles(ctx, raIDs); err != nil {
		resp.AddWarning(err.Error())
	}

	// removing group membership is effectively a garbage collection
	// operation. Errors will be noted but won't fail the revocation process.
	// Deleting the app, however, *is* required to consider the secret revoked.
	if err := c.removeGroupMemberships(ctx, spObjectID, gmIDs); err != nil {
		resp.AddWarning(err.Error())
	}

	// removing the service principal is effectively a garbage collection
	// operation. Errors will be noted but won't fail the revocation process.
	// Deleting the app, however, *is* required to consider the secret revoked.
	if err := c.deleteServicePrincipal(ctx, spObjectID, permanentlyDelete); err != nil {
		resp.AddWarning(err.Error())
	}

	err = c.deleteApp(ctx, appObjectID, permanentlyDelete)
	return resp, err
}

func (b *azureSecretBackend) staticSPRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	appObjectIDRaw, ok := req.Secret.InternalData["app_object_id"]
	if !ok {
		return nil, errors.New("internal data 'app_object_id' not found")
	}

	appObjectID := appObjectIDRaw.(string)

	c, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error during revoke: %w", err)
	}

	keyIDRaw, ok := req.Secret.InternalData["key_id"]
	if !ok {
		return nil, errors.New("internal data 'key_id' not found")
	}

	lock := locksutil.LockForKey(b.appLocks, appObjectID)
	lock.Lock()
	defer lock.Unlock()

	return nil, c.deleteAppPassword(ctx, appObjectID, keyIDRaw.(string))
}

const pathServicePrincipalHelpSyn = `
Request Service Principal credentials for a given Vault role.
`

const pathServicePrincipalHelpDesc = `
This path creates or updates dynamic Service Principal credentials.
The associated role can be configured to create a new App/Service Principal,
or add a new password to an existing App. The Service Principal or password
will be automatically deleted when the lease has expired.
`
