package azurekeyvault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"

	wrapping "github.com/hashicorp/go-kms-wrapping"
)

const (
	EnvAzureKeyVaultWrapperVaultName = "AZUREKEYVAULT_WRAPPER_VAULT_NAME"
	EnvVaultAzureKeyVaultVaultName   = "VAULT_AZUREKEYVAULT_VAULT_NAME"

	EnvAzureKeyVaultWrapperKeyName = "AZUREKEYVAULT_WRAPPER_KEY_NAME"
	EnvVaultAzureKeyVaultKeyName   = "VAULT_AZUREKEYVAULT_KEY_NAME"
)

// Wrapper is an Wrapper that uses Azure Key Vault
// for crypto operations.  Azure Key Vault currently does not support
// keys that can encrypt long data (RSA keys).  Due to this fact, we generate
// and AES key and wrap the key using Key Vault and store it with the
// data
type Wrapper struct {
	tenantID     string
	clientID     string
	clientSecret string
	vaultName    string
	keyName      string

	currentKeyID *atomic.Value

	environment azure.Environment
	client      *keyvault.BaseClient
}

// Ensure that we are implementing Wrapper
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new wrapper with the given options
func NewWrapper(opts *wrapping.WrapperOptions) *Wrapper {
	if opts == nil {
		opts = new(wrapping.WrapperOptions)
	}
	v := &Wrapper{
		currentKeyID: new(atomic.Value),
	}
	v.currentKeyID.Store("")
	return v
}

// SetConfig sets the fields on the Wrapper object based on
// values from the config parameter.
//
// Order of precedence:
// * Environment variable
// * Value from Vault configuration file
// * Managed Service Identity for instance
func (v *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	switch {
	case os.Getenv("AZURE_TENANT_ID") != "":
		v.tenantID = os.Getenv("AZURE_TENANT_ID")
	case config["tenant_id"] != "":
		v.tenantID = config["tenant_id"]
	}

	switch {
	case os.Getenv("AZURE_CLIENT_ID") != "":
		v.clientID = os.Getenv("AZURE_CLIENT_ID")
	case config["client_id"] != "":
		v.clientID = config["client_id"]
	}

	switch {
	case os.Getenv("AZURE_CLIENT_SECRET") != "":
		v.clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	case config["client_secret"] != "":
		v.clientSecret = config["client_secret"]
	}

	envName := os.Getenv("AZURE_ENVIRONMENT")
	if envName == "" {
		envName = config["environment"]
	}
	if envName == "" {
		v.environment = azure.PublicCloud
	} else {
		var err error
		v.environment, err = azure.EnvironmentFromName(envName)
		if err != nil {
			return nil, err
		}
	}

	switch {
	case os.Getenv(EnvAzureKeyVaultWrapperVaultName) != "":
		v.vaultName = os.Getenv(EnvAzureKeyVaultWrapperVaultName)
	case os.Getenv(EnvVaultAzureKeyVaultVaultName) != "":
		v.vaultName = os.Getenv(EnvVaultAzureKeyVaultVaultName)
	case config["vault_name"] != "":
		v.vaultName = config["vault_name"]
	default:
		return nil, errors.New("vault name is required")
	}

	switch {
	case os.Getenv(EnvAzureKeyVaultWrapperKeyName) != "":
		v.keyName = os.Getenv(EnvAzureKeyVaultWrapperKeyName)
	case os.Getenv(EnvVaultAzureKeyVaultKeyName) != "":
		v.keyName = os.Getenv(EnvVaultAzureKeyVaultKeyName)
	case config["key_name"] != "":
		v.keyName = config["key_name"]
	default:
		return nil, errors.New("key name is required")
	}

	if v.client == nil {
		client, err := v.getKeyVaultClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing Azure Key Vault wrapper client: %w", err)
		}

		// Test the client connection using provided key ID
		keyInfo, err := client.GetKey(context.Background(), v.buildBaseURL(), v.keyName, "")
		if err != nil {
			return nil, fmt.Errorf("error fetching Azure Key Vault wrapper key information: %w", err)
		}
		if keyInfo.Key == nil {
			return nil, errors.New("no key information returned")
		}
		v.currentKeyID.Store(parseKeyVersion(to.String(keyInfo.Key.Kid)))

		v.client = client
	}

	// Map that holds non-sensitive configuration info
	wrapperInfo := make(map[string]string)
	wrapperInfo["environment"] = v.environment.Name
	wrapperInfo["vault_name"] = v.vaultName
	wrapperInfo["key_name"] = v.keyName

	return wrapperInfo, nil
}

