// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/maps"
)

type TokenCapabilityRefreshBehaviour int

const (
	TokenCapabilityRefreshBehaviourOptimistic TokenCapabilityRefreshBehaviour = iota
	TokenCapabilityRefreshBehaviourPessimistic
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
	workerPool                                 *workerpool.WorkerPool
	staticSecretTokenCapabilityRefreshInterval time.Duration
	tokenCapabilityRefreshBehaviour            TokenCapabilityRefreshBehaviour
}

// StaticSecretCapabilityManagerConfig is the configuration for initializing a new
// StaticSecretCapabilityManager.
type StaticSecretCapabilityManagerConfig struct {
	LeaseCache                                  *LeaseCache
	Logger                                      hclog.Logger
	Client                                      *api.Client
	StaticSecretTokenCapabilityRefreshInterval  time.Duration
	StaticSecretTokenCapabilityRefreshBehaviour string
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

	if conf.Client == nil {
		return nil, fmt.Errorf("nil Client (a required parameter): %v", conf)
	}

	if conf.StaticSecretTokenCapabilityRefreshInterval == 0 {
		conf.StaticSecretTokenCapabilityRefreshInterval = DefaultStaticSecretTokenCapabilityRefreshInterval
	}

	behaviour := TokenCapabilityRefreshBehaviourOptimistic
	if conf.StaticSecretTokenCapabilityRefreshBehaviour != "" {
		switch conf.StaticSecretTokenCapabilityRefreshBehaviour {
		case "optimistic":
			behaviour = TokenCapabilityRefreshBehaviourOptimistic
		case "pessimistic":
			behaviour = TokenCapabilityRefreshBehaviourPessimistic
		default:
			return nil, fmt.Errorf("TokenCapabilityRefreshBehaviour must be either \"optimistic\" or \"pessimistic\"")
		}
	}

	workerPool := workerpool.New(DefaultWorkers)

	return &StaticSecretCapabilityManager{
		client:     conf.Client,
		leaseCache: conf.LeaseCache,
		logger:     conf.Logger,
		workerPool: workerPool,
		staticSecretTokenCapabilityRefreshInterval: conf.StaticSecretTokenCapabilityRefreshInterval,
		tokenCapabilityRefreshBehaviour:            behaviour,
	}, nil
}

// submitWorkToPoolAfterInterval submits work to the pool after the defined
// staticSecretTokenCapabilityRefreshInterval
func (sscm *StaticSecretCapabilityManager) submitWorkToPoolAfterInterval(work func()) {
	time.AfterFunc(sscm.staticSecretTokenCapabilityRefreshInterval, func() {
		if !sscm.workerPool.Stopped() {
			sscm.workerPool.Submit(work)
		}
	})
}

// Stop stops all ongoing jobs and ensures future jobs will not
// get added to the worker pool.
func (sscm *StaticSecretCapabilityManager) Stop() {
	sscm.workerPool.Stop()
}

