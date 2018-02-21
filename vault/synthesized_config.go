package vault

import (
	"context"
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/imdario/mergo"
)

// SynthesizableConfig holds configuration values that can be synthesized.
// Synthesized values are those that can be configured at the sys/config and
// mount/tune endpoints.
type SynthesizableConfig struct {
	AuditRequestHMACValues []string `json:"audit_request_hmac_values,omitempty" structs:"audit_request_hmac_values" mapstructure:"audit_request_hmac_values"`
}

// ConfigKeys returns the list of field names in the struct as specified in the
// 'structs' tag.
func (s SynthesizableConfig) ConfigKeys() []string {
	st := structs.New(s)
	m := st.Map()

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
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
	var sysConfig SynthesizableConfig
	if entry != nil {
		err := entry.DecodeJSON(&sysConfig)
		if err != nil {
			return fmt.Errorf("failed to decode %s config entry: %v", configKey, err)
		}
	}

	// Iterage through all entries and perform proper merging of sys + mount values
	for _, entry := range table.Entries {
		// TODO: Handle plugin type
		if entry.Type == backendType {
			mergedConfig := sysConfig
			err := mergo.Merge(&mergedConfig, entry.Config.SynthesizableConfig, mergo.WithOverride)
			if err != nil {
				return fmt.Errorf("failed to merge values: %v", err)
			}

			// If there's anything to store, it's an update, otherwise it's a delete
			if !reflect.DeepEqual(mergedConfig, SynthesizableConfig{}) {
				entry.synthesizedConfigCache.Store(configKey, mergedConfig)
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

	for _, entry := range c.mounts.Entries {
		keys := SynthesizableConfig{}.ConfigKeys()

		for _, key := range keys {
			entryKey := fmt.Sprintf("%s/%s/%s", key, "secret", entry.Type)
			l := locksutil.LockForKey(c.synthesizedConfigLocks, entryKey)
			l.Lock()
			err := c.updateMountEntriesCache(ctx, "secret", entry.Type, key)
			if err != nil {
				l.Unlock()
				return err
			}
			l.Unlock()
		}
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
