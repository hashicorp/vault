package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

// reloadPluginMounts reloads provided mounts, regardless of
// plugin name, as long as the backend type is plugin.
func (c *Core) reloadMatchingPluginMounts(ctx context.Context, mounts []string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	var errors error
	for _, mount := range mounts {
		entry := c.router.MatchingMountEntry(mount)
		if entry == nil {
			errors = multierror.Append(errors, fmt.Errorf("cannot fetch mount entry on %q", mount))
			continue
		}

		var isAuth bool
		fullPath := c.router.MatchingMount(mount)
		if strings.HasPrefix(fullPath, credentialRoutePrefix) {
			isAuth = true
		}

		if entry.Type == "plugin" {
			err := c.reloadBackendCommon(ctx, entry, isAuth)
			if err != nil {
				errors = multierror.Append(errors, errwrap.Wrapf(fmt.Sprintf("cannot reload plugin on %q: {{err}}", mount), err))
				continue
			}
			c.logger.Info("successfully reloaded plugin", "plugin", entry.Config.PluginName, "path", entry.Path)
		}
	}
	return errors
}

// reloadPlugin reloads all mounted backends that are of
// plugin pluginName (name of the plugin as registered in
// the plugin catalog).
func (c *Core) reloadMatchingPlugin(ctx context.Context, pluginName string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Filter mount entries that only matches the plugin name
	for _, entry := range c.mounts.Entries {
		if entry.Config.PluginName == pluginName && entry.Type == "plugin" {
			err := c.reloadBackendCommon(ctx, entry, false)
			if err != nil {
				return err
			}
			c.logger.Info("successfully reloaded plugin", "plugin", pluginName, "path", entry.Path)
		}
	}

	// Filter auth mount entries that ony matches the plugin name
	for _, entry := range c.auth.Entries {
		if entry.Config.PluginName == pluginName && entry.Type == "plugin" {
			err := c.reloadBackendCommon(ctx, entry, true)
			if err != nil {
				return err
			}
			c.logger.Info("successfully reloaded plugin", "plugin", pluginName, "path", entry.Path)
		}
	}

	return nil
}

// reloadBackendCommon is a generic method to reload a backend provided a
// MountEntry.
func (c *Core) reloadBackendCommon(ctx context.Context, entry *MountEntry, isAuth bool) error {
	// We don't want to reload the singleton mounts. They often have specific
	// inmemory elements and we don't want to touch them here.
	if strutil.StrListContains(singletonMounts, entry.Type) {
		c.logger.Debug("Skipping reload of singleton mount", "type", entry.Type)
		return nil
	}

	path := entry.Path

	if isAuth {
		path = credentialRoutePrefix + path
	}

	// Fast-path out if the backend doesn't exist
	raw, ok := c.router.root.Get(path)
	if !ok {
		return nil
	}

	re := raw.(*routeEntry)

	// Grab the lock, this allows requests to drain before we cleanup the
	// client.
	re.l.Lock()
	defer re.l.Unlock()

	// Only call Cleanup if backend is initialized
	if re.backend != nil {
		// Call backend's Cleanup routine
		re.backend.Cleanup(ctx)
	}

	view := re.storageView

	sysView := c.mountEntrySysView(entry)

	var backend logical.Backend
	var err error
	if !isAuth {
		// Dispense a new backend
		backend, err = c.newLogicalBackend(ctx, entry, sysView, view)
	} else {
		backend, err = c.newCredentialBackend(ctx, entry, sysView, view)
	}
	if err != nil {
		return err
	}
	if backend == nil {
		return fmt.Errorf("nil backend of type %q returned from creation function", entry.Type)
	}

	// Set the backend back
	re.backend = backend

	// Set paths as well
	paths := backend.SpecialPaths()
	if paths != nil {
		re.rootPaths.Store(pathsToRadix(paths.Root))
		re.loginPaths.Store(pathsToRadix(paths.Unauthenticated))
	}

	return nil
}
