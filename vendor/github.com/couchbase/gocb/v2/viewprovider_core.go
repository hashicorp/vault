package gocb

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

type jsonViewResponse struct {
	TotalRows uint64      `json:"total_rows,omitempty"`
	DebugInfo interface{} `json:"debug_info,omitempty"`
}

type jsonViewRow struct {
	ID    string          `json:"id"`
	Key   json.RawMessage `json:"key"`
	Value json.RawMessage `json:"value"`
}

type viewProviderCore struct {
	provider viewProviderCoreProvider

	bucketName           string
	retryStrategyWrapper *coreRetryStrategyWrapper
	transcoder           Transcoder
	timeouts             TimeoutsConfig
	tracer               RequestTracer
	meter                *meterWrapper
}

// ViewQuery performs a view query and returns a list of rows or an error.
func (v *viewProviderCore) ViewQuery(designDoc string, viewName string, opts *ViewOptions) (*ViewResult, error) {
	start := time.Now()
	defer v.meter.ValueRecord(meterValueServiceViews, "views", start)

	designDoc = v.maybePrefixDevDocument(opts.Namespace, designDoc)

	span := createSpan(v.tracer, opts.ParentSpan, "views", "views")
	span.SetAttribute("db.name", v.bucketName)
	span.SetAttribute("db.operation", designDoc+"/"+viewName)
	defer span.End()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = v.timeouts.ViewTimeout
	}
	deadline := time.Now().Add(timeout)

	retryWrapper := v.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryWrapper = newCoreRetryStrategyWrapper(opts.RetryStrategy)
	}

	urlValues, err := opts.toURLValues()
	if err != nil {
		return nil, wrapError(err, "could not parse query options")
	}

	return v.execViewQuery(opts.Context, span.Context(), "_view", designDoc, viewName, *urlValues, deadline,
		retryWrapper, opts.Internal.User)
}

func (v *viewProviderCore) execViewQuery(
	ctx context.Context,
	span RequestSpanContext,
	viewType, ddoc, viewName string,
	options url.Values,
	deadline time.Time,
	wrapper *coreRetryStrategyWrapper,
	user string,
) (*ViewResult, error) {
	res, err := v.provider.ViewQuery(ctx, gocbcore.ViewQueryOptions{
		DesignDocumentName: ddoc,
		ViewType:           viewType,
		ViewName:           viewName,
		Options:            options,
		RetryStrategy:      wrapper,
		Deadline:           deadline,
		TraceContext:       span,
		User:               user,
	})
	if err != nil {
		return nil, maybeEnhanceViewError(err)
	}

	return newViewResult(res), nil
}

func (v *viewProviderCore) maybePrefixDevDocument(namespace DesignDocumentNamespace, ddoc string) string {
	designDoc := ddoc
	if namespace == DesignDocumentNamespaceProduction {
		designDoc = strings.TrimPrefix(ddoc, "dev_")
	} else {
		if !strings.HasPrefix(ddoc, "dev_") {
			designDoc = "dev_" + ddoc
		}
	}

	return designDoc
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

type viewProviderCoreProvider interface {
	ViewQuery(ctx context.Context, opts gocbcore.ViewQueryOptions) (viewRowReader, error)
}
