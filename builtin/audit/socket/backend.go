// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package socket

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

var _ audit.Backend = (*Backend)(nil)

// Backend is the audit backend for the socket audit transport.
type Backend struct {
	fallback   bool
	name       string
	nodeIDList []eventlogger.NodeID
	nodeMap    map[eventlogger.NodeID]eventlogger.Node
	salt       *salt.Salt
	saltConfig *salt.Config
	saltMutex  sync.RWMutex
	saltView   logical.Storage
}

func Factory(_ context.Context, conf *audit.BackendConfig, headersConfig audit.HeaderFormatter) (audit.Backend, *audit.Error) {
	const op = "socket.Factory"

	if conf.SaltConfig == nil {
		return nil, audit.NewError(op, "nil salt config", audit.ErrInvalidParameter)
	}

	if conf.SaltView == nil {
		return nil, audit.NewError(op, "nil salt view", audit.ErrInvalidParameter)
	}

	if conf.Logger == nil || reflect.ValueOf(conf.Logger).IsNil() {
		return nil, audit.NewError(op, "nil logger", audit.ErrInvalidParameter)
	}

	address, ok := conf.Config["address"]
	if !ok {
		return nil, audit.NewError(op, "address is required", audit.ErrInvalidParameter)
	}

	socketType, ok := conf.Config["socket_type"]
	if !ok {
		socketType = "tcp"
	}

	writeDeadline, ok := conf.Config["write_timeout"]
	if !ok {
		writeDeadline = "2s"
	}

	// The config options 'fallback' and 'filter' are mutually exclusive, a fallback
	// device catches everything, so it cannot be allowed to filter.
	var fallback bool
	var err error
	if fallbackRaw, ok := conf.Config["fallback"]; ok {
		fallback, err = parseutil.ParseBool(fallbackRaw)
		if err != nil {
			return nil, audit.NewError(op, "unable to parse 'fallback", audit.ErrInvalidParameter).Detail(err)
		}
	}

	if _, ok := conf.Config["filter"]; ok && fallback {
		return nil, audit.NewError(op, "cannot configure a fallback device with a filter", audit.ErrInvalidParameter)
	}

	b := &Backend{
		fallback:   fallback,
		name:       conf.MountPath,
		saltConfig: conf.SaltConfig,
		saltView:   conf.SaltView,
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    make(map[eventlogger.NodeID]eventlogger.Node),
	}

	err = b.configureFilterNode(conf.Config["filter"])
	if err != nil {
		return nil, audit.NewError(op, "error configuring filter node", audit.ErrFilterParameter).Detail(err)
	}

	cfg, cfgErr := newFormatterConfig(headersConfig, conf.Config)
	if cfgErr != nil {
		return nil, audit.NewError(op, "failed to create formatter config", audit.ErrConfiguration).Wrap(cfgErr)
	}

	fmtNodeErr := b.configureFormatterNode(conf.MountPath, cfg, conf.Logger)
	if fmtNodeErr != nil {
		return nil, audit.NewError(op, "error configuring formatter node", audit.ErrConfiguration).Wrap(fmtNodeErr)
	}

	sinkOpts := []event.Option{
		event.WithSocketType(socketType),
		event.WithMaxDuration(writeDeadline),
	}

	err = b.configureSinkNode(conf.MountPath, address, cfg.RequiredFormat.String(), sinkOpts...)
	if err != nil {
		return nil, audit.NewError(op, "error configuring sink node", err)
	}

	return b, nil
}

func (b *Backend) LogTestMessage(ctx context.Context, in *logical.LogInput) error {
	if len(b.nodeIDList) > 0 {
		return audit.ProcessManual(ctx, in, b.nodeIDList, b.nodeMap)
	}

	return nil
}

func (b *Backend) Reload(ctx context.Context) error {
	for _, n := range b.nodeMap {
		if n.Type() == eventlogger.NodeTypeSink {
			return n.Reopen()
		}
	}

	return nil
}

func (b *Backend) Salt(ctx context.Context) (*salt.Salt, error) {
	b.saltMutex.RLock()
	if b.salt != nil {
		defer b.saltMutex.RUnlock()
		return b.salt, nil
	}
	b.saltMutex.RUnlock()
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	if b.salt != nil {
		return b.salt, nil
	}
	s, err := salt.NewSalt(ctx, b.saltView, b.saltConfig)
	if err != nil {
		return nil, err
	}
	b.salt = s
	return s, nil
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt = nil
}

