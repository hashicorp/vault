// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
)

// RegistrationPolicy is used to specify what kind of policy should apply when
// registering components (e.g. Pipeline, Node) with the Broker
type RegistrationPolicy string

const (
	AllowOverwrite RegistrationPolicy = "AllowOverwrite"
	DenyOverwrite  RegistrationPolicy = "DenyOverwrite"
)

// Broker is the top-level entity used in the library for configuring the system
// and for sending events.
//
// Brokers have registered Nodes which may be composed into registered Pipelines
// for EventTypes.
//
// A Node may be a filter, formatter or sink (see NodeType).
//
// A Broker may have multiple Pipelines.
//
// EventTypes may have multiple Pipelines.
//
// A Pipeline for an EventType may contain multiple filters, one formatter and
// one sink.
//
// If a Pipeline does not have a formatter, then the event will not be written
// to the Sink.
//
// A Node can be shared across multiple pipelines.
type Broker struct {
	nodes  map[NodeID]*nodeUsage
	graphs map[EventType]*graph
	lock   sync.RWMutex

	*clock
}

// nodeUsage tracks how many times a Node is referenced by registered pipelines.
type nodeUsage struct {
	node               Node
	referenceCount     int
	registrationPolicy RegistrationPolicy
}

// Option allows options to be passed as arguments.
type Option func(*options) error

// options are used to represent configuration for the broker.
type options struct {
	withPipelineRegistrationPolicy RegistrationPolicy
	withNodeRegistrationPolicy     RegistrationPolicy
}

// getDefaultOptions returns a set of default options
func getDefaultOptions() options {
	return options{
		withPipelineRegistrationPolicy: AllowOverwrite,
		withNodeRegistrationPolicy:     AllowOverwrite,
	}
}

// getOpts iterates the inbound Options and returns a struct.
// Each Option is applied in the order it appears in the argument list, so it is
// possible to supply the same Option numerous times and the 'last write wins'.
func getOpts(opt ...Option) (options, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(&opts); err != nil {
			return options{}, err
		}
	}
	return opts, nil
}

// WithPipelineRegistrationPolicy configures the option that determines the pipeline registration policy.
func WithPipelineRegistrationPolicy(policy RegistrationPolicy) Option {
	return func(o *options) error {
		var err error

		switch policy {
		case AllowOverwrite, DenyOverwrite:
			o.withPipelineRegistrationPolicy = policy
		default:
			err = fmt.Errorf("'%s' is not a valid pipeline registration policy: %w", policy, ErrInvalidParameter)
		}

		return err
	}
}

// WithNodeRegistrationPolicy configures the option that determines the node registration policy.
func WithNodeRegistrationPolicy(policy RegistrationPolicy) Option {
	return func(o *options) error {
		var err error

		switch policy {
		case AllowOverwrite, DenyOverwrite:
			o.withNodeRegistrationPolicy = policy
		default:
			err = fmt.Errorf("'%s' is not a valid node registration policy: %w", policy, ErrInvalidParameter)
		}

		return err
	}
}

// NewBroker creates a new Broker applying any relevant supplied options.
// Options are currently accepted, but none are applied.
func NewBroker(_ ...Option) (*Broker, error) {
	b := &Broker{
		nodes:  make(map[NodeID]*nodeUsage),
		graphs: make(map[EventType]*graph),
	}

	return b, nil
}

// clock only exists to make testing simpler.
type clock struct {
	now time.Time
}

// Now returns the current time
func (c *clock) Now() time.Time {
	if c == nil {
		return time.Now()
	}
	return c.now
}

// StopTimeAt allows you to "stop" the Broker's timestamp clock at a predicable
// point in time, so timestamps are predictable for testing.
func (b *Broker) StopTimeAt(now time.Time) {
	b.clock = &clock{now: now}
}

// Status describes the result of a Send.
type Status struct {
	// complete lists the IDs of 'filter' and 'sink' type nodes that successfully
	// processed the Event, resulting in immediate completion of a particular Pipeline.
	complete []NodeID
	// completeSinks lists the IDs of 'sink' type nodes that successfully processed
	// the Event, resulting in immediate completion of a particular Pipeline.
	completeSinks []NodeID
	// Warnings lists any non-fatal errors that occurred while sending an Event.
	Warnings []error
}

