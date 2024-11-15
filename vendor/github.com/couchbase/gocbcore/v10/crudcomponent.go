package gocbcore

import (
	"encoding/binary"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type crudComponent struct {
	cidMgr               *collectionsComponent
	defaultRetryStrategy RetryStrategy
	tracer               *tracerComponent
	errMapManager        *errMapComponent
	featureVerifier      bucketCapabilityVerifier
	clientProvider       clientProvider
	disableDecompression bool
}

func newCRUDComponent(cidMgr *collectionsComponent, defaultRetryStrategy RetryStrategy, tracerCmpt *tracerComponent,
	errMapManager *errMapComponent, featureVerifier bucketCapabilityVerifier, clientProvider clientProvider,
	disableDecompression bool) *crudComponent {
	return &crudComponent{
		cidMgr:               cidMgr,
		defaultRetryStrategy: defaultRetryStrategy,
		tracer:               tracerCmpt,
		errMapManager:        errMapManager,
		featureVerifier:      featureVerifier,
		disableDecompression: disableDecompression,
		clientProvider:       clientProvider,
	}
}

func (crud *crudComponent) Get(opts GetOptions, cb GetCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "Get", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if len(resp.Extras) != 4 {
			tracer.Finish()
			cb(nil, errProtocol)
			return
		}

		res := GetResult{}
		res.Value = resp.Value
		res.Flags = binary.BigEndian.Uint32(resp.Extras[0:])
		res.Cas = Cas(resp.Cas)
		res.Datatype = resp.Datatype
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(&res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdGet,
			Datatype:               0,
			Cas:                    0,
			Extras:                 nil,
			Key:                    opts.Key,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "Get", errUnambiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) GetAndTouch(opts GetAndTouchOptions, cb GetAndTouchCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "GetAndTouch", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if len(resp.Extras) != 4 {
			tracer.Finish()
			cb(nil, errProtocol)
			return
		}

		flags := binary.BigEndian.Uint32(resp.Extras[0:])

		res := &GetAndTouchResult{
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	extraBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Expiry)

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdGAT,
			Datatype:               0,
			Cas:                    0,
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "GetAndTouch", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) GetAndLock(opts GetAndLockOptions, cb GetAndLockCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "GetAndLock", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if len(resp.Extras) != 4 {
			tracer.Finish()
			cb(nil, errProtocol)
			return
		}

		flags := binary.BigEndian.Uint32(resp.Extras[0:])
		res := &GetAndLockResult{
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	extraBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extraBuf[0:], opts.LockTime)

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdGetLocked,
			Datatype:               0,
			Cas:                    0,
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "GetAndLock", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) GetOneReplica(opts GetOneReplicaOptions, cb GetReplicaCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "GetOneReplica", opts.TraceContext)

	if opts.ReplicaIdx <= 0 {
		tracer.Finish()
		return nil, errInvalidReplica
	}

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if len(resp.Extras) != 4 {
			tracer.Finish()
			cb(nil, errProtocol)
			return
		}

		flags := binary.BigEndian.Uint32(resp.Extras[0:])
		res := &GetReplicaResult{
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdGetReplica,
			Datatype:               0,
			Cas:                    0,
			Extras:                 nil,
			Key:                    opts.Key,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		ReplicaIdx:       opts.ReplicaIdx,
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
		ServerGroup:      opts.ServerGroup,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "GetOneReplica", errUnambiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) Touch(opts TouchOptions, cb TouchCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "Touch", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}

		res := &TouchResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	extraBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Expiry)

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdTouch,
			Datatype:               0,
			Cas:                    0,
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "Touch", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) Unlock(opts UnlockOptions, cb UnlockCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "Unlock", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}

		res := &UnlockResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdUnlockKey,
			Datatype:               0,
			Cas:                    uint64(opts.Cas),
			Extras:                 nil,
			Key:                    opts.Key,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "Unlock", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) Delete(opts DeleteOptions, cb DeleteCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "Delete", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}

		res := &DeleteResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.featureVerifier.HasBucketCapabilityStatus(BucketCapabilityDurableWrites, CapabilityStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: opts.DurabilityLevel,
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
		}
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdDelete,
			Datatype:               0,
			Cas:                    uint64(opts.Cas),
			Extras:                 nil,
			Key:                    opts.Key,
			Value:                  nil,
			DurabilityLevelFrame:   duraLevelFrame,
			DurabilityTimeoutFrame: duraTimeoutFrame,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "Delete", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) store(opName string, opcode memd.CmdCode, opts storeOptions, cb StoreCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, opName, opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}

		res := &StoreResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.featureVerifier.HasBucketCapabilityStatus(BucketCapabilityDurableWrites, CapabilityStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: opts.DurabilityLevel,
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
		}
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	var preserveExpiryFrame *memd.PreserveExpiryFrame
	if opts.PreserveExpiry {
		preserveExpiryFrame = &memd.PreserveExpiryFrame{}
	}

	extraBuf := make([]byte, 8)
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Flags)
	binary.BigEndian.PutUint32(extraBuf[4:], opts.Expiry)
	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                opcode,
			Datatype:               opts.Datatype,
			Cas:                    uint64(opts.Cas),
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  opts.Value,
			DurabilityLevelFrame:   duraLevelFrame,
			DurabilityTimeoutFrame: duraTimeoutFrame,
			UserImpersonationFrame: userFrame,
			CollectionID:           opts.CollectionID,
			PreserveExpiryFrame:    preserveExpiryFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, opName, errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) Set(opts SetOptions, cb StoreCallback) (PendingOp, error) {
	return crud.store("Set", memd.CmdSet, storeOptions{
		Key:                    opts.Key,
		CollectionName:         opts.CollectionName,
		ScopeName:              opts.ScopeName,
		RetryStrategy:          opts.RetryStrategy,
		Value:                  opts.Value,
		Flags:                  opts.Flags,
		Datatype:               opts.Datatype,
		Cas:                    0,
		Expiry:                 opts.Expiry,
		TraceContext:           opts.TraceContext,
		DurabilityLevel:        opts.DurabilityLevel,
		DurabilityLevelTimeout: opts.DurabilityLevelTimeout,
		CollectionID:           opts.CollectionID,
		Deadline:               opts.Deadline,
		User:                   opts.User,
		PreserveExpiry:         opts.PreserveExpiry,
	}, cb)
}

