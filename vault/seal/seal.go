// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync/atomic"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/internalshared/configutil"

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

type SealGenerationInfo struct {
	Generation uint64
	Seals      []*configutil.KMS
	rewrapped  atomic.Bool
}

// Validate is used to sanity check the seal generation info being created
func (sgi *SealGenerationInfo) Validate(existingSgi *SealGenerationInfo, hasPartiallyWrappedPaths bool) error {
	existingSealsLen := 0
	previousShamirConfigured := false
	if existingSgi != nil {
		if sgi.Generation == existingSgi.Generation {
			if !cmp.Equal(sgi.Seals, existingSgi.Seals) {
				return errors.New("existing seal generation is the same, but the configured seals are different")
			}
			return nil
		}

		existingSealsLen = len(existingSgi.Seals)
		for _, sealKmsConfig := range existingSgi.Seals {
			if sealKmsConfig.Type == wrapping.WrapperTypeShamir.String() {
				previousShamirConfigured = true
				break
			}
		}

		if !previousShamirConfigured && (!existingSgi.IsRewrapped() || hasPartiallyWrappedPaths) {
			return errors.New("cannot make seal config changes while seal re-wrap is in progress, please revert any seal configuration changes")
		}
	}

	numSealsToAdd := 0
	// With a previously configured shamir seal, we are either going from [shamir]->[auto]
	// or [shamir]->[another shamir] (since we do not allow multiple shamir
	// seals, and, mixed shamir and auto seals). Also, we do not allow shamir seals to
	// be set disabled, so, the number of seals to add is always going to be the length
	// of new seal configs.
	if previousShamirConfigured {
		numSealsToAdd = len(sgi.Seals)
	} else {
		numSealsToAdd = len(sgi.Seals) - existingSealsLen
	}

	numSealsToDelete := existingSealsLen - len(sgi.Seals)
	switch {
	case numSealsToAdd > 1:
		return errors.New("cannot add more than one seal")

	case numSealsToDelete > 1:
		return errors.New("cannot delete more than one seal")

	case !previousShamirConfigured && existingSgi != nil && !haveCommonSeal(existingSgi.Seals, sgi.Seals):
		// With a previously configured shamir seal, we are either going from [shamir]->[auto] or [shamir]->[another shamir],
		// in which case we cannot have a common seal because shamir seals cannot be set to disabled, they can only be deleted.
		return errors.New("must have at least one seal in common with the old generation")
	}
	return nil
}

func haveCommonSeal(existingSealKmsConfigs, newSealKmsConfigs []*configutil.KMS) (result bool) {
	for _, existingSealKmsConfig := range existingSealKmsConfigs {
		for _, newSealKmsConfig := range newSealKmsConfigs {
			// Clone the existing seal config and set 'Disabled' and 'Priority' fields same as the
			// new seal config, because there might be a case where a seal might be disabled in
			// current config, but might be stored as enabled previously, and this still needs to
			// be considered as a common seal.
			clonedSgi := existingSealKmsConfig.Clone()
			clonedSgi.Disabled = newSealKmsConfig.Disabled
			clonedSgi.Priority = newSealKmsConfig.Priority
			if cmp.Equal(clonedSgi, newSealKmsConfig.Clone()) {
				return true
			}
		}
	}
	return false
}

// SetRewrapped updates the SealGenerationInfo's rewrapped status to the provided value.
func (sgi *SealGenerationInfo) SetRewrapped(value bool) {
	sgi.rewrapped.Store(value)
}

// IsRewrapped returns the SealGenerationInfo's rewrapped status.
func (sgi *SealGenerationInfo) IsRewrapped() bool {
	return sgi.rewrapped.Load()
}

type sealGenerationInfoJson struct {
	Generation uint64
	Seals      []*configutil.KMS
	Rewrapped  bool
}

func (sgi *SealGenerationInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(sealGenerationInfoJson{
		Generation: sgi.Generation,
		Seals:      sgi.Seals,
		Rewrapped:  sgi.IsRewrapped(),
	})
}

