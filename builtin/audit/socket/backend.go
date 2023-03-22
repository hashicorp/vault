// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package socket

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *audit.BackendConfig) (audit.Backend, error) {
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

	format, ok := conf.Config["format"]
	if !ok {
		format = "json"
	}
	switch format {
	case "json", "jsonx":
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

	b := &Backend{
		saltConfig: conf.SaltConfig,
		saltView:   conf.SaltView,
		formatConfig: audit.FormatterConfig{
			Raw:                logRaw,
			HMACAccessor:       hmacAccessor,
			ElideListResponses: elideListResponses,
		},

		writeDuration: writeDuration,
		address:       address,
		socketType:    socketType,
	}

	switch format {
	case "json":
		b.formatter.AuditFormatWriter = &audit.JSONFormatWriter{
			Prefix:   conf.Config["prefix"],
			SaltFunc: b.Salt,
		}
	case "jsonx":
		b.formatter.AuditFormatWriter = &audit.JSONxFormatWriter{
			Prefix:   conf.Config["prefix"],
			SaltFunc: b.Salt,
		}
	}

	return b, nil
}

// Backend is the audit backend for the socket audit transport.
type Backend struct {
	connection net.Conn

	formatter    audit.AuditFormatter
	formatConfig audit.FormatterConfig

	writeDuration time.Duration
	address       string
	socketType    string

	sync.Mutex

	saltMutex  sync.RWMutex
	salt       *salt.Salt
	saltConfig *salt.Config
	saltView   logical.Storage
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
	if err := b.formatter.FormatRequest(ctx, &buf, b.formatConfig, in); err != nil {
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
	if err := b.formatter.FormatResponse(ctx, &buf, b.formatConfig, in); err != nil {
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
	var buf bytes.Buffer
	temporaryFormatter := audit.NewTemporaryFormatter(config["format"], config["prefix"])
	if err := temporaryFormatter.FormatRequest(ctx, &buf, b.formatConfig, in); err != nil {
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
