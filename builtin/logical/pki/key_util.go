package pki

import (
	"context"
	"crypto"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type managedKeyContext struct {
	ctx        context.Context
	b          *backend
	mountPoint string
}

func newManagedKeyContext(ctx context.Context, b *backend, mountPoint string) managedKeyContext {
	return managedKeyContext{
		ctx:        ctx,
		b:          b,
		mountPoint: mountPoint,
	}
}

func comparePublicKey(ctx managedKeyContext, key *keyEntry, publicKey crypto.PublicKey) (bool, error) {
	publicKeyForKeyEntry, err := getPublicKey(ctx, key)
	if err != nil {
		return false, err
	}

	return certutil.ComparePublicKeysAndType(publicKeyForKeyEntry, publicKey)
}

func getPublicKey(mkc managedKeyContext, key *keyEntry) (crypto.PublicKey, error) {
	if key.PrivateKeyType == certutil.ManagedPrivateKey {
		keyId, err := extractManagedKeyId([]byte(key.PrivateKey))
		if err != nil {
			return nil, err
		}
		return getManagedKeyPublicKey(mkc, keyId)
	}

	signer, _, _, err := getSignerFromKeyEntryBytes(key)
	if err != nil {
		return nil, err
	}
	return signer.Public(), nil
}

func getSignerFromKeyEntryBytes(key *keyEntry) (crypto.Signer, certutil.BlockType, *pem.Block, error) {
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

func getManagedKeyPublicKey(mkc managedKeyContext, keyId managedKeyId) (crypto.PublicKey, error) {
	// Determine key type and key bits from the managed public key
	var pubKey crypto.PublicKey
	err := withManagedPKIKey(mkc.ctx, mkc.b, keyId, mkc.mountPoint, func(ctx context.Context, key logical.ManagedSigningKey) error {
		var myErr error
		pubKey, myErr = key.GetPublicKey(ctx)
		if myErr != nil {
			return myErr
		}

		return nil
	})
	if err != nil {
		return nil, errors.New("failed to lookup public key from managed key: " + err.Error())
	}
	return pubKey, nil
}

func importKeyFromBytes(mkc managedKeyContext, s logical.Storage, keyValue string, keyName string) (*keyEntry, bool, error) {
	signer, _, _, err := getSignerFromBytes([]byte(keyValue))
	if err != nil {
		return nil, false, err
	}
	privateKeyType := certutil.GetPrivateKeyTypeFromSigner(signer)
	if privateKeyType == certutil.UnknownPrivateKey {
		return nil, false, errors.New("unsupported private key type within pem bundle")
	}

	key, existed, err := importKey(mkc, s, keyValue, keyName, privateKeyType)
	if err != nil {
		return nil, false, err
	}
	return key, existed, nil
}
