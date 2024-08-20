// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
)

const (
	stdout  = "stdout"
	discard = "discard"

	optionFilePath = "file_path"
	optionMode     = "mode"
)

var _ Backend = (*fileBackend)(nil)

type fileBackend struct {
	*backend
}

// NewFileBackend provides a wrapper to support the expectation elsewhere in Vault that
// all audit backends can be created via a factory that returns an interface (Backend).
func NewFileBackend(conf *BackendConfig, headersConfig HeaderFormatter) (be Backend, err error) {
	be, err = newFileBackend(conf, headersConfig)
	return
}

// newFileBackend creates a backend and configures all nodes including a file sink.
func newFileBackend(conf *BackendConfig, headersConfig HeaderFormatter) (*fileBackend, error) {
	if headersConfig == nil || reflect.ValueOf(headersConfig).IsNil() {
		return nil, fmt.Errorf("nil header formatter: %w", ErrInvalidParameter)
	}
	if conf == nil {
		return nil, fmt.Errorf("nil config: %w", ErrInvalidParameter)
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	// Get file path from config or fall back to the old option ('path') for compatibility
	// (see commit bac4fe0799a372ba1245db642f3f6cd1f1d02669).
	var filePath string
	if p, ok := conf.Config[optionFilePath]; ok {
		filePath = p
	} else if p, ok = conf.Config["path"]; ok {
		filePath = p
	} else {
		return nil, fmt.Errorf("%q is required: %w", optionFilePath, ErrExternalOptions)
	}

	bec, err := newBackend(headersConfig, conf)
	if err != nil {
		return nil, err
	}
	b := &fileBackend{backend: bec}

	// normalize file path if configured for stdout
	if strings.EqualFold(filePath, stdout) {
		filePath = stdout
	}
	if strings.EqualFold(filePath, discard) {
		filePath = discard
	}

	// Configure the sink.
	cfg, err := newFormatterConfig(headersConfig, conf.Config)
	if err != nil {
		return nil, err
	}

	sinkOpts := []event.Option{event.WithLogger(conf.Logger)}
	if mode, ok := conf.Config[optionMode]; ok {
		sinkOpts = append(sinkOpts, event.WithFileMode(mode))
	}

	err = b.configureSinkNode(conf.MountPath, filePath, cfg.requiredFormat, sinkOpts...)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// configureSinkNode is used internally by fileBackend to create and configure the
// sink node on the backend.
func (b *fileBackend) configureSinkNode(name string, filePath string, format format, opt ...event.Option) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required: %w", ErrExternalOptions)
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return fmt.Errorf("file path is required: %w", ErrExternalOptions)
	}

	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("error generating random NodeID for sink node: %w: %w", ErrInternal, err)
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
		sinkNode, err = event.NewStdoutSinkNode(format.String())
	case discard:
		sinkName = discard
		sinkNode = event.NewNoopSink()
	default:
		// The NewFileSink function attempts to open the file and will return an error if it can't.
		sinkName = name
		sinkNode, err = event.NewFileSink(filePath, format.String(), opt...)
	}
	if err != nil {
		return fmt.Errorf("file sink creation failed for path %q: %w", filePath, err)
	}

	// Wrap the sink node with metrics middleware
	err = b.wrapMetrics(sinkName, sinkNodeID, sinkNode)
	if err != nil {
		return err
	}

	return nil
}

// Reload will trigger the reload action on the sink node for this backend.
func (b *fileBackend) Reload() error {
	for _, n := range b.nodeMap {
		if n.Type() == eventlogger.NodeTypeSink {
			return n.Reopen()
		}
	}

	return nil
}
