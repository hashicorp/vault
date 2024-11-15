package gocb

import (
	"context"
	"errors"
	"fmt"
	"github.com/couchbaselabs/gocbconnstr/v2"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"github.com/couchbase/gocbcore/v10"
	"github.com/couchbase/gocbcoreps"
)

type psConnectionMgr struct {
	host   string
	lock   sync.Mutex
	config *gocbcoreps.DialOptions
	agent  *gocbcoreps.RoutingClient

	timeouts     TimeoutsConfig
	tracer       RequestTracer
	meter        *meterWrapper
	defaultRetry RetryStrategy

	closed      atomic.Bool
	activeOpsWg sync.WaitGroup
}

func (c *psConnectionMgr) connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	client, err := gocbcoreps.Dial(c.host, c.config)
	if err != nil {
		return err
	}

	c.agent = client

	return nil
}

func (c *psConnectionMgr) openBucket(bucketName string) error {
	return nil
}

func (c *psConnectionMgr) buildConfig(cluster *Cluster) error {
	resolved, err := gocbconnstr.Resolve(cluster.connSpec())
	if err != nil {
		return err
	}

	c.host = fmt.Sprintf("%s:%d", resolved.Couchbase2Host.Host, resolved.Couchbase2Host.Port)

	creds, err := cluster.authenticator().Credentials(AuthCredsRequest{})
	if err != nil {
		return err
	}

	logger := newZapLogger()

	var tp trace.TracerProvider
	if c.tracer != nil {
		if tracer, ok := c.tracer.(OtelAwareRequestTracer); ok {
			tp = tracer.Provider()
		}
	}
	var mp metric.MeterProvider
	if c.meter != nil {
		if meter, ok := c.meter.meter.(OtelAwareMeter); ok {
			mp = meter.Provider()
		}
	}

	c.config = &gocbcoreps.DialOptions{
		Username:           creds[0].Username,
		Password:           creds[0].Password,
		InsecureSkipVerify: cluster.securityConfig.TLSSkipVerify,
		ClientCertificate:  cluster.securityConfig.TLSRootCAs,
		Logger:             logger,
		TracerProvider:     tp,
		MeterProvider:      mp,
	}

	return nil
}

func (c *psConnectionMgr) canPerformOp() error {
	if c.closed.Load() {
		return ErrShutdown
	}

	return nil
}

func (c *psConnectionMgr) MarkOpBeginning() {
	c.activeOpsWg.Add(1)
}

func (c *psConnectionMgr) MarkOpCompleted() {
	c.activeOpsWg.Done()
}

func (c *psConnectionMgr) getKvProvider(bucketName string) (kvProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	kv := c.agent.KvV1()
	return &kvProviderPs{
		client: kv,

		tracer: c.tracer,
		meter:  c.meter,
	}, nil
}

func (c *psConnectionMgr) getKvBulkProvider(bucketName string) (kvBulkProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	kv := c.agent.KvV1()
	return &kvBulkProviderPs{
		client: kv,

		tracer: c.tracer,
		meter:  c.meter,
	}, nil
}

