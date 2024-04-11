// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package file

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	stdout  = "stdout"
	discard = "discard"
)

var _ audit.Backend = (*Backend)(nil)

// Backend is the audit backend for the file-based audit store.
//
// NOTE: This audit backend is currently very simple: it appends to a file.
// It doesn't do anything more at the moment to assist with rotation
// or reset the write cursor, this should be done in the future.
type Backend struct {
	fallback   bool
	name       string
	nodeIDList []eventlogger.NodeID
	nodeMap    map[eventlogger.NodeID]eventlogger.Node
	salt       *atomic.Value
	saltConfig *salt.Config
	saltMutex  sync.RWMutex
	saltView   logical.Storage
}

func Factory(_ context.Context, conf *audit.BackendConfig, headersConfig audit.HeaderFormatter) (audit.Backend, error) {
	if conf.SaltConfig == nil {
		return nil, fmt.Errorf("nil salt config: %w", audit.ErrInvalidParameter)
	}

	if conf.SaltView == nil {
		return nil, fmt.Errorf("nil salt view: %w", audit.ErrInvalidParameter)
	}

	if conf.Logger == nil || reflect.ValueOf(conf.Logger).IsNil() {
		return nil, fmt.Errorf("nil logger: %w", audit.ErrInvalidParameter)
	}

	if conf.MountPath == "" {
		return nil, fmt.Errorf("mount path cannot be empty: %w", audit.ErrInvalidParameter)
	}

	// The config options 'fallback' and 'filter' are mutually exclusive, a fallback
	// device catches everything, so it cannot be allowed to filter.
	var fallback bool
	var err error
	if fallbackRaw, ok := conf.Config["fallback"]; ok {
		fallback, err = parseutil.ParseBool(fallbackRaw)
		if err != nil {
			return nil, fmt.Errorf("unable to parse 'fallback': %w", audit.ErrExternalOptions)
		}
	}

	if _, ok := conf.Config["filter"]; ok && fallback {
		return nil, fmt.Errorf("cannot configure a fallback device with a filter: %w", audit.ErrExternalOptions)
	}

	// Get file path from config or fall back to the old option name ('path') for compatibility
	// (see commit bac4fe0799a372ba1245db642f3f6cd1f1d02669).
	var filePath string
	if p, ok := conf.Config["file_path"]; ok {
		filePath = p
	} else if p, ok = conf.Config["path"]; ok {
		filePath = p
	} else {
		return nil, fmt.Errorf("file_path is required: %w", audit.ErrExternalOptions)
	}

	// normalize file path if configured for stdout
	if strings.EqualFold(filePath, stdout) {
		filePath = stdout
	}
	if strings.EqualFold(filePath, discard) {
		filePath = discard
	}

	cfg, err := newFormatterConfig(headersConfig, conf.Config)
	if err != nil {
		return nil, err
	}

	b := &Backend{
		fallback:   fallback,
		name:       conf.MountPath,
		saltConfig: conf.SaltConfig,
		saltView:   conf.SaltView,
		salt:       new(atomic.Value),
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    make(map[eventlogger.NodeID]eventlogger.Node),
	}

	// Ensure we are working with the right type by explicitly storing a nil of
	// the right type
	b.salt.Store((*salt.Salt)(nil))

	err = b.configureFilterNode(conf.Config["filter"])
	if err != nil {
		return nil, err
	}

	err = b.configureFormatterNode(conf.MountPath, cfg, conf.Logger)
	if err != nil {
		return nil, err
	}

	err = b.configureSinkNode(conf.MountPath, filePath, conf.Config["mode"], cfg.RequiredFormat.String())
	if err != nil {
		return nil, fmt.Errorf("error configuring sink node: %w", err)
	}

	return b, nil
}

func (b *Backend) Salt(ctx context.Context) (*salt.Salt, error) {
	s := b.salt.Load().(*salt.Salt)
	if s != nil {
		return s, nil
	}

	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()

	s = b.salt.Load().(*salt.Salt)
	if s != nil {
		return s, nil
	}

	newSalt, err := salt.NewSalt(ctx, b.saltView, b.saltConfig)
	if err != nil {
		b.salt.Store((*salt.Salt)(nil))
		return nil, err
	}

	b.salt.Store(newSalt)
	return newSalt, nil
}

func (b *Backend) LogTestMessage(ctx context.Context, in *logical.LogInput) error {
	if len(b.nodeIDList) > 0 {
		return audit.ProcessManual(ctx, in, b.nodeIDList, b.nodeMap)
	}

	return nil
}

