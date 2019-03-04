package seal

import (
	"context"

	"github.com/hashicorp/vault/physical"
)

const (
	Shamir        = "shamir"
	PKCS11        = "pkcs11"
	AliCloudKMS   = "alicloudkms"
	AWSKMS        = "awskms"
	GCPCKMS       = "gcpckms"
	AzureKeyVault = "azurekeyvault"
	Transit       = "transit"
	Test          = "test-auto"

	// HSMAutoDeprecated is a deprecated seal type prior to 0.9.0.
	// It is still referenced in certain code paths for upgrade purporses
	HSMAutoDeprecated = "hsm-auto"
)

// Access is the embedded implemention of autoSeal that contains logic
// specific to encrypting and decrypting data, or in this case keys.
type Access interface {
	SealType() string
	KeyID() string

	Init(context.Context) error
	Finalize(context.Context) error

	Encrypt(context.Context, []byte) (*physical.EncryptedBlobInfo, error)
	Decrypt(context.Context, *physical.EncryptedBlobInfo) ([]byte, error)
}
