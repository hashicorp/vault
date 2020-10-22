package gocb

import (
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

type httpProvider interface {
	DoHTTPRequest(req *gocbcore.HTTPRequest) (*gocbcore.HTTPResponse, error)
}

type viewProvider interface {
	ViewQuery(opts gocbcore.ViewQueryOptions) (viewRowReader, error)
}

type queryProvider interface {
	N1QLQuery(opts gocbcore.N1QLQueryOptions) (queryRowReader, error)
	PreparedN1QLQuery(opts gocbcore.N1QLQueryOptions) (queryRowReader, error)
}

type analyticsProvider interface {
	AnalyticsQuery(opts gocbcore.AnalyticsQueryOptions) (analyticsRowReader, error)
}

type searchProvider interface {
	SearchQuery(opts gocbcore.SearchQueryOptions) (searchRowReader, error)
}

type waitUntilReadyProvider interface {
	WaitUntilReady(deadline time.Time, opts gocbcore.WaitUntilReadyOptions) error
}

type gocbcoreWaitUntilReadyProvider interface {
	WaitUntilReady(deadline time.Time, opts gocbcore.WaitUntilReadyOptions,
		cb gocbcore.WaitUntilReadyCallback) (gocbcore.PendingOp, error)
}

type diagnosticsProvider interface {
	Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error)
	Ping(opts gocbcore.PingOptions) (*gocbcore.PingResult, error)
}

type gocbcoreDiagnosticsProvider interface {
	Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error)
	Ping(opts gocbcore.PingOptions, cb gocbcore.PingCallback) (gocbcore.PendingOp, error)
}

type waitUntilReadyProviderWrapper struct {
	provider gocbcoreWaitUntilReadyProvider
}

func (wpw *waitUntilReadyProviderWrapper) WaitUntilReady(deadline time.Time, opts gocbcore.WaitUntilReadyOptions) (errOut error) {
	opm := newAsyncOpManager()
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
	}

	return
}

type diagnosticsProviderWrapper struct {
	provider gocbcoreDiagnosticsProvider
}

func (dpw *diagnosticsProviderWrapper) Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error) {
	return dpw.provider.Diagnostics(opts)
}

func (dpw *diagnosticsProviderWrapper) Ping(opts gocbcore.PingOptions) (pOut *gocbcore.PingResult, errOut error) {
	opm := newAsyncOpManager()
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
	provider *gocbcore.AgentGroup
}

func (hpw *httpProviderWrapper) DoHTTPRequest(req *gocbcore.HTTPRequest) (respOut *gocbcore.HTTPResponse, errOut error) {
	opm := newAsyncOpManager()
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

func (apw *analyticsProviderWrapper) AnalyticsQuery(opts gocbcore.AnalyticsQueryOptions) (aOut analyticsRowReader, errOut error) {
	opm := newAsyncOpManager()
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

func (apw *queryProviderWrapper) N1QLQuery(opts gocbcore.N1QLQueryOptions) (qOut queryRowReader, errOut error) {
	opm := newAsyncOpManager()
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

func (apw *queryProviderWrapper) PreparedN1QLQuery(opts gocbcore.N1QLQueryOptions) (qOut queryRowReader, errOut error) {
	opm := newAsyncOpManager()
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

func (apw *searchProviderWrapper) SearchQuery(opts gocbcore.SearchQueryOptions) (sOut searchRowReader, errOut error) {
	opm := newAsyncOpManager()
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
	provider *gocbcore.AgentGroup
}

func (apw *viewProviderWrapper) ViewQuery(opts gocbcore.ViewQueryOptions) (vOut viewRowReader, errOut error) {
	opm := newAsyncOpManager()
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
