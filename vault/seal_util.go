// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
	"google.golang.org/protobuf/proto"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Seal Wrapping

// SealWrapValue creates a BlobInfo wrapper with the entryValue being optionally encrypted with the give seal Access.
func SealWrapValue(ctx context.Context, access seal.Access, encrypt bool, entryValue []byte) (*wrapping.BlobInfo, error) {
	wrappedEntryValue := &wrapping.BlobInfo{
		Wrapped:    false,
		Ciphertext: entryValue,
	}

	if access != nil && encrypt {
		swi, err := access.Encrypt(ctx, entryValue, nil)
		if err != nil {
			return nil, err
		}
		wrappedEntryValue.Wrapped = true
		wrappedEntryValue.Ciphertext = swi.Ciphertext
		wrappedEntryValue.Iv = swi.Iv
		wrappedEntryValue.Hmac = swi.Hmac
		wrappedEntryValue.KeyInfo = swi.KeyInfo
	}

	return wrappedEntryValue, nil
}

// MarshalSealWrappedValue marshals a BlobInfo into a byte slice.
func MarshalSealWrappedValue(wrappedEntryValue *wrapping.BlobInfo) ([]byte, error) {
	return proto.Marshal(wrappedEntryValue)
}

// UnmarshalSealWrappedValue attempts to unmarshal a BlobInfo.
func UnmarshalSealWrappedValue(value []byte) (*wrapping.BlobInfo, error) {
	wrappedEntryValue := &wrapping.BlobInfo{}
	err := proto.Unmarshal(value, wrappedEntryValue)
	if err != nil {
		return nil, err
	}
	return wrappedEntryValue, nil
}

// UnmarshalSealWrappedValueWithCanary unmarshalls a byte array into a BlobInfo, taking care of
// removing the 's' canary value. Note that if the value does not end with the canary value,
// or a BlobInfo cannot be unmarshalled, nil is returned.
func UnmarshalSealWrappedValueWithCanary(value []byte) *wrapping.BlobInfo {
	eLen := len(value)
	if eLen > 0 && value[eLen-1] == 's' {
		if wrappedEntryValue, err := UnmarshalSealWrappedValue(value[:eLen-1]); err == nil {
			return wrappedEntryValue
		}
		// Else, note that having the canary value present is not a guarantee that
		// the value is wrapped, so if there is an error we will simply return a nil BlobInfo.
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Stored Barrier Keys (a.k.a. Root Key)

// SealWrapStoredBarrierKeys takes the json-marshalled barriers (root) keys, encrypts them using the seal access,
// and returns a physical.Entry for storage.
func SealWrapStoredBarrierKeys(ctx context.Context, access seal.Access, keys [][]byte) (*physical.Entry, error) {
	buf, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to encode keys for storage: %w", err)
	}

	blobInfo, err := SealWrapValue(ctx, access, true, buf)
	if err != nil {
		return nil, &ErrEncrypt{Err: fmt.Errorf("failed to encrypt keys for storage: %w", err)}
	}

	// Watch out, Wrapped has to be false for StoredBarrierKeysPath, since it used to be that the BlobInfo
	// returned by access.Encrypt() was marshalled directly. It probably would not matter if the value
	// was true, but setting if to false here makes TestSealWrapBackend_StorageBarrierKeyUpgrade_FromIVEntry
	// pass (maybe other tests as well?).
	blobInfo.Wrapped = false

	wrappedValue, err := MarshalSealWrappedValue(blobInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value for storage: %w", err)
	}
	return &physical.Entry{
		Key:   StoredBarrierKeysPath,
		Value: wrappedValue,
	}, nil
}

// UnsealWrapStoredBarrierKeys is the counterpart to SealWrapStoredBarrierKeys.
func UnsealWrapStoredBarrierKeys(ctx context.Context, access seal.Access, pe *physical.Entry) ([][]byte, error) {
	blobInfo, err := UnmarshalSealWrappedValue(pe.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to proto decode stored keys: %w", err)
	}

	pt, err := access.Decrypt(ctx, blobInfo, nil)
	if err != nil {
		if strings.Contains(err.Error(), "message authentication failed") {
			return nil, &ErrInvalidKey{Reason: fmt.Sprintf("failed to decrypt keys from storage: %v", err)}
		}
		return nil, &ErrDecrypt{Err: fmt.Errorf("failed to decrypt keys from storage: %w", err)}
	}

	// Decode the barrier entry
	var keys [][]byte
	if err := json.Unmarshal(pt, &keys); err != nil {
		return nil, fmt.Errorf("failed to decode stored keys: %v", err)
	}
	return keys, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Recovery Key

// SealWrapRecoveryKey encrypts the recovery key using the given seal access and returns a physical.Entry for storage.
func SealWrapRecoveryKey(ctx context.Context, access seal.Access, key []byte) (*physical.Entry, error) {
	blobInfo, err := SealWrapValue(ctx, access, true, key)
	if err != nil {
		return nil, &ErrEncrypt{Err: fmt.Errorf("failed to encrypt keys for storage: %w", err)}
	}

	// Not that we set Wrapped to false since it used to be that the BlobInfo returned by access.Encrypt()
	// was marshalled directly. It probably would not matter if the value was true, it doesn't seem to
	// break any tests.
	blobInfo.Wrapped = false

	wrappedValue, err := MarshalSealWrappedValue(blobInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value for storage: %w", err)
	}
	return &physical.Entry{
		Key:   recoveryKeyPath,
		Value: wrappedValue,
	}, nil
}

// UnsealWrapRecoveryKey is the counterpart to SealWrapRecoveryKey.
func UnsealWrapRecoveryKey(ctx context.Context, access seal.Access, pe *physical.Entry) ([]byte, *wrapping.BlobInfo, error) {
	blobInfo, err := UnmarshalSealWrappedValue(pe.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to proto decode recevory key: %w", err)
	}

	pt, err := access.Decrypt(ctx, blobInfo, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt recovery key from storage: %w", err)
	}
	return pt, blobInfo, nil
}
