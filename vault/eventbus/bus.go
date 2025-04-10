// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/eventlogger/formatter_filters/cloudevents"
	"github.com/hashicorp/go-bexpr"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// eventTypeAll is purely internal to the event bus. We use it to send all
	// events down one big firehose, and pipelines define their own filtering
	// based on what each subscriber is interested in.
	eventTypeAll   = "*"
	defaultTimeout = 60 * time.Second
)

var (
	ErrNotStarted = errors.New("event broker has not been started")
	subscriptions atomic.Int64 // keeps track of event subscription count in all event buses

	// these metadata fields will have the plugin mount path prepended to them
	metadataPrependPathFields = []string{
		logical.EventMetadataPath,
		logical.EventMetadataDataPath,
	}
)

// EventBus contains the main logic of running an event broker for Vault.
// Start() must be called before the EventBus will accept events for sending.
type EventBus struct {
	logger                     hclog.Logger
	broker                     *eventlogger.Broker
	started                    atomic.Bool
	formatterNodeID            eventlogger.NodeID
	timeout                    time.Duration
	filters                    *Filters
	cloudEventsFormatterFilter *cloudevents.FormatterFilter
}

type pluginEventBus struct {
	bus        *EventBus
	namespace  *namespace.Namespace
	pluginInfo *logical.EventPluginInfo
}

type asyncChanNode struct {
	// TODO: add bounded deque buffer of *EventReceived
	ctx    context.Context
	ch     chan *eventlogger.Event
	logger hclog.Logger

	// used to close the connection
	closeOnce      sync.Once
	cancelFunc     context.CancelFunc
	pipelineID     eventlogger.PipelineID
	removeFilter   func()
	removePipeline func(ctx context.Context, t eventlogger.EventType, id eventlogger.PipelineID) (bool, error)
}

var (
	_ eventlogger.Node    = (*asyncChanNode)(nil)
	_ logical.EventSender = (*pluginEventBus)(nil)
)

// Start starts the event bus, allowing events to be written.
// It is not possible to stop or restart the event bus.
// It is safe to call Start() multiple times.
func (bus *EventBus) Start() {
	wasStarted := bus.started.Swap(true)
	if !wasStarted {
		bus.logger.Info("Starting event system")
	}
}

// patchMountPath patches the event data's metadata "secret_path" field, if present, to include the mount path prepended.
func patchMountPath(data *logical.EventData, pluginInfo *logical.EventPluginInfo) *logical.EventData {
	if pluginInfo == nil || pluginInfo.MountPath == "" || data.Metadata == nil {
		return data
	}

	for _, field := range metadataPrependPathFields {
		if data.Metadata.Fields[field] != nil {
			newPath := path.Join(pluginInfo.MountPath, data.Metadata.Fields[field].GetStringValue())
			if pluginInfo.MountClass == "auth" {
				newPath = path.Join("auth", newPath)
			}
			data.Metadata.Fields[field] = structpb.NewStringValue(newPath)
		}
	}

	return data
}

// SendEventInternal sends an event to the event bus and routes it to all relevant subscribers.
// This function does *not* wait for all subscribers to acknowledge before returning.
// This function is meant to be used by trusted internal code, so it can specify details like the namespace
// and plugin info. Events from plugins should be routed through WithPlugin(), which will populate
// the namespace and plugin info automatically.
// The context passed in is currently ignored to ensure that the event is sent if the context is short-lived,
// such as with an HTTP request context.
func (bus *EventBus) SendEventInternal(_ context.Context, ns *namespace.Namespace, pluginInfo *logical.EventPluginInfo, eventType logical.EventType, forwarded bool, data *logical.EventData) error {
	if ns == nil {
		return namespace.ErrNoNamespace
	}
	if !bus.started.Load() {
		return ErrNotStarted
	}
	eventReceived := &logical.EventReceived{
		Namespace:  ns.Path,
		EventType:  string(eventType),
		PluginInfo: pluginInfo,
	}
	// If the event has been forwarded downstream, no need to patch the mount
	// path again
	if forwarded {
		eventReceived.Event = data
	} else {
		eventReceived.Event = patchMountPath(data, pluginInfo)
	}

	// We can't easily know when the SendEvent is complete, so we can't call the cancel function.
	// But, it is called automatically after bus.timeout, so there won't be any leak as long as bus.timeout is not too long.
	ctx, _ := context.WithTimeout(context.Background(), bus.timeout)
	_, err := bus.broker.Send(ctx, eventTypeAll, eventReceived)
	if err != nil {
		// if no listeners for this event type are registered, that's okay, the event
		// will just not be sent anywhere
		if strings.Contains(strings.ToLower(err.Error()), "no graph for eventtype") {
			return nil
		}
	}
	return err
}

