// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

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
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// timeout is the duration which should be used for context related timeouts.
	timeout = 5 * time.Second
)

type backendEntry struct {
	backend audit.Backend
	local   bool
}

// AuditBroker is used to provide a single ingest interface to auditable
// events given that multiple backends may be configured.
type AuditBroker struct {
	sync.RWMutex
	backends map[string]backendEntry

	// broker is used to register pipelines for all devices except a fallback device.
	broker *eventlogger.Broker

	// fallbackBroker is used to register a pipeline to be used as a fallback
	// in situations where we cannot use the eventlogger.Broker to guarantee that
	// the required number of sinks were successfully written to. This situation
	// occurs when all the audit devices registered with the broker use filtering.
	// NOTE: there should only ever be a single device registered on the fallbackBroker.
	fallbackBroker *eventlogger.Broker

	// fallbackName stores the name (path) of the audit device which has been configured
	// as the fallback pipeline (its eventlogger.PipelineID).
	fallbackName string
	logger       hclog.Logger
}

// NewAuditBroker creates a new audit broker
func NewAuditBroker(log hclog.Logger) (*AuditBroker, error) {
	eventBroker, err := eventlogger.NewBroker()
	if err != nil {
		return nil, fmt.Errorf("error creating event broker for audit events: %w", err)
	}

	// Set up the broker that will support a single fallback device.
	fallbackEventBroker, err := eventlogger.NewBroker()
	if err != nil {
		return nil, fmt.Errorf("error creating event fallback broker for audit event: %w", err)
	}

	broker := &AuditBroker{
		backends:       make(map[string]backendEntry),
		logger:         log,
		broker:         eventBroker,
		fallbackBroker: fallbackEventBroker,
	}

	return broker, nil
}

// Register is used to add new audit backend to the broker
func (a *AuditBroker) Register(name string, backend audit.Backend, local bool) error {
	a.Lock()
	defer a.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required: %w", audit.ErrInvalidParameter)
	}

	if backend == nil || reflect.ValueOf(backend).IsNil() {
		return fmt.Errorf("backend cannot be nil: %w", audit.ErrInvalidParameter)
	}

	// If the backend is already registered, we cannot re-register it.
	if a.isRegistered(name) {
		return fmt.Errorf("backend already registered '%s': %w", name, audit.ErrExternalOptions)
	}

	// Fallback devices are singleton instances, we cannot register more than one or overwrite the existing one.
	if backend.IsFallback() && hasAuditPipelines(a.fallbackBroker) {
		// Get the name of the fallback device which is registered with the broker.
		var existing string
		for _, be := range a.backends {
			if be.backend.IsFallback() {
				existing = be.backend.Name()
			}
		}
		if existing == "" {
			// We expected an existing fallback device but didn't find it.
			return fmt.Errorf("cannot determine name of existing registered fallback device: %w", audit.ErrInternal)
		}

		return fmt.Errorf("existing fallback device already registered %q: %w", existing, audit.ErrInvalidParameter)
	}

	if name != backend.Name() {
		return fmt.Errorf("audit registration failed due to device name mismatch: %q, %q: %w", name, backend.Name(), audit.ErrInternal)
	}

	switch {
	case backend.IsFallback():
		err := a.registerFallback(name, backend)
		if err != nil {
			return fmt.Errorf("unable to register fallback device for %q: %w: %w", name, err, audit.ErrInternal)
		}
	default:
		err := a.register(name, backend)
		if err != nil {
			return fmt.Errorf("unable to register device for %q: %w", name, err)
		}
	}

	a.backends[name] = backendEntry{
		backend: backend,
		local:   local,
	}

	return nil
}

// Deregister is used to remove an audit backend from the broker
func (a *AuditBroker) Deregister(ctx context.Context, name string) error {
	a.Lock()
	defer a.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required: %w", audit.ErrInvalidParameter)
	}

	// If the backend isn't actually registered, then there's nothing to do.
	// We don't return any error so that Deregister can be idempotent.
	if !a.isRegistered(name) {
		return nil
	}

	// Remove the Backend from the map first, so that if an error occurs while
	// removing the pipeline and nodes, we can quickly exit this method with
	// the error.
	delete(a.backends, name)

	var err error
	switch {
	case name == a.fallbackName:
		err = a.deregisterFallback(ctx, name)
	default:
		err = a.deregister(ctx, name)
	}

	if err != nil {
		return fmt.Errorf("deregistration failed for audit device %q: %w: %w", name, err, audit.ErrInternal)
	}

	return nil
}

