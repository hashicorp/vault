package gocbcore

import (
	"encoding/binary"
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
)

type crudComponent struct {
	cidMgr               *collectionsComponent
	defaultRetryStrategy RetryStrategy
	tracer               *tracerComponent
	errMapManager        *errMapComponent
}

func newCRUDComponent(cidMgr *collectionsComponent, defaultRetryStrategy RetryStrategy, tracerCmpt *tracerComponent,
	errMapManager *errMapComponent) *crudComponent {
	return &crudComponent{
		cidMgr:               cidMgr,
		defaultRetryStrategy: defaultRetryStrategy,
		tracer:               tracerCmpt,
		errMapManager:        errMapManager,
	}
}

func (crud *crudComponent) Get(opts GetOptions, cb GetCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("Get", opts.TraceContext)

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

		tracer.Finish()
		cb(&res, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdGet,
			Datatype:     0,
			Cas:          0,
			Extras:       nil,
			Key:          opts.Key,
			Value:        nil,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errUnambiguousTimeout,
				OperationID:        "Get",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) GetAndTouch(opts GetAndTouchOptions, cb GetAndTouchCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("GetAndTouch", opts.TraceContext)

	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
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

		tracer.Finish()
		cb(&GetAndTouchResult{
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	extraBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Expiry)

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdGAT,
			Datatype:     0,
			Cas:          0,
			Extras:       extraBuf,
			Key:          opts.Key,
			Value:        nil,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        "GetAndTouch",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) GetAndLock(opts GetAndLockOptions, cb GetAndLockCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("GetAndLock", opts.TraceContext)

	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
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

		tracer.Finish()
		cb(&GetAndLockResult{
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	extraBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extraBuf[0:], opts.LockTime)

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdGetLocked,
			Datatype:     0,
			Cas:          0,
			Extras:       extraBuf,
			Key:          opts.Key,
			Value:        nil,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        "GetAndLock",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) GetOneReplica(opts GetOneReplicaOptions, cb GetReplicaCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("GetOneReplica", opts.TraceContext)

	if opts.ReplicaIdx <= 0 {
		tracer.Finish()
		return nil, errInvalidReplica
	}

	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
		if err != nil {
			cb(nil, err)
			return
		}

		if len(resp.Extras) != 4 {
			cb(nil, errProtocol)
			return
		}

		flags := binary.BigEndian.Uint32(resp.Extras[0:])

		cb(&GetReplicaResult{
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdGetReplica,
			Datatype:     0,
			Cas:          0,
			Extras:       nil,
			Key:          opts.Key,
			Value:        nil,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		ReplicaIdx:       opts.ReplicaIdx,
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errUnambiguousTimeout,
				OperationID:        "GetOneReplica",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) Touch(opts TouchOptions, cb TouchCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("Touch", opts.TraceContext)

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

		tracer.Finish()
		cb(&TouchResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
	}

	extraBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Expiry)

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdTouch,
			Datatype:     0,
			Cas:          0,
			Extras:       extraBuf,
			Key:          opts.Key,
			Value:        nil,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        "Touch",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) Unlock(opts UnlockOptions, cb UnlockCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("Unlock", opts.TraceContext)

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

		tracer.Finish()
		cb(&UnlockResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdUnlockKey,
			Datatype:     0,
			Cas:          uint64(opts.Cas),
			Extras:       nil,
			Key:          opts.Key,
			Value:        nil,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        "Unlock",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) Delete(opts DeleteOptions, cb DeleteCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("Delete", opts.TraceContext)

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

		tracer.Finish()
		cb(&DeleteResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.cidMgr.HasDurabilityLevelStatus(durabilityLevelStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: memd.DurabilityLevel(opts.DurabilityLevel),
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
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
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        "Delete",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) store(opName string, opcode memd.CmdCode, opts storeOptions, cb StoreCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace(opName, opts.TraceContext)

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

		tracer.Finish()
		cb(&StoreResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.cidMgr.HasDurabilityLevelStatus(durabilityLevelStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: memd.DurabilityLevel(opts.DurabilityLevel),
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
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
			CollectionID:           opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        opName,
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
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
	}, cb)
}

func (crud *crudComponent) Replace(opts ReplaceOptions, cb StoreCallback) (PendingOp, error) {
	return crud.store("Replace", memd.CmdReplace, storeOptions{
		Key:                    opts.Key,
		CollectionName:         opts.CollectionName,
		ScopeName:              opts.ScopeName,
		RetryStrategy:          opts.RetryStrategy,
		Value:                  opts.Value,
		Flags:                  opts.Flags,
		Datatype:               opts.Datatype,
		Cas:                    opts.Cas,
		Expiry:                 opts.Expiry,
		TraceContext:           opts.TraceContext,
		DurabilityLevel:        opts.DurabilityLevel,
		DurabilityLevelTimeout: opts.DurabilityLevelTimeout,
		CollectionID:           opts.CollectionID,
		Deadline:               opts.Deadline,
	}, cb)
}

func (crud *crudComponent) adjoin(opName string, opcode memd.CmdCode, opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace(opName, opts.TraceContext)

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

		tracer.Finish()
		cb(&AdjoinResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.cidMgr.HasDurabilityLevelStatus(durabilityLevelStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: memd.DurabilityLevel(opts.DurabilityLevel),
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
		}
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
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        opName,
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) Append(opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	return crud.adjoin("Append", memd.CmdAppend, opts, cb)
}

func (crud *crudComponent) Prepend(opts AdjoinOptions, cb AdjoinCallback) (PendingOp, error) {
	return crud.adjoin("Prepend", memd.CmdPrepend, opts, cb)
}

func (crud *crudComponent) counter(opName string, opcode memd.CmdCode, opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace(opName, opts.TraceContext)

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

		tracer.Finish()
		cb(&CounterResult{
			Value:         intVal,
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
	}

	// You cannot have an expiry when you do not want to create the document.
	if opts.Initial == uint64(0xFFFFFFFFFFFFFFFF) && opts.Expiry != 0 {
		return nil, errInvalidArgument
	}

	var duraLevelFrame *memd.DurabilityLevelFrame
	var duraTimeoutFrame *memd.DurabilityTimeoutFrame
	if opts.DurabilityLevel > 0 {
		if crud.cidMgr.HasDurabilityLevelStatus(durabilityLevelStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
		duraLevelFrame = &memd.DurabilityLevelFrame{
			DurabilityLevel: memd.DurabilityLevel(opts.DurabilityLevel),
		}
		duraTimeoutFrame = &memd.DurabilityTimeoutFrame{
			DurabilityTimeout: opts.DurabilityLevelTimeout,
		}
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
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        opName,
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) Increment(opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	return crud.counter("Increment", memd.CmdIncrement, opts, cb)
}

func (crud *crudComponent) Decrement(opts CounterOptions, cb CounterCallback) (PendingOp, error) {
	return crud.counter("Decrement", memd.CmdDecrement, opts, cb)
}

func (crud *crudComponent) GetRandom(opts GetRandomOptions, cb GetRandomCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("GetRandom", opts.TraceContext)

	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
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

		tracer.Finish()
		cb(&GetRandomResult{
			Key:      resp.Key,
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Datatype: resp.Datatype,
		}, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdGetRandom,
			Datatype: 0,
			Cas:      0,
			Extras:   nil,
			Key:      nil,
			Value:    nil,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errUnambiguousTimeout,
				OperationID:        "GetRandom",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) GetMeta(opts GetMetaOptions, cb GetMetaCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("GetMeta", nil)

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

		deleted := binary.BigEndian.Uint32(resp.Extras[0:])
		flags := binary.BigEndian.Uint32(resp.Extras[4:])
		expTime := binary.BigEndian.Uint32(resp.Extras[8:])
		seqNo := SeqNo(binary.BigEndian.Uint64(resp.Extras[12:]))
		dataType := resp.Extras[20]

		tracer.Finish()
		cb(&GetMetaResult{
			Value:    resp.Value,
			Flags:    flags,
			Cas:      Cas(resp.Cas),
			Expiry:   expTime,
			SeqNo:    seqNo,
			Datatype: dataType,
			Deleted:  deleted,
		}, nil)
	}

	extraBuf := make([]byte, 1)
	extraBuf[0] = 2

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdGetMeta,
			Datatype:     0,
			Cas:          0,
			Extras:       extraBuf,
			Key:          opts.Key,
			Value:        nil,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errUnambiguousTimeout,
				OperationID:        "GetMeta",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) SetMeta(opts SetMetaOptions, cb SetMetaCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("SetMeta", nil)

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

		tracer.Finish()
		cb(&SetMetaResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
	}

	extraBuf := make([]byte, 30+len(opts.Extra))
	binary.BigEndian.PutUint32(extraBuf[0:], opts.Flags)
	binary.BigEndian.PutUint32(extraBuf[4:], opts.Expiry)
	binary.BigEndian.PutUint64(extraBuf[8:], uint64(opts.RevNo))
	binary.BigEndian.PutUint64(extraBuf[16:], uint64(opts.Cas))
	binary.BigEndian.PutUint32(extraBuf[24:], opts.Options)
	binary.BigEndian.PutUint16(extraBuf[28:], uint16(len(opts.Extra)))
	copy(extraBuf[30:], opts.Extra)

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdSetMeta,
			Datatype:     opts.Datatype,
			Cas:          0,
			Extras:       extraBuf,
			Key:          opts.Key,
			Value:        opts.Value,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        "SetMeta",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}

func (crud *crudComponent) DeleteMeta(opts DeleteMetaOptions, cb DeleteMetaCallback) (PendingOp, error) {
	tracer := crud.tracer.CreateOpTrace("DeleteMeta", nil)

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

		tracer.Finish()
		cb(&DeleteMetaResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
		}, nil)
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
			Magic:        memd.CmdMagicReq,
			Command:      memd.CmdDelMeta,
			Datatype:     opts.Datatype,
			Cas:          0,
			Extras:       extraBuf,
			Key:          opts.Key,
			Value:        opts.Value,
			CollectionID: opts.CollectionID,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.Timer = time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallback(&TimeoutError{
				InnerError:         errAmbiguousTimeout,
				OperationID:        "DeleteMeta",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Now().Sub(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			})
		})
	}

	return crud.cidMgr.Dispatch(req)
}
