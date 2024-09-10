// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	nshelper "github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// timeout is the duration which should be used for context related timeouts.
	timeout = 10 * time.Second
)

var (
	_ Registrar = (*Broker)(nil)
	_ Auditor   = (*Broker)(nil)
)

// Registrar interface describes a means to register and deregister audit devices.
type Registrar interface {
	Register(backend Backend, local bool) error
	Deregister(ctx context.Context, name string) error
	IsRegistered(name string) bool
	IsLocal(name string) (bool, error)
}

// Auditor interface describes methods which can be used to perform auditing.
type Auditor interface {
	LogRequest(ctx context.Context, input *logical.LogInput) error
	LogResponse(ctx context.Context, input *logical.LogInput) error
	GetHash(ctx context.Context, name string, input string) (string, error)
	Invalidate(ctx context.Context, key string)
}

// backendEntry composes a backend with additional settings.
type backendEntry struct {
	// backend is the underlying audit backend.
	backend Backend

	// local indicates whether this audit backend should be local to the Vault cluster.
	local bool
}

// Broker represents an audit broker which performs actions such as registering/de-registering
// backends and logging audit entries for a request or response.
// NOTE: NewBroker should be used to initialize the Broker struct.
type Broker struct {
	*brokerEnt

	sync.RWMutex

	logger hclog.Logger

	// backends is the map of audit device name to {thing}
	backends map[string]backendEntry

	// broker is used to register pipelines for audit devices.
	broker *eventlogger.Broker
}

// NewBroker initializes a broker, which can be used to perform audit logging.
func NewBroker(logger hclog.Logger) (*Broker, error) {
	if logger == nil || reflect.ValueOf(logger).IsNil() {
		return nil, fmt.Errorf("cannot create a new audit broker with nil logger: %w", ErrInvalidParameter)
	}

	eventBroker, err := eventlogger.NewBroker()
	if err != nil {
		return nil, fmt.Errorf("error creating event broker for audit events: %w", err)
	}

	ent, err := newBrokerEnt()
	if err != nil {
		return nil, fmt.Errorf("error creating audit broker extentions: %w", err)
	}

	return &Broker{
		backends:  make(map[string]backendEntry),
		broker:    eventBroker,
		brokerEnt: ent,
		logger:    logger,
	}, nil
}

// hasAuditPipelines can be used as a shorthand to check if a broker has any
// registered pipelines that are for the audit event type.
func hasAuditPipelines(broker *eventlogger.Broker) bool {
	return broker.IsAnyPipelineRegistered(event.AuditType.AsEventType())
}

// isRegistered is used to check if a given audit backend is registered.
// This method should be used within the broker to prevent locking issues.
func (b *Broker) isRegistered(backend Backend) error {
	if b.isRegisteredByName(backend.Name()) {
		return fmt.Errorf("backend already registered '%s': %w", backend.Name(), ErrExternalOptions)
	}

	if err := b.validateRegistrationRequest(backend); err != nil {
		return err
	}

	return nil
}

// isRegisteredByName returns a boolean to indicate whether an audit backend is
// registered with the broker.
func (b *Broker) isRegisteredByName(name string) bool {
	_, ok := b.backends[name]
	return ok
}

// register can be used to register a normal audit device, it will also calculate
// and configure the success threshold required for sinks.
// NOTE: register assumes that the backend which is being registered has not yet
// been added to the broker's backends.
func (b *Broker) register(backend Backend) error {
	err := registerNodesAndPipeline(b.broker, backend)
	if err != nil {
		return fmt.Errorf("audit pipeline registration error: %w", err)
	}

	threshold := 0
	if !backend.HasFiltering() {
		threshold = 1
	} else {
		threshold = b.requiredSuccessThresholdSinks()
	}

	// Update the success threshold now that the pipeline is registered.
	err = b.broker.SetSuccessThresholdSinks(event.AuditType.AsEventType(), threshold)
	if err != nil {
		return fmt.Errorf("unable to configure sink success threshold (%d): %w", threshold, err)
	}

	return nil
}

