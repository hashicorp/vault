// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-kms-wrapping/v2/aead"
)

type StoredKeysSupport int

const (
	// The 0 value of StoredKeysSupport is an invalid option
	StoredKeysInvalid StoredKeysSupport = iota
	StoredKeysNotSupported
	StoredKeysSupportedGeneric
	StoredKeysSupportedShamirRoot
)

func (s StoredKeysSupport) String() string {
	switch s {
	case StoredKeysNotSupported:
		return "Old-style Shamir"
	case StoredKeysSupportedGeneric:
		return "AutoUnseal"
	case StoredKeysSupportedShamirRoot:
		return "New-style Shamir"
	default:
		return "Invalid StoredKeys type"
	}
}

type SealInfo struct {
	wrapping.Wrapper
	Priority int
	Name     string
	Healthy  bool
}

func (si *SealInfo) keyId(ctx context.Context) string {
	if id, err := si.Wrapper.KeyId(ctx); err == nil {
		return id
	}
	return ""
}

// Access is the embedded implementation of autoSeal that contains logic
// specific to encrypting and decrypting data, or in this case keys.
type Access interface {
	wrapping.InitFinalizer

	Generation() uint64

	// Encrypt encrypts the given byte slice and stores the resulting
	// information in the returned blob info. Which options are used depends on
	// the underlying wrapper. Supported options: WithAad.
	// Returns a MultiWrapValue as long as at least one seal Access wrapper encrypted the data successfully, and
	// if this is the case errors may still be returned if any wrapper failed. The error map is keyd by seal name.
	Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) (*wrapping.MultiWrapValue, map[string]error)

	// Decrypt decrypts the given byte slice and stores the resulting information in the
	// returned byte slice. Which options are used depends on the underlying wrapper.
	// Supported options: WithAad.
	// Returns the plaintext, a flag indicating whether the ciphertext is up-to-date
	// (according to IsUpToDate), and an error.
	Decrypt(ctx context.Context, ciphertext *wrapping.MultiWrapValue, options ...wrapping.Option) ([]byte, bool, error)

	// IsUpToDate returns true if a MultiWrapValue is up-to-date. An MultiWrapValue is
	// considered to be up-to-date if its generation matches the Access generation, and if
	// it has a slot with a key ID that match the current key ID of each of the Access
	// wrappers.
	IsUpToDate(ctx context.Context, value *wrapping.MultiWrapValue, forceKeyIdRefresh bool) (bool, error)

	GetWrapper() wrapping.Wrapper
	SetShamirSealKey([]byte) error
	GetShamirKeyBytes(ctx context.Context) ([]byte, error)
	SealType(ctx context.Context) (SealType, error)
	// GetSealInfoByPriority the returned slice should be sorted in priority.
	GetSealInfoByPriority() []SealInfo
}

type access struct {
	generation         uint64
	wrappersByPriority []SealInfo
	keyIdSet           keyIdSet
}

var _ Access = (*access)(nil)

func NewAccess(sealInfos []SealInfo) Access {
	if len(sealInfos) == 0 {
		panic("cannot create a seal.Access without any seal info")
	}
	a := &access{
		generation:         1, // FIXME: Introduce Generation argument
		wrappersByPriority: sealInfos,
	}

	sort.Slice(a.wrappersByPriority, func(i int, j int) bool { return a.wrappersByPriority[i].Priority < a.wrappersByPriority[j].Priority })

	return a
}

func (a *access) GetSealInfoByPriority() []SealInfo {
	return a.wrappersByPriority
}

func (a *access) Generation() uint64 {
	return a.generation
}

func (a *access) KeyId(ctx context.Context) (string, error) {
	w := a.getDefaultWrapper()
	if w != nil {
		return w.KeyId(ctx)
	}
	return "", errors.New("no wrapper configured")
}

func (a *access) SetConfig(ctx context.Context, options ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	w := a.getDefaultWrapper()
	if w != nil {
		return w.SetConfig(ctx, options...)
	}

	return nil, errors.New("no wrapper configured")
}

func (a *access) GetWrapper() wrapping.Wrapper {
	return a.getDefaultWrapper()
}

func (a *access) Init(ctx context.Context, options ...wrapping.Option) error {
	var keyIds []string
	for _, sealInfo := range a.wrappersByPriority {
		if initWrapper, ok := sealInfo.Wrapper.(wrapping.InitFinalizer); ok {
			if err := initWrapper.Init(ctx, options...); err != nil {
				return err
			}
			keyId, err := sealInfo.Wrapper.KeyId(ctx)
			if err != nil {
				return fmt.Errorf("cannod determine key ID for seal %s: %w", sealInfo.Name, err)
			}
			keyIds = append(keyIds, keyId)
		}
	}
	a.keyIdSet.setIds(keyIds)
	return nil
}

