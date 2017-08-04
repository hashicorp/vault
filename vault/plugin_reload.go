package vault

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
)

// reloadPluginMounts reloads provided mounts, regardless of
// plugin name, as long as the backend type is a plugin.
func (c *Core) reloadPluginMounts(mounts []string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	for _, mount := range mounts {
		c.unloadMounts
	}

	return nil
}

// reloadPlugin reloads all mounted backends that are of
// plugin pluginName (name of the plugin as registered in
// the plugin catalog).
func (c *Core) reloadPlugin(pluginName string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Filter mount entries that only matches the plugin name
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

			// Create a barrier view using the UUID
			view = NewBarrierView(c.barrier, barrierPath)

			sysView := c.mountEntrySysView(entry)
			conf := make(map[string]string)
			if entry.Config.PluginName != "" {
				conf["plugin_name"] = entry.Config.PluginName
			}

			// Dispense a new backend
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
				return fmt.Errorf("cannot reload '%s' of type '%s' as a logical backend", entry.Config.PluginName, backendType)
			}

			// Call initialize; this takes care of init tasks that must be run after
			// the ignore paths are collected.
			if err := backend.Initialize(); err != nil {
				return err
			}

			// Set the backend back
			re.backend = backend

			c.logger.Info("core: successfully reloaded '%s' plugin on %s", pluginName, entry.Path)
		}
	}
	return nil
}
