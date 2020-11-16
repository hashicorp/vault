module go.opentelemetry.io/otel/exporters/trace/jaeger

go 1.14

replace (
	go.opentelemetry.io/otel => ../../..
	go.opentelemetry.io/otel/sdk => ../../../sdk
)

require (
	github.com/apache/thrift v0.13.0
	github.com/google/go-cmp v0.5.2
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
	google.golang.org/api v0.32.0
)
