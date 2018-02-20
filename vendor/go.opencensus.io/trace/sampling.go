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
	"encoding/binary"
)

const defaultSamplingProbability = 1e-4

var defaultSampler Sampler

func init() {
	defaultSampler = newDefaultSampler()
}

func newDefaultSampler() Sampler {
	return ProbabilitySampler(defaultSamplingProbability)
}

// SetDefaultSampler sets the default sampler used when creating new spans.
func SetDefaultSampler(sampler Sampler) {
	if sampler == nil {
		sampler = newDefaultSampler()
	}
	mu.Lock()
	defaultSampler = sampler
	mu.Unlock()
}

// Sampler is an interface for values that have a method that the trace library
// can call to determine whether to export a trace's spans.
type Sampler interface {
	Sample(p SamplingParameters) SamplingDecision
}

// SamplingParameters contains the values passed to a Sampler.
type SamplingParameters struct {
	ParentContext   SpanContext
	TraceID         TraceID
	SpanID          SpanID
	Name            string
	HasRemoteParent bool
}

// SamplingDecision is the value returned by a Sampler.
type SamplingDecision struct {
	Sample bool
}

// ProbabilitySampler returns a Sampler that samples a given fraction of traces.
//
// It also samples spans whose parents are sampled.
func ProbabilitySampler(fraction float64) Sampler {
	if !(fraction >= 0) {
		fraction = 0
	} else if fraction >= 1 {
		return AlwaysSample()
	}
	return probabilitySampler{
		traceIDUpperBound: uint64(fraction * (1 << 63)),
	}
}

type probabilitySampler struct {
	traceIDUpperBound uint64
}

var _ Sampler = (*probabilitySampler)(nil)

func (s probabilitySampler) Sample(p SamplingParameters) (d SamplingDecision) {
	if p.ParentContext.IsSampled() {
		return SamplingDecision{Sample: true}
	}
	x := binary.BigEndian.Uint64(p.TraceID[0:8]) >> 1
	return SamplingDecision{Sample: x < s.traceIDUpperBound}
}

// AlwaysSample returns a Sampler that samples every trace.
func AlwaysSample() Sampler {
	return always{}
}

type always struct{}

var _ Sampler = always{}

func (a always) Sample(p SamplingParameters) SamplingDecision {
	return SamplingDecision{Sample: true}
}

// NeverSample returns a Sampler that samples no traces.
func NeverSample() Sampler {
	return never{}
}

type never struct{}

var _ Sampler = never{}

func (n never) Sample(p SamplingParameters) SamplingDecision {
	return SamplingDecision{Sample: false}
}
