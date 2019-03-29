package azurekeyvault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// AzureKeyVaultSeal is an auto-seal that uses Azure Key Vault
// for crypto operations.  Azure Key Vault currently does not support
// keys that can encrypt long data (RSA keys).  Due to this fact, we generate
// and AES key and wrap the key using Key Vault and store it with the
// data
type AzureKeyVaultSeal struct {
	tenantID     string
	clientID     string
	clientSecret string
	vaultName    string
	keyName      string

	currentKeyID *atomic.Value

	environment azure.Environment
	client      *keyvault.BaseClient

	logger log.Logger
}

// Ensure that we are implementing AutoSealAccess
var _ seal.Access = (*AzureKeyVaultSeal)(nil)

func NewSeal(logger log.Logger) *AzureKeyVaultSeal {
	v := &AzureKeyVaultSeal{
		logger:       logger,
		currentKeyID: new(atomic.Value),
	}
	v.currentKeyID.Store("")
	return v
}

// SetConfig sets the fields on the AzureKeyVaultSeal object based on
// values from the config parameter.
//
// Order of precedence:
// * Environment variable
// * Value from Vault configuration file
// * Managed Service Identity for instance
func (v *AzureKeyVaultSeal) SetConfig(config map[string]string) (map[string]string, error) {
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
	case os.Getenv("VAULT_AZUREKEYVAULT_VAULT_NAME") != "":
		v.vaultName = os.Getenv("VAULT_AZUREKEYVAULT_VAULT_NAME")
	case config["vault_name"] != "":
		v.vaultName = config["vault_name"]
	default:
		return nil, errors.New("vault name is required")
	}

	switch {
	case os.Getenv("VAULT_AZUREKEYVAULT_KEY_NAME") != "":
		v.keyName = os.Getenv("VAULT_AZUREKEYVAULT_KEY_NAME")
	case config["key_name"] != "":
		v.keyName = config["key_name"]
	default:
		return nil, errors.New("key name is required")
	}

	if v.client == nil {
		client, err := v.getKeyVaultClient()
		if err != nil {
			return nil, errwrap.Wrapf("error initializing Azure Key Vault seal client: {{err}}", err)
		}

		// Test the client connection using provided key ID
		keyInfo, err := client.GetKey(context.Background(), v.buildBaseURL(), v.keyName, "")
		if err != nil {
			return nil, errwrap.Wrapf("error fetching Azure Key Vault seal key information: {{err}}", err)
		}
		if keyInfo.Key == nil {
			return nil, errors.New("no key information returned")
		}
		v.currentKeyID.Store(parseKeyVersion(to.String(keyInfo.Key.Kid)))

		v.client = client
	}

	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)
	sealInfo["environment"] = v.environment.Name
	sealInfo["vault_name"] = v.vaultName
	sealInfo["key_name"] = v.keyName

	return sealInfo, nil
}

// Init is called during core.Initialize.  This is a no-op.
func (v *AzureKeyVaultSeal) Init(context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op.
func (v *AzureKeyVaultSeal) Finalize(context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (v *AzureKeyVaultSeal) SealType() string {
	return seal.AzureKeyVault
}

// KeyID returns the last known key id.
func (v *AzureKeyVaultSeal) KeyID() string {
	return v.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt using Azure Key Vault.
// This returns the ciphertext, and/or any errors from this
// call.
func (v *AzureKeyVaultSeal) Encrypt(ctx context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "azurekeyvault", "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "azurekeyvault", "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "azurekeyvault", "encrypt"}, 1)

	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := seal.NewEnvelope().Encrypt(plaintext)
	if err != nil {
		return nil, errwrap.Wrapf("error wrapping dat: {{err}}", err)
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

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &physical.SealKeyInfo{
			KeyID:      keyVersion,
			WrappedKey: []byte(to.String(resp.Result)),
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext.
func (v *AzureKeyVaultSeal) Decrypt(ctx context.Context, in *physical.EncryptedBlobInfo) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "azurekeyvault", "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "azurekeyvault", "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "azurekeyvault", "decrypt"}, 1)

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
	envInfo := &seal.EnvelopeInfo{
		Key:        keyBytes,
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}
	return seal.NewEnvelope().Decrypt(envInfo)
}

func (v *AzureKeyVaultSeal) buildBaseURL() string {
	return fmt.Sprintf("https://%s.%s/", v.vaultName, v.environment.KeyVaultDNSSuffix)
}

func (v *AzureKeyVaultSeal) getKeyVaultClient() (*keyvault.BaseClient, error) {
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
