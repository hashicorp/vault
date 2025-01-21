// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package activationflags

import (
	"context"
	"fmt"
	"maps"
	"sync"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	storagePathActivationFlags = "activation-flags"
	IdentityDeduplication      = "force-identity-deduplication"
)

type FeatureActivationFlags struct {
	activationFlagsLock sync.RWMutex
	storage             logical.Storage
	activationFlags     map[string]bool
}

func NewFeatureActivationFlags() *FeatureActivationFlags {
	return &FeatureActivationFlags{
		activationFlags: map[string]bool{},
	}
}

func (f *FeatureActivationFlags) Initialize(ctx context.Context, storage logical.Storage) error {
	f.activationFlagsLock.Lock()
	defer f.activationFlagsLock.Unlock()

	if storage == nil {
		return fmt.Errorf("unable to access storage")
	}

	f.storage = storage

	entry, err := f.storage.Get(ctx, storagePathActivationFlags)
	if err != nil {
		return fmt.Errorf("failed to get activation flags from storage: %w", err)
	}
	if entry == nil {
		f.activationFlags = map[string]bool{}
		return nil
	}

	var activationFlags map[string]bool
	if err := entry.DecodeJSON(&activationFlags); err != nil {
		return fmt.Errorf("failed to decode activation flags from storage: %w", err)
	}

	f.activationFlags = activationFlags

	return nil
}

// Get is the helper function called by the activation-flags API read endpoint. This reads the
// actual values from storage, then updates the in-memory cache of the activation-flags. It
// returns a slice of the feature names which have already been activated.
func (f *FeatureActivationFlags) Get(ctx context.Context) ([]string, error) {
	// Don't use nil slice declaration, we want the JSON to show "[]" instead of null
	activated := []string{}

	_, err := f.ReloadFlagsFromStorage(ctx)
	if err != nil {
		return activated, err
	}

	f.activationFlagsLock.Lock()
	defer f.activationFlagsLock.Unlock()

	for flag, set := range f.activationFlags {
		if set {
			activated = append(activated, flag)
		}
	}

	return activated, nil
}

func (f *FeatureActivationFlags) ReloadFlagsFromStorage(ctx context.Context) (map[string]bool, error) {
	f.activationFlagsLock.Lock()
	defer f.activationFlagsLock.Unlock()

	if f.storage == nil {
		return map[string]bool{}, nil
	}

	entry, err := f.storage.Get(ctx, storagePathActivationFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to get activation flags from storage: %w", err)
	}
	if entry == nil {
		return map[string]bool{}, nil
	}

	var storageActivationFlags map[string]bool
	if err := entry.DecodeJSON(&storageActivationFlags); err != nil {
		return nil, fmt.Errorf("failed to decode activation flags from storage: %w", err)
	}

	// State Change Logic for Flags
	//
	// This logic determines changes to flags, but it does NOT account for flags that have been deleted.
	// As of this writing, flag removal is not supported for activation flags.
	//
	// Valid State Transitions:
	// 1. Unset (new flag) -> Active
	// 2. Active -> Inactive
	// 3. Inactive -> Active
	//
	// Behavior notes:
	// - If a flag does not exist in-memory (`!ok`), it is treated as a new flag.
	//   Nodes should only react to the new flag if its state is being set to "Active".
	// - If a flag exists in-memory, any change in its value (e.g., Active -> Inactive) is considered valid
	//   and is marked as a state change.
	//
	// The resulting `changedFlags` map will store the flags with their new values if they meet the above criteria.
	changedFlags := map[string]bool{}
	for flg, v := range storageActivationFlags {
		oldValue, ok := f.activationFlags[flg]

		switch {
		// New flag: handle only if transitioning to "Active (true)"
		case !ok && v:
			changedFlags[flg] = v
		default:
			// Existing flag: detect state change
			if oldValue != v {
				changedFlags[flg] = v
			}
		}
	}

	// Update the in-memory flags after loading the latest values from storage
	f.activationFlags = storageActivationFlags

	return changedFlags, nil
}

// Write is the helper function called by the activation-flags API write endpoint. This stores
// the boolean value for the activation-flag feature name into Vault storage across the cluster
// and updates the in-memory cache upon success.
func (f *FeatureActivationFlags) Write(ctx context.Context, featureName string, activate bool) (err error) {
	f.activationFlagsLock.Lock()
	defer f.activationFlagsLock.Unlock()

	if f.storage == nil {
		return fmt.Errorf("unable to access storage")
	}

	activationFlags := f.activationFlags

	clonedFlags := maps.Clone(f.activationFlags)
	clonedFlags[featureName] = activate
	// The cloned flags are updated but the in-memory state is only updated on success of the storage update.
	defer func() {
		if err == nil {
			activationFlags[featureName] = activate
		}
	}()

	entry, err := logical.StorageEntryJSON(storagePathActivationFlags, clonedFlags)
	if err != nil {
		return fmt.Errorf("failed to marshal object to JSON: %w", err)
	}

	err = f.storage.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to save object in storage: %w", err)
	}

	return nil
}

// IsActivationFlagEnabled is true if the specified flag is enabled in the core.
func (f *FeatureActivationFlags) IsActivationFlagEnabled(featureName string) bool {
	f.activationFlagsLock.RLock()
	defer f.activationFlagsLock.RUnlock()

	activated, ok := f.activationFlags[featureName]

	return ok && activated
}
