package gcpauth

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault-plugin-auth-gcp/plugin/util"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const warningACLReadAccess string = "Read access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords."

func pathConfig(b *GcpAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"credentials": {
				Type: framework.TypeString,
				Description: `
Google credentials JSON that Vault will use to verify users against GCP APIs.
If not specified, will use application default credentials`,
			},
			"google_certs_endpoint": {
				Type: framework.TypeString,
				Description: `
Base endpoint url that Vault will use to get Google certificates.
If not specified, will use the OAuth2 library default. Useful for testing.`,
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

func (b *GcpAuthBackend) pathConfigWrite(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(req.Storage)

	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &gcpConfig{}
	}

	if err := config.Update(data); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("could not update config: %v", err)), nil
	}

	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *GcpAuthBackend) pathConfigRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"client_email":          config.Credentials.ClientEmail,
			"client_id":             config.Credentials.ClientId,
			"private_key_id":        config.Credentials.PrivateKeyId,
			"private_key":           config.Credentials.PrivateKey,
			"project_id":            config.Credentials.ProjectId,
			"google_certs_endpoint": config.GoogleCertsEndpoint,
		},
	}

	resp.AddWarning(warningACLReadAccess)
	return resp, nil
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
	Credentials         *util.GcpCredentials `json:"credentials" structs:"credentials" mapstructure:"credentials"`
	GoogleCertsEndpoint string               `json:"google_certs_endpoint" structs:"google_certs_endpoint" mapstructure:"google_certs_endpoint"`
}

// Update sets gcpConfig values parsed from the FieldData.
func (config *gcpConfig) Update(data *framework.FieldData) error {
	credentialsJson := data.Get("credentials").(string)
	if credentialsJson != "" {
		creds, err := util.Credentials(credentialsJson)
		if err != nil {
			return fmt.Errorf("error reading google credentials from given JSON: %v", err)
		}
		if len(creds.PrivateKeyId) == 0 {
			return errors.New("google credentials not found from given JSON")
		}
		config.Credentials = creds
	}

	certsEndpoint := data.Get("google_certs_endpoint").(string)
	if len(certsEndpoint) > 0 {
		config.GoogleCertsEndpoint = certsEndpoint
	}

	return nil
}

// config reads the backend's gcpConfig from storage.
// This assumes the caller has already obtained the backend's config lock.
func (b *GcpAuthBackend) config(s logical.Storage) (*gcpConfig, error) {
	config := &gcpConfig{}
	entry, err := s.Get("config")

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
