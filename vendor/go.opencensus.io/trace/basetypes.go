// Copyright 2017, OpenCensus Authors
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

package trace

import (
	"fmt"
	"time"
)

type (
	// TraceID is a 16-byte identifier for a set of spans.
	TraceID [16]byte
	// SpanID is an 8-byte identifier for a single span.
	SpanID [8]byte
)

func (t TraceID) String() string {
	return fmt.Sprintf("%02x", [16]byte(t))
}

func (s SpanID) String() string {
	return fmt.Sprintf("%02x", [8]byte(s))
}

// Annotation represents a text annotation with a set of attributes and a timestamp.
type Annotation struct {
	Time       time.Time
	Message    string
	Attributes map[string]interface{}
}

// Attribute is an interface for attributes;
// it is implemented by BoolAttribute, IntAttribute, and StringAttribute.
type Attribute interface {
	isAttribute()
}

// BoolAttribute represents a bool-valued attribute.
type BoolAttribute struct {
	Key   string
	Value bool
}

func (b BoolAttribute) isAttribute() {}

// Int64Attribute represents an int64-valued attribute.
type Int64Attribute struct {
	Key   string
	Value int64
}

func (i Int64Attribute) isAttribute() {}

// StringAttribute represents a string-valued attribute.
type StringAttribute struct {
	Key   string
	Value string
}

func (s StringAttribute) isAttribute() {}

// LinkType specifies the relationship between the span that had the link
// added, and the linked span.
type LinkType int32

// LinkType values.
const (
	LinkTypeUnspecified LinkType = iota // The relationship of the two spans is unknown.
	LinkTypeChild                       // The current span is a child of the linked span.
	LinkTypeParent                      // The current span is the parent of the linked span.
)

// Link represents a reference from one span to another span.
type Link struct {
	TraceID TraceID
	SpanID  SpanID
	Type    LinkType
	// Attributes is a set of attributes on the link.
	Attributes map[string]interface{}
}

// MessageEventType specifies the type of message event.
type MessageEventType int32

// MessageEventType values.
const (
	MessageEventTypeUnspecified MessageEventType = iota // Unknown event type.
	MessageEventTypeSent                                // Indicates a sent RPC message.
	MessageEventTypeRecv                                // Indicates a received RPC message.
)

// MessageEvent represents an event describing a message sent or received on the network.
type MessageEvent struct {
	Time                 time.Time
	EventType            MessageEventType
	MessageID            int64
	UncompressedByteSize int64
	CompressedByteSize   int64
}

// Status is the status of a Span.
type Status struct {
	// Code is a status code.  Zero indicates success.
	//
	// If Code will be propagated to Google APIs, it ideally should be a value from
	// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto .
	Code    int32
	Message string
}
