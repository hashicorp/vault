package plugincatalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var ErrPluginRepositoryNotFound = errors.New("plugin repository not found")

type PluginRepositoryConfig struct {
	Name string `json:"name" structs:"name"`
	URL  string `json:"url" structs:"url"`
	// TODO how to take into account different types and authentication config
}

type PluginRepositoryCatalog struct {
	catalogView logical.Storage
	logger      log.Logger

	lock sync.RWMutex
}

func SetupPluginRepositoryCatalog(_ context.Context, logger log.Logger, catalogView logical.Storage) (*PluginRepositoryCatalog, error) {
	pluginRepositoryCatalog := &PluginRepositoryCatalog{
		catalogView: catalogView,
		logger:      logger,
	}

	logger.Info("successfully setup plugin repository catalog")

	return pluginRepositoryCatalog, nil
}

func (c *PluginRepositoryCatalog) Get(ctx context.Context, name string) (*PluginRepositoryConfig, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	entry, err := c.catalogView.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve plugin repository %q: %w", name, err)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to retrieve plugin repository %q: %w", name, err)
	}
	repo := new(PluginRepositoryConfig)
	if err := jsonutil.DecodeJSON(entry.Value, repo); err != nil {
		return nil, fmt.Errorf("failed to decode plugin repository entry: %w", err)
	}
	return repo, nil
}

func (c *PluginRepositoryCatalog) Set(ctx context.Context, conf *PluginRepositoryConfig) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if conf == nil {
		return fmt.Errorf("plugin runtime config reference is nil")
	}

	buf, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to encode plugin entry: %w", err)
	}

	logicalEntry := logical.StorageEntry{
		Key:   conf.Name,
		Value: buf,
	}

	if err := c.catalogView.Put(ctx, &logicalEntry); err != nil {
		return fmt.Errorf("failed to persist plugin repository entry: %w", err)
	}
	return err
}

func (c *PluginRepositoryCatalog) Delete(ctx context.Context, name string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	out, err := c.catalogView.Get(ctx, name)
	if err != nil || out == nil {
		return ErrPluginRuntimeNotFound
	}

	return c.catalogView.Delete(ctx, name)
}

func (c *PluginRepositoryCatalog) List(ctx context.Context) ([]*PluginRepositoryConfig, error) {
	c.logger.Info("Listing plugin repositories")
	c.lock.RLock()
	defer c.lock.RUnlock()

	var retList []*PluginRepositoryConfig
	keys, err := logical.CollectKeys(ctx, c.catalogView)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		entry, err := c.catalogView.Get(ctx, key)
		if err != nil || entry == nil {
			continue
		}

		conf := new(PluginRepositoryConfig)
		if err := jsonutil.DecodeJSON(entry.Value, conf); err != nil {
			return nil, fmt.Errorf("failed to decode plugin repository entry: %w", err)
		}

		retList = append(retList, conf)
	}
	return retList, nil
}
