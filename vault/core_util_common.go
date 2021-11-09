package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

const vaultVersionPath string = "core/versions/"

// storeVersionTimestamp will store the version and timestamp pair to storage only if no entry
// for that version already exists in storage.
func (c *Core) storeVersionTimestamp(ctx context.Context, version string, currentTime time.Time) (bool, error) {
	timeStamp, err := c.barrier.Get(ctx, vaultVersionPath+version)
	if err != nil {
		return false, err
	}

	if timeStamp != nil {
		return false, nil
	}

	vaultVersion := VaultVersion{TimestampInstalled: currentTime, Version: version}
	marshalledVaultVersion, err := json.Marshal(vaultVersion)
	if err != nil {
		return false, err
	}

	err = c.barrier.Put(ctx, &logical.StorageEntry{
		Key:   vaultVersionPath + version,
		Value: marshalledVaultVersion,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindOldestVersionTimestamp searches for the vault version with the oldest
// upgrade timestamp from storage. The earliest version this can be (barring
// downgrades) is 1.9.0.
func (c *Core) FindOldestVersionTimestamp() (string, time.Time, error) {
	if c.VersionTimestamps == nil || len(c.VersionTimestamps) == 0 {
		return "", time.Time{}, fmt.Errorf("version timestamps are not initialized")
	}

	// initialize oldestUpgradeTime to current time
	oldestUpgradeTime := time.Now()
	var oldestVersion string
	for version, upgradeTime := range c.VersionTimestamps {
		if upgradeTime.Before(oldestUpgradeTime) {
			oldestVersion = version
			oldestUpgradeTime = upgradeTime
		}
	}
	return oldestVersion, oldestUpgradeTime, nil
}

// loadVersionTimestamps loads all the vault versions and associated
// upgrade timestamps from storage.
func (c *Core) loadVersionTimestamps(ctx context.Context) (retErr error) {
	vaultVersions, err := c.barrier.List(ctx, vaultVersionPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve vault versions from storage: %+w", err)
	}

	for _, versionPath := range vaultVersions {
		version, err := c.barrier.Get(ctx, vaultVersionPath+versionPath)
		if err != nil {
			return fmt.Errorf("unable to read vault version at path %s: err %+w", versionPath, err)
		}
		if version == nil {
			return fmt.Errorf("nil version stored at path %s", versionPath)
		}
		var vaultVersion VaultVersion
		err = json.Unmarshal(version.Value, &vaultVersion)
		if err != nil {
			return fmt.Errorf("unable to unmarshal vault version for path %s: err %w", versionPath, err)
		}
		if vaultVersion.Version == "" || vaultVersion.TimestampInstalled.IsZero() {
			return fmt.Errorf("found empty serialized vault version at path %s", versionPath)
		}
		c.VersionTimestamps[vaultVersion.Version] = vaultVersion.TimestampInstalled
	}
	return nil
}