func (bus *EventBus) WithPlugin(ns *namespace.Namespace, eventPluginInfo *logical.EventPluginInfo) (*pluginEventBus, error) {
	if ns == nil {
		return nil, namespace.ErrNoNamespace
	}
	return &pluginEventBus{
		bus:        bus,
		namespace:  ns,
		pluginInfo: eventPluginInfo,
	}, nil
}

// SendEvent sends an event to the event bus and routes it to all relevant subscribers.
// This function does *not* wait for all subscribers to acknowledge before returning.
// The context passed in is currently ignored.
func (bus *pluginEventBus) SendEvent(ctx context.Context, eventType logical.EventType, data *logical.EventData) error {
	return bus.bus.SendEventInternal(ctx, bus.namespace, bus.pluginInfo, eventType, false, data)
}

func NewEventBus(localNodeID string, logger hclog.Logger) (*EventBus, error) {
	broker, err := eventlogger.NewBroker()
	if err != nil {
		return nil, err
	}

	formatterID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	formatterNodeID := eventlogger.NodeID(formatterID)

	if logger == nil {
		logger = hclog.Default().Named("events")
	}

	sourceUrl, err := url.Parse("vault://" + localNodeID)
	if err != nil {
		return nil, err
	}

	cloudEventsFormatterFilter := &cloudevents.FormatterFilter{
		Source: sourceUrl,
		Predicate: func(_ context.Context, e interface{}) (bool, error) {
			return true, nil
		},
	}

	return &EventBus{
		logger:                     logger,
		broker:                     broker,
		formatterNodeID:            formatterNodeID,
		timeout:                    defaultTimeout,
		cloudEventsFormatterFilter: cloudEventsFormatterFilter,
		filters:                    NewFilters(localNodeID),
	}, nil
}

// Subscribe subscribes to events in the given namespace matching the event type pattern and after
// applying the optional go-bexpr filter.
func (bus *EventBus) Subscribe(ctx context.Context, ns *namespace.Namespace, pattern string, bexprFilter string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.SubscribeMultipleNamespaces(ctx, []string{strings.Trim(ns.Path, "/")}, pattern, bexprFilter)
}

// SubscribeMultipleNamespaces subscribes to events in the given namespace matching the event type
// pattern and after applying the optional go-bexpr filter.
func (bus *EventBus) SubscribeMultipleNamespaces(ctx context.Context, namespacePathPatterns []string, pattern string, bexprFilter string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.subscribeInternal(ctx, namespacePathPatterns, pattern, bexprFilter, nil)
}

