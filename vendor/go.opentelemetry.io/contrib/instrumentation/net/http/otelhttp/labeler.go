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

package otelhttp

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/label"
)

// Labeler is used to allow instrumented HTTP handlers to add custom labels to
// the metrics recorded by the net/http instrumentation.
type Labeler struct {
	mu     sync.Mutex
	labels []label.KeyValue
}

// Add labels to a Labeler.
func (l *Labeler) Add(ls ...label.KeyValue) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.labels = append(l.labels, ls...)
}

// Labels returns a copy of the labels added to the Labeler.
func (l *Labeler) Get() []label.KeyValue {
	l.mu.Lock()
	defer l.mu.Unlock()
	ret := make([]label.KeyValue, len(l.labels))
	copy(ret, l.labels)
	return ret
}

type labelerContextKeyType int

const lablelerContextKey labelerContextKeyType = 0

func injectLabeler(ctx context.Context, l *Labeler) context.Context {
	return context.WithValue(ctx, lablelerContextKey, l)
}

// LabelerFromContext retrieves a Labeler instance from the provided context if
// one is available.  If no Labeler was found in the provided context a new, empty
// Labeler is returned and the second return value is false.  In this case it is
// safe to use the Labeler but any labels added to it will not be used.
func LabelerFromContext(ctx context.Context) (*Labeler, bool) {
	l, ok := ctx.Value(lablelerContextKey).(*Labeler)
	if !ok {
		l = &Labeler{}
	}
	return l, ok
}
