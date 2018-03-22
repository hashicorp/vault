package audit

import (
	"context"
	"io"
)

// Formatter is an interface that is responsible for formating a
// request/response into some format. Formatters write their output
// to an io.Writer.
//
// It is recommended that you pass data through Hash prior to formatting it.
type Formatter interface {
	FormatRequest(context.Context, io.Writer, FormatterConfig, *LogInput) error
	FormatResponse(context.Context, io.Writer, FormatterConfig, *LogInput) error
}

type FormatterConfig struct {
	Raw          bool
	HMACAccessor bool

	// This should only ever be used in a testing context
	OmitTime bool
}