// deregister can be used to deregister an audit device, it will also configure
// the success threshold required for sinks.
// NOTE: deregister assumes that the backend which is being deregistered has already
// been removed from the broker's backends.
func (b *Broker) deregister(ctx context.Context, name string) error {
	threshold := b.requiredSuccessThresholdSinks()

	err := b.broker.SetSuccessThresholdSinks(event.AuditType.AsEventType(), threshold)
	if err != nil {
		return fmt.Errorf("unable to reconfigure sink success threshold (%d): %w", threshold, err)
	}

	// The first return value, a bool, indicates whether
	// RemovePipelineAndNodes encountered the error while evaluating
	// pre-conditions (false) or once it started removing the pipeline and
	// the nodes (true). This code doesn't care either way.
	_, err = b.broker.RemovePipelineAndNodes(ctx, event.AuditType.AsEventType(), eventlogger.PipelineID(name))
	if err != nil {
		return fmt.Errorf("unable to remove pipeline and nodes: %w", err)
	}

	return nil
}

// registerNodesAndPipeline registers eventlogger nodes and a pipeline with the
// backend's name, on the specified eventlogger.Broker using the Backend to supply them.
func registerNodesAndPipeline(broker *eventlogger.Broker, b Backend) error {
	for id, node := range b.Nodes() {
		err := broker.RegisterNode(id, node)
		if err != nil {
			return fmt.Errorf("unable to register nodes for %q: %w", b.Name(), err)
		}
	}

	pipeline := eventlogger.Pipeline{
		PipelineID: eventlogger.PipelineID(b.Name()),
		EventType:  b.EventType(),
		NodeIDs:    b.NodeIDs(),
	}

	err := broker.RegisterPipeline(pipeline)
	if err != nil {
		return fmt.Errorf("unable to register pipeline for %q: %w", b.Name(), err)
	}

	return nil
}

func (b *Broker) Register(backend Backend, local bool) error {
	b.Lock()
	defer b.Unlock()

	if backend == nil || reflect.ValueOf(backend).IsNil() {
		return fmt.Errorf("backend cannot be nil: %w", ErrInvalidParameter)
	}

	// If the backend is already registered, we cannot re-register it.
	err := b.isRegistered(backend)
	if err != nil {
		return err
	}

	if err := b.handlePipelineRegistration(backend); err != nil {
		return err
	}

	b.backends[backend.Name()] = backendEntry{
		backend: backend,
		local:   local,
	}

	return nil
}

func (b *Broker) Deregister(ctx context.Context, name string) error {
	b.Lock()
	defer b.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required: %w", ErrInvalidParameter)
	}

	// If the backend isn't actually registered, then there's nothing to do.
	// We don't return any error so that Deregister can be idempotent.
	if !b.isRegisteredByName(name) {
		return nil
	}

	// Remove the Backend from the map first, so that if an error occurs while
	// removing the pipeline and nodes, we can quickly exit this method with
	// the error.
	delete(b.backends, name)

	if err := b.handlePipelineDeregistration(ctx, name); err != nil {
		return err
	}

	return nil
}

// LogRequest is used to ensure all the audit backends have an opportunity to
// log the given request and that *at least one* succeeds.
func (b *Broker) LogRequest(ctx context.Context, in *logical.LogInput) (retErr error) {
	b.RLock()
	defer b.RUnlock()

	// If no backends are registered then we have no devices to log the request.
	if len(b.backends) < 1 {
		return nil
	}

	defer metrics.MeasureSince([]string{"audit", "log_request"}, time.Now())
	defer func() {
		metricVal := float32(0.0)
		if retErr != nil {
			metricVal = 1.0
		}
		metrics.IncrCounter([]string{"audit", "log_request_failure"}, metricVal)
	}()

	e, err := newEvent(RequestType)
	if err != nil {
		return err
	}

	e.Data = in

	// Get a context to use for auditing.
	auditContext, auditCancel, err := getAuditContext(ctx)
	if err != nil {
		return err
	}
	defer auditCancel()

	var status eventlogger.Status
	if hasAuditPipelines(b.broker) {
		status, err = b.broker.Send(auditContext, event.AuditType.AsEventType(), e)
		if err != nil {
			return errors.Join(append([]error{err}, status.Warnings...)...)
		}
	}

	// Audit event ended up in at least 1 sink.
	if len(status.CompleteSinks()) > 0 {
		// We should log warnings to the operational logs regardless of whether
		// we consider the overall auditing attempt to be successful.
		if len(status.Warnings) > 0 {
			b.logger.Error("log request underlying pipeline error(s)", "error", errors.Join(status.Warnings...))
		}

		return nil
	}

	// There were errors from inside the pipeline and we didn't write to a sink.
	if len(status.Warnings) > 0 {
		return fmt.Errorf("error during audit pipeline processing: %w", errors.Join(status.Warnings...))
	}

	// Handle any additional audit that is required (Enterprise/CE dependant).
	err = b.handleAdditionalAudit(auditContext, e)
	if err != nil {
		return err
	}

	return nil
}