// Complete returns the IDs of 'filter' and 'sink' type nodes that successfully
// processed the Event, resulting in immediate completion of a particular Pipeline.
func (s Status) Complete() []NodeID {
	return s.complete
}

// CompleteSinks returns the IDs of 'sink' type nodes that successfully processed
// the Event, resulting in immediate completion of a particular Pipeline.
func (s Status) CompleteSinks() []NodeID {
	return s.completeSinks
}

func (s Status) getError(ctxErr error, threshold, thresholdSinks int) error {
	var err error
	switch {
	case len(s.complete) < threshold:
		err = fmt.Errorf("event not processed by enough 'filter' and 'sink' nodes")
	case len(s.completeSinks) < thresholdSinks:
		err = fmt.Errorf("event not processed by enough 'sink' nodes")
	default:
		return nil
	}

	return errors.Join(err, ctxErr)
}

// Send writes an event of type t to all registered pipelines concurrently and
// reports on the result.  An error will only be returned if a pipeline's delivery
// policies could not be satisfied.
func (b *Broker) Send(ctx context.Context, t EventType, payload interface{}) (Status, error) {
	b.lock.RLock()
	g, ok := b.graphs[t]
	b.lock.RUnlock()

	if !ok {
		return Status{}, fmt.Errorf("no graph for EventType %s", t)
	}

	e := &Event{
		Type:      t,
		CreatedAt: b.clock.Now(),
		Formatted: make(map[string][]byte),
		Payload:   payload,
	}

	return g.process(ctx, e)
}

// Reopen calls every registered Node's Reopen() function.  The intention is to
// ask all nodes to reopen any files they have open.  This is typically used as
// part of log rotation: after rotating, the rotator sends a signal to the
// application, which then would invoke this method.  Another typically use-case
// is to have all Nodes reevaluated any external configuration they might have.
func (b *Broker) Reopen(ctx context.Context) error {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for _, g := range b.graphs {
		if err := g.reopen(ctx); err != nil {
			return err
		}
	}

	return nil
}

// NodeID is a string that uniquely identifies a Node.
type NodeID string

// RegisterNode assigns a node ID to a node.  Node IDs should be unique. A Node
// may be a filter, formatter or sink (see NodeType). Nodes can be shared across
// multiple pipelines.
// Accepted options: WithNodeRegistrationPolicy (default: AllowOverwrite).
func (b *Broker) RegisterNode(id NodeID, node Node, opt ...Option) error {
	if id == "" {
		return fmt.Errorf("unable to register node, node ID cannot be empty: %w", ErrInvalidParameter)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return fmt.Errorf("cannot register node: %w", err)
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	nr := &nodeUsage{
		node:               node,
		referenceCount:     0,
		registrationPolicy: opts.withNodeRegistrationPolicy,
	}

	// Check if this node is already registered, if so maintain reference count
	r, exists := b.nodes[id]
	if exists {
		switch r.registrationPolicy {
		case AllowOverwrite:
			nr.referenceCount = r.referenceCount
		case DenyOverwrite:
			return fmt.Errorf("node ID %q is already registered, configured policy prevents overwriting", id)
		}
	}

	b.nodes[id] = nr

	return nil
}

// RemoveNode will remove a node from the broker, if it is not currently  in use
// This is useful if RegisterNode was used successfully prior to a failed RegisterPipeline call
// referencing those nodes
func (b *Broker) RemoveNode(ctx context.Context, id NodeID) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.removeNode(ctx, id, false)
}

