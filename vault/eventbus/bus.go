// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventbus

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/eventlogger/formatter_filters/cloudevents"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

var ErrNotStarted = errors.New("event broker has not been started")

var cloudEventsFormatterFilter *cloudevents.FormatterFilter

// EventBus contains the main logic of running an event broker for Vault.
// Start() must be called before the EventBus will accept events for sending.
type EventBus struct {
	logger          hclog.Logger
	broker          *eventlogger.Broker
	started         atomic.Bool
	formatterNodeID eventlogger.NodeID
}

type pluginEventBus struct {
	bus        *EventBus
	namespace  *namespace.Namespace
	pluginInfo *logical.EventPluginInfo
}

type asyncChanNode struct {
	// TODO: add bounded deque buffer of *EventReceived
	ch        chan *logical.EventReceived
	namespace *namespace.Namespace
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

// SendInternal sends an event to the event bus and routes it to all relevant subscribers.
// This function does *not* wait for all subscribers to acknowledge before returning.
// This function is meant to be used by trusted internal code, so it can specify details like the namespace
// and plugin info. Events from plugins should be routed through WithPlugin(), which will populate
// the namespace and plugin info automatically.
func (bus *EventBus) SendInternal(ctx context.Context, ns *namespace.Namespace, pluginInfo *logical.EventPluginInfo, eventType logical.EventType, data *logical.EventData) error {
	if ns == nil {
		return namespace.ErrNoNamespace
	}
	if !bus.started.Load() {
		return ErrNotStarted
	}
	eventReceived := &logical.EventReceived{
		Event:      data,
		Namespace:  ns.Path,
		EventType:  string(eventType),
		PluginInfo: pluginInfo,
	}
	bus.logger.Info("Sending event", "event", eventReceived)
	_, err := bus.broker.Send(ctx, eventlogger.EventType(eventType), eventReceived)
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

// Send sends an event to the event bus and routes it to all relevant subscribers.
// This function does *not* wait for all subscribers to acknowledge before returning.
func (bus *pluginEventBus) Send(ctx context.Context, eventType logical.EventType, data *logical.EventData) error {
	return bus.bus.SendInternal(ctx, bus.namespace, bus.pluginInfo, eventType, data)
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
	broker := eventlogger.NewBroker()

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
	}, nil
}

func (bus *EventBus) Subscribe(_ context.Context, ns *namespace.Namespace, eventType logical.EventType) (chan *logical.EventReceived, error) {
	// subscriptions are still stored even if the bus has not been started
	pipelineID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	nodeID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// TODO: should we have just one node per namespace, and handle all the routing ourselves?
	asyncNode := newAsyncNode(ns)
	err = bus.broker.RegisterNode(eventlogger.NodeID(nodeID), asyncNode)
	if err != nil {
		defer asyncNode.Close()
		return nil, err
	}

	nodes := []eventlogger.NodeID{bus.formatterNodeID, eventlogger.NodeID(nodeID)}

	pipeline := eventlogger.Pipeline{
		PipelineID: eventlogger.PipelineID(pipelineID),
		EventType:  eventlogger.EventType(eventType),
		NodeIDs:    nodes,
	}
	err = bus.broker.RegisterPipeline(pipeline)
	if err != nil {
		defer asyncNode.Close()
		return nil, err
	}
	return asyncNode.ch, nil
}

func newAsyncNode(namespace *namespace.Namespace) *asyncChanNode {
	return &asyncChanNode{
		ch:        make(chan *logical.EventReceived),
		namespace: namespace,
	}
}

func (node *asyncChanNode) Close() error {
	close(node.ch)
	return nil
}

func (node *asyncChanNode) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	// TODO: add timeout on sending to node.ch
	// sends to the channel async in another goroutine
	go func() {
		eventRecv := e.Payload.(*logical.EventReceived)
		// drop if event is not in our namespace
		// TODO: add wildcard processing here in some cases?
		if eventRecv.Namespace != node.namespace.Path {
			return
		}
		select {
		case node.ch <- eventRecv:
		case <-ctx.Done():
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
