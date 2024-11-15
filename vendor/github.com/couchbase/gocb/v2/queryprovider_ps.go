package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/couchbase/goprotostellar/genproto/kv_v1"
	"github.com/couchbase/goprotostellar/genproto/query_v1"

	"google.golang.org/protobuf/types/known/durationpb"
)

type queryProviderPs struct {
	provider query_v1.QueryServiceClient

	managerProvider *psOpManagerProvider
}

func (qpc *queryProviderPs) Query(statement string, s *Scope, opts *QueryOptions) (resOut *QueryResult, errOut error) {
	attribs := map[string]interface{}{
		"db.statement": statement,
	}
	if s != nil {
		attribs["db.name"] = s.BucketName()
		attribs["db.couchbase.scope"] = s.Name()
	}

	manager := qpc.managerProvider.NewManager(opts.ParentSpan, "query", attribs)
	// Spans in couchbase2 mode need to live for the lifetime of the response body as any underlying
	// grpc span will do so.
	defer manager.ValueRecord()
	defer func() {
		if errOut != nil {
			manager.Finish(true)
		}
	}()

	manager.SetIsIdempotent(opts.Readonly)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)
	manager.SetOperationID(opts.ClientContextID)

	if err := manager.CheckReadyForOp(); err != nil {
		return nil, err
	}

	prepared := !opts.Adhoc
	req := &query_v1.QueryRequest{
		Statement: statement,
		Prepared:  &prepared,
	}
	if s != nil {
		req.BucketName = &s.bucket.bucketName
		req.ScopeName = &s.scopeName
	}
	if opts.Readonly {
		req.ReadOnly = &opts.Readonly
	}
	req.TuningOptions = &query_v1.QueryRequest_TuningOptions{}
	if opts.MaxParallelism > 0 {
		req.TuningOptions.MaxParallelism = &opts.MaxParallelism
	}
	if opts.PipelineBatch > 0 {
		req.TuningOptions.PipelineBatch = &opts.PipelineBatch
	}
	if opts.PipelineCap > 0 {
		req.TuningOptions.PipelineCap = &opts.PipelineCap
	}
	if opts.ScanWait > 0 {
		req.TuningOptions.ScanWait = durationpb.New(opts.ScanWait)
	}
	if opts.ScanCap > 0 {
		req.TuningOptions.ScanCap = &opts.ScanCap
	}
	disableMetrics := !opts.Metrics
	if disableMetrics {
		req.TuningOptions.DisableMetrics = &disableMetrics
	}
	if opts.ClientContextID == "" {
		clientContextID := uuid.NewString()
		req.ClientContextId = &clientContextID
	} else {
		req.ClientContextId = &opts.ClientContextID
	}
	manager.SetOperationID(*req.ClientContextId)

	if opts.ScanConsistency != 0 {
		var consistency query_v1.QueryRequest_ScanConsistency
		if opts.ScanConsistency == QueryScanConsistencyNotBounded {
			consistency = query_v1.QueryRequest_SCAN_CONSISTENCY_NOT_BOUNDED
		} else if opts.ScanConsistency == QueryScanConsistencyRequestPlus {
			consistency = query_v1.QueryRequest_SCAN_CONSISTENCY_REQUEST_PLUS
		} else {
			return nil, makeInvalidArgumentsError("unexpected consistency option")
		}
		req.ScanConsistency = &consistency
	}

	if len(opts.PositionalParameters) > 0 {
		params := make([][]byte, len(opts.PositionalParameters))
		for i, param := range opts.PositionalParameters {
			b, err := json.Marshal(param)
			if err != nil {
				return nil, err
			}

			params[i] = b
		}

		req.PositionalParameters = params
	}
	if len(opts.NamedParameters) > 0 {
		params := make(map[string][]byte, len(opts.NamedParameters))
		for k, param := range opts.NamedParameters {
			b, err := json.Marshal(param)
			if err != nil {
				return nil, err
			}

			params[k] = b
		}

		req.NamedParameters = params
	}
	if opts.FlexIndex {
		req.FlexIndex = &opts.FlexIndex
	}
	if opts.PreserveExpiry {
		req.PreserveExpiry = &opts.PreserveExpiry
	}

	if opts.ConsistentWith != nil {
		tokens := make([]*kv_v1.MutationToken, len(opts.ConsistentWith.tokens))
		for i, tok := range opts.ConsistentWith.tokens {
			tokens[i] = &kv_v1.MutationToken{
				BucketName:  tok.BucketName(),
				VbucketId:   uint32(tok.PartitionID()),
				VbucketUuid: tok.PartitionUUID(),
				SeqNo:       tok.SequenceNumber(),
			}
		}
		req.ConsistentWith = tokens
	}

	if opts.Profile != "" {
		var profileMode query_v1.QueryRequest_ProfileMode
		switch opts.Profile {
		case QueryProfileModeNone:
			profileMode = query_v1.QueryRequest_PROFILE_MODE_OFF
		case QueryProfileModePhases:
			profileMode = query_v1.QueryRequest_PROFILE_MODE_PHASES
		case QueryProfileModeTimings:
			profileMode = query_v1.QueryRequest_PROFILE_MODE_TIMINGS
		default:
			manager.Finish(true)
			return nil, makeInvalidArgumentsError("unexpected profile mode option")
		}
		req.ProfileMode = &profileMode
	}

	timeout := manager.Timeout()
	userCtx := opts.Context
	if userCtx == nil {
		userCtx = context.Background()
	}
	// We create a context with a timeout which will control timing out the initial request portion
	// of the operation. We can defer the cancel for this as we aren't applying this context directly
	// to the request so cancellation will not terminate any streams.
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
	defer timeoutCancel()

	var cancellationIsTimeout uint32
	// This second context has no real parent and will be cancelled if the user context is cancelled or the timeout
	// is reached. However, if the user context does not get cancelled during the initial request portion of the
	// operation then this context will live for the lifetime of the op and be used for cancelled if the user calls
	// Close on the result.
	doneCh := make(chan struct{})
	reqCtx, reqCancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-userCtx.Done():
			if errors.Is(userCtx.Err(), context.DeadlineExceeded) {
				atomic.StoreUint32(&cancellationIsTimeout, 1)
			}
			reqCancel()
		case <-timeoutCtx.Done():
			atomic.StoreUint32(&cancellationIsTimeout, 1)
			reqCancel()
		case <-doneCh:
		}
	}()

	var firstRows *query_v1.QueryResponse
	res, err := wrapPSOpCtxWithPeek(reqCtx, manager, req, manager.TraceSpan(), qpc.provider.Query, func(client query_v1.QueryService_QueryClient) error {
		var err error
		firstRows, err = client.Recv()
		if err != nil {
			return err
		}

		return nil
	})
	close(doneCh)
	if err != nil {
		reqCancel()
		manager.Finish(true)
		return nil, qpc.makeError(err, statement, opts.Readonly, atomic.LoadUint32(&cancellationIsTimeout) == 1,
			manager.ElapsedTime(), manager.RetryInfo())
	}

	reader := &queryProviderPsRowReader{
		cli:        res,
		cancelFunc: reqCancel,

		statement: statement,
		readOnly:  opts.Readonly,

		nextRows: firstRows.Rows,
		meta:     firstRows.MetaData,

		manager: manager,
	}
	return newQueryResult(reader), nil
}

