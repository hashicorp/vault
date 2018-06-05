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
	"math"
)

// AggregationData represents an aggregated value from a collection.
// They are reported on the view data during exporting.
// Mosts users won't directly access aggregration data.
type AggregationData interface {
	isAggregationData() bool
	addSample(v interface{})
	addOther(other AggregationData)
	multiplyByFraction(fraction float64) AggregationData
	clear()
	clone() AggregationData
	equal(other AggregationData) bool
}

const epsilon = 1e-9

// CountData is the aggregated data for a CountAggregation.
// A count aggregation processes data and counts the recordings.
//
// Most users won't directly access count data.
type CountData int64

func newCountData(v int64) *CountData {
	tmp := CountData(v)
	return &tmp
}

func (a *CountData) isAggregationData() bool { return true }

func (a *CountData) addSample(v interface{}) {
	*a = *a + 1
}

func (a *CountData) clone() AggregationData {
	return newCountData(int64(*a))
}

func (a *CountData) multiplyByFraction(fraction float64) AggregationData {
	return newCountData(int64(float64(int64(*a))*fraction + 0.5)) // adding 0.5 because go runtime will take floor instead of rounding
}

func (a *CountData) addOther(av AggregationData) {
	other, ok := av.(*CountData)
	if !ok {
		return
	}
	*a = *a + *other
}

func (a *CountData) clear() {
	*a = 0
}

func (a *CountData) equal(other AggregationData) bool {
	a2, ok := other.(*CountData)
	if !ok {
		return false
	}

	return int64(*a) == int64(*a2)
}

// SumData is the aggregated data for a SumAggregation.
// A sum aggregation processes data and sums up the recordings.
//
// Most users won't directly access sum data.
type SumData float64

func newSumData(v float64) *SumData {
	tmp := SumData(v)
	return &tmp
}

func (a *SumData) isAggregationData() bool { return true }

func (a *SumData) addSample(v interface{}) {
	// Both float64 and int64 values will be cast to float64
	var f float64
	switch x := v.(type) {
	case int64:
		f = float64(x)
	case float64:
		f = x
	default:
		return
	}
	*a += SumData(f)
}

func (a *SumData) multiplyByFraction(fraction float64) AggregationData {
	return newSumData(float64(*a) * fraction)
}

func (a *SumData) clone() AggregationData {
	return newSumData(float64(*a))
}

func (a *SumData) addOther(av AggregationData) {
	other, ok := av.(*SumData)
	if !ok {
		return
	}
	*a = *a + *other
}

func (a *SumData) clear() {
	*a = 0
}

func (a *SumData) equal(other AggregationData) bool {
	a2, ok := other.(*SumData)
	if !ok {
		return false
	}
	return math.Pow(float64(*a)-float64(*a2), 2) < epsilon
}

// MeanData is the aggregated data for a MeanAggregation.
// A mean aggregation processes data and maintains the mean value.
//
// Most users won't directly access mean data.
type MeanData struct {
	Count float64 // number of data points aggregated
	Mean  float64 // mean of all data points
}

func newMeanData(mean float64, count float64) *MeanData {
	return &MeanData{
		Mean:  mean,
		Count: count,
	}
}

// Sum returns the sum of all samples collected.
func (a *MeanData) Sum() float64 { return a.Mean * float64(a.Count) }

func (a *MeanData) isAggregationData() bool { return true }

func (a *MeanData) addSample(v interface{}) {
	var f float64
	switch x := v.(type) {
	case int64:
		f = float64(x)
	case float64:
		f = x
	default:
		return
	}

	a.Count++
	if a.Count == 1 {
		a.Mean = f
		return
	}
	a.Mean = a.Mean + (f-a.Mean)/float64(a.Count)
}

func (a *MeanData) clone() AggregationData {
	return newMeanData(a.Mean, a.Count)
}

// Only Count will be mutiplied by the fraction, Mean will remain the same.
func (a *MeanData) multiplyByFraction(fraction float64) AggregationData {
	return newMeanData(a.Mean, a.Count*fraction)
}

func (a *MeanData) addOther(av AggregationData) {
	other, ok := av.(*MeanData)
	if !ok {
		return
	}

	if other.Count == 0 {
		return
	}

	a.Mean = (a.Sum() + other.Sum()) / (a.Count + other.Count)
	a.Count = a.Count + other.Count
}

func (a *MeanData) clear() {
	a.Count = 0
	a.Mean = 0
}

func (a *MeanData) equal(other AggregationData) bool {
	a2, ok := other.(*MeanData)
	if !ok {
		return false
	}
	return a.Count == a2.Count && math.Pow(a.Mean-a2.Mean, 2) < epsilon
}

