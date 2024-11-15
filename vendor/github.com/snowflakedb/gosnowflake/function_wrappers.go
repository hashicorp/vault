package gosnowflake

import "context"

// GoroutineWrapperFunc is used to wrap goroutines. This is useful if the caller wants
// to recover panics, rather than letting panics cause a system crash. A suggestion would be to
// use use the recover functionality, and log the panic as is most useful to you
type GoroutineWrapperFunc func(ctx context.Context, f func())

// The default GoroutineWrapperFunc; this does nothing. With this default wrapper
// panics will take down binary as expected
var noopGoroutineWrapper = func(_ context.Context, f func()) {
	f()
}

// GoroutineWrapper is used to hold the GoroutineWrapperFunc set by the client, or to
// store the default goroutine wrapper which does nothing
var GoroutineWrapper GoroutineWrapperFunc = noopGoroutineWrapper
