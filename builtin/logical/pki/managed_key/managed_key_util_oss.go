// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package managed_key

import (
	"context"
	"crypto"
	"errors"
	"io"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func GetPublicKeyFromKeyBytes(ctx context.Context, mkv PkiManagedKeyView, keyBytes []byte) (crypto.PublicKey, error) {
	return nil, errEntOnly
}

func GenerateManagedKeyCABundle(ctx context.Context, b PkiManagedKeyView, keyId managedKeyId, data *certutil.CreationBundle, randomSource io.Reader) (bundle *certutil.ParsedCertBundle, err error) {
	return nil, errEntOnly
}

func GenerateManagedKeyCSRBundle(ctx context.Context, b PkiManagedKeyView, keyId managedKeyId, data *certutil.CreationBundle, addBasicConstraints bool, randomSource io.Reader) (bundle *certutil.ParsedCSRBundle, err error) {
	return nil, errEntOnly
}

func GetManagedKeyPublicKey(ctx context.Context, b PkiManagedKeyView, keyId managedKeyId) (crypto.PublicKey, error) {
	return nil, errEntOnly
}

func ParseManagedKeyCABundle(ctx context.Context, mkv PkiManagedKeyView, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	return nil, errEntOnly
}

func ExtractManagedKeyId(privateKeyBytes []byte) (UUIDKey, error) {
	return "", errEntOnly
}

func CreateKmsKeyBundle(ctx context.Context, mkv PkiManagedKeyView, keyId managedKeyId) (certutil.KeyBundle, certutil.PrivateKeyType, error) {
	return certutil.KeyBundle{}, certutil.UnknownPrivateKey, errEntOnly
}

func GetManagedKeyInfo(ctx context.Context, mkv PkiManagedKeyView, keyId managedKeyId) (*ManagedKeyInfo, error) {
	return nil, errEntOnly
}
