package vault

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestManagedKeyConfiguration_initParameters(t *testing.T) {
	key := &ManagedKeyConfiguration{
		Type: ManagedKeyTypePkcs11,
		RawParameters: map[string]interface{}{
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

	expected := Pkcs11ManagedKeyParameters{
		Slot:       1,
		TokenLabel: "the token label",
		Pin:        "the pin",
		KeyId:      "the key ID",
		KeyLabel:   "the key label",
		Mechanism:  0xa,
	}

	require.Equal(t, expected, key.Parameters)
}

func TestManagedKeyConfiguration_initParameters_unsupported_type(t *testing.T) {
	key := &ManagedKeyConfiguration{
		Type: "awskms",
	}

	err := key.initParameters()
	require.Error(t, err, "awskms is not a supported key type")
}
