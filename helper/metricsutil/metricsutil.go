package metricsutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/logical"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"strings"
)

const (
	OpenMetricsMIMEType = "application/openmetrics-text"
)

const (
	PrometheusMetricFormat = "prometheus"
)

type MetricsHelper struct {
	inMemSink         *metrics.InmemSink
	PrometheusEnabled bool
}

func NewMetricsHelper(inMem *metrics.InmemSink, enablePrometheus bool) *MetricsHelper {
	return &MetricsHelper{inMem, enablePrometheus}
}

func FormatFromRequest(req *logical.Request) string {
	acceptHeaders := req.Headers["Accept"]
	if len(acceptHeaders) > 0 {
		acceptHeader := acceptHeaders[0]
		if strings.HasPrefix(acceptHeader, OpenMetricsMIMEType) {
			return "prometheus"
		}
	}
	return ""
}

func (m *MetricsHelper) ResponseForFormat(format string) (*logical.Response, error) {
	switch format {
	case PrometheusMetricFormat:
		return m.PrometheusResponse()
	case "":
		return m.GenericResponse()
	default:
		return nil, fmt.Errorf("metric response format \"%s\" unknown", format)
	}
}

func (m *MetricsHelper) PrometheusResponse() (*logical.Response, error) {
	if !m.PrometheusEnabled {
		return &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPContentType: "text/plain",
				logical.HTTPRawBody:     "prometheus is not enabled",
				logical.HTTPStatusCode:  400,
			},
		}, nil
	}
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil && len(metricsFamilies) == 0 {
		return nil, fmt.Errorf("no prometheus metrics could be decoded: %s", err)
	}

	// Initialize a byte buffer.
	buf := &bytes.Buffer{}
	defer buf.Reset()

	e := expfmt.NewEncoder(buf, expfmt.FmtText)
	for _, mf := range metricsFamilies {
		err := e.Encode(mf)
		if err != nil {
			return nil, fmt.Errorf("error during the encoding of metrics: %s", err)
		}
	}
	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: string(expfmt.FmtText),
			logical.HTTPRawBody:     buf.Bytes(),
			logical.HTTPStatusCode:  200,
		},
	}, nil
}

func (m *MetricsHelper) GenericResponse() (*logical.Response, error) {
	summary, err := m.inMemSink.DisplayMetrics(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error while fetching the in-memory metrics: %s", err)
	}
	content, err := json.Marshal(summary)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling the in-memory metrics: %s", err)
	}
	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "application/json",
			logical.HTTPRawBody:     content,
			logical.HTTPStatusCode:  200,
		},
	}, nil
}
