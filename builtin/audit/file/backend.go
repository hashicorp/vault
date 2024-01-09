// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/eventlogger"
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
	f            *os.File
	fallback     bool
	fileLock     sync.RWMutex
	formatter    *audit.EntryFormatterWriter
	formatConfig audit.FormatterConfig
	mode         os.FileMode
	name         string
	nodeIDList   []eventlogger.NodeID
	nodeMap      map[eventlogger.NodeID]eventlogger.Node
	filePath     string
	salt         *atomic.Value
	saltConfig   *salt.Config
	saltMutex    sync.RWMutex
	saltView     logical.Storage
}

func Factory(_ context.Context, conf *audit.BackendConfig, useEventLogger bool, headersConfig audit.HeaderFormatter) (audit.Backend, error) {
	const op = "file.Factory"

	if conf.SaltConfig == nil {
		return nil, fmt.Errorf("%s: nil salt config", op)
	}
	if conf.SaltView == nil {
		return nil, fmt.Errorf("%s: nil salt view", op)
	}

	// The config options 'fallback' and 'filter' are mutually exclusive, a fallback
	// device catches everything, so it cannot be allowed to filter.
	var fallback bool
	var err error
	if fallbackRaw, ok := conf.Config["fallback"]; ok {
		fallback, err = parseutil.ParseBool(fallbackRaw)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse 'fallback': %w", op, err)
		}
	}

	if _, ok := conf.Config["filter"]; ok && fallback {
		return nil, fmt.Errorf("%s: cannot configure a fallback device with a filter: %w", op, event.ErrInvalidParameter)
	}

	// Get file path from config or fall back to the old option name ('path') for compatibility
	// (see commit bac4fe0799a372ba1245db642f3f6cd1f1d02669).
	var filePath string
	if p, ok := conf.Config["file_path"]; ok {
		filePath = p
	} else if p, ok = conf.Config["path"]; ok {
		filePath = p
	} else {
		return nil, fmt.Errorf("%s: file_path is required", op)
	}

	// normalize file path if configured for stdout
	if strings.EqualFold(filePath, stdout) {
		filePath = stdout
	}
	if strings.EqualFold(filePath, discard) {
		filePath = discard
	}

	mode := os.FileMode(0o600)
	if modeRaw, ok := conf.Config["mode"]; ok {
		m, err := strconv.ParseUint(modeRaw, 8, 32)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse 'mode': %w", op, err)
		}
		switch m {
		case 0:
			// if mode is 0000, then do not modify file mode
			if filePath != stdout && filePath != discard {
				fileInfo, err := os.Stat(filePath)
				if err != nil {
					return nil, fmt.Errorf("%s: unable to stat %q: %w", op, filePath, err)
				}
				mode = fileInfo.Mode()
			}
		default:
			mode = os.FileMode(m)
		}
	}

	cfg, err := formatterConfig(conf.Config)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create formatter config: %w", op, err)
	}

	b := &Backend{
		fallback:     fallback,
		filePath:     filePath,
		formatConfig: cfg,
		mode:         mode,
		name:         conf.MountPath,
		saltConfig:   conf.SaltConfig,
		saltView:     conf.SaltView,
		salt:         new(atomic.Value),
	}

	// Ensure we are working with the right type by explicitly storing a nil of
	// the right type
	b.salt.Store((*salt.Salt)(nil))

	// Configure the formatter for either case.
	f, err := audit.NewEntryFormatter(b.formatConfig, b, audit.WithHeaderFormatter(headersConfig), audit.WithPrefix(conf.Config["prefix"]))
	if err != nil {
		return nil, fmt.Errorf("%s: error creating formatter: %w", op, err)
	}

	var w audit.Writer
	switch b.formatConfig.RequiredFormat {
	case audit.JSONFormat:
		w = &audit.JSONWriter{Prefix: conf.Config["prefix"]}
	case audit.JSONxFormat:
		w = &audit.JSONxWriter{Prefix: conf.Config["prefix"]}
	default:
		return nil, fmt.Errorf("%s: unknown format type %q", op, b.formatConfig.RequiredFormat)
	}

	fw, err := audit.NewEntryFormatterWriter(b.formatConfig, f, w)
	if err != nil {
		return nil, fmt.Errorf("%s: error creating formatter writer: %w", op, err)
	}
	b.formatter = fw

	if useEventLogger {
		b.nodeIDList = []eventlogger.NodeID{}
		b.nodeMap = make(map[eventlogger.NodeID]eventlogger.Node)

		err := b.configureFilterNode(conf.Config["filter"])
		if err != nil {
			return nil, fmt.Errorf("%s: error configuring filter node: %w", op, err)
		}

		formatterOpts := []audit.Option{
			audit.WithHeaderFormatter(headersConfig),
			audit.WithPrefix(conf.Config["prefix"]),
		}

		err = b.configureFormatterNode(cfg, formatterOpts...)
		if err != nil {
			return nil, fmt.Errorf("%s: error configuring formatter node: %w", op, err)
		}

		err = b.configureSinkNode(conf.MountPath, filePath, conf.Config["mode"], cfg.RequiredFormat.String())
		if err != nil {
			return nil, fmt.Errorf("%s: error configuring sink node: %w", op, err)
		}
	} else {
		switch filePath {
		case stdout:
		case discard:
		default:
			// Ensure that the file can be successfully opened for writing;
			// otherwise it will be too late to catch later without problems
			// (ref: https://github.com/hashicorp/vault/issues/550)
			if err := b.open(); err != nil {
				return nil, fmt.Errorf("%s: sanity check failed; unable to open %q for writing: %w", op, filePath, err)
			}
		}
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

// Deprecated: Use eventlogger.
func (b *Backend) LogRequest(ctx context.Context, in *logical.LogInput) error {
	var writer io.Writer
	switch b.filePath {
	case stdout:
		writer = os.Stdout
	case discard:
		return nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, 2000))
	err := b.formatter.FormatAndWriteRequest(ctx, buf, in)
	if err != nil {
		return err
	}

	return b.log(ctx, buf, writer)
}

