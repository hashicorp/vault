// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package keysutil

import (
	"crypto/x509"
	"fmt"
	"io"

	"github.com/hashicorp/vault/sdk/helper/errutil"
)

type entKeyEntry struct{}

func (e entKeyEntry) IsEntPrivateKeyMissing() bool {
	return true
}

func entSignWithOptions(p *Policy, input, context []byte, ver int, hashAlgorithm HashType, options *SigningOptions) ([]byte, error) {
	return nil, fmt.Errorf("unsupported key type %v", p.Type)
}

func entVerifySignatureWithOptions(p *Policy, input, context []byte, sigBytes []byte, ver int, options *SigningOptions) (bool, error) {
	return false, errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
}

func entRotateInMemory(p *Policy, keyType KeyType, entry *KeyEntry, rand io.Reader, config *KeyConfig) error {
	return fmt.Errorf("unsupported key type %v", keyType)
}

func entEncryptWithOptions(p *Policy, opts EncryptionOptions, value []byte) ([]byte, error) {
	return nil, fmt.Errorf("unsupported key type %v", p.Type)
}

func entDecryptWithOptions(p *Policy, opts EncryptionOptions, value []byte) ([]byte, error) {
	return nil, errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
}

func (p *Policy) CreateCsrWithManagedKeyVersion(params ManagedKeyParameters) CsrCreator {
	return func(_ int, _ *x509.CertificateRequest) ([]byte, error) {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
	}
}

func (p *Policy) ValidateLeafCertKeyMatchWithManagedKeyVersion(keyVersion int, certPublicKeyAlgorithm x509.PublicKeyAlgorithm, certPublicKey any, params ManagedKeyParameters) (bool, error) {
	return false, errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
}
