// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/vault/sdk/helper/cryptoutil"

	"golang.org/x/exp/maps"

	"github.com/mitchellh/mapstructure"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
)

const (
	// DefaultWorkers is the default number of workers for the worker pool.
	DefaultWorkers = 5

	// DefaultStaticSecretTokenCapabilityRefreshInterval is the default time
	// between each capability poll. This is configured with the following config value:
	// static_secret_token_capability_refresh_interval
	DefaultStaticSecretTokenCapabilityRefreshInterval = 5 * time.Minute
)

// StaticSecretCapabilityManager is a struct that utilizes
// a worker pool to keep capabilities up to date.
type StaticSecretCapabilityManager struct {
	client                                     *api.Client
	leaseCache                                 *LeaseCache
	logger                                     hclog.Logger
	tokenSink                                  sink.Sink
	workerPool                                 *workerpool.WorkerPool
	staticSecretTokenCapabilityRefreshInterval time.Duration
}

// StaticSecretCapabilityManagerConfig is the configuration for initializing a new
// StaticSecretCapabilityManager.
type StaticSecretCapabilityManagerConfig struct {
	LeaseCache *LeaseCache
	Logger     hclog.Logger
	// TokenSink is a token sync that will have the latest
	// token from auto-auth in it, to be used in event system
	// connections.
	TokenSink sink.Sink
}

// NewStaticSecretCapabilityManager creates a new instance of a StaticSecretCapabilityManager.
func NewStaticSecretCapabilityManager(conf *StaticSecretCapabilityManagerConfig) (*StaticSecretCapabilityManager, error) {
	if conf == nil {
		return nil, errors.New("nil configuration provided")
	}

	if conf.LeaseCache == nil {
		return nil, fmt.Errorf("nil Lease Cache (a required parameter): %v", conf)
	}

	if conf.Logger == nil {
		return nil, fmt.Errorf("nil Logger (a required parameter): %v", conf)
	}

	if conf.TokenSink == nil {
		return nil, fmt.Errorf("nil token sink (a required parameter): %v", conf)
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	workerPool := workerpool.New(DefaultWorkers)

	return &StaticSecretCapabilityManager{
		client:     client,
		leaseCache: conf.LeaseCache,
		logger:     conf.Logger,
		tokenSink:  conf.TokenSink,
		workerPool: workerPool,
		staticSecretTokenCapabilityRefreshInterval: DefaultStaticSecretTokenCapabilityRefreshInterval,
	}, nil
}

type PollingJob struct {
	ctx    context.Context
	client *api.Client
	// index is the index of the cache entry of the permissions we want
	// to check and keep up to date
	index *cachememdb.Index
}

// SubmitWorkToPoolAfterInterval submits work to the pool after the defined
// staticSecretTokenCapabilityRefreshInterval
func (sscm *StaticSecretCapabilityManager) SubmitWorkToPoolAfterInterval(work func()) {
	time.AfterFunc(sscm.staticSecretTokenCapabilityRefreshInterval, func() {
		sscm.workerPool.Submit(work)
	})
}

// StartRenewing takes a polling job and submits a constant renewal to the worker pool.
func (sscm *StaticSecretCapabilityManager) StartRenewing(pj *PollingJob) error {
	var work func()
	work = func() {
		capabilitiesIndex, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, pj.index.ID)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			// This cache entry no longer exists, so there is no more work to do.
			return
		}
		if err != nil {
			sscm.logger.Error("error when attempting to refresh token capabilities", "index.ID", pj.index.ID, "err", err)
			sscm.SubmitWorkToPoolAfterInterval(work)
			return
		}

		capabilitiesIndex.IndexLock.RLock()
		token := capabilitiesIndex.Token
		indexReadablePathsMap := capabilitiesIndex.ReadablePaths
		capabilitiesIndex.IndexLock.RUnlock()
		indexReadablePaths := maps.Keys(indexReadablePathsMap)

		client, err := sscm.client.Clone()
		if err != nil {
			sscm.logger.Error("error when attempting to refresh token capabilities", "index.ID", pj.index.ID, "err", err)
			sscm.SubmitWorkToPoolAfterInterval(work)
			return
		}

		client.SetToken(token)

		capabilities, err := getCapabilities(indexReadablePaths, client)
		if err != nil {
			sscm.logger.Error("error when attempting to refresh token capabilities", "index.ID", pj.index.ID, "err", err)
			sscm.SubmitWorkToPoolAfterInterval(work)
			return
		}

		newReadablePaths := reconcileCapabilities(indexReadablePaths, capabilities)
		if maps.Equal(indexReadablePathsMap, newReadablePaths) {
			// there's nothing to update!
			sscm.SubmitWorkToPoolAfterInterval(work)
			return
		}

		// before updating or evicting the index, we must update the tokens on
		// for each path, update the corresponding index with the diff
		for _, path := range indexReadablePaths {
			// If the old path isn't contained in the new readable paths,
			// we must delete it from the tokens map for its corresponding
			// path index.
			if _, ok := newReadablePaths[path]; !ok {
				// TODO: replace with hashStaticSecretIndex(path)
				indexId := hex.EncodeToString(cryptoutil.Blake2b256Hash(path))
				index, err := sscm.leaseCache.db.Get(cachememdb.IndexNameID, indexId)
				if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
					// Nothing to update!
					continue
				}
				if err != nil {
					sscm.logger.Error("error when attempting to update corresponding paths for capabilities index", "index.ID", pj.index.ID, "err", err)
					sscm.SubmitWorkToPoolAfterInterval(work)
					return
				}
				index.IndexLock.Lock()
				delete(index.Tokens, path)
				err = sscm.leaseCache.db.Set(index)
				if err != nil {
					sscm.logger.Error("error when attempting to update index in cache", "index.ID", index.ID, "err", err)
				}
				index.IndexLock.Unlock()
			}
		}

		// Lastly, we should update the capabilities index, either evicting or updating it
		capabilitiesIndex.IndexLock.Lock()
		defer capabilitiesIndex.IndexLock.Unlock()
		if len(newReadablePaths) == 0 {
			err := sscm.leaseCache.db.Evict(cachememdb.IndexNameID, pj.index.ID)
			if err != nil {
				sscm.logger.Error("error when attempting to evict capabilities from cache", "index.ID", pj.index.ID, "err", err)
				sscm.SubmitWorkToPoolAfterInterval(work)
				return
			}
			// If we successfully evicted the index, no need to re-submit the work to the pool.
			return
		}

		// The token still has some capabilities, so, update the capabilities index:
		capabilitiesIndex.ReadablePaths = newReadablePaths
		err = sscm.leaseCache.db.SetCapabilitiesIndex(capabilitiesIndex)
		if err != nil {
			sscm.logger.Error("error when attempting to update capabilities from cache", "index.ID", pj.index.ID, "err", err)
		}

		// Finally, put ourselves back on the work pool after
		sscm.SubmitWorkToPoolAfterInterval(work)
		return
	}

	return nil
}

