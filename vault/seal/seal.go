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
	Test          = "test-auto"
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
