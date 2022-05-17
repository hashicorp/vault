package pki

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestPKI_PathManageKeys_GenerateKeys(t *testing.T) {
	b, s := createBackendWithStorage(t)

	tests := []struct {
		name           string
		keyType        string
		keyBits        []int
		wantLogicalErr bool
	}{
		{"all-defaults", "", []int{0}, false},
		{"rsa", "rsa", []int{0, 2048, 3072, 4096}, false},
		{"ec", "ec", []int{0, 224, 256, 384, 521}, false},
		{"ed25519", "ed25519", []int{0}, false},
		{"error-rsa", "rsa", []int{-1, 343444}, true},
		{"error-ec", "ec", []int{-1, 3434324}, true},
		{"error-bad-type", "dskjfkdsfjdkf", []int{0}, true},
	}
	for _, tt := range tests {
		for _, keyBitParam := range tt.keyBits {
			keyName := fmt.Sprintf("%s-%d", tt.name, keyBitParam)
			t.Run(keyName, func(t *testing.T) {
				data := make(map[string]interface{})
				if tt.keyType != "" {
					data["key_type"] = tt.keyType
				}
				if keyBitParam != 0 {
					data["key_bits"] = keyBitParam
				}
				resp, err := b.HandleRequest(context.Background(), &logical.Request{
					Operation:  logical.UpdateOperation,
					Path:       "keys/generate/internal",
					Storage:    s,
					Data:       data,
					MountPoint: "pki/",
				})
				require.NoError(t, err,
					"Failed generating key with values key_type:%s key_bits:%d key_name:%s", tt.keyType, keyBitParam, keyName)
				require.NotNil(t, resp,
					"Got nil response generating key with values key_type:%s key_bits:%d key_name:%s", tt.keyType, keyBitParam, keyName)
				if tt.wantLogicalErr {
					require.True(t, resp.IsError(), "expected logical error but the request passed:\n%#v", resp)
				} else {
					require.False(t, resp.IsError(),
						"Got logical error response when not expecting one, "+
							"generating key with values key_type:%s key_bits:%d key_name:%s\n%s", tt.keyType, keyBitParam, keyName, resp.Error())
				}
			})
		}
	}
}
