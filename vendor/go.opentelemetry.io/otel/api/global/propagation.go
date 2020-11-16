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

package global

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global/internal"
)

// TextMapPropagator returns the global TextMapPropagator. If none has been
// set, a No-Op TextMapPropagator is returned.
func TextMapPropagator() otel.TextMapPropagator {
	return internal.TextMapPropagator()
}

// SetTextMapPropagator sets propagator as the global TSetTextMapPropagator.
func SetTextMapPropagator(propagator otel.TextMapPropagator) {
	internal.SetTextMapPropagator(propagator)
}