// Deprecated: Use eventlogger.
func (b *Backend) log(_ context.Context, buf *bytes.Buffer, writer io.Writer) error {
	reader := bytes.NewReader(buf.Bytes())

	b.fileLock.Lock()

	if writer == nil {
		if err := b.open(); err != nil {
			b.fileLock.Unlock()
			return err
		}
		writer = b.f
	}

	if _, err := reader.WriteTo(writer); err == nil {
		b.fileLock.Unlock()
		return nil
	} else if b.filePath == stdout {
		b.fileLock.Unlock()
		return err
	}

	// If writing to stdout there's no real reason to think anything would have
	// changed so return above. Otherwise, opportunistically try to re-open the
	// FD, once per call.
	b.f.Close()
	b.f = nil

	if err := b.open(); err != nil {
		b.fileLock.Unlock()
		return err
	}

	reader.Seek(0, io.SeekStart)
	_, err := reader.WriteTo(writer)
	b.fileLock.Unlock()
	return err
}

// Deprecated: Use eventlogger.
func (b *Backend) LogResponse(ctx context.Context, in *logical.LogInput) error {
	var writer io.Writer
	switch b.filePath {
	case stdout:
		writer = os.Stdout
	case discard:
		return nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, 6000))
	err := b.formatter.FormatAndWriteResponse(ctx, buf, in)
	if err != nil {
		return err
	}

	return b.log(ctx, buf, writer)
}

func (b *Backend) LogTestMessage(ctx context.Context, in *logical.LogInput, config map[string]string) error {
	// Event logger behavior - manually Process each node
	if len(b.nodeIDList) > 0 {
		return audit.ProcessManual(ctx, in, b.nodeIDList, b.nodeMap)
	}

	// Old behavior
	var writer io.Writer
	switch b.filePath {
	case stdout:
		writer = os.Stdout
	case discard:
		return nil
	}

	var buf bytes.Buffer

	temporaryFormatter, err := audit.NewTemporaryFormatter(config["format"], config["prefix"])
	if err != nil {
		return err
	}

	if err = temporaryFormatter.FormatAndWriteRequest(ctx, &buf, in); err != nil {
		return err
	}

	return b.log(ctx, &buf, writer)
}

