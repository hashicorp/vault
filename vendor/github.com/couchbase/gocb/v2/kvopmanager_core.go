package gocb

import (
	"context"
	"errors"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
	"github.com/couchbase/gocbcore/v10/memd"
)

// Contains information only useful to gocbcore
type kvOpManagerCore struct {
	parent *Collection
	kv     *kvProviderCore
	signal chan struct{}

	err           error
	wasResolved   bool
	mutationToken *MutationToken

	span            RequestSpan
	documentID      string
	transcoder      Transcoder
	timeout         time.Duration
	deadline        time.Time
	bytes           []byte
	flags           uint32
	persistTo       uint
	replicateTo     uint
	durabilityLevel memd.DurabilityLevel
	retryStrategy   *coreRetryStrategyWrapper
	cancelCh        chan struct{}
	impersonate     string

	operationName string
	createdTime   time.Time
	meter         *meterWrapper
	preserveTTL   bool

	ctx context.Context
}

func (m *kvOpManagerCore) getTimeout() time.Duration {
	if m.timeout > 0 {
		if m.durabilityLevel > 0 && m.timeout < durabilityTimeoutFloor {
			m.timeout = durabilityTimeoutFloor
			logWarnf("Durable operation in use so timeout value coerced up to %s", m.timeout.String())
		}
		return m.timeout
	}

	defaultTimeout := m.parent.timeoutsConfig.KVTimeout
	if m.durabilityLevel > memd.DurabilityLevelMajority || m.persistTo > 0 {
		defaultTimeout = m.parent.timeoutsConfig.KVDurableTimeout
	}

	if m.durabilityLevel > 0 && defaultTimeout < durabilityTimeoutFloor {
		defaultTimeout = durabilityTimeoutFloor
		logWarnf("Durable operation in user so timeout value coerced up to %s", defaultTimeout.String())
	}

	return defaultTimeout
}

func (m *kvOpManagerCore) SetDocumentID(id string) {
	m.documentID = id
}

func (m *kvOpManagerCore) SetCancelCh(cancelCh chan struct{}) {
	m.cancelCh = cancelCh
}

func (m *kvOpManagerCore) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func (m *kvOpManagerCore) SetTranscoder(transcoder Transcoder) {
	if transcoder == nil {
		transcoder = m.parent.transcoder
	}
	m.transcoder = transcoder
}

func (m *kvOpManagerCore) SetValue(val interface{}) {
	if m.err != nil {
		return
	}
	if m.transcoder == nil {
		m.err = errors.New("expected a transcoder to be specified first")
		return
	}

	espan := m.kv.StartKvOpTrace(m.parent, "request_encoding", m.span.Context(), true)
	defer espan.End()

	bytes, flags, err := m.transcoder.Encode(val)
	if err != nil {
		m.err = err
		return
	}

	m.bytes = bytes
	m.flags = flags
}

func (m *kvOpManagerCore) SetDuraOptions(persistTo, replicateTo uint, level DurabilityLevel) {
	if persistTo != 0 || replicateTo != 0 {
		if !m.parent.useMutationTokens {
			m.err = makeInvalidArgumentsError("cannot use observe based durability without mutation tokens")
			return
		}

		if level > 0 {
			m.err = makeInvalidArgumentsError("cannot mix observe based durability and synchronous durability")
			return
		}
	}

	if level == DurabilityLevelUnknown {
		level = DurabilityLevelNone
	}

	m.persistTo = persistTo
	m.replicateTo = replicateTo
	durabilityLevel, err := level.toMemd()
	if err != nil {
		m.err = err
		return
	}

	m.durabilityLevel = durabilityLevel

	if level > DurabilityLevelNone {
		levelStr, err := level.toManagementAPI()
		if err != nil {
			logDebugf("Could not convert durability level to string: %v", err)
			return
		}
		m.span.SetAttribute(spanAttribDBDurability, levelStr)
	}
}

func (m *kvOpManagerCore) SetRetryStrategy(retryStrategy RetryStrategy) {
	wrapper := m.parent.retryStrategyWrapper
	if retryStrategy != nil {
		wrapper = newCoreRetryStrategyWrapper(retryStrategy)
	}
	m.retryStrategy = wrapper
}

func (m *kvOpManagerCore) SetImpersonate(user string) {
	m.impersonate = user
}

func (m *kvOpManagerCore) SetContext(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	m.ctx = ctx
}

func (m *kvOpManagerCore) SetPreserveExpiry(preserveTTL bool) {
	m.preserveTTL = preserveTTL
}

func (m *kvOpManagerCore) Finish(noMetrics bool) {
	m.span.End()

	if !noMetrics {
		m.meter.ValueRecord(meterValueServiceKV, m.operationName, m.createdTime)
	}
}

func (m *kvOpManagerCore) TraceSpanContext() RequestSpanContext {
	return m.span.Context()
}

