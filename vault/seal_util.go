// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Seal Wrapping

// SealWrapValue creates a SealWrappedValue wrapper with the entryValue being optionally encrypted with the give seal Access.
func SealWrapValue(ctx context.Context, access seal.Access, encrypt bool, entryValue []byte) (*SealWrappedValue, error) {
	blobInfo := &wrapping.BlobInfo{
		Wrapped:    false,
		Ciphertext: entryValue,
	}
	if access != nil && encrypt {
		swi, err := access.Encrypt(ctx, entryValue, nil)
		if err != nil {
			return nil, err
		}
		blobInfo = &wrapping.BlobInfo{
			Wrapped:    true,
			Ciphertext: swi.Ciphertext,
			Iv:         swi.Iv,
			Hmac:       swi.Hmac,
			KeyInfo:    swi.KeyInfo,
		}
	}

	return NewSealWrappedValue(blobInfo), nil
}

// MarshalSealWrappedValue marshals a SealWrappedValue into a byte slice. If the seal wrapped value contains
// a single wrapping.BlobInfo, the BlobInfo will be marshalled directly; otherwise the SealWrappedValue
// will be.
func MarshalSealWrappedValue(wrappedEntryValue *SealWrappedValue) ([]byte, error) {
	if len(wrappedEntryValue.value.Slots) > 1 {
		return wrappedEntryValue.marshal()
	}

	return proto.Marshal(wrappedEntryValue.GetUniqueBlobInfo())
}

// UnmarshalSealWrappedValue attempts to unmarshal a SealWrappedValue. This method can unmarshal marshalled
// SealWrappedValues as well as wrapping.BlobInfos. When a BlobInfo is encountered, a "transitory"
// SealWrappedValue will be returned.
func UnmarshalSealWrappedValue(value []byte) (*SealWrappedValue, error) {
	swv := &SealWrappedValue{}
	swvErr := swv.unmarshal(value)
	if swvErr == nil {
		return swv, nil
	}

	blobInfo := &wrapping.BlobInfo{}
	blobInfoErr := proto.Unmarshal(value, blobInfo)
	if blobInfoErr == nil {
		return newTransitorySealWrappedValue(blobInfo), nil
	}

	return nil, fmt.Errorf("error unmarshalling seal wrapped value: %w, %w", swvErr, blobInfoErr)
}

// UnmarshalSealWrappedValueWithCanary unmarshalls a byte array into a SealWrappedValue, taking care of
// removing the 's' canary value.
// This method returns true if a SealWrappedValue was successfully unmarshaled.
func UnmarshalSealWrappedValueWithCanary(value []byte) (*SealWrappedValue, bool) {
	eLen := len(value)
	if eLen > 0 && value[eLen-1] == 's' {
		if wrappedEntryValue, err := UnmarshalSealWrappedValue(value[:eLen-1]); err == nil {
			return wrappedEntryValue, true
		}
		// Else, note that having the canary value present is not a guarantee that
		// the value is wrapped, so if there is an error we will simply return a nil BlobInfo.
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Stored Barrier Keys (a.k.a. Root Key)

// SealWrapStoredBarrierKeys takes the barrier (root) keys, encrypts them using the seal access,
// and returns a physical.Entry for storage.
func SealWrapStoredBarrierKeys(ctx context.Context, access seal.Access, keys [][]byte) (*physical.Entry, error) {
	// Note that even though keys is a slice, it seems to always contain a single key.
	buf, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to encode keys for storage: %w", err)
	}

	wrappedEntryValue, err := SealWrapValue(ctx, access, true, buf)
	if err != nil {
		return nil, &ErrEncrypt{Err: fmt.Errorf("failed to encrypt keys for storage: %w", err)}
	}

	// Watch out, Wrapped has to be false for StoredBarrierKeysPath, since it used to be that the BlobInfo
	// returned by access.Encrypt() was marshalled directly. It probably would not matter if the value
	// was true, but setting if to false here makes TestSealWrapBackend_StorageBarrierKeyUpgrade_FromIVEntry
	// pass (maybe other tests as well?).
	wrappedEntryValue.GetUniqueBlobInfo().Wrapped = false

	wrappedValue, err := MarshalSealWrappedValue(wrappedEntryValue)
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
	wrappedEntryValue, err := UnmarshalSealWrappedValue(pe.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to proto decode stored keys: %w", err)
	}

	return decodeBarrierKeys(ctx, access, wrappedEntryValue.GetUniqueBlobInfo())
}

func decodeBarrierKeys(ctx context.Context, access seal.Access, blobInfo *wrapping.BlobInfo) ([][]byte, error) {
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
	wrappedEntryValue, err := SealWrapValue(ctx, access, true, key)
	if err != nil {
		return nil, &ErrEncrypt{Err: fmt.Errorf("failed to encrypt recovery key for storage: %w", err)}
	}

	// Not that we set Wrapped to false since it used to be that the BlobInfo returned by access.Encrypt()
	// was marshalled directly. It probably would not matter if the value was true, it doesn't seem to
	// break any tests.
	wrappedEntryValue.GetUniqueBlobInfo().Wrapped = false

	wrappedValue, err := MarshalSealWrappedValue(wrappedEntryValue)
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
	wrappedEntryValue, err := UnmarshalSealWrappedValue(pe.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to proto decode recevory key: %w", err)
	}

	blobInfo := wrappedEntryValue.GetUniqueBlobInfo()
	pt, err := access.Decrypt(ctx, blobInfo, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt recovery key from storage: %w", err)
	}
	return pt, blobInfo, nil
}
