// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"fmt"
)

const (
	FileSink   SinkType = "file"
	SocketSink SinkType = "socket"
	SyslogSink SinkType = "syslog"
)

// SinkType defines the type of sink
type SinkType string

// Validate ensures that SinkType is one of the set of allowed sink types.
func (t SinkType) Validate() error {
	const op = "event.(SinkType).Validate"
	switch t {
	case FileSink, SocketSink, SyslogSink:
		return nil
	default:
		return fmt.Errorf("%s: '%s' is not a valid sink type: %w", op, t, ErrInvalidParameter)
	}
}
