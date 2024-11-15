package gocb

type diagnosticsProvider interface {
	Diagnostics(opts *DiagnosticsOptions) (*DiagnosticsResult, error)
	Ping(opts *PingOptions) (*PingResult, error)
}