// StartRenewingCapabilities takes a polling job and submits a constant renewal of capabilities to the worker pool.
// indexToRenew is the capabilities index we'll renew the capabilities for.
func (sscm *StaticSecretCapabilityManager) StartRenewingCapabilities(indexToRenew *cachememdb.CapabilitiesIndex) {
	var work func()
	work = func() {
		if sscm.workerPool.Stopped() {
			sscm.logger.Trace("worker pool stopped, stopping renewal")
			return
		}

		capabilitiesIndex, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, indexToRenew.ID)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			// This cache entry no longer exists, so there is no more work to do.
			sscm.logger.Trace("cache item not found for capabilities refresh, stopping the process")
			return
		}
		if err != nil {
			sscm.logger.Error("error when attempting to get capabilities index to refresh token capabilities", "indexToRenew.ID", indexToRenew.ID, "err", err)
			sscm.submitWorkToPoolAfterInterval(work)
			return
		}

		capabilitiesIndex.IndexLock.RLock()
		token := capabilitiesIndex.Token
		indexReadablePathsMap := capabilitiesIndex.ReadablePaths
		capabilitiesIndex.IndexLock.RUnlock()
		indexReadablePaths := maps.Keys(indexReadablePathsMap)

		client, err := sscm.client.Clone()
		if err != nil {
			sscm.logger.Error("error when attempting clone client to refresh token capabilities", "indexToRenew.ID", indexToRenew.ID, "err", err)
			sscm.submitWorkToPoolAfterInterval(work)
			return
		}

		client.SetToken(token)

		capabilities, err := getCapabilities(indexReadablePaths, client)
		if err != nil {
			sscm.logger.Warn("error when attempting to retrieve updated token capabilities", "indexToRenew.ID", indexToRenew.ID, "err", err)
			if sscm.tokenCapabilityRefreshBehaviour == TokenCapabilityRefreshBehaviourPessimistic {
				// Vault is be sealed or unreachable. If pessimistic, assume we might have
				// lost access. Set capabilities to an empty set, so they are all removed.
				capabilities = make(map[string][]string)
			} else {
				sscm.submitWorkToPoolAfterInterval(work)
				return
			}
		}

		newReadablePaths := reconcileCapabilities(indexReadablePaths, capabilities)
		if maps.Equal(indexReadablePathsMap, newReadablePaths) {
			sscm.logger.Trace("capabilities were the same for index, nothing to do", "indexToRenew.ID", indexToRenew.ID)
			// there's nothing to update!
			sscm.submitWorkToPoolAfterInterval(work)
			return
		}

		// before updating or evicting the index, we must update the tokens on
		// for each path, update the corresponding index with the diff
		for _, path := range indexReadablePaths {
			// If the old path isn't contained in the new readable paths,
			// we must delete it from the tokens map for its corresponding
			// path index.
			if _, ok := newReadablePaths[path]; !ok {
				indexId := hashStaticSecretIndex(path)
				index, err := sscm.leaseCache.db.Get(cachememdb.IndexNameID, indexId)
				if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
					// Nothing to update!
					continue
				}
				if err != nil {
					sscm.logger.Error("error when attempting to update corresponding paths for capabilities index", "indexToRenew.ID", indexToRenew.ID, "err", err)
					sscm.submitWorkToPoolAfterInterval(work)
					return
				}
				sscm.logger.Trace("updating tokens for index, as capability has been lost", "index.ID", index.ID, "request_path", index.RequestPath)
				index.IndexLock.Lock()
				delete(index.Tokens, capabilitiesIndex.Token)
				err = sscm.leaseCache.Set(context.Background(), index)
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
			err := sscm.leaseCache.db.EvictCapabilitiesIndex(cachememdb.IndexNameID, indexToRenew.ID)
			if err != nil {
				sscm.logger.Error("error when attempting to evict capabilities from cache", "index.ID", indexToRenew.ID, "err", err)
				sscm.submitWorkToPoolAfterInterval(work)
				return
			}
			sscm.logger.Debug("successfully evicted capabilities index from cache", "index.ID", indexToRenew.ID)
			// If we successfully evicted the index, no need to re-submit the work to the pool.
			return
		}

		// The token still has some capabilities, so, update the capabilities index:
		capabilitiesIndex.ReadablePaths = newReadablePaths
		err = sscm.leaseCache.SetCapabilitiesIndex(context.Background(), capabilitiesIndex)
		if err != nil {
			sscm.logger.Error("error when attempting to update capabilities from cache", "index.ID", indexToRenew.ID, "err", err)
		}

		// Finally, put ourselves back on the work pool after
		sscm.submitWorkToPoolAfterInterval(work)
		return
	}

	sscm.submitWorkToPoolAfterInterval(work)
}

// getCapabilities is a wrapper around a /sys/capabilities-self call that returns
// capabilities as a map with paths as keys, and capabilities as values.
func getCapabilities(paths []string, client *api.Client) (map[string][]string, error) {
	body := make(map[string]interface{})
	body["paths"] = paths
	capabilities := make(map[string][]string)

	secret, err := client.Logical().Write("sys/capabilities-self", body)
	if err != nil && strings.Contains(err.Error(), "permission denied") {
		// Token has expired. Return an empty set of capabilities:
		return capabilities, nil
	}
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

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