// LogResponse is used to ensure all the audit backends have an opportunity to
// log the given response and that *at least one* succeeds.
func (b *Broker) LogResponse(ctx context.Context, in *logical.LogInput) (retErr error) {
	b.RLock()
	defer b.RUnlock()

	// If no backends are registered then we have no devices to send audit entries to.
	if len(b.backends) < 1 {
		return nil
	}

	defer metrics.MeasureSince([]string{"audit", "log_response"}, time.Now())
	defer func() {
		metricVal := float32(0.0)
		if retErr != nil {
			metricVal = 1.0
		}
		metrics.IncrCounter([]string{"audit", "log_response_failure"}, metricVal)
	}()

	e, err := newEvent(ResponseType)
	if err != nil {
		return err
	}

	e.Data = in

	// Get a context to use for auditing.
	auditContext, auditCancel, err := getAuditContext(ctx)
	if err != nil {
		return err
	}
	defer auditCancel()

	var status eventlogger.Status
	if hasAuditPipelines(b.broker) {
		status, err = b.broker.Send(auditContext, event.AuditType.AsEventType(), e)
		if err != nil {
			return errors.Join(append([]error{err}, status.Warnings...)...)
		}
	}

	// Audit event ended up in at least 1 sink.
	if len(status.CompleteSinks()) > 0 {
		// We should log warnings to the operational logs regardless of whether
		// we consider the overall auditing attempt to be successful.
		if len(status.Warnings) > 0 {
			b.logger.Error("log response underlying pipeline error(s)", "error", errors.Join(status.Warnings...))
		}

		return nil
	}

	// There were errors from inside the pipeline and we didn't write to a sink.
	if len(status.Warnings) > 0 {
		return fmt.Errorf("error during audit pipeline processing: %w", errors.Join(status.Warnings...))
	}

	// Handle any additional audit that is required (Enterprise/CE dependant).
	err = b.handleAdditionalAudit(auditContext, e)
	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) Invalidate(ctx context.Context, _ string) {
	// For now, we ignore the key as this would only apply to salts.
	// We just sort of brute force it on each one.
	b.Lock()
	defer b.Unlock()

	for _, be := range b.backends {
		be.backend.Invalidate(ctx)
	}
}

// IsLocal is used to check if a given audit backend is registered
func (b *Broker) IsLocal(name string) (bool, error) {
	b.RLock()
	defer b.RUnlock()

	be, ok := b.backends[name]
	if ok {
		return be.local, nil
	}

	return false, fmt.Errorf("unknown audit backend %q", name)
}

// GetHash returns a hash using the salt of the given backend
func (b *Broker) GetHash(ctx context.Context, name string, input string) (string, error) {
	b.RLock()
	defer b.RUnlock()

	be, ok := b.backends[name]
	if !ok {
		return "", fmt.Errorf("unknown audit backend %q", name)
	}

	return hashString(ctx, be.backend, input)
}

// IsRegistered is used to check if a given audit backend is registered.
func (b *Broker) IsRegistered(name string) bool {
	b.RLock()
	defer b.RUnlock()

	return b.isRegisteredByName(name)
}

// getAuditContext extracts the namespace from the specified context and returns
// a new context and cancelation function, completely detached from the original
// with a timeout.
// NOTE: When error is nil, the context.CancelFunc returned from this function
// should be deferred immediately by the caller to prevent resource leaks.
func getAuditContext(ctx context.Context) (context.Context, context.CancelFunc, error) {
	ns, err := nshelper.FromContext(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("namespace missing from context: %w", err)
	}

	tempContext := nshelper.ContextWithNamespace(context.Background(), ns)
	auditContext, auditCancel := context.WithTimeout(tempContext, timeout)

	return auditContext, auditCancel, nil
}
