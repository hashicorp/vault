// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/eventlogger/formatter_filters/cloudevents"
	"github.com/hashicorp/go-bexpr"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
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
	eventTypeAll            = "*"
	defaultTimeout          = 60 * time.Second
	defaultSubscriberBuffer = 16
	maxSubscriberBuffer     = 1000
	eventMetadataVaultIndex = "vault_index"

	// EnvVaultEventNotificationsBoundedQueueSize is the environment variable to configure bounded event queues.
	// Set to a positive integer to enable buffered subscriber channels of that size.
	// Slow consumers will be dropped when their queue is full instead of blocking event processing.
	// Set to 0 or leave unset for unbuffered channels (default, backward compatible).
	// Values above 1000 will be capped at 1000 to prevent excessive memory usage.
	//
	// Recommended starting value: 16
	// This provides enough buffer to handle reasonable event bursts while keeping memory usage low
	// (~1.6-3KB per subscriber). Increase this value (e.g., 64-128) if you have high event rates
	// and subscribers that need more time to process events. Decrease it (e.g., 8) if you want
	// faster detection of slow subscribers or have memory constraints.
	EnvVaultEventNotificationsBoundedQueueSize = "VAULT_EVENT_NOTIFICATIONS_BOUNDED_QUEUE_SIZE"
)

var (
	ErrNotStarted   = errors.New("event broker has not been started")
	ErrSlowConsumer = errors.New("event subscriber is too slow")
	subscriptions   atomic.Int64 // keeps track of event subscription count in all event buses

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
	storageInfoGetter          StorageInfoGetter
	subscriberBufferSize       int // cached buffer size from VAULT_BOUNDED_EVENT_QUEUE env var (0 = unbuffered)
}

// StorageInfoGetter is an interface used to access some storage-related core
// functions without importing the core package
type StorageInfoGetter interface {
	GetCurrentWALHeader() string
	IsReplicated(secondaryID, namespaceName, mountPathRelative string) bool
}

type pluginEventBus struct {
	bus        *EventBus
	namespace  *namespace.Namespace
	pluginInfo *logical.EventPluginInfo
}

