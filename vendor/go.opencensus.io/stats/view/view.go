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

package view

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"sync/atomic"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/internal"
	"go.opencensus.io/tag"
)

// View allows users to filter and aggregate the recorded events.
// Each view has to be registered to enable data retrieval. Use New to
// initiate new views. Unregister views once you don't want to collect any more
// events.
type View struct {
	name        string // name of View. Must be unique.
	description string

	// tagKeys to perform the aggregation on.
	tagKeys []tag.Key

	// Examples of measures are cpu:tickCount, diskio:time...
	m stats.Measure

	subscribed uint32 // 1 if someone is subscribed and data need to be exported, use atomic to access

	collector *collector
}

// New creates a new view with the given name and description.
// View names need to be unique globally in the entire system.
//
// Data collection will only filter measurements recorded by the given keys.
// Collected data will be processed by the given aggregation algorithm.
//
// Views need to be subscribed toin order to retrieve collection data.
//
// Once the view is no longer required, the view can be unregistered.
func New(name, description string, keys []tag.Key, measure stats.Measure, agg Aggregation) (*View, error) {
	if err := checkViewName(name); err != nil {
		return nil, err
	}
	var ks []tag.Key
	if len(keys) > 0 {
		ks = make([]tag.Key, len(keys))
		copy(ks, keys)
		sort.Slice(ks, func(i, j int) bool { return ks[i].Name() < ks[j].Name() })
	}
	return &View{
		name:        name,
		description: description,
		tagKeys:     ks,
		m:           measure,
		collector:   &collector{make(map[string]AggregationData), agg},
	}, nil
}

// Name returns the name of the view.
func (v *View) Name() string {
	return v.name
}

// Description returns the name of the view.
func (v *View) Description() string {
	return v.description
}

func (v *View) subscribe() {
	atomic.StoreUint32(&v.subscribed, 1)
}

func (v *View) unsubscribe() {
	atomic.StoreUint32(&v.subscribed, 0)
}

// isSubscribed returns true if the view is exporting
// data by subscription.
func (v *View) isSubscribed() bool {
	return atomic.LoadUint32(&v.subscribed) == 1
}

func (v *View) clearRows() {
	v.collector.clearRows()
}

// TagKeys returns the list of tag keys associated with this view.
func (v *View) TagKeys() []tag.Key {
	return v.tagKeys
}

// Aggregation returns the data aggregation method used to aggregate
// the measurements collected by this view.
func (v *View) Aggregation() Aggregation {
	return v.collector.a
}

// Measure returns the measure the view is collecting measurements for.
func (v *View) Measure() stats.Measure {
	return v.m
}

func (v *View) collectedRows(now time.Time) []*Row {
	return v.collector.collectedRows(v.tagKeys, now)
}

func (v *View) addSample(m *tag.Map, val interface{}, now time.Time) {
	if !v.isSubscribed() {
		return
	}
	sig := string(encodeWithKeys(m, v.tagKeys))
	v.collector.addSample(sig, val, now)
}

// A Data is a set of rows about usage of the single measure associated
// with the given view. Each row is specific to a unique set of tags.
type Data struct {
	View       *View
	Start, End time.Time
	Rows       []*Row
}

// Row is the collected value for a specific set of key value pairs a.k.a tags.
type Row struct {
	Tags []tag.Tag
	Data AggregationData
}

func (r *Row) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{ ")
	buffer.WriteString("{ ")
	for _, t := range r.Tags {
		buffer.WriteString(fmt.Sprintf("{%v %v}", t.Key.Name(), t.Value))
	}
	buffer.WriteString(" }")
	buffer.WriteString(fmt.Sprintf("%v", r.Data))
	buffer.WriteString(" }")
	return buffer.String()
}

// Equal returns true if both Rows are equal. Tags are expected to be ordered
// by the key name. Even both rows have the same tags but the tags appear in
// different orders it will return false.
func (r *Row) Equal(other *Row) bool {
	if r == other {
		return true
	}
	return reflect.DeepEqual(r.Tags, other.Tags) && r.Data.equal(other.Data)
}

func checkViewName(name string) error {
	if len(name) > internal.MaxNameLength {
		return fmt.Errorf("view name cannot be larger than %v", internal.MaxNameLength)
	}
	if !internal.IsPrintable(name) {
		return fmt.Errorf("view name needs to be an ASCII string")
	}
	return nil
}
