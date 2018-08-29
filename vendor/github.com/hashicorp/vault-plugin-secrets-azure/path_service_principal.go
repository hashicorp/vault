package azuresecrets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	SecretTypeSP = "service_principal"
)

func secretServicePrincipal(b *azureSecretBackend) *framework.Secret {
	return &framework.Secret{
		Type:   SecretTypeSP,
		Renew:  b.spRenew,
		Revoke: b.spRevoke,
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

// pathSPRead generates Azure an service principal and credentials.
//
// This is a multistep process of:
//   1. Create an Azure application
//   2. Create a service principal associated with the new App
//   3. Assign roles
func (b *azureSecretBackend) pathSPRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	roleName := d.Get("role").(string)
	role, err := getRole(ctx, roleName, req.Storage)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role '%s' does not exists", roleName)), nil
	}

	// Create the App, which is the top level object to be tracked in the secret
	// and deleted upon revocation. If any subsequent step fails, the App is deleted.
	app, err := c.createApp(ctx)
	if err != nil {
		return nil, err
	}
	appID := to.String(app.AppID)
	appObjID := to.String(app.ObjectID)

	// Create the SP. A far future credential expiration is set on the Azure side.
	sp, password, err := c.createSP(ctx, app, 10*365*24*time.Hour)
	if err != nil {
		c.deleteApp(ctx, appObjID)
		return nil, err
	}

	raIDs, err := c.assignRoles(ctx, sp, role.AzureRoles)
	if err != nil {
		c.deleteApp(ctx, appObjID)
		return nil, err
	}

	resp := b.Secret(SecretTypeSP).Response(map[string]interface{}{
		"client_id":     appID,
		"client_secret": password,
	}, map[string]interface{}{
		"app_object_id":       appObjID,
		"role_assignment_ids": raIDs,
		"role":                roleName,
	})

	resp.Secret.TTL = role.TTL
	resp.Secret.MaxTTL = role.MaxTTL

	return resp, nil
}

func (b *azureSecretBackend) spRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return nil, errors.New("internal data not found")
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
		return nil, errors.New("internal data not found")
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

const pathServicePrincipalHelpSyn = `
Request Service Principal credentials for a given Vault role.
`

const pathServicePrincipalHelpDesc = `
This path creates a Service Principal and assigns Azure roles for a
given Vault role, returning the associated login credentials. The
Service Principal will be automatically deleted when the lease has expired.
`
