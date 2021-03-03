package wrapping

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

// These values define known types of Wrappers
const (
	AEAD            = "aead"
	AliCloudKMS     = "alicloudkms"
	AWSKMS          = "awskms"
	AzureKeyVault   = "azurekeyvault"
	GCPCKMS         = "gcpckms"
	HuaweiCloudKMS  = "huaweicloudkms"
	MultiWrapper    = "multiwrapper"
	OCIKMS          = "ocikms"
	PKCS11          = "pkcs11"
	Shamir          = "shamir"
	TencentCloudKMS = "tencentcloudkms"
	Transit         = "transit"
	YandexCloudKMS  = "yandexcloudkms"
	Test            = "test-auto"

	// HSMAutoDeprecated is a deprecated type relevant to Vault prior to 0.9.0.
	// It is still referenced in certain code paths for upgrade purporses
	HSMAutoDeprecated = "hsm-auto"
)

// Wrapper is the embedded implementation of autoSeal that contains logic
// specific to encrypting and decrypting data, or in this case keys.
type Wrapper interface {
	// Type is the type of Wrapper
	Type() string

	// KeyID is the ID of the key currently used for encryption
	KeyID() string
	// HMACKeyID is the ID of the key currently used for HMACing (if any)
	HMACKeyID() string

	// Init allows performing any necessary setup calls before using this Wrapper
	Init(context.Context) error
	// Finalize should be called when all usage of this Wrapper is done
	Finalize(context.Context) error

	// Encrypt encrypts the given byte slice and puts information about the final result in the returned value. The second byte slice is to pass any additional authenticated data; this may or may not be used depending on the particular implementation.
	Encrypt(context.Context, []byte, []byte) (*EncryptedBlobInfo, error)
	// Decrypt takes in the value and decrypts it into the byte slice.  The byte slice is to pass any additional authenticated data; this may or may not be used depending on the particular implementation.
	Decrypt(context.Context, *EncryptedBlobInfo, []byte) ([]byte, error)
}

// WrapperOptions contains options used when creating a Wrapper
type WrapperOptions struct {
	Logger hclog.Logger
}
