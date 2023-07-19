// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package syslog

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/hashicorp/eventlogger"
	gsyslog "github.com/hashicorp/go-syslog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *audit.BackendConfig, useEventLogger bool) (audit.Backend, error) {
	if conf.SaltConfig == nil {
		return nil, fmt.Errorf("nil salt config")
	}
	if conf.SaltView == nil {
		return nil, fmt.Errorf("nil salt view")
	}

	// Get facility or default to AUTH
	facility, ok := conf.Config["facility"]
	if !ok {
		facility = "AUTH"
	}

	// Get tag or default to 'vault'
	tag, ok := conf.Config["tag"]
	if !ok {
		tag = "vault"
	}

	auditFormat := audit.JSONFormat
	format, ok := conf.Config["format"]
	if !ok {
		format = audit.JSONFormat.String()
	}
	switch format {
	case audit.JSONFormat.String():
	case audit.JSONxFormat.String():
		auditFormat = audit.JSONxFormat
	default:
		return nil, fmt.Errorf("unknown format type %q", format)
	}

	// Check if hashing of accessor is disabled
	hmacAccessor := true
	if hmacAccessorRaw, ok := conf.Config["hmac_accessor"]; ok {
		value, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return nil, err
		}
		hmacAccessor = value
	}

	// Check if raw logging is enabled
	logRaw := false
	if raw, ok := conf.Config["log_raw"]; ok {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		logRaw = b
	}

	elideListResponses := false
	if elideListResponsesRaw, ok := conf.Config["elide_list_responses"]; ok {
		value, err := strconv.ParseBool(elideListResponsesRaw)
		if err != nil {
			return nil, err
		}
		elideListResponses = value
	}

	// Get the logger
	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, facility, tag)
	if err != nil {
		return nil, err
	}

	cfg := audit.FormatterConfig{
		Raw:                logRaw,
		HMACAccessor:       hmacAccessor,
		ElideListResponses: elideListResponses,
		RequiredFormat:     auditFormat,
	}

	b := &Backend{
		logger:       logger,
		saltConfig:   conf.SaltConfig,
		saltView:     conf.SaltView,
		formatConfig: cfg,

		nodeIDList: make([]eventlogger.NodeID, 0),
		nodeMap:    make(map[eventlogger.NodeID]eventlogger.Node),
	}

	// Configure the formatter for either case.
	f, err := audit.NewEventFormatter(b.formatConfig, b)
	if err != nil {
		return nil, fmt.Errorf("error creating formatter: %w", err)
	}

	var w audit.Writer
	switch format {
	case "json":
		w = &audit.JSONWriter{Prefix: conf.Config["prefix"]}
	case "jsonx":
		w = &audit.JSONxWriter{Prefix: conf.Config["prefix"]}
	}

	formatterNodeID := event.GenerateNodeID()
	b.nodeIDList = append(b.nodeIDList, formatterNodeID)
	b.nodeMap[formatterNodeID] = f

	fw, err := audit.NewEventFormatterWriter(b.formatConfig, f, w)
	if err != nil {
		return nil, fmt.Errorf("error creating formatter writer: %w", err)
	}

	b.formatter = fw

	sinkNode, err := event.NewSyslogSink(format, event.WithFacility(facility), event.WithTag(tag))
	if err != nil {
		return nil, fmt.Errorf("error creating syslog sink node: %w", err)
	}

	sinkNodeID := event.GenerateNodeID()
	b.nodeIDList = append(b.nodeIDList, sinkNodeID)
	b.nodeMap[sinkNodeID] = sinkNode

	return b, nil
}

// Backend is the audit backend for the syslog-based audit store.
type Backend struct {
	logger gsyslog.Syslogger

	formatter    *audit.EventFormatterWriter
	formatConfig audit.FormatterConfig

	saltMutex  sync.RWMutex
	salt       *salt.Salt
	saltConfig *salt.Config
	saltView   logical.Storage

	nodeIDList []eventlogger.NodeID
	nodeMap    map[eventlogger.NodeID]eventlogger.Node
}

var _ audit.Backend = (*Backend)(nil)

func (b *Backend) GetHash(ctx context.Context, data string) (string, error) {
	salt, err := b.Salt(ctx)
	if err != nil {
		return "", err
	}
	return audit.HashString(salt, data), nil
}

func (b *Backend) LogRequest(ctx context.Context, in *logical.LogInput) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatAndWriteRequest(ctx, &buf, in); err != nil {
		return err
	}

	// Write out to syslog
	_, err := b.logger.Write(buf.Bytes())
	return err
}

func (b *Backend) LogResponse(ctx context.Context, in *logical.LogInput) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatAndWriteResponse(ctx, &buf, in); err != nil {
		return err
	}

	// Write out to syslog
	_, err := b.logger.Write(buf.Bytes())
	return err
}

func (b *Backend) LogTestMessage(ctx context.Context, in *logical.LogInput, config map[string]string) error {
	var buf bytes.Buffer
	temporaryFormatter := audit.NewTemporaryFormatter(config["format"], config["prefix"])
	if err := temporaryFormatter.FormatAndWriteRequest(ctx, &buf, in); err != nil {
		return err
	}

	// Send to syslog
	_, err := b.logger.Write(buf.Bytes())
	return err
}

func (b *Backend) Reload(_ context.Context) error {
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
	salt, err := salt.NewSalt(ctx, b.saltView, b.saltConfig)
	if err != nil {
		return nil, err
	}
	b.salt = salt
	return salt, nil
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt = nil
}

func (b *Backend) RegisterNodesAndPipeline(broker *eventlogger.Broker, name string) error {
	for id, node := range b.nodeMap {
		if err := broker.RegisterNode(id, node); err != nil {
			return err
		}
	}

	pipeline := eventlogger.Pipeline{
		PipelineID: eventlogger.PipelineID(name),
		EventType:  eventlogger.EventType("audit"),
		NodeIDs:    b.nodeIDList,
	}

	return broker.RegisterPipeline(pipeline)
}
