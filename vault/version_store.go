package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	vaultVersionPath string = "core/versions/"
)

// storeVersionTimestamp will store the version and timestamp pair to storage
// only if no entry for that version already exists in storage. Version
// timestamps were initially stored in local time. UTC should be used. Existing
// entries can be overwritten via the force flag. A bool will be returned
// denoting whether the entry was updated
func (c *Core) storeVersionTimestamp(ctx context.Context, version string, timestampInstalled time.Time, force bool) (bool, error) {
	key := vaultVersionPath + version

	vaultVersion := VaultVersion{
		TimestampInstalled: timestampInstalled.UTC(),
		Version:            version,
	}

	marshalledVaultVersion, err := json.Marshal(vaultVersion)
	if err != nil {
		return false, err
	}

	newEntry := &logical.StorageEntry{
		Key:   key,
		Value: marshalledVaultVersion,
	}

	if force {
		// avoid storage lookup and write immediately
		err = c.barrier.Put(ctx, newEntry)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	existingEntry, err := c.barrier.Get(ctx, key)
	if err != nil {
		return false, err
	}

	if existingEntry != nil {
		return false, nil
	}

	err = c.barrier.Put(ctx, newEntry)

	if err != nil {
		return false, err
	}

	return true, nil
}

// FindOldestVersionTimestamp searches for the vault version with the oldest
// upgrade timestamp from storage. The earliest version this can be is 1.9.0.
func (c *Core) FindOldestVersionTimestamp() (string, time.Time, error) {
	if c.versionTimestamps == nil {
		return "", time.Time{}, fmt.Errorf("version timestamps are not initialized")
	}

	oldestUpgradeTime := time.Now().UTC()
	var oldestVersion string

	for version, upgradeTime := range c.versionTimestamps {
		if upgradeTime.Before(oldestUpgradeTime) {
			oldestVersion = version
			oldestUpgradeTime = upgradeTime
		}
	}
	return oldestVersion, oldestUpgradeTime, nil
}

func (c *Core) FindNewestVersionTimestamp() (string, time.Time, error) {
	if c.versionTimestamps == nil {
		return "", time.Time{}, fmt.Errorf("version timestamps are not initialized")
	}

	var newestUpgradeTime time.Time
	var newestVersion string

	for version, upgradeTime := range c.versionTimestamps {
		if upgradeTime.After(newestUpgradeTime) {
			newestVersion = version
			newestUpgradeTime = upgradeTime
		}
	}

	return newestVersion, newestUpgradeTime, nil
}

// loadVersionTimestamps loads all the vault versions and associated upgrade
// timestamps from storage. Version timestamps were originally stored in local
// time. A timestamp that is not in UTC will be rewritten to storage as UTC.
func (c *Core) loadVersionTimestamps(ctx context.Context) error {
	vaultVersions, err := c.barrier.List(ctx, vaultVersionPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve vault versions from storage: %w", err)
	}

	for _, versionPath := range vaultVersions {
		version, err := c.barrier.Get(ctx, vaultVersionPath+versionPath)
		if err != nil {
			return fmt.Errorf("unable to read vault version at path %s: err %w", versionPath, err)
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

		timestampInstalled := vaultVersion.TimestampInstalled

		// self-heal entries that were not stored in UTC
		if timestampInstalled.Location() != time.UTC {
			timestampInstalled = timestampInstalled.UTC()
			isUpdated, err := c.storeVersionTimestamp(ctx, vaultVersion.Version, timestampInstalled, true)
			if err != nil {
				c.logger.Warn("failed to rewrite vault version timestamp as UTC", "error", err)
			}

			if isUpdated {
				c.logger.Info("self-healed pre-existing vault version in UTC",
					"vault version", vaultVersion.Version, "UTC time", timestampInstalled)
			}
		}

		c.versionTimestamps[vaultVersion.Version] = timestampInstalled
	}
	return nil
}

func IsJWT(token string) bool {
	return len(token) > 3 && strings.Count(token, ".") == 2 &&
		(token[3] != '.' && token[1] != '.')
}

func IsSSCToken(token string) bool {
	return len(token) > MaxNsIdLength+TokenLength+TokenPrefixLength &&
		strings.HasPrefix(token, consts.ServiceTokenPrefix)
}

func IsServiceToken(token string) bool {
	return strings.HasPrefix(token, consts.ServiceTokenPrefix) ||
		strings.HasPrefix(token, consts.LegacyServiceTokenPrefix)
}

func IsBatchToken(token string) bool {
	return strings.HasPrefix(token, consts.LegacyBatchTokenPrefix) ||
		strings.HasPrefix(token, consts.BatchTokenPrefix)
}
