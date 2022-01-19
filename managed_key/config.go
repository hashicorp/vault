package managed_key

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type typedParams interface{}

// Configuration holds all the data needed for keys that are managed through an
// external KMS library. The ManagedKeyRegistry keeps track of them.
type Configuration struct {
	Type Type `json:"type" mapstructure:"type"`

	// Library is the name of the KMS library to use to access the key. It
	// should match the 'name' field of the 'kms_library' configuration stanza.
	Library string `json:"library" mapstructure:"library"`

	// Name to identify the key. Must be unique within the namespace in which it is created.
	Name string `json:"name" mapstructure:"name"`

	// AnyMount if true, indicates that any mount in the namespace may use the key,
	// without needing to tune the mount to specify the key name in the 'allowed_keys' field.
	AnyMount bool `json:"any_mount" mapstructure:"any_mount"`

	// AllowGenerateKey if true, allows users of the key to trigger the key generation.
	AllowGenerateKey bool `json:"allow_generate_key" mapstructure:"allow_generate_key"`

	// AllowStoreKey if true, allows users of the key to add key material when none was
	// previously present
	AllowStoreKey bool `json:"allow_store_key" mapstructure:"allow_store_key"`

	// AllowReplaceKey if true, allows users of the key to add key material even when
	// a key was previously present
	AllowReplaceKey bool `json:"allow_replace_key" mapstructure:"allow_replace_key"`

	// Parameters are the parameters required by the KMS library to manage the key
	Parameters map[string]interface{} `json:"parameters" mapstructure:"-"`

	// Parameters required by the KMS library in a key type specific struct, set by method initParameters()
	TypedParams typedParams `json:"-" mapstructure:"-"`
}

// initParameters converts the RawParameters field into the key type-specific Parameter
func (keyConfig *Configuration) initParameters() error {
	var err error
	factory, err := GetFactory(keyConfig.Type)
	if err != nil {
		return err
	}
	keyConfig.TypedParams, err = factory.parseParameters(keyConfig.Parameters)
	if err != nil {
		return err
	}

	return nil
}

func (keyConfig *Configuration) MarshalJSON() ([]byte, error) {
	type config2 Configuration
	// Preprocessing: Serialize the TypedParams
	if keyConfig.Parameters == nil {
		if err := mapstructure.Decode(keyConfig.TypedParams, &keyConfig.Parameters); err != nil {
			return nil, err
		}
	}
	return json.Marshal((*config2)(keyConfig))
}

func (keyConfig *Configuration) UnmarshalJSON(data []byte) error {
	type config2 Configuration
	if err := json.Unmarshal(data, (*config2)(keyConfig)); err != nil {
		return err
	}

	// Post processing
	return keyConfig.initParameters()
}
