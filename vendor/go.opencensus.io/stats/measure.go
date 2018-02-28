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
//

package stats

import (
	"errors"
	"fmt"
	"sync"

	"go.opencensus.io/stats/internal"
)

// Measure represents a type of metric to be tracked and recorded.
// For example, latency, request Mb/s, and response Mb/s are measures
// to collect from a server.
//
// Each measure needs to be registered before being used.
// Measure constructors such as Int64 and
// Float64 automatically registers the measure
// by the given name.
// Each registered measure needs to be unique by name.
// Measures also have a description and a unit.
type Measure interface {
	Name() string
	Description() string
	Unit() string
}

var (
	mu           sync.RWMutex
	measures     = make(map[string]Measure)
	errDuplicate = errors.New("duplicate measure name")
)

func FindMeasure(name string) Measure {
	mu.RLock()
	defer mu.RUnlock()
	if m, ok := measures[name]; ok {
		return m
	}
	return nil
}

func register(m Measure) (Measure, error) {
	key := m.Name()
	mu.Lock()
	defer mu.Unlock()
	if stored, ok := measures[key]; ok {
		return stored, errDuplicate
	}
	measures[key] = m
	return m, nil
}

// Measurement is the numeric value measured when recording stats. Each measure
// provides methods to create measurements of their kind. For example, Int64Measure
// provides M to convert an int64 into a measurement.
type Measurement struct {
	Value   interface{} // int64 or float64
	Measure Measure
}

func checkName(name string) error {
	if len(name) > internal.MaxNameLength {
		return fmt.Errorf("measure name cannot be larger than %v", internal.MaxNameLength)
	}
	if !internal.IsPrintable(name) {
		return fmt.Errorf("measure name needs to be an ASCII string")
	}
	return nil
}
