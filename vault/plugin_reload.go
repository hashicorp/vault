// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/internalshared/namespace"
	sdkbackoff "github.com/hashicorp/vault/sdk/helper/backoff"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin"
)

// concurrentReloadWorkers is the maximum number of mounts reloaded in
// parallel within a single plugin reload request.
//
// This count appears to strike a practical balance: it allows many mounts to
// be processed at once without overwhelming storage during the post-reload
// persist phase. Higher worker counts have tended to produce more overload
// retries and more persist failures after retries.
const concurrentReloadWorkers = 32

// dispatchReloads fans out work across entries using a bounded semaphore of
// concurrentReloadWorkers slots. If ctx is cancelled while the coordinator is
// waiting for a free slot, the remaining un-dispatched entries receive the zero
// value of R in their result slot, and the index of the first skipped entry is
// returned so the caller can stamp an appropriate error into those slots.
// Already-dispatched goroutines are always awaited before the function returns.
func dispatchReloads[E any, R any](
	ctx context.Context,
	entries []E,
	work func(ctx context.Context, e E) R,
) (results []R, skippedFrom int) {
	results = make([]R, len(entries))
	skippedFrom = len(entries) // default: nothing skipped

	sem := make(chan struct{}, concurrentReloadWorkers)
	var wg sync.WaitGroup

loop:
	for i, e := range entries {
		wg.Add(1)
		select {
		case sem <- struct{}{}:
			go func(i int, e E) {
				defer wg.Done()
				defer func() { <-sem }()
				results[i] = work(ctx, e)
			}(i, e)
		case <-ctx.Done():
			wg.Done() // balance Add; no goroutine was launched
			skippedFrom = i
			break loop
		}
	}
	wg.Wait()
	return results, skippedFrom
}

const (
	pluginReloadPluginsType = "plugins"
	pluginReloadMountsType  = "mounts"
)

// reloadPersistMaxRetries is the maximum number of times a failed persist
// (persistAuth / persistMounts) is retried when the error is ErrOverloaded.
const reloadPersistMaxRetries = 5

// reloadPersistRetryMin and reloadPersistRetryMax are declared as vars so
// tests can override them without real sleep durations.
var (
	reloadPersistRetryMin = 50 * time.Millisecond
	reloadPersistRetryMax = 5 * time.Second
)

// retryPersistOnOverload calls persist and retries with exponential backoff
// when the returned error is consts.ErrOverloaded. All other errors are logged
// and dropped because the process swap has already succeeded and keeping the
// new backend available takes precedence over immediately persisting its SHA.
//
// Retries run under the plugin reload request context, so this backoff shares
// the same overall request timeout as the reload itself. It is intentionally
// best-effort: enough to smooth transient storage pressure from concurrent
// reloads, but not intended to keep retrying indefinitely once request time
// becomes the tighter bound.
//
// Uses sdk/helper/backoff for consistent backoff behaviour across the codebase.
func retryPersistOnOverload(ctx context.Context, logger log.Logger, persist func() error) {
	b := sdkbackoff.NewBackoff(reloadPersistMaxRetries, reloadPersistRetryMin, reloadPersistRetryMax)
	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	for {
		err := persist()
		if err == nil {
			return
		}
		if !errors.Is(err, consts.ErrOverloaded) {
			logger.Warn("plugin reload: persist failed, SHA metadata may be stale until next reload", "error", err)
			return
		}
		d, backoffErr := b.Next()
		if backoffErr != nil {
			// sdkbackoff.ErrMaxRetry — retries exhausted
			logger.Warn("plugin reload: persist failed after max retries due to storage overload, SHA metadata may be stale until next reload", "error", err)
			return
		}
		logger.Warn("plugin reload: storage overloaded during persist, retrying", "backoff_ms", d.Milliseconds())
		if timer == nil {
			timer = time.NewTimer(d)
		} else {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(d)
		}
		select {
		case <-ctx.Done():
			logger.Warn("plugin reload: context cancelled during persist retry, SHA metadata may be stale until next reload")
			return
		case <-timer.C:
		}
	}
}