func (a *access) getDefaultWrapper() wrapping.Wrapper {
	if len(a.wrappersByPriority) > 0 {
		return a.wrappersByPriority[0].Wrapper
	}
	return nil
}

func (a *access) Type(ctx context.Context) (wrapping.WrapperType, error) {
	return a.getDefaultWrapper().Type(ctx)
}

func (a *access) SealType(ctx context.Context) (SealType, error) {
	if len(a.wrappersByPriority) > 1 {
		return SealTypeMultiSeal, nil
	}

	wrapperType, err := a.getDefaultWrapper().Type(ctx)
	if err != nil {
		return "", err
	}

	return SealType(wrapperType), nil
}

func (a *access) IsUpToDate(ctx context.Context, value *wrapping.MultiWrapValue, forceKeyIdRefresh bool) (bool, error) {
	// TODO(SEALHA): Enable Generation checking
	//if a.Generation() != value.Generation {
	//	return false, nil
	//}
	if forceKeyIdRefresh {
		test, errs := a.Encrypt(ctx, []byte{0})
		if test == nil {
			return false, JoinSealWrapErrors("cannot determine key IDs of Access wrappers", errs)
		}
		// TODO(SEALHA): What to do if there are partial failures?
		a.keyIdSet.set(test)
	}

	return a.keyIdSet.equal(value), nil
}

// Encrypt uses the underlying seal to encrypt the plaintext and returns it.
func (a *access) Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) (*wrapping.MultiWrapValue, map[string]error) {
	var slots []*wrapping.BlobInfo
	errs := make(map[string]error)

	for _, sealInfo := range a.wrappersByPriority {
		var encryptErr error
		defer func(now time.Time) {
			metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
			metrics.MeasureSince([]string{"seal", sealInfo.Name, "encrypt", "time"}, now)

			if encryptErr != nil {
				metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
				metrics.IncrCounter([]string{"seal", sealInfo.Name, "encrypt", "error"}, 1)
			}
		}(time.Now())

		metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
		metrics.IncrCounter([]string{"seal", sealInfo.Name, "encrypt"}, 1)

		ciphertext, encryptErr := sealInfo.Wrapper.Encrypt(ctx, plaintext, options...)
		if encryptErr != nil {
			// TODO (multiseal): logic for failures and setting the sentinel for retrying
			errs[sealInfo.Name] = encryptErr
		} else {
			slots = append(slots, ciphertext)
		}
	}

	if len(slots) == 0 {
		return nil, errs
	}

	ret := &wrapping.MultiWrapValue{
		Generation: a.Generation(),
		Slots:      slots,
	}

	// cache key IDs
	a.keyIdSet.set(ret)

	return ret, errs
}

// Decrypt uses the underlying seal to decrypt the ciphertext and returns it.
// Note that it is possible depending on the wrapper used that both pt and err
// are populated.
// Returns the plaintext, a flag indicating whether the ciphertext is up-to-date
// (according to IsUpToDate), and an error.
func (a *access) Decrypt(ctx context.Context, ciphertext *wrapping.MultiWrapValue, options ...wrapping.Option) ([]byte, bool, error) {
	blobInfoMap := slotsByKeyId(ciphertext)

	isUpToDate, err := a.IsUpToDate(ctx, ciphertext, false)
	if err != nil {
		return nil, false, err
	}

	// First, lets try the wrappers in order of priority and look for an exact key ID match
	for _, sealInfo := range a.wrappersByPriority {
		if keyId, err := sealInfo.Wrapper.KeyId(ctx); err == nil {
			if blobInfo, ok := blobInfoMap[keyId]; ok {
				pt, oldKey, err := a.tryDecrypt(ctx, sealInfo, blobInfo, options)
				if oldKey {
					return pt, false, err
				}
				if err == nil {
					return pt, isUpToDate, nil
				}
				// If there is an error, keep trying with the other wrappers
			}
		}
	}

	// No key ID match, so try each wrapper with all slots
	errs := make(map[string]error)
	for _, sealInfo := range a.wrappersByPriority {
		for _, blobInfo := range ciphertext.Slots {
			pt, oldKey, err := a.tryDecrypt(ctx, sealInfo, blobInfo, options)
			if oldKey {
				return pt, false, err
			}
			if err == nil {
				return pt, isUpToDate, nil
			}
			errs[sealInfo.Name] = err
		}
	}

	return nil, false, JoinSealWrapErrors("error decrypting seal wrapped value", errs)
}

func (a *access) tryDecrypt(ctx context.Context, sealInfo SealInfo, ciphertext *wrapping.BlobInfo, options []wrapping.Option) ([]byte, bool, error) {
	var decryptErr error
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", sealInfo.Name, "decrypt", "time"}, now)

		if decryptErr != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", sealInfo.Name, "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", sealInfo.Name, "decrypt"}, 1)

	pt, err := sealInfo.Wrapper.Decrypt(ctx, ciphertext, options...)
	isOldKey := false
	if err != nil && err.Error() == "decrypted with old key" {
		// This is for compatibility with sealWrapMigration
		isOldKey = true
	}
	return pt, isOldKey, err
}

