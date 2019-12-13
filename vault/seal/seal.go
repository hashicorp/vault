package seal

import (
	"context"

	"github.com/hashicorp/vault/sdk/physical"
)

const (
	Shamir        = "shamir"
	PKCS11        = "pkcs11"
	AliCloudKMS   = "alicloudkms"
	AWSKMS        = "awskms"
	GCPCKMS       = "gcpckms"
	AzureKeyVault = "azurekeyvault"
	OCIKMS        = "ocikms"
	Transit       = "transit"
	Test          = "test-auto"

	// HSMAutoDeprecated is a deprecated seal type prior to 0.9.0.
	// It is still referenced in certain code paths for upgrade purporses
	HSMAutoDeprecated = "hsm-auto"
)

type Encryptor interface {
	Encrypt(context.Context, []byte) (*physical.EncryptedBlobInfo, error)
	Decrypt(context.Context, *physical.EncryptedBlobInfo) ([]byte, error)
}

// Access is the embedded implementation of autoSeal that contains logic
// specific to encrypting and decrypting data, or in this case keys.
type Access interface {
	SealType() string
	KeyID() string

	Init(context.Context) error
	Finalize(context.Context) error

	Encryptor
}