// removeNode will remove a node from the broker, if it is not currently  in use.
// This is useful if RegisterNode was used successfully prior to a failed RegisterPipeline call
// referencing those nodes
// The force option can be used to decrement the count for the node if it's still in use by pipelines
// This function assumes that the caller holds a lock
func (b *Broker) removeNode(ctx context.Context, id NodeID, force bool) error {
	if id == "" {
		return fmt.Errorf("unable to remove node, node ID cannot be empty: %w", ErrInvalidParameter)
	}

	nodeUsage, ok := b.nodes[id]
	if !ok {
		return fmt.Errorf("%w: %q", ErrNodeNotFound, id)
	}

	// if force is passed, then decrement the count for this node instead of failing
	if nodeUsage.referenceCount > 0 && !force {
		return fmt.Errorf("cannot remove node, as it is still in use by 1 or more pipelines: %q", id)
	}

	var err error
	switch nodeUsage.referenceCount {
	case 0, 1:
		nc := NewNodeController(nodeUsage.node)
		if err = nc.Close(ctx); err != nil {
			err = fmt.Errorf("unable to close node ID %q: %w", id, err)
		}
		delete(b.nodes, id)
	default:
		nodeUsage.referenceCount--
	}

	return err
}

// PipelineID is a string that uniquely identifies a Pipeline within a given EventType.
type PipelineID string

// Pipeline defines a pipe: its ID, the EventType it's for, and the nodes
// that it contains. Nodes can be shared across multiple pipelines.
type Pipeline struct {
	// PipelineID uniquely identifies the Pipeline
	PipelineID PipelineID

	// EventType defines the type of event the Pipeline processes
	EventType EventType

	// NodeIDs defines Pipeline's the list of nodes
	NodeIDs []NodeID
}

