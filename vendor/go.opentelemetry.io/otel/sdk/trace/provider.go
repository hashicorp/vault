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
	"sync/atomic"

	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/api/trace"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

const (
	defaultTracerName = "go.opentelemetry.io/otel/sdk/tracer"
)

// TODO (MrAlias): unify this API option design:
// https://github.com/open-telemetry/opentelemetry-go/issues/536

// TracerProviderConfig
type TracerProviderConfig struct {
	processors []SpanProcessor
	config     Config
}

type TracerProviderOption func(*TracerProviderConfig)

type TracerProvider struct {
	mu             sync.Mutex
	namedTracer    map[instrumentation.Library]*tracer
	spanProcessors atomic.Value
	config         atomic.Value // access atomically
}

var _ apitrace.TracerProvider = &TracerProvider{}

// NewTracerProvider creates an instance of trace provider. Optional
// parameter configures the provider with common options applicable
// to all tracer instances that will be created by this provider.
func NewTracerProvider(opts ...TracerProviderOption) *TracerProvider {
	o := &TracerProviderConfig{}

	for _, opt := range opts {
		opt(o)
	}

	tp := &TracerProvider{
		namedTracer: make(map[instrumentation.Library]*tracer),
	}
	tp.config.Store(&Config{
		DefaultSampler:       ParentBased(AlwaysSample()),
		IDGenerator:          defIDGenerator(),
		MaxAttributesPerSpan: DefaultMaxAttributesPerSpan,
		MaxEventsPerSpan:     DefaultMaxEventsPerSpan,
		MaxLinksPerSpan:      DefaultMaxLinksPerSpan,
	})

	for _, sp := range o.processors {
		tp.RegisterSpanProcessor(sp)
	}

	tp.ApplyConfig(o.config)

	return tp
}

// Tracer with the given name. If a tracer for the given name does not exist,
// it is created first. If the name is empty, DefaultTracerName is used.
func (p *TracerProvider) Tracer(name string, opts ...apitrace.TracerOption) apitrace.Tracer {
	c := trace.NewTracerConfig(opts...)

	p.mu.Lock()
	defer p.mu.Unlock()
	if name == "" {
		name = defaultTracerName
	}
	il := instrumentation.Library{
		Name:    name,
		Version: c.InstrumentationVersion,
	}
	t, ok := p.namedTracer[il]
	if !ok {
		t = &tracer{
			provider:               p,
			instrumentationLibrary: il,
		}
		p.namedTracer[il] = t
	}
	return t
}

// RegisterSpanProcessor adds the given SpanProcessor to the list of SpanProcessors
func (p *TracerProvider) RegisterSpanProcessor(s SpanProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	new := spanProcessorStates{}
	if old, ok := p.spanProcessors.Load().(spanProcessorStates); ok {
		new = append(new, old...)
	}
	newSpanSync := &spanProcessorState{
		sp:    s,
		state: &sync.Once{},
	}
	new = append(new, newSpanSync)
	p.spanProcessors.Store(new)
}

// UnregisterSpanProcessor removes the given SpanProcessor from the list of SpanProcessors
func (p *TracerProvider) UnregisterSpanProcessor(s SpanProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	new := spanProcessorStates{}
	old, ok := p.spanProcessors.Load().(spanProcessorStates)
	if !ok || len(old) == 0 {
		return
	}
	new = append(new, old...)

	// stop the span processor if it is started and remove it from the list
	var stopOnce *spanProcessorState
	var idx int
	for i, sps := range new {
		if sps.sp == s {
			stopOnce = sps
			idx = i
		}
	}
	if stopOnce != nil {
		stopOnce.state.Do(func() {
			s.Shutdown()
		})
	}
	if len(new) > 1 {
		copy(new[idx:], new[idx+1:])
	}
	new[len(new)-1] = nil
	new = new[:len(new)-1]

	p.spanProcessors.Store(new)
}

// ApplyConfig changes the configuration of the provider.
// If a field in the configuration is empty or nil then its original value is preserved.
func (p *TracerProvider) ApplyConfig(cfg Config) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c := *p.config.Load().(*Config)
	if cfg.DefaultSampler != nil {
		c.DefaultSampler = cfg.DefaultSampler
	}
	if cfg.IDGenerator != nil {
		c.IDGenerator = cfg.IDGenerator
	}
	if cfg.MaxEventsPerSpan > 0 {
		c.MaxEventsPerSpan = cfg.MaxEventsPerSpan
	}
	if cfg.MaxAttributesPerSpan > 0 {
		c.MaxAttributesPerSpan = cfg.MaxAttributesPerSpan
	}
	if cfg.MaxLinksPerSpan > 0 {
		c.MaxLinksPerSpan = cfg.MaxLinksPerSpan
	}
	if cfg.Resource != nil {
		c.Resource = cfg.Resource
	}
	p.config.Store(&c)
}

// WithSyncer registers the exporter with the TracerProvider using a
// SimpleSpanProcessor.
func WithSyncer(e export.SpanExporter) TracerProviderOption {
	return WithSpanProcessor(NewSimpleSpanProcessor(e))
}

// WithBatcher registers the exporter with the TracerProvider using a
// BatchSpanProcessor configured with the passed opts.
func WithBatcher(e export.SpanExporter, opts ...BatchSpanProcessorOption) TracerProviderOption {
	return WithSpanProcessor(NewBatchSpanProcessor(e, opts...))
}

// WithSpanProcessor registers the SpanProcessor with a TracerProvider.
func WithSpanProcessor(sp SpanProcessor) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.processors = append(opts.processors, sp)
	}
}

// WithConfig option sets the configuration to provider.
func WithConfig(config Config) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.config = config
	}
}

// WithResource option attaches a resource to the provider.
// The resource is added to the span when it is started.
func WithResource(r *resource.Resource) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.config.Resource = r
	}
}
