// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package b3

import (
	"context"
	"errors"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/trace"
)

const (
	// Default B3 Header names.
	b3ContextHeader      = "b3"
	b3DebugFlagHeader    = "x-b3-flags"
	b3TraceIDHeader      = "x-b3-traceid"
	b3SpanIDHeader       = "x-b3-spanid"
	b3SampledHeader      = "x-b3-sampled"
	b3ParentSpanIDHeader = "x-b3-parentspanid"

	b3TraceIDPadding = "0000000000000000"

	// B3 Single Header encoding widths.
	separatorWidth      = 1       // Single "-" character.
	samplingWidth       = 1       // Single hex character.
	traceID64BitsWidth  = 64 / 4  // 16 hex character Trace ID.
	traceID128BitsWidth = 128 / 4 // 32 hex character Trace ID.
	spanIDWidth         = 16      // 16 hex character ID.
	parentSpanIDWidth   = 16      // 16 hex character ID.
)

var (
	empty = trace.EmptySpanContext()

	errInvalidSampledByte        = errors.New("invalid B3 Sampled found")
	errInvalidSampledHeader      = errors.New("invalid B3 Sampled header found")
	errInvalidTraceIDHeader      = errors.New("invalid B3 traceID header found")
	errInvalidSpanIDHeader       = errors.New("invalid B3 spanID header found")
	errInvalidParentSpanIDHeader = errors.New("invalid B3 ParentSpanID header found")
	errInvalidScope              = errors.New("require either both traceID and spanID or none")
	errInvalidScopeParent        = errors.New("ParentSpanID requires both traceID and spanID to be available")
	errInvalidScopeParentSingle  = errors.New("ParentSpanID requires traceID, spanID and Sampled to be available")
	errEmptyContext              = errors.New("empty request context")
	errInvalidTraceIDValue       = errors.New("invalid B3 traceID value found")
	errInvalidSpanIDValue        = errors.New("invalid B3 spanID value found")
	errInvalidParentSpanIDValue  = errors.New("invalid B3 ParentSpanID value found")
)

// Encoding is a bitmask representation of the B3 encoding type.
type Encoding uint8

// supports returns if e has o bit(s) set.
func (e Encoding) supports(o Encoding) bool {
	return e&o == o
}

const (
	// B3MultipleHeader is a B3 encoding that uses multiple headers to
	// transmit tracing information all prefixed with `x-b3-`.
	B3MultipleHeader Encoding = 1 << iota
	// B3SingleHeader is a B3 encoding that uses a single header named `b3`
	// to transmit tracing information.
	B3SingleHeader
	// B3Unspecified is an unspecified B3 encoding.
	B3Unspecified Encoding = 0
)

// B3 propagator serializes SpanContext to/from B3 Headers.
// This propagator supports both versions of B3 headers,
//  1. Single Header:
//    b3: {TraceId}-{SpanId}-{SamplingState}-{ParentSpanId}
//  2. Multiple Headers:
//    x-b3-traceid: {TraceId}
//    x-b3-parentspanid: {ParentSpanId}
//    x-b3-spanid: {SpanId}
//    x-b3-sampled: {SamplingState}
//    x-b3-flags: {DebugFlag}
type B3 struct {
	// InjectEncoding are the B3 encodings used when injecting trace
	// information. If no encoding is specified (i.e. `B3Unspecified`)
	// `B3MultipleHeader` will be used as the default.
	InjectEncoding Encoding
}

var _ otel.TextMapPropagator = B3{}

