package vault

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
)

// reloadPluginMounts reloads provided mount, regardless of
// plugin type, as long as the backend is a plugin backend.
func (c *Core) reloadPluginMounts(mounts []string) error {
	return nil
}

// reloadPlugin reloads all mounted backends that are of
// plugin type pluginName.
func (c *Core) reloadPlugin(pluginName string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Get all mounts for the plugin
	for _, entry := range c.mounts.Entries {
		if entry.Type == "plugin" && entry.Config.PluginName == pluginName {
			path := entry.Path

			// Fast-path out if the backend doesn't exist
			raw, ok := c.router.root.Get(path)
			if !ok {
				return nil
			}

			// Call backend's Cleanup routine
			re := raw.(*routeEntry)
			re.backend.Cleanup()

			var view *BarrierView

			// Initialize the backend, special casing for system
			barrierPath := backendBarrierPrefix + entry.UUID + "/"
			if entry.Type == "system" {
				barrierPath = systemBarrierPrefix
			}
			// Create a barrier view using the UUID
			view = NewBarrierView(c.barrier, barrierPath)

			// Dispense a new backend
			sysView := c.mountEntrySysView(entry)
			conf := make(map[string]string)
			if entry.Config.PluginName != "" {
				conf["plugin_name"] = entry.Config.PluginName
			}

			// Consider having plugin name under entry.Options
			backend, err := c.newLogicalBackend(entry.Type, sysView, view, conf)
			if err != nil {
				return err
			}
			if backend == nil {
				return fmt.Errorf("nil backend of type %q returned from creation function", entry.Type)
			}

			// Check for the correct backend type
			backendType := backend.Type()
			if entry.Type == "plugin" && backendType != logical.TypeLogical {
				return fmt.Errorf("cannot mount '%s' of type '%s' as a logical backend", entry.Config.PluginName, backendType)
			}

			// Call initialize; this takes care of init tasks that must be run after
			// the ignore paths are collected.
			if err := backend.Initialize(); err != nil {
				return err
			}

			// Set the backend back
			re.backend = backend
		}
	}

	return nil
}