func (c *psConnectionMgr) getKvCapabilitiesProvider(bucketName string) (kvCapabilityVerifier, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getViewProvider(bucketName string) (viewProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getViewIndexProvider(bucketName string) (viewIndexProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getQueryProvider() (queryProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	provider := c.agent.QueryV1()
	return &queryProviderPs{
		provider: provider,

		managerProvider: newPsOpManagerProvider(c.defaultRetry, c.tracer, c.timeouts.QueryTimeout, c.meter, meterValueServiceQuery),
	}, nil
}

func (c *psConnectionMgr) getQueryIndexProvider() (queryIndexProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	provider := c.agent.QueryAdminV1()
	return &queryIndexProviderPs{
		provider: provider,

		managerProvider: newPsOpManagerProvider(c.defaultRetry, c.tracer, c.timeouts.QueryTimeout, c.meter, meterValueServiceManagement),
	}, nil
}

func (c *psConnectionMgr) getSearchIndexProvider() (searchIndexProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	provider := c.agent.SearchAdminV1()
	return &searchIndexProviderPs{
		provider: provider,

		managerProvider: newPsOpManagerProvider(c.defaultRetry, c.tracer, c.timeouts.QueryTimeout, c.meter, meterValueServiceManagement),
	}, nil
}

func (c *psConnectionMgr) getSearchCapabilitiesProvider() (searchCapabilityVerifier, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getCollectionsManagementProvider(bucketName string) (collectionsManagementProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	return &collectionsManagementProviderPs{
		provider:   c.agent.CollectionV1(),
		bucketName: bucketName,

		managerProvider: newPsOpManagerProvider(c.defaultRetry, c.tracer, c.timeouts.QueryTimeout, c.meter, meterValueServiceManagement),
	}, nil
}

func (c *psConnectionMgr) getBucketManagementProvider() (bucketManagementProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	return &bucketManagementProviderPs{
		provider: c.agent.BucketV1(),

		managerProvider: newPsOpManagerProvider(c.defaultRetry, c.tracer, c.timeouts.QueryTimeout, c.meter, meterValueServiceManagement),
	}, nil
}

func (c *psConnectionMgr) getAnalyticsProvider() (analyticsProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getAnalyticsIndexProvider() (analyticsIndexProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getSearchProvider() (searchProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	return &searchProviderPs{
		provider: c.agent.SearchV1(),

		managerProvider: newPsOpManagerProvider(c.defaultRetry, c.tracer, c.timeouts.QueryTimeout, c.meter, meterValueServiceSearch),
	}, nil
}

func (c *psConnectionMgr) getHTTPProvider(bucketName string) (httpProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getDiagnosticsProvider(bucketName string) (diagnosticsProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getWaitUntilReadyProvider(bucketName string) (waitUntilReadyProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	return &waitUntilReadyProviderPs{
		defaultRetryStrategy: c.defaultRetry,
		client:               c.agent,
	}, nil
}

func (c *psConnectionMgr) getEventingManagementProvider() (eventingManagementProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getUserManagerProvider() (userManagerProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) getInternalProvider() (internalProvider, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) initTransactions(config TransactionsConfig, cluster *Cluster) error {
	// We don't return feature not available here as this function gets called from initTransactions which is part of
	// cluster setup, if we returned error here then cluster setup would always fail in couchbase2 mode.
	return nil
}

func (c *psConnectionMgr) getTransactionsProvider() (transactionsProvider, error) {
	if err := c.canPerformOp(); err != nil {
		return nil, err
	}

	return &transactionsProviderPs{}, nil
}

func (c *psConnectionMgr) connection(bucketName string) (*gocbcore.Agent, error) {
	return nil, ErrFeatureNotAvailable
}

func (c *psConnectionMgr) close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return ErrShutdown
	}

	err := c.agent.Close()

	logDebugf("Waiting for any active requests to complete")
	c.activeOpsWg.Wait()

	if c.tracer != nil {
		tracerDecRef(c.tracer)
		c.tracer = nil
	}
	if c.meter != nil {
		if meter, ok := c.meter.meter.(*LoggingMeter); ok {
			meter.close()
		}
		c.meter = nil
	}

	return err
}

type psOpManagerProvider struct {
	defaultRetryStrategy RetryStrategy
	tracer               RequestTracer
	defaultTimeout       time.Duration
	meter                *meterWrapper
	service              string
}

func newPsOpManagerProvider(retry RetryStrategy, tracer RequestTracer, timeout time.Duration, meter *meterWrapper,
	service string) *psOpManagerProvider {
	return &psOpManagerProvider{
		defaultRetryStrategy: retry,
		tracer:               tracer,
		defaultTimeout:       timeout,
		meter:                meter,
		service:              service,
	}
}

func (p *psOpManagerProvider) NewManager(parentSpan RequestSpan, opName string, attribs map[string]interface{}) *psOpManagerDefault {
	m := &psOpManagerDefault{
		defaultRetryStrategy: p.defaultRetryStrategy,
		tracer:               p.tracer,
		defaultTimeout:       p.defaultTimeout,
		meter:                p.meter,
		service:              p.service,
		createdTime:          time.Now(),

		span:   createSpan(p.tracer, parentSpan, opName, p.service),
		opName: opName,
	}

	for key, value := range attribs {
		m.span.SetAttribute(key, value)
	}

	return m
}

type psOpManager interface {
	IsIdempotent() bool
	TraceSpan() RequestSpan
	OperationID() string
	RetryStrategy() RetryStrategy
	OpName() string
	CreatedAt() time.Time
	Tracer() RequestTracer
	RetryInfo() retriedRequestInfo
	SetRetryRequest(ps *retriableRequestPs)
	RetryReasonFor(error) RetryReason
	Timeout() time.Duration
	Context() context.Context
}

type psOpManagerDefault struct {
	defaultRetryStrategy RetryStrategy
	tracer               RequestTracer
	defaultTimeout       time.Duration
	meter                *meterWrapper
	service              string
	createdTime          time.Time

	span   RequestSpan
	opName string

	operationID   string
	retryStrategy RetryStrategy
	timeout       time.Duration
	ctx           context.Context
	isIdempotent  bool

	req *retriableRequestPs
}

func (m *psOpManagerDefault) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func (m *psOpManagerDefault) SetRetryStrategy(retryStrategy RetryStrategy) {
	strat := m.defaultRetryStrategy
	if retryStrategy != nil {
		strat = retryStrategy
	}
	m.retryStrategy = strat
}

func (m *psOpManagerDefault) SetContext(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	m.ctx = ctx
}

func (m *psOpManagerDefault) SetIsIdempotent(idempotent bool) {
	m.isIdempotent = idempotent
}

func (m *psOpManagerDefault) SetOperationID(id string) {
	m.operationID = id
}

func (m *psOpManagerDefault) NewSpan(name string) RequestSpan {
	return createSpan(m.tracer, m.span, name, m.service)
}

func (m *psOpManagerDefault) Finish(noMetrics bool) {
	retries := m.RetryInfo().RetryAttempts()
	m.span.SetAttribute(spanAttribRetries, retries)
	m.span.End()

	if !noMetrics {
		m.meter.ValueRecord(m.service, m.opName, m.createdTime)
	}
}

func (m *psOpManagerDefault) ValueRecord() {
	m.meter.ValueRecord(m.service, m.opName, m.createdTime)
}

func (m *psOpManagerDefault) TraceSpan() RequestSpan {
	return m.span
}

func (m *psOpManagerDefault) RetryStrategy() RetryStrategy {
	return m.retryStrategy
}

func (m *psOpManagerDefault) IsIdempotent() bool {
	return m.isIdempotent
}

func (m *psOpManagerDefault) Timeout() time.Duration {
	if m.timeout > 0 {
		return m.timeout
	}

	return m.defaultTimeout
}

func (m *psOpManagerDefault) CheckReadyForOp() error {
	if m.Timeout() == 0 {
		return errors.New("op manager had no timeout specified")
	}

	if m.span == nil {
		return errors.New("op manager had no span specified")
	}

	return nil
}

func (m *psOpManagerDefault) ElapsedTime() time.Duration {
	return time.Since(m.createdTime)
}

func (m *psOpManagerDefault) OperationID() string {
	return m.operationID
}

func (m *psOpManagerDefault) OpName() string {
	return m.opName
}

func (m *psOpManagerDefault) RetryInfo() retriedRequestInfo {
	return m.req
}

func (m *psOpManagerDefault) SetRetryRequest(req *retriableRequestPs) {
	m.req = req
}

func (m *psOpManagerDefault) Tracer() RequestTracer {
	return m.tracer
}

func (m *psOpManagerDefault) CreatedAt() time.Time {
	return m.createdTime
}

func (m *psOpManagerDefault) RetryReasonFor(err error) RetryReason {
	if errors.Is(err, ErrServiceNotAvailable) {
		return ServiceNotAvailableRetryReason
	}

	return nil
}

func (m *psOpManagerDefault) Context() context.Context {
	return m.ctx
}

func wrapPSOpCtx[ReqT any, RespT any](ctx context.Context, m psOpManager,
	req ReqT,
	fn func(context.Context, ReqT, ...grpc.CallOption) (RespT, error)) (RespT, error) {
	return wrapPSOpCtxWithPeek(ctx, m, req, m.TraceSpan(), fn, nil)
}

func wrapPSOp[ReqT any, RespT any](m psOpManager, req ReqT,
	fn func(context.Context, ReqT, ...grpc.CallOption) (RespT, error)) (RespT, error) {
	ctx, cancel := context.WithTimeout(m.Context(), m.Timeout())
	defer cancel()

	return wrapPSOpCtx(ctx, m, req, fn)
}

func wrapPSOpCtxWithPeek[ReqT any, RespT any](ctx context.Context,
	m psOpManager,
	req ReqT,
	parentSpan RequestSpan,
	fn func(context.Context, ReqT, ...grpc.CallOption) (RespT, error),
	peekResult func(RespT) error) (RespT, error) {
	retryReq := newRetriableRequestPS(m.OpName(), m.IsIdempotent(), parentSpan, m.OperationID(), m.RetryStrategy())
	m.SetRetryRequest(retryReq)

	res, err := handleRetriableRequest(ctx, m.CreatedAt(), m.Tracer(), req, retryReq, fn, m.RetryReasonFor, peekResult)
	if err != nil {
		var emptyResp RespT
		return emptyResp, err
	}

	return res, nil
}
