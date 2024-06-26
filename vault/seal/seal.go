// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-kms-wrapping/v2/aead"
	"github.com/hashicorp/vault/internalshared/configutil"
)

type StoredKeysSupport int

const (
	// The 0 value of StoredKeysSupport is an invalid option
	StoredKeysInvalid StoredKeysSupport = iota
	StoredKeysNotSupported
	StoredKeysSupportedGeneric
	StoredKeysSupportedShamirRoot
)

var (
	ErrUnconfiguredWrapper  = errors.New("unconfigured wrapper")
	ErrNoHealthySeals       = errors.New("no healthy seals!")
	ErrNoConfiguredSeals    = errors.New("no configured seals")
	ErrNoSealGenerationInfo = errors.New("no seal generation info")
	ErrNoSeals              = errors.New("no seals provided in the configuration")
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
	Enabled    bool
}

// Validate is used to sanity check the seal generation info being created
func (sgi *SealGenerationInfo) Validate(existingSgi *SealGenerationInfo, hasPartiallyWrappedPaths bool) error {
	existingSealsLen := 0
	numConfiguredSeals := len(sgi.Seals)
	configuredSealNameAndType := sealNameAndTypeAsStr(sgi.Seals)

	// If no previous generation info exists, make sure we perform the initial migration/setup
	// check for enabled configured seals to allow an old style seal migration configuration
	if existingSgi == nil {
		if numConfiguredSeals > 1 {
			return fmt.Errorf("Initializing a cluster or enabling multi-seal on an existing "+
				"cluster must occur with a single seal before adding additional seals\n"+
				"Configured seals: %v", configuredSealNameAndType)
		}

		// No point in comparing anything more as we don't have any information around the
		// existing seal if any actually existed
		return nil
	}

	// Validate that we're in a safe spot with respect to disabling multiseal
	if existingSgi.Enabled && !sgi.Enabled {
		if len(existingSgi.Seals) > 1 {
			return fmt.Errorf("multi-seal is disabled but previous configuration had multiple seals.  re-enable and migrate to a single seal before disabling multi-seal")
		} else if !existingSgi.IsRewrapped() {
			return fmt.Errorf("multi-seal is disabled but previous storage was not fully re-wrapped, re-enable multi-seal and allow rewrapping to complete before disabling multi-seal")
		}
	}

	existingSealNameAndType := sealNameAndTypeAsStr(existingSgi.Seals)
	previousShamirConfigured := false

	if sgi.Generation == existingSgi.Generation {
		if !haveMatchingSeals(sgi.Seals, existingSgi.Seals) {
			return fmt.Errorf("existing seal generation is the same, but the configured seals are different\n"+
				"Existing seals: %v\n"+
				"Configured seals: %v", existingSealNameAndType, configuredSealNameAndType)
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

	if !previousShamirConfigured && (!existingSgi.IsRewrapped() || hasPartiallyWrappedPaths) && os.Getenv("VAULT_SEAL_REWRAP_SAFETY") != "disable" {
		return errors.New("cannot make seal config changes while seal re-wrap is in progress, please revert any seal configuration changes")
	}

	numSealsToAdd := 0
	// With a previously configured shamir seal, we are either going from [shamir]->[auto]
	// or [shamir]->[another shamir] (since we do not allow multiple shamir
	// seals, and, mixed shamir and auto seals). Also, we do not allow shamir seals to
	// be set disabled, so, the number of seals to add is always going to be the length
	// of new seal configs.
	if previousShamirConfigured {
		numSealsToAdd = numConfiguredSeals
	} else {
		numSealsToAdd = numConfiguredSeals - existingSealsLen
	}

	numSealsToDelete := existingSealsLen - numConfiguredSeals
	switch {
	case numSealsToAdd > 1:
		return fmt.Errorf("cannot add more than one seal\n"+
			"Existing seals: %v\n"+
			"Configured seals: %v", existingSealNameAndType, configuredSealNameAndType)

	case numSealsToDelete > 1:
		return fmt.Errorf("cannot delete more than one seal\n"+
			"Existing seals: %v\n"+
			"Configured seals: %v", existingSealNameAndType, configuredSealNameAndType)

	case !previousShamirConfigured && existingSgi != nil && !haveCommonSeal(existingSgi.Seals, sgi.Seals):
		// With a previously configured shamir seal, we are either going from [shamir]->[auto] or [shamir]->[another shamir],
		// in which case we cannot have a common seal because shamir seals cannot be set to disabled, they can only be deleted.
		return fmt.Errorf("must have at least one seal in common with the old generation\n"+
			"Existing seals: %v\n"+
			"Configured seals: %v", existingSealNameAndType, configuredSealNameAndType)
	}
	return nil
}

func sealNameAndTypeAsStr(seals []*configutil.KMS) string {
	info := []string{}
	for _, seal := range seals {
		info = append(info, fmt.Sprintf("Name: %s Type: %s", seal.Name, seal.Type))
	}
	return fmt.Sprintf("[%s]", strings.Join(info, ", "))
}

// haveMatchingSeals verifies that we have the corresponding matching seals by name and type, config and other
// properties are ignored in the comparison
func haveMatchingSeals(existingSealKmsConfigs, newSealKmsConfigs []*configutil.KMS) bool {
	if len(existingSealKmsConfigs) != len(newSealKmsConfigs) {
		return false
	}

	for _, existingSealKmsConfig := range existingSealKmsConfigs {
		found := false
		for _, newSealKmsConfig := range newSealKmsConfigs {
			if cmp.Equal(existingSealKmsConfig, newSealKmsConfig, compareKMSConfigByNameAndType()) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

// haveCommonSeal verifies that we have at least one matching seal across
// the inputs by name and type, config and other properties are ignored in
// the comparison
func haveCommonSeal(existingSealKmsConfigs, newSealKmsConfigs []*configutil.KMS) bool {
	for _, existingSealKmsConfig := range existingSealKmsConfigs {
		for _, newSealKmsConfig := range newSealKmsConfigs {
			// Technically we might be matching the "wrong" seal if the old seal was renamed to
			// "transit-disabled" and we have a new seal named transit. There isn't any way for
			// us to properly distinguish between them
			if cmp.Equal(existingSealKmsConfig, newSealKmsConfig, compareKMSConfigByNameAndType()) {
				return true
			}
		}
	}

	// We might have renamed a disabled seal that was previously used so attempt to match by
	// removing the "-disabled" suffix
	for _, seal := range findRenamedDisabledSeals(newSealKmsConfigs) {
		clonedSeal := seal.Clone()
		clonedSeal.Name = strings.TrimSuffix(clonedSeal.Name, configutil.KmsRenameDisabledSuffix)

		for _, existingSealKmsConfig := range existingSealKmsConfigs {
			if cmp.Equal(existingSealKmsConfig, clonedSeal, compareKMSConfigByNameAndType()) {
				return true
			}
		}
	}

	return false
}

func findRenamedDisabledSeals(configs []*configutil.KMS) []*configutil.KMS {
	disabledSeals := []*configutil.KMS{}
	for _, seal := range configs {
		if seal.Disabled && strings.HasSuffix(seal.Name, configutil.KmsRenameDisabledSuffix) {
			disabledSeals = append(disabledSeals, seal)
		}
	}
	return disabledSeals
}

func compareKMSConfigByNameAndType() cmp.Option {
	// We only match based on name and type to avoid configuration changes such
	// as a Vault token change in the config map from eliminating the match and
	// preventing startup on a matching seal.
	return cmp.Comparer(func(a, b *configutil.KMS) bool {
		return a.Name == b.Name && a.Type == b.Type
	})
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
	Enabled    bool
}

func (sgi *SealGenerationInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(sealGenerationInfoJson{
		Generation: sgi.Generation,
		Seals:      sgi.Seals,
		Rewrapped:  sgi.IsRewrapped(),
		Enabled:    sgi.Enabled,
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
	sgi.Enabled = value.Enabled

	return nil
}

// OldKey is used as a return value from Decrypt to indicate that the old
// key was used for decryption and that the value should be re-encrypted
// with the new key and saved. It is not returned as an error by any
// function.
var OldKey = errors.New("decrypted with old key")

func IsOldKeyError(err error) bool {
	return errors.Is(err, OldKey)
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

	// GetAllSealWrappersByPriority returns all the SealWrappers including disabled and unconfigured wrappers.
	GetAllSealWrappersByPriority() []*SealWrapper

	// GetConfiguredSealWrappersByPriority returns all the configured SealWrappers for all the seal wrappers, including disabled ones.
	GetConfiguredSealWrappersByPriority() []*SealWrapper

	// GetEnabledSealWrappersByPriority returns the SealWrappers for the enabled seal wrappers.
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

func NewAccess(logger hclog.Logger, sealGenerationInfo *SealGenerationInfo, sealWrappers []*SealWrapper) (Access, error) {
	if logger == nil {
		logger = hclog.NewNullLogger()
	}
	if sealGenerationInfo == nil {
		logger.Error("cannot create a seal.Access without a SealGenerationInfo")
		return nil, ErrNoSealGenerationInfo
	}
	if len(sealWrappers) == 0 {
		logger.Error("cannot create a seal.Access without any seal wrappers")
		return nil, ErrNoSeals
	}
	a := &access{
		sealGenerationInfo: sealGenerationInfo,
		logger:             logger,
	}
	a.wrappersByPriority = make([]*SealWrapper, len(sealWrappers))
	for i, sw := range sealWrappers {
		a.wrappersByPriority[i] = sw
	}

	configuredSealWrappers := a.GetConfiguredSealWrappersByPriority()
	if len(configuredSealWrappers) == 0 {
		a.logger.Error("cannot create a seal.Access without any configured seal wrappers")
		return nil, ErrNoConfiguredSeals
	}

	sort.Slice(a.wrappersByPriority, func(i int, j int) bool { return a.wrappersByPriority[i].Priority < a.wrappersByPriority[j].Priority })

	return a, nil
}

func NewAccessFromSealWrappers(logger hclog.Logger, generation uint64, rewrapped bool, sealWrappers []*SealWrapper) (Access, error) {
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
	return NewAccess(logger, sealGenerationInfo, sealWrappers)
}

// NewAccessFromWrapper creates an enabled Access for a single wrapping.Wrapper.
// The Access has generation set to 1 and the rewrapped flag set to true.
// The SealWrapper created uses the seal config type as the name, has priority set to 1 and the
// disabled flag set to false.
func NewAccessFromWrapper(logger hclog.Logger, wrapper wrapping.Wrapper, sealConfigType string) (Access, error) {
	sealWrapper := NewSealWrapper(wrapper, 1, sealConfigType, sealConfigType, false, true)

	return NewAccessFromSealWrappers(logger, 1, true, []*SealWrapper{sealWrapper})
}

func (a *access) GetAllSealWrappersByPriority() []*SealWrapper {
	return a.filterSealWrappers(allWrappers)
}

func (a *access) GetConfiguredSealWrappersByPriority() []*SealWrapper {
	return a.filterSealWrappers(configuredWrappers, allWrappers)
}

func (a *access) GetEnabledSealWrappersByPriority() []*SealWrapper {
	return a.filterSealWrappers(configuredWrappers, enabledWrappers)
}

func (a *access) AllSealWrappersHealthy() bool {
	return len(a.wrappersByPriority) == len(a.filterSealWrappers(configuredWrappers, healthyWrappers))
}

type sealWrapperFilter func(*SealWrapper) bool

func allWrappers(wrapper *SealWrapper) bool {
	return true
}

func healthyWrappers(wrapper *SealWrapper) bool {
	return wrapper.IsHealthy()
}

func unhealthyWrappers(wrapper *SealWrapper) bool {
	return !wrapper.IsHealthy()
}

func enabledWrappers(wrapper *SealWrapper) bool {
	return !wrapper.Disabled
}

func configuredWrappers(wrapper *SealWrapper) bool {
	return wrapper.Configured
}

// Returns a slice of wrappers who satisfy all filters
func (a *access) filterSealWrappers(filters ...sealWrapperFilter) []*SealWrapper {
	return filterSealWrappers(a.wrappersByPriority, filters...)
}

func filterSealWrappers(wrappers []*SealWrapper, filters ...sealWrapperFilter) []*SealWrapper {
	ret := make([]*SealWrapper, 0, len(wrappers))
outer:
	for _, sw := range wrappers {
		for _, f := range filters {
			if !f(sw) {
				continue outer
			}
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
	for _, sealWrapper := range a.GetConfiguredSealWrappersByPriority() {
		if initWrapper, ok := sealWrapper.Wrapper.(wrapping.InitFinalizer); ok {
			if err := initWrapper.Init(ctx, options...); err != nil {
				return err
			}
			keyId, err := sealWrapper.Wrapper.KeyId(ctx)
			if err != nil {
				a.logger.Warn("cannot determine key ID for seal", "seal", sealWrapper.Name, "err", err)
				return fmt.Errorf("cannod determine key ID for seal %s: %w", sealWrapper.Name, err)
			}
			if keyId != "" {
				// Some wrappers may not yet know their key id. For emample, see gcpkms.Wrapper.
				keyIds = append(keyIds, keyId)
			}
		}
	}
	a.keyIdSet.setIds(keyIds)
	return nil
}

func (a *access) IsUpToDate(ctx context.Context, value *MultiWrapValue, forceKeyIdRefresh bool) (bool, error) {
	// Note that we don't compare generations when the value is transitory, since all single-blobInfo
	// values (i.e. not yet upgraded to MultiWrapValues) are unmarshalled as transitory values.
	if value.Generation != 0 && value.Generation != a.Generation() {
		return false, nil
	}
	if forceKeyIdRefresh {
		test, errs := a.Encrypt(ctx, []byte{0})
		if test == nil {
			a.logger.Error("error refreshing seal key IDs")
			return false, JoinSealWrapErrors("cannot determine key IDs of Access wrappers", errs)
		}
		if len(errs) > 0 {
			a.logger.Warn("cannot determine if seal wrapped entry needs update: there were errors determining the key IDs for one or more seals")
			a.logger.Debug("cannot determine if seal wrapped entry needs update", "err", JoinSealWrapErrors("error refreshing key IDs of Access wrappers", errs))

			// Return true, since the encrypted values cannot be re-encrypted without
			// losing the ciphertext of unhealthy wrappers.
			return true, nil
		}
	} else if !a.keyIdSet.initialized() {
		// Since the key ID set is not initialized, we cannot determine if the value is up-to-date, so assume it is.
		// Note that we cannot just force an update, since that breaks migrations to a Shamir defaultSeal.
		return true, nil
	}

	return a.keyIdSet.equal(value), nil
}

const (
	// wrapperEncryptTimeout is the duration we will wait for seal wrappers to return from an encrypt call.
	// After the timeout, we return any successful results and errors for the rest of the wrappers, so
	// that a partial seal wrap entry can be recorded.
	wrapperEncryptTimeout = 10 * time.Second

	// wrapperDecryptHighPriorityHeadStart is the duration we wait for the highest priority wrapper
	// to return from a decrypt call before we try decrypting with any additional wrappers.
	wrapperDecryptHighPriorityHeadStart = 2 * time.Second
)

// Encrypt uses the underlying seal to encrypt the plaintext and returns it.
func (a *access) Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) (*MultiWrapValue, map[string]error) {
	errs := make(map[string]error)

	// Note that we do not encrypt with disabled wrappers. Disabled wrappers are only used to decrypt.
	candidateWrappers := a.filterSealWrappers(enabledWrappers, healthyWrappers)
	if len(candidateWrappers) > 0 {
		// As there are healthy wrappers, add errors for any unhealthy ones, so that it
		// it is clear that the resulting MultiWrapValue is missing ciphertext for some seals.
		for i, unhealthyWrapper := range a.filterSealWrappers(enabledWrappers, unhealthyWrappers) {
			var keyId string
			if unhealthyWrapper.Wrapper != nil {
				// Annoying, apparently Wrapper may be null, see setSeal() in server.go,
				// in the config seal loop.
				keyId, _ = unhealthyWrapper.Wrapper.KeyId(ctx)
			}
			if keyId == "" {
				keyId = unhealthyWrapper.Name
				if _, duplicated := errs[keyId]; duplicated {
					keyId = fmt.Sprintf("%s-%d", keyId, i)
				}
			}
			errs[keyId] = errors.New("seal is unhealthy")
		}
	} else {
		// If all seals are unhealthy, try with all of them since a seal may have recovered.
		candidateWrappers = a.filterSealWrappers(enabledWrappers)
	}
	enabledWrappersByPriority := filterSealWrappers(candidateWrappers, configuredWrappers)

	type result struct {
		name       string
		ciphertext *wrapping.BlobInfo
		err        error
	}
	resultCh := make(chan *result)

	encryptCtx, cancelEncryptCtx := context.WithTimeout(ctx, wrapperEncryptTimeout)
	defer cancelEncryptCtx()

	// Start goroutines to encrypt the value using each of the wrappers.
	for _, sealWrapper := range enabledWrappersByPriority {
		go func(sealWrapper *SealWrapper) {
			ciphertext, err := a.tryEncrypt(encryptCtx, sealWrapper, plaintext, options...)
			resultCh <- &result{
				name:       sealWrapper.Name,
				ciphertext: ciphertext,
				err:        err,
			}
		}(sealWrapper)
	}

	results := make(map[string]*result)
GATHER_RESULTS:
	for {
		select {
		case result := <-resultCh:
			results[result.name] = result
			if len(results) == len(enabledWrappersByPriority) {
				break GATHER_RESULTS
			}
		case <-encryptCtx.Done():
			break GATHER_RESULTS
		case <-ctx.Done():
			cancelEncryptCtx()
			break GATHER_RESULTS
		}
	}

	{
		// Check for duplicate Key IDs.
		// If any wrappers produce duplicated IDs, their BlobInfo will be replaced by an error.

		keyIdToSealWrapperNameMap := make(map[string]string)
		for _, sealWrapper := range enabledWrappersByPriority {
			wrapperName := sealWrapper.Name
			if result, ok := results[wrapperName]; ok {
				if result.err != nil {
					continue
				}
				if result.ciphertext.KeyInfo == nil {
					// Can this really happen? Probably not?
					continue
				}
				keyId := result.ciphertext.KeyInfo.KeyId
				duplicateWrapperName, isDuplicate := keyIdToSealWrapperNameMap[keyId]
				if isDuplicate {
					for _, name := range []string{wrapperName, duplicateWrapperName} {
						results[name].err = fmt.Errorf("seal %s has returned duplicate key ID %s, key IDs must be unique", name, keyId)
						results[name].ciphertext = nil
					}
				}
				keyIdToSealWrapperNameMap[keyId] = wrapperName
			}
		}
	}

	// Sort out the successful results from the errors
	var slots []*wrapping.BlobInfo
	for _, sealWrapper := range enabledWrappersByPriority {
		if result, ok := results[sealWrapper.Name]; ok {
			if result.err != nil {
				errs[sealWrapper.Name] = result.err
			} else {
				slots = append(slots, result.ciphertext)
			}
		} else {
			if encryptCtx.Err() != nil {
				errs[sealWrapper.Name] = encryptCtx.Err()
			} else {
				// Just being paranoid, encryptCtx.Err() should never be nil in this case
				errs[sealWrapper.Name] = errors.New("context timeout exceeded")
			}
			// This failure did not happen on tryEncrypt, so we must log it here
			a.logger.Trace("error encrypting with seal", "seal", sealWrapper.Name, "err", errs[sealWrapper.Name])
		}
	}

	if len(slots) == 0 {
		a.logger.Error("failed to encrypt value using any seal wrappers")
		return nil, errs
	}

	a.logger.Trace("successfully encrypted value", "encryption seal wrappers", len(slots), "total enabled seal wrappers",
		len(a.GetEnabledSealWrappersByPriority()))

	ret := &MultiWrapValue{
		Generation: a.Generation(),
		Slots:      slots,
	}

	if len(errs) == 0 {
		// cache key IDs
		a.keyIdSet.set(ret)
	}

	// Add errors for unconfigured wrappers
	for _, sw := range candidateWrappers {
		if !sw.Configured {
			errs[sw.Name] = ErrUnconfiguredWrapper
		}
	}

	return ret, errs
}

func (a *access) tryEncrypt(ctx context.Context, sealWrapper *SealWrapper, plaintext []byte, options ...wrapping.Option) (*wrapping.BlobInfo, error) {
	now := time.Now()
	var encryptErr error
	mLabels := []metrics.Label{{Name: "seal_wrapper_name", Value: sealWrapper.Name}}

	defer func(now time.Time) {
		metrics.MeasureSinceWithLabels([]string{"seal", "encrypt", "time"}, now, mLabels)

		if encryptErr != nil {
			metrics.IncrCounterWithLabels([]string{"seal", "encrypt", "error"}, 1, mLabels)
		}
	}(now)

	metrics.IncrCounterWithLabels([]string{"seal", "encrypt"}, 1, mLabels)

	ciphertext, encryptErr := sealWrapper.Wrapper.Encrypt(ctx, plaintext, options...)
	if encryptErr != nil {
		a.logger.Warn("error encrypting with seal", "seal", sealWrapper.Name)
		a.logger.Trace("error encrypting with seal", "seal", sealWrapper.Name, "err", encryptErr)

		sealWrapper.SetHealthy(false, now)
		return nil, encryptErr
	}
	a.logger.Trace("encrypted value using seal", "seal", sealWrapper.Name, "keyId", ciphertext.KeyInfo.KeyId)

	sealWrapper.SetHealthy(true, now)
	return ciphertext, nil
}

// Decrypt uses the underlying seal to decrypt the ciphertext and returns it.
// Note that it is possible depending on the wrapper used that both pt and err
// are populated.
// Returns the plaintext, a flag indicating whether the ciphertext is up-to-date
// (according to IsUpToDate), and an error.
func (a *access) Decrypt(ctx context.Context, ciphertext *MultiWrapValue, options ...wrapping.Option) ([]byte, bool, error) {
	isUpToDate, err := a.IsUpToDate(ctx, ciphertext, false)
	if err != nil {
		return nil, false, err
	}

	wrappersByPriority := a.filterSealWrappers(configuredWrappers, healthyWrappers)
	if len(wrappersByPriority) == 0 {
		// If all seals are unhealthy, try any way since a seal may have recovered
		wrappersByPriority = a.filterSealWrappers(configuredWrappers)
	}
	if len(wrappersByPriority) == 0 {
		return nil, false, ErrNoHealthySeals
	}

	type result struct {
		name   string
		pt     []byte
		oldKey bool
		err    error
	}

	resultCh := make(chan *result)
	var resultWg sync.WaitGroup
	defer func() {
		// Consume all the discarded results
		go func() {
			for range resultCh {
			}
		}()
		resultWg.Wait()
		close(resultCh)
	}()

	reportResult := func(name string, plaintext []byte, oldKey bool, err error) {
		resultCh <- &result{
			name:   name,
			pt:     plaintext,
			oldKey: oldKey,
			err:    err,
		}
		resultWg.Done()
	}

	decrypt := func(sealWrapper *SealWrapper) {
		pt, oldKey, err := a.tryDecrypt(ctx, sealWrapper, ciphertext, options)
		reportResult(sealWrapper.Name, pt, oldKey, err)
	}

	// Start goroutines to decrypt the value
	first := wrappersByPriority[0]
	found := false
outer:
	// This loop finds the highest priority seal with a keyId in common with the blobInfoMap,
	// and ensures we'll use it first.  This should equal the highest priority wrapper in the nominal
	// case, but may not if a seal is unhealthy.  This ensures we try the highest priority healthy
	// seal first if available, and warn if we don't think we have one in common.
	for _, sealWrapper := range wrappersByPriority {
		keyId, err := sealWrapper.Wrapper.KeyId(ctx)
		if err != nil {
			resultWg.Add(1)
			go reportResult(sealWrapper.Name, nil, false, err)
			continue
		}
		if bi := ciphertext.BlobInfoForKeyId(keyId); bi != nil {
			found = true
			first = sealWrapper
			break outer
		}
	}

	if !found {
		a.logger.Warn("while unwrapping, value has no key-id in common with currently healthy seals.  Trying all healthy seals")
	}

	resultWg.Add(1)
	go decrypt(first)
	for _, sealWrapper := range wrappersByPriority {
		sealWrapper := sealWrapper
		if sealWrapper != first {
			timer := time.AfterFunc(wrapperDecryptHighPriorityHeadStart, func() {
				resultWg.Add(1)
				decrypt(sealWrapper)
			})
			defer timer.Stop()
		}
	}

	// Gathering failures, but return right away if there is a successful result
	errs := make(map[string]error)
GATHER_RESULTS:
	for {
		select {
		case result := <-resultCh:
			switch {
			case result.err != nil:
				errs[result.name] = result.err
				if len(errs) == len(wrappersByPriority) {
					break GATHER_RESULTS
				}

			case result.oldKey:
				return result.pt, false, OldKey
			default:
				return result.pt, isUpToDate, nil
			}
		case <-ctx.Done():
			break GATHER_RESULTS
		}
	}

	// No wrapper was able to decrypt the value, return an error
	if len(errs) > 0 {
		return nil, false, JoinSealWrapErrors("error decrypting seal wrapped value", errs)
	}

	if ctx.Err() != nil {
		return nil, false, ctx.Err()
	}
	// Just being paranoid, ctx.Err() should never be nil in this case
	return nil, false, errors.New("context timeout exceeded")
}

// tryDecrypt returns the plaintext and a flag indicating whether the decryption was done by the "unwrapSeal" (see
// sealWrapMigration.Decrypt).
func (a *access) tryDecrypt(ctx context.Context, sealWrapper *SealWrapper, value *MultiWrapValue, options []wrapping.Option) ([]byte, bool, error) {
	now := time.Now()
	var decryptErr error
	mLabels := []metrics.Label{{Name: "seal_wrapper_name", Value: sealWrapper.Name}}

	defer func(now time.Time) {
		metrics.MeasureSinceWithLabels([]string{"seal", "decrypt", "time"}, now, mLabels)

		if decryptErr != nil {
			metrics.IncrCounterWithLabels([]string{"seal", "decrypt", "error"}, 1, mLabels)
		}
	}(now)

	metrics.IncrCounterWithLabels([]string{"seal", "decrypt"}, 1, mLabels)

	var pt []byte

	// First, let's look for an exact key ID match
	var keyId string
	if id, err := sealWrapper.Wrapper.KeyId(ctx); err == nil {
		keyId = id
		if ciphertext := value.BlobInfoForKeyId(keyId); ciphertext != nil {
			pt, decryptErr = sealWrapper.Wrapper.Decrypt(ctx, ciphertext, options...)

			sealWrapper.SetHealthy(decryptErr == nil || IsOldKeyError(decryptErr), now)
		}
	}
	// If we don't get a result, try all the slots
	if pt == nil && decryptErr == nil {
		for _, ciphertext := range value.Slots {
			pt, decryptErr = sealWrapper.Wrapper.Decrypt(ctx, ciphertext, options...)
			if decryptErr == nil {
				// Note that we only update wrapper health for failures on exact key ID match,
				// otherwise we would have false negatives.
				sealWrapper.SetHealthy(true, now)
				break
			}
		}
	}

	switch {
	case decryptErr != nil && IsOldKeyError(decryptErr):
		// an OldKey error is not an actual error, it just means that the decryption was done
		// by the "unwrapSeal" of a seal migration (see sealWrapMigration.Decrypt).
		a.logger.Trace("decrypted using OldKey", "seal_name", sealWrapper.Name)
		return pt, true, nil

	case decryptErr != nil:
		// Note that if there are more than one ciphertext, the error may be misleading...
		a.logger.Trace("error decrypting with seal, this may be a harmless mismatch between wrapper and ciphertext", "seal_name", sealWrapper.Name, "keyId", keyId, "err", decryptErr)
		return nil, false, decryptErr

	default:
		a.logger.Trace("decrypted value using seal", "seal_name", sealWrapper.Name)
		return pt, false, nil
	}
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

	for _, w := range a.GetConfiguredSealWrappersByPriority() {
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

func (v *MultiWrapValue) BlobInfoForKeyId(keyId string) *wrapping.BlobInfo {
	for _, blobInfo := range v.Slots {
		if blobInfo.KeyInfo != nil && blobInfo.KeyInfo.KeyId == keyId {
			return blobInfo
		}
	}
	return nil
}

type keyIdSet struct {
	keyIds atomic.Pointer[[]string]
}

func (s *keyIdSet) initialized() bool {
	return len(s.get()) > 0
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