// reloadMatchingPluginMounts reloads provided mounts, regardless of
// plugin name, as long as the backend type is plugin.
func (c *Core) reloadMatchingPluginMounts(ctx context.Context, ns *namespace.Namespace, mounts []string) error {
	if c.pluginCatalog != nil {
		c.pluginCatalog.BeginReloadOperation()
		defer c.pluginCatalog.EndReloadOperation()
	}

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()
	c.authLock.RLock()
	defer c.authLock.RUnlock()

	type mountReloadTarget struct {
		mountPath string
		entry     *MountEntry
		isAuth    bool
	}

	// Collect all matching targets first under the locks, then reload
	// concurrently to reduce total wall time.
	var errors error
	var targets []mountReloadTarget
	for _, mount := range mounts {
		var isAuth bool
		// allow any of
		//   - sys/auth/foo/
		//   - sys/auth/foo
		//   - auth/foo/
		//   - auth/foo
		if strings.HasPrefix(mount, credentialRoutePrefix) {
			isAuth = true
		} else if strings.HasPrefix(mount, mountPathSystem+credentialRoutePrefix) {
			isAuth = true
			mount = strings.TrimPrefix(mount, mountPathSystem)
		}
		if !strings.HasSuffix(mount, "/") {
			mount += "/"
		}

		entry := c.router.MatchingMountEntry(ctx, mount)
		if entry == nil {
			errors = multierror.Append(errors, fmt.Errorf("cannot fetch mount entry on %q", mount))
			continue
		}

		// We dont reload mounts that are not in the same namespace
		if ns.ID != entry.Namespace().ID {
			continue
		}

		targets = append(targets, mountReloadTarget{
			mountPath: mount,
			entry:     entry,
			isAuth:    isAuth,
		})
	}

	if len(targets) == 0 {
		return errors
	}

	type reloadResult struct {
		mountPath string
		entry     *MountEntry
		err       error
	}
	results, skippedFrom := dispatchReloads(ctx, targets,
		func(ctx context.Context, t mountReloadTarget) reloadResult {
			return reloadResult{
				mountPath: t.mountPath,
				entry:     t.entry,
				err:       c.reloadBackendCommon(ctx, t.entry, t.isAuth),
			}
		},
	)
	for i := skippedFrom; i < len(results); i++ {
		if results[i].err == nil {
			results[i].err = ctx.Err()
		}
	}

	for _, r := range results {
		if r.err != nil {
			errors = multierror.Append(errors, fmt.Errorf("cannot reload plugin on %q: %w", r.mountPath, r.err))
			continue
		}
		c.logger.Info("successfully reloaded plugin", "plugin", r.entry.Accessor, "path", r.entry.Path, "version", r.entry.RunningVersion)
	}
	return errors
}

