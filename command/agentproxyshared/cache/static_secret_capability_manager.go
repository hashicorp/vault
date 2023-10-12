// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	client     *api.Client
	leaseCache *LeaseCache
	logger     hclog.Logger
	tokenSink  sink.Sink
	workerPool *workerpool.WorkerPool
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
	}, nil
}

type PollingJob struct {
	ctx    context.Context
	client *api.Client
	// index is the index of the cache entry of the permissions we want
	// to check and keep up to date
	index *cachememdb.Index
}

// StartRenewing takes a polling job and submits a constant renewal to the worker pool.
func (sscm *StaticSecretCapabilityManager) StartRenewing(pj *PollingJob) error {
	work := func() {
		index, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, pj.index.ID)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			// This cache entry no longer exists, so there is no more work to do.
			return
		}
		if err != nil {
			sscm.logger.Error("error when attempting to refresh token capabilities", "index.ID", pj.index.ID)
			// TODO what do we do here?
			return
		}

		token := index.Token
		client, err := sscm.client.Clone()
		if err != nil {
			sscm.logger.Error("error when attempting to refresh token capabilities", "index.ID", pj.index.ID)
			// TODO what do we do here?
			return
		}

		client.SetToken(token)

		capabilities, err := getCapabilities(maps.Keys(index.ReadablePaths), client)
		if err != nil {
			sscm.logger.Error("error when attempting to refresh token capabilities", "index.ID", pj.index.ID)
			// TODO what do we do here?
			return
		}
		for _, capabilities := range capabilities {
			if len(capabilities) > 0 {
				// TODO
			}
		}
	}

	time.AfterFunc(DefaultStaticSecretTokenCapabilityRefreshInterval, func() {
		sscm.workerPool.Submit(work)
	})

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