func (crud *crudComponent) Add(opts AddOptions, cb StoreCallback) (PendingOp, error) {
	return crud.store("Add", memd.CmdAdd, storeOptions{
		Key:                    opts.Key,
		CollectionName:         opts.CollectionName,
		ScopeName:              opts.ScopeName,
		RetryStrategy:          opts.RetryStrategy,
		Value:                  opts.Value,
		Flags:                  opts.Flags,
		Datatype:               opts.Datatype,
		Cas:                    0,
		Expiry:                 opts.Expiry,
		TraceContext:           opts.TraceContext,
		DurabilityLevel:        opts.DurabilityLevel,
		DurabilityLevelTimeout: opts.DurabilityLevelTimeout,
		CollectionID:           opts.CollectionID,
		Deadline:               opts.Deadline,
		User:                   opts.User,
	}, cb)
}

func (crud *crudComponent) Replace(opts ReplaceOptions, cb StoreCallback) (PendingOp, error) {
	if opts.PreserveExpiry && opts.Expiry > 0 {
		return nil, wrapError(errInvalidArgument, "cannot use preserve expiry and an expiry > 0 for replace")
	}
	return crud.store("Replace", memd.CmdReplace, storeOptions(opts), cb)
}

func (crud *crudComponent) adjoin(opName string, opcode memd.CmdCode, opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, opName, opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}
		res := &AdjoinResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.featureVerifier.HasBucketCapabilityStatus(BucketCapabilityDurableWrites, CapabilityStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: opts.DurabilityLevel,
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
		}
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	var preserveExpiryFrame *memd.PreserveExpiryFrame
	if opts.PreserveExpiry {
		preserveExpiryFrame = &memd.PreserveExpiryFrame{}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                opcode,
			Datatype:               0,
			Cas:                    uint64(opts.Cas),
			Extras:                 nil,
			Key:                    opts.Key,
			Value:                  opts.Value,
			DurabilityLevelFrame:   duraLevelFrame,
			DurabilityTimeoutFrame: duraTimeoutFrame,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
			PreserveExpiryFrame:    preserveExpiryFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, opName, errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) Append(opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	return crud.adjoin("Append", memd.CmdAppend, opts, cb)
}

func (crud *crudComponent) Prepend(opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	return crud.adjoin("Prepend", memd.CmdPrepend, opts, cb)
}