func JoinSealWrapErrors(msg string, errorMap map[string]error) error {
	errs := []error{errors.New(msg)}
	for name, err := range errorMap {
		errs = append(errs, fmt.Errorf("error decrypting using seal %s: %w", name, err))
	}
	return errors.Join(errs...)
}

func (a *access) Finalize(ctx context.Context, options ...wrapping.Option) error {
	var errs []error

	for _, w := range a.wrappersByPriority {
		if finalizeWrapper, ok := w.Wrapper.(wrapping.InitFinalizer); ok {
			if err := finalizeWrapper.Finalize(ctx, options...); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (a *access) SetShamirSealKey(key []byte) error {
	if len(a.wrappersByPriority) == 0 {
		return errors.New("no wrapper configured")
	}

	wrapper := a.wrappersByPriority[0].Wrapper

	shamirWrapper, ok := wrapper.(*aead.ShamirWrapper)
	if !ok {
		return errors.New("seal is not a Shamir seal")
	}

	return shamirWrapper.SetAesGcmKeyBytes(key)
}

func (a *access) GetShamirKeyBytes(ctx context.Context) ([]byte, error) {
	if len(a.wrappersByPriority) == 0 {
		return nil, errors.New("no wrapper configured")
	}

	wrapper := a.wrappersByPriority[0].Wrapper

	shamirWrapper, ok := wrapper.(*aead.ShamirWrapper)
	if !ok {
		return nil, errors.New("seal is not a shamir seal")
	}

	return shamirWrapper.KeyBytes(ctx)
}

func slotsByKeyId(value *wrapping.MultiWrapValue) map[string]*wrapping.BlobInfo {
	ret := make(map[string]*wrapping.BlobInfo)
	for _, blobInfo := range value.Slots {
		keyId := ""
		if blobInfo.KeyInfo != nil {
			keyId = blobInfo.KeyInfo.KeyId
		}
		ret[keyId] = blobInfo
	}
	return ret
}

type keyIdSet struct {
	keyIds atomic.Pointer[[]string]
}

func (s *keyIdSet) set(value *wrapping.MultiWrapValue) {
	keyIds := s.collect(value)
	s.setIds(keyIds)
}

func (s *keyIdSet) setIds(keyIds []string) {
	keyIds = s.deduplicate(keyIds)
	s.keyIds.Store(&keyIds)
}

func (s *keyIdSet) get() []string {
	pids := s.keyIds.Load()
	if pids == nil {
		return nil
	}
	return *pids
}

func (s *keyIdSet) equal(value *wrapping.MultiWrapValue) bool {
	keyIds := s.collect(value)
	expected := s.get()
	return reflect.DeepEqual(keyIds, expected)
}

func (s *keyIdSet) collect(value *wrapping.MultiWrapValue) []string {
	var keyIds []string
	for _, blobInfo := range value.Slots {
		if blobInfo.KeyInfo != nil {
			// Ideally we should always have a KeyInfo.KeyId, but:
			// 1) plaintext entries are stored on a blob info with Wrapped == false
			// 2) some unit test wrappers do not return a blob info
			keyIds = append(keyIds, blobInfo.KeyInfo.KeyId)
		}
	}
	return s.deduplicate(keyIds)
}

func (s *keyIdSet) deduplicate(ids []string) []string {
	m := make(map[string]struct{})
	for _, id := range ids {
		m[id] = struct{}{}
	}
	deduplicated := make([]string, 0, len(m))
	for id := range m {
		deduplicated = append(deduplicated, id)
	}
	sort.Strings(deduplicated)
	return deduplicated
}

type SealType string

const (
	SealTypeMultiSeal         SealType = "multiseal"
	SealTypeAliCloudKms                = SealType(wrapping.WrapperTypeAliCloudKms)
	SealTypeAwsKms                     = SealType(wrapping.WrapperTypeAwsKms)
	SealTypeAzureKeyVault              = SealType(wrapping.WrapperTypeAzureKeyVault)
	SealTypeGcpCkms                    = SealType(wrapping.WrapperTypeGcpCkms)
	SealTypePkcs11                     = SealType(wrapping.WrapperTypePkcs11)
	SealTypeOciKms                     = SealType(wrapping.WrapperTypeOciKms)
	SealTypeShamir                     = SealType(wrapping.WrapperTypeShamir)
	SealTypeTransit                    = SealType(wrapping.WrapperTypeTransit)
	SealTypeHsmAutoDeprecated          = SealType(wrapping.WrapperTypeHsmAuto)
)

func (s SealType) String() string {
	return string(s)
}