// The file lock must be held before calling this
// Deprecated: Use eventlogger.
func (b *Backend) open() error {
	if b.f != nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(b.filePath), b.mode); err != nil {
		return err
	}

	var err error
	b.f, err = os.OpenFile(b.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, b.mode)
	if err != nil {
		return err
	}

	// Change the file mode in case the log file already existed. We special
	// case /dev/null since we can't chmod it and bypass if the mode is zero
	switch b.filePath {
	case "/dev/null":
	default:
		if b.mode != 0 {
			err = os.Chmod(b.filePath, b.mode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Backend) Reload(_ context.Context) error {
	// When there are nodes created in the map, use the eventlogger behavior.
	if len(b.nodeMap) > 0 {
		for _, n := range b.nodeMap {
			if n.Type() == eventlogger.NodeTypeSink {
				return n.Reopen()
			}
		}

		return nil
	} else {
		// old non-eventlogger behavior
		switch b.filePath {
		case stdout, discard:
			return nil
		}

		b.fileLock.Lock()
		defer b.fileLock.Unlock()

		if b.f == nil {
			return b.open()
		}

		err := b.f.Close()
		// Set to nil here so that even if we error out, on the next access open()
		// will be tried
		b.f = nil
		if err != nil {
			return err
		}

		return b.open()
	}
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt.Store((*salt.Salt)(nil))
}

// formatterConfig creates the configuration required by a formatter node using
// the config map supplied to the factory.
func formatterConfig(config map[string]string) (audit.FormatterConfig, error) {
	const op = "file.formatterConfig"

	var opts []audit.Option

	if format, ok := config["format"]; ok {
		opts = append(opts, audit.WithFormat(format))
	}

	// Check if hashing of accessor is disabled
	if hmacAccessorRaw, ok := config["hmac_accessor"]; ok {
		v, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return audit.FormatterConfig{}, fmt.Errorf("%s: unable to parse 'hmac_accessor': %w", op, err)
		}
		opts = append(opts, audit.WithHMACAccessor(v))
	}

	// Check if raw logging is enabled
	if raw, ok := config["log_raw"]; ok {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return audit.FormatterConfig{}, fmt.Errorf("%s: unable to parse 'log_raw': %w", op, err)
		}
		opts = append(opts, audit.WithRaw(v))
	}

	if elideListResponsesRaw, ok := config["elide_list_responses"]; ok {
		v, err := strconv.ParseBool(elideListResponsesRaw)
		if err != nil {
			return audit.FormatterConfig{}, fmt.Errorf("%s: unable to parse 'elide_list_responses': %w", op, err)
		}
		opts = append(opts, audit.WithElision(v))
	}

	return audit.NewFormatterConfig(opts...)
}

// configureFilterNode is used to configure a filter node and associated ID on the Backend.
func (b *Backend) configureFilterNode(filter string) error {
	const op = "file.(Backend).configureFilterNode"

	filter = strings.TrimSpace(filter)
	if filter == "" {
		return nil
	}

	filterNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("%s: error generating random NodeID for filter node: %w", op, err)
	}

	filterNode, err := audit.NewEntryFilter(filter)
	if err != nil {
		return fmt.Errorf("%s: error creating filter node: %w", op, err)
	}

	b.nodeIDList = append(b.nodeIDList, filterNodeID)
	b.nodeMap[filterNodeID] = filterNode
	return nil
}

// configureFormatterNode is used to configure a formatter node and associated ID on the Backend.
func (b *Backend) configureFormatterNode(formatConfig audit.FormatterConfig, opts ...audit.Option) error {
	const op = "file.(Backend).configureFormatterNode"

	formatterNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("%s: error generating random NodeID for formatter node: %w", op, err)
	}

	formatterNode, err := audit.NewEntryFormatter(formatConfig, b, opts...)
	if err != nil {
		return fmt.Errorf("%s: error creating formatter: %w", op, err)
	}

	b.nodeIDList = append(b.nodeIDList, formatterNodeID)
	b.nodeMap[formatterNodeID] = formatterNode
	return nil
}

// configureSinkNode is used to configure a sink node and associated ID on the Backend.
func (b *Backend) configureSinkNode(name string, filePath string, mode string, format string) error {
	const op = "file.(Backend).configureSinkNode"

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("%s: name is required: %w", op, event.ErrInvalidParameter)
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return fmt.Errorf("%s: file path is required: %w", op, event.ErrInvalidParameter)
	}

	format = strings.TrimSpace(format)
	if format == "" {
		return fmt.Errorf("%s: format is required: %w", op, event.ErrInvalidParameter)
	}

	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("%s: error generating random NodeID for sink node: %w", op, err)
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
		return fmt.Errorf("%s: file sink creation failed for path %q: %w", op, filePath, err)
	}

	sinkNode = &audit.SinkWrapper{Name: sinkName, Sink: sinkNode}

	b.nodeIDList = append(b.nodeIDList, sinkNodeID)
	b.nodeMap[sinkNodeID] = sinkNode
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
	return len(b.nodeIDList) > 0 && b.nodeMap[b.nodeIDList[0]].Type() == eventlogger.NodeTypeFilter
}

// IsFallback can be used to determine if this audit backend device is intended to
// be used as a fallback to catch all events that are not written when only using
// filtered pipelines.
func (b *Backend) IsFallback() bool {
	return b.fallback
}
