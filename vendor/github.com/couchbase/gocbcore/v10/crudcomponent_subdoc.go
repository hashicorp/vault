package gocbcore

import (
	"encoding/binary"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type subdocOpList struct {
	ops     []SubDocOp
	indexes []int
}

func (sol *subdocOpList) Reorder(ops []SubDocOp) {
	var xAttrOps []SubDocOp
	var xAttrIndexes []int
	var sops []SubDocOp
	var opIndexes []int
	for i, op := range ops {
		if op.Flags&memd.SubdocFlagXattrPath != 0 {
			xAttrOps = append(xAttrOps, op)
			xAttrIndexes = append(xAttrIndexes, i)
		} else {
			sops = append(sops, op)
			opIndexes = append(opIndexes, i)
		}
	}

	sol.ops = append(xAttrOps, sops...)
	sol.indexes = append(xAttrIndexes, opIndexes...)
}
func (crud *crudComponent) LookupIn(opts LookupInOptions, cb LookupInCallback) (PendingOp, error) {
	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "LookupIn", opts.TraceContext)

	results := make([]SubDocResult, len(opts.Ops))
	var subdocs subdocOpList

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil &&
			!isErrorStatus(err, memd.StatusSubDocMultiPathFailureDeleted) &&
			!isErrorStatus(err, memd.StatusSubDocSuccessDeleted) &&
			!isErrorStatus(err, memd.StatusSubDocBadMulti) {
			tracer.Finish()
			cb(nil, err)
			return
		}

		respIter := 0
		for i := range results {
			if respIter+6 > len(resp.Value) {
				tracer.Finish()
				cb(nil, errProtocol)
				return
			}

			resError := memd.StatusCode(binary.BigEndian.Uint16(resp.Value[respIter+0:]))
			resValueLen := int(binary.BigEndian.Uint32(resp.Value[respIter+2:]))

			if respIter+6+resValueLen > len(resp.Value) {
				tracer.Finish()
				cb(nil, errProtocol)
				return
			}

			if resError != memd.StatusSuccess {
				results[subdocs.indexes[i]].Err = crud.makeSubDocError(i, resError, req, resp)
			}

			results[subdocs.indexes[i]].Value = resp.Value[respIter+6 : respIter+6+resValueLen]
			respIter += 6 + resValueLen
		}
		res := &LookupInResult{
			Cas: Cas(resp.Cas),
			Ops: results,
		}
		res.Internal.IsDeleted = isErrorStatus(err, memd.StatusSubDocSuccessDeleted) ||
			isErrorStatus(err, memd.StatusSubDocMultiPathFailureDeleted)
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(res, nil)
	}

	subdocs.Reorder(opts.Ops)

	pathBytesList := make([][]byte, len(opts.Ops))
	pathBytesTotal := 0
	for i, op := range subdocs.ops {
		pathBytes := []byte(op.Path)
		pathBytesList[i] = pathBytes
		pathBytesTotal += len(pathBytes)
	}

	valueBuf := make([]byte, len(opts.Ops)*4+pathBytesTotal)

	valueIter := 0
	for i, op := range subdocs.ops {
		if op.Op != memd.SubDocOpGet && op.Op != memd.SubDocOpExists &&
			op.Op != memd.SubDocOpGetDoc && op.Op != memd.SubDocOpGetCount {
			return nil, errInvalidArgument
		}
		if op.Value != nil {
			return nil, errInvalidArgument
		}

		pathBytes := pathBytesList[i]
		pathBytesLen := len(pathBytes)

		valueBuf[valueIter+0] = uint8(op.Op)
		valueBuf[valueIter+1] = uint8(op.Flags)
		binary.BigEndian.PutUint16(valueBuf[valueIter+2:], uint16(pathBytesLen))
		copy(valueBuf[valueIter+4:], pathBytes)
		valueIter += 4 + pathBytesLen
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	var extraBuf []byte
	if opts.Flags != 0 {
		if opts.Flags&memd.SubdocDocFlagReplicaRead != 0 {
			// We can get here before support status is actually known, we'll send the request unless we know for a fact
			// that this is unsupported.
			if crud.featureVerifier.HasBucketCapabilityStatus(BucketCapabilityReplicaRead, CapabilityStatusUnsupported) {
				return nil, errFeatureNotAvailable
			}
		}

		extraBuf = append(extraBuf, uint8(opts.Flags))
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdSubDocMultiLookup,
			Datatype:               0,
			Cas:                    0,
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  valueBuf,
			CollectionID:           opts.CollectionID,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RootTraceContext: tracer.RootContext(),
		CollectionName:   opts.CollectionName,
		ScopeName:        opts.ScopeName,
		RetryStrategy:    opts.RetryStrategy,
		ReplicaIdx:       opts.ReplicaIdx,
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
				makeTimeoutError(start, "LookupIn", errUnambiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) MutateIn(opts MutateInOptions, cb MutateInCallback) (PendingOp, error) {
	if len(opts.Ops) == 0 {
		return nil, wrapError(errInvalidArgument, "at least one op must be present")
	}

	tracer := crud.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "MutateIn", opts.TraceContext)

	results := make([]SubDocResult, len(opts.Ops))
	var subdocs subdocOpList

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		// GOCBC-1356: memcached can return a NOT_STORED response when inserting a doc with sub-doc.
		if isErrorStatus(err, memd.StatusNotStored) && opts.Flags&memd.SubdocDocFlagAddDoc != 0 {
			tracer.Finish()
			cb(nil, crud.errMapManager.EnhanceKvError(errDocumentExists, resp, req))
			return
		}

		if err != nil &&
			!isErrorStatus(err, memd.StatusSubDocSuccessDeleted) &&
			!isErrorStatus(err, memd.StatusSubDocBadMulti) {
			tracer.Finish()
			cb(nil, err)
			return
		}

		if isErrorStatus(err, memd.StatusSubDocBadMulti) {
			if len(resp.Value) != 3 {
				tracer.Finish()
				cb(nil, errProtocol)
				return
			}

			opIndex := int(resp.Value[0])
			resError := memd.StatusCode(binary.BigEndian.Uint16(resp.Value[1:]))

			err := crud.makeSubDocError(opIndex, resError, req, resp)
			tracer.Finish()
			cb(nil, err)
			return
		}

		for readPos := uint32(0); readPos < uint32(len(resp.Value)); {
			opIndex := int(resp.Value[readPos+0])
			opStatus := memd.StatusCode(binary.BigEndian.Uint16(resp.Value[readPos+1:]))

			results[subdocs.indexes[opIndex]].Err = crud.makeSubDocError(opIndex, opStatus, req, resp)
			readPos += 3

			if opStatus == memd.StatusSuccess {
				valLength := binary.BigEndian.Uint32(resp.Value[readPos:])
				results[subdocs.indexes[opIndex]].Value = resp.Value[readPos+4 : readPos+4+valLength]
				readPos += 4 + valLength
			}
		}

		mutToken := MutationToken{}
		if len(resp.Extras) >= 16 {
			mutToken.VbID = req.Vbucket
			mutToken.VbUUID = VbUUID(binary.BigEndian.Uint64(resp.Extras[0:]))
			mutToken.SeqNo = SeqNo(binary.BigEndian.Uint64(resp.Extras[8:]))
		}
		res := &MutateInResult{
			Cas:           Cas(resp.Cas),
			MutationToken: mutToken,
			Ops:           results,
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
		if opts.Flags|memd.SubdocDocFlagAddDoc == 1 {
			return nil, wrapError(errInvalidArgument, "cannot use preserve expiry with add doc flags")
		}
		if opts.Expiry != 0 && opts.PreserveExpiry && opts.Flags|memd.SubdocDocFlagNone == 1 {
			return nil, wrapError(errInvalidArgument, "cannot use preserve expiry with expiry and no doc flags")
		}
		preserveExpiryFrame = &memd.PreserveExpiryFrame{}
	}

	if opts.Flags&memd.SubdocDocFlagCreateAsDeleted != 0 {
		// We can get here before support status is actually known, we'll send the request unless we know for a fact
		// that this is unsupported.
		if crud.featureVerifier.HasBucketCapabilityStatus(BucketCapabilityCreateAsDeleted, CapabilityStatusUnsupported) {
			return nil, errFeatureNotAvailable
		}
	}

	subdocs.Reorder(opts.Ops)

	pathBytesList := make([][]byte, len(opts.Ops))
	pathBytesTotal := 0
	valueBytesTotal := 0
	for i, op := range subdocs.ops {
		pathBytes := []byte(op.Path)
		pathBytesList[i] = pathBytes
		pathBytesTotal += len(pathBytes)
		valueBytesTotal += len(op.Value)
	}

	valueBuf := make([]byte, len(opts.Ops)*8+pathBytesTotal+valueBytesTotal)

	valueIter := 0
	for i, op := range subdocs.ops {
		if op.Op != memd.SubDocOpDictAdd && op.Op != memd.SubDocOpDictSet &&
			op.Op != memd.SubDocOpDelete && op.Op != memd.SubDocOpReplace &&
			op.Op != memd.SubDocOpArrayPushLast && op.Op != memd.SubDocOpArrayPushFirst &&
			op.Op != memd.SubDocOpArrayInsert && op.Op != memd.SubDocOpArrayAddUnique &&
			op.Op != memd.SubDocOpCounter && op.Op != memd.SubDocOpSetDoc &&
			op.Op != memd.SubDocOpAddDoc && op.Op != memd.SubDocOpDeleteDoc &&
			op.Op != memd.SubDocOpReplaceBodyWithXattr {
			return nil, errInvalidArgument
		}

		if op.Op == memd.SubDocOpReplaceBodyWithXattr {
			// We can get here before support status is actually known, we'll send the request unless we know for a fact
			// that this is unsupported.
			if crud.featureVerifier.HasBucketCapabilityStatus(BucketCapabilityReplaceBodyWithXattr, CapabilityStatusUnsupported) {
				return nil, errFeatureNotAvailable
			}
		}

		pathBytes := pathBytesList[i]
		pathBytesLen := len(pathBytes)
		valueBytesLen := len(op.Value)

		valueBuf[valueIter+0] = uint8(op.Op)
		valueBuf[valueIter+1] = uint8(op.Flags)
		binary.BigEndian.PutUint16(valueBuf[valueIter+2:], uint16(pathBytesLen))
		binary.BigEndian.PutUint32(valueBuf[valueIter+4:], uint32(valueBytesLen))
		copy(valueBuf[valueIter+8:], pathBytes)
		copy(valueBuf[valueIter+8+pathBytesLen:], op.Value)
		valueIter += 8 + pathBytesLen + valueBytesLen
	}

	var extraBuf []byte
	if opts.Expiry != 0 {
		tmpBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(tmpBuf[0:], opts.Expiry)
		extraBuf = append(extraBuf, tmpBuf...)
	}
	if opts.Flags != 0 {
		extraBuf = append(extraBuf, uint8(opts.Flags))
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = crud.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdSubDocMultiMutation,
			Datatype:               0,
			Cas:                    uint64(opts.Cas),
			Extras:                 extraBuf,
			Key:                    opts.Key,
			Value:                  valueBuf,
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
				makeTimeoutError(start, "MutateIn", errAmbiguousTimeout, req),
				tracer,
			)
		}))
	}

	return op, nil
}

func (crud *crudComponent) makeSubDocError(index int, code memd.StatusCode, req *memdQRequest, resp *memdQResponse) error {
	err := getKvStatusCodeError(code)
	err = translateMemdError(err, req)
	err = crud.errMapManager.EnhanceKvError(err, resp, req)

	return SubDocumentError{
		Index:      index,
		InnerError: err,
	}
}