// getCapabilities is a wrapper around a /sys/capabilities-self call that returns
// capabilities as a map with paths as keys, and capabilities as values.
func getCapabilities(paths []string, client *api.Client) (map[string][]string, error) {
	body := make(map[string]interface{})
	body["paths"] = paths
	secret, err := client.Logical().Write("/sys/capabilities-self", body)
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	capabilities := make(map[string][]string)

	for _, path := range paths {
		var res []string
		err = mapstructure.Decode(secret.Data[path], &res)
		if err != nil {
			return nil, err
		}

		capabilities[path] = res
	}

	return capabilities, nil
}

// reconcileCapabilities takes a set of known readable paths, and a set of capabilities (a response from the
// sys/capabilities-self endpoint) and returns a subset of the readablePaths after taking into account any updated
// capabilities as a set, represented by a map of strings to structs.
// It will delete any path in readablePaths if it does not have a "root" or "read" capability listed in the
// capabilities map.
func reconcileCapabilities(readablePaths []string, capabilities map[string][]string) map[string]struct{} {
	newReadablePaths := make(map[string]struct{})
	for pathName, permissions := range capabilities {
		if slices.Contains(permissions, "read") || slices.Contains(permissions, "root") {
			// We do this as an additional sanity check. We never want to
			// add permissions that weren't there before.
			if slices.Contains(readablePaths, pathName) {
				newReadablePaths[pathName] = struct{}{}
			}
		}
	}

	return newReadablePaths
}