// IsRegistered is used to check if a given audit backend is registered.
func (a *AuditBroker) IsRegistered(name string) bool {
	a.RLock()
	defer a.RUnlock()

	return a.isRegistered(name)
}

// isRegistered is used to check if a given audit backend is registered.
// This method should be used within the AuditBroker to prevent locking issues.
func (a *AuditBroker) isRegistered(name string) bool {
	_, ok := a.backends[name]
	return ok
}

// IsLocal is used to check if a given audit backend is registered
func (a *AuditBroker) IsLocal(name string) (bool, error) {
	a.RLock()
	defer a.RUnlock()

	be, ok := a.backends[name]
	if ok {
		return be.local, nil
	}

	return false, fmt.Errorf("unknown audit backend %q", name)
}

// GetHash returns a hash using the salt of the given backend
func (a *AuditBroker) GetHash(ctx context.Context, name string, input string) (string, error) {
	a.RLock()
	defer a.RUnlock()

	be, ok := a.backends[name]
	if !ok {
		return "", fmt.Errorf("unknown audit backend %q", name)
	}

	return audit.HashString(ctx, be.backend, input)
}

// LogRequest is used to ensure all the audit backends have an opportunity to
// log the given request and that *at least one* succeeds.
func (a *AuditBroker) LogRequest(ctx context.Context, in *logical.LogInput) (ret error) {
	a.RLock()
	defer a.RUnlock()

	// If no backends are registered then we have no devices to log the request.
	if len(a.backends) < 1 {
		return nil
	}

	defer metrics.MeasureSince([]string{"audit", "log_request"}, time.Now())
	defer func() {
		metricVal := float32(0.0)
		if ret != nil {
			metricVal = 1.0
		}
		metrics.IncrCounter([]string{"audit", "log_request_failure"}, metricVal)
	}()

	var retErr *multierror.Error

	e, err := audit.NewEvent(audit.RequestType)
	if err != nil {
		retErr = multierror.Append(retErr, err)
		return retErr.ErrorOrNil()
	}

	e.Data = in

	// Evaluate whether we need a new context for auditing.
	var auditContext context.Context
	if isContextViable(ctx) {
		auditContext = ctx
	} else {
		// In cases where we are trying to audit the request, and the existing
		// context is not viable due to a short deadline, we detach ourselves from
		// the original context (keeping only the namespace).
		// This is so that we get a fair run at writing audit entries if Vault
		// has taken up a lot of time handling the request before audit (request)
		// is triggered. Pipeline nodes and the eventlogger.Broker may check for a
		// cancelled context and refuse to process the nodes further.
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("namespace missing from context: %w", err))
			return retErr.ErrorOrNil()
		}

		tempContext, auditCancel := context.WithTimeout(context.Background(), timeout)
		defer auditCancel()
		auditContext = namespace.ContextWithNamespace(tempContext, ns)
	}

	var status eventlogger.Status
	if hasAuditPipelines(a.broker) {
		status, err = a.broker.Send(auditContext, event.AuditType.AsEventType(), e)
		if err != nil {
			retErr = multierror.Append(retErr, multierror.Append(err, status.Warnings...))
			return retErr.ErrorOrNil()
		}
	}

	// Audit event ended up in at least 1 sink.
	if len(status.CompleteSinks()) > 0 {
		return retErr.ErrorOrNil()
	}

	// There were errors from inside the pipeline and we didn't write to a sink.
	if len(status.Warnings) > 0 {
		retErr = multierror.Append(retErr, multierror.Append(errors.New("error during audit pipeline processing"), status.Warnings...))
		return retErr.ErrorOrNil()
	}

	// If a fallback device is registered we can rely on that to 'catch all'
	// and also the broker level guarantee for completed sinks.
	if a.fallbackBroker.IsAnyPipelineRegistered(event.AuditType.AsEventType()) {
		status, err = a.fallbackBroker.Send(auditContext, event.AuditType.AsEventType(), e)
		if err != nil {
			retErr = multierror.Append(retErr, multierror.Append(fmt.Errorf("auditing request to fallback device failed: %w", err), status.Warnings...))
		}
	} else {
		// This audit event won't make it to any devices, we class this as a 'miss' for auditing.
		metrics.IncrCounter(audit.MetricLabelsFallbackMiss(), 1)
	}

	return retErr.ErrorOrNil()
}

