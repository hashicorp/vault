// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/vault/internal/observability/event"
)

const (
	optionFacility = "facility"
	optionTag      = "tag"
)

var _ Backend = (*syslogBackend)(nil)

type syslogBackend struct {
	*backend
}

// NewSyslogBackend provides a wrapper to support the expectation elsewhere in Vault that
// all audit backends can be created via a factory that returns an interface (Backend).
func NewSyslogBackend(conf *BackendConfig, headersConfig HeaderFormatter) (be Backend, err error) {
	be, err = newSyslogBackend(conf, headersConfig)
	return
}

// newSyslogBackend creates a backend and configures all nodes including a socket sink.
func newSyslogBackend(conf *BackendConfig, headersConfig HeaderFormatter) (*syslogBackend, error) {
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

	// Get facility or default to AUTH
	facility, ok := conf.Config[optionFacility]
	if !ok {
		facility = "AUTH"
	}

	// Get tag or default to 'vault'
	tag, ok := conf.Config[optionTag]
	if !ok {
		tag = "vault"
	}

	sinkOpts := []event.Option{
		event.WithFacility(facility),
		event.WithTag(tag),
		event.WithLogger(conf.Logger),
	}

	err = event.ValidateOptions(sinkOpts...)
	if err != nil {
		return nil, err
	}

	b := &syslogBackend{backend: bec}

	// Configure the sink.
	cfg, err := newFormatterConfig(headersConfig, conf.Config)
	if err != nil {
		return nil, err
	}

	err = b.configureSinkNode(conf.MountPath, cfg.requiredFormat, sinkOpts...)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *syslogBackend) configureSinkNode(name string, format format, opts ...event.Option) error {
	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return fmt.Errorf("error generating random NodeID for sink node: %w: %w", ErrInternal, err)
	}

	n, err := event.NewSyslogSink(format.String(), opts...)
	if err != nil {
		return fmt.Errorf("error creating syslog sink node: %w", err)
	}

	err = b.wrapMetrics(name, sinkNodeID, n)
	if err != nil {
		return err
	}

	return nil
}

// Reload will trigger the reload action on the sink node for this backend.
func (b *syslogBackend) Reload() error {
	return nil
}
