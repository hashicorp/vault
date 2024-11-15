package gocbcore

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/couchbase/gocbcore/v10/memd"
)

type dcpComponent struct {
	kvMux           *kvMux
	streamIDEnabled bool
}

func newDcpComponent(kvMux *kvMux, streamIDEnabled bool) *dcpComponent {
	return &dcpComponent{
		kvMux:           kvMux,
		streamIDEnabled: streamIDEnabled,
	}
}

func (dcp *dcpComponent) OpenStream(vbID uint16, flags memd.DcpStreamAddFlag, vbUUID VbUUID, startSeqNo,
	endSeqNo, snapStartSeqNo, snapEndSeqNo SeqNo, evtHandler StreamObserver, opts OpenStreamOptions,
	cb OpenStreamCallback) (PendingOp, error) {
	var req *memdQRequest
	var openHandled uint32
	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
		if resp == nil && err == nil {
			logWarnf("DCP event occurred with no error and no response")
			return
		}

		if err != nil {
			if resp == nil {
				if atomic.CompareAndSwapUint32(&openHandled, 0, 1) {
					// If open hasn't been handled and there's no response then it's reasonably safe to assume that
					// this occurring for the open stream request.
					cb(nil, err)
					return
				}
			}

			if resp != nil && resp.Magic == memd.CmdMagicRes {
				// CmdMagicRes means that this must be the open stream request response.
				atomic.StoreUint32(&openHandled, 1)
				// We need to decorate rollback errors with extra information that the server returns to us.
				// Unforunately we have to check for the memd due to earlier oversights where we missed converting
				// it to a proper gocbcore error.
				if errors.Is(err, ErrMemdRollback) {
					err = DCPRollbackError{
						InnerError: err,
						SeqNo:      SeqNo(binary.BigEndian.Uint64(resp.Value)),
					}
				}
				cb(nil, err)
				return
			}

			var streamID uint16
			if opts.StreamOptions != nil {
				streamID = opts.StreamOptions.StreamID
			}
			evtHandler.End(DcpStreamEnd{vbID, streamID}, err)
			return
		}

		if resp.Magic == memd.CmdMagicRes {
			atomic.StoreUint32(&openHandled, 1)
			// This is the response to the open stream request.
			numEntries := len(resp.Value) / 16
			entries := make([]FailoverEntry, numEntries)
			for i := 0; i < numEntries; i++ {
				entries[i] = FailoverEntry{
					VbUUID: VbUUID(binary.BigEndian.Uint64(resp.Value[i*16+0:])),
					SeqNo:  SeqNo(binary.BigEndian.Uint64(resp.Value[i*16+8:])),
				}
			}

			cb(entries, nil)
			return
		}

		// This is one of the stream events
		switch resp.Command {
		case memd.CmdDcpSnapshotMarker:
			snapShotmarker := DcpSnapshotMarker{VbID: resp.Vbucket}
			if resp.StreamIDFrame != nil {
				snapShotmarker.StreamID = resp.StreamIDFrame.StreamID
			}
			if len(resp.Extras) == 20 {
				// Length of 20 indicates a v1 packet
				snapShotmarker.StartSeqNo = binary.BigEndian.Uint64(resp.Extras[0:])
				snapShotmarker.EndSeqNo = binary.BigEndian.Uint64(resp.Extras[8:])
				snapShotmarker.SnapshotType = SnapshotState(binary.BigEndian.Uint32(resp.Extras[16:]))
			} else if len(resp.Extras) == 1 {
				// Length of 1 indicates a v2 packet
				snapShotmarker.StartSeqNo = binary.BigEndian.Uint64(resp.Value[0:])
				snapShotmarker.EndSeqNo = binary.BigEndian.Uint64(resp.Value[8:])
				snapShotmarker.SnapshotType = SnapshotState(binary.BigEndian.Uint32(resp.Value[16:]))
				snapShotmarker.MaxVisibleSeqNo = binary.BigEndian.Uint64(resp.Value[20:])
				snapShotmarker.HighCompletedSeqNo = binary.BigEndian.Uint64(resp.Value[28:])
				version := int(resp.Extras[0])
				if version == 1 {
					// v2.1 includes the snapshot TimeStamp
					snapShotmarker.SnapshotTimeStamp = binary.BigEndian.Uint64(resp.Value[36:])
				}
			}
			evtHandler.SnapshotMarker(snapShotmarker)
		case memd.CmdDcpMutation:
			mutation := DcpMutation{
				SeqNo:        binary.BigEndian.Uint64(resp.Extras[0:]),
				RevNo:        binary.BigEndian.Uint64(resp.Extras[8:]),
				Flags:        binary.BigEndian.Uint32(resp.Extras[16:]),
				Expiry:       binary.BigEndian.Uint32(resp.Extras[20:]),
				LockTime:     binary.BigEndian.Uint32(resp.Extras[24:]),
				Cas:          resp.Cas,
				Datatype:     resp.Datatype,
				VbID:         resp.Vbucket,
				CollectionID: resp.CollectionID,
				Key:          resp.Key,
				Value:        resp.Value,
			}
			if resp.StreamIDFrame != nil {
				mutation.StreamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.Mutation(mutation)
		case memd.CmdDcpDeletion:
			deletion := DcpDeletion{
				SeqNo:        binary.BigEndian.Uint64(resp.Extras[0:]),
				RevNo:        binary.BigEndian.Uint64(resp.Extras[8:]),
				Cas:          resp.Cas,
				Datatype:     resp.Datatype,
				VbID:         resp.Vbucket,
				CollectionID: resp.CollectionID,
				Key:          resp.Key,
				Value:        resp.Value,
			}
			if len(resp.Extras) == 21 {
				// Length of 21 indicates a v2 packet
				deletion.DeleteTime = binary.BigEndian.Uint32(resp.Extras[16:])
			}
			if resp.StreamIDFrame != nil {
				deletion.StreamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.Deletion(deletion)
		case memd.CmdDcpExpiration:
			expiration := DcpExpiration{
				SeqNo:        binary.BigEndian.Uint64(resp.Extras[0:]),
				RevNo:        binary.BigEndian.Uint64(resp.Extras[8:]),
				Cas:          resp.Cas,
				VbID:         resp.Vbucket,
				CollectionID: resp.CollectionID,
				Key:          resp.Key,
			}
			if len(resp.Extras) > 16 {
				expiration.DeleteTime = binary.BigEndian.Uint32(resp.Extras[16:])
			}
			if resp.StreamIDFrame != nil {
				expiration.StreamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.Expiration(expiration)
		case memd.CmdDcpEvent:
			vbID := resp.Vbucket
			seqNo := binary.BigEndian.Uint64(resp.Extras[0:])
			eventCode := memd.StreamEventCode(binary.BigEndian.Uint32(resp.Extras[8:]))
			version := resp.Extras[12]
			var streamID uint16
			if resp.StreamIDFrame != nil {
				streamID = resp.StreamIDFrame.StreamID
			}

			switch eventCode {
			case memd.StreamEventCollectionCreate:
				creation := DcpCollectionCreation{
					SeqNo:        seqNo,
					Version:      version,
					VbID:         vbID,
					ManifestUID:  binary.BigEndian.Uint64(resp.Value[0:]),
					ScopeID:      binary.BigEndian.Uint32(resp.Value[8:]),
					CollectionID: binary.BigEndian.Uint32(resp.Value[12:]),
					StreamID:     streamID,
					Key:          resp.Key,
				}
				if version == 1 {
					creation.Ttl = binary.BigEndian.Uint32(resp.Value[16:])
				}
				evtHandler.CreateCollection(creation)
			case memd.StreamEventCollectionDelete:
				deleteion := DcpCollectionDeletion{
					SeqNo:        seqNo,
					Version:      version,
					VbID:         vbID,
					ManifestUID:  binary.BigEndian.Uint64(resp.Value[0:]),
					ScopeID:      binary.BigEndian.Uint32(resp.Value[8:]),
					CollectionID: binary.BigEndian.Uint32(resp.Value[12:]),
					StreamID:     streamID,
				}
				evtHandler.DeleteCollection(deleteion)
			case memd.StreamEventCollectionFlush:
				flush := DcpCollectionFlush{
					SeqNo:        seqNo,
					Version:      version,
					VbID:         vbID,
					ManifestUID:  binary.BigEndian.Uint64(resp.Value[0:]),
					CollectionID: binary.BigEndian.Uint32(resp.Value[8:]),
					StreamID:     streamID,
				}
				evtHandler.FlushCollection(flush)
			case memd.StreamEventScopeCreate:
				creation := DcpScopeCreation{
					SeqNo:       seqNo,
					Version:     version,
					VbID:        vbID,
					ManifestUID: binary.BigEndian.Uint64(resp.Value[0:]),
					ScopeID:     binary.BigEndian.Uint32(resp.Value[8:]),
					StreamID:    streamID,
					Key:         resp.Key,
				}
				evtHandler.CreateScope(creation)
			case memd.StreamEventScopeDelete:
				deletion := DcpScopeDeletion{
					SeqNo:       seqNo,
					Version:     version,
					VbID:        vbID,
					ManifestUID: binary.BigEndian.Uint64(resp.Value[0:]),
					ScopeID:     binary.BigEndian.Uint32(resp.Value[8:]),
					StreamID:    streamID,
				}
				evtHandler.DeleteScope(deletion)
			case memd.StreamEventCollectionChanged:
				modification := DcpCollectionModification{
					SeqNo:        seqNo,
					Version:      version,
					VbID:         vbID,
					ManifestUID:  binary.BigEndian.Uint64(resp.Value[0:]),
					CollectionID: binary.BigEndian.Uint32(resp.Value[8:]),
					Ttl:          binary.BigEndian.Uint32(resp.Value[12:]),
					StreamID:     streamID,
				}
				evtHandler.ModifyCollection(modification)
			}
		case memd.CmdDcpStreamEnd:
			code := memd.StreamEndStatus(binary.BigEndian.Uint32(resp.Extras[0:]))
			end := DcpStreamEnd{
				VbID: resp.Vbucket,
			}
			if resp.StreamIDFrame != nil {
				end.StreamID = resp.StreamIDFrame.StreamID
			}
			if req.internalCancel(err) {
				evtHandler.End(end, getStreamEndStatusError(code))
			}
		case memd.CmdDcpOsoSnapshot:
			snapshot := DcpOSOSnapshot{
				VbID:         resp.Vbucket,
				SnapshotType: binary.BigEndian.Uint32(resp.Extras[0:]),
			}
			if resp.StreamIDFrame != nil {
				snapshot.StreamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.OSOSnapshot(snapshot)
		case memd.CmdDcpSeqNoAdvanced:
			seqNoAdvanced := DcpSeqNoAdvanced{
				SeqNo: binary.BigEndian.Uint64(resp.Extras[0:]),
				VbID:  resp.Vbucket,
			}
			if resp.StreamIDFrame != nil {
				seqNoAdvanced.StreamID = resp.StreamIDFrame.StreamID
			}
			evtHandler.SeqNoAdvanced(seqNoAdvanced)
		}
	}

	extraBuf := make([]byte, 48)
	binary.BigEndian.PutUint32(extraBuf[0:], uint32(flags))
	binary.BigEndian.PutUint32(extraBuf[4:], 0)
	binary.BigEndian.PutUint64(extraBuf[8:], uint64(startSeqNo))
	binary.BigEndian.PutUint64(extraBuf[16:], uint64(endSeqNo))
	binary.BigEndian.PutUint64(extraBuf[24:], uint64(vbUUID))
	binary.BigEndian.PutUint64(extraBuf[32:], uint64(snapStartSeqNo))
	binary.BigEndian.PutUint64(extraBuf[40:], uint64(snapEndSeqNo))

	var val []byte
	val = nil
	if opts.StreamOptions != nil || opts.FilterOptions != nil || opts.ManifestOptions != nil {
		convertedFilter := streamFilter{}

		if opts.FilterOptions != nil {
			// If there are collection IDs then we can assume that scope ID of 0 actually means no scope ID
			if len(opts.FilterOptions.CollectionIDs) > 0 {
				for _, cid := range opts.FilterOptions.CollectionIDs {
					convertedFilter.Collections = append(convertedFilter.Collections, fmt.Sprintf("%x", cid))
				}
			} else {
				// No collection IDs but the filter was set so even if scope ID is 0 then we use it
				convertedFilter.Scope = fmt.Sprintf("%x", opts.FilterOptions.ScopeID)
			}

		}
		if opts.ManifestOptions != nil {
			convertedFilter.ManifestUID = fmt.Sprintf("%x", opts.ManifestOptions.ManifestUID)
		}
		if opts.StreamOptions != nil {
			convertedFilter.StreamID = opts.StreamOptions.StreamID
		}

		var err error
		val, err = json.Marshal(convertedFilter)
		if err != nil {
			return nil, err
		}
	}

	req = &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdDcpStreamReq,
			Datatype: 0,
			Cas:      0,
			Extras:   extraBuf,
			Key:      nil,
			Value:    val,
			Vbucket:  vbID,
		},
		Callback:   handler,
		ReplicaIdx: 0,
		Persistent: true,
	}
	return dcp.kvMux.DispatchDirect(req)
}

func (dcp *dcpComponent) CloseStream(vbID uint16, opts CloseStreamOptions, cb CloseStreamCallback) (PendingOp, error) {
	handler := func(_ *memdQResponse, _ *memdQRequest, err error) {
		cb(err)
	}

	var streamFrame *memd.StreamIDFrame
	if opts.StreamOptions != nil {
		if !dcp.streamIDEnabled {
			return nil, errStreamIDNotEnabled
		}

		streamFrame = &memd.StreamIDFrame{
			StreamID: opts.StreamOptions.StreamID,
		}
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:         memd.CmdMagicReq,
			Command:       memd.CmdDcpCloseStream,
			Datatype:      0,
			Cas:           0,
			Extras:        nil,
			Key:           nil,
			Value:         nil,
			Vbucket:       vbID,
			StreamIDFrame: streamFrame,
		},
		Callback:      handler,
		ReplicaIdx:    0,
		Persistent:    false,
		RetryStrategy: newFailFastRetryStrategy(),
	}

	return dcp.kvMux.DispatchDirect(req)
}