func (m *kvOpManagerCore) TraceSpan() RequestSpan {
	return m.span
}

func (m *kvOpManagerCore) DocumentID() []byte {
	return []byte(m.documentID)
}

func (m *kvOpManagerCore) CollectionName() string {
	return m.parent.name()
}

func (m *kvOpManagerCore) ScopeName() string {
	return m.parent.ScopeName()
}

func (m *kvOpManagerCore) BucketName() string {
	return m.parent.bucketName()
}

func (m *kvOpManagerCore) ValueBytes() []byte {
	return m.bytes
}

func (m *kvOpManagerCore) ValueFlags() uint32 {
	return m.flags
}

func (m *kvOpManagerCore) Transcoder() Transcoder {
	return m.transcoder
}

func (m *kvOpManagerCore) DurabilityLevel() memd.DurabilityLevel {
	return m.durabilityLevel
}

func (m *kvOpManagerCore) DurabilityTimeout() time.Duration {
	if m.durabilityLevel == 0 {
		return 0
	}

	timeout := m.getTimeout()

	duraTimeout := time.Duration(float64(timeout) * 0.9)

	if duraTimeout < durabilityTimeoutFloor {
		duraTimeout = durabilityTimeoutFloor
	}

	return duraTimeout
}

func (m *kvOpManagerCore) Deadline() time.Time {
	if m.deadline.IsZero() {
		timeout := m.getTimeout()
		m.deadline = time.Now().Add(timeout)
	}

	return m.deadline
}

func (m *kvOpManagerCore) RetryStrategy() *coreRetryStrategyWrapper {
	return m.retryStrategy
}

func (m *kvOpManagerCore) Impersonate() string {
	return m.impersonate
}

func (m *kvOpManagerCore) PreserveExpiry() bool {
	return m.preserveTTL
}

func (m *kvOpManagerCore) CheckReadyForOp() error {
	if m.err != nil {
		return m.err
	}

	if m.getTimeout() == 0 {
		return errors.New("op manager had no timeout specified")
	}

	return nil
}

func (m *kvOpManagerCore) NeedsObserve() bool {
	return m.persistTo > 0 || m.replicateTo > 0
}

func (m *kvOpManagerCore) EnhanceErr(err error) error {
	return maybeEnhanceCollKVErr(err, m.parent, m.documentID)
}

func (m *kvOpManagerCore) EnhanceMt(token gocbcore.MutationToken) *MutationToken {
	if token.VbUUID != 0 {
		return &MutationToken{
			token:      token,
			bucketName: m.BucketName(),
		}
	}

	return nil
}

func (m *kvOpManagerCore) Reject() {
	m.signal <- struct{}{}
}

func (m *kvOpManagerCore) Resolve(token *MutationToken) {
	m.wasResolved = true
	m.mutationToken = token
	m.signal <- struct{}{}
}

func (m *kvOpManagerCore) Wait(op gocbcore.PendingOp, err error) error {
	if err != nil {
		return err
	}
	if m.err != nil {
		op.Cancel()
	}

	select {
	case <-m.signal:
		// Good to go
	case <-m.cancelCh:
		op.Cancel()
		<-m.signal
	case <-m.ctx.Done():
		op.Cancel()
		<-m.signal
	}

	if m.wasResolved && (m.persistTo > 0 || m.replicateTo > 0) {
		if m.mutationToken == nil {
			return errors.New("expected a mutation token")
		}

		return m.kv.waitForDurability(
			m.ctx,
			m.parent,
			m.span,
			m.documentID,
			m.mutationToken.token,
			m.replicateTo,
			m.persistTo,
			m.Deadline(),
			m.cancelCh,
			m.impersonate,
		)
	}

	return nil
}

func newKvOpManagerCore(c *Collection, opName string, parentSpan RequestSpan, kv *kvProviderCore) *kvOpManagerCore {
	var tracectx RequestSpanContext
	if parentSpan != nil {
		tracectx = parentSpan.Context()
	}

	span := kv.StartKvOpTrace(c, opName, tracectx, false)

	return &kvOpManagerCore{
		parent:        c,
		signal:        make(chan struct{}, 1),
		span:          span,
		operationName: opName,
		createdTime:   time.Now(),
		meter:         kv.meter,
		kv:            kv,
	}
}

func durationToExpiry(dura time.Duration) uint32 {
	// If the duration is 0, that indicates never-expires
	if dura == 0 {
		return 0
	}

	// If the duration is less than one second, we must force the
	// value to 1 to avoid accidentally making it never expire.
	if dura < 1*time.Second {
		return 1
	}

	if dura < 30*24*time.Hour {
		// Translate into a uint32 in seconds.
		return uint32(dura / time.Second)
	}

	// Send the duration as a unix timestamp of now plus duration.
	return uint32(time.Now().Add(dura).Unix())
}
