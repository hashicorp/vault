package gcpsecrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"google.golang.org/api/iam/v1"
)

const (
	SecretTypeKey      = "service_account_key"
	keyAlgorithmRSA2k  = "KEY_ALG_RSA_2048"
	privateKeyTypeJson = "TYPE_GOOGLE_CREDENTIALS_FILE"
)

func secretServiceAccountKey(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTypeKey,
		Fields: map[string]*framework.FieldSchema{
			"private_key_data": {
				Type:        framework.TypeString,
				Description: "Base-64 encoded string. Private key data for a service account key",
			},
			"key_algorithm": {
				Type:        framework.TypeString,
				Description: "Which type of key and algorithm to use for the key (defaults to 2K RSA). Valid values are GCP enum(ServiceAccountKeyAlgorithm)",
			},
			"key_type": {
				Type:        framework.TypeString,
				Description: "Type of the private key (i.e. whether it is JSON or P12). Valid values are GCP enum(ServiceAccountPrivateKeyType)",
			},
		},

		Renew:  b.secretKeyRenew,
		Revoke: secretKeyRevoke,
	}
}

func pathSecretServiceAccountKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("key/%s", framework.GenericNameRegex("roleset")),
		Fields: map[string]*framework.FieldSchema{
			"roleset": {
				Type:        framework.TypeString,
				Description: "Required. Name of the role set.",
			},
			"key_algorithm": {
				Type:        framework.TypeString,
				Description: fmt.Sprintf(`Private key algorithm for service account key - defaults to %s"`, keyAlgorithmRSA2k),
				Default:     keyAlgorithmRSA2k,
			},
			"key_type": {
				Type:        framework.TypeString,
				Description: fmt.Sprintf(`Private key type for service account key - defaults to %s"`, privateKeyTypeJson),
				Default:     privateKeyTypeJson,
			},
		},
		ExistenceCheck: b.pathRoleSetExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathServiceAccountKey,
			logical.UpdateOperation: b.pathServiceAccountKey,
		},
		HelpSynopsis:    pathServiceAccountKeySyn,
		HelpDescription: pathServiceAccountKeyDesc,
	}
}

func (b *backend) pathServiceAccountKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rsName := d.Get("roleset").(string)
	keyType := d.Get("key_type").(string)
	keyAlg := d.Get("key_algorithm").(string)

	rs, err := getRoleSet(rsName, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return logical.ErrorResponse(fmt.Sprintf("role set '%s' does not exists", rsName)), nil
	}

	if rs.SecretType != SecretTypeKey {
		return logical.ErrorResponse(fmt.Sprintf("role set '%s' cannot generate service account keys (has secret type %s)", rsName, rs.SecretType)), nil
	}

	return b.getSecretKey(ctx, req.Storage, rs, keyType, keyAlg)
}

func (b *backend) secretKeyRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp, err := b.verifySecretServiceKeyExists(ctx, req)
	if err != nil {
		return resp, err
	}
	if resp == nil {
		resp = &logical.Response{}
	}
	cfg, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &config{}
	}

	resp.Secret = req.Secret
	resp.Secret.TTL = cfg.TTL
	resp.Secret.MaxTTL = cfg.MaxTTL
	return resp, nil
}