func (dcp *dcpComponent) GetFailoverLog(vbID uint16, cb GetFailoverLogCallback) (PendingOp, error) {
	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
		if err != nil {
			cb(nil, err)
			return
		}

		numEntries := len(resp.Value) / 16
		entries := make([]FailoverEntry, numEntries)
		for i := 0; i < numEntries; i++ {
			entries[i] = FailoverEntry{
				VbUUID: VbUUID(binary.BigEndian.Uint64(resp.Value[i*16+0:])),
				SeqNo:  SeqNo(binary.BigEndian.Uint64(resp.Value[i*16+8:])),
			}
		}
		cb(entries, nil)
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdDcpGetFailoverLog,
			Datatype: 0,
			Cas:      0,
			Extras:   nil,
			Key:      nil,
			Value:    nil,
			Vbucket:  vbID,
		},
		Callback:      handler,
		ReplicaIdx:    0,
		Persistent:    false,
		RetryStrategy: newFailFastRetryStrategy(),
	}
	return dcp.kvMux.DispatchDirect(req)
}

func (dcp *dcpComponent) GetVbucketSeqnos(serverIdx int, state memd.VbucketState, opts GetVbucketSeqnoOptions, cb GetVBucketSeqnosCallback) (PendingOp, error) {
	handler := func(resp *memdQResponse, _ *memdQRequest, err error) {
		if err != nil {
			cb(nil, err)
			return
		}

		var vbs []VbSeqNoEntry

		numVbs := len(resp.Value) / 10
		for i := 0; i < numVbs; i++ {
			vbs = append(vbs, VbSeqNoEntry{
				VbID:  binary.BigEndian.Uint16(resp.Value[i*10:]),
				SeqNo: SeqNo(binary.BigEndian.Uint64(resp.Value[i*10+2:])),
			})
		}

		cb(vbs, nil)
	}

	var extraBuf []byte

	if opts.FilterOptions == nil {
		extraBuf = make([]byte, 4)
		binary.BigEndian.PutUint32(extraBuf[0:], uint32(state))
	} else {
		if !dcp.kvMux.SupportsCollections() {
			return nil, errCollectionsUnsupported
		}

		extraBuf = make([]byte, 8)
		binary.BigEndian.PutUint32(extraBuf[0:], uint32(state))
		binary.BigEndian.PutUint32(extraBuf[4:], opts.FilterOptions.CollectionID)
	}
	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdGetAllVBSeqnos,
			Datatype: 0,
			Cas:      0,
			Extras:   extraBuf,
			Key:      nil,
			Value:    nil,
			Vbucket:  0,
		},
		Callback:      handler,
		ReplicaIdx:    -serverIdx,
		Persistent:    false,
		RetryStrategy: newFailFastRetryStrategy(),
	}

	return dcp.kvMux.DispatchDirect(req)
}
