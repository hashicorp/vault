package gcpsecrets

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
	"time"
)

func pathSecretAccessToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("token/%s", framework.GenericNameRegex("roleset")),
		Fields: map[string]*framework.FieldSchema{
			"roleset": {
				Type:        framework.TypeString,
				Description: "Required. Name of the role set.",
			},
		},
		ExistenceCheck: b.pathRoleSetExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathAccessToken,
			logical.UpdateOperation: b.pathAccessToken,
		},
		HelpSynopsis:    pathTokenHelpSyn,
		HelpDescription: pathTokenHelpDesc,
	}
}

func (b *backend) pathAccessToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rsName := d.Get("roleset").(string)

	rs, err := getRoleSet(rsName, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return logical.ErrorResponse(fmt.Sprintf("role set '%s' does not exists", rsName)), nil
	}

	if rs.SecretType != SecretTypeAccessToken {
		return logical.ErrorResponse(fmt.Sprintf("role set '%s' cannot generate access tokens (has secret type %s)", rsName, rs.SecretType)), nil
	}

	return b.secretAccessTokenResponse(ctx, req.Storage, rs)
}

func (b *backend) secretAccessTokenResponse(ctx context.Context, s logical.Storage, rs *RoleSet) (*logical.Response, error) {
	iamC, err := newIamAdmin(ctx, s)
	if err != nil {
		return nil, errwrap.Wrapf("could not create IAM Admin client: {{err}}", err)
	}

	// Verify account still exists
	_, err = rs.getServiceAccount(iamC)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("could not get role set service account: %v", err)), nil
	}

	if rs.TokenGen == nil || rs.TokenGen.KeyName == "" {
		return logical.ErrorResponse(fmt.Sprintf("invalid role set has no service account key, must be updated (path roleset/%s/rotate-key) before generating new secrets", rs.Name)), nil
	}

	token, err := rs.TokenGen.getAccessToken(ctx, iamC)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("could not generate token: %v", err)), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token":              token.AccessToken,
			"token_ttl":          token.Expiry.UTC().Sub(time.Now().UTC()) / (time.Second),
			"expires_at_seconds": token.Expiry.Unix(),
		},
	}, nil
}

func (tg *TokenGenerator) getAccessToken(ctx context.Context, iamAdmin *iam.Service) (*oauth2.Token, error) {
	key, err := iamAdmin.Projects.ServiceAccounts.Keys.Get(tg.KeyName).Do()
	if err != nil {
		return nil, errwrap.Wrapf("could not verify key used to generate tokens: {{err}}", err)
	}
	if key == nil {
		return nil, errors.New("could not find key used to generate tokens, must update role set")
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(tg.B64KeyJSON)
	if err != nil {
		return nil, errwrap.Wrapf("could not b64-decode key data: {{err}}", err)
	}

	cfg, err := google.JWTConfigFromJSON(jsonBytes, tg.Scopes...)
	if err != nil {
		return nil, errwrap.Wrapf("could not generate token JWT config: {{err}}", err)
	}

	tkn, err := cfg.TokenSource(ctx).Token()
	if err != nil {
		return nil, errwrap.Wrapf("could not generate token: {{err}}", err)
	}
	return tkn, err
}

const deprecationWarning = `
This endpoint no longer generates leases due to limitations of the GCP API, as OAuth2 tokens belonging to Service 
Accounts cannot be revoked. This access_token and lease were created by a previous version of the GCP secrets 
engine and will be cleaned up now. Note that there is the chance that this access_token, if not already expired, 
will still be valid up to one hour. 
`

const pathTokenHelpSyn = `Generate an OAuth2 access token under a specific role set.`
const pathTokenHelpDesc = `
This path will generate a new OAuth2 access token for accessing GCP APIs.
A role set, binding IAM roles to specific GCP resources, will be specified 
by name - for example, if this backend is mounted at "gcp",
then "gcp/token/deploy" would generate tokens for the "deploy" role set.

On the backend, each roleset is associated with a service account. 
The token will be associated with this service account. Tokens have a 
short-term lease (1-hour) associated with them but cannot be renewed.

Please see backend documentation for more information: 
https://www.vaultproject.io/docs/secrets/gcp/index.html
`

// EVERYTHING USING THIS SECRET TYPE IS CURRENTLY DEPRECATED.
// We keep it to allow for clean up of access_token secrets/leases that may have be left over
// by older versions of Vault.
const SecretTypeAccessToken = "access_token"

func secretAccessToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTypeAccessToken,
		Fields: map[string]*framework.FieldSchema{
			"token": {
				Type:        framework.TypeString,
				Description: "OAuth2 token",
			},
		},
		Renew:  b.secretAccessTokenRenew,
		Revoke: b.secretAccessTokenRevoke,
	}
}

// Renewal will still return an error, but return the warning in case as well.
func (b *backend) secretAccessTokenRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := logical.ErrorResponse("short-term access tokens cannot be renewed - request new access token instead")
	resp.AddWarning(deprecationWarning)
	return resp, nil
}

// Revoke will no-op and pass but warn the user. This is mostly to clean up old leases.
// Any associated secret (access_token) has already expired and thus doesn't need to
// actually be revoked,  or will expire within an hour and currently can't actually be revoked anyways.
func (b *backend) secretAccessTokenRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := &logical.Response{}
	resp.AddWarning(deprecationWarning)
	return resp, nil
}
