// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"context"

	"go.opencensus.io/trace"
	"google.golang.org/api/googleapi"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/status"
)

// EndSpanFunc is a function that ends a span, reporting an error if necessary.
type EndSpanFunc func(error)

// Attribute annotates a span with additional data.
type Attribute struct {
	key   string
	value interface{}
}

func (a Attribute) traceAttr() trace.Attribute {
	// always use a string attribute for now
	// if need for additional types arise, this can be expanded.
	return trace.StringAttribute(a.key, a.value.(string))
}

// AddInstanceName creates an attribute with the Cloud SQL instance name.
func AddInstanceName(name string) Attribute {
	return Attribute{key: "/cloudsql/instance", value: name}
}

// AddDialerID creates an attribute to identify a particular dialer.
func AddDialerID(dialerID string) Attribute {
	return Attribute{key: "/cloudsql/dialer_id", value: dialerID}
}

// StartSpan begins a span with the provided name and returns a context and a
// function to end the created span.
func StartSpan(ctx context.Context, name string, attrs ...Attribute) (context.Context, EndSpanFunc) {
	var span *trace.Span
	ctx, span = trace.StartSpan(ctx, name)
	as := make([]trace.Attribute, 0, len(attrs))
	for _, a := range attrs {
		as = append(as, a.traceAttr())
	}
	span.AddAttributes(as...)
	return ctx, func(err error) {
		if err != nil {
			span.SetStatus(toStatus(err))
		}
		span.End()
	}
}

// toStatus interrogates an error and converts it to an appropriate
// OpenCensus status.
// Note: this function is borrowed from
// https://github.com/googleapis/google-cloud-go/blob/master/internal/trace/trace.go
func toStatus(err error) trace.Status {
	if err2, ok := err.(*googleapi.Error); ok {
		return trace.Status{Code: httpStatusCodeToOCCode(err2.Code), Message: err2.Message}
	}
	if s, ok := status.FromError(err); ok {
		return trace.Status{Code: int32(s.Code()), Message: s.Message()}
	}
	return trace.Status{Code: int32(code.Code_UNKNOWN), Message: err.Error()}
}

// Reference: https://github.com/googleapis/googleapis/blob/26b634d2724ac5dd30ae0b0cbfb01f07f2e4050e/google/rpc/code.proto
func httpStatusCodeToOCCode(httpStatusCode int) int32 {
	switch httpStatusCode {
	case 200:
		return int32(code.Code_OK)
	case 499:
		return int32(code.Code_CANCELLED)
	case 500:
		return int32(code.Code_UNKNOWN) // Could also be Code_INTERNAL, Code_DATA_LOSS
	case 400:
		return int32(code.Code_INVALID_ARGUMENT) // Could also be Code_OUT_OF_RANGE
	case 504:
		return int32(code.Code_DEADLINE_EXCEEDED)
	case 404:
		return int32(code.Code_NOT_FOUND)
	case 409:
		return int32(code.Code_ALREADY_EXISTS) // Could also be Code_ABORTED
	case 403:
		return int32(code.Code_PERMISSION_DENIED)
	case 401:
		return int32(code.Code_UNAUTHENTICATED)
	case 429:
		return int32(code.Code_RESOURCE_EXHAUSTED)
	case 501:
		return int32(code.Code_UNIMPLEMENTED)
	case 503:
		return int32(code.Code_UNAVAILABLE)
	default:
		return int32(code.Code_UNKNOWN)
	}
}