// Init is called during core.Initialize.  This is a no-op.
func (v *Wrapper) Init(context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op.
func (v *Wrapper) Finalize(context.Context) error {
	return nil
}

// Type returns the type for this particular Wrapper implementation
func (v *Wrapper) Type() string {
	return wrapping.AzureKeyVault
}

// KeyID returns the last known key id
func (v *Wrapper) KeyID() string {
	return v.currentKeyID.Load().(string)
}

// HMACKeyID returns the last known HMAC key id
func (v *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt using Azure Key Vault.
// This returns the ciphertext, and/or any errors from this
// call.
func (v *Wrapper) Encrypt(ctx context.Context, plaintext, aad []byte) (blob *wrapping.EncryptedBlobInfo, err error) {
	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := wrapping.NewEnvelope(nil).Encrypt(plaintext, aad)
	if err != nil {
		return nil, fmt.Errorf("error wrapping dat: %w", err)
	}

	// Encrypt the DEK using Key Vault
	params := keyvault.KeyOperationsParameters{
		Algorithm: keyvault.RSAOAEP256,
		Value:     to.StringPtr(base64.URLEncoding.EncodeToString(env.Key)),
	}
	// Wrap key with the latest version for the key name
	resp, err := v.client.WrapKey(ctx, v.buildBaseURL(), v.keyName, "", params)
	if err != nil {
		return nil, err
	}

	// Store the current key version
	keyVersion := parseKeyVersion(to.String(resp.Kid))
	v.currentKeyID.Store(keyVersion)

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &wrapping.KeyInfo{
			KeyID:      keyVersion,
			WrappedKey: []byte(to.String(resp.Result)),
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext
func (v *Wrapper) Decrypt(ctx context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) (pt []byte, err error) {
	if in == nil {
		return nil, errors.New("given input for decryption is nil")
	}

	if in.KeyInfo == nil {
		return nil, errors.New("key info is nil")
	}

	// Unwrap the key
	params := keyvault.KeyOperationsParameters{
		Algorithm: keyvault.RSAOAEP256,
		Value:     to.StringPtr(string(in.KeyInfo.WrappedKey)),
	}
	resp, err := v.client.UnwrapKey(ctx, v.buildBaseURL(), v.keyName, in.KeyInfo.KeyID, params)
	if err != nil {
		return nil, err
	}

	keyBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(to.String(resp.Result))
	if err != nil {
		return nil, err
	}
	envInfo := &wrapping.EnvelopeInfo{
		Key:        keyBytes,
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}
	return wrapping.NewEnvelope(nil).Decrypt(envInfo, aad)
}

func (v *Wrapper) buildBaseURL() string {
	return fmt.Sprintf("https://%s.%s/", v.vaultName, v.environment.KeyVaultDNSSuffix)
}

func (v *Wrapper) getKeyVaultClient() (*keyvault.BaseClient, error) {
	var authorizer autorest.Authorizer
	var err error

	switch {
	case v.clientID != "" && v.clientSecret != "":
		config := auth.NewClientCredentialsConfig(v.clientID, v.clientSecret, v.tenantID)
		config.AADEndpoint = v.environment.ActiveDirectoryEndpoint
		config.Resource = strings.TrimSuffix(v.environment.KeyVaultEndpoint, "/")
		authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	// By default use MSI
	default:
		config := auth.NewMSIConfig()
		config.Resource = strings.TrimSuffix(v.environment.KeyVaultEndpoint, "/")
		authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	}

	client := keyvault.New()
	client.Authorizer = authorizer
	return &client, nil
}

// Kid gets returned as a full URL, get the last bit which is just
// the version
func parseKeyVersion(kid string) string {
	keyVersionParts := strings.Split(kid, "/")
	return keyVersionParts[len(keyVersionParts)-1]
}
