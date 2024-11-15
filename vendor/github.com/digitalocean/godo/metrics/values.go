// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// A SampleValue is a representation of a value for a given sample at a given time.
type SampleValue float64

// UnmarshalJSON implements json.Unmarshaler.
func (v *SampleValue) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("sample value must be a quoted string")
	}
	f, err := strconv.ParseFloat(string(b[1:len(b)-1]), 64)
	if err != nil {
		return err
	}
	*v = SampleValue(f)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (v SampleValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v SampleValue) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// Equal returns true if the value of v and o is equal or if both are NaN. Note
// that v==o is false if both are NaN. If you want the conventional float
// behavior, use == to compare two SampleValues.
func (v SampleValue) Equal(o SampleValue) bool {
	if v == o {
		return true
	}
	return math.IsNaN(float64(v)) && math.IsNaN(float64(o))
}

// SamplePair pairs a SampleValue with a Timestamp.
type SamplePair struct {
	Timestamp Time
	Value     SampleValue
}

func (s SamplePair) String() string {
	return fmt.Sprintf("%s @[%s]", s.Value, s.Timestamp)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SamplePair) UnmarshalJSON(b []byte) error {
	v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
	return json.Unmarshal(b, &v)
}

// MarshalJSON implements json.Marshaler.
func (s SamplePair) MarshalJSON() ([]byte, error) {
	t, err := json.Marshal(s.Timestamp)
	if err != nil {
		return nil, err
	}
	v, err := json.Marshal(s.Value)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
}

// SampleStream is a stream of Values belonging to an attached COWMetric.
type SampleStream struct {
	Metric Metric       `json:"metric"`
	Values []SamplePair `json:"values"`
}

func (ss SampleStream) String() string {
	vals := make([]string, len(ss.Values))
	for i, v := range ss.Values {
		vals[i] = v.String()
	}
	return fmt.Sprintf("%s =>\n%s", ss.Metric, strings.Join(vals, "\n"))
}
