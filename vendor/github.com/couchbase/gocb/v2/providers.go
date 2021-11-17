package gocb

import (
	"context"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

// NOTE: context in these provider functions can be passed as a nil value.
// The async op manager will check for a nil context.Context, the context values should never be assumed to be non-nil.

type httpProvider interface {
	DoHTTPRequest(ctx context.Context, req *gocbcore.HTTPRequest) (*gocbcore.HTTPResponse, error)
}

type viewProvider interface {
	ViewQuery(ctx context.Context, opts gocbcore.ViewQueryOptions) (viewRowReader, error)
}

type queryProvider interface {
	N1QLQuery(ctx context.Context, opts gocbcore.N1QLQueryOptions) (queryRowReader, error)
	PreparedN1QLQuery(ctx context.Context, opts gocbcore.N1QLQueryOptions) (queryRowReader, error)
}

type analyticsProvider interface {
	AnalyticsQuery(ctx context.Context, opts gocbcore.AnalyticsQueryOptions) (analyticsRowReader, error)
}

type searchProvider interface {
	SearchQuery(ctx context.Context, opts gocbcore.SearchQueryOptions) (searchRowReader, error)
}

type waitUntilReadyProvider interface {
	WaitUntilReady(ctx context.Context, deadline time.Time, opts gocbcore.WaitUntilReadyOptions) error
}

type gocbcoreWaitUntilReadyProvider interface {
	WaitUntilReady(deadline time.Time, opts gocbcore.WaitUntilReadyOptions,
		cb gocbcore.WaitUntilReadyCallback) (gocbcore.PendingOp, error)
}

type diagnosticsProvider interface {
	Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error)
	Ping(ctx context.Context, opts gocbcore.PingOptions) (*gocbcore.PingResult, error)
}

type gocbcoreDiagnosticsProvider interface {
	Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error)
	Ping(opts gocbcore.PingOptions, cb gocbcore.PingCallback) (gocbcore.PendingOp, error)
}

type gocbcoreHTTPProvider interface {
	DoHTTPRequest(req *gocbcore.HTTPRequest, cb gocbcore.DoHTTPRequestCallback) (gocbcore.PendingOp, error)
}

type waitUntilReadyProviderWrapper struct {
	provider gocbcoreWaitUntilReadyProvider
}

func (wpw *waitUntilReadyProviderWrapper) WaitUntilReady(ctx context.Context, deadline time.Time,
	opts gocbcore.WaitUntilReadyOptions) (errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(wpw.provider.WaitUntilReady(deadline, opts, func(res *gocbcore.WaitUntilReadyResult, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		opm.Resolve()
	}))
	if err != nil {
		errOut = err
		return
	}

	return
}

type diagnosticsProviderWrapper struct {
	provider gocbcoreDiagnosticsProvider
}

func (dpw *diagnosticsProviderWrapper) Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error) {
	return dpw.provider.Diagnostics(opts)
}

func (dpw *diagnosticsProviderWrapper) Ping(ctx context.Context, opts gocbcore.PingOptions) (pOut *gocbcore.PingResult, errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(dpw.provider.Ping(opts, func(res *gocbcore.PingResult, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		pOut = res
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}

type httpProviderWrapper struct {
	provider gocbcoreHTTPProvider
}

func (hpw *httpProviderWrapper) DoHTTPRequest(ctx context.Context, req *gocbcore.HTTPRequest) (respOut *gocbcore.HTTPResponse, errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(hpw.provider.DoHTTPRequest(req, func(res *gocbcore.HTTPResponse, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		respOut = res
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}

type analyticsProviderWrapper struct {
	provider *gocbcore.AgentGroup
}

func (apw *analyticsProviderWrapper) AnalyticsQuery(ctx context.Context, opts gocbcore.AnalyticsQueryOptions) (aOut analyticsRowReader, errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(apw.provider.AnalyticsQuery(opts, func(reader *gocbcore.AnalyticsRowReader, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		aOut = reader
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}

type queryProviderWrapper struct {
	provider *gocbcore.AgentGroup
}

func (apw *queryProviderWrapper) N1QLQuery(ctx context.Context, opts gocbcore.N1QLQueryOptions) (qOut queryRowReader, errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(apw.provider.N1QLQuery(opts, func(reader *gocbcore.N1QLRowReader, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		qOut = reader
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}

func (apw *queryProviderWrapper) PreparedN1QLQuery(ctx context.Context, opts gocbcore.N1QLQueryOptions) (qOut queryRowReader, errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(apw.provider.PreparedN1QLQuery(opts, func(reader *gocbcore.N1QLRowReader, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		qOut = reader
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}

type searchProviderWrapper struct {
	provider *gocbcore.AgentGroup
}

func (apw *searchProviderWrapper) SearchQuery(ctx context.Context, opts gocbcore.SearchQueryOptions) (sOut searchRowReader, errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(apw.provider.SearchQuery(opts, func(reader *gocbcore.SearchRowReader, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		sOut = reader
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}

type viewProviderWrapper struct {
	provider *gocbcore.Agent
}

func (apw *viewProviderWrapper) ViewQuery(ctx context.Context, opts gocbcore.ViewQueryOptions) (vOut viewRowReader, errOut error) {
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(apw.provider.ViewQuery(opts, func(reader *gocbcore.ViewQueryRowReader, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		vOut = reader
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}