func (qpc *queryProviderPs) makeError(err error, statement string, readonly, hasTimedOut bool, elapsed time.Duration,
	retryInfo retriedRequestInfo) error {
	var gocbErr *GenericError
	if !errors.As(err, &gocbErr) {
		return err
	}

	if errors.Is(err, ErrRequestCanceled) && hasTimedOut {
		var innerErr error
		if readonly {
			innerErr = ErrUnambiguousTimeout
		} else {
			innerErr = ErrAmbiguousTimeout
		}

		return &TimeoutError{
			InnerError:    innerErr,
			TimeObserved:  elapsed,
			OperationID:   retryInfo.Operation(),
			Opaque:        retryInfo.Identifier(),
			RetryReasons:  retryInfo.RetryReasons(),
			RetryAttempts: retryInfo.RetryAttempts(),
		}
	}

	if gocbErr.Context == nil {
		gocbErr.Context = make(map[string]interface{})
	}
	gocbErr.Context["statement"] = statement

	return gocbErr
}

type queryProviderPsRowReader struct {
	cli        query_v1.QueryService_QueryClient
	cancelFunc context.CancelFunc

	statement string
	readOnly  bool

	nextRowsIndex int
	nextRows      [][]byte
	err           error
	meta          *query_v1.QueryResponse_MetaData

	manager *psOpManagerDefault
}

func (q *queryProviderPsRowReader) NextRow() []byte {
	if q.nextRowsIndex < len(q.nextRows) {
		row := q.nextRows[q.nextRowsIndex]
		q.nextRowsIndex++
		return row
	}

	res, err := q.cli.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			q.finishWithoutError()
			return nil
		}
		q.finishWithError(err)
		return nil
	}

	q.nextRows = res.Rows
	q.nextRowsIndex = 1
	q.meta = res.MetaData

	if len(res.Rows) > 0 {
		return res.Rows[0]
	}

	return nil
}