func (sgi *SealGenerationInfo) UnmarshalJSON(b []byte) error {
	var value sealGenerationInfoJson
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	sgi.Generation = value.Generation
	sgi.Seals = value.Seals
	sgi.SetRewrapped(value.Rewrapped)

	return nil
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
	// if this is the case errors may still be returned if any wrapper failed. The error map is keyed by seal name.
	Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) (*MultiWrapValue, map[string]error)

	// Decrypt decrypts the given byte slice and stores the resulting information in the
	// returned byte slice. Which options are used depends on the underlying wrapper.
	// Supported options: WithAad.
	// Returns the plaintext, a flag indicating whether the ciphertext is up-to-date
	// (according to IsUpToDate), and an error.
	Decrypt(ctx context.Context, ciphertext *MultiWrapValue, options ...wrapping.Option) ([]byte, bool, error)

	// IsUpToDate returns true if a MultiWrapValue is up-to-date. An MultiWrapValue is
	// considered to be up-to-date if its generation matches the Access generation, and if
	// it has a slot with a key ID that match the current key ID of each of the Access
	// wrappers.
	IsUpToDate(ctx context.Context, value *MultiWrapValue, forceKeyIdRefresh bool) (bool, error)

	// GetEnabledWrappers returns all the enabled seal Wrappers, in order of priority.
	GetEnabledWrappers() []wrapping.Wrapper

	SetShamirSealKey([]byte) error
	GetShamirKeyBytes(ctx context.Context) ([]byte, error)

	// GetAllSealWrappersByPriority returns all the SealWrapper for all the seal wrappers, including disabled ones.
	GetAllSealWrappersByPriority() []*SealWrapper

	// GetEnabledSealWrappersByPriority returns the SealWrapper for the enabled seal wrappers.
	GetEnabledSealWrappersByPriority() []*SealWrapper

	// AllSealsWrappersHealthy returns whether all enabled SealWrappers are currently healthy.
	AllSealWrappersHealthy() bool

	GetSealGenerationInfo() *SealGenerationInfo
}

type access struct {
	sealGenerationInfo *SealGenerationInfo
	wrappersByPriority []*SealWrapper
	keyIdSet           keyIdSet
	logger             hclog.Logger
}

var _ Access = (*access)(nil)

func NewAccess(logger hclog.Logger, sealGenerationInfo *SealGenerationInfo, sealWrappers []SealWrapper) Access {
	if logger == nil {
		logger = hclog.NewNullLogger()
	}
	if sealGenerationInfo == nil {
		panic("cannot create a seal.Access without a SealGenerationInfo")
	}
	if len(sealWrappers) == 0 {
		panic("cannot create a seal.Access without any seal wrappers")
	}
	a := &access{
		sealGenerationInfo: sealGenerationInfo,
		logger:             logger,
	}
	a.wrappersByPriority = make([]*SealWrapper, len(sealWrappers))
	for i, sw := range sealWrappers {
		v := sw
		a.wrappersByPriority[i] = &v
		v.Healthy = true
		v.LastSeenHealthy = time.Now()
	}

	sort.Slice(a.wrappersByPriority, func(i int, j int) bool { return a.wrappersByPriority[i].Priority < a.wrappersByPriority[j].Priority })

	return a
}

func NewAccessFromSealWrappers(logger hclog.Logger, generation uint64, rewrapped bool, sealWrappers []SealWrapper) (Access, error) {
	sealGenerationInfo := &SealGenerationInfo{
		Generation: generation,
	}
	sealGenerationInfo.SetRewrapped(rewrapped)
	ctx := context.Background()
	for _, sw := range sealWrappers {
		typ, err := sw.Wrapper.Type(ctx)
		if err != nil {
			return nil, err
		}
		sealGenerationInfo.Seals = append(sealGenerationInfo.Seals, &configutil.KMS{
			Type:     typ.String(),
			Priority: sw.Priority,
			Name:     sw.Name,
		})
	}
	return NewAccess(logger, sealGenerationInfo, sealWrappers), nil
}

func (a *access) GetAllSealWrappersByPriority() []*SealWrapper {
	return copySealWrappers(a.wrappersByPriority, false)
}

