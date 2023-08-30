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
	ErrNotStarted              = errors.New("event broker has not been started")
	cloudEventsFormatterFilter *cloudevents.FormatterFilter
	subscriptions              atomic.Int64 // keeps track of event subscription count in all event buses

	// these metadata fields will have the plugin mount path prepended to them
	metadataPrependPathFields = []string{
		"path",
		logical.EventMetadataDataPath,
	}
)

// EventBus contains the main logic of running an event broker for Vault.
// Start() must be called before the EventBus will accept events for sending.
type EventBus struct {
	logger          hclog.Logger
	broker          *eventlogger.Broker
	started         atomic.Bool
	formatterNodeID eventlogger.NodeID
	timeout         time.Duration
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
	closeOnce  sync.Once
	cancelFunc context.CancelFunc
	pipelineID eventlogger.PipelineID
	broker     *eventlogger.Broker
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
func (bus *EventBus) SendEventInternal(ctx context.Context, ns *namespace.Namespace, pluginInfo *logical.EventPluginInfo, eventType logical.EventType, data *logical.EventData) error {
	if ns == nil {
		return namespace.ErrNoNamespace
	}
	if !bus.started.Load() {
		return ErrNotStarted
	}
	eventReceived := &logical.EventReceived{
		Event:      patchMountPath(data, pluginInfo),
		Namespace:  ns.Path,
		EventType:  string(eventType),
		PluginInfo: pluginInfo,
	}
	bus.logger.Info("Sending event", "event", eventReceived)

	// We can't easily know when the SendEvent is complete, so we can't call the cancel function.
	// But, it is called automatically after bus.timeout, so there won't be any leak as long as bus.timeout is not too long.
	ctx, _ = context.WithTimeout(ctx, bus.timeout)
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
func (bus *pluginEventBus) SendEvent(ctx context.Context, eventType logical.EventType, data *logical.EventData) error {
	return bus.bus.SendEventInternal(ctx, bus.namespace, bus.pluginInfo, eventType, data)
}

func init() {
	// TODO: maybe this should relate to the Vault core somehow?
	sourceUrl, err := url.Parse("https://vaultproject.io/")
	if err != nil {
		panic(err)
	}
	cloudEventsFormatterFilter = &cloudevents.FormatterFilter{
		Source: sourceUrl,
		Predicate: func(_ context.Context, e interface{}) (bool, error) {
			return true, nil
		},
	}
}

func NewEventBus(logger hclog.Logger) (*EventBus, error) {
	broker, err := eventlogger.NewBroker()
	if err != nil {
		return nil, err
	}

	formatterID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	formatterNodeID := eventlogger.NodeID(formatterID)
	err = broker.RegisterNode(formatterNodeID, cloudEventsFormatterFilter)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = hclog.Default().Named("events")
	}

	return &EventBus{
		logger:          logger,
		broker:          broker,
		formatterNodeID: formatterNodeID,
		timeout:         defaultTimeout,
	}, nil
}

func (bus *EventBus) Subscribe(ctx context.Context, ns *namespace.Namespace, pattern string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	return bus.SubscribeMultipleNamespaces(ctx, []string{strings.Trim(ns.Path, "/")}, pattern)
}

func (bus *EventBus) SubscribeMultipleNamespaces(ctx context.Context, namespacePathPatterns []string, pattern string) (<-chan *eventlogger.Event, context.CancelFunc, error) {
	// subscriptions are still stored even if the bus has not been started
	pipelineID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	filterNodeID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	filterNode := newFilterNode(namespacePathPatterns, pattern)
	err = bus.broker.RegisterNode(eventlogger.NodeID(filterNodeID), filterNode)
	if err != nil {
		return nil, nil, err
	}

	sinkNodeID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	asyncNode := newAsyncNode(ctx, bus.logger)
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

func newFilterNode(namespacePatterns []string, pattern string) *eventlogger.Filter {
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

			// Filter for correct event type, including wildcards.
			if !glob.Glob(pattern, eventRecv.EventType) {
				return false, nil
			}

			return true, nil
		},
	}
}

func newAsyncNode(ctx context.Context, logger hclog.Logger) *asyncChanNode {
	return &asyncChanNode{
		ctx:    ctx,
		ch:     make(chan *eventlogger.Event),
		logger: logger,
	}
}

// Close tells the bus to stop sending us events.
func (node *asyncChanNode) Close(ctx context.Context) {
	node.closeOnce.Do(func() {
		defer node.cancelFunc()
		if node.broker != nil {
			isPipelineRemoved, err := node.broker.RemovePipelineAndNodes(ctx, eventTypeAll, node.pipelineID)

			switch {
			case err != nil && isPipelineRemoved:
				msg := fmt.Sprintf("Error removing nodes referenced by pipeline %q", node.pipelineID)
				node.logger.Warn(msg, err)
			case err != nil:
				msg := fmt.Sprintf("Error removing pipeline %q", node.pipelineID)
				node.logger.Warn(msg, err)
			}
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