// subscribeInternal creates the pipeline and connects it to the event bus to receive events. If the
// clusterNode is specified, then the namespacePathPatterns, pattern, and bexprFilter are ignored,
// and instead this subscription will be tied to the given cluster node's filter.
func (bus *EventBus) subscribeInternal(ctx context.Context, namespacePathPatterns []string, pattern string, bexprFilter string, clusterNode *string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	// subscriptions are still stored even if the bus has not been started
	pipelineID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	err = bus.broker.RegisterNode(bus.formatterNodeID, bus.cloudEventsFormatterFilter)
	if err != nil {
		return nil, nil, err
	}

	filterNodeID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	var filterNode *eventlogger.Filter
	if clusterNode != nil {
		filterNode, err = newClusterNodeFilterNode(bus.filters, clusterNodeID(*clusterNode))
		if err != nil {
			return nil, nil, err
		}
	} else {
		filterNode, err = newFilterNode(namespacePathPatterns, pattern, bexprFilter)
		if err != nil {
			return nil, nil, err
		}
		bus.filters.addPattern(bus.filters.self, namespacePathPatterns, pattern)
	}
	err = bus.broker.RegisterNode(eventlogger.NodeID(filterNodeID), filterNode)
	if err != nil {
		return nil, nil, err
	}

	sinkNodeID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	asyncNode := newAsyncNode(ctx, bus.logger, bus.broker, func() {
		if clusterNode == nil {
			bus.filters.removePattern(bus.filters.self, namespacePathPatterns, pattern)
		}
	})
	err = bus.broker.RegisterNode(eventlogger.NodeID(sinkNodeID), asyncNode)
	if err != nil {
		defer cancel()
		return nil, nil, err
	}

	nodes := []eventlogger.NodeID{eventlogger.NodeID(filterNodeID), bus.formatterNodeID, eventlogger.NodeID(sinkNodeID)}

	pipeline := eventlogger.Pipeline{
		PipelineID: eventlogger.PipelineID(pipelineID),
		EventType:  eventTypeAll,
		NodeIDs:    nodes,
	}
	err = bus.broker.RegisterPipeline(pipeline)
	if err != nil {
		defer cancel()
		return nil, nil, err
	}

	addSubscriptions(1)
	// add info needed to cancel the subscription
	asyncNode.pipelineID = eventlogger.PipelineID(pipelineID)
	asyncNode.cancelFunc = cancel
	// Capture context in a closure for the cancel func
	return asyncNode.ch, func() { asyncNode.Close(ctx) }, nil
}

// SetSendTimeout sets the timeout of sending events. If the events are not accepted by the
// underlying channel before this timeout, then the channel closed.
func (bus *EventBus) SetSendTimeout(timeout time.Duration) {
	bus.timeout = timeout
}

// GlobalMatch returns true if the given namespace and event type match the current global filter.
func (bus *EventBus) GlobalMatch(ns *namespace.Namespace, eventType logical.EventType) bool {
	return bus.filters.globalMatch(ns, eventType)
}

// ApplyClusterNodeFilterChanges applies the given filter changes to the cluster node's filters.
func (bus *EventBus) ApplyClusterNodeFilterChanges(c string, changes []FilterChange) {
	bus.filters.applyChanges(clusterNodeID(c), changes)
}

// ApplyGlobalFilterChanges applies the given filter changes to the global filters.
func (bus *EventBus) ApplyGlobalFilterChanges(changes []FilterChange) {
	bus.filters.applyChanges(globalCluster, changes)
}

// ClearGlobalFilter removes all entries from the current global filter.
func (bus *EventBus) ClearGlobalFilter() {
	bus.filters.clearGlobalPatterns()
}

// ClearClusterNodeFilter removes all entries from the given cluster node's filter.
func (bus *EventBus) ClearClusterNodeFilter(id string) {
	bus.filters.clearClusterNodePatterns(clusterNodeID(id))
}

// NotifyOnGlobalFilterChanges returns a channel that receives changes to the global filter.
func (bus *EventBus) NotifyOnGlobalFilterChanges(ctx context.Context) (<-chan []FilterChange, context.CancelFunc, error) {
	return bus.filters.watch(ctx, globalCluster)
}

// NotifyOnLocalFilterChanges returns a channel that receives changes to the filter for the current cluster node.
func (bus *EventBus) NotifyOnLocalFilterChanges(ctx context.Context) (<-chan []FilterChange, context.CancelFunc, error) {
	return bus.NotifyOnClusterNodeFilterChanges(ctx, string(bus.filters.self))
}

// NotifyOnClusterNodeFilterChanges returns a channel that receives changes to the filter for the given cluster node.
func (bus *EventBus) NotifyOnClusterNodeFilterChanges(ctx context.Context, clusterNode string) (<-chan []FilterChange, context.CancelFunc, error) {
	return bus.filters.watch(ctx, clusterNodeID(clusterNode))
}

// NewAllEventsSubscription creates a new subscription to all events.
func (bus *EventBus) NewAllEventsSubscription(ctx context.Context) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.subscribeInternal(ctx, nil, "*", "", nil)
}