// DistributionData is the aggregated data for an
// DistributionAggregation.
//
// Most users won't directly access distribution data.
type DistributionData struct {
	Count           int64     // number of data points aggregated
	Min             float64   // minimum value in the distribution
	Max             float64   // max value in the distribution
	Mean            float64   // mean of the distribution
	SumOfSquaredDev float64   // sum of the squared deviation from the mean
	CountPerBucket  []int64   // number of occurrences per bucket
	bounds          []float64 // histogram distribution of the values
}

func newDistributionData(bounds []float64) *DistributionData {
	return &DistributionData{
		CountPerBucket: make([]int64, len(bounds)+1),
		bounds:         bounds,
		Min:            math.MaxFloat64,
		Max:            math.SmallestNonzeroFloat64,
	}
}

// Sum returns the sum of all samples collected.
func (a *DistributionData) Sum() float64 { return a.Mean * float64(a.Count) }

func (a *DistributionData) variance() float64 {
	if a.Count <= 1 {
		return 0
	}
	return a.SumOfSquaredDev / float64(a.Count-1)
}

func (a *DistributionData) isAggregationData() bool { return true }

func (a *DistributionData) addSample(v interface{}) {
	var f float64
	switch x := v.(type) {
	case int64:
		f = float64(x)
	case float64:
		f = x
	default:
		return
	}

	if f < a.Min {
		a.Min = f
	}
	if f > a.Max {
		a.Max = f
	}
	a.Count++
	a.incrementBucketCount(f)

	if a.Count == 1 {
		a.Mean = f
		return
	}

	oldMean := a.Mean
	a.Mean = a.Mean + (f-a.Mean)/float64(a.Count)
	a.SumOfSquaredDev = a.SumOfSquaredDev + (f-oldMean)*(f-a.Mean)
}

func (a *DistributionData) incrementBucketCount(f float64) {
	if len(a.bounds) == 0 {
		a.CountPerBucket[0]++
		return
	}

	for i, b := range a.bounds {
		if f < b {
			a.CountPerBucket[i]++
			return
		}
	}
	a.CountPerBucket[len(a.bounds)]++
}

// DistributionData will not multiply by the fraction for this type
// of aggregation. The 'fraction' argument is there just to satisfy the
// interface 'AggregationData'. For simplicity, we include the oldest partial
// bucket in its entirety when the aggregation is a distribution. We do not try
//  to multiply it by the fraction as it would make the calculation too complex
// and will create inconsistencies between sumOfSquaredDev, min, max and the
// various buckets of the histogram.
func (a *DistributionData) multiplyByFraction(fraction float64) AggregationData {
	ret := newDistributionData(a.bounds)
	copy(ret.CountPerBucket, a.CountPerBucket)
	ret.Count = a.Count
	ret.Min = a.Min
	ret.Max = a.Max
	ret.Mean = a.Mean
	ret.SumOfSquaredDev = a.SumOfSquaredDev
	return ret
}

func (a *DistributionData) addOther(av AggregationData) {
	other, ok := av.(*DistributionData)
	if !ok {
		return
	}
	if other.Count == 0 {
		return
	}
	if other.Min < a.Min {
		a.Min = other.Min
	}
	if other.Max > a.Max {
		a.Max = other.Max
	}
	delta := other.Mean - a.Mean
	a.SumOfSquaredDev = a.SumOfSquaredDev + other.SumOfSquaredDev + math.Pow(delta, 2)*float64(a.Count*other.Count)/(float64(a.Count+other.Count))

	a.Mean = (a.Sum() + other.Sum()) / float64(a.Count+other.Count)
	a.Count = a.Count + other.Count
	for i := range other.CountPerBucket {
		a.CountPerBucket[i] = a.CountPerBucket[i] + other.CountPerBucket[i]
	}
}

func (a *DistributionData) clear() {
	a.Count = 0
	a.Min = math.MaxFloat64
	a.Max = math.SmallestNonzeroFloat64
	a.Mean = 0
	a.SumOfSquaredDev = 0
	for i := range a.CountPerBucket {
		a.CountPerBucket[i] = 0
	}
}

func (a *DistributionData) clone() AggregationData {
	counts := make([]int64, len(a.CountPerBucket))
	copy(counts, a.CountPerBucket)
	c := *a
	c.CountPerBucket = counts
	return &c
}

func (a *DistributionData) equal(other AggregationData) bool {
	a2, ok := other.(*DistributionData)
	if !ok {
		return false
	}
	if a2 == nil {
		return false
	}
	if len(a.CountPerBucket) != len(a2.CountPerBucket) {
		return false
	}
	for i := range a.CountPerBucket {
		if a.CountPerBucket[i] != a2.CountPerBucket[i] {
			return false
		}
	}
	return a.Count == a2.Count && a.Min == a2.Min && a.Max == a2.Max && math.Pow(a.Mean-a2.Mean, 2) < epsilon && math.Pow(a.variance()-a2.variance(), 2) < epsilon
}
