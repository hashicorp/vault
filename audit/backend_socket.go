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
	optionAddress      = "address"
	optionSocketType   = "socket_type"
	optionWriteTimeout = "write_timeout"
)

var _ Backend = (*socketBackend)(nil)

type socketBackend struct {
	*backend
}

// NewSocketBackend provides a means to create socket backend audit devices that
// satisfy the Factory pattern expected elsewhere in Vault.
func NewSocketBackend(conf *BackendConfig, headersConfig HeaderFormatter) (be Backend, err error) {
	be, err = newSocketBackend(conf, headersConfig)
	return
}

// newSocketBackend creates a backend and configures all nodes including a socket sink.
func newSocketBackend(conf *BackendConfig, headersConfig HeaderFormatter) (*socketBackend, error) {
	if headersConfig == nil || reflect.ValueOf(headersConfig).IsNil() {
		return nil, fmt.Errorf("nil header formatter: %w", ErrInvalidParameter)
	}
	if conf == nil {
		return nil, fmt.Errorf("nil config: %w", ErrInvalidParameter)
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	bec, err := newBackend(headersConfig, conf)
	if err != nil {
		return nil, err
	}

	address, ok := conf.Config[optionAddress]
	if !ok {
		return nil, fmt.Errorf("%q is required: %w", optionAddress, ErrExternalOptions)
	}
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("%q cannot be empty: %w", optionAddress, ErrExternalOptions)
	}

	socketType, ok := conf.Config[optionSocketType]
	if !ok {
		socketType = "tcp"
	}

	writeDeadline, ok := conf.Config[optionWriteTimeout]
	if !ok {
		writeDeadline = "2s"
	}

	sinkOpts := []event.Option{
		event.WithSocketType(socketType),
		event.WithMaxDuration(writeDeadline),
		event.WithLogger(conf.Logger),
	}

	err = event.ValidateOptions(sinkOpts...)
	if err != nil {
		return nil, err
	}

	b := &socketBackend{backend: bec}

	// Configure the sink.
	cfg, err := newFormatterConfig(headersConfig, conf.Config)
	if err != nil {
		return nil, err
	}

	err = b.configureSinkNode(conf.MountPath, address, cfg.requiredFormat, sinkOpts...)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *socketBackend) configureSinkNode(name string, address string, format format, opts ...event.Option) error {
	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("error generating random NodeID for sink node: %w", err)
	}

	n, err := event.NewSocketSink(address, format.String(), opts...)
	if err != nil {
		return err
	}

	// Wrap the sink node with metrics middleware
	err = b.wrapMetrics(name, sinkNodeID, n)
	if err != nil {
		return err
	}

	return nil
}

// Reload will trigger the reload action on the sink node for this backend.
func (b *socketBackend) Reload() error {
	for _, n := range b.nodeMap {
		if n.Type() == eventlogger.NodeTypeSink {
			return n.Reopen()
		}
	}

	return nil
}
