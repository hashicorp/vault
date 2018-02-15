package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/locksutil"
)

// SynthesizableConfig holds configuration values that can be synthesized.
// Synthesized values are those that can be configured at the sys/config and
// mount/tune endpoints.
type SynthesizableConfig struct {
	AuditRequestHMACValues []string `json:"audit_request_hmac_values,omitempty" structs:"audit_request_hmac_values,omitempty" mapstructure:"audit_request_hmac_values"`
}

// BackendsConfig holds synthesizable backend configuration
type BackendsConfig struct {
	Auth   map[string]SynthesizableConfig `json:"auth"`
	Secret map[string]SynthesizableConfig `json:"secret"`
	Audit  map[string]SynthesizableConfig `json:"audit"`
}

// updateMountEntriesCache updates the cached values on the mount. Since this
// method fetches value from storage, this should be called with the proper
// locks held, more especifically c.synthesizeConfigLocks and c.mountsLock/c.authLock.
func (c *Core) updateMountEntriesCache(ctx context.Context, logicalType, backendType, configKey string) error {
	var table *MountTable
	var lType string
	switch logicalType {
	case "secret", mountTableType:
		lType = "secret"
		table = c.mounts
	case credentialTableType:
		lType = credentialTableType
		table = c.auth
	}

	// Get sys/config value
	view := c.systemBarrierView.SubView("config/")
	entryKey := fmt.Sprintf("%s/%s/%s", configKey, lType, backendType)

	entry, err := view.Get(ctx, entryKey)
	if err != nil {
		return fmt.Errorf("failed to save %s config: %v", configKey, err)
	}

	// Decode entry into SynthesizableConfig
	var config SynthesizableConfig
	if entry != nil {
		err := entry.DecodeJSON(&config)
		if err != nil {
			return fmt.Errorf("failed to decode %s config entry: %v", configKey, err)
		}
	}

	for _, entry := range table.Entries {
		// TODO: Handle plugin type
		valToStore := config.AuditRequestHMACValues
		if entry.Type == backendType {
			// Get MountConfigValue
			switch configKey {
			case "audit_request_hmac_values":
				mountVal := entry.Config.AuditRequestHMACValues
				if len(mountVal) > 0 {
					// TODO: Proper merging of any typed value
					valToStore = append(valToStore, mountVal...)
				}
			}

			// If there's anything to store, it's an update, otherwise it's a delete
			if len(valToStore) > 0 {
				entry.synthesizedConfigCache.Store(configKey, valToStore)
			} else {
				entry.synthesizedConfigCache.Delete(configKey)
			}
			c.logger.Debug("core: updated cached value in mount entry", "path", entry.Path)
		}
	}

	return nil
}

// setupMountsConfigCache build the sythensized config cache for each of the mount entries
func (c *Core) setupMountsConfigCache(ctx context.Context) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Get the sys-level config value locks
	entryKey := fmt.Sprintf("%s/%s/%s", "audit_request_hmac_values", "secret", "kv")

	// Acquire the lock
	lock := locksutil.LockForKey(c.synthesizedConfigLocks, entryKey)
	lock.Lock()
	defer lock.Unlock()

	for _, entry := range c.mounts.Entries {
		c.updateMountEntriesCache(ctx, "secret", entry.Type, "audit_request_hmac_values")
	}

	if c.logger.IsInfo() {
		c.logger.Info("core: successfully setup cached config on mount entries")
	}

	return nil
}

// setupMountsConfigCache build the sythensized config cache for each of the auth entries
func (c *Core) setupCredentialsConfigCache(ctx context.Context) error {
	return nil
}
