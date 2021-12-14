package vault

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type ManagedKeyType string

const (
	ManagedKeyTypePkcs11 ManagedKeyType = "pkcs11"
)

// ManagedKeyConfiguration holds all the data needed for keys that are managed through an
// external KMS library. The ManagedKeyRegistry keeps track of them.
type ManagedKeyConfiguration struct {
	Type ManagedKeyType

	// Library is the name of the KMS library to use to access the key. It
	// should match the 'name' field of the 'kms_library' configuration stanza.
	Library string

	// Name to identify the key. Must be unique within the namespace in which it is created.
	Name string

	// AnyMount if true, indicates that any mount in the namespace may use the key,
	// without needing to tune the mount to specify the key name in the 'allowed_keys' field.
	AnyMount bool

	// AllowGenerateKey if true, allows users of the key to trigger the key generation.
	AllowGenerateKey bool

	// AllowStoreKey if true, allows users of the key to add key material when none was
	// previously present
	AllowStoreKey bool

	// AllowReplaceKey if true, allows users of the key to add key material even when
	// a key was previously present
	AllowReplaceKey bool

	// RawParameters are the parameters required by the KMS library to manage the key
	RawParameters map[string]interface{}

	// Parameters required by the KMS library in a key type specific struct, set by method initParameters()
	Parameters interface{}
}

type Pkcs11ManagedKeyParameters struct {
	Slot       uint
	TokenLabel string `mapstructure:"token_label"`
	Pin        string
	KeyId      string `mapstructure:"key_id"`
	KeyLabel   string `mapstructure:"key_label"`
	Mechanism  uint
}

// initParameters converts the RawParameters field into the key type-specific Parameter
func (keyConfig *ManagedKeyConfiguration) initParameters() error {
	switch keyConfig.Type {
	case ManagedKeyTypePkcs11:
		var parameters Pkcs11ManagedKeyParameters
		err := mapstructure.WeakDecode(keyConfig.RawParameters, &parameters)
		if err != nil {
			return err
		}
		keyConfig.Parameters = parameters

	default:
		return fmt.Errorf("cannot get PKCS#11 parameters for a managed key of type %s", keyConfig.Type)
	}

	return nil
}
