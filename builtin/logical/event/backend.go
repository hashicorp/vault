// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"fmt"
	"net/rpc"
	"strings"
	"sync"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/syncmap"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/event/evplugin"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	operationPrefixEvent        = "event"
	operationSuffixSubscription = "subscription"
	eventSubscriptionPath       = "subscription/"
)

type subscriptionInstance struct {
	sync.RWMutex
	eventPlugin evplugin.EventPlugin

	id     string
	name   string
	closed bool
}

func (s *subscriptionInstance) ID() string {
	return s.id
}

func (s *subscriptionInstance) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.closed {
		return nil
	}
	s.closed = true

	return s.eventPlugin.Close()
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf).(*eventBackend)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	// collect metrics on number of plugin instances
	var err error
	b.gaugeCollectionProcess, err = metricsutil.NewGaugeCollectionProcess(
		[]string{"event", "pluginInstances", "count"},
		[]metricsutil.Label{},
		b.collectPluginInstanceGaugeValues,
		metrics.Default(),
		configutil.UsageGaugeDefaultPeriod, // TODO: add config settings for these, or add plumbing to the main config settings
		configutil.MaximumGaugeCardinalityDefault,
		b.logger)
	if err != nil {
		return nil, err
	}
	go b.gaugeCollectionProcess.Run()
	return b, nil
}

func Backend(conf *logical.BackendConfig) logical.Backend {
	var b eventBackend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				eventSubscriptionPath + "*",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathListSubscription(&b),
				pathConfigureSubscription(&b),
				pathResetSubscription(&b),
			},
		),

		Clean:       b.clean,
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	b.logger = conf.Logger
	b.subscriptions = syncmap.NewSyncMap[string, *subscriptionInstance]()
	b.queueCtx, b.cancelQueueCtx = context.WithCancel(context.Background())
	return &b
}

func (b *eventBackend) collectPluginInstanceGaugeValues(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	// copy the map so we can release the lock
	subscriptionsCopy := b.subscriptions.Values()
	counts := map[string]int{}
	for _, v := range subscriptionsCopy {
		dbType, err := v.eventPlugin.Type(ctx)
		if err != nil {
			// there's a chance this will already be closed since we don't hold the lock
			continue
		}
		if _, ok := counts[dbType]; !ok {
			counts[dbType] = 0
		}
		counts[dbType] += 1
	}
	var gauges []metricsutil.GaugeLabelValues
	for k, v := range counts {
		gauges = append(gauges, metricsutil.GaugeLabelValues{Labels: []metricsutil.Label{{Name: "eventPluginType", Value: k}}, Value: float32(v)})
	}
	return gauges, nil
}

type eventBackend struct {
	// subscriptions holds configured subscriptions by config name
	subscriptions *syncmap.SyncMap[string, *subscriptionInstance]
	logger        log.Logger

	*framework.Backend
	// queueCtx is the context for the priority queue
	queueCtx context.Context
	// cancelQueueCtx is used to terminate the background ticker
	cancelQueueCtx context.CancelFunc

	// the running gauge collection process
	gaugeCollectionProcess     *metricsutil.GaugeCollectionProcess
	gaugeCollectionProcessStop sync.Once
}

func (b *eventBackend) SubscriptionConfig(ctx context.Context, s logical.Storage, name string) (*SubscriptionConfig, error) {
	entry, err := s.Get(ctx, fmt.Sprintf("config/%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to read subscription configuration: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to find entry for subscription with name: %q", name)
	}

	var config SubscriptionConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (b *eventBackend) invalidate(_ context.Context, key string) {
	switch {
	case strings.HasPrefix(key, eventSubscriptionPath):
		name := strings.TrimPrefix(key, eventSubscriptionPath)
		b.ClearSubscription(name)
	}
}

func (b *eventBackend) GetSubscription(ctx context.Context, s logical.Storage, name string) (*subscriptionInstance, error) {
	config, err := b.SubscriptionConfig(ctx, s, name)
	if err != nil {
		return nil, err
	}

	return b.GetSubscriptionWithConfig(ctx, name, config)
}

func (b *eventBackend) GetSubscriptionWithConfig(ctx context.Context, name string, config *SubscriptionConfig) (*subscriptionInstance, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// ClearSubscription closes the subscription connection and
// removes it from the b.subscriptions map.
func (b *eventBackend) ClearSubscription(name string) error {
	db := b.subscriptions.Pop(name)
	if db != nil {
		// Ignore error here since the subscription client is always killed
		db.Close()
	}
	return nil
}

// ClearSubscriptionId closes the subscription connection with a specific id and
// removes it from the b.subscriptions map.
func (b *eventBackend) ClearSubscriptionId(name, id string) error {
	db := b.subscriptions.PopIfEqual(name, id)
	if db != nil {
		// Ignore error here since the subscription client is always killed
		db.Close()
	}
	return nil
}

func (b *eventBackend) CloseIfShutdown(db *subscriptionInstance, err error) {
	// Plugin has shutdown, close it so next call can reconnect.
	switch err {
	case rpc.ErrShutdown:
		// Put this in a goroutine so that requests can run with the read or write lock
		// and simply defer the unlock.  Since we are attaching the instance and matching
		// the id in the subscriptions map, we can safely do this.
		go func() {
			db.Close()

			// Delete the subscription if it is still active.
			b.subscriptions.PopIfEqual(db.name, db.id)
		}()
	}
}

// clean closes all subscriptions
func (b *eventBackend) clean(_ context.Context) {
	// kill the queue and terminate the background ticker
	if b.cancelQueueCtx != nil {
		b.cancelQueueCtx()
	}

	subscriptions := b.subscriptions.Clear()
	for _, db := range subscriptions {
		go db.Close()
	}
	b.gaugeCollectionProcessStop.Do(func() {
		if b.gaugeCollectionProcess != nil {
			b.gaugeCollectionProcess.Stop()
		}
		b.gaugeCollectionProcess = nil
	})
}

const backendHelp = `
The event backend supports routing events to many different
subscriptions.

After mounting this backend, configure it using the endpoints within
the "event/subscription/" path.
`
