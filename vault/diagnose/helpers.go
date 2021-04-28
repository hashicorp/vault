package diagnose

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"io"
)

const (
	warningEventName        = "warning"
	actionKey               = "actionKey"
	spotCheckOkEventName    = "spot-check-ok"
	spotCheckWarnEventName  = "spot-check-warn"
	spotCheckErrorEventName = "spot-check-error"
	errorMessageKey         = attribute.Key("error.message")
	nameKey                 = attribute.Key("name")
	messageKey              = attribute.Key("message")
)

var (
	MainSection = trace.WithAttributes(attribute.Key("diagnose").String("main-section"))
)

var tp *sdktrace.TracerProvider
var tracer trace.Tracer
var tc *TelemetryCollector

// Init initializes a Diagnose tracing session.  In particular this wires a TelemetryCollector, which
// synchronously receives and tracks OpenTelemetry spans in order to provide a tree structure of results
// when the outermost span ends.
func Init(w io.Writer) {
	tc = NewTelemetryCollector(w)
	//so, _ := stdout.NewExporter(stdout.WithPrettyPrint())
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		//sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(so)),
		sdktrace.WithSpanProcessor(tc),
	)
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("vault-diagnose")
}

// Ends the Diagnose session, returning the root of the result tree.  This will be empty until
// the outermost span ends.
func Shutdown() *Result {
	return tc.RootResult
}

// Start a "diagnose" span, which is really just an Otel Tracing span.
func StartSpan(ctx context.Context, spanName string, options ...trace.SpanOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, options...)
}

// Fail records a failure in the current span
func Fail(ctx context.Context, message string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, message)
}

// Error records an error in the current span (but unlike Fail, doesn't set the overall span status to Error)
func Error(ctx context.Context, err error, options ...trace.EventOption) error {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, options...)
	return err
}

// Warn records a warning on the current span
func Warn(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(warningEventName, trace.WithAttributes(messageKey.String(msg)))
}

func SpotOk(ctx context.Context, checkName, message string, options ...trace.EventOption) {
	addSpotCheckResult(ctx, spotCheckOkEventName, checkName, message, options...)
}

func SpotWarn(ctx context.Context, checkName, message string, options ...trace.EventOption) {
	addSpotCheckResult(ctx, spotCheckWarnEventName, checkName, message, options...)
}

func SpotError(ctx context.Context, checkName string, err error, options ...trace.EventOption) error {
	var message string
	if err != nil {
		message = err.Error()
	}
	addSpotCheckResult(ctx, spotCheckErrorEventName, checkName, message, options...)
	return err
}

func addSpotCheckResult(ctx context.Context, eventName, checkName, message string, options ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	attrs := append(options, trace.WithAttributes(nameKey.String(checkName)))
	if message != "" {
		attrs = append(attrs, trace.WithAttributes(messageKey.String(message)))
	}
	span.AddEvent(eventName, attrs...)
}

func SpotCheck(ctx context.Context, checkName string, f func() error) error {
	err := f()
	if err != nil {
		SpotError(ctx, checkName, err)
		return err
	} else {
		SpotOk(ctx, checkName, "")
	}
	return nil
}

// Test creates a new named span, and executes the provided function within it.  If the function returns an error,
// the span is considered to have failed.
func Test(ctx context.Context, spanName string, function func(context.Context) error, options ...trace.SpanOption) error {
	ctx, span := tracer.Start(ctx, spanName, options...)
	defer span.End()

	err := function(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