// Inject injects a context into the carrier as B3 headers.
// The parent span ID is omitted because it is not tracked in the
// SpanContext.
func (b3 B3) Inject(ctx context.Context, carrier otel.TextMapCarrier) {
	sc := trace.SpanFromContext(ctx).SpanContext()

	if b3.InjectEncoding.supports(B3SingleHeader) {
		header := []string{}
		if sc.TraceID.IsValid() && sc.SpanID.IsValid() {
			header = append(header, sc.TraceID.String(), sc.SpanID.String())
		}

		if sc.TraceFlags&trace.FlagsDebug == trace.FlagsDebug {
			header = append(header, "d")
		} else if !(sc.TraceFlags&trace.FlagsDeferred == trace.FlagsDeferred) {
			if sc.IsSampled() {
				header = append(header, "1")
			} else {
				header = append(header, "0")
			}
		}

		carrier.Set(b3ContextHeader, strings.Join(header, "-"))
	}

	if b3.InjectEncoding.supports(B3MultipleHeader) || b3.InjectEncoding == B3Unspecified {
		if sc.TraceID.IsValid() && sc.SpanID.IsValid() {
			carrier.Set(b3TraceIDHeader, sc.TraceID.String())
			carrier.Set(b3SpanIDHeader, sc.SpanID.String())
		}

		if sc.TraceFlags&trace.FlagsDebug == trace.FlagsDebug {
			// Since Debug implies deferred, don't also send "X-B3-Sampled".
			carrier.Set(b3DebugFlagHeader, "1")
		} else if !(sc.TraceFlags&trace.FlagsDeferred == trace.FlagsDeferred) {
			if sc.IsSampled() {
				carrier.Set(b3SampledHeader, "1")
			} else {
				carrier.Set(b3SampledHeader, "0")
			}
		}
	}
}