func (crud *crudComponent) counter(opName string, opcode memd.CmdCode, opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, opName, opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if len(resp.Value) != 8 {
			tracer.Finish()
			cb(nil, errProtocol)
			return
		}
		intVal := binary.BigEndian.Uint64(resp.Value)

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}
		res := &CounterResult{
			Value:         intVal,
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	// You cannot have an expiry when you do not want to create the document.
	if opts.Initial == uint64(0xFFFFFFFFFFFFFFFF) && opts.Expiry != 0 {
		return nil, errInvalidArgument
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.featureVerifier.HasBucketCapabilityStatus(BucketCapabilityDurableWrites, CapabilityStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: opts.DurabilityLevel,
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
		}
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}
	var preserveExpiryFrame *memd.PreserveExpiryFrame
	if opts.PreserveExpiry {
		preserveExpiryFrame = &memd.PreserveExpiryFrame{}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	extraBuf := make([]byte, 20)
	binary.BigEndian.PutUint64(extraBuf[0:], opts.Delta)
	if opts.Initial != uint64(0xFFFFFFFFFFFFFFFF) {
		binary.BigEndian.PutUint64(extraBuf[8:], opts.Initial)
		binary.BigEndian.PutUint32(extraBuf[16:], opts.Expiry)
	} else {
		binary.BigEndian.PutUint64(extraBuf[8:], 0x0000000000000000)
		binary.BigEndian.PutUint32(extraBuf[16:], 0xFFFFFFFF)
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                opcode,
			Datatype:               0,
			Cas:                    uint64(opts.Cas),
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  nil,
			DurabilityLevelFrame:   duraLevelFrame,
			DurabilityTimeoutFrame: duraTimeoutFrame,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
			PreserveExpiryFrame:    preserveExpiryFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, opName, errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) Increment(opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	return crud.counter("Increment", memd.CmdIncrement, opts, cb)
}

func (crud *crudComponent) Decrement(opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	return crud.counter("Decrement", memd.CmdDecrement, opts, cb)
}

func (crud *crudComponent) GetRandom(opts GetRandomOptions, cb GetRandomCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "GetRandom", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if len(resp.Extras) != 4 {
			tracer.Finish()
			cb(nil, errProtocol)
			return
		}

		flags := binary.BigEndian.Uint32(resp.Extras[0:])
		res := &GetRandomResult{
			Key:      resp.Key,
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdGetRandom,
			Datatype:               0,
			Cas:                    0,
			Extras:                 nil,
			Key:                    nil,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		RetryStrategy:    opts.RetryStrategy,
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "GetRandom", errUnambiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) GetMeta(opts GetMetaOptions, cb GetMetaCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "GetMeta", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if len(resp.Extras) != 21 {
			tracer.Finish()
			cb(nil, errProtocol)
			return
		}

		res := &GetMetaResult{
			Value: resp.Value,
			Cas:   Cas(resp.Cas),
		}
		res.Deleted = binary.BigEndian.Uint32(resp.Extras[0:])
		res.Flags = binary.BigEndian.Uint32(resp.Extras[4:])
		res.Expiry = binary.BigEndian.Uint32(resp.Extras[8:])
		res.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[12:]))
		res.Datatype = resp.Extras[20]
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	extraBuf := make([]byte, 1)
	extraBuf[0] = 2

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdGetMeta,
			Datatype:               0,
			Cas:                    0,
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  nil,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "GetMeta", errUnambiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) SetMeta(opts SetMetaOptions, cb SetMetaCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "SetMeta", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}
		res := &SetMetaResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	extraBuf := make([]byte, 30+len(opts.Extra))
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Flags)
	binary.BigEndian.PutUint32(extraBuf[4:], opts.Expiry)
	binary.BigEndian.PutUint64(extraBuf[8:], opts.RevNo)
	binary.BigEndian.PutUint64(extraBuf[16:], uint64(opts.Cas))
	binary.BigEndian.PutUint32(extraBuf[24:], opts.Options)
	binary.BigEndian.PutUint16(extraBuf[28:], uint16(len(opts.Extra)))
	copy(extraBuf[30:], opts.Extra)

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdSetMeta,
			Datatype:               opts.Datatype,
			Cas:                    0,
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  opts.Value,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "SetMeta", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) DeleteMeta(opts DeleteMetaOptions, cb DeleteMetaCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "DeleteMeta", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}
		res := &DeleteMetaResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	extraBuf := make([]byte, 30+len(opts.Extra))
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Flags)
	binary.BigEndian.PutUint32(extraBuf[4:], opts.Expiry)
	binary.BigEndian.PutUint64(extraBuf[8:], opts.RevNo)
	binary.BigEndian.PutUint64(extraBuf[16:], uint64(opts.Cas))
	binary.BigEndian.PutUint32(extraBuf[24:], opts.Options)
	binary.BigEndian.PutUint16(extraBuf[28:], uint16(len(opts.Extra)))
	copy(extraBuf[30:], opts.Extra)

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdDelMeta,
			Datatype:               opts.Datatype,
			Cas:                    0,
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  opts.Value,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	op, err := crud.cidMgr.Dispatch(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			req.cancelWithCallbackAndFinishTracer(
				makeTimeoutError(start, "DeleteMeta", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}
