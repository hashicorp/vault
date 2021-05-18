package metricsutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

const (
	OpenMetricsMIMEType = "application/openmetrics-text"

	PrometheusSchemaMIMEType = "prometheus/telemetry"

	// ErrorContentType is the content type returned by an error response.
	ErrorContentType = "text/plain"
)

const (
	PrometheusMetricFormat = "prometheus"
)

// PhysicalTableSizeName is a set of gauge metric keys for physical mount table sizes
var PhysicalTableSizeName []string = []string{"core", "mount_table", "size"}

// LogicalTableSizeName is a set of gauge metric keys for logical mount table sizes
var LogicalTableSizeName []string = []string{"core", "mount_table", "num_entries"}

type MetricsHelper struct {
	inMemSink         *metrics.InmemSink
	PrometheusEnabled bool
	LoopMetrics       GaugeMetrics
}

type GaugeMetrics struct {
	// Metrics is a map from keys concatenated by "." to the metric.
	// It is a map because although we do not care about distinguishing
	// these loop metrics during emission, we must distinguish them
	// when we update a metric.
	Metrics sync.Map
}

type GaugeMetric struct {
	Value  float32
	Labels []Label
	Key    []string
}

func NewMetricsHelper(inMem *metrics.InmemSink, enablePrometheus bool) *MetricsHelper {
	return &MetricsHelper{inMem, enablePrometheus, GaugeMetrics{Metrics: sync.Map{}}}
}

func FormatFromRequest(req *logical.Request) string {
	acceptHeaders := req.Headers["Accept"]
	if len(acceptHeaders) > 0 {
		acceptHeader := acceptHeaders[0]
		if strings.HasPrefix(acceptHeader, OpenMetricsMIMEType) {
			return PrometheusMetricFormat
		}

		// Look for prometheus accept header
		for _, header := range acceptHeaders {
			if strings.Contains(header, PrometheusSchemaMIMEType) {
				return PrometheusMetricFormat
			}
		}
	}
	return ""
}

func (m *MetricsHelper) AddGaugeLoopMetric(key []string, val float32, labels []Label) {
	mapKey := m.CreateMetricsCacheKeyName(key, val, labels)
	m.LoopMetrics.Metrics.Store(mapKey,
		GaugeMetric{
			Key:    key,
			Value:  val,
			Labels: labels,
		})
}

func (m *MetricsHelper) CreateMetricsCacheKeyName(key []string, val float32, labels []Label) string {
	var keyJoin string = strings.Join(key, ".")
	labelJoinStr := ""
	for _, label := range labels {
		labelJoinStr = labelJoinStr + label.Name + "|" + label.Value + "||"
	}
	keyJoin = keyJoin + "." + labelJoinStr
	return keyJoin
}

func (m *MetricsHelper) ResponseForFormat(format string) *logical.Response {
	switch format {
	case PrometheusMetricFormat:
		return m.PrometheusResponse()
	case "":
		return m.GenericResponse()
	default:
		return &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPContentType: ErrorContentType,
				logical.HTTPRawBody:     fmt.Sprintf("metric response format \"%s\" unknown", format),
				logical.HTTPStatusCode:  http.StatusBadRequest,
			},
		}
	}
}

func (m *MetricsHelper) PrometheusResponse() *logical.Response {
	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: ErrorContentType,
			logical.HTTPStatusCode:  http.StatusBadRequest,
		},
	}

	if !m.PrometheusEnabled {
		resp.Data[logical.HTTPRawBody] = "prometheus is not enabled"
		return resp
	}
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil && len(metricsFamilies) == 0 {
		resp.Data[logical.HTTPRawBody] = fmt.Sprintf("no prometheus metrics could be decoded: %s", err)
		return resp
	}

	// Initialize a byte buffer.
	buf := &bytes.Buffer{}
	defer buf.Reset()

	e := expfmt.NewEncoder(buf, expfmt.FmtText)
	for _, mf := range metricsFamilies {
		err := e.Encode(mf)
		if err != nil {
			resp.Data[logical.HTTPRawBody] = fmt.Sprintf("error during the encoding of metrics: %s", err)
			return resp
		}
	}
	resp.Data[logical.HTTPContentType] = string(expfmt.FmtText)
	resp.Data[logical.HTTPRawBody] = buf.Bytes()
	resp.Data[logical.HTTPStatusCode] = http.StatusOK
	return resp
}

func (m *MetricsHelper) GenericResponse() *logical.Response {
	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: ErrorContentType,
			logical.HTTPStatusCode:  http.StatusBadRequest,
		},
	}

	summary, err := m.inMemSink.DisplayMetrics(nil, nil)
	if err != nil {
		resp.Data[logical.HTTPRawBody] = fmt.Sprintf("error while fetching the in-memory metrics: %s", err)
		return resp
	}
	content, err := json.Marshal(summary)
	if err != nil {
		resp.Data[logical.HTTPRawBody] = fmt.Sprintf("error while marshalling the in-memory metrics: %s", err)
		return resp
	}
	resp.Data[logical.HTTPContentType] = "application/json"
	resp.Data[logical.HTTPRawBody] = content
	resp.Data[logical.HTTPStatusCode] = http.StatusOK
	return resp
}
