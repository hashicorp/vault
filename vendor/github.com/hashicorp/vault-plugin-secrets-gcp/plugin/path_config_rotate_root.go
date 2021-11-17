package gcpsecrets

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/iam/v1"
)

func pathConfigRotateRoot(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/rotate-root",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigRotateRootWrite,
			},
		},

		HelpSynopsis:    pathConfigRotateRootHelpSyn,
		HelpDescription: pathConfigRotateRootHelpDesc,
	}
}

func (b *backend) pathConfigRotateRootWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the current configuration
	cfg, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("no configuration")
	}
	if cfg.CredentialsRaw == "" {
		return nil, fmt.Errorf("configuration does not have credentials - this " +
			"endpoint only works with user-provided JSON credentials explicitly " +
			"provided via the config/ endpoint")
	}

	// Parse the credential JSON to extract the email (we need it for the API
	// call)
	creds, err := gcputil.Credentials(cfg.CredentialsRaw)
	if err != nil {
		return nil, errwrap.Wrapf("credentials are invalid: {{err}}", err)
	}

	// Generate a new service account key
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create iam client: {{err}}", err)
	}

	saName := "projects/-/serviceAccounts/" + creds.ClientEmail
	newKey, err := iamAdmin.Projects.ServiceAccounts.Keys.
		Create(saName, &iam.CreateServiceAccountKeyRequest{
			KeyAlgorithm:   keyAlgorithmRSA2k,
			PrivateKeyType: privateKeyTypeJson,
		}).
		Context(ctx).
		Do()
	if err != nil {
		return nil, errwrap.Wrapf("failed to create new key: {{err}}", err)
	}

	// Base64-decode the private key data (it's the JSON file)
	newCredsJSON, err := base64.StdEncoding.DecodeString(newKey.PrivateKeyData)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode credentials: {{err}}", err)
	}

	// Verify creds are valid
	newCreds, err := gcputil.Credentials(string(newCredsJSON))
	if err != nil {
		return nil, errwrap.Wrapf("api returned invalid credentials: {{err}}", err)
	}

	// Update the configuration
	cfg.CredentialsRaw = string(newCredsJSON)
	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, errwrap.Wrapf("failed to generate new configuration: {{err}}", err)
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, errwrap.Wrapf("failed to save new configuration: {{err}}", err)
	}

	// Clear caches to pick up the new credentials
	b.ClearCaches()

	// Delete the old service account key
	oldKeyName := fmt.Sprintf("projects/%s/serviceAccounts/%s/keys/%s",
		creds.ProjectId,
		creds.ClientEmail,
		creds.PrivateKeyId)
	if _, err := iamAdmin.Projects.ServiceAccounts.Keys.
		Delete(oldKeyName).
		Context(ctx).
		Do(); err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf(
			"failed to delete old service account key (%q) - the new service "+
				"account key (%q) is active, but the old one still exists: {{err}}",
			creds.PrivateKeyId, newCreds.PrivateKeyId), err)
	}

	// We did it!
	return &logical.Response{
		Data: map[string]interface{}{
			"private_key_id": newCreds.PrivateKeyId,
		},
	}, nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the GCP credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the GCP service account credentials used by Vault
for this mount. It does this by generating a new key for the service account,
replacing the internal value, and then scheduling a deletion of the old service
account key. Note that it does not create a new service account, only a new
version of the service account key.

This path is only valid if Vault has been configured to use GCP credentials via
the config/ endpoint where "credentials" were specified. Additionally, the
provided service account must have permissions to create and delete service
account keys.
`
