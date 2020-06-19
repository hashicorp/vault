package gcpauth

import (
	"context"
	"errors"

	"encoding/json"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *GcpAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"credentials": {
				Type: framework.TypeString,
				Description: `
Google credentials JSON that Vault will use to verify users against GCP APIs.
If not specified, will use application default credentials`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Credentials",
				},
			},
			"google_certs_endpoint": {
				Type: framework.TypeString,
				Description: `
Deprecated. This field does nothing and be removed in a future release`,
				Deprecated: true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.UpdateOperation: b.pathConfigWrite,
		},

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

func (b *GcpAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, d); err != nil {
		return nil, logical.CodedError(422, err.Error())
	}

	c, err := b.config(ctx, req.Storage)

	if err != nil {
		return nil, err
	}
	if c == nil {
		c = &gcpConfig{}
	}

	changed, err := c.Update(d)
	if err != nil {
		return nil, logical.CodedError(400, err.Error())
	}

	// Only do the following if the config is different
	if changed {
		// Generate a new storage entry
		entry, err := logical.StorageEntryJSON("config", c)
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate JSON configuration: {{err}}", err)
		}

		// Save the storage entry
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, errwrap.Wrapf("failed to persist configuration to storage: {{err}}", err)
		}

		// Invalidate existing client so it reads the new configuration
		b.ClearCaches()
	}

	return nil, nil
}

func (b *GcpAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, d); err != nil {
		return nil, logical.CodedError(422, err.Error())
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	resp := make(map[string]interface{})

	if v := config.Credentials.ClientEmail; v != "" {
		resp["client_email"] = v
	}
	if v := config.Credentials.ClientId; v != "" {
		resp["client_id"] = v
	}
	if v := config.Credentials.PrivateKeyId; v != "" {
		resp["private_key_id"] = v
	}
	if v := config.Credentials.ProjectId; v != "" {
		resp["project_id"] = v
	}

	return &logical.Response{
		Data: resp,
	}, nil
}

const confHelpSyn = `Configure credentials used to query the GCP IAM API to verify authenticating service accounts`
const confHelpDesc = `
The GCP IAM auth backend makes queries to the GCP IAM auth backend to verify a service account
attempting login. It verifies the service account exists and retrieves a public key to verify
signed JWT requests passed in on login. The credentials should have the following permissions:

iam AUTH:
* iam.serviceAccountKeys.get
`

// gcpConfig contains all config required for the GCP backend.
type gcpConfig struct {
	Credentials *gcputil.GcpCredentials `json:"credentials"`
}

// standardizedCreds wraps gcputil.GcpCredentials with a type to allow
// parsing through Google libraries, since the google libraries struct is not
// exposed.
type standardizedCreds struct {
	*gcputil.GcpCredentials
	CredType string `json:"type"`
}

const serviceAccountCredsType = "service_account"

// formatAsCredentialJSON converts and marshals the config credentials
// into a parsable format by Google libraries.
func (config *gcpConfig) formatAndMarshalCredentials() ([]byte, error) {
	if config == nil || config.Credentials == nil {
		return []byte{}, nil
	}

	return json.Marshal(standardizedCreds{
		GcpCredentials: config.Credentials,
		CredType:       serviceAccountCredsType,
	})
}

// Update sets gcpConfig values parsed from the FieldData.
func (c *gcpConfig) Update(d *framework.FieldData) (bool, error) {
	if d == nil {
		return false, nil
	}

	changed := false

	if v, ok := d.GetOk("credentials"); ok {
		creds, err := gcputil.Credentials(v.(string))
		if err != nil {
			return false, errwrap.Wrapf("failed to read credentials: {{err}}", err)
		}

		if len(creds.PrivateKeyId) == 0 {
			return false, errors.New("missing private key in credentials")
		}

		c.Credentials = creds
		changed = true
	}

	return changed, nil
}

// config reads the backend's gcpConfig from storage.
// This assumes the caller has already obtained the backend's config lock.
func (b *GcpAuthBackend) config(ctx context.Context, s logical.Storage) (*gcpConfig, error) {
	config := &gcpConfig{}
	entry, err := s.Get(ctx, "config")

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}
