package diagnose

import (
	"context"
	"github.com/mitchellh/cli"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	trace3 "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"sync"
)

var currentPhase = struct{}{}

type Status int

const (
	StatusOk    = 0
	StatusWarn  = iota
	StatusError = iota
	StatusFatal = iota
)

type SpanProcessor struct {
	UI     cli.Ui
	Indent bool
	Spans  map[trace.SpanID]trace3.ReadWriteSpan
	mu     sync.RWMutex
}

func (d *SpanProcessor) Shutdown(_ context.Context) error {
	return nil
}

func (d *SpanProcessor) ForceFlush(_ context.Context) error {
	return nil
}

func NewDiagnoseSpanProcessor(ui cli.Ui) *SpanProcessor {
	return &SpanProcessor{
		UI:    ui,
		Spans: make(map[trace.SpanID]trace3.ReadWriteSpan),
	}
}

// OnStart does nothing.
func (d *SpanProcessor) OnStart(_ context.Context, sp trace3.ReadWriteSpan) {
	d.Spans[sp.SpanContext().SpanID()] = sp
}

// OnEnd immediately exports a ReadOnlySpan.
func (d *SpanProcessor) OnEnd(s trace3.ReadOnlySpan) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if s.SpanContext().IsSampled() {
		span, ok := d.Spans[s.SpanContext().SpanID()]
		if ok {

			var sb strings.Builder

			if d.Indent {
				cp := span.Parent()
				for cp.HasSpanID() {
					ps := d.Spans[cp.SpanID()]
					cp = ps.Parent()
					sb.WriteRune('\t')
				}
			}

			switch span.StatusCode() {
			case codes.Ok:
				sb.WriteString(status_ok)
			case codes.Unset:
				sb.WriteString(status_warn)
			case codes.Error:
				sb.WriteString(status_failed)
			}
			sb.WriteString(span.Name())
			if len(span.StatusMessage()) > 0 {
				sb.WriteString(": ")
				sb.WriteString(span.StatusMessage())
			}
			d.UI.Output(sb.String())
		}
	}
}

type Span struct {
	name     string
	tracer   trace.Tracer
	span     trace.Span
	Status   Status
	messages []string
	children []*Span
	options  []PhaseOption
	parent   *Span
}

type PhaseOption interface {
	Apply(*Span)
}

func StartPhase(ctx context.Context, phaseName string, options ...PhaseOption) (context.Context, *Span) {
	cPhase := CurrentPhase(ctx)
	phase := Span{
		name:    phaseName,
		options: options,
		parent:  cPhase,
	}
	if cPhase != nil {
		cPhase.children = append(cPhase.children, &phase)
		phase.tracer = cPhase.tracer
	} else {
		phase.tracer = otel.GetTracerProvider().Tracer(phaseName)
	}
	ctx, phase.span = phase.tracer.Start(ctx, phaseName)
	return context.WithValue(ctx, currentPhase, &phase), &phase
}

func CurrentPhase(ctx context.Context) *Span {
	cPhaseVal := ctx.Value(currentPhase)
	if cPhaseVal == nil {
		return nil
	}
	return cPhaseVal.(*Span)
}

func Test(ctx context.Context, phaseName string, f func(context.Context) error) error {
	ctx, phase := StartPhase(ctx, phaseName)
	defer phase.End()
	err := f(ctx)
	if err != nil {
		phase.Error(err)
	}
	return err
}

func PassFail(ctx context.Context, phaseName string, err error) error {
	ctx, phase := StartPhase(ctx, phaseName)
	if err != nil {
		phase.Error(err)
	}
	phase.End()
	return err
}

func Error(ctx context.Context, err error, options ...trace.EventOption) error {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, options...)
	return err
}

func (p *Span) Error(err error) error {
	span,
	p.setStatus(StatusError, err.Error())
	return err
}

func (p *Span) Warn(msg string) {
	p.setStatus(StatusWarn, msg)
}

func (p *Span) Fatal(err error) error {
	p.setStatus(StatusFatal, err.Error())
	return err
}

func (p *Span) End(spanOptions ...trace.SpanOption) {
	p.span.End(spanOptions...)
	if p.Status == StatusOk {
		p.span.SetStatus(codes.Ok, "")
	}
	for _, opt := range p.options {
		opt.Apply(p)
	}
}

func (p *Span) setStatus(status Status, msg string) {
	p.messages = append(p.messages, msg)
	if p.Status < status {
		p.Status = status
	}
	if p.span != nil {
		switch p.Status {
		case StatusError:
			p.span.SetStatus(codes.Error, strings.Join(p.messages, "; "))
		case StatusOk:
			p.span.SetStatus(codes.Ok, strings.Join(p.messages, "; "))
		}
	}
}