func (q *queryProviderPsRowReader) Err() error {
	if q.err == nil {
		return nil
	}

	gocbErr := mapPsErrorToGocbError(q.err, q.readOnly)
	if gocbErr.Context == nil {
		gocbErr.Context = make(map[string]interface{})
	}
	gocbErr.Context["statement"] = q.statement

	return gocbErr
}

func (q *queryProviderPsRowReader) MetaData() ([]byte, error) {
	if err := q.Err(); err != nil {
		return nil, err
	}
	if q.cli != nil {
		return nil, errors.New("the result must be fully read before accessing the meta-data")
	}
	if q.meta == nil {
		return nil, errors.New("an error occurred during querying which has made the meta-data unavailable")
	}

	meta := jsonQueryResponse{
		RequestID:       q.meta.RequestId,
		ClientContextID: q.meta.ClientContextId,
		Profile:         json.RawMessage(q.meta.Profile),
		Signature:       json.RawMessage(q.meta.Signature),
	}
	switch q.meta.Status {
	case query_v1.QueryResponse_MetaData_STATUS_RUNNING:
		meta.Status = QueryStatusRunning
	case query_v1.QueryResponse_MetaData_STATUS_SUCCESS:
		meta.Status = QueryStatusSuccess
	case query_v1.QueryResponse_MetaData_STATUS_ERRORS:
		meta.Status = QueryStatusErrors
	case query_v1.QueryResponse_MetaData_STATUS_COMPLETED:
		meta.Status = QueryStatusCompleted
	case query_v1.QueryResponse_MetaData_STATUS_STOPPED:
		meta.Status = QueryStatusStopped
	case query_v1.QueryResponse_MetaData_STATUS_TIMEOUT:
		meta.Status = QueryStatusTimeout
	case query_v1.QueryResponse_MetaData_STATUS_CLOSED:
		meta.Status = QueryStatusClosed
	case query_v1.QueryResponse_MetaData_STATUS_FATAL:
		meta.Status = QueryStatusFatal
	case query_v1.QueryResponse_MetaData_STATUS_ABORTED:
		meta.Status = QueryStatusAborted
	case query_v1.QueryResponse_MetaData_STATUS_UNKNOWN:
		meta.Status = QueryStatusUnknown
	default:
		meta.Status = QueryStatusUnknown
	}

	if len(q.meta.Warnings) > 0 {
		meta.Warnings = make([]jsonQueryWarning, len(q.meta.Warnings))
		for i, warning := range q.meta.Warnings {
			meta.Warnings[i] = jsonQueryWarning{
				Code:    warning.Code,
				Message: warning.Message,
			}
		}
	}

	if q.meta.Metrics != nil {
		meta.Metrics = &jsonQueryMetrics{
			ElapsedTime:   q.meta.Metrics.ElapsedTime.AsDuration().String(),
			ExecutionTime: q.meta.Metrics.ExecutionTime.AsDuration().String(),
			ResultCount:   q.meta.Metrics.ResultCount,
			ResultSize:    q.meta.Metrics.ResultSize,
			MutationCount: q.meta.Metrics.MutationCount,
			SortCount:     q.meta.Metrics.SortCount,
			ErrorCount:    q.meta.Metrics.ErrorCount,
			WarningCount:  q.meta.Metrics.WarningCount,
		}
	}

	return json.Marshal(meta)
}

func (q *queryProviderPsRowReader) Close() error {
	if q.err != nil {
		return q.err
	}
	// if the client is nil then we must be closed already.
	if q.cli == nil {
		return nil
	}
	q.cancelFunc()
	err := q.cli.CloseSend()

	q.manager.Finish(true)

	q.cli = nil
	return err
}

func (q *queryProviderPsRowReader) PreparedName() (string, error) {
	return "", nil
}

func (q *queryProviderPsRowReader) Endpoint() string {
	return ""
}

func (r *queryProviderPsRowReader) finishWithoutError() {
	r.cancelFunc()
	// Close the stream now that we are done with it
	err := r.cli.CloseSend()
	if err != nil {
		logWarnf("query stream close failed after meta-data: %s", err)
	}

	r.manager.Finish(true)

	r.cli = nil
}

func (r *queryProviderPsRowReader) finishWithError(err error) {
	// Lets record the error that happened
	r.err = err
	r.cancelFunc()

	// Lets Close the underlying stream
	closeErr := r.cli.CloseSend()
	if closeErr != nil {
		// We log this at debug level, but its almost always going to be an
		// error since thats the most likely reason we are in finishWithError
		logDebugf("query stream close failed after error: %s", closeErr)
	}

	r.manager.Finish(true)

	// Our client is invalidated as soon as an error occurs
	r.cli = nil
}