func (b *backend) verifySecretServiceKeyExists(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	keyName, ok := req.Secret.InternalData["key_name"]
	if !ok {
		return nil, fmt.Errorf("invalid secret, internal data is missing key name")
	}

	rsName, ok := req.Secret.InternalData["role_set"]
	if !ok {
		return nil, fmt.Errorf("invalid secret, internal data is missing role set name")
	}

	bindingSum, ok := req.Secret.InternalData["role_set_bindings"]
	if !ok {
		return nil, fmt.Errorf("invalid secret, internal data is missing role set checksum")
	}

	// Verify role set was not deleted.
	rs, err := getRoleSet(rsName.(string), ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("could not find role set '%v' for secret", rsName)), nil
	}

	// Verify role set bindings have not changed since secret was generated.
	if rs.bindingHash() != bindingSum.(string) {
		return logical.ErrorResponse(fmt.Sprintf("role set '%v' bindings were updated since secret was generated, cannot renew", rsName)), nil
	}

	// Verify service account key still exists.
	iamAdmin, err := newIamAdmin(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse("could not confirm key still exists in GCP"), nil
	}
	if k, err := iamAdmin.Projects.ServiceAccounts.Keys.Get(keyName.(string)).Do(); err != nil || k == nil {
		return logical.ErrorResponse(fmt.Sprintf("could not confirm key still exists in GCP: %v", err)), nil
	}
	return nil, nil
}

func secretKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyNameRaw, ok := req.Secret.InternalData["key_name"]
	if !ok {
		return nil, fmt.Errorf("secret is missing key_name internal data")
	}

	iamAdmin, err := newIamAdmin(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	_, err = iamAdmin.Projects.ServiceAccounts.Keys.Delete(keyNameRaw.(string)).Do()
	if err != nil && !isGoogleApi404Error(err) {
		return logical.ErrorResponse(fmt.Sprintf("unable to delete service account key: %v", err)), nil
	}

	return nil, nil
}

func (b *backend) getSecretKey(ctx context.Context, s logical.Storage, rs *RoleSet, keyType, keyAlgorithm string) (*logical.Response, error) {
	cfg, err := getConfig(ctx, s)
	if err != nil {
		return nil, errwrap.Wrapf("could not read backend config: {{err}}", err)
	}
	if cfg == nil {
		cfg = &config{}
	}

	iamC, err := newIamAdmin(ctx, s)
	if err != nil {
		return nil, errwrap.Wrapf("could not create IAM Admin client: {{err}}", err)
	}

	account, err := rs.getServiceAccount(iamC)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("roleset service account was removed - role set must be updated (write to roleset/%s/rotate) before generating new secrets", rs.Name)), nil
	}

	key, err := iamC.Projects.ServiceAccounts.Keys.Create(
		account.Name, &iam.CreateServiceAccountKeyRequest{
			KeyAlgorithm:   keyAlgorithm,
			PrivateKeyType: keyType,
		}).Do()
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	secretD := map[string]interface{}{
		"private_key_data": key.PrivateKeyData,
		"key_algorithm":    key.KeyAlgorithm,
		"key_type":         key.PrivateKeyType,
	}
	internalD := map[string]interface{}{
		"key_name":          key.Name,
		"role_set":          rs.Name,
		"role_set_bindings": rs.bindingHash(),
	}

	resp := b.Secret(SecretTypeKey).Response(secretD, internalD)
	resp.Secret.TTL = cfg.TTL
	resp.Secret.MaxTTL = cfg.MaxTTL
	resp.Secret.Renewable = true
	return resp, nil
}

const pathTokenHelpSyn = `Generate an OAuth2 access token under a specific role set.`
const pathTokenHelpDesc = `
This path will generate a new OAuth2 access token for accessing GCP APIs.
A role set, binding IAM roles to specific GCP resources, will be specified 
by name - for example, if this backend is mounted at "gcp",
then "gcp/token/deploy" would generate tokens for the "deploy" role set.

On the backend, each roleset is associated with a service account. 
The token will be associated with this service account. Tokens have a 
short-term lease (1-hour) associated with them but cannot be renewed.
`

const pathServiceAccountKeySyn = `Generate an service account private key under a specific role set.`
const pathServiceAccountKeyDesc = `
This path will generate a new service account private key for accessing GCP APIs.
A role set, binding IAM roles to specific GCP resources, will be specified 
by name - for example, if this backend is mounted at "gcp", then "gcp/key/deploy" 
would generate service account keys for the "deploy" role set.

On the backend, each roleset is associated with a service account under
which secrets/keys are created.
`
