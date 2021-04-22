package diagnose

import (
	"github.com/mitchellh/cli"
)

const status_unknown = "[      ] "
const status_ok = "\u001b[32m[  ok  ]\u001b[0m "
const status_failed = "\u001b[31m[failed]\u001b[0m "
const status_warn = "\u001b[33m[ warn ]\u001b[0m "
const same_line = "\u001b[F"

type AnsiTracer struct {
	UI     cli.Ui
	Indent bool
}

func (t *AnsiTracer) Apply(phase *Span) {

}
