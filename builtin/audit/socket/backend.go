// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package socket

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *audit.BackendConfig, useEventLogger bool, headersConfig audit.HeaderFormatter) (audit.Backend, error) {
	if conf.SaltConfig == nil {
		return nil, fmt.Errorf("nil salt config")
	}
	if conf.SaltView == nil {
		return nil, fmt.Errorf("nil salt view")
	}

	address, ok := conf.Config["address"]
	if !ok {
		return nil, fmt.Errorf("address is required")
	}

	socketType, ok := conf.Config["socket_type"]
	if !ok {
		socketType = "tcp"
	}
	writeDeadline, ok := conf.Config["write_timeout"]
	if !ok {
		writeDeadline = "2s"
	}
	writeDuration, err := parseutil.ParseDurationSecond(writeDeadline)
	if err != nil {
		return nil, err
	}

	var cfgOpts []audit.Option

	if format, ok := conf.Config["format"]; ok {
		cfgOpts = append(cfgOpts, audit.WithFormat(format))
	}

	// Check if hashing of accessor is disabled
	if hmacAccessorRaw, ok := conf.Config["hmac_accessor"]; ok {
		v, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return nil, err
		}
		cfgOpts = append(cfgOpts, audit.WithHMACAccessor(v))
	}

	// Check if raw logging is enabled
	if raw, ok := conf.Config["log_raw"]; ok {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		cfgOpts = append(cfgOpts, audit.WithRaw(v))
	}

	if elideListResponsesRaw, ok := conf.Config["elide_list_responses"]; ok {
		v, err := strconv.ParseBool(elideListResponsesRaw)
		if err != nil {
			return nil, err
		}
		cfgOpts = append(cfgOpts, audit.WithElision(v))
	}

	cfg, err := audit.NewFormatterConfig(cfgOpts...)
	if err != nil {
		return nil, err
	}

	b := &Backend{
		saltConfig:   conf.SaltConfig,
		saltView:     conf.SaltView,
		formatConfig: cfg,

		writeDuration: writeDuration,
		address:       address,
		socketType:    socketType,
	}

	// Configure the formatter for either case.
	f, err := audit.NewEntryFormatter(b.formatConfig, b, audit.WithHeaderFormatter(headersConfig))
	if err != nil {
		return nil, fmt.Errorf("error creating formatter: %w", err)
	}
	var w audit.Writer
	switch b.formatConfig.RequiredFormat {
	case audit.JSONFormat:
		w = &audit.JSONWriter{Prefix: conf.Config["prefix"]}
	case audit.JSONxFormat:
		w = &audit.JSONxWriter{Prefix: conf.Config["prefix"]}
	}

	fw, err := audit.NewEntryFormatterWriter(b.formatConfig, f, w)
	if err != nil {
		return nil, fmt.Errorf("error creating formatter writer: %w", err)
	}

	b.formatter = fw

	if useEventLogger {
		var opts []event.Option

		if socketType, ok := conf.Config["socket_type"]; ok {
			opts = append(opts, event.WithSocketType(socketType))
		}

		if writeDeadline, ok := conf.Config["write_timeout"]; ok {
			opts = append(opts, event.WithMaxDuration(writeDeadline))
		}

		b.nodeIDList = make([]eventlogger.NodeID, 2)
		b.nodeMap = make(map[eventlogger.NodeID]eventlogger.Node)

		formatterNodeID, err := event.GenerateNodeID()
		if err != nil {
			return nil, fmt.Errorf("error generating random NodeID for formatter node: %w", err)
		}
		b.nodeIDList[0] = formatterNodeID
		b.nodeMap[formatterNodeID] = f

		n, err := event.NewSocketSink(b.formatConfig.RequiredFormat.String(), address, opts...)
		if err != nil {
			return nil, fmt.Errorf("error creating socket sink node: %w", err)
		}
		sinkNode := &audit.SinkWrapper{Name: conf.MountPath, Sink: n}
		sinkNodeID, err := event.GenerateNodeID()
		if err != nil {
			return nil, fmt.Errorf("error generating random NodeID for sink node: %w", err)
		}
		b.nodeIDList[1] = sinkNodeID
		b.nodeMap[sinkNodeID] = sinkNode
	}

	return b, nil
}

// Backend is the audit backend for the socket audit transport.
type Backend struct {
	connection net.Conn

	formatter    *audit.EntryFormatterWriter
	formatConfig audit.FormatterConfig

	writeDuration time.Duration
	address       string
	socketType    string

	sync.Mutex

	saltMutex  sync.RWMutex
	salt       *salt.Salt
	saltConfig *salt.Config
	saltView   logical.Storage

	nodeIDList []eventlogger.NodeID
	nodeMap    map[eventlogger.NodeID]eventlogger.Node
}

var _ audit.Backend = (*Backend)(nil)

func (b *Backend) LogRequest(ctx context.Context, in *logical.LogInput) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatAndWriteRequest(ctx, &buf, in); err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()

	err := b.write(ctx, buf.Bytes())
	if err != nil {
		rErr := b.reconnect(ctx)
		if rErr != nil {
			err = multierror.Append(err, rErr)
		} else {
			// Try once more after reconnecting
			err = b.write(ctx, buf.Bytes())
		}
	}

	return err
}

func (b *Backend) LogResponse(ctx context.Context, in *logical.LogInput) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatAndWriteResponse(ctx, &buf, in); err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()

	err := b.write(ctx, buf.Bytes())
	if err != nil {
		rErr := b.reconnect(ctx)
		if rErr != nil {
			err = multierror.Append(err, rErr)
		} else {
			// Try once more after reconnecting
			err = b.write(ctx, buf.Bytes())
		}
	}

	return err
}

func (b *Backend) LogTestMessage(ctx context.Context, in *logical.LogInput, config map[string]string) error {
	// Event logger behavior - manually Process each node
	if len(b.nodeIDList) > 0 {
		return audit.ProcessManual(ctx, in, b.nodeIDList, b.nodeMap)
	}

	// Old behavior
	var buf bytes.Buffer

	temporaryFormatter, err := audit.NewTemporaryFormatter(config["format"], config["prefix"])
	if err != nil {
		return err
	}

	if err = temporaryFormatter.FormatAndWriteRequest(ctx, &buf, in); err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()

	err = b.write(ctx, buf.Bytes())
	if err != nil {
		rErr := b.reconnect(ctx)
		if rErr != nil {
			err = multierror.Append(err, rErr)
		} else {
			// Try once more after reconnecting
			err = b.write(ctx, buf.Bytes())
		}
	}

	return err
}

func (b *Backend) write(ctx context.Context, buf []byte) error {
	if b.connection == nil {
		if err := b.reconnect(ctx); err != nil {
			return err
		}
	}

	err := b.connection.SetWriteDeadline(time.Now().Add(b.writeDuration))
	if err != nil {
		return err
	}

	_, err = b.connection.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) reconnect(ctx context.Context) error {
	if b.connection != nil {
		b.connection.Close()
		b.connection = nil
	}

	timeoutContext, cancel := context.WithTimeout(ctx, b.writeDuration)
	defer cancel()

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(timeoutContext, b.socketType, b.address)
	if err != nil {
		return err
	}

	b.connection = conn

	return nil
}

func (b *Backend) Reload(ctx context.Context) error {
	b.Lock()
	defer b.Unlock()

	err := b.reconnect(ctx)

	return err
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

// RegisterNodesAndPipeline registers the nodes and a pipeline as required by
// the audit.Backend interface.
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