func (b *Backend) Reload(_ context.Context) error {
	for _, n := range b.nodeMap {
		if n.Type() == eventlogger.NodeTypeSink {
			return n.Reopen()
		}
	}

	return nil
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt.Store((*salt.Salt)(nil))
}

// newFormatterConfig creates the configuration required by a formatter node using
// the config map supplied to the factory.
func newFormatterConfig(headerFormatter audit.HeaderFormatter, config map[string]string) (audit.FormatterConfig, error) {
	var opts []audit.Option

	if format, ok := config["format"]; ok {
		if !audit.IsValidFormat(format) {
			return audit.FormatterConfig{}, fmt.Errorf("unsupported 'format': %w", audit.ErrExternalOptions)
		}

		opts = append(opts, audit.WithFormat(format))
	}

	// Check if hashing of accessor is disabled
	if hmacAccessorRaw, ok := config["hmac_accessor"]; ok {
		v, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return audit.FormatterConfig{}, fmt.Errorf("unable to parse 'hmac_accessor': %w", audit.ErrExternalOptions)
		}
		opts = append(opts, audit.WithHMACAccessor(v))
	}

	// Check if raw logging is enabled
	if raw, ok := config["log_raw"]; ok {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return audit.FormatterConfig{}, fmt.Errorf("unable to parse 'log_raw: %w", audit.ErrExternalOptions)
		}
		opts = append(opts, audit.WithRaw(v))
	}

	if elideListResponsesRaw, ok := config["elide_list_responses"]; ok {
		v, err := strconv.ParseBool(elideListResponsesRaw)
		if err != nil {
			return audit.FormatterConfig{}, fmt.Errorf("unable to parse 'elide_list_responses': %w", audit.ErrExternalOptions)
		}
		opts = append(opts, audit.WithElision(v))
	}

	if prefix, ok := config["prefix"]; ok {
		opts = append(opts, audit.WithPrefix(prefix))
	}

	return audit.NewFormatterConfig(headerFormatter, opts...)
}

// configureFormatterNode is used to configure a formatter node and associated ID on the Backend.
func (b *Backend) configureFormatterNode(name string, formatConfig audit.FormatterConfig, logger hclog.Logger) error {
	formatterNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("error generating random NodeID for formatter node: %w: %w", audit.ErrInternal, err)
	}

	formatterNode, err := audit.NewEntryFormatter(name, formatConfig, b, logger)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	b.nodeIDList = append(b.nodeIDList, formatterNodeID)
	b.nodeMap[formatterNodeID] = formatterNode

	return nil
}

// configureSinkNode is used to configure a sink node and associated ID on the Backend.
func (b *Backend) configureSinkNode(name string, filePath string, mode string, format string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required: %w", audit.ErrExternalOptions)
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return fmt.Errorf("file path is required: %w", audit.ErrExternalOptions)
	}

	format = strings.TrimSpace(format)
	if format == "" {
		return fmt.Errorf("format is required: %w", audit.ErrInvalidParameter)
	}

	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("error generating random NodeID for sink node: %w: %w", audit.ErrInternal, err)
	}

	// normalize file path if configured for stdout or discard
	if strings.EqualFold(filePath, stdout) {
		filePath = stdout
	} else if strings.EqualFold(filePath, discard) {
		filePath = discard
	}

	var sinkNode eventlogger.Node
	var sinkName string

	switch filePath {
	case stdout:
		sinkName = stdout
		sinkNode, err = event.NewStdoutSinkNode(format)
	case discard:
		sinkName = discard
		sinkNode = event.NewNoopSink()
	default:
		// The NewFileSink function attempts to open the file and will return an error if it can't.
		sinkName = name
		sinkNode, err = event.NewFileSink(filePath, format, []event.Option{event.WithFileMode(mode)}...)
	}

	if err != nil {
		return fmt.Errorf("file sink creation failed for path %q: %w", filePath, err)
	}

	// Wrap the sink node with metrics middleware
	sinkMetricTimer, err := audit.NewSinkMetricTimer(sinkName, sinkNode)
	if err != nil {
		return fmt.Errorf("unable to add timing metrics to sink for path %q: %w", filePath, err)
	}

	// Decide what kind of labels we want and wrap the sink node inside a metrics counter.
	var metricLabeler event.Labeler
	switch {
	case b.fallback:
		metricLabeler = &audit.MetricLabelerAuditFallback{}
	default:
		metricLabeler = &audit.MetricLabelerAuditSink{}
	}

	sinkMetricCounter, err := event.NewMetricsCounter(sinkName, sinkMetricTimer, metricLabeler)
	if err != nil {
		return fmt.Errorf("unable to add counting metrics to sink for path %q: %w", filePath, err)
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
	return event.AuditType.AsEventType()
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
