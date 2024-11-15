package gocb

import (
	"context"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

// NOTE: context in these provider functions can be passed as a nil value.
// The async op manager will check for a nil context.Context, the context values should never be assumed to be non-nil.

type httpProvider interface {
	DoHTTPRequest(ctx context.Context, req *gocbcore.HTTPRequest) (*gocbcore.HTTPResponse, error)
}

type gocbcoreDiagnosticsProvider interface {
	Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error)
	Ping(opts gocbcore.PingOptions, cb gocbcore.PingCallback) (gocbcore.PendingOp, error)
}

type gocbcoreHTTPProvider interface {
	DoHTTPRequest(req *gocbcore.HTTPRequest, cb gocbcore.DoHTTPRequestCallback) (gocbcore.PendingOp, error)
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
