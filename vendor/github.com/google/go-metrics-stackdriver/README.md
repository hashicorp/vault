# go-metrics-stackdriver
[![godoc](https://godoc.org/github.com/google/go-metrics-stackdriver?status.svg)](http://godoc.org/github.com/google/go-metrics-stackdriver)

This library provides a stackdriver sink for applications instrumented with the
[go-metrics](https://github.com/armon/go-metrics) library.

## Details

[stackdriver.NewSink](https://godoc.org/github.com/google/go-metrics-stackdriver#NewSink)'s return value satisfies the go-metrics library's [MetricSink](https://godoc.org/github.com/armon/go-metrics#MetricSink) interface. When providing a `stackdriver.Sink` to libraries and applications instrumented against `MetricSink`, the metrics will be aggregated within this library and written to stackdriver as [Generic Task](https://cloud.google.com/monitoring/api/resources#tag_generic_task) timeseries metrics.

## Example

```go
import "github.com/google/go-metrics-stackdriver"
...
client, _ := monitoring.NewMetricClient(context.Background())
ss := stackdriver.NewSink(client, &stackdriver.Config{
  ProjectID: projectID,
})
...
ss.SetGauge([]string{"foo"}, 42)
ss.IncrCounter([]string{"baz"}, 1)
ss.AddSample([]string{"method", "const"}, 200)
```

The [full example](example/main.go) can be run from a cloud shell console to test how metrics are collected and displayed.


__This is not an officially supported Google product.__
