// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

	"github.com/mitchellh/go-wordwrap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	status_unknown = "[      ] "
	status_ok      = "\u001b[32m[ success ]\u001b[0m "
	status_failed  = "\u001b[31m[ failure ]\u001b[0m "
	status_warn    = "\u001b[33m[ warning ]\u001b[0m "
	status_skipped = "\u001b[90m[ skipped ]\u001b[0m "
	same_line      = "\x0d"
	ErrorStatus    = 2
	WarningStatus  = 1
	OkStatus       = 0
	SkippedStatus  = -1
)

var errUnimplemented = errors.New("unimplemented")

type status int

func (s status) String() string {
	switch s {
	case OkStatus:
		return "ok"
	case WarningStatus:
		return "warn"
	case ErrorStatus:
		return "fail"
	}
	return "invalid"
}

func (s status) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", s.String(), "\"")), nil
}

type Result struct {
	Time     time.Time `json:"time"`
	Name     string    `json:"name"`
	Status   status    `json:"status"`
	Warnings []string  `json:"warnings,omitempty"`
	Message  string    `json:"message,omitempty"`
	Advice   string
	Children []*Result `json:"children,omitempty"`
}

func (r *Result) finalize() status {
	maxStatus := r.Status
	if len(r.Children) > 0 {
		sort.SliceStable(r.Children, func(i, j int) bool {
			return r.Children[i].Time.Before(r.Children[j].Time)
		})
		for _, c := range r.Children {
			cms := c.finalize()
			if cms > maxStatus {
				maxStatus = cms
			}
		}
		if maxStatus > r.Status {
			r.Status = maxStatus
		}
	}
	return maxStatus
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

// NewTelemetryCollector creates a SpanProcessor that collects OpenTelemetry spans
// and aggregates them into a tree structure for use by Diagnose.
// It also outputs the status of main sections to that writer.
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
		fmt.Fprint(t.ui, status_unknown+s.Name())
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
		t.RootResult.finalize()
	} else if isMainSection(e) {
		r := t.getOrBuildResult(e.SpanContext().SpanID())
		if r != nil {
			fmt.Print(same_line)
			fmt.Fprintln(t.ui, r.String())
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
			Message: s.Status().Description,
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
			case skippedEventName:
				r.Status = SkippedStatus
			case "fail":
				message, action := findAttributes(e, errorMessageKey, actionKey)
				if message != "" && action != "" {
					r.Children = append(r.Children, &Result{
						Name:    action,
						Status:  ErrorStatus,
						Message: message,
					})
				}
			case spotCheckOkEventName:
				checkName, message := findAttributes(e, nameKey, messageKey)
				if checkName != "" {
					r.Children = append(r.Children,
						&Result{
							Name:    checkName,
							Status:  OkStatus,
							Message: message,
							Time:    e.Time,
							Advice:  findAttribute(e, adviceKey),
						})
				}
			case spotCheckWarnEventName:
				checkName, message := findAttributes(e, nameKey, messageKey)
				if checkName != "" {
					r.Children = append(r.Children,
						&Result{
							Name:    checkName,
							Status:  WarningStatus,
							Message: message,
							Time:    e.Time,
							Advice:  findAttribute(e, adviceKey),
						})
				}
			case spotCheckErrorEventName:
				checkName, message := findAttributes(e, nameKey, messageKey)
				if checkName != "" {
					r.Children = append(r.Children,
						&Result{
							Name:    checkName,
							Status:  ErrorStatus,
							Message: message,
							Time:    e.Time,
							Advice:  findAttribute(e, adviceKey),
						})
				}
			case spotCheckSkippedEventName:
				checkName, message := findAttributes(e, nameKey, messageKey)
				if checkName != "" {
					r.Children = append(r.Children,
						&Result{
							Name:    checkName,
							Status:  SkippedStatus,
							Message: message,
							Time:    e.Time,
							Advice:  findAttribute(e, adviceKey),
						})
				}
			case adviceEventName:
				message, _ := findAttributes(e, adviceKey, "")
				if message != "" {
					r.Advice = message
				}
			}
		}
		switch s.Status().Code {
		case codes.Unset:
			if len(r.Warnings) > 0 {
				r.Status = WarningStatus
			} else if r.Status != SkippedStatus {
				r.Status = OkStatus
			}
		case codes.Ok:
			if r.Status != SkippedStatus {
				r.Status = OkStatus
			}
		case codes.Error:
			if r.Status != SkippedStatus {
				r.Status = ErrorStatus
			}
		}
		t.results[id] = r
	}
	return r
}

func findAttribute(e sdktrace.Event, attr attribute.Key) string {
	for _, a := range e.Attributes {
		if a.Key == attr {
			return a.Value.AsString()
		}
	}
	return ""
}

func findAttributes(e sdktrace.Event, attr1, attr2 attribute.Key) (string, string) {
	var av1, av2 string
	for _, a := range e.Attributes {
		switch a.Key {
		case attr1:
			av1 = a.Value.AsString()
		case attr2:
			av2 = a.Value.AsString()
		}
	}
	return av1, av2
}

// Write outputs a human readable version of the results tree
func (r *Result) Write(writer io.Writer, wrapLimit int) error {
	var sb strings.Builder
	r.write(&sb, 0, wrapLimit)
	_, err := writer.Write([]byte(sb.String()))
	return err
}

const (
	indentString    = "  "
	statusPrefixLen = 9
)

func indent(sb *strings.Builder, depth int) {
	for i := 0; i < depth; i++ {
		sb.WriteString(indentString)
	}
}

func (r *Result) String() string {
	return r.StringWrapped(80)
}

func (r *Result) StringWrapped(wrapLimit int) string {
	var sb strings.Builder
	r.write(&sb, 0, wrapLimit)
	return sb.String()
}

func (r *Result) write(sb *strings.Builder, depth int, limit int) {
	indent(sb, depth)
	var prelude string
	switch r.Status {
	case OkStatus:
		prelude = status_ok
	case WarningStatus:
		prelude = status_warn
	case ErrorStatus:
		prelude = status_failed
	case SkippedStatus:
		prelude = status_skipped
	}
	prelude = prelude + r.Name

	if r.Message != "" {
		prelude = prelude + ": " + r.Message
	}
	warnings := r.Warnings
	if r.Message == "" && len(warnings) > 0 {
		prelude = status_warn + r.Name + ": "
		if len(warnings) == 1 {
			prelude = prelude + warnings[0]
			warnings = warnings[1:]
		}
	}

	writeWrapped(sb, prelude, depth+1, limit)
	for _, w := range warnings {
		sb.WriteRune('\n')
		indent(sb, depth+1)
		sb.WriteString(status_warn)
		writeWrapped(sb, w, depth+2, limit)
	}

	if r.Advice != "" {
		advice := "\u001b[35m" + r.Advice + "\u001b[0m"
		sb.WriteRune('\n')
		indent(sb, depth+1)
		writeWrapped(sb, advice, depth+1, limit)
	}
	sb.WriteRune('\n')
	for _, c := range r.Children {
		c.write(sb, depth+1, limit)
	}
}

func writeWrapped(sb *strings.Builder, msg string, depth int, limit int) {
	if limit > 0 {
		sz := uint(limit - depth*len(indentString))
		msg = wordwrap.WrapString(msg, sz)
		parts := strings.Split(msg, "\n")
		sb.WriteString(parts[0])
		for _, p := range parts[1:] {
			sb.WriteRune('\n')
			indent(sb, depth)
			sb.WriteString(p)
		}
	} else {
		sb.WriteString(msg)
	}
}
