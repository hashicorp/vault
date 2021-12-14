package vault

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/namespace"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// managedKeyRegistrySubPath is the storage prefix used by the registry.
	// Entries are stored under managedKeyRegistrySubPath / <namespace ID> / <key Type>
	managedKeyRegistrySubPath = "key-registry/"

	// managedKeyRegistryCacheSize is the number of policies that are kept cached
	managedKeyRegistryCacheSize = 1024
)

var (
	ErrNilMangedKeyConfiguration = errors.New("nil ManagedKeyConfiguration passed in for storage")
	ErrInvalidManagedKeyName     = errors.New("invalid name for a ManagedKeyConfiguration")
	ErrMangedKeyInUse            = errors.New("cannot delete a managed key whlie it is configured for use by a mount")
)

// ManagedKeyRegistry is used to provide durable storage for ManagedKeyConfigurations.
type ManagedKeyRegistry struct {
	core *Core

	// The view for storage used by the registry for the root namespace.
	// It should be accessed via method getBarrierView().
	// It is protected by modifyLock.
	view *BarrierView

	// The cache for the registry, it is protected by modifyLock.
	cache *lru.TwoQueueCache

	modifyLock *sync.RWMutex

	logger hclog.Logger
}

func NewManagedKeyRegistry(core *Core, baseView *BarrierView, system logical.SystemView, logger hclog.Logger) (*ManagedKeyRegistry, error) {
	r := &ManagedKeyRegistry{
		core:       core,
		view:       baseView.SubView(managedKeyRegistrySubPath),
		modifyLock: new(sync.RWMutex),
		logger:     logger,
	}

	if !system.CachingDisabled() {
		cache, err := lru.New2Q(managedKeyRegistryCacheSize)
		if err != nil {
			return nil, err
		}
		r.cache = cache
	}

	return r, nil
}

// setupManagedKeyRegistry creates the ManagedKeyRegistry.
func (c *Core) setupManagedKeyRegistry() error {
	sysView := &dynamicSystemView{core: c}
	logger := c.baseLogger.Named("managedKeyRegistry")
	c.AddLogger(logger)

	r, err := NewManagedKeyRegistry(c, c.systemBarrierView, sysView, logger)
	if err != nil {
		return err
	}

	c.managedKeyRegistry = r

	return nil
}

func (r *ManagedKeyRegistry) lock() {
	r.modifyLock.Lock()
}

func (r *ManagedKeyRegistry) unlock() {
	r.modifyLock.Unlock()
}

func (r *ManagedKeyRegistry) cacheIndex(ns *namespace.Namespace, sanitizedName string, keyType ManagedKeyType) string {
	return path.Join(ns.ID, string(keyType), sanitizedName)
}

func (r *ManagedKeyRegistry) cacheAdd(ns *namespace.Namespace, keyConfig *ManagedKeyConfiguration) {
	if r.cache == nil {
		return
	}

	r.cache.Add(r.cacheIndex(ns, keyConfig.Name, keyConfig.Type), keyConfig)
}

func (r *ManagedKeyRegistry) cacheRemove(ns *namespace.Namespace, sanitizedName string, keyType ManagedKeyType) {
	if r.cache == nil {
		return
	}

	r.cache.Remove(r.cacheIndex(ns, sanitizedName, keyType))
}

func (r *ManagedKeyRegistry) cacheGet(ns *namespace.Namespace, sanitizedName string, keyType ManagedKeyType) (*ManagedKeyConfiguration, bool) {
	if r.cache == nil {
		return nil, false
	}

	raw, ok := r.cache.Get(r.cacheIndex(ns, sanitizedName, keyType))
	if ! ok {
		return nil, false
	}

	return raw.(*ManagedKeyConfiguration), true
}

func sanitizeManagedKeyName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "", ErrInvalidManagedKeyName
	}

	return name, nil
}

// SetManagedKey is used to create or update a ManagedKeyConfiguration
func (r *ManagedKeyRegistry) SetManagedKey(ctx context.Context, keyConfig *ManagedKeyConfiguration) error {
	if keyConfig == nil {
		return ErrNilMangedKeyConfiguration
	}

	r.lock()
	defer r.unlock()

	view, ns, err := r.getBarrierView(ctx, keyConfig.Type)
	if err != nil {
		return err
	}

	name, err := sanitizeManagedKeyName(keyConfig.Name)
	if err != nil {
		return err
	}
	keyConfig.Name = name

	// Ensure the raw parameters can be converted to key type-specific parameters
	err = keyConfig.initParameters()
	if err != nil {
		return err
	}

	entry, err := logical.StorageEntryJSON(keyConfig.Name, keyConfig)
	if err != nil {
		return fmt.Errorf("failed to create entry: %w", err)
	}

	if err := view.Put(ctx, entry); err != nil {
		return fmt.Errorf("failed to persist managed key: %w", err)
	}
	r.cacheAdd(ns, keyConfig)

	return nil
}

