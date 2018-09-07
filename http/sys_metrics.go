package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/armon/go-metrics"
)

func handleSysMetrics(inm *metrics.InmemSink, prometheusRetention time.Duration) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if format := req.URL.Query().Get("format"); format == "prometheus" {

				if prometheusRetention.Nanoseconds() == 0 {
					res.WriteHeader(500)
					res.Write([]byte("prometheus support is not enabled"))
				}

				handlerOptions := promhttp.HandlerOpts{
					ErrorHandling:      promhttp.ContinueOnError,
					DisableCompression: true,
				}

				handler := promhttp.HandlerFor(prometheus.DefaultGatherer, handlerOptions)
				handler.ServeHTTP(res, req)
				return
			}
			summary, err := inm.DisplayMetrics(res, req)
			if err != nil {
				res.WriteHeader(500)
				res.Write([]byte(err.Error()))
			} else {
				content, err := json.Marshal(summary)
				if err != nil {
					res.WriteHeader(500)
					res.Write([]byte("can't marshal internal metrics into json"))
				} else {
					res.WriteHeader(200)
					res.Write(content)
				}
			}
		})
}
