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

// Int64Measure is a measure of type int64.
type Int64Measure struct {
	name        string
	unit        string
	description string
}

// Name returns the name of the measure.
func (m *Int64Measure) Name() string {
	return m.name
}

// Description returns the description of the measure.
func (m *Int64Measure) Description() string {
	return m.description
}

// Unit returns the unit of the measure.
func (m *Int64Measure) Unit() string {
	return m.unit
}

// M creates a new int64 measurement.
// Use Record to record measurements.
func (m *Int64Measure) M(v int64) Measurement {
	return Measurement{Measure: m, Value: v}
}

// Int64 creates a new measure of type Int64Measure. It returns an
// error if a measure with the same name already exists.
func Int64(name, description, unit string) (*Int64Measure, error) {
	if err := checkName(name); err != nil {
		return nil, err
	}
	m := &Int64Measure{
		name:        name,
		description: description,
		unit:        unit,
	}
	_, err := register(m)
	if err != nil {
		return nil, err
	} else {
		return m, err
	}
}
