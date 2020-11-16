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

package otel

import (
	"context"

	"go.opentelemetry.io/otel/internal/baggage"
	"go.opentelemetry.io/otel/label"
)

// Baggage returns a copy of the baggage in ctx.
func Baggage(ctx context.Context) label.Set {
	// TODO (MrAlias, #1222): The underlying storage, the Map, shares many of
	// the functional elements of the label.Set. These should be unified so
	// this conversion is unnecessary and there is no performance hit calling
	// this.
	m := baggage.MapFromContext(ctx)
	values := make([]label.KeyValue, 0, m.Len())
	m.Foreach(func(kv label.KeyValue) bool {
		values = append(values, kv)
		return true
	})
	return label.NewSet(values...)
}

// BaggageValue returns the value related to key in the baggage of ctx. If no
// value is set, the returned label.Value will be an uninitialized zero-value
// with type INVALID.
func BaggageValue(ctx context.Context, key label.Key) label.Value {
	v, _ := baggage.MapFromContext(ctx).Value(key)
	return v
}

// ContextWithBaggageValues returns a copy of parent with pairs updated in the baggage.
func ContextWithBaggageValues(parent context.Context, pairs ...label.KeyValue) context.Context {
	m := baggage.MapFromContext(parent).Apply(baggage.MapUpdate{
		MultiKV: pairs,
	})
	return baggage.ContextWithMap(parent, m)
}

// ContextWithoutBaggageValues returns a copy of parent in which the values related
// to keys have been removed from the baggage.
func ContextWithoutBaggageValues(parent context.Context, keys ...label.Key) context.Context {
	m := baggage.MapFromContext(parent).Apply(baggage.MapUpdate{
		DropMultiK: keys,
	})
	return baggage.ContextWithMap(parent, m)
}

// ContextWithoutBaggage returns a copy of parent without baggage.
func ContextWithoutBaggage(parent context.Context) context.Context {
	return baggage.ContextWithNoCorrelationData(parent)
}