type asyncChanNode struct {
	// TODO: add bounded deque buffer of *EventReceived
	ctx        context.Context
	ch         chan *eventlogger.Event
	logger     hclog.Logger
	bufferSize int // cached buffer size (0 = unbuffered, >0 = buffered)

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

// getIndexForEvent returns the storage index (wal header) for events with
// metadata.modified=true.
func (bus *EventBus) getIndexForEvent(event *logical.EventReceived) (string, error) {
	if event.Event == nil || event.Event.Metadata == nil || bus.storageInfoGetter == nil {
		return "", nil
	}
	eventMetadataModified := event.Event.Metadata.GetFields()[logical.EventMetadataModified]
	if eventMetadataModified != nil {
		isModified, err := parseutil.ParseBool(eventMetadataModified.GetStringValue())
		if err != nil {
			return "", fmt.Errorf("failed to parse event metadata modified: %w", err)
		}
		if isModified {
			return bus.storageInfoGetter.GetCurrentWALHeader(), nil
		}
	}
	return "", nil
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
		walStr, err := bus.getIndexForEvent(eventReceived)
		if err != nil {
			bus.logger.Warn("Failed to get index for event", "error", err)
		}
		if walStr != "" {
			eventReceived.Event.Metadata.Fields[eventMetadataVaultIndex] = structpb.NewStringValue(walStr)
		}
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

func NewEventBus(localNodeID string, logger hclog.Logger, c StorageInfoGetter) (*EventBus, error) {
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
		storageInfoGetter:          c,
		subscriberBufferSize:       getSubscriberBufferSize(),
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
// clusterNode is specified, then the namespacePathPatterns and pattern are ignored,
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
		filterNode, err = newClusterNodeFilterNode(bus.filters, clusterNodeID(*clusterNode), bexprFilter, bus.storageInfoGetter)
		if err != nil {
			return nil, nil, err
		}
	} else {
		filterNode, err = newFilterNode(namespacePathPatterns, pattern, bexprFilter)
		if err != nil {
			return nil, nil, err
		}
		// use filterNodeID as the "subscription id" when storing a subscriber
		// pattern
		bus.filters.addPattern(bus.filters.self, namespacePathPatterns, pattern, filterNodeID)
		bus.filters.addClusterWidePattern(namespacePathPatterns, pattern)
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
	asyncNode := newAsyncNode(ctx, bus.logger, bus.broker, bus.subscriberBufferSize, func() {
		if clusterNode == nil {
			// use filterNodeID as the "subscription id" when removing a
			// subscriber pattern
			bus.filters.removePattern(bus.filters.self, namespacePathPatterns, pattern, filterNodeID)
			bus.filters.makeClusterWideFilters()
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

// ClusterWideMatch returns true if the given namespace and event type match the
// current cluster-wide filter.
func (bus *EventBus) ClusterWideMatch(ns *namespace.Namespace, eventType logical.EventType) bool {
	return bus.filters.clusterWideMatch(ns, eventType)
}

// ApplyClusterNodeFilterChanges applies the given filter changes to the cluster node's filters.
func (bus *EventBus) ApplyClusterNodeFilterChanges(c string, changes []FilterChange) {
	bus.filters.applyChanges(clusterNodeID(c), changes)
}

// ApplyClusterWideFilterChanges applies the given filter changes to the
// cluster-wide filters.
func (bus *EventBus) ApplyClusterWideFilterChanges(changes []FilterChange) {
	bus.filters.applyChanges(clusterWide, changes)
}

// MakeClusterWideFilters populates the cluster-wide filter based on the current
// filters for all cluster nodes.
func (bus *EventBus) MakeClusterWideFilters() {
	bus.filters.makeClusterWideFilters()
}

// GetClusterWideFilterAdditions returns the current cluster-wide filter
// represented purely as additive changes. It is used by a secondary cluster's
// active node to resync its full set of event subscription filters to the
// primary.
func (bus *EventBus) GetClusterWideFilterAdditions() []FilterChange {
	return bus.filters.getFilterAdditions(clusterWide)
}

// GetLocalFilterAdditions returns the current local filter represented
// purely as additive changes. It is used by a performance standby node to
// resync its full set of event subscription filters to the active node.
func (bus *EventBus) GetLocalFilterAdditions() []FilterChange {
	return bus.filters.getFilterAdditions(bus.filters.self)
}

// ClearClusterWideFilter removes all entries from the current cluster-wide
// filter.
func (bus *EventBus) ClearClusterWideFilter() {
	bus.filters.clearClusterWidePatterns()
}

// ClearClusterNodeFilter removes all entries from the given cluster node's filter.
func (bus *EventBus) ClearClusterNodeFilter(id string) {
	bus.filters.clearClusterNodePatterns(clusterNodeID(id))
}

// NotifyOnClusterWideFilterChanges returns a channel that receives changes to
// the cluster-wide filter.
func (bus *EventBus) NotifyOnClusterWideFilterChanges(ctx context.Context) (<-chan []FilterChange, context.CancelFunc) {
	return bus.filters.watch(ctx, clusterWide)
}

// NotifyOnLocalFilterChanges returns a channel that receives changes to the filter for the current cluster node.
func (bus *EventBus) NotifyOnLocalFilterChanges(ctx context.Context) (<-chan []FilterChange, context.CancelFunc) {
	return bus.NotifyOnClusterNodeFilterChanges(ctx, string(bus.filters.self))
}

// NotifyOnClusterNodeFilterChanges returns a channel that receives changes to the filter for the given cluster node.
func (bus *EventBus) NotifyOnClusterNodeFilterChanges(ctx context.Context, clusterNode string) (<-chan []FilterChange, context.CancelFunc) {
	return bus.filters.watch(ctx, clusterNodeID(clusterNode))
}

// NewAllEventsSubscription creates a new subscription to all events.
func (bus *EventBus) NewAllEventsSubscription(ctx context.Context) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.subscribeInternal(ctx, nil, "*", "", nil)
}

// NewClusterWideSubscription creates a new subscription to all events that
// match the cluster-wide filter.
func (bus *EventBus) NewClusterWideSubscription(ctx context.Context) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	g := clusterWide
	return bus.subscribeInternal(ctx, nil, "", "", &g)
}

// NewClusterNodeSubscription creates a new subscription to all events that match the given cluster node's filter.
func (bus *EventBus) NewClusterNodeSubscription(ctx context.Context, clusterNode string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.subscribeInternal(ctx, nil, "", "", &clusterNode)
}

// NewClusterNodeSubscriptionNonLocal creates a new subscription to all events
// that match the given cluster node's filter, excluding local mount events.
func (bus *EventBus) NewClusterNodeSubscriptionNonLocal(ctx context.Context, clusterNode string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.subscribeInternal(ctx, nil, "", "source_plugin_is_local == false", &clusterNode)
}

// creates a new filter node that is tied to the filter for a given cluster node
func newClusterNodeFilterNode(filters *Filters, c clusterNodeID, bexprFilter string, coreInfo StorageInfoGetter) (*eventlogger.Filter, error) {
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
			if !filters.clusterNodeMatch(c, &namespace.Namespace{
				Path: eventNs,
			}, logical.EventType(eventRecv.EventType)) {
				return false, nil
			}
			if pluginInfo := eventRecv.GetPluginInfo(); pluginInfo != nil {
				if !coreInfo.IsReplicated(string(c), eventRecv.Namespace, pluginInfo.MountPath) {
					return false, nil
				}
			}

			// apply go-bexpr filter
			if evaluator != nil {
				return evaluator.Evaluate(eventRecv.BexprDatum())
			}
			return true, nil
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

// getSubscriberBufferSize reads the VAULT_EVENT_NOTIFICATIONS_BOUNDED_QUEUE_SIZE environment variable
// and returns the buffer size for subscriber channels. Returns 0 for unbuffered (default).
// Values above maxSubscriberBuffer (1000) are capped to prevent excessive memory usage.
func getSubscriberBufferSize() int {
	if v := os.Getenv(EnvVaultEventNotificationsBoundedQueueSize); v != "" {
		size, err := strconv.Atoi(v)
		if err != nil || size < 0 {
			// If parsing fails or negative, default to 0 (unbuffered)
			return 0
		}
		// Cap at maximum to prevent excessive memory usage
		if size > maxSubscriberBuffer {
			return maxSubscriberBuffer
		}
		return size
	}
	return 0
}

func newAsyncNode(ctx context.Context, logger hclog.Logger, broker *eventlogger.Broker, bufferSize int, removeFilter func()) *asyncChanNode {
	return &asyncChanNode{
		ctx:            ctx,
		ch:             make(chan *eventlogger.Event, bufferSize),
		logger:         logger,
		bufferSize:     bufferSize,
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

// getEventID safely extracts the event ID from an eventlogger.Event, returning "unknown" if extraction fails.
func getEventID(e *eventlogger.Event) string {
	if e == nil || e.Payload == nil {
		return "unknown"
	}
	eventReceived, ok := e.Payload.(*logical.EventReceived)
	if !ok || eventReceived == nil || eventReceived.Event == nil {
		return "unknown"
	}
	return eventReceived.Event.Id
}

func (node *asyncChanNode) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	if node.bufferSize > 0 {
		// Bounded queues: synchronous send with immediate slow consumer detection
		select {
		case node.ch <- e:
			return e, nil
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				node.logger.Info("Subscriber took too long to process event, closing", "id", getEventID(e))
				node.Close(ctx)
				return nil, ErrSlowConsumer
			}
			return e, nil
		case <-node.ctx.Done():
			return e, nil
		default:
			node.logger.Info("Subscriber queue is full, closing", "id", getEventID(e))
			node.Close(ctx)
			return nil, ErrSlowConsumer
		}
	}

	// Standard behavior: async send with goroutine per event
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
			node.logger.Info("Subscriber took too long to process event, closing", "id", getEventID(e))
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
