// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudevents

import (
	"fmt"

	"github.com/hashicorp/eventlogger"
)

const (
	DataContentTypeCloudEvents = "application/cloudevents" // the media type for FormatJSON cloudevent data
	DataContentTypeText        = "text/plain"              // the media type for FormatText cloudevent data
)

// Format defines a type for the supported encoding formats used by
// Formatter.Format and used when when calling eventlogger Event.Format(...)
// from other nodes
type Format string

var (
	FormatJSON        Format = "cloudevents-json" // JSON encoding which is accessible via Event.Format(...) in other nodes (like FileSinks)
	FormatText        Format = "cloudevents-text" // Text encoding which is accessible via Event.Format(...) in other nodes (like FileSinks)
	FormatUnspecified Format = ""                 // Unspecified format which defaults to FormatJSON
)

func (f Format) validate() error {
	const op = "cloudevents.(Format).validate"
	switch f {
	case FormatJSON, FormatText, FormatUnspecified:
		return nil
	default:
		return fmt.Errorf("%s: '%s' is not a valid format: %w", op, f, eventlogger.ErrInvalidParameter)
	}
}

func (f Format) convertToDataContentType() string {
	switch f {
	case FormatJSON, FormatUnspecified:
		return DataContentTypeCloudEvents
	default:
		return DataContentTypeText
	}
}
