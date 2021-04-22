package diagnose

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Result struct {
	Name     string
	Warnings []string
	Status   string
	Message  string
	Children []*Result
}

type TelemetryCollector struct {
	spans      map[trace.SpanID]sdktrace.ReadOnlySpan
	rootSpan   sdktrace.ReadOnlySpan
	results    map[trace.SpanID]*Result
	RootResult *Result
}

func NewTelemetryCollector() *TelemetryCollector {
	return &TelemetryCollector{
		spans:   make(map[trace.SpanID]sdktrace.ReadOnlySpan),
		results: make(map[trace.SpanID]*Result),
	}
}

func (t *TelemetryCollector) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
	t.spans[s.SpanContext().SpanID()] = s

}

func (t *TelemetryCollector) OnEnd(e sdktrace.ReadOnlySpan) {
	if !e.Parent().HasSpanID() {
		// Deep first walk the span structs to construct the top down tree results we want
		for _, s := range t.spans {
			r := t.getOrBuildResult(s.SpanContext().SpanID())
			if s.Parent().HasSpanID() {
				p := t.getOrBuildResult(s.Parent().SpanID())
				p.Children = append(p.Children, r)
			} else {
				t.RootResult = r
			}
		}
		fmt.Printf("%v", t.RootResult)
	}
}

func (t *TelemetryCollector) Shutdown(ctx context.Context) error {
	return nil
}

func (t *TelemetryCollector) ForceFlush(ctx context.Context) error {
	return nil
}

func (t *TelemetryCollector) getOrBuildResult(id trace.SpanID) *Result {
	s := t.spans[id]
	r, ok := t.results[id]
	if !ok {
		r = &Result{
			Name:    s.Name(),
			Message: s.StatusMessage(),
		}
		for _, e := range s.Events() {
			if e.Name == warningEventName {
				for _, a := range e.Attributes {
					if a.Key == "message" {
						r.Warnings = append(r.Warnings, a.Value.AsString())
					}
				}
			} else if e.Name == "error" {
				var message string
				var action string
				for _, a := range e.Attributes {
					switch a.Key {
					case actionKey:
						action = a.Value.AsString()
					case errorMessageKey:
						message = a.Value.AsString()
					}
				}
				if message != "" && action != "" {
					r.Children = append(r.Children, &Result{
						Name:    action,
						Status:  "error",
						Message: message,
					})

				}
			}
		}
		switch s.StatusCode() {
		case codes.Unset:
			if len(r.Warnings) > 0 {
				r.Status = "warning"
			} else {
				r.Status = codes.Ok.String()
			}
		default:
			r.Status = s.StatusCode().String()
		}
		t.results[id] = r
	}
	return r
}