// newFormatterConfig creates the configuration required by a formatter node using
// the config map supplied to the factory.
func newFormatterConfig(headerFormatter audit.HeaderFormatter, config map[string]string) (audit.FormatterConfig, *audit.Error) {
	const op = "socket.newFormatterConfig"

	var opts []audit.Option

	if format, ok := config["format"]; ok {
		opts = append(opts, audit.WithFormat(format))
	}

	// Check if hashing of accessor is disabled
	if hmacAccessorRaw, ok := config["hmac_accessor"]; ok {
		v, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return audit.FormatterConfig{}, audit.NewError(op, "unable to parse 'hmac_accessor'", audit.ErrInvalidParameter).Detail(err)
		}
		opts = append(opts, audit.WithHMACAccessor(v))
	}

	// Check if raw logging is enabled
	if raw, ok := config["log_raw"]; ok {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return audit.FormatterConfig{}, audit.NewError(op, "unable to parse 'log_raw'", audit.ErrInvalidParameter).Detail(err)
		}
		opts = append(opts, audit.WithRaw(v))
	}

	if elideListResponsesRaw, ok := config["elide_list_responses"]; ok {
		v, err := strconv.ParseBool(elideListResponsesRaw)
		if err != nil {
			return audit.FormatterConfig{}, audit.NewError(op, "unable to parse 'elide_list_responses'", audit.ErrInvalidParameter).Detail(err)
		}
		opts = append(opts, audit.WithElision(v))
	}

	return audit.NewFormatterConfig(headerFormatter, opts...)
}

// configureFormatterNode is used to configure a formatter node and associated ID on the Backend.
func (b *Backend) configureFormatterNode(name string, formatConfig audit.FormatterConfig, logger hclog.Logger) *audit.Error {
	const op = "socket.(Backend).configureFormatterNode"

	formatterNodeID, err := event.GenerateNodeID()
	if err != nil {
		return audit.NewError(op, "error generating random NodeID for formatter node", audit.ErrUnknown).Detail(err)
	}

	formatterNode, entryErr := audit.NewEntryFormatter(name, formatConfig, b, logger)
	if entryErr != nil {
		return audit.NewError(op, "error creating formatter", audit.ErrConfiguration).Wrap(entryErr)
	}

	b.nodeIDList = append(b.nodeIDList, formatterNodeID)
	b.nodeMap[formatterNodeID] = formatterNode

	return nil
}

// configureSinkNode is used to configure a sink node and associated ID on the Backend.
func (b *Backend) configureSinkNode(name string, address string, format string, opts ...event.Option) error {
	const op = "socket.(Backend).configureSinkNode"

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("%s: name is required: %w", op, audit.ErrInvalidParameter)
	}

	address = strings.TrimSpace(address)
	if address == "" {
		return fmt.Errorf("%s: address is required: %w", op, audit.ErrInvalidParameter)
	}

	format = strings.TrimSpace(format)
	if format == "" {
		return fmt.Errorf("%s: format is required: %w", op, audit.ErrInvalidParameter)
	}

	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("%s: error generating random NodeID for sink node: %w", op, err)
	}

	n, err := event.NewSocketSink(address, format, opts...)
	if err != nil {
		return fmt.Errorf("%s: error creating socket sink node: %w", op, err)
	}

	// Wrap the sink node with metrics middleware
	sinkMetricTimer, err := audit.NewSinkMetricTimer(name, n)
	if err != nil {
		return fmt.Errorf("%s: unable to add timing metrics to sink for path %q: %w", op, name, err)
	}

	// Decide what kind of labels we want and wrap the sink node inside a metrics counter.
	var metricLabeler event.Labeler
	switch {
	case b.fallback:
		metricLabeler = &audit.MetricLabelerAuditFallback{}
	default:
		metricLabeler = &audit.MetricLabelerAuditSink{}
	}

	sinkMetricCounter, err := event.NewMetricsCounter(name, sinkMetricTimer, metricLabeler)
	if err != nil {
		return fmt.Errorf("%s: unable to add counting metrics to sink for path %q: %w", op, name, err)
	}

	b.nodeIDList = append(b.nodeIDList, sinkNodeID)
	b.nodeMap[sinkNodeID] = sinkMetricCounter

	return nil
}

// Name for this backend, this would ideally correspond to the mount path for the audit device.
func (b *Backend) Name() string {
	return b.name
}

// Nodes returns the nodes which should be used by the event framework to process audit entries.
func (b *Backend) Nodes() map[eventlogger.NodeID]eventlogger.Node {
	return b.nodeMap
}

// NodeIDs returns the IDs of the nodes, in the order they are required.
func (b *Backend) NodeIDs() []eventlogger.NodeID {
	return b.nodeIDList
}

// EventType returns the event type for the backend.
func (b *Backend) EventType() eventlogger.EventType {
	return eventlogger.EventType(event.AuditType.String())
}

// HasFiltering determines if the first node for the pipeline is an eventlogger.NodeTypeFilter.
func (b *Backend) HasFiltering() bool {
	if b.nodeMap == nil {
		return false
	}

	return len(b.nodeIDList) > 0 && b.nodeMap[b.nodeIDList[0]].Type() == eventlogger.NodeTypeFilter
}

// IsFallback can be used to determine if this audit backend device is intended to
// be used as a fallback to catch all events that are not written when only using
// filtered pipelines.
func (b *Backend) IsFallback() bool {
	return b.fallback
}