// LogResponse is used to ensure all the audit backends have an opportunity to
// log the given response and that *at least one* succeeds.
func (a *AuditBroker) LogResponse(ctx context.Context, in *logical.LogInput) (ret error) {
	a.RLock()
	defer a.RUnlock()

	// If no backends are registered then we have no devices to send audit entries to.
	if len(a.backends) < 1 {
		return nil
	}

	defer metrics.MeasureSince([]string{"audit", "log_response"}, time.Now())
	defer func() {
		metricVal := float32(0.0)
		if ret != nil {
			metricVal = 1.0
		}
		metrics.IncrCounter([]string{"audit", "log_response_failure"}, metricVal)
	}()

	var retErr *multierror.Error

	e, err := audit.NewEvent(audit.ResponseType)
	if err != nil {
		retErr = multierror.Append(retErr, err)
		return retErr.ErrorOrNil()
	}

	e.Data = in

	// Evaluate whether we need a new context for auditing.
	var auditContext context.Context
	if isContextViable(ctx) {
		auditContext = ctx
	} else {
		// In cases where we are trying to audit the response, and the existing
		// context is not viable due to a short deadline, we detach ourselves from
		// the original context (keeping only the namespace).
		// This is so that we get a fair run at writing audit entries if Vault
		// has taken up a lot of time handling the request before audit (response)
		// is triggered. Pipeline nodes and the eventlogger.Broker may check for a
		// cancelled context and refuse to process the nodes further.
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("namespace missing from context: %w", err))
			return retErr.ErrorOrNil()
		}

		tempContext, auditCancel := context.WithTimeout(context.Background(), timeout)
		defer auditCancel()
		auditContext = namespace.ContextWithNamespace(tempContext, ns)
	}

	var status eventlogger.Status
	if a.broker.IsAnyPipelineRegistered(event.AuditType.AsEventType()) {
		status, err = a.broker.Send(auditContext, event.AuditType.AsEventType(), e)
		if err != nil {
			retErr = multierror.Append(retErr, multierror.Append(err, status.Warnings...))
			return retErr.ErrorOrNil()
		}
	}

	// Audit event ended up in at least 1 sink.
	if len(status.CompleteSinks()) > 0 {
		return retErr.ErrorOrNil()
	}

	// There were errors from inside the pipeline and we didn't write to a sink.
	if len(status.Warnings) > 0 {
		retErr = multierror.Append(retErr, multierror.Append(errors.New("error during audit pipeline processing"), status.Warnings...))
		return retErr.ErrorOrNil()
	}

	// If a fallback device is registered we can rely on that to 'catch all'
	// and also the broker level guarantee for completed sinks.
	if a.fallbackBroker.IsAnyPipelineRegistered(event.AuditType.AsEventType()) {
		status, err = a.fallbackBroker.Send(auditContext, event.AuditType.AsEventType(), e)
		if err != nil {
			retErr = multierror.Append(retErr, multierror.Append(fmt.Errorf("auditing response to fallback device failed: %w", err), status.Warnings...))
		}
	} else {
		// This audit event won't make it to any devices, we class this as a 'miss' for auditing.
		metrics.IncrCounter(audit.MetricLabelsFallbackMiss(), 1)
	}

	return retErr.ErrorOrNil()
}

func (a *AuditBroker) Invalidate(ctx context.Context, _ string) {
	// For now, we ignore the key as this would only apply to salts.
	// We just sort of brute force it on each one.
	a.Lock()
	defer a.Unlock()

	for _, be := range a.backends {
		be.backend.Invalidate(ctx)
	}
}

// requiredSuccessThresholdSinks examines backends that have already been registered,
// and returns the value that should be used for configuring success threshold sinks
// on the eventlogger broker.
// If all backends have nodes which provide filtering, then we cannot rely on the
// guarantee provided by setting the threshold to 1, and must set it to 0.
// If you are registering an audit device, you should first check if that backend
// does not have filtering before querying the backends via requiredSuccessThresholdSinks.
// backends may also contain a fallback device, which should be ignored as it is
// handled by the fallbackBroker.
func (a *AuditBroker) requiredSuccessThresholdSinks() int {
	threshold := 0

	// We might need to check over all the existing backends to discover if any
	// don't use filtering.
	for _, be := range a.backends {
		switch {
		case be.backend.IsFallback():
			// Ignore fallback devices as they're handled by a separate broker.
			continue
		case !be.backend.HasFiltering():
			threshold = 1
			break
		}
	}

	return threshold
}

