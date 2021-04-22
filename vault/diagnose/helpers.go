package diagnose

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

const (
	warningEventName = "warning"
	actionKey        = "actionKey"
	errorMessageKey  = attribute.Key("error.message")
)

func StartSpan(ctx context.Context, spanName string, options ...trace.SpanOption) (context.Context, trace.Span) {
	return tracer().Start(ctx, spanName, options...)
}

func Fail(ctx context.Context, message string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, message)
}

func Error(ctx context.Context, err error, options ...trace.EventOption) error {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, options...)
	return err
}

func Warn(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)

	span.AddEvent(warningEventName, trace.WithAttributes(attribute.String("message", msg)))
}

func Action(actionName string) trace.LifeCycleOption {
	return trace.WithAttributes(attribute.String(actionKey, actionName))
}

func Test(ctx context.Context, spanName string, function func(context.Context) error, options ...trace.SpanOption) error {
	ctx, span := tracer().Start(ctx, spanName, options...)
	defer span.End()

	err := function(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

var diagnoseInit sync.Once
var diagnoseTracer trace.Tracer

func tracer() trace.Tracer {
	diagnoseInit.Do(func() {
		diagnoseTracer = otel.GetTracerProvider().Tracer("vault-diagnose")
	})
	return diagnoseTracer
}
