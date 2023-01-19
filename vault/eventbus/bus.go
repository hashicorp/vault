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

type asyncChanNode struct {
	// TODO: add bounded deque buffer of *any
	ch chan *logical.EventData
}

var (
	_ eventlogger.Node    = (*asyncChanNode)(nil)
	_ logical.EventSender = (*EventBus)(nil)
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

// Send sends an event to the event bus and routes it to all relevant subscribers.
// This function does *not* wait for all subscribers to acknowledge before returning.
func (bus *EventBus) Send(ctx context.Context, eventType logical.EventType, data *logical.EventData) error {
	if !bus.started.Load() {
		return ErrNotStarted
	}
	bus.logger.Info("Sending event", "event", data)
	_, err := bus.broker.Send(ctx, eventlogger.EventType(eventType), data)
	if err != nil {
		// if no listeners for this event type are registered, that's okay, the event
		// will just not be sent anywhere
		if strings.Contains(strings.ToLower(err.Error()), "no graph for eventtype") {
			return nil
		}
	}
	return err
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

func (bus *EventBus) Subscribe(_ context.Context, eventType logical.EventType) (chan *logical.EventData, error) {
	// subscriptions are still stored even if the bus has not been started
	pipelineID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	nodeID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// TODO: should we have just one node, and handle all the routing ourselves?
	asyncNode := newAsyncNode()
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

func newAsyncNode() *asyncChanNode {
	return &asyncChanNode{
		ch: make(chan *logical.EventData),
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
		select {
		case node.ch <- e.Payload.(*logical.EventData):
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
