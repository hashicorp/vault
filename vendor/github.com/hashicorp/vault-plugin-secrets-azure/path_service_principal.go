package azuresecrets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	SecretTypeSP       = "service_principal"
	SecretTypeStaticSP = "static_service_principal"
)

// SPs will be created with a far-future expiration in Azure
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
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the Vault role",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathSPRead,
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
		return logical.ErrorResponse(fmt.Sprintf("role '%s' does not exists", roleName)), nil
	}

	var resp *logical.Response

	if role.ApplicationObjectID != "" {
		resp, err = b.createStaticSPSecret(ctx, client, roleName, role)
	} else {
		resp, err = b.createSPSecret(ctx, client, roleName, role)
	}

	if err != nil {
		return nil, err
	}

	resp.Secret.TTL = role.TTL
	resp.Secret.MaxTTL = role.MaxTTL

	return resp, nil
}

// createSPSecret generates a new App/Service Principal.
func (b *azureSecretBackend) createSPSecret(ctx context.Context, c *client, roleName string, role *roleEntry) (*logical.Response, error) {
	// Create the App, which is the top level object to be tracked in the secret
	// and deleted upon revocation. If any subsequent step fails, the App is deleted.
	app, err := c.createApp(ctx)
	if err != nil {
		return nil, err
	}
	appID := to.String(app.AppID)
	appObjID := to.String(app.ObjectID)

	// Create a service principal associated with the new App
	sp, password, err := c.createSP(ctx, app, spExpiration)
	if err != nil {
		c.deleteApp(ctx, appObjID)
		return nil, err
	}

	// Assign Azure roles to the new SP
	raIDs, err := c.assignRoles(ctx, sp, role.AzureRoles)
	if err != nil {
		c.deleteApp(ctx, appObjID)
		return nil, err
	}

	data := map[string]interface{}{
		"client_id":     appID,
		"client_secret": password,
	}
	internalData := map[string]interface{}{
		"app_object_id":       appObjID,
		"role_assignment_ids": raIDs,
		"role":                roleName,
	}

	return b.Secret(SecretTypeSP).Response(data, internalData), nil
}

// createStaticSPSecret adds a new password to the App associated with the role.
func (b *azureSecretBackend) createStaticSPSecret(ctx context.Context, c *client, roleName string, role *roleEntry) (*logical.Response, error) {
	lock := locksutil.LockForKey(b.appLocks, role.ApplicationObjectID)
	lock.Lock()
	defer lock.Unlock()

	keyID, password, err := c.addAppPassword(ctx, role.ApplicationObjectID, spExpiration)
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

	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = role.TTL
	resp.Secret.MaxTTL = role.MaxTTL

	return resp, nil
}

func (b *azureSecretBackend) spRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := new(logical.Response)

	appObjectIDRaw, ok := req.Secret.InternalData["app_object_id"]
	if !ok {
		return nil, errors.New("internal data 'app_object_id' not found")
	}

	appObjectID := appObjectIDRaw.(string)

	var raIDs []string
	if req.Secret.InternalData["role_assignment_ids"] != nil {
		for _, v := range req.Secret.InternalData["role_assignment_ids"].([]interface{}) {
			raIDs = append(raIDs, v.(string))
		}
	}

	c, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, errwrap.Wrapf("error during revoke: {{err}}", err)
	}

	// unassigning roles is effectively a garbage collection operation. Errors will be noted but won't fail the
	// revocation process. Deleting the app, however, *is* required to consider the secret revoked.
	if err := c.unassignRoles(ctx, raIDs); err != nil {
		resp.AddWarning(err.Error())
	}

	err = c.deleteApp(ctx, appObjectID)

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
		return nil, errwrap.Wrapf("error during revoke: {{err}}", err)
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