func (a *access) GetEnabledSealWrappersByPriority() []*SealWrapper {
	return copySealWrappers(a.wrappersByPriority, true)
}

func (a *access) AllSealWrappersHealthy() bool {
	for _, sw := range a.wrappersByPriority {
		// Ignore disabled seals
		if sw.Disabled {
			continue
		}
		sw.HcLock.RLock()
		defer sw.HcLock.RUnlock()
		if !sw.Healthy {
			return false
		}
	}
	return true
}

func copySealWrappers(sealWrappers []*SealWrapper, enabledOnly bool) []*SealWrapper {
	ret := make([]*SealWrapper, 0, len(sealWrappers))
	for _, sw := range sealWrappers {
		if enabledOnly && sw.Disabled {
			continue
		}
		ret = append(ret, sw)
	}
	return ret
}

func (a *access) GetSealGenerationInfo() *SealGenerationInfo {
	return a.sealGenerationInfo
}

func (a *access) Generation() uint64 {
	return a.sealGenerationInfo.Generation
}

func (a *access) GetEnabledWrappers() []wrapping.Wrapper {
	var ret []wrapping.Wrapper
	for _, si := range a.GetEnabledSealWrappersByPriority() {
		ret = append(ret, si.Wrapper)
	}
	return ret
}

func (a *access) Init(ctx context.Context, options ...wrapping.Option) error {
	var keyIds []string
	for _, sealWrapper := range a.GetAllSealWrappersByPriority() {
		if initWrapper, ok := sealWrapper.Wrapper.(wrapping.InitFinalizer); ok {
			if err := initWrapper.Init(ctx, options...); err != nil {
				return err
			}
			keyId, err := sealWrapper.Wrapper.KeyId(ctx)
			if err != nil {
				a.logger.Warn("cannot determine key ID for seal", "seal", sealWrapper.Name, "err", err)
				return fmt.Errorf("cannod determine key ID for seal %s: %w", sealWrapper.Name, err)
			}
			keyIds = append(keyIds, keyId)
		}
	}
	a.keyIdSet.setIds(keyIds)
	return nil
}

func (a *access) IsUpToDate(ctx context.Context, value *MultiWrapValue, forceKeyIdRefresh bool) (bool, error) {
	// Note that we don't compare generations when the value is transitory, since all single-blobInfo
	// values are unmarshalled as transitory values.
	if value.Generation != 0 && value.Generation != a.Generation() {
		return false, nil
	}
	if forceKeyIdRefresh {
		test, errs := a.Encrypt(ctx, []byte{0})
		if test == nil {
			a.logger.Error("error refreshing seal key IDs")
			return false, JoinSealWrapErrors("cannot determine key IDs of Access wrappers", errs)
		}
		// TODO(SEALHA): What to do if there are partial failures?
		if len(errs) > 0 {
			msg := "could not determine key IDs of some Access wrappers"
			a.logger.Warn(msg)
			a.logger.Trace("partial failure refreshing seal key IDs", "err", JoinSealWrapErrors(msg, errs))
		}
		a.keyIdSet.set(test)
	}

	return a.keyIdSet.equal(value), nil
}

