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

type noopTracerProvider struct{}

var _ TracerProvider = noopTracerProvider{}

// Tracer returns noop implementation of Tracer.
func (p noopTracerProvider) Tracer(_ string, _ ...TracerOption) Tracer {
	return noopTracer{}
}

// NoopTracerProvider returns a noop implementation of TracerProvider. The
// Tracer and Spans created from the noop provider will also be noop.
func NoopTracerProvider() TracerProvider {
	return noopTracerProvider{}
}
