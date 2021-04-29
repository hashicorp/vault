package diagnose

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

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
	ErrorStatus    = "error"
	WarningStatus  = "warn"
	OkStatus       = "ok"
)

var errUnimplemented = errors.New("unimplemented")

type Result struct {
	Time     time.Time
	Name     string
	Warnings []string
	Status   string
	Message  string
	Children []*Result
}

func (r *Result) sortChildren() {
	if len(r.Children) > 0 {
		sort.SliceStable(r.Children, func(i, j int) bool {
			return r.Children[i].Time.Before(r.Children[j].Time)
		})
		for _, c := range r.Children {
			c.sortChildren()
		}
	}
}

func (r *Result) ZeroTimes() {
	var zero time.Time
	r.Time = zero
	for _, c := range r.Children {
		c.ZeroTimes()
	}
}

// TelemetryCollector is an otel SpanProcessor that gathers spans and once the outermost
// span ends, walks the otel traces in order to produce a top-down tree of Diagnose results.
type TelemetryCollector struct {
	ui         io.Writer
	spans      map[trace.SpanID]sdktrace.ReadOnlySpan
	rootSpan   sdktrace.ReadOnlySpan
	results    map[trace.SpanID]*Result
	RootResult *Result
	mu         sync.Mutex
}

func NewTelemetryCollector(w io.Writer) *TelemetryCollector {
	return &TelemetryCollector{
		ui:      w,
		spans:   make(map[trace.SpanID]sdktrace.ReadOnlySpan),
		results: make(map[trace.SpanID]*Result),
	}
}

// OnStart tracks spans by id for later retrieval
func (t *TelemetryCollector) OnStart(_ context.Context, s sdktrace.ReadWriteSpan) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans[s.SpanContext().SpanID()] = s
	if isMainSection(s) {
		fmt.Fprintf(t.ui, status_unknown+s.Name())
	}
}

func isMainSection(s sdktrace.ReadOnlySpan) bool {
	for _, a := range s.Attributes() {
		if a.Key == "diagnose" && a.Value.AsString() == "main-section" {
			return true
		}
	}
	return false
}

func (t *TelemetryCollector) OnEnd(e sdktrace.ReadOnlySpan) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !e.Parent().HasSpanID() {
		// First walk the span structs to construct the top down tree results we want
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

		// Then walk the results sorting children by time
		t.RootResult.sortChildren()
	} else if isMainSection(e) {
		r := t.getOrBuildResult(e.SpanContext().SpanID())
		if r != nil {
			fmt.Fprintf(t.ui, same_line+r.String()+"\n\n")
		}
	}
}

// required to implement SpanProcessor, but noops for our purposes
func (t *TelemetryCollector) Shutdown(_ context.Context) error {
	return nil
}

// required to implement SpanProcessor, but noops for our purposes
func (t *TelemetryCollector) ForceFlush(_ context.Context) error {
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
			Time:    s.StartTime(),
		}
		for _, e := range s.Events() {
			switch e.Name {
			case warningEventName:
				for _, a := range e.Attributes {
					if a.Key == messageKey {
						r.Warnings = append(r.Warnings, a.Value.AsString())
					}
				}
			case ErrorStatus:
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
						Status:  ErrorStatus,
						Message: message,
					})

				}
			case spotCheckOkEventName:
				var checkName string
				var message string
				for _, a := range e.Attributes {
					switch a.Key {
					case nameKey:
						checkName = a.Value.AsString()
					case messageKey:
						message = a.Value.AsString()
					}
				}
				if checkName != "" {
					r.Children = append(r.Children,
						&Result{
							Name:    checkName,
							Status:  OkStatus,
							Message: message,
							Time:    e.Time,
						})
				}
			case spotCheckWarnEventName:
				var checkName string
				var message string
				for _, a := range e.Attributes {
					switch a.Key {
					case nameKey:
						checkName = a.Value.AsString()
					case messageKey:
						message = a.Value.AsString()
					}
				}
				if checkName != "" {
					r.Children = append(r.Children,
						&Result{
							Name:    checkName,
							Status:  WarningStatus,
							Message: message,
							Time:    e.Time,
						})
				}
			case spotCheckErrorEventName:
				var checkName string
				var message string
				for _, a := range e.Attributes {
					switch a.Key {
					case nameKey:
						checkName = a.Value.AsString()
					case messageKey:
						message = a.Value.AsString()
					}
				}
				if checkName != "" {
					r.Children = append(r.Children,
						&Result{
							Name:    checkName,
							Status:  ErrorStatus,
							Message: message,
							Time:    e.Time,
						})
				}
			}
		}
		switch s.StatusCode() {
		case codes.Unset:
			if len(r.Warnings) > 0 {
				r.Status = WarningStatus
			} else {
				r.Status = OkStatus
			}
		case codes.Ok:
			r.Status = OkStatus
		case codes.Error:
			r.Status = ErrorStatus
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
	for i := 0; i < depth; i++ {
		sb.WriteRune('\t')
	}
	sb.WriteString(r.String())
	sb.WriteRune('\n')
	for _, c := range r.Children {
		c.write(sb, depth+1)
	}
}

func (r *Result) String() string {
	var sb strings.Builder
	if r.Status != WarningStatus || (len(r.Warnings) == 0 && r.Message != "") {
		switch r.Status {
		case OkStatus:
			sb.WriteString(status_ok)
		case WarningStatus:
			sb.WriteString(status_warn)
		case ErrorStatus:
			sb.WriteString(status_failed)
		}
		sb.WriteString(r.Name)

		if r.Message != "" || len(r.Warnings) > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(r.Message)
	}
	for _, w := range r.Warnings {
		//TODO: Indentation
		sb.WriteRune('\n')
		sb.WriteString(status_warn)
		sb.WriteString(r.Name)
		sb.WriteString(": ")
		sb.WriteString(w)
	}
	return sb.String()

}