// NewGlobalSubscription creates a new subscription to all events that match the global filter.
func (bus *EventBus) NewGlobalSubscription(ctx context.Context) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	g := globalCluster
	return bus.subscribeInternal(ctx, nil, "", "", &g)
}

// NewClusterNodeSubscription creates a new subscription to all events that match the given cluster node's filter.
func (bus *EventBus) NewClusterNodeSubscription(ctx context.Context, clusterNode string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.subscribeInternal(ctx, nil, "", "", &clusterNode)
}

// creates a new filter node that is tied to the filter for a given cluster node
func newClusterNodeFilterNode(filters *Filters, c clusterNodeID) (*eventlogger.Filter, error) {
	return &eventlogger.Filter{
		Predicate: func(e *eventlogger.Event) (bool, error) {
			eventRecv := e.Payload.(*logical.EventReceived)
			eventNs := strings.Trim(eventRecv.Namespace, "/")
			if filters.clusterNodeMatch(c, &namespace.Namespace{
				Path: eventNs,
			}, logical.EventType(eventRecv.EventType)) {
				return true, nil
			}
			return false, nil
		},
	}, nil
}

func newFilterNode(namespacePatterns []string, pattern string, bexprFilter string) (*eventlogger.Filter, error) {
	var evaluator *bexpr.Evaluator
	if bexprFilter != "" {
		var err error
		evaluator, err = bexpr.CreateEvaluator(bexprFilter)
		if err != nil {
			return nil, err
		}
	}
	return &eventlogger.Filter{
		Predicate: func(e *eventlogger.Event) (bool, error) {
			eventRecv := e.Payload.(*logical.EventReceived)
			eventNs := strings.Trim(eventRecv.Namespace, "/")
			// Drop if event is not in namespace patterns namespace.
			if len(namespacePatterns) > 0 {
				allow := false
				for _, nsPattern := range namespacePatterns {
					if glob.Glob(nsPattern, eventNs) {
						allow = true
						break
					}
				}
				if !allow {
					return false, nil
				}
			}

			// ClusterFilter for correct event type, including wildcards.
			if !glob.Glob(pattern, eventRecv.EventType) {
				return false, nil
			}

			// apply go-bexpr filter
			if evaluator != nil {
				return evaluator.Evaluate(eventRecv.BexprDatum())
			}
			return true, nil
		},
	}, nil
}

func newAsyncNode(ctx context.Context, logger hclog.Logger, broker *eventlogger.Broker, removeFilter func()) *asyncChanNode {
	return &asyncChanNode{
		ctx:            ctx,
		ch:             make(chan *eventlogger.Event),
		logger:         logger,
		removeFilter:   removeFilter,
		removePipeline: broker.RemovePipelineAndNodes,
	}
}

// Close tells the bus to stop sending us events.
func (node *asyncChanNode) Close(ctx context.Context) {
	node.closeOnce.Do(func() {
		defer node.cancelFunc()
		node.removeFilter()
		removed, err := node.removePipeline(ctx, eventTypeAll, node.pipelineID)

		switch {
		case err != nil && removed:
			msg := fmt.Sprintf("Error removing nodes referenced by pipeline %q", node.pipelineID)
			node.logger.Warn(msg, err)
		case err != nil:
			msg := fmt.Sprintf("Error removing pipeline %q", node.pipelineID)
			node.logger.Warn(msg, err)
		}
		addSubscriptions(-1)
	})
}

func (node *asyncChanNode) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	// sends to the channel async in another goroutine
	go func() {
		var timeout bool
		select {
		case node.ch <- e:
		case <-ctx.Done():
			timeout = errors.Is(ctx.Err(), context.DeadlineExceeded)
		case <-node.ctx.Done():
			timeout = errors.Is(node.ctx.Err(), context.DeadlineExceeded)
		}
		if timeout {
			node.logger.Info("Subscriber took too long to process event, closing", "ID", e.Payload.(*logical.EventReceived).Event.Id)
			node.Close(ctx)
		}
	}()
	return e, nil
}

func (node *asyncChanNode) Reopen() error {
	return nil
}

func (node *asyncChanNode) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}

func addSubscriptions(delta int64) {
	metrics.SetGauge([]string{"events", "subscriptions"}, float32(subscriptions.Add(delta)))
}
