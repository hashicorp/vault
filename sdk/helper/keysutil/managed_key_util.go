//go:build !enterprise

package keysutil

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/logical"
)

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func (p *Policy) decryptWithManagedKey(params *ManagedKeyParameters, keyEntry KeyEntry, ciphertext []byte, nonce []byte) (plaintext []byte, err error) {
	return nil, errEntOnly
}

func (p *Policy) encryptWithManagedKey(params *ManagedKeyParameters, keyEntry KeyEntry, plaintext []byte, nonce []byte, ver int) (ciphertext []byte, err error) {
	return nil, errEntOnly
}

func (p *Policy) signWithManagedKey(params *ManagedKeyParameters, options *SigningOptions, keyEntry KeyEntry, input []byte) (sig []byte, err error) {
	return nil, errEntOnly
}

func (p *Policy) verifyWithManagedKey(params *ManagedKeyParameters, options *SigningOptions, keyEntry KeyEntry, input, sig []byte) (verified bool, err error) {
	return false, errEntOnly
}

func (p *Policy) HMACWithManagedKey(ctx context.Context, ver int, managedKeySystemView logical.ManagedKeySystemView, backendUUID string, algorithm string, data []byte) (hmacBytes []byte, err error) {
	return nil, errEntOnly
}

func GetManagedKeyUUID(params *ManagedKeyParameters, keyName string, keyId string) (uuid string, err error) {
	return "", errEntOnly
}