// RegisterPipeline adds a pipeline to the broker.
// Accepted options: WithPipelineRegistrationPolicy (default: AllowOverwrite).
func (b *Broker) RegisterPipeline(def Pipeline, opt ...Option) error {
	err := def.validate()
	if err != nil {
		return err
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return fmt.Errorf("cannot register pipeline: %w", err)
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	g, exists := b.graphs[def.EventType]
	if !exists {
		g = &graph{}
		b.graphs[def.EventType] = g
	}

	// Get the configured policy
	pol := AllowOverwrite
	g.roots.Range(func(key PipelineID, v *registeredPipeline) bool {
		if key == def.PipelineID {
			pol = v.registrationPolicy
			return false
		}
		return true
	})

	if pol == DenyOverwrite {
		return fmt.Errorf("pipeline ID %q is already registered, configured policy prevents overwriting", def.PipelineID)
	}

	// Gather the registered nodes, so they can be referenced for this pipeline.
	nodes := make([]Node, len(def.NodeIDs))
	for i, n := range def.NodeIDs {
		nodeUsage, ok := b.nodes[n]
		if !ok {
			return fmt.Errorf("node ID %q not registered", n)
		}
		nodes[i] = nodeUsage.node
	}

	root, err := linkNodes(nodes, def.NodeIDs)
	if err != nil {
		return err
	}

	err = g.doValidate(nil, root)
	if err != nil {
		return err
	}

	// Create the pipeline registration using the optional policy (or default).
	pipelineReg := &registeredPipeline{
		rootNode:           root,
		registrationPolicy: opts.withPipelineRegistrationPolicy,
	}

	// Store the pipeline and then update the reference count of the nodes in that pipeline.
	g.roots.Store(def.PipelineID, pipelineReg)
	for _, id := range def.NodeIDs {
		nodeUsage, ok := b.nodes[id]
		// We can be optimistic about this as we would have already errored above.
		if ok {
			nodeUsage.referenceCount++
		}
	}

	return nil
}

// RemovePipeline removes a pipeline from the broker.
func (b *Broker) RemovePipeline(t EventType, id PipelineID) error {
	switch {
	case t == "":
		return errors.New("event type cannot be empty")
	case id == "":
		return errors.New("pipeline ID cannot be empty")
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	g, ok := b.graphs[t]
	if !ok {
		return fmt.Errorf("no graph for EventType %s", t)
	}

	g.roots.Delete(id)
	return nil
}

// RemovePipelineAndNodes will attempt to remove all nodes referenced by the pipeline.
// Any nodes that are referenced by other pipelines will not be removed.
//
// Failed preconditions will result in a return of false with an error and
// neither the pipeline nor nodes will be deleted.
//
// Once we start deleting the pipeline and nodes, we will continue until completion,
// but we'll return true along with any errors encountered (as multierror.Error).
func (b *Broker) RemovePipelineAndNodes(ctx context.Context, t EventType, id PipelineID) (bool, error) {
	switch {
	case t == "":
		return false, errors.New("event type cannot be empty")
	case id == "":
		return false, errors.New("pipeline ID cannot be empty")
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	g, ok := b.graphs[t]
	if !ok {
		return false, fmt.Errorf("no graph for EventType %s", t)
	}

	nodes, err := g.roots.Nodes(id)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve all nodes referenced by pipeline ID %q: %w", id, err)
	}

	g.roots.Delete(id)

	var nodeErr error

	for _, nodeID := range nodes {
		err = b.removeNode(ctx, nodeID, true)
		if err != nil {
			nodeErr = multierror.Append(nodeErr, err)
		}
	}

	return true, nodeErr
}

// SetSuccessThreshold sets the success threshold per EventType.  For the
// overall processing of a given event to be considered a success, at least as
// many pipelines as the threshold value must successfully process the event.
// This means that a filter could of course filter an event before it reaches
// the pipeline's sink, but it would still count as success when it comes to
// meeting this threshold.  Use this when you want to allow the filtering of
// events without causing an error because an event was filtered.
func (b *Broker) SetSuccessThreshold(t EventType, successThreshold int) error {
	switch {
	case t == "":
		return errors.New("event type cannot be empty")
	case successThreshold < 0:
		return fmt.Errorf("successThreshold must be 0 or greater")
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	g, ok := b.graphs[t]
	if !ok {
		g = &graph{}
		b.graphs[t] = g
	}

	g.successThreshold = successThreshold
	return nil
}

// SetSuccessThresholdSinks sets the success threshold per EventType.  For the
// overall processing of a given event to be considered a success, at least as
// many sinks as the threshold value must successfully process the event.
func (b *Broker) SetSuccessThresholdSinks(t EventType, successThresholdSinks int) error {
	switch {
	case t == "":
		return errors.New("event type cannot be empty")
	case successThresholdSinks < 0:
		return fmt.Errorf("successThresholdSinks must be 0 or greater")
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	g, ok := b.graphs[t]
	if !ok {
		g = &graph{}
		b.graphs[t] = g
	}

	g.successThresholdSinks = successThresholdSinks
	return nil
}

// SuccessThreshold returns the configured success threshold per EventType.
// For the overall processing of a given event to be considered a success, at least
// as many filter or sink nodes as the threshold value must successfully process
// the event.
// The threshold is returned (default: 0), along with a boolean indicating whether
// the EventType was registered with the broker, if true, the threshold is accurate
// for the specified EventType.
func (b *Broker) SuccessThreshold(t EventType) (int, bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	g, ok := b.graphs[t]
	if ok {
		return g.successThreshold, true
	}

	return 0, false
}

// SuccessThresholdSinks returns the configured success threshold per EventType.
// For the overall processing of a given event to be considered a success, at least
// as many sink nodes as the threshold value must successfully process the event.
// The threshold is returned (default: 0), along with a boolean indicating whether
// the EventType was registered with the broker, if true, the threshold is accurate
// for the specified EventType.
func (b *Broker) SuccessThresholdSinks(t EventType) (int, bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	g, ok := b.graphs[t]
	if ok {
		return g.successThresholdSinks, true
	}

	return 0, false
}

// IsAnyPipelineRegistered returns whether a pipeline for a given event type is already registered or not.
func (b *Broker) IsAnyPipelineRegistered(e EventType) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	g, found := b.graphs[e]
	if !found {
		return false
	}

	found = false
	g.roots.Range(func(_ PipelineID, pipeline *registeredPipeline) bool {
		found = true
		return false
	})
	return found
}

// validate ensures that the Pipeline has the required configuration to allow
// registration, removal or usage, without issue.
func (p Pipeline) validate() error {
	var err error

	if p.PipelineID == "" {
		err = multierror.Append(err, errors.New("pipeline ID is required"))
	}

	if p.EventType == "" {
		err = multierror.Append(err, errors.New("event type is required"))
	}

	if len(p.NodeIDs) == 0 {
		err = multierror.Append(err, errors.New("node IDs are required"))
	}

	for _, n := range p.NodeIDs {
		if n == "" {
			err = multierror.Append(err, errors.New("node ID cannot be empty"))
			break
		}
	}

	return err
}
