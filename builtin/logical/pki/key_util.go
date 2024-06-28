// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/managed_key"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
)

func comparePublicKey(sc *storageContext, key *issuing.KeyEntry, publicKey crypto.PublicKey) (bool, error) {
	publicKeyForKeyEntry, err := getPublicKey(sc.Context, sc.GetPkiManagedView(), key)
	if err != nil {
		return false, err
	}

	return certutil.ComparePublicKeysAndType(publicKeyForKeyEntry, publicKey)
}

func getPublicKey(ctx context.Context, mkv managed_key.PkiManagedKeyView, key *issuing.KeyEntry) (crypto.PublicKey, error) {
	if key.PrivateKeyType == certutil.ManagedPrivateKey {
		return managed_key.GetPublicKeyFromKeyBytes(ctx, mkv, []byte(key.PrivateKey))
	}

	signer, _, _, err := getSignerFromKeyEntryBytes(key)
	if err != nil {
		return nil, err
	}
	return signer.Public(), nil
}

func getSignerFromKeyEntryBytes(key *issuing.KeyEntry) (crypto.Signer, certutil.BlockType, *pem.Block, error) {
	if key.PrivateKeyType == certutil.UnknownPrivateKey {
		return nil, certutil.UnknownBlock, nil, errutil.InternalError{Err: fmt.Sprintf("unsupported unknown private key type for key: %s (%s)", key.ID, key.Name)}
	}

	if key.PrivateKeyType == certutil.ManagedPrivateKey {
		return nil, certutil.UnknownBlock, nil, errutil.InternalError{Err: fmt.Sprintf("can not get a signer from a managed key: %s (%s)", key.ID, key.Name)}
	}

	bytes, blockType, blk, err := getSignerFromBytes([]byte(key.PrivateKey))
	if err != nil {
		return nil, certutil.UnknownBlock, nil, errutil.InternalError{Err: fmt.Sprintf("failed parsing key entry bytes for key id: %s (%s): %s", key.ID, key.Name, err.Error())}
	}

	return bytes, blockType, blk, nil
}

func getSignerFromBytes(keyBytes []byte) (crypto.Signer, certutil.BlockType, *pem.Block, error) {
	pemBlock, _ := pem.Decode(keyBytes)
	if pemBlock == nil {
		return nil, certutil.UnknownBlock, pemBlock, errutil.InternalError{Err: "no data found in PEM block"}
	}

	signer, blk, err := certutil.ParseDERKey(pemBlock.Bytes)
	if err != nil {
		return nil, certutil.UnknownBlock, pemBlock, errutil.InternalError{Err: fmt.Sprintf("failed to parse PEM block: %s", err.Error())}
	}
	return signer, blk, pemBlock, nil
}

func getPublicKeyFromBytes(keyBytes []byte) (crypto.PublicKey, error) {
	signer, _, _, err := getSignerFromBytes(keyBytes)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("failed parsing key bytes: %s", err.Error())}
	}

	return signer.Public(), nil
}

func importKeyFromBytes(sc *storageContext, keyValue string, keyName string) (*issuing.KeyEntry, bool, error) {
	signer, _, _, err := getSignerFromBytes([]byte(keyValue))
	if err != nil {
		return nil, false, err
	}
	privateKeyType := certutil.GetPrivateKeyTypeFromSigner(signer)
	if privateKeyType == certutil.UnknownPrivateKey {
		return nil, false, errors.New("unsupported private key type within pem bundle")
	}

	key, existed, err := sc.importKey(keyValue, keyName, privateKeyType)
	if err != nil {
		return nil, false, err
	}
	return key, existed, nil
}
