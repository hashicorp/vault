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

package trace

import (
	"sync"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

// SpanProcessor is interface to add hooks to start and end method invocations.
type SpanProcessor interface {

	// OnStart method is invoked when span is started. It is a synchronous call
	// and hence should not block.
	OnStart(sd *export.SpanData)

	// OnEnd method is invoked when span is finished. It is a synchronous call
	// and hence should not block.
	OnEnd(sd *export.SpanData)

	// Shutdown is invoked when SDK shutsdown. Use this call to cleanup any processor
	// data. No calls to OnStart and OnEnd method is invoked after Shutdown call is
	// made. It should not be blocked indefinitely.
	Shutdown()

	// ForceFlush exports all ended spans to the configured Exporter that have not yet
	// been exported.  It should only be called when absolutely necessary, such as when
	// using a FaaS provider that may suspend the process after an invocation, but before
	// the Processor can export the completed spans.
	ForceFlush()
}

type spanProcessorState struct {
	sp    SpanProcessor
	state *sync.Once
}
type spanProcessorStates []*spanProcessorState
