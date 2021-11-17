package api

import (
	"io/ioutil"
	"time"
)

// MetricsSummary holds a roll-up of metrics info for a given interval
type MetricsSummary struct {
	Timestamp string
	Gauges    []GaugeValue
	Points    []PointValue
	Counters  []SampledValue
	Samples   []SampledValue
}

type GaugeValue struct {
	Name  string
	Hash  string `json:"-"`
	Value float32

	Labels        []Label           `json:"-"`
	DisplayLabels map[string]string `json:"Labels"`
}

type PointValue struct {
	Name   string
	Points []float32
}

type SampledValue struct {
	Name string
	Hash string `json:"-"`
	*AggregateSample
	Mean   float64
	Stddev float64

	Labels        []Label           `json:"-"`
	DisplayLabels map[string]string `json:"Labels"`
}

// AggregateSample is used to hold aggregate metrics
// about a sample
type AggregateSample struct {
	Count       int       // The count of emitted pairs
	Rate        float64   // The values rate per time unit (usually 1 second)
	Sum         float64   // The sum of values
	SumSq       float64   `json:"-"` // The sum of squared values
	Min         float64   // Minimum value
	Max         float64   // Maximum value
	LastUpdated time.Time `json:"-"` // When value was last updated
}

type Label struct {
	Name  string
	Value string
}

// Metrics returns a slice of bytes containing metrics, optionally formatted as either json or prometheus
func (op *Operator) Metrics(q *QueryOptions) ([]byte, error) {
	if q == nil {
		q = &QueryOptions{}
	}

	metricsReader, err := op.c.rawQuery("/v1/metrics", q)
	if err != nil {
		return nil, err
	}

	metricsBytes, err := ioutil.ReadAll(metricsReader)
	if err != nil {
		return nil, err
	}

	return metricsBytes, nil
}

// MetricsSummary returns a MetricsSummary struct and query metadata
func (op *Operator) MetricsSummary(q *QueryOptions) (*MetricsSummary, *QueryMeta, error) {
	var resp *MetricsSummary
	qm, err := op.c.query("/v1/metrics", &resp, q)
	if err != nil {
		return nil, nil, err
	}

	return resp, qm, nil
}