// Extract extracts a context from the carrier if it contains B3 headers.
func (b3 B3) Extract(ctx context.Context, carrier otel.TextMapCarrier) context.Context {
	var (
		sc  trace.SpanContext
		err error
	)

	// Default to Single Header if a valid value exists.
	if h := carrier.Get(b3ContextHeader); h != "" {
		sc, err = extractSingle(h)
		if err == nil && sc.IsValid() {
			return trace.ContextWithRemoteSpanContext(ctx, sc)
		}
		// The Single Header value was invalid, fallback to Multiple Header.
	}

	var (
		traceID      = carrier.Get(b3TraceIDHeader)
		spanID       = carrier.Get(b3SpanIDHeader)
		parentSpanID = carrier.Get(b3ParentSpanIDHeader)
		sampled      = carrier.Get(b3SampledHeader)
		debugFlag    = carrier.Get(b3DebugFlagHeader)
	)
	sc, err = extractMultiple(traceID, spanID, parentSpanID, sampled, debugFlag)
	if err != nil || !sc.IsValid() {
		return ctx
	}
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (b3 B3) Fields() []string {
	header := []string{}
	if b3.InjectEncoding.supports(B3SingleHeader) {
		header = append(header, b3ContextHeader)
	}
	if b3.InjectEncoding.supports(B3MultipleHeader) || b3.InjectEncoding == B3Unspecified {
		header = append(header, b3TraceIDHeader, b3SpanIDHeader, b3SampledHeader, b3DebugFlagHeader)
	}
	return header
}

// extractMultiple reconstructs a SpanContext from header values based on B3
// Multiple header. It is based on the implementation found here:
// https://github.com/openzipkin/zipkin-go/blob/v0.2.2/propagation/b3/spancontext.go
// and adapted to support a SpanContext.
func extractMultiple(traceID, spanID, parentSpanID, sampled, flags string) (trace.SpanContext, error) {
	var (
		err           error
		requiredCount int
		sc            = trace.SpanContext{}
	)

	// correct values for an existing sampled header are "0" and "1".
	// For legacy support and  being lenient to other tracing implementations we
	// allow "true" and "false" as inputs for interop purposes.
	switch strings.ToLower(sampled) {
	case "0", "false":
		// Zero value for TraceFlags sample bit is unset.
	case "1", "true":
		sc.TraceFlags = trace.FlagsSampled
	case "":
		sc.TraceFlags = trace.FlagsDeferred
	default:
		return empty, errInvalidSampledHeader
	}

	// The only accepted value for Flags is "1". This will set Debug bitmask and
	// sampled bitmask to 1 since debug implicitly means sampled. All other
	// values and omission of header will be ignored. According to the spec. User
	// shouldn't send X-B3-Sampled header along with X-B3-Flags header. Thus we will
	// ignore X-B3-Sampled header when X-B3-Flags header is sent and valid.
	if flags == "1" {
		sc.TraceFlags |= trace.FlagsDebug | trace.FlagsSampled
		sc.TraceFlags &= ^trace.FlagsDeferred
	}

	if traceID != "" {
		requiredCount++
		id := traceID
		if len(traceID) == 16 {
			// Pad 64-bit trace IDs.
			id = b3TraceIDPadding + traceID
		}
		if sc.TraceID, err = trace.IDFromHex(id); err != nil {
			return empty, errInvalidTraceIDHeader
		}
	}

	if spanID != "" {
		requiredCount++
		if sc.SpanID, err = trace.SpanIDFromHex(spanID); err != nil {
			return empty, errInvalidSpanIDHeader
		}
	}

	if requiredCount != 0 && requiredCount != 2 {
		return empty, errInvalidScope
	}

	if parentSpanID != "" {
		if requiredCount == 0 {
			return empty, errInvalidScopeParent
		}
		// Validate parent span ID but we do not use it so do not save it.
		if _, err = trace.SpanIDFromHex(parentSpanID); err != nil {
			return empty, errInvalidParentSpanIDHeader
		}
	}

	return sc, nil
}

// extractSingle reconstructs a SpanContext from contextHeader based on a B3
// Single header. It is based on the implementation found here:
// https://github.com/openzipkin/zipkin-go/blob/v0.2.2/propagation/b3/spancontext.go
// and adapted to support a SpanContext.
func extractSingle(contextHeader string) (trace.SpanContext, error) {
	if contextHeader == "" {
		return empty, errEmptyContext
	}

	var (
		sc       = trace.SpanContext{}
		sampling string
	)

	headerLen := len(contextHeader)

	if headerLen == samplingWidth {
		sampling = contextHeader
	} else if headerLen == traceID64BitsWidth || headerLen == traceID128BitsWidth {
		// Trace ID by itself is invalid.
		return empty, errInvalidScope
	} else if headerLen >= traceID64BitsWidth+spanIDWidth+separatorWidth {
		pos := 0
		var traceID string
		if string(contextHeader[traceID64BitsWidth]) == "-" {
			// traceID must be 64 bits
			pos += traceID64BitsWidth // {traceID}
			traceID = b3TraceIDPadding + string(contextHeader[0:pos])
		} else if string(contextHeader[32]) == "-" {
			// traceID must be 128 bits
			pos += traceID128BitsWidth // {traceID}
			traceID = string(contextHeader[0:pos])
		} else {
			return empty, errInvalidTraceIDValue
		}
		var err error
		sc.TraceID, err = trace.IDFromHex(traceID)
		if err != nil {
			return empty, errInvalidTraceIDValue
		}
		pos += separatorWidth // {traceID}-

		sc.SpanID, err = trace.SpanIDFromHex(contextHeader[pos : pos+spanIDWidth])
		if err != nil {
			return empty, errInvalidSpanIDValue
		}
		pos += spanIDWidth // {traceID}-{spanID}

		if headerLen > pos {
			if headerLen == pos+separatorWidth {
				// {traceID}-{spanID}- is invalid.
				return empty, errInvalidSampledByte
			}
			pos += separatorWidth // {traceID}-{spanID}-

			if headerLen == pos+samplingWidth {
				sampling = string(contextHeader[pos])
			} else if headerLen == pos+parentSpanIDWidth {
				// {traceID}-{spanID}-{parentSpanID} is invalid.
				return empty, errInvalidScopeParentSingle
			} else if headerLen == pos+samplingWidth+separatorWidth+parentSpanIDWidth {
				sampling = string(contextHeader[pos])
				pos += samplingWidth + separatorWidth // {traceID}-{spanID}-{sampling}-

				// Validate parent span ID but we do not use it so do not
				// save it.
				_, err = trace.SpanIDFromHex(contextHeader[pos:])
				if err != nil {
					return empty, errInvalidParentSpanIDValue
				}
			} else {
				return empty, errInvalidParentSpanIDValue
			}
		}
	} else {
		return empty, errInvalidTraceIDValue
	}
	switch sampling {
	case "":
		sc.TraceFlags = trace.FlagsDeferred
	case "d":
		sc.TraceFlags = trace.FlagsDebug | trace.FlagsSampled
	case "1":
		sc.TraceFlags = trace.FlagsSampled
	case "0":
		// Zero value for TraceFlags sample bit is unset.
	default:
		return empty, errInvalidSampledByte
	}

	return sc, nil
}