// registerNodesAndPipeline registers eventlogger nodes and a pipeline with the
// backend's name, on the specified eventlogger.Broker using the audit.Backend
// to supply them.
func registerNodesAndPipeline(broker *eventlogger.Broker, b audit.Backend) error {
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

// registerFallback can be used to register a fallback device, it will also
// configure the success threshold required for sinks.
func (a *AuditBroker) registerFallback(name string, backend audit.Backend) error {
	err := registerNodesAndPipeline(a.fallbackBroker, backend)
	if err != nil {
		return fmt.Errorf("pipeline registration error for fallback device: %w", err)
	}

	// Store the name of the fallback audit device so that we can check when
	// deregistering if the device is the single fallback one.
	a.fallbackName = backend.Name()

	// We need to turn on the threshold for the fallback broker, so we can guarantee it ends up somewhere
	err = a.fallbackBroker.SetSuccessThresholdSinks(event.AuditType.AsEventType(), 1)
	if err != nil {
		return fmt.Errorf("unable to configure fallback sink success threshold (1): %w", err)
	}

	return nil
}

// deregisterFallback can be used to deregister a fallback audit device, it will
// also configure the success threshold required for sinks.
func (a *AuditBroker) deregisterFallback(ctx context.Context, name string) error {
	err := a.fallbackBroker.SetSuccessThresholdSinks(event.AuditType.AsEventType(), 0)
	if err != nil {
		return fmt.Errorf("unable to reconfigure fallback sink success threshold (0): %w", err)
	}

	_, err = a.fallbackBroker.RemovePipelineAndNodes(ctx, event.AuditType.AsEventType(), eventlogger.PipelineID(name))
	if err != nil {
		return fmt.Errorf("unable to deregister fallback device: %w", err)
	}

	// Clear the fallback device name now we've deregistered.
	a.fallbackName = ""

	return nil
}

// register can be used to register a normal audit device, it will also calculate
// and configure the success threshold required for sinks.
func (a *AuditBroker) register(name string, backend audit.Backend) error {
	err := registerNodesAndPipeline(a.broker, backend)
	if err != nil {
		return fmt.Errorf("audit pipeline registration error: %w", err)
	}

	// Establish if we ONLY have pipelines that include filter nodes.
	// Otherwise, we can rely on the eventlogger broker guarantee.
	// Check the backend we're working with first, then query the backends
	// that are already registered.
	threshold := 0
	if !backend.HasFiltering() {
		threshold = 1
	} else {
		threshold = a.requiredSuccessThresholdSinks()
	}

	// Update the success threshold now that the pipeline is registered.
	err = a.broker.SetSuccessThresholdSinks(event.AuditType.AsEventType(), threshold)
	if err != nil {
		return fmt.Errorf("unable to configure sink success threshold (%d): %w", threshold, err)
	}

	return nil
}

// deregister can be used to deregister a normal audit device, it will also
// calculate and configure the success threshold required for sinks.
func (a *AuditBroker) deregister(ctx context.Context, name string) error {
	// Establish if we ONLY have pipelines that include filter nodes.
	// Otherwise, we can rely on the eventlogger broker guarantee.
	threshold := a.requiredSuccessThresholdSinks()

	err := a.broker.SetSuccessThresholdSinks(event.AuditType.AsEventType(), threshold)
	if err != nil {
		return fmt.Errorf("unable to reconfigure sink success threshold (%d): %w", threshold, err)
	}

	// The first return value, a bool, indicates whether
	// RemovePipelineAndNodes encountered the error while evaluating
	// pre-conditions (false) or once it started removing the pipeline and
	// the nodes (true). This code doesn't care either way.
	_, err = a.broker.RemovePipelineAndNodes(ctx, event.AuditType.AsEventType(), eventlogger.PipelineID(name))
	if err != nil {
		return fmt.Errorf("unable to remove pipeline and nodes: %w", err)
	}

	return nil
}

// hasAuditPipelines can be used as a shorthand to check if a broker has any
// registered pipelines that are for the audit event type.
func hasAuditPipelines(broker *eventlogger.Broker) bool {
	return broker.IsAnyPipelineRegistered(event.AuditType.AsEventType())
}

// isContextViable examines the supplied context to see if its own deadline would
// occur later than a newly created context with a specific timeout.
// If the existing context is viable it can be used 'as-is', if not, the caller
// should consider creating a new context with the relevant deadline and associated
// context values (e.g. namespace) in order to reduce the likelihood that the
// audit system believes there is a failure in audit (and updating its metrics)
// when the root cause is elsewhere.
func isContextViable(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	deadline, hasDeadline := ctx.Deadline()

	// If there's no deadline on the context then we don't need to worry about
	// it being cancelled due to a timeout.
	if !hasDeadline {
		return true
	}

	return deadline.After(time.Now().Add(timeout))
}
