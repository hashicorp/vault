package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sync"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/pluginruntimeutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	pluginRuntimeCatalogPath           = "core/plugin-runtime-catalog/"
	ErrPluginRuntimeNotFound           = errors.New("plugin runtime not found")
	ErrPluginRuntimeBadType            = errors.New("unable to determine plugin runtime type")
	ErrPluginRuntimeBadContainerConfig = errors.New("bad container config")
)

// PluginRuntimeCatalog keeps a record of plugin runtimes. Plugin runtimes need
// to be registered to the catalog before they can be used in backends when registering plugins with runtimes
type PluginRuntimeCatalog struct {
	catalogView *BarrierView
	logger      log.Logger

	lock sync.RWMutex
}

func (c *Core) setupPluginRuntimeCatalog(ctx context.Context) error {
	c.pluginRuntimeCatalog = &PluginRuntimeCatalog{
		catalogView: NewBarrierView(c.barrier, pluginRuntimeCatalogPath),
		logger:      c.logger,
	}

	if c.logger.IsInfo() {
		c.logger.Info("successfully setup plugin runtime catalog")
	}

	return nil
}

// Get retrieves a plugin runtime with the specified name from the catalog
// It returns a PluginRuntimeConfig or an error if no plugin runtime was found.
func (c *PluginRuntimeCatalog) Get(ctx context.Context, name string, prt consts.PluginRuntimeType) (*pluginruntimeutil.PluginRuntimeConfig, error) {
	storageKey := path.Join(prt.String(), name)
	c.lock.RLock()
	defer c.lock.RUnlock()
	entry, err := c.catalogView.Get(ctx, storageKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve plugin runtime %q %q: %w", prt.String(), name, err)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to retrieve plugin %q %q: %w", prt.String(), name, err)
	}
	runner := new(pluginruntimeutil.PluginRuntimeConfig)
	if err := jsonutil.DecodeJSON(entry.Value, runner); err != nil {
		return nil, fmt.Errorf("failed to decode plugin runtime entry: %w", err)
	}
	if runner.Type != prt {
		return nil, nil
	}
	return runner, nil
}

// Set registers a new plugin with the catalog, or updates an existing plugin runtime
func (c *PluginRuntimeCatalog) Set(ctx context.Context, conf *pluginruntimeutil.PluginRuntimeConfig) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if conf == nil {
		return fmt.Errorf("plugin runtime config reference is nil")
	}

	buf, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to encode plugin entry: %w", err)
	}

	storageKey := path.Join(conf.Type.String(), conf.Name)
	logicalEntry := logical.StorageEntry{
		Key:   storageKey,
		Value: buf,
	}

	if err := c.catalogView.Put(ctx, &logicalEntry); err != nil {
		return fmt.Errorf("failed to persist plugin runtime entry: %w", err)
	}
	return err
}

// Delete is used to remove an external plugin from the catalog. Builtin plugins
// can not be deleted.
func (c *PluginRuntimeCatalog) Delete(ctx context.Context, name string, prt consts.PluginRuntimeType) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	storageKey := path.Join(prt.String(), name)
	out, err := c.catalogView.Get(ctx, storageKey)
	if err != nil || out == nil {
		return ErrPluginRuntimeNotFound
	}

	return c.catalogView.Delete(ctx, storageKey)
}

func (c *PluginRuntimeCatalog) List(ctx context.Context, prt consts.PluginRuntimeType) ([]*pluginruntimeutil.PluginRuntimeConfig, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var retList []*pluginruntimeutil.PluginRuntimeConfig
	keys, err := logical.CollectKeys(ctx, c.catalogView)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		entry, err := c.catalogView.Get(ctx, key)
		if err != nil || entry == nil {
			continue
		}

		conf := new(pluginruntimeutil.PluginRuntimeConfig)
		if err := jsonutil.DecodeJSON(entry.Value, conf); err != nil {
			return nil, fmt.Errorf("failed to decode plugin runtime entry: %w", err)
		}

		if conf.Type != prt {
			continue
		}

		retList = append(retList, conf)
	}
	return retList, nil
}
