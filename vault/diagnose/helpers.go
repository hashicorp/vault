package diagnose

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	warningEventName          = "warning"
	skippedEventName          = "skipped"
	actionKey                 = "actionKey"
	spotCheckOkEventName      = "spot-check-ok"
	spotCheckWarnEventName    = "spot-check-warn"
	spotCheckErrorEventName   = "spot-check-error"
	spotCheckSkippedEventName = "spot-check-skipped"
	adviceEventName           = "advice"
	errorMessageKey           = attribute.Key("error.message")
	nameKey                   = attribute.Key("name")
	messageKey                = attribute.Key("message")
	adviceKey                 = attribute.Key("advice")
)

var MainSection = trace.WithAttributes(attribute.Key("diagnose").String("main-section"))

var (
	diagnoseSession = struct{}{}
	noopTracer      = trace.NewNoopTracerProvider().Tracer("vault-diagnose")
)

type testFunction func(context.Context) error

type Session struct {
	tc          *TelemetryCollector
	tracer      trace.Tracer
	tp          *sdktrace.TracerProvider
	SkipFilters []string
}

// New initializes a Diagnose tracing session.  In particular this wires a TelemetryCollector, which
// synchronously receives and tracks OpenTelemetry spans in order to provide a tree structure of results
// when the outermost span ends.
func New(w io.Writer) *Session {
	tc := NewTelemetryCollector(w)
	// so, _ := stdout.NewExporter(stdout.WithPrettyPrint())
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(so)),
		sdktrace.WithSpanProcessor(tc),
	)
	tracer := tp.Tracer("vault-diagnose")
	sess := &Session{
		tp:     tp,
		tc:     tc,
		tracer: tracer,
	}
	return sess
}

// IsSkipped returns true if skipName is present in the SkipFilters list.  Can be used in combination with Skip to mark a
// span skipped and conditionally skips some logic.
func (s *Session) IsSkipped(spanName string) bool {
	return strutil.StrListContainsCaseInsensitive(s.SkipFilters, spanName)
}

// Context returns a new context with a defined diagnose session
func Context(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, diagnoseSession, sess)
}

// CurrentSession retrieves the active diagnose session from the context, or nil if none.
func CurrentSession(ctx context.Context) *Session {
	sessionCtxVal := ctx.Value(diagnoseSession)
	if sessionCtxVal != nil {
		return sessionCtxVal.(*Session)
	}
	return nil
}

// Finalize ends the Diagnose session, returning the root of the result tree.  This will be empty until
// the outermost span ends.
func (s *Session) Finalize(ctx context.Context) *Result {
	s.tp.ForceFlush(ctx)
	return s.tc.RootResult
}

// StartSpan starts a "diagnose" span, which is really just an OpenTelemetry Tracing span.
func StartSpan(ctx context.Context, spanName string, options ...trace.SpanStartOption) (context.Context, trace.Span) {
	session := CurrentSession(ctx)
	if session != nil {
		return session.tracer.Start(ctx, spanName, options...)
	} else {
		return noopTracer.Start(ctx, spanName, options...)
	}
}

// Success sets the span to Successful (overriding any previous status) and sets the message to the input.
func Success(ctx context.Context, message string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Ok, message)
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

// Skipped marks the current span skipped
func Skipped(ctx context.Context, message string) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(skippedEventName)
	span.SetStatus(codes.Error, message)
}

// Warn records a warning on the current span
func Warn(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(warningEventName, trace.WithAttributes(messageKey.String(msg)))
}

// SpotOk adds an Ok result without adding a new Span.  This should be used for instantaneous checks with no
// possible sub-spans
func SpotOk(ctx context.Context, checkName, message string, options ...trace.EventOption) {
	addSpotCheckResult(ctx, spotCheckOkEventName, checkName, message, options...)
}

// SpotWarn adds a Warning result without adding a new Span.  This should be used for instantaneous checks with no
// possible sub-spans
func SpotWarn(ctx context.Context, checkName, message string, options ...trace.EventOption) {
	addSpotCheckResult(ctx, spotCheckWarnEventName, checkName, message, options...)
}

// SpotError adds an Error result without adding a new Span.  This should be used for instantaneous checks with no
// possible sub-spans
func SpotError(ctx context.Context, checkName string, err error, options ...trace.EventOption) error {
	var message string
	if err != nil {
		message = err.Error()
	}
	addSpotCheckResult(ctx, spotCheckErrorEventName, checkName, message, options...)
	return err
}

// SpotSkipped adds a Skipped result without adding a new Span.
func SpotSkipped(ctx context.Context, checkName, message string, options ...trace.EventOption) {
	addSpotCheckResult(ctx, spotCheckSkippedEventName, checkName, message, options...)
}

// Advice builds an EventOption containing advice message.  Use to add to spot results.
func Advice(message string) trace.EventOption {
	return trace.WithAttributes(adviceKey.String(message))
}

// Advise adds advice to the current diagnose span
func Advise(ctx context.Context, message string) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(adviceEventName, Advice(message))
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
	sess := CurrentSession(ctx)
	if sess.IsSkipped(checkName) {
		SpotSkipped(ctx, checkName, "skipped as requested")
		return nil
	}

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
func Test(ctx context.Context, spanName string, function testFunction, options ...trace.SpanStartOption) error {
	ctx, span := StartSpan(ctx, spanName, options...)
	defer span.End()
	sess := CurrentSession(ctx)
	if sess.IsSkipped(spanName) {
		Skipped(ctx, "skipped as requested")
		return nil
	}

	err := function(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// WithTimeout wraps a context consuming function, and when called, returns an error if the sub-function does not
// complete within the timeout, e.g.
//
// diagnose.Test(ctx, "my-span", diagnose.WithTimeout(5 * time.Second, myTestFunc))
func WithTimeout(d time.Duration, f testFunction) testFunction {
	return func(ctx context.Context) error {
		rch := make(chan error)
		t := time.NewTimer(d)
		defer t.Stop()
		go func() { rch <- f(ctx) }()
		select {
		case <-t.C:
			return fmt.Errorf("Timeout after %s.", d.String())
		case err := <-rch:
			return err
		}
	}
}

// CapitalizeFirstLetter returns a string with the first letter capitalized
func CapitalizeFirstLetter(msg string) string {
	words := strings.Split(msg, " ")
	if len(words) == 0 {
		return ""
	}
	if len(words) > 1 {
		return strings.Title(words[0]) + " " + strings.Join(words[1:], " ")
	}
	return strings.Title(words[0])
}
