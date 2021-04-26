package diagnose

import (
	"context"
	"errors"
	"io"
	"strings"

	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	status_unknown = "[      ] "
	status_ok      = "\u001b[32m[  ok  ]\u001b[0m "
	status_failed  = "\u001b[31m[failed]\u001b[0m "
	status_warn    = "\u001b[33m[ warn ]\u001b[0m "
	same_line      = "\u001b[F"
	errorStatus    = "error"
	warnStatus     = "warn"
	okStatus       = "ok"
)

var errUnimplemented = errors.New("unimplemented")

type Result struct {
	Name     string
	Warnings []string
	Status   string
	Message  string
	Children []*Result
}

// TelemetryCollector is an otel SpanProcessor that gathers spans and once the outermost
// span ends, walks the otel traces in order to produce a top-down tree of Diagnose results.
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

// OnStart tracks spans by id for later retrieval
func (t *TelemetryCollector) OnStart(_ context.Context, s sdktrace.ReadWriteSpan) {
	t.spans[s.SpanContext().SpanID()] = s
}

func (t *TelemetryCollector) OnEnd(e sdktrace.ReadOnlySpan) {
	if !e.Parent().HasSpanID() {
		// Deep first walk the span structs to construct the top down tree results we want
		for _, s := range t.spans {
			r := t.getOrBuildResult(s.SpanContext().SpanID())
			if r != nil {
				if s.Parent().HasSpanID() {
					p := t.getOrBuildResult(s.Parent().SpanID())
					if p != nil {
						p.Children = append(p.Children, r)
					}
				} else {
					t.RootResult = r
				}
			}
		}
	}
}

// required to implement SpanProcessor, but noops for our purposes
func (t *TelemetryCollector) Shutdown(ctx context.Context) error {
	return nil
}

// required to implement SpanProcessor, but noops for our purposes
func (t *TelemetryCollector) ForceFlush(ctx context.Context) error {
	return nil
}

func (t *TelemetryCollector) getOrBuildResult(id trace.SpanID) *Result {
	s := t.spans[id]
	if s == nil {
		return nil
	}
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
			} else if e.Name == errorStatus {
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
						Status:  errorStatus,
						Message: message,
					})

				}
			}
		}
		switch s.StatusCode() {
		case codes.Unset:
			if len(r.Warnings) > 0 {
				r.Status = warnStatus
			} else {
				r.Status = okStatus
			}
		case codes.Ok:
			r.Status = okStatus
		case codes.Error:
			r.Status = errorStatus
		}
		t.results[id] = r
	}
	return r
}

// Write outputs a human readable version of the results tree
func (r *Result) Write(writer io.Writer) error {
	var sb strings.Builder
	r.write(&sb, 0)
	_, err := writer.Write([]byte(sb.String()))
	return err
}

func (r *Result) write(sb *strings.Builder, depth int) {
	if r.Status != warnStatus || (len(r.Warnings) == 0 && r.Message != "") {
		for i := 0; i < depth; i++ {
			sb.WriteRune('\t')
		}
		switch r.Status {
		case okStatus:
			sb.WriteString(status_ok)
		case warnStatus:
			sb.WriteString(status_warn)
		case errorStatus:
			sb.WriteString(status_failed)
		}
		sb.WriteString(r.Name)

		if r.Message != "" || len(r.Warnings) > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(r.Message)
	}
	for _, w := range r.Warnings {
		sb.WriteRune('\n')
		for i := 0; i < depth; i++ {
			sb.WriteRune('\t')
		}
		sb.WriteString(status_warn)
		sb.WriteString(r.Name)
		sb.WriteString(": ")
		sb.WriteString(w)
	}
	sb.WriteRune('\n')
	for _, c := range r.Children {
		c.write(sb, depth+1)
	}
}
