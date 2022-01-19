//go:build hsm

package managed_key

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestManagedKeyConfiguration_initParameters(t *testing.T) {
	key := &Configuration{
		Type: Pkcs11Type,
		Parameters: map[string]interface{}{
			"slot":        "1",
			"token_label": "the token label",
			"pin":         "the pin",
			"key_id":      "the key ID",
			"key_label":   "the key label",
			"mechanism":   "0xa",
		},
	}

	err := key.initParameters()
	if err != nil {
		t.Fatal(err)
	}

	expected := &pkcs11ManagedKeyParameters{
		Slot:       1,
		TokenLabel: "the token label",
		Pin:        "the pin",
		KeyId:      "the key ID",
		KeyLabel:   "the key label",
		Mechanism:  0xa,
	}

	require.Equal(t, expected, key.TypedParams)
}

func TestManagedKeyConfiguration_initParameters_unsupported_type(t *testing.T) {
	key := &Configuration{
		Type: "awskms",
	}

	err := key.initParameters()
	require.Error(t, err, "awskms is not a supported key type")
}