// DeleteManagedKey is used to delete a ManagedKeyConfiguration
func (r *ManagedKeyRegistry) DeleteManagedKey(ctx context.Context, name string, keyType ManagedKeyType) error {
	name, err := sanitizeManagedKeyName(name)
	if err != nil {
		return err
	}

	r.lock()
	defer r.unlock()

	view, ns, err := r.getBarrierView(ctx, keyType)
	if err != nil {
		return err
	}

	if r.isKeyInUse(ns, name) {
		return ErrMangedKeyInUse
	}

	if err := view.Delete(ctx, name); err != nil {
		return fmt.Errorf("failed to delete managed key: %w", err)
	}
	r.cacheRemove(ns, name, keyType)

	return nil
}

// isKeyInUse returns true if there is any mount within the namespace that has been tuned
// to specify the key name as an allowed_managed_keys.
// Warning: the name is assumed to be already sanitized.
func (r *ManagedKeyRegistry) isKeyInUse(ns *namespace.Namespace, name string) bool {
	r.core.mountsLock.RLock()
	defer r.core.mountsLock.RUnlock()

	for _, mount := range r.core.mounts.Entries {
		mountNamespace, ok := mount.synthesizedConfigCache.Load("namespace_id")
		if ok && ns.ID != mountNamespace {
			r.logger.Warn("cannot determine namespace for mount", "path", mount.Path)
			continue
		}

		if rawVal, ok := mount.synthesizedConfigCache.Load("allowed_managed_keys"); ok {
			allowedManagedKeys := rawVal.([]string)
			// sanitize the names
			for i := range allowedManagedKeys {
				sanitizedName, err := sanitizeManagedKeyName(allowedManagedKeys[i])
				if err != nil {
					continue
				}
				allowedManagedKeys[i] = sanitizedName
			}

			if strutil.StrListContains(allowedManagedKeys, name) {
				return true
			}
		}
	}

	return false
}

// GetManagedKey is used to read a ManagedKeyConfiguration from storage. Returns (nil, nil) if
// no key is found.
func (r *ManagedKeyRegistry) GetManagedKey(ctx context.Context, name string, keyType ManagedKeyType) (*ManagedKeyConfiguration, error) {
	name, err := sanitizeManagedKeyName(name)
	if err != nil {
		return nil, err
	}

	view, ns, err := r.getBarrierView(ctx, keyType)
	if err != nil {
		return nil, err
	}

	if cachedConfig, ok := r.cacheGet(ns, name, keyType); ok {
		return cachedConfig, nil
	}

	r.lock()
	defer r.unlock()

	// Check the cache again as the config might have been added while waiting to grab the lock
	if cachedConfig, ok := r.cacheGet(ns, name, keyType); ok {
		return cachedConfig, nil
	}

	entry, err := view.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	keyConfig := new(ManagedKeyConfiguration)
	err = entry.DecodeJSON(keyConfig)
	if err != nil {
		return nil, err
	}

	err = keyConfig.initParameters()
	if err != nil {
		return nil, err
	}

	r.cacheAdd(ns, keyConfig)

	return keyConfig, nil
}

// ListManagedKeys is used to list all the ManagedKeyConfiguration entries in storage
// for the namespace in the context.
func (r *ManagedKeyRegistry) ListManagedKeys(ctx context.Context, keyType ManagedKeyType) ([]string, error) {
	view, _, err := r.getBarrierView(ctx, keyType)
	if err != nil {
		return nil, err
	}

	keys, err := logical.CollectKeys(ctx, view)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// invalidate updates the cache. Argument keyPath is the full storage key for the
// storage entry, composed of managedKeyRegistrySubPath + <keyType sub path> + <key name>
func (r *ManagedKeyRegistry) invalidate(ctx context.Context, keyPath string) {
	rawKeyType, name := path.Split(strings.TrimPrefix(keyPath, managedKeyRegistrySubPath))
	keyType := ManagedKeyType(rawKeyType)

	_, ns, err := r.getBarrierView(ctx, keyType)
	if err != nil {
		r.logger.Error("Unable to invalidate the managed key, cannot determine the namespace", "name", name)
	}

	r.lock()
	defer r.unlock()

	r.cacheRemove(ns, name, keyType)
}