// reloadMatchingPlugin reloads all mounted backends that are named pluginName
// (name of the plugin as registered in the plugin catalog). It returns the
// number of plugins that were reloaded and an error if any.
func (c *Core) reloadMatchingPlugin(ctx context.Context, ns *namespace.Namespace, pluginType consts.PluginType, pluginName string) (reloaded int, err error) {
	if c.pluginCatalog != nil {
		c.pluginCatalog.BeginReloadOperation()
		defer c.pluginCatalog.EndReloadOperation()
	}

	var secrets, auth, database bool
	switch pluginType {
	case consts.PluginTypeSecrets:
		secrets = true
	case consts.PluginTypeCredential:
		auth = true
	case consts.PluginTypeDatabase:
		database = true
	case consts.PluginTypeUnknown:
		secrets = true
		auth = true
		database = true
	default:
		return reloaded, fmt.Errorf("unsupported plugin type %q", pluginType.String())
	}

	type reloadResult struct {
		entry *MountEntry
		err   error
	}

	// reloadErrs accumulates failures across secrets and auth sections so
	// every matching mount is attempted even when some fail.
	var reloadErrs error

	if secrets || database {
		c.mountsLock.RLock()
		defer c.mountsLock.RUnlock()

		// Collect all matching secret mounts first so we can reload them
		// concurrently. Database mounts use a different reload path (internal
		// router request) and retain early-return semantics intentionally —
		// the routed request delegates failure handling to the database plugin.
		var matchingSecretEntries []*MountEntry
		for _, entry := range c.mounts.Entries {
			// We don't reload mounts that are not in the same namespace
			if ns != nil && ns.ID != entry.Namespace().ID {
				continue
			}

			if secrets && (entry.Type == pluginName || (entry.Type == "plugin" && entry.Config.PluginName == pluginName)) {
				matchingSecretEntries = append(matchingSecretEntries, entry)
			} else if database && entry.Type == "database" {
				// The combined database plugin is itself a secrets engine, but
				// knowledge of whether a database plugin is in use within a
				// particular mount is internal to the combined database
				// plugin's storage, so we delegate the reload request with an
				// internally routed request.
				reqCtx := namespace.ContextWithNamespace(ctx, entry.namespace)
				req := &logical.Request{
					Operation: logical.UpdateOperation,
					Path:      entry.Path + "reload/" + pluginName,
				}
				resp, err := c.router.Route(reqCtx, req)
				if err != nil {
					return reloaded, err
				}
				if resp == nil {
					return reloaded, fmt.Errorf("failed to reload %q database plugin(s) mounted under %s", pluginName, entry.Path)
				}
				if resp.IsError() {
					return reloaded, fmt.Errorf("failed to reload %q database plugin(s) mounted under %s: %s", pluginName, entry.Path, resp.Error())
				}

				if count, ok := resp.Data["count"].(int); ok && count > 0 {
					c.logger.Info("successfully reloaded database plugin(s)", "plugin", pluginName, "namespace", entry.Namespace(), "path", entry.Path, "connections", resp.Data["connections"])
					reloaded += count
				}
			}
		}

		if len(matchingSecretEntries) > 0 {
			results, skippedFrom := dispatchReloads(ctx, matchingSecretEntries,
				func(ctx context.Context, entry *MountEntry) reloadResult {
					return reloadResult{
						entry: entry,
						err:   c.reloadBackendCommon(ctx, entry, false),
					}
				},
			)
			for i := skippedFrom; i < len(results); i++ {
				if results[i].err == nil {
					results[i].err = ctx.Err()
				}
			}

			for _, r := range results {
				if r.err != nil {
					reloadErrs = multierror.Append(reloadErrs, fmt.Errorf("cannot reload plugin on %q: %w", r.entry.Path, r.err))
					continue
				}
				reloaded++
				c.logger.Info("successfully reloaded plugin", "plugin", pluginName, "namespace", r.entry.Namespace(), "path", r.entry.Path, "version", r.entry.RunningVersion)
			}
		}
	}

	if auth {
		c.authLock.RLock()
		defer c.authLock.RUnlock()

		// Collect all matching auth mounts first so we can reload them
		// concurrently. Each reloadBackendCommon call holds a per-routeEntry
		// lock (re.l), so concurrent calls across different mounts do not
		// contend with each other. The plugin catalog lock is acquired only
		// for the brief kill and newPluginClient operations, with the reloading
		// sentinel ensuring the shared process is killed exactly once.
		var matchingAuthEntries []*MountEntry
		for _, entry := range c.auth.Entries {
			// We don't reload mounts that are not in the same namespace
			if ns != nil && ns.ID != entry.Namespace().ID {
				continue
			}

			if entry.Type == pluginName || (entry.Type == "plugin" && entry.Config.PluginName == pluginName) {
				matchingAuthEntries = append(matchingAuthEntries, entry)
			}
		}

		if len(matchingAuthEntries) > 0 {
			results, skippedFrom := dispatchReloads(ctx, matchingAuthEntries,
				func(ctx context.Context, entry *MountEntry) reloadResult {
					return reloadResult{
						entry: entry,
						err:   c.reloadBackendCommon(ctx, entry, true),
					}
				},
			)
			for i := skippedFrom; i < len(results); i++ {
				if results[i].err == nil {
					results[i].err = ctx.Err()
				}
			}

			for _, r := range results {
				if r.err != nil {
					reloadErrs = multierror.Append(reloadErrs, fmt.Errorf("cannot reload plugin on %q: %w", r.entry.Path, r.err))
					continue
				}
				reloaded++
				c.logger.Info("successfully reloaded plugin", "plugin", r.entry.Accessor, "path", r.entry.Path, "version", r.entry.RunningVersion)
			}
		}
	}

	return reloaded, reloadErrs
}