// Encrypt uses the underlying seal to encrypt the plaintext and returns it.
func (a *access) Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) (*MultiWrapValue, map[string]error) {
	var slots []*wrapping.BlobInfo
	errs := make(map[string]error)

	for _, sealWrapper := range a.GetEnabledSealWrappersByPriority() {
		var encryptErr error
		defer func(now time.Time) {
			metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
			metrics.MeasureSince([]string{"seal", sealWrapper.Name, "encrypt", "time"}, now)

			if encryptErr != nil {
				metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
				metrics.IncrCounter([]string{"seal", sealWrapper.Name, "encrypt", "error"}, 1)
			}
		}(time.Now())

		metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
		metrics.IncrCounter([]string{"seal", sealWrapper.Name, "encrypt"}, 1)

		ciphertext, encryptErr := sealWrapper.Wrapper.Encrypt(ctx, plaintext, options...)
		if encryptErr != nil {
			a.logger.Warn("error encrypting with seal", "seal", sealWrapper.Name)
			a.logger.Trace("error encrypting with seal", "seal", sealWrapper.Name, "err", encryptErr)

			errs[sealWrapper.Name] = encryptErr
			sealWrapper.Healthy = false
		} else {
			a.logger.Trace("encrypted value using seal", "seal", sealWrapper.Name, "keyId", ciphertext.KeyInfo.KeyId)

			slots = append(slots, ciphertext)
		}
	}

	if len(slots) == 0 {
		a.logger.Error("all seals failed to encrypt value")
		return nil, errs
	}

	a.logger.Trace("successfully encrypted value", "encryption seal wrappers", len(slots), "total enabled seal wrappers",
		len(a.GetEnabledSealWrappersByPriority()))
	ret := &MultiWrapValue{
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
func (a *access) Decrypt(ctx context.Context, ciphertext *MultiWrapValue, options ...wrapping.Option) ([]byte, bool, error) {
	blobInfoMap := slotsByKeyId(ciphertext)

	isUpToDate, err := a.IsUpToDate(ctx, ciphertext, false)
	if err != nil {
		return nil, false, err
	}

	// First, lets try the wrappers in order of priority and look for an exact key ID match
	for _, sealWrapper := range a.GetAllSealWrappersByPriority() {
		if keyId, err := sealWrapper.Wrapper.KeyId(ctx); err == nil {
			if blobInfo, ok := blobInfoMap[keyId]; ok {
				pt, oldKey, err := a.tryDecrypt(ctx, sealWrapper, blobInfo, options)
				if oldKey {
					a.logger.Trace("decrypted using OldKey", "seal", sealWrapper.Name)
					return pt, false, err
				}
				if err == nil {
					a.logger.Trace("decrypted value using seal", "seal", sealWrapper.Name)
					return pt, isUpToDate, nil
				}
				// If there is an error, keep trying with the other wrappers
				a.logger.Trace("error decrypting with seal, will try other seals", "seal", sealWrapper.Name, "keyId", keyId, "err", err)
			}
		}
	}

	// No key ID match, so try each wrapper with all slots
	errs := make(map[string]error)
	for _, sealWrapper := range a.GetAllSealWrappersByPriority() {
		for _, blobInfo := range ciphertext.Slots {
			pt, oldKey, err := a.tryDecrypt(ctx, sealWrapper, blobInfo, options)
			if oldKey {
				a.logger.Trace("decrypted using OldKey", "seal", sealWrapper.Name)
				return pt, false, err
			}
			if err == nil {
				a.logger.Trace("decrypted value using seal", "seal", sealWrapper.Name)
				return pt, isUpToDate, nil
			}
			errs[sealWrapper.Name] = err
		}
	}

	return nil, false, JoinSealWrapErrors("error decrypting seal wrapped value", errs)
}

func (a *access) tryDecrypt(ctx context.Context, sealWrapper *SealWrapper, ciphertext *wrapping.BlobInfo, options []wrapping.Option) ([]byte, bool, error) {
	var decryptErr error
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", sealWrapper.Name, "decrypt", "time"}, now)

		if decryptErr != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", sealWrapper.Name, "decrypt", "error"}, 1)
		}
		// TODO (multiseal): log an error?
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", sealWrapper.Name, "decrypt"}, 1)

	pt, err := sealWrapper.Wrapper.Decrypt(ctx, ciphertext, options...)
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

	for _, w := range a.GetAllSealWrappersByPriority() {
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
		return errors.New("no wrappers configured")
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

func slotsByKeyId(value *MultiWrapValue) map[string]*wrapping.BlobInfo {
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

func (s *keyIdSet) set(value *MultiWrapValue) {
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

func (s *keyIdSet) equal(value *MultiWrapValue) bool {
	keyIds := s.collect(value)
	expected := s.get()
	return reflect.DeepEqual(keyIds, expected)
}

func (s *keyIdSet) collect(value *MultiWrapValue) []string {
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
