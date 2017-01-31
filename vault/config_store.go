package vault

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/logical"
)

const (
	// configSubPath is the sub-path used for the config store
	// view. This is nested under the system view.
	configSubPath = "config/"

	// configCacheSize is the number of configs that are kept cached
	configCacheSize = 1024
)

var defaultConfigs = []string{
	"cors",
}

type Config struct {
	Name     string
	Settings map[string]string
}

// ConfigStore is used to provide durable storage of configuration items
type ConfigStore struct {
	view *BarrierView
	lru  *lru.TwoQueueCache
}

// ConfigEntry is used to store a configuration by name
type ConfigEntry struct {
	Version int
	Raw     string
}

// NewConfigStore creates a new ConfigStore that is backed using a given view.
// It used used to durably store and manage named configurations.
func NewConfigStore(view *BarrierView, system logical.SystemView) *ConfigStore {
	c := &ConfigStore{
		view: view,
	}
	if !system.CachingDisabled() {
		cache, _ := lru.New2Q(configCacheSize)
		c.lru = cache
	}

	return c
}

// setupConfigStore is used to initialize the config store
// when the vault is being unsealed.
func (c *Core) setupConfigStore() error {
	// Create a sub-view
	view := c.systemBarrierView.SubView(configSubPath)

	// Create the config store
	c.configStore = NewConfigStore(view, &dynamicSystemView{core: c})

	for _, name := range defaultConfigs {
		config, err := c.configStore.GetConfig(name)
		if err != nil {
			return err
		}

		if config != nil {
			switch config.Name {
			case "cors":
				c.corsConfig.Enable(config.Settings["allowed_origins"])
			}
		}
	}

	return nil
}

// teardownConfigStore is used to reverse setupConfigStore
// when the vault is being sealed.
func (c *Core) teardownConfigStore() error {
	c.configStore = nil
	return nil
}

// SetConfig is used to create or update the given config
func (cs *ConfigStore) SetConfig(c *Config) error {
	defer metrics.MeasureSince([]string{"config", "set_config"}, time.Now())
	if c.Name == "" {
		return fmt.Errorf("config name missing")
	}

	return cs.setConfigInternal(c)
}

func (cs *ConfigStore) setConfigInternal(c *Config) error {
	var entry *logical.StorageEntry
	var err error

	entry, err = logical.StorageEntryJSON(c.Name, c)
	if cs.lru != nil {
		// Update the LRU cache
		cs.lru.Add(c.Name, c)
	}

	if err != nil {
		return fmt.Errorf("failed to create entry: %v", err)
	}
	if err := cs.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist config: %v", err)
	}

	return nil
}

// GetConfig is used to fetch the named config
func (cs *ConfigStore) GetConfig(name string) (*Config, error) {
	defer metrics.MeasureSince([]string{"config", "get_config"}, time.Now())
	if cs.lru != nil {
		// Check for cached config
		if raw, ok := cs.lru.Get(name); ok {
			return raw.(*Config), nil
		}
	}

	// Load the config in
	out, err := cs.view.Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	if out == nil {
		return nil, nil
	}

	config := new(Config)
	err = out.DecodeJSON(config)
	if err != nil {
		return nil, err
	}

	if cs.lru != nil {
		// Update the LRU cache
		cs.lru.Add(name, config)
	}

	return config, nil
}

// ListConfigs is used to list the available configs
func (cs *ConfigStore) ListConfigs() ([]string, error) {
	defer metrics.MeasureSince([]string{"config", "list_configs"}, time.Now())
	// Scan the view, since the config names are the same as the
	// key names.
	keys, err := logical.CollectKeys(cs.view)

	return keys, err
}

// DeleteConfig is used to delete the named config
func (cs *ConfigStore) DeleteConfig(name string) error {
	defer metrics.MeasureSince([]string{"config", "delete_config"}, time.Now())
	if err := cs.view.Delete(name); err != nil {
		return fmt.Errorf("failed to delete config: %v", err)
	}

	if cs.lru != nil {
		// Clear the cache
		cs.lru.Remove(name)
	}
	return nil
}