// forceReloadBackend reloads any backend, even a singleton. It does not update the cache. Most users
// should use reloadBackendCommon.
//
// The reload is split into two phases:
//
//  1. Process swap (under re.l write lock): drain in-flight requests, clean up
//     the old backend, spawn/attach the new one, and make it live. The lock is
//     released immediately once the new backend is installed so live traffic
//     resumes without waiting for the storage write.
//
//  2. Persist (outside re.l): write the updated RunningSha256 to storage via
//     retryPersistOnOverload. Failure here is non-fatal — the process swap
//     already succeeded, so the in-memory backend remains available even if
//     its persisted SHA metadata is temporarily stale.
func (c *Core) forceReloadBackend(ctx context.Context, entry *MountEntry, isAuth bool) error {
	re, err := c.getRouteEntryForMount(entry, isAuth)
	if err != nil {
		return err
	}

	// Backend doesn't exist
	if re == nil {
		return nil
	}

	// Phase 1: swap the backend under the route lock.
	// defer re.l.Unlock() is scoped to the IIFE so the lock is released as
	// soon as the new backend is installed, before the persist call below.
	var runningSHAChanged bool
	swapErr := func() error {
		// Grab the lock, this allows requests to drain before we cleanup the
		// client.
		re.l.Lock()
		defer re.l.Unlock()

		oldSHA := entry.RunningSha256

		// Only call Cleanup if backend is initialized
		if re.backend != nil {
			// Pass a context value so that the plugin client will call the
			// appropriate cleanup method for reloading
			reloadCtx := context.WithValue(ctx, plugin.ContextKeyPluginReload, "reload")
			// Call backend's Cleanup routine
			re.backend.Cleanup(reloadCtx)
		}

		view := re.storageView
		viewPath := entry.UUID + "/"
		switch entry.Table {
		case mountTableType:
			viewPath = backendBarrierPrefix + viewPath
		case credentialTableType:
			viewPath = credentialBarrierPrefix + viewPath
		}

		removePathCheckers(c, entry, viewPath)

		sysView := c.mountEntrySysView(entry)

		nilMount, err := preprocessMount(c, entry, view.(*BarrierView))
		if err != nil {
			return err
		}

		var backend logical.Backend
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

		addPathCheckers(c, entry, backend, viewPath)

		if nilMount {
			backend.Cleanup(ctx)
			backend = nil
		}

		// Set the backend back
		re.backend = backend

		if backend != nil {
			// Initialize the backend after reload. This is a no-op for backends < v5 which
			// rely on lazy loading for initialization. v5 backends do not rely on lazy loading
			// for initialization unless the plugin process is killed. Reload of a v5 backend
			// results in a new plugin process, so we must initialize the backend here.
			err := backend.Initialize(ctx, &logical.InitializationRequest{
				Storage:             view,
				MountPoint:          entry.Path,
				MountType:           entry.Type,
				MountAccessor:       entry.Accessor,
				BackendUUID:         entry.BackendAwareUUID,
				MountRunningVersion: entry.RunningVersion,
			})
			if err != nil {
				return err
			}

			// Set paths as well
			paths := backend.SpecialPaths()
			if paths != nil {
				rootPathsEntry, err := parseSpecialPaths(paths.Root)
				if err != nil {
					return err
				}
				re.rootPaths.Store(rootPathsEntry)
				loginPathsEntry, err := parseSpecialPaths(paths.Unauthenticated)
				if err != nil {
					return err
				}
				re.loginPaths.Store(loginPathsEntry)
				binaryPathsEntry, err := parseSpecialPaths(paths.Binary)
				if err != nil {
					return err
				}
				re.binaryPaths.Store(binaryPathsEntry)
				allowSnapshotReadPathsEntry, err := parseSpecialPaths(paths.AllowSnapshotRead)
				if err != nil {
					return err
				}
				re.allowSnapshotReadPaths.Store(allowSnapshotReadPathsEntry)
			}
		}

		runningSHAChanged = oldSHA != entry.RunningSha256
		return nil
	}()
	if swapErr != nil {
		return swapErr
	}

	// Phase 2: persist the updated RunningSha256 to storage, outside re.l.
	// re.l is no longer held; live traffic flows through the new backend
	// during any retry backoff.
	if runningSHAChanged && MountTableUpdateStorage {
		retryPersistOnOverload(ctx, c.logger, func() error {
			if isAuth {
				return c.persistAuth(ctx, c.auth, &entry.Local)
			}
			return c.persistMounts(ctx, c.mounts, &entry.Local)
		})
	}

	return nil
}

func (c *Core) getRouteEntryForMount(entry *MountEntry, isAuth bool) (*routeEntry, error) {
	path := entry.Path

	if isAuth {
		path = credentialRoutePrefix + path
	}

	// Fast-path out if the backend doesn't exist
	raw, ok := c.router.root.Get(entry.Namespace().Path + path)
	if !ok {
		return nil, logical.ErrNotFound
	}

	re := raw.(*routeEntry)
	return re, nil
}

// reloadBackendCommon is a generic method to reload a backend provided a
// MountEntry.
func (c *Core) reloadBackendCommon(ctx context.Context, entry *MountEntry, isAuth bool) error {
	// Make sure our cache is up-to-date. Since some singleton mounts can be
	// tuned, we do this before the below check.
	entry.SyncCache()

	// We don't want to reload the singleton mounts. They often have specific
	// inmemory elements and we don't want to touch them here.
	if strutil.StrListContains(singletonMounts, entry.Type) {
		c.logger.Debug("skipping reload of singleton mount", "type", entry.Type)
		return nil
	}

	return c.forceReloadBackend(ctx, entry, isAuth)
}

func (c *Core) setupPluginReload() error {
	return handleSetupPluginReload(c)
}
